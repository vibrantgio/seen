package render

import (
	"testing"
	"github.com/reactivego/seen/document"
)

// Mocks

type MockRenderContextScene struct {
}

func (s *MockRenderContextScene) Paint(context PaintContext) {
	// Generate a RenderModel for every Surface in the Scene.

	// Sort the RenderModels based on z-depth back to front.

	// Paint the RenderModels on the PaintContext.
}

// Tests

func TestMakeRenderContext(t *testing.T) {
	document.Reset()
	svg := document.CreateElementNS(document.SVG_NS, "svg")
	if svg == nil {
		t.Error("Expected a valid svg element")
	}
	svg.SetAttribute("id", "my-3d-svg")

	s := &MockRenderContextScene{}

	c := MakeRenderContext("invalid", s)
	if c != nil {
		t.Error("Expected MakeRenderContext to return nil")
	}

	canvas := document.CreateElementNS("", "canvas")
	if canvas == nil {
		t.Error("Expected a valid canvas element")
	}
	canvas.SetAttribute("id", "my-3d-canvas")

	c = MakeRenderContext("my-3d-canvas", s)
	if c == nil {
		t.Error("Expected to get a render context for valid canvas element.")
	}

	c = MakeRenderContext("my-3d-svg", s)
	if c == nil {
		t.Error("Expected to get a render context for valid svg element.")
	}

	c.Render()
}
