package codegen

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/brennoo/go-gmp/internal/schema"
)

// validateOrSkip runs Validate for the family owning the named command,
// skipping (not failing) when no "go" toolchain is available to build and run
// the reflect helper GeneratedGoRoots depends on — the same accommodation
// TestGeneratedCodeCompiles makes.
//
// One family's commands, not the whole manifest: GeneratedGoRoots calls
// Generate directly, which is scoped to a single family by GenerateAll's
// contract, so validating per family is what keeps the "after" side identical
// to the bytes the golden holds for that family.
func validateOrSkip(t *testing.T, command string) (before, after []schema.Result) {
	t.Helper()
	fam, ok := manifestOrFail(t).FamilyOf(command)
	if !ok {
		t.Fatalf("The manifest no longer declares %s", command)
	}
	cmds, err := LoadSnapshots(snapshotDir, fam.Commands)
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}
	before, after, err = Validate(cmds, TypeOverrides, PointerOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist, nil)
	if err != nil {
		if errors.Is(err, ErrNoGoToolchain) {
			t.Skipf("Skipping: %v", err)
		}
		t.Fatalf("Validate failed: %v", err)
	}
	return before, after
}

// counts is a category -> count map, split the same way Report's totals are:
// actionable findings (what commands/ would need to change to fix) and
// suppressed ones (recorded, but not charged against commands/ — see
// schema.Comparer's checkOmitEmpty for why).
type counts struct{ actionable, suppressed map[string]int }

func countFindings(fs []schema.Finding) counts {
	c := counts{actionable: map[string]int{}, suppressed: map[string]int{}}
	for _, f := range fs {
		if f.Suppressed {
			c.suppressed[f.Category]++
		} else {
			c.actionable[f.Category]++
		}
	}
	return c
}

// TestValidateBeforeMatchesKnownConformance pins down what Validate's "before"
// side — the committed commands/ package, reflected by BeforeGoRoots — reports
// for the note family, cross-checked against internal/schema/testdata/
// CONFORMANCE.md's own create_note/modify_note/get_notes sections (committed,
// generated from the same schema.Comparer). It exists to catch a bug in
// BeforeGoRoots or in this package's snapshot loading — not in the generator —
// before that bug gets blamed on the generated side by
// TestValidateGeneratedMatchesOrImproves below.
//
// Since the swap-in (execution-order step 7 of the T-9 prototype plan),
// commands/{create_note,modify_note,get_notes,note_shared}.go *are* the
// generator's output, so "before" and "after" now describe the same source and
// these numbers equal wantAfter's. That is not a weakening: it makes this test
// a drift detector on the committed files — hand-edit a generated struct in a
// way that costs conformance and this fails, independently of whether anyone
// re-ran hack/generate.sh.
//
// For the record, the pre-swap hand-written numbers this used to pin, all of
// which generation closed (see TestValidateGeneratedMatchesOrImproves for the
// per-finding explanation):
//
//	create_note: actionable EXTRA 1
//	modify_note: actionable EXTRA 1, suppressed OMITEMPTY 1
//	get_notes:   actionable MISSING 19, TYPE 4, SCHEMA-DEFECT 2,
//	             suppressed OMITEMPTY 14
func TestValidateBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_note")

	want := map[string]counts{
		"create_note": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_note": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_notes": {
			actionable: map[string]int{schema.CatDefect: 2},
			suppressed: map[string]int{},
		},
	}
	for _, res := range before {
		got := countFindings(res.Findings)
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, got)
	}
}

// TestValidateGeneratedMatchesOrImproves is the pre-swap gate: T-9 acceptance
// criterion 6 requires that swapping in the generated types match-or-improve
// schema conformance for this family, with any change explained rather than
// silently regressed.
//
// What this run actually found, verified against the printed findings while
// writing this test: generation closes every actionable and suppressed
// finding for all three commands except the two pre-existing SCHEMA-DEFECT
// findings on get_notes (both "duplicate-ele: redefined in the same scope"),
// which are untouched byte-for-byte. That is expected and correct:
// SCHEMA-DEFECT findings come from cmd.Defects, the schema's own upstream
// bugs, not from comparing the Go model at all, so no amount of regenerating
// commands/ can change them — schema.Comparer.Compare copies them into the
// result before it ever looks at req/resp.
//
// The specific closures, all of them genuine (not category reshuffles — the
// path either has zero findings after, or keeps only its SCHEMA-DEFECT):
//   - create_note, modify_note: the EXTRA finding on .../nvt/name disappears
//     because the generated NoteNVT (create_note.json/modify_note.json's
//     schema shape) has no Name field to begin with — the hand-written NoteNVT
//     carried one with nothing behind it in the schema.
//   - modify_note: the suppressed OMITEMPTY on .../text disappears because the
//     generated ModifyNote.Text has no omitempty tag, matching the schema's
//     "required".
//   - get_notes: all 19 MISSING findings close (7 request attributes, 12
//     response elements including "notes" and the note/text/* shape, see
//     testdata/golden/get_notes.go for the fields that now exist), the 14
//     suppressed OMITEMPTY findings close (the generated struct's
//     omitempty/pointer choices track the schema's required/optional exactly,
//     per T-5's mapping table), and all 4 TYPE findings close because
//     Note.Active/InUse/Orphan/Writable are commands.GMPBool — the
//     TypeOverrides-driven fix for the T-9 defect #6 / T-5 boolean lexical
//     rule the family was chosen to exercise.
//
// No unexplained finding category appears anywhere in the "after" results:
// every category present after generation was already present before it, and
// no category's count goes up. This test asserts both the specific known-good
// numbers (tight regression protection, expected to need updating in lockstep
// with testdata/golden if the generator's output for this family changes) and
// the general "no unexplained regression" rule (robust to that kind of
// change, and the one that would actually catch a genuine future problem like
// the required-vs-pointer or recursion-cut deltas flagged as risks during
// design — neither of which, it turns out, produces a schema.Comparer finding
// at all for this family, checked directly against the printed diff above).
func TestValidateGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_note")

	wantAfter := map[string]counts{
		"create_note": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_note": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_notes":   {actionable: map[string]int{schema.CatDefect: 2}, suppressed: map[string]int{}},
	}

	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateOverrideBeforeMatchesKnownConformance is the override family's
// counterpart to TestValidateBeforeMatchesKnownConformance.
//
// Since the swap-in, commands/{create_override,modify_override,get_overrides,
// override_shared}.go *are* the generator's output, so "before" and "after"
// now describe the same source and these numbers equal wantAfter's — a drift
// detector on the committed files, not a weaker assertion.
//
// For the record, the pre-swap hand-written numbers this used to pin, all of
// which generation closed (see TestValidateOverrideGeneratedMatchesOrImproves
// for the per-finding explanation):
//
//	create_override, modify_override: actionable EXTRA 2, MISSING 2 each
//	get_overrides: actionable EXTRA 3, MISSING 28, TYPE 6, SCHEMA-DEFECT 2,
//	               suppressed OMITEMPTY 40
func TestValidateOverrideBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_override")

	want := map[string]counts{
		"create_override": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_override": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_overrides": {
			actionable: map[string]int{schema.CatDefect: 2},
			suppressed: map[string]int{},
		},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateOverrideGeneratedMatchesOrImproves is the override family's
// pre-swap gate, run at the same rigour the note family's was.
//
// What this run found while override was still review_only: generation
// closed all 84 findings (44 actionable, 40 suppressed) except the two
// SCHEMA-DEFECT ones on get_overrides, which are byte-identical before and
// after. Those two are cmd.Defects — upstream duplicate-ele bugs in the
// schema itself, copied into the result before Comparer looks at any Go
// type — so no generated output can move them, exactly as for get_notes.
//
//   - create_override, modify_override: the EXTRA 2/MISSING 2 pairs closed
//     because the generated OverrideNVT has no Name/Type fields the schema
//     doesn't declare, and Result/Task became ID-bearing structs instead of
//     bare strings.
//   - get_overrides: all 28 MISSING findings closed (2 request attributes,
//     25 response elements, and 1 chardata site — the schema's <sort>
//     declares bare text alongside its <field> child, which schema.Comparer
//     now understands via Node.Chardata, matched by the generated
//     OverrideSort's own Value field), the 40 suppressed OMITEMPTY findings
//     closed, the 6 TYPE findings closed (GMPBool for lexical booleans,
//     float64 for severity — cvss_base stays int, matching the schema's own
//     declaration), and the 3 EXTRA findings closed (the undeclared @trash
//     request attribute and result-NVT cve/bid are gone).
//
// Nothing regresses and no new category appears, so the assertion below is
// the strongest form: after == {SCHEMA-DEFECT: 2} for get_overrides and
// empty for the other two.
func TestValidateOverrideGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_override")

	wantAfter := map[string]counts{
		"create_override": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_override": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_overrides":   {actionable: map[string]int{schema.CatDefect: 2}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateRoleBeforeMatchesKnownConformance is the role family's
// counterpart to TestValidateBeforeMatchesKnownConformance.
//
// Since the swap-in, commands/{create_role,modify_role,get_roles}.go *are*
// the generator's output, so "before" and "after" now describe the same
// source and these numbers equal wantAfter's — a drift detector on the
// committed files, not a weaker assertion.
//
// For the record, the pre-swap hand-written numbers this used to pin, all of
// which generation closed (see TestValidateRoleGeneratedMatchesOrImproves for
// the per-finding explanation):
//
//	create_role, modify_role: no findings either side
//	get_roles: actionable MISSING 11, TYPE 2, suppressed OMITEMPTY 6
//
// Unlike note and override, get_roles carries zero SCHEMA-DEFECT findings —
// role's schema has no known upstream defect, which was one of the criteria
// that ranked it above the other create/get/modify candidates.
func TestValidateRoleBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_role")

	want := map[string]counts{
		"create_role": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_role": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_roles":   {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateRoleGeneratedMatchesOrImproves is the role family's pre-swap
// gate, run at the same rigour the note and override families were.
//
// What this run found while role was still review_only: generation closed
// every actionable and suppressed finding on all three commands, leaving all
// of them clean: the 11 MISSING findings on get_roles closed (request filter
// attributes plus response elements — see testdata/golden/get_roles.go), the
// 2 TYPE findings closed (Writable/InUse become commands.GMPBool), and the 6
// suppressed OMITEMPTY findings closed (the generated struct's
// omitempty/pointer choices track the schema's required/optional exactly).
// Since role has no SCHEMA-DEFECT findings to begin with, "after" is empty
// for every command in this family.
//
// "Empty" describes what schema.Comparer can see, which is bounded by what
// the schema's own <pattern> declares — it is not a claim that the generated
// model is a lossless mirror of every byte gvmd sends. role_count is a
// concrete counterexample: TestRoleCountBareTotalIsAnAcceptedLoss pins a
// case, identical in kind to create_override's result-as-text gap (see the
// "Accepted evidence" note in docs/review-2026-08.md), where the schema's own
// <pattern> for role_count declares only <e>filtered</e>/<e>page</e> with no
// text token anywhere in it, while both bundled command_get_roles.xml
// examples and a live capture send a bare total ("1", "6") ahead of those two
// children. Comparer has nothing to compare against for content the pattern
// never declares, so this loss is invisible to it — and, per the same
// generate-only-what-the-schema-declares principle that left boolean_atoi
// without an override, it is accepted and documented rather than patched
// into RoleCount as a field the schema itself never asked for.
func TestValidateRoleGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_role")

	wantAfter := map[string]counts{
		"create_role": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_role": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_roles":   {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// roleCountProbe mirrors testdata/golden/get_roles.go's RoleCount field for
// field, byte for byte — keep the two in sync. It exists to decode real
// get_roles fixtures with the exact shape the generator emits today, without
// pulling in the out-of-process GeneratedGoRoots machinery just to observe
// one field.
type roleCountProbe struct {
	XMLName   xml.Name `xml:"get_roles_response"`
	RoleCount struct {
		Filtered int `xml:"filtered"`
		Page     int `xml:"page"`
	} `xml:"role_count"`
}

// roleCountEvidence is roleCountProbe plus a ",chardata" field the generated
// shape does not have. It exists solely to prove what the fixture actually
// contains, independent of what the generated shape can retain — decoding
// the same bytes into roleCountProbe proves the drop; decoding them into this
// type proves there really was something to drop.
type roleCountEvidence struct {
	XMLName   xml.Name `xml:"get_roles_response"`
	RoleCount struct {
		Value    string `xml:",chardata"`
		Filtered int    `xml:"filtered"`
		Page     int    `xml:"page"`
	} `xml:"role_count"`
}

// TestRoleCountBareTotalIsAnAcceptedLoss pins a real, live-confirmed gap the
// role family's otherwise-clean TestValidateRoleGeneratedMatchesOrImproves
// result does not cover: command_get_roles.xml's own <role_count> pattern
// (both its bundled examples, and a live gvmd capture under
// internal/gmptest/testdata/live/get_roles) sends a bare total ahead of
// <filtered>/<page> — "1" in the upstream example, "6" live — but the
// schema's declared <pattern> for role_count is only <e>filtered</e>
// <e>page</e>, with no text token anywhere in it (contrast <sort>, whose
// pattern literally contains the bareword "text" next to <e>field</e>,
// which is exactly what Node.Chardata detects and RoleSort's Value field
// reflects). There is nothing for the schema-side parser to see here: the
// pattern the generator draws from simply never declares this content, so
// generation cannot recover it, and schema.Comparer — which only compares
// against what the pattern declares — correctly reports no finding for it
// either.
//
// This is the same class of gap as create_override's result-as-text loss
// (docs/review-2026-08.md's "Accepted evidence" note): the schema's own
// declared structure and its own examples disagree, and per the
// generate-only-what-the-schema-declares principle, the fix is to document
// the loss, not to add a field the schema never asked for.
//
// Each fixture is decoded twice, into two structurally identical types that
// differ only in one field, so the test proves both halves of the claim
// against the actual wire bytes rather than assuming either: roleCountEvidence
// (which has a ,chardata field) proves the bare total is really there;
// roleCountProbe (an exact mirror of the generated RoleCount) proves that,
// with only Filtered/Page to decode into, that same total has nowhere to go.
// If this test starts failing because Filtered/Page stop decoding, because
// the bare total disappears or changes, or because testdata/golden/get_roles.go
// grows a chardata field that recovers it, update this comment and
// docs/review-2026-08.md together.
func TestRoleCountBareTotalIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantTotal  string
		wantChild  int
	}{
		{"upstream", filepath.Join(gmptestDir, "upstream", "get_roles", "001_response.xml"), "1", 1},
		{"live", filepath.Join(gmptestDir, "live", "get_roles", "001_response.xml"), "6", 6},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			// First, prove the wire actually carries the bare total: decode
			// with a ,chardata field the generated shape lacks.
			var evidence roleCountEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into roleCountEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.RoleCount.Value); got != tc.wantTotal {
				t.Fatalf("%s: role_count bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantTotal)
			}
			if evidence.RoleCount.Filtered != tc.wantChild || evidence.RoleCount.Page != tc.wantChild {
				t.Fatalf("%s: roleCountEvidence children = %+v, want Filtered=Page=%d",
					tc.path, evidence.RoleCount, tc.wantChild)
			}

			// Second, prove the exact generated shape retains the children
			// but has nowhere to put that total.
			var got roleCountProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into roleCountProbe: %v", tc.path, err)
			}
			if got.RoleCount.Filtered != tc.wantChild || got.RoleCount.Page != tc.wantChild {
				t.Errorf("roleCountProbe = %+v, want Filtered=Page=%d (the generated shape changed unexpectedly)",
					got.RoleCount, tc.wantChild)
			}
			// roleCountProbe has no field that could have captured
			// evidence.RoleCount.Value — that silent drop, now proven against
			// the actual fixture bytes rather than assumed, is the accepted
			// loss this test exists to pin.
		})
	}
}

// TestValidateGroupBeforeMatchesKnownConformance is the group family's
// counterpart to TestValidateBeforeMatchesKnownConformance.
//
// Since the swap-in, commands/{create_group,modify_group,get_groups}.go
// *are* the generator's output, so "before" and "after" now describe the
// same source and these numbers equal wantAfter's — a drift detector on the
// committed files, not a weaker assertion.
//
// For the record, the pre-swap hand-written numbers this used to pin, both
// of which generation closed (see TestValidateGroupGeneratedMatchesOrImproves
// for the per-finding explanation):
//
//	create_group: actionable SCHEMA-DEFECT 1 (unchanged by the swap — schema
//	  defects are copied from cmd.Defects before any Go-type comparison ever
//	  runs, so generation cannot move them)
//	modify_group: no findings either side
//	get_groups: actionable TYPE 2, suppressed OMITEMPTY 1
func TestValidateGroupBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_group")

	want := map[string]counts{
		"create_group": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{},
		},
		"modify_group": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_groups":   {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateGroupGeneratedMatchesOrImproves is the group family's post-swap
// gate, run at the same rigour note, override, and role were.
//
// What this run found while group was still review_only: generation closed
// both the TYPE 2 and suppressed OMITEMPTY 1 findings on get_groups
// (Writable/InUse become commands.GMPBool; GroupFilters.Name gains
// omitempty, tracking the schema's optional <name>), leaving get_groups and
// modify_group clean. create_group keeps its SCHEMA-DEFECT 1 untouched:
// SCHEMA-DEFECT findings are copied from cmd.Defects before any Go-type
// comparison ever runs, so no generated output can move them — the same
// reasoning already on record for note, override, and role.
//
// One loss the comparer cannot see, recorded here because it is the same
// class already documented for role: the schema's declared <pattern> for
// group_count is only <e>filtered</e><e>page</e>, no text token anywhere,
// and <field> (inside <sort>) is likewise <e>order</e> with no text token —
// yet both command_get_groups.xml's own bundled example and a live capture
// send a bare total/name ahead of those children. See
// TestGroupCountBareTotalIsAnAcceptedLoss and
// TestGroupSortFieldNameIsAnAcceptedLoss below for both halves, each pinned
// against real fixture bytes the same way TestRoleCountBareTotalIsAnAcceptedLoss
// pins role's.
func TestValidateGroupGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_group")

	wantAfter := map[string]counts{
		"create_group": {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{}},
		"modify_group": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_groups":   {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// groupCountProbe mirrors testdata/golden/get_groups.go's GroupCount field,
// byte for byte — keep the two in sync. It exists to decode real get_groups
// fixtures with the exact shape the generator emits today, without pulling
// in the out-of-process GeneratedGoRoots machinery just to observe one
// field. See roleCountProbe above for the identical role-family precedent.
type groupCountProbe struct {
	XMLName    xml.Name `xml:"get_groups_response"`
	GroupCount struct {
		Filtered int `xml:"filtered"`
		Page     int `xml:"page"`
	} `xml:"group_count"`
}

// groupCountEvidence is groupCountProbe plus a ",chardata" field the
// generated shape does not have, proving what the fixture actually contains
// independent of what the generated shape can retain.
type groupCountEvidence struct {
	XMLName    xml.Name `xml:"get_groups_response"`
	GroupCount struct {
		Value    string `xml:",chardata"`
		Filtered int    `xml:"filtered"`
		Page     int    `xml:"page"`
	} `xml:"group_count"`
}

// TestGroupCountBareTotalIsAnAcceptedLoss is the group family's counterpart
// to TestRoleCountBareTotalIsAnAcceptedLoss: command_get_groups.xml's own
// bundled example, and a live gvmd capture under
// internal/gmptest/testdata/live/get_groups, both send a bare total ahead of
// <filtered>/<page> ("1" in both cases here), but the schema's declared
// <pattern> for group_count is only <e>filtered</e><e>page</e>, with no text
// token anywhere in it. There is nothing for the schema-side parser to see:
// the pattern the generator draws from simply never declares this content,
// so generation cannot recover it, and schema.Comparer — which only
// compares against what the pattern declares — correctly reports no finding
// for it either.
//
// This is the same class of gap as create_override's result-as-text loss and
// role_count's identical shape (docs/review-2026-08.md's "Accepted evidence"
// note): the schema's own declared structure and its own examples disagree,
// and per the generate-only-what-the-schema-declares principle, the fix is
// to document the loss, not to add a field the schema never asked for.
//
// Each fixture is decoded twice, into two structurally identical types that
// differ only in one field, so the test proves both halves of the claim
// against the actual wire bytes rather than assuming either: groupCountEvidence
// (which has a ,chardata field) proves the bare total is really there;
// groupCountProbe (an exact mirror of the generated GroupCount) proves that,
// with only Filtered/Page to decode into, that same total has nowhere to go.
func TestGroupCountBareTotalIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantTotal  string
		wantChild  int
	}{
		{"upstream", filepath.Join(gmptestDir, "upstream", "get_groups", "001_response.xml"), "1", 1},
		{"live", filepath.Join(gmptestDir, "live", "get_groups", "003_response.xml"), "1", 1},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			// First, prove the wire actually carries the bare total: decode
			// with a ,chardata field the generated shape lacks.
			var evidence groupCountEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into groupCountEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.GroupCount.Value); got != tc.wantTotal {
				t.Fatalf("%s: group_count bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantTotal)
			}
			if evidence.GroupCount.Filtered != tc.wantChild || evidence.GroupCount.Page != tc.wantChild {
				t.Fatalf("%s: groupCountEvidence children = %+v, want Filtered=Page=%d",
					tc.path, evidence.GroupCount, tc.wantChild)
			}

			// Second, prove the exact generated shape retains the children
			// but has nowhere to put that total.
			var got groupCountProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into groupCountProbe: %v", tc.path, err)
			}
			if got.GroupCount.Filtered != tc.wantChild || got.GroupCount.Page != tc.wantChild {
				t.Errorf("groupCountProbe = %+v, want Filtered=Page=%d (the generated shape changed unexpectedly)",
					got.GroupCount, tc.wantChild)
			}
			// groupCountProbe has no field that could have captured
			// evidence.GroupCount.Value — that silent drop, now proven
			// against the actual fixture bytes rather than assumed, is the
			// accepted loss this test exists to pin.
		})
	}
}

