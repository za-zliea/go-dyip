package dns

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"dyip-sync/meta"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Tencent struct {
}

func NewTencent() Tencent {
	return Tencent{}
}

type TencentRequest struct {
	Domain     string `json:"Domain"`
	Subdomain  string `json:"Subdomain"`
	RecordType string `json:"RecordType"`
}

type TencentQueryResponse struct {
	Response *TencentQueryResponseData `json:"Response"`
}

type TencentQueryResponseData struct {
	Error           *TencentErrorResponse   `json:"Error"`
	RequestId       string                  `json:"RequestId"`
	RecordCountInfo *TencentRecordCountInfo `json:"RecordCountInfo"`
	RecordList      []*TencentRecordData    `json:"RecordList"`
}

type TencentErrorResponse struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

type TencentRecordCountInfo struct {
	SubdomainCount int `json:"SubdomainCount"`
	TotalCount     int `json:"TotalCount"`
	ListCount      int `json:"ListCount"`
}

type TencentRecordData struct {
	RecordId      int64       `json:"RecordId"`
	Value         string      `json:"Value"`
	Status        string      `json:"Status"`
	UpdatedOn     string      `json:"UpdatedOn"`
	Name          string      `json:"Name"`
	Line          string      `json:"Line"`
	LineId        string      `json:"LineId"`
	Type          string      `json:"Type"`
	Weight        interface{} `json:"Weight"`
	MonitorStatus string      `json:"MonitorStatus"`
	Remark        string      `json:"Remark"`
	TTL           int         `json:"TTL"`
	MX            int         `json:"MX"`
}

type TencentRecordRequest struct {
	Domain       string `json:"Domain"`
	SubDomain    string `json:"SubDomain"`
	RecordType   string `json:"RecordType"`
	RecordLine   string `json:"RecordLine"`
	RecordLineId string `json:"RecordLineId"`
	Value        string `json:"Value"`
	TTL          int    `json:"TTL"`
	Status       string `json:"Status"`
	RecordId     int64  `json:"RecordId"`
}

type TencentRecordResponse struct {
	Response struct {
		Error     *TencentErrorResponse `json:"Error"`
		RequestId string                `json:"RequestId"`
		RecordId  int64                 `json:"RecordId"`
	} `json:"Response"`
}

const TENCENT_URL = "https://dnspod.tencentcloudapi.com"

func (t Tencent) Query(ipMeta *meta.IpMeta) (string, error) {
	recordData, err := t.query(ipMeta)

	if err != nil {
		return "", err
	}

	return recordData.Value, nil
}

func (t Tencent) query(ipMeta *meta.IpMeta) (*TencentRecordData, error) {
	data := TencentRequest{Domain: ipMeta.Domain, Subdomain: ipMeta.Subdomain, RecordType: "A"}

	reqBody, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", TENCENT_URL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	t.tencentSign(ipMeta, "DescribeRecordList", reqBody, req)

	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return nil, errors.New(rsp.Status)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	var response TencentQueryResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Response != nil && response.Response.Error != nil {
		return nil, errors.New(response.Response.Error.Message)
	}

	if response.Response.RecordList == nil || len(response.Response.RecordList) == 0 || response.Response.RecordList[0] == nil {
		return nil, errors.New("empty query response")
	}

	return response.Response.RecordList[0], nil
}

func (t Tencent) Sync(ipMeta *meta.IpMeta) error {
	recordData, err := t.query(ipMeta)

	if err != nil {
		return err
	}

	data := TencentRecordRequest{
		RecordId:     recordData.RecordId,
		Domain:       ipMeta.Domain,
		SubDomain:    ipMeta.Subdomain,
		RecordType:   "A",
		RecordLine:   "默认",
		RecordLineId: "0",
		Value:        *ipMeta.Ip,
		TTL:          600,
		Status:       "ENABLE",
	}

	reqBody, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", TENCENT_URL, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	t.tencentSign(ipMeta, "ModifyRecord", reqBody, req)

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return errors.New(rsp.Status)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	var response TencentQueryResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Response != nil && response.Response.Error != nil {
		return errors.New(response.Response.Error.Message)
	}

	return nil
}

func (Tencent) sha256hex(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}

func (Tencent) hmacsha256(content []byte, key []byte) []byte {
	hashed := hmac.New(sha256.New, key)
	hashed.Write(content)
	return hashed.Sum(nil)
}

func (t Tencent) tencentSign(ipMeta *meta.IpMeta, action string, payload []byte, request *http.Request) {
	host := "dnspod.tencentcloudapi.com"
	algorithm := "TC3-HMAC-SHA256"
	service := "dnspod"
	version := "2021-03-23"
	contentType := "application/json"
	var timestamp int64 = time.Now().Unix()

	httpRequestMethod := "POST"
	canonicalURI := "/"
	canonicalQueryString := ""
	canonicalHeaders := "content-type:application/json\n" + "host:" + host + "\n"
	signedHeaders := "content-type;host"
	payloadHash := t.sha256hex(payload)

	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		httpRequestMethod,
		canonicalURI,
		canonicalQueryString,
		canonicalHeaders,
		signedHeaders,
		payloadHash)

	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, service)
	hashedCanonicalRequest := t.sha256hex([]byte(canonicalRequest))
	string2sign := fmt.Sprintf("%s\n%d\n%s\n%s",
		algorithm,
		timestamp,
		credentialScope,
		hashedCanonicalRequest)

	secretDate := t.hmacsha256([]byte(date), []byte("TC3"+ipMeta.AccessKeySecret))
	secretService := t.hmacsha256([]byte(service), secretDate)
	secretSigning := t.hmacsha256([]byte("tc3_request"), secretService)
	signature := hex.EncodeToString([]byte(t.hmacsha256([]byte(string2sign), secretSigning)))

	authorization := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		algorithm,
		ipMeta.Accesskey,
		credentialScope,
		signedHeaders,
		signature)

	request.Header.Set("Host", host)
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("X-TC-Action", action)
	request.Header.Set("X-TC-Version", version)
	request.Header.Set("X-TC-Timestamp", strconv.FormatInt(timestamp, 10))
	request.Header.Set("Authorization", authorization)
}
