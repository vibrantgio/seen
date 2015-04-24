package seen

import (
	"strconv"
	"testing"
	"xpt.nl/document"
)

// Mocks

type MockRenderContextScene struct {
}

func (s *MockRenderContextScene) Render(context RenderLayerContext) {
}

// Tests

func TestFormatFloat(t *testing.T) {
	if strconv.FormatFloat(234.553343, 'g', -1, 64) != "234.553343" {
		t.Fail()
	}
}

func TestNewRenderContext(t *testing.T) {
	svg := document.CreateElementNS(SVG_NS, "svg")
	if svg == nil {
		t.Fail()
	}
	svg.SetAttribute("id", "my-3d-svg")

	s := &MockRenderContextScene{}

	c := NewRenderContext("invalid", s)
	if c != nil {
		t.Fail()
	}

	canvas := document.CreateElementNS("", "canvas")
	canvas.SetAttribute("id", "my-3d-canvas")
	c = NewRenderContext("my-3d-canvas", s)
	// canvas not yet supported, so should return nil
	if c != nil {
		t.Fail()
	}

	c = NewRenderContext("my-3d-svg", s)
	if c == nil {
		t.Fail()
	}

	c.Render()
}
