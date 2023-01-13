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

type Godaddy struct {
}

func NewGodaddy() Godaddy {
	return Godaddy{}
}

type GodaddyData struct {
	Data string `json:"data"`
	Name string `json:"name"`
	Ttl  int    `json:"ttl"`
	Type string `json:"type"`
}

type GodaddyErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (g Godaddy) Query(ipMeta *meta.IpMeta) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", g.godaddyUrl(ipMeta), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", g.godaddyAuthorization(ipMeta))

	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	if rsp.StatusCode == 200 {
		var response []GodaddyData
		err = json.Unmarshal(body, &response)
		if err != nil {
			return "", err
		}

		if response == nil || len(response) == 0 || response[0].Data == "" {
			return "", nil
		}

		return response[0].Data, nil
	} else {
		var response GodaddyErrorResponse
		err = json.Unmarshal(body, &response)

		if err != nil {
			return "", err
		}

		return "", errors.New(response.Message)
	}
}

func (g Godaddy) Sync(ipMeta *meta.IpMeta) error {
	godaddyData := GodaddyData{Data: *ipMeta.Ip, Name: ipMeta.Subdomain, Ttl: 600, Type: meta.GetProtocolDns(ipMeta.Protocol)}

	data := []GodaddyData{godaddyData}

	reqBody, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", g.godaddyUrl(ipMeta), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", g.godaddyAuthorization(ipMeta))

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode == 200 {
		return nil
	} else {

		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return err
		}

		var response GodaddyErrorResponse
		err = json.Unmarshal(body, &response)

		if err != nil {
			return err
		}

		return errors.New(response.Message)
	}
}

func (Godaddy) godaddyUrl(ipMeta *meta.IpMeta) string {
	return fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/%s/%s", ipMeta.Domain, meta.GetProtocolDns(ipMeta.Protocol), ipMeta.Subdomain)
}

func (Godaddy) godaddyAuthorization(ipMeta *meta.IpMeta) string {
	return fmt.Sprintf("sso-key %s:%s", ipMeta.Accesskey, ipMeta.AccessKeySecret)
}
