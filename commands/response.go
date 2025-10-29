package commands

import "strconv"

// ResponseWithStatus is implemented by responses that expose status fields.
type ResponseWithStatus interface {
	GetStatus() (string, string)
}

// ValidateResponse checks if a GMP response indicates an error and returns a GMPError if so.
func ValidateResponse(status, statusText string) error {
	if status == "" {
		return nil // No status means success
	}

	statusCode, convErr := strconv.Atoi(status)
	if convErr == nil && statusCode >= 400 {
		return NewGMPError(status, statusText)
	}
	return nil
}
