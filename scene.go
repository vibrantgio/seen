package seen

// Scene
type Scene struct {
	// Group is the root group for the scene, which contains Shapes, Lights, and
	// other Groups
	Group *Group

	// Camera which defines the projection transformation.
	// The default projection is perspective.
	Camera Camera

	// Viewport defines the projection from shape-space to
	// projection-space then to screen-space. The default viewport is on a
	// space from (0,0,0) to (1,1,1). To map more naturally to pixels, create a
	// viewport with the same width/height as the DOM element.
	Viewport Viewport

	// Shader determines which lighting model is used.
	Shader Shader

	// The ShowBackfaces bool can be used to turn on showing of backfaces
	// for the whole scene. Beware, turning this on can slow down a scene's
	// rendering by a factor of 2. You can also turn on backface showing for
	// individual surfaces with a boolean on those objects.
	ShowBackfaces bool

	// FractionalPoints bool determines if we round the surface
	// coordinates to the nearest integer. Rounding the coordinates before
	// display speeds up path drawing  especially when using an SVG context
	// since it cuts down on the length of path data. Anecdotally, my speedup
	// on a complex demo scene was 10 FPS. However, it does introduce a slight
	// jittering effect when animating.
	FractionalPoints bool

	// Regenerate is a bool that when set to true will force regeneration of render surfaces.
	// A render surface is generated for each surface in the scene. When Regenerate is set
	// to false (default), the generated render surfaces will be cached. The cache is a simple
	// map keyed by the surface's unique id. The cache has no eviction policy.
	// To flush the cache, call FlushCache()
	Regenerate bool
}

// EmptyScene returns a new Scene that has a default Camera, Viewport and Shader and
// an empty Group.
func EmptyScene() *Scene {
	return &Scene{
		Group:    EmptyGroup(),
		Camera:   DefaultCamera,
		Viewport: OriginViewport(0, 0, 1, 1),
		Shader:   PhongShader,
	}
}

// DefaultScene returns a new Scene that has a default Camera, Viewport and Shader and
// a Group with Hollywood-style 3-part lighting.
func DefaultScene() *Scene {
	return &Scene{
		Group:    GroupWith(DefaultLights()...),
		Camera:   DefaultCamera,
		Viewport: OriginViewport(0, 0, 1, 1),
		Shader:   PhongShader,
	}
}
