package nsort

import (
	"math"
	"sort"

	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/layer/bsort/bsp"
	"github.com/vibrantgio/seen/point"
)

// This file is the view-dependent ordering core: a Newell–Newell–Sancha
// depth sort. Polygons are emitted farthest-first; a polygon may be emitted
// once it provably occludes no other unemitted polygon FROM THE CURRENT EYE.
// Proofs escalate from cheap to precise per overlapping pair, and only a
// genuine occlusion cycle — mutual overlap no linear order can satisfy —
// forces a polygon to be cut. On scenes without interpenetration this
// degenerates to a plain depth sort with zero cuts, which is exactly why it
// suits per-frame dynamic geometry where the view-independent splitting BSP
// (layer/bsort) cuts almost everything it touches.

// SideEpsilon is the world-space on-plane tolerance, shared with the BSP
// (bsp.SideEpsilon): vertices within it count as lying on the plane, so
// edge-sharing neighbours pass the plane tests instead of reading as
// conflicts.
const SideEpsilon = bsp.SideEpsilon

// overlapSlackPx is the screen-space slack for the precise 2D overlap test:
// projections that touch by less than this (shared edges, shared vertices)
// do not count as overlapping area.
const overlapSlackPx = 1e-3

// record is one polygon in the working set: its world-space plane (the bsp
// type, so bsp.Plane.Split is reusable for cycle cuts), its screen-space
// projection, and the extents the pruning tests run on.
type record struct {
	plane  bsp.Plane
	scr    point.Points
	fill   *color.Color
	stroke *color.Color
	opts   map[string]string

	minX, maxX, minY, maxY float64 // screen AABB
	nearD, farD            float64 // squared world-space distance to the eye

	movedEpoch int  // last standoff epoch this record was moved-to-front in
	unverified bool // emitted without an occlusion proof (forced or capped)
}

// debugValidate, when set (tests only), is called at the top of every
// resolver iteration with the working list, the candidate index and the
// perturbed-prefix size, to assert the loop invariants from outside.
var debugValidate func(list []*record, i, pert int)

// stats counts what the resolver had to do for one ordering pass. Tests pin
// these; scenes without interpenetration must keep swaps == splits == 0.
type stats struct {
	swaps  int // conflicts resolved by reordering alone
	splits int // polygons cut because a cycle was detected
	forced int // unresolvable (coplanar/degenerate/capped) pairs emitted anyway
	capped bool
}

// computeExtents fills the screen AABB and the squared eye-distance range.
// Eye distance is measured in world space: along any single eye ray it
// agrees with depth, which is all the extent prune relies on.
func (r *record) computeExtents(eye point.Point) {
	r.minX, r.minY = math.Inf(1), math.Inf(1)
	r.maxX, r.maxY = math.Inf(-1), math.Inf(-1)
	for _, p := range r.scr {
		r.minX, r.maxX = math.Min(r.minX, p.X), math.Max(r.maxX, p.X)
		r.minY, r.maxY = math.Min(r.minY, p.Y), math.Max(r.maxY, p.Y)
	}
	r.nearD, r.farD = math.Inf(1), math.Inf(-1)
	for _, p := range r.plane.Points {
		d := p.Minus(eye)
		dd := d.Dot(d)
		r.nearD, r.farD = math.Min(r.nearD, dd), math.Max(r.farD, dd)
	}
}

