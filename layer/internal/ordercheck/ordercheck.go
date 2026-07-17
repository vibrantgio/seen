// Package ordercheck is a deterministic repro + regression harness for
// painter-order (visibility) correctness, shared by the ordering layers.
//
// Each scene is built from explicit world-space faces and rendered through
// the same SVG pipeline used by mocap/render_test.go, with the DEFAULT
// IDENTITY CAMERA (no camera translation) so ordering is isolated from any
// eye-position concern. The emitted SVG paths are parsed back (document
// order == paint order) and, for chosen screen sample points, the polygon
// painted LAST over that point is compared against the analytically-correct
// nearest face (computed by intersecting the eye ray with each covering
// face's plane).
//
// The scenes are the two classical painter's-algorithm killers plus a
// control:
//
//   - "cross": two quads interpenetrating in an X — a polygon must paint
//     in front of the other on one side and behind it on the other;
//   - "cycle": three bars in a non-transitive occlusion cycle — no linear
//     whole-polygon order exists;
//   - "coplanar-control": disjoint squares whose paint order is irrelevant,
//     guarding the harness against false positives.
//
// A layer passes only by CUTTING polygons where the view demands it: bsort
// does so view-independently at tree-build time (bsp.NewTree), nsort
// view-dependently when it detects a cycle from the current eye. zsort,
// which never cuts, is approximate by design and is not run against this
// harness.
package ordercheck

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/camera"
	"github.com/vibrantgio/seen/context/svg"
	"github.com/vibrantgio/seen/face"
	"github.com/vibrantgio/seen/layer"
	"github.com/vibrantgio/seen/point"
	"github.com/vibrantgio/seen/shader"
	"github.com/vibrantgio/seen/shape"
	"github.com/vibrantgio/seen/viewport"
)

const (
	viewSize = 500.0
)

// faceSpec is one quad described by its world-space points and the distinct
// flat fill color (a "#RRGGBB" with channels restricted to 00/FF so Color.Hex()
// round-trips exactly) used to recover which face was painted.
type faceSpec struct {
	name  string
	hex   string // e.g. "#FF0000"
	world []point.Point
}

// ---------------------------------------------------------------------------
// Projection math, mirrored from the layers' RenderOn.
// ---------------------------------------------------------------------------

// projectPoint maps a world point to screen coordinates exactly as the
// renderers do: projection -> clip (reject when clip-space z <= -2) ->
// perspective divide -> postscale.
func projectPoint(p point.Point) (sx, sy float64, ok bool) {
	cam := camera.Default
	vp := viewport.Center(0, 0, viewSize, viewSize)
	proj := cam.Projection.Mul(vp.Prescale).Mul(cam.Matrix())
	post := vp.Postscale
	x, y, z, w := proj.Transform(p.X, p.Y, p.Z, 1.0)
	if z <= -2 {
		return 0, 0, false
	}
	rx, ry, _ := post.Transform3(x/w, y/w, z/w)
	return rx, ry, true
}

// eyePoint recovers the eye location the renderers use for back-to-front
// order: the preimage of the eye-space origin under the view transform, with
// the legacy -1/proj[2][2] estimate as the same fallback the renderers keep
// for degenerate views.
func eyePoint() point.Point {
	cam := camera.Default
	vp := viewport.Center(0, 0, viewSize, viewSize)
	view := vp.Prescale.Mul(cam.Matrix())
	if inv, ok := view.Invert(); ok {
		x, y, z := inv.Transform3(0, 0, 0)
		return point.Pt(x, y, z)
	}
	proj := cam.Projection.Mul(vp.Prescale).Mul(cam.Matrix())
	return point.Pt(0, 0, -1.0/proj[2][2])
}

// projectPoly projects all the points of a face to screen space.
func projectPoly(world []point.Point) [][2]float64 {
	out := make([][2]float64, len(world))
	for i, p := range world {
		sx, sy, ok := projectPoint(p)
		if !ok {
			panic("harness: face point clipped; choose geometry inside the frustum")
		}
		out[i] = [2]float64{sx, sy}
	}
	return out
}

// ---------------------------------------------------------------------------
// Geometry helpers.
// ---------------------------------------------------------------------------

// planeNormalD returns the plane (unit normal, offset d with n.x == d) using the
// same winding convention as point.Points.Normal: cross(p1-p0, p_last-p0).
func planeNormalD(pts []point.Point) (n point.Point, d float64) {
	p0, p1, pl := pts[0], pts[1], pts[len(pts)-1]
	n = p1.Minus(p0).Cross(pl.Minus(p0)).Normalize()
	d = n.Dot(p0)
	return
}

