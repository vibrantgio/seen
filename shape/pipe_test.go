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
