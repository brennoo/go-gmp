package commands

import "encoding/xml"

// Authenticate represents an authenticate command request.
type Authenticate struct {
	XMLName     xml.Name                 `xml:"authenticate"`
	Credentials *AuthenticateCredentials `xml:"credentials"`
}

// AuthenticateResponse represents an authenticate command response.
type AuthenticateResponse struct {
	XMLName    xml.Name `xml:"authenticate_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	Role       string   `xml:"role,omitempty"`
	Timezone   string   `xml:"timezone,omitempty"`
}

// GetStatus returns the status and status text from the response.
func (r *AuthenticateResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}

// AuthenticateCredentials contains the username and password for authentication.
type AuthenticateCredentials struct {
	Username string `xml:"username"`
	Password string `xml:"password"`
}
