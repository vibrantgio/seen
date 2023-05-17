package seen

import (
	"strconv"

	"github.com/reactivego/seen/color"
)

type Surfaces []Surface

// SurfacesWith joins the points into surfaces using the coordinate map,
// which is a 2-dimensional array of index integers.
// Note that a point that is part of multiple surfaces will also be inserted
// into multiple surfaces. Because of this points shared by multiple surfaces
// have to transformed multiple times instead of only once.
// So this could be optimized by allowing surfaces to store pointers to
// points instead of the actual points.
func SurfacesWith(points Points, coordinateMap [][]int) (surfaces Surfaces) {
	surfaces = make(Surfaces, len(coordinateMap))
	for s, coords := range coordinateMap {
		for _, c := range coords {
			surfaces[s].Id = UniqueId()
			surfaces[s].Options = make(map[string]string)
			surfaces[s].Points = append(surfaces[s].Points, points[c])
		}
	}
	return
}

// SetColorFrom sets a color on every surface by reading it from
// the passed in colors.Source.
func (s Surfaces) SetColorFrom(source color.Source) (err error) {
	err = nil
	for i := range s {
		if err = s[i].SetFill(source.Read()); err != nil {
			return
		}
	}
	return
}

// SetFill applies the supplied fill Material to each surface
func (s Surfaces) SetFill(value interface{}) (err error) {
	for i := range s {
		if err = s[i].SetFill(value); err != nil {
			return
		}
	}
	return
}

// SetStroke applies the supplied stroke Material to each surface
func (s Surfaces) SetStroke(value interface{}) (err error) {
	for i := range s {
		if err = s[i].SetStroke(value); err != nil {
			return
		}
	}
	return
}

// SetStrokeWidth sets the supplied stroke-width option on each surface
func (s Surfaces) SetStrokeWidth(value int) {
	for i := range s {
		s[i].Options["stroke-width"] = strconv.Itoa(value)
	}
}

// SetShowBackfaces will set the ShowBackfaces bool on the surfaces.
func (s Surfaces) SetShowBackfaces(value bool) {
	for i := range s {
		s[i].ShowBackfaces = value
	}
}
