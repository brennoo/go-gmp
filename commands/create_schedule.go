package commands

import "encoding/xml"

// CreateSchedule represents a create_schedule command request.
type CreateSchedule struct {
	XMLName   xml.Name `xml:"create_schedule"`
	Name      string   `xml:"name"`
	Comment   string   `xml:"comment,omitempty"`
	Copy      string   `xml:"copy,omitempty"`
	ICalendar string   `xml:"icalendar,omitempty"`
	Timezone  string   `xml:"timezone,omitempty"`
}

// CreateScheduleResponse represents a create_schedule command response.
type CreateScheduleResponse struct {
	XMLName    xml.Name `xml:"create_schedule_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *CreateScheduleResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
