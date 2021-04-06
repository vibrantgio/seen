package seen

import "github.com/reactivego/seen/colors"

// Model is the object model class. It stores Shapes, Lights, and other Models as
// well as a transformation matrix.
//
// Notably, models are hierarchical, like a tree. This means you can isolate
// the transformation of groups of shapes in the scene, as well as create
// chains of transformations for creating, for example, articulated skeletons.
type Model struct {
	Object

	Lights   []*Light
	Children []Transformable
}

func MakeModel() *Model {
	m := &Model{}
	m.Init()
	return m
}

// MakeDefaultModel creates a default model that contains standard Hollywood-style 3-part lighting
func MakeDefaultModel() *Model {
	model := MakeModel()

	// Key light
	l := MakeDirectionalLight()
	l.Normal = Point{-1, 1, 1}.Normalize()
	l.Color = colors.ColorHsl(0.1, 0.3, 0.7, 1.0)
	l.Intensity = 1.0 // 0.004 * 255.0
	model.Add(l)

	// Back light
	l = MakeDirectionalLight()
	l.Normal = Point{1, 1, -1}.Normalize()
	l.Intensity = 0.765 // 0.003 * 255.0
	model.Add(l)

	// Fill light
	l = MakeAmbientLight()
	l.Intensity = 0.3825 // 0.0015 * 255.0
	model.Add(l)

	return model
}

// Add a `Shape`, `Light`, and other `Model` as a child of this `Model`
// Any number of children can by supplied as arguments.
// Add will return the model itself to facilitate method chaining.
func (m *Model) Add(childs ...Transformable) *Model {
	for _, child := range childs {
		switch c := child.(type) {
		case *Shape, *Model:
			m.Children = append(m.Children, c)
		case *Light:
			m.Lights = append(m.Lights, c)
		default:
			// skip
		}
	}
	return m
}

type LightFunc func(light *Light, transform Matrix) *LightRenderData
type ShapeFunc func(shape *Shape, lights []*LightRenderData, transform Matrix)

// EachRenderable visits each Light and Shape, accumulating the recursive transformation
// matrices along the way. The light callback will be called with each light
// and its accumulated transform and it should return a LightRenderData object.
// Each shape callback will be called with each shape and its accumulated
// transform as well as the list of light render datas that apply to that shape.
func (m *Model) EachRenderable(light LightFunc, shape ShapeFunc) {
	m.eachRenderable(light, shape, []*LightRenderData{}, m.Matrix())
}

// Go through the model depth first recursively and call either the light or shape function.
func (m *Model) eachRenderable(light LightFunc, shape ShapeFunc, lightModels []*LightRenderData, transform Matrix) {
	for _, l := range m.Lights {
		if l.Enabled {
			lightModels = append(lightModels, light(l, transform.Mul(l.Matrix())))
		}
	}

	for _, child := range m.Children {
		switch c := child.(type) {
		case *Shape:
			shape(c, lightModels, transform.Mul(c.Matrix()))
		case *Model:
			c.eachRenderable(light, shape, lightModels, transform.Mul(c.Matrix()))
		default:
			// skip
		}
	}
}
