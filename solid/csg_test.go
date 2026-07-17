package solid_test

import (
	"math"
	"testing"

	"github.com/vibrantgio/seen/solid"
)

// volume computes the signed volume enclosed by a CSG's polygons via the
// divergence theorem: each polygon is fanned into triangles and each
// triangle (a, b, c) contributes a·(b×c)/6. With the package's outward-facing
// winding the result is POSITIVE and equals the enclosed volume — one number
// that catches winding, flip, clip and split mistakes at once. It assumes
// nothing about connectivity, so the T-junctions CSG output legitimately
// contains do not disturb it.
func volume(s solid.CSG) float64 {
	v := 0.0
	for _, p := range s {
		for i := 1; i+1 < len(p.Vertices); i++ {
			a, b, c := p.Vertices[0].Pos, p.Vertices[i].Pos, p.Vertices[i+1].Pos
			v += a.Dot(b.Cross(c))
		}
	}
	return v / 6
}

// wantVolume asserts an exact expected volume within float tolerance —
// boolean results of plane-bounded solids are exact up to rounding.
func wantVolume(t *testing.T, name string, s solid.CSG, want float64) {
	t.Helper()
	if got := volume(s); math.Abs(got-want) > 1e-9 {
		t.Errorf("%s: volume = %v, want %v", name, got, want)
	}
}

// TestVolumePrimitives pins the generators' winding and closure. The
// tessellated sphere and cylinder are inscribed, so their volumes fall a few
// percent SHORT of the analytic solid — never over.
func TestVolumePrimitives(t *testing.T) {
	wantVolume(t, "cube 2x2x2", solid.Cube(), 8)
	wantVolume(t, "cuboid 1x2x3", solid.Cube(solid.Size(1, 2, 3)), 6)

	sphere := volume(solid.Sphere(solid.Slices(32), solid.Stacks(16)))
	exact := 4.0 / 3.0 * math.Pi
	if sphere < 0.95*exact || sphere > exact {
		t.Errorf("sphere volume = %v, want within [%v, %v] (inscribed tessellation of %v)",
			sphere, 0.95*exact, exact, exact)
	}

	cyl := volume(solid.Cylinder(solid.Slices(32)))
	exactCyl := math.Pi * 1 * 1 * 2 // r=1, height 2
	if cyl < 0.95*exactCyl || cyl > exactCyl {
		t.Errorf("cylinder volume = %v, want within [%v, %v] (inscribed tessellation of %v)",
			cyl, 0.95*exactCyl, exactCyl, exactCyl)
	}
}

// TestVolumeBooleans: two unit-ish cubes offset by half their width have an
// analytically known union/difference/intersection. All cuts run along
// axis-aligned planes, so the volumes are exact.
func TestVolumeBooleans(t *testing.T) {
	a := solid.Cube()                      // [-1,1]^3, volume 8
	b := solid.Cube(solid.Center(1, 0, 0)) // [0,2]x[-1,1]^2, overlap 4

	wantVolume(t, "A∪B", a.Union(b), 12)
	wantVolume(t, "A−B", a.Subtract(b), 4)
	wantVolume(t, "A∩B", a.Intersect(b), 4)
	wantVolume(t, "B−A", b.Subtract(a), 4)

	// Disjoint solids: union adds, intersection annihilates.
	c := solid.Cube(solid.Center(5, 0, 0))
	wantVolume(t, "A∪C disjoint", a.Union(c), 16)
	wantVolume(t, "A∩C disjoint", a.Intersect(c), 0)
	wantVolume(t, "A−C disjoint", a.Subtract(c), 8)
}

// TestVolumeSelfOperations: identical operands exercise the coplanar paths.
func TestVolumeSelfOperations(t *testing.T) {
	a := solid.Cube()
	wantVolume(t, "A∪A", a.Union(a), 8)
	wantVolume(t, "A∩A", a.Intersect(a), 8)
	wantVolume(t, "A−A", a.Subtract(a), 0)
}

// TestEmptyOperands: booleans with empty operands resolve algebraically and
// never panic — Subtract and Intersect used to dereference the empty tree's
// nil root plane in Invert. Chained emptiness (the result of A−A feeding the
// next operation) is the realistic route into that crash.
func TestEmptyOperands(t *testing.T) {
	a := solid.Cube()
	empty := solid.CSG{}

	wantVolume(t, "∅∪A", empty.Union(a), 8)
	wantVolume(t, "A∪∅", a.Union(empty), 8)
	wantVolume(t, "∅−A", empty.Subtract(a), 0)
	wantVolume(t, "A−∅", a.Subtract(empty), 8)
	wantVolume(t, "∅∩A", empty.Intersect(a), 0)
	wantVolume(t, "A∩∅", a.Intersect(empty), 0)

	wantVolume(t, "(A−A)−A", a.Subtract(a).Subtract(a), 0)
	wantVolume(t, "(A−A)∪A", a.Subtract(a).Union(a), 8)

	// The BSP-level guard directly: inverting an empty tree is a no-op.
	(&solid.BSP{}).Invert()
}

// TestDegeneratePivot: a zero-area polygon (collinear points, NaN plane
// normal) must never become a partition plane — a NaN pivot classifies every
// polygon COPLANAR, collapsing the whole set into one unsorted node and
// making later clips keep everything. The poisoned solid must produce the
// same boolean results as the clean one.
func TestDegeneratePivot(t *testing.T) {
	v := func(x, y, z float64) solid.Vertex {
		return solid.Vertex{Pos: solid.Vector{X: x, Y: y, Z: z}}
	}
	degenerate := solid.PolygonFromVertices(v(0, 0, 0), v(1, 0, 0), v(2, 0, 0))
	if !math.IsNaN(degenerate.Plane.Normal.X) {
		t.Fatal("harness expectation: collinear polygon should have a NaN plane normal")
	}

	// The degenerate polygon FIRST, so it is exactly the pivot AddPolygons
	// would have picked.
	poisoned := append(solid.CSG{degenerate}, solid.Cube()...)
	b := solid.Cube(solid.Center(1, 0, 0))

	wantVolume(t, "poisoned−B", poisoned.Subtract(b), 4)
	wantVolume(t, "poisoned∪B", poisoned.Union(b), 12)
	wantVolume(t, "poisoned∩B", poisoned.Intersect(b), 4)

	// All-degenerate input is zero-area debris: it must drop out entirely
	// rather than poison the tree it is added to.
	debris := solid.CSG{degenerate}
	wantVolume(t, "debris−B", debris.Subtract(b), 0)
	wantVolume(t, "debris∪B", debris.Union(b), 8)
}
