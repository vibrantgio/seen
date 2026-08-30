package camera

import (
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/point"
	"github.com/vibrantgio/seen/projection"
	"github.com/vibrantgio/seen/transform"
)

// Camera owns the full world→view pose: the world transform, the eye
// position, the view normalization, and the perspective projection.
//
// A world point p travels Projection · View() · p into clip space; the
// Viewport then maps the normalized device coordinates to pixels.
type Camera struct {
	// Transform is applied to WORLD points directly — it IS a view matrix,
	// not a pose that gets inverted. Rotation therefore orbits the world
	// origin.
	transform.Transform

	// Projection is the perspective (or orthographic) projection matrix,
	// applied after View().
	Projection matrix.Matrix

	// Eye is the eye position in world space. Its translation is applied
	// AFTER Transform (see View), so it does not disturb Transform's
	// world-transform semantics. The fitting helpers place it
	// above the fitted region; Default puts it at (0, 0, 1).
	Eye point.Point

	// Norm is the view normalization — the fit scale that maps the region
	// of interest into the canonical view volume. matrix.Identity when
	// unused.
	Norm matrix.Matrix
}

// Default is the zero-config camera: identity world transform, the default
// perspective projection, eye at (0, 0, 1), and identity normalization.
var Default = CameraWithProjection(projection.DefaultPerspective)

// CameraWithProjection returns Default with the given projection matrix.
func CameraWithProjection(projection matrix.Matrix) Camera {
	return Camera{
		Transform:  transform.Default,
		Projection: projection,
		Eye:        point.Pt(0, 0, 1),
		Norm:       matrix.Identity,
	}
}

// View returns the world→view matrix:
//
//	Norm · Translate(−Eye) · Transform.Matrix()
//
// Transform first (world-transform semantics), then the eye translation,
// then the fit normalization.
func (c Camera) View() matrix.Matrix {
	return c.Norm.Mul(matrix.Translate(-c.Eye.X, -c.Eye.Y, -c.Eye.Z)).Mul(c.Matrix())
}

// EyeInWorld returns the eye's world position: the preimage of the
// view-space origin under View().
//
// The projection matrix's row 3 of [0, 0, −1, 0] makes w_clip vanish exactly
// at the view-space origin, so the center of projection is the world point
// that View() maps to the origin, i.e. View()⁻¹ · 0. This is independent of
// the frustum's near/far and correctly accounts for camera dolly and fit
// offsets. Sorting layers need it to order faces from the eye outward.
//
// For a degenerate view (e.g. a zero-scale Norm) it falls back to the
// estimate −1/proj[2][2] against the fully composed projection
// Projection · View().
func (c Camera) EyeInWorld() point.Point {
	if inv, ok := c.View().Invert(); ok {
		ex, ey, ez := inv.Transform3(0, 0, 0)
		return point.Pt(ex, ey, ez)
	}
	// Degenerate view (e.g. zero scale). Fall back to the estimate.
	proj := c.Projection.Mul(c.View())
	return point.Pt(0, 0, -1.0/proj[2][2])
}
