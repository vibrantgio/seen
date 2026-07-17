package bsp_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/vibrantgio/seen/layer/bsort/bsp"
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/point"
)

// planeWith builds a Plane directly from world-space points.
func planeWith(id int, pts ...point.Point) bsp.Plane {
	return bsp.PlaneWith(id, pts, matrix.Identity)
}

// area returns the area of a planar polygon.
func area(pts point.Points) float64 {
	var sum point.Point
	for i := 1; i < len(pts)-1; i++ {
		sum = sum.Plus(pts[i].Minus(pts[0]).Cross(pts[i+1].Minus(pts[0])))
	}
	return sum.Length() / 2
}

// sideOf returns the signed distance of p to l's plane.
func sideOf(l bsp.Plane, p point.Point) float64 {
	return l.Normal.Dot(p) - l.Normal.Dot(l.Barycenter)
}

// assertPiece checks that every vertex of piece lies on the wanted side of
// l's plane (or on the plane itself), and that the piece kept the identity
// and orientation of the polygon it was cut from.
func assertPiece(t *testing.T, name string, l, target, piece bsp.Plane, want float64) {
	t.Helper()
	if piece.Id != target.Id {
		t.Errorf("%s: Id = %d, want %d", name, piece.Id, target.Id)
	}
	if !piece.Piece {
		t.Errorf("%s: Piece flag not set", name)
	}
	for i, p := range piece.Points {
		if s := sideOf(l, p); s*want < -bsp.SideEpsilon {
			t.Errorf("%s: point %d (%v) has signed distance %g, want side %g", name, i, p, s, want)
		}
	}
	// The piece lies in the same plane with the same winding: its recomputed
	// unit normal matches the target's.
	if n := piece.Points.Normal().Normalize(); n.Minus(target.Normal).Length() > 1e-9 {
		t.Errorf("%s: normal = %v, want %v", name, n, target.Normal)
	}
}

func TestSplitSquare(t *testing.T) {
	// Splitter: the x=0 plane with normal +x.
	splitter := planeWith(1,
		point.Pt(0, 0, 0), point.Pt(0, 1, 0), point.Pt(0, 1, 1), point.Pt(0, 0, 1))
	if want := point.Pt(1, 0, 0); splitter.Normal.Minus(want).Length() > 1e-12 {
		t.Fatalf("splitter normal = %v, want %v", splitter.Normal, want)
	}

	// Target: a 2x2 square in the z=0 plane straddling x=0.
	target := planeWith(2,
		point.Pt(-1, -1, 0), point.Pt(1, -1, 0), point.Pt(1, 1, 0), point.Pt(-1, 1, 0))

	front, back, ok := splitter.Split(target)
	if !ok {
		t.Fatal("Split returned !ok for a straddling square")
	}
	assertPiece(t, "front", splitter, target, front, -1)
	assertPiece(t, "back", splitter, target, back, +1)

	if got, want := area(front.Points)+area(back.Points), area(target.Points); math.Abs(got-want) > 1e-9 {
		t.Errorf("areas of pieces sum to %g, want %g", got, want)
	}

	// The cut edge x=0 runs from (0,-1,0) to (0,1,0) and is shared verbatim.
	shared := 0
	for _, p := range front.Points {
		for _, q := range back.Points {
			if p.Minus(q).Length() < 1e-12 {
				shared++
			}
		}
	}
	if shared != 2 {
		t.Errorf("pieces share %d points, want the 2 cut points", shared)
	}
}

