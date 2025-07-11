package client

import (
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestCreateNote(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.CreateNote{
		Text:   "This issue should be resolved after the upgrade.",
		NVT:    &commands.NoteNVT{OID: "1.3.6.1.4.1.25623.1.0.10330"},
		Result: &commands.CreateNoteResult{ID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"},
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

	cmd := &commands.ModifyNote{
		NoteID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Text:   "This issue should be resolved after the upgrade.",
		Result: &commands.ModifyNoteResult{ID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"},
	}
	resp, err := cli.ModifyNote(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyNote: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
}

func TestGetNotes(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetNotes{}
	resp, err := cli.GetNotes(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetNotes: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if len(resp.Notes) != 1 || resp.Notes[0].Text != "This is the full text of the note." {
		t.Errorf("Expected note text 'This is the full text of the note.', got %+v", resp.Notes)
	}
}

func TestDeleteNote(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DeleteNote{
		NoteID:   "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Ultimate: "0",
	}
	resp, err := cli.DeleteNote(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteNote: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status_text 'OK', got %s", resp.StatusText)
	}
}
