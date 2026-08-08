package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	RegisterRequestCase("create_note", "with_nvt_and_result", &commands.CreateNote{
		Text:   "False positive: this service is intentionally exposed for testing.",
		NVT:    &commands.NoteNVT{OID: "1.3.6.1.4.1.25623.1.0.100315"},
		Result: &commands.CreateNoteResult{ID: "9b7f6e5d-4c3b-2a1f-0e9d-8c7b6a5d4e3f"},
		Hosts:  "192.0.2.10",
		Port:   "443/tcp",
	})
}
