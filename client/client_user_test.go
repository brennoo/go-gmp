package client

import (
	"context"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestCreateUser(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &commands.CreateUser{
		Name:     "testuser",
		Password: "testpass",
		Roles:    []*commands.UserRole{{ID: "role-uuid"}},
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
	cmdFail := &commands.CreateUser{
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

func TestModifyUser(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &commands.ModifyUser{
		UserID:  "user-uuid",
		Name:    "testuser",
		NewName: "newuser",
		Comment: "Updated user",
		Password: &commands.UserPassword{
			Modify: "1",
			Text:   "newpass",
		},
		Roles: []*commands.UserRole{{ID: "role-uuid"}},
	}
	resp, err := cli.ModifyUser(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyUser: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}

	// Failure case
	cmdFail := &commands.ModifyUser{
		UserID: "",
		Name:   "",
	}
	respFail, err := cli.ModifyUser(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyUser (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
}

func TestGetUsers(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	ctx := context.Background()
	resp, err := cli.GetUsers(ctx, "name=test-user")
	if err != nil {
		t.Fatalf("Unexpected error during GetUsers: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if len(resp.Users) != 1 || resp.Users[0].Name != "foobar" {
		t.Errorf("Expected user 'foobar', got %+v", resp.Users)
	}

	// Failure case (simulate by using a special test flag)
	respFail, err := cli.GetUsers(ctx, "invalid-filter")
	if err != nil {
		t.Fatalf("Unexpected error during GetUsers (fail): %s", err)
	}
	// In this mock, failure is not simulated by input, so just check status
	if respFail.Status != "200" {
		t.Errorf("Expected status 200, got %s", respFail.Status)
	}
}

func TestDeleteUser(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &commands.DeleteUser{
		Name: "foobar",
	}
	resp, err := cli.DeleteUser(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteUser: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}

	// Failure case
	cmdFail := &commands.DeleteUser{
		Name: "",
	}
	respFail, err := cli.DeleteUser(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteUser (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
}
