package float

import "math"

// Epsilon is some tiny value that determines how precisely equal we want our floats to be
// This is exported and left as a variable in case you want to change the default threshold for the
// purposes of certain methods (e.g. Unproject uses the default epsilon when determining
// if the determinant is "close enough" to zero to mean there's no inverse).
//
// This is, obviously, not mutex protected so be **absolutely sure** that no functions using Epsilon
// are being executed when you change this.
var Epsilon float64 = 1e-10

// 1 / 2**(127 - 1)
var MinNormal float64 = 1.1754943508222875e-38
var MinValue float64 = math.SmallestNonzeroFloat64 // 4.940656458412465441765687928682213723651e-324
var MaxValue float64 = math.MaxFloat64             // 1.797693134862315708145274237317043567981e+308

// Equal is a safe utility function to compare floats.
// It's Taken from http://floating-point-gui.de/errors/comparison/
//
// It is slightly altered to not call Abs when not needed.
func Equal(a, b float64) bool {
	return EqualThreshold(a, b, Epsilon)
}

// EqualPairs is  a utility function to compare pairs of floats.
func EqualPairs(f ...float64) bool {
	flen := len(f)
	if flen%2 != 0 {
		return false // didn't supply an even number of floats
	}
	for i := 0; i < len(f); i += 2 {
		if !EqualThreshold(f[i], f[i+1], Epsilon) {
			return false // at least 1 pair didn't equate
		}
	}
	return true // all pairs were equal when considering the Epsilon threshold.
}

// EqualThreshold is a utility function to compare floats.
// It's inspired by code from http://floating-point-gui.de/errors/comparison/
//
// This differs from Equal in that it lets you pass in your comparison
// threshold, so that you can adjust the comparison value to your specific
// needs
//
// Parameter a represents the measured or approximated value being tested
// Parameter b represents the true value being tested against.
// Parameter epsilon is the allowed error specified as a percentage of the value of b.
func EqualThreshold(a, b, epsilon float64) bool {
	// Handle the case of inf or shortcuts the loop when no significant error has accumulated
	if a == b {
		return true
	}

	diff := math.Abs(a - b)

	// For any value a or b below 1 we test against the real value of epsilon.
	// Also if the difference is already extremely small, we also do this.
	if a < 1 || b < 1 || diff < MinNormal {
		// Mathematically identical to using relative error diff / 1.0 and comparing that against epsilon.
		return diff < epsilon
	}

	// Calculate relative error where epsilon is defined as a percentage (percentage=100 * epsilon) of the target value b
	// We use this to allow bigger error margins when the a and b values are getting huge.
	// Rationale being that an error of e.g. 2 on a value of 1e300 is acceptable while
	// an error of 2 on a value of 200 would be unacceptable.
	return diff/math.Min(math.Abs(b), MaxValue) < epsilon
}
