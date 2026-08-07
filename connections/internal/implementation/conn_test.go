package implementation

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/brennoo/go-gmp/commands"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// newPipeConn wires a Connection to one end of a net.Pipe and hands back the
// other (server) end.
//
// net.Pipe is fine for single-transaction tests. It is deliberately NOT used by
// the concurrency test: net.Pipe is synchronous and unbuffered, so it
// accidentally serializes callers and would mask exactly the bug that test
// exists to catch.
func newPipeConn(t *testing.T) (*Connection, net.Conn) {
	t.Helper()
	client, server := net.Pipe()
	t.Cleanup(func() {
		_ = client.Close()
		_ = server.Close()
	})
	conn := New(client)
	return conn, server
}

// mustFinish runs fn in a goroutine and fails if it has not returned within d.
// Every framing test uses it so that a regression of the "read loop hangs on an
// exact multiple of the buffer size" bug fails fast instead of hanging the
// whole suite forever.
//
// fn must not call t.Fatalf; it runs on another goroutine. Have it assign to
// captured variables and assert after mustFinish returns.
func mustFinish(t *testing.T, d time.Duration, what string, fn func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()
	select {
	case <-done:
	case <-time.After(d):
		t.Fatalf("%s did not return within %s (likely hang)", what, d)
	}
}

// xmlDocOfSize builds a well-formed XML document of exactly n bytes.
func xmlDocOfSize(t *testing.T, n int) string {
	t.Helper()
	const prefix, suffix = "<r><d>", "</d></r>"
	if n < len(prefix)+len(suffix) {
		t.Fatalf("cannot build an XML document of %d bytes", n)
	}
	return prefix + strings.Repeat("a", n-len(prefix)-len(suffix)) + suffix
}

// echoServer copies everything it reads straight back. Because transactions are
// serialized, a byte-for-byte echo is a valid per-request response.
func echoServer(server net.Conn) {
	go func() { _, _ = io.Copy(server, server) }()
}

// readOneRequest reads a single request off server and reports when it has
// landed. Once the returned channel is closed the client's write has completed,
// so the client is provably admitted and blocked in its read — which lets the
// cancellation tests be deterministic instead of sleep-based.
func readOneRequest(t *testing.T, server net.Conn) <-chan struct{} {
	t.Helper()
	arrived := make(chan struct{})
	go func() {
		buf := make([]byte, 65536)
		_, _ = server.Read(buf)
		close(arrived)
	}()
	return arrived
}

// followUpExecute issues a fresh Execute under a bound and returns its error.
//
// The bound matters: these call sites assert that a poisoned connection is
// refused at entry. Should poisoning regress, the call would instead block on a
// dead socket, so an unbounded assertion would hang the suite rather than fail
// it.
func followUpExecute(t *testing.T, conn *Connection, what string) error {
	t.Helper()
	var err error
	mustFinish(t, 5*time.Second, what, func() {
		err = conn.Execute(context.Background(), &simpleCmd{Foo: "follow-up"}, &simpleResp{})
	})
	return err
}

// newTCPEchoConn starts a loopback TCP echo server and returns a Connection
// dialed to it, plus a channel closed when the server goroutine exits.
//
// A real TCP socket is mandatory for the concurrency tests: net.Pipe is
// synchronous and unbuffered, so it accidentally serializes callers and would
// mask the very interleaving these tests exist to catch.
func newTCPEchoConn(t *testing.T) (*Connection, <-chan struct{}) {
	t.Helper()
	ctx := context.Background()

	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	t.Cleanup(func() { _ = ln.Close() })

	served := make(chan struct{})
	go func() {
		defer close(served)
		c, err := ln.Accept()
		if err != nil {
			return
		}
		defer c.Close()
		_, _ = io.Copy(c, c)
	}()

	var d net.Dialer
	client, err := d.DialContext(ctx, "tcp", ln.Addr().String())
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	conn := New(client)
	return conn, served
}

type simpleCmd struct {
	XMLName xml.Name `xml:"cmd"`
	Foo     string   `xml:"foo"`
}

type simpleResp struct {
	Foo string `xml:"foo"`
}

// countingConn records how many times Write was called, so tests can prove a
// refused transaction never reached the socket.
type countingConn struct {
	net.Conn
	writes atomic.Int32
}

func (c *countingConn) Write(p []byte) (int, error) {
	c.writes.Add(1)
	return c.Conn.Write(p)
}

// shortWriteConn returns n < len(p) with a nil error. A real socket does not do
// this — io.Writer's contract forbids it — but hand-written mocks in this repo
// can, and checking n is the caller's obligation regardless.
type shortWriteConn struct {
	net.Conn
	writes atomic.Int32
}

func (c *shortWriteConn) Write(p []byte) (int, error) {
	c.writes.Add(1)
	if len(p) == 0 {
		return 0, nil
	}
	return len(p) - 1, nil
}

// sideEffectCmd fails to marshal and records that it was asked to: the failure
// proves a lifecycle error outranks a marshalling error, the counter proves
// MarshalXML never runs at all.
type sideEffectCmd struct {
	marshalled atomic.Int32
}

func (c *sideEffectCmd) MarshalXML(*xml.Encoder, xml.StartElement) error {
	c.marshalled.Add(1)
	return errors.New("marshalling deliberately failed")
}

// cancelOnFinalReadConn deterministically reproduces a transaction completing
// and then being destroyed by a cancellation moments later. Its single Read
// hands over the whole response, fires the cancellation, then blocks until the
// watchdog closes it — so the connection is provably poisoned before Read
// returns, yet the response is complete and frames without reading again,
// making performRequest return nil.
type cancelOnFinalReadConn struct {
	net.Conn  // never set: every method the transport uses is overridden below
	doc       string
	cancel    context.CancelFunc
	closed    chan struct{}
	closeOnce sync.Once
	served    atomic.Bool
}

func newCancelOnFinalReadConn(doc string, cancel context.CancelFunc) *cancelOnFinalReadConn {
	return &cancelOnFinalReadConn{doc: doc, cancel: cancel, closed: make(chan struct{})}
}

func (c *cancelOnFinalReadConn) Write(p []byte) (int, error) { return len(p), nil }

func (c *cancelOnFinalReadConn) Read(p []byte) (int, error) {
	if c.served.Swap(true) {
		// Should be unreachable: the document is delivered whole, so framing
		// completes out of the buffered reader. Block rather than spin if the
		// framing logic ever regresses into asking for more.
		<-c.closed
		return 0, io.EOF
	}
	n := copy(p, c.doc)
	c.cancel()
	<-c.closed // the watchdog has now poisoned the connection and closed us
	return n, nil
}

