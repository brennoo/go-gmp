package client

import (
	"testing"

	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

func TestCreateUser(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.CreateUserCommand{
		Name:     "testuser",
		Password: "testpass",
		Roles:    []gmp.CreateUserRole{{ID: "role-uuid"}},
	}
	resp, err := cli.CreateUser(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateUser: %s", err)
	}
	if resp.Status != "201" {
		t.Errorf("Expected status 201, got %s", resp.Status)
	}
	if resp.ID != "created-user-id" {
		t.Errorf("Expected ID 'created-user-id', got %s", resp.ID)
	}

	// Failure case
	cmdFail := &gmp.CreateUserCommand{
		Name:     "",
		Password: "",
	}
	respFail, err := cli.CreateUser(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during CreateUser (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
}
