package lbclient

import (
	"encoding/json"
	"fmt"
)

type SortKey struct {
	Field      string
	Descending bool
}

func (s SortKey) String() string {
	var dir string
	if s.Descending {
		dir = "$desc"
	} else {
		dir = "$asc"
	}
	return fmt.Sprintf("{\"%s\":\"%s\"}", s.Field, dir)
}

func (s SortKey) MarshalJSON() ([]byte, error) {
	return []byte(s.String()), nil
}

type Sort struct {
	Keys []SortKey
}

func (s *Sort) Empty() bool {
	return s.Keys == nil || len(s.Keys) == 0
}

// Constructs a sort using a map field:dir where dir<0 means sort descending, and dir>=0 means sort ascending
func SortBy(fields map[string]int) Sort {
	var sort Sort
	for key, value := range fields {
		sk := SortKey{key, value < 0}
		sort.Keys = append(sort.Keys, sk)
	}
	return sort
}

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
