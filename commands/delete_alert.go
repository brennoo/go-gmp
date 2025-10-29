package commands

import "encoding/xml"

// DeleteAlert represents a delete_alert command request.
type DeleteAlert struct {
	XMLName  xml.Name `xml:"delete_alert"`
	AlertID  string   `xml:"alert_id,attr"`
	Ultimate bool     `xml:"ultimate,attr"`
}

// DeleteAlertResponse represents a delete_alert command response.
type DeleteAlertResponse struct {
	XMLName    xml.Name `xml:"delete_alert_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *DeleteAlertResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
