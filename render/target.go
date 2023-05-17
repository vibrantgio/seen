package render

import (
	"github.com/reactivego/seen"
)

type Target interface {
	SetLayers(layers ...Layer)
	Render()
	Animate() seen.Animator
	Drag(options ...seen.DragOption) seen.Dragger
	Zoom(options ...seen.ZoomOption) seen.Zoomer
}
