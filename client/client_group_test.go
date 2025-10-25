package client

import (
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestCreateGroup(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.CreateGroup{
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

	cmd := &commands.ModifyGroup{
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

func TestGetGroups(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetGroups{}
	resp, err := cli.GetGroups(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetGroups: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if len(resp.Groups) != 1 || resp.Groups[0].Name != "Management" {
		t.Errorf("Expected group 'Management', got %+v", resp.Groups)
	}
}

func TestDeleteGroup(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DeleteGroup{
		GroupID:  "d94211b6-ba40-11e3-bcb1-406186ea4fc5",
		Ultimate: true,
	}
	resp, err := cli.DeleteGroup(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteGroup: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
}
