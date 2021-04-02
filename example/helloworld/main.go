package main

import (
	"os"

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
	context := Setup(window)
	ops := &op.Ops{}
	for event := range window.Events() {
		if frame, ok := event.(system.FrameEvent); ok {
			ops.Reset()
			ppd := frame.Metric.PxPerDp
			op.Affine(f32.NewAffine2D(ppd, 0, 0, 0, ppd, 0)).Add(ops)
			context.HandleEvents(frame.Queue, ops)
			frame.Frame(ops)
		}
	}
	os.Exit(0)
}

func Setup(window *app.Window) *gio.Context {
	// Create sphere shape with randomly colored surfaces
	shape := shapes.MakeSphere(2)
	shape.SetScale(HeightDp*0.4, HeightDp*0.4, HeightDp*0.4)
	shape.ColorSurfaces(colors.MakeRandomSource2())

	// Create scene and add shape to model
	scene := seen.MakeScene()
	scene.Model = seen.MakeDefaultModel()
	scene.Model.Add(shape)
	scene.Viewport = seen.MakeCenterViewport(0, 0, WidthDp, HeightDp)

	// Create a render layer and render context
	layer := render.MakeSceneLayer(scene)
	context := gio.MakeContext(window, layer)

	// Slowly rotate sphere
	animator := context.Animate()
	animator.OnBefore(func(t, dt float64) {
		ryrx := transform.MakeQuatRotY(0.7 * dt * 1e-4).MulRotX(dt * 1e-4)
		shape.SetRotation(ryrx.Mul(shape.Rotation()))
	})
	animator.Start()

	// Enable drag-to-rotate
	drag := context.Drag(seen.Inertia(true))
	drag.On(func(e seen.DragEvent) {
		dx, dy := e.OffsetRelativeX/150, e.OffsetRelativeY/150
		ryrx := transform.MakeQuatRotY(dx).MulRotX(dy)
		shape.SetRotation(ryrx.Mul(shape.Rotation()))
		context.Render()
	})

	return context
}
