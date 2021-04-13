package seen

// Model is the object model class. It stores Shapes, Lights, and other Models as
// well as a transformation matrix.
//
// Notably, models are hierarchical, like a tree. This means you can isolate
// the transformation of groups of shapes in the scene, as well as create
// chains of transformations for creating, for example, articulated skeletons.
type Model struct {
	Transform
	Lights   []*Light
	Children []Transformable
}

func EmptyModel() *Model {
	return &Model{Transform: DefaultTransform}
}

func ModelWith(children ...Transformable) *Model {
	m := Model{Transform: DefaultTransform}
	m.Add(children...)
	return &m
}

// Add a `Shape`, `Light`, and other `Model` as a child of this `Model`
// Any number of children can by supplied as arguments.
func (m *Model) Add(children ...Transformable) {
	for _, child := range children {
		if light, ok := child.(*Light); ok {
			m.Lights = append(m.Lights, light)
		} else {
			m.Children = append(m.Children, child)
		}
	}
}

type ShapeFunc func(shape *Shape, lights []LightShaderData, transform Matrix)

// EachRenderable visits each Shape, accumulating the recursive transformation
// matrices along the way. Each shape callback will be called with each shape and
// its accumulated transform as well as the list of light render datas that apply
// to that shape.
func (m *Model) EachRenderable(shape ShapeFunc) {
	m.eachRenderable(shape, []LightShaderData{}, m.Matrix())
}

// Go through the model depth first recursively and call the shape function for every shape.
func (m *Model) eachRenderable(shape ShapeFunc, lsd []LightShaderData, transform Matrix) {
	for _, light := range m.Lights {
		if light.Enabled {
			t := transform.Mul(light.Matrix())
			lsd = append(lsd, light.ShaderData(t))
		}
	}
	for _, child := range m.Children {
		switch c := child.(type) {
		case *Shape:
			shape(c, lsd, transform.Mul(c.Matrix()))
		case *Model:
			c.eachRenderable(shape, lsd, transform.Mul(c.Matrix()))
		default:
			// skip
		}
	}
}

type ModelVisitor interface {
	Push()
	Pop()

	VisitLight(l *Light)

	VisitSurface(s *Surface)
	EnterShape(s *Shape)
	LeaveShape(s *Shape)

	EnterModel(m *Model)
	LeaveModel(m *Model)
}

func (m *Model) Accept(v ModelVisitor) {
	v.Push()
	v.EnterModel(m)
	for _, light := range m.Lights {
		v.VisitLight(light)
	}
	for _, child := range m.Children {
		switch c := child.(type) {
		case *Shape:
			v.Push()
			v.EnterShape(c)
			for i := range c.Surfaces {
				v.VisitSurface(&c.Surfaces[i])
			}
			v.LeaveShape(c)
			v.Pop()
		case *Model:
			c.Accept(v)
		default:
			// skip
		}
	}
	v.LeaveModel(m)
	v.Pop()
}
