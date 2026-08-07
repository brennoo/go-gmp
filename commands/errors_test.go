package commands

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func TestValidateResponse_Success(t *testing.T) {
	if err := ValidateResponse("200", "OK"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateResponse_ClientError(t *testing.T) {
	err := ValidateResponse("401", "Unauthorized")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	var ge *GMPError
	if !errors.As(err, &ge) {
		t.Fatalf("expected *GMPError, got %T", err)
	}
	if ge.Type != ErrorTypeAuthentication {
		t.Fatalf("expected authentication error, got %v", ge.Type)
	}
}

func TestValidateResponse_NotFound(t *testing.T) {
	err := ValidateResponse("404", "Not Found")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	var ge *GMPError
	if !errors.As(err, &ge) {
		t.Fatalf("expected *GMPError, got %T", err)
	}
	if ge.Type != ErrorTypeNotFound {
		t.Fatalf("expected not_found error, got %v", ge.Type)
	}
}

func TestValidateResponse_ServerError(t *testing.T) {
	err := ValidateResponse("500", "Internal Server Error")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	var ge *GMPError
	if !errors.As(err, &ge) {
		t.Fatalf("expected *GMPError, got %T", err)
	}
	if ge.Type != ErrorTypeServerError {
		t.Fatalf("expected server_error, got %v", ge.Type)
	}
}

func TestValidateResponse_InvalidRequest(t *testing.T) {
	err := ValidateResponse("400", "Bad Request")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	var ge *GMPError
	if !errors.As(err, &ge) {
		t.Fatalf("expected *GMPError, got %T", err)
	}
	if ge.Type != ErrorTypeInvalidRequest {
		t.Fatalf("expected invalid_request error, got %v", ge.Type)
	}
}

func TestValidateResponse_Permission(t *testing.T) {
	err := ValidateResponse("403", "Forbidden")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	var ge *GMPError
	if !errors.As(err, &ge) {
		t.Fatalf("expected *GMPError, got %T", err)
	}
	if ge.Type != ErrorTypePermission {
		t.Fatalf("expected permission error, got %v", ge.Type)
	}
}

func TestValidateResponse_NonNumericStatus(t *testing.T) {
	// Non-numeric status should not return an error
	if err := ValidateResponse("abc", "OK"); err != nil {
		t.Fatalf("expected no error for non-numeric status, got %v", err)
	}
}

func TestValidateResponse_EdgeCases(t *testing.T) {
	// Test edge cases
	testCases := []struct {
		status     string
		statusText string
		expectErr  bool
		errorType  ErrorType
	}{
		{"200", "OK", false, ErrorTypeUnknown},
		{"201", "Created", false, ErrorTypeUnknown},
		{"202", "Accepted", false, ErrorTypeUnknown},
		{"399", "Client Error", false, ErrorTypeUnknown},
		{"400", "Bad Request", true, ErrorTypeInvalidRequest},
		{"401", "Unauthorized", true, ErrorTypeAuthentication},
		{"403", "Forbidden", true, ErrorTypePermission},
		{"404", "Not Found", true, ErrorTypeNotFound},
		{"409", "Conflict", true, ErrorTypeUnknown},
		{"500", "Internal Server Error", true, ErrorTypeServerError},
		{"502", "Bad Gateway", true, ErrorTypeServerError},
		{"503", "Service Unavailable", true, ErrorTypeServerError},
		{"504", "Gateway Timeout", true, ErrorTypeServerError},
		{"599", "Server Error", true, ErrorTypeUnknown},
	}

	for _, tc := range testCases {
		t.Run(tc.status, func(t *testing.T) {
			err := ValidateResponse(tc.status, tc.statusText)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error for status %s, got nil", tc.status)
				}
				var ge *GMPError
				if !errors.As(err, &ge) {
					t.Fatalf("expected *GMPError for status %s, got %T", tc.status, err)
				}
				if ge.Type != tc.errorType {
					t.Fatalf("expected error type %v for status %s, got %v", tc.errorType, tc.status, ge.Type)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error for status %s, got %v", tc.status, err)
				}
			}
		})
	}
}

