package lbclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
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

func (c *HttpClient) DataCall(header *RequestHeader, body []byte, operation string, httpMethod string) (*Response, error) {
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
	var response Response
	json.Unmarshal(responseBody, &response)
	return &response, nil
}

func (c *HttpClient) Find(request *FindRequest) (*Response, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return c.DataCall(&request.RequestHeader, body, "find", "POST")
}
