package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/context/gio"
	"github.com/vibrantgio/seen/context/svg"
	"github.com/vibrantgio/seen/drag"
	"github.com/vibrantgio/seen/layer"
	"github.com/vibrantgio/seen/layer/bsort"
	"github.com/vibrantgio/seen/layer/zsort"
	"github.com/vibrantgio/seen/quaternion"
	"github.com/vibrantgio/seen/shape"
	"github.com/vibrantgio/seen/viewport"
	"github.com/vibrantgio/seen/zoom"
)

const should_use_bsp_sorter = true

const should_save_to_svg = false

func main() {
	go Text()
	app.Main()
}

func Text() {
	window := new(app.Window)
	window.Option(
		app.Title("Seen - Text"),
		app.Size(900, 500),
		app.MinSize(450, 250))

	// Generate some random data points
	data := make([]float64, 0, 10)
	for i := 0; i < 10; i++ {
		data = append(data, rand.Float64()*80.0+20.0)
	}

	scene := seen.NewDefaultScene()
	group := seen.NewGroup()
	scene.Group.Add(group)

	// Draw bars for data
	for i, d := range data {
		uc := shape.UnitCube()
		uc.Faces().SetFill("#0088FF")
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
		t.Faces().SetShowBackfaces(true)
		t.SetTranslation(float64(i)*30+10-160, d+10-50, 10)
		t.Faces().SetFill("#000000")
		group.Add(t)
	}

	group.SetScale(2, 2, 2)

	// Create a layer that renders a scene by sorting the polygons
	var layer layer.Layer
	if should_use_bsp_sorter {
		layer = bsort.NewLayerForScene(scene)
	} else {
		layer = zsort.NewLayerForScene(scene)
	}

	// Create a context that hooks seen into the gio window
	context := gio.NewContext(window, layer)

	// Slowly rotate the bar chart
	context.Animate().OnBefore(func(t, dt time.Duration) {
		dtms := float64(dt.Milliseconds())
		group.SetRotation(group.Rotation().RotY(0.7 * dtms * 1e-4))
	}).Start()

	// Enable drag-to-rotate
	context.Drag(drag.Inertia(true)).On(func(e drag.Event) {
		r := quaternion.RotX(e.Dy / 150).Mul(group.Rotation()).RotY(e.Dx / 150)
		group.SetRotation(r)
		context.Render()
	})

	// Enable mouse-wheel zoom
	context.Zoom().On(func(e zoom.Event) {
		sx, sy, sz := group.Scale()
		group.SetScale(sx*e.Zoom, sy*e.Zoom, sz*e.Zoom)
	})

	view := gio.Widget(context, func(w, h unit.Dp) {
		scene.Viewport = viewport.Center(0, 0, float64(w), float64(h))
	})

	ops := &op.Ops{}
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			gtx := app.NewContext(ops, e)
			view(gtx)
			e.Frame(ops)
		}
	}

	// Save scene to svg file
	if should_save_to_svg {
		const dx, dy, ppd = 900.0, 500.0, 1.0
		svgdoc, err := svg.NewSVG("seen-svg", int(dx/ppd), int(dy/ppd))
		if err != nil {
			log.Fatal(err)
		}
		if context := svg.NewContext(svgdoc.GetElementById("seen-svg"), layer); context != nil {
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
