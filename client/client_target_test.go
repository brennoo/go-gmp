package client

import (
	"testing"

	"github.com/brennoo/go-gmp"
)

func TestCreateTarget(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateTargetCommand{}
	cmd.Name = "New Target"
	resp, err := cli.CreateTarget(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateTarget: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

func TestModifyTarget(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.ModifyTargetCommand{}
	cmd.TargetID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
	cmd.Name = "Modified Target"
	resp, err := cli.ModifyTarget(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyTarget: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

func TestGetTargets(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetTargetsCommand{}
	cmd.TargetID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
	resp, err := cli.GetTargets(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetTargets: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

func TestDeleteTarget(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.DeleteTargetCommand{}
	cmd.TargetID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
	resp, err := cli.DeleteTarget(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteTarget: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}
