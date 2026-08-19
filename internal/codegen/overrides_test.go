package codegen

import (
	"strings"
	"testing"
	"time"

	"github.com/brennoo/go-gmp/internal/schema"
)

// testProvenance is a valid, complete schema.Provenance for tests that build
// synthetic commands: checkProvenance rejects anything less, and tests here
// are not exercising T-16a itself unless they say so explicitly.
func testProvenance() schema.Provenance {
	return schema.Provenance{
		Method:          schema.ExtractionFile,
		Source:          "testdata/GMP.xml.in",
		SourceSHA256:    "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		ParsedSHA256:    "fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210",
		DeclaredVersion: "@GMP_VERSION@",
		ExtractedAt:     time.Unix(0, 0).UTC(),
	}
}

// minimalCommand returns a single-command schema.Command tree just complete
// enough for Generate to run its validation pass: validateFieldOverrides
// checks the FieldOverrides table's own shape unconditionally, without ever
// needing the bad entry's key to match anything this tree declares.
func minimalCommand() schema.Command {
	return schema.Command{
		Name:       "test_cmd",
		Request:    &schema.Node{Name: "test_cmd"},
		Response:   &schema.Node{Name: "test_cmd_response"},
		Provenance: testProvenance(),
	}
}

// TestCheckProvenanceRejectsMissingFields is a table of every provenance
// field checkProvenance requires, proving each one's absence fails generation
// closed rather than silently generating from an unidentified schema source.
func TestCheckProvenanceRejectsMissingFields(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(p *schema.Provenance)
		wantSub string
	}{
		{"missing Method", func(p *schema.Provenance) { p.Method = "" }, "missing Method"},
		{"missing Source", func(p *schema.Provenance) { p.Source = "" }, "missing Source"},
		{"missing SourceSHA256", func(p *schema.Provenance) { p.SourceSHA256 = "" }, "missing SourceSHA256"},
		{"missing ParsedSHA256", func(p *schema.Provenance) { p.ParsedSHA256 = "" }, "missing ParsedSHA256"},
		{"missing DeclaredVersion", func(p *schema.Provenance) { p.DeclaredVersion = "" }, "missing DeclaredVersion"},
		{"missing ExtractedAt", func(p *schema.Provenance) { p.ExtractedAt = time.Time{} }, "missing ExtractedAt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := minimalCommand()
			tt.mutate(&cmd.Provenance)
			_, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, nil, nil, nil, "", nil)
			if err == nil {
				t.Fatalf("Generate succeeded with %s; want an error", tt.name)
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Fatalf("Generate error = %q, want it to contain %q", err.Error(), tt.wantSub)
			}
		})
	}
}

// TestCheckProvenanceRejectsMalformedFields is a table of provenance values
// that are present but wrong, proving checkProvenance validates content and
// not merely non-emptiness.
func TestCheckProvenanceRejectsMalformedFields(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(p *schema.Provenance)
		wantSub string
	}{
		{"unknown Method", func(p *schema.Provenance) { p.Method = "downloaded" }, "unknown Method"},
		{"whitespace-only Source", func(p *schema.Provenance) { p.Source = "   " }, "missing Source"},
		{"whitespace-only DeclaredVersion", func(p *schema.Provenance) { p.DeclaredVersion = "\t\n" }, "missing DeclaredVersion"},
		{"short SHA-256", func(p *schema.Provenance) { p.SourceSHA256 = "abcd" }, "not 64 lowercase hex"},
		{"malformed ParsedSHA256", func(p *schema.Provenance) { p.ParsedSHA256 = "abcd" }, "ParsedSHA256 \"abcd\" is not 64 lowercase hex"},
		{
			"uppercase SHA-256",
			func(p *schema.Provenance) { p.SourceSHA256 = strings.ToUpper(p.SourceSHA256) },
			"not 64 lowercase hex",
		},
		{
			"non-hex SHA-256",
			func(p *schema.Provenance) { p.SourceSHA256 = strings.Repeat("g", 64) },
			"not 64 lowercase hex",
		},
		{
			"non-UTC ExtractedAt",
			func(p *schema.Provenance) { p.ExtractedAt = p.ExtractedAt.In(time.FixedZone("CET", 3600)) },
			"is not UTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := minimalCommand()
			tt.mutate(&cmd.Provenance)
			_, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, nil, nil, nil, "", nil)
			if err == nil {
				t.Fatalf("Generate succeeded with %s; want an error", tt.name)
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Fatalf("Generate error = %q, want it to contain %q", err.Error(), tt.wantSub)
			}
		})
	}
}

