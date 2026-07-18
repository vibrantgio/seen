package camera_test

import (
	"testing"

	"github.com/vibrantgio/seen/camera"
	"github.com/vibrantgio/seen/float"
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/point"
)

// View() must compose Norm · Translate(−Eye) · Transform.Matrix() — exactly
// that factor order. Non-uniform Norm, non-zero Eye, and a rotated,
// translated transform make every reordering of the factors distinguishable.
func TestViewFactorOrder(t *testing.T) {
	cam := camera.Default
	cam.SetTranslation(2, -1, 4)
	cam.RotY(0.3)
	cam.Eye = point.Pt(3, 5, 7)
	cam.Norm = matrix.Scale(0.5, 0.25, 0.125)

	// Hand-built with the matrix builder (chaining right-multiplies):
	// Scale(...).Translate(...) = Norm · T(−Eye), then · Transform.Matrix().
	want := matrix.Scale(0.5, 0.25, 0.125).Translate(-3, -5, -7).Mul(cam.Matrix())
	if got := cam.View(); !got.Equal(want) {
		t.Errorf("View() = %.4v, want %.4v", got, want)
	}

	// Sanity: the reversed grouping (translate before scale) differs, so the
	// comparison above actually pins the factor order.
	reversed := matrix.Translate(-3, -5, -7).Scale(0.5, 0.25, 0.125).Mul(cam.Matrix())
	if reversed.Equal(want) {
		t.Fatal("test camera does not discriminate factor order")
	}

	// Independent scalar-path check: transform the point first (world
	// transform, ADR-003), then subtract the eye, then apply the fit scale
	// componentwise.
	px, py, pz := 1.0, -2.0, 3.0
	wx, wy, wz := cam.Transform.Transform(px, py, pz)
	ex, ey, ez := 0.5*(wx-3), 0.25*(wy-5), 0.125*(wz-7)
	gx, gy, gz := cam.View().Transform3(px, py, pz)
	if !float.EqualPairs(gx, ex, gy, ey, gz, ez) {
		t.Errorf("View()·p = (%v, %v, %v), want (%v, %v, %v)", gx, gy, gz, ex, ey, ez)
	}
}

// With an identity transform the eye's world position is Eye itself; the fit
// scale Norm preserves the origin and so must not move the recovered eye.
func TestEyeInWorldIdentityTransform(t *testing.T) {
	cam := camera.Default
	if got := cam.EyeInWorld(); !got.Equal(point.Pt(0, 0, 1)) {
		t.Errorf("Default EyeInWorld() = %v, want (0, 0, 1)", got)
	}

	cam.Eye = point.Pt(2, -3, 7)
	cam.Norm = matrix.Scale(0.25, 0.5, 0.125)
	if got := cam.EyeInWorld(); !got.Equal(cam.Eye) {
		t.Errorf("EyeInWorld() = %v, want %v", got, cam.Eye)
	}
}

// The mocap dolly case: SetTranslation moves WORLD points (ADR-003), so the
// world point that lands at the eye is Eye minus the camera translation.
func TestEyeInWorldTranslatedCamera(t *testing.T) {
	cam := camera.Default
	cam.SetTranslation(0, 0, 5) // dolly: world points move +5 in z

	if got := cam.EyeInWorld(); !got.Equal(point.Pt(0, 0, -4)) {
		t.Errorf("EyeInWorld() = %v, want (0, 0, -4)", got)
	}
}

// A zero-scale Norm makes View() singular; EyeInWorld must then fall back to
// the legacy −1/proj[2][2] estimate against the COMPOSED projection
// Projection · View(). The transform's z-scale of 2 doubles proj[2][2]
// relative to the raw Projection, so using the wrong matrix would be caught.
func TestEyeInWorldFallback(t *testing.T) {
	cam := camera.Default
	cam.SetScale(1, 1, 2)
	cam.Norm = matrix.Scale(0, 1, 1) // zero-scale: view loses the x axis

	if _, ok := cam.View().Invert(); ok {
		t.Fatal("test view unexpectedly invertible; fallback path not exercised")
	}

	// Projection = Frustum(1, 1, 1, 100): proj[2][2] of the composed matrix
	// is (f+n)/(n−f) · 2 = −202/99, so the fallback eye is (0, 0, 99/202).
	if got := cam.EyeInWorld(); !got.Equal(point.Pt(0, 0, 99.0/202.0)) {
		t.Errorf("EyeInWorld() = %v, want (0, 0, %v)", got, 99.0/202.0)
	}
}