func TestSplitThroughVertices(t *testing.T) {
	// Splitter: the x=0 plane. Target: a diamond whose top and bottom
	// vertices lie exactly on the cut, so no intersection points need to be
	// inserted -- the pieces are triangles sharing those two vertices.
	splitter := planeWith(1,
		point.Pt(0, 0, 0), point.Pt(0, 1, 0), point.Pt(0, 1, 1), point.Pt(0, 0, 1))
	target := planeWith(2,
		point.Pt(0, 1, 0), point.Pt(1, 0, 0), point.Pt(0, -1, 0), point.Pt(-1, 0, 0))

	front, back, ok := splitter.Split(target)
	if !ok {
		t.Fatal("Split returned !ok for a diamond straddling the plane")
	}
	if len(front.Points) != 3 || len(back.Points) != 3 {
		t.Fatalf("piece sizes = %d and %d, want 3 and 3", len(front.Points), len(back.Points))
	}
	assertPiece(t, "front", splitter, target, front, -1)
	assertPiece(t, "back", splitter, target, back, +1)
	if got, want := area(front.Points)+area(back.Points), area(target.Points); math.Abs(got-want) > 1e-9 {
		t.Errorf("areas of pieces sum to %g, want %g", got, want)
	}
}

func TestSplitTriangle(t *testing.T) {
	splitter := planeWith(1,
		point.Pt(0, 0, 0), point.Pt(0, 1, 0), point.Pt(0, 1, 1), point.Pt(0, 0, 1))
	target := planeWith(2,
		point.Pt(-1, -1, 0), point.Pt(2, -1, 0), point.Pt(-1, 2, 0))

	front, back, ok := splitter.Split(target)
	if !ok {
		t.Fatal("Split returned !ok for a straddling triangle")
	}
	// One vertex right of the cut: the back piece is a triangle, the front
	// piece a quad.
	if len(back.Points) != 3 || len(front.Points) != 4 {
		t.Fatalf("piece sizes = front %d, back %d; want 4 and 3", len(front.Points), len(back.Points))
	}
	assertPiece(t, "front", splitter, target, front, -1)
	assertPiece(t, "back", splitter, target, back, +1)
	if got, want := area(front.Points)+area(back.Points), area(target.Points); math.Abs(got-want) > 1e-9 {
		t.Errorf("areas of pieces sum to %g, want %g", got, want)
	}
}

func TestSplitNonStraddling(t *testing.T) {
	splitter := planeWith(1,
		point.Pt(0, 0, 0), point.Pt(0, 1, 0), point.Pt(0, 1, 1), point.Pt(0, 0, 1))

	for _, tc := range []struct {
		name   string
		target bsp.Plane
	}{
		{"one side", planeWith(2,
			point.Pt(1, -1, 0), point.Pt(3, -1, 0), point.Pt(3, 1, 0), point.Pt(1, 1, 0))},
		{"touching", planeWith(3,
			point.Pt(0, -1, 0), point.Pt(2, -1, 0), point.Pt(2, 1, 0), point.Pt(0, 1, 0))},
		{"coplanar", planeWith(4,
			point.Pt(0, 2, 2), point.Pt(0, 3, 2), point.Pt(0, 3, 3), point.Pt(0, 2, 3))},
		{"within epsilon", planeWith(5,
			point.Pt(bsp.SideEpsilon/2, -1, 0), point.Pt(-bsp.SideEpsilon/2, 1, 0),
			point.Pt(bsp.SideEpsilon/2, 1, 1), point.Pt(-bsp.SideEpsilon/2, -1, 1))},
	} {
		if _, _, ok := splitter.Split(tc.target); ok {
			t.Errorf("%s: Split returned ok, want !ok", tc.name)
		}
	}
}

// TestCompareAbsoluteEpsilon pins the classification tolerance to a world
// distance: a polygon jittering 1e-6 units about a plane 300 units from the
// origin is coplanar noise, not a straddler. The previous relative epsilon
// (1e-10 of a ~300-unit dot product, a ~3e-8 dead zone) classified this as
// Splits and sent the polygon through the conflict path.
func TestCompareAbsoluteEpsilon(t *testing.T) {
	const z, jitter = 300, 1e-6
	l := planeWith(1,
		point.Pt(0, 0, z), point.Pt(1, 0, z), point.Pt(1, 1, z), point.Pt(0, 1, z))
	r := planeWith(2,
		point.Pt(2, 0, z+jitter), point.Pt(3, 0, z-jitter), point.Pt(3, 1, z-jitter), point.Pt(2, 1, z+jitter))

	if got := bsp.Compare(l, r); got != bsp.Coplanar {
		t.Errorf("Compare = %v, want Coplanar for sub-epsilon jitter at world scale", got)
	}
}