// orderRecords emits every record (and any pieces cut from them) in painter's
// order for the given eye. project maps world points to screen points and
// reports false when the polygon leaves the view frustum (such pieces are
// dropped, matching how whole faces are culled). The recs slice is reordered
// and grown in place.
//
// Invariants of the working list list[i:]:
//   - the first pert entries ([i, i+pert)) are the "perturbed prefix":
//     records moved to the front by conflict resolution, in no particular
//     depth order — the scan always tests all of them;
//   - the rest ([i+pert, len)) is the sorted suffix, descending by farD —
//     the scan may stop at the first suffix entry wholly nearer than the
//     candidate, because everything after it is nearer still.
//
// Split pieces are inserted into the sorted suffix at their depth position,
// so the suffix invariant survives cuts. Emitting advances i past a prefix
// entry first whenever pert > 0 (moved records sit at the front by
// construction).
//
// Termination: between two emits, every move marks a record that was not yet
// marked in the current standoff epoch (marked conflicts split instead of
// moving), so moves per standoff are bounded by the list length; splits are
// bounded by maxSplits; emits strictly shrink the working set. The ops guard
// is a belt-and-suspenders backstop — if it ever fires, the remaining
// records are emitted in their current (approximate) order rather than
// looping.
func orderRecords(recs []*record, eye point.Point, project func(point.Points) (point.Points, bool), emit func(*record), st *stats) {
	sort.Slice(recs, func(i, j int) bool { return recs[i].farD > recs[j].farD })

	list := recs
	maxSplits := 4 * len(list)
	maxOps := 64*len(list) + 1024
	ops := 0

	epoch := 1
	pert := 0 // size of the perturbed prefix at [i, i+pert)

	i := 0
	for i < len(list) {
		if debugValidate != nil {
			debugValidate(list, i, pert)
		}
		if ops++; ops > maxOps {
			st.capped = true
			break
		}
		P := list[i]

		conflictAt := -1
		for j := i + 1; j < len(list); j++ {
			Q := list[j]
			if j >= i+pert && Q.farD <= P.nearD {
				// Sorted suffix and Q is wholly nearer than P: P cannot
				// occlude Q nor anything after it.
				break
			}
			if occludes(P, Q, eye) {
				conflictAt = j
				break
			}
		}

		if conflictAt < 0 {
			emit(P)
			i++
			if pert > 0 {
				pert--
			}
			epoch++ // progress: previous standoff's marks expire
			continue
		}

		Q := list[conflictAt]
		if Q.movedEpoch != epoch {
			// First conflict with Q this standoff: paint Q before P. Move it
			// to the front of the working list and reconsider it there.
			Q.movedEpoch = epoch
			copy(list[i+1:conflictAt+1], list[i:conflictAt])
			list[i] = Q
			if conflictAt >= i+pert {
				pert++ // Q left the sorted suffix
			}
			st.swaps++
			continue
		}

		// Q was already moved in this standoff: P and Q occlude each other
		// in turn — a cycle no linear order satisfies. Cut one of them by
		// the other's plane; each piece then lies wholly on one side and
		// orders cleanly. Prefer cutting the candidate.
		if st.splits < maxSplits {
			if !P.plane.NoSplit {
				if front, back, ok := Q.plane.Split(P.plane); ok {
					list, pert = replaceWithPieces(list, i, i, pert, P, front, back, eye, project)
					st.splits++
					continue
				}
			}
			if !Q.plane.NoSplit {
				if front, back, ok := P.plane.Split(Q.plane); ok {
					list, pert = replaceWithPieces(list, conflictAt, i, pert, Q, front, back, eye, project)
					st.splits++
					continue
				}
			}
		}

		// No cut possible (coplanar within epsilon, NoSplit text faces, or
		// the split cap): break the tie by emitting the candidate. The scan
		// stopped at its first conflict, so the candidate may occlude more
		// than one later polygon — it is emitted without any proof at all.
		st.forced++
		P.unverified = true
		emit(P)
		i++
		if pert > 0 {
			pert--
		}
		epoch++
	}

	// Ops-guard fallback: paint what is left in its current order.
	for ; i < len(list); i++ {
		list[i].unverified = true
		emit(list[i])
	}
}

// occludes reports whether P may occlude Q from the eye — i.e. whether
// painting P before Q could be wrong. It applies the Newell tests cheapest
// first; any passing test PROVES P occludes no part of Q.
func occludes(P, Q *record, eye point.Point) bool {
	// Q wholly nearer than P: along any shared eye ray P's hit is at least
	// P.nearD out and Q's at most Q.farD, so P is never strictly in front.
	// (The scan's sorted-suffix cutoff is this same proof; stating it here
	// keeps the test complete for perturbed-prefix pairs too.)
	if Q.farD <= P.nearD {
		return false
	}
	// Screen bounding boxes disjoint: no shared pixel at all.
	if P.maxX <= Q.minX || Q.maxX <= P.minX || P.maxY <= Q.minY || Q.maxY <= P.minY {
		return false
	}
	// P wholly on the far side of Q's plane: along any shared ray Q's
	// surface comes first, so P cannot be in front of Q anywhere.
	if wholeSideOpposite(P.plane.Points, Q.plane, eye) {
		return false
	}
	// Q wholly on the eye's side of P's plane: Q is in front of P's plane
	// everywhere, so P cannot poke in front of Q.
	if wholeSideSame(Q.plane.Points, P.plane, eye) {
		return false
	}
	// Precise test: do the screen projections actually share area?
	return polyOverlap(P.scr, Q.scr)
}

