package shape

import (
	"math"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/point"
)

type PipeConfig struct {
	// radius default 1
	radius float64
	// segments default 8
	segments int
}

type PipeOption func(*PipeConfig)

func Radius(radius float64) PipeOption {
	return func(cfg *PipeConfig) {
		cfg.radius = radius
	}
}
func Segments(segments int) PipeOption {
	return func(cfg *PipeConfig) {
		cfg.segments = segments
	}
}

// Pipe creates a 3D cylinder between two points in space. The cylinder's central axis
// extends from point1 to point2. The cylinder's appearance can be customized using:
//   - Radius: Sets the radius (default: 1)
//   - Segments: Sets the number of sides around the circumference (default: 8)
func Pipe(point1, point2 point.Point, options ...PipeOption) seen.Object {
	cfg := &PipeConfig{radius: 1, segments: 8}
	for _, option := range options {
		option(cfg)
	}
	// Compute a normal perpendicular to the axis point1->point2 and define the
	// rotations about the axis as a quaternion
	axis := point2.Minus(point1)
	perp := axis.Perpendicular().Times(cfg.radius)
	theta := -math.Pi * 2.0 / float64(cfg.segments)

	quat := axis.PointAngle(theta).Normalize().Matrix()

	// Apply the quaternion rotations to create one face
	points := make([]point.Point, cfg.segments)
	for i := range cfg.segments {
		points[i] = point1.Plus(perp)
		perp = perp.Mul(quat)
	}
	return Extrude(points, axis)
}
