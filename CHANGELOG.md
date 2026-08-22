# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
This project is pre-1.0; see [CONTRIBUTING.md](CONTRIBUTING.md) for the versioning policy.

> [!NOTE]
> `v0.1.0`–`v0.3.0` were tagged before this file existed in its current form. Between `v0.3.0`
> (2025-10-26) and this release, the library was rewritten around a schema-driven code generator
> covering all commands, along with the connection-lifecycle, error-model, and API-shape work
> below — none of it tagged or changelogged at the time. This entry is the first to document that
> work; it is not a diff against `v0.3.0` line by line, but the accumulated state since.

## [0.4.0] - 2026-08-22

### Breaking

- Fields backed by a GMP `<alts>` enumeration (e.g. credential type, filter/keyword relation
  operators) are now named string types with generated constants instead of plain `string`
  (T-19). Existing string literals still assign directly; slice literals need a type conversion.
- Presence semantics changed across several response and request shapes so that "the schema
  marks this optional" and "gvmd can actually omit it" agree: count-wrapper fields (e.g.
  `FilterCount`) are now pointers that stay `nil` when the wire omits them instead of fabricating
  a zero value (T-23b); `nvt`/`refs` on results, notes, and overrides are now pointers for the
  same reason (T-26).
- `GetInfoResponse.Info2` was added to hold `get_info`'s response-wrapper element, previously
  merged with the per-item `Info` due to a wire name collision (T-26).
- `BuildFilterFromArgs` (`commands/filtering`) now returns `(string, error)` instead of `string` —
  an unrecognized filter argument type is a build-time silent no-op no longer; it's a runtime
  error (T-29).
- The `agent`/`agent_group` command surface (`GetAgents`, `ModifyAgents`, `DeleteAgents`, and their
  types) was removed. It was never reachable against the pinned gvmd image (`ENABLE_AGENTS` is
  compiled out) and carried known presence bugs. See T-32 for the evidence required to bring it
  back.
- `Client` gained `Close() error` (T-36).
- `GetAggregates` and `GetResourceNames` moved from the filter-variadic form to the typed
  `(ctx, *commands.GetX)` form — both have a schema-required attribute the filter-only form could
  never set, so the previous signature made that attribute unreachable (T-36).
- `GetRoles`, `GetPermissions`, `GetGroups`, `GetCredentials`, `GetSchedules`, `GetAlerts`,
  `GetReports`, `GetReportFormats`, `GetReportConfigs`, and `GetNotes` moved from the typed
  `(ctx, *commands.GetX)` form to the filter-variadic `(ctx, ...filtering.FilterArg)` form,
  matching every other list-fetching method on `Client`. The old typed call is still available as
  `Get*Raw(ctx, *commands.GetX)` (T-36a).

### Added

- `Client` gained all 25 `Get*Raw(ctx, *commands.GetX)` methods (10 new, 15 backfilled from the
  previously interface-invisible concrete implementation), so every command that needs a
  non-filter parameter (an ID, a format) remains reachable through the public interface (T-36a).
- A committed public API snapshot (`internal/apisnapshot/testdata/PUBLIC_API.txt`) with CI
  enforcement — any change to the exported surface of this module's packages now shows up as a
  reviewable diff (T-27).
- Structured error types (`GMPError`, `TransportError`, `ProtocolError`, `RequestEncodingError`,
  `ConnectionNotInitializedError`) replacing a single, ambiguous network-error type; see the
  README's "Error Handling" section for how to match on them.
- `NewTLSConnection` gained functional options (`WithRootCAs`, `WithRootCAFile`,
  `WithClientCertificate`, `WithClientCertFiles`, `WithMinVersion`, `WithServerName`,
  `WithInsecureSkipVerify`, `WithTLSConfig`), replacing a two-argument `(address, insecure bool)`
  signature that could not express client certificates, a custom CA pool, or SNI.
- `context.Context` is now the first parameter of every method on `Client` and `Connection`.
- `Connection` is now safe for concurrent use, with an explicit lifecycle (open, broken, closed).

### Fixed

- A response whose length was an exact multiple of 1024 bytes no longer hangs the read loop;
  framing now follows XML document structure instead of comparing against a fixed buffer size.
- Several `commands/pagination` iterator constructors dereferenced a count-wrapper field with no
  nil check, which would have panicked once those fields became pointers; all now nil-check.

### Known upstream behaviors

Confirmed against the pinned gvmd 26.36.1 runtime; not bugs in this library, but sharp edges
worth knowing before calling these commands:

- **`modify_permission`** crashes the gvmd worker process on essentially any real invocation
  (root cause: a column-name mismatch in gvmd's own SQL). Independently corroborated by
  [python-gvm#1047](https://github.com/greenbone/python-gvm/issues/1047), which reports the same
  "Error in SUBJECT" / crash symptom. `Subject` is also required by gvmd despite the schema
  marking it optional.
- **`modify_user`** segfaults gvmd if `Sources` is sent as explicit-empty
  (`<sources><source></source></sources>`); the typed API cannot construct this payload, so this
  only matters if you hand-build XML via `RawXML`.
- **Full-replace `modify_*` commands.** `modify_role`, `modify_group`, `modify_port_list`,
  `modify_asset`, `modify_note`, and `modify_override` clear one or more fields to empty when that
  field is *omitted* from the request, not just when it's sent explicit-empty — modifying one
  field alone can silently blank the others. `modify_override` is asymmetric: `Hosts`/`Port`/
  `Threat` full-replace, but `NewThreat` does not. See each command's doc comment for the exact
  field list.

[0.4.0]: https://github.com/brennoo/go-gmp/releases/tag/v0.4.0
