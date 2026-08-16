#!/bin/sh
# Resolves every image: reference in hack/docker-compose.yml to its current
# content digest and rewrites the file in place, replacing the tag (or lack
# of one) with a `@sha256:...` pin and keeping what it was pinned from as a
# trailing comment. Already-pinned lines are re-resolved too, so a moved tag
# (e.g. `:stable` pointing at a new build) is reflected on the next run
# instead of being frozen forever — see internal/cmd/pinimages.
#
# This is the only sanctioned way hack/docker-compose.yml's pins move:
# nothing re-runs this automatically or on a schedule. Run it, review the
# diff like any other change (a moved digest for a `:stable` tag is exactly
# the kind of thing that deserves a human look before it's committed), and
# commit deliberately.
#
# Usage: hack/pin-images.sh
set -eu

HACK_DIR="$(cd "$(dirname "$0")" && pwd)"
repo_root="$(cd "$HACK_DIR/.." && pwd)"

cd "$repo_root"
exec go run ./internal/cmd/pinimages -file "$HACK_DIR/docker-compose.yml"
