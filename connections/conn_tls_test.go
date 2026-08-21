package connections

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/xml"
	"errors"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

var certPem = []byte(`-----BEGIN CERTIFICATE-----
MIIDCTCCAfGgAwIBAgIUaVye29+L133ysotFOBpR1o5Nvm0wDQYJKoZIhvcNAQEL
BQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTI1MDYxMDE1NTE1OVoXDTM1MDYw
ODE1NTE1OVowFDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEAl5fdmnLSWpzkvxtxJezU+f4+OCjZqrWLG/vb6KuYfEaX
Hms0q9iAeC1iGJ6U6HZ+SR+aekkFHEuVN98/W+UnhuQaS9hJvj97YmZNZNSf3XhP
1vnB1g2xNjCnsa8RHgz5x4QfZ7aql59J3+kqvzi24Xkz7jXbeiynC+URZv0DJvu7
dwp115gBbLU0NpPQjzIS4LgL8kFuSOxQrRfWoKW5AEugMMdHR+7z7bXtfmNMrBK7
XgIx+8kHvUCd2wxDMs07b2EmFtgtNfWWfu0WOnSQnnYNkjzNqCvaYtOGjFyxioGQ
ROYCS9rqR34ibdhv39zMrrrnIL9mJGsfiXC7F0pCGQIDAQABo1MwUTAdBgNVHQ4E
FgQUnZdYoEeBStYYGYdbOFLglsv18pwwHwYDVR0jBBgwFoAUnZdYoEeBStYYGYdb
OFLglsv18pwwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAWw4X
omJHfhJf9x2y3N3E9bfFNSG8QXe8OrfekM1UY/yP9mvxIuBedhQZxM5J/SLdjpeB
z9wYujFj/T4UJ79eXVsyxdtneGAQBEmYhOfxVUDiG3yrAqRqunIYLcFi7nQdfDYi
nfYxvwXv/VVAzF7Em3Z3UQp0NNksMzFEcRPPgXUZ7IN7dmMRZDod8uJm+JdZkGnw
D7bAJgViDuRUAuDUltBeV27SHv+WOirg5tVop/0+aiH/seUJ3z08JU2tbVgLojAI
Xap92QLN4NutCz8PXcdoo8nnnOrKRDcE/r966hQ96bWiEqFccPk63eFbU1ig+GEF
OWtJ+FTgdpf2SsIEnA==
-----END CERTIFICATE-----`)

