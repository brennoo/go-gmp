package commands

import "strconv"

// ResponseWithStatus is implemented by responses that expose status fields.
type ResponseWithStatus interface {
	GetStatus() (string, string)
}

// ValidateResponse checks if a GMP response indicates an error and returns a
// GMPError if so. It is the sole producer of *GMPError.
func ValidateResponse(status, statusText string) error {
	if status == "" {
		return nil // No status means success
	}

	// A non-numeric status is deliberately treated as success, not classified as
	// an error. Changing that would affect all ~123 command types at once with no
	// fixture corpus to verify against; tracked as T-7b, blocked on a live server.
	// Do not "fix" it here.
	statusCode, convErr := strconv.Atoi(status)
	if convErr == nil && statusCode >= 400 {
		return NewGMPError(status, statusText)
	}
	return nil
}
