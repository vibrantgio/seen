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
			if Compare(planej, planei) != Splits {
				// Situation: plane[i] splits plane[j] but not vice versa.
				if recursion < 16 {
					// Try  to partition space with plane[j] instead of plane[i]
					return Process(plane, j, recursion+1, report)
				} else {
					// Situation: Were are probably looping....
					// TBD: use plane[i] to split plane[j]
					behind = append(behind, planej)
					if report != nil {
						report("split loop", i, j)
					}
				}
			} else {
				// Situation: planes[i] and planes[j] split each other.
				// TBD: use plane[i] to split plane[j]
				behind = append(behind, planej)
				if report != nil {
					report("split conflict", i, j)
					// plane.Face.FillMaterial, _ = shader.NewMaterialWith("#ff0000")
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
