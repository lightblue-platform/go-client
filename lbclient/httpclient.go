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
	"strconv"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type CrudOperation string

var CRUD_UPDATE CrudOperation = "update"
var CRUD_FIND CrudOperation = "find"
var CRUD_DELETE CrudOperation = "delete"
var CRUD_SAVE CrudOperation = "save"
var CRUD_INSERT CrudOperation = "insert"

// PEMCertAuthConfig implements the AuthConfig interface for a
// certificate and private key encoded in PEM format in files. The
// certificate and the private key can be in the same file, or in
// different files. The private key must be an unencrypted private key
//
type PEMCertAuthConfig struct {
	// Certificate file, PEM
	PEMCertFile string
	// Unencrypted private key file, PEM, can be the same file as the PEMCertFile
	PEMPrivateKeyFile string
}

// AuthConfig interface defines the BuildTransport method that
// configures the HTTP transport using the authorization scheme
type AuthConfig interface {
	BuildTransport(t *http.Transport) error
}

// HttpClientConfig is the client configuration for the HTTP client
type HttpClientConfig struct {
	// The URI for the data service, it should contain the host, port,
	// and the context root for the CRUD app
	DataServiceURI string
	// The URI for the metadata service, it should contain the host,
	// port, and the context root for the metadata app
	MetadataServiceURI string
	// Authoziation configuration implementation
	AuthConfig
	// Default read preference for the client
	ReadPreference string
	// Default write concernt for the client
	WriteConcern string
	// Maximum query time
	MaxQueryTimeMS int
	// Execution options
	ExecutionOptions interface{}
}

// HttpClient can be initialised once, and shared by multiple threads
type HttpClient struct {
	Config    *HttpClientConfig
	Transport *http.Transport
	Client    *http.Client
}

// BuildTransport builds the transport using the certificate and
// unencrypted private key
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

// NewHttpClient creates and initialized a new client using the HttpClientConfig
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

// Call performs an HTTP method on the url, and returns the result
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

