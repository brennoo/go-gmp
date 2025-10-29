package commands

import "encoding/xml"

// ModifySetting represents a modify_setting command request.
type ModifySetting struct {
	XMLName   xml.Name `xml:"modify_setting"`
	SettingID string   `xml:"setting_id,attr,omitempty"`
	Name      string   `xml:"name"`
	Value     string   `xml:"value"`
}

// ModifySettingResponse represents a modify_setting command response.
type ModifySettingResponse struct {
	XMLName    xml.Name `xml:"modify_setting_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *ModifySettingResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