// groupSortFieldProbe mirrors testdata/golden/get_groups.go's GroupSort/
// GroupField fields, byte for byte — keep the two in sync. Same rationale as
// groupCountProbe: decode real fixtures with the exact shape the generator
// emits today, without the out-of-process GeneratedGoRoots machinery.
type groupSortFieldProbe struct {
	XMLName xml.Name `xml:"get_groups_response"`
	Sort    struct {
		Field struct {
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// groupSortFieldEvidence is groupSortFieldProbe plus a ",chardata" field on
// Field that the generated GroupField does not have, proving what the
// fixture actually contains independent of what the generated shape can
// retain.
type groupSortFieldEvidence struct {
	XMLName xml.Name `xml:"get_groups_response"`
	Sort    struct {
		Field struct {
			Value string `xml:",chardata"`
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// TestGroupSortFieldNameIsAnAcceptedLoss is TestGroupCountBareTotalIsAnAcceptedLoss's
// counterpart for the other half of the same contradiction: command_get_groups.xml's
// own bundled example, and a live gvmd capture, both send the sort field's name
// as bare text ahead of <order> ("name" in both fixtures here), but the
// schema's declared <pattern> for <field> (inside <sort>) is only
// <e>order</e>, with no text token anywhere in it — unlike <sort> itself,
// whose pattern does contain the bareword text and is why GroupSort.Value
// exists. There is nothing for the schema-side parser to see: the pattern
// the generator draws from simply never declares this content, so
// generation cannot recover it, and schema.Comparer — which only compares
// against what the pattern declares — correctly reports no finding for it
// either.
//
// This is the same class of gap as group_count's bare total above and
// role_count's identical shape (docs/review-2026-08.md's "Accepted evidence"
// note): the schema's own declared structure and its own examples disagree,
// and per the generate-only-what-the-schema-declares principle, the fix is
// to document the loss, not to add a field the schema never asked for.
//
// Each fixture is decoded twice, into two structurally identical types that
// differ only in one field, so the test proves both halves of the claim
// against the actual wire bytes rather than assuming either:
// groupSortFieldEvidence (which has a ,chardata field on Field) proves the
// bare field name is really there; groupSortFieldProbe (an exact mirror of
// the generated GroupSort/GroupField) proves that, with only Order to decode
// into, that same name has nowhere to go.
func TestGroupSortFieldNameIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantName   string
		wantOrder  string
	}{
		{"upstream", filepath.Join(gmptestDir, "upstream", "get_groups", "001_response.xml"), "name", "ascending"},
		{"live", filepath.Join(gmptestDir, "live", "get_groups", "003_response.xml"), "name", "ascending"},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			// First, prove the wire actually carries the bare field name:
			// decode with a ,chardata field the generated shape lacks.
			var evidence groupSortFieldEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into groupSortFieldEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.Sort.Field.Value); got != tc.wantName {
				t.Fatalf("%s: sort/field bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantName)
			}
			if evidence.Sort.Field.Order != tc.wantOrder {
				t.Fatalf("%s: groupSortFieldEvidence order = %q, want %q",
					tc.path, evidence.Sort.Field.Order, tc.wantOrder)
			}

			// Second, prove the exact generated shape retains Order but has
			// nowhere to put the field name.
			var got groupSortFieldProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into groupSortFieldProbe: %v", tc.path, err)
			}
			if got.Sort.Field.Order != tc.wantOrder {
				t.Errorf("groupSortFieldProbe order = %q, want %q (the generated shape changed unexpectedly)",
					got.Sort.Field.Order, tc.wantOrder)
			}
			// groupSortFieldProbe has no field that could have captured
			// evidence.Sort.Field.Value — that silent drop, now proven
			// against the actual fixture bytes rather than assumed, is the
			// accepted loss this test exists to pin.
		})
	}
}

// TestValidateTagBeforeMatchesKnownConformance is the tag family's
// counterpart to TestValidateBeforeMatchesKnownConformance.
//
// Since the swap-in, commands/{create_tag,modify_tag,get_tags}.go *are* the
// generator's output, so "before" and "after" now describe the same source
// and these numbers equal wantAfter's — a drift detector on the committed
// files, not a weaker assertion.
//
// For the record, the pre-swap hand-written numbers this used to pin, all of
// which generation closed except get_tags's one SCHEMA-DEFECT (see
// TestValidateTagGeneratedMatchesOrImproves for the per-finding explanation):
//
//	create_tag: actionable EXTRA 3, OMITEMPTY 1, suppressed OMITEMPTY 1
//	get_tags:   actionable MISSING 8, EXTRA 2, CARDINALITY 4, TYPE 4,
//	            SCHEMA-DEFECT 1, suppressed OMITEMPTY 2
//	modify_tag: actionable MISSING 1, EXTRA 3, OMITEMPTY 1
func TestValidateTagBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_tag")

	want := map[string]counts{
		"create_tag": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_tags":   {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{}},
		"modify_tag": {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateTagGeneratedMatchesOrImproves is the tag family's post-swap
// gate, run at the same rigour note, override, role, and group were.
//
// What this run found while tag was still review_only: generation closed
// every actionable and suppressed finding on all three commands except
// get_tags's one pre-existing SCHEMA-DEFECT (/get_tags_response/tag/
// resources/resource/permissions, an empty-pattern site — the same defect
// class note/override/create_group already carry), which is untouched
// byte-for-byte since SCHEMA-DEFECT findings are copied from cmd.Defects
// before any Go-type comparison runs.
//
// The specific closures:
//   - create_tag: the 3 EXTRA findings (resources/resource's name/trash/type)
//     close because the generated TagResource (shared with modify_tag, see
//     testdata/golden/tag_shared.go) only has an ID field — the schema's
//     request-side resource pattern never declares those three. The
//     actionable OMITEMPTY on .../resources (schema says required) closes
//     because the generated CreateTag.Resources is a non-pointer
//     TagResources. The suppressed OMITEMPTY on .../resources/resource/@id
//     closes because the generated TagResource.ID keeps its omitempty,
//     matching the schema's optional attribute.
//   - get_tags: all 8 MISSING findings close (3 request attributes —
//     filt_id, trash, names_only — and 5 response elements/wrappers —
//     filters, sort, tags, tag_count, and the tag/resources wrapper, see
//     testdata/golden/get_tags.go), both EXTRA findings close (the
//     hand-written Tag.Resources's flattened shape and GetTagsResponse's
//     unused Truncated field are both gone — Truncated has no schema
//     counterpart at all and is dropped rather than replaced), all 4
//     CARDINALITY and all 4 TYPE findings close together on the same four
//     fields (Writable/InUse/Active/Orphan become *commands.GMPBool,
//     closing the optional-lexical-boolean pair at once), and both
//     suppressed OMITEMPTY findings close (tag/@id and tag/permissions both
//     gain the omitempty the schema's optional declaration calls for).
//   - modify_tag: the MISSING finding (resources/@action) closes because the
//     generated ModifyTagResources gains an Action field; the 3 EXTRA
//     findings close the same way create_tag's did, via the shared
//     TagResource; the OMITEMPTY finding closes the same way, via
//     TagResource.ID's omitempty.
//
// Nothing regresses and no new category appears: after == {SCHEMA-DEFECT: 1}
// for get_tags and empty for the other two.
func TestValidateTagGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_tag")

	wantAfter := map[string]counts{
		"create_tag": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_tags":   {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{}},
		"modify_tag": {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// tagCountProbe mirrors testdata/golden/get_tags.go's TagCount field, byte
// for byte — keep the two in sync. Same rationale as groupCountProbe: decode
// real get_tags fixtures with the exact shape the generator emits today,
// without pulling in the out-of-process GeneratedGoRoots machinery.
type tagCountProbe struct {
	XMLName  xml.Name `xml:"get_tags_response"`
	TagCount struct {
		Filtered int `xml:"filtered"`
		Page     int `xml:"page"`
	} `xml:"tag_count"`
}

// tagCountEvidence is tagCountProbe plus a ",chardata" field the generated
// shape does not have, proving what the fixture actually contains
// independent of what the generated shape can retain.
type tagCountEvidence struct {
	XMLName  xml.Name `xml:"get_tags_response"`
	TagCount struct {
		Value    string `xml:",chardata"`
		Filtered int    `xml:"filtered"`
		Page     int    `xml:"page"`
	} `xml:"tag_count"`
}

// TestTagCountBareTotalIsAnAcceptedLoss is the tag family's counterpart to
// TestGroupCountBareTotalIsAnAcceptedLoss: live gvmd captures under
// internal/gmptest/testdata/live/get_tags send a bare total ahead of
// <filtered>/<page>, but the schema's declared <pattern> for tag_count is
// only <e>filtered</e><e>page</e>, with no text token anywhere in it. There
// is nothing for the schema-side parser to see: the pattern the generator
// draws from simply never declares this content, so generation cannot
// recover it, and schema.Comparer — which only compares against what the
// pattern declares — correctly reports no finding for it either.
//
// Unlike role_count/group_count, no bundled upstream example exists for
// get_tags (internal/gmptest/testdata/upstream/get_tags has only one
// example and it has no tag_count at all — a filtered-empty result), so this
// test is pinned against all three live captures instead of an
// upstream/live pair. All three matter: 001 and 002 send a bare "0" ahead of
// <filtered>0</filtered>/<page>0</page>, proving the loss is real evidence of
// gvmd's own wire shape rather than an artifact of a nonzero total; 003 sends
// "1", the nonzero case.
//
// Each fixture is decoded twice, into two structurally identical types that
// differ only in one field, so the test proves both halves of the claim
// against the actual wire bytes rather than assuming either: tagCountEvidence
// (which has a ,chardata field) proves the bare total is really there;
// tagCountProbe (an exact mirror of the generated TagCount) proves that,
// with only Filtered/Page to decode into, that same total has nowhere to go.
func TestTagCountBareTotalIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantTotal  string
		wantChild  int
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_tags", "001_response.xml"), "0", 0},
		{"live_002", filepath.Join(gmptestDir, "live", "get_tags", "002_response.xml"), "0", 0},
		{"live_003", filepath.Join(gmptestDir, "live", "get_tags", "003_response.xml"), "1", 1},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence tagCountEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into tagCountEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.TagCount.Value); got != tc.wantTotal {
				t.Fatalf("%s: tag_count bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantTotal)
			}
			if evidence.TagCount.Filtered != tc.wantChild || evidence.TagCount.Page != tc.wantChild {
				t.Fatalf("%s: tagCountEvidence children = %+v, want Filtered=Page=%d",
					tc.path, evidence.TagCount, tc.wantChild)
			}

			var got tagCountProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into tagCountProbe: %v", tc.path, err)
			}
			if got.TagCount.Filtered != tc.wantChild || got.TagCount.Page != tc.wantChild {
				t.Errorf("tagCountProbe = %+v, want Filtered=Page=%d (the generated shape changed unexpectedly)",
					got.TagCount, tc.wantChild)
			}
			// tagCountProbe has no field that could have captured
			// evidence.TagCount.Value — that silent drop, now proven against
			// the actual fixture bytes rather than assumed, is the accepted
			// loss this test exists to pin.
		})
	}
}

