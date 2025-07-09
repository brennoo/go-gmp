package client

import (
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestGetNvtsWithDetails(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetNvtsCommand{
		Details: "1",
	}
	resp, err := cli.GetNvts(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetNvts: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if len(resp.Nvts) == 0 {
		t.Fatalf("Expected at least one NVT, got %d", len(resp.Nvts))
	}
	if resp.Nvts[0].OID != "1.3.6.1.4.1.25623.1.7.13005" {
		t.Fatalf("Unexpected first NVT OID. Expected: 1.3.6.1.4.1.25623.1.7.13005 Got: %s", resp.Nvts[0].OID)
	}
	if resp.Nvts[0].Name != "Services" {
		t.Fatalf("Unexpected first NVT name. Expected: Services Got: %s", resp.Nvts[0].Name)
	}
}

func TestGetNvtsSingle(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetNvtsCommand{
		NvtOID: "1.3.6.1.4.1.25623.1.0.10330",
	}
	resp, err := cli.GetNvts(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetNvts: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if len(resp.Nvts) != 1 {
		t.Fatalf("Expected exactly one NVT, got %d", len(resp.Nvts))
	}
	if resp.Nvts[0].OID != "1.3.6.1.4.1.25623.1.0.10330" {
		t.Fatalf("Unexpected NVT OID. Expected: 1.3.6.1.4.1.25623.1.0.10330 Got: %s", resp.Nvts[0].OID)
	}
	if resp.Nvts[0].Name != "Services" {
		t.Fatalf("Unexpected NVT name. Expected: Services Got: %s", resp.Nvts[0].Name)
	}
}

func TestGetNvtFamilies(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetNvtFamiliesCommand{}
	resp, err := cli.GetNvtFamilies(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetNvtFamilies: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if len(resp.Families) == 0 {
		t.Fatalf("Expected at least one family, got %d", len(resp.Families))
	}
	if resp.Families[0].Name != "Credentials" {
		t.Fatalf("Unexpected first family name. Expected: Credentials Got: %s", resp.Families[0].Name)
	}
	if resp.Families[0].MaxNvtCount != 8 {
		t.Fatalf("Unexpected first family NVT count. Expected: 8 Got: %d", resp.Families[0].MaxNvtCount)
	}
	if resp.Families[1].Name != "Port scanners" {
		t.Fatalf("Unexpected second family name. Expected: Port scanners Got: %s", resp.Families[1].Name)
	}
	if resp.Families[1].MaxNvtCount != 12 {
		t.Fatalf("Unexpected second family NVT count. Expected: 12 Got: %d", resp.Families[1].MaxNvtCount)
	}
}
