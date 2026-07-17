package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/context/gio"
	"github.com/vibrantgio/seen/drag"
	"github.com/vibrantgio/seen/layer/nsort"
	"github.com/vibrantgio/seen/point"
	"github.com/vibrantgio/seen/quaternion"
	"github.com/vibrantgio/seen/shape"
	"github.com/vibrantgio/seen/viewport"
	"github.com/vibrantgio/seen/zoom"
)

func main() {
	go Rectangle()
	app.Main()
}

func Rectangle() {
	window := new(app.Window)
	window.Option(
		app.Title("Seen - Rectangle"),
		app.Size(1000, 1000))

	scene := Scene(gio.NewContext(window), 1000)

	ops := &op.Ops{}
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

func Scene(context *gio.Context, size float64) layout.Widget {
	// Create a rectangle that covers the entire viewport exactly.
	// The viewport shows a world where (0,0) is centered in the
	// middle of the screen with the positive x-axis pointing to
	// the right and the positive y-axis pointing up. The positive
	// z-axis is pointing out of the screen towards the viewer.
	// This makes it a right-handed coordinate system.
	min, max := point.P(-0.5*size, -0.5*size, 0.0), point.P(0.5*size, 0.5*size, -20.0)
	shape := shape.Rectangle(min, max)
	shape.Faces().SetColorFrom(color.NewDriftingSourceWith(color.Opacity(0.9)))

	// Create a scene and add the shape to group
	scene := seen.NewDefaultScene()
	scene.ShowBackfaces = true
	scene.Group.Add(shape)

	// Create a layer that renders a scene by sorting the polygons
	context.SetLayers(nsort.NewLayerForScene(scene))

	// Enable drag-to-rotate
	context.Drag(drag.Inertia(true)).On(func(e drag.Event) {
		if true {
			shape.SetRotation(quaternion.RotY(e.Dx / 150).RotX(e.Dy / 150).Mul(shape.Rotation()))
		} else {
			shape.RotX(e.Dy / 150).RotY(e.Dx / 150)
		}
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
