package bsp

import (
	"fmt"
	"os"
	"strings"

	"github.com/vibrantgio/seen/point"
)

type Tree struct {
	Plane []Plane
	Front *Tree
	Back  *Tree
}

func NewTree(planes []Plane) *Tree {
	if len(planes) == 0 {
		return nil
	}

	const NoReporting = true
	if NoReporting {
		return Process(planes, len(planes)/2, 0, nil)
	}

	fmt.Fprintln(os.Stderr, "Building BSP...")
	bsp := Process(planes, len(planes)/2, 0, func(p ...any) {
		fmt.Fprintln(os.Stderr, p...)
	})
	fmt.Fprintln(os.Stderr, "BSP Done")
	return bsp
}

func (tree *Tree) Display(eye point.Point, f func([]Plane)) {
	if tree == nil || len(tree.Plane) == 0 {
		return
	}
	if tree.Plane[0].Normal.Dot(eye) < tree.Plane[0].Normal.Dot(tree.Plane[0].Barycenter) {
		// eye in front of bsp.Plane
		if tree.Back != nil {
			tree.Back.Display(eye, f)
		}
		f(tree.Plane)
		if tree.Front != nil {
			tree.Front.Display(eye, f)
		}
	} else {
		// eye in behind of bsp.Plane
		if tree.Front != nil {
			tree.Front.Display(eye, f)
		}
		f(tree.Plane)
		if tree.Back != nil {
			tree.Back.Display(eye, f)
		}
	}
}

func (tree *Tree) print(level int, sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("%*s%s:%d\n", level*2, "", "face" /*tree.Plane[0].Kind*/, tree.Plane[0].Id))
	if tree.Front != nil {
		tree.Front.print(level+1, sb)
	}
	if tree.Back != nil {
		tree.Back.print(level+1, sb)
	}
}

func (tree *Tree) String() string {
	var sb strings.Builder
	tree.print(0, &sb)
	return sb.String()
}
