package client

import (
	"context"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestCreateGroup(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.CreateGroup{
		Name:  "Managers",
		Users: "sarah, bob",
	}
	resp, err := cli.CreateGroup(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateGroup: %s", err)
	}
	if resp.Status != "201" {
		t.Errorf("Expected status 201, got %s", resp.Status)
	}
	if resp.ID != "d94211b6-ba40-11e3-bcb1-406186ea4fc5" {
		t.Errorf("Expected group ID 'd94211b6-ba40-11e3-bcb1-406186ea4fc5', got %s", resp.ID)
	}
}

func TestModifyGroup(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.ModifyGroup{
		GroupID: "d94211b6-ba40-11e3-bcb1-406186ea4fc5",
		Name:    "Line Managers",
	}
	resp, err := cli.ModifyGroup(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyGroup: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
}

func TestGetGroups(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetGroups{}
	resp, err := cli.GetGroups(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetGroups: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if len(resp.Group) != 1 || resp.Group[0].Name != "Management" {
		t.Errorf("Expected group 'Management', got %+v", resp.Group)
	}
}

func TestDeleteGroup(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DeleteGroup{
		GroupID:  "d94211b6-ba40-11e3-bcb1-406186ea4fc5",
		Ultimate: true,
	}
	resp, err := cli.DeleteGroup(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteGroup: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
}
