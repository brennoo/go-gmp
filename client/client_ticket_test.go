package client

import (
	"context"
	"testing"

	"github.com/brennoo/go-gmp/commands"
	"github.com/brennoo/go-gmp/commands/filtering"
)

func TestCreateTicket(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.CreateTicket{
		Result: commands.TicketResult{ID: "138c1216-4acb-4ded-bef3-7fab80eac8c7"},
		Assigned: commands.TicketAssignedTo{
			User: &commands.TicketUser{ID: "33e92d3e-a379-4c46-a4cf-88c8201ab710"},
		},
		OpenNote: "Please fix today.",
	}
	resp, err := cli.CreateTicket(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateTicket: %s", err)
	}
	if resp.Status != "201" {
		t.Errorf("Expected status 201, got %s", resp.Status)
	}
	if resp.StatusText != "OK, resource created" {
		t.Errorf("Expected status_text 'OK, resource created', got %s", resp.StatusText)
	}
	if resp.ID != "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
		t.Errorf("Expected ticket ID '254cd3ef-bbe1-4d58-859d-21b8d0c046c6', got %s", resp.ID)
	}
}

func TestGetTickets(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()
	resp, err := cli.GetTickets(ctx)
	if err != nil {
		t.Fatalf("Unexpected error during GetTickets: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status_text 'OK', got %s", resp.StatusText)
	}
	if len(resp.Tickets) != 1 {
		t.Fatalf("Expected 1 ticket, got %d", len(resp.Tickets))
	}
	ticket := resp.Tickets[0]
	if ticket.ID != "93cd2f71-48c3-4cf2-b542-5b256f59cae0" {
		t.Errorf("Expected ticket ID '93cd2f71-48c3-4cf2-b542-5b256f59cae0', got %s", ticket.ID)
	}
	if ticket.Name != "OpenSSH Denial of Service Vulnerability - Jan16" {
		t.Errorf("Expected ticket name 'OpenSSH Denial of Service Vulnerability - Jan16', got %s", ticket.Name)
	}
	if ticket.Status != "Open" {
		t.Errorf("Expected ticket status 'Open', got %s", ticket.Status)
	}
	if ticket.OpenNote != "Probably the new version fixes this" {
		t.Errorf("Expected open_note 'Probably the new version fixes this', got %s", ticket.OpenNote)
	}
}

func TestModifyTicket(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.ModifyTicket{
		TicketID:   "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Status:     "Closed",
		ClosedNote: "Resolved with January update",
	}
	resp, err := cli.ModifyTicket(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyTicket: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status_text 'OK', got %s", resp.StatusText)
	}
}

func TestDeleteTicket(t *testing.T) {
	cli := New(MockedConnection())
	cmd := &commands.DeleteTicket{
		TicketID: "ticket-uuid",
		Ultimate: true,
	}
	resp, err := cli.DeleteTicket(cmd)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "200" {
		t.Errorf("expected status 200, got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("expected status text OK, got: %s", resp.StatusText)
	}
}

// TestGetTicketsConsolidated tests the GetTickets method.
func TestGetTicketsConsolidated(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()

	// Test with no filters
	resp, err := cli.GetTickets(ctx)
	if err != nil {
		t.Fatalf("Unexpected error during GetTickets: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with string filters
	resp, err = cli.GetTickets(ctx, "status=Open")
	if err != nil {
		t.Fatalf("Unexpected error during GetTickets with string filter: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with functional options
	resp, err = cli.GetTickets(ctx, filtering.WithStatus("Open"))
	if err != nil {
		t.Fatalf("Unexpected error during GetTickets with functional options: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}
