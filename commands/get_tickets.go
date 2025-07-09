package commands

import "encoding/xml"

// GetTickets represents a get_tickets command request.
type GetTickets struct {
	XMLName xml.Name `xml:"get_tickets"`
}

// GetTicketsResponse represents a get_tickets command response.
type GetTicketsResponse struct {
	XMLName    xml.Name      `xml:"get_tickets_response"`
	Status     string        `xml:"status,attr"`
	StatusText string        `xml:"status_text,attr"`
	Tickets    []TicketEntry `xml:"ticket"`
}

// TicketEntry represents a ticket entry in the get_tickets response.
type TicketEntry struct {
	ID               string             `xml:"id,attr"`
	Owner            *TicketOwner       `xml:"owner,omitempty"`
	Name             string             `xml:"name"`
	Comment          string             `xml:"comment"`
	CreationTime     string             `xml:"creation_time"`
	ModificationTime string             `xml:"modification_time"`
	Writable         string             `xml:"writable"`
	InUse            string             `xml:"in_use"`
	Permissions      *TicketPermissions `xml:"permissions,omitempty"`
	AssignedTo       *TicketAssignedTo  `xml:"assigned_to,omitempty"`
	TaskID           string             `xml:"task,attr,omitempty"`
	ReportID         string             `xml:"report,attr,omitempty"`
	Severity         string             `xml:"severity,omitempty"`
	Host             string             `xml:"host,omitempty"`
	Location         string             `xml:"location,omitempty"`
	SolutionType     string             `xml:"solution_type,omitempty"`
	Status           string             `xml:"status"`
	OpenTime         string             `xml:"open_time,omitempty"`
	OpenNote         string             `xml:"open_note,omitempty"`
	Result           *TicketResult      `xml:"result,omitempty"`
}

type TicketOwner struct {
	Name string `xml:"name"`
}

type TicketPermissions struct {
	Permission []TicketPermission `xml:"permission"`
}

type TicketPermission struct {
	Name string `xml:"name"`
}

type TicketAssignedTo struct {
	User *TicketUser `xml:"user"`
}

type TicketUser struct {
	ID string `xml:"id,attr"`
}

type TicketResult struct {
	ID string `xml:"id,attr"`
}