func (c *cancelOnFinalReadConn) Close() error {
	c.closeOnce.Do(func() { close(c.closed) })
	return nil
}

func (c *cancelOnFinalReadConn) SetDeadline(time.Time) error      { return nil }
func (c *cancelOnFinalReadConn) SetReadDeadline(time.Time) error  { return nil }
func (c *cancelOnFinalReadConn) SetWriteDeadline(time.Time) error { return nil }

// ---------------------------------------------------------------------------
// existing tests (preserved; call sites adapted to the new internals)
// ---------------------------------------------------------------------------

func TestNew(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := New(c1)

	var wg sync.WaitGroup
	var nRead int
	var errRead error

	wg.Add(1)
	go func() {
		buf := make([]byte, 1024)
		nRead, errRead = c2.Read(buf)
		wg.Done()
	}()

	n, err := conn.rawConn.Write([]byte("0123456789"))
	if err != nil {
		t.Fatalf("Unexpected error during write: %s", err)
	}

	if n != 10 {
		t.Fatalf("Expected to write 10 bytes, wrote: %d", n)
	}

	wg.Wait()

	if errRead != nil {
		t.Fatalf("Unexpected error during read: %s", errRead)
	}

	if nRead != 10 {
		t.Fatalf("Expected to read 10 bytes, got: %d", nRead)
	}

	// New must also install the persistent buffered reader; without it every
	// transaction would strand over-read bytes in a discarded buffer.
	if conn.reader == nil {
		t.Fatalf("Expected New to install a persistent reader")
	}

	errClose := conn.Close()
	if errClose != nil {
		t.Fatalf("Unexpected error during Close: %s", errClose)
	}
}

// TestPerformRequest round-trips a large payload.
//
// Changed from the original: the payload is now a large XML document rather
// than 50000 bytes of arbitrary binary. Framing is done on the XML document
// itself (T-1), so a response that is not XML has no end-of-message signal at
// all. The intent — a payload far larger than any single socket read completes
// correctly — is unchanged, and is in fact tested harder, since the response
// now spans many reads of the buffered reader.
func TestPerformRequest(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := New(c1)

	doc := xmlDocOfSize(t, 50000)

	go func() {
		buf2 := make([]byte, 150000)
		_, _ = c2.Read(buf2)
		_, _ = c2.Write([]byte(doc))
	}()

	var bufRead []byte
	var err error
	mustFinish(t, 10*time.Second, "performRequest", func() {
		bufRead, err = performRequest(context.Background(), conn.rawConn, conn.reader, []byte(doc))
	})
	if err != nil {
		t.Fatalf("Unexpected error during performRequest: %s", err)
	}

	if len(bufRead) != 50000 {
		t.Fatalf("Expected to read 50000 bytes, got: %d", len(bufRead))
	}

	if string(bufRead) != doc {
		t.Fatalf("Round-tripped document did not match the original")
	}
}

func TestPerformRequestWriteFail(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := New(c1)

	c1.Close()
	go func() {
		buf2 := make([]byte, 150000)
		nRead, _ := c2.Read(buf2)
		c2.Write(buf2[:nRead])
	}()

	buf := make([]byte, 50000)
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i % 255)
	}

	_, err := performRequest(context.Background(), conn.rawConn, conn.reader, buf)
	expectedError := "io: read/write on closed pipe"
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Unexpected error during performRequest.\nExpected: %s\n     Got: %s", expectedError, err)
	}
}

func TestPerformRequestReadFail(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := New(c1)

	go func() {
		buf2 := make([]byte, 150000)
		c2.Read(buf2)
		c2.Close()
	}()

	buf := make([]byte, 50000)
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i % 255)
	}

	_, err := performRequest(context.Background(), conn.rawConn, conn.reader, buf)
	expectedError := "EOF"
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Unexpected error during performRequest.\nExpected: %s\n     Got: %s", expectedError, err)
	}
}

// Test that performRequest honors context deadlines via SetWriteDeadline.
func TestPerformRequestWriteDeadline(t *testing.T) {
	c1, _ := net.Pipe()

	conn := New(c1)

	// Very short timeout; no reader on the other end, so Write should time out
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	buf := []byte("hello")
	_, err := performRequest(ctx, conn.rawConn, conn.reader, buf)
	if err == nil {
		t.Fatalf("Expected write timeout error, got nil")
	}
	var nerr net.Error
	if !errors.As(err, &nerr) || !nerr.Timeout() {
		t.Fatalf("Expected net.Error timeout on write, got: %v", err)
	}
}

// Test that performRequest honors context deadlines via SetReadDeadline.
func TestPerformRequestReadDeadline(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := New(c1)

	go func() {
		// Read the request to let write succeed, but never write a response
		buf := make([]byte, 1024)
		_, _ = c2.Read(buf)
		// do not write back -> read on c1 should block until deadline
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := performRequest(ctx, conn.rawConn, conn.reader, []byte("request"))
	if err == nil {
		t.Fatalf("Expected read timeout error, got nil")
	}
	var nerr net.Error
	if !errors.As(err, &nerr) || !nerr.Timeout() {
		t.Fatalf("Expected net.Error timeout on read, got: %v", err)
	}
}

// TestAlreadyCancelledContextIsRefusedWithoutTouchingTheSocket supersedes an
// earlier test asserting the opposite. A pre-cancelled context used to arm a
// watchdog that fired immediately, so one caller's stale context permanently
// destroyed a connection shared by everyone else. Nothing has been written at
// that point, so the correct behaviour matches cancel-while-waiting: refuse the
// caller, change no state, leave the socket alone.
func TestAlreadyCancelledContextIsRefusedWithoutTouchingTheSocket(t *testing.T) {
	client, server := net.Pipe()
	t.Cleanup(func() {
		_ = client.Close()
		_ = server.Close()
	})

	counter := &countingConn{Conn: client}
	conn := New(counter)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var err error
	mustFinish(t, 5*time.Second, "Execute with a pre-cancelled context", func() {
		err = conn.Execute(ctx, &simpleCmd{Foo: "x"}, &simpleResp{})
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Expected context.Canceled, got: %v", err)
	}
	if errors.Is(err, ErrConnectionBroken) || errors.Is(err, ErrConnectionClosed) {
		t.Fatalf("A refused caller must not report a lifecycle error, got: %v", err)
	}
	if n := counter.writes.Load(); n != 0 {
		t.Fatalf("Expected the socket never to be written, saw %d writes", n)
	}

	conn.mu.Lock()
	state := conn.state
	conn.mu.Unlock()
	if state != stateOpen {
		t.Fatalf("One caller's stale context must not disturb the connection; got state %v", state)
	}

	// The decisive assertion: the connection is still usable by everyone else.
	echoServer(server)
	resp := &simpleResp{}
	mustFinish(t, 5*time.Second, "Execute after the refusal", func() {
		err = conn.Execute(context.Background(), &simpleCmd{Foo: "survivor"}, resp)
	})
	if err != nil {
		t.Fatalf("The connection should have survived untouched, got: %v", err)
	}
	if resp.Foo != "survivor" {
		t.Fatalf("Expected Foo=survivor, got %q", resp.Foo)
	}
}

func TestExecute(t *testing.T) {
	conn, server := newPipeConn(t)
	echoServer(server)

	expectedValue := "0123456789"
	cmd := &simpleCmd{Foo: expectedValue}
	response := &simpleResp{}

	var err error
	mustFinish(t, 5*time.Second, "Execute", func() {
		err = conn.Execute(context.Background(), cmd, response)
	})
	if err != nil {
		t.Fatalf("Unexpected error during Execute: %s", err)
	}

	if response.Foo != expectedValue {
		t.Fatalf("Unexpected response value: %s\nExpected: %s", response.Foo, expectedValue)
	}
}

func TestExecuteXMLMarshallFail(t *testing.T) {
	conn, server := newPipeConn(t)
	echoServer(server)

	cmd := func() {}
	response := &simpleResp{}

	err := conn.Execute(context.Background(), cmd, response)
	expectedError := "xml: unsupported type: func()"
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Unexpected error during Execute.\nExpected: %s\n     Got: %s", expectedError, err)
	}
}

func TestExecuteXMLUnmarshallFail(t *testing.T) {
	conn, server := newPipeConn(t)
	echoServer(server)

	cmd := &simpleCmd{Foo: "0123456789"}
	response := simpleResp{} // not a pointer

	var err error
	mustFinish(t, 5*time.Second, "Execute", func() {
		err = conn.Execute(context.Background(), cmd, response)
	})
	expectedError := "non-pointer passed to Unmarshal"
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Unexpected error during Execute.\nExpected: %s\n     Got: %s", expectedError, err)
	}
}

