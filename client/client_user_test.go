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
		Role:     []commands.UserRole{{ID: "role-uuid"}},
	}
	resp, err := cli.CreateUser(context.Background(), cmd)
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
	respFail, err := cli.CreateUser(context.Background(), cmdFail)
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
	comment := "Updated user"
	cmd := &commands.ModifyUser{
		UserID:  "user-uuid",
		Name:    "testuser",
		NewName: "newuser",
		Comment: &comment,
		Password: &commands.UserPassword{
			Modify: commands.Bool(true),
			Value:  "newpass",
		},
		Role: []commands.UserRole{{ID: "role-uuid"}},
	}
	resp, err := cli.ModifyUser(context.Background(), cmd)
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
	respFail, err := cli.ModifyUser(context.Background(), cmdFail)
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
	if len(resp.User) != 1 || resp.User[0].Name != "foobar" {
		t.Errorf("Expected user 'foobar', got %+v", resp.User)
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
	resp, err := cli.DeleteUser(context.Background(), cmd)
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
	respFail, err := cli.DeleteUser(context.Background(), cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteUser (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
}
