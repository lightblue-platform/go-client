package lbclient

import (
	"fmt"
)

func ExampleSortKey() {
	s := SortKey{Field: "fld", Descending: true}
	fmt.Println(s.String())
	// Output: {"fld":"$desc"}
}
