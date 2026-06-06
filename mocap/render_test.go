package mocap_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/bvh"
	"github.com/vibrantgio/seen/mocap"
	"github.com/vibrantgio/seen/context/svg"
	"github.com/vibrantgio/seen/layer/bsort"
	"github.com/vibrantgio/seen/shader"
	"github.com/vibrantgio/seen/viewport"
)

// TestRenderFrame poses the skeleton and renders one frame through the real
// layer pipeline (bsort sort -> Phong shade -> SVG canvas) to confirm the
// render path holds and emits geometry. It uses the CPU-only SVG context, so
// it needs no GPU or display. This complements TestNew, which only exercises
// model construction.
func TestRenderFrame(t *testing.T) {
	h, err := bvh.Load("../bvh/testdata/01_06.bvh")
	if err != nil {
		t.Fatal(err)
	}

	m := mocap.New(h, nil)
	m.Apply(m.Frames() / 2) // pose mid-capture, not the rest pose

	scene := seen.NewDefaultScene()
	scene.Shader = shader.Phong
	scene.Group.Add(m.Group)
	scene.Viewport = viewport.Center(0, 0, 500, 500)

	doc, err := svg.NewSVG("seen-svg", 500, 500)
	if err != nil {
		t.Fatal(err)
	}
	ctx := svg.NewContext(doc.GetElementById("seen-svg"), bsort.NewLayerForScene(scene))
	if ctx == nil {
		t.Fatal("nil svg render context")
	}
	ctx.Render()

	var buf bytes.Buffer
	if _, err := doc.WriteTo(&buf); err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(buf.String(), "<path"); n < 10 {
		t.Fatalf("rendered SVG has %d <path> elements, want many (output %d bytes)", n, buf.Len())
	}
}
