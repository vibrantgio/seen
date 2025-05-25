package svg

import "github.com/vibrantgio/seen/canvas"

// Canvas
type Canvas struct {
	group  *Element
	path   *Path
	text   *Text
	circle *Circle
	rect   *Rect
	i      int
}

var _ canvas.Canvas = (*Canvas)(nil)

func NewCanvas(group *Element) *Canvas {
	c := &Canvas{}
	c.group = group
	c.path = newPath(c.elementFactory)
	c.text = newText(c.elementFactory)
	c.circle = newCircle(c.elementFactory)
	c.rect = newRect(c.elementFactory)
	return c
}

// Returns an element with tagname `type`.
//
// This method uses an internal iterator to add new elements as they are
// drawn. If there is no child element at the current index, we append one
// and return it. If an element exists at the current index and it is the
// same tag, we return that. If the element is a different type, we create
// one and replace it then return it.
func (c *Canvas) elementFactory(tag string) *Element {
	children := c.group.ChildNodes
	if c.i >= len(children) {
		path := c.group.CreateElementNS(SVG_NS, tag)
		c.group.AppendChild(path)
		c.i++
		return path
	}

	current := children[c.i]
	if current.Tag == tag {
		c.i++
		return current
	}

	path := c.group.CreateElementNS(SVG_NS, tag)
	c.group.ReplaceChild(path, current)
	c.i++
	return path
}

func (c *Canvas) Path() canvas.PathPainter {
	c.path.Clear()
	return c.path
}

func (c *Canvas) Rect() canvas.RectPainter {
	c.rect.Clear()
	return c.rect
}

func (c *Canvas) Circle() canvas.CirclePainter {
	c.circle.Clear()
	return c.circle
}

func (c *Canvas) Text() canvas.TextPainter {
	return c.text
}

func (c *Canvas) Reset() {
	c.i = 0
}

func (c *Canvas) Cleanup() {
	children := c.group.ChildNodes
	for c.i < len(children) {
		children[c.i].SetAttribute("style", "display: none;")
		c.i++
	}
}
