package gmp

import "encoding/xml"

// ResumeTaskCommand represents a resume_task command request.
type ResumeTaskCommand struct {
	XMLName xml.Name `xml:"resume_task"`
	TaskID  string   `xml:"task_id,attr"`
}

// ResumeTaskResponse represents a resume_task command response.
type ResumeTaskResponse struct {
	XMLName    xml.Name `xml:"resume_task_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ReportID   string   `xml:"report_id,omitempty"`
}