func TestErrorType_String(t *testing.T) {
	testCases := []struct {
		errorType ErrorType
		expected  string
	}{
		{ErrorTypeUnknown, "unknown"},
		{ErrorTypeAuthentication, "authentication"},
		{ErrorTypeNotFound, "not_found"},
		{ErrorTypePermission, "permission"},
		{ErrorTypeInvalidRequest, "invalid_request"},
		{ErrorTypeServerError, "server_error"},
		{ErrorType(999), "unknown"}, // Test unknown error type
	}

	for _, tc := range testCases {
		if got := tc.errorType.String(); got != tc.expected {
			t.Errorf("ErrorType(%d).String() = %q, want %q", tc.errorType, got, tc.expected)
		}
	}
}

func TestGMPError_Error(t *testing.T) {
	err := &GMPError{
		Type:       ErrorTypeNotFound,
		Status:     "404",
		StatusText: "Not Found",
	}
	expected := "GMP error [not_found]: Not Found (status: 404)"
	if got := err.Error(); got != expected {
		t.Errorf("GMPError.Error() = %q, want %q", got, expected)
	}
}

func TestRequestEncodingError_Error(t *testing.T) {
	cause := errors.New("underlying error")
	err := &RequestEncodingError{Cause: cause}
	expected := "encode request: underlying error"
	if got := err.Error(); got != expected {
		t.Errorf("RequestEncodingError.Error() = %q, want %q", got, expected)
	}
}

func TestRequestEncodingError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &RequestEncodingError{Cause: cause}
	if got := err.Unwrap(); !errors.Is(got, cause) {
		t.Errorf("RequestEncodingError.Unwrap() = %v, want %v", got, cause)
	}
}

func TestRequestEncodingError_ErrorsIs(t *testing.T) {
	sentinel := errors.New("sentinel")
	err := fmt.Errorf("wrapped: %w", &RequestEncodingError{Cause: sentinel})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected errors.Is to reach sentinel through RequestEncodingError")
	}
}

func TestProtocolError_Error(t *testing.T) {
	cause := errors.New("underlying error")
	err := &ProtocolError{Cause: cause}
	expected := "decode response: underlying error"
	if got := err.Error(); got != expected {
		t.Errorf("ProtocolError.Error() = %q, want %q", got, expected)
	}
}

func TestProtocolError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &ProtocolError{Cause: cause}
	if got := err.Unwrap(); !errors.Is(got, cause) {
		t.Errorf("ProtocolError.Unwrap() = %v, want %v", got, cause)
	}
}

func TestProtocolError_ErrorsIs(t *testing.T) {
	sentinel := errors.New("sentinel")
	err := fmt.Errorf("wrapped: %w", &ProtocolError{Cause: sentinel})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected errors.Is to reach sentinel through ProtocolError")
	}
}

func TestTransportError_Error(t *testing.T) {
	cause := errors.New("underlying error")
	err := &TransportError{Cause: cause}
	expected := "transport: underlying error"
	if got := err.Error(); got != expected {
		t.Errorf("TransportError.Error() = %q, want %q", got, expected)
	}
}

func TestTransportError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &TransportError{Cause: cause}
	if got := err.Unwrap(); !errors.Is(got, cause) {
		t.Errorf("TransportError.Unwrap() = %v, want %v", got, cause)
	}
}

func TestTransportError_ErrorsIs(t *testing.T) {
	sentinel := errors.New("sentinel")
	err := fmt.Errorf("wrapped: %w", &TransportError{Cause: sentinel})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected errors.Is to reach sentinel through TransportError")
	}
}

// A nil Cause must render as "<nil>" rather than panic, which is why Error()
// uses fmt.Sprintf("%v", ...) instead of e.Cause.Error().
func TestNewTypes_NilCauseDoesNotPanic(t *testing.T) {
	if got := (&RequestEncodingError{}).Error(); got != "encode request: <nil>" {
		t.Errorf("RequestEncodingError.Error() with nil cause = %q", got)
	}
	if got := (&ProtocolError{}).Error(); got != "decode response: <nil>" {
		t.Errorf("ProtocolError.Error() with nil cause = %q", got)
	}
	if got := (&TransportError{}).Error(); got != "transport: <nil>" {
		t.Errorf("TransportError.Error() with nil cause = %q", got)
	}
}

