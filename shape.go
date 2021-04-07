package seen

import "github.com/reactivego/seen/colors"

// Shape contains a collection of surfaces. They may create a closed 3D
// shape, but not necessarily. For example, a cube is a closed shape, but a
// patch is not.
type Shape struct {
	Object
	Kind     string
	Surfaces []Surface
}

func ShapeWith(kind string, surfaces []Surface) Shape {
	return Shape{DefaultObject, kind, surfaces}
}

// ColorSurfaces sets a color on every surface of the Shape by
// reading it from the passed in colors.Source.
func (s *Shape) ColorSurfaces(source colors.Source) (err error) {
	err = nil
	for i := range s.Surfaces {
		if err = s.Surfaces[i].SetFill(source.Read()); err != nil {
			return
		}
	}
	return
}

// SetFill applies the supplied fill Material to each surface
func (s *Shape) SetFill(value interface{}) (err error) {
	for i := range s.Surfaces {
		if err = s.Surfaces[i].SetFill(value); err != nil {
			return
		}
	}
	return
}

// SetStroke applies the supplied stroke Material to each surface
func (s *Shape) SetStroke(value interface{}) (err error) {
	for i := range s.Surfaces {
		if err = s.Surfaces[i].SetStroke(value); err != nil {
			return
		}
	}
	return
}

// SetShowBackfaces will set the ShowBackfaces bool on the surfaces of this shape.
func (s *Shape) SetShowBackfaces(value bool) {
	for i := range s.Surfaces {
		s.Surfaces[i].ShowBackfaces = value
	}
}
