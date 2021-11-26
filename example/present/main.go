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
	"github.com/reactivego/seen/render/zsort"
)

const WidthDp = 1024
const HeightDp = 1024
const MinWidthDp = 1024
const MinHeightDp = 1024

func main() {
	go Present()
	app.Main()
}

func Present() {
	window := app.NewWindow(
		app.Title("Seen - Present"),
		app.Size(unit.Dp(WidthDp), unit.Dp(HeightDp)),
		app.MinSize(unit.Dp(MinWidthDp), unit.Dp(MinHeightDp)),
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

	// QuartRibbon
	qrib := QuartRibbon()
	qrib.SetFill(ribbonfill)
	qrib.SetStroke(ribbonfill)
	qrib.SetScale(50, 50, 50)

	qrib2 := QuartRibbon()
	qrib2.SetFill(ribbonfill)
	qrib2.SetStroke(ribbonfill)
	qrib2.SetScale(50, 50, 50)
	qrib2.SetRotation(quat.RotY(math.Pi / 2.0))

	qrib3 := QuartRibbon()
	qrib3.SetFill(ribbonfill)
	qrib3.SetStroke(ribbonfill)
	qrib3.SetScale(50, 50, 50)
	qrib3.SetRotation(quat.RotY(math.Pi))

	qrib4 := QuartRibbon()
	qrib4.SetFill(ribbonfill)
	qrib4.SetStroke(ribbonfill)
	qrib4.SetScale(50, 50, 50)
	qrib4.SetRotation(quat.RotY(-math.Pi / 2.0))

	qcorn := QuartCorner()
	qcorn.SetFill(boxfill)
	qcorn.SetStroke(boxfill)
	qcorn.SetScale(50, 50, 50)

	qcorn2 := QuartCorner()
	qcorn2.SetFill(boxfill)
	qcorn2.SetStroke(boxfill)
	qcorn2.SetScale(50, 50, 50)
	qcorn2.SetRotation(quat.RotY(math.Pi / 2.0))

	qcorn3 := QuartCorner()
	qcorn3.SetFill(boxfill)
	qcorn3.SetStroke(boxfill)
	qcorn3.SetScale(50, 50, 50)
	qcorn3.SetRotation(quat.RotY(math.Pi))

	qcorn4 := QuartCorner()
	qcorn4.SetFill(boxfill)
	qcorn4.SetStroke(boxfill)
	qcorn4.SetScale(50, 50, 50)
	qcorn4.SetRotation(quat.RotY(1.5 * math.Pi))

	// Present is a box with a lid and a ribbon around it
	//present := seen.GroupWith(box , ribbon)
	present := seen.GroupWith(qrib, qrib2, qrib3, qrib4, qcorn, qcorn2, qcorn3, qcorn4)
	present.SetRotation(quat.RotX(0.25 * math.Pi).RotY(-0.25 * math.Pi))

	// Create scene and add shape to group
	scene := seen.DefaultScene()
	scene.Shader = seen.PhongShader
	scene.FractionalPoints = true
	scene.Group.Add(present)
	scene.Group.SetScale(3, 3, 3)
	scene.Viewport = seen.CenterViewport(0, 0, WidthDp, HeightDp)

	// Create separate layers for the stage.
	backdrop := render.FillLayerWith(WidthDp, HeightDp, 0, 0, backdropfill)
	curtain := render.FillLayerWith(WidthDp, HeightDp/2, 0, 0, curtainfill)

	var foreground render.Layer
	if true {
		foreground = bsp.LayerWith(scene)
	} else {
		foreground = zsort.LayerWith(scene)
	}

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

const chamfer = 0.015

var QuartRibbonPoints = [...]seen.Point{
	// lid-top
	0: {X: 0, Y: 1, Z: 0},
	1: {X: 0.2, Y: 1, Z: -0.2},
	2: {X: 0.2, Y: 1, Z: 0.2},
	3: {X: 1.0 - chamfer, Y: 1, Z: -0.2},
	4: {X: 1.0 - chamfer, Y: 1, Z: 0.2},
	// lid-top-chamfer
	5: {X: 1.0, Y: 1.0 - chamfer, Z: -0.2},
	6: {X: 1.0, Y: 1.0 - chamfer, Z: 0.2},
	// lid-side
	7: {X: 1.0, Y: 0.7 + chamfer, Z: -0.2},
	8: {X: 1.0, Y: 0.7 + chamfer, Z: 0.2},
	// lid-bottom-chamfer
	9:  {X: 1.0 - chamfer, Y: 0.7, Z: -0.2},
	10: {X: 1.0 - chamfer, Y: 0.7, Z: 0.2},
	// lid-bottom
	11: {X: 0.85, Y: 0.7, Z: -0.2},
	12: {X: 0.85, Y: 0.7, Z: 0.2},
	// side
	13: {X: 0.85, Y: -1, Z: -0.2},
	14: {X: 0.85, Y: -1, Z: 0.2},
	// bottom
	15: {X: 0.2, Y: -1, Z: -0.2},
	16: {X: 0.2, Y: -1, Z: 0.2},
	17: {X: 0, Y: -1, Z: 0},
}

var QuartRibbonMap = [...][]int{
	0: {1, 0, 2, 4, 3},      // lid-top
	1: {3, 4, 6, 5},         // lid-top-chamfer
	2: {5, 6, 8, 7},         // lid-side
	3: {7, 8, 10, 9},        // lid-bottom-chamfer
	4: {9, 10, 12, 11},      // lid-bottom
	5: {11, 12, 14, 13},     // side
	6: {13, 14, 16, 17, 15}, // bottom
}

func QuartRibbon() *seen.Shape {
	return &seen.Shape{
		Type:      "qribbon",
		Transform: seen.DefaultTransform,
		Surfaces:  seen.SurfacesWith(QuartRibbonPoints[:], QuartRibbonMap[:]),
	}
}

var QuartCornerPoints = [...]seen.Point{
	// lid top
	0: {X: 0.2, Y: 1, Z: 1.0 - chamfer},
	1: {X: 0.2, Y: 1, Z: 0.2},
	2: {X: 1.0 - chamfer, Y: 1, Z: 1.0 - chamfer},
	3: {X: 1.0 - chamfer, Y: 1, Z: 0.2},
	// chamfer
	4: {X: 1.0, Y: 1 - chamfer, Z: 1.0 - chamfer},
	5: {X: 1.0, Y: 1 - chamfer, Z: 0.2},
	// lid side
	6: {X: 1.0, Y: 0.7 + chamfer, Z: 1.0 - chamfer},
	7: {X: 1.0, Y: 0.7 + chamfer, Z: 0.2},
	// chamfer
	8: {X: 1.0 - chamfer, Y: 0.7, Z: 1.0 - chamfer},
	9: {X: 1.0 - chamfer, Y: 0.7, Z: 0.2},
	// lid bottom
	10: {X: 0.85, Y: 0.7, Z: 0.85},
	11: {X: 0.85, Y: 0.7, Z: 0.2},
	// box side
	12: {X: 0.85, Y: -1.0, Z: 0.85},
	13: {X: 0.85, Y: -1.0, Z: 0.2},
	// box bottom
	14: {X: 0.2, Y: -1, Z: 0.85},
	15: {X: 0.2, Y: -1, Z: 0.2},
	//
	16: {X: 0.2, Y: 1 - chamfer, Z: 1.0},
	17: {X: 1.0 - chamfer, Y: 1 - chamfer, Z: 1.0},
	//
	18: {X: 0.2, Y: 0.7 + chamfer, Z: 1.0},
	19: {X: 1.0 - chamfer, Y: 0.7 + chamfer, Z: 1.0},
	//
	20: {X: 0.2, Y: 0.7, Z: 1.0 - chamfer},
	21: {X: 1.0 - chamfer, Y: 0.7, Z: 1.0 - chamfer},
	//
	22: {X: 0.2, Y: 0.7, Z: 0.85},
	23: {X: 0.85, Y: 0.7, Z: 0.85},
	//
	24: {X: 0.2, Y: -1, Z: 0.85},
	25: {X: 0.85, Y: -1, Z: 0.85},
}

var QuartCornerMap = [...][]int{
	0:  {1, 0, 2, 3},
	1:  {3, 2, 4, 5},
	2:  {5, 4, 6, 7},
	3:  {7, 6, 8, 9},
	4:  {9, 8, 10, 11},
	5:  {11, 10, 12, 13},
	6:  {13, 12, 14, 15},
	7:  {2, 0, 16, 17},
	8:  {17, 16, 18, 19},
	9:  {19, 18, 20, 21},
	10: {21, 20, 22, 23},
	11: {23, 22, 24, 25},
	12: {2, 17, 4},
	13: {4, 17, 19, 6},
	14: {6, 19, 8},
}

func QuartCorner() *seen.Shape {
	return &seen.Shape{
		Type:      "qcorner",
		Transform: seen.DefaultTransform,
		Surfaces:  seen.SurfacesWith(QuartCornerPoints[:], QuartCornerMap[:]),
	}
}
