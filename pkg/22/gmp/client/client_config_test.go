package client

import (
	"testing"

	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

func TestGetConfigs(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetConfigsCommand{}
	cmd.ConfigID = "bde773f3-2b3d-4fe6-81cb-6321ae2cc629"
	resp, err := cli.GetConfigs(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetConfigs: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestCreateConfig(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateConfigCommand{}
	cmd.Name = "New Config"
	resp, err := cli.CreateConfig(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateConfig: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestModifyConfig(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.ModifyConfigCommand{}
	cmd.Name = "Modified Config"
	resp, err := cli.ModifyConfig(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyConfig: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestDeleteConfig(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.DeleteConfigCommand{
		ConfigID: "267a3405-e84a-47da-97b2-5fa0d2e8995e",
		Ultimate: "1",
	}
	resp, err := cli.DeleteConfig(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteConfig: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
}

func TestSyncConfig(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.SyncConfigCommand{}
	resp, err := cli.SyncConfig(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during SyncConfig: %s", err)
	}
	if resp.Status != "202" {
		t.Fatalf("Unexpected status. Expected: 202 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK, request submitted" {
		t.Fatalf("Unexpected status text. Expected: OK, request submitted Got: %s", resp.StatusText)
	}
}
