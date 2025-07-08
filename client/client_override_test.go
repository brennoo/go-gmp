package client

import (
	"testing"

	"github.com/brennoo/go-gmp"
)

func TestCreateOverride(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateOverrideCommand{
		Text:      "This is actually of little concern.",
		NVT:       gmp.CreateOverrideNVT{OID: "1.3.6.1.4.1.25623.1.0.10330"},
		NewThreat: "Low",
		Result:    "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
	}
	resp, err := cli.CreateOverride(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateOverride: %s", err)
	}
	if resp.Status != "201" {
		t.Errorf("Expected status 201, got %s", resp.Status)
	}
	if resp.StatusText != "OK, resource created" {
		t.Errorf("Expected status text 'OK, resource created', got '%s'", resp.StatusText)
	}
	if resp.ID != "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
		t.Errorf("Expected ID '254cd3ef-bbe1-4d58-859d-21b8d0c046c6', got '%s'", resp.ID)
	}
}

func TestDeleteOverride(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.DeleteOverrideCommand{
		OverrideID: "267a3405-e84a-47da-97b2-5fa0d2e8995e",
		Ultimate:   "0",
	}
	resp, err := cli.DeleteOverride(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteOverride: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text OK, got %s", resp.StatusText)
	}
}

func TestGetOverrides(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetOverridesCommand{
		OverrideID: "b76b81a7-9df8-42df-afff-baa9d4620128",
	}
	resp, err := cli.GetOverrides(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetOverrides: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text OK, got %s", resp.StatusText)
	}
	if len(resp.Overrides) != 1 {
		t.Errorf("Expected 1 override, got %d", len(resp.Overrides))
	}
	if resp.Overrides[0].ID != "b76b81a7-9df8-42df-afff-baa9d4620128" {
		t.Errorf("Expected override ID b76b81a7-9df8-42df-afff-baa9d4620128, got %s", resp.Overrides[0].ID)
	}
}

func TestModifyOverride(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.ModifyOverrideCommand{
		OverrideID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Text:       "This issue is less important in our setup.",
		NewThreat:  "Low",
	}
	resp, err := cli.ModifyOverride(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyOverride: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text OK, got %s", resp.StatusText)
	}
}
