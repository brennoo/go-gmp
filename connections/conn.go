package connections

import (
	"net"

	"github.com/brennoo/go-gmp"
	"github.com/brennoo/go-gmp/connections/internal/implementation"
)

// New returns an instance of `gmp.Connection` that uses `conn` as underlying transport to communicate with Openvas GVMD.
func New(conn net.Conn) gmp.Connection {
	c := &implementation.Connection{}
	c.SetRawConn(conn)
	return c
}
