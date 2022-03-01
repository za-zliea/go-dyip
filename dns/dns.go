package dns

import (
	"dyip-sync/meta"
	"errors"
	"strings"
)

type DnsOperate interface {
	Query(ipMeta *meta.IpMeta) (string, error)
	Sync(ipMeta *meta.IpMeta) error
}

type Dns struct {
}

func NewDns() Dns {
	return Dns{}
}

func (d Dns) Query(ipMeta *meta.IpMeta) (string, error) {
	var ip string
	var err error
	switch strings.ToUpper(ipMeta.Provider) {
	case "NONE":
		ip, err = NewNone().Query(ipMeta)
	case "GODADDY":
		ip, err = NewGodaddy().Query(ipMeta)
	case "TENCENT":
		ip, err = NewTencent().Query(ipMeta)
	case "ALIYUN":
		ip, err = NewAliyun().Query(ipMeta)
	case "GOOGLE":
		ip, err = NewGoogle().Query(ipMeta)
	case "CLOUDFLARE":
		ip, err = NewCloudflare().Query(ipMeta)
	default:
		ip = ""
		err = errors.New("provider not support")
	}
	return ip, err
}

func (d Dns) Sync(ipMeta *meta.IpMeta) error {
	var err error
	switch strings.ToUpper(ipMeta.Provider) {
	case "NONE":
		err = NewNone().Sync(ipMeta)
	case "GODADDY":
		err = NewGodaddy().Sync(ipMeta)
	case "TENCENT":
		err = NewTencent().Sync(ipMeta)
	case "ALIYUN":
		err = NewAliyun().Sync(ipMeta)
	case "GOOGLE":
		err = NewGoogle().Sync(ipMeta)
	case "CLOUDFLARE":
		err = NewCloudflare().Sync(ipMeta)
	default:
		err = errors.New("provider not support")
	}
	return err
}
