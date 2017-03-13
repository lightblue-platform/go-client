package lbclient

// import (
// 	"testing"
// 	"time"
// )

// func TestWithValueString(t *testing.T) {
// 	m := WithValue("field", EQ, "string value")
// 	if m["field"] != "field" ||
// 		m["op"] != EQ ||
// 		m["rvalue"] != "string value" {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestWithValueNumber(t *testing.T) {
// 	m := WithValue("field", EQ, 123)
// 	if m["field"] != "field" ||
// 		m["op"] != EQ ||
// 		m["rvalue"] != 123 {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestWithValueBool(t *testing.T) {
// 	m := WithValue("field", EQ, false)
// 	if m["field"] != "field" ||
// 		m["op"] != EQ ||
// 		m["rvalue"] != false {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestWithValueDate(t *testing.T) {
// 	m := WithValue("field", EQ, time.Date(2017, time.January, 15, 13, 11, 15, 123000000, time.UTC))
// 	if m["field"] != "field" ||
// 		m["op"] != EQ ||
// 		m["rvalue"] != "20170115T13:11:15.123+0000" {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestWithValueNil(t *testing.T) {
// 	m := WithValue("field", EQ, nil)
// 	if m["field"] != "field" ||
// 		m["op"] != EQ ||
// 		m["rvalue"] != nil {
// 		t.Errorf("%q", m)
// 	}
// }

// func arrayCheck(t *testing.T, value interface{}, s ...interface{}) bool {
// 	v, ok := value.([]interface{})
// 	if !ok {
// 		t.Errorf("cannot get as interface{}: %q", value)
// 		return false
// 	}

// 	if len(v) != len(s) {
// 		t.Errorf("different lengths")
// 		return false
// 	}
// 	for i := 0; i < len(v); i++ {
// 		if v[i] != s[i] {
// 			t.Errorf("Not equal %q %q", v[i], s[i])
// 			return false
// 		}
// 	}
// 	return true
// }

// func TestWithValuesListString(t *testing.T) {
// 	m := WithValuesList("field", IN, []interface{}{"s1", "s2", "s3"})

// 	if m["field"] != "field" ||
// 		m["op"] != IN ||
// 		!arrayCheck(t, m["values"], "s1", "s2", "s3") {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestWithValuesListNumber(t *testing.T) {
// 	m := WithValuesList("field", IN, []interface{}{1, 2, 3})
// 	if m["field"] != "field" ||
// 		m["op"] != IN ||
// 		!arrayCheck(t, m["values"], 1, 2, 3) {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestWithValuesListBool(t *testing.T) {
// 	m := WithValuesList("field", IN, []interface{}{true, false})
// 	if m["field"] != "field" ||
// 		m["op"] != IN ||
// 		!arrayCheck(t, m["values"], true, false) {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestWithValuesListDate(t *testing.T) {
// 	m := WithValuesList("field", IN, []interface{}{time.Date(2017, time.January, 15, 13, 11, 15, 123000000, time.UTC),
// 		time.Date(2018, time.January, 15, 13, 11, 15, 123000000, time.UTC)})
// 	if m["field"] != "field" ||
// 		m["op"] != IN ||
// 		!arrayCheck(t, m["values"], "20170115T13:11:15.123+0000", "20180115T13:11:15.123+0000") {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestWithValuesListNil(t *testing.T) {
// 	m := WithValuesList("field", IN, []interface{}{"s1", nil})
// 	if m["field"] != "field" ||
// 		m["op"] != IN ||
// 		!arrayCheck(t, m["values"], "s1", nil) {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestWithValuesString(t *testing.T) {
// 	m := WithValues("field", IN, "s1", "s2", "s3")

