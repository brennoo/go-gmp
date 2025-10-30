package connections

import (
	"context"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"testing"
)

func TestNewUnixConnection(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "testunixsocket")
	if err != nil {
		log.Fatal(err)
	}
	unixSocketFilename := tmpfile.Name()
	os.Remove(unixSocketFilename)
	defer os.Remove(unixSocketFilename)

	lc := &net.ListenConfig{}
	ln, err := lc.Listen(context.Background(), "unix", unixSocketFilename)
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	errCh := make(chan error, 1)
	go func() {
		for {
			conn, acceptErr := ln.Accept()
			if acceptErr != nil {
				errCh <- acceptErr
				return
			}
			buf2 := make([]byte, 150000)
			nRead, _ := conn.Read(buf2)
			conn.Write(buf2[:nRead])
		}
	}()

	conn, err := NewUnixConnection(unixSocketFilename)
	if err != nil {
		t.Fatalf("Unexpected error during NewUnixConnection: %s", err)
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

func TestNewUnixConnectionUnixSocketIsNotThere(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "testunixsocket")
	if err != nil {
		log.Fatal(err)
	}
	unixSocketFilename := tmpfile.Name()
	os.Remove(unixSocketFilename)

	lc := &net.ListenConfig{}
	ln, err := lc.Listen(context.Background(), "unix", unixSocketFilename)
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	errCh := make(chan error, 1)
	go func() {
		for {
			conn, acceptErr := ln.Accept()
			if acceptErr != nil {
				errCh <- acceptErr
				return
			}
			buf2 := make([]byte, 150000)
			nRead, _ := conn.Read(buf2)
			conn.Write(buf2[:nRead])
		}
	}()

	os.Remove(unixSocketFilename)
	_, err = NewUnixConnection(unixSocketFilename)
	expectedError := "dial unix FILEPATH: connect: no such file or directory"
	expectedError = strings.ReplaceAll(expectedError, "FILEPATH", unixSocketFilename)
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Unexpected error during Execute.\nExpected: %s\n     Got: %s", expectedError, err)
	}

	select {
	case err := <-errCh:
		t.Fatalf("Unexpected error during Accept: %s", err)
	default:
	}
}
