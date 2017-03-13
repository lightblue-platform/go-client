package lbclient

import (
	"encoding/json"
	"fmt"
)

// Regular expression query options
type RegexOptions struct {
	CaseInsensitive bool
	Extended        bool
	Multiline       bool
	Dotall          bool
}

type RelationalOp string
type NaryOp string
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

type Query struct {
	q map[string]interface{}
}

func (q *Query) Empty() bool {
	return len(q.q) == 0
}

// Returns a query of the form
//
//    { field: <field>, op:<op>, rvalue:<rvalue> }
//
// op is one of EQ, NEQ, LT, LTE, GT, GTE
func CmpValue(field string, op RelationalOp, rvalue Literal) Query {
	var ret Query
	ret.q = map[string]interface{}{"field": field, "op": op, "rvalue": rvalue}
	return ret
}

// Returns a query of the form
//
//    { field: <field>, op:<op>, values:[<values>] }
//
// op is one of IN, NIN
func CmpValueList(field string, op NaryOp, values []Literal) Query {
	var ret Query
	ret.q = map[string]interface{}{"field": field, "op": op, "values": values}
	return ret
}

// Returns a query of the form
//
//    { field: <field>, op:<op>, values:[<values>] }
//
// op is one of IN, NIN
func CmpValues(field string, op NaryOp, values ...Literal) Query {
	var ret Query
	ret.q = map[string]interface{}{"field": field, "op": op, "values": values}
	return ret
}

// Returns a query of the form
//
//    { field: <field>, op:<op>, rfield:<rfield> }
//
// op is one of EQ, NEQ, LT, LTE, GT, GTE
func CmpField(field string, op RelationalOp, rfield string) Query {
	var ret Query
	ret.q = map[string]interface{}{"field": field, "op": op, "rfield": rfield}
	return ret
}

// Returns a query of the form
//
//    { field: <field>, op:<op>, rfield:<rfield> }
//
// op is one of IN, NIN
func CmpFieldValues(field string, op NaryOp, rfield string) Query {
	var ret Query
	ret.q = map[string]interface{}{"field": field, "op": op, "rfield": rfield}
	return ret
}

// Returns a query of the form
//
//    { $not: {q} }
//
func Not(q Query) Query {
	var ret Query
	ret.q = map[string]interface{}{"$not": q}
	return ret
}

// Returns a query of the form
//
//    { $and: [ <query> ] }
//
func And(query ...Query) Query {
	var ret Query
	ret.q = map[string]interface{}{"$and": query}
	return ret
}

// Returns a query of the form
//
//    { $and: [ <query> ] }
//
func AndList(query []Query) Query {
	var ret Query
	ret.q = map[string]interface{}{"$and": query}
	return ret
}

// Returns a query of the form
//
//    { $or: [ <query> ] }
//
func Or(query ...Query) Query {
	var ret Query
	ret.q = map[string]interface{}{"$or": query}
	return ret
}

// Returns a query of the form
//
//    { $or: [ <query> ] }
//
func OrList(query []Query) Query {
	var ret Query
	ret.q = map[string]interface{}{"$or": query}
	return ret
}

// Returns a query of the form
//
// { field:<field>, regex: <pattern>, caseInsensitive: <option>, extended:<option>,multiline:<option>, dotall:<option>}
//
func CmpRegex(field string, pattern string, options RegexOptions) Query {
	var ret Query
	ret.q = map[string]interface{}{"field": field,
		"regex":           pattern,
		"caseInsensitive": options.CaseInsensitive,
		"extended":        options.Extended,
		"multiline":       options.Multiline,
		"dotall":          options.Dotall}
	return ret
}

// Returns a query of the form
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

// Returns a query of the form
//
//  { array:<array>, elemMatch: <elemMatch> }
//
func ArrayMatch(array string, elemMatch Query) Query {
	var ret Query
	ret.q = map[string]interface{}{"array": array, "elemMatch": elemMatch}
	return ret
}

func (q Query) String() string {
	return fmt.Sprintf("%s", q.q)
}

func (q Query) MarshalJSON() ([]byte, error) {
	return json.Marshal(q.q)
}
