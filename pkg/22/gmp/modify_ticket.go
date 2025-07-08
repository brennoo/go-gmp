package gmp

import "encoding/xml"

// ModifyTicketCommand represents a modify_ticket command request.
type ModifyTicketCommand struct {
	XMLName    xml.Name                `xml:"modify_ticket"`
	TicketID   string                  `xml:"ticket_id,attr"`
	Comment    string                  `xml:"comment,omitempty"`
	Status     string                  `xml:"status,omitempty"`
	OpenNote   string                  `xml:"open_note,omitempty"`
	FixedNote  string                  `xml:"fixed_note,omitempty"`
	ClosedNote string                  `xml:"closed_note,omitempty"`
	Assigned   *ModifyTicketAssignedTo `xml:"assigned_to,omitempty"`
}

type ModifyTicketAssignedTo struct {
	User ModifyTicketUser `xml:"user"`
}

type ModifyTicketUser struct {
	ID string `xml:"id,attr"`
}

// ModifyTicketResponse represents a modify_ticket command response.
type ModifyTicketResponse struct {
	XMLName    xml.Name `xml:"modify_ticket_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
