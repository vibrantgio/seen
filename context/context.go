package context

import (
	"github.com/vibrantgio/seen/animation"
	"github.com/vibrantgio/seen/drag"
	"github.com/vibrantgio/seen/layer"
	"github.com/vibrantgio/seen/zoom"
)

type Context interface {
	SetLayers(layers ...layer.Layer)
	Render()
	Animate() animation.Animator
	Drag(options ...drag.Option) drag.Dragger
	Zoom(options ...zoom.Option) zoom.Zoomer
}