// TestCheckProvenanceRejectsMixedBatch proves that commands disagreeing on
// schema provenance within one Generate call fail closed: they were not
// extracted together, so nothing can vouch for them describing one upstream
// schema.
func TestCheckProvenanceRejectsMixedBatch(t *testing.T) {
	cmdA := minimalCommand()
	cmdA.Name = "cmd_a"
	cmdA.Request = &schema.Node{Name: "cmd_a"}
	cmdA.Response = &schema.Node{Name: "cmd_a_response"}

	cmdB := minimalCommand()
	cmdB.Name = "cmd_b"
	cmdB.Request = &schema.Node{Name: "cmd_b"}
	cmdB.Response = &schema.Node{Name: "cmd_b_response"}
	cmdB.Provenance.SourceSHA256 = strings.Repeat("f", 64)

	_, err := Generate([]schema.Command{cmdA, cmdB}, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, nil, nil, nil, "", nil)
	if err == nil {
		t.Fatal("Generate succeeded with mixed schema provenance; want an error")
	}
	if !strings.Contains(err.Error(), "mixed schema provenance") {
		t.Fatalf("Generate error = %q, want it to mention mixed schema provenance", err)
	}
	if !strings.Contains(err.Error(), "cmd_a") || !strings.Contains(err.Error(), "cmd_b") {
		t.Fatalf("Generate error = %q, want it to name both disagreeing commands", err)
	}
}

// TestCheckProvenanceAllowsDifferentExtractionTimes proves the gate reads
// schema identity, not the clock. Adding a family re-extracts only its own
// snapshots, so every adoption after the first presents a batch of mixed
// timestamps over one unchanged schema.
func TestCheckProvenanceAllowsDifferentExtractionTimes(t *testing.T) {
	cmdA := minimalCommand()
	cmdA.Name = "cmd_a"
	cmdA.Request = &schema.Node{Name: "cmd_a"}
	cmdA.Response = &schema.Node{Name: "cmd_a_response"}

	cmdB := minimalCommand()
	cmdB.Name = "cmd_b"
	cmdB.Request = &schema.Node{Name: "cmd_b"}
	cmdB.Response = &schema.Node{Name: "cmd_b_response"}
	cmdB.Provenance.ExtractedAt = cmdA.Provenance.ExtractedAt.Add(time.Hour)

	if _, err := Generate([]schema.Command{cmdA, cmdB}, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, nil, nil, nil, "", nil); err != nil {
		t.Fatalf("Generate failed on one schema extracted at two times: %v", err)
	}
}

// TestCheckProvenanceAcceptsConsistentBatch proves the mirror case: identical,
// complete provenance across every command in the batch passes.
func TestCheckProvenanceAcceptsConsistentBatch(t *testing.T) {
	cmdA := minimalCommand()
	cmdA.Name = "cmd_a"
	cmdA.Request = &schema.Node{Name: "cmd_a"}
	cmdA.Response = &schema.Node{Name: "cmd_a_response"}

	cmdB := minimalCommand()
	cmdB.Name = "cmd_b"
	cmdB.Request = &schema.Node{Name: "cmd_b"}
	cmdB.Response = &schema.Node{Name: "cmd_b_response"}

	if _, err := Generate([]schema.Command{cmdA, cmdB}, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, nil, nil, nil, "", nil); err != nil {
		t.Fatalf("Generate failed on a consistent provenance batch: %v", err)
	}
}

// TestValidateFieldOverridesRejectsBadEntries is a table of every shape
// validateFieldOverrides must reject, exercised through Generate so the wiring
// from Generate into the check is covered too, not only the check in
// isolation.
func TestValidateFieldOverridesRejectsBadEntries(t *testing.T) {
	tests := []struct {
		name    string
		fo      map[string][]FieldOverride
		wantSub string
	}{
		{
			name: "empty GoType",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Name: "n", GoType: "", Reason: "r"}},
			},
			wantSub: "GoType is empty",
		},
		{
			name: "empty Reason",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Name: "n", GoType: "string", Reason: ""}},
			},
			wantSub: "Reason is empty",
		},
		{
			name: "ordinary override with empty Name",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Name: "", GoType: "string", Reason: "r"}},
			},
			wantSub: "Name is empty",
		},
		{
			name: "Chardata combined with Name",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Chardata: true, Name: "n", GoType: "string", Reason: "r"}},
			},
			wantSub: "mutually exclusive",
		},
		{
			name: "Chardata combined with Attr",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Chardata: true, Attr: true, GoType: "string", Reason: "r"}},
			},
			wantSub: "mutually exclusive",
		},
		{
			name: "Chardata combined with OmitEmpty",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Chardata: true, OmitEmpty: true, GoType: "string", Reason: "r"}},
			},
			wantSub: "mutually exclusive",
		},
		{
			name: "two chardata overrides for one key",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {
					{Chardata: true, GoType: "string", Reason: "r1"},
					{Chardata: true, GoType: "string", Reason: "r2"},
				},
			},
			wantSub: "more than one chardata override",
		},
		{
			name: "duplicate ordinary (Name,Attr) pair for one key",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {
					{Name: "n", GoType: "string", Reason: "r1"},
					{Name: "n", GoType: "int", Reason: "r2"},
				},
			},
			wantSub: "duplicate wire name",
		},
		{
			name: "syntactically invalid GoType func(",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Name: "n", GoType: "func(", Reason: "r"}},
			},
			wantSub: "not a valid Go type expression",
		},
		{
			name: "syntactically invalid GoType 1+",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Name: "n", GoType: "1+", Reason: "r"}},
			},
			wantSub: "not a valid Go type expression",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Generate([]schema.Command{minimalCommand()}, TypeOverrides, PointerOverrides, nil, tt.fo, nil, nil, nil, nil, nil, "", nil)
			if err == nil {
				t.Fatalf("Generate succeeded; want an error containing %q", tt.wantSub)
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Fatalf("Generate error = %q, want it to contain %q", err.Error(), tt.wantSub)
			}
		})
	}
}

