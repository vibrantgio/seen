package render

import (
	"strconv"
	"testing"
)

func TestFormatFloat(t *testing.T) {
	if strconv.FormatFloat(234.553343, 'g', -1, 64) != "234.553343" {
		t.Fail()
	}
}
