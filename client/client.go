package client

import (
	"context"
	"dyip-sync/meta"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var MetaData meta.ClientMeta

func Sync() error {
	var fullUrl string
	if strings.HasSuffix(MetaData.Server, "/") {
		fullUrl = MetaData.Server + "sync"
	} else {
		fullUrl = MetaData.Server + "/sync"
	}

	params := url.Values{}
	URI, err := url.Parse(fullUrl)
	if err != nil {
		return err
	}

	params.Set("domain", MetaData.Domain)
	params.Set("auth", MetaData.Auth)
	params.Set("protocol", string(MetaData.Protocol))

	if MetaData.Local {
		iface, err := net.InterfaceByName(MetaData.Interface)
		if err != nil {
			return err
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if MetaData.Protocol == meta.IPV4 && ipnet.IP.To4() != nil {
					params.Set("localip", ipnet.IP.String())
					break
				}

				if MetaData.Protocol == meta.IPV6 && ipnet.IP.To16() != nil {
					params.Set("localip", ipnet.IP.String())
					break
				}
			}
		}
	}

	URI.RawQuery = params.Encode()
	finalUrl := URI.String()

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, _, address string) (net.Conn, error) {
				dialer := net.Dialer{}
				return dialer.DialContext(ctx, meta.GetHttpDial(MetaData.Protocol), address)
			},
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", MetaData.Token)

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	var response ResponseDTO
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if !response.IsSuccess() {
		return errors.New(response.Message)
	}

	return nil
}
