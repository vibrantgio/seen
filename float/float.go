package float

import "math"

// Epsilon is some tiny value that determines how precisely equal we want our floats to be.
// We've choosen 0.0000000001 (1e-10) as the threshold value.
const Epsilon float64 = 1e-10

// EqualPairs is  a utility function to compare pairs of floats.
func EqualPairs(f ...float64) bool {
	return AlmostEqualPairs(Epsilon, f...)
}

// Equal is a safe utility function to compare floats.
// It uses AlmostEqual under the hood with an epsilon of 1e-10
func Equal(a, b float64) bool {
	return AlmostEqual(Epsilon, a, b)
}

// AlmostEqualPairs is  a utility function to compare pairs of floats.
func AlmostEqualPairs(eps float64, f ...float64) bool {
	flen := len(f)
	if flen%2 != 0 {
		return false // didn't supply an even number of floats
	}
	for i := 0; i < len(f); i += 2 {
		if !AlmostEqual(eps, f[i], f[i+1]) {
			return false // at least 1 pair didn't equate
		}
	}
	return true // all pairs were equal when considering the Epsilon threshold.
}

// AlmostEqual determines if a and b are “close enough” using a combination
// of absolute and relative checks. This version also guards against
// NaN/Inf edge cases more explicitly, which can avoid subtle pitfalls.
// It's inspired by code from http://floating-point-gui.de/errors/comparison/
func AlmostEqual(eps, a, b float64) bool {
	// If exactly equal or both infinities of the same sign:
	if a == b {
		return true
	}

	// If either is NaN, they're never equal.
	if math.IsNaN(a) || math.IsNaN(b) {
		return false
	}

	// If one is Inf and the other isn't the same Inf, they're not equal (the a == b check above
	// already handles the case where both are the same Inf).
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return false
	}

	// Compute absolute difference.
	diff := math.Abs(a - b)

	// 1) Absolute-check for near-zero values.
	//    If both are small or close to zero,
	//    we say diff < eps is good enough.
	if math.Abs(a) < 1.0 && math.Abs(b) < 1.0 {
		return diff < eps
	}

	// 2) Otherwise do a relative comparison. The idea:
	//    “diff / max(|a|, |b|) < eps”
	//    or “diff / min(|a|, |b|) < eps,”
	//    depending on which logic you prefer for large/small scales.
	//    Here we choose max, so that relative error is small if
	//    a and b are both large numbers.
	denom := math.Max(math.Abs(a), math.Abs(b))
	return diff/denom < eps
}
