package lbclient

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func print3(q *Query) {
	var valueName string
	var value interface{}
	if v, ok := q.q["rvalue"]; ok {
		valueName = "rvalue"
		value = v
	} else if v, ok := q.q["values"]; ok {
		valueName = "values"
		value = v
	} else {
		valueName = "rfield"
		value = q.q["rfield"]
	}
	b, _ := json.Marshal(value)
	fmt.Printf("{\"field\":\"%s\",\"op\":\"%s\",\"%s\":%s}\n",
		q.q["field"],
		q.q["op"],
		valueName,
		string(b))
}

func ExampleCmpValue() {
	print3(CmpValue("field", EQ, LitStr("string value")))
	print3(CmpValue("field", EQ, LitInt(123)))
	print3(CmpValue("field", EQ, LitBool(false)))
	print3(CmpValue("field", EQ, LitNull()))
	print3(CmpValue("field", EQ, LitDate(time.Date(2017, time.January, 15, 13, 11, 15, 123000000, time.UTC))))
	// Output:
	// {"field":"field","op":"=","rvalue":"string value"}
	// {"field":"field","op":"=","rvalue":123}
	// {"field":"field","op":"=","rvalue":false}
	// {"field":"field","op":"=","rvalue":null}
	// {"field":"field","op":"=","rvalue":"20170115T13:11:15.123+0000"}
}

func ExampleCmpValueList() {
	print3(CmpValueList("field", IN, []Literal{LitStr("s1"), LitStr("s2"), LitStr("s3")}))
	print3(CmpValueList("field", IN, LitStrs("s1", "s2", "s3")))
	print3(CmpValueList("field", IN, []Literal{LitInt(1), LitInt(2), LitInt(3)}))
	print3(CmpValueList("field", IN, LitInts(1, 2, 3)))
	print3(CmpValueList("field", IN, []Literal{LitBool(true), LitBool(false)}))
	print3(CmpValueList("field", IN, LitBools(true, false)))
	print3(CmpValueList("field", IN, []Literal{LitDate(time.Date(2017, time.January, 15, 13, 11, 15, 123000000, time.UTC)),
		LitDate(time.Date(2018, time.January, 15, 13, 11, 15, 123000000, time.UTC))}))
	print3(CmpValueList("field", IN, LitDates(time.Date(2017, time.January, 15, 13, 11, 15, 123000000, time.UTC),
		time.Date(2018, time.January, 15, 13, 11, 15, 123000000, time.UTC))))
	print3(CmpValueList("field", IN, []Literal{LitNull(), LitNull()}))
	// Output:
	// {"field":"field","op":"$in","values":["s1","s2","s3"]}
	// {"field":"field","op":"$in","values":["s1","s2","s3"]}
	// {"field":"field","op":"$in","values":[1,2,3]}
	// {"field":"field","op":"$in","values":[1,2,3]}
	// {"field":"field","op":"$in","values":[true,false]}
	// {"field":"field","op":"$in","values":[true,false]}
	// {"field":"field","op":"$in","values":["20170115T13:11:15.123+0000","20180115T13:11:15.123+0000"]}
	// {"field":"field","op":"$in","values":["20170115T13:11:15.123+0000","20180115T13:11:15.123+0000"]}
	// {"field":"field","op":"$in","values":[null,null]}
}

func ExampleCmpValues() {
	print3(CmpValues("field", IN, LitStr("s1"), LitStr("s2"), LitStr("s3")))
	print3(CmpValues("field", IN, LitInt(1), LitInt(2), LitInt(3)))
	print3(CmpValues("field", IN, LitBool(true), LitBool(false)))
	print3(CmpValues("field", IN, LitDate(time.Date(2017, time.January, 15, 13, 11, 15, 123000000, time.UTC)),
		LitDate(time.Date(2018, time.January, 15, 13, 11, 15, 123000000, time.UTC))))
	print3(CmpValues("field", IN, LitNull(), LitNull()))
	// Output:
	// {"field":"field","op":"$in","values":["s1","s2","s3"]}
	// {"field":"field","op":"$in","values":[1,2,3]}
	// {"field":"field","op":"$in","values":[true,false]}
	// {"field":"field","op":"$in","values":["20170115T13:11:15.123+0000","20180115T13:11:15.123+0000"]}
	// {"field":"field","op":"$in","values":[null,null]}
}

