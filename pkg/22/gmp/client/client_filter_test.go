package client

import (
	"testing"

	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

func TestCreateFilter(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateFilterCommand{
		Name:    "Single Targets",
		Comment: "Targets with only one host",
		Term:    "ips=1 first=1 rows=-2",
		Type:    "target",
	}
	resp, err := cli.CreateFilter(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateFilter: %s", err)
	}
	if resp.Status != "201" {
		t.Errorf("Expected status 201, got %s", resp.Status)
	}
	if resp.StatusText != "OK, resource created" {
		t.Errorf("Expected status_text 'OK, resource created', got %s", resp.StatusText)
	}
	if resp.ID != "254cd3ef-bbe1-4d58-859d-21b8d0c046c7" {
		t.Errorf("Expected filter ID '254cd3ef-bbe1-4d58-859d-21b8d0c046c7', got %s", resp.ID)
	}
}
