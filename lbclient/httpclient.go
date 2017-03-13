package lbclient

import (
	"bytes"
	"crypto/tls"
	//	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"golang.org/x/crypto/pkcs12"
	"io/ioutil"
	"net/http"
	"net/url"
)

type BasicAuthConfig struct {
	UserName string
	Password string
}

type CertAuthConfig struct {
	// Certificate file, PEM
	PEMCertFile string
	// Private key file,  PKCS12 if not contained in certificate
	PrivateKeyFile string
	CertAlias      string
	CertPassword   string
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

func (c *CertAuthConfig) BuildTransport(t *http.Transport) error {
}

func ReadPEMFile(file string) ([]pem.Block, error) {
	blocks := make([]pem.Block, 0)

	inp, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	done := false
	for !done {
		var blk *pem.Block
		blk, inp = pem.Decode(inp)
		if blk != nil {
			blocks = append(blocks, *blk)
		} else {
			done = true
		}
	}
	return blocks, nil
}

func ReadPKCS12File(file string, password string) ([]*pem.Block, error) {
	inp, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return pkcs12.ToPEM(inp, password)
}

func NewHttpClient(config *HttpClientConfig) *HttpClient {
	var cli HttpClient
	cli.Config = config
	cli.Transport = &http.Transport{}
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
	cli := &http.Client{}
	var transport http.Transport
	if c.Config.AuthConfig != nil {
		if err = c.Config.AuthConfig.BuildTransport(&transport); err != nil {
			return nil, err
		}
		cli.Transport = &transport
	}
	resp, err := cli.Do(req)
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