// TestProcessInvariant builds trees from interpenetrating geometry and
// verifies the defining BSP invariant on every node: planes in the Front
// subtree lie entirely on the partition plane's negative-normal side and
// planes in the Back subtree entirely on the positive side. Before polygon
// splitting was implemented, straddlers were dumped wholesale into the
// behind list and this invariant did not hold.
func TestProcessInvariant(t *testing.T) {
	quad := func(id int, c point.Point, u, v point.Point) bsp.Plane {
		return planeWith(id,
			c.Minus(u).Minus(v), c.Plus(u).Minus(v), c.Plus(u).Plus(v), c.Minus(u).Plus(v))
	}

	scenes := map[string][]bsp.Plane{
		// The artifact-test cross: two quads interpenetrating in an X.
		"cross": {
			planeWith(1, point.Pt(-120, -120, -120), point.Pt(120, -120, 120),
				point.Pt(120, 120, 120), point.Pt(-120, 120, -120)),
			planeWith(2, point.Pt(-120, -120, 120), point.Pt(120, -120, -120),
				point.Pt(120, 120, -120), point.Pt(-120, 120, 120)),
		},
		// Three mutually intersecting orthogonal squares.
		"orthogonal": {
			quad(1, point.Pt(0, 0, 0), point.Pt(100, 0, 0), point.Pt(0, 100, 0)),
			quad(2, point.Pt(0, 0, 0), point.Pt(0, 100, 0), point.Pt(0, 0, 100)),
			quad(3, point.Pt(0, 0, 0), point.Pt(0, 0, 100), point.Pt(100, 0, 0)),
		},
		"random": randomTriangles(40),
	}

	for name, planes := range scenes {
		t.Run(name, func(t *testing.T) {
			tree := bsp.NewTree(planes)
			if tree == nil {
				t.Fatal("NewTree returned nil")
			}
			leaves := 0
			ids := map[int]bool{}
			checkInvariant(t, tree, nil, &leaves, ids)
			if leaves < len(planes) {
				t.Errorf("tree holds %d planes, want at least the %d inputs", leaves, len(planes))
			}
			for _, p := range planes {
				if !ids[p.Id] {
					t.Errorf("face %d disappeared from the tree", p.Id)
				}
			}
		})
	}
}

func randomTriangles(n int) []bsp.Plane {
	rnd := rand.New(rand.NewSource(42))
	planes := make([]bsp.Plane, 0, n)
	for i := range n {
		c := point.Pt(rnd.Float64()*200-100, rnd.Float64()*200-100, rnd.Float64()*200-100)
		pts := make(point.Points, 3)
		for j := range pts {
			pts[j] = c.Plus(point.Pt(rnd.Float64()*120-60, rnd.Float64()*120-60, rnd.Float64()*120-60))
		}
		planes = append(planes, bsp.PlaneWith(i+1, pts, matrix.Identity))
	}
	return planes
}

// checkInvariant walks the tree and asserts that every plane in node's Front
// subtree is on (or within epsilon of) the negative-normal side of node's
// partition plane, and every plane in the Back subtree on the positive side.
// ancestors carries the (partition plane, required side) constraints
// accumulated on the path from the root.
type constraint struct {
	plane bsp.Plane
	want  float64
}

