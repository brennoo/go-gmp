package implementation

import (
	"bytes"
	"context"
	"encoding/xml"
	"net"
	"time"

	"github.com/brennoo/go-gmp/commands"
)

type Connection struct {
	rawConn net.Conn
}

func (conn *Connection) SetRawConn(rawConn net.Conn) {
	conn.rawConn = rawConn
}

func (conn *Connection) performRequest(ctx context.Context, buffer []byte) ([]byte, error) {
	// Set deadlines if context has a deadline
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.rawConn.SetWriteDeadline(deadline)
	}
	_, err := conn.rawConn.Write(buffer)
	if err != nil {
		return nil, err
	}

	var resp bytes.Buffer
	for true {
		if deadline, ok := ctx.Deadline(); ok {
			_ = conn.rawConn.SetReadDeadline(deadline)
		} else {
			// Clear previous deadlines
			_ = conn.rawConn.SetReadDeadline(time.Time{})
		}
		buf := make([]byte, 1024)
		n, err := conn.rawConn.Read(buf[:])
		if err != nil {
			return nil, err
		}
		resp.Write(buf[:n])

		if n < 1024 {
			break
		}
	}

	return resp.Bytes(), nil
}

func (conn *Connection) Execute(ctx context.Context, command interface{}, response interface{}) error {
	cmdBuf, err := xml.Marshal(command)
	if err != nil {
		return &commands.NetworkError{Type: commands.ErrorTypeInvalidRequest, Cause: err}
	}

	// Handle context cancellation by closing the connection in a separate goroutine
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			_ = conn.rawConn.Close()
		case <-done:
		}
	}()

	resp, err := conn.performRequest(ctx, cmdBuf)
	if err != nil {
		return &commands.NetworkError{Type: commands.ErrorTypeNetwork, Cause: err}
	}

	err = xml.Unmarshal(resp, response)
	if err != nil {
		return &commands.NetworkError{Type: commands.ErrorTypeInvalidRequest, Cause: err}
	}

	// Check response status if available
	if rs, ok := response.(commands.ResponseWithStatus); ok {
		status, statusText := rs.GetStatus()
		if gmpErr := commands.ValidateResponse(status, statusText); gmpErr != nil {
			return gmpErr
		}
	}

	return nil
}

func (conn *Connection) RawXML(xmlStr string) (string, error) {
	resp, err := conn.performRequest(context.Background(), []byte(xmlStr))
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func (conn *Connection) Close() error {
	return conn.rawConn.Close()
}
