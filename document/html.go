package document

import (
	"os"
	"io"
	"errors"
	"strconv"
)

// HTML
type HTML struct {
	*Element
	body *Element
}

// MakeHTML
func MakeHTML() (*HTML, error) {
	html := CreateElementNS("", "html")
	if html == nil {
		return nil,errors.New("Expected to be able to create a html element")
	}
	head := CreateElementNS("", "head")
	if head == nil {
		return nil,errors.New("Expected to be able to create a head element")
	}
	html.AppendChild(head)
	body := CreateElementNS("", "body")
	if body == nil {
		return nil,errors.New("Expected to be able to create a body element")
	}
	html.AppendChild(body)
	return &HTML{html,body}, nil
}

// AddSVG
func (html *HTML) AddSVG(id string, width, height int) (*SVG,error) {
	svg := CreateElementNS(SVG_NS, "svg")
	if svg == nil {
		return nil,errors.New("Expected to be able to create a svg element")
	}
	svg.SetAttribute("width", strconv.Itoa(width))
	svg.SetAttribute("height", strconv.Itoa(height))
	svg.SetAttribute("id", id)
	html.body.AppendChild(svg)

	// Put a colored background inside the svg
	rect := CreateElementNS(SVG_NS, "rect")
	rect.SetAttribute("width", "100%")
	rect.SetAttribute("height", "100%")
	rect.SetAttribute("rx","5")
	rect.SetAttribute("ry","5")
	rect.SetAttribute("style", "fill: #eeddff")
	svg.AppendChild(rect)

	return &SVG{svg},nil
}

// AddCanvas
func (html *HTML) AddCanvas(id string, width, height int) (*Element, error) {
	canvas := CreateElementNS("","canvas")
	if canvas == nil {
		return nil, errors.New("Expected to be able to create a canvas element")
	}
	canvas.SetAttribute("width",strconv.Itoa(width))
	canvas.SetAttribute("height",strconv.Itoa(height))
	canvas.SetAttribute("id", id)
	html.body.AppendChild(canvas)
	return canvas,nil
}

// SaveToFile
func (html *HTML) SaveToFile(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.WriteString(f, `<!DOCTYPE html>`+"\n")
	if err != nil {
		return err
	}
	_, err = html.WriteTo(f)
	if err != nil {
		return err
	}
	return nil
}
