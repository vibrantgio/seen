package seen

import (
	"slices"

	"github.com/vibrantgio/seen/light"
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/transform"
)

// Group is the object collection class.
// It stores Shapes, Lights, and other Groups as well as a transformation matrix.
//
// Notably, groups are hierarchical, like a tree. This means you can isolate
// the transformation of groups of shapes in the scene, as well as create
// chains of transformations for creating, for example, articulated skeletons.
type Group struct {
	transform.Transform
	Lights   []Light
	Children []Node
}

var _ Node = (*Group)(nil)

func NewGroup(children ...Node) *Group {
	return &Group{Transform: transform.Default, Children: children}
}

func NewGroupWithLights(lights ...*light.Light) *Group {
	g := &Group{Transform: transform.Default}
	for light := range slices.Values(lights) {
		g.Lights = append(g.Lights, light)
	}
	return g
}

func (m *Group) Kind() string {
	return "group"
}

// Add lights to this `Group`
func (m *Group) AddLights(lights ...Light) {
	m.Lights = append(m.Lights, lights...)
}

// Add Nodes as children of this `Group`
// Any number of children can by supplied as arguments.
func (m *Group) Add(children ...Node) {
	m.Children = append(m.Children, children...)
}

type ObjectFunc func(object Object, lights []light.ShaderData, model matrix.Matrix)

// EachRenderable visits each Node, recursively accumulating the model matrix
// along the way. The model matrix is used to transform the points of a node
// into the global World space. For every Node the callback is invoked passing
// in the node and the model matrix as well as the list of light render datas
// that apply to the node.
func (m *Group) EachRenderable(cb ObjectFunc) {
	m.eachRenderable(cb, nil, m.Matrix())
}

// Go through the group depth first recursively and call the shape function for every shape.
func (m *Group) eachRenderable(cb ObjectFunc, lsd []light.ShaderData, model matrix.Matrix) {
	for _, light := range m.Lights {
		if light.IsEnabled() {
			lsd = append(lsd, light.ShaderData(model.Mul(light.Matrix())))
		}
	}
	for _, c := range m.Children {
		switch child := c.(type) {
		case Object:
			cb(child, lsd, model.Mul(child.Matrix()))
		case *Group:
			child.eachRenderable(cb, lsd, model.Mul(child.Matrix()))
		default:
			// skip
		}
	}
}

func (m *Group) Accept(v Visitor) {
	v.VisitGroup(m)
}
