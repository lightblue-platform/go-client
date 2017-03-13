package lbclient

import (
	"encoding/json"
)

type Projection struct {
	p []projectionPart
}

func (p *Projection) Empty() bool {
	return p.p == nil || len(p.p) == 0
}

func (p Projection) MarshalJSON() ([]byte, error) {
	var x []map[string]interface{}
	x = make([]map[string]interface{}, len(p.p))
	for i, item := range p.p {
		x[i] = item.GetProjection()
	}
	return json.Marshal(x)
}

type projectionPart interface {
	GetProjection() map[string]interface{}
}

type baseProjection struct {
	field   string
	include bool
}

type FieldProjection struct {
	baseProjection
	recursive bool
}

func (p FieldProjection) GetProjection() map[string]interface{} {
	return map[string]interface{}{"field": p.field,
		"include":   p.include,
		"recursive": p.recursive}
}

func MakeProjection(parts ...projectionPart) Projection {
	var p Projection
	p.Add(parts...)
	return p
}

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

func ProjectField(fld string, inc bool, rec bool) FieldProjection {
	var f FieldProjection
	f.field = fld
	f.include = inc
	f.recursive = rec
	return f
}

func IncludeField(fld string, rec bool) FieldProjection {
	return ProjectField(fld, true, rec)
}

func IncludeTree(fld string) FieldProjection {
	return IncludeField(fld, true)
}

func ExcludeField(fld string, rec bool) FieldProjection {
	return ProjectField(fld, false, rec)
}

func ExcludeTree(fld string) FieldProjection {
	return ExcludeField(fld, true)
}

type ArrayProjection struct {
	baseProjection
	s *Sort
	p *Projection
}

type RangeProjection struct {
	ArrayProjection
	rng [2]int
}

func ProjectRange(fld string, inc bool, r [2]int, projection *Projection, sort *Sort) RangeProjection {
	var p RangeProjection
	p.field = fld
	p.include = inc
	p.rng = r
	p.s = sort
	p.p = projection
	return p
}

func IncludeRange(fld string, r [2]int, projection *Projection, sort *Sort) RangeProjection {
	return ProjectRange(fld, true, r, projection, sort)
}

func ExcludeRange(fld string, r [2]int, projection *Projection, sort *Sort) RangeProjection {
	return ProjectRange(fld, false, r, projection, sort)
}

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

type MatchProjection struct {
	ArrayProjection
	q Query
}

func ProjectMatching(fld string, inc bool, match Query, projection *Projection, sort *Sort) MatchProjection {
	var p MatchProjection
	p.field = fld
	p.include = inc
	p.q = match
	p.s = sort
	p.p = projection
	return p
}

func IncludeMatching(fld string, match Query, projection *Projection, sort *Sort) MatchProjection {
	return ProjectMatching(fld, true, match, projection, sort)
}

func ExcludeMatching(fld string, match Query, projection *Projection, sort *Sort) MatchProjection {
	return ProjectMatching(fld, false, match, projection, sort)
}

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
