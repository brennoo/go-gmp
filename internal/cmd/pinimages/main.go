package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "pinimages:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("pinimages", flag.ContinueOnError)
	file := fs.String("file", "hack/docker-compose.yml",
		"the Compose file to resolve and pin image: references in; rewritten in place")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("unexpected argument %q; pinimages takes flags only", fs.Arg(0))
	}

	if err := exec.CommandContext(ctx, "docker", "buildx", "version").Run(); err != nil {
		return fmt.Errorf("docker buildx is required to resolve image digests without pulling: %w", err)
	}

	data, err := os.ReadFile(*file)
	if err != nil {
		return err
	}

	resolve := cachedResolver(func(ref string) (string, error) { return resolveDigest(ctx, ref) })
	lines := strings.Split(string(data), "\n")
	changed := 0
	for i, line := range lines {
		out, err := refreshLine(line, resolve)
		if err != nil {
			return fmt.Errorf("%s line %d: %w", *file, i+1, err)
		}
		if out != line {
			fmt.Fprintf(os.Stderr, "%s\n", out)
			changed++
		}
		lines[i] = out
	}

	if err := os.WriteFile(*file, []byte(strings.Join(lines, "\n")), 0o644); err != nil { //nolint:gosec // a Compose file is not sensitive
		return err
	}
	fmt.Fprintf(os.Stderr, "wrote %s (%d line(s) changed)\n", *file, changed)
	return nil
}

// cachedResolver wraps resolve so a reference shared by multiple lines
// (e.g. several services all pinned to openvas-scanner:stable) is resolved
// only once per run — resolving each occurrence independently could pin a
// moving tag to different digests within the same run.
func cachedResolver(resolve func(ref string) (string, error)) func(ref string) (string, error) {
	cache := map[string]string{}
	return func(ref string) (string, error) {
		if digest, ok := cache[ref]; ok {
			return digest, nil
		}
		digest, err := resolve(ref)
		if err != nil {
			return "", err
		}
		cache[ref] = digest
		return digest, nil
	}
}

// resolveDigest returns ref's manifest digest (bare hex, no "sha256:"
// prefix) without pulling it, by hashing the raw manifest
// `docker buildx imagetools inspect` prints.
func resolveDigest(ctx context.Context, ref string) (string, error) {
	out, err := exec.CommandContext(ctx, "docker", "buildx", "imagetools", "inspect", "--raw", ref).Output() //nolint:gosec // ref comes from this repo's own compose file
	if err != nil {
		return "", fmt.Errorf("docker buildx imagetools inspect --raw %s: %w", ref, err)
	}
	sum := sha256.Sum256(out)
	return hex.EncodeToString(sum[:]), nil
}
