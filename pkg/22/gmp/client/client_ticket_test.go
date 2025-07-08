package client

import (
	"testing"

	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

func TestCreateTicket(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateTicketCommand{
		Result: gmp.CreateTicketResult{ID: "138c1216-4acb-4ded-bef3-7fab80eac8c7"},
		Assigned: gmp.CreateTicketAssignedTo{
			User: gmp.CreateTicketUser{ID: "33e92d3e-a379-4c46-a4cf-88c8201ab710"},
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
