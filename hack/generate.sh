#!/bin/sh
# Regenerates the note-family command types (create_note, modify_note,
# get_notes) from the committed GMP schema snapshots under
# internal/codegen/testdata/schema, per docs/review-2026-08.md T-9.
#
# This is a manual, human-reviewed step, not something CI runs on its own:
# CI only checks that commands/ is not *stale* relative to what this script
# would produce (a separate `git diff --exit-code` step), it never runs
# gengmp and commits the result itself. Run this script yourself, then:
#
#   git diff -- commands/create_note.go commands/modify_note.go commands/get_notes.go commands/note_shared.go
#
# and review the diff before committing. It writes directly into commands/,
# overwriting those four files (the fourth, note_shared.go, holds types the
# generator deduplicated across more than one command in the family).
#
# Usage:
#   hack/generate.sh
set -eu

repo_root="$(cd "$(dirname "$0")/.." && pwd)"

cd "$repo_root"
go run ./internal/cmd/gengmp -schema internal/codegen/testdata/schema -out ./commands

echo "Review the diff before committing:" >&2
echo "  git diff -- commands/create_note.go commands/modify_note.go commands/get_notes.go commands/note_shared.go" >&2
