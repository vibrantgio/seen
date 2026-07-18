package seen

import (
	"math"
	"testing"

	"github.com/vibrantgio/seen/matrix"
)

// TestFitEquivalence proves that FitCenter/FitOrigin reproduce the legacy
// viewport Center/Origin pipeline exactly. The legacy formulas are hardcoded
// here (they no longer exist in the viewport package):
//
//	prescale         = Scale(1/W, 1/H, 1/D).Translate(-x, -y, -D)
//	postscale center = Translate(x+w/2, y+h/2, D).Scale(W, -H, D)
//	postscale origin = Translate(x, y, D).Scale(W, -H, D)
//
// with W, H, D = w, h, h — or dist, dist, dist when dist is given — and the
// legacy render pipeline composing Projection · prescale · Camera.Matrix().
// The new pipeline composes Projection · Camera.View() and maps to pixels
// with Viewport.Screen.
func TestFitEquivalence(t *testing.T) {
	tuples := []struct {
		name       string
		x, y, w, h float64
		dist       []float64
	}{
		{"launcher", 0, 0, 1100, 760, []float64{2200}},
		{"non-square", 10, 20, 600, 300, nil},
		{"unit", 0, 0, 1, 1, nil},
		{"offset-locked", 5, -3, 400, 250, []float64{100}},
	}
	fits := []struct {
		name      string
		fit       func(s *Scene, x, y, w, h float64, dist ...float64)
		postscale func(x, y, w, h, W, H, D float64) matrix.Matrix
	}{
		{
			"FitCenter",
			(*Scene).FitCenter,
			func(x, y, w, h, W, H, D float64) matrix.Matrix {
				return matrix.Translate(x+w/2, y+h/2, D).Scale(W, -H, D)
			},
		},
		{
			"FitOrigin",
			(*Scene).FitOrigin,
			func(x, y, w, h, W, H, D float64) matrix.Matrix {
				return matrix.Translate(x, y, D).Scale(W, -H, D)
			},
		},
	}
	cameras := []struct {
		name  string
		setup func(s *Scene)
	}{
		{"default camera", func(s *Scene) {}},
		{"moved camera", func(s *Scene) {
			s.Camera.RotY(0.3).RotX(-0.2)
			s.Camera.SetTranslation(3, -7, 42)
		}},
	}

	for _, tc := range tuples {
		for _, fit := range fits {
			for _, cam := range cameras {
				t.Run(tc.name+"/"+fit.name+"/"+cam.name, func(t *testing.T) {
					s := NewScene()
					cam.setup(s)
					transformBefore := s.Camera.Transform
					projectionBefore := s.Camera.Projection

					fit.fit(s, tc.x, tc.y, tc.w, tc.h, tc.dist...)

					if s.Camera.Transform != transformBefore {
						t.Error("fit changed Camera.Transform")
					}
					if s.Camera.Projection != projectionBefore {
						t.Error("fit changed Camera.Projection")
					}

					W, H, D := tc.w, tc.h, tc.h
					if len(tc.dist) > 0 {
						W, H, D = tc.dist[0], tc.dist[0], tc.dist[0]
					}
					prescale := matrix.Scale(1/W, 1/H, 1/D).Translate(-tc.x, -tc.y, -D)
					oldProj := s.Camera.Projection.Mul(prescale).Mul(s.Camera.Matrix())
					newProj := s.Camera.Projection.Mul(s.Camera.View())
					for i := range oldProj {
						for j := range oldProj[i] {
							if d := math.Abs(newProj[i][j] - oldProj[i][j]); d > 1e-12 {
								t.Errorf("projection[%d][%d] = %v, legacy %v (|diff| %v > 1e-12)",
									i, j, newProj[i][j], oldProj[i][j], d)
							}
						}
					}

					postscale := fit.postscale(tc.x, tc.y, tc.w, tc.h, W, H, D)
					if s.Viewport.Screen != postscale {
						t.Errorf("Screen = %v, want legacy postscale %v", s.Viewport.Screen, postscale)
					}
				})
			}
		}
	}
}

// TestDefaultViewportEquivalence pins viewport.Default to the screen mapping
// of the legacy Origin(0, 0, 1, 1) viewport's postscale.
func TestDefaultViewportEquivalence(t *testing.T) {
	s := NewScene()
	legacy := matrix.Translate(0, 0, 1).Scale(1, -1, 1)
	if s.Viewport.Screen != legacy {
		t.Errorf("default Screen = %v, want legacy Origin(0,0,1,1) postscale %v", s.Viewport.Screen, legacy)
	}
}
