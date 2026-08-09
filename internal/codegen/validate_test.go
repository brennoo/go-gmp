package codegen

import (
	"errors"
	"fmt"
	"sort"
	"testing"

	"github.com/brennoo/go-gmp/internal/schema"
)

// validateOrSkip runs Validate for the committed snapshot family, skipping
// (not failing) when no "go" toolchain is available to build and run the
// reflect helper GeneratedGoRoots depends on — the same accommodation
// TestGeneratedCodeCompiles makes.
func validateOrSkip(t *testing.T) (before, after []schema.Result) {
	t.Helper()
	cmds, err := LoadSnapshots(snapshotDir)
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
	before, _ := validateOrSkip(t)

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
	before, after := validateOrSkip(t)

	wantAfter := map[string]counts{
		"create_note": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_note": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_notes":   {actionable: map[string]int{schema.CatDefect: 2}, suppressed: map[string]int{}},
	}

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
