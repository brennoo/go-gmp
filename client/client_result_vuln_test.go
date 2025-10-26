package client

import (
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestGetResults(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetResults{TaskID: "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31"}
	resp, err := cli.GetResults(cmd)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Expected status 200, got %s", resp.Status)
	}
}

func TestGetVulns(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetVulns{VulnID: "1.3.6.1.4.1.25623.1.0.808160"}
	resp, err := cli.GetVulns(cmd)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Expected status 200, got %s", resp.Status)
	}
}
