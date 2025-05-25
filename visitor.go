package seen

import (
	"github.com/vibrantgio/seen/matrix"
)

type Handler interface {
	EnterGroup()
	LeaveGroup()
	VisitLight(Light, matrix.Matrix)
	VisitObject(Object, matrix.Matrix)
}

// Visitor manages scene graph coordinate transformations. It maintains the
// 'model transform' the object-space to world-space coordinate transform that
// allows local object coordinates to be transformed into world coordinates. It
// supports push/pop operations to correctly compose transformations as you
// enter and exit shapes or groups.
type Visitor interface {
	VisitGroup(*Group)
	VisitLight(Light)
	VisitObject(Object)
}

func NewVisitor(handler Handler) Visitor {
	g := new(graphVisitor)
	g.Handler = handler
	g.model.Top = matrix.Identity
	return g
}

type graphVisitor struct {
	Handler Handler
	model   struct {
		Top   matrix.Matrix
		Stack []matrix.Matrix
	}
}

var _ Visitor = (*graphVisitor)(nil)

func (v *graphVisitor) Push() {
	v.model.Stack = append(v.model.Stack, v.model.Top)
}

func (v *graphVisitor) Pop() {
	n := len(v.model.Stack)
	v.model.Top = v.model.Stack[n-1]
	v.model.Stack = v.model.Stack[:n-1]
}

func (v *graphVisitor) VisitGroup(g *Group) {
	v.Handler.EnterGroup()
	v.Push()
	v.model.Top = v.model.Top.Mul(g.Matrix())
	for _, light := range g.Lights {
		v.VisitLight(light)
	}
	for _, c := range g.Children {
		switch child := c.(type) {
		case Object:
			v.VisitObject(child)
		case *Group:
			child.Accept(v)
		default:
			// skip
		}
	}
	v.Pop()
	v.Handler.LeaveGroup()
}

func (v *graphVisitor) VisitLight(l Light) {
	v.Handler.VisitLight(l, v.model.Top.Mul(l.Matrix()))
}

func (v *graphVisitor) VisitObject(o Object) {
	v.Handler.VisitObject(o, v.model.Top.Mul(o.Matrix()))
}