func TestExecutePerformRequestFail(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := New(c1)

	c1.Close()
	go func() {
		buf2 := make([]byte, 150000)
		nRead, _ := c2.Read(buf2)
		c2.Write(buf2[:nRead])
	}()

	cmd := &simpleCmd{Foo: "0123456789"}
	response := &simpleResp{}

	err := conn.Execute(context.Background(), cmd, response)
	expectedError := "io: read/write on closed pipe"
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Unexpected error during Execute.\nExpected: %s\n     Got: %s", expectedError, err)
	}
}

func TestRawXML(t *testing.T) {
	conn, server := newPipeConn(t)
	echoServer(server)

	doc := `<raw_test_command foo="bar"/>`

	var resp string
	var err error
	mustFinish(t, 5*time.Second, "RawXML", func() {
		resp, err = conn.RawXML(context.Background(), doc)
	})
	if err != nil {
		t.Fatalf("Unexpected error from RawXML: %v", err)
	}
	if resp != doc {
		t.Errorf("RawXML response did not match input.\nExpected: %s\nGot: %s", doc, resp)
	}
}

func TestRawXMLEmptyResponse(t *testing.T) {
	conn, server := newPipeConn(t)

	go func() {
		// Do not write anything back, simulating an empty response
		_, _ = server.Read(make([]byte, 1024))
		server.Close()
	}()

	var resp string
	var err error
	mustFinish(t, 5*time.Second, "RawXML", func() {
		resp, err = conn.RawXML(context.Background(), `<raw_test_command foo="bar"/>`)
	})
	if err == nil {
		t.Errorf("Expected error for empty response, got nil")
	}
	if resp != "" {
		t.Errorf("Expected empty string for response, got: %q", resp)
	}
}

// TestExecuteWithResponseWithStatus tests the new error handling for ResponseWithStatus.
func TestExecuteWithResponseWithStatus(t *testing.T) {
	conn, server := newPipeConn(t)
	echoServer(server)

	cmd := &simpleCmd{Foo: "test"}

	// Test with response that implements ResponseWithStatus and returns success
	response := &mockResponseWithStatus{
		Status:     "200",
		StatusText: "OK",
	}

	var err error
	mustFinish(t, 5*time.Second, "Execute", func() {
		err = conn.Execute(context.Background(), cmd, response)
	})
	if err != nil {
		t.Fatalf("Unexpected error during Execute with success response: %s", err)
	}
}

// TestExecuteWithResponseWithStatusError tests the new error handling for ResponseWithStatus with error.
func TestExecuteWithResponseWithStatusError(t *testing.T) {
	conn, server := newPipeConn(t)
	echoServer(server)

	cmd := &simpleCmd{Foo: "test"}

	// Test with response that implements ResponseWithStatus and returns error
	response := &mockResponseWithStatus{
		Status:     "401",
		StatusText: "Unauthorized",
	}

	var err error
	mustFinish(t, 5*time.Second, "Execute", func() {
		err = conn.Execute(context.Background(), cmd, response)
	})
	if err == nil {
		t.Fatalf("Expected error during Execute with error response, got nil")
	}

	// Check that it's a GMPError
	var gmpErr *commands.GMPError
	if !errors.As(err, &gmpErr) {
		t.Fatalf("Expected GMPError, got %T", err)
	}

	if gmpErr.Type != commands.ErrorTypeAuthentication {
		t.Fatalf("Expected authentication error, got %v", gmpErr.Type)
	}
}

// TestExecuteWithResponseWithoutStatus tests that non-ResponseWithStatus responses work normally.
func TestExecuteWithResponseWithoutStatus(t *testing.T) {
	conn, server := newPipeConn(t)
	echoServer(server)

	cmd := &simpleCmd{Foo: "test"}
	response := &simpleResp{}

	var err error
	mustFinish(t, 5*time.Second, "Execute", func() {
		err = conn.Execute(context.Background(), cmd, response)
	})
	if err != nil {
		t.Fatalf("Unexpected error during Execute with non-ResponseWithStatus: %s", err)
	}
}

// mockResponseWithStatus implements ResponseWithStatus for testing.
type mockResponseWithStatus struct {
	Status     string `xml:"status,attr"`
	StatusText string `xml:"status_text,attr"`
}

func (r *mockResponseWithStatus) GetStatus() (string, string) {
	return r.Status, r.StatusText
}

// ---------------------------------------------------------------------------
// T-1: framing
// ---------------------------------------------------------------------------

