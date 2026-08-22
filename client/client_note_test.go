package client

import (
	"context"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestCreateNote(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.CreateNote{
		Text:   "This issue should be resolved after the upgrade.",
		NVT:    commands.NoteNVT{OID: "1.3.6.1.4.1.25623.1.0.10330"},
		Result: &commands.NoteResult{ID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"},
	}
	resp, err := cli.CreateNote(context.Background(), cmd)
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
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.ModifyNote{
		NoteID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Text:   "This issue should be resolved after the upgrade.",
		Result: &commands.NoteResult{ID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"},
	}
	resp, err := cli.ModifyNote(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyNote: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
}

func TestGetNotesRaw(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetNotes{}
	resp, err := cli.GetNotesRaw(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetNotesRaw: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	// <note> is a repeated child of <get_notes_response>, not of a <notes>
	// wrapper: resp.Notes is the schema's <notes start=".." max=".."> paging
	// element. Note.Text carries only @excerpt (the schema DSL cannot declare
	// an element with both attributes and text), so the note body is not
	// assertable here; assert the note's identity and flags instead.
	if len(resp.Note) != 1 {
		t.Fatalf("Expected 1 note, got %d: %+v", len(resp.Note), resp.Note)
	}
	note := resp.Note[0]
	if note.ID != "b76b81a7-9df8-42df-afff-baa9d4620128" {
		t.Errorf("Expected note ID 'b76b81a7-9df8-42df-afff-baa9d4620128', got %s", note.ID)
	}
	if note.NVT.OID != "1.3.6.1.4.1.25623.1.0.75" || note.NVT.Name != "Test NVT: long lines" {
		t.Errorf("Expected NVT 1.3.6.1.4.1.25623.1.0.75/'Test NVT: long lines', got %+v", note.NVT)
	}
	if !bool(note.Writable) || bool(note.InUse) || !bool(note.Active) || !bool(note.Orphan) {
		t.Errorf("Expected writable/active/orphan true and in_use false, got %+v", note)
	}
}

func TestGetNotesVariadic(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	resp, err := cli.GetNotes(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error during GetNotes: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
}

func TestDeleteNote(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DeleteNote{
		NoteID:   "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Ultimate: false,
	}
	resp, err := cli.DeleteNote(context.Background(), cmd)
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