// TestValidateFieldOverridesAcceptsCommittedTable is a regression net: the
// real, committed FieldOverrides table must always pass the table-shape
// checks, so any future entry that fails this test has a defect in the entry
// itself, not in the checker.
func TestValidateFieldOverridesAcceptsCommittedTable(t *testing.T) {
	g := &gen{fo: FieldOverrides}
	if err := g.validateFieldOverrides(); err != nil {
		t.Fatalf("The committed FieldOverrides table failed validation: %v", err)
	}
}

// TestValidateFieldSuppressionsAcceptsCommittedTable is
// TestValidateFieldOverridesAcceptsCommittedTable's counterpart for the
// committed FieldSuppressions/RelocatedFrom relocation pairs.
func TestValidateFieldSuppressionsAcceptsCommittedTable(t *testing.T) {
	if err := validateFieldRelocations(FieldOverrides, FieldSuppressions); err != nil {
		t.Fatalf("The committed FieldSuppressions table failed validation: %v", err)
	}
}

// TestValidateFieldSuppressionsRejectsBadEntries is
// TestValidateFieldOverridesRejectsBadEntries's counterpart for
// FieldSuppressions' own chardata shape checks, exercised directly through
// validateFieldRelocations the same way TestFieldRelocationsRejectMissingCounterpart is.
func TestValidateFieldSuppressionsRejectsBadEntries(t *testing.T) {
	tests := []struct {
		name    string
		fo      map[string][]FieldOverride
		fs      map[string][]FieldSuppression
		wantSub string
	}{
		{
			name: "Chardata combined with Name",
			fs: map[string][]FieldSuppression{
				"test_cmd:/x": {{Chardata: true, Name: "n", RelocatedTo: "test_cmd:/x/y", Reason: "r"}},
			},
			wantSub: "mutually exclusive",
		},
		{
			name: "Chardata combined with Attr",
			fs: map[string][]FieldSuppression{
				"test_cmd:/x": {{Chardata: true, Attr: true, RelocatedTo: "test_cmd:/x/y", Reason: "r"}},
			},
			wantSub: "mutually exclusive",
		},
		{
			name: "two chardata suppressions for one key",
			fo: map[string][]FieldOverride{
				"test_cmd:/x/y": {{Chardata: true, GoType: "string", RelocatedFrom: "test_cmd:/x", Reason: "r1"}},
			},
			fs: map[string][]FieldSuppression{
				"test_cmd:/x": {
					{Chardata: true, RelocatedTo: "test_cmd:/x/y", Reason: "r1"},
					{Chardata: true, RelocatedTo: "test_cmd:/x/z", Reason: "r2"},
				},
			},
			wantSub: "more than one chardata suppression",
		},
		{
			name: "chardata suppression with no FieldOverrides counterpart",
			fs: map[string][]FieldSuppression{
				"test_cmd:/x": {{Chardata: true, RelocatedTo: "test_cmd:/x/y", Reason: "r"}},
			},
			wantSub: "no FieldOverrides entry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFieldRelocations(tt.fo, tt.fs)
			if err == nil {
				t.Fatalf("validateFieldRelocations succeeded; want an error containing %q", tt.wantSub)
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Fatalf("validateFieldRelocations error = %q, want it to contain %q", err.Error(), tt.wantSub)
			}
		})
	}
}