// TestFramingExactBufferMultiples is the direct regression test for the
// original bug: the old loop treated a read of fewer than 1024 bytes as
// end-of-message, so a response whose length was an exact multiple of 1024
// issued one more Read that blocked forever.
func TestFramingExactBufferMultiples(t *testing.T) {
	for _, size := range []int{1024, 2048, 3072, 4096} {
		t.Run(fmt.Sprintf("%d_bytes", size), func(t *testing.T) {
			conn, server := newPipeConn(t)
			doc := xmlDocOfSize(t, size)

			go func() {
				buf := make([]byte, 65536)
				_, _ = server.Read(buf)
				_, _ = server.Write([]byte(doc))
			}()

			var got string
			var err error
			mustFinish(t, 5*time.Second, "RawXML", func() {
				got, err = conn.RawXML(context.Background(), "<get_version/>")
			})
			if err != nil {
				t.Fatalf("Unexpected error for a %d-byte response: %v", size, err)
			}
			if len(got) != size {
				t.Fatalf("Expected %d bytes back, got %d", size, len(got))
			}
			if got != doc {
				t.Fatalf("Response did not match the document the server sent")
			}
		})
	}
}

// TestFramingAcrossReadBoundaries writes the document to the socket in small
// chunks so that the document boundary never coincides with a read boundary.
func TestFramingAcrossReadBoundaries(t *testing.T) {
	conn, server := newPipeConn(t)
	doc := `<get_version_response status="200"><version>22.4</version></get_version_response>`

	go func() {
		buf := make([]byte, 65536)
		_, _ = server.Read(buf)
		for i := 0; i < len(doc); i += 7 {
			end := i + 7
			if end > len(doc) {
				end = len(doc)
			}
			_, _ = server.Write([]byte(doc[i:end]))
			time.Sleep(time.Millisecond)
		}
	}()

	var got string
	var err error
	mustFinish(t, 10*time.Second, "RawXML", func() {
		got, err = conn.RawXML(context.Background(), "<get_version/>")
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got != doc {
		t.Fatalf("Chunked document was framed incorrectly.\nExpected: %s\n     Got: %s", doc, got)
	}
}

// TestFramingLeadingPrologAndTrailingWhitespace pins down the "wire segment"
// contract documented on RawXML: a decoder legally emits ProcInst, Comment,
// Directive and whitespace CharData before the root while depth is still zero,
// so a depth counter that starts counting immediately would terminate before
// the document even began. Those leading bytes are captured; trailing
// whitespace after the root's closing tag is not, because it is never consumed.
func TestFramingLeadingPrologAndTrailingWhitespace(t *testing.T) {
	conn, server := newPipeConn(t)

	const (
		leading  = "\n\t<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!-- a comment -->\n<!DOCTYPE get_version_response>\n"
		root     = `<get_version_response status="200"><version>22.4</version></get_version_response>`
		trailing = "\n\n   "
	)

	go func() {
		buf := make([]byte, 65536)
		_, _ = server.Read(buf)
		_, _ = server.Write([]byte(leading + root + trailing))
	}()

	var got string
	var err error
	mustFinish(t, 5*time.Second, "RawXML", func() {
		got, err = conn.RawXML(context.Background(), "<get_version/>")
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got != leading+root {
		t.Fatalf("Wire segment mismatch.\nExpected: %q\n     Got: %q", leading+root, got)
	}
	if strings.HasSuffix(got, " ") || strings.HasSuffix(got, "\n") {
		t.Fatalf("Trailing whitespace should not have been consumed, got: %q", got)
	}
}

// TestFramingPrologDecodesViaExecute confirms the captured wire segment,
// leading prolog and all, still unmarshals into a typed response.
func TestFramingPrologDecodesViaExecute(t *testing.T) {
	conn, server := newPipeConn(t)

	go func() {
		buf := make([]byte, 65536)
		_, _ = server.Read(buf)
		_, _ = server.Write([]byte("  <?xml version=\"1.0\"?><!-- hi -->\n<cmd><foo>prolog</foo></cmd>"))
	}()

	response := &simpleResp{}
	var err error
	mustFinish(t, 5*time.Second, "Execute", func() {
		err = conn.Execute(context.Background(), &simpleCmd{Foo: "x"}, response)
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if response.Foo != "prolog" {
		t.Fatalf("Expected Foo=prolog, got %q", response.Foo)
	}
}

// TestFramingNestedAndSelfClosingElements checks the depth counter against
// shapes that could confuse it: a self-closing root, self-closing children
// (which the decoder synthesizes an EndElement for), deep nesting, and an
// empty-but-not-self-closing element.
func TestFramingNestedAndSelfClosingElements(t *testing.T) {
	docs := map[string]string{
		"self_closing_root":  `<get_version/>`,
		"nested":             `<a><b><c><d>deep</d></c></b></a>`,
		"self_closing_child": `<a><b/><c/><d/></a>`,
		"empty_element":      `<a><b></b></a>`,
		"mixed":              `<a attr="1"><b/>text<c><d/></c><e></e></a>`,
	}

	for name, doc := range docs {
		t.Run(name, func(t *testing.T) {
			conn, server := newPipeConn(t)
			go func() {
				buf := make([]byte, 65536)
				_, _ = server.Read(buf)
				_, _ = server.Write([]byte(doc))
			}()

			var got string
			var err error
			mustFinish(t, 5*time.Second, "RawXML", func() {
				got, err = conn.RawXML(context.Background(), "<get_version/>")
			})
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got != doc {
				t.Fatalf("Framing mismatch.\nExpected: %s\n     Got: %s", doc, got)
			}
		})
	}
}

// TestFramingTwoDocumentsInASingleWrite checks framing correctness when the
// fake server emits two documents in one Write: each sequential call on the
// same Connection must frame exactly one of them. This is about the decoder
// consuming precisely the root element and no more; it is NOT an assertion that
// GMP servers pipeline responses, which is not established protocol behavior.
func TestFramingTwoDocumentsInASingleWrite(t *testing.T) {
	conn, server := newPipeConn(t)

	const first = `<r><n>1</n></r>`
	const second = `<r><n>2</n></r>`

	go func() {
		buf := make([]byte, 65536)
		_, _ = server.Read(buf)
		_, _ = server.Write([]byte(first + second))
		_, _ = server.Read(buf)
	}()

	var got1, got2 string
	var err1, err2 error
	mustFinish(t, 5*time.Second, "RawXML pair", func() {
		got1, err1 = conn.RawXML(context.Background(), "<a/>")
		got2, err2 = conn.RawXML(context.Background(), "<b/>")
	})
	if err1 != nil || err2 != nil {
		t.Fatalf("Unexpected errors: %v / %v", err1, err2)
	}
	if got1 != first {
		t.Fatalf("First document framed incorrectly.\nExpected: %s\n     Got: %s", first, got1)
	}
	if got2 != second {
		t.Fatalf("Second document framed incorrectly.\nExpected: %s\n     Got: %s", second, got2)
	}
}

// TestFramingMalformedXMLErrorsAndPoisons: unframeable input leaves the
// message-boundary position unknown, so the connection must be poisoned rather
// than reused.
func TestFramingMalformedXMLErrorsAndPoisons(t *testing.T) {
	conn, server := newPipeConn(t)

	go func() {
		buf := make([]byte, 65536)
		_, _ = server.Read(buf)
		_, _ = server.Write([]byte(`<a></b>`))
	}()

	var err error
	mustFinish(t, 5*time.Second, "RawXML", func() {
		_, err = conn.RawXML(context.Background(), "<get_version/>")
	})
	if err == nil {
		t.Fatalf("Expected an error for malformed XML, got nil")
	}

	// The connection must now refuse further use rather than write into a
	// stream whose position it no longer knows.
	second := followUpExecute(t, conn, "Execute after malformed XML")
	if !errors.Is(second, ErrConnectionBroken) {
		t.Fatalf("Expected ErrConnectionBroken on the next call, got: %v", second)
	}
}

// TestFramingTruncatedXMLErrorsAndPoisons: EOF before the root closes is also
// boundary-unknown.
func TestFramingTruncatedXMLErrorsAndPoisons(t *testing.T) {
	conn, server := newPipeConn(t)

	go func() {
		buf := make([]byte, 65536)
		_, _ = server.Read(buf)
		_, _ = server.Write([]byte(`<a><b>truncated`))
		_ = server.Close()
	}()

	var err error
	mustFinish(t, 5*time.Second, "RawXML", func() {
		_, err = conn.RawXML(context.Background(), "<get_version/>")
	})
	if err == nil {
		t.Fatalf("Expected an error for truncated XML, got nil")
	}

	second := followUpExecute(t, conn, "Execute after truncated XML")
	if !errors.Is(second, ErrConnectionBroken) {
		t.Fatalf("Expected ErrConnectionBroken on the next call, got: %v", second)
	}
}

// ---------------------------------------------------------------------------
// T-13: poisoning classification
// ---------------------------------------------------------------------------

// TestUnmarshalFailureLeavesConnectionUsable is the other half of the
// classification. A document that framed cleanly leaves the stream position
// known; only the decode into a concrete Go type failed (schema drift). That
// must not poison the connection.
func TestUnmarshalFailureLeavesConnectionUsable(t *testing.T) {
	conn, server := newPipeConn(t)
	echoServer(server)

	// The echoed document is well-formed, but <foo> holds text that cannot be
	// decoded into an int.
	typedResp := &struct {
		Foo int `xml:"foo"`
	}{}
	var err error
	mustFinish(t, 5*time.Second, "Execute (mismatched type)", func() {
		err = conn.Execute(context.Background(), &simpleCmd{Foo: "not-a-number"}, typedResp)
	})
	if err == nil {
		t.Fatalf("Expected a decode error, got nil")
	}
	if errors.Is(err, ErrConnectionBroken) || errors.Is(err, ErrConnectionClosed) {
		t.Fatalf("A clean frame that failed to decode must not poison the connection, got: %v", err)
	}
	var netErr *commands.NetworkError
	if !errors.As(err, &netErr) || netErr.Type != commands.ErrorTypeInvalidRequest {
		t.Fatalf("Expected an invalid_request NetworkError, got %#v", err)
	}

	// The connection must still be open and usable.
	good := &simpleResp{}
	mustFinish(t, 5*time.Second, "Execute (reuse)", func() {
		err = conn.Execute(context.Background(), &simpleCmd{Foo: "second"}, good)
	})
	if err != nil {
		t.Fatalf("Expected the connection to remain usable, got: %v", err)
	}
	if good.Foo != "second" {
		t.Fatalf("Expected Foo=second, got %q", good.Foo)
	}
}

// TestShortWritePoisonsWithoutRetry covers the write-count check. A short write
// with a nil error is ambiguous about what the server received, so it poisons
// and must not resume the remainder.
func TestShortWritePoisonsWithoutRetry(t *testing.T) {
	client, server := net.Pipe()
	t.Cleanup(func() {
		_ = client.Close()
		_ = server.Close()
	})

	sw := &shortWriteConn{Conn: client}
	conn := New(sw)

	var err error
	mustFinish(t, 5*time.Second, "Execute", func() {
		err = conn.Execute(context.Background(), &simpleCmd{Foo: "x"}, &simpleResp{})
	})
	if err == nil {
		t.Fatalf("Expected a short-write error, got nil")
	}
	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("Expected io.ErrShortWrite, got: %v", err)
	}
	if n := sw.writes.Load(); n != 1 {
		t.Fatalf("Expected exactly one Write (no retry of the remainder), got %d", n)
	}

	second := followUpExecute(t, conn, "Execute after short write")
	if !errors.Is(second, ErrConnectionBroken) {
		t.Fatalf("Expected ErrConnectionBroken after a short write, got: %v", second)
	}
	if n := sw.writes.Load(); n != 1 {
		t.Fatalf("A broken connection must be refused before writing; Write count is %d", n)
	}
}

// ---------------------------------------------------------------------------
// T-2 / T-2b: concurrency, cancellation and lifecycle
// ---------------------------------------------------------------------------

type echoID struct {
	XMLName xml.Name `xml:"echo"`
	ID      int      `xml:"id"`
}

// TestConcurrentExecuteOverRealTCPSocket is the T-2 regression test.
//
// It must use a real TCP socket: net.Pipe is synchronous and unbuffered, so it
// accidentally serializes callers and hides the interleaving entirely. With a
// real socket and no serialization, 50 goroutines sharing one Connection
// deadlock permanently and read each other's replies.
func TestConcurrentExecuteOverRealTCPSocket(t *testing.T) {
	conn, serverDone := newTCPEchoConn(t)

	const executors = 50
	const rawCallers = 10

	failures := make(chan error, executors+rawCallers)
	var wg sync.WaitGroup

	for i := 0; i < executors; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			resp := &echoID{}
			if err := conn.Execute(context.Background(), &echoID{ID: id}, resp); err != nil {
				failures <- fmt.Errorf("Execute(%d): %w", id, err)
				return
			}
			if resp.ID != id {
				failures <- fmt.Errorf("Execute(%d) got another goroutine's response: id=%d", id, resp.ID)
			}
		}(i)
	}

	for i := 0; i < rawCallers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			want := fmt.Sprintf("<raw><id>%d</id></raw>", id)
			got, err := conn.RawXML(context.Background(), want)
			if err != nil {
				failures <- fmt.Errorf("RawXML(%d): %w", id, err)
				return
			}
			if got != want {
				failures <- fmt.Errorf("RawXML(%d) got another caller's response: %s", id, got)
			}
		}(1000 + i)
	}

	mustFinish(t, 60*time.Second, "50 concurrent Execute + 10 RawXML", wg.Wait)
	close(failures)

	for err := range failures {
		t.Error(err)
	}

	if err := conn.Close(); err != nil {
		t.Fatalf("Unexpected error from Close: %v", err)
	}
	<-serverDone
}

// TestCancelWhileQueuedForAdmission: a caller cancelled before it is ever
// admitted must return ctx.Err() without touching the socket and without
// changing connection state.
//
// This is deterministic, not sleep-based: the holder's request has provably
// reached the server before the queued caller runs, so the holder provably owns
// admission.
func TestCancelWhileQueuedForAdmission(t *testing.T) {
	conn, server := newPipeConn(t)

	admitted := make(chan struct{})
	release := make(chan struct{})

	go func() {
		buf := make([]byte, 65536)
		_, _ = server.Read(buf) // holder's request
		close(admitted)
		<-release
		_, _ = server.Write([]byte(`<r><n>1</n></r>`))

		_, _ = server.Read(buf) // later request
		_, _ = server.Write([]byte(`<r><n>3</n></r>`))
	}()

	holderErr := make(chan error, 1)
	go func() {
		_, err := conn.RawXML(context.Background(), "<hold/>")
		holderErr <- err
	}()

	<-admitted // the holder owns the transaction

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var queuedErr error
	mustFinish(t, 5*time.Second, "queued Execute", func() {
		queuedErr = conn.Execute(ctx, &simpleCmd{Foo: "queued"}, &simpleResp{})
	})
	if !errors.Is(queuedErr, context.Canceled) {
		t.Fatalf("Expected the queued caller to get context.Canceled, got: %v", queuedErr)
	}
	if errors.Is(queuedErr, ErrConnectionBroken) || errors.Is(queuedErr, ErrConnectionClosed) {
		t.Fatalf("A queued caller must not report a lifecycle error, got: %v", queuedErr)
	}

	// No state change: the holder completes, and a later call still works.
	close(release)
	if err := <-holderErr; err != nil {
		t.Fatalf("The in-flight holder should have been undisturbed, got: %v", err)
	}

	conn.mu.Lock()
	state := conn.state
	conn.mu.Unlock()
	if state != stateOpen {
		t.Fatalf("Expected the connection to still be open, got state %v", state)
	}

	var thirdErr error
	mustFinish(t, 5*time.Second, "post-cancel RawXML", func() {
		_, thirdErr = conn.RawXML(context.Background(), "<again/>")
	})
	if thirdErr != nil {
		t.Fatalf("Expected the connection to remain usable, got: %v", thirdErr)
	}
}

// TestCancelWhileInFlight: cancellation during real I/O leaves the boundary
// ambiguous, so the connection is poisoned and the caller's error wraps
// ctx.Err() rather than the bare platform "closed connection" error the
// watchdog provoked.
func TestCancelWhileInFlight(t *testing.T) {
	conn, server := newPipeConn(t)

	arrived := readOneRequest(t, server) // server never responds

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-arrived
		cancel()
	}()

	var err error
	mustFinish(t, 5*time.Second, "Execute", func() {
		err = conn.Execute(ctx, &simpleCmd{Foo: "x"}, &simpleResp{})
	})
	if err == nil {
		t.Fatalf("Expected a cancellation error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Expected the error to wrap context.Canceled, got: %v", err)
	}
	if err.Error() == "io: read/write on closed pipe" {
		t.Fatalf("Expected the cancellation to be reported, not the bare platform error")
	}

	conn.mu.Lock()
	state := conn.state
	conn.mu.Unlock()
	if state != stateBroken {
		t.Fatalf("Expected the connection to be broken, got state %v", state)
	}

	second := followUpExecute(t, conn, "Execute after in-flight cancellation")
	if !errors.Is(second, ErrConnectionBroken) {
		t.Fatalf("Expected ErrConnectionBroken on the next call, got: %v", second)
	}
}

