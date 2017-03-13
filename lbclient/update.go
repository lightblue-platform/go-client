package lbclient

import (
	"encoding/json"
	"fmt"
)

type updatePart interface {
	fmt.Stringer
	GetMap() map[string]interface{}
}

// Represents a {$set:{field:rvalue}} operation
type SetOperation struct {
	field string
	value RValue
}

func (s SetOperation) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"$set": map[string]interface{}{s.field: s.value}}
}

// Returns {$set:{field:value}}
func (s SetOperation) String() string {
	return fmt.Sprintf("{\"$set\":{\"%s\":%s}}", s.field, s.value)
}

// Represents an {$unset:field} operation
type UnsetOperation struct {
	field string
}

// Returns $unset
func (s UnsetOperation) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"$unset": s.field}
}

// Returns {$unset:field}
func (s UnsetOperation) String() string {
	return fmt.Sprintf("{\"$unset\":\"%s\"}", s.field)
}

// Represents a {$add:{field:rvalue}} operation
type AddOperation struct {
	field string
	value RValue
}

// Returns $add
func (s AddOperation) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"$add": map[string]interface{}{s.field: s.value}}
}

// Returns {$add:{field:value}}
func (s AddOperation) String() string {
	return fmt.Sprintf("{\"$add\":{\"%s\":%s}}", s.field, s.value)
}

// Represents an array append operation {$append:{field:[values]}}
type AppendOperation struct {
	field  string
	values []RValue
}

func (s AppendOperation) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"$append": map[string]interface{}{s.field: s.values}}
}

// Returns {$append:{field:[values]}}
func (s AppendOperation) String() string {
	return fmt.Sprintf("{\"$append\":{\"%s\":%s}}", s.field, s.values)
}

// Represents an array insert operation {$insert:{field:[values]}
type InsertOperation struct {
	field  string
	values []RValue
}

func (s InsertOperation) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"$insert": map[string]interface{}{s.field: s.values}}
}

// Returns {$insert:{field:[values]}}
func (s InsertOperation) String() string {
	return fmt.Sprintf("{\"$insert\":{\"%s\":%s}}", s.field, s.values)
}

// Represents a foreach operation
type ForEachOperation struct {
	field string
	q     struct {
		q     Query
		isAll bool
	}
	u struct {
		u        Update
		isRemove bool
	}
}

func (s ForEachOperation) GetMap() map[string]interface{} {
	m := make(map[string]interface{})
	if s.q.isAll {
		m["field"] = "$all"
	} else {
		m["field"] = s.q.q
	}
	if s.u.isRemove {
		m["$update"] = "$remove"
	} else {
		m["$update"] = s.u.u
	}
	return map[string]interface{}{
		"$foreach": m}
}

// Returns {$foreach:{field:q, $update:u}}
func (s ForEachOperation) String() string {
	return fmt.Sprintf("%s", s.GetMap())
}

// An opaque Update structure containing the update expression parts
type Update struct {
	u []updatePart
}

func (u Update) Empty() bool {
	return u.u == nil || len(u.u) == 0
}

func (u *Update) add(part updatePart) *Update {
	if u.u == nil {
		u.u = make([]updatePart, 0, 4)
	}
	u.u = append(u.u, part)
	return u
}

// Adds a $set operation to the update expression
func (u *Update) Set(fld string, val RValue) *Update {
	return u.add(SetOperation{field: fld, value: val})
}

// Adds an $unset operation to the update expression
func (u *Update) Unset(fld string) *Update {
	return u.add(UnsetOperation{fld})
}

// Adds a $add operation to the update expression
func (u *Update) Add(fld string, val RValue) *Update {
	return u.add(AddOperation{field: fld, value: val})
}

// Adds an $append operation to the update epxression
func (u *Update) Append(fld string, val ...RValue) *Update {
	return u.add(AppendOperation{field: fld, values: val})
}

// Adds an $append operation to the update expression
func (u *Update) AppendList(fld string, val []RValue) *Update {
	return u.add(AppendOperation{field: fld, values: val})
}

// Adds an $insert operation to the update expression
func (u *Update) Insert(fld string, val ...RValue) *Update {
	return u.add(InsertOperation{field: fld, values: val})
}

// Adds an $insert operation to the update expression
func (u *Update) InsertList(fld string, val []RValue) *Update {
	return u.add(InsertOperation{field: fld, values: val})
}

// Adds a $foreach operation to the update expression
func (u *Update) ForEach(fld string, query Query, all bool, update Update, remove bool) *Update {
	f := ForEachOperation{field: fld}
	f.q.q = query
	f.q.isAll = all
	f.u.u = update
	f.u.isRemove = remove
	return u.add(f)
}

func (u Update) MarshalJSON() ([]byte, error) {
	switch len(u.u) {
	case 0:
		return []byte("[]"), nil
	case 1:
		return json.Marshal(u.u[0].GetMap())
	default:
		v := make([]map[string]interface{}, len(u.u))
		for i, x := range u.u {
			v[i] = x.GetMap()
		}
		return json.Marshal(v)
	}
}

func (u *Update) String() string {
	x, _ := u.MarshalJSON()
	return string(x)
}
