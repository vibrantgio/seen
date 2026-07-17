package nsort

import (
	"math"
	"math/rand"
	"testing"

	"github.com/vibrantgio/seen/layer/bsort/bsp"
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/point"
)

// The algorithm tests run orthographic: screen x/y ARE world x/y and the eye
// sits at (0, 0, eyeZ) with every polygon well below it, so LARGER world z is
// NEARER the eye. project is the identity — pieces keep their world x/y.

const eyeZ = 1000.0

var testEye = point.Pt(0, 0, eyeZ)

func identityProject(pts point.Points) (point.Points, bool) {
	scr := make(point.Points, len(pts))
	copy(scr, pts)
	return scr, true
}

func rec(id int, pts ...point.Point) *record {
	r := &record{plane: bsp.PlaneWith(id, pts, matrix.Identity)}
	r.scr = r.plane.Points
	r.computeExtents(testEye)
	return r
}

// quad builds the rectangle [x0,x1]×[y0,y1] with z = zAt(x, y).
func quad(id int, x0, y0, x1, y1 float64, zAt func(x, y float64) float64) *record {
	return rec(id,
		point.Pt(x0, y0, zAt(x0, y0)),
		point.Pt(x1, y0, zAt(x1, y0)),
		point.Pt(x1, y1, zAt(x1, y1)),
		point.Pt(x0, y1, zAt(x0, y1)),
	)
}

func run(t *testing.T, recs []*record) (order []*record, st stats) {
	t.Helper()
	orderRecords(recs, testEye, identityProject, func(r *record) {
		order = append(order, r)
	}, &st)
	return order, st
}

// pointInPoly: winding-free convex containment — the point must be on one
// consistent side of every edge (within slack, so edges count as inside).
func pointInPoly(poly point.Points, x, y float64) bool {
	sign := 0.0
	for i := range poly {
		j := (i + 1) % len(poly)
		cross := (poly[j].X-poly[i].X)*(y-poly[i].Y) - (poly[j].Y-poly[i].Y)*(x-poly[i].X)
		if math.Abs(cross) < 1e-9 {
			continue
		}
		if sign == 0 {
			sign = cross
		} else if (cross > 0) != (sign > 0) {
			return false
		}
	}
	return sign != 0
}

// topAt returns the id of the polygon painted last over screen point (x, y).
func topAt(order []*record, x, y float64) (int, bool) {
	id, found := 0, false
	for _, r := range order {
		if pointInPoly(r.scr, x, y) {
			id, found = r.plane.Id, true
		}
	}
	return id, found
}

// verifyNoViolations asserts the painter invariant on the emitted order:
// nothing painted earlier may occlude anything painted later. The only
// tolerated exceptions are pairs whose EARLIER member was emitted without a
// proof (forced tie-break or ops-cap fallback) — a later unverified record
// excuses nothing, since verified emits proved themselves against everything
// still unemitted at their time.
func verifyNoViolations(t *testing.T, order []*record, st stats) {
	t.Helper()
	violations := 0
	for i := 0; i < len(order); i++ {
		for j := i + 1; j < len(order); j++ {
			if occludes(order[i], order[j], testEye) && !order[i].unverified {
				violations++
				if violations <= 5 {
					t.Logf("violation: verified face %d (emitted #%d) occludes face %d (emitted #%d)",
						order[i].plane.Id, i, order[j].plane.Id, j)
				}
			}
		}
	}
	if violations > 0 {
		t.Errorf("%d occlusion violations by verified emits (stats %+v)", violations, st)
	}
}

// TestOrderDisjoint: screen-disjoint polygons never interact.
func TestOrderDisjoint(t *testing.T) {
	flat := func(z float64) func(x, y float64) float64 { return func(_, _ float64) float64 { return z } }
	order, st := run(t, []*record{
		quad(1, -100, -50, -20, 50, flat(0)),
		quad(2, 20, -50, 100, 50, flat(80)),
	})
	if len(order) != 2 || st.swaps != 0 || st.splits != 0 || st.forced != 0 {
		t.Fatalf("order %d, stats %+v — want 2 polygons untouched", len(order), st)
	}
}