// tagSortFieldProbe mirrors testdata/golden/get_tags.go's TagSort/TagField
// fields, byte for byte — keep the two in sync. Same rationale as
// groupSortFieldProbe: decode real fixtures with the exact shape the
// generator emits today, without the out-of-process GeneratedGoRoots
// machinery.
type tagSortFieldProbe struct {
	XMLName xml.Name `xml:"get_tags_response"`
	Sort    struct {
		Field struct {
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// tagSortFieldEvidence is tagSortFieldProbe plus a ",chardata" field on Field
// that the generated TagField does not have, proving what the fixture
// actually contains independent of what the generated shape can retain.
type tagSortFieldEvidence struct {
	XMLName xml.Name `xml:"get_tags_response"`
	Sort    struct {
		Field struct {
			Value string `xml:",chardata"`
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// TestTagSortFieldNameIsAnAcceptedLoss is TestTagCountBareTotalIsAnAcceptedLoss's
// counterpart for the other half of the same contradiction: all three live
// gvmd captures under internal/gmptest/testdata/live/get_tags send the sort
// field's name as bare text ahead of <order>, but the schema's declared
// <pattern> for <field> (inside <sort>) is only <e>order</e>, with no text
// token anywhere in it — unlike <sort> itself, whose pattern does contain
// the bareword text and is why TagSort.Value exists on the golden struct
// already. There is nothing for the schema-side parser to see: the pattern
// the generator draws from simply never declares this content, so
// generation cannot recover it, and schema.Comparer — which only compares
// against what the pattern declares — correctly reports no finding for it
// either.
func TestTagSortFieldNameIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantName   string
		wantOrder  string
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_tags", "001_response.xml"), "name", "ascending"},
		{"live_002", filepath.Join(gmptestDir, "live", "get_tags", "002_response.xml"), "name", "ascending"},
		{"live_003", filepath.Join(gmptestDir, "live", "get_tags", "003_response.xml"), "name", "ascending"},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence tagSortFieldEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into tagSortFieldEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.Sort.Field.Value); got != tc.wantName {
				t.Fatalf("%s: sort/field bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantName)
			}
			if evidence.Sort.Field.Order != tc.wantOrder {
				t.Fatalf("%s: tagSortFieldEvidence order = %q, want %q",
					tc.path, evidence.Sort.Field.Order, tc.wantOrder)
			}

			var got tagSortFieldProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into tagSortFieldProbe: %v", tc.path, err)
			}
			if got.Sort.Field.Order != tc.wantOrder {
				t.Errorf("tagSortFieldProbe order = %q, want %q (the generated shape changed unexpectedly)",
					got.Sort.Field.Order, tc.wantOrder)
			}
			// tagSortFieldProbe has no field that could have captured
			// evidence.Sort.Field.Value — that silent drop, now proven
			// against the actual fixture bytes rather than assumed, is the
			// accepted loss this test exists to pin.
		})
	}
}

// TestValidateScheduleBeforeMatchesKnownConformance is the schedule family's
// counterpart to TestValidateBeforeMatchesKnownConformance. schedule is
// adopted, so commands/{create_schedule,modify_schedule,get_schedules}.go are
// already the generator's output — this is a drift detector on the committed
// files, the same role group and tag's post-adoption "before" tests already
// are.
func TestValidateScheduleBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_schedule")

	want := map[string]counts{
		"create_schedule": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_schedule": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_schedules":   {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{}},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateScheduleGeneratedMatchesOrImproves is the schedule family's
// post-swap gate, run at the same rigour note, override, role, group, and tag
// were.
//
// What this run found while schedule was still review_only: generation
// closed the EXTRA 1 finding on get_schedules (ScheduleField drops its
// hand-written `,chardata` field, which the schema's <field> pattern never
// declared), closed both TYPE 2 findings (Writable/InUse become
// commands.GMPBool), and closed all 17 suppressed OMITEMPTY findings
// (Filters/Sort/Schedules/ScheduleCount and their required children become
// non-pointer, non-omitempty values; the optional <tasks> wrapper becomes a
// pointer). create_schedule and modify_schedule stayed clean throughout.
// get_schedules keeps its SCHEMA-DEFECT 1 untouched: SCHEMA-DEFECT findings
// are copied from cmd.Defects before any Go-type comparison ever runs, so no
// generated output can move them — the same reasoning already on record for
// every prior family.
//
// Two losses the comparer cannot see, recorded here for the same reason
// role's and group's were:
//
//   - ScheduleCount: the schema's declared <pattern> for schedule_count is
//     only <e>filtered</e><e>page</e>, no text token anywhere, yet a live
//     capture sends a bare total ahead of those children. See
//     TestScheduleCountBareTotalIsAnAcceptedLoss below.
//   - ScheduleField: the schema's declared <pattern> for <field> (inside
//     <sort>) is likewise <e>order</e> with no text token, yet the same live
//     captures send the sort field's name as bare text ahead of <order>. See
//     TestScheduleSortFieldNameIsAnAcceptedLoss below.
//
// A third gap worth naming even though it closed no finding and needed no
// test: get_schedules's request-side Trash/Tasks attrs are schema-declared
// booleans (Type:"boolean") but were hand-written as plain string, not bool —
// and schema.Comparer's checkType only raises a TYPE finding when the Go
// field's base type is bool (see conform.go's checkType), never when it is
// already string. The generator still improved these to *GMPBool for typed,
// schema-correct wire encoding, but that improvement tracked no finding
// either side, since none was ever raised for the string case.
func TestValidateScheduleGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_schedule")

	wantAfter := map[string]counts{
		"create_schedule": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_schedule": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_schedules":   {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// scheduleCountProbe mirrors testdata/golden/get_schedules.go's
// ScheduleCount field, byte for byte — keep the two in sync. It exists to
// decode real get_schedules fixtures with the exact shape the generator
// emits today, without pulling in the out-of-process GeneratedGoRoots
// machinery just to observe one field. See roleCountProbe/groupCountProbe
// above for the identical precedent.
type scheduleCountProbe struct {
	XMLName       xml.Name `xml:"get_schedules_response"`
	ScheduleCount struct {
		Filtered int `xml:"filtered"`
		Page     int `xml:"page"`
	} `xml:"schedule_count"`
}

// scheduleCountEvidence is scheduleCountProbe plus a ",chardata" field the
// generated shape does not have, proving what the fixture actually contains
// independent of what the generated shape can retain.
type scheduleCountEvidence struct {
	XMLName       xml.Name `xml:"get_schedules_response"`
	ScheduleCount struct {
		Value    string `xml:",chardata"`
		Filtered int    `xml:"filtered"`
		Page     int    `xml:"page"`
	} `xml:"schedule_count"`
}

// TestScheduleCountBareTotalIsAnAcceptedLoss is the schedule family's
// counterpart to TestRoleCountBareTotalIsAnAcceptedLoss and
// TestGroupCountBareTotalIsAnAcceptedLoss: live gvmd captures under
// internal/gmptest/testdata/live/get_schedules send a bare total ahead of
// <filtered>/<page>, but the schema's declared <pattern> for schedule_count
// is only <e>filtered</e><e>page</e>, with no text token anywhere in it.
// There is nothing for the schema-side parser to see: the pattern the
// generator draws from simply never declares this content, so generation
// cannot recover it, and schema.Comparer — which only compares against what
// the pattern declares — correctly reports no finding for it either.
//
// This is the same class of gap as create_override's result-as-text loss and
// role_count/group_count's identical shape (docs/review-2026-08.md's
// "Accepted evidence" note): the schema's own declared structure and its own
// examples disagree, and per the generate-only-what-the-schema-declares
// principle, the fix is to document the loss, not to add a field the schema
// never asked for.
//
// Each fixture is decoded twice, into two structurally identical types that
// differ only in one field, so the test proves both halves of the claim
// against the actual wire bytes rather than assuming either:
// scheduleCountEvidence (which has a ,chardata field) proves the bare total
// is really there; scheduleCountProbe (an exact mirror of the generated
// ScheduleCount) proves that, with only Filtered/Page to decode into, that
// same total has nowhere to go.
func TestScheduleCountBareTotalIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantTotal  string
		wantChild  int
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_schedules", "001_response.xml"), "0", 0},
		{"live_002", filepath.Join(gmptestDir, "live", "get_schedules", "002_response.xml"), "0", 0},
		{"live_003", filepath.Join(gmptestDir, "live", "get_schedules", "003_response.xml"), "1", 1},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			// First, prove the wire actually carries the bare total: decode
			// with a ,chardata field the generated shape lacks.
			var evidence scheduleCountEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into scheduleCountEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.ScheduleCount.Value); got != tc.wantTotal {
				t.Fatalf("%s: schedule_count bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantTotal)
			}
			if evidence.ScheduleCount.Filtered != tc.wantChild || evidence.ScheduleCount.Page != tc.wantChild {
				t.Fatalf("%s: scheduleCountEvidence children = %+v, want Filtered=Page=%d",
					tc.path, evidence.ScheduleCount, tc.wantChild)
			}

			// Second, prove the exact generated shape retains the children
			// but has nowhere to put that total.
			var got scheduleCountProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into scheduleCountProbe: %v", tc.path, err)
			}
			if got.ScheduleCount.Filtered != tc.wantChild || got.ScheduleCount.Page != tc.wantChild {
				t.Errorf("scheduleCountProbe = %+v, want Filtered=Page=%d (the generated shape changed unexpectedly)",
					got.ScheduleCount, tc.wantChild)
			}
			// scheduleCountProbe has no field that could have captured
			// evidence.ScheduleCount.Value — that silent drop, now proven
			// against the actual fixture bytes rather than assumed, is the
			// accepted loss this test exists to pin.
		})
	}
}

// scheduleSortFieldProbe mirrors testdata/golden/get_schedules.go's
// ScheduleSort/ScheduleField fields, byte for byte — keep the two in sync.
// Same rationale as scheduleCountProbe: decode real fixtures with the exact
// shape the generator emits today, without the out-of-process
// GeneratedGoRoots machinery.
type scheduleSortFieldProbe struct {
	XMLName xml.Name `xml:"get_schedules_response"`
	Sort    struct {
		Field struct {
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// scheduleSortFieldEvidence is scheduleSortFieldProbe plus a ",chardata"
// field on Field that the generated ScheduleField does not have, proving
// what the fixture actually contains independent of what the generated
// shape can retain.
type scheduleSortFieldEvidence struct {
	XMLName xml.Name `xml:"get_schedules_response"`
	Sort    struct {
		Field struct {
			Value string `xml:",chardata"`
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// TestScheduleSortFieldNameIsAnAcceptedLoss is
// TestScheduleCountBareTotalIsAnAcceptedLoss's counterpart for the other
// half of the same contradiction: live gvmd captures under
// internal/gmptest/testdata/live/get_schedules send the sort field's name as
// bare text ahead of <order> ("name" in each fixture here), but the schema's
// declared <pattern> for <field> (inside <sort>) is only <e>order</e>, with
// no text token anywhere in it. There is nothing for the schema-side parser
// to see: the pattern the generator draws from simply never declares this
// content, so generation cannot recover it, and schema.Comparer — which only
// compares against what the pattern declares — correctly reports no finding
// for it either.
//
// This is the same class of gap as schedule_count's bare total above and
// role/group/tag's identical shape (docs/review-2026-08.md's "Accepted
// evidence" note): the schema's own declared structure and its own examples
// disagree, and per the generate-only-what-the-schema-declares principle,
// the fix is to document the loss, not to add a field the schema never
// asked for.
//
// Each fixture is decoded twice, into two structurally identical types that
// differ only in one field, so the test proves both halves of the claim
// against the actual wire bytes rather than assuming either:
// scheduleSortFieldEvidence (which has a ,chardata field on Field) proves
// the bare field name is really there; scheduleSortFieldProbe (an exact
// mirror of the generated ScheduleSort/ScheduleField) proves that, with only
// Order to decode into, that same name has nowhere to go.
func TestScheduleSortFieldNameIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantName   string
		wantOrder  string
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_schedules", "001_response.xml"), "name", "ascending"},
		{"live_002", filepath.Join(gmptestDir, "live", "get_schedules", "002_response.xml"), "name", "ascending"},
		{"live_003", filepath.Join(gmptestDir, "live", "get_schedules", "003_response.xml"), "name", "ascending"},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence scheduleSortFieldEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into scheduleSortFieldEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.Sort.Field.Value); got != tc.wantName {
				t.Fatalf("%s: sort/field bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantName)
			}
			if evidence.Sort.Field.Order != tc.wantOrder {
				t.Fatalf("%s: scheduleSortFieldEvidence order = %q, want %q",
					tc.path, evidence.Sort.Field.Order, tc.wantOrder)
			}

			var got scheduleSortFieldProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into scheduleSortFieldProbe: %v", tc.path, err)
			}
			if got.Sort.Field.Order != tc.wantOrder {
				t.Errorf("scheduleSortFieldProbe order = %q, want %q (the generated shape changed unexpectedly)",
					got.Sort.Field.Order, tc.wantOrder)
			}
			// scheduleSortFieldProbe has no field that could have captured
			// evidence.Sort.Field.Value — that silent drop, now proven
			// against the actual fixture bytes rather than assumed, is the
			// accepted loss this test exists to pin.
		})
	}
}

// TestValidateAlertBeforeMatchesKnownConformance is the alert family's
// counterpart to TestValidateBeforeMatchesKnownConformance. alert is adopted,
// so commands/{create_alert,modify_alert,get_alerts}.go are already the
// generator's output — this is a drift detector on the committed files, the
// same role, group, tag, and schedule's post-adoption "before" tests already
// are.
func TestValidateAlertBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_alert")

	want := map[string]counts{
		"create_alert": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_alert": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_alerts":   {actionable: map[string]int{schema.CatDefect: 2}, suppressed: map[string]int{}},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateAlertGeneratedMatchesOrImproves is the alert family's post-swap
// gate, run at the same rigour every prior family's was.
//
// What this run found while alert was still review_only: generation closed
// create_alert's suppressed OMITEMPTY (Filter becomes a required, non-pointer
// value); closed get_alerts's 2 MISSING findings by modeling
// method/data/credential and method/data/certificate_info (as
// AlertCredential/AlertCertificateInfo on a get_alerts-specific
// AlertMethodData, since create_alert's and modify_alert's simpler <data>
// shape stays the shared AlertData); closed all 4 TYPE findings via
// InUse/Writable/Active/filter.Trash becoming commands.GMPBool; and closed
// all 7 suppressed OMITEMPTY findings by making Filters/Sort/Alerts/AlertCount
// and their required children non-pointer, non-omitempty values, while the
// optional <tasks> and <user_tags> wrappers became pointers. modify_alert's
// Name/Comment also dropped their hand-written *string in favor of a plain
// string (the convention every other adopted family's Name/Comment already
// follows), a change schema.Comparer has no category for either side of, so
// it moved no number here. create_alert and modify_alert stayed otherwise
// clean throughout. get_alerts keeps its SCHEMA-DEFECT 2 untouched:
// SCHEMA-DEFECT findings are copied from cmd.Defects before any Go-type
// comparison ever runs, so no generated output can move them — the same
// reasoning already on record for every prior family.
//
// Two losses the comparer cannot see, recorded here for the same reason
// role's, group's, tag's, and schedule's were:
//
//   - AlertCount: the schema's declared <pattern> for alert_count is only
//     <e>filtered</e><e>page</e>, no text token anywhere, yet live captures
//     send a bare total ahead of those children. See
//     TestAlertCountBareTotalIsAnAcceptedLoss below.
//   - AlertField: the schema's declared <pattern> for <field> (inside <sort>)
//     is likewise <e>order</e> with no text token, yet the same live captures
//     send the sort field's name as bare text ahead of <order>. See
//     TestAlertSortFieldNameIsAnAcceptedLoss below.
func TestValidateAlertGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_alert")

	wantAfter := map[string]counts{
		"create_alert": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_alert": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_alerts":   {actionable: map[string]int{schema.CatDefect: 2}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// alertCountProbe mirrors testdata/golden/get_alerts.go's AlertCount field,
// byte for byte — keep the two in sync. It exists to decode real get_alerts
// fixtures with the exact shape the generator emits today, without pulling in
// the out-of-process GeneratedGoRoots machinery just to observe one field.
// See scheduleCountProbe above for the identical precedent.
type alertCountProbe struct {
	XMLName    xml.Name `xml:"get_alerts_response"`
	AlertCount struct {
		Filtered int `xml:"filtered"`
		Page     int `xml:"page"`
	} `xml:"alert_count"`
}

// alertCountEvidence is alertCountProbe plus a ",chardata" field the
// generated shape does not have, proving what the fixture actually contains
// independent of what the generated shape can retain.
type alertCountEvidence struct {
	XMLName    xml.Name `xml:"get_alerts_response"`
	AlertCount struct {
		Value    string `xml:",chardata"`
		Filtered int    `xml:"filtered"`
		Page     int    `xml:"page"`
	} `xml:"alert_count"`
}

// TestAlertCountBareTotalIsAnAcceptedLoss is the alert family's counterpart
// to TestScheduleCountBareTotalIsAnAcceptedLoss: live gvmd captures under
// internal/gmptest/testdata/live/get_alerts send a bare total ahead of
// <filtered>/<page>, but the schema's declared <pattern> for alert_count is
// only <e>filtered</e><e>page</e>, with no text token anywhere in it. There
// is nothing for the schema-side parser to see: the pattern the generator
// draws from simply never declares this content, so generation cannot
// recover it, and schema.Comparer — which only compares against what the
// pattern declares — correctly reports no finding for it either.
//
// This is the same class of gap as create_override's result-as-text loss and
// role_count/group_count/tag_count/schedule_count's identical shape
// (docs/review-2026-08.md's "Accepted evidence" note): the schema's own
// declared structure and its own examples disagree, and per the
// generate-only-what-the-schema-declares principle, the fix is to document
// the loss, not to add a field the schema never asked for.
//
// Each fixture is decoded twice, into two structurally identical types that
// differ only in one field, so the test proves both halves of the claim
// against the actual wire bytes rather than assuming either: alertCountEvidence
// (which has a ,chardata field) proves the bare total is really there;
// alertCountProbe (an exact mirror of the generated AlertCount) proves that,
// with only Filtered/Page to decode into, that same total has nowhere to go.
func TestAlertCountBareTotalIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantTotal  string
		wantChild  int
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_alerts", "001_response.xml"), "0", 0},
		{"live_002", filepath.Join(gmptestDir, "live", "get_alerts", "002_response.xml"), "0", 0},
		{"live_003", filepath.Join(gmptestDir, "live", "get_alerts", "003_response.xml"), "1", 1},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			// First, prove the wire actually carries the bare total: decode
			// with a ,chardata field the generated shape lacks.
			var evidence alertCountEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into alertCountEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.AlertCount.Value); got != tc.wantTotal {
				t.Fatalf("%s: alert_count bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantTotal)
			}
			if evidence.AlertCount.Filtered != tc.wantChild || evidence.AlertCount.Page != tc.wantChild {
				t.Fatalf("%s: alertCountEvidence children = %+v, want Filtered=Page=%d",
					tc.path, evidence.AlertCount, tc.wantChild)
			}

			// Second, prove the exact generated shape retains the children
			// but has nowhere to put that total.
			var got alertCountProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into alertCountProbe: %v", tc.path, err)
			}
			if got.AlertCount.Filtered != tc.wantChild || got.AlertCount.Page != tc.wantChild {
				t.Errorf("alertCountProbe = %+v, want Filtered=Page=%d (the generated shape changed unexpectedly)",
					got.AlertCount, tc.wantChild)
			}
			// alertCountProbe has no field that could have captured
			// evidence.AlertCount.Value — that silent drop, now proven
			// against the actual fixture bytes rather than assumed, is the
			// accepted loss this test exists to pin.
		})
	}
}