// TestCancelAfterCloseIsANoOp: closed is terminal and outranks everything.
func TestCancelAfterCloseIsANoOp(t *testing.T) {
	conn, _ := newPipeConn(t)

	if err := conn.Close(); err != nil {
		t.Fatalf("Unexpected error from Close: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var err error
	mustFinish(t, 5*time.Second, "Execute after Close", func() {
		err = conn.Execute(ctx, &simpleCmd{Foo: "x"}, &simpleResp{})
	})
	if !errors.Is(err, ErrConnectionClosed) {
		t.Fatalf("Expected ErrConnectionClosed, got: %v", err)
	}

	conn.mu.Lock()
	state := conn.state
	conn.mu.Unlock()
	if state != stateClosed {
		t.Fatalf("closed must outrank broken; got state %v", state)
	}
}

// TestCancelDuringFlightAfterCloseReportsClosed: when Close has already
// published stateClosed, a cancellation landing afterwards must not downgrade
// the reported state to broken.
func TestCancelDuringFlightAfterCloseReportsClosed(t *testing.T) {
	conn, server := newPipeConn(t)

	arrived := readOneRequest(t, server) // server never responds

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- conn.Execute(ctx, &simpleCmd{Foo: "x"}, &simpleResp{})
	}()

	<-arrived
	if err := conn.Close(); err != nil {
		t.Fatalf("Unexpected error from Close: %v", err)
	}
	cancel() // strictly after Close published stateClosed

	select {
	case err := <-errCh:
		if !errors.Is(err, ErrConnectionClosed) {
			t.Fatalf("Expected ErrConnectionClosed, got: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("The in-flight request did not return after Close")
	}

	conn.mu.Lock()
	state := conn.state
	conn.mu.Unlock()
	if state != stateClosed {
		t.Fatalf("Expected stateClosed, got %v", state)
	}
}

// TestCloseDuringInFlightRequest: Close must not wait for the in-flight
// request, and the in-flight request must fail rather than hang.
func TestCloseDuringInFlightRequest(t *testing.T) {
	conn, server := newPipeConn(t)

	arrived := readOneRequest(t, server) // server never responds

	errCh := make(chan error, 1)
	go func() {
		errCh <- conn.Execute(context.Background(), &simpleCmd{Foo: "x"}, &simpleResp{})
	}()

	<-arrived

	start := time.Now()
	if err := conn.Close(); err != nil {
		t.Fatalf("Unexpected error from Close: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("Close blocked on the in-flight request for %s", elapsed)
	}

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatalf("Expected the in-flight request to fail")
		}
		if !errors.Is(err, ErrConnectionClosed) {
			t.Fatalf("Expected ErrConnectionClosed, got: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("The in-flight request hung after Close")
	}
}

// TestCloseIsIdempotent: repeated Close never panics and returns consistently,
// including when the socket was already closed underneath it.
func TestCloseIsIdempotent(t *testing.T) {
	conn, _ := newPipeConn(t)

	for i := 0; i < 3; i++ {
		if err := conn.Close(); err != nil {
			t.Fatalf("Close call %d returned %v, want nil", i+1, err)
		}
	}

	// A Connection that never had a socket installed must also be closeable.
	bare := &Connection{}
	if err := bare.Close(); err != nil {
		t.Fatalf("Close on a bare Connection returned %v, want nil", err)
	}

	// A socket already closed out from under Close is not a Close failure.
	client, server := net.Pipe()
	defer server.Close()
	other := New(client)
	_ = client.Close()
	if err := other.Close(); err != nil {
		t.Fatalf("Close over an already-closed socket returned %v, want nil", err)
	}
}

// TestConcurrentCloseIsRaceFree hammers Close from many goroutines alongside
// in-flight requests. Its value is under -race.
func TestConcurrentCloseIsRaceFree(t *testing.T) {
	conn, _ := newTCPEchoConn(t)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Errors are expected and irrelevant here; the point is that
			// nothing races, panics or hangs.
			_ = conn.Execute(context.Background(), &echoID{ID: id}, &echoID{})
		}(i)
	}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := conn.Close(); err != nil {
				t.Errorf("Close returned %v, want nil", err)
			}
		}()
	}

	mustFinish(t, 30*time.Second, "concurrent Execute + Close", wg.Wait)
}

