#!/bin/sh
# Regenerates every command family declared in internal/codegen/manifest.json
# from the committed GMP schema snapshots under internal/codegen/testdata/schema,
# per docs/review-2026-08.md T-9.
#
# This is a manual, human-reviewed step, not something CI runs on its own: CI
# only checks that commands/ is not *stale* relative to what this script would
# produce (`go run ./internal/cmd/gengmp -check`, which writes nothing), it
# never runs gengmp and commits the result itself. Run this script yourself,
# review `git diff -- commands/`, then commit.
#
# It writes directly into commands/, overwriting whatever files the manifest's
# families currently generate (gengmp prints the exact list and byte counts
# below). It does not delete files left over from a family that was since
# removed from the manifest — `gengmp -check` reports those as obsolete so
# they aren't missed.
#
# Usage:
#   hack/generate.sh
set -eu

repo_root="$(cd "$(dirname "$0")/.." && pwd)"

cd "$repo_root"
go run ./internal/cmd/gengmp \
  -manifest internal/codegen/manifest.json \
  -schema internal/codegen/testdata/schema \
  -out ./commands

echo "Review the diff before committing:" >&2
echo "  git diff -- commands/" >&2
