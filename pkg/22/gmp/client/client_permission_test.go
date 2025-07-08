package client

import (
	"testing"

	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

func TestCreatePermission(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.CreatePermissionCommand{
		Name: "get_targets",
		Subject: &gmp.CreatePermissionSubject{
			ID:   "66abe5ce-c011-11e3-b96e-406186ea4fc5",
			Type: "user",
		},
		Resource: &gmp.CreatePermissionResource{
			ID: "b493b7a8-7489-11df-a3ec-002264764cea",
		},
	}
	resp, err := cli.CreatePermission(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreatePermission: %s", err)
	}
	if resp.Status != "201" {
		t.Errorf("Expected status 201, got %s", resp.Status)
	}
	if resp.ID != "created-permission-id" {
		t.Errorf("Expected ID 'created-permission-id', got %s", resp.ID)
	}

	// Failure case
	cmdFail := &gmp.CreatePermissionCommand{
		Name:    "",
		Subject: &gmp.CreatePermissionSubject{},
	}
	respFail, err := cli.CreatePermission(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during CreatePermission (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
}