// TestFieldRelocationsRejectMissingCounterpart proves validateFieldRelocations
// fails closed in both directions: a FieldSuppressions entry with no matching
// FieldOverrides.RelocatedFrom back-reference, and a FieldOverride.RelocatedFrom
// with no matching FieldSuppressions entry. Either half deleted alone must
// break the build, not silently leave a suppressed field with nowhere to
// reappear or an addition claiming a relocation that never happened.
func TestFieldRelocationsRejectMissingCounterpart(t *testing.T) {
	t.Run("orphaned suppression", func(t *testing.T) {
		fs := map[string][]FieldSuppression{
			"test_cmd:/x/y": {{Name: "n", Attr: true, RelocatedTo: "test_cmd:/x", Reason: "r"}},
		}
		err := validateFieldRelocations(nil, fs)
		if err == nil || !strings.Contains(err.Error(), "no FieldOverrides entry") {
			t.Fatalf("validateFieldRelocations error = %v, want a missing-FieldOverrides error", err)
		}
	})

	t.Run("orphaned RelocatedFrom", func(t *testing.T) {
		fo := map[string][]FieldOverride{
			"test_cmd:/x": {{Name: "n", Attr: true, GoType: "string", RelocatedFrom: "test_cmd:/x/y", Reason: "r"}},
		}
		err := validateFieldRelocations(fo, nil)
		if err == nil || !strings.Contains(err.Error(), "no FieldSuppressions entry") {
			t.Fatalf("validateFieldRelocations error = %v, want a missing-FieldSuppressions error", err)
		}
	})

	t.Run("matching pair passes", func(t *testing.T) {
		fo := map[string][]FieldOverride{
			"test_cmd:/x": {{Name: "n", Attr: true, GoType: "string", RelocatedFrom: "test_cmd:/x/y", Reason: "r"}},
		}
		fs := map[string][]FieldSuppression{
			"test_cmd:/x/y": {{Name: "n", Attr: true, RelocatedTo: "test_cmd:/x", Reason: "r"}},
		}
		if err := validateFieldRelocations(fo, fs); err != nil {
			t.Fatalf("validateFieldRelocations failed on a correctly-paired relocation: %v", err)
		}
	})

	t.Run("matching chardata pair passes", func(t *testing.T) {
		fo := map[string][]FieldOverride{
			"test_cmd:/x/y": {{Chardata: true, GoType: "string", RelocatedFrom: "test_cmd:/x", Reason: "r"}},
		}
		fs := map[string][]FieldSuppression{
			"test_cmd:/x": {{Chardata: true, RelocatedTo: "test_cmd:/x/y", Reason: "r"}},
		}
		if err := validateFieldRelocations(fo, fs); err != nil {
			t.Fatalf("validateFieldRelocations failed on a correctly-paired chardata relocation: %v", err)
		}
	})
}

// commandWithDefects returns minimalCommand() carrying the given defects, so
// checkDefects has something to match against without depending on any real
// committed snapshot.
func commandWithDefects(defects ...schema.Defect) schema.Command {
	c := minimalCommand()
	c.Defects = defects
	return c
}

// commandWithResponseChild returns a minimal command whose response declares
// exactly one top-level child, for exercising fieldOf's RequiredOverrides
// branch through the real Generate() pipeline — needed for any node with
// Children, since goTypeOf resolves those through g.byHash, which only
// g.collect() (run inside Generate) populates.
func commandWithResponseChild(child *schema.Node) schema.Command {
	c := minimalCommand()
	c.Response.Children = []*schema.Node{child}
	return c
}

// TestCheckDefectsRejectsUnallowlisted proves the fail-closed gate: a
// snapshot defect whose {Command, Class, Path} has no DefectAllowlist entry
// must refuse generation, regardless of what Detail says.
func TestCheckDefectsRejectsUnallowlisted(t *testing.T) {
	cmd := commandWithDefects(schema.Defect{Class: schema.DefectEmptyPattern, Path: "/test_cmd/x", Detail: "pattern declares no content"})
	_, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, nil, nil, nil, "", nil)
	if err == nil {
		t.Fatal("Generate succeeded with an unallowlisted defect; want an error")
	}
	if !strings.Contains(err.Error(), "not in DefectAllowlist") {
		t.Fatalf("Generate error = %q, want it to mention DefectAllowlist", err)
	}
}

// TestCheckDefectsAcceptsMatchingEntry proves a {Command, Class, Path} match
// permits generation, and that matching ignores Detail: the allowlist key
// carries no Detail field, so two occurrences that differ only in Detail
// still match the same entry.
func TestCheckDefectsAcceptsMatchingEntry(t *testing.T) {
	cmd := commandWithDefects(schema.Defect{Class: schema.DefectEmptyPattern, Path: "/test_cmd/x", Detail: "some detail text"})
	da := map[DefectKey]DefectAllowlistEntry{
		{Command: "test_cmd", Class: schema.DefectEmptyPattern, Path: "/test_cmd/x"}: {Reason: "test"},
	}
	if _, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, da, nil, nil, "", nil); err != nil {
		t.Fatalf("Generate failed on an allowlisted defect: %v", err)
	}

	cmd2 := commandWithDefects(schema.Defect{Class: schema.DefectEmptyPattern, Path: "/test_cmd/x", Detail: "totally different detail"})
	if _, err := Generate([]schema.Command{cmd2}, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, da, nil, nil, "", nil); err != nil {
		t.Fatalf("Generate failed on an allowlisted defect with a different Detail: %v", err)
	}
}