// TestOrderStack: overlapping parallel quads emit far to near with no
// resolver work — the plane tests prove the sorted order.
func TestOrderStack(t *testing.T) {
	flat := func(z float64) func(x, y float64) float64 { return func(_, _ float64) float64 { return z } }
	order, st := run(t, []*record{
		quad(2, -40, -40, 40, 40, flat(100)), // nearer
		quad(1, -50, -50, 50, 50, flat(0)),   // farther
		quad(3, -30, -30, 30, 30, flat(200)), // nearest
	})
	if st.swaps != 0 || st.splits != 0 || st.forced != 0 {
		t.Fatalf("stats %+v — want a pure sort", st)
	}
	if got := []int{order[0].plane.Id, order[1].plane.Id, order[2].plane.Id}; got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Fatalf("order %v, want [1 2 3] (far to near)", got)
	}
}

// TestOrderNeighbors: edge-sharing tilted triangles — the height-field case —
// must resolve by plane tests alone.
func TestOrderNeighbors(t *testing.T) {
	a := rec(1, point.Pt(0, 0, 10), point.Pt(100, 0, -5), point.Pt(0, 100, 3))
	b := rec(2, point.Pt(100, 0, -5), point.Pt(100, 100, 12), point.Pt(0, 100, 3))
	order, st := run(t, []*record{a, b})
	if len(order) != 2 || st.swaps != 0 || st.splits != 0 || st.forced != 0 {
		t.Fatalf("order %d, stats %+v — neighbours must not conflict", len(order), st)
	}
}

// TestOrderSwapResolves: a slanted quad whose FAR corner sorts it first even
// though it lies wholly in FRONT of a small quad where they overlap. One
// move fixes the order; no cut.
func TestOrderSwapResolves(t *testing.T) {
	small := quad(1, -10, -10, 10, 10, func(_, _ float64) float64 { return 0 })
	slanted := quad(2, -50, -50, 50, 50, func(x, _ float64) float64 { return -2*x + 50 })
	order, st := run(t, []*record{slanted, small})
	if st.splits != 0 || st.forced != 0 {
		t.Fatalf("stats %+v — want swap only", st)
	}
	if st.swaps == 0 {
		t.Fatal("expected at least one swap")
	}
	if order[0].plane.Id != 1 || order[1].plane.Id != 2 {
		t.Fatalf("order [%d %d], want [1 2] (small quad is behind everywhere they overlap)",
			order[0].plane.Id, order[1].plane.Id)
	}
	verifyNoViolations(t, order, st)
}

// TestOrderCross: two quads interpenetrating in an X — the artifact-test
// scene. No linear order of whole polygons is correct; exactly one cut must
// resolve it, and each side of the intersection must show the right face.
func TestOrderCross(t *testing.T) {
	a := quad(1, -100, -100, 100, 100, func(x, _ float64) float64 { return x })  // nearer at +x
	b := quad(2, -100, -100, 100, 100, func(x, _ float64) float64 { return -x }) // nearer at -x
	order, st := run(t, []*record{a, b})
	if st.splits == 0 {
		t.Fatalf("stats %+v — a cross cannot resolve without a cut", st)
	}
	if id, ok := topAt(order, 50, 0); !ok || id != 1 {
		t.Errorf("top at (+50,0) = %d (found %v), want 1", id, ok)
	}
	if id, ok := topAt(order, -50, 0); !ok || id != 2 {
		t.Errorf("top at (-50,0) = %d (found %v), want 2", id, ok)
	}
	verifyNoViolations(t, order, st)
}

