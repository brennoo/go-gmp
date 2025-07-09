package client

import (
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestCreateCredential(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &commands.CreateCredential{
		Name:    "Test Credential",
		KDCs:    commands.NewCredentialKDCs("kdc1.example.com", "kdc2.example.com"),
		Key:     commands.NewCredentialKey("passphrase", "private-key", "public-key"),
		Privacy: commands.NewCredentialPrivacy("AES", "snmp-password"),
	}
	resp, err := cli.CreateCredential(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateCredential: %s", err)
	}
	if resp.Status != "201" {
		t.Errorf("Expected status 201, got %s", resp.Status)
	}
	if resp.ID != "created-credential-id" {
		t.Errorf("Expected ID 'created-credential-id', got %s", resp.ID)
	}

	// Failure case
	cmdFail := &commands.CreateCredential{
		Name:    "",
		KDCs:    commands.NewCredentialKDCs(),
		Key:     commands.NewCredentialKey("", "", ""),
		Privacy: commands.NewCredentialPrivacy("", ""),
	}
	respFail, err := cli.CreateCredential(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during CreateCredential (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
}

func TestModifyCredential(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &commands.ModifyCredential{
		CredentialID: "cred-uuid",
		Name:         "Updated Credential",
	}
	resp, err := cli.ModifyCredential(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyCredential: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got '%s'", resp.StatusText)
	}

	// Failure case
	cmdFail := &commands.ModifyCredential{
		CredentialID: "",
	}
	respFail, err := cli.ModifyCredential(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyCredential (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
	if respFail.StatusText != "Bad request" {
		t.Errorf("Expected status text 'Bad request', got '%s'", respFail.StatusText)
	}
}

func TestGetCredentials(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// List all credentials
	cmd := &commands.GetCredentials{}
	resp, err := cli.GetCredentials(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetCredentials: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if len(resp.Credentials) != 2 {
		t.Errorf("Expected 2 credentials, got %d", len(resp.Credentials))
	}
	if resp.Credentials[0].ID != "cred-uuid-1" || resp.Credentials[1].ID != "cred-uuid-2" {
		t.Errorf("Unexpected credential IDs: %+v", resp.Credentials)
	}

	// Fetch single credential
	cmdSingle := &commands.GetCredentials{CredentialID: "cred-uuid-1"}
	respSingle, err := cli.GetCredentials(cmdSingle)
	if err != nil {
		t.Fatalf("Unexpected error during GetCredentials (single): %s", err)
	}
	if respSingle.Status != "200" {
		t.Errorf("Expected status 200, got %s", respSingle.Status)
	}
	if len(respSingle.Credentials) != 1 || respSingle.Credentials[0].ID != "cred-uuid-1" {
		t.Errorf("Unexpected credential: %+v", respSingle.Credentials)
	}

	// Failure case
	cmdFail := &commands.GetCredentials{CredentialID: "notfound"}
	respFail, err := cli.GetCredentials(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during GetCredentials (fail): %s", err)
	}
	if respFail.Status != "404" {
		t.Errorf("Expected status 404, got %s", respFail.Status)
	}
}

func TestDeleteCredential(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DeleteCredential{CredentialID: "cred-uuid-1", Ultimate: "1"}
	resp, err := cli.DeleteCredential(cmd)
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
