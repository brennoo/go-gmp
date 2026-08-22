# Releasing

## Versioning policy

This project is pre-1.0 and uses [semver](https://semver.org/) with the pre-1.0 relaxation that
implies: breaking changes are allowed in a **minor** release, not just a major one. Patch releases
never break the public API.

- A breaking change ships in a minor bump (`0.4.0` → `0.5.0`), never a patch, and must have a
  `### Breaking` entry in [CHANGELOG.md](CHANGELOG.md) describing the change and the migration.
- A non-breaking addition or fix can ship in either a minor or a patch bump, whichever the
  maintainer judges fits the accumulated changes since the last tag.
- There is no deprecation cycle requirement pre-1.0. `AGENTS.md`'s policy stands: make the
  necessary breaking change now rather than carrying a compatibility shim for an API that hasn't
  stabilized yet.

**The enforcement hook.** `internal/apisnapshot/testdata/PUBLIC_API.txt` (T-27) is a committed,
CI-checked snapshot of every exported declaration this module ships. Any PR whose diff includes a
change to that file is a semver decision, not a mechanical fix — regenerating it with `-update` to
make a test pass, without stating in the PR description whether the diff is breaking and what
version bump it implies, is not an acceptable way to resolve a failing `TestPublicAPISnapshot`.

## What "supported" means

Per [SECURITY.md](SECURITY.md): the latest tag and the default branch (`main`) are supported.
Older tags receive no backports; if you need a fix on an older line, update to latest.

## Release procedure

1. Confirm the validation suite is green on the commit being released:
   ```
   go build ./...
   go test ./... -race -count=1
   go run ./internal/cmd/gengmp -check
   gofmt -l .
   golangci-lint run ./...
   go build -tags integration ./...
   go vet -tags integration ./...
   ```
2. Add a dated entry to `CHANGELOG.md` under a new version heading, moving anything relevant out
   of `[Unreleased]` if present. Breaking changes get their own `### Breaking` subsection.
3. Confirm `internal/apisnapshot/testdata/PUBLIC_API.txt` is up to date
   (`go test ./internal/apisnapshot -run TestPublicAPISnapshot`) and that its diff, if any, matches
   what the CHANGELOG's Breaking section claims.
4. Tag the merge commit on `main` with an annotated tag: `git tag -a vX.Y.Z -m "..."`.
5. Push the tag: `git push origin vX.Y.Z`.
6. Create a GitHub release from the tag. Release notes should point at the README's
   [Compatibility](README.md#compatibility) section rather than restating the compatibility
   statement, and summarize the CHANGELOG entry rather than duplicating it verbatim.
