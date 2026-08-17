#!/bin/sh
# Resolves gvmd's current upstream digest (never the pinned one — see
# docs/review-2026-08.md, T-17, for why diffing against the pinned image
# would never show drift), boots a candidate stack from it, and reports how
# its served schema differs from the committed testdata/schema snapshots.
#
# Requires `skopeo` on PATH. Does not tear the stack down; run
# hack/gvmd-down.sh --purge afterward.
#
# Usage: hack/schema-drift.sh [schemadrift flags, e.g. -gmp.insecure]
set -eu

HACK_DIR="$(cd "$(dirname "$0")" && pwd)"
repo_root="$(cd "$HACK_DIR/.." && pwd)"
cd "$repo_root"

SCHEMA_DRIFT_GVMD_IMAGE="$(go run ./internal/cmd/schemadrift -resolve-only)"
export SCHEMA_DRIFT_GVMD_IMAGE
log() { printf '%s\n' "$*" >&2; }
log "Candidate gvmd image: $SCHEMA_DRIFT_GVMD_IMAGE"

export COMPOSE_OVERRIDE_FILE="$HACK_DIR/docker-compose.schema-drift-override.yml"
hack/gvmd-up.sh

exec go run ./internal/cmd/schemadrift "$@"
