package commands

import (
	"encoding/xml"
	"strconv"
)

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

// genericEnvelope reads the status attributes off a <gmp_response/> element,
// the generic root gvmd reports some failures under instead of the command's
// own typed "..._response" root. XMLName is left unconstrained so Unmarshal
// doesn't reject other roots itself; ErrorFromDocument checks it explicitly.
type genericEnvelope struct {
	XMLName    xml.Name
	Status     string `xml:"status,attr"`
	StatusText string `xml:"status_text,attr"`
}

// gvmdGenericResponseRoot is the only root name confirmed live (T-7b in
// docs/review-2026-08.md).
const gvmdGenericResponseRoot = "gmp_response"

// ErrorFromDocument recovers a *GMPError from a document that failed to
// unmarshal into its expected typed response, by re-decoding it as a bare
// status envelope. Only a <gmp_response/> root is accepted; any other root
// means the failure happened elsewhere in a correctly-rooted document and
// must keep surfacing as the caller's original decode error.
//
// A nil result means doc is not a recoverable <gmp_response/> error — either
// it didn't decode as one, or it did but wasn't an error — not that doc was fine.
func ErrorFromDocument(doc []byte) *GMPError {
	var env genericEnvelope
	if xml.Unmarshal(doc, &env) != nil || env.XMLName.Local != gvmdGenericResponseRoot {
		return nil
	}
	if err := ValidateResponse(env.Status, env.StatusText); err != nil {
		return err.(*GMPError)
	}
	return nil
}
