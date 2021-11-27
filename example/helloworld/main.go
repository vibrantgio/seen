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
	"github.com/reactivego/seen/color"
	"github.com/reactivego/seen/quat"
	"github.com/reactivego/seen/render"
	"github.com/reactivego/seen/render/bsp"
	"github.com/reactivego/seen/render/gio"
	"github.com/reactivego/seen/render/zsort"
	"github.com/reactivego/seen/shape"
)

const should_use_bsp_sorter = true

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
	shape := shape.Sphere(2)
	shape.SetScale(HeightDp*0.4, HeightDp*0.4, HeightDp*0.4)
	shape.SetColorFrom(color.RandomSource2With(color.Opacity(0.9)))

	// Create scene and add shape to group
	scene := seen.DefaultScene()
	scene.ShowBackfaces = true
	scene.Group.Add(shape)

	// Create a layer that renders a scene by sorting the polygons
	var layer render.Layer
	if should_use_bsp_sorter {
		layer = bsp.LayerWith(scene)
	} else {
		layer = zsort.LayerWith(scene)
	}

	// Create a context that hooks seen into the gio window
	context := gio.ContextWith(window, layer)

	// Slowly rotate sphere
	animator := context.Animate()
	animator.OnBefore(func(t, dt time.Duration) {
		dtms := float64(dt.Milliseconds())
		r := quat.RotY(0.7 * dtms * 1e-4).RotX(dtms * 1e-4).Mul(shape.Rotation())
		shape.SetRotation(r)
	})
	animator.Start()

	// Enable drag-to-rotate
	drag := context.Drag(seen.Inertia(true))
	drag.On(func(e seen.DragEvent) {
		r := quat.RotY(e.Dx / 150).RotX(e.Dy / 150).Mul(shape.Rotation())
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
			ppd, dx, dy := frame.Metric.PxPerDp, float32(frame.Size.X), float32(frame.Size.Y)
			ops.Reset()
			scene.Viewport = seen.CenterViewport(0, 0, float64(dx/ppd), float64(dy/ppd))
			op.Affine(f32.NewAffine2D(ppd, 0, 0, 0, ppd, 0)).Add(ops)
			context.Draw(ops, frame.Queue)
			frame.Frame(ops)
		}
	}
	os.Exit(0)
}
