package commands

import "encoding/xml"

// GetInfo represents a get_info command request.
type GetInfo struct {
	XMLName xml.Name `xml:"get_info"`
	Name    string   `xml:"name,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Details string   `xml:"details,attr,omitempty"`
}

// GetInfoResponse represents a get_info command response.
type GetInfoResponse struct {
	XMLName    xml.Name      `xml:"get_info_response"`
	Status     string        `xml:"status,attr"`
	StatusText string        `xml:"status_text,attr"`
	Infos      []GetInfoInfo `xml:"info"`
}

type GetInfoInfo struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}
