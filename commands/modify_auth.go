package commands

import "encoding/xml"

// ModifyAuth represents a modify_auth command request.
type ModifyAuth struct {
	XMLName xml.Name  `xml:"modify_auth"`
	Group   AuthGroup `xml:"group"`
}

// ModifyAuthResponse represents a modify_auth command response.
type ModifyAuthResponse struct {
	XMLName    xml.Name `xml:"modify_auth_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
