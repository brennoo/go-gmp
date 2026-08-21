#!/bin/sh
# Mirrors every image hack/image-mirror.lock.json's previous run recorded
# into ghcr.io/brennoo/gvm-mirror using Skopeo, updates that lock file with
# the result, and rewrites hack/docker-compose.yml's pins to match — see
# internal/cmd/mirrorimages and T-22 in docs/review-2026-08.md for why this
# exists.
#
# Requires `skopeo` on PATH and a registry session already authenticated for
# push (e.g. `skopeo login ghcr.io` or `docker login ghcr.io`).
#
# This is the only sanctioned way the lock file and Compose pins move:
# nothing commits their changes automatically. Run it, review the diff, and
# commit deliberately.
#
# Usage: hack/mirror-images.sh [flags]
set -eu

HACK_DIR="$(cd "$(dirname "$0")" && pwd)"
repo_root="$(cd "$HACK_DIR/.." && pwd)"

cd "$repo_root"
exec go run ./internal/cmd/mirrorimages "$@"