// TestCheckDefectsMultipleFailuresAreDeterministic proves that more than one
// unallowlisted defect is reported together, sorted, rather than stopping at
// the first — the same "report every problem in one pass" property
// checkFieldOverridesConsumed already gives FieldOverrides.
func TestCheckDefectsMultipleFailuresAreDeterministic(t *testing.T) {
	cmd := commandWithDefects(
		schema.Defect{Class: schema.DefectEmptyPattern, Path: "/test_cmd/z"},
		schema.Defect{Class: schema.DefectEmptyPattern, Path: "/test_cmd/a"},
	)
	_, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, nil, nil, nil, "", nil)
	if err == nil {
		t.Fatal("Generate succeeded with two unallowlisted defects; want an error")
	}
	got := err.Error()
	posA := strings.Index(got, "/test_cmd/a")
	posZ := strings.Index(got, "/test_cmd/z")
	if posA == -1 || posZ == -1 {
		t.Fatalf("Generate error = %q, want both defect paths mentioned", got)
	}
	if posA > posZ {
		t.Fatalf("Generate error = %q, want /test_cmd/a to sort before /test_cmd/z", got)
	}

	second := err
	_, err = Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, nil, nil, nil, "", nil)
	if err == nil || err.Error() != second.Error() {
		t.Fatalf("Generate error was not deterministic across runs: %q vs %q", second, err)
	}
}

// TestValidateDefectAllowlistRejectsBadEntries is a table of every shape
// validateDefectAllowlist must reject, mirroring
// TestValidateFieldOverridesRejectsBadEntries for the FieldOverrides table.
func TestValidateDefectAllowlistRejectsBadEntries(t *testing.T) {
	tests := []struct {
		name    string
		da      map[DefectKey]DefectAllowlistEntry
		wantSub string
	}{
		{
			name:    "empty Command",
			da:      map[DefectKey]DefectAllowlistEntry{{Command: "", Class: schema.DefectEmptyPattern, Path: "/x"}: {Reason: "r"}},
			wantSub: "Command is empty",
		},
		{
			name:    "empty Path",
			da:      map[DefectKey]DefectAllowlistEntry{{Command: "test_cmd", Class: schema.DefectEmptyPattern, Path: ""}: {Reason: "r"}},
			wantSub: "Path is empty",
		},
		{
			name:    "unknown Class",
			da:      map[DefectKey]DefectAllowlistEntry{{Command: "test_cmd", Class: "not-a-real-class", Path: "/x"}: {Reason: "r"}},
			wantSub: "not a known schema.Defect",
		},
		{
			name:    "empty Reason",
			da:      map[DefectKey]DefectAllowlistEntry{{Command: "test_cmd", Class: schema.DefectEmptyPattern, Path: "/x"}: {Reason: ""}},
			wantSub: "Reason is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Generate([]schema.Command{minimalCommand()}, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, tt.da, nil, nil, "", nil)
			if err == nil {
				t.Fatalf("Generate succeeded; want an error containing %q", tt.wantSub)
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Fatalf("Generate error = %q, want it to contain %q", err.Error(), tt.wantSub)
			}
		})
	}
}

// TestCheckDefectAllowlistConsumedRejectsStaleEntry proves the two-sided
// guarantee: an allowlist entry for a command in scope, whose class or path
// no longer matches any generated defect, fails generation just as
// checkFieldOverridesConsumed fails on a stale FieldOverrides entry.
func TestCheckDefectAllowlistConsumedRejectsStaleEntry(t *testing.T) {
	cmd := commandWithDefects(schema.Defect{Class: schema.DefectEmptyPattern, Path: "/test_cmd/x"})

	t.Run("stale path", func(t *testing.T) {
		da := map[DefectKey]DefectAllowlistEntry{
			{Command: "test_cmd", Class: schema.DefectEmptyPattern, Path: "/test_cmd/x"}:    {Reason: "real"},
			{Command: "test_cmd", Class: schema.DefectEmptyPattern, Path: "/test_cmd/gone"}: {Reason: "stale"},
		}
		_, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, da, nil, nil, "", nil)
		if err == nil || !strings.Contains(err.Error(), "never matched a generated defect") {
			t.Fatalf("Generate error = %v, want a stale-entry error", err)
		}
	})

	t.Run("stale class", func(t *testing.T) {
		da := map[DefectKey]DefectAllowlistEntry{
			{Command: "test_cmd", Class: schema.DefectEmptyPattern, Path: "/test_cmd/x"}: {Reason: "real"},
			{Command: "test_cmd", Class: schema.DefectNoPattern, Path: "/test_cmd/x"}:    {Reason: "stale"},
		}
		_, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, da, nil, nil, "", nil)
		if err == nil || !strings.Contains(err.Error(), "never matched a generated defect") {
			t.Fatalf("Generate error = %v, want a stale-entry error", err)
		}
	})
}

