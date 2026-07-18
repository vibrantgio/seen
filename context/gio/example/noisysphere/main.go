package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
	"golang.org/x/exp/slices"

	"github.com/vibrantgio/noise"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/context/gio"
	"github.com/vibrantgio/seen/context/svg"
	"github.com/vibrantgio/seen/drag"
	"github.com/vibrantgio/seen/layer"
	"github.com/vibrantgio/seen/layer/nsort"
	"github.com/vibrantgio/seen/point"
	"github.com/vibrantgio/seen/quaternion"
	"github.com/vibrantgio/seen/shape"
	"github.com/vibrantgio/seen/zoom"
)

const should_save_to_svg = false

func main() {
	go NoisySphere()
	app.Main()
}

func NoisySphere() {
	width := unit.Dp(900)
	height := unit.Dp(500)

	window := new(app.Window)
	window.Option(
		app.Title("Seen - Noisy Sphere"),
		app.Size(width, height))

	context := gio.NewContext(window)

	scene, layer := Scene(context)
	widget := gio.Widget(context, func(w, h unit.Dp) {
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
			widget(gtx)
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
		err = doc.SaveToFile("noisysphere.svg")
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(0)
}

func Scene(context *gio.Context) (*seen.Scene, layer.Layer) {
	// Create a simple sphere
	shape := shape.Sphere(2)
	shape.SetScale(150, 150, 150)
	shape.Faces().SetColorFrom(color.NewDriftingSourceWith(color.Opacity(0.9)))

	// Create scene and add shape to group
	scene := seen.NewDefaultScene()
	scene.ShowBackfaces = true
	scene.Group.Add(shape)

	// Create a copy of the face points so we can manipulate them later
	points := []point.Points{}
	faces := shape.Faces()
	for i := range faces {
		points = append(points, slices.Clone(faces[i].Points))
	}

	// Create a layer that renders a scene by depth-sorting the polygons for the current eye
	layer := nsort.NewLayerForScene(scene)
	context.SetLayers(layer)

	// Create a 3D simplex noise generator
	noiser := noise.NewSimplex3D(rand.Float64())

	// Dynamically change the shape of the spheres
	context.Animate().OnBefore(func(t, dt time.Duration) {
		tms := float64(t.Milliseconds())
		dtms := float64(dt.Milliseconds())

		if true {
			faces := shape.Faces()
			for i := range faces {
				// Apply noise to sphere vertices
				for j := range faces[i].Points {
					p := points[i][j]
					n := noiser.Noise(p.X, p.Y, p.Z+tms*1e-4)
					faces[i].Points[j] = p.Times(1 + n/8)
				}
				// Since we're modifying the point directly, we need to mark the face dirty
				// to make sure the cache doesn't ignore the change
				faces[i].Dirty = true
			}
		}

		if false {
			shape.SetRotation(quaternion.RotZ(-dtms * 1e-4).RotX(dtms * 1e-4).Mul(shape.Rotation()))
		}
	}).Start()

	// Enable drag-to-rotate
	context.Drag(drag.Inertia(true)).On(func(e drag.Event) {
		shape.SetRotation(quaternion.RotY(e.Dx / 150).RotX(e.Dy / 150).Mul(shape.Rotation()))
		context.Render()
	})

	// Enable mouse-wheel zoom
	context.Zoom().On(func(e zoom.Event) {
		sx, sy, sz := shape.Scale()
		shape.SetScale(sx*e.Zoom, sy*e.Zoom, sz*e.Zoom)
	})

	return scene, layer
}
