// go:build pigeon

package main

import (
	"fmt"
	"os"
)

func main() {
	result, err := ParseFile("./testdata/01_06.bvh")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(result.(Hierarchy).Root)

	// simple := &Grammar{Buffer: example1}
	// simple.Init()
	// if err := simple.Parse(); err != nil {
	// 	println(err.Error())
	// }
}
