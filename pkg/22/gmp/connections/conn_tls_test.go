package connections

import (
	"crypto/tls"
	"encoding/xml"
	"log"
	"net"
	"testing"
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
		MinVersion:         tls.VersionSSL30,
		MaxVersion:         tls.VersionSSL30,
	})
	if err != nil {
		ln, err = tls.Listen("tcp6", "[::1]:0", nil)
	}
	if err != nil {
		t.Fatal(err)
	}
	return ln
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

	conn, err := NewTLSConnection(ln.Addr().String(), true)
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

	err = conn.Execute(cmd, response)
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

	_, err := NewTLSConnection(ln.Addr().String(), true)

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
