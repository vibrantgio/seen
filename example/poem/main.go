package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/reactivego/seen"
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
	go Poem()
	app.Main()
}

func Poem() {
	window := app.NewWindow(
		app.Title("Seen - Poem"),
		app.Size(unit.Dp(WidthDp), unit.Dp(HeightDp)),
		app.MinSize(unit.Dp(640), unit.Dp(480)),
	)

	// Draw text
	opts := map[string]string{
		"font-family": "Roboto",
		"font-weight": "normal", // normal | bold
		"font-size":   "14px",
		"anchor":      "middle", // start | middle | end
		"inline-size": "250px",
	}
	t := shape.Text(
		`Two roads diverged in a yellow wood,
And sorry I could not travel both
And be one traveler, long I stood
And looked down one as far as I could
To where it bent in the undergrowth;`, opts)
	t.SetShowBackfaces(true)
	t.SetTranslation(0, 0, 0)
	t.SetFill("#000000")

	// Create scene and add shape to group
	scene := seen.DefaultScene()
	scene.Group.Add(t)
	scene.Group.SetScale(2, 2, 2)

	// Create a layer that renders a scene by sorting the polygons
	var layer render.Layer
	if should_use_bsp_sorter {
		layer = bsp.LayerWith(scene)
	} else {
		layer = zsort.LayerWith(scene)
	}

	// Create a context that hooks seen into the gio window
	context := gio.ContextWith(window, layer)

	// Enable drag-to-rotate
	drag := context.Drag(seen.Inertia(true))
	drag.On(func(e seen.DragEvent) {
		r := quat.RotX(e.Dy / 150).Mul(t.Rotation()).RotY(e.Dx / 150)
		t.SetRotation(r)
		context.Render()
	})

	// Enable mouse-wheel zoom
	zoom := context.Zoom()
	zoom.On(func(e seen.ZoomEvent) {
		sx, sy, sz := t.Scale()
		t.SetScale(sx*e.Zoom, sy*e.Zoom, sz*e.Zoom)
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
