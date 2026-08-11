package codegen

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
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
	before, after, err = Validate(cmds, TypeOverrides, PointerOverrides, nil)
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
// review_only, so commands/{create_filter,modify_filter,get_filters}.go are
// the hand-written originals — this pins the pre-generation gate the family
// was selected to exercise. create_filter and modify_filter are already
// conformant. get_filters carries the family's entire actionable/suppressed
// load, plus its own SCHEMA-DEFECT 1 (see
// TestValidateFilterGeneratedMatchesOrImproves for what each finding is).
//
// MISSING/EXTRA are 1/1, not the 2/5 an occurrence-blind comparer once
// reported: internal/schema/conform.go's walk now pairs the schema's two
// same-named "filters" children against Go fields by occurrence instead of
// both resolving to the same node (see
// TestValidateFilterGeneratedMatchesOrImproves), which was inflating these
// counts even on the hand-written side.
func TestValidateFilterBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_filter")

	want := map[string]counts{
		"create_filter": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_filter": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_filters": {
			actionable: map[string]int{
				schema.CatType:    2,
				schema.CatDefect:  1,
				schema.CatMissing: 1,
				schema.CatExtra:   1,
			},
			suppressed: map[string]int{schema.CatOmitEmpty: 1},
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
// pre-swap gate, run at the same rigour note through alert were.
//
// Generation closes the TYPE 2 findings (Filter.InUse/.Writable become
// commands.GMPBool), the suppressed OMITEMPTY finding (Filters.Name gains
// omitempty), and the MISSING/EXTRA findings the hand-written FilterInfo2
// field produced by being wrongly tagged xml:"filters2" instead of the
// wire's actual "filters". SCHEMA-DEFECT 1 is untouched (the generic
// empty-pattern class already on record for a dozen other commands, at
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
		// before it (no brand-new problem class), and no category's count may
		// increase (no problem the hand-written code didn't already have,
		// getting worse).
		b := countFindings(beforeByCmd[res.Command].Findings)
		assertNoRegression(t, res.Command, "actionable", b.actionable, got.actionable)
		assertNoRegression(t, res.Command, "suppressed", b.suppressed, got.suppressed)

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
