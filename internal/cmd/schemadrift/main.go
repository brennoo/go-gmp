// Command schemadrift compares a live gvmd's served GMP schema against the
// committed internal/codegen/testdata/schema snapshots. See docs/review-2026-08.md,
// T-17, and hack/schema-drift.sh for how it's meant to be invoked.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/brennoo/go-gmp/internal/dockerstack"
)

func main() {
	if err := run(context.Background(), os.Args[1:], Exec, dialLive, dockerstack.ImageDigest, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "schemadrift: %v\n", err)
		os.Exit(1)
	}
}
