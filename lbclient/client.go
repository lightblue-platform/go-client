package lbclient

import (
	"net/url"
)

type LBClient interface {
	Call(url *url.URL, httpMethod string, body []byte) ([]byte, error)
	DataCall(header *RequestHeader, body []byte, operation string, httpMethod string) (*Response, error)
	Find(request *FindRequest) (*Response, error)
}
