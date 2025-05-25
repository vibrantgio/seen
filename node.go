package seen

// Node represents the fundamental element in a 3D scene graph. This Node
// interface defines a transformable object that encapsulates geometric
// information and face data for rendering.
type Node interface {
	Transformer
	Kind() string
}
