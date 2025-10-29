package commands

import "encoding/xml"

// CreateConfig represents a create_config command request.
type CreateConfig struct {
	XMLName            xml.Name            `xml:"create_config"`
	Comment            string              `xml:"comment,omitempty"`
	Copy               string              `xml:"copy,omitempty"`
	GetConfigsResponse *GetConfigsResponse `xml:"get_configs_response,omitempty"`
	Name               string              `xml:"name,omitempty"`
	UsageType          string              `xml:"usage_type,omitempty"`
}

// CreateConfigResponse represents a create_config command response.
type CreateConfigResponse struct {
	XMLName    xml.Name `xml:"create_config_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *CreateConfigResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