// TestOrderCycle: three bars in a rock-paper-scissors occlusion cycle. A cut
// must break the cycle and every corner must show the bar that is on top
// there.
func TestOrderCycle(t *testing.T) {
	// Corners of a triangle; bar k runs corner k -> corner k+1, high (near,
	// z=+50) at its head and low (far, z=-50) at its tail, so each bar
	// crosses OVER the next bar's tail: 0 over 1 over 2 over 0.
	corners := []point.Point{
		point.Pt(100, 0, 0),
		point.Pt(-50, 85, 0),
		point.Pt(-50, -85, 0),
	}
	bar := func(id int, from, to point.Point) *record {
		dx, dy := to.X-from.X, to.Y-from.Y
		l := math.Hypot(dx, dy)
		ux, uy := dx/l, dy/l
		px, py := -uy*9, ux*9 // half-width 9
		over := 30.0          // extend past both corners so the overlap area is real
		sx, sy := from.X-ux*over, from.Y-uy*over
		ex, eyy := to.X+ux*over, to.Y+uy*over
		zr := 100 / l // z ramps -50 (tail) .. +50 (head) across the bar
		zs, ze := -50-zr*over, 50+zr*over
		return rec(id,
			point.Pt(sx+px, sy+py, zs), point.Pt(ex+px, eyy+py, ze),
			point.Pt(ex-px, eyy-py, ze), point.Pt(sx-px, sy-py, zs),
		)
	}
	order, st := run(t, []*record{
		bar(1, corners[0], corners[1]),
		bar(2, corners[1], corners[2]),
		bar(3, corners[2], corners[0]),
	})
	if st.splits == 0 {
		t.Fatalf("stats %+v — a cycle cannot resolve without a cut", st)
	}
	// At each corner the ARRIVING bar's head lies over the LEAVING bar's tail.
	for k, want := range map[int]int{0: 3, 1: 1, 2: 2} {
		if id, ok := topAt(order, corners[k].X, corners[k].Y); !ok || id != want {
			t.Errorf("top at corner %d = %d (found %v), want bar %d", k, id, ok, want)
		}
	}
	verifyNoViolations(t, order, st)
}

// TestOrderCoplanar: coplanar overlapping quads are mutually non-occluding —
// the plane tests must pass them through as a plain sort, not a conflict.
func TestOrderCoplanar(t *testing.T) {
	flat := func(_, _ float64) float64 { return 0 }
	order, st := run(t, []*record{
		quad(1, -50, -50, 20, 20, flat),
		quad(2, -20, -20, 50, 50, flat),
	})
	if len(order) != 2 || st.swaps != 0 || st.splits != 0 || st.forced != 0 {
		t.Fatalf("order %d, stats %+v — coplanar faces must not conflict", len(order), st)
	}
}

// TestOrderNoSplitForced: the cross again, but both faces are NoSplit (text
// faces must stay whole). The resolver has to fall back to a forced emit and
// still terminate with both polygons whole.
func TestOrderNoSplitForced(t *testing.T) {
	a := quad(1, -100, -100, 100, 100, func(x, _ float64) float64 { return x })
	b := quad(2, -100, -100, 100, 100, func(x, _ float64) float64 { return -x })
	a.plane.NoSplit = true
	b.plane.NoSplit = true
	order, st := run(t, []*record{a, b})
	if len(order) != 2 {
		t.Fatalf("emitted %d polygons, want the 2 whole NoSplit faces", len(order))
	}
	if st.splits != 0 {
		t.Errorf("stats %+v — NoSplit faces were cut", st)
	}
	if st.forced == 0 {
		t.Errorf("stats %+v — expected a forced emit", st)
	}
}

