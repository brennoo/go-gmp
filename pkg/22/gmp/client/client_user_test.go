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

func TestModifyUser(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.ModifyUserCommand{
		UserID:  "user-uuid",
		Name:    "testuser",
		NewName: "newuser",
		Comment: "Updated user",
		Password: &gmp.ModifyUserPassword{
			Modify: "1",
			Text:   "newpass",
		},
		Roles: []gmp.CreateUserRole{{ID: "role-uuid"}},
	}
	resp, err := cli.ModifyUser(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyUser: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}

	// Failure case
	cmdFail := &gmp.ModifyUserCommand{
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
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.GetUsersCommand{}
	resp, err := cli.GetUsers(cmd)
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
	cmdFail := &gmp.GetUsersCommand{}
	respFail, err := cli.GetUsers(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during GetUsers (fail): %s", err)
	}
	// In this mock, failure is not simulated by input, so just check status
	if respFail.Status != "200" {
		t.Errorf("Expected status 200, got %s", respFail.Status)
	}
}

func TestDeleteUser(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.DeleteUserCommand{
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
	cmdFail := &gmp.DeleteUserCommand{
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
