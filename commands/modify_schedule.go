package commands

import "encoding/xml"

// ModifySchedule represents a modify_schedule command request.
type ModifySchedule struct {
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

// GetStatus returns the status and status text from the response.
func (r *ModifyScheduleResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
