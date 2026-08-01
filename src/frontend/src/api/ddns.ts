import { request } from '@/utils/request'

// Matches src/meta/base.go Protocol.
export type DdnsProtocol = 'IPV4' | 'IPV6'

// GET /front/api/domain  (array element) — matches src/server/front.go
// FrontDomainEntry. `ip` / `update_time` are omitted when the record has
// never been synced. NOTE: this endpoint does NO DNS lookup, so there is no
// `dip` / `consistent` here — those only come from ddnsInfo().
export interface DdnsDomainEntry {
  domain: string
  subdomain: string
  // FQDN = subdomain + '.' + domain, built server-side.
  name: string
  provider: string
  protocol: DdnsProtocol
  // Latest synced IP (omitted when never synced).
  ip?: string
  // "yyyy-MM-dd HH:mm:ss", omitted when no history.
  update_time?: string
  // Raw 2-digit code, e.g. "11".
  synctype?: string
  // synctype[0] === '1' — web console may push an IP.
  console_enabled: boolean
  // synctype[last] === '1' — client uploads localip.
  client_upload: boolean
}

// History entry inside ddnsInfo() — newest first, capped at 5.
export interface DdnsHistoryEntry {
  ip: string
  // "yyyy-MM-dd HH:mm:ss".
  time: string
}

// GET /front/api/{domain}/{subdomain}/{protocol}/info  (single object) —
// matches src/server/front.go FrontInfoEntry. `dip` comes from a live DNS
// query and is omitted on DNS error; `consistent = dip !== '' && ip !== ''
// && dip === ip` (so an empty dip ⇒ consistent:false, NOT a real mismatch).
export interface DdnsInfo {
  domain: string
  subdomain: string
  provider: string
  protocol: DdnsProtocol
  ip?: string
  update_time?: string
  // Live DNS query result; '' on DNS error.
  dip?: string
  consistent: boolean
  history: DdnsHistoryEntry[]
  // Declared in the struct but never populated server-side today.
  dns_error?: string
  synctype?: string
  console_enabled: boolean
  client_upload: boolean
}

// GET /front/api/ip/self — exactly one of ipv4/ipv6 is set per request family.
export interface DdnsSelfIp {
  ipv4?: string
  ipv6?: string
}

// POST /front/api/{domain}/{subdomain}/{protocol}/sync body.
export interface DdnsSyncRequest {
  ip: string
}

// GET /front/api/domain → DdnsDomainEntry[]
export function listDdnsDomains(): Promise<DdnsDomainEntry[]> {
  return request<DdnsDomainEntry[]>({
    url: '/front/api/domain',
    method: 'GET'
  })
}

// GET /front/api/{domain}/{subdomain}/{protocol}/info → DdnsInfo
export function ddnsInfo(
  domain: string,
  subdomain: string,
  protocol: DdnsProtocol
): Promise<DdnsInfo> {
  const d = encodeURIComponent(domain)
  const s = encodeURIComponent(subdomain)
  const p = encodeURIComponent(protocol)
  return request<DdnsInfo>({
    url: `/front/api/${d}/${s}/${p}/info`,
    method: 'GET'
  })
}

// GET /front/api/ip/self → DdnsSelfIp
export function ddnsSelfIp(): Promise<DdnsSelfIp> {
  return request<DdnsSelfIp>({
    url: '/front/api/ip/self',
    method: 'GET'
  })
}

// POST /front/api/{domain}/{subdomain}/{protocol}/sync → null
//
// On success the server returns {code:0, message:"SUCCESS"|"same ip, skip",
// data:null}. The response interceptor unwraps to `data` (null) and discards
// `message`, so the "same ip, skip" case is indistinguishable here and is
// handled client-side (see views/Ddns/Update). Errors (403/400/500) are
// surfaced by the global interceptor with the server's message.
export function syncDdns(
  domain: string,
  subdomain: string,
  protocol: DdnsProtocol,
  body: DdnsSyncRequest
): Promise<void> {
  const d = encodeURIComponent(domain)
  const s = encodeURIComponent(subdomain)
  const p = encodeURIComponent(protocol)
  return request<void>({
    url: `/front/api/${d}/${s}/${p}/sync`,
    method: 'POST',
    data: body
  })
}
