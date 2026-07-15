package bsp

func Process(plane []Plane, i int, recursion int, report func(...any)) *Tree {
	bsp := Tree{Plane: []Plane{plane[i]}}
	planei := bsp.Plane[0]
	var before, behind []Plane
	for j, planej := range plane {
		if j == i {
			continue
		}
		switch Compare(planei, planej) {
		case Coplanar:
			bsp.Plane = append(bsp.Plane, planej)
		case Before:
			before = append(before, planej)
		case Behind:
			behind = append(behind, planej)
		case Splits:
			if recursion < 16 && Compare(planej, planei) != Splits {
				// Situation: plane[i] splits plane[j] but not vice versa.
				// Try to partition space with plane[j] instead of plane[i],
				// which avoids cutting any polygon.
				return Process(plane, j, recursion+1, report)
			}
			// Either the planes split each other or re-rooting kept
			// looping: cut plane[j] by plane[i]'s plane and route the
			// pieces to their own sides, so back-to-front traversal stays
			// correct on both sides of the partition.
			if planej.NoSplit {
				// This polygon must stay whole (e.g. a text face): keep it
				// on its barycenter's side. It may sort slightly wrong
				// against the partition it straddles, but cutting it would
				// paint its content once per piece.
				if planei.Normal.Dot(planej.Barycenter) < planei.Normal.Dot(planei.Barycenter) {
					before = append(before, planej)
				} else {
					behind = append(behind, planej)
				}
				if report != nil {
					report("nosplit", i, j)
				}
			} else if front, back, ok := planei.Split(planej); ok {
				before = append(before, front)
				behind = append(behind, back)
				if report != nil {
					report("split", i, j)
				}
			} else {
				// Degenerate cut (Compare and Split disagree within
				// epsilon). Keep the polygon whole on its barycenter's side.
				if planei.Normal.Dot(planej.Barycenter) < planei.Normal.Dot(planei.Barycenter) {
					before = append(before, planej)
				} else {
					behind = append(behind, planej)
				}
				if report != nil {
					report("split failed", i, j)
				}
			}
		}
	}
	if len(before) > 0 {
		bsp.Front = Process(before, len(before)/2, 0, report)
	}
	if len(behind) > 0 {
		bsp.Back = Process(behind, len(behind)/2, 0, report)
	}
	return &bsp
}
