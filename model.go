package seen

// Group is the object collection class.
// It stores Shapes, Lights, and other Groups as well as a transformation matrix.
//
// Notably, groups are hierarchical, like a tree. This means you can isolate
// the transformation of groups of shapes in the scene, as well as create
// chains of transformations for creating, for example, articulated skeletons.
type Group struct {
	Transform
	Lights   []*Light
	Children []Transformable
}

func EmptyGroup() *Group {
	return &Group{Transform: DefaultTransform}
}

func GroupWith(children ...Transformable) *Group {
	m := Group{Transform: DefaultTransform}
	m.Add(children...)
	return &m
}

// Add a `Shape`, `Light`, and other `Group` as a child of this `Group`
// Any number of children can by supplied as arguments.
func (m *Group) Add(children ...Transformable) {
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
func (m *Group) EachRenderable(shape ShapeFunc) {
	m.eachRenderable(shape, []LightShaderData{}, m.Matrix())
}

// Go through the group depth first recursively and call the shape function for every shape.
func (m *Group) eachRenderable(shape ShapeFunc, lsd []LightShaderData, transform Matrix) {
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
		case *Group:
			c.eachRenderable(shape, lsd, transform.Mul(c.Matrix()))
		default:
			// skip
		}
	}
}

type GroupVisitor interface {
	Push()
	Pop()

	VisitLight(l *Light)

	VisitSurface(s *Surface)
	EnterShape(s *Shape)
	LeaveShape(s *Shape)

	EnterGroup(m *Group)
	LeaveGroup(m *Group)
}

func (m *Group) Accept(v GroupVisitor) {
	v.Push()
	v.EnterGroup(m)
	for _, light := range m.Lights {
		v.VisitLight(light)
	}
	for _, child := range m.Children {
		switch c := child.(type) {
		case *Shape:
			v.Push()
			v.EnterShape(c)
			for i := range c.Surfaces {
				c.Surfaces[i].Shape = c
				v.VisitSurface(&c.Surfaces[i])
			}
			v.LeaveShape(c)
			v.Pop()
		case *Group:
			c.Accept(v)
		default:
			// skip
		}
	}
	v.LeaveGroup(m)
	v.Pop()
}