// TestCheckDefectAllowlistConsumedIgnoresOutOfScopeEntries proves
// checkDefectAllowlistConsumed's scoping: an allowlist entry for a command
// that is not part of this Generate call is not "stale" for this call, the
// same per-family scoping checkFieldOverridesConsumed already gives
// FieldOverrides — otherwise family-local generation (GenerateAll's whole
// premise) could never pass with a shared allowlist.
func TestCheckDefectAllowlistConsumedIgnoresOutOfScopeEntries(t *testing.T) {
	cmd := commandWithDefects(schema.Defect{Class: schema.DefectEmptyPattern, Path: "/test_cmd/x"})
	da := map[DefectKey]DefectAllowlistEntry{
		{Command: "test_cmd", Class: schema.DefectEmptyPattern, Path: "/test_cmd/x"}:        {Reason: "real"},
		{Command: "some_other_command", Class: schema.DefectEmptyPattern, Path: "/other/y"}: {Reason: "out of scope"},
	}
	if _, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, da, nil, nil, "", nil); err != nil {
		t.Fatalf("Generate failed because of an out-of-scope allowlist entry: %v", err)
	}
}

// TestValidateDefectAllowlistAcceptsCommittedTable is a regression net: the
// real, committed DefectAllowlist table must always pass the table-shape
// checks, so any future entry that fails this test has a defect in the entry
// itself, not in the checker.
func TestValidateDefectAllowlistAcceptsCommittedTable(t *testing.T) {
	g := &gen{da: DefectAllowlist}
	if err := g.validateDefectAllowlist(); err != nil {
		t.Fatalf("The committed DefectAllowlist table failed validation: %v", err)
	}
}

// TestDefectAllowlistHasNoOrphans checks the whole corpus at once, which
// checkDefectAllowlistConsumed structurally cannot: it runs per family, so an
// entry naming a command in no family is never anyone's stale entry. That
// gap matters here more than for the other tables — an orphan is a standing
// pre-approval for a defect nothing is generating through any more.
func TestDefectAllowlistHasNoOrphans(t *testing.T) {
	cmds, err := LoadSnapshots(snapshotDir, manifestOrFail(t).AllCommands())
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}

	corpus := map[DefectKey]bool{}
	for _, c := range cmds {
		for _, d := range c.Defects {
			corpus[DefectKey{Command: c.Name, Class: d.Class, Path: d.Path}] = true
		}
	}

	for key := range DefectAllowlist {
		if !corpus[key] {
			t.Errorf("DefectAllowlist entry %s/%s at %s matches no defect in the committed snapshots; "+
				"remove it or re-extract", key.Command, key.Class, key.Path)
		}
	}
}

// TestValidateRequiredOverridesRejectsBadEntries proves RequiredOverrides'
// table-shape checks fail closed the same way validateFieldOverrides and
// validateDefectAllowlist do: a malformed key or missing Reason is a defect
// in the table itself, caught before generation ever looks for a matching
// field.
func TestValidateRequiredOverridesRejectsBadEntries(t *testing.T) {
	tests := []struct {
		name    string
		ro      map[string]RequiredOverride
		wantSub string
	}{
		{
			name:    "key with no dot separator",
			ro:      map[string]RequiredOverride{"test_cmdreport_id": {Reason: "r"}},
			wantSub: `key must be "<command>.<field>"`,
		},
		{
			name:    "empty command",
			ro:      map[string]RequiredOverride{".report_id": {Reason: "r"}},
			wantSub: `key must be "<command>.<field>"`,
		},
		{
			name:    "empty field",
			ro:      map[string]RequiredOverride{"test_cmd.": {Reason: "r"}},
			wantSub: `key must be "<command>.<field>"`,
		},
		{
			name:    "empty Reason",
			ro:      map[string]RequiredOverride{"test_cmd.report_id": {Reason: ""}},
			wantSub: "Reason is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Generate([]schema.Command{minimalCommand()}, TypeOverrides, PointerOverrides, tt.ro, nil, nil, nil, nil, nil, nil, "", nil)
			if err == nil {
				t.Fatalf("Generate succeeded; want an error containing %q", tt.wantSub)
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Fatalf("Generate error = %q, want it to contain %q", err.Error(), tt.wantSub)
			}
		})
	}
}

// TestValidateRequiredOverridesAcceptsCommittedTable is a regression net: the
// real, committed RequiredOverrides table must always pass the table-shape
// checks, so any future entry that fails this test has a defect in the entry
// itself, not in the checker.
func TestValidateRequiredOverridesAcceptsCommittedTable(t *testing.T) {
	g := &gen{ro: RequiredOverrides}
	if err := g.validateRequiredOverrides(); err != nil {
		t.Fatalf("The committed RequiredOverrides table failed validation: %v", err)
	}
}

