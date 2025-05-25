package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const header = `
// pointlander/peg
package main
type Grammar Peg {
	Test string
}
`

func usage() {
	fmt.Println("Usage: go run ./peg <grammar.peg>")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	frompeg := os.Args[1]
	to, found := strings.CutSuffix(frompeg, ".peg")
	if !found {
		usage()
	}
	togo := to + ".go"
	topeg := to + ".peg.peg"

	src, err := os.Open(frompeg)
	if err != nil {
		fmt.Println("Error:", err, frompeg)
	}
	defer src.Close()
	dst, err := os.Create(topeg)
	if err != nil {
		fmt.Println("Error:", err, topeg)
	}
	defer dst.Close()

	command := fmt.Sprintf("//go:generate peg -noast -switch -inline -strict -output %s %s\n", togo, topeg)
	_, err = dst.WriteString(command)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	_, err = dst.WriteString(header)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
