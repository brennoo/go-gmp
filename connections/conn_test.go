package connections

import (
	"context"
	"encoding/xml"
	"errors"
	"net"
	"testing"
)

func TestNewConn(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := New(c1)

	go func() {
		buf2 := make([]byte, 150000)
		nRead, _ := c2.Read(buf2)
		c2.Write(buf2[:nRead])
	}()

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

	err := conn.Execute(context.Background(), cmd, response)
	if err != nil {
		t.Fatalf("Unexpected error during Execute: %s", err)
	}

	if response.Foo != expectedValue {
		t.Fatalf("Unexpected response value: %s\nExpected: %s", response.Foo, expectedValue)
	}
}

// TestReExportedSentinelsMatchViaErrorsIs proves that the sentinels
// re-exported from the internal implementation package retain their
// identity: an error produced through the public API (New, then Close, then
// Execute) must satisfy errors.Is against the connections-package spelling of
// the sentinel, not just the internal one.
func TestReExportedSentinelsMatchViaErrorsIs(t *testing.T) {
	c1, _ := net.Pipe()
	conn := New(c1)

	if err := conn.Close(); err != nil {
		t.Fatalf("Unexpected error during Close: %s", err)
	}

	cmd := &struct {
		XMLName xml.Name `xml:"cmd"`
		Foo     string   `xml:"foo"`
	}{Foo: "x"}
	response := &struct {
		Foo string `xml:"foo"`
	}{}

	err := conn.Execute(context.Background(), cmd, response)
	if !errors.Is(err, ErrConnectionClosed) {
		t.Fatalf("Expected errors.Is to match the re-exported ErrConnectionClosed, got: %v", err)
	}
}
