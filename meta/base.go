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