func TestNewProtocolError_TruncatesAndCopies(t *testing.T) {
	source := make([]byte, maxRetainedDocument+1000)
	for i := range source {
		source[i] = byte('a' + i%26)
	}
	cause := errors.New("decode failed")

	err := NewProtocolError(source, cause)

	if len(err.Document) != maxRetainedDocument {
		t.Fatalf("expected Document to be truncated to %d bytes, got %d", maxRetainedDocument, len(err.Document))
	}
	if err.DocumentSize != len(source) {
		t.Fatalf("expected DocumentSize %d, got %d", len(source), err.DocumentSize)
	}

	// Genuine copy: mutating the source after the fact must not affect Document.
	before := make([]byte, len(err.Document))
	copy(before, err.Document)
	for i := range source {
		source[i] = 0
	}
	if !bytes.Equal(before, err.Document) {
		t.Fatalf("expected Document to be unaffected by mutating the source slice")
	}
}

func TestNewGMPError(t *testing.T) {
	testCases := []struct {
		status       string
		statusText   string
		expectedType ErrorType
		expectedCode int
	}{
		{"400", "Bad Request", ErrorTypeInvalidRequest, 400},
		{"401", "Unauthorized", ErrorTypeAuthentication, 401},
		{"403", "Forbidden", ErrorTypePermission, 403},
		{"404", "Not Found", ErrorTypeNotFound, 404},
		{"500", "Internal Server Error", ErrorTypeServerError, 500},
		{"502", "Bad Gateway", ErrorTypeServerError, 502},
		{"503", "Service Unavailable", ErrorTypeServerError, 503},
		{"504", "Gateway Timeout", ErrorTypeServerError, 504},
		{"999", "Unknown Error", ErrorTypeUnknown, 999},
		{"abc", "Invalid Status", ErrorTypeUnknown, 0}, // Non-numeric status
	}

	for _, tc := range testCases {
		t.Run(tc.status, func(t *testing.T) {
			err := NewGMPError(tc.status, tc.statusText)
			if err.Type != tc.expectedType {
				t.Errorf("NewGMPError(%q, %q).Type = %v, want %v", tc.status, tc.statusText, err.Type, tc.expectedType)
			}
			if err.StatusCode != tc.expectedCode {
				t.Errorf("NewGMPError(%q, %q).StatusCode = %d, want %d", tc.status, tc.statusText, err.StatusCode, tc.expectedCode)
			}
			if err.Status != tc.status {
				t.Errorf("NewGMPError(%q, %q).Status = %q, want %q", tc.status, tc.statusText, err.Status, tc.status)
			}
			if err.StatusText != tc.statusText {
				t.Errorf("NewGMPError(%q, %q).StatusText = %q, want %q", tc.status, tc.statusText, err.StatusText, tc.statusText)
			}
		})
	}
}

func TestValidateResponse_EdgeCases_Response(t *testing.T) {
	// Test the ValidateResponse function from response.go
	testCases := []struct {
		status     string
		statusText string
		expectErr  bool
		errorType  ErrorType
	}{
		{"200", "OK", false, ErrorTypeUnknown},
		{"201", "Created", false, ErrorTypeUnknown},
		{"202", "Accepted", false, ErrorTypeUnknown},
		{"399", "Client Error", false, ErrorTypeUnknown},
		{"400", "Bad Request", true, ErrorTypeInvalidRequest},
		{"401", "Unauthorized", true, ErrorTypeAuthentication},
		{"403", "Forbidden", true, ErrorTypePermission},
		{"404", "Not Found", true, ErrorTypeNotFound},
		{"409", "Conflict", true, ErrorTypeUnknown},
		{"500", "Internal Server Error", true, ErrorTypeServerError},
		{"502", "Bad Gateway", true, ErrorTypeServerError},
		{"503", "Service Unavailable", true, ErrorTypeServerError},
		{"504", "Gateway Timeout", true, ErrorTypeServerError},
		{"599", "Server Error", true, ErrorTypeUnknown},
		{"abc", "Invalid Status", false, ErrorTypeUnknown}, // Non-numeric status
		{"", "Empty Status", false, ErrorTypeUnknown},      // Empty status
	}

	for _, tc := range testCases {
		t.Run(tc.status, func(t *testing.T) {
			err := ValidateResponse(tc.status, tc.statusText)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error for status %s, got nil", tc.status)
				}
				var ge *GMPError
				if !errors.As(err, &ge) {
					t.Fatalf("expected *GMPError for status %s, got %T", tc.status, err)
				}
				if ge.Type != tc.errorType {
					t.Fatalf("expected error type %v for status %s, got %v", tc.errorType, tc.status, ge.Type)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error for status %s, got %v", tc.status, err)
				}
			}
		})
	}
}