// TestValidateNameOverridesRejectsBadEntries mirrors
// TestValidateFieldOverridesRejectsBadEntries for the NameOverrides table:
// a malformed key, missing Reason, or invalid GoName is a defect in the
// table itself, caught before generation ever looks for a matching type.
func TestValidateNameOverridesRejectsBadEntries(t *testing.T) {
	tests := []struct {
		name    string
		no      map[string]NameOverride
		wantSub string
	}{
		{
			name:    "key with no colon-slash separator",
			no:      map[string]NameOverride{"test_cmdthing": {GoName: "Thing", Reason: "r"}},
			wantSub: `key must be "<command>:/<path>"`,
		},
		{
			name:    "empty command",
			no:      map[string]NameOverride{":/thing": {GoName: "Thing", Reason: "r"}},
			wantSub: `key must be "<command>:/<path>"`,
		},
		{
			name:    "empty path",
			no:      map[string]NameOverride{"test_cmd:/": {GoName: "Thing", Reason: "r"}},
			wantSub: `key must be "<command>:/<path>"`,
		},
		{
			name:    "empty Reason",
			no:      map[string]NameOverride{"test_cmd:/thing": {GoName: "Thing", Reason: ""}},
			wantSub: "Reason is empty",
		},
		{
			name:    "empty GoName",
			no:      map[string]NameOverride{"test_cmd:/thing": {GoName: "", Reason: "r"}},
			wantSub: "not a valid exported Go identifier",
		},
		{
			name:    "unexported GoName",
			no:      map[string]NameOverride{"test_cmd:/thing": {GoName: "thing", Reason: "r"}},
			wantSub: "not a valid exported Go identifier",
		},
		{
			name:    "GoName is not a valid identifier",
			no:      map[string]NameOverride{"test_cmd:/thing": {GoName: "Thing Two", Reason: "r"}},
			wantSub: "not a valid exported Go identifier",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Generate([]schema.Command{minimalCommand()}, TypeOverrides, PointerOverrides, nil, nil, nil, nil, nil, tt.no, nil, "", nil)
			if err == nil {
				t.Fatalf("Generate succeeded; want an error containing %q", tt.wantSub)
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Fatalf("Generate error = %q, want it to contain %q", err.Error(), tt.wantSub)
			}
		})
	}
}

// TestValidateNameOverridesAcceptsCommittedTable is a regression net: the
// real, committed NameOverrides table must always pass the table-shape
// checks, so any future entry that fails this test has a defect in the entry
// itself, not in the checker.
func TestValidateNameOverridesAcceptsCommittedTable(t *testing.T) {
	g := &gen{no: NameOverrides}
	if err := g.validateNameOverrides(); err != nil {
		t.Fatalf("The committed NameOverrides table failed validation: %v", err)
	}
}

// TestValidateFieldTypeOverridesRejectsBadEntries mirrors
// TestValidateFieldOverridesRejectsBadEntries for the FieldTypeOverrides
// table: GoType is emitted directly into generated Go source, so an invalid
// value must fail here rather than surface later as a compile error with
// Generate having already reported success.
func TestValidateFieldTypeOverridesRejectsBadEntries(t *testing.T) {
	tests := []struct {
		name    string
		tov     map[string][]FieldTypeOverride
		wantSub string
	}{
		{
			name:    "key with no colon-slash separator",
			tov:     map[string][]FieldTypeOverride{"test_cmdthing": {{Name: "n", GoType: "string", Reason: "r"}}},
			wantSub: `key must be "<command>:/<path>"`,
		},
		{
			name:    "empty command",
			tov:     map[string][]FieldTypeOverride{":/thing": {{Name: "n", GoType: "string", Reason: "r"}}},
			wantSub: `key must be "<command>:/<path>"`,
		},
		{
			name:    "empty path",
			tov:     map[string][]FieldTypeOverride{"test_cmd:/": {{Name: "n", GoType: "string", Reason: "r"}}},
			wantSub: `key must be "<command>:/<path>"`,
		},
		{
			name:    "empty Name",
			tov:     map[string][]FieldTypeOverride{"test_cmd:/thing": {{Name: "", GoType: "string", Reason: "r"}}},
			wantSub: "Name is empty",
		},
		{
			name:    "empty Reason",
			tov:     map[string][]FieldTypeOverride{"test_cmd:/thing": {{Name: "n", GoType: "string", Reason: ""}}},
			wantSub: "Reason is empty",
		},
		{
			name:    "empty GoType",
			tov:     map[string][]FieldTypeOverride{"test_cmd:/thing": {{Name: "n", GoType: "", Reason: "r"}}},
			wantSub: "GoType is empty",
		},
		{
			name:    "syntactically invalid GoType",
			tov:     map[string][]FieldTypeOverride{"test_cmd:/thing": {{Name: "n", GoType: "func(", Reason: "r"}}},
			wantSub: "not a valid Go type expression",
		},
		{
			name: "duplicate (Name,Attr) pair for one key",
			tov: map[string][]FieldTypeOverride{
				"test_cmd:/thing": {
					{Name: "n", GoType: "string", Reason: "r1"},
					{Name: "n", GoType: "float64", Reason: "r2"},
				},
			},
			wantSub: "duplicate wire name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Generate([]schema.Command{minimalCommand()}, TypeOverrides, PointerOverrides, nil, nil, nil, nil, nil, nil, tt.tov, "", nil)
			if err == nil {
				t.Fatalf("Generate succeeded; want an error containing %q", tt.wantSub)
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Fatalf("Generate error = %q, want it to contain %q", err.Error(), tt.wantSub)
			}
		})
	}
}

