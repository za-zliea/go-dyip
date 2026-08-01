package meta

type Protocol string

const (
	IPV4 Protocol = "IPV4"
	IPV6 Protocol = "IPV6"
)

func GetHttpDial(protocol Protocol) string {
	var dnsProtocol string
	switch protocol {
	case IPV4:
		dnsProtocol = "tcp4"
		break
	case IPV6:
		dnsProtocol = "tcp6"
		break
	default:
		dnsProtocol = "tcp4"
		break
	}
	return dnsProtocol
}

func GetProtocolDns(protocol Protocol) string {
	var dnsProtocol string
	switch protocol {
	case IPV4:
		dnsProtocol = "A"
		break
	case IPV6:
		dnsProtocol = "AAAA"
		break
	default:
		dnsProtocol = "A"
		break
	}
	return dnsProtocol
}

// IsLocalIp reports whether the record's IP is provided locally (SyncType last digit == '1'):
// the server reads the submitted localip instead of deriving the IP from the connection,
// and the client reads the local NIC and sends localip. An empty/malformed SyncType decodes
// to false (equivalent to "00").
func IsLocalIp(syncType string) bool {
	return len(syncType) > 0 && syncType[len(syncType)-1] == '1'
}

// ConsoleEnabled reports whether the web console may push an IP for this record
// (SyncType first digit == '1'). Server-side only; the client ignores the first digit.
// An empty/malformed SyncType decodes to false (equivalent to "00").
func ConsoleEnabled(syncType string) bool {
	return len(syncType) > 0 && syncType[0] == '1'
}
