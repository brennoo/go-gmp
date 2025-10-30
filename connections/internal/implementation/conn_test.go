package implementation

import (
	"context"
	"encoding/xml"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/brennoo/go-gmp/commands"
)

func TestSetRawConn(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := Connection{}
	conn.SetRawConn(c1)

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
		t.Fatalf("Unexpected error during read: %s", err)
	}

	if nRead != 10 {
		t.Fatalf("Expected to read 10 bytes, got: %d", nRead)
	}

	errClose := conn.Close()
	if errClose != nil {
		t.Fatalf("Unexpected error during Close: %s", errClose)
	}
}

func TestPerformRequest(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := Connection{}
	conn.SetRawConn(c1)

	go func() {
		buf2 := make([]byte, 150000)
		nRead, _ := c2.Read(buf2)
		c2.Write(buf2[:nRead])
	}()

	buf := make([]byte, 50000)
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i % 255)
	}

	bufRead, err := conn.performRequest(context.Background(), buf)
	if err != nil {
		t.Fatalf("Unexpected error during write: %s", err)
	}

	if len(bufRead) != 50000 {
		t.Fatalf("Expected to write 50000 bytes, wrote: %d", len(bufRead))
	}

	for i := 0; i < 50000; i++ {
		if bufRead[i] != byte(i%255) {
			t.Fatalf("Incorrect value at position %d: %d", i, bufRead[i])
		}
	}
}

func TestPerformRequestWriteFail(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := Connection{}
	conn.SetRawConn(c1)

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

	_, err := conn.performRequest(context.Background(), buf)
	expectedError := "io: read/write on closed pipe"
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Unexpected error during performRequest.\nExpected: %s\n     Got: %s", expectedError, err)
	}
}

func TestPerformRequestReadFail(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := Connection{}
	conn.SetRawConn(c1)

	go func() {
		buf2 := make([]byte, 150000)
		c2.Read(buf2)
		c2.Close()
	}()

	buf := make([]byte, 50000)
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i % 255)
	}

	_, err := conn.performRequest(context.Background(), buf)
	expectedError := "EOF"
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Unexpected error during performRequest.\nExpected: %s\n     Got: %s", expectedError, err)
	}
}

// Test that performRequest honors context deadlines via SetWriteDeadline.
func TestPerformRequestWriteDeadline(t *testing.T) {
	c1, _ := net.Pipe()

	conn := Connection{}
	conn.SetRawConn(c1)

	// Very short timeout; no reader on the other end, so Write should time out
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	buf := []byte("hello")
	_, err := conn.performRequest(ctx, buf)
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

	conn := Connection{}
	conn.SetRawConn(c1)

	go func() {
		// Read the request to let write succeed, but never write a response
		_ = make([]byte, 1024)
		buf := make([]byte, 1024)
		_, _ = c2.Read(buf)
		// do not write back -> read on c1 should block until deadline
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	_, err := conn.performRequest(ctx, []byte("request"))
	if err == nil {
		t.Fatalf("Expected read timeout error, got nil")
	}
	var nerr net.Error
	if !errors.As(err, &nerr) || !nerr.Timeout() {
		t.Fatalf("Expected net.Error timeout on read, got: %v", err)
	}
}

// Test that Execute closes the raw connection when context is cancelled.
func TestExecuteClosesOnCancel(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := Connection{}
	conn.SetRawConn(c1)

	// Backend just blocks on read to simulate a stalled peer
	done := make(chan struct{})
	go func() {
		defer close(done)
		buf := make([]byte, 1024)
		_, _ = c2.Read(buf)
	}()

	// Create a context that is cancelled immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cmd := &struct {
		XMLName xml.Name `xml:"cmd"`
		Foo     string   `xml:"foo"`
	}{Foo: "x"}
	resp := &struct {
		Foo string `xml:"foo"`
	}{}

	// Execute should return quickly and close the underlying connection
	_ = conn.Execute(ctx, cmd, resp)

	// Subsequent writes should fail due to closed connection
	if _, err := conn.rawConn.Write([]byte("test")); err == nil {
		t.Fatalf("Expected error writing to closed connection, got nil")
	}

	// ensure backend goroutine exits
	<-done
}

