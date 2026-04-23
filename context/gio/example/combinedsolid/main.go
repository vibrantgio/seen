package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/context/gio"
	"github.com/vibrantgio/seen/context/svg"
	"github.com/vibrantgio/seen/drag"
	"github.com/vibrantgio/seen/layer/bsort"
	"github.com/vibrantgio/seen/quaternion"
	"github.com/vibrantgio/seen/viewport"
	"github.com/vibrantgio/seen/zoom"

	. "github.com/vibrantgio/seen/solid"
)

const should_save_to_svg = false

func main() {
	go CombinedSolid()
	app.Main()
}

func CombinedSolid() {
	width := unit.Dp(900)
	height := unit.Dp(900)

	window := app.NewWindow(
		app.Title("Seen - Combined Solid"),
		app.Size(width, height))

	size := float64(height)
	if width < height {
		size = float64(width)
	}

	// Create a shape by subtracting a sphere from a cube

	const redhue = 0.0
	const greenhue = 0.333333
	const bluehue = 0.666666
	_, _, _ = redhue, greenhue, bluehue

	// a := solid.Cube()
	// b := solid.Sphere(solid.Radius(1.35), solid.Stacks(12))
	// c := solid.Cylinder(solid.Radius(0.7), solid.Start(-1, 0, 0), solid.End(1, 0, 0))
	// d := solid.Cylinder(solid.Radius(0.7), solid.Start(0, -1, 0), solid.End(0, 1, 0))
	// e := solid.Cylinder(solid.Radius(0.7), solid.Start(0, 0, -1), solid.End(0, 0, 1))

	// cylinder := Cylinder(Start(0, .7, 0), End(0, -.7, 0), Radius(.7))
	// core := Cylinder(Start(0, .8, 0), End(0, -.8, 0), Radius(.1))
	// cube := Cube(Center(-0.5, 0, 0), Size(0.5, 1.5, 0.7))
	// sphere := Sphere(Center(0.7, 0, 0), Radius(0.4))
	// solid := cylinder.Subtract(core).Subtract(cube).Subtract(sphere)

	a := Cube()
	b := Sphere(Radius(1.35), Stacks(12))
	c := Cylinder(Radius(0.7), Start(-1, 0, 0), End(1, 0, 0))
	d := Cylinder(Radius(0.7), Start(0, -1, 0), End(0, 1, 0))
	e := Cylinder(Radius(0.7), Start(0, 0, -1), End(0, 0, 1))
	csg := a.Intersect(b).Subtract(c.Union(d).Union(e))

	hue := greenhue

	fmt.Println("POLYCOUNT", len(csg))

	node := NewSolid("csg", csg)
	node.SetScale(size*0.2, size*0.2, size*0.2)
	node.Faces().SetColorFrom(color.NewDriftingSourceWith(color.Hue(hue), color.Lit(0.4), color.Drift(0.002)))

	// Create scene and add shape to group
	scene := seen.NewDefaultScene()
	scene.ShowBackfaces = true
	scene.FractionalPoints = true
	scene.Group.Add(node)

	// Create a layer that renders a scene by bsp-sorting the polygons
	layer := bsort.NewLayerForScene(scene)

	// Create a render context that hooks seen into the gio window
	context := gio.NewContext(window, layer)

	// Slowly rotate sphere
	context.Animate().OnBefore(func(t, dt time.Duration) {
		dtms := float64(dt.Milliseconds())
		r := quaternion.RotY(0.7 * dtms * 1e-4).RotX(dtms * 1e-4).Mul(node.Rotation())
		node.SetRotation(r)
	}).Start()

	// Enable dragger-to-rotate
	context.Drag(drag.Inertia(true)).On(func(e drag.Event) {
		r := quaternion.RotY(e.Dx / 150).RotX(e.Dy / 150).Mul(node.Rotation())
		node.SetRotation(r)
		context.Render()
	})

	// Enable mouse-wheel zoom
	context.Zoom().On(func(e zoom.Event) {
		sx, sy, sz := node.Scale()
		node.SetScale(sx*e.Zoom, sy*e.Zoom, sz*e.Zoom)
	})

	widget := gio.Widget(context, func(w, h unit.Dp) {
		width, height = w, h
		scene.Viewport = viewport.Center(0, 0, float64(w), float64(h))
	})

	ops := &op.Ops{}
	for event := range window.Events() {
		if frame, ok := event.(system.FrameEvent); ok {
			gtx := layout.NewContext(ops, frame)
			widget(gtx)
			frame.Frame(ops)
		}
	}

	// Save scene to svg file
	if should_save_to_svg {
		doc, err := svg.NewSVG("seen-svg", int(width), int(height))
		if err != nil {
			log.Fatal(err)
		}
		if context := svg.NewContext(doc.GetElementById("seen-svg"), layer); context != nil {
			context.Render()
		} else {
			log.Fatal("Render context is nil")
		}
		err = doc.SaveToFile("combinedsolid.svg")
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(0)
}
