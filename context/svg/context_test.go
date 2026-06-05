package svg_test

import (
	"math"
	"math/rand"
	"path"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/affine"
	"github.com/vibrantgio/seen/animation"
	"github.com/vibrantgio/seen/camera"
	"github.com/vibrantgio/seen/canvas"
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/context/svg"
	"github.com/vibrantgio/seen/layer"
	"github.com/vibrantgio/seen/layer/backdrop"
	"github.com/vibrantgio/seen/layer/zsort"
	"github.com/vibrantgio/seen/point"
	"github.com/vibrantgio/seen/quaternion"
	"github.com/vibrantgio/seen/shader"
	"github.com/vibrantgio/seen/shape"
	"github.com/vibrantgio/seen/viewport"
)

// Helpers

func SourceFile(filename string) string {
	_, sourcepath, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(sourcepath), filename)
}

// Mocks

type MockSceneLayer struct {
}

func (s *MockSceneLayer) RenderOn(canvas canvas.Canvas) {
	// Generate a Render for every face in the Scene.

	// Sort the Render list based on z-depth back to front.

	// Paint the Render list on the PaintContext.
}

// Tests

func TestContextWith(t *testing.T) {
	dom := svg.NewDom()
	root := dom.CreateElementNS(svg.SVG_NS, "svg")
	if root == nil {
		t.Error("Expected a valid svg element")
	}
	root.SetAttribute("id", "my-3d-svg")

	s := &MockSceneLayer{}

	c := svg.NewContext(nil, s)
	if c != nil {
		t.Error("Expected ContextWith to return nil")
	}

	if root != nil {
		c = svg.NewContext(root.GetElementById("my-3d-svg"), s)
		if c == nil {
			t.Error("Expected to get a render context for valid svg element.")
		}

		c.Render()
	}
}

// Mock Canvas Layer

type MockCanvasLayer struct {
	Width, Height, Radius float64
	RectFill              string
	CircleFill            string
	Text                  string
	TextStyle             canvas.Style
}

func (l MockCanvasLayer) RenderOn(c canvas.Canvas) {
	c.Rect().Rect(l.Width, l.Height).Fill(canvas.Style{"fill": l.RectFill})

	center := point.Pt(l.Width/2, l.Height/2, 0.0)
	c.Circle().Circle(center, l.Radius).Fill(canvas.Style{"fill": l.CircleFill})

	transform := affine.Matrix{A: 1, D: 1, E: l.Width / 2, F: l.Height / 2}
	c.Text().FillText(transform, l.Text, l.TextStyle)
}

