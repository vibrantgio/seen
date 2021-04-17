package main

import (
	"image"
	"image/png"
	"log"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/gpu/headless"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"
	"golang.org/x/image/draw"

	"github.com/reactivego/seen"
	"github.com/reactivego/seen/bsp"
	"github.com/reactivego/seen/color"
	"github.com/reactivego/seen/document"
	"github.com/reactivego/seen/quat"
	"github.com/reactivego/seen/render"
	"github.com/reactivego/seen/render/gio"
	"github.com/reactivego/seen/render/svg"
)

const WidthDp = 1024
const HeightDp = 1024

func main() {
	go Present()
	app.Main()
}

func Present() {
	window := app.NewWindow(
		app.Title("Seen - Present"),
		app.Size(unit.Dp(WidthDp), unit.Dp(HeightDp)),
		app.MinSize(unit.Dp(WidthDp), unit.Dp(HeightDp)),
	)

	// Colors to use
	lightblue := color.Color{R: 1.0 / 255.0, G: 202.0 / 255.0, B: 252.0 / 255.0, A: 1.0}
	darkblue := color.Color{R: 0.0, G: 130.0 / 255.0, B: 193.0 / 255.0, A: 1.0}
	hardblue, _ := color.ColorWithString("#cceeff")
	lightgrey, _ := color.ColorWithString("#eeeeee")
	lightorange := color.Color{R: 247.0 / 255.0, G: 148.0 / 255.0, B: 29.0 / 255.0, A: 1.0}
	darkorange := color.Color{R: 224.0 / 255.0, G: 134.0 / 255.0, B: 26.0 / 255.0, A: 1.0}

	_, _, _, _, _, _ = lightblue, darkblue, hardblue, lightgrey, lightorange, darkorange

	backdropfill := color.White
	curtainfill := color.White
	boxfill := lightblue
	ribbonfill := darkorange

	// Box with a lid
	box := Box()
	box.SetFill(boxfill)
	box.SetScale(49, 49, 49)
	box.SetTranslation(0, 0, 0)

	// Ribbon around box
	ribbon := Ribbon()
	ribbon.SetFill(ribbonfill)
	ribbon.SetScale(50, 50, 50)
	ribbon.SetTranslation(0, 0, 0)

	// Present is a box with a lid and a ribbon around it
	present := seen.GroupWith(box, ribbon)
	present.SetRotation(quat.RotX(0.25 * math.Pi).RotY(-0.25 * math.Pi))

	// Create scene and add shape to group
	scene := seen.DefaultScene()
	scene.FractionalPoints = true
	scene.Group.Add(present)
	scene.Group.SetScale(2, 2, 2)
	scene.Viewport = seen.CenterViewport(0, 0, WidthDp, HeightDp)

	// Create separate layers for the stage.
	backdrop := render.FillLayerWith(WidthDp, HeightDp, 0, 0, backdropfill)
	curtain := render.FillLayerWith(WidthDp, HeightDp/2, 0, 0, curtainfill)
	foreground := bsp.SceneLayerWith(scene)

	context := gio.ContextWith(window, backdrop, curtain, foreground)

	// Enable drag-to-rotate
	drag := context.Drag(seen.Inertia(true))
	drag.On(func(e seen.DragEvent) {
		r := present.Rotation()
		r = quat.RotX(e.Dy / 150).Mul(r).RotY(e.Dx / 150)
		present.SetRotation(r)
		context.Render()
	})

	// Enable mouse-wheel zoom
	zoom := context.Zoom()
	zoom.On(func(e seen.ZoomEvent) {
		sx, sy, sz := present.Scale()
		present.SetScale(sx*e.Zoom, sy*e.Zoom, sz*e.Zoom)
	})

	ops := &op.Ops{}
	ppd := float32(1.0)
	for event := range window.Events() {
		if frame, ok := event.(system.FrameEvent); ok {
			ops.Reset()
			ppd = frame.Metric.PxPerDp
			op.Affine(f32.NewAffine2D(ppd, 0, 0, 0, ppd, 0)).Add(ops)
			context.Draw(ops, frame.Queue)
			frame.Frame(ops)
		}
	}

	// Save scene to png file
	if window, err := headless.NewWindow(int(ppd*WidthDp), int(ppd*HeightDp)); err == nil {
		window.Frame(ops)
		if src, err := window.Screenshot(); err == nil {
			dst := src
			// Scale image when ppdp (pixels per device pixel) is not 1.0
			if ppd != 1.0 {
				sb := src.Bounds()
				w, h := int(float32(sb.Dx())/ppd), int(float32(sb.Dy())/ppd)
				dst := image.NewRGBA(image.Rect(sb.Min.X, sb.Min.Y, w, h))
				draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
			}
			f, err := os.Create("present.png")
			if err != nil {
				log.Fatal(err)
			}
			if err := png.Encode(f, dst); err != nil {
				f.Close()
				log.Fatal(err)
			}
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}

	// Save scene to svg file
	svgdoc, err := document.MakeSVG("seen-svg", WidthDp, HeightDp)
	if err != nil {
		log.Fatal(err)
	}
	if context := svg.ContextWith(svgdoc.GetElementById("seen-svg"), backdrop, curtain, foreground); context != nil {
		context.Render()
	} else {
		log.Fatal("Render context is nil")
	}
	err = svgdoc.SaveToFile("present.svg")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

var BoxPoints = [...]seen.Point{
	// lid
	0: {X: -1, Y: 0.75, Z: -1},
	1: {X: -1, Y: 0.75, Z: 1},
	2: {X: -1, Y: 1, Z: -1},
	3: {X: -1, Y: 1, Z: 1},
	4: {X: 1, Y: 0.75, Z: -1},
	5: {X: 1, Y: 0.75, Z: 1},
	6: {X: 1, Y: 1, Z: -1},
	7: {X: 1, Y: 1, Z: 1},
	// box
	8:  {X: -0.9, Y: -1, Z: -0.9},   // 0'
	9:  {X: -0.9, Y: -1, Z: 0.9},    // 1'
	10: {X: -0.9, Y: 0.75, Z: -0.9}, // 2'
	11: {X: -0.9, Y: 0.75, Z: 0.9},  // 3'
	12: {X: 0.9, Y: -1, Z: -0.9},    // 4'
	13: {X: 0.9, Y: -1, Z: 0.9},     // 5'
	14: {X: 0.9, Y: 0.75, Z: -0.9},  // 6'
	15: {X: 0.9, Y: 0.75, Z: 0.9},   // 7'
}

// Map to points in the surfaces of a cube
var BoxMap = [...][]int{
	// lid
	{0, 1, 3, 2}, // left
	{5, 4, 6, 7}, // right
	{2, 3, 7, 6}, // top
	{3, 1, 5, 7}, // front
	{0, 2, 6, 4}, // back
	// Lid-bottom & box-top
	{0, 2 + 8, 3 + 8, 1}, // a
	{0, 4, 6 + 8, 2 + 8}, // b
	{4, 5, 7 + 8, 6 + 8}, // c
	{1, 3 + 8, 7 + 8, 5}, // d
	// box
	{0 + 8, 1 + 8, 3 + 8, 2 + 8}, // left'
	{5 + 8, 4 + 8, 6 + 8, 7 + 8}, // right'
	{3 + 8, 1 + 8, 5 + 8, 7 + 8}, // front'
	{0 + 8, 2 + 8, 6 + 8, 4 + 8}, // back'
	{1 + 8, 0 + 8, 4 + 8, 5 + 8}, // bottom'
}

func Box() *seen.Shape {
	return &seen.Shape{
		Type:      "box",
		Transform: seen.DefaultTransform,
		Surfaces:  seen.SurfacesWith(BoxPoints[:], BoxMap[:]),
	}
}

var RibbonPoints = [...]seen.Point{
	0:  {X: 1.0, Y: 1.0, Z: 0.2},
	1:  {X: 1.0, Y: 1.0, Z: -0.2},
	2:  {X: 0.2, Y: 1.0, Z: -1.0},
	3:  {X: -0.2, Y: 1.0, Z: -1.0},
	4:  {X: -1.0, Y: 1.0, Z: -0.2},
	5:  {X: -1.0, Y: 1.0, Z: 0.2},
	6:  {X: -0.2, Y: 1.0, Z: 1.0},
	7:  {X: 0.2, Y: 1.0, Z: 1.0},
	8:  {X: 0.2, Y: 1.0, Z: 0.2},
	9:  {X: 0.2, Y: 1.0, Z: -0.2},
	10: {X: -0.2, Y: 1.0, Z: -0.2},
	11: {X: -0.2, Y: 1.0, Z: 0.2},
}

func Ribbon() *seen.Shape {
	points := RibbonPoints[:]
	for i := range RibbonPoints {
		p := RibbonPoints[i]
		p.Y *= -1.0
		points = append(points, p)
	}

	faces := [][]int(nil)
	face := func(p ...int) {
		faces = append(faces, p)
	}
	extrude := func(p, q int, d seen.Point) (r, s int) {
		points = append(points, points[p].Add(d))
		r = len(points) - 1
		points = append(points, points[q].Add(d))
		s = len(points) - 1
		face(p, q, s, r)
		return
	}
	move := func(n, o, p, q int, d seen.Point) (r, s, t, u int) {
		points = append(points, points[n].Add(d))
		r = len(points) - 1
		points = append(points, points[o].Add(d))
		s = len(points) - 1
		points = append(points, points[p].Add(d))
		t = len(points) - 1
		points = append(points, points[q].Add(d))
		u = len(points) - 1
		return
	}

	face(0, 1, 9, 2, 3, 10, 4, 5, 11, 6, 7, 8)

	p1, p0 := extrude(1, 0, seen.Pt(0.0, -0.30, 0.0))
	p1, p0 = extrude(p1, p0, seen.Pt(-0.1, 0.0, 0.0))
	p1, p0 = extrude(p1, p0, seen.Pt(0.0, -1.70, 0.0))

	p3, p2 := extrude(3, 2, seen.Pt(0.0, -0.30, 0.0))
	p3, p2 = extrude(p3, p2, seen.Pt(0.0, 0.0, 0.1))
	p3, p2 = extrude(p3, p2, seen.Pt(0.0, -1.70, 0.0))

	p5, p4 := extrude(5, 4, seen.Pt(0.0, -0.30, 0.0))
	p5, p4 = extrude(p5, p4, seen.Pt(0.1, 0.0, 0.0))
	p5, p4 = extrude(p5, p4, seen.Pt(0.0, -1.70, 0.0))

	p7, p6 := extrude(7, 6, seen.Pt(0.0, -0.30, 0.0))
	p7, p6 = extrude(p7, p6, seen.Pt(0.0, 0.0, -0.1))
	p7, p6 = extrude(p7, p6, seen.Pt(0.0, -1.70, 0.0))

	p8, p9, p10, p11 := move(8, 9, 10, 11, seen.Pt(0.0, -2.0, 0.0))
	face(p1, p0, p8, p7, p6, p11, p5, p4, p10, p3, p2, p9)

	return &seen.Shape{
		Type:      "ribbon",
		Transform: seen.DefaultTransform,
		Surfaces:  seen.SurfacesWith(points, faces),
	}
}
