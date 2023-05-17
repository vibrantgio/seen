package color

// Source is an interface to a color source. It is used for generating
// a sequence of random colors that are slightly different.
type Source interface {
	Read() Color
}
