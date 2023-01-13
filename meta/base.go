package meta

type Protocol string

const (
	IPV4 Protocol = "tcp4"
	IPV6 Protocol = "tcp6"
)

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