func TestPaintLayer(t *testing.T) {
	// create simple svg file
	const width = 450
	const height = 400

	// Create the context to render to (svg in this case) and put a colored background in.
	svgId := "my-3d-svg"
	doc, err := svg.NewSVG(svgId, width, height)
	if err != nil {
		t.Error(err)
		return
	}

	context := svg.NewContext(doc.GetElementById(svgId))
	if context == nil {
		t.Error("Expected to be able to create render.Context")
		return
	}

	layer := MockCanvasLayer{width, height, height / 3, "#ccccff", "#ff8888", "Hello, World!", canvas.Style{
		"font-family": "Roboto",
		"font-weight": "bold",
		"font-size":   "24px",
		"fill":        "#000000",
		"text-anchor": "middle",
		"inline-size": "100px",
	}}
	context.SetLayers(layer)
	context.Render()

	// Save the generated svg to file
	err = doc.SaveToFile(SourceFile("TestPaintLayer.svg"))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDemoEmpty(t *testing.T) {
	const width = 450
	const height = 200

	doc, err := svg.NewSVG("my-3d-svg", width, height)
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

	context := svg.NewContext(doc.GetElementById("my-3d-svg"), l)
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
	doc, err := svg.NewSVG(svgId, width, height)
	if err != nil {
		t.Error(err)
		return
	}

	context := svg.NewContext(doc.GetElementById(svgId))
	if context == nil {
		t.Error("Expected to be able to create render.Context")
		return
	}

	blueish, _ := color.ColorWithString("#eeddff")
	backdrop := backdrop.NewLayer(width, height, 8, 8, blueish)

	// Create the scene to render
	s := seen.NewDefaultScene()
	s.Shader = shader.Phong
	s.Camera = camera.Default
	s.Camera.SetTranslation(0, 0, -550)

	source := color.NewDriftingSourceWith(color.Drift(0.03), color.Sat(0.5))

	// Add icosahedron to the scene
	icosahedron := shape.Icosahedron()
	scale := float64(400) * 0.3
	icosahedron.SetScale(scale, scale, scale)
	icosahedron.SetRotation(quaternion.AxisAngle(1, 1, 0, 0.25*math.Pi))
	err = icosahedron.Faces().SetColorFrom(source)
	if err != nil {
		t.Error(err)
		return
	}
	s.Group.Add(icosahedron)

	// Add a cube to the scene
	cube := shape.UnitCube()
	cube.SetScale(scale, scale, scale)
	cube.SetRotation(quaternion.AxisAngle(0.1, 1, 0, 0.1*math.Pi))
	cube.SetTranslation(-350, 0, 0)
	err = cube.Faces().SetColorFrom(source)
	if err != nil {
		t.Error(err)
		return
	}
	s.Group.Add(cube)

	s.Viewport = viewport.Center(0, 0, width, height)
	// Add scene as a layer to the render context
	context.SetLayers(backdrop, zsort.NewLayerForScene(s))

	// Actually render the scene on the context
	context.Render()

	// Save the generated svg to file
	err = doc.SaveToFile(SourceFile("TestDemoSimple.svg"))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDemoSvgCanvas(t *testing.T) {
	const width = 450
	const height = 200

	html, err := svg.NewHTML()
	if err != nil {
		t.Error(err)
		return
	}
	for i := range 4 {
		html.AddCanvas("seen-canvas-"+strconv.Itoa(i), width, height)
		html.AddSVG("seen-svg-"+strconv.Itoa(i), width, height)
	}

	// Create one shape to be shared between the SVG and Canvas
	spheres := []seen.Node(nil)
	for i := range 4 {
		sphere := shape.Sphere(i)
		scale := float64(height) * 0.4
		sphere.SetScale(scale, scale, scale)
		source := color.NewDriftingSourceWith(color.Drift(0.03), color.Sat(0.5))
		err := sphere.Faces().SetColorFrom(source)
		if err != nil {
			t.Error(err)
			return
		}
		spheres = append(spheres, sphere)
	}

	// Create one scene for each shape
	scenes := []layer.Layer{}
	for _, sphere := range spheres {
		scene := seen.NewDefaultScene()
		scene.Shader = shader.Phong
		scene.Group.Add(sphere)
		scene.Viewport = viewport.Center(0, 0, width, height)
		scenes = append(scenes, zsort.NewLayerForScene(scene))
	}

	// Create a render context for each SVG and Canvas
	contexts := []*svg.Context{}
	for i, scene := range scenes {
		for _, kind := range []string{ /*"canvas",*/ "svg"} {
			elementId := "seen-" + kind + "-" + strconv.Itoa(i)
			context := svg.NewContext(html.GetElementById(elementId), scene)
			if context == nil {
				t.Errorf("Expected %q to be present", elementId)
				return
			}
			context.Render()
			contexts = append(contexts, context)
		}
	}

	// Slowly rotate shapes
	a := &animation.Animation{}
	a.OnFrame(func(t, dt time.Duration) {
		for _, sphere := range spheres {
			dtms := float64(dt.Milliseconds())
			ryrx := quaternion.RotY(dtms * 2e-4).RotX(dtms * 3e-4)
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

	doc, err := svg.NewSVG("seen-svg", width, height)
	if err != nil {
		t.Error(err)
		return
	}

	// Generate some random data points
	var data []float64
	for range 10 {
		data = append(data, rand.Float64()*80.0+20.0)
	}

	// Create scene
	scene := seen.NewDefaultScene()

	// Draw bars for data
	for i, d := range data {
		uc := shape.UnitCube()
		uc.SetScale(20.0, d, 20.0)
		uc.SetTranslation(float64(i)*30.0, 0, 0)
		uc.Faces().SetFill("#0088FF")
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
		t.Faces().SetShowBackfaces(true)
		t.SetTranslation(float64(i)*30+10, d+10, 10)
		t.Faces().SetFill("#000000")
		scene.Group.Add(t)
	}

	// Create scene
	scene.Group.SetTranslation(-150, -50, 0)
	scene.Group.SetScale(2, 2, 2)
	scene.Viewport = viewport.Center(0, 0, width, height)
	scene.Camera.SetRotation(quaternion.AxisAngle(0.1, 1, 0, math.Pi*0.2))

	// Create render context from canvas
	context := svg.NewContext(doc.GetElementById("seen-svg"), zsort.NewLayerForScene(scene))
	if context == nil {
		t.Error("Render context is nil")
		return
	}
	context.Render()

	// Save the generated svg to file
	err = doc.SaveToFile(SourceFile("TestDemoText.svg"))
	if err != nil {
		t.Error(err)
		return
	}
}
