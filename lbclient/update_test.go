package lbclient

import (
	"strings"
	"testing"
)

func TestUpdateEmpty(t *testing.T) {
	var u Update
	cmp(t, "[]", u)
}

func TestUpdateSet(t *testing.T) {
	var u Update
	u.Set("field", LitStr("string"))
	cmp(t, strings.Replace("{'$set':{'field':'string'}}", "'", "\"", -1), u)
	var x Update
	x.Set("field", ValueOfField("f"))
	cmp(t, strings.Replace("{'$set':{'field':{'$valueof':'f'}}}", "'", "\"", -1), x)
}