func TestExecute(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := Connection{}
	conn.SetRawConn(c1)

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

func TestExecuteXMLMarshallFail(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := Connection{}
	conn.SetRawConn(c1)

	go func() {
		buf2 := make([]byte, 150000)
		nRead, _ := c2.Read(buf2)
		c2.Write(buf2[:nRead])
	}()

	cmd := func() {}
	response := &struct {
		Foo string `xml:"foo"`
	}{}

	err := conn.Execute(context.Background(), cmd, response)
	expectedError := "xml: unsupported type: func()"
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Unexpected error during Execute.\nExpected: %s\n     Got: %s", expectedError, err)
	}
}

func TestExecuteXMLUnmarshallFail(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := Connection{}
	conn.SetRawConn(c1)

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
	response := struct {
		Foo string `xml:"foo"`
	}{}

	err := conn.Execute(context.Background(), cmd, response)
	expectedError := "non-pointer passed to Unmarshal"
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Unexpected error during Execute.\nExpected: %s\n     Got: %s", expectedError, err)
	}
}

func TestExecutePerformRequestFail(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := Connection{}
	conn.SetRawConn(c1)

	c1.Close()
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
	expectedError := "io: read/write on closed pipe"
	if err == nil || err.Error() != expectedError {
		t.Fatalf("Unexpected error during Execute.\nExpected: %s\n     Got: %s", expectedError, err)
	}
}

func TestRawXML(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := Connection{}
	conn.SetRawConn(c1)

	go func() {
		buf := make([]byte, 1024)
		n, _ := c2.Read(buf)
		// Echo back the same bytes as a response
		c2.Write(buf[:n])
	}()

	xml := `<raw_test_command foo="bar"/>`
	resp, err := conn.RawXML(xml)
	if err != nil {
		t.Fatalf("Unexpected error from RawXML: %v", err)
	}
	if resp != xml {
		t.Errorf("RawXML response did not match input.\nExpected: %s\nGot: %s", xml, resp)
	}
}

func TestRawXMLEmptyResponse(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := Connection{}
	conn.SetRawConn(c1)

	go func() {
		// Do not write anything back, simulating an empty response
		c2.Read(make([]byte, 1024)) // Read the request but don't respond
		c2.Close()
	}()

	xml := `<raw_test_command foo="bar"/>`
	resp, err := conn.RawXML(xml)
	if err == nil {
		t.Errorf("Expected error for empty response, got nil")
	}
	if resp != "" {
		t.Errorf("Expected empty string for response, got: %q", resp)
	}
}

// TestExecuteWithResponseWithStatus tests the new error handling for ResponseWithStatus.
func TestExecuteWithResponseWithStatus(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := Connection{}
	conn.SetRawConn(c1)

	go func() {
		buf2 := make([]byte, 150000)
		nRead, _ := c2.Read(buf2)
		c2.Write(buf2[:nRead])
	}()

	cmd := &struct {
		XMLName xml.Name `xml:"cmd"`
		Foo     string   `xml:"foo"`
	}{
		Foo: "test",
	}

	// Test with response that implements ResponseWithStatus and returns success
	response := &mockResponseWithStatus{
		Status:     "200",
		StatusText: "OK",
	}

	err := conn.Execute(context.Background(), cmd, response)
	if err != nil {
		t.Fatalf("Unexpected error during Execute with success response: %s", err)
	}
}

// TestExecuteWithResponseWithStatusError tests the new error handling for ResponseWithStatus with error.
func TestExecuteWithResponseWithStatusError(t *testing.T) {
	c1, c2 := net.Pipe()

	conn := Connection{}
	conn.SetRawConn(c1)

	go func() {
		buf2 := make([]byte, 150000)
		nRead, _ := c2.Read(buf2)
		c2.Write(buf2[:nRead])
	}()

	cmd := &struct {
		XMLName xml.Name `xml:"cmd"`
		Foo     string   `xml:"foo"`
	}{
		Foo: "test",
	}

	// Test with response that implements ResponseWithStatus and returns error
	response := &mockResponseWithStatus{
		Status:     "401",
		StatusText: "Unauthorized",
	}

	err := conn.Execute(context.Background(), cmd, response)
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
	c1, c2 := net.Pipe()

	conn := Connection{}
	conn.SetRawConn(c1)

	go func() {
		buf2 := make([]byte, 150000)
		nRead, _ := c2.Read(buf2)
		c2.Write(buf2[:nRead])
	}()

	cmd := &struct {
		XMLName xml.Name `xml:"cmd"`
		Foo     string   `xml:"foo"`
	}{
		Foo: "test",
	}

	// Test with response that doesn't implement ResponseWithStatus
	response := &struct {
		Foo string `xml:"foo"`
	}{}

	err := conn.Execute(context.Background(), cmd, response)
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