// TestOrderFuzzTerminates: dense random interpenetrating triangles. The
// resolver must terminate inside its caps, keep every face represented, and
// leave no unexcused occlusion violation in the emitted order.
func TestOrderFuzzTerminates(t *testing.T) {
	rnd := rand.New(rand.NewSource(7))
	var recs []*record
	for id := 1; id <= 80; id++ {
		c := point.Pt(rnd.Float64()*200-100, rnd.Float64()*200-100, rnd.Float64()*200-100)
		pts := make(point.Points, 3)
		for j := range pts {
			pts[j] = c.Plus(point.Pt(rnd.Float64()*120-60, rnd.Float64()*120-60, rnd.Float64()*120-60))
		}
		if bsp.PlaneWith(id, pts, matrix.Identity).Normal.X != bsp.PlaneWith(id, pts, matrix.Identity).Normal.X {
			continue // NaN normal (degenerate) — the layer's scene walk skips these too
		}
		recs = append(recs, rec(id, pts...))
	}
	seen := map[int]bool{}
	var order []*record
	var st stats
	orderRecords(recs, testEye, identityProject, func(r *record) {
		order = append(order, r)
		seen[r.plane.Id] = true
	}, &st)
	t.Logf("fuzz stats: %+v, emitted %d polygons from %d faces", st, len(order), len(recs))
	if st.capped {
		t.Error("resolver hit the ops cap on 80 triangles")
	}
	for _, r := range recs {
		if !seen[r.plane.Id] {
			t.Errorf("face %d disappeared from the output", r.plane.Id)
		}
	}
	verifyNoViolations(t, order, st)
}

// TestPolyOverlap pins the precise test's edge cases.
func TestPolyOverlap(t *testing.T) {
	sq := func(x0, y0, x1, y1 float64) point.Points {
		return point.Points{point.Pt(x0, y0, 0), point.Pt(x1, y0, 0), point.Pt(x1, y1, 0), point.Pt(x0, y1, 0)}
	}
	if !polyOverlap(sq(0, 0, 10, 10), sq(5, 5, 15, 15)) {
		t.Error("overlapping squares reported disjoint")
	}
	if polyOverlap(sq(0, 0, 10, 10), sq(10, 0, 20, 10)) {
		t.Error("edge-sharing squares reported overlapping")
	}
	if polyOverlap(sq(0, 0, 10, 10), sq(10, 10, 20, 20)) {
		t.Error("corner-touching squares reported overlapping")
	}
	if polyOverlap(sq(0, 0, 10, 10), sq(11, 0, 21, 10)) {
		t.Error("separated squares reported overlapping")
	}
	// Shared-edge triangles (field neighbours).
	a := point.Points{point.Pt(0, 0, 0), point.Pt(10, 0, 0), point.Pt(0, 10, 0)}
	b := point.Points{point.Pt(10, 0, 0), point.Pt(10, 10, 0), point.Pt(0, 10, 0)}
	if polyOverlap(a, b) {
		t.Error("edge-sharing triangles reported overlapping")
	}
}

// BenchmarkOrderField measures the full ordering pass on launcher-like
// geometry: a 30×20-cell triangulated field (1200 faces), tilted back, with
// sine relief standing in for the noise displacement. This is the scene class
// nsort exists for; it must stay a plain sort — zero swaps, zero splits.
func BenchmarkOrderField(b *testing.B) {
	build := func() []*record {
		var recs []*record
		id := 0
		const cell = 70.0
		for gy := 0; gy < 20; gy++ {
			for gx := 0; gx < 30; gx++ {
				z := func(x, y float64) float64 {
					return 10 * math.Sin(x/97) * math.Cos(y/83)
				}
				x0, y0 := float64(gx)*cell-1050, float64(gy)*cell-700
				x1, y1 := x0+cell, y0+cell
				// Tilt: recede in z as y grows, like the field's RotX.
				tilt := func(x, y float64) point.Point {
					return point.Pt(x, y*0.94, z(x, y)+y*0.34-700)
				}
				id++
				recs = append(recs, rec(id, tilt(x0, y0), tilt(x1, y0), tilt(x0, y1)))
				id++
				recs = append(recs, rec(id, tilt(x1, y0), tilt(x1, y1), tilt(x0, y1)))
			}
		}
		return recs
	}

	base := build()
	work := make([]*record, len(base))
	b.ResetTimer()
	var st stats
	for range b.N {
		copy(work, base)
		st = stats{}
		orderRecords(work, testEye, identityProject, func(*record) {}, &st)
	}
	b.StopTimer()
	if st.swaps != 0 || st.splits != 0 || st.forced != 0 {
		b.Fatalf("field scene needed resolver work: %+v", st)
	}
}
