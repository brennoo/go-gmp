package client

import (
	"testing"

	"github.com/brennoo/go-gmp"
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

func TestModifyPermission(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.ModifyPermissionCommand{
		PermissionID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Subject: &gmp.CreatePermissionSubject{
			ID:   "76e47468-c095-11e3-9285-406186ea4fc5",
			Type: "user",
		},
	}
	resp, err := cli.ModifyPermission(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyPermission: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got '%s'", resp.StatusText)
	}

	// Failure case
	cmdFail := &gmp.ModifyPermissionCommand{
		PermissionID: "",
		Subject:      &gmp.CreatePermissionSubject{},
	}
	respFail, err := cli.ModifyPermission(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyPermission (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
}

func TestGetPermissions(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.GetPermissionsCommand{}
	resp, err := cli.GetPermissions(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetPermissions: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if len(resp.Permissions) != 1 || resp.Permissions[0].Name != "Management" {
		t.Errorf("Expected permission 'Management', got %+v", resp.Permissions)
	}
}

func TestDeletePermission(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.DeletePermissionCommand{
		PermissionID: "267a3405-e84a-47da-97b2-5fa0d2e8995e",
		Ultimate:     "1",
	}
	resp, err := cli.DeletePermission(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeletePermission: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
}
