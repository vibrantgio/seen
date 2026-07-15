package bsp_test

import (
	"testing"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/bvh"
	"github.com/vibrantgio/seen/layer/bsort"
	"github.com/vibrantgio/seen/layer/bsort/bsp"
	"github.com/vibrantgio/seen/mocap"
)

// TestMocapInvariant builds a BSP tree from a posed motion-capture skeleton
// — organic geometry whose bones genuinely interpenetrate at every joint —
// and verifies the splitter upholds the BSP invariant there too: no polygon
// in the tree straddles any ancestor partition plane.
func TestMocapInvariant(t *testing.T) {
	h, err := bvh.Load("../../../bvh/testdata/01_06.bvh")
	if err != nil {
		t.Fatal(err)
	}
	m := mocap.New(h, nil)
	m.Apply(m.Frames() / 2)

	scene := seen.NewDefaultScene()
	scene.Group.Add(m.Group)

	var planes bsort.Planes
	scene.Accept(&planes)

	splits, failed := 0, 0
	tree := bsp.Process(planes, len(planes)/2, 0, func(p ...any) {
		switch p[0] {
		case "split":
			splits++
		case "split failed":
			failed++
		}
	})

	leaves := 0
	ids := map[int]bool{}
	checkInvariant(t, tree, nil, &leaves, ids)

	if splits == 0 {
		t.Error("no splits on interpenetrating skeleton geometry; expected many")
	}
	if failed > 0 {
		t.Errorf("%d degenerate-cut fallbacks; expected none", failed)
	}
	if leaves < len(planes) {
		t.Errorf("tree holds %d polygons, want at least the %d inputs", leaves, len(planes))
	}
	for _, p := range planes {
		if !ids[p.Id] {
			t.Errorf("face %d disappeared from the tree", p.Id)
		}
	}
	t.Logf("%d planes in, %d polygons out, %d splits", len(planes), leaves, splits)
}
