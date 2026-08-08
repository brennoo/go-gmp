# Contributing

Thanks for looking at go-gmp. This is a pre-1.0 library — expect the API to
keep moving, and feel free to propose breaking changes when they make the
library correct or clearer rather than working around them.

## Building and testing

```bash
go build ./...
go test ./... -race -count=1
golangci-lint run ./...
```

- Always run tests with `-race`. This library's `Connection` type is
  explicitly designed to be safe for concurrent use, and races in that area
  are the kind of bug the race detector catches and a plain `go test` won't.
- Lint configuration lives in `.golangci.yml`. CI runs `golangci-lint run
  ./...` and requires zero issues — please run it locally before opening a
  PR rather than relying on CI to find problems for you.

## Manual testing against a real Greenbone instance

Unit tests cover the client and transport logic without a live server, but
changes that touch protocol behavior are worth exercising against a real
gvmd. `hack/docker-compose.yml` brings up a Greenbone Community Edition stack
(`gvmd`, the vulnerability feed containers, PostgreSQL, etc.) for that
purpose. Bringing it up from cold takes 20-45 minutes (the SCAP/CERT/NVT feed
sync), so budget for that before you need it.

```bash
hack/gvmd-up.sh                    # generate certs, start gvmd, create admin
SCAN_DATA=1 hack/gvmd-up.sh        # also start vulnerability-tests, for get_nvts
hack/gvmd-down.sh                  # stop, keeping the feed sync for next time
hack/gvmd-down.sh --purge          # stop and delete volumes, losing the sync
```

Both scripts work from any directory, and detect whether you have the
`docker compose` plugin or the standalone `docker-compose` binary.
`hack/gvmd-up.sh` prints, at the end, the server
certificate's CN (needed for TLS verification — see `hack/README.md` for
why 127.0.0.1 doesn't work) and the exact command to run the integration
suite against the stack it just started. See `hack/README.md` for the
manual version of every step, useful when something in the script fails and
you need to run the underlying `docker compose` / `gvmd` commands by hand.

### Running the integration suite

Code behind `//go:build integration` (currently `internal/gmptest/live.go`
and `live_test.go`) only compiles with the tag; a plain `go build ./...` or
`go test ./...` never touches it, and neither does `golangci-lint run
./...` without `run.build-tags: [integration]` in `.golangci.yml` (already
set). CI separately runs `go build -tags integration ./...` and `go vet
-tags integration ./...` on every PR specifically so this code cannot bit-rot
between the rare occasions someone actually runs it against a live server.

```bash
go build -tags integration ./...
go vet -tags integration ./...
go test -tags integration ./internal/gmptest/... \
  -gmp.addr=127.0.0.1:9390 -gmp.user=admin -gmp.pass=<password> \
  -gmp.servername=<CN from hack/gvmd-up.sh's output>
```

`requireLive` in `live_test.go` polls a readiness ladder (TLS dial +
`get_version`, then `authenticate`, then `get_feeds` showing every feed
synced) for up to an hour on first use, which is what makes the suite safe to
point at a stack that just cold-booted. It runs at most once per test binary,
and only when a test that needs a live server actually runs — so a selection
that skips (`-run TestCaptureLive` with no `-capture`) costs nothing.

### Capturing fixtures from a live gvmd

`TestCaptureLive` records the `testdata/live/` half of the fixture corpus. It
is gated behind `-capture` on top of the build tag, because unlike the rest of
the suite it creates and deletes resources on the server and writes files into
the repository.

```bash
go test -tags integration ./internal/gmptest -run TestCaptureLive -capture \
  -gmp.addr=127.0.0.1:9390 -gmp.pass=<password> -gmp.servername=<CN> \
  -timeout 60m
```

Everything it creates is named `go-gmp-fixture-<runid>-…` and torn down with
`ultimate="1"`, so a run that dies before its teardowns leaves behind
something identifiable you can sweep by name. One command failing never aborts
the run; the end-of-run summary lists what was captured, what failed and what
was skipped, and it is worth reading in full.

**Redaction is deny-by-default.** `internal/gmptest/redact.go` holds a
per-command list of canonical paths whose values are rewritten before anything
touches disk, plus an explicit set of commands considered secret-bearing —
capture *refuses* to write a fixture for a command in that set that has no
spec, because an absent spec must never be read as "nothing to redact". After
redacting, every byte about to be written is swept for the harness's own
password and any secret the run generated; a hit aborts the whole run. Read
the diff before committing, and grep it for `[REDACTED-BY-GO-GMP-CAPTURE]` to
see what was rewritten.

Captured `NNN_request.xml` and `NNN_response.xml` are evidence, never
regenerated: `WriteGolden` refuses to write them, so no `-update` path can
rewrite a fixture to match a regression.

### Integration in CI

`.github/workflows/integration.yml` runs the tagged suite (the readiness
ladder and the read-only tests, never `-capture`) against a real stack. It is
`workflow_dispatch`-only and prints `df -h` and timings, because neither the
disk cost nor the wall-clock cost of a cold Greenbone boot on a GitHub-hosted
runner has been measured yet — disk is the binding constraint at roughly 14 GB
free.

Do not add a `schedule:` trigger before a manual run has produced those
numbers, and know what you are signing up for if you do: **GitHub disables
`schedule:` triggers after 60 days of repository inactivity**, silently and
with no notification. A scheduled integration job fails open, not closed, so
it can never be the only thing keeping this code alive — that job belongs to
the `go build -tags integration` / `go vet -tags integration` steps that run
on every PR.

The job never commits fixtures and has read-only `contents` permission.
Capture stays a deliberate human operation: an auto-committing bot is exactly
how an unredacted credential reaches git history.

## Code style

- Match the style of the surrounding code rather than introducing a new
  pattern for the same problem — this codebase leans on generics (see
  `commands/pagination`) and functional options (see `commands/filtering` and
  `connections`) in specific, established places.
- Run `gofmt` before committing; CI checks that `gofmt -l .` is empty.
- Comments that read as full sentences must end in a period — the `godot`
  linter enforces this.

## Tests

Any change to behavior — not just new features, but bug fixes and refactors
too — should come with a test that would have failed before the change.
Table-driven tests are the norm here; look at an existing `_test.go` file in
the package you're touching for the expected shape before adding a new one.

## Tracking work

Outstanding work is tracked in
[GitHub Issues](https://github.com/brennoo/go-gmp/issues). Before starting
something substantial, check whether an issue already covers it, and comment
there so effort isn't duplicated. If nothing covers it, open an issue first —
for anything beyond a small fix it's worth agreeing on the approach before the
code exists.

Please don't use issues to report security vulnerabilities; see
[SECURITY.md](SECURITY.md) for the private disclosure path.

## Background reading

[`docs/review-2026-08.md`](docs/review-2026-08.md) is a point-in-time design
review from August 2026. It is a historical record rather than a tracker — its
statuses are frozen — but it explains the reasoning behind several decisions
that are otherwise hard to reconstruct from the code, particularly around
connection lifecycle, cancellation, and how errors are classified. Worth
skimming before changing anything in `connections/`.
