package commands

import "encoding/xml"

// ResumeTask represents a resume_task command request.
type ResumeTask struct {
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

// GetStatus returns the status and status text from the response.
func (r *ResumeTaskResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
