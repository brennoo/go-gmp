package client

import (
	"testing"

	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

func TestModifyAgents(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.ModifyAgentsCommand{
		Agents: []gmp.Agent{
			{ID: "fb6451bf-ec5a-45a8-8bab-5cf4b862e51b"},
		},
		Authorized:        "1",
		MinInterval:       "1000",
		HeartbeatInterval: "0",
		Schedule:          "@every 12h",
		Comment:           "example update",
	}
	resp, err := cli.ModifyAgents(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyAgents: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
}