// alertSortFieldProbe mirrors testdata/golden/get_alerts.go's
// AlertSort/AlertField fields, byte for byte — keep the two in sync. Same
// rationale as alertCountProbe: decode real fixtures with the exact shape
// the generator emits today, without the out-of-process GeneratedGoRoots
// machinery.
type alertSortFieldProbe struct {
	XMLName xml.Name `xml:"get_alerts_response"`
	Sort    struct {
		Field struct {
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// alertSortFieldEvidence is alertSortFieldProbe plus a ",chardata" field on
// Field that the generated AlertField does not have, proving what the
// fixture actually contains independent of what the generated shape can
// retain.
type alertSortFieldEvidence struct {
	XMLName xml.Name `xml:"get_alerts_response"`
	Sort    struct {
		Field struct {
			Value string `xml:",chardata"`
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// TestAlertSortFieldNameIsAnAcceptedLoss is
// TestAlertCountBareTotalIsAnAcceptedLoss's counterpart for the other half of
// the same contradiction: live gvmd captures under
// internal/gmptest/testdata/live/get_alerts send the sort field's name as
// bare text ahead of <order> ("name" in each fixture here), but the schema's
// declared <pattern> for <field> (inside <sort>) is only <e>order</e>, with
// no text token anywhere in it. There is nothing for the schema-side parser
// to see: the pattern the generator draws from simply never declares this
// content, so generation cannot recover it, and schema.Comparer — which only
// compares against what the pattern declares — correctly reports no finding
// for it either.
//
// This is the same class of gap as alert_count's bare total above and
// role/group/tag/schedule's identical shape (docs/review-2026-08.md's
// "Accepted evidence" note): the schema's own declared structure and its own
// examples disagree, and per the generate-only-what-the-schema-declares
// principle, the fix is to document the loss, not to add a field the schema
// never asked for.
//
// Each fixture is decoded twice, into two structurally identical types that
// differ only in one field, so the test proves both halves of the claim
// against the actual wire bytes rather than assuming either:
// alertSortFieldEvidence (which has a ,chardata field on Field) proves the
// bare field name is really there; alertSortFieldProbe (an exact mirror of
// the generated AlertSort/AlertField) proves that, with only Order to decode
// into, that same name has nowhere to go.
func TestAlertSortFieldNameIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantName   string
		wantOrder  string
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_alerts", "001_response.xml"), "name", "ascending"},
		{"live_002", filepath.Join(gmptestDir, "live", "get_alerts", "002_response.xml"), "name", "ascending"},
		{"live_003", filepath.Join(gmptestDir, "live", "get_alerts", "003_response.xml"), "name", "ascending"},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence alertSortFieldEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into alertSortFieldEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.Sort.Field.Value); got != tc.wantName {
				t.Fatalf("%s: sort/field bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantName)
			}
			if evidence.Sort.Field.Order != tc.wantOrder {
				t.Fatalf("%s: alertSortFieldEvidence order = %q, want %q",
					tc.path, evidence.Sort.Field.Order, tc.wantOrder)
			}

			var got alertSortFieldProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into alertSortFieldProbe: %v", tc.path, err)
			}
			if got.Sort.Field.Order != tc.wantOrder {
				t.Errorf("alertSortFieldProbe order = %q, want %q (the generated shape changed unexpectedly)",
					got.Sort.Field.Order, tc.wantOrder)
			}
			// alertSortFieldProbe has no field that could have captured
			// evidence.Sort.Field.Value — that silent drop, now proven against
			// the actual fixture bytes rather than assumed, is the accepted
			// loss this test exists to pin.
		})
	}
}

// TestValidateFilterBeforeMatchesKnownConformance is the filter family's
// counterpart to TestValidateBeforeMatchesKnownConformance. filter is
// adopted, so commands/{create_filter,modify_filter,get_filters}.go are
// already the generator's output — this is a drift detector on the committed
// files, the same role every other adopted family's post-adoption "before"
// test plays.
//
// For the record, the pre-adoption hand-written numbers this used to pin, all
// of which generation closed (see TestValidateFilterGeneratedMatchesOrImproves
// for the per-finding explanation): get_filters carried actionable TYPE 2,
// SCHEMA-DEFECT 1, MISSING 1, EXTRA 1, suppressed OMITEMPTY 1; create_filter
// and modify_filter were already conformant.
func TestValidateFilterBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_filter")

	want := map[string]counts{
		"create_filter": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_filter": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_filters": {
			actionable: map[string]int{
				schema.CatDefect: 1,
			},
			suppressed: map[string]int{},
		},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateFilterGeneratedMatchesOrImproves is the filter family's
// post-swap gate, run at the same rigour every prior family's was.
//
// What this run found while filter was still review_only: generation closed
// the TYPE 2 findings (Filter.InUse/.Writable become commands.GMPBool), the
// suppressed OMITEMPTY finding (Filters.Name gains omitempty), and the
// MISSING/EXTRA findings the hand-written FilterInfo2 field produced by being
// wrongly tagged xml:"filters2" instead of the wire's actual "filters".
// SCHEMA-DEFECT 1 stayed untouched (the generic empty-pattern class already
// on record for a dozen other commands, at
// /get_filters_response/filter/alerts/alert/permissions): those findings
// are copied from cmd.Defects before any Go-type comparison runs.
//
// get_filters's <filters> collision (see get_filters.go's generated
// UnmarshalXML/MarshalXML and TestGetFiltersCollisionRoundTrips) can't give
// the second field a real xml:"filters" tag — go vet's structtag analyzer
// statically rejects two fields sharing one tag. It carries xml:"-" instead,
// plus a private gmpwire:"filters" tag recording its true wire identity.
// internal/schema/gomodel.go's addField reads xml:"-" fields' gmpwire tag
// (when present) to make them visible again, by name, with no import or
// method added to commands/*.go; internal/schema/conform.go's walk pairs
// same-name schema children against same-name Go fields by occurrence
// instead of collapsing every occurrence onto the first match. Together
// these let the comparer see both Filters and Filters2 correctly, closing
// the gap to zero MISSING/EXTRA.
func TestValidateFilterGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_filter")

	wantAfter := map[string]counts{
		"create_filter": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_filter": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_filters": {
			actionable: map[string]int{
				schema.CatDefect: 1,
			},
			suppressed: map[string]int{},
		},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// filterCountProbe mirrors testdata/golden/get_filters.go's FilterCount
// field, byte for byte — keep the two in sync. It exists to decode real
// get_filters fixtures with the exact shape the generator emits today,
// without pulling in the out-of-process GeneratedGoRoots machinery just to
// observe one field. See roleCountProbe/groupCountProbe/scheduleCountProbe
// above for the identical precedent.
type filterCountProbe struct {
	XMLName     xml.Name `xml:"get_filters_response"`
	FilterCount struct {
		Filtered int `xml:"filtered"`
		Page     int `xml:"page"`
	} `xml:"filter_count"`
}

// filterCountEvidence is filterCountProbe plus a ",chardata" field the
// generated shape does not have, proving what the fixture actually contains
// independent of what the generated shape can retain.
type filterCountEvidence struct {
	XMLName     xml.Name `xml:"get_filters_response"`
	FilterCount struct {
		Value    string `xml:",chardata"`
		Filtered int    `xml:"filtered"`
		Page     int    `xml:"page"`
	} `xml:"filter_count"`
}

// TestFilterCountBareTotalIsAnAcceptedLoss is the filter family's
// counterpart to TestScheduleCountBareTotalIsAnAcceptedLoss: live gvmd
// captures under internal/gmptest/testdata/live/get_filters send a bare
// total ahead of <filtered>/<page>, but the schema's declared <pattern> for
// filter_count is only <e>filtered</e><e>page</e>, with no text token
// anywhere in it. There is nothing for the schema-side parser to see: the
// pattern the generator draws from simply never declares this content, so
// generation cannot recover it, and schema.Comparer — which only compares
// against what the pattern declares — correctly reports no finding for it
// either.
//
// This is the same class of gap as create_override's result-as-text loss and
// role_count/group_count/tag_count/schedule_count/alert_count's identical
// shape (docs/review-2026-08.md's "Accepted evidence" note): the schema's
// own declared structure and its own examples disagree, and per the
// generate-only-what-the-schema-declares principle, the fix is to document
// the loss, not to add a field the schema never asked for.
//
// Each fixture is decoded twice, into two structurally identical types that
// differ only in one field, so the test proves both halves of the claim
// against the actual wire bytes rather than assuming either:
// filterCountEvidence (which has a ,chardata field) proves the bare total is
// really there; filterCountProbe (an exact mirror of the generated
// FilterCount) proves that, with only Filtered/Page to decode into, that
// same total has nowhere to go.
func TestFilterCountBareTotalIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantTotal  string
		wantChild  int
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_filters", "001_response.xml"), "0", 0},
		{"live_002", filepath.Join(gmptestDir, "live", "get_filters", "002_response.xml"), "0", 0},
		{"live_003", filepath.Join(gmptestDir, "live", "get_filters", "003_response.xml"), "1", 1},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence filterCountEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into filterCountEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.FilterCount.Value); got != tc.wantTotal {
				t.Fatalf("%s: filter_count bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantTotal)
			}
			if evidence.FilterCount.Filtered != tc.wantChild || evidence.FilterCount.Page != tc.wantChild {
				t.Fatalf("%s: filterCountEvidence children = %+v, want Filtered=Page=%d",
					tc.path, evidence.FilterCount, tc.wantChild)
			}

			var got filterCountProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into filterCountProbe: %v", tc.path, err)
			}
			if got.FilterCount.Filtered != tc.wantChild || got.FilterCount.Page != tc.wantChild {
				t.Errorf("filterCountProbe = %+v, want Filtered=Page=%d (the generated shape changed unexpectedly)",
					got.FilterCount, tc.wantChild)
			}
			// filterCountProbe has no field that could have captured
			// evidence.FilterCount.Value — that silent drop, now proven
			// against the actual fixture bytes rather than assumed, is the
			// accepted loss this test exists to pin.
		})
	}
}

// filterSortFieldProbe mirrors testdata/golden/get_filters.go's
// FilterSort/FilterField fields, byte for byte — keep the two in sync. Same
// rationale as filterCountProbe: decode real fixtures with the exact shape
// the generator emits today, without the out-of-process GeneratedGoRoots
// machinery.
type filterSortFieldProbe struct {
	XMLName xml.Name `xml:"get_filters_response"`
	Sort    struct {
		Field struct {
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// filterSortFieldEvidence is filterSortFieldProbe plus a ",chardata" field
// on Field that the generated FilterField does not have, proving what the
// fixture actually contains independent of what the generated shape can
// retain.
type filterSortFieldEvidence struct {
	XMLName xml.Name `xml:"get_filters_response"`
	Sort    struct {
		Field struct {
			Value string `xml:",chardata"`
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// TestFilterSortFieldNameIsAnAcceptedLoss is
// TestFilterCountBareTotalIsAnAcceptedLoss's counterpart for the other half
// of the same contradiction: live gvmd captures under
// internal/gmptest/testdata/live/get_filters send the sort field's name as
// bare text ahead of <order> ("name" in each fixture here), but the schema's
// declared <pattern> for <field> (inside <sort>) is only <e>order</e>, with
// no text token anywhere in it. There is nothing for the schema-side parser
// to see: the pattern the generator draws from simply never declares this
// content, so generation cannot recover it, and schema.Comparer — which only
// compares against what the pattern declares — correctly reports no finding
// for it either.
//
// This is the same class of gap as filter_count's bare total above and
// role/group/tag/schedule/alert's identical shape (docs/review-2026-08.md's
// "Accepted evidence" note): the schema's own declared structure and its own
// examples disagree, and per the generate-only-what-the-schema-declares
// principle, the fix is to document the loss, not to add a field the schema
// never asked for.
//
// Each fixture is decoded twice, into two structurally identical types that
// differ only in one field, so the test proves both halves of the claim
// against the actual wire bytes rather than assuming either:
// filterSortFieldEvidence (which has a ,chardata field on Field) proves the
// bare field name is really there; filterSortFieldProbe (an exact mirror of
// the generated FilterSort/FilterField) proves that, with only Order to
// decode into, that same name has nowhere to go.
func TestFilterSortFieldNameIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantName   string
		wantOrder  string
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_filters", "001_response.xml"), "name", "ascending"},
		{"live_002", filepath.Join(gmptestDir, "live", "get_filters", "002_response.xml"), "name", "ascending"},
		{"live_003", filepath.Join(gmptestDir, "live", "get_filters", "003_response.xml"), "name", "ascending"},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence filterSortFieldEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into filterSortFieldEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.Sort.Field.Value); got != tc.wantName {
				t.Fatalf("%s: sort/field bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantName)
			}
			if evidence.Sort.Field.Order != tc.wantOrder {
				t.Fatalf("%s: filterSortFieldEvidence order = %q, want %q",
					tc.path, evidence.Sort.Field.Order, tc.wantOrder)
			}

			var got filterSortFieldProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into filterSortFieldProbe: %v", tc.path, err)
			}
			if got.Sort.Field.Order != tc.wantOrder {
				t.Errorf("filterSortFieldProbe order = %q, want %q (the generated shape changed unexpectedly)",
					got.Sort.Field.Order, tc.wantOrder)
			}
			// filterSortFieldProbe has no field that could have captured
			// evidence.Sort.Field.Value — that silent drop, now proven
			// against the actual fixture bytes rather than assumed, is the
			// accepted loss this test exists to pin.
		})
	}
}

// TestValidateTLSCertificateBeforeMatchesKnownConformance is the
// tls_certificate family's counterpart to
// TestValidateBeforeMatchesKnownConformance. tls_certificate is adopted, so
// commands/{create_tls_certificate,modify_tls_certificate,
// get_tls_certificates}.go are already the generator's output — this is a
// drift detector on the committed files, the same role every other adopted
// family's post-adoption "before" test plays.
//
// For the record, the pre-adoption hand-written numbers this used to pin, all
// of which generation closed (see TestValidateTLSCertificateGeneratedMatchesOrImproves
// for the per-finding explanation): get_tls_certificates carried actionable
// MISSING 5, EXTRA 5, KIND 2, TYPE 3, suppressed OMITEMPTY 33; create_
// and modify_tls_certificate were already conformant.
func TestValidateTLSCertificateBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_tls_certificate")

	want := map[string]counts{
		"create_tls_certificate": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_tls_certificate": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_tls_certificates":   {actionable: map[string]int{}, suppressed: map[string]int{schema.CatExtra: 2}},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateTLSCertificateGeneratedMatchesOrImproves is the tls_certificate
// family's post-swap gate, run at the same rigour every prior family's was.
//
// What this run found while tls_certificate was still review_only: generation
// closed three of get_tls_certificates's five EXTRA findings and all of
// MISSING/KIND/TYPE, leaving only two EXTRA findings — both suppressed, via
// FieldOverrides — on all three commands:
//
//   - MISSING 5: @filt_id on the request (GetTLSCertificates gains FiltID),
//     and /tls_certificate/sources/source/{location,origin,timestamp,
//     tls_versions}. The hand-written TLSCertificateSource only models
//     Type/ID attrs; the schema's actual source pattern is id(attr) plus
//     required timestamp/tls_versions/origin and optional location, none of
//     which the hand-written struct carries at all.
//   - EXTRA 5: /filters/comment (schema's filters pattern has no comment
//     child), /tls_certificate/sources/source/@type (the hand-written Type
//     attr the schema's source pattern never declares), /truncated (schema's
//     get_tls_certificates_response pattern never declares it), plus the two
//     discussed below — closed differently from the other three.
//   - KIND 2: tls_certificate_count/filtered and /page are schema elements
//     that the hand-written TLSCertificateCount models as attrs; generation
//     fixes both to plain int elements.
//   - TYPE 3: in_use/trust/writable become commands.GMPBool, the same
//     lexical-boolean fix every other family needed.
//   - OMITEMPTY 33 (suppressed): the hand-written struct tags most response
//     fields ,omitempty despite the schema marking them required; inert on
//     decode, so suppressed rather than actionable, per checkOmitEmpty.
//
// Two of the five EXTRA findings deserve individual comment beyond "closed":
//
// /tls_certificate/certificate/#text (the hand-written TLSCertificateData's
// chardata field) is not dropped: FieldOverrides adds a chardata override at
// "get_tls_certificates:/tls_certificate/certificate", so
// TLSCertificateCertificate keeps a Value field alongside Format. The
// schema's certificate pattern truly has no chardata token (Chardata: false
// in the committed snapshot), yet the upstream bundled example
// (internal/gmptest/testdata/upstream/get_tls_certificates/001_response.xml)
// sends real base64 certificate bytes as that element's text — evidence
// strong enough to keep the field rather than silently drop it. The finding
// still fires (the schema still declares no chardata here) but comes back
// Suppressed, with Reason citing that evidence. See
// TestTLSCertificateDataRoundTrips below.
//
// /tls_certificate/valid is handled the same way, via a FieldOverride at
// "get_tls_certificates:/tls_certificate" that adds a Valid GMPBool field.
// Unlike every other EXTRA above, this is not the generated shape losing
// content the schema declares elsewhere — the schema's tls_certificate
// pattern never references a valid child at all (confirmed against both the
// parsed snapshot and the raw command_get_tls_certificates.xml). Live
// traffic proves gvmd sends it regardless
// (internal/gmptest/testdata/live/get_tls_certificates/003_response.xml has
// <valid>1</valid>, matched by the upstream example above), so the override
// keeps the field, and the resulting EXTRA finding is suppressed with a
// Reason citing both fixtures.
//
// Also worth flagging as a review artifact rather than a schema-conformance
// finding: the generated SubjectDN/IssuerDN field names now match the
// hand-written SubjectDN/IssuerDN (and get_reports.go's own SubjectDN)
// casing, since internal/codegen/generate.go's initialisms table gained a
// "dn" entry as part of the same pass. schema.Comparer only compares wire
// tags, not Go identifier casing, so this produces no finding in either
// count above.
func TestValidateTLSCertificateGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_tls_certificate")

	wantAfter := map[string]counts{
		"create_tls_certificate": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_tls_certificate": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_tls_certificates":   {actionable: map[string]int{}, suppressed: map[string]int{schema.CatExtra: 2}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// tlsCertificateCountProbe mirrors testdata/golden/get_tls_certificates.go's
// TLSCertificateCount field, byte for byte — keep the two in sync. Same
// rationale as filterCountProbe: decode real fixtures with the exact shape
// the generator emits today, without the out-of-process GeneratedGoRoots
// machinery.
type tlsCertificateCountProbe struct {
	XMLName             xml.Name `xml:"get_tls_certificates_response"`
	TLSCertificateCount struct {
		Filtered int `xml:"filtered"`
		Page     int `xml:"page"`
	} `xml:"tls_certificate_count"`
}

// tlsCertificateCountEvidence is tlsCertificateCountProbe plus a ",chardata"
// field the generated shape does not have, proving what the fixture
// actually contains independent of what the generated shape can retain.
type tlsCertificateCountEvidence struct {
	XMLName             xml.Name `xml:"get_tls_certificates_response"`
	TLSCertificateCount struct {
		Value    string `xml:",chardata"`
		Filtered int    `xml:"filtered"`
		Page     int    `xml:"page"`
	} `xml:"tls_certificate_count"`
}

// TestTLSCertificateCountBareTotalIsAnAcceptedLoss is the tls_certificate
// family's counterpart to TestFilterCountBareTotalIsAnAcceptedLoss: live
// gvmd captures under internal/gmptest/testdata/live/get_tls_certificates
// send a bare total ahead of <filtered>/<page>, but the schema's declared
// <pattern> for tls_certificate_count is only <e>filtered</e><e>page</e>,
// with no text token anywhere in it. There is nothing for the schema-side
// parser to see: the pattern the generator draws from simply never declares
// this content, so generation cannot recover it, and schema.Comparer —
// which only compares against what the pattern declares — correctly
// reports no finding for it either.
//
// This is the same class of gap as role_count/group_count/tag_count/
// schedule_count/alert_count/filter_count's identical shape
// (docs/review-2026-08.md's "Accepted evidence" note): the schema's own
// declared structure and its own examples disagree, and per the
// generate-only-what-the-schema-declares principle, the fix is to document
// the loss, not to add a field the schema never asked for.
//
// Each fixture is decoded twice, into two structurally identical types that
// differ only in one field, so the test proves both halves of the claim
// against the actual wire bytes rather than assuming either:
// tlsCertificateCountEvidence (which has a ,chardata field) proves the bare
// total is really there; tlsCertificateCountProbe (an exact mirror of the
// generated TLSCertificateCount) proves that, with only Filtered/Page to
// decode into, that same total has nowhere to go.
func TestTLSCertificateCountBareTotalIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantTotal  string
		wantChild  int
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_tls_certificates", "001_response.xml"), "0", 0},
		{"live_002", filepath.Join(gmptestDir, "live", "get_tls_certificates", "002_response.xml"), "0", 0},
		{"live_003", filepath.Join(gmptestDir, "live", "get_tls_certificates", "003_response.xml"), "1", 1},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence tlsCertificateCountEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into tlsCertificateCountEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.TLSCertificateCount.Value); got != tc.wantTotal {
				t.Fatalf("%s: tls_certificate_count bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantTotal)
			}
			if evidence.TLSCertificateCount.Filtered != tc.wantChild || evidence.TLSCertificateCount.Page != tc.wantChild {
				t.Fatalf("%s: tlsCertificateCountEvidence children = %+v, want Filtered=Page=%d",
					tc.path, evidence.TLSCertificateCount, tc.wantChild)
			}

			var got tlsCertificateCountProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into tlsCertificateCountProbe: %v", tc.path, err)
			}
			if got.TLSCertificateCount.Filtered != tc.wantChild || got.TLSCertificateCount.Page != tc.wantChild {
				t.Errorf("tlsCertificateCountProbe = %+v, want Filtered=Page=%d (the generated shape changed unexpectedly)",
					got.TLSCertificateCount, tc.wantChild)
			}
			// tlsCertificateCountProbe has no field that could have captured
			// evidence.TLSCertificateCount.Value — that silent drop, now
			// proven against the actual fixture bytes rather than assumed,
			// is the accepted loss this test exists to pin.
		})
	}
}

