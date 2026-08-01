package server

import (
	"dyip-sync/src/dns"
	dymeta "dyip-sync/src/meta"
	"dyip-sync/src/util"
	"encoding/json"
	"log"
	"time"

	"github.com/savsgio/atreugo/v11"
)

// pathParam reads an atreugo path parameter as a string, tolerating either a string or
// []byte backing value.
func pathParam(ctx *atreugo.RequestCtx, key string) string {
	switch v := ctx.UserValue(key).(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	return ""
}

// FrontIpResponse reports the caller's IP split by family. A single HTTP request arrives
// on exactly one family, so only the matching field is populated.
type FrontIpResponse struct {
	Ipv4 string `json:"ipv4,omitempty"`
	Ipv6 string `json:"ipv6,omitempty"`
}

// FrontIpHandler returns the web user's IP (resolved via MetaData.RealIp), classified into
// IPV4 and IPV6.
// GET /front/api/ip/self
func FrontIpHandler(ctx *atreugo.RequestCtx) error {
	if _, ok := authJWT(ctx); !ok {
		return ctx.JSONResponse(FailedWithS("auth failed", 401), 401)
	}

	resp := FrontIpResponse{}
	if ip := clientIP(ctx); ip != "" {
		if family, err := util.GetIpFamily(ip); err == nil {
			switch family {
			case dymeta.IPV4:
				resp.Ipv4 = ip
			case dymeta.IPV6:
				resp.Ipv6 = ip
			}
		}
	}
	return ctx.JSONResponse(SuccessWithD(resp), 200)
}

// frontTimeLayout is the format used for all times returned by the web API.
const frontTimeLayout = "2006-01-02 15:04:05"

// formatFrontTime formats t as yyyy-MM-dd HH:mm:ss; the zero time becomes "" so it
// can be omitted via `json:",omitempty"`.
func formatFrontTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(frontTimeLayout)
}

// FrontHistoryEntry is one history record formatted for the web UI.
type FrontHistoryEntry struct {
	Ip   string `json:"ip"`
	Time string `json:"time"` // yyyy-MM-dd HH:mm:ss, empty if zero
}

// FrontInfoEntry is the full record for one subdomain.domain.PROTOCOL entry.
type FrontInfoEntry struct {
	Domain     string              `json:"domain"`
	Subdomain  string              `json:"subdomain"`
	Provider   string              `json:"provider"`
	Protocol   dymeta.Protocol     `json:"protocol"`
	Ip         string              `json:"ip,omitempty"`          // 记录的最新IP
	UpdateTime string              `json:"update_time,omitempty"` // yyyy-MM-dd HH:mm:ss，未同步过则省略
	Dip        string              `json:"dip,omitempty"`         // DNS 对应IP（实时查询）
	Consistent bool                `json:"consistent"`            // 记录IP == DNS IP
	History    []FrontHistoryEntry `json:"history"`               // 历史更新记录
	DnsError   string              `json:"dns_error,omitempty"`   // 非致命，按条记录
}

// FrontInfoHandler returns the full info for one subdomain.domain.PROTOCOL record.
// GET /front/api/:domain/:subdomain/:protocol/info
func FrontInfoHandler(ctx *atreugo.RequestCtx) error {
	if _, ok := authJWT(ctx); !ok {
		return ctx.JSONResponse(FailedWithS("auth failed", 401), 401)
	}

	domain := pathParam(ctx, "domain")
	subdomain := pathParam(ctx, "subdomain")
	protocol, ok := parseProtocol(pathParam(ctx, "protocol"))
	if !ok {
		return ctx.JSONResponse(FailedWithS("invalid protocol", 400), 400)
	}

	ipMeta, ok := MetaData.MetaMap[subdomain+"."+domain+"."+string(protocol)]
	if !ok {
		return ctx.JSONResponse(FailedWithS("domain not found", 404), 404)
	}

	history := make([]FrontHistoryEntry, len(ipMeta.History))
	for i, h := range ipMeta.History {
		history[i] = FrontHistoryEntry{Ip: h.Ip, Time: formatFrontTime(h.Time)}
	}

	entry := FrontInfoEntry{
		Domain:    ipMeta.Domain,
		Subdomain: ipMeta.Subdomain,
		Provider:  ipMeta.Provider,
		Protocol:  ipMeta.Protocol,
		History:   history,
	}
	if ipMeta.Ip != nil {
		entry.Ip = *ipMeta.Ip
	}
	if len(ipMeta.History) > 0 {
		entry.UpdateTime = formatFrontTime(ipMeta.History[0].Time)
	}

	if dip, err := dns.NewDns().Query(ipMeta); err != nil {
		log.Printf("front info %s.%s/%s query error: %v", subdomain, domain, protocol, err)
		entry.DnsError = err.Error()
	} else {
		entry.Dip = dip
		entry.Consistent = dip != "" && entry.Ip != "" && dip == entry.Ip
	}

	return ctx.JSONResponse(SuccessWithD(entry), 200)
}

