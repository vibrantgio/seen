package main

import (
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/context/gio"
	"github.com/vibrantgio/seen/drag"
	"github.com/vibrantgio/seen/layer/bsort"
	"github.com/vibrantgio/seen/shape"
	"github.com/vibrantgio/seen/viewport"
	"github.com/vibrantgio/seen/zoom"
)

func main() {
	go HelloWorld()
	app.Main()
}

func HelloWorld() {
	window := new(app.Window)
	window.Option(
		app.Title("Seen - Hello, World!"),
		app.Size(1000, 1000))

	scene := Scene(gio.NewContext(window), 1000, 1000)

	ops := new(op.Ops)
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			gtx := app.NewContext(ops, e)
			scene(gtx)
			e.Frame(ops)
		}
	}
}

func Scene(context *gio.Context, width, height unit.Dp) layout.Widget {
	size := float64(height)
	if width < height {
		size = float64(width)
	}

	// Create sphere shape with randomly colored faces
	shape := shape.Sphere(2)
	size *= 0.45
	shape.SetScale(size, size, size)
	shape.Faces().SetColorFrom(color.NewDriftingSourceWith(color.Opacity(0.9)))

	// Create a scene and add the shape to group
	scene := seen.NewDefaultScene()
	scene.ShowBackfaces = true
	scene.Group.Add(shape)

	// Create a layer that renders a scene by sorting the polygons
	context.SetLayers(bsort.NewLayerForScene(scene))

	// Slowly rotate sphere
	context.Animate().OnBefore(func(t, dt time.Duration) {
		dtms := float64(dt.Milliseconds())
		shape.RotX(dtms * 1e-4).RotY(0.7 * dtms * 1e-4)
	}).Start()

	// Enable drag-to-rotate
	context.Drag(drag.Inertia(true)).On(func(e drag.Event) {
		shape.RotX(e.Dy / 150).RotY(e.Dx / 150)
		context.Render()
	})

	// Enable mouse-wheel zoom
	context.Zoom().On(func(e zoom.Event) {
		sx, sy, sz := shape.Scale()
		shape.SetScale(sx*e.Zoom, sy*e.Zoom, sz*e.Zoom)
	})

	return gio.Widget(context, func(w, h unit.Dp) {
		scene.Viewport = viewport.Center(0, 0, float64(w), float64(h))
	})
}