// tlsCertificateSortFieldProbe mirrors testdata/golden/get_tls_certificates.
// go's TLSCertificateSort/TLSCertificateField fields, byte for byte — keep
// the two in sync. Same rationale as filterSortFieldProbe: decode real
// fixtures with the exact shape the generator emits today, without the
// out-of-process GeneratedGoRoots machinery.
type tlsCertificateSortFieldProbe struct {
	XMLName xml.Name `xml:"get_tls_certificates_response"`
	Sort    struct {
		Field struct {
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// tlsCertificateSortFieldEvidence is tlsCertificateSortFieldProbe plus a
// ",chardata" field on Field that the generated TLSCertificateField does not
// have, proving what the fixture actually contains independent of what the
// generated shape can retain.
type tlsCertificateSortFieldEvidence struct {
	XMLName xml.Name `xml:"get_tls_certificates_response"`
	Sort    struct {
		Field struct {
			Value string `xml:",chardata"`
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// TestTLSCertificateSortFieldNameIsAnAcceptedLoss is
// TestTLSCertificateCountBareTotalIsAnAcceptedLoss's counterpart for the
// other half of the same contradiction: live gvmd captures under
// internal/gmptest/testdata/live/get_tls_certificates send the sort field's
// name as bare text ahead of <order> ("name" in each fixture here), but the
// schema's declared <pattern> for <field> (inside <sort>) is only
// <e>order</e>, with no text token anywhere in it. There is nothing for the
// schema-side parser to see: the pattern the generator draws from simply
// never declares this content, so generation cannot recover it, and
// schema.Comparer — which only compares against what the pattern declares —
// correctly reports no finding for it either.
//
// This is the same class of gap as tls_certificate_count's bare total above
// and role/group/tag/schedule/alert/filter's identical shape
// (docs/review-2026-08.md's "Accepted evidence" note): the schema's own
// declared structure and its own examples disagree, and per the
// generate-only-what-the-schema-declares principle, the fix is to document
// the loss, not to add a field the schema never asked for.
//
// Each fixture is decoded twice, into two structurally identical types that
// differ only in one field, so the test proves both halves of the claim
// against the actual wire bytes rather than assuming either:
// tlsCertificateSortFieldEvidence (which has a ,chardata field on Field)
// proves the bare field name is really there; tlsCertificateSortFieldProbe
// (an exact mirror of the generated TLSCertificateSort/TLSCertificateField)
// proves that, with only Order to decode into, that same name has nowhere
// to go.
func TestTLSCertificateSortFieldNameIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantName   string
		wantOrder  string
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_tls_certificates", "001_response.xml"), "name", "ascending"},
		{"live_002", filepath.Join(gmptestDir, "live", "get_tls_certificates", "002_response.xml"), "name", "ascending"},
		{"live_003", filepath.Join(gmptestDir, "live", "get_tls_certificates", "003_response.xml"), "name", "ascending"},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence tlsCertificateSortFieldEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into tlsCertificateSortFieldEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.Sort.Field.Value); got != tc.wantName {
				t.Fatalf("%s: sort/field bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantName)
			}
			if evidence.Sort.Field.Order != tc.wantOrder {
				t.Fatalf("%s: tlsCertificateSortFieldEvidence order = %q, want %q",
					tc.path, evidence.Sort.Field.Order, tc.wantOrder)
			}

			var got tlsCertificateSortFieldProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into tlsCertificateSortFieldProbe: %v", tc.path, err)
			}
			if got.Sort.Field.Order != tc.wantOrder {
				t.Errorf("tlsCertificateSortFieldProbe order = %q, want %q (the generated shape changed unexpectedly)",
					got.Sort.Field.Order, tc.wantOrder)
			}
			// tlsCertificateSortFieldProbe has no field that could have
			// captured evidence.Sort.Field.Value — that silent drop, now
			// proven against the actual fixture bytes rather than assumed,
			// is the accepted loss this test exists to pin.
		})
	}
}

// tlsCertificateDataProbeCert mirrors testdata/golden/get_tls_certificates.
// go's TLSCertificateCertificate, byte for byte — keep the two in sync,
// including the FieldOverride-added Value field.
type tlsCertificateDataProbeCert struct {
	Value  string `xml:",chardata"`
	Format string `xml:"format,attr"`
}

// tlsCertificateDataProbe mirrors testdata/golden/get_tls_certificates.go's
// TLSCertificate.Certificate field. Unlike the count/sort-field probes
// above, this is the first probe shaped around a repeated element
// (tls_certificate is a slice); the fixture used here has exactly one
// entry, so index 0 is always the one under test.
type tlsCertificateDataProbe struct {
	XMLName        xml.Name `xml:"get_tls_certificates_response"`
	TLSCertificate []struct {
		Certificate tlsCertificateDataProbeCert `xml:"certificate"`
	} `xml:"tls_certificate"`
}

// TestTLSCertificateDataRoundTrips proves the FieldOverride at
// "get_tls_certificates:/tls_certificate/certificate" (internal/codegen/
// overrides.go) actually recovers the base64 certificate bytes gvmd sends,
// rather than merely silencing the schema-conformance finding for the gap.
//
// The schema's certificate pattern declares only a format attribute, no
// chardata token anywhere (Chardata: false in the committed snapshot and in
// the raw command_get_tls_certificates.xml), yet gvmd's own bundled example
// sends real base64-encoded certificate bytes as that element's text. None
// of the three live captures under
// internal/gmptest/testdata/live/get_tls_certificates set
// include_certificate_data, so <certificate> is always empty there — the
// only fixture that demonstrates this is gvmd's own upstream example
// (internal/gmptest/testdata/upstream/get_tls_certificates/001_response.xml).
//
// The test decodes that fixture, confirms Value actually captured non-empty
// base64 content and Format is "DER", then marshals the decoded certificate
// back out and decodes it again, confirming the same content and format
// attribute survive the round trip.
func TestTLSCertificateDataRoundTrips(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	path := filepath.Join(dir, "..", "gmptest", "testdata", "upstream", "get_tls_certificates", "001_response.xml")

	data, err := os.ReadFile(path) //nolint:gosec // fixed test-fixture path under this package's own tree
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}

	var got tlsCertificateDataProbe
	if err := xml.Unmarshal(data, &got); err != nil {
		t.Fatalf("decoding %s into tlsCertificateDataProbe: %v", path, err)
	}
	if len(got.TLSCertificate) != 1 {
		t.Fatalf("%s: got %d tls_certificate entries, want 1", path, len(got.TLSCertificate))
	}
	cert := got.TLSCertificate[0].Certificate
	value := strings.TrimSpace(cert.Value)
	if value == "" {
		t.Fatalf("%s: certificate bare text is empty — fixture no longer demonstrates the FieldOverride's evidence", path)
	}
	if cert.Format != "DER" {
		t.Fatalf("%s: certificate format = %q, want %q", path, cert.Format, "DER")
	}

	out, err := xml.Marshal(cert)
	if err != nil {
		t.Fatalf("marshaling certificate back out: %v", err)
	}
	var roundTripped tlsCertificateDataProbeCert
	if err := xml.Unmarshal(out, &roundTripped); err != nil {
		t.Fatalf("decoding round-tripped certificate: %v", err)
	}
	if got := strings.TrimSpace(roundTripped.Value); got != value {
		t.Errorf("round-tripped certificate value = %q, want %q", got, value)
	}
	if roundTripped.Format != "DER" {
		t.Errorf("round-tripped certificate format = %q, want %q", roundTripped.Format, "DER")
	}
}

// TestValidatePermissionBeforeMatchesKnownConformance is the permission
// family's counterpart to TestValidateRoleBeforeMatchesKnownConformance.
// permission is adopted, so commands/{create_permission,get_permissions,
// modify_permission}.go are already the generator's output — this is a drift
// detector on the committed files, the same role every other adopted
// family's post-adoption "before" test plays.
//
// For the record, the pre-adoption hand-written numbers this used to pin, all
// of which generation closed except the two SCHEMA-DEFECT findings (see
// TestValidatePermissionGeneratedMatchesOrImproves for the per-finding
// explanation): create_ and modify_permission each carried actionable
// EXTRA 8, OMITEMPTY 2, suppressed OMITEMPTY 3/2; get_permissions carried
// actionable MISSING 1, EXTRA 1, TYPE 5, DEFECT 2, suppressed OMITEMPTY 29.
func TestValidatePermissionBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_permission")

	want := map[string]counts{
		"create_permission": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_permission": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_permissions": {
			actionable: map[string]int{
				schema.CatDefect: 2,
			},
			suppressed: map[string]int{},
		},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidatePermissionGeneratedMatchesOrImproves is the permission family's
// post-swap gate, run at the same rigour every prior family's was.
//
// What this run found while permission was still review_only: generation
// drops the request-side subject/resource's
// fabricated Name/Permissions/Trash/Deleted fields (the family-local
// PermissionSubject/PermissionResource in testdata/golden/
// permission_shared.go carry only the schema's own id+type), closing all 8
// EXTRA findings on both create_permission and modify_permission; their
// ID fields gain the omitempty the schema's optional id attribute calls for,
// and their Type fields lose the omitempty the hand-written struct added
// despite the schema marking it required — closing every OMITEMPTY finding
// on both commands too. On get_permissions, PermissionSort gains a
// `,chardata` field (closing the MISSING sort/#text finding, since unlike
// every prior family's <sort>, get_permissions's own schema pattern declares
// "text" directly alongside <field> — see
// TestPermissionSortFieldNameIsAnAcceptedLoss below for the field-level gap
// that does NOT close), GetPermissionsSubject drops the schema-absent
// Deleted field (closing the one EXTRA finding), every response field's
// omitempty tag now matches the schema's own required/optional split
// (closing all 29 suppressed OMITEMPTY findings), and in_use/resource.
// deleted/resource.trash/subject.trash/writable become commands.GMPBool
// (closing all 5 TYPE findings). The 2 SCHEMA-DEFECT findings on resource/
// permissions and subject/permissions stay untouched — copied from
// cmd.Defects before any Go-type comparison runs, same as every other
// family's empty-pattern class.
func TestValidatePermissionGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_permission")

	wantAfter := map[string]counts{
		"create_permission": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_permission": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_permissions": {
			actionable: map[string]int{
				schema.CatDefect: 2,
			},
			suppressed: map[string]int{},
		},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// permissionCountProbe mirrors testdata/golden/get_permissions.go's
// PermissionCount field, byte for byte — keep the two in sync. See
// filterCountProbe above for the identical precedent.
type permissionCountProbe struct {
	XMLName         xml.Name `xml:"get_permissions_response"`
	PermissionCount struct {
		Filtered int `xml:"filtered"`
		Page     int `xml:"page"`
	} `xml:"permission_count"`
}

// permissionCountEvidence is permissionCountProbe plus a ",chardata" field
// the generated shape does not have, proving what the fixture actually
// contains independent of what the generated shape can retain.
type permissionCountEvidence struct {
	XMLName         xml.Name `xml:"get_permissions_response"`
	PermissionCount struct {
		Value    string `xml:",chardata"`
		Filtered int    `xml:"filtered"`
		Page     int    `xml:"page"`
	} `xml:"permission_count"`
}

// TestPermissionCountBareTotalIsAnAcceptedLoss is the permission family's
// counterpart to TestFilterCountBareTotalIsAnAcceptedLoss: live gvmd
// captures under internal/gmptest/testdata/live/get_permissions, and gvmd's
// own bundled/upstream example under
// internal/gmptest/testdata/upstream/get_permissions, all send a bare total
// ahead of <filtered>/<page>, but the schema's declared <pattern> for
// permission_count is only <e>filtered</e><e>page</e>, with no text token
// anywhere in it. There is nothing for the schema-side parser to see: the
// pattern the generator draws from simply never declares this content, so
// generation cannot recover it, and schema.Comparer — which only compares
// against what the pattern declares — correctly reports no finding for it
// either.
//
// This is the same class of gap as role_count/group_count/tag_count/
// schedule_count/alert_count/filter_count/tls_certificate_count's identical
// shape (docs/review-2026-08.md's "Accepted evidence" note): the schema's
// own declared structure and its own examples disagree, and per the
// generate-only-what-the-schema-declares principle, the fix is to document
// the loss, not to add a field the schema never asked for.
//
// Each fixture is decoded twice, into two structurally identical types that
// differ only in one field, so the test proves both halves of the claim
// against the actual wire bytes rather than assuming either:
// permissionCountEvidence (which has a ,chardata field) proves the bare
// total is really there; permissionCountProbe (an exact mirror of the
// generated PermissionCount) proves that, with only Filtered/Page to decode
// into, that same total has nowhere to go.
func TestPermissionCountBareTotalIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantTotal  string
		wantFilter int
		wantPage   int
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_permissions", "001_response.xml"), "260", 260, 10},
		{"live_002", filepath.Join(gmptestDir, "live", "get_permissions", "002_response.xml"), "260", 260, 2},
		{"live_003", filepath.Join(gmptestDir, "live", "get_permissions", "003_response.xml"), "262", 1, 1},
		{"upstream_001", filepath.Join(gmptestDir, "upstream", "get_permissions", "001_response.xml"), "1", 1, 1},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence permissionCountEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into permissionCountEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.PermissionCount.Value); got != tc.wantTotal {
				t.Fatalf("%s: permission_count bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantTotal)
			}
			if evidence.PermissionCount.Filtered != tc.wantFilter || evidence.PermissionCount.Page != tc.wantPage {
				t.Fatalf("%s: permissionCountEvidence children = %+v, want Filtered=%d Page=%d",
					tc.path, evidence.PermissionCount, tc.wantFilter, tc.wantPage)
			}

			var got permissionCountProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into permissionCountProbe: %v", tc.path, err)
			}
			if got.PermissionCount.Filtered != tc.wantFilter || got.PermissionCount.Page != tc.wantPage {
				t.Errorf("permissionCountProbe = %+v, want Filtered=%d Page=%d (the generated shape changed unexpectedly)",
					got.PermissionCount, tc.wantFilter, tc.wantPage)
			}
			// permissionCountProbe has no field that could have captured
			// evidence.PermissionCount.Value — that silent drop, now proven
			// against the actual fixture bytes rather than assumed, is the
			// accepted loss this test exists to pin.
		})
	}
}

