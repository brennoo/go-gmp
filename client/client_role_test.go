package client

import (
	"context"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestCreateRole(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &commands.CreateRole{
		Name:    "Test Role",
		Comment: "A test role",
	}
	resp, err := cli.CreateRole(context.Background(), cmd)
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
	cmdFail := &commands.CreateRole{
		Name: "",
	}
	respFail, err := cli.CreateRole(context.Background(), cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during CreateRole (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
}

func TestModifyRole(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &commands.ModifyRole{
		RoleID:  "role-uuid",
		Name:    "Updated Role",
		Comment: "Updated comment",
		Users:   "user1,user2",
	}
	resp, err := cli.ModifyRole(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyRole: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got '%s'", resp.StatusText)
	}

	// Failure case
	cmdFail := &commands.ModifyRole{
		RoleID: "",
		Name:   "",
	}
	respFail, err := cli.ModifyRole(context.Background(), cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyRole (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
	if respFail.StatusText != "Bad request" {
		t.Errorf("Expected status text 'Bad request', got '%s'", respFail.StatusText)
	}
}

func TestGetRolesRaw(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetRoles{}
	resp, err := cli.GetRolesRaw(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetRolesRaw: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if len(resp.Role) != 1 || resp.Role[0].Name != "Management" {
		t.Errorf("Expected role 'Management', got %+v", resp.Role)
	}
}

func TestGetRolesVariadic(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	resp, err := cli.GetRoles(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error during GetRoles: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if len(resp.Role) != 1 || resp.Role[0].Name != "Management" {
		t.Errorf("Expected role 'Management', got %+v", resp.Role)
	}
}

func TestDeleteRole(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &commands.DeleteRole{
		RoleID:   "b64c81b2-b9de-11e3-a2e9-406186ea4fc5",
		Ultimate: true,
	}
	resp, err := cli.DeleteRole(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteRole: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text OK, got %s", resp.StatusText)
	}

	// Failure case
	cmdFail := &commands.DeleteRole{
		RoleID:   "",
		Ultimate: true,
	}
	respFail, err := cli.DeleteRole(context.Background(), cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteRole (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
}