// rayPlaneT intersects the ray E + t*dir with the plane n.x == d.
func rayPlaneT(E, dir, n point.Point, d float64) (t float64, ok bool) {
	denom := n.Dot(dir)
	if math.Abs(denom) < 1e-12 {
		return 0, false
	}
	return (d - n.Dot(E)) / denom, true
}

// pointInPoly reports whether (px,py) is inside the screen polygon by ray
// casting. Sample points are chosen well away from edges/vertices.
func pointInPoly(px, py float64, poly [][2]float64) bool {
	inside := false
	j := len(poly) - 1
	for i := 0; i < len(poly); i++ {
		xi, yi := poly[i][0], poly[i][1]
		xj, yj := poly[j][0], poly[j][1]
		if (yi > py) != (yj > py) {
			xint := (xj-xi)*(py-yi)/(yj-yi) + xi
			if px < xint {
				inside = !inside
			}
		}
		j = i
	}
	return inside
}

// ---------------------------------------------------------------------------
// Analytic ground truth.
// ---------------------------------------------------------------------------

// expectedFace returns the name of the face that is actually nearest the eye
// along the ray through sampleWorld (which must lie on one of the faces). It also
// returns how many faces cover the sample's screen projection (a harness sanity
// figure). The ray direction is the exact eye->sampleWorld direction, so no
// screen unprojection is needed and the answer matches the real 3D geometry.
func expectedFace(sampleWorld point.Point, faces []faceSpec) (name string, covering int) {
	E := eyePoint()
	sx, sy, ok := projectPoint(sampleWorld)
	if !ok {
		return "", 0
	}
	dir := sampleWorld.Minus(E)
	bestT := math.Inf(1)
	for _, f := range faces {
		if !pointInPoly(sx, sy, projectPoly(f.world)) {
			continue
		}
		covering++
		n, d := planeNormalD(f.world)
		t, okp := rayPlaneT(E, dir, n, d)
		if okp && t > 0 && t < bestT {
			bestT = t
			name = f.name
		}
	}
	return name, covering
}

// ---------------------------------------------------------------------------
// Rendering + SVG parsing.
// ---------------------------------------------------------------------------

// renderScene builds a scene from the face specs and renders it through the
// SVG pipeline with the layer under test, returning the document. Every face
// gets a flat fill so it is always emitted; ShowBackfaces is on so backface
// culling never drops a face (isolating the ordering behaviour). The Flat
// shader ignores lights, so the emitted fill color equals the face's material
// color exactly.
func renderScene(t *testing.T, id string, faces []faceSpec, layerFor func(*seen.Scene) layer.Layer) *svg.SVG {
	t.Helper()

	scene := seen.NewScene() // no lights needed with the Flat shader
	scene.Shader = shader.Flat
	scene.ShowBackfaces = true
	scene.Camera = camera.Default // explicit: default identity-transform camera
	scene.Viewport = viewport.Center(0, 0, viewSize, viewSize)

	for _, fs := range faces {
		pts := make(point.Points, len(fs.world))
		copy(pts, fs.world)
		f := face.FaceWith(pts)
		if err := f.SetFill(fs.hex); err != nil {
			t.Fatalf("harness: SetFill(%q): %v", fs.hex, err)
		}
		obj := shape.NewShapeWithFaces(fs.name, face.Faces{f})
		scene.Group.Add(obj)
	}

	doc, err := svg.NewSVG(id, int(viewSize), int(viewSize))
	if err != nil {
		t.Fatalf("harness: NewSVG: %v", err)
	}
	ctx := svg.NewContext(doc.GetElementById(id), layerFor(scene))
	if ctx == nil {
		t.Fatalf("harness: nil svg render context")
	}
	ctx.Render()
	return doc
}

// paintedPath is one emitted SVG <path>: its fill color (lowercased "#rrggbb")
// and its screen polygon. Slice order == document order == paint order.
type paintedPath struct {
	fill string
	poly [][2]float64
}

// parsePaths walks the SVG DOM (svg -> layer <g> -> <path> children) in document
// order, extracting fill color and polygon points from each path. Walking the
// ordered DOM tree (rather than regexing the serialized string, whose attribute
// order is map-random) keeps paint order exact.
func parsePaths(t *testing.T, doc *svg.SVG) []paintedPath {
	t.Helper()
	var out []paintedPath
	var visit func(e *svg.Element)
	visit = func(e *svg.Element) {
		if e == nil {
			return
		}
		if strings.EqualFold(e.Tag, "path") {
			style, _ := e.Attribute("style")
			if strings.Contains(style, "display: none") {
				return // hidden leftover from canvas.Cleanup, not painted
			}
			d, hasD := e.Attribute("d")
			if !hasD {
				return
			}
			out = append(out, paintedPath{
				fill: fillFromStyle(style),
				poly: parsePathD(t, d),
			})
			return
		}
		for _, c := range e.ChildNodes {
			visit(c)
		}
	}
	visit(doc.Element)
	return out
}