// permissionSortFieldProbe mirrors testdata/golden/get_permissions.go's
// PermissionSort/PermissionField fields, byte for byte, including the
// schema-derived Value chardata field on Sort itself — keep the two in
// sync. Same rationale as permissionCountProbe.
type permissionSortFieldProbe struct {
	XMLName xml.Name `xml:"get_permissions_response"`
	Sort    struct {
		Value string `xml:",chardata"`
		Field struct {
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// permissionSortFieldEvidence is permissionSortFieldProbe plus a ",chardata"
// field on Field that the generated PermissionField does not have, proving
// what the fixture actually contains independent of what the generated
// shape can retain.
type permissionSortFieldEvidence struct {
	XMLName xml.Name `xml:"get_permissions_response"`
	Sort    struct {
		Field struct {
			Value string `xml:",chardata"`
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// TestPermissionSortFieldNameIsAnAcceptedLoss is
// TestPermissionCountBareTotalIsAnAcceptedLoss's counterpart for the other
// half of the same contradiction: live gvmd captures under
// internal/gmptest/testdata/live/get_permissions, and gvmd's own
// bundled/upstream example, send the sort field's name as bare text inside
// <field> ("name" in every fixture here), but the schema's declared
// <pattern> for <field> (inside <sort>) is only <e>order</e>, with no text
// token anywhere in it. There is nothing for the schema-side parser to see:
// the pattern the generator draws from simply never declares this content,
// so generation cannot recover it, and schema.Comparer — which only
// compares against what the pattern declares — correctly reports no
// finding for it either.
//
// get_permissions's schema does declare a text token one level up, directly
// on <sort> itself (<pattern>text<e>field</e></pattern> in
// command_get_permissions.xml), so PermissionSort does get a real,
// schema-derived Value field, closing the MISSING sort/#text finding. But
// that field and the field-name loss pinned here are not the same thing:
// Go's encoding/xml assigns a ",chardata" field only the text nodes that
// are direct children of that element, and in every fixture on record the
// only direct-child text of <sort> is indentation whitespace (or nothing at
// all in the live fixtures' compact XML) — "name" is a child of <field>,
// one level down, not of <sort>. TestPermissionSortValueIsWhitespaceOnly
// below proves PermissionSort.Value never holds anything but that
// whitespace; do not read the MISSING finding closing as evidence that
// PermissionSort.Value captures the observed sort field name.
//
// This is the same class of gap as permission_count's bare total above and
// role/group/tag/schedule/alert/filter/tls_certificate's identical shape
// (docs/review-2026-08.md's "Accepted evidence" note): the schema's own
// declared structure and its own examples disagree, and per the
// generate-only-what-the-schema-declares principle, the fix is to document
// the loss, not to add a field the schema never asked for.
//
// Each fixture is decoded twice, into two structurally identical types that
// differ only in one field, so the test proves both halves of the claim
// against the actual wire bytes rather than assuming either:
// permissionSortFieldEvidence (which has a ,chardata field on Field) proves
// the bare field name is really there; permissionSortFieldProbe (an exact
// mirror of the generated PermissionSort/PermissionField) proves that, with
// only Order to decode into, that same name has nowhere to go.
func TestPermissionSortFieldNameIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantName   string
		wantOrder  string
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_permissions", "001_response.xml"), "name", "ascending"},
		{"live_002", filepath.Join(gmptestDir, "live", "get_permissions", "002_response.xml"), "name", "ascending"},
		{"live_003", filepath.Join(gmptestDir, "live", "get_permissions", "003_response.xml"), "name", "ascending"},
		{"upstream_001", filepath.Join(gmptestDir, "upstream", "get_permissions", "001_response.xml"), "name", "ascending"},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence permissionSortFieldEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into permissionSortFieldEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.Sort.Field.Value); got != tc.wantName {
				t.Fatalf("%s: sort/field bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantName)
			}
			if evidence.Sort.Field.Order != tc.wantOrder {
				t.Fatalf("%s: permissionSortFieldEvidence order = %q, want %q",
					tc.path, evidence.Sort.Field.Order, tc.wantOrder)
			}

			var got permissionSortFieldProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into permissionSortFieldProbe: %v", tc.path, err)
			}
			if got.Sort.Field.Order != tc.wantOrder {
				t.Errorf("permissionSortFieldProbe order = %q, want %q (the generated shape changed unexpectedly)",
					got.Sort.Field.Order, tc.wantOrder)
			}
			// permissionSortFieldProbe.Sort.Field has no field that could
			// have captured evidence.Sort.Field.Value — that silent drop,
			// now proven against the actual fixture bytes rather than
			// assumed, is the accepted loss this test exists to pin.
			if got := strings.TrimSpace(got.Sort.Value); got != "" {
				t.Errorf("permissionSortFieldProbe.Sort.Value = %q, want empty — PermissionSort.Value was expected to hold only indentation whitespace, not %q",
					got, got)
			}
		})
	}
}

// TestPermissionSortValueIsWhitespaceOnly guards against the same
// misreading TestPermissionSortFieldNameIsAnAcceptedLoss's doc comment
// warns about: generation closes the MISSING sort/#text finding by giving
// PermissionSort a real ",chardata" field, because get_permissions's schema
// — uniquely among every family reviewed so far — declares a text token
// directly on <sort>'s own pattern. That closure is schema-accurate, but it
// is not evidence PermissionSort.Value carries the sort field's name: in
// every fixture on record, the only text that is a direct child of <sort>
// (as opposed to a child of the nested <field>) is XML indentation
// whitespace, or nothing at all in the live fixtures' compact encoding. The
// meaningful "name" text lives one level down, inside <field>, which is
// exactly the accepted loss TestPermissionSortFieldNameIsAnAcceptedLoss
// pins — PermissionSort.Value and that loss are not the same field and
// closing the former says nothing about the latter.
func TestPermissionSortValueIsWhitespaceOnly(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	paths := []string{
		filepath.Join(gmptestDir, "live", "get_permissions", "001_response.xml"),
		filepath.Join(gmptestDir, "live", "get_permissions", "002_response.xml"),
		filepath.Join(gmptestDir, "live", "get_permissions", "003_response.xml"),
		filepath.Join(gmptestDir, "upstream", "get_permissions", "001_response.xml"),
	}
	for _, path := range paths {
		t.Run(filepath.Base(filepath.Dir(path))+"/"+filepath.Base(path), func(t *testing.T) {
			data, err := os.ReadFile(path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", path, err)
			}

			var got permissionSortFieldProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into permissionSortFieldProbe: %v", path, err)
			}
			if v := strings.TrimSpace(got.Sort.Value); v != "" {
				t.Fatalf("%s: PermissionSort.Value = %q after trimming, want empty — fixture no longer demonstrates that this field captures only formatting whitespace, not the sort field name",
					path, v)
			}
		})
	}
}

// TestValidateTicketBeforeMatchesKnownConformance is the ticket family's
// counterpart to TestValidatePermissionBeforeMatchesKnownConformance. ticket
// is adopted, so commands/{create_ticket,get_tickets,modify_ticket}.go are
// already the generator's output — this is a drift detector on the
// committed files.
//
// For the record, the pre-adoption hand-written numbers this used to pin,
// all of which generation closed (see TestValidateTicketGeneratedMatchesOrImproves
// for the per-finding explanation): get_tickets carried actionable
// MISSING 6, EXTRA 4, KIND 3, TYPE 3, suppressed OMITEMPTY 23.
func TestValidateTicketBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_ticket")

	want := map[string]counts{
		"create_ticket": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_ticket": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_tickets":   {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateTicketGeneratedMatchesOrImproves is the ticket family's
// post-swap gate. See docs/review-2026-08.md's ticket section for what each
// closed finding was and why.
func TestValidateTicketGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_ticket")

	wantAfter := map[string]counts{
		"create_ticket": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_ticket": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_tickets":   {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// ticketCountProbe mirrors TicketCount in testdata/golden/get_tickets.go.
type ticketCountProbe struct {
	XMLName     xml.Name `xml:"get_tickets_response"`
	TicketCount struct {
		Filtered int `xml:"filtered"`
		Page     int `xml:"page"`
	} `xml:"ticket_count"`
}

// ticketCountEvidence is ticketCountProbe plus a ",chardata" field to prove
// the bare total the generated shape can't retain.
type ticketCountEvidence struct {
	XMLName     xml.Name `xml:"get_tickets_response"`
	TicketCount struct {
		Value    string `xml:",chardata"`
		Filtered int    `xml:"filtered"`
		Page     int    `xml:"page"`
	} `xml:"ticket_count"`
}

// TestTicketCountBareTotalIsAnAcceptedLoss mirrors
// TestPermissionCountBareTotalIsAnAcceptedLoss: fixtures send a bare total
// ahead of <filtered>/<page>, but ticket_count's schema pattern declares no
// text token, so generation can't recover it. See docs/review-2026-08.md.
func TestTicketCountBareTotalIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantTotal  string
		wantFilter int
		wantPage   int
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_tickets", "001_response.xml"), "0", 0, 0},
		{"live_002", filepath.Join(gmptestDir, "live", "get_tickets", "002_response.xml"), "0", 0, 0},
		{"upstream_001", filepath.Join(gmptestDir, "upstream", "get_tickets", "001_response.xml"), "2", 2, 2},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence ticketCountEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into ticketCountEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.TicketCount.Value); got != tc.wantTotal {
				t.Fatalf("%s: ticket_count bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantTotal)
			}
			if evidence.TicketCount.Filtered != tc.wantFilter || evidence.TicketCount.Page != tc.wantPage {
				t.Fatalf("%s: ticketCountEvidence children = %+v, want Filtered=%d Page=%d",
					tc.path, evidence.TicketCount, tc.wantFilter, tc.wantPage)
			}

			var got ticketCountProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into ticketCountProbe: %v", tc.path, err)
			}
			if got.TicketCount.Filtered != tc.wantFilter || got.TicketCount.Page != tc.wantPage {
				t.Errorf("ticketCountProbe = %+v, want Filtered=%d Page=%d (the generated shape changed unexpectedly)",
					got.TicketCount, tc.wantFilter, tc.wantPage)
			}
			// ticketCountProbe has no field that could have captured
			// evidence.TicketCount.Value — that silent drop, now proven
			// against the actual fixture bytes rather than assumed, is the
			// accepted loss this test exists to pin.
		})
	}
}

// ticketSortFieldProbe mirrors TicketSort/TicketField in
// testdata/golden/get_tickets.go.
type ticketSortFieldProbe struct {
	XMLName xml.Name `xml:"get_tickets_response"`
	Sort    struct {
		Value string `xml:",chardata"`
		Field struct {
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// ticketSortFieldEvidence is ticketSortFieldProbe plus a ",chardata" field
// on Field to prove the bare sort field name the generated shape can't
// retain.
type ticketSortFieldEvidence struct {
	XMLName xml.Name `xml:"get_tickets_response"`
	Sort    struct {
		Field struct {
			Value string `xml:",chardata"`
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// TestTicketSortFieldNameIsAnAcceptedLoss mirrors
// TestPermissionSortFieldNameIsAnAcceptedLoss: fixtures send the sort
// field's name as bare text inside <field>, but <field>'s schema pattern
// declares only <order>, so generation can't recover it. This is distinct
// from TicketSort.Value (see TestTicketSortValueIsWhitespaceOnly), which
// closes the schema-level sort/#text finding but never holds this name. See
// docs/review-2026-08.md.
func TestTicketSortFieldNameIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantName   string
		wantOrder  string
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_tickets", "001_response.xml"), "name", "ascending"},
		{"live_002", filepath.Join(gmptestDir, "live", "get_tickets", "002_response.xml"), "name", "ascending"},
		{"upstream_001", filepath.Join(gmptestDir, "upstream", "get_tickets", "001_response.xml"), "name", "ascending"},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence ticketSortFieldEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into ticketSortFieldEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.Sort.Field.Value); got != tc.wantName {
				t.Fatalf("%s: sort/field bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantName)
			}
			if evidence.Sort.Field.Order != tc.wantOrder {
				t.Fatalf("%s: ticketSortFieldEvidence order = %q, want %q",
					tc.path, evidence.Sort.Field.Order, tc.wantOrder)
			}

			var got ticketSortFieldProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into ticketSortFieldProbe: %v", tc.path, err)
			}
			if got.Sort.Field.Order != tc.wantOrder {
				t.Errorf("ticketSortFieldProbe order = %q, want %q (the generated shape changed unexpectedly)",
					got.Sort.Field.Order, tc.wantOrder)
			}
			// ticketSortFieldProbe.Sort.Field has no field that could have
			// captured evidence.Sort.Field.Value — that silent drop, now
			// proven against the actual fixture bytes rather than assumed,
			// is the accepted loss this test exists to pin.
			if got := strings.TrimSpace(got.Sort.Value); got != "" {
				t.Errorf("ticketSortFieldProbe.Sort.Value = %q, want empty — TicketSort.Value was expected to hold only indentation whitespace, not %q",
					got, got)
			}
		})
	}
}

// TestTicketSortValueIsWhitespaceOnly proves TicketSort.Value (the field
// that closes the schema-level sort/#text finding) only ever holds
// formatting whitespace, not the sort field name pinned by
// TestTicketSortFieldNameIsAnAcceptedLoss.
func TestTicketSortValueIsWhitespaceOnly(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	paths := []string{
		filepath.Join(gmptestDir, "live", "get_tickets", "001_response.xml"),
		filepath.Join(gmptestDir, "live", "get_tickets", "002_response.xml"),
		filepath.Join(gmptestDir, "upstream", "get_tickets", "001_response.xml"),
	}
	for _, path := range paths {
		t.Run(filepath.Base(filepath.Dir(path))+"/"+filepath.Base(path), func(t *testing.T) {
			data, err := os.ReadFile(path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", path, err)
			}

			var got ticketSortFieldProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into ticketSortFieldProbe: %v", path, err)
			}
			if v := strings.TrimSpace(got.Sort.Value); v != "" {
				t.Fatalf("%s: TicketSort.Value = %q after trimming, want empty — fixture no longer demonstrates that this field captures only formatting whitespace, not the sort field name",
					path, v)
			}
		})
	}
}

// TestValidateCredentialBeforeMatchesKnownConformance is the credential
// family's counterpart to TestValidatePermissionBeforeMatchesKnownConformance.
// credential is adopted, so commands/{create_credential,get_credentials,
// modify_credential}.go are already the generator's output — this is a drift
// detector on the committed files, the same role every other adopted
// family's post-adoption "before" test plays.
//
// For the record, the pre-adoption hand-written numbers this used to pin, all
// of which generation closed except the 2 SCHEMA-DEFECT findings (see
// TestValidateCredentialGeneratedMatchesOrImproves for the per-finding
// explanation): create_credential and modify_credential each carried
// suppressed OMITEMPTY 2; get_credentials carried actionable MISSING 2,
// EXTRA 1, TYPE 3, DEFECT 2, suppressed OMITEMPTY 2.
func TestValidateCredentialBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_credential")

	want := map[string]counts{
		"create_credential": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_credential": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_credentials": {
			actionable: map[string]int{schema.CatDefect: 2},
			suppressed: map[string]int{},
		},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateCredentialGeneratedMatchesOrImproves is the credential family's
// post-swap gate, run at the same rigour every prior family's was.
//
// What this run found while credential was still review_only: generation
// drops the request-side privacy/algorithm and privacy/password omitempty
// tags (the family-local CredentialPrivacy in
// testdata/golden/credential_shared.go matches the schema's "required"),
// closing the suppressed OMITEMPTY finding on both create_credential and
// modify_credential. On get_credentials, CertificateInfo and Privacy gain
// full generated struct shapes instead of a bare chardata scalar (closing
// both MISSING findings), the schema-absent Truncated field disappears
// (closing the EXTRA finding), allow_insecure/writable/in_use become
// commands.GMPBool (closing all 3 TYPE findings), and Formats/Filters.Name's
// omitempty tags now track the schema's required/optional split exactly
// (closing both suppressed OMITEMPTY findings). The 2 SCHEMA-DEFECT findings
// on scanner/permissions and target/permissions stay untouched — copied from
// cmd.Defects before any Go-type comparison runs, same as every other
// family's empty-pattern class.
func TestValidateCredentialGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_credential")

	wantAfter := map[string]counts{
		"create_credential": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_credential": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_credentials": {
			actionable: map[string]int{schema.CatDefect: 2},
			suppressed: map[string]int{},
		},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// credentialCountProbe mirrors CredentialCount in
// testdata/golden/get_credentials.go.
type credentialCountProbe struct {
	XMLName         xml.Name `xml:"get_credentials_response"`
	CredentialCount struct {
		Filtered int `xml:"filtered"`
		Page     int `xml:"page"`
	} `xml:"credential_count"`
}

// credentialCountEvidence is credentialCountProbe plus a ",chardata" field to
// prove the bare total the generated shape can't retain.
type credentialCountEvidence struct {
	XMLName         xml.Name `xml:"get_credentials_response"`
	CredentialCount struct {
		Value    string `xml:",chardata"`
		Filtered int    `xml:"filtered"`
		Page     int    `xml:"page"`
	} `xml:"credential_count"`
}

// TestCredentialCountBareTotalIsAnAcceptedLoss mirrors
// TestTicketCountBareTotalIsAnAcceptedLoss: fixtures send a bare total ahead
// of <filtered>/<page>, but credential_count's schema pattern declares no
// text token, so generation can't recover it.
func TestCredentialCountBareTotalIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantTotal  string
		wantFilter int
		wantPage   int
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_credentials", "001_response.xml"), "0", 0, 0},
		{"live_003", filepath.Join(gmptestDir, "live", "get_credentials", "003_response.xml"), "1", 1, 1},
		{"live_004", filepath.Join(gmptestDir, "live", "get_credentials", "004_response.xml"), "2", 1, 1},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence credentialCountEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into credentialCountEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.CredentialCount.Value); got != tc.wantTotal {
				t.Fatalf("%s: credential_count bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantTotal)
			}
			if evidence.CredentialCount.Filtered != tc.wantFilter || evidence.CredentialCount.Page != tc.wantPage {
				t.Fatalf("%s: credentialCountEvidence children = %+v, want Filtered=%d Page=%d",
					tc.path, evidence.CredentialCount, tc.wantFilter, tc.wantPage)
			}

			var got credentialCountProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into credentialCountProbe: %v", tc.path, err)
			}
			if got.CredentialCount.Filtered != tc.wantFilter || got.CredentialCount.Page != tc.wantPage {
				t.Errorf("credentialCountProbe = %+v, want Filtered=%d Page=%d (the generated shape changed unexpectedly)",
					got.CredentialCount, tc.wantFilter, tc.wantPage)
			}
			// credentialCountProbe has no field that could have captured
			// evidence.CredentialCount.Value — that silent drop, now proven
			// against the actual fixture bytes rather than assumed, is the
			// accepted loss this test exists to pin.
		})
	}
}