func ExampleCmpField() {
	print3(CmpField("field1", EQ, "field2"))
	// Output: {"field":"field1","op":"=","rfield":"field2"}
}

func ExampleCmpFieldValues() {
	print3(CmpFieldValues("field1", IN, "field2"))
	// Output: {"field":"field1","op":"$in","rfield":"field2"}
}

func ExampleCmpRegex() {
	q := CmpRegex("field", "pattern", RegexOptions{CaseInsensitive: true})
	fmt.Printf("{\"field\":\"%s\",\"regex\":\"%s\",\"caseInsensitive\":%t}", q.q["field"], q.q["regex"], q.q["caseInsensitive"])
	// Output: {"field":"field","regex":"pattern","caseInsensitive":true}
}

func TestNot(t *testing.T) {
	m := Not(CmpFieldValues("field", IN, "field2"))
	n, ok := m.q["$not"].(Query)
	if !ok {
		t.Error("Cannot get $not")
	}
	if n.q["field"] != "field" ||
		n.q["op"] != IN ||
		n.q["rfield"] != "field2" {
		t.Errorf("%q", n)
	}
}

func TestAnd(t *testing.T) {
	m := And(CmpFieldValues("field", IN, "field2"), CmpValue("field", EQ, LitStr("string value")))
	n, ok := m.q["$and"].([]Query)
	if !ok {
		t.Error("Cannot get $and")
	}
	if len(n) != 2 {
		t.Error("Expecting 2 subqueries")
	}
	if n[0].q["field"] != "field" ||
		n[0].q["op"] != IN ||
		n[0].q["rfield"] != "field2" {
		t.Errorf("%q", m)
	}
	if n[1].q["field"] != "field" ||
		n[1].q["op"] != EQ ||
		n[1].q["rvalue"].(Literal).String() != "\"string value\"" {
		t.Errorf("%q", m)
	}
}

func TestOr(t *testing.T) {
	m := Or(CmpFieldValues("field", IN, "field2"), CmpValue("field", EQ, LitStr("string value")))
	n, ok := m.q["$or"].([]Query)
	if !ok {
		t.Error("Cannot get $pr")
	}
	if len(n) != 2 {
		t.Error("Expecting 2 subqueries")
	}
	if n[0].q["field"] != "field" ||
		n[0].q["op"] != IN ||
		n[0].q["rfield"] != "field2" {
		t.Errorf("%q", m)
	}
	if n[1].q["field"] != "field" ||
		n[1].q["op"] != EQ ||
		n[1].q["rvalue"].(Literal).String() != "\"string value\"" {
		t.Errorf("%q", m)
	}
}

func TestArrayContainsStrings(t *testing.T) {
	m := ArrayContains("field", ANY, LitStr("s1"), LitStr("s2"))
	if m.q["array"] != "field" ||
		m.q["contains"] != ANY ||
		m.q["values"].([]Literal)[0].String() != "\"s1\"" ||
		m.q["values"].([]Literal)[1].String() != "\"s2\"" {
		t.Errorf("%q", m)
	}
}

func TestArrayContainsStringsList(t *testing.T) {
	m := ArrayContainsList("field", ANY, LitInts(1, 2, 3))
	if m.q["array"] != "field" ||
		m.q["contains"] != ANY ||
		m.q["values"].([]Literal)[0].String() != "1" ||
		m.q["values"].([]Literal)[1].String() != "2" ||
		m.q["values"].([]Literal)[2].String() != "3" {
		t.Errorf("%q", m)
	}
}

func TestArrayMatch(t *testing.T) {
	m := ArrayMatch("field", CmpValue("field", EQ, LitInt(123)))
	if m.q["array"] != "field" {
		t.Errorf("%q", m)
	}
	x, _ := m.q["elemMatch"].(Query)
	if x.q["field"] != "field" ||
		x.q["op"] != EQ ||
		x.q["rvalue"].(Literal).String() != "123" {
		t.Errorf("%q", m)
	}
}