// fillFromStyle pulls the lowercased fill color out of a "k:v;k:v;" style string.
func fillFromStyle(style string) string {
	for _, part := range strings.Split(style, ";") {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 && strings.TrimSpace(kv[0]) == "fill" {
			return strings.ToLower(strings.TrimSpace(kv[1]))
		}
	}
	return ""
}

// parsePathD parses an SVG path "d" emitted by context/svg/path.go. The exact
// format is "M{x} {y}L{x} {y}L{x} {y}..." -- the command letter is glued to the
// x of each coordinate pair and the pair's two numbers are space separated
// (precision 3); there is no explicit Z, so the polygon is implicitly closed.
// We replace the M/L commands with spaces, leaving a flat space-separated list
// of numbers, then read them in (x,y) pairs.
func parsePathD(t *testing.T, d string) [][2]float64 {
	t.Helper()
	flat := strings.NewReplacer("M", " ", "L", " ", "Z", " ", "z", " ").Replace(d)
	fields := strings.Fields(flat)
	var nums []float64
	for _, f := range fields {
		v, err := strconv.ParseFloat(f, 64)
		if err != nil {
			t.Fatalf("harness: cannot parse path coordinate %q in d=%q", f, d)
		}
		nums = append(nums, v)
	}
	if len(nums)%2 != 0 {
		t.Fatalf("harness: odd number of path coordinates in d=%q", d)
	}
	poly := make([][2]float64, 0, len(nums)/2)
	for i := 0; i+1 < len(nums); i += 2 {
		poly = append(poly, [2]float64{nums[i], nums[i+1]})
	}
	return poly
}

// gotFace returns the color of the LAST painted path (document order) whose
// polygon covers (sx,sy), and whether any path covered it at all.
func gotFace(paths []paintedPath, sx, sy float64) (fill string, covered bool) {
	for _, p := range paths {
		if pointInPoly(sx, sy, p.poly) {
			fill = p.fill
			covered = true
		}
	}
	return
}

// colorToName builds the fill->name map for a scene (lowercased hex).
func colorToName(faces []faceSpec) map[string]string {
	m := make(map[string]string, len(faces))
	for _, f := range faces {
		m[strings.ToLower(f.hex)] = f.name
	}
	return m
}

// ---------------------------------------------------------------------------
// The scenes.
// ---------------------------------------------------------------------------

// crossFaces: two quads interpenetrating in an X about the line x==0. Quad A
// slopes z=+x (front for x>0), quad B slopes z=-x (front for x<0). Each straddles
// the other's plane; A is nearer for x>0, B for x<0. No single whole-polygon
// paint order is right on both sides.
func crossFaces() []faceSpec {
	const w, h = 120.0, 120.0
	A := []point.Point{point.Pt(-w, -h, -w), point.Pt(w, -h, w), point.Pt(w, h, w), point.Pt(-w, h, -w)}
	B := []point.Point{point.Pt(-w, -h, w), point.Pt(w, -h, -w), point.Pt(w, h, -w), point.Pt(-w, h, w)}
	return []faceSpec{
		{name: "A", hex: "#FF0000", world: A},
		{name: "B", hex: "#00FF00", world: B},
	}
}

// cycleFaces: the three edges of an equilateral triangle, each a thin quad whose
// world z ramps from -60 at its start vertex to +60 at its end vertex. Eye is at
// +z, so the larger-z (ending) bar is in front at each corner. The three corner
// overlaps therefore form a non-transitive occlusion cycle: 0>1, 1>2, 2>0. No
// linear paint order can satisfy all three corners.
func cycleFaces() []faceSpec {
	V0 := point.Pt(120, 0, 0)
	V1 := point.Pt(-60, 103.9, 0)
	V2 := point.Pt(-60, -103.9, 0)
	const halfW, over, zLow, zHigh = 18.0, 35.0, -60.0, 60.0

	bar := func(a, b point.Point) []point.Point {
		dx, dy := b.X-a.X, b.Y-a.Y
		L := math.Hypot(dx, dy)
		ux, uy := dx/L, dy/L
		px, py := -uy, ux // perpendicular
		sx, sy := a.X-ux*over, a.Y-uy*over
		ex, ey := b.X+ux*over, b.Y+uy*over
		zr := (zHigh - zLow) / L
		zs := zLow - zr*over
		ze := zHigh + zr*over
		return []point.Point{
			point.Pt(sx+px*halfW, sy+py*halfW, zs),
			point.Pt(ex+px*halfW, ey+py*halfW, ze),
			point.Pt(ex-px*halfW, ey-py*halfW, ze),
			point.Pt(sx-px*halfW, sy-py*halfW, zs),
		}
	}
	return []faceSpec{
		{name: "0", hex: "#FF0000", world: bar(V0, V1)},
		{name: "1", hex: "#00FF00", world: bar(V1, V2)},
		{name: "2", hex: "#0000FF", world: bar(V2, V0)},
	}
}

