package bsp

import (
	"fmt"
	"strings"

	"github.com/reactivego/seen"
)

type BSP struct {
	Plane []seen.Plane
	Front *BSP
	Back  *BSP
}

func (bsp *BSP) Display(eye seen.Point, f func([]seen.Plane)) {
	if bsp == nil || len(bsp.Plane) == 0 {
		return
	}
	if bsp.Plane[0].Normal.Dot(eye) < bsp.Plane[0].Normal.Dot(bsp.Plane[0].Barycenter) {
		// eye in front of bsp.Plane
		if bsp.Back != nil {
			bsp.Back.Display(eye, f)
		}
		f(bsp.Plane)
		if bsp.Front != nil {
			bsp.Front.Display(eye, f)
		}
	} else {
		// eye in behind of bsp.Plane
		if bsp.Front != nil {
			bsp.Front.Display(eye, f)
		}
		f(bsp.Plane)
		if bsp.Back != nil {
			bsp.Back.Display(eye, f)
		}
	}
}

func (bsp *BSP) print(level int, sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("%*s%s:%d\n", level*2, "", bsp.Plane[0].Type, bsp.Plane[0].Id))
	if bsp.Front != nil {
		bsp.Front.print(level+1, sb)
	}
	if bsp.Back != nil {
		bsp.Back.print(level+1, sb)
	}
}

func (bsp *BSP) String() string {
	var sb strings.Builder
	bsp.print(0, &sb)
	return sb.String()
}
