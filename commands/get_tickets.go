package commands

import "encoding/xml"

// GetTickets represents a get_tickets command request.
type GetTickets struct {
	XMLName xml.Name `xml:"get_tickets"`
}

// GetTicketsResponse represents a get_tickets command response.
type GetTicketsResponse struct {
	XMLName    xml.Name `xml:"get_tickets_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	Tickets    []Ticket `xml:"ticket"`
}

// Ticket represents a <ticket> element in the get_tickets response.
type Ticket struct {
	ID               string             `xml:"id,attr"`
	Owner            *TicketOwner       `xml:"owner,omitempty"`
	Name             string             `xml:"name"`
	Comment          string             `xml:"comment"`
	CreationTime     string             `xml:"creation_time"`     // ISO 8601 string
	ModificationTime string             `xml:"modification_time"` // ISO 8601 string
	Writable         string             `xml:"writable"`          // "0"/"1"
	InUse            string             `xml:"in_use"`            // "0"/"1"
	Permissions      *TicketPermissions `xml:"permissions,omitempty"`
	AssignedTo       *TicketAssignedTo  `xml:"assigned_to,omitempty"`
	TaskID           string             `xml:"task,attr,omitempty"`
	ReportID         string             `xml:"report,attr,omitempty"`
	Severity         string             `xml:"severity,omitempty"`
	Host             string             `xml:"host,omitempty"`
	Location         string             `xml:"location,omitempty"`
	SolutionType     string             `xml:"solution_type,omitempty"`
	Status           string             `xml:"status"`
	OpenTime         string             `xml:"open_time,omitempty"` // ISO 8601 string
	OpenNote         string             `xml:"open_note,omitempty"`
	Result           *TicketResult      `xml:"result,omitempty"`
}

// TicketOwner represents the <owner> element for a ticket.
type TicketOwner struct {
	Name string `xml:"name"`
}

// TicketPermissions represents the <permissions> element for a ticket.
type TicketPermissions struct {
	Permission []TicketPermission `xml:"permission"`
}

// TicketPermission represents a <permission> element for a ticket.
type TicketPermission struct {
	Name string `xml:"name"`
}

// TicketAssignedTo represents the <assigned_to> element for a ticket.
type TicketAssignedTo struct {
	User *TicketUser `xml:"user"`
}

// TicketUser represents the <user> element in assigned_to.
type TicketUser struct {
	ID string `xml:"id,attr"`
}

// TicketResult represents the <result> element for a ticket.
type TicketResult struct {
	ID string `xml:"id,attr"`
}
