package bvh

import "testing"

func TestLoad(t *testing.T) {
	files := []string{
		"testdata/Example1.bvh",
		"testdata/01_06.bvh",
		"testdata/05_11.bvh",
	}
	for _, f := range files {
		t.Run(f, func(t *testing.T) {
			h, err := Load(f)
			if err != nil {
				t.Fatalf("Load(%q): %v", f, err)
			}
			if h.Root.Id == "" {
				t.Error("root joint has empty Id")
			}
			if len(h.Motion.Frames) == 0 {
				t.Error("no motion frames parsed")
			}
			if h.Motion.FrameTime <= 0 {
				t.Errorf("non-positive FrameTime: %s", h.Motion.FrameTime)
			}
			t.Logf("root=%s joints=%d frames=%d frametime=%s",
				h.Root.Id, len(h.Root.Joints), len(h.Motion.Frames), h.Motion.FrameTime)
		})
	}
}
