package client

import (
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestAuthenticate(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.Authenticate{}
	cmd.Credentials.Username = "openvas"
	cmd.Credentials.Password = "123"
	resp, err := cli.Authenticate(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during Authenticate: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}
