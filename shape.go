package seen

import (
	"github.com/reactivego/seen/colors"
)

// Shape contains a collection of surfaces. They may create a closed 3D
// shape, but not necessarily. For example, a cube is a closed shape, but a
// patch is not.
type Shape struct {
	Object
	Kind string
	Surfaces []Surface
}

func (s *Shape) Init(kind string, surfaces []Surface) {
	s.Object.Init()
	s.Kind = kind
	s.Surfaces = surfaces
}

// ColorSurfaces sets a color on every surface of the Shape by 
// reading it from the passed in ColorReader.
func (s *Shape) ColorSurfaces(reader colors.ColorReader) (err error) {
	err = nil
	for i := range s.Surfaces {
		if err = s.Surfaces[i].SetFillMaterial(reader.ReadColor()); err != nil {
			return
		}
	}
	return
}

// SetFillMaterial applies the supplied fill Material to each surface
func (s *Shape) SetFillMaterial(value interface{}) {
	for i := range s.Surfaces {
		s.Surfaces[i].SetFillMaterial(value)
	}
}

// SetStrokeMaterial applies the supplied stroke Material to each surface
func (s *Shape) SetStrokeMaterial(value interface{}) {
	for i := range s.Surfaces {
		s.Surfaces[i].SetStrokeMaterial(value)
	}
}

// SetShowBackfaces will set the ShowBackfaces bool on the surfaces of this shape.
func (s *Shape) SetShowBackfaces(value bool) {
	for i := range s.Surfaces {
		s.Surfaces[i].ShowBackfaces = value
	}
}
