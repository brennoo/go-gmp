package commands

import "encoding/xml"

// ModifyAlert represents a modify_alert command request.
type ModifyAlert struct {
	XMLName   xml.Name        `xml:"modify_alert"`
	AlertID   string          `xml:"alert_id,attr"`
	Name      *string         `xml:"name,omitempty"`
	Comment   *string         `xml:"comment,omitempty"`
	Filter    *AlertFilter    `xml:"filter,omitempty"`
	Event     *AlertEvent     `xml:"event,omitempty"`
	Condition *AlertCondition `xml:"condition,omitempty"`
	Method    *AlertMethod    `xml:"method,omitempty"`
}

// ModifyAlertResponse represents a modify_alert command response.
type ModifyAlertResponse struct {
	XMLName    xml.Name `xml:"modify_alert_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
