package commands

import "encoding/xml"

// TestAlert represents a test_alert command request.
type TestAlert struct {
	XMLName xml.Name `xml:"test_alert"`
	AlertID string   `xml:"alert_id,attr"`
}

// TestAlertResponse represents a test_alert command response.
type TestAlertResponse struct {
	XMLName       xml.Name `xml:"test_alert_response"`
	Status        string   `xml:"status,attr"`
	StatusText    string   `xml:"status_text,attr"`
	StatusDetails string   `xml:"status_details,omitempty"`
}
