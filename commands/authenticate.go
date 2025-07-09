package commands

import "encoding/xml"

// Authenticate represents an authenticate command request.
type Authenticate struct {
	XMLName     xml.Name     `xml:"authenticate"`
	Credentials *credentials `xml:"credentials"`
}

// credentials contains the username and password for authentication.
type credentials struct {
	Username string `xml:"username"`
	Password string `xml:"password"`
}

// AuthenticateResponse represents an authenticate command response.
type AuthenticateResponse struct {
	XMLName    xml.Name `xml:"authenticate_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	Role       string   `xml:"role,omitempty"`
	Timezone   string   `xml:"timezone,omitempty"`
}

// NewAuthenticateCredentials creates a new Authenticate command with the given username and password.
func NewAuthenticateCredentials(username, password string) *Authenticate {
	return &Authenticate{
		Credentials: &credentials{
			Username: username,
			Password: password,
		},
	}
}
