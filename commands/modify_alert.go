package commands

import "encoding/xml"

// ModifyAlert represents a modify_alert command request.
type ModifyAlert struct {
	XMLName   xml.Name           `xml:"modify_alert"`
	AlertID   string             `xml:"alert_id,attr"`
	Name      *string            `xml:"name,omitempty"`
	Comment   *string            `xml:"comment,omitempty"`
	Filter    *ModifyAlertFilter `xml:"filter,omitempty"`
	Event     *event             `xml:"event,omitempty"`
	Condition *condition         `xml:"condition,omitempty"`
	Method    *method            `xml:"method,omitempty"`
}

// ModifyAlertFilter represents the filter element for modify_alert.
type ModifyAlertFilter struct {
	ID string `xml:"id,attr"`
}

// ModifyAlertResponse represents a modify_alert command response.
type ModifyAlertResponse struct {
	XMLName    xml.Name `xml:"modify_alert_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
