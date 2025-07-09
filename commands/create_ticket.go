package commands

import "encoding/xml"

// CreateTicketCommand represents a create_ticket command request.
type CreateTicketCommand struct {
	XMLName  xml.Name               `xml:"create_ticket"`
	Comment  string                 `xml:"comment,omitempty"`
	Copy     string                 `xml:"copy,omitempty"`
	Result   CreateTicketResult     `xml:"result"`
	Assigned CreateTicketAssignedTo `xml:"assigned_to"`
	OpenNote string                 `xml:"open_note"`
}

type CreateTicketResult struct {
	ID string `xml:"id,attr"`
}

type CreateTicketAssignedTo struct {
	User CreateTicketUser `xml:"user"`
}

type CreateTicketUser struct {
	ID string `xml:"id,attr"`
}

// CreateTicketResponse represents a create_ticket command response.
type CreateTicketResponse struct {
	XMLName    xml.Name `xml:"create_ticket_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
