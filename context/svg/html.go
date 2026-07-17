package svg

import (
	"errors"
	"io"
	"os"
	"strconv"
)

// HTML
type HTML struct {
	*Element
	head *Element
	body *Element
}

// NewHTML
func NewHTML() (*HTML, error) {
	html := NewDom().CreateElementNS("", "html")
	if html == nil {
		return nil, errors.New("failed to create html element")
	}
	head := html.CreateElementNS("", "head")
	if head == nil {
		return nil, errors.New("failed to create head element")
	}
	html.AppendChild(head)
	body := html.CreateElementNS("", "body")
	if body == nil {
		return nil, errors.New("failed to create body element")
	}
	html.AppendChild(body)
	return &HTML{html, head, body}, nil
}

// AddSVG
func (html *HTML) AddSVG(id string, width, height int) (*Element, error) {
	svg := html.CreateElementNS(SVG_NS, "svg")
	if svg == nil {
		return nil, errors.New("failed to create svg element")
	}
	svg.SetAttribute("width", strconv.Itoa(width))
	svg.SetAttribute("height", strconv.Itoa(height))
	svg.SetAttribute("id", id)
	html.body.AppendChild(svg)

	// Put a colored background inside the svg
	rect := html.CreateElementNS(SVG_NS, "rect")
	rect.SetAttribute("width", "100%")
	rect.SetAttribute("height", "100%")
	rect.SetAttribute("rx", "5")
	rect.SetAttribute("ry", "5")
	rect.SetAttribute("style", "fill: #eeddff")
	svg.AppendChild(rect)

	return svg, nil
}

// AddCanvas
func (html *HTML) AddCanvas(id string, width, height int) (*Element, error) {
	canvas := html.CreateElementNS("", "canvas")
	if canvas == nil {
		return nil, errors.New("failed to create canvas element")
	}
	canvas.SetAttribute("width", strconv.Itoa(width))
	canvas.SetAttribute("height", strconv.Itoa(height))
	canvas.SetAttribute("id", id)
	html.body.AppendChild(canvas)
	return canvas, nil
}

// WriteDocumentTo writes the complete html document, doctype included, to w.
func (html *HTML) WriteDocumentTo(w io.Writer) (n int64, err error) {
	sn, err := io.WriteString(w, `<!DOCTYPE html>`+"\n")
	n += int64(sn)
	if err != nil {
		return
	}
	en, err := html.WriteTo(w)
	n += en
	return
}

// SaveToFile
func (html *HTML) SaveToFile(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = html.WriteDocumentTo(f)
	return err
}
