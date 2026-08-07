package connections

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
)

// TLSOption represents a functional option that configures the `tls.Config` used by `NewTLSConnection`.
type TLSOption func(*tlsSettings)

// tlsSettings accumulates the effect of applied TLSOption values before dialling.
// A functional option cannot return an error, so options that can fail (ex: reading a
// file from disk) record the first error here instead. The first error wins: once set,
// later options must not overwrite it.
type tlsSettings struct {
	config *tls.Config
	err    error
}

func newTLSSettings() *tlsSettings {
	return &tlsSettings{
		config: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}
}

// setErr records err if this is the first error seen, preserving whichever option failed first.
func (s *tlsSettings) setErr(err error) {
	if s.err == nil {
		s.err = err
	}
}

// WithTLSConfig replaces the base `tls.Config` wholesale. It is an escape hatch for callers who
// need full control over the config. If combined with other options, put `WithTLSConfig` first:
// later options mutate whatever config is current, so applying it after other options discards them.
//
// cfg is cloned, so later mutations by the caller do not affect the connection. A nil cfg is an
// error rather than a silent fallback to the defaults, which would otherwise be indistinguishable
// from not passing the option at all.
func WithTLSConfig(cfg *tls.Config) TLSOption {
	return func(s *tlsSettings) {
		if cfg == nil {
			s.setErr(errors.New("WithTLSConfig: nil *tls.Config"))
			return
		}
		s.config = cfg.Clone()
	}
}

// WithInsecureSkipVerify disables certificate verification. This makes the connection vulnerable
// to man-in-the-middle attacks and must not be used in production; prefer WithRootCAFile or
// WithRootCAs with a trusted certificate authority.
func WithInsecureSkipVerify() TLSOption {
	return func(s *tlsSettings) {
		s.config.InsecureSkipVerify = true // nolint:gosec
	}
}

// WithRootCAs sets the pool of root certificate authorities used to verify the server certificate.
// A nil pool means the host's system roots, matching `tls.Config.RootCAs`.
//
// The trust set is snapshotted when the option is created, not when it is applied, so mutating pool
// at any point afterwards has no effect. `x509.CertPool` is mutable, and cloning only on apply would
// leave the window between creating the option and passing it to NewTLSConnection.
func WithRootCAs(pool *x509.CertPool) TLSOption {
	if pool == nil {
		return func(s *tlsSettings) { s.config.RootCAs = nil }
	}
	snapshot := pool.Clone()
	// Clone again per application so two connections built from one option never share a pool.
	return func(s *tlsSettings) { s.config.RootCAs = snapshot.Clone() }
}

// WithRootCAFile reads the PEM-encoded certificate(s) at path and uses them as the pool of root
// certificate authorities used to verify the server certificate. It is a convenience wrapper
// around building an `x509.CertPool` by hand.
func WithRootCAFile(path string) TLSOption {
	return func(s *tlsSettings) {
		pem, err := os.ReadFile(path) // nolint:gosec // path is supplied by the caller of this option by design
		if err != nil {
			s.setErr(fmt.Errorf("read root CA file %q: %w", path, err))
			return
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pem) {
			s.setErr(fmt.Errorf("read root CA file %q: no certificates found", path))
			return
		}

		s.config.RootCAs = pool
	}
}

// WithClientCertificate sets the client certificate presented during the TLS handshake, for
// servers that require mutual TLS.
func WithClientCertificate(cert tls.Certificate) TLSOption {
	return func(s *tlsSettings) {
		s.config.Certificates = []tls.Certificate{cert}
	}
}

// WithClientCertFiles reads the certificate and private key at certFile and keyFile and sets them
// as the client certificate presented during the TLS handshake. It is a convenience wrapper around
// `tls.LoadX509KeyPair`.
func WithClientCertFiles(certFile, keyFile string) TLSOption {
	return func(s *tlsSettings) {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			s.setErr(fmt.Errorf("load client certificate %q / %q: %w", certFile, keyFile, err))
			return
		}

		s.config.Certificates = []tls.Certificate{cert}
	}
}

// WithMinVersion sets the minimum TLS version accepted, e.g. `tls.VersionTLS12`.
func WithMinVersion(v uint16) TLSOption {
	return func(s *tlsSettings) {
		s.config.MinVersion = v
	}
}

// WithServerName sets the server name used for SNI and certificate verification, overriding the
// host derived from the dial address.
func WithServerName(name string) TLSOption {
	return func(s *tlsSettings) {
		s.config.ServerName = name
	}
}