// TestClosedConnectionRefusesBeforeWriting: queued and subsequent callers get a
// stable error BEFORE any write is attempted.
func TestClosedConnectionRefusesBeforeWriting(t *testing.T) {
	client, server := net.Pipe()
	t.Cleanup(func() {
		_ = client.Close()
		_ = server.Close()
	})
	echoServer(server)

	counter := &countingConn{Conn: client}
	conn := New(counter)

	var err error
	mustFinish(t, 5*time.Second, "Execute", func() {
		err = conn.Execute(context.Background(), &simpleCmd{Foo: "first"}, &simpleResp{})
	})
	if err != nil {
		t.Fatalf("Unexpected error on the first call: %v", err)
	}
	if n := counter.writes.Load(); n != 1 {
		t.Fatalf("Expected 1 write, got %d", n)
	}

	if err := conn.Close(); err != nil {
		t.Fatalf("Unexpected error from Close: %v", err)
	}

	execErr := followUpExecute(t, conn, "Execute after Close")
	if !errors.Is(execErr, ErrConnectionClosed) {
		t.Fatalf("Expected ErrConnectionClosed from Execute, got: %v", execErr)
	}
	var rawErr error
	mustFinish(t, 5*time.Second, "RawXML after Close", func() {
		_, rawErr = conn.RawXML(context.Background(), "<third/>")
	})
	if !errors.Is(rawErr, ErrConnectionClosed) {
		t.Fatalf("Expected ErrConnectionClosed from RawXML, got: %v", rawErr)
	}
	if n := counter.writes.Load(); n != 1 {
		t.Fatalf("A closed connection must be refused before writing; Write count is %d", n)
	}
}

