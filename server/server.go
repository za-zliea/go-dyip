package server

import (
	"dyip-sync/config"
	"dyip-sync/dns"
	"dyip-sync/meta"
	"dyip-sync/util"
	"errors"
	"fmt"
	"github.com/savsgio/atreugo/v11"
	"log"
)

var ConfigFileServer string
var MetaData meta.ServerMeta

type IpResponse struct {
	Ip string `json:"ip"`
}

type IpmResponse struct {
	Domain    string   `json:"domain"`
	Subdomain string   `json:"subdomain"`
	Ip        *string  `json:"ip,omitempty"`
	Dip       *string  `json:"dip,omitempty"`
	History   []string `json:"history"`
}

func IndexHandler(ctx *atreugo.RequestCtx) error {
	return ctx.JSONResponse(Success())
}

func IpHandler(ctx *atreugo.RequestCtx) error {
	auth := authGlobal(ctx)
	if !auth.IsSuccess() {
		return ctx.JSONResponse(auth, auth.status)
	}

	var ip string
	if MetaData.RealIp == nil || *MetaData.RealIp == "" {
		ip = ctx.RemoteIP().String()
	} else {
		ip = string(ctx.Request.Header.Peek(*MetaData.RealIp))
	}

	return ctx.JSONResponse(SuccessWithD(IpResponse{Ip: ip}))
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
	if ipMeta.Local {
		ip = string(ctx.QueryArgs().Peek("localip"))
	} else {
		if MetaData.RealIp == nil || *MetaData.RealIp == "" {
			ip = ctx.RemoteIP().String()
		} else {
			ip = string(ctx.Request.Header.Peek(*MetaData.RealIp))
		}
		if ip == "" {
			log.Printf("sync %s.%s-%s error: no ip\n", ipMeta.Subdomain, ipMeta.Domain, ip)
			return ctx.JSONResponse(Failed("no ip"), 200)
		}
	}

	protocol, err := util.GetIpFamily(ip)
	if err != nil {
		log.Printf("sync %s.%s-%s error: check ip protocol empty\n", ipMeta.Subdomain, ipMeta.Domain, ip)
		return ctx.JSONResponse(Failed("protocol empty"), 200)
	}

	if ipMeta.Protocol != protocol {
		log.Printf("sync %s.%s-%s error: check ip protocol not match\n", ipMeta.Subdomain, ipMeta.Domain, ip)
		return ctx.JSONResponse(Failed("protocol not match"), 200)
	}

	if ipMeta.Ip != nil && *ipMeta.Ip == ip {
		log.Printf("sync %s.%s-%s same ip, skip\n", ipMeta.Subdomain, ipMeta.Domain, ip)
		return ctx.JSONResponse(Success(), 200)
	}

	ipMeta.Ip = &ip
	var length int
	if ipMeta.History == nil {
		length = 0
	} else {
		length = len(ipMeta.History)
	}

	length = length + 1
	if length > 5 {
		length = 5
	}

	history := make([]string, length)
	history[0] = ip
	for i := 1; i < length; i++ {
		history[i] = ipMeta.History[i-1]
	}

	ipMeta.History = history

	newdns := dns.NewDns()
	dip, err := newdns.Query(ipMeta)
	if err != nil {
		message := fmt.Sprintf("sync %s.%s-%s query provider error: %v\n", ipMeta.Subdomain, ipMeta.Domain, ip, err)
		log.Printf(message)
		return ctx.JSONResponse(Failed(message), 200)
	}

	if dip != ip {
		err = newdns.Sync(ipMeta)
		if err != nil {
			message := fmt.Sprintf("sync %s.%s-%s sync provider error: %v", ipMeta.Subdomain, ipMeta.Domain, ip, err)
			log.Println(message)
			return ctx.JSONResponse(Failed(message), 200)
		}
	} else {
		log.Printf("sync %s.%s-%s provider same ip, skip\n", ipMeta.Subdomain, ipMeta.Domain, ip)
	}

	err = config.WriteConfig(ConfigFileServer, &MetaData)
	if err != nil {
		message := fmt.Sprintf("sync %s.%s-%s error: write to config file error %v\n", ipMeta.Subdomain, ipMeta.Domain, ip, err)
		log.Printf(message)
		return ctx.JSONResponse(Failed(message), 200)
	}

	log.Printf("sync %s.%s-%s success\n", ipMeta.Subdomain, ipMeta.Domain, ip)
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

func authDomain(ctx *atreugo.RequestCtx) (*meta.IpMeta, error) {
	domain := string(ctx.QueryArgs().Peek("domain"))
	domainAuth := string(ctx.QueryArgs().Peek("auth"))
	protocolBytes := ctx.QueryArgs().Peek("protocol")

	protocol := meta.IPV4
	if protocolBytes != nil {
		protocol = meta.Protocol(protocolBytes)
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
