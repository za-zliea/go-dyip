package util

import (
	"dyip-sync/meta"
	"errors"
	"fmt"
	"net"
)

func GetIpFamily(address string) (meta.Protocol, error) {
	ip := net.ParseIP(address)
	if ip == nil {
		fmt.Printf("Invalid IP Address: %s\n", address)
		return "", errors.New(fmt.Sprintf("Invalid IP Address: %s\n", address))
	}
	if ip.To4() != nil {
		return meta.IPV4, nil
	} else {
		return meta.IPV6, nil
	}
}
