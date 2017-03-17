package lbclient

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

const MAXRANGE = 4294967295

// Empty range, [0,-1]
var EMPTYRANGE = Range{0, -1}

// Range that covers everything, [0,MAXRANGE]
var ALLRANGE = Range{0, MAXRANGE}

// Ramge structure contains the resultset range, both from and to inclusive
// If from<to, then the range is empty
type Range struct {
	from, to int
}

// NewRange resutns a new Range with the give to/from values
func NewRange(f, t int) *Range {
	r := Range{from: f, to: t}
	return &r
}

// RangeFrom returns a range that starts from f to the end of the resultset
func RangeFrom(f int) *Range {
	r := Range{from: f, to: MAXRANGE}
	return &r
}

// RangeTo returns a range that starts from 0 up to and including t
func RangeTo(t int) *Range {
	r := Range{from: 0, to: t}
	return &r
}

// IsAll returns if the range is the ALLRANGE
func (r Range) IsAll() bool {
	return r.from == 0 && r.to == MAXRANGE
}

// Returns JSON representation of the range
func (r Range) MarshalJSON() ([]byte, error) {
	var to string

	if r.to == MAXRANGE {
		to = "null"
	} else {
		to = strconv.Itoa(r.to)
	}

	return []byte(fmt.Sprintf("[%d,%s]", r.from, to)), nil
}

// RequestHeader is the common portion of all Lightblue  requests.
type RequestHeader struct {
	EntityName       string
	EntityVersion    string
	ClientId         interface{}
	ExecutionOptions interface{}
}

type FindRequest struct {
	RequestHeader
	Q *Query
	P *Projection
	S *Sort
	R *Range
}

type projectionAndRange interface {
	getProjection() *Projection
	getRange() *Range
}

type MapConvertable interface {
	AsMap() interface{}
}

func (f *FindRequest) getProjection() *Projection { return f.P }

func (f *FindRequest) getRange() *Range { return f.R }

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
	if r.getProjection() != nil && !r.getProjection().Empty() {
		m["projection"] = *(r.getProjection())
	}
	if r.getRange() != nil && !r.getRange().IsAll() {
		m["range"] = *(r.getRange())
	}
}

// MarshalJSON returns the JSON representation of a Find request
func (f *FindRequest) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	f.RequestHeader.marshal(m)
	marshalProjectionAndRange(f, m)
	if f.Q != nil && !f.Q.Empty() {
		m["query"] = *(f.Q)
	}
	if f.S != nil && !f.S.Empty() {
		m["sort"] = *(f.S)
	}
	b, e := json.Marshal(m)
	return b, e
}

// InsertRequest is used to insert documents to the database.
type InsertRequest struct {
	RequestHeader
	P       *Projection
	R       *Range
	DocData json.RawMessage
}

// MakeDocData builds document data needed for InsertRequest and
// SaveRequest from actual objects, maps, or raw json ([]byte)
//
// * If v is nil, MakeDocData returns nil
// * If v is a []byte, it assumes the documents are already JSON encoded, and returns v
// * If v is a map[string]interface{} or []map[string]interface{}, it assumes the document or the documents
//   are converted to a msp, or maps, and marshals that into JSON
// * If v is a struct array, or a struct, then it marshals the structs and returns that. If there is only one
//   struct, it is placed into an array of 1 before JSON encoding
func MakeDocData(v interface{}) json.RawMessage {
	if v == nil {
		return nil
	} else if m, ok := v.(map[string]interface{}); ok {
		s := []map[string]interface{}{m}
		ret, err := json.Marshal(s)
		if err != nil {
			panic(err.Error())
		}
		return ret
	} else if m, ok := v.([]map[string]interface{}); ok {
		ret, err := json.Marshal(m)
		if err != nil {
			panic(err.Error())
		}
		return ret
	} else if m, ok := v.([]byte); ok {
		return json.RawMessage(m)
	} else {
		var sliceData interface{}
		if reflect.TypeOf(v).Kind() != reflect.Slice {
			sliceVal := reflect.New(reflect.SliceOf(reflect.TypeOf(v))).Elem()
			sliceVal = reflect.Append(sliceVal, reflect.ValueOf(v))
			sliceData = sliceVal.Interface()
		} else {
			sliceData = v
		}
		ret, err := json.Marshal(sliceData)
		if err != nil {
			panic(err.Error())
		}
		return ret
	}
}

func (r *InsertRequest) getProjection() *Projection { return r.P }

func (r *InsertRequest) getRange() *Range { return r.R }

func (r *InsertRequest) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	r.RequestHeader.marshal(m)
	marshalProjectionAndRange(r, m)
	m["data"] = r.DocData
	return json.Marshal(m)
}

type SaveRequest struct {
	RequestHeader
	P                *Projection
	R                *Range
	DocData          json.RawMessage
	Upsert           bool
	IfCurrentOnly    bool
	DocumentVersions []string
}

func (r *SaveRequest) getProjection() *Projection { return r.P }

func (r *SaveRequest) getRange() *Range { return r.R }

func (r *SaveRequest) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	r.RequestHeader.marshal(m)
	marshalProjectionAndRange(r, m)
	m["data"] = r.DocData
	m["upsert"] = r.Upsert
	if r.IfCurrentOnly {
		m["onlyIfCurrent"] = true
		m["documentVersions"] = r.DocumentVersions
	}
	return json.Marshal(m)
}

type DeleteRequest struct {
	RequestHeader
	Q *Query
}

func (r *DeleteRequest) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	r.RequestHeader.marshal(m)
	m["query"] = *(r.Q)
	return json.Marshal(m)
}

type UpdateRequest struct {
	RequestHeader
	Q                *Query
	U                *Update
	P                *Projection
	R                *Range
	IfCurrentOnly    bool
	DocumentVersions []string
}

func (r *UpdateRequest) getProjection() *Projection { return r.P }

func (r *UpdateRequest) getRange() *Range { return r.R }

func (r *UpdateRequest) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	r.RequestHeader.marshal(m)
	marshalProjectionAndRange(r, m)
	m["query"] = *(r.Q)
	m["update"] = *(r.U)
	if r.IfCurrentOnly {
		m["onlyIfCurrent"] = true
		m["documentVersions"] = r.DocumentVersions
	}
	return json.Marshal(m)
}
