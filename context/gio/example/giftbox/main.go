package main

import (
	"image"
	"image/png"
	"log"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/gpu/headless"
	"gioui.org/op"
	"gioui.org/unit"
	"golang.org/x/image/draw"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/context/gio"
	"github.com/vibrantgio/seen/context/svg"
	"github.com/vibrantgio/seen/drag"
	"github.com/vibrantgio/seen/face"
	"github.com/vibrantgio/seen/layer/backdrop"
	"github.com/vibrantgio/seen/layer/nsort"
	"github.com/vibrantgio/seen/point"
	"github.com/vibrantgio/seen/quaternion"
	"github.com/vibrantgio/seen/shader"
	"github.com/vibrantgio/seen/shape"
	"github.com/vibrantgio/seen/viewport"
	"github.com/vibrantgio/seen/zoom"
)

const should_save_to_png = false
const should_save_to_svg = false

func main() {
	go GiftBox()
	app.Main()
}

func GiftBox() {
	const SIZE = 1024

	window := new(app.Window)
	window.Option(
		app.Title("Seen - Gift Box"),
		app.Size(SIZE, SIZE),
		app.MinSize(SIZE/2, SIZE/2))

	// Colors to use
	lightblue := color.Color{R: 1.0 / 255.0, G: 202.0 / 255.0, B: 252.0 / 255.0, A: 1.0}
	darkblue := color.Color{R: 0.0, G: 130.0 / 255.0, B: 193.0 / 255.0, A: 1.0}
	hardblue, _ := color.ColorWithString("#cceeff")
	lightgrey, _ := color.ColorWithString("#eeeeee")
	lightorange := color.Color{R: 247.0 / 255.0, G: 148.0 / 255.0, B: 29.0 / 255.0, A: 1.0}
	darkorange := color.Color{R: 224.0 / 255.0, G: 134.0 / 255.0, B: 26.0 / 255.0, A: 1.0}

	_, _, _, _, _, _ = lightblue, darkblue, hardblue, lightgrey, lightorange, darkorange

	backdropfill := darkblue
	curtainfill := color.White
	boxfill := lightblue
	ribbonfill := darkorange

	// QuartRibbon
	qrib := QuartRibbon()
	qrib.Faces().SetFill(ribbonfill)
	qrib.Faces().SetStroke(ribbonfill)
	qrib.SetScale(50, 50, 50)

	qrib2 := QuartRibbon()
	qrib2.Faces().SetFill(ribbonfill)
	qrib2.Faces().SetStroke(ribbonfill)
	qrib2.SetScale(50, 50, 50)
	qrib2.SetRotation(quaternion.RotY(math.Pi / 2.0))

	qrib3 := QuartRibbon()
	qrib3.Faces().SetFill(ribbonfill)
	qrib3.Faces().SetStroke(ribbonfill)
	qrib3.SetScale(50, 50, 50)
	qrib3.SetRotation(quaternion.RotY(math.Pi))

	qrib4 := QuartRibbon()
	qrib4.Faces().SetFill(ribbonfill)
	qrib4.Faces().SetStroke(ribbonfill)
	qrib4.SetScale(50, 50, 50)
	qrib4.SetRotation(quaternion.RotY(-math.Pi / 2.0))

	qcorn := QuartCorner()
	qcorn.Faces().SetFill(boxfill)
	qcorn.Faces().SetStroke(boxfill)
	qcorn.SetScale(50, 50, 50)

	qcorn2 := QuartCorner()
	qcorn2.Faces().SetFill(boxfill)
	qcorn2.Faces().SetStroke(boxfill)
	qcorn2.SetScale(50, 50, 50)
	qcorn2.SetRotation(quaternion.RotY(math.Pi / 2.0))

	qcorn3 := QuartCorner()
	qcorn3.Faces().SetFill(boxfill)
	qcorn3.Faces().SetStroke(boxfill)
	qcorn3.SetScale(50, 50, 50)
	qcorn3.SetRotation(quaternion.RotY(math.Pi))

	qcorn4 := QuartCorner()
	qcorn4.Faces().SetFill(boxfill)
	qcorn4.Faces().SetStroke(boxfill)
	qcorn4.SetScale(50, 50, 50)
	qcorn4.SetRotation(quaternion.RotY(1.5 * math.Pi))

	// A Gift Box has a lid and a ribbon around it
	//giftbox := seen.GroupWith(box , ribbon)
	giftbox := seen.NewGroup(qrib, qrib2, qrib3, qrib4, qcorn, qcorn2, qcorn3, qcorn4)
	giftbox.SetRotation(quaternion.RotX(0.25 * math.Pi).RotY(-0.25 * math.Pi))

	// Create scene and add shape to group
	scene := seen.NewDefaultScene()
	scene.Shader = shader.Phong
	scene.Group.Add(giftbox)
	scene.Group.SetScale(3, 3, 3)

	// Create separate layers for the stage.
	background := backdrop.NewLayer(SIZE, SIZE, 0, 0, backdropfill)
	curtain := backdrop.NewLayer(SIZE, SIZE/2, 0, 0, curtainfill)

	// Create a layer that renders a scene by depth-sorting the polygons for
	// the current eye
	foreground := nsort.NewLayerForScene(scene)

	// Create a context that renders into the gio window
	context := gio.NewContext(window, background, curtain, foreground)

	// Enable dragger-to-rotate
	context.Drag(drag.Inertia(true)).On(func(e drag.Event) {
		r := giftbox.Rotation()
		r = quaternion.RotX(e.Dy / 150).Mul(r).RotY(e.Dx / 150)
		giftbox.SetRotation(r)
		context.Render()
	})

	// Enable mouse-wheel zoom
	context.Zoom().On(func(e zoom.Event) {
		sx, sy, sz := giftbox.Scale()
		giftbox.SetScale(sx*e.Zoom, sy*e.Zoom, sz*e.Zoom)
	})

	widget := gio.Widget(context, func(w, h unit.Dp) {
		background.Width, background.Height = float64(w), float64(h)
		curtain.Width, curtain.Height = float64(w), float64(h)/2
		scene.Viewport = viewport.Center(0, 0, float64(w), float64(h))
	})

	metric := unit.Metric{PxPerDp: 1.0, PxPerSp: 1.0}
	size := image.Point{X: SIZE, Y: SIZE}
	ops := new(op.Ops)
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			metric = e.Metric
			size = e.Size
			gtx := app.NewContext(ops, e)
			widget(gtx)
			e.Frame(ops)
		}
	}

	// Save scene to png file
	if should_save_to_png {
		if window, err := headless.NewWindow(size.X, size.Y); err == nil {
			window.Frame(ops)
			src := image.NewRGBA(image.Rectangle{Max: window.Size()})
			if err := window.Screenshot(src); err == nil {
				dst := src
				// Scale image when ppdp (pixels per device pixel) is not 1.0
				if metric.PxPerDp != 1.0 {
					sb := src.Bounds()
					w, h := int(metric.PxToDp(sb.Dx())), int(metric.PxToDp(sb.Dy()))
					dst = image.NewRGBA(image.Rect(sb.Min.X, sb.Min.Y, sb.Min.X+w, sb.Min.Y+h))
					draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
				}
				f, err := os.Create("giftbox.png")
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
	}

	// Save scene to svg file
	if should_save_to_svg {
		svgdoc, err := svg.NewSVG("seen-svg", int(background.Width), int(background.Height))
		if err != nil {
			log.Fatal(err)
		}
		if context := svg.NewContext(svgdoc.GetElementById("seen-svg"), background, curtain, foreground); context != nil {
			context.Render()
		} else {
			log.Fatal("Render context is nil")
		}
		err = svgdoc.SaveToFile("giftbox.svg")
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(0)
}

const chamfer = 0.015

var QuartRibbonPoints = [...]point.Point{
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

var QuartRibbonFacets = [...]face.Facet{
	0: {1, 0, 2, 4, 3},      // lid-top
	1: {3, 4, 6, 5},         // lid-top-chamfer
	2: {5, 6, 8, 7},         // lid-side
	3: {7, 8, 10, 9},        // lid-bottom-chamfer
	4: {9, 10, 12, 11},      // lid-bottom
	5: {11, 12, 14, 13},     // side
	6: {13, 14, 16, 17, 15}, // bottom
}

func QuartRibbon() seen.Object {
	return shape.NewShape("qribbon", QuartRibbonPoints[:], QuartRibbonFacets[:])
}

var QuartCornerPoints = [...]point.Point{
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

var QuartCornerFacets = [...]face.Facet{
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

func QuartCorner() seen.Object {
	return shape.NewShape("qcorner", QuartCornerPoints[:], QuartCornerFacets[:])
}
