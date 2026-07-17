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

// NewTreeNoSplit builds the tree in whole-polygon mode: no polygon is ever
// cut. Every polygon is treated the way Process treats a NoSplit text face —
// a straddler stays whole on its barycenter's side of the partition (Process
// still re-roots onto a cut-free partition first when one exists). Paint
// order is then only approximate where polygons genuinely interpenetrate or
// occlude cyclically; in exchange no cut edges exist, so antialiased fills
// never show seams where a polygon was split. Use it for scenes that cannot
// sort wrong by construction (height fields, convex hulls) or where crawling
// split seams cost more than an occasional misordered pixel; NewTree remains
// the exact, splitting builder. The input slice is left unmodified.
func NewTreeNoSplit(planes []Plane) *Tree {
	if len(planes) == 0 {
		return nil
	}
	whole := make([]Plane, len(planes))
	copy(whole, planes)
	for i := range whole {
		whole[i].NoSplit = true
	}
	return Process(whole, len(whole)/2, 0, nil)
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
