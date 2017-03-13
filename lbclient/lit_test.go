package lbclient

import (
	"encoding/json"
	"testing"
	"time"
)

func TestLit(t *testing.T) {
	cmp(t, "1", LitInt(1))
	cmp(t, "\"str\"", LitStr("str"))
	cmp(t, "123.123", LitDouble(123.123))
	cmp(t, "true", LitBool(true))
	cmp(t, "\"20170102T13:14:15.123+0000\"", LitDate(time.Date(2017, 1, 2, 13, 14, 15, 123000000, time.UTC)))
	cmp(t, "null", LitNull())
}

func TestRValue(t *testing.T) {
	cmp(t, "{\"$valueof\":\"fld\"}", ValueOfField("fld"))
}

func TestValueOfSubst(t *testing.T) {
	b, _ := litToStr(ValueOfField("fld"))
	if string(b) != "{\"$valueof\":\"fld\"}" {
		t.Errorf("err: %q", b)
	}
}

func litToStr(l RValue) ([]byte, error) {
	return json.Marshal(l)
}

func cmp(t *testing.T, expected string, v interface{}) {
	var j []byte
	j, _ = json.Marshal(v)
	s := string(j)
	if expected != s {
		t.Errorf("expected: %s got %s", expected, s)
	}
}
