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
purpose. See `hack/README.md` for the certificate setup required before gvmd
will accept TLS connections (`gvm-manage-certs -af` on first run, and
`gvm-manage-certs -V` to confirm they're valid).

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
