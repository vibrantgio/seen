// Command multipleangles shows the same scene from five angles at once: one
// shared model — the Stanford bunny from bunny-low.obj — rendered through
// five scenes that differ only in camera and viewport. The main 600x300 view
// on top looks at the bunny straight on through a camera scaled 2x; beneath
// it four 150x150 minis show the default camera plus the camera pitched a
// quarter turn up and down and yawed a quarter turn sideways. Drag to rotate
// the shared model: all five views update in lockstep. It is a Go port of
// the original seen.js multipleangles demo (see main.coffee).
package main

import (
	"bytes"
	_ "embed"
	"math"
	"math/rand"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/context/gio"
	"github.com/vibrantgio/seen/drag"
	"github.com/vibrantgio/seen/face"
	"github.com/vibrantgio/seen/layer"
	"github.com/vibrantgio/seen/layer/nsort"
	"github.com/vibrantgio/seen/light"
	"github.com/vibrantgio/seen/obj"
	"github.com/vibrantgio/seen/point"
	"github.com/vibrantgio/seen/quaternion"
	"github.com/vibrantgio/seen/shape"
	"github.com/vibrantgio/seen/viewport"
)

//go:embed bunny-low.obj
var bunnyObj []byte

func main() {
	go MultipleAngles()
	app.Main()
}

func MultipleAngles() {
	window := new(app.Window)
	window.Option(
		app.Title("Seen - Same Scene, Multiple Angles"),
		app.Size(600, 450))

	// One model shared between all the scenes.
	bunny := Bunny()
	model := seen.NewGroupWithLights(light.DefaultLights()...)
	model.Add(bunny)

	// Five scenes render the shared model, each through its own camera and
	// viewport; one nsort layer per scene composites them onto one canvas.
	scenes := make([]*seen.Scene, 5)
	layers := make([]layer.Layer, len(scenes))
	for i := range scenes {
		scenes[i] = seen.NewScene()
		scenes[i].Group = model
		layers[i] = nsort.NewLayerForScene(scenes[i])
	}

	// Zoom the main camera in 2x and angle the last three mini cameras to
	// look from above, below and the side.
	scenes[0].Camera.SetScale(2, 2, 2)
	scenes[2].Camera.RotX(math.Pi / 2)
	scenes[3].Camera.RotX(-math.Pi / 2)
	scenes[4].Camera.RotY(-math.Pi / 2)

	context := gio.NewContext(window, layers...)

	// Drag-to-rotate the shared model; every view follows.
	context.Drag(drag.Inertia(true)).On(func(e drag.Event) {
		rotation := quaternion.RotY(e.Dx / 150).RotX(e.Dy / 150)
		bunny.SetRotation(rotation.Mul(bunny.Rotation()))
		context.Render()
	})

	widget := gio.Widget(context, func(w, h unit.Dp) {
		// The main view fills the window above a row of four minis.
		mini := float64(w) / 4
		for i, scene := range scenes {
			ox, oy, vw, vh := 0.0, 0.0, float64(w), float64(h)-mini
			if i > 0 {
				ox, oy, vw, vh = float64(i-1)*mini, float64(h)-mini, mini, mini
			}
			// Center(ox, oy, …) parks the eye above the world point
			// (ox, oy), so steering each camera to that same point
			// puts the model — at the world origin — in the middle
			// of the region whose top-left corner is at (ox, oy).
			scene.Viewport = viewport.Center(ox, oy, vw, vh)
			scene.Camera.SetTranslation(ox, oy, 0)
		}
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

// Bunny builds the bunny model from the embedded OBJ file, posed and colored
// like the original demo: scaled 8x, lowered 30 units, pitched, yawed and
// rolled a quarter turn onto its feet, with a random walk of green hues
// across the faces.
func Bunny() *seen.Group {
	var faces face.Faces
	for pts := range obj.Parse(bytes.NewReader(bunnyObj)) {
		points := make(point.Points, 0, len(pts))
		for _, p := range pts {
			points = append(points, point.Pt(p[0], p[1], p[2]))
		}
		faces = append(faces, face.FaceWith(points))
	}
	faces.SetShowBackfaces(true) // the coffee loads the obj with cullBackfaces=false
	faces.SetColorFrom(&hueWalk{rng: rand.New(rand.NewSource(6)), hue: 0.35, step: 0.02})

	// The coffee scales and translates the shape before rotating it, an
	// order Transform's T*R*S cannot express in one node: bake the scale
	// and translation into the shape and the rotations into a wrapping
	// group.
	node := shape.NewShapeWithFaces("bunny", faces)
	node.SetScale(8, 8, 8)
	node.SetTranslation(0, -30, 0)

	group := seen.NewGroup(node)
	group.RotX(math.Pi / 4)
	group.RotY(-math.Pi / 4)
	group.RotZ(-math.Pi / 4)
	return group
}

// hueWalk mirrors seen.Colors.randomSurfaces2: each face's hue takes a small
// random step away from the previous face's, here from a fixed seed so the
// coloring is the same on every run, and with a smaller step than the
// original's 0.1 so 400 faces of bunny stay within one family of hues.
type hueWalk struct {
	rng  *rand.Rand
	hue  float64
	step float64
}

var _ color.Source = (*hueWalk)(nil)

func (walk *hueWalk) Read() color.Color {
	walk.hue += (walk.rng.Float64() - 0.5) * walk.step
	switch {
	case walk.hue < 0:
		walk.hue++
	case walk.hue > 1:
		walk.hue--
	}
	return color.ColorHSL(walk.hue, 0.5, 0.4, 1)
}