// TestValidateFieldTypeOverridesAcceptsCommittedTable is a regression net:
// the real, committed FieldTypeOverrides table must always pass the
// table-shape checks, so any future entry that fails this test has a defect
// in the entry itself, not in the checker.
func TestValidateFieldTypeOverridesAcceptsCommittedTable(t *testing.T) {
	g := &gen{tov: FieldTypeOverrides}
	if err := g.validateFieldTypeOverrides(); err != nil {
		t.Fatalf("The committed FieldTypeOverrides table failed validation: %v", err)
	}
}

// TestCheckRequiredOverridesConsumedRejectsStaleEntry proves the fail-closed
// guarantee finding 1 asked for: an override entry for a command in scope
// that never matched a generated required field — a typo'd field name here —
// fails generation rather than silently doing nothing.
func TestCheckRequiredOverridesConsumedRejectsStaleEntry(t *testing.T) {
	cmd := commandWithResponseChild(&schema.Node{
		Name: "report_id", Kind: schema.KindElement, Card: schema.CardRequired, Type: "text",
	})
	ro := map[string]RequiredOverride{"test_cmd.reportid": {Reason: "typo'd field name"}}
	_, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, ro, nil, nil, nil, nil, nil, nil, "", nil)
	if err == nil || !strings.Contains(err.Error(), "never matched a generated required field") {
		t.Fatalf("Generate error = %v, want a stale-entry error", err)
	}
}

// TestCheckRequiredOverridesConsumedRejectsWrongCardinality proves an entry
// aimed at a field that is not schema-required (fieldOf only ever consults
// RequiredOverrides from its CardRequired branch) is stale in exactly the
// same way a typo is, rather than being silently ignored.
func TestCheckRequiredOverridesConsumedRejectsWrongCardinality(t *testing.T) {
	cmd := commandWithResponseChild(&schema.Node{
		Name: "report_id", Kind: schema.KindElement, Card: schema.CardOptional, Type: "text",
	})
	ro := map[string]RequiredOverride{"test_cmd.report_id": {Reason: "field is optional, not required"}}
	_, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, ro, nil, nil, nil, nil, nil, nil, "", nil)
	if err == nil || !strings.Contains(err.Error(), "never matched a generated required field") {
		t.Fatalf("Generate error = %v, want a stale-entry error", err)
	}
}

// TestCheckRequiredOverridesConsumedIgnoresOutOfScopeEntries mirrors
// TestCheckDefectAllowlistConsumedIgnoresOutOfScopeEntries: an entry for a
// command that is not part of this Generate call is not "stale" for this
// call, the per-family scoping GenerateAll's family-local generation depends
// on for a table shared across families.
func TestCheckRequiredOverridesConsumedIgnoresOutOfScopeEntries(t *testing.T) {
	cmd := commandWithResponseChild(&schema.Node{
		Name: "report_id", Kind: schema.KindElement, Card: schema.CardRequired, Type: "text",
	})
	ro := map[string]RequiredOverride{
		"test_cmd.report_id":         {Reason: "real"},
		"some_other_command.field_x": {Reason: "out of scope for this call"},
	}
	if _, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, ro, nil, nil, nil, nil, nil, nil, "", nil); err != nil {
		t.Fatalf("Generate failed because of an out-of-scope RequiredOverrides entry: %v", err)
	}
}

// TestRequiredOverrideRejectsUnsupportedTargets is finding 2's regression
// test: RequiredOverrides is only evidenced for a plain required string leaf
// (resume_task/start_task's report_id). A non-string scalar or a structured
// field would additionally need a pointer to stay distinguishable from its
// zero value — a bigger, unevidenced API change — so Generate must refuse
// rather than silently reshaping the field into a pointer.
func TestRequiredOverrideRejectsUnsupportedTargets(t *testing.T) {
	tests := []struct {
		name  string
		child *schema.Node
	}{
		{
			name:  "non-string scalar",
			child: &schema.Node{Name: "count", Kind: schema.KindElement, Card: schema.CardRequired, Type: "integer"},
		},
		{
			name: "structured field with children",
			child: &schema.Node{
				Name: "nested", Kind: schema.KindElement, Card: schema.CardRequired,
				Children: []*schema.Node{{Name: "leaf", Kind: schema.KindElement, Card: schema.CardRequired, Type: "text"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := commandWithResponseChild(tt.child)
			ro := map[string]RequiredOverride{"test_cmd." + tt.child.Name: {Reason: "r"}}
			_, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, ro, nil, nil, nil, nil, nil, nil, "", nil)
			if err == nil || !strings.Contains(err.Error(), "only evidenced for plain string fields") {
				t.Fatalf("Generate error = %v, want a plain-string-only rejection", err)
			}
		})
	}
}
