// Package dockerstack resolves the content-addressed image digest backing a
// running Docker Compose service. It exists so schema-extraction provenance
// (T-16b) and live test fixture metadata (T-16c) share one implementation
// instead of independently drifting.
//
// The package has no build tag and compiles anywhere; only callers that
// actually invoke ImageDigest need Docker present at runtime.
package dockerstack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/brennoo/go-gmp/internal/imageref"
)

// ImageDigest resolves the content-addressed digest reference
// (repo@sha256:...) of the image backing service's currently running
// container under the Compose project defined by composeFile. It fails
// closed — an error, never a tag or empty string — when the service has no
// running container, resolves to more than one, or the running image has no
// RepoDigests entry at all (for example, a locally built image that was
// never pulled or pushed by digest).
func ImageDigest(ctx context.Context, composeFile, service string) (string, error) {
	composeCmd, err := composeCommand(ctx)
	if err != nil {
		return "", err
	}

	containerID, err := runningContainer(ctx, composeCmd, composeFile, service)
	if err != nil {
		return "", err
	}

	configured, err := configuredImage(ctx, containerID)
	if err != nil {
		return "", err
	}
	repo := imageref.Repository(configured)

	imageID, err := inspectFormat(ctx, containerID, "{{.Image}}")
	if err != nil {
		return "", fmt.Errorf("dockerstack: resolving image id for container %s: %w", containerID, err)
	}

	digestsJSON, err := inspectFormat(ctx, imageID, "{{json .RepoDigests}}")
	if err != nil {
		return "", fmt.Errorf("dockerstack: resolving repo digests for image %s: %w", imageID, err)
	}
	digests, err := parseRepoDigests([]byte(digestsJSON))
	if err != nil {
		return "", fmt.Errorf("dockerstack: parsing repo digests for image %s: %w", imageID, err)
	}
	if len(digests) == 0 {
		return "", fmt.Errorf(
			"dockerstack: image %s (service %q) has no resolvable repo digest — "+
				"it looks locally built and was never pulled or pushed by digest, so it cannot be pinned",
			imageID, service)
	}

	digest, err := selectDigest(digests, repo)
	if err != nil {
		return "", fmt.Errorf("dockerstack: image %s (service %q, configured image %q): %w",
			imageID, service, configured, err)
	}
	return digest, nil
}

// configuredImage returns the image reference containerID was started from,
// via `docker inspect` on the container itself rather than `compose config
// --images <service>`: the latter resolves the full dependency subtree
// (every transitive depends_on's image too) with no single-service filter.
func configuredImage(ctx context.Context, containerID string) (string, error) {
	return inspectFormat(ctx, containerID, "{{.Config.Image}}")
}

// selectDigest returns the one entry in digests whose repository is repo,
// failing closed if none or more than one match — either is an ambiguous
// result RepoDigests' order cannot be trusted to break the tie for.
func selectDigest(digests []string, repo string) (string, error) {
	var matches []string
	for _, d := range digests {
		if imageref.Repository(d) == repo {
			matches = append(matches, d)
		}
	}
	switch len(matches) {
	case 0:
		return "", fmt.Errorf("no RepoDigests entry matches repository %q (have %v)", repo, digests)
	case 1:
		return matches[0], nil
	default:
		return "", fmt.Errorf("%d RepoDigests entries match repository %q, ambiguous: %v", len(matches), repo, matches)
	}
}

// composeCommand picks the Docker Compose invocation available on this
// machine, mirroring hack/compose-lib.sh's own v2-plugin-then-v1-binary
// fallback so this package and the shell scripts agree on what "Docker
// Compose" means. The returned slice is the full argv prefix, e.g.
// ["docker", "compose"] or ["docker-compose"].
func composeCommand(ctx context.Context) ([]string, error) {
	if exec.CommandContext(ctx, "docker", "compose", "version").Run() == nil {
		return []string{"docker", "compose"}, nil
	}
	if _, err := exec.LookPath("docker-compose"); err == nil {
		return []string{"docker-compose"}, nil
	}
	return nil, errors.New(`dockerstack: neither "docker compose" nor "docker-compose" is available`)
}

// psEntry is the subset of `docker compose ps --format json`'s per-container
// object this package needs.
type psEntry struct {
	Service string `json:"Service"`
	ID      string `json:"ID"`
	State   string `json:"State"`
}

// runningContainer returns the container ID for service's sole running
// container. hack/docker-compose.yml runs one container per service; a
// service resolving to zero or more than one is treated as a failure rather
// than guessed at.
func runningContainer(ctx context.Context, composeCmd []string, composeFile, service string) (string, error) {
	args := append(append([]string{}, composeCmd[1:]...), "-f", composeFile, "ps", "--format", "json", service)
	out, err := runCommand(ctx, composeCmd[0], args...)
	if err != nil {
		return "", fmt.Errorf("dockerstack: listing %q containers: %w", service, err)
	}

	entries, err := parsePSEntries(out)
	if err != nil {
		return "", fmt.Errorf("dockerstack: parsing compose ps output for %q: %w", service, err)
	}

	switch len(entries) {
	case 0:
		return "", fmt.Errorf("dockerstack: service %q has no running container under %s", service, composeFile)
	case 1:
		if entries[0].State != "running" {
			return "", fmt.Errorf("dockerstack: service %q container is %s, not running", service, entries[0].State)
		}
		return entries[0].ID, nil
	default:
		return "", fmt.Errorf("dockerstack: service %q resolved to %d containers, expected exactly one", service, len(entries))
	}
}

// parsePSEntries decodes `docker compose ps --format json`'s output, which is
// JSON Lines (one object per line) rather than a JSON array.
func parsePSEntries(out []byte) ([]psEntry, error) {
	var entries []psEntry
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		var e psEntry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// parseRepoDigests decodes `docker inspect --format '{{json .RepoDigests}}'`'s
// output.
func parseRepoDigests(out []byte) ([]string, error) {
	var digests []string
	if err := json.Unmarshal(bytes.TrimSpace(out), &digests); err != nil {
		return nil, err
	}
	return digests, nil
}

func inspectFormat(ctx context.Context, ref, format string) (string, error) {
	out, err := runCommand(ctx, "docker", "inspect", "--format", format, ref)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func runCommand(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...) //nolint:gosec // name/args are this package's own fixed docker/compose invocations
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s %s: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return stdout.Bytes(), nil
}