// 	if m["field"] != "field" ||
// 		m["op"] != IN ||
// 		!arrayCheck(t, m["values"], "s1", "s2", "s3") {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestWithValuesNumber(t *testing.T) {
// 	m := WithValues("field", IN, 1, 2, 3)
// 	if m["field"] != "field" ||
// 		m["op"] != IN ||
// 		!arrayCheck(t, m["values"], 1, 2, 3) {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestWithValuesBool(t *testing.T) {
// 	m := WithValues("field", IN, true, false)
// 	if m["field"] != "field" ||
// 		m["op"] != IN ||
// 		!arrayCheck(t, m["values"], true, false) {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestWithValuesDate(t *testing.T) {
// 	m := WithValues("field", IN, time.Date(2017, 1, 15, 13, 11, 15, 123000000, time.UTC),
// 		time.Date(2018, time.January, 15, 13, 11, 15, 123000000, time.UTC))
// 	if m["field"] != "field" ||
// 		m["op"] != IN ||
// 		!arrayCheck(t, m["values"], "20170115T13:11:15.123+0000", "20180115T13:11:15.123+0000") {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestWithValuesNil(t *testing.T) {
// 	m := WithValues("field", IN, "s1", nil)
// 	if m["field"] != "field" ||
// 		m["op"] != IN ||
// 		!arrayCheck(t, m["values"], "s1", nil) {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestWithField(t *testing.T) {
// 	m := WithField("field", EQ, "field2")
// 	if m["field"] != "field" ||
// 		m["op"] != EQ ||
// 		m["rfield"] != "field2" {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestWithFieldValues(t *testing.T) {
// 	m := WithFieldValues("field", IN, "field2")
// 	if m["field"] != "field" ||
// 		m["op"] != IN ||
// 		m["rfield"] != "field2" {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestRegex(t *testing.T) {
// 	m := WithPattern("field", "pattern", RegexOptions{CaseInsensitive: true})
// 	if m["field"] != "field" ||
// 		m["regex"] != "pattern" ||
// 		m["caseInsensitive"] != true {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestNot(t *testing.T) {
// 	m := Not(WithFieldValues("field", IN, "field2"))
// 	n, ok := m["$not"].(Query)
// 	if !ok {
// 		t.Error("Cannot get $not")
// 	}
// 	m = n
// 	if m["field"] != "field" ||
// 		m["op"] != IN ||
// 		m["rfield"] != "field2" {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestAnd(t *testing.T) {
// 	m := And(WithFieldValues("field", IN, "field2"), WithValue("field", EQ, "string value"))
// 	n, ok := m["$and"].([]Query)
// 	if !ok {
// 		t.Error("Cannot get $and")
// 	}
// 	if len(n) != 2 {
// 		t.Error("Expecting 2 subqueries")
// 	}
// 	m = n[0]
// 	if m["field"] != "field" ||
// 		m["op"] != IN ||
// 		m["rfield"] != "field2" {
// 		t.Errorf("%q", m)
// 	}
// 	m = n[1]
// 	if m["field"] != "field" ||
// 		m["op"] != EQ ||
// 		m["rvalue"] != "string value" {
// 		t.Errorf("%q", m)
// 	}
// }
// func TestOr(t *testing.T) {
// 	m := Or(WithFieldValues("field", IN, "field2"), WithValue("field", EQ, "string value"))
// 	n, ok := m["$or"].([]Query)
// 	if !ok {
// 		t.Error("Cannot get $or")
// 	}
// 	if len(n) != 2 {
// 		t.Error("Expecting 2 subqueries")
// 	}
// 	m = n[0]
// 	if m["field"] != "field" ||
// 		m["op"] != IN ||
// 		m["rfield"] != "field2" {
// 		t.Errorf("%q", m)
// 	}
// 	m = n[1]
// 	if m["field"] != "field" ||
// 		m["op"] != EQ ||
// 		m["rvalue"] != "string value" {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestArrayContainsStrings(t *testing.T) {
// 	m := ArrayContains("field", ANY, "s1", "s2")
// 	if m["array"] != "field" ||
// 		m["contains"] != ANY ||
// 		!arrayCheck(t, m["values"], "s1", "s2") {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestArrayContainsStringsList(t *testing.T) {
// 	m := ArrayContainsList("field", ANY, []interface{}{"s1", "s2"})
// 	if m["array"] != "field" ||
// 		m["contains"] != ANY ||
// 		!arrayCheck(t, m["values"], "s1", "s2") {
// 		t.Errorf("%q", m)
// 	}
// }

// func TestArrayMatch(t *testing.T) {
// 	m := ArrayMatch("field", WithValue("field", EQ, 123))
// 	if m["array"] != "field" {
// 		t.Errorf("%q", m)
// 	}
// 	m, _ = m["elemMatch"].(Query)
// 	x := m
// 	if x["field"] != "field" ||
// 		x["op"] != EQ ||
// 		x["rvalue"] != 123 {
// 		t.Errorf("%q", x)
// 	}
// }
