package client

import (
	"testing"

	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

func TestCreateTLSCertificate(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateTLSCertificateCommand{
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
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.ModifyTLSCertificateCommand{
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
