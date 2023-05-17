package svg

import (
	"github.com/reactivego/seen/document"
	"github.com/reactivego/seen/render"
)

// Painter
type Painter struct {
	group         *document.Element
	pathPainter   render.PathPainter
	textPainter   render.TextPainter
	circlePainter render.CirclePainter
	rectPainter   render.RectPainter
	i             int
}

var _ render.Painter = (*Painter)(nil)

func NewPainter(group *document.Element) *Painter {
	c := &Painter{}
	c.group = group
	c.pathPainter = NewPathPainter(c.elementFactory)
	c.textPainter = NewTextPainter(c.elementFactory)
	c.circlePainter = NewCirclePainter(c.elementFactory)
	c.rectPainter = NewRectPainter(c.elementFactory)
	return c
}

// Returns an element with tagname `type`.
//
// This method uses an internal iterator to add new elements as they are
// drawn. If there is no child element at the current index, we append one
// and return it. If an element exists at the current index and it is the
// same tag, we return that. If the element is a different type, we create
// one and replace it then return it.
func (c *Painter) elementFactory(tag string) *document.Element {
	children := c.group.ChildNodes
	if c.i >= len(children) {
		path := c.group.CreateElementNS(document.SVG_NS, tag)
		c.group.AppendChild(path)
		c.i++
		return path
	}

	current := children[c.i]
	if current.Tag == tag {
		c.i++
		return current
	}

	path := c.group.CreateElementNS(document.SVG_NS, tag)
	c.group.ReplaceChild(path, current)
	c.i++
	return path
}

func (c *Painter) Path() render.PathPainter {
	return c.pathPainter
}

func (c *Painter) Rect() render.RectPainter {
	return c.rectPainter
}

func (c *Painter) Circle() render.CirclePainter {
	return c.circlePainter
}

func (c *Painter) Text() render.TextPainter {
	return c.textPainter
}

func (c *Painter) Reset() {
	c.i = 0
}

func (c *Painter) Cleanup() {
	children := c.group.ChildNodes
	for c.i < len(children) {
		children[c.i].SetAttribute("style", "display: none;")
		c.i++
	}
}