// TestExecuteWithoutRawConn: a bare Connection literal, built without going
// through New, must report the missing socket clearly instead of panicking on
// a nil one.
func TestExecuteWithoutRawConn(t *testing.T) {
	conn := &Connection{}

	err := conn.Execute(context.Background(), &simpleCmd{Foo: "x"}, &simpleResp{})
	if !errors.Is(err, ErrNoRawConn) {
		t.Fatalf("Expected ErrNoRawConn, got: %v", err)
	}

	if _, err := conn.RawXML(context.Background(), "<a/>"); !errors.Is(err, ErrNoRawConn) {
		t.Fatalf("Expected ErrNoRawConn from RawXML, got: %v", err)
	}
}

// TestContextDeadlineMidDocument: the deadline expires after a partial response
// has already arrived but before the root element closes. The call must return
// promptly and poison the connection.
//
// Which error surfaces is a genuine race between the socket read deadline and
// the close-on-cancel watchdog — both are armed from the same context — so the
// assertion accepts either, but requires the failure to be recognisable as a
// deadline rather than an arbitrary error.
func TestContextDeadlineMidDocument(t *testing.T) {
	conn, server := newPipeConn(t)

	go func() {
		buf := make([]byte, 65536)
		_, _ = server.Read(buf)
		_, _ = server.Write([]byte(`<get_version_response status="200"><version>22`))
		// never closes the root
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	var err error
	start := time.Now()
	mustFinish(t, 5*time.Second, "Execute", func() {
		err = conn.Execute(ctx, &simpleCmd{Foo: "x"}, &simpleResp{})
	})
	if err == nil {
		t.Fatalf("Expected a deadline error, got nil")
	}
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Fatalf("Execute took %s; the deadline should have cut it short", elapsed)
	}

	var netErr net.Error
	isTimeout := errors.As(err, &netErr) && netErr.Timeout()
	if !errors.Is(err, context.DeadlineExceeded) && !isTimeout {
		t.Fatalf("Expected a deadline/timeout error, got: %v", err)
	}

	conn.mu.Lock()
	state := conn.state
	conn.mu.Unlock()
	if state != stateBroken {
		t.Fatalf("A partial read leaves the boundary unknown and must poison; state is %v", state)
	}

	second := followUpExecute(t, conn, "Execute after a mid-document deadline")
	if !errors.Is(second, ErrConnectionBroken) {
		t.Fatalf("Expected ErrConnectionBroken on the next call, got: %v", second)
	}
}

// TestSuccessfulTransactionClearsStaleDeadline: a deadline set by one
// transaction must not leak into the next one, which may have none.
func TestSuccessfulTransactionClearsStaleDeadline(t *testing.T) {
	conn, server := newPipeConn(t)
	echoServer(server)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	var err error
	mustFinish(t, 5*time.Second, "Execute with deadline", func() {
		err = conn.Execute(ctx, &simpleCmd{Foo: "first"}, &simpleResp{})
	})
	cancel()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Well after the first deadline would have expired, an un-deadlined call
	// must still work.
	time.Sleep(50 * time.Millisecond)

	resp := &simpleResp{}
	mustFinish(t, 5*time.Second, "Execute without deadline", func() {
		err = conn.Execute(context.Background(), &simpleCmd{Foo: "second"}, resp)
	})
	if err != nil {
		t.Fatalf("A stale deadline leaked into the next transaction: %v", err)
	}
	if resp.Foo != "second" {
		t.Fatalf("Expected Foo=second, got %q", resp.Foo)
	}
}

