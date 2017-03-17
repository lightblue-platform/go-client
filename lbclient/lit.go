package lbclient

import (
	"encoding/json"
	"fmt"
	"time"
)

// DATEFORMAT is the date format used by lightblue server
const DATEFORMAT = "20060102T15:04:05.000-0700"

// Literal is an opaque structure for literal values
type Literal struct {
	i     int
	s     string
	d     float32
	b     bool
	j     []byte
	which int
}

// LitInt makes a literal integer
func LitInt(n int) Literal {
	return Literal{i: n, which: 0}
}

// LitInts returns a literal int array
func LitInts(values ...int) []Literal {
	ret := make([]Literal, len(values))
	for i, v := range values {
		ret[i] = LitInt(v)
	}
	return ret
}

// LitStr makes a literal string
func LitStr(str string) Literal {
	return Literal{s: str, which: 1}
}

// LitStrs returns a literal string array
func LitStrs(str ...string) []Literal {
	ret := make([]Literal, len(str))
	for i, s := range str {
		ret[i] = LitStr(s)
	}
	return ret
}

// LitDouble makes a literal double number
func LitDouble(n float32) Literal {
	return Literal{d: n, which: 2}
}

// LitDoubles returns a literal double array
func LitDoubles(values ...float32) []Literal {
	ret := make([]Literal, len(values))
	for i, v := range values {
		ret[i] = LitDouble(v)
	}
	return ret
}

// LitBool makes a literal boolean
func LitBool(v bool) Literal {
	return Literal{b: v, which: 3}
}

// LitBools returns a literal boolean array
func LitBools(values ...bool) []Literal {
	ret := make([]Literal, len(values))
	for i, b := range values {
		ret[i] = LitBool(b)
	}
	return ret
}

// LitDate makes a literal date
func LitDate(t time.Time) Literal {
	return Literal{s: t.Format(DATEFORMAT), which: 1}
}

// LitDates returns a literal date array
func LitDates(values ...time.Time) []Literal {
	ret := make([]Literal, len(values))
	for i, v := range values {
		ret[i] = LitDate(v)
	}
	return ret
}

// LitJson makes a literal JSON value
func LitJson(v []byte) Literal {
	return Literal{j: v, which: 4}
}

// LitNull makes a null value
func LitNull() Literal {
	return Literal{which: -1}
}

// MarshalJSON returns JSON representation of a  literal value
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

// String returns a string representation of literal value
func (l Literal) String() string {
	b, _ := l.MarshalJSON()
	return string(b)
}

// ValueOf represents a {valueof: field} construct. Use ValueOfField to instantiate
type ValueOf struct {
	field string
}

// ValueOfField makes a {valueof:f} construct
func ValueOfField(f string) ValueOf {
	return ValueOf{field: f}
}

// RValue is either a literal value, or a {valueof:field}
type RValue interface {

	// IsLiteral returns if this is a literal, false if this is a {valueof:f}
	IsLiteral() bool
	// GetLiteral returns the literal value, or panic if this is a {valueof:f}
	GetLiteral() Literal
	// GetValueOfField returns the field of {valueof:f}, or panics if this is  a literal
	GetValueOfField() string
	// String returns the string representation of RValue
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
