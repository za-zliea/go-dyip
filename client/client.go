package client

import (
	"dyip-sync/meta"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

	URI.RawQuery = params.Encode()
	finalUrl := URI.String()

	client := &http.Client{}
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
