package light

type Kind int

const (
	AmbientKind Kind = iota
	PointKind
	DirectionalKind
)
