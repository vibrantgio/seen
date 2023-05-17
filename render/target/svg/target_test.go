package svg

import (
	"math"
	"math/rand"
	"path"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/reactivego/seen"
	"github.com/reactivego/seen/color"
	"github.com/reactivego/seen/document"
	"github.com/reactivego/seen/quat"
	"github.com/reactivego/seen/render"
	"github.com/reactivego/seen/render/layer/zsort"
	"github.com/reactivego/seen/shape"
)

// Helpers

func SourceFile(filename string) string {
	_, sourcepath, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(sourcepath), filename)
}

// Mocks

type MockSceneLayer struct {
}

func (s *MockSceneLayer) Paint(painter render.Painter) {
	// Generate a render.Surface for every Surface in the Scene.

	// Sort the render.Surface list based on z-depth back to front.

	// Paint the render.Surface list on the PaintContext.
}

// Tests

func TestContextWith(t *testing.T) {
	dom := document.NewDom()
	svg := dom.CreateElementNS(document.SVG_NS, "svg")
	if svg == nil {
		t.Error("Expected a valid svg element")
	}
	svg.SetAttribute("id", "my-3d-svg")

	s := &MockSceneLayer{}

	c := NewTarget(nil, s)
	if c != nil {
		t.Error("Expected ContextWith to return nil")
	}

	if svg != nil {
		c = NewTarget(svg.GetElementById("my-3d-svg"), s)
		if c == nil {
			t.Error("Expected to get a render context for valid svg element.")
		}

		c.Render()
	}
}

func TestDemoEmpty(t *testing.T) {
	const width = 450
	const height = 200

	svg, err := document.NewSVG("my-3d-svg", width, height)
	if err != nil {
		t.Error(err)
		return
	}

	scene := seen.NewScene()
	l := zsort.NewLayerForScene(scene)
	if l == nil {
		t.Error("unable to create scene layer")
		return
	}

	context := NewTarget(svg.GetElementById("my-3d-svg"), l)
	if context == nil {
		t.Error("unable to find element my-3d-svg")
		return
	}

	context.Render()
}

