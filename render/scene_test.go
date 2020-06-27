package render

import (
	"os"
	"math"
	"math/rand"
	"path"
	"strconv"
	"testing"
	"github.com/reactivego/seen/document"
	"github.com/reactivego/seen"
	"github.com/reactivego/seen/colors"
	"github.com/reactivego/seen/shapes"
	"github.com/reactivego/seen/transform"
)

func TestDemoEmpty(t *testing.T) {
	// Clear the current document (really needed!)
	document.Reset()

	const width = 450
	const height = 200

	_, err := document.MakeSVG("my-3d-svg", width, height)
	if  err != nil {
		t.Error(err)
		return
	}

	s := MakeRenderScene()
	if s == nil {
		t.Error("unable to create render scene")
		return
	}

	c := MakeRenderContext("my-3d-svg", s)
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

	context := MakeRenderContext(svgId, MakeFillLayer(width,height,8,8,"#eeddff"))
	if context == nil {
		t.Error("Expected to be able to create RenderContext")
		return
	}

	// Create the scene to render
	s := MakeRenderScene()
	s.FractionalPoints = true
	s.Model = seen.MakeDefaultModel()
	s.Shader = seen.MakePhongShader()
	s.Camera = seen.MakeCameraWithProjection(seen.MakeDefaultPerspectiveProjection())
	s.Camera.SetTranslation(0,0,-550)

	colorReader := colors.MakeRandomColorReader2(colors.ColorDrift(0.03), colors.ColorSat(0.5))

	// Add icosahedron to the scene
	icosahedron := shapes.MakeIcosahedron()
	if icosahedron == nil {
		t.Fail()
		return
	}
	scale := float64(400) * 0.3
	icosahedron.SetScale(scale, scale, scale)
	icosahedron.SetRotation(transform.MakeQuatAxisAngle(1,1,0,0.25 * math.Pi))
	err = icosahedron.ColorSurfaces(colorReader)
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
	cube.SetScale(scale,scale,scale)
	cube.SetRotation(transform.MakeQuatAxisAngle(0.1,1,0,0.1 * math.Pi))
	cube.SetTranslation(-350,0,0)
	err = cube.ColorSurfaces(colorReader)
	if err != nil {
		t.Error(err)
		return
	}
	s.Model.Add(cube)

	s.Viewport = seen.MakeCenterViewport(0, 0, width, height)
	// Add scene as a layer to the render context
	context.Layer(s)

	// Actually render the scene on the context
	context.Render()

	// Save the generated svg to file
	err = svg.SaveToFile(path.Join(os.Getenv("HOME"), "TestDemoSimple.svg"))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDemoSimpleInteractive(t *testing.T) {
/*
  width  = 900
  height = 500

  # Create sphere shape with randomly colored surfaces
  shape = seen.Shapes.sphere(2).scale(height * 0.4)
  seen.Colors.randomSurfaces2(shape)

  # Create scene and add shape to model
  scene = new seen.Scene
    model    : seen.Models.default().add(shape)
    viewport : seen.Viewports.center(width, height)

  # Create render context from canvas
  context = seen.Context('seen-canvas', scene).render()

  # Slowly rotate sphere
  context.animate()
    .onBefore((t, dt) -> shape.rotx(dt*1e-4).roty(0.7*dt*1e-4))
    .start()

  # Enable drag-to-rotate on the canvas
  dragger = new seen.Drag('seen-canvas', {inertia : true})
  dragger.on('drag.rotate', (e) ->
    xform = seen.Quaternion.xyToTransform(e.offsetRelative...)
    shape.transform(xform)
    context.render()
  )
*/
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
	for i:=0; i<4; i++ {
		html.AddCanvas("seen-canvas-"+strconv.Itoa(i),width,height)
		html.AddSVG("seen-svg-"+strconv.Itoa(i),width,height)
	}

	// Create one shape to be shared between the SVG and Canvas
	spheres := []*seen.Shape{}
	for i := 0; i < 4; i++ {
		sphere := shapes.MakeSphere(i)
		scale := float64(height) * 0.4
		sphere.SetScale(scale, scale, scale)
		colorReader := colors.MakeRandomColorReader2(colors.ColorDrift(0.03), colors.ColorSat(0.5))
		err := sphere.ColorSurfaces(colorReader)
		if err != nil {
			t.Error(err)
			return
		}
		spheres = append(spheres, sphere)
	}

	// Create one scene for each shape
	scenes := []*RenderScene{}
	for _, sphere := range spheres {
		s := MakeRenderScene()
		s.Shader = seen.MakePhongShader()
		s.FractionalPoints = true
		s.Model = seen.MakeDefaultModel()
		s.Model.Add(sphere)
		s.Viewport = seen.MakeCenterViewport(0, 0, width, height)
		scenes = append(scenes, s)
	}

	// Create a render context for each SVG and Canvas
	contexts := []RenderContext{}
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

func TestDemoText(t *testing.T) {
	// Clear the current document (really needed!)
	document.Reset()

	const width  = 900
	const height = 500

	svg, err := document.MakeSVG("seen-svg", width, height)
	if  err != nil {
		t.Error(err)
		return
	}

	// Generate some random data points
	var data []float64 
	for  i:=0; i<10; i++ {
		data = append(data, rand.Float64() * 80.0 + 20.0)
	}

	// Create scene model
	model := seen.MakeDefaultModel()

	// Draw bars for data
	for i,d := range data {
		uc := shapes.MakeUnitCube()
		uc.SetScale(20.0, d, 20.0)
		uc.SetTranslation(float64(i) * 30.0, 0, 0)
		uc.SetFillMaterial("#0088FF")
		model.Add(uc)
	}

	// Draw text above bars
	for i,d := range data {
		opts := map[string]string {
			"font": "10px Roboto",
			"anchor": "middle",
		}
		t := shapes.MakeText(strconv.FormatFloat(d, 'f', 1, 64), opts)
		t.SetShowBackfaces(true)
		t.SetTranslation(float64(i) * 30 + 10, d + 10, 10)
		t.SetFillMaterial("#000000")
		model.Add(t)
	}

	// Create scene
	scene := MakeRenderScene()
	model.SetTranslation(-150, -50, 0)
	model.SetScale(2,2,2)
	scene.Model = model
	scene.Viewport = seen.MakeCenterViewport(0, 0, width, height)
	scene.Camera.SetRotation(transform.MakeQuatAxisAngle(0.1,1,0,math.Pi *0.2))

	// Create render context from canvas
	context := MakeRenderContext("seen-svg", scene)
	if context == nil {
		t.Error("Render context is nil")
		return
	}
	context.Render()

/*
	// Enable drag-to-rotate on the canvas
	dragger = new seen.Drag('seen-canvas', {inertia : true})
	dragger.on('drag.rotate', (e) ->
		xform = seen.Quaternion.xyToTransform(e.offsetRelative...)
		model.transform(xform)
		context.render()
	)
*/

	// Save the generated svg to file
	err = svg.SaveToFile(path.Join(os.Getenv("HOME"), "TestDemoText.svg"))
	if err != nil {
		t.Error(err)
		return
	}
}
