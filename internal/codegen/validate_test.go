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
	before, after, err = Validate(cmds, TypeOverrides, nil)
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
