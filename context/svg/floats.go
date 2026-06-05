package svg

import "strconv"

// Ftoa formats v with the given decimal precision: 0 rounds to the nearest
// integer, N keeps N decimal places, and -1 emits the shortest string that
// round-trips. Rounding here — in serialization, not in the geometry — keeps
// SVG path data short without disturbing the sub-pixel points the raster layers
// draw, which is why precision lives on the SVG context, not on the scene.
func Ftoa(precision int, v float64) string {
	return strconv.FormatFloat(v, 'f', precision, 64)
}

func Fjoin(precision int, v ...float64) string {
	if len(v) == 0 {
		return ""
	}
	s := []byte(Ftoa(precision, v[0]))
	for _, f := range v[1:] {
		s = append(append(s, ' '), Ftoa(precision, f)...)
	}
	return string(s)
}
