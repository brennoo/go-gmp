package client

import (
	"context"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestCreateTLSCertificate(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.CreateTLSCertificate{
		Name:        "Example Certificate",
		Certificate: "MIIDNjCCAp+gAwIBAgIBATANBgkqhkiG9w0BAQQFADCBqTELM[...]",
	}
	resp, err := cli.CreateTLSCertificate(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateTLSCertificate: %s", err)
	}
	if resp.Status != "201" {
		t.Fatalf("Unexpected status. Expected: 201 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK, resource created" {
		t.Fatalf("Unexpected status text. Expected: OK, resource created Got: %s", resp.StatusText)
	}
	if resp.ID != "8a925978-59d0-494b-a837-40b271569727" {
		t.Fatalf("Unexpected ID. Expected: 8a925978-59d0-494b-a837-40b271569727 Got: %s", resp.ID)
	}
}

func TestModifyTLSCertificate(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.ModifyTLSCertificate{
		TLSCertificateID: "8a925978-59d0-494b-a837-40b271569727",
		Name:             "Renamed Example Certificate",
	}
	resp, err := cli.ModifyTLSCertificate(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyTLSCertificate: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
}

func TestGetTLSCertificates(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()
	resp, err := cli.GetTLSCertificates(ctx, "name=test-cert")
	if err != nil {
		t.Fatalf("Unexpected error during GetTLSCertificates: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if len(resp.TLSCertificates) != 1 {
		t.Fatalf("Expected 1 TLS certificate, got %d", len(resp.TLSCertificates))
	}
	cert := resp.TLSCertificates[0]
	if cert.ID != "ba36ed15-92fa-4ae0-af53-bad8ce472f18" {
		t.Fatalf("Unexpected certificate ID. Expected: ba36ed15-92fa-4ae0-af53-bad8ce472f18 Got: %s", cert.ID)
	}
	if cert.Name != "Example Certificate" {
		t.Fatalf("Unexpected certificate name. Expected: Example Certificate Got: %s", cert.Name)
	}
}
