package connections

import (
	"context"
	"crypto/tls"

	"github.com/brennoo/go-gmp"
	"github.com/brennoo/go-gmp/connections/internal/implementation"
)

// NewTLSConnection returns an instance of `gmp.Connection` that uses a TLS connection as underlying trasport to communicate with Openvas GVMD.
// The `address` parameter refers to the host and port where Openvas GVMD is listening.
// The `insecure` parameter allows the client to accept invalid certificates (ex: self-signed).
func NewTLSConnection(address string, insecure bool) (gmp.Connection, error) {
	d := &tls.Dialer{Config: &tls.Config{
		InsecureSkipVerify: insecure, // nolint:gosec
		MinVersion:         tls.VersionTLS13,
	}}
	conn, err := d.DialContext(context.Background(), "tcp", address)
	if err != nil {
		return nil, err
	}

	c := &implementation.Connection{}
	c.SetRawConn(conn)

	return c, nil
}
