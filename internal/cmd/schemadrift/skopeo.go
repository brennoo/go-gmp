package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"

	"github.com/brennoo/go-gmp/internal/imageref"
)

// Runner executes name with args and returns its stdout, or an error
// embedding stderr — the seam the one skopeo invocation below goes through
// so it can be unit-tested against a fake.
type Runner func(ctx context.Context, name string, args ...string) ([]byte, error)

// Exec is the real Runner, shelling out via os/exec.
func Exec(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...) //nolint:gosec // name/args are this package's own fixed skopeo invocation
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s %s: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return []byte(stdout.String()), nil
}

// resolveDigest resolves ref (a mutable channel or anything Skopeo accepts)
// to its content digest, by hashing the raw manifest bytes rather than
// trusting a registry-reported digest header. github.com/brennoo/gvm-mirror's
// own ResolveDigest must compute it the same way, or comparing this against
// a lock entry's DestinationDigest silently stops detecting drift.
func resolveDigest(ctx context.Context, run Runner, ref string) (digestRef string, err error) {
	out, err := run(ctx, "skopeo", "inspect", "--raw", "docker://"+ref)
	if err != nil {
		return "", fmt.Errorf("resolving digest for %s: %w", ref, err)
	}
	sum := sha256.Sum256(out)
	return imageref.Repository(ref) + "@sha256:" + hex.EncodeToString(sum[:]), nil
}
