package face

var NEXT_UNIQUE_ID int

// Returns an ID which is unique to this instance of the library
func UniqueId() int {
	NEXT_UNIQUE_ID++
	return NEXT_UNIQUE_ID
}
