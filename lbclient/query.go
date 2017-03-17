package lbclient

import (
	"encoding/json"
)

// RegexOptions is a structure used to construct regular expression search predicates using additional options.
type RegexOptions struct {
	// CaseInsensitive=true means regex comparison is case insensitive, otherwise, comparison is case sensitive
	CaseInsensitive bool
	// Extended=true means ignore all white spaces in pattern unless escaped
	Extended bool
	// Multiline=true means treat the string as a multiline string, matching anchors to the beginning/end of lines
	Multiline bool
	// Dotall=true means dot matches everything, including newline
	Dotall bool
}

// RelationalOp is the relational operator type, one of =, !=, <, <=, >, >=
type RelationalOp string

// NaryOp is the n-ary operator type, one of $in, $nin
type NaryOp string

// ArrayOp is the array operator type, one of $any, $all, $none
type ArrayOp string

const (
	EQ  RelationalOp = "="
	NEQ RelationalOp = "!="
	LT  RelationalOp = "<"
	LTE RelationalOp = "<="
	GT  RelationalOp = ">"
	GTE RelationalOp = ">="
)

const (
	IN  NaryOp = "$in"
	NIN NaryOp = "$nin"
)

const (
	ANY  ArrayOp = "$any"
	ALL  ArrayOp = "$all"
	NONE ArrayOp = "$none"
)

// Query is an opaque structure representing a query. The underlying structure is  a map[string]interface{}
type Query struct {
	q map[string]interface{}
}

// Empty returns if a query is empty
func (q *Query) Empty() bool {
	return len(q.q) == 0
}

// CmpValue returns a query of the form
//
//    { field: <field>, op:<op>, rvalue:<rvalue> }
//
// op is one of EQ, NEQ, LT, LTE, GT, GTE
func CmpValue(field string, op RelationalOp, rvalue Literal) Query {
	var ret Query
	ret.q = map[string]interface{}{"field": field, "op": op, "rvalue": rvalue}
	return ret
}

// CmpValueList returns a query of the form
//
//    { field: <field>, op:<op>, values:[<values>] }
//
// op is one of IN, NIN
func CmpValueList(field string, op NaryOp, values []Literal) Query {
	var ret Query
	ret.q = map[string]interface{}{"field": field, "op": op, "values": values}
	return ret
}

// CmpValues returns a query of the form
//
//    { field: <field>, op:<op>, values:[<values>] }
//
// op is one of IN, NIN
func CmpValues(field string, op NaryOp, values ...Literal) Query {
	var ret Query
	ret.q = map[string]interface{}{"field": field, "op": op, "values": values}
	return ret
}

// CmpField returns a query of the form
//
//    { field: <field>, op:<op>, rfield:<rfield> }
//
// op is one of EQ, NEQ, LT, LTE, GT, GTE
func CmpField(field string, op RelationalOp, rfield string) Query {
	var ret Query
	ret.q = map[string]interface{}{"field": field, "op": op, "rfield": rfield}
	return ret
}

// CmpFieldValues returns a query of the form
//
//    { field: <field>, op:<op>, rfield:<rfield> }
//
// op is one of IN, NIN
func CmpFieldValues(field string, op NaryOp, rfield string) Query {
	var ret Query
	ret.q = map[string]interface{}{"field": field, "op": op, "rfield": rfield}
	return ret
}

// Not returns a query of the form
//
//    { $not: {q} }
//
func Not(q Query) Query {
	var ret Query
	ret.q = map[string]interface{}{"$not": q}
	return ret
}

// And returns a query of the form
//
//    { $and: [ <query> ] }
//
func And(query ...Query) Query {
	var ret Query
	ret.q = map[string]interface{}{"$and": query}
	return ret
}

// AndList returns a query of the form
//
//    { $and: [ <query> ] }
//
func AndList(query []Query) Query {
	var ret Query
	ret.q = map[string]interface{}{"$and": query}
	return ret
}

// Or returns a query of the form
//
//    { $or: [ <query> ] }
//
func Or(query ...Query) Query {
	var ret Query
	ret.q = map[string]interface{}{"$or": query}
	return ret
}

// OrList returns a query of the form
//
//    { $or: [ <query> ] }
//
func OrList(query []Query) Query {
	var ret Query
	ret.q = map[string]interface{}{"$or": query}
	return ret
}

// CmpRegex returns a query of the form
//
// { field:<field>, regex: <pattern>, caseInsensitive: <option>, extended:<option>,multiline:<option>, dotall:<option>}
//
func CmpRegex(field string, pattern string, options RegexOptions) Query {
	var ret Query
	ret.q = map[string]interface{}{"field": field,
		"regex": pattern}
	if options.CaseInsensitive {
		ret.q["caseInsensitive"] = true
	}
	if options.Extended {
		ret.q["extended"] = true
	}
	if options.Multiline {
		ret.q["multiline"] = true
	}
	if options.Dotall {
		ret.q["dotall"] = true
	}
	return ret
}

// ArryaContains returns a query of the form
//
//  { array:<array>, contains: <op>, values:[<values>] }
//
// where op is one of ANY, ALL, NONE
func ArrayContains(array string, op ArrayOp, values ...Literal) Query {
	var ret Query
	ret.q = map[string]interface{}{"array": array, "contains": op, "values": values}
	return ret
}

// Returns a query of the form
//
//  { array:<array>, contains: <op>, values:[<values>] }
//
// where op is one of ANY, ALL, NONE
func ArrayContainsList(array string, op ArrayOp, values []Literal) Query {
	var ret Query
	ret.q = map[string]interface{}{"array": array, "contains": op, "values": values}
	return ret
}

// ArrayMatch returns a query of the form
//
//  { array:<array>, elemMatch: <elemMatch> }
//
func ArrayMatch(array string, elemMatch Query) Query {
	var ret Query
	ret.q = map[string]interface{}{"array": array, "elemMatch": elemMatch}
	return ret
}

// String returns a string representation of athe query
func (q Query) String() string {
	x, _ := q.MarshalJSON()
	return string(x)
}

// MarshalJSON returns the JSON representation for the query
func (q Query) MarshalJSON() ([]byte, error) {
	return json.Marshal(q.q)
}
