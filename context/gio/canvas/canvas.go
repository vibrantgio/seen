package canvas

import (
	"gioui.org/op"
	"github.com/vibrantgio/seen/canvas"
)

type Canvas struct {
	*op.Ops
	path   Path
	rect   Rect
	circle Circle
	text   Text
}

func (c *Canvas) Path() canvas.PathPainter {
	c.path.Ops = c.Ops
	return &c.path
}

func (c *Canvas) Rect() canvas.RectPainter {
	c.rect.Ops = c.Ops
	return &c.rect
}

func (c *Canvas) Circle() canvas.CirclePainter {
	c.circle.Ops = c.Ops
	return &c.circle
}

func (c *Canvas) Text() canvas.TextPainter {
	c.text.Ops = c.Ops
	return &c.text
}

func (c *Canvas) Reset() {

}

func (c *Canvas) Cleanup() {

}
