package client

import (
	"testing"

	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

func TestCreateNote(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateNoteCommand{
		Text:   "This issue should be resolved after the upgrade.",
		NVT:    &gmp.CreateNoteNVT{OID: "1.3.6.1.4.1.25623.1.0.10330"},
		Result: &gmp.CreateNoteResult{ID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"},
	}
	resp, err := cli.CreateNote(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateNote: %s", err)
	}
	if resp.Status != "202" {
		t.Errorf("Expected status 202, got %s", resp.Status)
	}
	if resp.ID != "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
		t.Errorf("Expected note ID '254cd3ef-bbe1-4d58-859d-21b8d0c046c6', got %s", resp.ID)
	}
}

func TestModifyNote(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.ModifyNoteCommand{
		NoteID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Text:   "This issue should be resolved after the upgrade.",
		Result: &gmp.ModifyNoteResult{ID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"},
	}
	resp, err := cli.ModifyNote(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyNote: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
}
