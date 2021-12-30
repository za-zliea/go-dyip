package dns

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"dyip-sync/meta"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"
)

type Aliyun struct {
}

func NewAliyun() Aliyun {
	return Aliyun{}
}

type AliyunQueryResponse struct {
	HostId        string                     `json:"HostId"`
	Code          *string                    `json:"Code"`
	Message       *string                    `json:"Message"`
	RequestId     string                     `json:"RequestId"`
	TotalCount    int                        `json:"TotalCount"`
	PageSize      int                        `json:"PageSize"`
	DomainRecords *AliyunQueryRecordResponse `json:"DomainRecords"`
	PageNumber    int                        `json:"PageNumber"`
}

type AliyunQueryRecordResponse struct {
	Record []*AliyunQueryRecordDataResponse `json:"Record"`
}

type AliyunQueryRecordDataResponse struct {
	RR         string `json:"RR"`
	Line       string `json:"Line"`
	Status     string `json:"Status"`
	Locked     bool   `json:"Locked"`
	Type       string `json:"Type"`
	DomainName string `json:"DomainName"`
	Value      string `json:"Value"`
	RecordId   string `json:"RecordId"`
	TTL        int    `json:"TTL"`
	Weight     int    `json:"Weight"`
}

type AliyunRecordResponse struct {
	HostId    string  `json:"HostId"`
	Code      *string `json:"Code"`
	Message   *string `json:"Message"`
	RequestId string  `json:"RequestId"`
	RecordId  string  `json:"RecordId"`
}

const ALIYUN_URL = "https://alidns.aliyuncs.com"

func (a Aliyun) Query(ipMeta *meta.IpMeta) (string, error) {
	recordData, err := a.query(ipMeta)

	if err != nil {
		return "", err
	}

	return recordData.Value, nil
}

func (a Aliyun) query(ipMeta *meta.IpMeta) (*AliyunQueryRecordDataResponse, error) {
	params := make(map[string]string)

	params["SubDomain"] = fmt.Sprintf("%s.%s", ipMeta.Subdomain, ipMeta.Domain)
	params["Type"] = "A"
	params["DomainName"] = ipMeta.Domain
	params["Line"] = "default"

	urlParams := a.aliyunSign(ipMeta, "DescribeSubDomainRecords", "GET", params)

	URI, err := url.Parse(ALIYUN_URL)
	if err != nil {
		return nil, err
	}

	URI.RawQuery = urlParams.Encode()
	finalUrl := URI.String()

	client := &http.Client{}
	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		return nil, err
	}

	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	var response AliyunQueryResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != nil {
		return nil, errors.New(*response.Message)
	}

	if response.DomainRecords == nil || response.DomainRecords.Record == nil || len(response.DomainRecords.Record) == 0 || response.DomainRecords.Record[0] == nil {
		return nil, errors.New("empty query response")
	}

	return response.DomainRecords.Record[0], nil
}

func (a Aliyun) Sync(ipMeta *meta.IpMeta) error {
	recordData, err := a.query(ipMeta)

	params := make(map[string]string)

	params["RecordId"] = recordData.RecordId
	params["RR"] = ipMeta.Subdomain
	params["Type"] = "A"
	params["Value"] = *ipMeta.Ip
	params["TTL"] = "600"
	params["Priority"] = "1"
	params["Line"] = "default"

	urlParams := a.aliyunSign(ipMeta, "UpdateDomainRecord", "GET", params)

	URI, err := url.Parse(ALIYUN_URL)
	if err != nil {
		return err
	}

	URI.RawQuery = urlParams.Encode()
	finalUrl := URI.String()

	client := &http.Client{}
	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		return err
	}

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	var response AliyunQueryResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Code != nil {
		return errors.New(*response.Message)
	}

	return nil
}

func (t Aliyun) aliyunSign(ipMeta *meta.IpMeta, action string, method string, params map[string]string) *url.Values {
	params["Action"] = action
	params["Version"] = "2015-01-09"
	params["Format"] = "JSON"
	params["Timestamp"] = time.Now().UTC().Format(time.RFC3339)
	params["AccessKeyId"] = ipMeta.Accesskey
	params["SignatureMethod"] = "HMAC-SHA1"
	params["SignatureVersion"] = "1.0"
	params["SignatureNonce"] = uuid.NewString()

	keys := make([]string, 0, len(params))

	for key, _ := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	canonicalRequestBuilder := bytes.Buffer{}

	for idx, key := range keys {
		canonicalRequestBuilder.WriteString(url.QueryEscape(key) + "=" + url.QueryEscape(params[key]))
		if idx != len(keys)-1 {
			canonicalRequestBuilder.WriteString("&")
		}
	}

	canonicalRequest := canonicalRequestBuilder.String()

	string2sign := fmt.Sprintf("%s&%s&%s", method, url.QueryEscape("/"), url.QueryEscape(canonicalRequest))

	signature := base64.StdEncoding.EncodeToString(t.hmacsha1([]byte(string2sign), []byte(ipMeta.AccessKeySecret+"&")))

	params["Signature"] = signature

	urlParams := url.Values{}

	for key, value := range params {
		urlParams.Add(key, value)
	}

	return &urlParams
}

func (Aliyun) hmacsha1(content []byte, key []byte) []byte {
	hashed := hmac.New(sha1.New, key)
	hashed.Write(content)
	return hashed.Sum(nil)
}
