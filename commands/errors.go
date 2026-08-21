package commands

import (
	"fmt"
	"strconv"
)

// ErrorType represents the category of a GMP-related error.
type ErrorType int

const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeAuthentication
	ErrorTypeNotFound
	ErrorTypePermission
	ErrorTypeInvalidRequest
	ErrorTypeServerError
)

func (et ErrorType) String() string {
	switch et {
	case ErrorTypeAuthentication:
		return "authentication"
	case ErrorTypeNotFound:
		return "not_found"
	case ErrorTypePermission:
		return "permission"
	case ErrorTypeInvalidRequest:
		return "invalid_request"
	case ErrorTypeServerError:
		return "server_error"
	default:
		return "unknown"
	}
}

// GMPError represents an error returned by the GMP server.
//
// It is a leaf error: it carries no cause and does not implement Unwrap. The
// server's status and status text are the entire failure — there is nothing
// underneath to unwrap to.
type GMPError struct {
	Type       ErrorType
	Status     string
	StatusText string
	StatusCode int
}

func (e *GMPError) Error() string {
	return fmt.Sprintf("GMP error [%s]: %s (status: %s)", e.Type.String(), e.StatusText, e.Status)
}

// NewGMPError creates a GMPError from status and statusText.
func NewGMPError(status, statusText string) *GMPError {
	statusCode, _ := strconv.Atoi(status)

	errorType := ErrorTypeUnknown
	switch statusCode {
	case 400:
		errorType = ErrorTypeInvalidRequest
	case 401:
		errorType = ErrorTypeAuthentication
	case 403:
		errorType = ErrorTypePermission
	case 404:
		errorType = ErrorTypeNotFound
	case 500, 502, 503, 504:
		errorType = ErrorTypeServerError
	}

	return &GMPError{
		Type:       errorType,
		Status:     status,
		StatusText: statusText,
		StatusCode: statusCode,
	}
}

// ConnectionNotInitializedError means a Connection was used without an
// underlying socket — built as a bare struct literal instead of via its
// constructor. It's a misuse of the library, not a runtime failure, and
// retrying can never help.
//
// Deliberately not a TransportError: no transport was ever established, so
// neither answer to "is this connection usable" would be honest.
type ConnectionNotInitializedError struct {
	Cause error
}

func (e *ConnectionNotInitializedError) Error() string {
	return fmt.Sprintf("connection not initialized: %v", e.Cause)
}

func (e *ConnectionNotInitializedError) Unwrap() error { return e.Cause }

// RequestEncodingError means a command failed to marshal into the outbound
// XML request. The connection is never touched, so it is unaffected.
type RequestEncodingError struct {
	Cause error
}

func (e *RequestEncodingError) Error() string {
	return fmt.Sprintf("encode request: %v", e.Cause)
}

func (e *RequestEncodingError) Unwrap() error { return e.Cause }

// maxRetainedDocument bounds how much of a response ProtocolError retains.
const maxRetainedDocument = 4096

// ProtocolError means a response framed as a complete XML document but could
// not be unmarshalled into the caller's type — for example schema or version
// drift between client and server.
//
// The connection is always left usable: the document framed cleanly, so the
// stream's message-boundary position is known. Contrast TransportError, where
// it depends on ConnectionUsable.
type ProtocolError struct {
	Cause error
	// Document is a capped copy of the response that failed to decode.
	// Deliberately excluded from Error(): GMP responses routinely carry
	// sensitive material (credentials, report bodies) that shouldn't leak into
	// logs by accident. It's a fixed-length byte prefix, not a boundary-aware
	// cut, so it may be neither valid UTF-8 nor well-formed XML.
	Document []byte
	// DocumentSize is the full on-wire length of the response, even when
	// Document was truncated. Truncation is exactly len(Document) < DocumentSize.
	DocumentSize int
}

func (e *ProtocolError) Error() string {
	return fmt.Sprintf("decode response: %v", e.Cause)
}

func (e *ProtocolError) Unwrap() error { return e.Cause }

// NewProtocolError builds a ProtocolError over a copy of document (at most
// maxRetainedDocument bytes): aliasing the caller's slice would otherwise pin
// a potentially multi-megabyte backing array alive for the error's lifetime.
func NewProtocolError(document []byte, cause error) *ProtocolError {
	retained := make([]byte, min(len(document), maxRetainedDocument))
	copy(retained, document)

	return &ProtocolError{
		Cause:        cause,
		Document:     retained,
		DocumentSize: len(document),
	}
}

// TransportError means the transport layer did not complete the transaction:
// the socket-level exchange failed, a response could not be framed as a
// complete XML document, or the request was refused outright.
//
// Check ConnectionUsable before discarding the connection. Not every
// TransportError is fatal to it, and connections are shared.
type TransportError struct {
	Cause error

	// ConnectionUsable reports whether the connection can still carry further
	// requests. True only when the request was refused before the socket was
	// touched (e.g. an already-done context), leaving it exactly as it was.
	// False means replace it: the stream's message-boundary position is
	// unknown or the connection is already closed/broken, and neither can be
	// revived.
	//
	// The zero value is deliberately the conservative "do not reuse" answer.
	ConnectionUsable bool
}

func (e *TransportError) Error() string {
	return fmt.Sprintf("transport: %v", e.Cause)
}

func (e *TransportError) Unwrap() error { return e.Cause }
