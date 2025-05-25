package face

import (
	"strconv"

	"github.com/vibrantgio/seen/color"
)

type Faces []Face

// SetColorFrom sets a color on every face by reading it from
// the passed in colors.Source.
func (s Faces) SetColorFrom(source color.Source) (err error) {
	err = nil
	for i := range s {
		if err = s[i].SetFill(source.Read()); err != nil {
			return
		}
	}
	return
}

// SetFill applies the supplied fill Material to each face
func (s Faces) SetFill(value any) (err error) {
	for i := range s {
		if err = s[i].SetFill(value); err != nil {
			return
		}
	}
	return
}

// SetStroke applies the supplied stroke Material to each face
func (s Faces) SetStroke(value any) (err error) {
	for i := range s {
		if err = s[i].SetStroke(value); err != nil {
			return
		}
	}
	return
}

// SetStrokeWidth sets the supplied stroke-width option on each face
func (s Faces) SetStrokeWidth(value int) {
	for i := range s {
		s[i].Options["stroke-width"] = strconv.Itoa(value)
	}
}

// SetShowBackfaces will set the ShowBackfaces bool on the faces.
func (s Faces) SetShowBackfaces(value bool) {
	for i := range s {
		s[i].ShowBackfaces = value
	}
}
