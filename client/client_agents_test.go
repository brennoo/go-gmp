package client

import (
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestModifyAgents(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.ModifyAgents{
		Agents: []commands.Agent{
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

func TestGetAgents(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetAgents{}
	resp, err := cli.GetAgents(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetAgents: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if len(resp.Agents) != 1 {
		t.Fatalf("Expected 1 agent, got %d", len(resp.Agents))
	}
	agent := resp.Agents[0]
	if agent.ID != "62462fe0-5834-4630-afc2-0d040c63487c" {
		t.Fatalf("Unexpected agent ID. Expected: 62462fe0-5834-4630-afc2-0d040c63487c Got: %s", agent.ID)
	}
	if agent.AgentID != "GAT-29-p0MPX0FT" {
		t.Fatalf("Unexpected agent_id. Expected: GAT-29-p0MPX0FT Got: %s", agent.AgentID)
	}
}

func TestDeleteAgents(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DeleteAgents{
		Agents: []commands.Agent{
			{ID: "c6f1cdc3-8c2c-4b2e-9f43-139d23c7cfd4"},
		},
	}
	resp, err := cli.DeleteAgents(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteAgents: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
}
