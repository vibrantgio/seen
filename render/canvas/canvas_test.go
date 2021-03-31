package canvas

import (
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/reactivego/seen"
	"github.com/reactivego/seen/colors"
	"github.com/reactivego/seen/document"
	"github.com/reactivego/seen/render"
	"github.com/reactivego/seen/shapes"
	"github.com/reactivego/seen/transform"
)

// Mocks

type MockRenderContextScene struct {
}

func (s *MockRenderContextScene) Paint(context render.PaintContext) {
	// Generate a RenderModel for every Surface in the Scene.

	// Sort the RenderModels based on z-depth back to front.

	// Paint the RenderModels on the PaintContext.
}

// Tests

func TestMakeRenderContext(t *testing.T) {
	document.Reset()

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

	c.Render()
}

func TestDemoSvgCanvas(t *testing.T) {
	// Clear the current document (really needed!)
	document.Reset()

	const width = 450
	const height = 200

	html, err := document.MakeHTML()
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 4; i++ {
		html.AddCanvas("seen-canvas-"+strconv.Itoa(i), width, height)
		html.AddSVG("seen-svg-"+strconv.Itoa(i), width, height)
	}

	// Create one shape to be shared between the SVG and Canvas
	spheres := []*seen.Shape{}
	for i := 0; i < 4; i++ {
		sphere := shapes.MakeSphere(i)
		scale := float64(height) * 0.4
		sphere.SetScale(scale, scale, scale)
		source := colors.MakeRandomSource2(colors.Drift(0.03), colors.Sat(0.5))
		err := sphere.ColorSurfaces(source)
		if err != nil {
			t.Error(err)
			return
		}
		spheres = append(spheres, sphere)
	}

	// Create one scene for each shape
	scenes := []*render.SceneLayer{}
	for _, sphere := range spheres {
		s := seen.MakeScene()
		s.Shader = seen.MakePhongShader()
		s.FractionalPoints = true
		s.Model = seen.MakeDefaultModel()
		s.Model.Add(sphere)
		s.Viewport = seen.MakeCenterViewport(0, 0, width, height)
		scenes = append(scenes, render.MakeSceneLayer(s))
	}

	// Create a render context for each SVG and Canvas
	contexts := []render.RenderContext{}
	for i, scene := range scenes {
		for _, kind := range []string{ /*"canvas",*/ "svg"} {
			elementId := "seen-" + kind + "-" + strconv.Itoa(i)
			context := MakeRenderContext(elementId, scene)
			if context == nil {
				t.Errorf("Expected %q to be present", elementId)
				return
			}
			context.Render()
			contexts = append(contexts, context)
		}
	}

	// Slowly rotate shapes
	a := seen.MakeAnimator()
	a.OnFrame(func(t, dt float64) {
		for _, sphere := range spheres {
			ryrx := transform.MakeQuatRotY(dt * 2e-4).MulRotX(dt * 3e-4)
			sphere.SetRotation(ryrx.Mul(sphere.Rotation()))
		}
		for _, context := range contexts {
			context.Render()
		}
	})
	a.Start()

	// Save the generated html element to file
	html.SaveToFile(path.Join(os.Getenv("HOME"), "TestDemoSvgCanvas.html"))
}
