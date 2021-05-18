package svg

import "strconv"

func Ftoa(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func Fjoin(v ...float64) string {
	if len(v) == 0 {
		return ""
	}
	s := []byte(Ftoa(v[0]))
	for _, f := range v[1:] {
		s = append(append(s, ' '), Ftoa(f)...)
	}
	return string(s)
}
