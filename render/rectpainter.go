package render

// RectPainter
type RectPainter interface {
	Size(width, height float64)
	CornerRadius(rx, ry float64)

	// Fill the rect
	Fill(Style)
}
