package lbclient

import (
	"encoding/json"
	"fmt"
	"time"
)

const DATEFORMAT = "20060102T15:04:05.000-0700"

// Opaque structure for literal values
type Literal struct {
	i     int
	s     string
	d     float32
	b     bool
	j     []byte
	which int
}

// Makes a literal integer
func LitInt(n int) Literal {
	return Literal{i: n, which: 0}
}

// Makes a literal string
func LitStr(str string) Literal {
	return Literal{s: str, which: 1}
}

// Makes a literal double number
func LitDouble(n float32) Literal {
	return Literal{d: n, which: 2}
}

// Makes a literal boolean
func LitBool(v bool) Literal {
	return Literal{b: v, which: 3}
}

// Makes a literal date
func LitDate(t time.Time) Literal {
	return Literal{s: t.Format(DATEFORMAT), which: 1}
}

// Makes a literal JSON value
func LitJson(v []byte) Literal {
	return Literal{j: v, which: 4}
}

// Makes a null value
func LitNull() Literal {
	return Literal{which: -1}
}

// Marshals a literal value to json
func (l Literal) MarshalJSON() ([]byte, error) {
	switch l.which {
	case 0:
		return json.Marshal(l.i)
	case 1:
		return json.Marshal(l.s)
	case 2:
		return json.Marshal(l.d)
	case 3:
		return json.Marshal(l.b)
	case 4:
		return l.j, nil
	default:
		return []byte("null"), nil
	}
}

// Returns a string representation of literal value
func (l Literal) String() string {
	b, _ := l.MarshalJSON()
	return string(b)
}

// Represents a {valueof: field} construct
type ValueOf struct {
	field string
}

// Makes a {valueof:f} construct
func ValueOfField(f string) ValueOf {
	return ValueOf{field: f}
}

// Either a literal value, or a {valueof:field}
type RValue interface {

	// Returns if this is a literal, false if this is a {valueof:f}
	IsLiteral() bool
	// Returns the literal value, or panic if this is a {valueof:f}
	GetLiteral() Literal
	// Returns the field of {valueof:f}, or panics if this is  a literal
	GetValueOfField() string

	String() string
}

// RValue.IsLiteral: returns true
func (l Literal) IsLiteral() bool {
	return true
}

// RValue.IsLiteral: returns false
func (v ValueOf) IsLiteral() bool {
	return false
}

// RValue.GetLiteral, returns the literal value
func (l Literal) GetLiteral() Literal {
	return l
}

// RValue.GetLiteral, panics
func (v ValueOf) GetLiteral() Literal {
	panic("GetLitaral called for ValueOf")
}

// RValue.GetValueOfField, panics
func (l Literal) GetValueOfField() string {
	panic("GetValueOfField called for literal")
}

// RValue.GetValueOfField, returns the field of {valueof:f}
func (v ValueOf) GetValueOfField() string {
	return v.field
}

// Marshals a {valueof:f} to JSON
func (r ValueOf) MarshalJSON() ([]byte, error) {
	return []byte(r.String()), nil
}

// Returns the string representation of {valueof:f}
func (r ValueOf) String() string {
	return fmt.Sprintf("{\"$valueof\":\"%s\"}", r.field)
}
