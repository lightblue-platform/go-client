package lbclient

import (
	"encoding/json"
	"fmt"
)

// SortKey contains a field and sort direction
type SortKey struct {
	// The field to sort on
	Field string
	// Descending determines the sort direction
	Descending bool
}

// String returns the string representation of a sort key
func (s SortKey) String() string {
	var dir string
	if s.Descending {
		dir = "$desc"
	} else {
		dir = "$asc"
	}
	return fmt.Sprintf("{\"%s\":\"%s\"}", s.Field, dir)
}

// MarshalJSON returns the JSON representation of a sort key
func (s SortKey) MarshalJSON() ([]byte, error) {
	return []byte(s.String()), nil
}

// Sort type contains one or more sort keys
type Sort struct {
	// Sort keys
	Keys []SortKey
}

// Empty returns true if the sort is empty
func (s *Sort) Empty() bool {
	return s.Keys == nil || len(s.Keys) == 0
}

// SortBy constructs a sort using a map field:dir where dir<0 means
// sort descending, and dir>=0 means sort ascending
func SortBy(fields map[string]int) Sort {
	var sort Sort
	for key, value := range fields {
		sk := SortKey{key, value < 0}
		sort.Keys = append(sort.Keys, sk)
	}
	return sort
}

// MarshalJSON returns the JSON representation of a Sort object
func (s Sort) MarshalJSON() ([]byte, error) {
	switch len(s.Keys) {
	case 0:
		return []byte("[]"), nil
	case 1:
		return json.Marshal(s.Keys[0])
	default:
		return json.Marshal(s.Keys)
	}
}