// FrontSyncRequest is the body of FrontSyncHandler.
type FrontSyncRequest struct {
	Ip string `json:"ip"`
}

// FrontSyncHandler pushes a submitted IP to DNS for one subdomain.domain.PROTOCOL record.
// POST /front/api/:domain/:subdomain/:protocol/sync
func FrontSyncHandler(ctx *atreugo.RequestCtx) error {
	if _, ok := authJWT(ctx); !ok {
		return ctx.JSONResponse(FailedWithS("auth failed", 401), 401)
	}

	domain := pathParam(ctx, "domain")
	subdomain := pathParam(ctx, "subdomain")
	protocol, ok := parseProtocol(pathParam(ctx, "protocol"))
	if !ok {
		return ctx.JSONResponse(FailedWithS("invalid protocol", 400), 400)
	}

	var req FrontSyncRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return ctx.JSONResponse(FailedWithS("invalid body", 400), 400)
	}
	if req.Ip == "" {
		return ctx.JSONResponse(FailedWithS("ip required", 400), 400)
	}

	family, err := util.GetIpFamily(req.Ip)
	if err != nil {
		return ctx.JSONResponse(FailedWithS("invalid ip", 400), 400)
	}

	key := subdomain + "." + domain + "." + string(protocol)
	ipMeta, ok := MetaData.MetaMap[key]
	if !ok {
		return ctx.JSONResponse(FailedWithS("domain not found", 404), 404)
	}

	if !dymeta.ConsoleEnabled(ipMeta.SyncType) {
		return ctx.JSONResponse(FailedWithS("console sync not supported for this domain", 403), 403)
	}

	if family != protocol {
		return ctx.JSONResponse(FailedWithS("ip family does not match protocol", 400), 400)
	}

	msg, err := performSync(ipMeta, req.Ip)
	if err != nil {
		log.Printf("front sync %s error: %v", key, err)
		return ctx.JSONResponse(FailedWithS(msg, 500), 500)
	}

	log.Printf("front sync %s %s", key, msg)
	if msg == "same ip, skip" {
		return ctx.JSONResponse(SuccessWithM(msg), 200)
	}
	return ctx.JSONResponse(Success(), 200)
}

// FrontDomainEntry is the summary of one configured domain record.
type FrontDomainEntry struct {
	Domain     string          `json:"domain"`                // 根域名 (apex)
	Subdomain  string          `json:"subdomain"`             // 主机标签
	Name       string          `json:"name"`                  // 完整域名 = subdomain+"."+domain
	Provider   string          `json:"provider"`              // 域名提供商
	Protocol   dymeta.Protocol `json:"protocol"`              // IPV4/IPV6
	Ip         string          `json:"ip,omitempty"`          // 记录的最新IP
	UpdateTime string          `json:"update_time,omitempty"` // yyyy-MM-dd HH:mm:ss
}

// FrontDomainHandler returns a summary of every configured domain record (no DNS lookup).
// GET /front/api/domain
func FrontDomainHandler(ctx *atreugo.RequestCtx) error {
	if _, ok := authJWT(ctx); !ok {
		return ctx.JSONResponse(FailedWithS("auth failed", 401), 401)
	}

	entries := make([]FrontDomainEntry, 0, len(MetaData.Metas))
	for _, m := range MetaData.Metas {
		entry := FrontDomainEntry{
			Domain:    m.Domain,
			Subdomain: m.Subdomain,
			Name:      m.Subdomain + "." + m.Domain,
			Provider:  m.Provider,
			Protocol:  m.Protocol,
		}
		if m.Ip != nil {
			entry.Ip = *m.Ip
		}
		if len(m.History) > 0 {
			entry.UpdateTime = formatFrontTime(m.History[0].Time)
		}
		entries = append(entries, entry)
	}

	return ctx.JSONResponse(SuccessWithD(entries), 200)
}
