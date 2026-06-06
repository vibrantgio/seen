package bvh

import "testing"

// knownChannels is the set of valid channel identifiers. Parsed channels must
// match these exactly — in particular with no surrounding whitespace, which a
// grammar bug once left attached (e.g. "Zrotation ", "Yrotation\r\n").
var knownChannels = map[Channel]bool{
	Xposition: true, Yposition: true, Zposition: true,
	Xrotation: true, Yrotation: true, Zrotation: true,
}

func TestChannelsAreClean(t *testing.T) {
	h, err := Load("testdata/05_11.bvh")
	if err != nil {
		t.Fatal(err)
	}
	var walk func(j Joint)
	walk = func(j Joint) {
		for _, ch := range j.Channels {
			if !knownChannels[ch] {
				t.Errorf("joint %q has unrecognised channel %q", j.Id, string(ch))
			}
		}
		for _, c := range j.Joints {
			walk(c)
		}
	}
	walk(h.Root)
}

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
