package lbclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

type PEMCertAuthConfig struct {
	// Certificate file, PEM
	PEMCertFile string
	// Private key file, PEM
	PEMPrivateKeyFile string
}

type AuthConfig interface {
	BuildTransport(t *http.Transport) error
}

type HttpClientConfig struct {
	DataServiceURI     string
	MetadataServiceURI string
	AuthConfig
	Compression      string
	ReadPreference   string
	WriteConcern     string
	MaxQueryTimeMS   int
	ExecutionOptions interface{}
}

type HttpClient struct {
	Config    *HttpClientConfig
	Transport *http.Transport
	Client    *http.Client
}

func (c *PEMCertAuthConfig) BuildTransport(t *http.Transport) error {
	cert, err := tls.LoadX509KeyPair(c.PEMCertFile, c.PEMPrivateKeyFile)
	if err != nil {
		return err
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	tlsConfig.BuildNameToCertificate()
	t.TLSClientConfig = tlsConfig
	return nil
}

func NewHttpClient(config *HttpClientConfig) *HttpClient {
	var cli HttpClient
	cli.Config = config
	cli.Transport = &http.Transport{}
	if cli.Config.AuthConfig != nil {
		if err := cli.Config.AuthConfig.BuildTransport(cli.Transport); err != nil {
			panic("Cannot build transport:" + err.Error())
		}
	}
	cli.Client = &http.Client{Transport: cli.Transport}
	return &cli
}

func (c *HttpClient) Call(url *url.URL, httpMethod string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(httpMethod, url.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rdbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return rdbody, nil
}

type marshalResponse struct {
	EntityName     string   `json:"entity"`
	EntityVersion  string   `json:"entityVersion"`
	HostName       string   `json:"hostname"`
	Status         OpStatus `json:"status"`
	ModifiedCount  int      `json:"modifiedCount"`
	MatchCount     int      `json:"matchCount"`
	TaskHandle     string
	Session        string
	ResultMetadata []map[string]interface{} `json:"resultMetadata"`
	DataErrors     []map[string]interface{} `json:"dataErrors"`
	Errors         []map[string]interface{} `json:"errors"`
	EntityData     interface{}              `json:"processed"`
}

func (c *HttpClient) DataCall(header *RequestHeader, body []byte, data reflect.Type, operation string, httpMethod string) (*Response, error) {
	b := bytes.Buffer{}
	b.WriteString(c.Config.DataServiceURI)
	if c.Config.DataServiceURI[len(c.Config.DataServiceURI)-1] != '/' {
		b.WriteRune('/')
	}
	b.WriteString(operation)
	b.WriteRune('/')
	b.WriteString((*header).EntityName)
	if len((*header).EntityVersion) > 0 {
		b.WriteRune('/')
		b.WriteString((*header).EntityVersion)
	}
	url, err := url.Parse(b.String())
	if err != nil {
		panic("Invalid URI:" + b.String())
	}
	responseBody, err := c.Call(url, httpMethod, body)
	if err != nil {
		return nil, err
	}

	var mr marshalResponse
	singleResult := false
	if data != nil {
		if data.Kind() != reflect.Slice {
			mr.EntityData = reflect.New(reflect.SliceOf(data)).Interface()
			singleResult = true
		} else {
			mr.EntityData = reflect.New(data).Interface()
		}
	}
	err = json.Unmarshal(responseBody, &mr)
	if err != nil {
		return nil, err
	}
	var response Response
	response.EntityName = mr.EntityName
	response.EntityVersion = mr.EntityVersion
	response.HostName = mr.HostName
	response.Status = mr.Status
	response.ModifiedCount = mr.ModifiedCount
	response.MatchCount = mr.MatchCount
	response.TaskHandle = mr.TaskHandle
	response.Session = mr.Session
	response.Errors = unmarshalErrors(mr.Errors)
	response.DataErrors = unmarshalDataErrors(mr.DataErrors)
	response.ResultMetadata = unmarshalRmd(mr.ResultMetadata)
	if mr.EntityData != nil {
		t, ok := mr.EntityData.([]map[string]interface{})
		if ok {
			response.EntityData = t
		} else if data == nil {
			response.EntityData = mr.EntityData
		} else {
			// We are using allocated slice
			ed := reflect.ValueOf(mr.EntityData).Elem().Interface()
			// If single result is expected, make sure there is at most one
			if singleResult {
				edValue := reflect.ValueOf(ed)
				if edValue.Len() > 1 {
					return nil, errors.New("More than one results for a non-array resultset")
				}
				response.EntityData = edValue.Index(0).Interface()
			} else {
				response.EntityData = ed
			}
		}
	}
	return &response, nil
}

func unmarshalErrors(errors []map[string]interface{}) []RequestError {
	if errors == nil || len(errors) == 0 {
		return nil
	}
	ret := make([]RequestError, len(errors))
	for i, err := range errors {
		ret[i] = UnmarshalError(err)
	}
	return ret
}

func UnmarshalError(error map[string]interface{}) RequestError {
	return RequestError{Context: error["context"].(string),
		ErrorCode: error["errorCode"].(string),
		Msg:       error["msg"].(string)}
}

func unmarshalDataErrors(errors []map[string]interface{}) []DataError {
	if errors == nil || len(errors) == 0 {
		return nil
	}
	ret := make([]DataError, len(errors))
	for i, err := range errors {
		ret[i] = UnmarshalDataError(err)
	}
	return ret
}

func UnmarshalDataError(error map[string]interface{}) DataError {
	return DataError{Errors: unmarshalErrors(error["errors"].([]map[string]interface{})),
		EntityData: error["entityData"].([]map[string]interface{})}
}

func unmarshalRmd(rmd []map[string]interface{}) []ResultMd {
	if rmd == nil || len(rmd) == 0 {
		return nil
	}
	ret := make([]ResultMd, len(rmd))
	for i, r := range rmd {
		ret[i] = UnmarshalRmd(r)
	}
	return ret
}

func UnmarshalRmd(r map[string]interface{}) ResultMd {
	return ResultMd{DocumentVersion: r["documentVersion"].(string)}
}

func (c *HttpClient) Find(request *FindRequest, data interface{}) (*Response, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	var t reflect.Type
	if data != nil {
		var ok bool
		t, ok = data.(reflect.Type)
		if !ok {
			t = reflect.TypeOf(data)
		}
	}
	return c.DataCall(&request.RequestHeader, body, t, "find", "POST")
}
