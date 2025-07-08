package client

import (
	"testing"

	gmp "github.com/brennoo/go-gmp/pkg"
)

func TestCreateScanner(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateScannerCommand{
		Name:  "Default Scanner",
		Host:  "localhost",
		Port:  "9391",
		Type:  "2",
		CAPub: "dummy-ca-pub",
		Credential: &gmp.CreateScannerCredential{
			ID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		},
	}
	resp, err := cli.CreateScanner(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateScanner: %s", err)
	}
	if resp.Status != "201" {
		t.Fatalf("Unexpected status. Expected: 201 Got: %s", resp.Status)
	}
	if resp.ID != "814cd30f-dee1-4d58-851d-21b8d0c048e3" {
		t.Fatalf("Unexpected ID. Expected: 814cd30f-dee1-4d58-851d-21b8d0c048e3 Got: %s", resp.ID)
	}
}

func TestGetScanners(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetScannersCommand{}
	cmd.ScannerID = "ee0311e7-3247-4425-bb9c-866d59f1e0e9"
	resp, err := cli.GetScanners(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetScanners: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestCreateScannerFail(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateScannerCommand{
		Name:  "",
		Host:  "localhost",
		Port:  "9391",
		Type:  "2",
		CAPub: "dummy-ca-pub",
		Credential: &gmp.CreateScannerCredential{
			ID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		},
	}
	resp, err := cli.CreateScanner(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateScanner (fail): %s", err)
	}
	if resp.Status != "400" {
		t.Fatalf("Unexpected status. Expected: 400 Got: %s", resp.Status)
	}
}

func TestModifyScanner(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.ModifyScannerCommand{
		ScannerID: "scanner-uuid",
		Name:      "Updated Scanner",
	}
	resp, err := cli.ModifyScanner(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyScanner: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
}

func TestModifyScannerFail(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.ModifyScannerCommand{
		ScannerID: "",
		Name:      "",
	}
	resp, err := cli.ModifyScanner(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyScanner (fail): %s", err)
	}
	if resp.Status != "400" {
		t.Fatalf("Unexpected status. Expected: 400 Got: %s", resp.Status)
	}
	if resp.StatusText != "Bad request" {
		t.Fatalf("Unexpected status text. Expected: Bad request Got: %s", resp.StatusText)
	}
}

func TestDeleteScanner(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.DeleteScannerCommand{ScannerID: "scanner-uuid", Ultimate: "1"}
	resp, err := cli.DeleteScanner(cmd)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Expected status text 'OK', got '%s'", resp.StatusText)
	}
}

func TestVerifyScanner(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.VerifyScannerCommand{ScannerID: "scanner-uuid"}
	resp, err := cli.VerifyScanner(cmd)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Expected status text 'OK', got '%s'", resp.StatusText)
	}
	if resp.Version != "OTP/2.0" {
		t.Fatalf("Expected version 'OTP/2.0', got '%s'", resp.Version)
	}
}
