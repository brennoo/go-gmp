package connections

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/brennoo/go-gmp"
	"github.com/brennoo/go-gmp/connections/internal/implementation"
)

// NewTLSConnection returns an instance of `gmp.Connection` that uses a TLS connection as underlying
// transport to communicate with Openvas GVMD. The `address` parameter refers to the host and port
// where Openvas GVMD is listening.
//
// By default the connection verifies the server certificate and requires TLS 1.3, which is
// stricter than the Go standard library default (TLS 1.2). This is the correct posture for a
// security library: callers that need to reach a server limited to TLS 1.2 must opt in explicitly
// with `WithMinVersion(tls.VersionTLS12)`. Certificate verification can be disabled for development
// with `WithInsecureSkipVerify`, and a custom certificate authority can be trusted with
// `WithRootCAs` or `WithRootCAFile`.
//
// `ctx` governs the dial itself; cancelling it aborts an in-progress connection attempt.
//
// Errors here are configuration or raw dial/TLS handshake errors, deliberately not wrapped in
// `*commands.TransportError` — no connection exists yet, so there is nothing to decide about retiring.
func NewTLSConnection(ctx context.Context, address string, opts ...TLSOption) (gmp.Connection, error) {
	settings := newTLSSettings()
	for i, opt := range opts {
		if opt == nil {
			return nil, fmt.Errorf("nil TLSOption at index %d", i)
		}
		// Stop at the first failure: later options would do file I/O whose
		// result is already discarded.
		if settings.err != nil {
			return nil, settings.err
		}
		opt(settings)
	}
	if settings.err != nil {
		return nil, settings.err
	}

	d := &tls.Dialer{Config: settings.config}
	conn, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, err
	}

	return implementation.New(conn), nil
}
