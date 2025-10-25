package commands

import "encoding/xml"

// DeleteSchedule represents a delete_schedule command request.
type DeleteSchedule struct {
	XMLName    xml.Name `xml:"delete_schedule"`
	ScheduleID string   `xml:"schedule_id,attr"`
	Ultimate   bool     `xml:"ultimate,attr"`
}

// DeleteScheduleResponse represents a delete_schedule command response.
type DeleteScheduleResponse struct {
	XMLName    xml.Name `xml:"delete_schedule_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
