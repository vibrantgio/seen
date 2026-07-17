package bsort_test

import (
	"testing"

	"github.com/vibrantgio/seen/layer/bsort"
	"github.com/vibrantgio/seen/layer/internal/ordercheck"
)

// TestArtifact pins the splitting bsort layer against the shared
// painter-order harness (layer/internal/ordercheck): interpenetration
// ("cross") and non-transitive occlusion ("cycle") must render correctly,
// which the view-independent BSP achieves by splitting straddling polygons.
//
// The no-split mode (NewNoSplitLayerForScene) is deliberately NOT run here:
// keeping straddlers whole is its documented trade-off, and these scenes are
// exactly the ones it gives up on.
func TestArtifact(t *testing.T) {
	ordercheck.Run(t, "bsort", bsort.NewLayerForScene)
}
