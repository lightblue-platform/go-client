package lbclient

import (
	"fmt"
)

type OpStatus string

const (
	COMPLETE OpStatus = "COMPLETE"
	PARTIAL  OpStatus = "PARTIAL"
	ASYNC    OpStatus = "ASYNC"
	ERROR    OpStatus = "ERROR"
)

type Response struct {
	EntityName     string   `json:"entity"`
	EntityVersion  string   `json:"entityVersion"`
	HostName       string   `json:"hostname"`
	Status         OpStatus `json:"status"`
	ModifiedCount  int      `json:"modifiedCount"`
	MatchCount     int      `json:"matchCount"`
	TaskHandle     string
	Session        string
	ResultMetadata []ResultMd               `json:"resultMetadata"`
	EntityData     []map[string]interface{} `json:"processed"`
	DataErrors     []DataError              `json:"dataErrors"`
	Errors         []RequestError           `json:"errors"`
}

func (r *Response) String() string {
	return fmt.Sprintf("entity: %s, ver: %s, hostname: %s, status: %s, matchCount: %d,"+
		" modifiedCount: %d, rmd: %q, entityData: %q, dataErrors: %q, errors: %q",
		r.EntityName, r.EntityVersion, r.HostName, r.Status, r.MatchCount, r.ModifiedCount,
		r.ResultMetadata, r.EntityData, r.DataErrors, r.Errors)
}

type ResultMd struct {
	DocumentVersion string `json:"documentVersion"`
}

func (r *ResultMd) String() string {
	return fmt.Sprintf("docver: %s", r.DocumentVersion)
}

type DataError struct {
	EntityData []map[string]interface{} `json:"entityData"`
	Errors     []RequestError           `json:"errors"`
}

type RequestError struct {
	Context   string `json:"context"`
	ErrorCode string `json:"errorCode"`
	Msg       string `json:"msg"`
}

func (r *RequestError) String() string {
	return fmt.Sprintf("ctx: %s, err: %s, msg: %s", r.Context, r.ErrorCode, r.Msg)
}