// DataCall performs a data service call
//
//  * entityName is the entity name for the call
//  * entityVersion is empty (for default version), or the required entity version
//  * body is the JSON-marshaled request body
//  * returnDataType is the type of the EntityData to be returned in the response
//  * operation is the CRUD operation to be performed,
//  * httpMethod GET, POST, etc
//
// Returns the response, and if there is an error or not. If returnDataType is a struct, then
// the JSON documents are unmarshaled to that type. Othrwise, the documents are returned
// as a slice/map tree
func (c *HttpClient) DataCall(entityName, entityVersion string, body []byte, returnDataType reflect.Type, operation CrudOperation, httpMethod string) (*Response, error) {
	b := bytes.Buffer{}
	b.WriteString(c.Config.DataServiceURI)
	if c.Config.DataServiceURI[len(c.Config.DataServiceURI)-1] != '/' {
		b.WriteRune('/')
	}
	b.WriteString(string(operation))
	b.WriteRune('/')
	b.WriteString(entityName)
	if len(entityVersion) > 0 {
		b.WriteRune('/')
		b.WriteString(entityVersion)
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
	if returnDataType != nil {
		if returnDataType.Kind() != reflect.Slice {
			mr.EntityData = reflect.New(reflect.SliceOf(returnDataType)).Interface()
			singleResult = true
		} else {
			mr.EntityData = reflect.New(returnDataType).Interface()
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
		} else if returnDataType == nil {
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

func (c *HttpClient) docCall(request interface{}, data interface{}, entityName, entityVersion string, op CrudOperation, mth string) (*Response, error) {
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
	return c.DataCall(entityName, entityVersion, body, t, op, mth)
}

// Find issues a find request to the server
//
//   * request The Find request
//   * data If nil, the documents will be returned as slice/map tree.
//     If data is a reflect.Type, then the documents will be unmarshaled as a slice
//     of that type. If data is any other struct, then the result documents will be
//     unmarshaled of that type
func (c *HttpClient) Find(request *FindRequest, data interface{}) (*Response, error) {
	return c.docCall(request, data, request.EntityName, request.EntityVersion, CRUD_FIND, POST)
}

// Insert adds documents to a database
//
//   * request The insert request
//   * returnData This can be JSON marshalable object (or array of such objects), or
//     a map[string]interface{} (or an array of it)
func (c *HttpClient) Insert(request *InsertRequest, returnData interface{}) (*Response, error) {
	return c.docCall(request, returnData, request.EntityName, request.EntityVersion, CRUD_INSERT, PUT)
}

// Save adds/updates documents in a database
//
//   * request The save request
//   * returnData This can be JSON marshalable object (or array of such objects), or
//     a map[string]interface{} (or an array of it)
func (c *HttpClient) Save(request *SaveRequest, returnData interface{}) (*Response, error) {
	return c.docCall(request, returnData, request.EntityName, request.EntityVersion, CRUD_SAVE, POST)
}

// Update modifies documents in a database
//
//   * request The update request
//   * returnData This can be JSON marshalable object (or array of such objects), or
//     a map[string]interface{} (or an array of it)
func (c *HttpClient) Update(request *UpdateRequest, returnData interface{}) (*Response, error) {
	return c.docCall(request, returnData, request.EntityName, request.EntityVersion, CRUD_UPDATE, POST)
}

// Delete removes documents from a database
//
//   * request The delete request
func (c *HttpClient) Delete(request *DeleteRequest) (*Response, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return c.DataCall(request.EntityName, request.EntityVersion, body, nil, CRUD_DELETE, POST)
}

func buildLockRequest(operation, domain, callerId, resourceId string) map[string]string {
	return map[string]string{"operation": operation,
		"domain":     domain,
		"callerId":   callerId,
		"resourceId": resourceId}
}

func (c *HttpClient) lock(req map[string]string) ([]byte, error) {
	b := bytes.Buffer{}
	b.WriteString(c.Config.DataServiceURI)
	if c.Config.DataServiceURI[len(c.Config.DataServiceURI)-1] != '/' {
		b.WriteRune('/')
	}
	b.WriteString("lock")
	url, err := url.Parse(b.String())
	if err != nil {
		panic("Invalid URI:" + b.String())
	}
	body, _ := json.Marshal(req)
	return c.Call(url, POST, body)
}

func parseLockResult(body []byte, err error) (string, error) {
	if err == nil {
		var resultMap map[string]interface{}
		err = json.Unmarshal(body, &resultMap)
		if err == nil {
			if resultMap["status"] == "ERROR" {
				return "", UnmarshalError(resultMap["errors"].([]map[string]interface{})[0])
			} else {
				return resultMap["result"].(string), nil
			}
		}
	}
	return "", err
}

func (c *HttpClient) Acquire(domain, callerId, resourceId string, ttl int) (bool, error) {
	req := buildLockRequest("acquire", domain, callerId, resourceId)
	if ttl > 0 {
		req["ttl"] = strconv.Itoa(ttl)
	}
	res, err := parseLockResult(c.lock(req))
	if err == nil {
		return res == "true", nil
	} else {
		return false, err
	}
}

func (c *HttpClient) Release(domain, callerId, resourceId string) (bool, error) {
	res, err := parseLockResult(c.lock(buildLockRequest("release", domain, callerId, resourceId)))
	if err == nil {
		return res == "true", nil
	} else {
		return false, err
	}
}

func (c *HttpClient) GetLockCount(domain, callerId, resourceId string) (int, error) {
	res, err := parseLockResult(c.lock(buildLockRequest("count", domain, callerId, resourceId)))
	if err == nil {
		return strconv.Atoi(res)
	} else {
		return -1, err
	}
}

func (c *HttpClient) Ping(domain, callerId, resourceId string) (bool, error) {
	res, err := parseLockResult(c.lock(buildLockRequest("ping", domain, callerId, resourceId)))
	if err == nil {
		return res == "true", nil
	} else {
		return false, err
	}
}
