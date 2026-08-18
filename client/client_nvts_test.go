package client

import (
	"context"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestGetNVTsWithDetails(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetNVTs{
		Details: commands.Bool(true),
	}
	resp, err := cli.GetNVTs(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetNVTs: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if len(resp.NVT) == 0 {
		t.Fatalf("Expected at least one NVT, got %d", len(resp.NVT))
	}
	if resp.NVT[0].OID != "1.3.6.1.4.1.25623.1.7.13005" {
		t.Fatalf("Unexpected first NVT OID. Expected: 1.3.6.1.4.1.25623.1.7.13005 Got: %s", resp.NVT[0].OID)
	}
	if resp.NVT[0].Name != "Services" {
		t.Fatalf("Unexpected first NVT name. Expected: Services Got: %s", resp.NVT[0].Name)
	}
}

func TestGetNVTsSingle(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetNVTs{
		NVTOID: "1.3.6.1.4.1.25623.1.0.10330",
	}
	resp, err := cli.GetNVTs(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetNVTs: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if len(resp.NVT) != 1 {
		t.Fatalf("Expected exactly one NVT, got %d", len(resp.NVT))
	}
	if resp.NVT[0].OID != "1.3.6.1.4.1.25623.1.0.10330" {
		t.Fatalf("Unexpected NVT OID. Expected: 1.3.6.1.4.1.25623.1.0.10330 Got: %s", resp.NVT[0].OID)
	}
	if resp.NVT[0].Name != "Services" {
		t.Fatalf("Unexpected NVT name. Expected: Services Got: %s", resp.NVT[0].Name)
	}
}

func TestGetNVTFamilies(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetNVTFamilies{}
	resp, err := cli.GetNVTFamilies(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetNVTFamilies: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if len(resp.Families.Family) == 0 {
		t.Fatalf("Expected at least one family, got %d", len(resp.Families.Family))
	}
	if resp.Families.Family[0].Name != "Credentials" {
		t.Fatalf("Unexpected first family name. Expected: Credentials Got: %s", resp.Families.Family[0].Name)
	}
	if resp.Families.Family[0].MaxNVTCount != 8 {
		t.Fatalf("Unexpected first family NVT count. Expected: 8 Got: %d", resp.Families.Family[0].MaxNVTCount)
	}
	if resp.Families.Family[1].Name != "Port scanners" {
		t.Fatalf("Unexpected second family name. Expected: Port scanners Got: %s", resp.Families.Family[1].Name)
	}
	if resp.Families.Family[1].MaxNVTCount != 12 {
		t.Fatalf("Unexpected second family NVT count. Expected: 12 Got: %d", resp.Families.Family[1].MaxNVTCount)
	}
}
