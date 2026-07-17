package nsort_test

import (
	"testing"

	"github.com/vibrantgio/seen/layer/internal/ordercheck"
	"github.com/vibrantgio/seen/layer/nsort"
)

// TestArtifact pins the nsort layer against the shared painter-order harness
// (layer/internal/ordercheck): interpenetration ("cross") and non-transitive
// occlusion ("cycle") must render correctly, which the view-dependent depth
// sort achieves by cutting a polygon exactly when it detects an occlusion
// cycle from the current eye.
func TestArtifact(t *testing.T) {
	ordercheck.Run(t, "nsort", nsort.NewLayerForScene)
}
