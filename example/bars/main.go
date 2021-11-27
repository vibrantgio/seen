package main

import (
	"log"
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
	"github.com/reactivego/seen/document"
	"github.com/reactivego/seen/quat"
	"github.com/reactivego/seen/render"
	"github.com/reactivego/seen/render/bsp"
	"github.com/reactivego/seen/render/gio"
	"github.com/reactivego/seen/render/svg"
	"github.com/reactivego/seen/render/zsort"
	"github.com/reactivego/seen/shape"
)

const should_use_bsp_sorter = true
const should_save_to_svg = false

const WidthDp = 900
const HeightDp = 500

func main() {
	go Bars()
	app.Main()
}

func Bars() {
	window := app.NewWindow(
		app.Title("Seen - Bars"),
		app.Size(unit.Dp(WidthDp), unit.Dp(HeightDp)),
		app.MinSize(unit.Dp(640), unit.Dp(480)),
	)

	// Generate some random data points
	data := make([]float64, 0, 10)
	for i := 0; i < 10; i++ {
		data = append(data, rand.Float64()*80.0+20.0)
	}

	scene := seen.DefaultScene()
	group := seen.EmptyGroup()
	scene.Group.Add(group)

	// Draw bars for data
	for i, d := range data {
		uc := shape.UnitCube()
		uc.SetFill("#0088FF")
		uc.SetScale(20, d, 20)
		uc.SetTranslation(float64(i*30)-160, -50, 0)
		group.Add(uc)
	}

	// Draw text above bars
	for i, d := range data {
		opts := map[string]string{
			"font-family": "Roboto",
			"font-weight": "normal", // normal | bold
			"font-size":   "10px",
			"anchor":      "middle",
			"inline-size": "200px",
		}
		t := shape.Text(strconv.FormatFloat(d, 'f', 1, 64), opts)
		t.SetShowBackfaces(true)
		t.SetTranslation(float64(i)*30+10-160, d+10-50, 10)
		t.SetFill("#000000")
		group.Add(t)
	}

	group.SetScale(2, 2, 2)

	// Create a layer that renders a scene by sorting the polygons
	var layer render.Layer
	if should_use_bsp_sorter {
		layer = bsp.LayerWith(scene)
	} else {
		layer = zsort.LayerWith(scene)
	}

	// Create a context that hooks seen into the gio window
	context := gio.ContextWith(window, layer)

	// Slowly rotate the bar chart
	animator := context.Animate()
	animator.OnBefore(func(t, dt time.Duration) {
		dtms := float64(dt.Milliseconds())
		group.SetRotation(group.Rotation().RotY(0.7 * dtms * 1e-4))
	})
	animator.Start()

	// Enable drag-to-rotate
	drag := context.Drag(seen.Inertia(true))
	drag.On(func(e seen.DragEvent) {
		r := quat.RotX(e.Dy / 150).Mul(group.Rotation()).RotY(e.Dx / 150)
		group.SetRotation(r)
		context.Render()
	})

	// Enable mouse-wheel zoom
	zoom := context.Zoom()
	zoom.On(func(e seen.ZoomEvent) {
		sx, sy, sz := group.Scale()
		group.SetScale(sx*e.Zoom, sy*e.Zoom, sz*e.Zoom)
	})

	ops := &op.Ops{}
	ppd, dx, dy := float32(1.0), float32(WidthDp), float32(HeightDp)
	for event := range window.Events() {
		if frame, ok := event.(system.FrameEvent); ok {
			ppd, dx, dy = frame.Metric.PxPerDp, float32(frame.Size.X), float32(frame.Size.Y)
			ops.Reset()
			scene.Viewport = seen.CenterViewport(0, 0, float64(dx/ppd), float64(dy/ppd))
			op.Affine(f32.NewAffine2D(ppd, 0, 0, 0, ppd, 0)).Add(ops)
			context.Draw(ops, frame.Queue)
			frame.Frame(ops)
		}
	}

	// Save scene to svg file
	if should_save_to_svg {
		svgdoc, err := document.MakeSVG("seen-svg", int(dx/ppd), int(dy/ppd))
		if err != nil {
			log.Fatal(err)
		}
		if context := svg.ContextWith(svgdoc.GetElementById("seen-svg"), layer); context != nil {
			context.Render()
		} else {
			log.Fatal("Render context is nil")
		}
		err = svgdoc.SaveToFile("bars.svg")
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(0)
}
