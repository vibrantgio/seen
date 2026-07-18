package main

import (
	"log"
	"math"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/vibrantgio/noise"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/context/gio"
	"github.com/vibrantgio/seen/context/svg"
	"github.com/vibrantgio/seen/layer/nsort"
	"github.com/vibrantgio/seen/quaternion"
	"github.com/vibrantgio/seen/shape"
)

const should_save_to_svg = false

func main() {
	go NoisyWavePatch()
	app.Main()
}

func NoisyWavePatch() {
	width := unit.Dp(900)
	height := unit.Dp(500)

	window := new(app.Window)
	window.Option(
		app.Title("Seen - Noisy Wave Patch"),
		app.Size(width, height))

	// Create a context that hooks seen into the gio window
	context := gio.NewContext(window)

	// Create patch of triangles that spans the view
	equilateralAltitude := math.Sqrt(3.0) / 2.0
	triangleScale := 70.0
	patch_width := float64(width) * 1.5
	patch_height := float64(height) * 1.5

	nx := patch_width / triangleScale / equilateralAltitude
	ny := patch_height / triangleScale
	shape := shape.Patch(nx, ny)
	shape.SetScale(triangleScale, triangleScale, triangleScale)
	shape.SetRotation(quaternion.RotX(-0.3))
	shape.SetTranslation(-patch_width/2, -patch_height/2+80, 0)
	shape.Faces().SetColorFrom(color.NewDriftingSource())

	// Create scene and add shape to group
	scene := seen.NewDefaultScene()
	scene.ShowBackfaces = true
	scene.Group.Add(shape)

	// Create a layer that renders a scene by depth-sorting the polygons for the current eye
	layer := nsort.NewLayerForScene(scene)
	context.SetLayers(layer)

	// Apply animated 3D simplex noise to patch vertices
	noiser := noise.NewSimplex3D(0)
	context.Animate().OnBefore(func(t, dt time.Duration) {
		faces := shape.Faces()
		for i, surf := range faces {
			for j, p := range surf.Points {
				tms := float64(t.Milliseconds())
				//shape.Faces[i].Points[j].Z = 4.0 * noiser.Noise(p.X/8.0, p.Y/8.0, tms*1e-4)
				faces[i].Points[j].Z = noiser.Noise(p.X/8.0, p.Y/8.0, tms*1e-4) / 2
			}
			// Since we're modifying the point directly, we need to mark the face dirty
			// to make sure the cache doesn't ignore the change
			faces[i].Dirty = true
		}
	}).Start()

	view := gio.Widget(context, func(w, h unit.Dp) {
		width, height = w, h
		scene.FitCenter(0, 0, float64(w), float64(h))
	})

	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			view(gtx)
			e.Frame(gtx.Ops)
		}
	}

	// Save scene to svg file
	if should_save_to_svg {
		doc, err := svg.NewSVG("seen-svg", int(width), int(height))
		if err != nil {
			log.Fatal(err)
		}
		if context := svg.NewContext(doc.GetElementById("seen-svg"), layer); context != nil {
			context.Render()
		} else {
			log.Fatal("Render context is nil")
		}
		err = doc.SaveToFile("noisywavepatch.svg")
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(0)
}
