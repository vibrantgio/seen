package seen

import "github.com/vibrantgio/seen/face"

type Object interface {
	Node
	Faces() face.Faces
}
