package float_test

import (
	"math"
	"testing"

	"github.com/vibrantgio/seen/float"
)

// TestAlmostEqual exercises many edge cases relevant to comparing floats.
func TestAlmostEqual(t *testing.T) {
	const eps = 1e-10

	tests := []struct {
		name    string
		a, b    float64
		epsilon float64
		want    bool
	}{
		{
			name:    "Exactly the same float",
			a:       1.2345,
			b:       1.2345,
			epsilon: eps,
			want:    true,
		},
		{
			name:    "Both zero",
			a:       0.0,
			b:       0.0,
			epsilon: eps,
			want:    true,
		},
		{
			name:    "Small absolute difference (within eps) near zero",
			a:       1e-11,
			b:       0.0,
			epsilon: 1e-10,
			want:    true,
		},
		{
			name:    "Small absolute difference (greater than eps) near zero",
			a:       1e-9,
			b:       0.0,
			epsilon: 1e-10,
			want:    false,
		},
		{
			name:    "Relative difference within eps (larger values)",
			a:       1000.0,
			b:       1000.0000000001, // very small % difference
			epsilon: eps,
			want:    true,
		},
		{
			name:    "Relative difference exceeds eps (larger values)",
			a:       1000.0,
			b:       1000.001, // bigger difference
			epsilon: 1e-7,
			want:    false,
		},
		{
			name:    "Negative vs. positive same magnitude (large) - should fail",
			a:       -1e5,
			b:       1e5,
			epsilon: 1e-3,
			want:    false,
		},
		{
			name:    "NaN vs. normal",
			a:       math.NaN(),
			b:       1.0,
			epsilon: eps,
			want:    false,
		},
		{
			name:    "Both are NaN",
			a:       math.NaN(),
			b:       math.NaN(),
			epsilon: eps,
			want:    false, // by design, we return false on NaN
		},
		{
			name:    "Both +Inf",
			a:       math.Inf(1),
			b:       math.Inf(1),
			epsilon: eps,
			want:    true, // a==b passes the first check
		},
		{
			name:    "Both -Inf",
			a:       math.Inf(-1),
			b:       math.Inf(-1),
			epsilon: eps,
			want:    true,
		},
		{
			name:    "+Inf vs. -Inf",
			a:       math.Inf(1),
			b:       math.Inf(-1),
			epsilon: eps,
			want:    false,
		},
		{
			name:    "One Inf, one finite",
			a:       math.Inf(1),
			b:       1e10,
			epsilon: eps,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := float.AlmostEqual(tt.epsilon, tt.a, tt.b)
			if got != tt.want {
				t.Errorf("EqualThreshold(%g, %g, %g) = %v; want %v",
					tt.a, tt.b, tt.epsilon, got, tt.want)
			}
		})
	}
}
