package client

import (
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestGetPortLists(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetPortLists{}
	cmd.PortListID = "33d0cd82-57c6-11e1-8ed1-406186ea4fc5"
	resp, err := cli.GetPortLists(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetPortLists: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestCreatePortList(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.CreatePortList{
		Name:      "All TCP",
		Comment:   "All possible TCP ports",
		PortRange: "T:1-65535",
	}
	resp, err := cli.CreatePortList(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreatePortList: %s", err)
	}
	if resp.Status != "201" {
		t.Fatalf("Unexpected status. Expected: 201 Got: %s", resp.Status)
	}
	if resp.ID != "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
		t.Fatalf("Unexpected ID. Expected: 254cd3ef-bbe1-4d58-859d-21b8d0c046c6 Got: %s", resp.ID)
	}
}

func TestModifyPortList(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.ModifyPortList{
		PortListID: "27140836-05ae-4e8b-9abf-f725ddc2888f",
		Name:       "PL-WS1",
		Comment:    "Port List for Web Server 1",
	}
	resp, err := cli.ModifyPortList(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyPortList: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

func TestDeletePortList(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DeletePortList{
		PortListID: "267a3405-e84a-47da-97b2-5fa0d2e8995e",
		Ultimate:   true,
	}
	resp, err := cli.DeletePortList(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeletePortList: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

func TestCreatePortRange(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.CreatePortRange{
		PortListID: "354cd3ef-bbe1-4d58-859d-21b8d0c046c4",
		Start:      "777",
		End:        "779",
		Type:       "TCP",
	}
	resp, err := cli.CreatePortRange(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreatePortRange: %s", err)
	}
	if resp.Status != "201" {
		t.Fatalf("Unexpected status. Expected: 201 Got: %s", resp.Status)
	}
	if resp.ID != "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
		t.Fatalf("Unexpected ID. Expected: 254cd3ef-bbe1-4d58-859d-21b8d0c046c6 Got: %s", resp.ID)
	}
}

func TestDeletePortRange(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DeletePortRange{
		PortRangeID: "267a3405-e84a-47da-97b2-5fa0d2e8995e",
	}
	resp, err := cli.DeletePortRange(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeletePortRange: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}
