package client

import (
	"testing"

	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

func TestCreateRole(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.CreateRoleCommand{
		Name:    "Test Role",
		Comment: "A test role",
	}
	resp, err := cli.CreateRole(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateRole: %s", err)
	}
	if resp.Status != "201" {
		t.Errorf("Expected status 201, got %s", resp.Status)
	}
	if resp.ID != "created-role-id" {
		t.Errorf("Expected ID 'created-role-id', got %s", resp.ID)
	}

	// Failure case
	cmdFail := &gmp.CreateRoleCommand{
		Name: "",
	}
	respFail, err := cli.CreateRole(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during CreateRole (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
}
