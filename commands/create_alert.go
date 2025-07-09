package commands

import "encoding/xml"

// CreateAlert represents a create_alert command request.
type CreateAlert struct {
	XMLName   xml.Name  `xml:"create_alert"`
	Name      string    `xml:"name"`
	Comment   string    `xml:"comment,omitempty"`
	Copy      string    `xml:"copy,omitempty"`
	Condition condition `xml:"condition"`
	Event     event     `xml:"event"`
	Method    method    `xml:"method"`
	Filter    *filter   `xml:"filter,omitempty"`
	Active    *bool     `xml:"active,omitempty"`
}

// CreateAlertResponse represents a create_alert command response.
type CreateAlertResponse struct {
	XMLName    xml.Name `xml:"create_alert_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// NewCreateAlert creates a new CreateAlert command.
func NewCreateAlert(
	name, comment, copy string,
	condition condition,
	event event,
	method method,
	filter *filter,
	active *bool,
) *CreateAlert {
	return &CreateAlert{
		Name:      name,
		Comment:   comment,
		Copy:      copy,
		Condition: condition,
		Event:     event,
		Method:    method,
		Filter:    filter,
		Active:    active,
	}
}

// AlertData represents data for alert condition, event, or method.
type AlertData struct {
	Name string `xml:"name"`
	Text string `xml:",chardata"`
}

// NewAlertData creates a new AlertData.
func NewAlertData(name, text string) AlertData {
	return AlertData{Name: name, Text: text}
}

// condition represents the condition that must be satisfied for the alert to occur.
type condition struct {
	Text string      `xml:",chardata"`
	Data []AlertData `xml:"data,omitempty"`
}

// NewAlertCondition creates a new condition for alerts.
func NewAlertCondition(text string, data []AlertData) condition {
	return condition{Text: text, Data: data}
}

// event represents the event that must happen for the alert to occur.
type event struct {
	Text string      `xml:",chardata"`
	Data []AlertData `xml:"data,omitempty"`
}

// NewAlertEvent creates a new event for alerts.
func NewAlertEvent(text string, data []AlertData) event {
	return event{Text: text, Data: data}
}

// method represents the method by which the user is alerted.
type method struct {
	Text string      `xml:",chardata"`
	Data []AlertData `xml:"data,omitempty"`
}

// NewAlertMethod creates a new method for alerts.
func NewAlertMethod(text string, data []AlertData) method {
	return method{Text: text, Data: data}
}

// filter represents a filter to apply when executing alert.
type filter struct {
	ID string `xml:"id,attr"`
}