func checkInvariant(t *testing.T, tree *bsp.Tree, ancestors []constraint, leaves *int, ids map[int]bool) {
	t.Helper()
	if tree == nil {
		return
	}
	for _, p := range tree.Plane {
		*leaves++
		ids[p.Id] = true
		for _, c := range ancestors {
			for i, pt := range p.Points {
				if s := sideOf(c.plane, pt); s*c.want < -bsp.SideEpsilon {
					t.Errorf("face %d point %d (%v) is on the wrong side of ancestor partition %d: signed distance %g, want side %g",
						p.Id, i, pt, c.plane.Id, s, c.want)
				}
			}
		}
	}
	partition := tree.Plane[0]
	checkInvariant(t, tree.Front, append(ancestors, constraint{partition, -1}), leaves, ids)
	checkInvariant(t, tree.Back, append(ancestors, constraint{partition, +1}), leaves, ids)
}

// TestDegenerateFaceNotCoplanar: a collinear face has no plane — its normal
// is NaN. The side classification must not silently read NaN as "on the
// plane" and group unrelated faces as coplanar with it.
func TestDegenerateFaceNotCoplanar(t *testing.T) {
	partition := planeWith(1, point.Pt(-1, -1, -3), point.Pt(1, -1, -3), point.Pt(0, 1, -3))
	collinear := planeWith(2, point.Pt(0, 0, 0), point.Pt(1, 0, 0), point.Pt(2, 0, 0))
	if !math.IsNaN(collinear.Normal.X) {
		t.Fatalf("collinear face normal = %v, expected NaN", collinear.Normal)
	}
	if got := bsp.Compare(partition, collinear); got == bsp.Coplanar {
		t.Error("collinear face 3 units from the partition classified Coplanar")
	}
}

// TestProcessDegenerateNoCollapse: a degenerate face landing at the pivot
// index must not become a partition that swallows every other face into one
// unsorted node.
func TestProcessDegenerateNoCollapse(t *testing.T) {
	planes := []bsp.Plane{
		planeWith(10, point.Pt(-1, -1, -3), point.Pt(1, -1, -3), point.Pt(0, 1, -3)),
		planeWith(99, point.Pt(0, 0, 0), point.Pt(1, 0, 0), point.Pt(2, 0, 0)), // pivot: collinear, NaN normal
		planeWith(20, point.Pt(-1, -1, 3), point.Pt(1, -1, 3), point.Pt(0, 1, 3)),
	}
	tree := bsp.NewTree(planes)
	if tree.Front == nil && tree.Back == nil && len(tree.Plane) == len(planes) {
		t.Error("faces at z=-3 and z=+3 collapsed into one unsorted node behind a degenerate partition")
	}
}

// TestProcessNoSplit: a NoSplit plane that straddles the partition must stay
// whole — exactly one node holds it, uncut and unmarked — because pieces of
// it would each repaint its full content (text faces).
func TestProcessNoSplit(t *testing.T) {
	straddler := planeWith(2,
		point.Pt(-1, -1, 0), point.Pt(1, -1, 0), point.Pt(1, 1, 0), point.Pt(-1, 1, 0))
	straddler.NoSplit = true
	planes := []bsp.Plane{
		straddler,
		planeWith(1, point.Pt(0, 0, -2), point.Pt(0, 1, -2), point.Pt(0, 1, 2), point.Pt(0, 0, 2)),
	}

	found := 0
	var walk func(tr *bsp.Tree)
	walk = func(tr *bsp.Tree) {
		if tr == nil {
			return
		}
		for _, p := range tr.Plane {
			if p.Id == straddler.Id {
				found++
				if p.Piece {
					t.Error("NoSplit plane came back as a Piece")
				}
				if len(p.Points) != len(straddler.Points) {
					t.Errorf("NoSplit plane has %d points, want %d", len(p.Points), len(straddler.Points))
				}
			}
		}
		walk(tr.Front)
		walk(tr.Back)
	}
	walk(bsp.NewTree(planes))
	if found != 1 {
		t.Errorf("NoSplit plane appears in %d nodes, want exactly 1", found)
	}
}
