// Command solids showcases the three CSG boolean operations from the csg
// package README (now seen/solid): the union, difference and intersection of
// the same cube and sphere, side by side —
//
//	cube := Cube()
//	sphere := Sphere(Radius(1.3))
//	cube.Union(sphere)     // left,   red
//	cube.Subtract(sphere)  // middle, green
//	cube.Intersect(sphere) // right,  blue
//
// All three solids rotate in lockstep so the differences between the
// operations stay comparable. Drag to rotate them yourself, scroll to zoom.
package main

import (
	"fmt"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/context/gio"
	"github.com/vibrantgio/seen/drag"
	"github.com/vibrantgio/seen/layer/nsort"
	"github.com/vibrantgio/seen/quaternion"
	"github.com/vibrantgio/seen/viewport"
	"github.com/vibrantgio/seen/zoom"

	. "github.com/vibrantgio/seen/solid"
)

func main() {
	go Solids()
	app.Main()
}

func Solids() {
	width := unit.Dp(1200)
	height := unit.Dp(500)

	window := new(app.Window)
	window.Option(
		app.Title("Seen - CSG Solids: union · subtract · intersect"),
		app.Size(width, height))

	size := float64(height)

	// The README's operands: a 2x2x2 cube and a sphere that pokes through
	// its faces. All three operations combine the SAME two solids.
	cube := Cube()
	sphere := Sphere(Radius(1.3))

	const (
		redhue   = 0.0
		greenhue = 0.333333
		bluehue  = 0.666666
	)
	ops := []struct {
		name string
		csg  CSG
		hue  float64
		dx   float64
	}{
		{"cube ∪ sphere", cube.Union(sphere), redhue, -0.62},
		{"cube − sphere", cube.Subtract(sphere), greenhue, 0},
		{"cube ∩ sphere", cube.Intersect(sphere), bluehue, +0.62},
	}

	scene := seen.NewDefaultScene()
	scene.ShowBackfaces = true

	nodes := make([]seen.Object, len(ops))
	for i, op := range ops {
		fmt.Printf("%-14s %4d polygons\n", op.name, len(op.csg))
		node := NewSolid(op.name, op.csg)
		node.SetScale(size*0.16, size*0.16, size*0.16)
		node.SetTranslation(op.dx*size, 0, 0)
		node.Faces().SetColorFrom(color.NewDriftingSourceWith(
			color.Hue(op.hue), color.Lit(0.4), color.Drift(0.002)))
		scene.Group.Add(node)
		nodes[i] = node
	}

	// nsort: the solids' model matrices change every frame (rotation), and
	// concave CSG output orders exactly under the per-eye depth sort without
	// the per-frame tree rebuild and split seams the BSP layer would cost.
	layer := nsort.NewLayerForScene(scene)

	// Create a render context that hooks seen into the gio window
	context := gio.NewContext(window, layer)

	// Rotate all three solids in lockstep so the ops stay comparable.
	rotate := func(dy, dx float64) {
		for _, node := range nodes {
			node.SetRotation(quaternion.RotY(dy).RotX(dx).Mul(node.Rotation()))
		}
	}

	context.Animate().OnBefore(func(t, dt time.Duration) {
		dtms := float64(dt.Milliseconds())
		rotate(0.7*dtms*1e-4, dtms*1e-4)
	}).Start()

	// Enable drag-to-rotate
	context.Drag(drag.Inertia(true)).On(func(e drag.Event) {
		rotate(e.Dx/150, e.Dy/150)
		context.Render()
	})

	// Enable mouse-wheel zoom
	context.Zoom().On(func(e zoom.Event) {
		for _, node := range nodes {
			sx, sy, sz := node.Scale()
			node.SetScale(sx*e.Zoom, sy*e.Zoom, sz*e.Zoom)
		}
	})

	widget := gio.Widget(context, func(w, h unit.Dp) {
		scene.Viewport = viewport.Center(0, 0, float64(w), float64(h))
	})

	ops2 := &op.Ops{}
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			gtx := app.NewContext(ops2, e)
			widget(gtx)
			e.Frame(ops2)
		}
	}
}
