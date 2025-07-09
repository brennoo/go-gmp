package commands

import "encoding/xml"

// ModifyAlertCommand represents a modify_alert command request.
type ModifyAlertCommand struct {
	XMLName   xml.Name           `xml:"modify_alert"`
	AlertID   string             `xml:"alert_id,attr"`
	Name      *string            `xml:"name,omitempty"`
	Comment   *string            `xml:"comment,omitempty"`
	Filter    *ModifyAlertFilter `xml:"filter,omitempty"`
	Event     *AlertEvent        `xml:"event,omitempty"`
	Condition *AlertCondition    `xml:"condition,omitempty"`
	Method    *AlertMethod       `xml:"method,omitempty"`
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