// controlFaces: two non-overlapping side-by-side squares at different z. Their
// screen projections are disjoint, so paint order is irrelevant and this scene
// renders correctly with any layer. It guards the harness against false
// positives.
func controlFaces() []faceSpec {
	L := []point.Point{point.Pt(-200, -50, 0), point.Pt(-100, -50, 0), point.Pt(-100, 50, 0), point.Pt(-200, 50, 0)}
	R := []point.Point{point.Pt(100, -50, 80), point.Pt(200, -50, 80), point.Pt(200, 50, 80), point.Pt(100, 50, 80)}
	return []faceSpec{
		{name: "L", hex: "#FF0000", world: L},
		{name: "R", hex: "#00FF00", world: R},
	}
}

// sample is one screen sample point obtained by projecting a known world point
// that lies on (the front face at) that location.
type sample struct {
	label  string
	world  point.Point // a world point on the intended-visible face
	expect string      // analytic ground truth face name
}

// ---------------------------------------------------------------------------
// The check.
// ---------------------------------------------------------------------------

// Run renders the cross, cycle and control scenes through the given layer
// constructor and asserts that the face painted last over every sample point
// is the analytically nearest one. tag distinguishes the layer under test in
// SVG document ids and failure output.
func Run(t *testing.T, tag string, layerFor func(*seen.Scene) layer.Layer) {
	cases := []struct {
		scene   string
		faces   []faceSpec
		samples []sample
	}{
		{
			scene: "cross",
			faces: crossFaces(),
			samples: []sample{
				// +x side: quad A (z=+x) is nearer.
				{label: "+x side", world: point.Pt(60, 0, 60), expect: "A"},
				// -x side: quad B (z=-x) is nearer.
				{label: "-x side", world: point.Pt(-60, 0, -60), expect: "B"},
			},
		},
		{
			scene: "cycle",
			faces: cycleFaces(),
			samples: []sample{
				// Near vertex V1: bar 0 (ending here, high z) over bar 1.
				{label: "corner V1", world: point.Pt(-49.2, 85.198, 0), expect: "0"},
				// Near vertex V2: bar 1 over bar 2.
				{label: "corner V2", world: point.Pt(-49.2, -85.198, 0), expect: "1"},
				// Near vertex V0: bar 2 over bar 0.
				{label: "corner V0", world: point.Pt(98.4, 0, 0), expect: "2"},
			},
		},
		{
			scene: "coplanar-control",
			faces: controlFaces(),
			samples: []sample{
				{label: "left square", world: point.Pt(-150, 0, 0), expect: "L"},
				{label: "right square", world: point.Pt(150, 0, 80), expect: "R"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.scene, func(t *testing.T) {
			doc := renderScene(t, "artifact-"+tag+"-"+tc.scene, tc.faces, layerFor)
			paths := parsePaths(t, doc)
			c2n := colorToName(tc.faces)

			for _, s := range tc.samples {
				sx, sy, ok := projectPoint(s.world)
				if !ok {
					t.Fatalf("harness bug: %s/%s sample world %v projects outside the frustum",
						tc.scene, s.label, s.world)
				}

				// Cross-check the declared expectation against the independent
				// analytic ground truth. A mismatch is a harness bug, not an
				// artifact; surface it as Fatal so it can never masquerade as a
				// visibility failure.
				analytic, covering := expectedFace(s.world, tc.faces)
				if analytic != s.expect {
					t.Fatalf("harness bug: %s/%s analytic-nearest=%q but declared expect=%q (covering=%d)",
						tc.scene, s.label, analytic, s.expect, covering)
				}

				gotFill, covered := gotFace(paths, sx, sy)
				if !covered {
					t.Fatalf("harness bug: %s/%s screen point (%.3f,%.3f) covered by NO painted path "+
						"(expected face %q); face may have been clipped or culled",
						tc.scene, s.label, sx, sy, s.expect)
				}
				got := c2n[gotFill]
				if got == "" {
					got = "unknown(" + gotFill + ")"
				}

				if got != s.expect {
					t.Errorf("%s[%s]/%s: at screen (%.3f,%.3f) painter's order is wrong: "+
						"expected face %q nearest the eye, but face %q is painted last (covers the pixel). "+
						"covering=%d.",
						tc.scene, tag, s.label, sx, sy, s.expect, got, covering)
				}
			}
		})
	}
}
