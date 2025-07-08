package gmp

import "encoding/xml"

// DeleteAlertCommand represents a delete_alert command request.
type DeleteAlertCommand struct {
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
