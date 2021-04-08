package main

import (
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/reactivego/seen"
	"github.com/reactivego/seen/colors"
	"github.com/reactivego/seen/render"
	"github.com/reactivego/seen/render/gio"
	"github.com/reactivego/seen/shapes"
	"github.com/reactivego/seen/transform"
)

const WidthDp = 900
const HeightDp = 500

func main() {
	go HelloWorld()
	app.Main()
}

func HelloWorld() {
	window := app.NewWindow(
		app.Title("Seen - Hello, World!"),
		app.Size(unit.Dp(WidthDp), unit.Dp(HeightDp)),
		app.MinSize(unit.Dp(640), unit.Dp(480)),
	)

	// Create sphere shape with randomly colored surfaces
	shape := shapes.MakeSphere(2)
	shape.SetScale(HeightDp*0.4, HeightDp*0.4, HeightDp*0.4)
	shape.ColorSurfaces(colors.MakeRandomSource2())

	// Create scene and add shape to model
	scene := seen.DefaultScene()
	scene.Model.Add(shape)
	scene.Viewport = seen.CenterViewport(0, 0, WidthDp, HeightDp)

	// Create a render layer and render context
	layer := render.SceneLayerWith(&scene)
	context := gio.MakeContext(window, layer)

	// Slowly rotate sphere
	animator := context.Animate()
	animator.OnBefore(func(t, dt time.Duration) {
		dtms := float64(dt.Milliseconds())
		r := transform.QuatRotY(0.7 * dtms * 1e-4).MulRotX(dtms * 1e-4).Mul(shape.Rotation())
		shape.SetRotation(r)
	})
	animator.Start()

	// Enable drag-to-rotate
	drag := context.Drag(seen.Inertia(true))
	drag.On(func(e seen.DragEvent) {
		r := transform.QuatRotY(e.Dx / 150).MulRotX(e.Dy / 150).Mul(shape.Rotation())
		shape.SetRotation(r)
		context.Render()
	})

	// Enable mouse-wheel zoom
	zoom := context.Zoom()
	zoom.On(func(e seen.ZoomEvent) {
		sx, sy, sz := shape.Scale()
		shape.SetScale(sx*e.Zoom, sy*e.Zoom, sz*e.Zoom)
	})

	ops := &op.Ops{}
	for event := range window.Events() {
		if frame, ok := event.(system.FrameEvent); ok {
			ops.Reset()
			ppd := frame.Metric.PxPerDp
			op.Affine(f32.NewAffine2D(ppd, 0, 0, 0, ppd, 0)).Add(ops)
			context.Draw(ops, frame.Queue)
			frame.Frame(ops)
		}
	}
	os.Exit(0)
}
