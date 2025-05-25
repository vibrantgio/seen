package face

import (
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/point"
)

// Coordinates contains the transformed and projected points of a face as well
// as a reference back to the original face.
//
// Once initialized, the object will have a constant memory footprint down to
// `Number` primitives. Also, we compare each transform and projection to
// prevent unnecessary re-computation.
//
// If you need to force a re-computation, mark the face as 'dirty'.
type Coordinates struct {
	// Face is a reference to the Face that is being rendered.
	// The reference is retained so it can be checked for the Dirty flag.
	// When the Dirty flag is set, the rendering needs to be regenerated.
	Face *Face

	// Model to World space transformation
	Model      matrix.Matrix
	WorldSpace struct {
		Points     point.Points
		Barycenter point.Point
		Normal     point.Point
	}

	// World to View space transformation
	Projection    matrix.Matrix
	InViewFrustum bool

	// View to Screen space transformation
	Viewport    matrix.Matrix
	ScreenSpace struct {
		Points     point.Points
		Barycenter point.Point
		Normal     point.Point
	}
}

func (r *Coordinates) MaybeUpdate(points point.Points, model, projection, viewport matrix.Matrix) bool {
	if !model.Equal(r.Model) || !projection.Equal(r.Projection) || !viewport.Equal(r.Viewport) {
		r.Update(points, model, projection, viewport)
		return true
	}
	return false
}

func (r *Coordinates) Update(points point.Points, model, projection, viewport matrix.Matrix) {
	r.Model = model
	r.Projection = projection
	r.Viewport = viewport
	if len(r.WorldSpace.Points) != len(points) {
		r.WorldSpace.Points = make([]point.Point, len(points))
	}
	if len(r.ScreenSpace.Points) != len(points) {
		r.ScreenSpace.Points = make([]point.Point, len(points))
	}

	// Apply transform to points. Calculates transformed points and barycenter
	// Initialize the shader data with the baryCenter and the normal of the
	// transformed points.
	r.WorldSpace.Barycenter = points.MulB(r.Model, r.WorldSpace.Points)
	r.WorldSpace.Normal = r.WorldSpace.Points.Normal().Normalize()

	var clippedPoints = make(point.Points, len(r.WorldSpace.Points))
	if r.InViewFrustum = r.WorldSpace.Points.Clip(r.Projection, -2, clippedPoints); r.InViewFrustum {
		// Project camera space points into screen space
		r.ScreenSpace.Barycenter = clippedPoints.MulB(r.Viewport, r.ScreenSpace.Points)
		r.ScreenSpace.Normal = r.ScreenSpace.Points.Normal().Normalize()
	}
}
