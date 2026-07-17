package nsort

import (
	"math/rand"
	"testing"

	"github.com/vibrantgio/seen/layer/bsort/bsp"
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/point"
)

// TestDebugInvariants re-runs the fuzz scene with the loop invariants
// asserted from outside on every iteration: the suffix beyond the perturbed
// prefix must be sorted descending by farD, and the emit decision must agree
// with an unpruned full scan.
func TestDebugInvariants(t *testing.T) {
	rnd := rand.New(rand.NewSource(7))
	var recs []*record
	for id := 1; id <= 80; id++ {
		c := point.Pt(rnd.Float64()*200-100, rnd.Float64()*200-100, rnd.Float64()*200-100)
		pts := make(point.Points, 3)
		for j := range pts {
			pts[j] = c.Plus(point.Pt(rnd.Float64()*120-60, rnd.Float64()*120-60, rnd.Float64()*120-60))
		}
		if pl := bsp.PlaneWith(id, pts, matrix.Identity); pl.Normal.X != pl.Normal.X {
			continue
		}
		recs = append(recs, rec(id, pts...))
	}

	violations := 0
	debugValidate = func(list []*record, i, pert int) {
		if violations > 3 {
			return
		}
		for k := i + pert; k+1 < len(list); k++ {
			if list[k].farD < list[k+1].farD {
				violations++
				t.Errorf("suffix unsorted at %d (i=%d pert=%d): farD %v < %v (ids %d, %d)",
					k, i, pert, list[k].farD, list[k+1].farD, list[k].plane.Id, list[k+1].plane.Id)
				return
			}
		}
		// Would the pruned scan and a full scan disagree about a conflict?
		P := list[i]
		pruned, full := -1, -1
		for j := i + 1; j < len(list); j++ {
			if j >= i+pert && list[j].farD <= P.nearD {
				break
			}
			if occludes(P, list[j], testEye) {
				pruned = j
				break
			}
		}
		for j := i + 1; j < len(list); j++ {
			if occludes(P, list[j], testEye) {
				full = j
				break
			}
		}
		if (pruned < 0) != (full < 0) {
			violations++
			t.Errorf("prune skipped a conflict: i=%d pert=%d pruned=%d full=%d (candidate %d vs %d)",
				i, pert, pruned, full, P.plane.Id, list[max(full, 0)].plane.Id)
		}
	}
	defer func() { debugValidate = nil }()

	var st stats
	orderRecords(recs, testEye, identityProject, func(*record) {}, &st)
	t.Logf("stats: %+v", st)
}
