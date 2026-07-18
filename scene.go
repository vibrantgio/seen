package seen

import (
	"github.com/vibrantgio/seen/camera"
	"github.com/vibrantgio/seen/light"
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/point"
	"github.com/vibrantgio/seen/shader"
	"github.com/vibrantgio/seen/viewport"
)

// Scene
type Scene struct {
	// Group is the root group for the scene, which contains Shapes, Lights, and
	// other Groups
	Group *Group

	// Camera which defines the projection transformation.
	// The default projection is perspective.
	Camera camera.Camera

	// Viewport maps normalized device coordinates to screen space. The
	// default maps to a space from (0,0,0) to (1,1,1). To map more
	// naturally to pixels, configure the scene with FitCenter or FitOrigin
	// using the width/height of the view.
	Viewport viewport.Viewport

	// Shader determines which lighting model is used.
	Shader shader.Shader

	// The ShowBackfaces bool can be used to turn on showing of backfaces
	// for the whole scene. Beware, turning this on can slow down a scene's
	// rendering by a factor of 2. You can also turn on backface showing for
	// individual faces with a boolean on those objects.
	ShowBackfaces bool

	// Regenerate is a bool that when set to true will force regeneration of render faces.
	// A render face is generated for each face in the scene. When Regenerate is set
	// to false (default), the generated render faces will be cached. The cache is a simple
	// map keyed by the face's unique id. The cache has no eviction policy.
	// To flush the cache, call FlushCache()
	Regenerate bool
}

// NewScene returns a new Scene that has a default Camera, Viewport and Shader and
// an empty Group. So, there are no lights present in the scene.
func NewScene() *Scene {
	return &Scene{
		Group:    NewGroup(),
		Camera:   camera.Default,
		Viewport: viewport.Default,
		Shader:   shader.Default,
	}
}

// NewDefaultScene returns a new Scene that has a default Camera, Viewport and Shader and
// a Group with Hollywood-style 3-part lighting.
func NewDefaultScene() *Scene {
	return &Scene{
		Group:    NewGroupWithLights(light.DefaultLights()...),
		Camera:   camera.Default,
		Viewport: viewport.Default,
		Shader:   shader.Default,
	}
}

func (scene *Scene) Accept(handler Handler) {
	scene.Group.Accept(NewVisitor(handler))
}

// FitCenter fits the scene to a view region so that the scene's origin maps
// to the centre of the region: it places the camera eye above world (x, y)
// at the fitting distance, sets the camera's view normalization to the
// region's scale, and maps the result to the region's pixels with world
// (x, y, 0) at the centre.
//
// By default the projection's scale and camera distance follow the view's
// width and height, so the scene fills the view and rescales as the view is
// resized.
//
// Pass an optional dist to lock the projection to that fixed reference
// distance instead. The scale and camera distance then no longer depend on
// width and height — those only position the centre — so one world unit
// projects to a constant number of pixels at any view size. Use this to keep
// on-screen size from changing as the window resizes; only the first dist
// value is used. FitCenter(x, y, s, s) and FitCenter(x, y, w, h, s) coincide
// when w == h == s.
//
// Camera.Transform and Camera.Projection are left untouched.
func (s *Scene) FitCenter(x, y, w, h float64, dist ...float64) {
	projW, projH, projD := fitReference(w, h, dist)
	s.Camera.Eye = point.Pt(x, y, projD)
	s.Camera.Norm = matrix.Scale(1/projW, 1/projH, 1/projD)
	// The centring Translate always uses the real width/height, never the
	// reference.
	s.Viewport.Screen = matrix.Translate(x+w/2, y+h/2, projD).Scale(projW, -projH, projD)
}

// FitOrigin fits the scene to a view region so that the scene's origin
// aligns with the region's origin ([x, y]), which is usually the top left.
//
// As with FitCenter, the projection by default follows the view's width and
// height. Pass an optional dist to lock the scale and camera distance to
// that fixed reference instead, so on-screen size stays constant as the view
// is resized (width and height then only place the origin); only the first
// dist value is used.
//
// Camera.Transform and Camera.Projection are left untouched.
func (s *Scene) FitOrigin(x, y, w, h float64, dist ...float64) {
	projW, projH, projD := fitReference(w, h, dist)
	s.Camera.Eye = point.Pt(x, y, projD)
	s.Camera.Norm = matrix.Scale(1/projW, 1/projH, 1/projD)
	s.Viewport.Screen = matrix.Translate(x, y, projD).Scale(projW, -projH, projD)
}

// fitReference returns the reference the projection is built against: the
// view's own size by default, or the fixed dist when given.
func fitReference(w, h float64, dist []float64) (projW, projH, projD float64) {
	if len(dist) > 0 {
		return dist[0], dist[0], dist[0]
	}
	return w, h, h
}
