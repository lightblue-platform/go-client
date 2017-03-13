package lbclient

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFindRequest(t *testing.T) {
	q := And(CmpValue("field1", EQ, LitInt(1)),
		CmpValue("field1", EQ, LitStr("str")))

	p := MakeProjection(IncludeTree("*"))
	req := FindRequest{RequestHeader: RequestHeader{EntityName: "test"}, Q: q, P: p, R: EMPTYRANGE}

	b, err := json.Marshal(&req)
	if err != nil {
		t.Errorf("%s", err)
	}
	s := string(b)

	expected := strings.Replace("{'entity':'test','projection':[{'field':'*','include':true,'recursive':true}],'query':{'$and':[{'field':'field1','op':'=','rvalue':1},{'field':'field1','op':'=','rvalue':'str'}]},'range':[0,-1]}", "'", "\"", -1)
	if expected != s {
		t.Errorf("Expected: %s, got: %s", expected, s)
	}

}