// wholeSideOpposite reports whether every point lies on the opposite side of
// the plane from the eye (on-plane within SideEpsilon counts as compatible).
// False when the eye is on the plane — the test is then inconclusive.
func wholeSideOpposite(pts point.Points, plane bsp.Plane, eye point.Point) bool {
	d := plane.Normal.Dot(plane.Barycenter)
	eyeDist := plane.Normal.Dot(eye) - d
	if math.Abs(eyeDist) <= SideEpsilon {
		return false
	}
	for _, p := range pts {
		dist := plane.Normal.Dot(p) - d
		if math.Abs(dist) <= SideEpsilon {
			continue
		}
		if (dist > 0) == (eyeDist > 0) {
			return false
		}
	}
	return true
}

// wholeSideSame reports whether every point lies on the same side of the
// plane as the eye (on-plane within SideEpsilon counts as compatible).
func wholeSideSame(pts point.Points, plane bsp.Plane, eye point.Point) bool {
	d := plane.Normal.Dot(plane.Barycenter)
	eyeDist := plane.Normal.Dot(eye) - d
	if math.Abs(eyeDist) <= SideEpsilon {
		return false
	}
	for _, p := range pts {
		dist := plane.Normal.Dot(p) - d
		if math.Abs(dist) <= SideEpsilon {
			continue
		}
		if (dist > 0) != (eyeDist > 0) {
			return false
		}
	}
	return true
}

// polyOverlap is a separating-axis test on the screen-space x/y projections.
// A separating axis among the edge normals PROVES the polygons are disjoint
// (valid for any point sets); absence of one reports overlap, which is exact
// for convex polygons and errs conservative for concave ones. Touching by no
// more than overlapSlackPx (shared edges and vertices) counts as disjoint.
func polyOverlap(a, b point.Points) bool {
	if len(a) < 3 || len(b) < 3 {
		return false
	}
	project := func(pts point.Points, nx, ny float64) (min, max float64) {
		min, max = math.Inf(1), math.Inf(-1)
		for _, p := range pts {
			d := nx*p.X + ny*p.Y
			min, max = math.Min(min, d), math.Max(max, d)
		}
		return min, max
	}
	axesOf := func(pts point.Points) bool { // true when a separating axis exists
		for i := range pts {
			j := (i + 1) % len(pts)
			nx, ny := -(pts[j].Y - pts[i].Y), pts[j].X-pts[i].X
			// A degenerate edge spans no axis; skipping it is mandatory —
			// its zero-length projections would fake a separation.
			if nx*nx+ny*ny < 1e-12 {
				continue
			}
			// Normalise so the slack is in screen units regardless of edge length.
			il := 1 / math.Hypot(nx, ny)
			nx, ny = nx*il, ny*il
			minA, maxA := project(a, nx, ny)
			minB, maxB := project(b, nx, ny)
			if maxA-minB <= overlapSlackPx || maxB-minA <= overlapSlackPx {
				return true
			}
		}
		return false
	}
	return !axesOf(a) && !axesOf(b)
}

// replaceWithPieces swaps list[at] for the two pieces cut from it. Pieces
// inherit the parent's fill/stroke/options (they shade from the same face),
// are projected through the caller's pipeline, and are inserted into the
// sorted suffix at their depth position so the suffix invariant survives. A
// piece that leaves the frustum is dropped. Returns the updated list and
// prefix size.
func replaceWithPieces(list []*record, at, i, pert int, parent *record, front, back bsp.Plane, eye point.Point, project func(point.Points) (point.Points, bool)) ([]*record, int) {
	// Remove the parent. It sits in the perturbed prefix whenever it is
	// marked or is the candidate during a standoff; positionally that is
	// exactly at < i+pert.
	if at < i+pert {
		pert--
	}
	list = append(list[:at], list[at+1:]...)

	for _, piece := range []bsp.Plane{front, back} {
		scr, ok := project(piece.Points)
		if !ok {
			continue
		}
		r := &record{
			plane:  piece,
			scr:    scr,
			fill:   parent.fill,
			stroke: parent.stroke,
			opts:   parent.opts,
		}
		r.computeExtents(eye)
		// Insert into the sorted suffix [i+pert, len) by descending farD.
		lo := i + pert
		pos := lo + sort.Search(len(list)-lo, func(k int) bool { return list[lo+k].farD <= r.farD })
		list = append(list, nil)
		copy(list[pos+1:], list[pos:])
		list[pos] = r
	}
	return list, pert
}
