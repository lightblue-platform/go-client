package lbclient

import ()

type OpStatus string

const (
	COMPLETE OpStatus = "COMPLETE"
	PARTIAL  OpStatus = "PARTIAL"
	ASYNC    OpStatus = "ASYNC"
	ERROR    OpStatus = "ERROR"
)

type Response struct {
	EntityName     string
	EntityVersion  string
	HostName       string
	Status         OpStatus
	ModifiedCount  int
	MatchCount     int
	TaskHandle     string
	Session        string
	ResultMetadata []ResultMd
	EntityData     []byte
	DataErrors     []DataError
	Errors         []RequestError
}

type ResultMd struct {
	DocumentVersion string `json:"documentVersion"`
}

type DataError struct {
	EntityData []byte
	Errors     []RequestError
}

type RequestError struct {
	Context   string `json:"context"`
	ErrorCode string `json:"errorCode"`
	Msg       string `json:"msg"`
}
