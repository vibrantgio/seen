package svg

import (
	"math"
	"math/rand"
	"os"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/reactivego/seen"
	"github.com/reactivego/seen/colors"
	"github.com/reactivego/seen/document"
	"github.com/reactivego/seen/render"
	"github.com/reactivego/seen/shapes"
	"github.com/reactivego/seen/transform"
)

// Mocks

type MockSceneLayer struct {
}

func (s *MockSceneLayer) Paint(painter render.Painter) {
	// Generate a RenderModel for every Surface in the Scene.

	// Sort the RenderModels based on z-depth back to front.

	// Paint the RenderModels on the PaintContext.
}

// Tests

func TestMakeContext(t *testing.T) {
	document.Reset()
	svg := document.CreateElementNS(document.SVG_NS, "svg")
	if svg == nil {
		t.Error("Expected a valid svg element")
	}
	svg.SetAttribute("id", "my-3d-svg")

	s := &MockSceneLayer{}

	c := MakeContext("invalid", s)
	if c != nil {
		t.Error("Expected MakeContext to return nil")
	}

	c = MakeContext("my-3d-svg", s)
	if c == nil {
		t.Error("Expected to get a render context for valid svg element.")
	}

	c.Render()
}

func TestDemoEmpty(t *testing.T) {
	// Clear the current document (really needed!)
	document.Reset()

	const width = 450
	const height = 200

	_, err := document.MakeSVG("my-3d-svg", width, height)
	if err != nil {
		t.Error(err)
		return
	}

	s := seen.EmptyScene()
	l := render.MakeSceneLayer(&s)
	if l == nil {
		t.Error("unable to create scene layer")
		return
	}

	c := MakeContext("my-3d-svg", l)
	if c == nil {
		t.Error("unable to find element my-3d-svg")
		return
	}

	c.Render()
}

func TestDemoSimple(t *testing.T) {
	// Clear the current document (really needed!)
	document.Reset()

	// create simple svg file
	const width = 450
	const height = 400

	// Create the context to render to (svg in this case) and put a colored background in.
	svgId := "my-3d-svg"
	svg, err := document.MakeSVG(svgId, width, height)
	if err != nil {
		t.Error(err)
		return
	}

	context := MakeContext(svgId, render.MakeFillLayer(width, height, 8, 8, "#eeddff"))
	if context == nil {
		t.Error("Expected to be able to create RenderContext")
		return
	}

	// Create the scene to render
	s := seen.DefaultScene()
	s.FractionalPoints = true
	s.Shader = seen.PhongShader
	s.Camera = seen.DefaultCamera
	s.Camera.SetTranslation(0, 0, -550)

	source := colors.MakeRandomSource2(colors.Drift(0.03), colors.Sat(0.5))

	// Add icosahedron to the scene
	icosahedron := shapes.MakeIcosahedron()
	if icosahedron == nil {
		t.Fail()
		return
	}
	scale := float64(400) * 0.3
	icosahedron.SetScale(scale, scale, scale)
	icosahedron.SetRotation(transform.QuatAxisAngle(1, 1, 0, 0.25*math.Pi))
	err = icosahedron.ColorSurfaces(source)
	if err != nil {
		t.Error(err)
		return
	}
	s.Model.Add(icosahedron)

	// Add a cube to the scene
	cube := shapes.MakeUnitCube()
	if cube == nil {
		t.Fail()
		return
	}
	cube.SetScale(scale, scale, scale)
	cube.SetRotation(transform.QuatAxisAngle(0.1, 1, 0, 0.1*math.Pi))
	cube.SetTranslation(-350, 0, 0)
	err = cube.ColorSurfaces(source)
	if err != nil {
		t.Error(err)
		return
	}
	s.Model.Add(cube)

	s.Viewport = seen.CenterViewport(0, 0, width, height)
	// Add scene as a layer to the render context
	context.Layer(render.MakeSceneLayer(&s))

	// Actually render the scene on the context
	context.Render()

	// Save the generated svg to file
	err = svg.SaveToFile(path.Join(os.Getenv("HOME"), "TestDemoSimple.svg"))
	if err != nil {
		t.Error(err)
		return
	}
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
		s := seen.DefaultScene()
		s.Shader = seen.PhongShader
		s.FractionalPoints = true
		s.Model.Add(sphere)
		s.Viewport = seen.CenterViewport(0, 0, width, height)
		scenes = append(scenes, render.MakeSceneLayer(&s))
	}

	// Create a render context for each SVG and Canvas
	contexts := []render.RenderContext{}
	for i, scene := range scenes {
		for _, kind := range []string{ /*"canvas",*/ "svg"} {
			elementId := "seen-" + kind + "-" + strconv.Itoa(i)
			context := MakeContext(elementId, scene)
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
	a.OnFrame(func(t, dt time.Duration) {
		for _, sphere := range spheres {
			dtms := float64(dt.Milliseconds())
			ryrx := transform.QuatRotY(dtms * 2e-4).MulRotX(dtms * 3e-4)
			sphere.SetRotation(ryrx.Mul(sphere.Rotation()))
		}
		for _, context := range contexts {
			context.Render()
		}
	})
	a.Start()
	time.Sleep(100 * time.Millisecond)
	a.Stop()

	// Save the generated html element to file
	html.SaveToFile(path.Join(os.Getenv("HOME"), "TestDemoSvgCanvas.html"))
}

func TestDemoText(t *testing.T) {
	// Clear the current document (really needed!)
	document.Reset()

	const width = 900
	const height = 500

	svg, err := document.MakeSVG("seen-svg", width, height)
	if err != nil {
		t.Error(err)
		return
	}

	// Generate some random data points
	var data []float64
	for i := 0; i < 10; i++ {
		data = append(data, rand.Float64()*80.0+20.0)
	}

	// Create scene model
	scene := seen.DefaultScene()

	// Draw bars for data
	for i, d := range data {
		uc := shapes.MakeUnitCube()
		uc.SetScale(20.0, d, 20.0)
		uc.SetTranslation(float64(i)*30.0, 0, 0)
		uc.SetFill("#0088FF")
		scene.Model.Add(uc)
	}

	// Draw text above bars
	for i, d := range data {
		opts := map[string]string{
			//"font": "10px Roboto",
			"font-family": "Roboto",
			"font-size":   "10px",
			"anchor":      "middle",
		}
		t := shapes.MakeText(strconv.FormatFloat(d, 'f', 1, 64), opts)
		t.SetShowBackfaces(true)
		t.SetTranslation(float64(i)*30+10, d+10, 10)
		t.SetFill("#000000")
		scene.Model.Add(t)
	}

	// Create scene
	scene.Model.SetTranslation(-150, -50, 0)
	scene.Model.SetScale(2, 2, 2)
	scene.Viewport = seen.CenterViewport(0, 0, width, height)
	scene.Camera.SetRotation(transform.QuatAxisAngle(0.1, 1, 0, math.Pi*0.2))

	// Create render context from canvas
	context := MakeContext("seen-svg", render.MakeSceneLayer(&scene))
	if context == nil {
		t.Error("Render context is nil")
		return
	}
	context.Render()

	// Save the generated svg to file
	err = svg.SaveToFile(path.Join(os.Getenv("HOME"), "TestDemoText.svg"))
	if err != nil {
		t.Error(err)
		return
	}
}