// ---------------------------------------------------------------------------
// entry-guard precedence: lifecycle and context outrank marshalling
// ---------------------------------------------------------------------------

// TestClosedConnectionOutranksMarshallingFailure: the marshalling failure is a
// property of one call, the closed connection a property of everything that
// follows — so the latter is what the caller needs to hear.
func TestClosedConnectionOutranksMarshallingFailure(t *testing.T) {
	conn, _ := newPipeConn(t)
	if err := conn.Close(); err != nil {
		t.Fatalf("Unexpected error from Close: %v", err)
	}

	cmd := &sideEffectCmd{}
	var err error
	mustFinish(t, 5*time.Second, "Execute on a closed connection", func() {
		err = conn.Execute(context.Background(), cmd, &simpleResp{})
	})

	if !errors.Is(err, ErrConnectionClosed) {
		t.Fatalf("Expected ErrConnectionClosed to outrank the marshalling error, got: %v", err)
	}
	if n := cmd.marshalled.Load(); n != 0 {
		t.Fatalf("A refused connection must not marshal the command; MarshalXML ran %d times", n)
	}
}

// TestBrokenConnectionOutranksMarshallingFailure is the same precedence rule for
// the other terminal state.
func TestBrokenConnectionOutranksMarshallingFailure(t *testing.T) {
	conn, server := newPipeConn(t)

	// Break it: a truncated document leaves the boundary unknown.
	go func() {
		buf := make([]byte, 65536)
		_, _ = server.Read(buf)
		_, _ = server.Write([]byte("<r><d>unterminated"))
		_ = server.Close()
	}()

	var first error
	mustFinish(t, 5*time.Second, "Execute against a truncated response", func() {
		first = conn.Execute(context.Background(), &simpleCmd{Foo: "x"}, &simpleResp{})
	})
	if first == nil {
		t.Fatalf("Expected the truncated response to fail")
	}

	conn.mu.Lock()
	state := conn.state
	conn.mu.Unlock()
	if state != stateBroken {
		t.Fatalf("Expected the connection to be broken, got state %v", state)
	}

	cmd := &sideEffectCmd{}
	var err error
	mustFinish(t, 5*time.Second, "Execute on a broken connection", func() {
		err = conn.Execute(context.Background(), cmd, &simpleResp{})
	})
	if !errors.Is(err, ErrConnectionBroken) {
		t.Fatalf("Expected ErrConnectionBroken to outrank the marshalling error, got: %v", err)
	}
	if n := cmd.marshalled.Load(); n != 0 {
		t.Fatalf("A refused connection must not marshal the command; MarshalXML ran %d times", n)
	}
}

// TestCancelledContextOutranksMarshallingFailure: same precedence for an
// expired context. Marshalling can be expensive and can run arbitrary user
// code, so it must not happen for a caller who has already given up.
func TestCancelledContextOutranksMarshallingFailure(t *testing.T) {
	conn, _ := newPipeConn(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cmd := &sideEffectCmd{}
	var err error
	mustFinish(t, 5*time.Second, "Execute with a pre-cancelled context", func() {
		err = conn.Execute(ctx, cmd, &simpleResp{})
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Expected context.Canceled to outrank the marshalling error, got: %v", err)
	}
	if n := cmd.marshalled.Load(); n != 0 {
		t.Fatalf("A cancelled caller must not marshal the command; MarshalXML ran %d times", n)
	}
}

// ---------------------------------------------------------------------------
// cancellation racing completion
// ---------------------------------------------------------------------------

// TestCancelRacingCompletionIsNotReportedAsSuccess pins the one cell of the
// outcome table that used to lie. The watchdog can only be joined after
// performRequest returns, so a cancellation landing in between poisons the
// connection while ioErr is still nil — and classifying on ioErr alone reported
// success while handing back a transport every later call would reject. The
// window cannot be closed, so the outcome is reconciled instead.
func TestCancelRacingCompletionIsNotReportedAsSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	raced := newCancelOnFinalReadConn(`<r><foo>done</foo></r>`, cancel)
	conn := New(raced)

	var err error
	resp := &simpleResp{}
	mustFinish(t, 5*time.Second, "Execute racing a cancellation", func() {
		err = conn.Execute(ctx, &simpleCmd{Foo: "x"}, resp)
	})

	if err == nil {
		t.Fatalf("The watchdog destroyed the connection; reporting success hides that entirely")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Expected the error to wrap context.Canceled, got: %v", err)
	}

	// The connection really is gone, which is the whole reason the error matters.
	conn.mu.Lock()
	state := conn.state
	conn.mu.Unlock()
	if state != stateBroken {
		t.Fatalf("Expected the watchdog to have poisoned the connection, got state %v", state)
	}

	second := followUpExecute(t, conn, "Execute after a cancellation raced completion")
	if !errors.Is(second, ErrConnectionBroken) {
		t.Fatalf("Expected ErrConnectionBroken on the next call, got: %v", second)
	}
}

// ---------------------------------------------------------------------------
// lifecycle is one-way
// ---------------------------------------------------------------------------

// TestClosedConnectionCannotBeReopened covers the invariant that replaced the
// old SetRawConn setter: a retired Connection cannot take a fresh socket, so
// recovery is always a new Connection. Resuming an ambiguously-failed stream is
// what would let the next call read the tail of a previous response.
func TestClosedConnectionCannotBeReopened(t *testing.T) {
	conn, _ := newPipeConn(t)
	if err := conn.Close(); err != nil {
		t.Fatalf("Unexpected error from Close: %v", err)
	}
	if err := followUpExecute(t, conn, "Execute after Close"); !errors.Is(err, ErrConnectionClosed) {
		t.Fatalf("Expected ErrConnectionClosed, got: %v", err)
	}

	// Recovery: a brand new Connection over a brand new socket, entirely
	// independent of the retired one.
	replacement, server := newPipeConn(t)
	echoServer(server)

	resp := &simpleResp{}
	var err error
	mustFinish(t, 5*time.Second, "Execute on the replacement", func() {
		err = replacement.Execute(context.Background(), &simpleCmd{Foo: "fresh"}, resp)
	})
	if err != nil {
		t.Fatalf("The replacement connection should work, got: %v", err)
	}
	if resp.Foo != "fresh" {
		t.Fatalf("Expected Foo=fresh, got %q", resp.Foo)
	}

	// And the retired one stays retired.
	if err := followUpExecute(t, conn, "Execute on the retired connection"); !errors.Is(err, ErrConnectionClosed) {
		t.Fatalf("Expected the retired connection to stay closed, got: %v", err)
	}
}