var keyPem = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCXl92actJanOS/
G3El7NT5/j44KNmqtYsb+9voq5h8RpceazSr2IB4LWIYnpTodn5JH5p6SQUcS5U3
3z9b5SeG5BpL2Em+P3tiZk1k1J/deE/W+cHWDbE2MKexrxEeDPnHhB9ntqqXn0nf
6Sq/OLbheTPuNdt6LKcL5RFm/QMm+7t3CnXXmAFstTQ2k9CPMhLguAvyQW5I7FCt
F9agpbkAS6Awx0dH7vPtte1+Y0ysErteAjH7yQe9QJ3bDEMyzTtvYSYW2C019ZZ+
7RY6dJCedg2SPM2oK9pi04aMXLGKgZBE5gJL2upHfiJt2G/f3Myuuucgv2Ykax+J
cLsXSkIZAgMBAAECggEAClvOAH0Kmzb/2YA4S9liVz/UhBVkhEn5+rxL01sQZSOU
tXcpZvm5E500SA7kCXt5VylR7ru5He0bQfFfHLosqIbDNj4OEfI4DeoELVjrIKDV
QdEYufX/Dz8lYlqeXpGP6t6AH9/nScuQany+F0l6k68qyYMEH5XYB1PMi99LnRj5
D3X68DvOnXoMv4ld50Kb8+LKQXFj4PgXY67T4QJi/F8i0Ev/26g6cn896RKf8Pdr
Qbjj+Dq1yDc0eDGWdF7nGBLr8UEov6YCCVlVB2s9pfaymTvvATVk1LRQmZSIPtng
ZOWCXO3f6ldKyVdCM9KSsu1U47wKab5u1LyTOGQbgQKBgQDLWEshmMLTKYzJz3qm
lsdDueKVZwbAejSvWyj5pS8FykeT6NM65SSJNRHkuTGmYJUpGSd6NXVCfZTSntc4
twEqUQgtEDu7ThYhH4qnEkEhv4Z1mdVWJ4pttY8W/K4zGZ3KIBrvPtod+NX2oOmJ
o9zLYLjfWKKj/vHq6/3fqe9E2QKBgQC+2PXtdul7pdRxABgnFy+H8D3TBhjyhzDq
CGRVied1ajqv8ub2hormf+S8esqbyRrORSG8Ec23rR1T4UwG9Lz910aqKAhbvJWM
MI83P9oUfsr5zFP1tw7ALN2xSv12mgqs2MiLWBrGNPMFiP1jZJX6C3OdIuFenqwm
VA3OS5+fQQKBgDy0jjFodnktxU8WNqp3BuPsIX2ytSxPD11uXIwyyRf0bGv2rOC/
OaoQgtgChZUC7mTHkqaEQQ5piT1bmCrhT2K/sC4r03k2dZrL3MvFwlX95HlFRJAk
28mm0yHsq2Cr/BQ0g3X7EIi0GcQ4A6BIoAjDRk2/G0bUPIqCi1bzV2/5AoGAZYRp
RPik+BXc7IL9VgXaTbg4WsD6kj/hgAwYuvuROH7aWu+ddfNZfYT4el3i2n1eGezi
JfVedV7Lo+vLEkQrJ+fUefyzfYDSF/FYrS305kQP2lKbMrA+U3FKN93I98a2+PyT
qLOkAlz2DLKy/qTptklTJXoBYeCLeBPsCBln/kECgYB7x6NtIuk3O1lPvOAx7Fr3
wO1lY1EEk1GLQ3vAtRPmU5SCvyiDoqzMvAtXX9wwLJXtmbUBb+V/Zf9uDAtUEnP+
KQQ5Z7kaD8C+bxu4/d6psEsc/69SQyMZH3stt0DE8zzB3oW4hXgzlOySN/AttlJ4
MhCg59fKjTw+CzRzSc5b2A==
-----END PRIVATE KEY-----`)

func newLocalListener(t testing.TB) net.Listener {
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	})
	if err != nil {
		ln, err = tls.Listen("tcp6", "[::1]:0", nil)
	}
	if err != nil {
		t.Fatal(err)
	}
	return ln
}

func newLocalListenerProtocolNotSupported(t testing.TB) net.Listener {
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
		MaxVersion:         tls.VersionTLS10,
	})
	if err != nil {
		ln, err = tls.Listen("tcp6", "[::1]:0", nil)
	}
	if err != nil {
		t.Fatal(err)
	}
	return ln
}

// newLocalListenerTLS12Only starts a listener that only accepts TLS 1.2, to prove WithMinVersion
// can reach it and the TLS 1.3 default cannot.
func newLocalListenerTLS12Only(t testing.TB) net.Listener {
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		t.Fatal(err)
	}

	ln, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS12,
	})
	if err != nil {
		t.Fatal(err)
	}
	return ln
}

// newLocalListenerRequireClientCert starts a listener that requires the client to present a
// certificate, to prove WithClientCertificate/WithClientCertFiles work.
func newLocalListenerRequireClientCert(t testing.TB) net.Listener {
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		t.Fatal(err)
	}

	ln, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAnyClientCert,
	})
	if err != nil {
		t.Fatal(err)
	}
	return ln
}

// generateSelfSignedCert creates a self-signed certificate with SANs, unlike the checked-in
// certPem/keyPem fixture (CN only), since modern Go no longer falls back to CN for verification.
func generateSelfSignedCert(t testing.TB, dnsNames ...string) tls.Certificate {
	t.Helper()

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: dnsNames[0]},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:              dnsNames,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatal(err)
	}

	leaf, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}

	return tls.Certificate{
		Certificate: [][]byte{der},
		PrivateKey:  priv,
		Leaf:        leaf,
	}
}

// serveEcho accepts a single connection on ln and echoes back whatever it reads.
func serveEcho(ln net.Listener, errCh chan<- error) {
	conn, err := ln.Accept()
	if err != nil {
		errCh <- err
		return
	}
	buf := make([]byte, 150000)
	nRead, _ := conn.Read(buf)
	conn.Write(buf[:nRead])
}

func TestNewTLSConnection(t *testing.T) {
	ln := newLocalListener(t)
	errCh := make(chan error, 1)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				errCh <- err
				return
			}
			buf2 := make([]byte, 150000)
			nRead, _ := conn.Read(buf2)
			conn.Write(buf2[:nRead])
		}
	}()

	conn, err := NewTLSConnection(context.Background(), ln.Addr().String(), WithInsecureSkipVerify())
	if err != nil {
		t.Fatalf("Unexpected error during NewTLSConnection: %s", err)
	}

	expectedValue := "0123456789"
	cmd := &struct {
		XMLName xml.Name `xml:"cmd"`
		Foo     string   `xml:"foo"`
	}{
		Foo: expectedValue,
	}
	response := &struct {
		Foo string `xml:"foo"`
	}{}

	err = conn.Execute(context.Background(), cmd, response)
	if err != nil {
		t.Fatalf("Unexpected error during Execute: %s", err)
	}

	if response.Foo != expectedValue {
		t.Fatalf("Unexpected response value: %s\nExpected: %s", response.Foo, expectedValue)
	}

	select {
	case err := <-errCh:
		t.Fatalf("Unexpected error during Accept: %s", err)
	default:
	}
}

func TestNewTLSConnectionFail(t *testing.T) {
	ln := newLocalListenerProtocolNotSupported(t)
	errCh := make(chan error, 1)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				errCh <- err
				return
			}
			buf2 := make([]byte, 150000)
			nRead, _ := conn.Read(buf2)
			conn.Write(buf2[:nRead])
		}
	}()

	_, err := NewTLSConnection(context.Background(), ln.Addr().String(), WithInsecureSkipVerify())

	expectedError := "remote error: tls: protocol version not supported"
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Unexpected error during Execute.\nExpected: %s\n     Got: %s", expectedError, err)
	}

	select {
	case err := <-errCh:
		t.Fatalf("Unexpected error during Accept: %s", err)
	default:
	}
}

// TestNewTLSConnectionCustomRootCA proves that secure-by-default actually works: verification
// succeeds against a custom root CA pool without WithInsecureSkipVerify.
func TestNewTLSConnectionCustomRootCA(t *testing.T) {
	cert := generateSelfSignedCert(t, "localhost")

	ln, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{Certificates: []tls.Certificate{cert}})
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	errCh := make(chan error, 1)
	go serveEcho(ln, errCh)

	pool := x509.NewCertPool()
	pool.AddCert(cert.Leaf)

	conn, err := NewTLSConnection(context.Background(), ln.Addr().String(), WithRootCAs(pool), WithServerName("localhost"))
	if err != nil {
		t.Fatalf("Unexpected error during NewTLSConnection: %s", err)
	}
	defer conn.Close() //nolint:errcheck

	expectedValue := "0123456789"
	cmd := &struct {
		XMLName xml.Name `xml:"cmd"`
		Foo     string   `xml:"foo"`
	}{
		Foo: expectedValue,
	}
	response := &struct {
		Foo string `xml:"foo"`
	}{}

	if err := conn.Execute(context.Background(), cmd, response); err != nil {
		t.Fatalf("Unexpected error during Execute: %s", err)
	}

	if response.Foo != expectedValue {
		t.Fatalf("Unexpected response value: %s\nExpected: %s", response.Foo, expectedValue)
	}

	select {
	case err := <-errCh:
		t.Fatalf("Unexpected error during Accept: %s", err)
	default:
	}
}

// TestNewTLSConnectionMinVersion proves that WithMinVersion(tls.VersionTLS12) can reach a
// TLS-1.2-only server, and that the TLS 1.3 default cannot.
func TestNewTLSConnectionMinVersion(t *testing.T) {
	ln := newLocalListenerTLS12Only(t)
	defer ln.Close()
	errCh := make(chan error, 1)
	go serveEcho(ln, errCh)

	conn, err := NewTLSConnection(context.Background(), ln.Addr().String(),
		WithInsecureSkipVerify(),
		WithMinVersion(tls.VersionTLS12),
	)
	if err != nil {
		t.Fatalf("Unexpected error connecting with WithMinVersion(tls.VersionTLS12): %s", err)
	}
	defer conn.Close() //nolint:errcheck

	select {
	case err := <-errCh:
		t.Fatalf("Unexpected error during Accept: %s", err)
	default:
	}
}

func TestNewTLSConnectionMinVersionDefaultRefused(t *testing.T) {
	ln := newLocalListenerTLS12Only(t)
	defer ln.Close()
	errCh := make(chan error, 1)
	go serveEcho(ln, errCh)

	_, err := NewTLSConnection(context.Background(), ln.Addr().String(), WithInsecureSkipVerify())
	if err == nil {
		t.Fatal("expected error connecting to a TLS-1.2-only server with the default (TLS 1.3) minimum version")
	}
}

// TestNewTLSConnectionClientCertificate proves WithClientCertificate presents a client
// certificate that the server actually receives.
func TestNewTLSConnectionClientCertificate(t *testing.T) {
	clientCert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		t.Fatalf("Unexpected error building client certificate: %s", err)
	}

	assertServerReceivesClientCert(t, WithClientCertificate(clientCert))
}

// TestNewTLSConnectionClientCertFiles proves WithClientCertFiles reads a certificate/key pair
// from disk and presents it, and that the server actually receives it.
func TestNewTLSConnectionClientCertFiles(t *testing.T) {
	dir := t.TempDir()
	certFile := dir + "/client-cert.pem"
	keyFile := dir + "/client-key.pem"
	if err := os.WriteFile(certFile, certPem, 0o600); err != nil {
		t.Fatalf("Unexpected error writing client cert file: %s", err)
	}
	if err := os.WriteFile(keyFile, keyPem, 0o600); err != nil {
		t.Fatalf("Unexpected error writing client key file: %s", err)
	}

	assertServerReceivesClientCert(t, WithClientCertFiles(certFile, keyFile))
}

func assertServerReceivesClientCert(t *testing.T, clientCertOpt TLSOption) {
	t.Helper()

	ln := newLocalListenerRequireClientCert(t)
	defer ln.Close()

	gotClientCert := make(chan bool, 1)
	errCh := make(chan error, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			errCh <- err
			return
		}
		tlsConn, ok := conn.(*tls.Conn)
		if !ok {
			errCh <- errors.New("accepted connection is not a *tls.Conn")
			return
		}
		if err := tlsConn.HandshakeContext(context.Background()); err != nil {
			errCh <- err
			return
		}
		gotClientCert <- len(tlsConn.ConnectionState().PeerCertificates) > 0

		buf := make([]byte, 150000)
		nRead, _ := tlsConn.Read(buf)
		tlsConn.Write(buf[:nRead])
	}()

	conn, err := NewTLSConnection(context.Background(), ln.Addr().String(), WithInsecureSkipVerify(), clientCertOpt)
	if err != nil {
		t.Fatalf("Unexpected error during NewTLSConnection: %s", err)
	}
	defer conn.Close() //nolint:errcheck

	select {
	case got := <-gotClientCert:
		if !got {
			t.Fatal("server did not receive a client certificate")
		}
	case err := <-errCh:
		t.Fatalf("Unexpected error during server handshake: %s", err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for server handshake")
	}
}

// TestNewTLSConnectionRootCAFileMissing proves a failing option returns an error without ever
// dialling: the address has nothing listening, so a "connection refused" error rather than the
// expected root-CA-file error would mean a dial was attempted despite the earlier failure.
func TestNewTLSConnectionRootCAFileMissing(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := NewTLSConnection(ctx, "127.0.0.1:1", WithRootCAFile("/nonexistent/path/to/ca.pem"))
	if err == nil {
		t.Fatal("expected error for nonexistent root CA file")
	}
	if !strings.Contains(err.Error(), "root CA file") {
		t.Fatalf("expected error to mention the failing option, got: %s", err)
	}
}

// TestNewTLSConnectionContextCancelled proves an already-cancelled context aborts the dial.
func TestNewTLSConnectionContextCancelled(t *testing.T) {
	ln := newLocalListener(t)
	defer ln.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := NewTLSConnection(ctx, ln.Addr().String(), WithInsecureSkipVerify())
	if err == nil {
		t.Fatal("expected error from already-cancelled context")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected error to wrap context.Canceled, got: %s", err)
	}
}

// TestWithTLSConfigNilIsAnError: a nil config must be reported rather than
// panicking when a later option dereferences it.
func TestWithTLSConfigNilIsAnError(t *testing.T) {
	_, err := NewTLSConnection(
		context.Background(),
		"127.0.0.1:1",
		WithTLSConfig(nil),
		WithInsecureSkipVerify(),
	)
	if err == nil {
		t.Fatalf("Expected an error for a nil *tls.Config")
	}
	if !strings.Contains(err.Error(), "nil *tls.Config") {
		t.Fatalf("Expected the error to name the nil config, got: %v", err)
	}
}

// TestWithTLSConfigIsCloned: mutating the caller's config after passing it must
// not reach the connection.
func TestWithTLSConfigIsCloned(t *testing.T) {
	ln := newLocalListener(t)
	defer ln.Close()
	go func() {
		c, err := ln.Accept()
		if err == nil {
			_ = c.Close()
		}
	}()

	cfg := &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS12} // nolint:gosec
	opt := WithTLSConfig(cfg)

	settings := newTLSSettings()
	opt(settings)
	cfg.ServerName = "mutated-after-the-fact"

	if settings.err != nil {
		t.Fatalf("Unexpected error: %v", settings.err)
	}
	if settings.config.ServerName != "" {
		t.Fatalf("Config was aliased, not cloned: ServerName leaked as %q", settings.config.ServerName)
	}
}

// TestWithRootCAsSnapshotsAtOptionCreation: x509.CertPool is mutable, and the
// snapshot is taken when the option is built, so the window between building it
// and passing it to NewTLSConnection is covered too.
func TestWithRootCAsSnapshotsAtOptionCreation(t *testing.T) {
	extra := generateSelfSignedCert(t, "other.example.com")
	parsed, err := x509.ParseCertificate(extra.Certificate[0])
	if err != nil {
		t.Fatalf("Failed to parse the generated certificate: %v", err)
	}

	pool := x509.NewCertPool()
	opt := WithRootCAs(pool)

	pool.AddCert(parsed) // after the option exists, before it is applied

	settings := newTLSSettings()
	opt(settings)
	if n := len(settings.config.RootCAs.Subjects()); n != 0 { // nolint:staticcheck // pool size only
		t.Fatalf("Expected the trust set snapshotted at option creation (0 certs), got %d", n)
	}
}

// TestWithRootCAsDoesNotShareAcrossApplications: one option applied twice must
// not hand both configs the same mutable pool.
func TestWithRootCAsDoesNotShareAcrossApplications(t *testing.T) {
	pool := x509.NewCertPool()
	opt := WithRootCAs(pool)

	a, b := newTLSSettings(), newTLSSettings()
	opt(a)
	opt(b)

	if a.config.RootCAs == b.config.RootCAs {
		t.Fatalf("Two applications of one option share a pool")
	}
}

func TestWithRootCAsNilMeansSystemRoots(t *testing.T) {
	settings := newTLSSettings()
	WithRootCAs(nil)(settings)
	if settings.config.RootCAs != nil {
		t.Fatalf("A nil pool must leave RootCAs nil, meaning the system roots")
	}
}

// TestNilTLSOptionIsRejected: TLSOption is a function type, so nil is reachable
// from a caller and must be an error rather than a panic.
func TestNilTLSOptionIsRejected(t *testing.T) {
	_, err := NewTLSConnection(context.Background(), "127.0.0.1:1", WithInsecureSkipVerify(), nil)
	if err == nil {
		t.Fatalf("Expected an error for a nil TLSOption")
	}
	if !strings.Contains(err.Error(), "nil TLSOption") {
		t.Fatalf("Expected the error to name the nil option, got: %v", err)
	}
}

func TestOptionsStopAtTheFirstFailure(t *testing.T) {
	ran := false
	probe := func(s *tlsSettings) { ran = true }

	_, err := NewTLSConnection(context.Background(), "127.0.0.1:1",
		WithRootCAFile("/nonexistent/ca.pem"),
		TLSOption(probe),
	)
	if err == nil {
		t.Fatalf("Expected the missing CA file to fail the constructor")
	}
	if !strings.Contains(err.Error(), "root CA file") {
		t.Fatalf("Expected the first error to win, got: %v", err)
	}
	if ran {
		t.Fatalf("An option after a failed one must not run")
	}
}
