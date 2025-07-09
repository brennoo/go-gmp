package client

import (
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestCreateFilter(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.CreateFilter{
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

func TestModifyFilter(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.ModifyFilter{
		FilterID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c7",
		Term:     "name=local",
	}
	resp, err := cli.ModifyFilter(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyFilter: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status_text 'OK', got %s", resp.StatusText)
	}
}

func TestGetFilters(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetFilters{}
	resp, err := cli.GetFilters(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetFilters: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status_text 'OK', got %s", resp.StatusText)
	}
	if len(resp.Filters) != 1 {
		t.Fatalf("Expected 1 filter, got %d", len(resp.Filters))
	}
	f := resp.Filters[0]
	if f.ID != "b493b7a8-7489-11df-a3ec-001164764cea" {
		t.Errorf("Expected filter ID 'b493b7a8-7489-11df-a3ec-001164764cea', got %s", f.ID)
	}
	if f.Name != "Single Targets" {
		t.Errorf("Expected filter name 'Single Targets', got %s", f.Name)
	}
	if f.Comment != "Targets with only one host" {
		t.Errorf("Expected filter comment 'Targets with only one host', got %s", f.Comment)
	}
	if f.Term != "ips=1 first=1 rows=-2" {
		t.Errorf("Expected filter term 'ips=1 first=1 rows=-2', got %s", f.Term)
	}
	if f.Type != "target" {
		t.Errorf("Expected filter type 'target', got %s", f.Type)
	}
}

func TestDeleteFilter(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DeleteFilter{
		FilterID: "b493b7a8-7489-11df-a3ec-001164764cea",
		Ultimate: "0",
	}
	resp, err := cli.DeleteFilter(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteFilter: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status_text 'OK', got %s", resp.StatusText)
	}
}
