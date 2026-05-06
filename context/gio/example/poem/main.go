package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/context/gio"
	"github.com/vibrantgio/seen/drag"
	"github.com/vibrantgio/seen/layer/bsort"
	"github.com/vibrantgio/seen/quaternion"
	"github.com/vibrantgio/seen/shape"
	"github.com/vibrantgio/seen/viewport"
	"github.com/vibrantgio/seen/zoom"
)

func main() {
	go Poem()
	app.Main()
}

func Poem() {
	window := new(app.Window)
	window.Option(
		app.Title("Seen - Poem"),
		app.Size(1600, 900),
		app.MinSize(640, 480))

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
	t.Faces().SetShowBackfaces(true)
	t.SetTranslation(0, 0, 0)
	t.Faces().SetFill("#000000")

	// Create scene and add shape to its group
	scene := seen.NewDefaultScene()
	scene.Group.Add(t)
	scene.Group.SetScale(2, 2, 2)

	// Create a layer that renders a scene by sorting the polygons
	foreground := bsort.NewLayerForScene(scene)

	// Create a context that hooks seen into the gio window
	context := gio.NewContext(window, foreground)

	// Enable drag-to-rotate
	context.Drag(drag.Inertia(true)).On(func(e drag.Event) {
		t.SetRotation(quaternion.RotX(e.Dy / 150).Mul(t.Rotation()).RotY(e.Dx / 150))
		context.Render()
	})

	// Enable mouse-wheel zoom
	context.Zoom().On(func(e zoom.Event) {
		sx, sy, sz := t.Scale()
		t.SetScale(sx*e.Zoom, sy*e.Zoom, sz*e.Zoom)
	})

	widget := gio.Widget(context, func(w, h unit.Dp) {
		scene.Viewport = viewport.Center(0, 0, float64(w), float64(h))
	})

	ops := new(op.Ops)
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			gtx := app.NewContext(ops, e)
			widget(gtx)
			e.Frame(ops)
		}
	}
}
