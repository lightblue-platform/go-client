package lbclient

import (
	"encoding/json"
)

// Projection is an opaque structure containing zero or more
// projectionPart objects. To make a Projection, use MakeProjection
// function
type Projection struct {
	p []projectionPart
}

// Empty returns true if projection is empty
func (p *Projection) Empty() bool {
	return p.p == nil || len(p.p) == 0
}

func (p *Projection) AsMap() interface{} {
	var x []map[string]interface{}
	x = make([]map[string]interface{}, len(p.p))
	for i, item := range p.p {
		x[i] = item.GetProjection()
	}
	return x
}

// MarshalJSON returns a JSON representation of the projection
func (p Projection) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.AsMap())
}

// The internal representation of each part of a projection
type projectionPart interface {
	GetProjection() map[string]interface{}
}

type baseProjection struct {
	field   string
	include bool
}

// FieldProjection contains a field, include flag, and recursive flag
type FieldProjection struct {
	baseProjection
	recursive bool
}

// GetProjection returns a nap representation of a FieldProjection
func (p FieldProjection) GetProjection() map[string]interface{} {
	return map[string]interface{}{"field": p.field,
		"include":   p.include,
		"recursive": p.recursive}
}

// MakeProjections makes a Projection object fron parts
func MakeProjection(parts ...projectionPart) *Projection {
	var p Projection
	p.Add(parts...)
	return &p
}

// Add adds new parts to a Projection object
func (p *Projection) Add(parts ...projectionPart) {
	if p.p == nil {
		p.p = make([]projectionPart, len(parts))
		for i, part := range parts {
			p.p[i] = part
		}
	} else {
		p.p = append(p.p, parts...)
	}
}

// ProjectField returns a projection of the form
//
//    { field: fld, include: inc, recursive: rec}
//
func ProjectField(fld string, inc bool, rec bool) FieldProjection {
	var f FieldProjection
	f.field = fld
	f.include = inc
	f.recursive = rec
	return f
}

// IncludeField returns a projection of the form
//
//   { field: fld, include: true, recursive: rec}
//
func IncludeField(fld string, rec bool) FieldProjection {
	return ProjectField(fld, true, rec)
}

// IncludeTree returns a projection of the form
//
//   { field: fld, include: true, recursive: true}
//
func IncludeTree(fld string) FieldProjection {
	return IncludeField(fld, true)
}

// ExcludeField returns a projection of the form
//
//   { field: fld, include: false, recursive: rec}
//
func ExcludeField(fld string, rec bool) FieldProjection {
	return ProjectField(fld, false, rec)
}

// ExcludeTree returns a projection of the form
//
//  { field: fld, include: false, recursive: true }
//
func ExcludeTree(fld string) FieldProjection {
	return ExcludeField(fld, true)
}

// ArrayProjection is the common base for RangeProjection and
// MatchProjection. Contains field, include flag, sort and projection
// fields
type ArrayProjection struct {
	baseProjection
	s *Sort
	p *Projection
}

// RangeProjection projects a given range of an array
type RangeProjection struct {
	ArrayProjection
	rng [2]int
}

// ProjectRange returns a projection that projects a range of an array
//
//   { field: fld, include: inc, range: [r], projection: *projection, sort: *sort }
//
// The projection and sort fields are optional, and can be nil
func ProjectRange(fld string, inc bool, r [2]int, projection *Projection, sort *Sort) RangeProjection {
	var p RangeProjection
	p.field = fld
	p.include = inc
	p.rng = r
	p.s = sort
	p.p = projection
	return p
}

// IncludeRange returns a range projection that includes the given range
func IncludeRange(fld string, r [2]int, projection *Projection, sort *Sort) RangeProjection {
	return ProjectRange(fld, true, r, projection, sort)
}

// ExcludeRange returns a range projection that excludes the given range
func ExcludeRange(fld string, r [2]int, projection *Projection, sort *Sort) RangeProjection {
	return ProjectRange(fld, false, r, projection, sort)
}

// GetProjection returns the map representation of the range projection
func (p RangeProjection) GetProjection() map[string]interface{} {
	var ret map[string]interface{}
	ret["field"] = p.field
	ret["include"] = p.include
	ret["range"] = p.rng
	if p.s != nil {
		ret["sort"] = p.s
	}
	if p.p != nil {
		ret["projection"] = p.p
	}
	return ret
}

// MatchProjection contains an array elemMatch projection that
// projects elements of an array that matches a query
type MatchProjection struct {
	ArrayProjection
	q Query
}

// ProjectMatching returns an elemMatch projection
//
//   {field: fld, include: inc, match: q, projection: *projection, sort: *sort }
//
// Projection and sort fields are optional, and can be nil
func ProjectMatching(fld string, inc bool, match Query, projection *Projection, sort *Sort) MatchProjection {
	var p MatchProjection
	p.field = fld
	p.include = inc
	p.q = match
	p.s = sort
	p.p = projection
	return p
}

// IncludeMatching returns an array match projection that includes the array elements that match the query
func IncludeMatching(fld string, match Query, projection *Projection, sort *Sort) MatchProjection {
	return ProjectMatching(fld, true, match, projection, sort)
}

// ExcludeMatching returns an array match projection that excludes the array elements that match the query
func ExcludeMatching(fld string, match Query, projection *Projection, sort *Sort) MatchProjection {
	return ProjectMatching(fld, false, match, projection, sort)
}

// GetProjection returns the map representation of the match projection
func (p MatchProjection) GetProjection() map[string]interface{} {
	var ret map[string]interface{}
	ret["field"] = p.field
	ret["include"] = p.include
	ret["match"] = p.q
	if p.s != nil {
		ret["sort"] = p.s
	}
	if p.p == nil {
		ret["projection"] = p.p
	}
	return ret
}
