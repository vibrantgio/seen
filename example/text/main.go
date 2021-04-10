package main

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/reactivego/seen"
	"github.com/reactivego/seen/quat"
	"github.com/reactivego/seen/render"
	"github.com/reactivego/seen/render/gio"
	"github.com/reactivego/seen/shapes"
)

const WidthDp = 900
const HeightDp = 500

func main() {
	go Text()
	app.Main()
}

func Text() {
	window := app.NewWindow(
		app.Title("Seen - Text"),
		app.Size(unit.Dp(WidthDp), unit.Dp(HeightDp)),
		app.MinSize(unit.Dp(640), unit.Dp(480)),
	)

	// Generate some random data points
	data := make([]float64, 0, 10)
	for i := 0; i < 10; i++ {
		data = append(data, rand.Float64()*80.0+20.0)
	}

	scene := seen.DefaultScene()
	model := seen.EmptyModel()
	scene.Model.Add(model)

	// Draw bars for data
	for i, d := range data {
		uc := shapes.UnitCube()
		uc.SetFill("#0088FF")
		uc.SetScale(20, d, 20)
		uc.SetTranslation(float64(i*30)-160, -50, 0)
		model.Add(uc)
	}

	// Draw text above bars
	for i, d := range data {
		opts := map[string]string{
			"font-family": "Roboto",
			"font-weight": "normal", // normal | bold
			"font-size":   "10px",
			"anchor":      "middle",
			"textLength":  "200px",
		}
		t := shapes.Text(strconv.FormatFloat(d, 'f', 1, 64), opts)
		t.SetShowBackfaces(true)
		t.SetTranslation(float64(i)*30+10-160, d+10-50, 10)
		t.SetFill("#000000")
		model.Add(t)
	}

	model.SetScale(2, 2, 2)

	// Create scene and add shape to model
	scene.Viewport = seen.CenterViewport(0, 0, WidthDp, HeightDp)

	// Create a render layer and render context
	layer := render.SceneLayerWith(&scene)
	layer.FractionalPoints = true
	context := gio.MakeContext(window, layer)

	// Slowly rotate the bar chart
	animator := context.Animate()
	animator.OnBefore(func(t, dt time.Duration) {
		dtms := float64(dt.Milliseconds())
		model.SetRotation(model.Rotation().RotY(0.7 * dtms * 1e-4))
	})
	animator.Start()

	// Enable drag-to-rotate
	drag := context.Drag(seen.Inertia(true))
	drag.On(func(e seen.DragEvent) {
		r := quat.RotX(e.Dy / 150).Mul(model.Rotation()).RotY(e.Dx / 150)
		model.SetRotation(r)
		context.Render()
	})

	// Enable mouse-wheel zoom
	zoom := context.Zoom()
	zoom.On(func(e seen.ZoomEvent) {
		sx, sy, sz := model.Scale()
		model.SetScale(sx*e.Zoom, sy*e.Zoom, sz*e.Zoom)
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
