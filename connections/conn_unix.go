package connections

import (
	"context"
	"net"

	"github.com/brennoo/go-gmp"
	"github.com/brennoo/go-gmp/connections/internal/implementation"
)

// NewUnixConnection returns an instance of `gmp.Connection`. The `socket` parameter refers to the
// file path of the Unix Socket where Openvas GVMD is listening. Ex: "/var/run/gvmd.sock".
//
// `ctx` governs the dial itself and expresses the deadline, which is why there is no options
// parameter: the remaining `net.Dialer` knobs are meaningless for a Unix domain socket.
//
// As with NewTLSConnection, a dial failure is returned raw, not wrapped in `*commands.TransportError`.
func NewUnixConnection(ctx context.Context, socket string) (gmp.Connection, error) {
	d := &net.Dialer{}
	conn, err := d.DialContext(ctx, "unix", socket)
	if err != nil {
		return nil, err
	}

	return implementation.New(conn), nil
}
