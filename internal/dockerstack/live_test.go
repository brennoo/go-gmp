//go:build integration

package dockerstack

import (
	"context"
	"regexp"
	"testing"
	"time"
)

var repoDigestPattern = regexp.MustCompile(`^[^@]+@sha256:[0-9a-f]{64}$`)

// TestImageDigestLive resolves gvmd's image digest against a real stack
// brought up by hack/gvmd-up.sh (or already running under the same compose
// project). It is the manual sanity check for ImageDigest itself, ahead of
// T-16c/T-16b wiring it into capture and extraction.
func TestImageDigestLive(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	digest, err := ImageDigest(ctx, "../../hack/docker-compose.yml", "gvmd")
	if err != nil {
		t.Fatalf("ImageDigest: %v", err)
	}
	if !repoDigestPattern.MatchString(digest) {
		t.Fatalf("ImageDigest = %q, want a repo@sha256:<64 hex> reference", digest)
	}
	t.Logf("gvmd image digest: %s", digest)
}
