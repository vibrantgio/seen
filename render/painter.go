package render

// Painter
type Painter interface {
	Path() PathPainter
	Rect() RectPainter
	Circle() CirclePainter
	Text() TextPainter

	Reset()
	Cleanup()
}
