package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// CreateNote.Active/Severity are *int/*float64: unset (this case), a
	// nonzero value, and an explicit zero all round-trip distinctly.
	RegisterRequestCase("create_note", "with_nvt_and_result", &commands.CreateNote{
		Text:   "False positive: this service is intentionally exposed for testing.",
		NVT:    commands.NoteNVT{OID: "1.3.6.1.4.1.25623.1.0.100315"},
		Result: &commands.NoteResult{ID: "9b7f6e5d-4c3b-2a1f-0e9d-8c7b6a5d4e3f"},
		Hosts:  "192.0.2.10",
		Port:   "443/tcp",
	})

	RegisterRequestCase("create_note", "active_and_severity_set", &commands.CreateNote{
		Text:     "Active for 7 days at severity 7.5.",
		NVT:      commands.NoteNVT{OID: "1.3.6.1.4.1.25623.1.0.100315"},
		Active:   commands.Int(7),
		Severity: commands.Float64(7.5),
	})

	RegisterRequestCase("create_note", "active_and_severity_zero", &commands.CreateNote{
		Text:     "Active off, severity explicitly zero.",
		NVT:      commands.NoteNVT{OID: "1.3.6.1.4.1.25623.1.0.100315"},
		Active:   commands.Int(0),
		Severity: commands.Float64(0),
	})

	// ModifyNote.Active/Severity are *int/*float64: unset, a nonzero value,
	// and an explicit zero all round-trip distinctly.
	RegisterRequestCase("modify_note", "unset", &commands.ModifyNote{
		NoteID: "9b7f6e5d-4c3b-2a1f-0e9d-8c7b6a5d4e3f",
		Text:   "Updated note text.",
	})

	RegisterRequestCase("modify_note", "active_and_severity_set", &commands.ModifyNote{
		NoteID:   "9b7f6e5d-4c3b-2a1f-0e9d-8c7b6a5d4e3f",
		Active:   commands.Int(-1),
		Severity: commands.Float64(4.2),
	})

	RegisterRequestCase("modify_note", "active_and_severity_zero", &commands.ModifyNote{
		NoteID:   "9b7f6e5d-4c3b-2a1f-0e9d-8c7b6a5d4e3f",
		Active:   commands.Int(0),
		Severity: commands.Float64(0),
	})
}
