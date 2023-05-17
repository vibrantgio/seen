package shape

import (
	"math"

	"github.com/reactivego/seen"
)

// ALTITUDE of equilateral triangle for computing triangular patch size is sqrt(3)/2
// which approximates to 0.86602540378443864676372317075293618347140262690519.
const ALTITUDE = 0.86602540378443864676372317075293618347140262690519

// Patch generates a triangular patch with the specified number of columns (nx) and rows (ny).
// The patch is made up of equilateral triangles and is returned as a seen.Shape object.
func Patch(nx, ny float64) *seen.Shape {
	nx = math.Round(nx)
	ny = math.Round(ny)
	var surfaces []seen.Points
	for x := 0.0; x < nx; x++ {
		var triangularPatch []seen.Points
		for y := 0.0; y < ny; y++ {
			triangles := []seen.Points{{
				{X: x, Y: y},
				{X: x + 1, Y: y - 0.5},
				{X: x + 1, Y: y + 0.5},
			}, {
				{X: x, Y: y},
				{X: x + 1, Y: y + 0.5},
				{X: x, Y: y + 1},
			}}
			for i := range triangles {
				for j := range triangles[i] {
					triangles[i][j].X *= ALTITUDE
					if int(x)%2 == 0 {
						triangles[i][j].Y += 0.5
					}
				}
				triangularPatch = append(triangularPatch, triangles[i])
			}
		}
		if int(x)%2 != 0 {
			for i := range triangularPatch[0] {
				triangularPatch[0][i].Y += ny
			}
			triangularPatch = append(triangularPatch, triangularPatch[0])
			triangularPatch = triangularPatch[1:]
		}
		surfaces = append(surfaces, triangularPatch...)
	}

	var surfaceObjects seen.Surfaces
	for _, s := range surfaces {
		surfaceObjects = append(surfaceObjects, *seen.NewSurfaceWith(s))
	}

	return &seen.Shape{Type: "patch", Transform: seen.DefaultTransform, Surfaces: surfaceObjects}
}
