package mocap

import (
	"testing"

	"github.com/vibrantgio/seen/bvh"
	"github.com/vibrantgio/seen/quaternion"
)

// TestApplyRotates guards against a frozen skeleton: a mid-capture frame must
// actually rotate joints. (A parser bug once left whitespace on channel names
// so no rotation channel matched, posing every frame in the rest position.)
func TestApplyRotates(t *testing.T) {
	h, err := bvh.Load("../bvh/testdata/05_11.bvh")
	if err != nil {
		t.Fatal(err)
	}
	m := New(h, nil)
	m.Apply(m.Frames() / 2)

	rotated := 0
	for _, j := range m.joints {
		if !j.node.Rotation().Equal(quaternion.Identity) {
			rotated++
		}
	}
	if rotated == 0 {
		t.Fatal("no joint rotated on a mid-capture frame; skeleton is frozen in its rest pose")
	}
	t.Logf("%d/%d joints rotated", rotated, len(m.joints))
}

func TestNew(t *testing.T) {
	h, err := bvh.Load("../bvh/testdata/01_06.bvh")
	if err != nil {
		t.Fatal(err)
	}
	m := New(h, nil)

	if got := m.Frames(); got != len(h.Motion.Frames) {
		t.Fatalf("Frames() = %d, want %d", got, len(h.Motion.Frames))
	}
	if len(m.joints) == 0 {
		t.Fatal("no animated joints collected")
	}

	// The flat per-frame values must be fully consumed: the sum of every
	// joint's channel count equals the frame width. A mismatch means the
	// joint traversal order is out of step with the motion data layout.
	total := 0
	for _, j := range m.joints {
		total += len(j.channels)
	}
	if got := len(h.Motion.Frames[0]); got != total {
		t.Fatalf("frame width = %d, sum of channels = %d", got, total)
	}

	// Every frame produces one pose per animated joint.
	for i, frame := range m.frames {
		if len(frame) != len(m.joints) {
			t.Fatalf("frame %d: %d poses, want %d", i, len(frame), len(m.joints))
		}
	}

	// Apply wraps modulo the frame count and is deterministic.
	m.Apply(0)
	x0, y0, z0 := m.joints[0].node.Translation()
	r0 := m.joints[0].node.Rotation()
	m.Apply(m.Frames()) // wraps back to frame 0
	x1, y1, z1 := m.joints[0].node.Translation()
	r1 := m.joints[0].node.Rotation()
	if x0 != x1 || y0 != y1 || z0 != z1 || !r0.Equal(r1) {
		t.Fatalf("Apply not deterministic across wrap: (%v,%v,%v,%v) vs (%v,%v,%v,%v)",
			x0, y0, z0, r0, x1, y1, z1, r1)
	}
}
