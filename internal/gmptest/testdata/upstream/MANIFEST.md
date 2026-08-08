# Upstream fixture corpus

154 request/response examples transcribed from gvmd's GMP schema by
`go test ./internal/gmptest -run TestTranscribeUpstream -transcribe`.

## Provenance

- Source: `vendor/XML/commands/command_*.xml`, derived from gvmd's `GMP.xml.in`.
- Revision: GMP 22.7 (newest entry in vendor/XML/commands/changes.xml)
- License: AGPL-3.0-or-later, Copyright (C) Greenbone AG

This directory is **AGPL-derived** and is deliberately isolated so it can be
removed as a unit without touching any code — the fixture loader discovers
files by walking the tree and never names one. The nested `go.mod` keeps it out
of the published module zip.

## These are not wire bytes

The examples are documentation. Their indentation and self-closing forms are
the schema author's, not gvmd's, and marker resolution edits them further, so
every fixture here records `byte_faithful: false`. **No byte-exactness
assertion may run against this tree** — only path-level fidelity, which is
insensitive to formatting.

GMP.xml.in is a CMake template and its substitution markers appear inside the
example payloads, so the raw examples are not valid XML. Transcription resolves
them as a gvmd built with CVSS3 ratings and agents enabled, and each fixture's
`meta.json` records exactly which transforms were applied.
