# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

This project is pre-1.0. Its public API has no stability guarantee yet, and
breaking changes should be expected in any release until 1.0. There are no
tagged releases yet, so everything to date is listed under Unreleased.

## [Unreleased]

### Changed

- **Breaking:** replaced `commands.NetworkError` with three distinct error
  types — `commands.RequestEncodingError`, `commands.ProtocolError` and
  `commands.TransportError` — because the single type conflated socket
  faults, connection-lifecycle refusals and XML marshalling bugs into one
  case callers could not usefully switch on. See the "Error handling"
  section of the README for how to match on the new types.
- `commands.TransportError` carries `ConnectionUsable`, distinguishing a
  request refused before the socket was touched — an already-cancelled
  context, most commonly — from a failure that leaves the connection
  permanently unusable. Connections are shared, so without this a caller
  could discard a working transport over someone else's stale context. The
  zero value is the conservative answer, "do not reuse".
- Added `commands.ConnectionNotInitializedError` for a `Connection` used
  without an underlying socket. It is not a `TransportError`: no transport was
  ever established, so neither answer to "is this connection usable" is
  honest.
- **Breaking:** `connections.NewTLSConnection` is now
  `NewTLSConnection(ctx context.Context, address string, opts ...TLSOption)`,
  replacing `NewTLSConnection(address string, insecure bool)`. It is secure by
  default (certificate verification on, TLS 1.3 minimum) and supports client
  certificates, custom CA pools, SNI override and a configurable minimum TLS
  version through the new `TLSOption` functions: `WithRootCAs`,
  `WithRootCAFile`, `WithClientCertificate`, `WithClientCertFiles`,
  `WithMinVersion`, `WithServerName`, `WithInsecureSkipVerify`,
  `WithTLSConfig`.
- **Breaking:** `connections.NewUnixConnection` now takes a `context.Context`
  as its first parameter.
- `connections` now re-exports the connection sentinel errors
  (`ErrConnectionClosed`, `ErrConnectionBroken`, `ErrNoRawConn`) for use with
  `errors.Is`. They previously lived in an internal package external callers
  could not reach.
- **Breaking:** `context.Context` is now the first parameter of every method
  on `Client` and `Connection`, including `RawXML`. Previously roughly 105
  methods hardcoded `context.Background()` internally, so cancellation and
  deadlines set by a caller were silently ignored.
- **Breaking:** `implementation.New(conn)` replaces the previous
  `SetRawConn` setter; the underlying socket can now only be installed at
  construction. This closes a race where a socket swapped into a live
  connection could collide with an in-flight transaction.
- `Connection` is now safe for concurrent use: transactions are serialized
  internally. Previously, concurrent use of a single `Connection` could
  deadlock or cross-deliver one caller's response to another.
- A connection now has an explicit lifecycle (open, broken, closed). A
  transaction that fails in a way that leaves the message boundary
  ambiguous permanently retires the connection, instead of letting a later
  call read the tail of a previous, unrelated response.

### Removed

- **Breaking:** removed `commands.ErrorTypeNetwork`. It was unreachable —
  `NewGMPError` could never produce it — and existed only as dead code
  alongside the old `NetworkError` type.
- **Breaking:** removed `GMPError.Message` (a pre-rendered duplicate of
  `Error()`) and `GMPError.Unwrap` (it always returned `nil`).

### Added

- `commands.ProtocolError` retains a capped copy of the response that failed
  to decode, to make schema drift between client and server diagnosable.
- Enabled `staticcheck` in CI linting.

### Fixed

- Fixed a hang on any response whose length was an exact multiple of 1024
  bytes. Framing is now driven by the XML document structure (tracking
  element depth as the decoder consumes tokens) rather than by comparing
  read lengths against a fixed buffer size.
