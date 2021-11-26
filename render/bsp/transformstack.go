package bsp

import "github.com/reactivego/seen"

type TransformStack struct {
	Transform seen.Matrix
	Stack     []seen.Matrix
}

func (v *TransformStack) Push() {
	if len(v.Stack) == 0 {
		v.Transform = seen.IdentityMatrix
		v.Stack = []seen.Matrix{v.Transform}
	} else {
		v.Stack = append(v.Stack, v.Transform)
	}
}

func (v *TransformStack) Pop() {
	v.Transform = v.Stack[len(v.Stack)-1]
	v.Stack = v.Stack[:len(v.Stack)-1]
}

func (v *TransformStack) VisitLight(l *seen.Light) {
}

func (v *TransformStack) VisitSurface(s *seen.Surface) {
}

func (v *TransformStack) EnterShape(s *seen.Shape) {
	// fmt.Printf("Enter Shape %s\n", s.Type)
	v.Transform = v.Transform.Mul(s.Matrix())
}

func (v *TransformStack) LeaveShape(s *seen.Shape) {
	// fmt.Println("Leave Shape")
}

func (v *TransformStack) EnterGroup(m *seen.Group) {
	// fmt.Println("Enter Group")
	v.Transform = v.Transform.Mul(m.Matrix())
}

func (v *TransformStack) LeaveGroup(m *seen.Group) {
	// fmt.Println("Leave Group")
}