func TestDemoSimple(t *testing.T) {
	// create simple svg file
	const width = 450
	const height = 400

	// Create the context to render to (svg in this case) and put a colored background in.
	svgId := "my-3d-svg"
	svg, err := document.NewSVG(svgId, width, height)
	if err != nil {
		t.Error(err)
		return
	}

	context := NewTarget(svg.GetElementById(svgId))
	if context == nil {
		t.Error("Expected to be able to create render.Context")
		return
	}

	blueish, _ := color.ColorWithString("#eeddff")
	backdrop := render.NewBackdrop(width, height, 8, 8, blueish)

	// Create the scene to render
	s := seen.NewDefaultScene()
	s.FractionalPoints = true
	s.Shader = seen.PhongShader
	s.Camera = seen.DefaultCamera
	s.Camera.SetTranslation(0, 0, -550)

	source := color.NewDriftingSourceWith(color.Drift(0.03), color.Sat(0.5))

	// Add icosahedron to the scene
	icosahedron := shape.Icosahedron()
	scale := float64(400) * 0.3
	icosahedron.SetScale(scale, scale, scale)
	icosahedron.SetRotation(quat.AxisAngle(1, 1, 0, 0.25*math.Pi))
	err = icosahedron.SetColorFrom(source)
	if err != nil {
		t.Error(err)
		return
	}
	s.Group.Add(icosahedron)

	// Add a cube to the scene
	cube := shape.UnitCube()
	cube.SetScale(scale, scale, scale)
	cube.SetRotation(quat.AxisAngle(0.1, 1, 0, 0.1*math.Pi))
	cube.SetTranslation(-350, 0, 0)
	err = cube.SetColorFrom(source)
	if err != nil {
		t.Error(err)
		return
	}
	s.Group.Add(cube)

	s.Viewport = seen.CenterViewport(0, 0, width, height)
	// Add scene as a layer to the render context
	context.SetLayers(backdrop, zsort.NewLayerForScene(s))

	// Actually render the scene on the context
	context.Render()

	// Save the generated svg to file
	err = svg.SaveToFile(SourceFile("TestDemoSimple.svg"))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDemoSvgCanvas(t *testing.T) {
	const width = 450
	const height = 200

	html, err := document.NewHTML()
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
		sphere := shape.Sphere(i)
		scale := float64(height) * 0.4
		sphere.SetScale(scale, scale, scale)
		source := color.NewDriftingSourceWith(color.Drift(0.03), color.Sat(0.5))
		err := sphere.SetColorFrom(source)
		if err != nil {
			t.Error(err)
			return
		}
		spheres = append(spheres, sphere)
	}

	// Create one scene for each shape
	scenes := []*zsort.Layer{}
	for _, sphere := range spheres {
		scene := seen.NewDefaultScene()
		scene.Shader = seen.PhongShader
		scene.FractionalPoints = true
		scene.Group.Add(sphere)
		scene.Viewport = seen.CenterViewport(0, 0, width, height)
		scenes = append(scenes, zsort.NewLayerForScene(scene))
	}

	// Create a render context for each SVG and Canvas
	contexts := []*Target{}
	for i, scene := range scenes {
		for _, kind := range []string{ /*"canvas",*/ "svg"} {
			elementId := "seen-" + kind + "-" + strconv.Itoa(i)
			context := NewTarget(html.GetElementById(elementId), scene)
			if context == nil {
				t.Errorf("Expected %q to be present", elementId)
				return
			}
			context.Render()
			contexts = append(contexts, context)
		}
	}

	// Slowly rotate shapes
	a := &seen.Animation{}
	a.OnFrame(func(t, dt time.Duration) {
		for _, sphere := range spheres {
			dtms := float64(dt.Milliseconds())
			ryrx := quat.RotY(dtms * 2e-4).RotX(dtms * 3e-4)
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
	html.SaveToFile(SourceFile("TestDemoSvgCanvas.html"))
}

func TestDemoText(t *testing.T) {
	const width = 900
	const height = 500

	svg, err := document.NewSVG("seen-svg", width, height)
	if err != nil {
		t.Error(err)
		return
	}

	// Generate some random data points
	var data []float64
	for i := 0; i < 10; i++ {
		data = append(data, rand.Float64()*80.0+20.0)
	}

	// Create scene
	scene := seen.NewDefaultScene()

	// Draw bars for data
	for i, d := range data {
		uc := shape.UnitCube()
		uc.SetScale(20.0, d, 20.0)
		uc.SetTranslation(float64(i)*30.0, 0, 0)
		uc.SetFill("#0088FF")
		scene.Group.Add(uc)
	}

	// Draw text above bars
	for i, d := range data {
		opts := map[string]string{
			//"font": "10px Roboto",
			"font-family": "Roboto",
			"font-size":   "10px",
			"anchor":      "middle",
		}
		t := shape.Text(strconv.FormatFloat(d, 'f', 1, 64), opts)
		t.SetShowBackfaces(true)
		t.SetTranslation(float64(i)*30+10, d+10, 10)
		t.SetFill("#000000")
		scene.Group.Add(t)
	}

	// Create scene
	scene.Group.SetTranslation(-150, -50, 0)
	scene.Group.SetScale(2, 2, 2)
	scene.Viewport = seen.CenterViewport(0, 0, width, height)
	scene.Camera.SetRotation(quat.AxisAngle(0.1, 1, 0, math.Pi*0.2))

	// Create render context from canvas
	context := NewTarget(svg.GetElementById("seen-svg"), zsort.NewLayerForScene(scene))
	if context == nil {
		t.Error("Render context is nil")
		return
	}
	context.Render()

	// Save the generated svg to file
	err = svg.SaveToFile(SourceFile("TestDemoText.svg"))
	if err != nil {
		t.Error(err)
		return
	}
}
