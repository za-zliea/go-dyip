package server

import (
	"dyip-sync/src/config"
	"dyip-sync/src/dns"
	dymeta "dyip-sync/src/meta"
	"dyip-sync/src/util"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/savsgio/atreugo/v11"
)

var ConfigFileServer string
var MetaData dymeta.ServerMeta

type IpResponse struct {
	Ip string `json:"ip"`
}

type IpmResponse struct {
	Domain    string                `json:"domain"`
	Subdomain string                `json:"subdomain"`
	Ip        *string               `json:"ip,omitempty"`
	Dip       *string               `json:"dip,omitempty"`
	History   []dymeta.HistoryEntry `json:"history"`
}

func IpHandler(ctx *atreugo.RequestCtx) error {
	auth := authGlobal(ctx)
	if !auth.IsSuccess() {
		return ctx.JSONResponse(auth, auth.status)
	}

	return ctx.JSONResponse(SuccessWithD(IpResponse{Ip: clientIP(ctx)}))
}

func SyncHandler(ctx *atreugo.RequestCtx) error {
	auth := authGlobal(ctx)
	if !auth.IsSuccess() {
		return ctx.JSONResponse(auth, auth.status)
	}

	ipMeta, err := authDomain(ctx)
	if err != nil {
		return ctx.JSONResponse(Failed(err.Error()), 401)
	}

	var ip string
	if dymeta.IsLocalIp(ipMeta.SyncType) {
		ip = string(ctx.QueryArgs().Peek("localip"))
	} else {
		ip = clientIP(ctx)
		if ip == "" {
			log.Printf("sync %s.%s error: no ip\n", ipMeta.Subdomain, ipMeta.Domain)
			return ctx.JSONResponse(Failed("no ip"), 200)
		}
	}

	protocol, err := util.GetIpFamily(ip)
	if err != nil {
		log.Printf("sync %s.%s-%s error: invalid ip\n", ipMeta.Subdomain, ipMeta.Domain, ip)
		return ctx.JSONResponse(Failed("protocol empty"), 200)
	}

	if ipMeta.Protocol != protocol {
		log.Printf("sync %s.%s-%s error: ip protocol not match\n", ipMeta.Subdomain, ipMeta.Domain, ip)
		return ctx.JSONResponse(Failed("protocol not match"), 200)
	}

	msg, err := performSync(ipMeta, ip)
	if err != nil {
		message := fmt.Sprintf("sync %s.%s-%s %s", ipMeta.Subdomain, ipMeta.Domain, ip, msg)
		log.Println(message)
		return ctx.JSONResponse(Failed(message), 200)
	}

	log.Printf("sync %s.%s-%s %s\n", ipMeta.Subdomain, ipMeta.Domain, ip, msg)
	return ctx.JSONResponse(Success(), 200)
}

func LoadHandler(ctx *atreugo.RequestCtx) error {
	auth := authGlobal(ctx)
	if !auth.IsSuccess() {
		return ctx.JSONResponse(auth, auth.status)
	}

	ipMeta, err := authDomain(ctx)
	if err != nil {
		return ctx.JSONResponse(Failed(err.Error()), 401)
	}

	dip, err := dns.NewDns().Query(ipMeta)
	if err != nil {
		message := fmt.Sprintf("load %s.%s provider error: %v", ipMeta.Subdomain, ipMeta.Domain, err)
		log.Println(message)
		return ctx.JSONResponse(Failed(message))
	}

	var _dip *string
	if dip != "" {
		_dip = &dip
	} else {
		_dip = nil
	}

	return ctx.JSONResponse(SuccessWithD(IpmResponse{Domain: ipMeta.Domain, Subdomain: ipMeta.Subdomain, Ip: ipMeta.Ip, Dip: _dip, History: ipMeta.History}), 200)
}

func authGlobal(ctx *atreugo.RequestCtx) ResponseDTO {
	var response ResponseDTO

	token := string(ctx.Request.Header.Peek("Authorization"))
	if token == "" {
		response = FailedWithS("auth failed", 401)
	} else if token != MetaData.Token {
		response = FailedWithS("auth failed", 401)
	} else {
		response = Success()
	}

	return response
}

func authDomain(ctx *atreugo.RequestCtx) (*dymeta.IpMeta, error) {
	domain := string(ctx.QueryArgs().Peek("domain"))
	domainAuth := string(ctx.QueryArgs().Peek("auth"))
	protocolBytes := ctx.QueryArgs().Peek("protocol")

	protocol := dymeta.IPV4
	if protocolBytes != nil {
		protocol = dymeta.Protocol(protocolBytes)
	}

	ipMeta, ok := MetaData.MetaMap[domain+"."+string(protocol)]
	if !ok {
		return nil, errors.New(fmt.Sprintf("domain %s auth failed", domain))
	}
	if ipMeta.Auth != domainAuth {
		return nil, errors.New(fmt.Sprintf("domain %s.%s auth failed", ipMeta.Subdomain, ipMeta.Domain))
	}

	return ipMeta, nil
}

// clientIP returns the caller's IP, preferring the header named by MetaData.RealIp
// (e.g. x-real-ip for reverse-proxy deployments) and falling back to the TCP remote addr.
func clientIP(ctx *atreugo.RequestCtx) string {
	if MetaData.RealIp == nil || *MetaData.RealIp == "" {
		return ctx.RemoteIP().String()
	}
	return string(ctx.Request.Header.Peek(*MetaData.RealIp))
}

// parseProtocol upper-cases and validates a protocol string ("IPV4"/"IPV6").
func parseProtocol(s string) (dymeta.Protocol, bool) {
	switch strings.ToUpper(s) {
	case string(dymeta.IPV4):
		return dymeta.IPV4, true
	case string(dymeta.IPV6):
		return dymeta.IPV6, true
	}
	return "", false
}

// recordHistory prepends {ip, now} to ipMeta.History, capping at 5 entries (newest first).
func recordHistory(ipMeta *dymeta.IpMeta, ip string, now time.Time) {
	prev := ipMeta.History
	length := len(prev) + 1
	if length > 5 {
		length = 5
	}

	history := make([]dymeta.HistoryEntry, length)
	history[0] = dymeta.HistoryEntry{Ip: ip, Time: now}
	for i := 1; i < length; i++ {
		history[i] = prev[i-1]
	}

	ipMeta.History = history
}

// performSync runs the post-IP-resolution sync pipeline shared by SyncHandler and
// FrontSyncHandler. Callers MUST guarantee ip != "" and that ipMeta.Protocol matches the
// family of ip. Returns (message, err); a nil err means success, with message describing
// the outcome ("same ip, skip" or "success").
func performSync(ipMeta *dymeta.IpMeta, ip string) (string, error) {
	if ipMeta.Ip != nil && *ipMeta.Ip == ip {
		return "same ip, skip", nil
	}

	ipMeta.Ip = &ip
	recordHistory(ipMeta, ip, time.Now())

	newdns := dns.NewDns()
	dip, err := newdns.Query(ipMeta)
	if err != nil {
		return fmt.Sprintf("query provider error: %v", err), err
	}

	if dip != ip {
		if err := newdns.Sync(ipMeta); err != nil {
			return fmt.Sprintf("sync provider error: %v", err), err
		}
	}

	if err := config.WriteConfig(ConfigFileServer, &MetaData); err != nil {
		return fmt.Sprintf("write config error: %v", err), err
	}

	return "success", nil
}
