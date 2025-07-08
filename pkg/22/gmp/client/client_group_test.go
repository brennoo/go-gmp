package client

import (
	"testing"

	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

func TestCreateGroup(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateGroupCommand{
		Name:  "Managers",
		Users: "sarah, bob",
	}
	resp, err := cli.CreateGroup(cmd)
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
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.ModifyGroupCommand{
		GroupID: "d94211b6-ba40-11e3-bcb1-406186ea4fc5",
		Name:    "Line Managers",
	}
	resp, err := cli.ModifyGroup(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyGroup: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
}
