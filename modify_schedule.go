package gmp

import "encoding/xml"

// ModifyScheduleCommand represents a modify_schedule command request.
type ModifyScheduleCommand struct {
	XMLName    xml.Name `xml:"modify_schedule"`
	ScheduleID string   `xml:"schedule_id,attr"`
	Comment    string   `xml:"comment,omitempty"`
	Name       string   `xml:"name,omitempty"`
	ICalendar  string   `xml:"icalendar,omitempty"`
	Timezone   string   `xml:"timezone,omitempty"`
}

// ModifyScheduleResponse represents a modify_schedule command response.
type ModifyScheduleResponse struct {
	XMLName    xml.Name `xml:"modify_schedule_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
