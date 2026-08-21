package connections

import "github.com/brennoo/go-gmp/connections/internal/implementation"

// Sentinel errors describing why a Connection refused to run a transaction.
// They are re-exported here, as var aliases rather than new errors.New
// values, so that error identity is preserved: errors.Is matches against
// either spelling, and external callers (who cannot import the internal
// implementation package) can still reach them.
var (
	// ErrConnectionClosed is returned once Close has been called.
	ErrConnectionClosed = implementation.ErrConnectionClosed
	// ErrConnectionBroken is returned once a transaction has left the stream
	// position ambiguous relative to message boundaries.
	ErrConnectionBroken = implementation.ErrConnectionBroken
	// ErrNoRawConn is returned when a Connection is used without a socket,
	// which only happens if one is built as a bare literal instead of through
	// a constructor.
	ErrNoRawConn = implementation.ErrNoRawConn
)
