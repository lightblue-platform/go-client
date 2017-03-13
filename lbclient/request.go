package lbclient

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const MAXRANGE = 4294967295

var EMPTYRANGE = Range{0, -1}
var ALLRANGE = Range{0, MAXRANGE}

type Range struct {
	from, to int
}

func NewRange(f, t int) *Range {
	r := Range{from: f, to: t}
	return &r
}

func RangeFrom(f int) *Range {
	r := Range{from: f, to: MAXRANGE}
	return &r
}

func RangeTo(t int) *Range {
	r := Range{from: 0, to: t}
	return &r
}

func (r Range) IsAll() bool {
	return r.from == 0 && r.to == MAXRANGE
}

func (r Range) MarshalJSON() ([]byte, error) {
	var to string

	if r.to == MAXRANGE {
		to = "null"
	} else {
		to = strconv.Itoa(r.to)
	}

	return []byte(fmt.Sprintf("[%d,%s]", r.from, to)), nil
}

type RequestHeader struct {
	EntityName       string
	EntityVersion    string
	ClientId         interface{}
	ExecutionOptions interface{}
}

type FindRequest struct {
	RequestHeader
	Q Query
	P Projection
	S Sort
	R Range
}

type projectionAndRange interface {
	getProjection() *Projection
	getRange() *Range
}

func (f *FindRequest) getProjection() *Projection { return &f.P }
func (f *FindRequest) getRange() *Range           { return &f.R }

func (r *RequestHeader) marshal(m map[string]interface{}) {
	if len(r.EntityName) != 0 {
		m["entity"] = r.EntityName
	}
	if len(r.EntityVersion) != 0 {
		m["entityVersion"] = r.EntityVersion
	}
	if r.ClientId != nil {
		m["client"] = r.ClientId
	}
	if r.ExecutionOptions != nil {
		m["execution"] = r.ExecutionOptions
	}
}

func marshalProjectionAndRange(r projectionAndRange, m map[string]interface{}) {
	if !r.getProjection().Empty() {
		m["projection"] = *r.getProjection()
	}
	if !r.getRange().IsAll() {
		m["range"] = *r.getRange()
	}
}

func (f *FindRequest) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	f.RequestHeader.marshal(m)
	marshalProjectionAndRange(f, m)
	if !f.Q.Empty() {
		m["query"] = f.Q
	}
	if !f.S.Empty() {
		m["sort"] = f.S
	}
	return json.Marshal(m)
}

type InsertRequest struct {
	RequestHeader
	P    Projection
	R    Range
	Docs []byte
}

type SaveRequest struct {
	RequestHeader
	P                Projection
	R                Range
	Docs             []byte
	Upsert           bool
	IfCurrentOnly    bool
	DocumentVersions []string
}

type DeleteRequest struct {
	RequestHeader
	Q Query
}

type UpdateRequest struct {
	RequestHeader
	Q                Query
	U                Update
	P                Projection
	R                Range
	IfCurrentOnly    bool
	DocumentVersions []string
}