// credentialSortFieldProbe mirrors CredentialSort/CredentialField in
// testdata/golden/get_credentials.go.
type credentialSortFieldProbe struct {
	XMLName xml.Name `xml:"get_credentials_response"`
	Sort    struct {
		Value string `xml:",chardata"`
		Field struct {
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// credentialSortFieldEvidence is credentialSortFieldProbe plus a ",chardata"
// field on Field to prove the bare sort field name the generated shape can't
// retain.
type credentialSortFieldEvidence struct {
	XMLName xml.Name `xml:"get_credentials_response"`
	Sort    struct {
		Field struct {
			Value string `xml:",chardata"`
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// TestCredentialSortFieldNameIsAnAcceptedLoss mirrors
// TestTicketSortFieldNameIsAnAcceptedLoss: fixtures send the sort field's
// name as bare text inside <field>, but <field>'s schema pattern declares
// only <order>, so generation can't recover it. This is distinct from
// CredentialSort.Value (see TestCredentialSortValueIsWhitespaceOnly), which
// closes the schema-level sort/#text finding but never holds this name.
func TestCredentialSortFieldNameIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantName   string
		wantOrder  string
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_credentials", "001_response.xml"), "name", "ascending"},
		{"live_002", filepath.Join(gmptestDir, "live", "get_credentials", "002_response.xml"), "name", "ascending"},
		{"live_003", filepath.Join(gmptestDir, "live", "get_credentials", "003_response.xml"), "name", "ascending"},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence credentialSortFieldEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into credentialSortFieldEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.Sort.Field.Value); got != tc.wantName {
				t.Fatalf("%s: sort/field bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantName)
			}
			if evidence.Sort.Field.Order != tc.wantOrder {
				t.Fatalf("%s: credentialSortFieldEvidence order = %q, want %q",
					tc.path, evidence.Sort.Field.Order, tc.wantOrder)
			}

			var got credentialSortFieldProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into credentialSortFieldProbe: %v", tc.path, err)
			}
			if got.Sort.Field.Order != tc.wantOrder {
				t.Errorf("credentialSortFieldProbe order = %q, want %q (the generated shape changed unexpectedly)",
					got.Sort.Field.Order, tc.wantOrder)
			}
			// credentialSortFieldProbe.Sort.Field has no field that could
			// have captured evidence.Sort.Field.Value — that silent drop,
			// now proven against the actual fixture bytes rather than
			// assumed, is the accepted loss this test exists to pin.
			if got := strings.TrimSpace(got.Sort.Value); got != "" {
				t.Errorf("credentialSortFieldProbe.Sort.Value = %q, want empty — CredentialSort.Value was expected to hold only indentation whitespace, not %q",
					got, got)
			}
		})
	}
}

// TestCredentialSortValueIsWhitespaceOnly proves CredentialSort.Value (the
// field that closes the schema-level sort/#text finding) only ever holds
// formatting whitespace, not the sort field name pinned by
// TestCredentialSortFieldNameIsAnAcceptedLoss.
func TestCredentialSortValueIsWhitespaceOnly(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	paths := []string{
		filepath.Join(gmptestDir, "live", "get_credentials", "001_response.xml"),
		filepath.Join(gmptestDir, "live", "get_credentials", "002_response.xml"),
		filepath.Join(gmptestDir, "live", "get_credentials", "003_response.xml"),
	}
	for _, path := range paths {
		t.Run(filepath.Base(filepath.Dir(path))+"/"+filepath.Base(path), func(t *testing.T) {
			data, err := os.ReadFile(path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", path, err)
			}

			var got credentialSortFieldProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into credentialSortFieldProbe: %v", path, err)
			}
			if v := strings.TrimSpace(got.Sort.Value); v != "" {
				t.Fatalf("%s: CredentialSort.Value = %q after trimming, want empty — fixture no longer demonstrates that this field captures only formatting whitespace, not the sort field name",
					path, v)
			}
		})
	}
}

// TestValidatePortListBeforeMatchesKnownConformance is the port_list family's
// counterpart to TestValidateCredentialBeforeMatchesKnownConformance.
// port_list is adopted, so this is a drift detector on the committed
// commands/ files. Pre-adoption findings are recorded in
// docs/review-2026-08.md.
func TestValidatePortListBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_port_list")

	want := map[string]counts{
		"create_port_list": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{},
		},
		"modify_port_list": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_port_lists": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{},
		},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidatePortListGeneratedMatchesOrImproves is the port_list family's
// post-swap gate, run at the same rigour every prior family's was. See
// docs/review-2026-08.md's port_list section for the per-finding
// explanation of what generation closed while port_list was still
// review_only.
func TestValidatePortListGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_port_list")

	wantAfter := map[string]counts{
		"create_port_list": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{},
		},
		"modify_port_list": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_port_lists": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{},
		},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateAssetBeforeMatchesKnownConformance is the asset family's
// counterpart to TestValidatePortListBeforeMatchesKnownConformance. asset is
// adopted, so this is a drift detector on the committed commands/ files.
// Pre-adoption findings are recorded in docs/review-2026-08.md.
func TestValidateAssetBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_asset")

	want := map[string]counts{
		"create_asset": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_asset": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_assets": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{},
		},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateAssetGeneratedMatchesOrImproves is the asset family's
// post-swap gate, run at the same rigour every prior family's was. See
// docs/review-2026-08.md's asset section for the per-finding explanation of
// what generation closed while asset was still review_only.
func TestValidateAssetGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_asset")

	wantAfter := map[string]counts{
		"create_asset": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_asset": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_assets": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{},
		},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateReportFormatBeforeMatchesKnownConformance is report_format's
// counterpart to TestValidateAssetBeforeMatchesKnownConformance. report_format
// is adopted, so this is a drift detector on the committed commands/ files.
// Pre-adoption findings are recorded in docs/review-2026-08.md.
func TestValidateReportFormatBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_report_format")

	want := map[string]counts{
		"create_report_format": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{},
		},
		"modify_report_format": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_report_formats": {
			actionable: map[string]int{},
			suppressed: map[string]int{schema.CatExtra: 1},
		},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateReportFormatGeneratedMatchesOrImproves is the report_format
// family's post-swap gate, run at the same rigour every prior family's was.
// Generation closed every actionable finding on get_report_formats and
// modify_report_format; create_report_format keeps its allowlisted
// unsupported-r SCHEMA-DEFECT, which no generated shape can resolve. See
// docs/review-2026-08.md's report_format section.
//
// get_report_formats' one remaining EXTRA comes back Suppressed, the same
// disposition tls_certificate's /tls_certificate/valid gets and for the same
// reason: a FieldOverride at "get_report_formats:/report_format" keeps a
// Comment field the schema's pattern never declares, because gvmd emits
// <comment> on every report format in all three committed fixtures. The
// finding still fires — the schema really does not declare it — but carries
// the evidence in its Reason instead of counting against commands/. See
// TestReportFormatCommentRoundTrips.
func TestValidateReportFormatGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_report_format")

	wantAfter := map[string]counts{
		"create_report_format": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{},
		},
		"modify_report_format": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_report_formats": {
			actionable: map[string]int{},
			suppressed: map[string]int{schema.CatExtra: 1},
		},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// reportFormatCommentProbe mirrors the generated ReportFormat's comment
// handling in testdata/golden/get_report_formats.go — keep the two in sync.
type reportFormatCommentProbe struct {
	XMLName      xml.Name                   `xml:"get_report_formats_response"`
	ReportFormat []reportFormatCommentEntry `xml:"report_format"`
}

type reportFormatCommentEntry struct {
	XMLName xml.Name `xml:"report_format"`
	Name    string   `xml:"name"`
	Comment string   `xml:"comment"`
}

// TestReportFormatCommentRoundTrips is the fixture-backed half of the
// FieldOverride's evidence: the override's Reason cites three captures, and
// this decodes all three to prove the claim still holds and that the
// generated shape retains what they carry.
//
// Presence is the whole point. Every captured <comment> is empty, so the
// assertion that matters is that the element survives a decode/marshal
// round trip rather than vanishing — which is what would happen under
// omitempty, and what happens today when the field is absent entirely.
func TestReportFormatCommentRoundTrips(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	fixtures := filepath.Join(dir, "..", "gmptest", "testdata")

	for _, tc := range []struct {
		path  string
		count int
	}{
		{filepath.Join(fixtures, "live", "get_report_formats", "001_response.xml"), 5},
		{filepath.Join(fixtures, "live", "get_report_formats", "002_response.xml"), 5},
		{filepath.Join(fixtures, "upstream", "get_report_formats", "002_response.xml"), 1},
	} {
		data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture paths under this package's own tree
		if err != nil {
			t.Fatalf("reading %s: %v", tc.path, err)
		}
		if got := strings.Count(string(data), "<comment"); got != tc.count {
			t.Fatalf("%s: %d <comment> elements on the wire, want %d — the FieldOverride's cited evidence has changed",
				tc.path, got, tc.count)
		}

		var got reportFormatCommentProbe
		if err := xml.Unmarshal(data, &got); err != nil {
			t.Fatalf("decoding %s: %v", tc.path, err)
		}
		if len(got.ReportFormat) != tc.count {
			t.Fatalf("%s: decoded %d report_format entries, want %d", tc.path, len(got.ReportFormat), tc.count)
		}

		out, err := xml.Marshal(got.ReportFormat[0])
		if err != nil {
			t.Fatalf("%s: marshaling a report_format back out: %v", tc.path, err)
		}
		if !strings.Contains(string(out), "<comment>") {
			t.Errorf("%s: round-tripped report_format dropped <comment>: %s", tc.path, out)
		}
	}
}

// portListCountProbe mirrors PortListCount in testdata/golden/get_port_lists.go.
type portListCountProbe struct {
	XMLName       xml.Name `xml:"get_port_lists_response"`
	PortListCount struct {
		Filtered int `xml:"filtered"`
		Page     int `xml:"page"`
	} `xml:"port_list_count"`
}

// portListCountEvidence is portListCountProbe plus a ",chardata" field to
// prove the bare total the generated shape can't retain.
type portListCountEvidence struct {
	XMLName       xml.Name `xml:"get_port_lists_response"`
	PortListCount struct {
		Value    string `xml:",chardata"`
		Filtered int    `xml:"filtered"`
		Page     int    `xml:"page"`
	} `xml:"port_list_count"`
}

// TestPortListCountBareTotalIsAnAcceptedLoss mirrors
// TestTicketCountBareTotalIsAnAcceptedLoss: live fixtures send a bare total
// ahead of <filtered>/<page>, but port_list_count's schema pattern declares
// no text token, so generation can't recover it. See docs/review-2026-08.md.
func TestPortListCountBareTotalIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantTotal  string
		wantFilter int
		wantPage   int
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_port_lists", "001_response.xml"), "3", 3, 3},
		{"live_002", filepath.Join(gmptestDir, "live", "get_port_lists", "002_response.xml"), "3", 3, 2},
		{"live_003", filepath.Join(gmptestDir, "live", "get_port_lists", "003_response.xml"), "4", 1, 1},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence portListCountEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into portListCountEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.PortListCount.Value); got != tc.wantTotal {
				t.Fatalf("%s: port_list_count bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantTotal)
			}
			if evidence.PortListCount.Filtered != tc.wantFilter || evidence.PortListCount.Page != tc.wantPage {
				t.Fatalf("%s: portListCountEvidence children = %+v, want Filtered=%d Page=%d",
					tc.path, evidence.PortListCount, tc.wantFilter, tc.wantPage)
			}

			var got portListCountProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into portListCountProbe: %v", tc.path, err)
			}
			if got.PortListCount.Filtered != tc.wantFilter || got.PortListCount.Page != tc.wantPage {
				t.Errorf("portListCountProbe = %+v, want Filtered=%d Page=%d (the generated shape changed unexpectedly)",
					got.PortListCount, tc.wantFilter, tc.wantPage)
			}
			// portListCountProbe has no field that could have captured
			// evidence.PortListCount.Value — that silent drop, now proven
			// against the actual fixture bytes rather than assumed, is the
			// accepted loss this test exists to pin.
		})
	}
}

