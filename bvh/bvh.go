package bvh

import "fmt"

// Load reads and parses the BVH file at path and returns the parsed Hierarchy.
// It is a typed convenience wrapper around ParseFile, which returns an untyped
// any for compatibility with the pigeon-generated API.
func Load(path string) (Hierarchy, error) {
	res, err := ParseFile(path)
	if err != nil {
		return Hierarchy{}, err
	}
	h, ok := res.(Hierarchy)
	if !ok {
		return Hierarchy{}, fmt.Errorf("bvh: unexpected parse result of type %T", res)
	}
	return h, nil
}
