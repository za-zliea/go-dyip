package dns

import (
	"dyip-sync/meta"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type Google struct {
}

func NewGoogle() Google {
	return Google{}
}

const GOOGLE_URL = "https://domains.google.com/nic/update"

func (Google) Query(ipMeta *meta.IpMeta) (string, error) {
	iprecords, _ := net.LookupIP(fmt.Sprintf("%s.%s", ipMeta.Subdomain, ipMeta.Domain))

	if iprecords == nil || len(iprecords) == 0 || iprecords[0] == nil {
		return "", errors.New("empty query response")
	}

	return iprecords[0].String(), nil
}

func (Google) Sync(ipMeta *meta.IpMeta) error {
	urlParams := url.Values{}
	URI, err := url.Parse(GOOGLE_URL)
	if err != nil {
		return err
	}

	urlParams.Set("hostname", fmt.Sprintf("%s.%s", ipMeta.Subdomain, ipMeta.Domain))
	urlParams.Set("myip", *ipMeta.Ip)

	URI.RawQuery = urlParams.Encode()
	finalUrl := URI.String()

	client := &http.Client{}
	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(ipMeta.Accesskey, ipMeta.AccessKeySecret)

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	bodyStr := string(body)

	if !strings.HasPrefix(bodyStr, "good ") {
		return errors.New(bodyStr)
	}

	return nil
}
