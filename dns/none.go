package dns

import (
	"dyip-sync/meta"
	"fmt"
	"net"
	"strings"
)

type None struct {
}

func NewNone() None {
	return None{}
}

func (g None) Query(ipMeta *meta.IpMeta) (string, error) {
	if strings.HasSuffix(ipMeta.Domain, ".internal") {
		return "", nil
	}

	iprecords, _ := net.LookupIP(fmt.Sprintf("%s.%s", ipMeta.Subdomain, ipMeta.Domain))

	if iprecords == nil || len(iprecords) == 0 || iprecords[0] == nil {
		return "", nil
	} else {
		return iprecords[0].String(), nil
	}
}

func (g None) Sync(ipMeta *meta.IpMeta) error {
	return nil
}
