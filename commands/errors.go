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
	ErrorTypeNetwork
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
	case ErrorTypeNetwork:
		return "network"
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
type GMPError struct {
	Type       ErrorType
	Status     string
	StatusText string
	StatusCode int
	Message    string
}

func (e *GMPError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("GMP error [%s]: %s (status: %s)", e.Type.String(), e.StatusText, e.Status)
}

func (e *GMPError) Unwrap() error { return nil }

// NetworkError represents a client-side networking or transport error.
type NetworkError struct {
	Type    ErrorType
	Message string
	Cause   error
}

func (e *NetworkError) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	if e.Message != "" {
		return e.Message
	}
	return e.Type.String()
}

func (e *NetworkError) Unwrap() error { return e.Cause }

// NewGMPError creates a GMPError from status and statusText.
func NewGMPError(status, statusText string) *GMPError {
	statusCode, _ := strconv.Atoi(status)

	var errorType ErrorType = ErrorTypeUnknown
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
		Message:    fmt.Sprintf("GMP error [%s]: %s (status: %s)", errorType.String(), statusText, status),
	}
}
