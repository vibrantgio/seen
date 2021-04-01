package seen

import "strconv"

var NEXT_UNIQUE_ID int

// Returns an ID which is unique to this instance of the library
func UniqueId(prefix string) string {
	NEXT_UNIQUE_ID++
	return prefix + strconv.Itoa(NEXT_UNIQUE_ID)
}
