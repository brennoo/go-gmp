// Package connections dials the transports go-gmp can speak to gvmd over —
// TLS and Unix domain sockets — or adapts an already-open net.Conn, and
// returns a [github.com/brennoo/go-gmp.Connection] in each case.
//
// [NewTLSConnection] dials a TCP+TLS connection and is secure by default: it
// verifies the server certificate and requires TLS 1.3. Both are meaningful
// restrictions to relax, so each has an explicit, separately named opt-out:
// [WithMinVersion] for a server limited to an older TLS version, and
// [WithInsecureSkipVerify] to disable certificate verification entirely,
// which is for local development only and must never be used against a
// production server. [WithRootCAs] and [WithRootCAFile] trust a custom
// certificate authority; [WithClientCertificate] and [WithClientCertFiles]
// present a client certificate for mutual TLS; [WithServerName] overrides SNI;
// [WithTLSConfig] is an escape hatch onto the underlying tls.Config.
//
//	conn, err := connections.NewTLSConnection(ctx, "gmp.example.com:9390",
//		connections.WithRootCAFile("/path/to/ca.pem"),
//	)
//
// [NewUnixConnection] dials a local Unix domain socket, as gvmd typically
// exposes when running alongside the caller:
//
//	conn, err := connections.NewUnixConnection(ctx, "/var/run/gvmd.sock")
//
// [New] wraps an already-established net.Conn directly, for transports these
// two constructors don't cover.
//
// In all three cases ctx governs only the dial: cancelling it aborts an
// in-progress connection attempt but has no effect on the returned
// [github.com/brennoo/go-gmp.Connection] afterward.
//
// # Concurrency and lifecycle
//
// A [github.com/brennoo/go-gmp.Connection] serializes transactions: it is safe
// for concurrent use by multiple goroutines, but not parallel — one
// transaction (request write plus response read) runs at a time, and other
// callers wait for admission. Admission is not FIFO; a waiting caller has no
// starvation bound. Use one connection per goroutine, each separately
// authenticated, if genuine parallelism is required.
//
// A connection moves through open, broken and closed states. Any transaction
// that fails in a way that leaves the message boundary ambiguous — a
// mid-write socket error, a response that never finished framing — poisons
// the connection permanently rather than risk the next call reading the tail
// of a previous response. Close is always safe to call and is idempotent.
//
// [ErrConnectionClosed], [ErrConnectionBroken] and [ErrNoRawConn] are
// re-exported sentinel values for use with errors.Is against the
// *commands.TransportError a poisoned or closed connection returns.
package connections
