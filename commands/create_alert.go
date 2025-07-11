package commands

import "encoding/xml"

// CreateAlert represents a create_alert command request.
type CreateAlert struct {
	XMLName   xml.Name       `xml:"create_alert"`
	Name      string         `xml:"name"`
	Comment   string         `xml:"comment,omitempty"`
	Copy      string         `xml:"copy,omitempty"`
	Condition AlertCondition `xml:"condition"`
	Event     AlertEvent     `xml:"event"`
	Method    AlertMethod    `xml:"method"`
	Filter    *AlertFilter   `xml:"filter,omitempty"`
	Active    string         `xml:"active,omitempty"`
}

// CreateAlertResponse represents a create_alert command response.
type CreateAlertResponse struct {
	XMLName    xml.Name `xml:"create_alert_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// AlertData represents data for alert condition, event, or method.
type AlertData struct {
	Name string `xml:"name"`
	Text string `xml:",chardata"`
}

// AlertCondition represents the condition that must be satisfied for the alert to occur.
type AlertCondition struct {
	Text string      `xml:",chardata"`
	Data []AlertData `xml:"data,omitempty"`
}

// AlertEvent represents the event that must happen for the alert to occur.
type AlertEvent struct {
	Text string      `xml:",chardata"`
	Data []AlertData `xml:"data,omitempty"`
}

// AlertMethod represents the method by which the user is alerted.
type AlertMethod struct {
	Text string      `xml:",chardata"`
	Data []AlertData `xml:"data,omitempty"`
}

// AlertFilter represents a filter to apply when executing alert.
type AlertFilter struct {
	ID string `xml:"id,attr"`
}
