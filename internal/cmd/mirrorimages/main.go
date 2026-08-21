// Command mirrorimages copies every image its lock file's previous run
// recorded into a registry this project controls, using Skopeo, records
// what it did back into that lock file, and rewrites a Docker Compose
// file's pins to match. See hack/mirror-images.sh, the thin wrapper around
// this command, and T-22 in docs/review-2026-08.md for why this exists.
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/brennoo/go-gmp/internal/mirror"
)

func main() {
	if err := run(context.Background(), os.Args[1:], mirror.Exec, time.Now, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "mirrorimages:", err)
		os.Exit(1)
	}
}

func checkSkopeoAvailable() error {
	if _, err := exec.LookPath("skopeo"); err != nil {
		return fmt.Errorf("skopeo is required to mirror images without pulling: %w", err)
	}
	return nil
}
