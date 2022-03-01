package dns

import (
	"bytes"
	"dyip-sync/meta"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Cloudflare struct {
}

func NewCloudflare() Cloudflare {
	return Cloudflare{}
}

var cloudflareIdMap map[string]CloudflareId

type CloudflareResponse struct {
	Success bool              `json:"success"`
	Errors  []CloudflareError `json:"errors"`
}

type CloudflareError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CloudflareZone struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type CloudflareZoneResponse struct {
	CloudflareResponse
	Result []*CloudflareZone `json:"result"`
}

type CloudflareData struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

type CloudflareDataQueryResponse struct {
	CloudflareResponse
	Result []*CloudflareData `json:"result"`
}

type CloudflareDataResponse struct {
	CloudflareResponse
	Result CloudflareData `json:"result"`
}

type CloudflareId struct {
	ZoneId   string
	RecordId string
}

func (c Cloudflare) Query(ipMeta *meta.IpMeta) (string, error) {
	client := &http.Client{}

	url, err := c.cloudflareUrl(ipMeta)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", c.cloudflareAuthorization(ipMeta))

	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	var response CloudflareDataResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if response.Success {
		return response.Result.Content, nil
	} else {
		return "", errors.New(response.Errors[0].Message)
	}
}

func (c Cloudflare) Sync(ipMeta *meta.IpMeta) error {
	cloudflareData := CloudflareData{Content: *ipMeta.Ip, Name: fmt.Sprintf("%s.%s", ipMeta.Subdomain, ipMeta.Domain), Ttl: 300, Type: "A", Proxied: false}

	reqBody, err := json.Marshal(&cloudflareData)
	if err != nil {
		return err
	}

	client := &http.Client{}

	url, err := c.cloudflareUrl(ipMeta)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.cloudflareAuthorization(ipMeta))

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	var response CloudflareResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Success {
		return nil
	} else {
		return errors.New(response.Errors[0].Message)
	}
}

func (c Cloudflare) cloudflareUrl(ipMeta *meta.IpMeta) (string, error) {
	cloudflareId, err := c.cloudflareId(ipMeta)

	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", cloudflareId.ZoneId, cloudflareId.RecordId), nil
}

func (Cloudflare) cloudflareAuthorization(ipMeta *meta.IpMeta) string {
	return fmt.Sprintf("Bearer %s", ipMeta.AccessKeySecret)
}

func (c Cloudflare) cloudflareId(ipMeta *meta.IpMeta) (*CloudflareId, error) {
	fullDomain := fmt.Sprintf("%s.%s", ipMeta.Subdomain, ipMeta.Domain)
	cloudflareId, ok := cloudflareIdMap[fullDomain]

	if !ok {
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/zones", nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", c.cloudflareAuthorization(ipMeta))

		query := req.URL.Query()
		query.Add("name", ipMeta.Domain)
		query.Add("status", "active")
		req.URL.RawQuery = query.Encode()

		rsp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()

		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		var zoneResponse CloudflareZoneResponse
		err = json.Unmarshal(body, &zoneResponse)
		if err != nil {
			return nil, err
		}

		if zoneResponse.Result[0] == nil {
			return nil, errors.New("empty zone")
		}

		zoneId := zoneResponse.Result[0].Id

		req, err = http.NewRequest("GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneId), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", c.cloudflareAuthorization(ipMeta))

		query = req.URL.Query()
		query.Add("name", fullDomain)
		query.Add("type", "A")
		req.URL.RawQuery = query.Encode()

		rsp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()

		body, err = ioutil.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		var dataQueryResponse CloudflareDataQueryResponse
		err = json.Unmarshal(body, &dataQueryResponse)
		if err != nil {
			return nil, err
		}

		if dataQueryResponse.Result[0] == nil {
			return nil, errors.New("empty dns")
		}

		recordId := dataQueryResponse.Result[0].Id
		cloudflareId = CloudflareId{
			ZoneId:   zoneId,
			RecordId: recordId,
		}
	}
	return &cloudflareId, nil
}
