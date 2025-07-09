package commands

import (
	"testing"
)

func TestNewAuthenticateCredentials(t *testing.T) {
	username := "testuser"
	password := "testpass"
	cmd := NewAuthenticateCredentials(username, password)

	if cmd.Credentials == nil {
		t.Fatal("Credentials should not be nil")
	}
	if cmd.Credentials.Username != username {
		t.Errorf("expected username %q, got %q", username, cmd.Credentials.Username)
	}
	if cmd.Credentials.Password != password {
		t.Errorf("expected password %q, got %q", password, cmd.Credentials.Password)
	}
}
