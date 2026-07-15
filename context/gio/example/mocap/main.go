// Command mocap renders a BVH motion-capture skeleton in a gio window,
// animating through the captured frames. It is a Go port of the original
// seen.js mocap demo (see main.coffee).
package main

import (
	_ "embed"
	"log"
	"os"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/bvh"
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/context/gio"
	"github.com/vibrantgio/seen/drag"
	"github.com/vibrantgio/seen/layer/backdrop"
	"github.com/vibrantgio/seen/layer/zsort"
	"github.com/vibrantgio/seen/mocap"
	"github.com/vibrantgio/seen/point"
	"github.com/vibrantgio/seen/quaternion"
	"github.com/vibrantgio/seen/shader"
	"github.com/vibrantgio/seen/shape"
	"github.com/vibrantgio/seen/viewport"
	"github.com/vibrantgio/seen/zoom"
)

//go:embed 05_11.bvh
var motion05_11 []byte

//go:embed 01_06.bvh
var motion01_06 []byte

const size = 900

func main() {
	go Mocap()
	app.Main()
}

func Mocap() {
	window := new(app.Window)
	window.Option(
		app.Title("Seen - Motion Capture"),
		app.Size(size, size))

	scene := Scene(gio.NewContext(window))

	ops := new(op.Ops)
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			gtx := app.NewContext(ops, e)
			scene(gtx)
			e.Frame(ops)
		}
	}
}

func Scene(context *gio.Context) layout.Widget {
	// Parse the embedded BVH file and build an animatable skeleton.
	result, err := bvh.Parse("05_11.bvh", motion05_11)
	if err != nil {
		log.Fatal(err)
	}
	model := mocap.New(result.(bvh.Hierarchy), shapeFactory)

	// Add the skeleton to a scene, lifting it so it sits in view.
	skeleton := seen.NewGroup(model.Group)
	skeleton.SetTranslation(0, -50, 0)

	scene := seen.NewDefaultScene()
	scene.Shader = shader.Phong
	scene.Group.Add(skeleton)

	// A dark backdrop behind the sorted skeleton polygons.
	background := backdrop.NewLayer(size, size, 0, 0, mustColor("#444444"))
	context.SetLayers(background, zsort.NewLayerForScene(scene))

	// Drag-to-rotate the whole skeleton.
	context.Drag(drag.Inertia(true)).On(func(e drag.Event) {
		r := model.Group.Rotation()
		r = quaternion.RotX(e.Dy / 150).Mul(r).RotY(e.Dx / 150)
		model.Group.SetRotation(r)
		context.Render()
	})

	// Mouse-wheel zoom.
	context.Zoom().On(func(e zoom.Event) {
		sx, sy, sz := scene.Group.Scale()
		scene.Group.SetScale(sx*e.Zoom, sy*e.Zoom, sz*e.Zoom)
	})

	// Play the captured frames in real time: the animation loop ticks slower
	// than the capture rate, so pick the frame by elapsed time rather than
	// stepping one frame per tick.
	var start time.Duration
	context.Animate().OnBefore(func(t, dt time.Duration) {
		if start == 0 {
			start = t
		}
		model.Apply(int((t - start) / model.FrameTime))
	}).Start()

	return gio.Widget(context, func(w, h unit.Dp) {
		background.Width, background.Height = float64(w), float64(h)
		scene.Viewport = viewport.Center(0, 0, float64(w), float64(h))
		// Dolly the camera in: the viewport puts the eye at distance h, so
		// translating the world +h/2 towards the eye halves that distance.
		scene.Camera.SetTranslation(0, 0, float64(h)/2)
	})
}

// shapeFactory sizes and colours each bone by joint name, mirroring the
// original seen.js demo.
func shapeFactory(joint bvh.Joint, endpoint point.Point) seen.Object {
	if endpoint.Length() < 1e-9 {
		return nil
	}
	id := strings.ToLower(joint.Id)

	radius := 2.0
	switch {
	case containsAny(id, "thumb", "index", "mid", "ring", "pinky"):
		radius = 0.4
	case strings.Contains(id, "hand"):
		radius = 1
	case containsAny(id, "abdomen", "chest"):
		radius = 6
	case strings.Contains(id, "hip"):
		radius = 4
	}

	fill := "#FFD2A6"
	switch {
	case containsAny(id, "foot", "hip", "abdomen", "chest"):
		fill = "#FF88FF"
	case strings.Contains(id, "shin"):
		fill = "#FFFFFF"
	}

	pipe := shape.Pipe(point.Pt(0, 0, 0), endpoint, shape.Radius(radius), shape.Segments(8))
	pipe.Faces().SetFill(mustColor(fill))
	return pipe
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func mustColor(s string) color.Color {
	c, err := color.ColorWithString(s)
	if err != nil {
		log.Fatal(err)
	}
	return c
}
