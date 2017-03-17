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
	req := FindRequest{RequestHeader: RequestHeader{EntityName: "test"}, Q: q, P: p, R: &EMPTYRANGE}

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

type MakeDocDataTestStruct struct {
	StringValue string                  `json:"stringValue"`
	IntValue    int                     `json:"intValue"`
	StringArray []string                `json:"stringArray"`
	ObjArray    []MakeDocDataTestStruct `json:"objArray"`
}

func TestMakeDocData(t *testing.T) {
	td := MakeDocDataTestStruct{StringValue: "strvalue",
		IntValue:    123,
		StringArray: []string{"1", "2", "3"},
		ObjArray: []MakeDocDataTestStruct{MakeDocDataTestStruct{StringValue: "nestedStr",
			IntValue: 234}}}
	raw := MakeDocData(td)
	var r []map[string]interface{}
	json.Unmarshal([]byte(raw), &r)
	if r[0]["stringValue"].(string) != "strvalue" ||
		r[0]["intValue"].(float64) != 123 ||
		r[0]["stringArray"].([]interface{})[0].(string) != "1" ||
		r[0]["stringArray"].([]interface{})[1].(string) != "2" ||
		r[0]["stringArray"].([]interface{})[2].(string) != "3" ||
		r[0]["objArray"].([]interface{})[0].(map[string]interface{})["stringValue"].(string) != "nestedStr" ||
		r[0]["objArray"].([]interface{})[0].(map[string]interface{})["intValue"].(float64) != 234 {
		t.Errorf("%q\n", r)
	}
}

func TestMakeDocDataWithMap(t *testing.T) {
	td := map[string]string{"string1": "str1",
		"string2": "str2"}
	raw := MakeDocData(td)
	var r []map[string]string
	json.Unmarshal([]byte(raw), &r)
	if r[0]["string1"] != "str1" ||
		r[0]["string2"] != "str2" {
		t.Errorf("%q\n", r)
	}
}
