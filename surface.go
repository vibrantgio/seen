package seen

// Surface is a defined as a planar object in 3D space. These paths don't
// necessarily need to be convex, but they should be non-degenerate. This
// library does not support shapes with holes.
type Surface struct {
	// Points contain a list of vertices of the planar polygon that defines the
	// outline of the surface.
	Points Points

	// Id holds a unique identifier for the surface.
	// We store a unique Id for every surface so we can look them up quickly
	// with the render model cache.
	Id string

	// ShowBackfaces when set to true will override backface culling, which is useful if your
	// material is transparent. See comment in Scene.
	ShowBackfaces bool

	// FillMaterial may be a Material object which defines the color and
	// finish of the object and are rendered using the scene's shader.
	// If not material is set a Material(C.gray) will be used.
	FillMaterial *Material

	// StrokeMaterial may be a Material object that defines the color when
	// an object is stroked. By default no stroke material will be set.
	StrokeMaterial *Material

	// Dirty flag can be set whenever the RenderSurface generated from
	// the Surface needs to be regenerated.
	Dirty bool

	// Options is a map of additional options that can be specified for a surface.
	// The option with key "stroke-width" is passed in the style map parameter to
	// PathPainter.Stroke() call.
	// The keys "font" and "anchor" are passed in as keys "font" and "text-anchor" in
	// the style map parameter to TextPainter.FillText() call.
	Options map[string]string
}

// SurfacesWith joins the points into surfaces using the coordinate map,
// which is a 2-dimensional array of index integers.
// Note that a point that is part of multiple surfaces will also be inserted
// into multiple surfaces. Because of this points shared by multiple surfaces
// have to transformed multiple times instead of only once.
// So this could be optimized by allowing surfaces to store pointers to
// points instead of the actual points.
func SurfacesWith(points Points, coordinateMap [][]int) (surfaces []Surface) {
	surfaces = make([]Surface, len(coordinateMap))
	for s, coords := range coordinateMap {
		for _, c := range coords {
			surfaces[s].Id = UniqueId("s")
			surfaces[s].Options = make(map[string]string)
			surfaces[s].Points = append(surfaces[s].Points, points[c])
		}
	}
	return
}

func SurfaceWith(points Points) *Surface {
	s := &Surface{}
	s.Id = UniqueId("s")
	s.Options = make(map[string]string)
	s.Points = append(Points(nil), points...)
	return s
}

func (s *Surface) SetFill(value interface{}) (err error) {
	s.FillMaterial, err = MaterialWith(value)
	return
}

func (s *Surface) SetStroke(value interface{}) (err error) {
	s.StrokeMaterial, err = MaterialWith(value)
	return
}
