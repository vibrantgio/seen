package shape

import (
	"testing"

	"github.com/vibrantgio/seen/point"
)

// TestPipe guards against the Extrude length bug that previously made Pipe
// panic: a pipe with s segments extrudes 2*s points into s side faces plus a
// front and back cap.
func TestPipe(t *testing.T) {
	const segments = 8
	obj := Pipe(point.Pt(0, 0, 0), point.Pt(0, 10, 0), Radius(2), Segments(segments))
	if got, want := len(obj.Faces()), segments+2; got != want {
		t.Fatalf("Pipe faces = %d, want %d", got, want)
	}
}

// TestPipeOutwardNormals guards against the Extrude cap bug that placed both
// caps at the near ring (leaving one cap facing inward and the far end
// uncapped). Every face of a pipe must face outward: side faces radially, and
// each cap away from the solid along the axis.
func TestPipeOutwardNormals(t *testing.T) {
	const segments = 8
	const length = 10.0
	obj := Pipe(point.Pt(0, 0, 0), point.Pt(0, 0, length), Radius(2), Segments(segments))

	for k, f := range obj.Faces() {
		bc := f.Points.Barycenter()
		normal := f.Points.Normal().Normalize()

		var outward point.Point
		if len(f.Points) == segments { // cap
			if bc.Z < length/2 {
				outward = point.Pt(0, 0, -1) // near cap faces -z
			} else {
				outward = point.Pt(0, 0, 1) // far cap faces +z
			}
		} else { // side face: radial
			outward = point.Pt(bc.X, bc.Y, 0).Normalize()
		}

		if normal.Dot(outward) <= 0 {
			t.Errorf("face %d (bc.z=%.1f) faces inward: normal=%v, expected outward=%v",
				k, bc.Z, normal, outward)
		}
	}
}
