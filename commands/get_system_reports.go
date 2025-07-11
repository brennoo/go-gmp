package commands

import "encoding/xml"

// GetSystemReports represents a get_system_reports command request.
type GetSystemReports struct {
	XMLName   xml.Name `xml:"get_system_reports"`
	Name      string   `xml:"name,attr,omitempty"`
	Duration  string   `xml:"duration,attr,omitempty"`
	StartTime string   `xml:"start_time,attr,omitempty"`
	EndTime   string   `xml:"end_time,attr,omitempty"`
	Brief     string   `xml:"brief,attr,omitempty"`
	SlaveID   string   `xml:"slave_id,attr,omitempty"`
}

// GetSystemReportsResponse represents a get_system_reports command response.
type GetSystemReportsResponse struct {
	XMLName       xml.Name       `xml:"get_system_reports_response"`
	Status        string         `xml:"status,attr"`
	StatusText    string         `xml:"status_text,attr"`
	SystemReports []SystemReport `xml:"system_report"`
}

// SystemReport represents a <system_report> element in the get_system_reports response.
type SystemReport struct {
	Name   string              `xml:"name"`
	Title  string              `xml:"title"`
	Report *SystemReportReport `xml:"report,omitempty"`
}

// SystemReportReport represents the content of a system report.
type SystemReportReport struct {
	Format   string `xml:"format,attr"`
	Duration string `xml:"duration,attr"`
	Value    string `xml:",chardata"`
}
