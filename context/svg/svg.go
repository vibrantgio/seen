package svg

import (
	"errors"
	"io"
	"os"
	"strconv"
)

const SVG_NS = "http://www.w3.org/2000/svg"

// SVG
type SVG struct{ *Element }

// NewSVG creates a simple document with just an svg element.
// <?xml version="1.0" standalone="yes"?>
// <!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
// <svg xmlns="http://www.w3.org/2000/svg" version="1.1" id="my-3d-svg" width="500" height="400">
// </svg>
func NewSVG(id string, width, height int) (*SVG, error) {
	svg := NewDom().CreateElementNS(SVG_NS, "svg")
	if svg == nil {
		return nil, errors.New("failed to create svg element")
	}
	svg.SetAttribute("xmlns", SVG_NS)
	svg.SetAttribute("version", "1.1")
	if id != "" {
		svg.SetAttribute("id", id)
	}
	svg.SetAttribute("width", strconv.Itoa(width))
	svg.SetAttribute("height", strconv.Itoa(height))
	return &SVG{svg}, nil
}

// SaveToFile
func (svg *SVG) SaveToFile(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.WriteString(f, `<?xml version="1.0" standalone="yes"?>`+"\n")
	if err != nil {
		return err
	}
	_, err = io.WriteString(f, `<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">`+"\n")
	if err != nil {
		return err
	}
	_, err = svg.WriteTo(f)
	if err != nil {
		return err
	}
	return nil
}