// portListSortFieldProbe mirrors PortListSort/PortListField in
// testdata/golden/get_port_lists.go.
type portListSortFieldProbe struct {
	XMLName xml.Name `xml:"get_port_lists_response"`
	Sort    struct {
		Value string `xml:",chardata"`
		Field struct {
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// portListSortFieldEvidence is portListSortFieldProbe plus a ",chardata"
// field on Field to prove the bare sort field name the generated shape
// can't retain.
type portListSortFieldEvidence struct {
	XMLName xml.Name `xml:"get_port_lists_response"`
	Sort    struct {
		Field struct {
			Value string `xml:",chardata"`
			Order string `xml:"order"`
		} `xml:"field"`
	} `xml:"sort"`
}

// TestPortListSortFieldNameIsAnAcceptedLoss mirrors
// TestTicketSortFieldNameIsAnAcceptedLoss: fixtures send the sort field's
// name as bare text inside <field>, but <field>'s schema pattern declares
// only <order>, so generation can't recover it. This is distinct from
// PortListSort.Value (see TestPortListSortValueIsWhitespaceOnly), which
// closes the schema-level sort/#text finding but never holds this name. See
// docs/review-2026-08.md.
func TestPortListSortFieldNameIsAnAcceptedLoss(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	cases := []struct {
		provenance string
		path       string
		wantName   string
		wantOrder  string
	}{
		{"live_001", filepath.Join(gmptestDir, "live", "get_port_lists", "001_response.xml"), "name", "ascending"},
		{"live_002", filepath.Join(gmptestDir, "live", "get_port_lists", "002_response.xml"), "name", "ascending"},
		{"live_003", filepath.Join(gmptestDir, "live", "get_port_lists", "003_response.xml"), "name", "ascending"},
	}
	for _, tc := range cases {
		t.Run(tc.provenance, func(t *testing.T) {
			data, err := os.ReadFile(tc.path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", tc.path, err)
			}

			var evidence portListSortFieldEvidence
			if err := xml.Unmarshal(data, &evidence); err != nil {
				t.Fatalf("decoding %s into portListSortFieldEvidence: %v", tc.path, err)
			}
			if got := strings.TrimSpace(evidence.Sort.Field.Value); got != tc.wantName {
				t.Fatalf("%s: sort/field bare text = %q, want %q — fixture no longer demonstrates the documented contradiction",
					tc.path, got, tc.wantName)
			}
			if evidence.Sort.Field.Order != tc.wantOrder {
				t.Fatalf("%s: portListSortFieldEvidence order = %q, want %q",
					tc.path, evidence.Sort.Field.Order, tc.wantOrder)
			}

			var got portListSortFieldProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into portListSortFieldProbe: %v", tc.path, err)
			}
			if got.Sort.Field.Order != tc.wantOrder {
				t.Errorf("portListSortFieldProbe order = %q, want %q (the generated shape changed unexpectedly)",
					got.Sort.Field.Order, tc.wantOrder)
			}
			// portListSortFieldProbe.Sort.Field has no field that could have
			// captured evidence.Sort.Field.Value — that silent drop, now
			// proven against the actual fixture bytes rather than assumed,
			// is the accepted loss this test exists to pin.
			if got := strings.TrimSpace(got.Sort.Value); got != "" {
				t.Errorf("portListSortFieldProbe.Sort.Value = %q, want empty — PortListSort.Value was expected to hold only indentation whitespace, not %q",
					got, got)
			}
		})
	}
}

// TestPortListSortValueIsWhitespaceOnly proves PortListSort.Value (the field
// that closes the schema-level sort/#text finding) only ever holds
// formatting whitespace, not the sort field name pinned by
// TestPortListSortFieldNameIsAnAcceptedLoss.
func TestPortListSortValueIsWhitespaceOnly(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	gmptestDir := filepath.Join(dir, "..", "gmptest", "testdata")

	paths := []string{
		filepath.Join(gmptestDir, "live", "get_port_lists", "001_response.xml"),
		filepath.Join(gmptestDir, "live", "get_port_lists", "002_response.xml"),
		filepath.Join(gmptestDir, "live", "get_port_lists", "003_response.xml"),
	}
	for _, path := range paths {
		t.Run(filepath.Base(filepath.Dir(path))+"/"+filepath.Base(path), func(t *testing.T) {
			data, err := os.ReadFile(path) //nolint:gosec // fixed test-fixture path under this package's own tree
			if err != nil {
				t.Fatalf("reading %s: %v", path, err)
			}

			var got portListSortFieldProbe
			if err := xml.Unmarshal(data, &got); err != nil {
				t.Fatalf("decoding %s into portListSortFieldProbe: %v", path, err)
			}
			if v := strings.TrimSpace(got.Sort.Value); v != "" {
				t.Fatalf("%s: PortListSort.Value = %q after trimming, want empty — fixture no longer demonstrates that this field captures only formatting whitespace, not the sort field name",
					path, v)
			}
		})
	}
}

// TestValidateReportConfigBeforeMatchesKnownConformance is report_config's
// counterpart to TestValidateReportFormatBeforeMatchesKnownConformance.
// report_config is adopted, so this is a drift detector on the committed
// commands/ files. Pre-adoption findings are recorded in
// docs/review-2026-08.md.
func TestValidateReportConfigBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_report_config")

	want := map[string]counts{
		"create_report_config": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_report_config": {
			actionable: map[string]int{},
			suppressed: map[string]int{schema.CatExtra: 1},
		},
		"get_report_configs": {
			actionable: map[string]int{schema.CatDefect: 1, schema.CatMissing: 1},
			suppressed: map[string]int{schema.CatExtra: 1},
		},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateReportConfigGeneratedMatchesOrImproves is the report_config
// family's post-swap gate, run at the same rigour every prior family's was.
// Generation closes every actionable finding on create_report_config and
// modify_report_config, and on get_report_configs except the allowlisted
// empty-pattern SCHEMA-DEFECT (the schema itself declares
// report_config/report_format/permissions with no content) and the
// using_default relocation: FieldSuppressions removes
// param/value/@using_default, which schema.Comparer honestly reports as one
// MISSING finding — schema.CatMissing has no suppression mechanism, unlike
// CatExtra's gmpoverride tag — and FieldOverrides' matching addition on
// <param> is a real, schema-undeclared field, reported as one suppressed
// EXTRA via that same tag. Both counts are the fix surfacing a previously
// hidden truth (see TestReportConfigUsingDefaultOnParamAgreesWithFixture and
// TestReportConfigUsingDefaultDecodesOnGeneratedParam below), not a
// regression.
//
// modify_report_config gains its own suppressed EXTRA finding for the same
// reason: FieldOverrides["modify_report_config:/modify_report_config"] adds
// ReportConfigID, a live-gvmd-evidenced field the schema's <pattern> never
// declares at all (see the FieldOverride's Reason and the modify_report_config
// CommandWarning in overrides.go). schema.Comparer honestly reports the
// schema-undeclared attribute as EXTRA; the gmpoverride tag suppresses it.
func TestValidateReportConfigGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_report_config")

	wantAfter := map[string]counts{
		"create_report_config": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_report_config": {
			actionable: map[string]int{},
			suppressed: map[string]int{schema.CatExtra: 1},
		},
		"get_report_configs": {
			actionable: map[string]int{schema.CatDefect: 1, schema.CatMissing: 1},
			suppressed: map[string]int{schema.CatExtra: 1},
		},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// reportConfigUsingDefaultProbe decodes <param using_default="..."> exactly
// as command_get_report_configs.xml's own <example> and
// internal/gmptest/testdata/upstream/get_report_configs/001_response.xml
// send it — on <param>, not nested under <value> the way the schema's
// <pattern> (and hence testdata/golden/get_report_configs.go's
// ReportConfigParamValue) declares it.
type reportConfigUsingDefaultProbe struct {
	XMLName xml.Name `xml:"get_report_configs_response"`
	Report  struct {
		Param []struct {
			UsingDefault bool `xml:"using_default,attr"`
		} `xml:"param"`
	} `xml:"report_config"`
}

// TestReportConfigUsingDefaultOnParamAgreesWithFixture is the fixture-backed
// half of the get_report_configs_response CommandWarning's evidence: it
// proves the upstream fixture still places using_default on <param>, so a
// future schema or fixture update that moved it under <value> (matching the
// generated shape) would be caught here rather than leaving a stale warning
// on generated code that no longer needs it.
func TestReportConfigUsingDefaultOnParamAgreesWithFixture(t *testing.T) {
	dir, err := packageDir()
	if err != nil {
		t.Fatalf("packageDir: %v", err)
	}
	path := filepath.Join(dir, "..", "gmptest", "testdata", "upstream", "get_report_configs", "001_response.xml")

	data, err := os.ReadFile(path) //nolint:gosec // fixed test-fixture path under this package's own tree
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	if got := strings.Count(string(data), "using_default="); got != 2 {
		t.Fatalf("%s: %d using_default attributes on the wire, want 2 — the CommandWarning's cited evidence has changed", path, got)
	}

	var got reportConfigUsingDefaultProbe
	if err := xml.Unmarshal(data, &got); err != nil {
		t.Fatalf("decoding %s: %v", path, err)
	}
	if len(got.Report.Param) != 2 {
		t.Fatalf("%s: decoded %d param entries, want 2", path, len(got.Report.Param))
	}
	if got.Report.Param[0].UsingDefault {
		t.Errorf("%s: param[0]'s using_default decoded true reading it off <param>, want false", path)
	}
	if !got.Report.Param[1].UsingDefault {
		t.Errorf("%s: param[1]'s using_default did not decode to true reading it off <param> — the fixture no longer demonstrates the schema/wire mismatch", path)
	}
}

// TestReportConfigUsingDefaultDecodesOnGeneratedParam is
// TestReportConfigUsingDefaultOnParamAgreesWithFixture's counterpart against
// the actual generated type, not a private probe struct: it proves the
// FieldSuppressions/FieldOverrides relocation actually lands UsingDefault on
// the real GetReportConfigsParam this package emits, using the same
// scratch-package/`go run` machinery as TestGetFiltersCollisionRoundTrips.
func TestReportConfigUsingDefaultDecodesOnGeneratedParam(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("Skipping: no go toolchain on PATH (%v)", err)
	}
	got := generateOrFail(t)
	if _, ok := got[CommandFile("get_report_configs")]; !ok {
		t.Fatal("get_report_configs was not generated, so this test would not compile it")
	}

	fixture, err := filepath.Abs("../gmptest/testdata/upstream/get_report_configs/001_response.xml")
	if err != nil {
		t.Fatalf("Failed to resolve the fixture path: %v", err)
	}
	if _, err := os.Stat(fixture); err != nil {
		t.Fatalf("Fixture not found: %v", err)
	}

	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to resolve the module root: %v", err)
	}
	dir := t.TempDir()
	pkgDir := filepath.Join(dir, "genproto")
	if err := os.MkdirAll(pkgDir, 0o750); err != nil {
		t.Fatalf("Failed to create the scratch package: %v", err)
	}

	goVersion, err := goDirective(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatalf("Failed to read the module's go directive: %v", err)
	}
	gomod := "module gogmpcodegencheck\n\ngo " + goVersion + "\n\n" +
		"require github.com/brennoo/go-gmp v0.0.0\n\n" +
		"replace github.com/brennoo/go-gmp => " + root + "\n"
	mustWriteFile(t, filepath.Join(dir, "go.mod"), gomod)
	mustWriteFile(t, filepath.Join(pkgDir, "support.go"),
		strings.Replace(supportShim(), "package "+targetPackage+"\n", "package genproto\n", 1))
	for name, src := range got {
		rewritten := strings.Replace(src, "\npackage "+targetPackage+"\n", "\npackage genproto\n", 1)
		if rewritten == src {
			t.Fatalf("%s: expected a %q package clause to rewrite", name, targetPackage)
		}
		mustWriteFile(t, filepath.Join(pkgDir, name), rewritten)
	}
	mustWriteFile(t, filepath.Join(dir, "main.go"), reportConfigUsingDefaultProgram)

	cmd := exec.CommandContext(context.Background(), goBin, "run", ".", fixture)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Decode program failed: %v\n%s", err, out)
	}
	const want = "param[0].UsingDefault=false\nparam[1].UsingDefault=true\n"
	if got := string(out); got != want {
		t.Fatalf("Decode program printed:\n%s\nwant:\n%s", got, want)
	}
}

// reportConfigUsingDefaultProgram is a `go run` main for
// TestReportConfigUsingDefaultDecodesOnGeneratedParam. It decodes the fixture
// named on argv[1] through the actual generated GetReportConfigsResponse and
// prints each <param>'s UsingDefault, read off the generated
// GetReportConfigsParam field the relocation adds.
const reportConfigUsingDefaultProgram = `package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"gogmpcodegencheck/genproto"
)

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var decoded genproto.GetReportConfigsResponse
	if err := xml.Unmarshal(data, &decoded); err != nil {
		fmt.Fprintln(os.Stderr, "decode:", err)
		os.Exit(1)
	}
	if len(decoded.ReportConfig) != 1 || len(decoded.ReportConfig[0].Param) != 2 {
		fmt.Fprintf(os.Stderr, "unexpected shape: %d report_config(s), want 1\n", len(decoded.ReportConfig))
		os.Exit(1)
	}
	for i, p := range decoded.ReportConfig[0].Param {
		fmt.Printf("param[%d].UsingDefault=%v\n", i, bool(p.UsingDefault))
	}
}
`

// TestValidateUserBeforeMatchesKnownConformance is user's counterpart to
// TestValidateReportConfigBeforeMatchesKnownConformance. user is adopted, so
// this is a drift detector on the committed commands/ files. Pre-adoption
// findings are recorded in docs/review-2026-08.md.
func TestValidateUserBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_user")

	want := map[string]counts{
		"create_user": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_user": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_users": {
			actionable: map[string]int{schema.CatDefect: 2},
			suppressed: map[string]int{},
		},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateUserGeneratedMatchesOrImproves is user's pre-adoption gate, run
// at the same rigour every prior family's was. Generation closes every
// actionable finding on create_user and modify_user (role/name and
// groups/group/name have no schema counterpart; sources/source narrows from
// a hand-written slice to the schema's declared scalar; hosts/@allow gains
// omitempty) and every actionable finding on get_users except the two
// allowlisted empty-pattern SCHEMA-DEFECT findings on role/permissions and
// groups/group/permissions, which the schema itself declares with no
// content. No FieldOverride/FieldSuppression/PointerOverride/CommandWarning
// was needed for this family — unlike report_config, no fixture or
// live-gvmd evidence surfaced a schema/wire mismatch here.
func TestValidateUserGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_user")

	wantAfter := map[string]counts{
		"create_user": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_user": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_users": {
			actionable: map[string]int{schema.CatDefect: 2},
			suppressed: map[string]int{},
		},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// assertMatchesOrImproves is the shared body of the two match-or-improve gates:
// the exact expected "after" counts, plus the two general rules that hold for
// any family — no category may grow, and the schema's own defects may not move.
func assertMatchesOrImproves(t *testing.T, before, after []schema.Result, wantAfter map[string]counts) {
	t.Helper()
	beforeByCmd := map[string]schema.Result{}
	for _, res := range before {
		beforeByCmd[res.Command] = res
	}

	for _, res := range after {
		w, ok := wantAfter[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's after results", res.Command)
			continue
		}
		got := countFindings(res.Findings)
		assertCounts(t, "after", res.Command, w, got)

		// The general, snapshot-independent rule: nothing regresses. Every
		// category present after generation must already have been present
		// before it (no brand-new problem class), and no category's total
		// count (actionable plus suppressed) may increase. The two are
		// combined, rather than checked independently, because a finding can
		// legitimately move from actionable to suppressed within a category
		// (e.g. a FieldOverride turning a bare EXTRA into a suppressed one)
		// without either side alone reflecting "no regression".
		b := countFindings(beforeByCmd[res.Command].Findings)
		assertNoRegression(t, res.Command, "actionable+suppressed",
			addCounts(b.actionable, b.suppressed), addCounts(got.actionable, got.suppressed))

		// SCHEMA-DEFECT findings are copied from cmd.Defects verbatim,
		// independent of req/resp — they must be byte-identical before and
		// after, since nothing about regenerating commands/ can fix an
		// upstream schema bug.
		assertSameDefects(t, res.Command, beforeByCmd[res.Command].Findings, res.Findings)
	}
	if len(after) != len(wantAfter) {
		t.Errorf("Validate returned %d after-results, want %d", len(after), len(wantAfter))
	}
}

func assertCounts(t *testing.T, phase, cmd string, want, got counts) {
	t.Helper()
	for _, cat := range schema.Categories {
		if want.actionable[cat] != got.actionable[cat] {
			t.Errorf("%s %s: actionable %s = %d, want %d", phase, cmd, cat, got.actionable[cat], want.actionable[cat])
		}
		if want.suppressed[cat] != got.suppressed[cat] {
			t.Errorf("%s %s: suppressed %s = %d, want %d", phase, cmd, cat, got.suppressed[cat], want.suppressed[cat])
		}
	}
}

// addCounts sums two per-category count maps, used to compare actionable and
// suppressed totals combined rather than as two independent ledgers.
func addCounts(a, b map[string]int) map[string]int {
	out := make(map[string]int, len(a)+len(b))
	for c, n := range a {
		out[c] += n
	}
	for c, n := range b {
		out[c] += n
	}
	return out
}

func assertNoRegression(t *testing.T, cmd, phase string, before, after map[string]int) {
	t.Helper()
	cats := map[string]bool{}
	for c := range before {
		cats[c] = true
	}
	for c := range after {
		cats[c] = true
	}
	names := make([]string, 0, len(cats))
	for c := range cats {
		names = append(names, c)
	}
	sort.Strings(names)
	for _, c := range names {
		if after[c] > before[c] {
			t.Errorf("%s: %s %s findings went from %d (before) to %d (after) — unexplained regression",
				cmd, phase, c, before[c], after[c])
		}
	}
}

func assertSameDefects(t *testing.T, cmd string, before, after []schema.Finding) {
	t.Helper()
	b := defectDetails(before)
	a := defectDetails(after)
	if fmt.Sprint(b) != fmt.Sprint(a) {
		t.Errorf("%s: SCHEMA-DEFECT findings changed across generation, which should be impossible\nbefore: %v\nafter:  %v",
			cmd, b, a)
	}
}

func defectDetails(fs []schema.Finding) []string {
	var out []string
	for _, f := range fs {
		if f.Category == schema.CatDefect {
			out = append(out, f.Path+": "+f.Detail)
		}
	}
	sort.Strings(out)
	return out
}

// TestValidateTargetBeforeMatchesKnownConformance is target's counterpart to
// TestValidateUserBeforeMatchesKnownConformance. target is adopted, so this
// is a drift detector on the committed commands/ files. Pre-adoption
// findings are recorded in docs/review-2026-08.md.
func TestValidateTargetBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_target")

	want := map[string]counts{
		"create_target": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_target": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_targets": {
			actionable: map[string]int{schema.CatDefect: 9},
			suppressed: map[string]int{},
		},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateTargetGeneratedMatchesOrImproves is target's pre-adoption gate,
// run at the same rigour every prior family's was. Generation closes every
// actionable finding on create_target and modify_target, and every
// actionable finding on get_targets except the nine allowlisted
// SCHEMA-DEFECT findings the schema itself carries (see
// internal/codegen/testdata/schema/get_targets.json's Defects array and the
// matching DefectAllowlist entries in overrides.go). No
// FieldOverride/FieldSuppression/TypeOverride was needed for this family —
// every schema-vs-hand-written divergence found during review traced
// directly to Card/Kind values already in the schema snapshot. The one
// PointerOverride this family does carry (modify_target.comment) and the
// modify_target CommandWarning (Name/Hosts/ExcludeHosts/HostsOrdering
// presence hazards) came from live gvmd evidence, not from a
// generator-vs-schema mismatch, so neither shows up as a suppressed finding
// here.
func TestValidateTargetGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_target")

	wantAfter := map[string]counts{
		"create_target": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_target": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_targets": {
			actionable: map[string]int{schema.CatDefect: 9},
			suppressed: map[string]int{},
		},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateTaskBeforeMatchesKnownConformance is task's counterpart to
// TestValidateTargetBeforeMatchesKnownConformance. task is review-only, so
// this pins the pre-swap gate against the still-hand-written
// commands/{create_task,get_tasks,modify_task}.go. All three commands share
// one hand-written set of nested types (TaskConfig, TaskTarget, TaskScanner,
// TaskAlert, TaskSchedule, TaskObservers, TaskPreferences), each shaped for
// get_tasks' rich response; reused on create_task/modify_task's request side
// they carry Name/Trash/Type fields the schema never declares there, which is
// most of both commands' EXTRA counts. modify_task is also missing the
// schema's <file> element entirely. get_tasks itself is the heaviest: bool
// where the wire sends a GMP tri-state, SchedulePeriods modeled as a
// sub-element container instead of the schema's bare integer, ResultCount as
// int instead of the schema's text, Scanner.Type as string instead of the
// schema's integer, ReportCount/UserTags/CurrentReport/LastReport's
// severity-and-compliance breakdowns missing or under-detailed, and an
// invented Reports field the schema never declares at all — see
// docs/review-2026-08.md's `task` review-only section for the full
// field-by-field list. The ten SCHEMA-DEFECT findings are the schema's own
// get_tasks defects (six duplicate-ref, four empty-pattern), unrelated to the
// hand-written code and unfixable by generation.
func TestValidateTaskBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_task")

	want := map[string]counts{
		"create_task": {
			actionable: map[string]int{schema.CatCardinality: 1, schema.CatExtra: 15, schema.CatMissing: 1},
			suppressed: map[string]int{schema.CatOmitEmpty: 4},
		},
		"modify_task": {
			actionable: map[string]int{schema.CatExtra: 12, schema.CatMissing: 1},
			suppressed: map[string]int{},
		},
		"get_tasks": {
			actionable: map[string]int{
				schema.CatCardinality: 7,
				schema.CatExtra:       14,
				schema.CatMissing:     9,
				schema.CatDefect:      10,
				schema.CatType:        8,
			},
			suppressed: map[string]int{schema.CatOmitEmpty: 7},
		},
	}
	for _, res := range before {
		w, ok := want[res.Command]
		if !ok {
			t.Errorf("%s: unexpected command in Validate's before results", res.Command)
			continue
		}
		assertCounts(t, "before", res.Command, w, countFindings(res.Findings))
	}
}

// TestValidateTaskGeneratedMatchesOrImproves is task's pre-adoption gate, run
// at the same rigour every prior family's was. Generation closes every
// actionable finding on create_task and modify_task, and every actionable
// finding on get_tasks except the ten allowlisted SCHEMA-DEFECT findings the
// schema itself carries (see
// internal/codegen/testdata/schema/get_tasks.json's Defects array and the
// matching DefectAllowlist entries in overrides.go). The two PointerOverrides
// this family carries (modify_task.comment, modify_task.observers) came from
// live gvmd evidence, not from a generator-vs-schema mismatch, so neither
// shows up as a suppressed finding here.
func TestValidateTaskGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_task")

	wantAfter := map[string]counts{
		"create_task": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_task": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_tasks": {
			actionable: map[string]int{schema.CatDefect: 10},
			suppressed: map[string]int{},
		},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}
