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
// skipping (not failing) when no "go" toolchain is available — the same
// accommodation TestGeneratedCodeCompiles makes. It validates one family at a
// time because GeneratedGoRoots's Generate call is scoped to a single family
// by GenerateAll's contract, keeping "after" identical to that family's golden
// bytes.
func validateOrSkip(t *testing.T, command string) (before, after []schema.Result) {
	t.Helper()
	fam, ok := manifestOrFail(t).FamilyOf(command)
	if !ok {
		t.Fatalf("The manifest no longer declares %s", command)
	}
	cmds, err := LoadGenerationModels(snapshotDir, fam.Commands)
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}
	before, after, err = Validate(cmds, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist, NameOverrides, FieldTypeOverrides, nil, fam.Name, CardinalityOverrides, EnumConstantNames)
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

// TestValidateBeforeMatchesKnownConformance pins what Validate's "before" side
// reports for the note family, cross-checked against internal/schema/testdata/
// CONFORMANCE.md's committed create_note/modify_note/get_notes sections. It
// exists to catch a bug in BeforeGoRoots or snapshot loading — not the
// generator — before it gets blamed on the generated side by
// TestValidateGeneratedMatchesOrImproves below.
//
// commands/{create_note,modify_note,get_notes,note_shared}.go are now the
// generator's own output, so "before" and "after" describe the same source and
// these numbers equal wantAfter's: this is a drift detector on the committed
// files, not a weaker assertion. Pre-swap numbers and per-finding history are
// in docs/review-2026-08.md's T-9 section.
func TestValidateBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_note")

	want := map[string]counts{
		"create_note": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_note": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_notes": {
			actionable: map[string]int{schema.CatDefect: 3},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 2},
		},
		"delete_note": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
// Generation closes every actionable and suppressed finding for all three
// commands except get_notes's three pre-existing SCHEMA-DEFECT findings
// ("duplicate-ele: redefined in the same scope"), untouched byte-for-byte:
// those come from cmd.Defects, the schema's own upstream bugs, copied into the
// result before Comparer ever looks at req/resp, so no regeneration can move
// them. See docs/review-2026-08.md's T-9 note section for the per-finding
// breakdown of what closed and why.
//
// note_count's bare leading total is recovered via overrides.go's
// FieldOverrides["get_notes:/note_count"] (T-23), adding one suppressed
// EXTRA below.
func TestValidateGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_note")

	wantAfter := map[string]counts{
		"create_note": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_note": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_notes":   {actionable: map[string]int{schema.CatDefect: 3}, suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 2}},
		"delete_note": {actionable: map[string]int{}, suppressed: map[string]int{}},
	}

	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateOverrideBeforeMatchesKnownConformance is the override family's
// counterpart to TestValidateBeforeMatchesKnownConformance: commands/
// {create_override,modify_override,get_overrides,override_shared}.go are now
// the generator's own output, so this is a drift detector on the committed
// files. Pre-swap numbers and the per-finding closure history are in
// docs/review-2026-08.md's T-9 override section.
func TestValidateOverrideBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_override")

	want := map[string]counts{
		"create_override": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_override": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_overrides": {
			actionable: map[string]int{schema.CatDefect: 3},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 2},
		},
		"delete_override": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
// pre-swap gate, run at the same rigour the note family's was. Generation
// closed all 84 findings except get_overrides's three SCHEMA-DEFECT ones
// (upstream duplicate-ele bugs, byte-identical before and after — no
// generated output can move a cmd.Defects finding, same as get_notes). See
// docs/review-2026-08.md's T-9 override section for the per-finding breakdown.
//
// override_count's bare leading total is recovered via overrides.go's
// FieldOverrides["get_overrides:/override_count"] (T-23), adding one
// suppressed EXTRA below.
func TestValidateOverrideGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_override")

	wantAfter := map[string]counts{
		"create_override": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_override": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_overrides":   {actionable: map[string]int{schema.CatDefect: 3}, suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 2}},
		"delete_override": {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateRoleBeforeMatchesKnownConformance is the role family's
// counterpart to TestValidateBeforeMatchesKnownConformance.
// commands/{create_role,modify_role,get_roles}.go are the generator's own
// output, so this is a drift detector on the committed files, not a weaker
// assertion. Unlike note and override, get_roles carries no SCHEMA-DEFECT
// findings — role's schema has no known upstream defect.
func TestValidateRoleBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_role")

	want := map[string]counts{
		"create_role": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_role": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_roles":   {actionable: map[string]int{}, suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1}},
		"delete_role": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
// gate. Generation closed every actionable and suppressed finding on all
// three commands; role carries no SCHEMA-DEFECT findings to begin with.
//
// role_count's bare leading total is recovered via overrides.go's
// FieldOverrides["get_roles:/role_count"] (T-23), adding one suppressed
// EXTRA below.
func TestValidateRoleGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_role")

	wantAfter := map[string]counts{
		"create_role": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_role": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_roles":   {actionable: map[string]int{}, suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1}},
		"delete_role": {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateGroupBeforeMatchesKnownConformance is the group family's
// counterpart to TestValidateBeforeMatchesKnownConformance.
// commands/{create_group,modify_group,get_groups}.go are the generator's own
// output, so this is a drift detector on the committed files. create_group's
// actionable SCHEMA-DEFECT is copied from cmd.Defects and can't be moved by
// regeneration.
func TestValidateGroupBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_group")

	want := map[string]counts{
		"create_group": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{},
		},
		"modify_group": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_groups":   {actionable: map[string]int{}, suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1}},
		"delete_group": {actionable: map[string]int{}, suppressed: map[string]int{}},
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

// TestValidateGroupGeneratedMatchesOrImproves is the group family's
// post-swap gate. Generation closed every finding except create_group's
// SCHEMA-DEFECT (copied from cmd.Defects, untouched by regeneration).
//
// group_count's bare leading total is recovered via overrides.go's
// FieldOverrides["get_groups:/group_count"] (T-23), adding one suppressed
// EXTRA below. The sort field-name MISSING/EXTRA pair on get_groups is
// T-18a's relocation of GroupSort's chardata onto GroupField — see
// FieldSuppressions["get_groups:/sort"].
func TestValidateGroupGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_group")

	wantAfter := map[string]counts{
		"create_group": {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{}},
		"modify_group": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_groups":   {actionable: map[string]int{}, suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1}},
		"delete_group": {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateTagBeforeMatchesKnownConformance is the tag family's
// counterpart to TestValidateBeforeMatchesKnownConformance.
// commands/{create_tag,modify_tag,get_tags}.go are the generator's own
// output, so this is a drift detector on the committed files. get_tags's
// actionable SCHEMA-DEFECT is copied from cmd.Defects, untouched by
// regeneration.
func TestValidateTagBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_tag")

	want := map[string]counts{
		"create_tag": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_tags":   {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1}},
		"modify_tag": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"delete_tag": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
// gate. Generation closed every actionable and suppressed finding on all
// three commands except get_tags's one pre-existing SCHEMA-DEFECT
// (resources/resource/permissions, an empty-pattern site), copied from
// cmd.Defects and untouched by regeneration.
func TestValidateTagGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_tag")

	wantAfter := map[string]counts{
		"create_tag": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_tags":   {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1}},
		"modify_tag": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"delete_tag": {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
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
		"get_schedules":   {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1}},
		"delete_schedule": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
// post-swap gate. Generation closed every finding on get_schedules except
// its SCHEMA-DEFECT (copied from cmd.Defects, untouched by regeneration).
//
// schedule_count's bare leading total is recovered via overrides.go's
// FieldOverrides["get_schedules:/schedule_count"] (T-23), adding one
// suppressed EXTRA below. The sort field-name MISSING/EXTRA pair is T-18a's
// relocation of ScheduleSort's chardata onto ScheduleField — see
// FieldSuppressions["get_schedules:/sort"].
func TestValidateScheduleGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_schedule")

	wantAfter := map[string]counts{
		"create_schedule": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_schedule": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_schedules":   {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1}},
		"delete_schedule": {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
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
		"get_alerts":   {actionable: map[string]int{schema.CatDefect: 2}, suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1}},
		"delete_alert": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"test_alert":   {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{}},
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

// TestValidateAlertGeneratedMatchesOrImproves is the alert family's
// post-swap gate. Generation closed every finding except get_alerts's two
// SCHEMA-DEFECT findings (copied from cmd.Defects, untouched by
// regeneration).
//
// alert_count's bare leading total is recovered via overrides.go's
// FieldOverrides["get_alerts:/alert_count"] (T-23), adding one suppressed
// EXTRA below. The sort field-name MISSING/EXTRA pair is T-18a's relocation
// of AlertSort's chardata onto AlertField — see
// FieldSuppressions["get_alerts:/sort"].
func TestValidateAlertGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_alert")

	wantAfter := map[string]counts{
		"create_alert": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_alert": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_alerts":   {actionable: map[string]int{schema.CatDefect: 2}, suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1}},
		"delete_alert": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"test_alert":   {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateFilterBeforeMatchesKnownConformance is the filter family's
// counterpart to TestValidateBeforeMatchesKnownConformance. filter is
// adopted, so commands/{create_filter,modify_filter,get_filters}.go are
// already the generator's output — this is a drift detector on the committed
// files. Pre-adoption findings are recorded in docs/review-2026-08.md.
func TestValidateFilterBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_filter")

	want := map[string]counts{
		"create_filter": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_filter": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_filters": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_filter": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
// post-swap gate, run at the same rigour every prior family's was. Generation
// closed every finding except get_filters's SCHEMA-DEFECT (copied from
// cmd.Defects, untouched by regeneration). See docs/review-2026-08.md's
// filter section for the per-finding explanation, including the
// xml:"-"/gmpwire:"filters" resolution of the <filters> field collision (see
// get_filters.go's generated UnmarshalXML/MarshalXML and
// TestGetFiltersCollisionRoundTrips).
//
// filter_count's bare leading total, once pinned as a permanent loss by
// TestFilterCountBareTotalIsAnAcceptedLoss, is now recovered via
// overrides.go's FieldOverrides["get_filters:/filter_count"] (T-23); that
// adds one more suppressed EXTRA below.
func TestValidateFilterGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_filter")

	wantAfter := map[string]counts{
		"create_filter": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_filter": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_filters": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_filter": {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateTLSCertificateBeforeMatchesKnownConformance is the
// tls_certificate family's counterpart to
// TestValidateBeforeMatchesKnownConformance. tls_certificate is adopted, so
// commands/{create_tls_certificate,modify_tls_certificate,
// get_tls_certificates}.go are already the generator's output — this is a
// drift detector on the committed files. Pre-adoption findings are recorded
// in docs/review-2026-08.md.
func TestValidateTLSCertificateBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_tls_certificate")

	want := map[string]counts{
		"create_tls_certificate": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_tls_certificate": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_tls_certificates":   {actionable: map[string]int{}, suppressed: map[string]int{schema.CatExtra: 4, schema.CatMissing: 1, schema.CatOmitEmpty: 1}},
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
// See docs/review-2026-08.md's tls_certificate section for the per-finding
// breakdown of what generation closed while it was still review_only.
//
// Two EXTRA findings survive, both suppressed via FieldOverrides on
// live/upstream evidence rather than a schema-declared shape:
// /tls_certificate/certificate/#text keeps a Value field via the override at
// "get_tls_certificates:/tls_certificate/certificate" — the schema declares
// no chardata there, but the upstream bundled example sends real base64
// certificate bytes as that element's text (see
// TestTLSCertificateDataRoundTrips). /tls_certificate/valid keeps a Valid
// GMPBool field via the override at "get_tls_certificates:/tls_certificate"
// — the schema never references it at all, but live captures send
// <valid>1</valid> regardless.
//
// tls_certificate_count's bare leading total, once pinned as a permanent
// loss by TestTLSCertificateCountBareTotalIsAnAcceptedLoss, is now recovered
// via overrides.go's FieldOverrides["get_tls_certificates:/tls_certificate_count"]
// (T-23); that adds one more suppressed EXTRA below.
func TestValidateTLSCertificateGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_tls_certificate")

	wantAfter := map[string]counts{
		"create_tls_certificate": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_tls_certificate": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_tls_certificates":   {actionable: map[string]int{}, suppressed: map[string]int{schema.CatExtra: 4, schema.CatMissing: 1, schema.CatOmitEmpty: 1}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
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
// detector on the committed files. Pre-adoption findings are recorded in
// docs/review-2026-08.md.
func TestValidatePermissionBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_permission")

	want := map[string]counts{
		"create_permission": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_permission": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_permissions": {
			actionable: map[string]int{schema.CatDefect: 2},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_permission": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
// post-swap gate, run at the same rigour every prior family's was. Generation
// closed every finding except get_permissions's two SCHEMA-DEFECT findings on
// resource/permissions and subject/permissions (copied from cmd.Defects,
// untouched by regeneration). See docs/review-2026-08.md's permission
// section for the per-finding breakdown.
//
// permission_count's bare leading total, once pinned as a permanent loss by
// TestPermissionCountBareTotalIsAnAcceptedLoss, is now recovered via
// overrides.go's FieldOverrides["get_permissions:/permission_count"] (T-23);
// that adds one more suppressed EXTRA below.
func TestValidatePermissionGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_permission")

	wantAfter := map[string]counts{
		"create_permission": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_permission": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_permissions": {
			actionable: map[string]int{schema.CatDefect: 2},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_permission": {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateTicketBeforeMatchesKnownConformance is the ticket family's
// counterpart to TestValidatePermissionBeforeMatchesKnownConformance. ticket
// is adopted, so commands/{create_ticket,get_tickets,modify_ticket}.go are
// already the generator's output — this is a drift detector on the
// committed files. Pre-adoption findings are recorded in
// docs/review-2026-08.md.
func TestValidateTicketBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_ticket")

	want := map[string]counts{
		"create_ticket": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_ticket": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_tickets":   {actionable: map[string]int{}, suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1}},
		"delete_ticket": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
//
// ticket_count's bare leading total, once pinned as a permanent loss by
// TestTicketCountBareTotalIsAnAcceptedLoss, is now recovered via
// overrides.go's FieldOverrides["get_tickets:/ticket_count"] (T-23); that
// adds one more suppressed EXTRA below.
func TestValidateTicketGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_ticket")

	wantAfter := map[string]counts{
		"create_ticket": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_ticket": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_tickets":   {actionable: map[string]int{}, suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1}},
		"delete_ticket": {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateCredentialBeforeMatchesKnownConformance is the credential
// family's counterpart to TestValidatePermissionBeforeMatchesKnownConformance.
// credential is adopted, so commands/{create_credential,get_credentials,
// modify_credential}.go are already the generator's output — this is a drift
// detector on the committed files. Pre-adoption findings are recorded in
// docs/review-2026-08.md.
func TestValidateCredentialBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_credential")

	want := map[string]counts{
		"create_credential": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_credential": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_credentials": {
			actionable: map[string]int{schema.CatDefect: 2},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_credential": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
// post-swap gate, run at the same rigour every prior family's was. Generation
// closed every finding except get_credentials's two SCHEMA-DEFECT findings on
// scanner/permissions and target/permissions (copied from cmd.Defects,
// untouched by regeneration).
//
// credential_count's bare leading total, once pinned as a permanent loss by
// TestCredentialCountBareTotalIsAnAcceptedLoss, is now recovered via
// overrides.go's FieldOverrides["get_credentials:/credential_count"] (T-23);
// that adds one more suppressed EXTRA below.
func TestValidateCredentialGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_credential")

	wantAfter := map[string]counts{
		"create_credential": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_credential": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_credentials": {
			actionable: map[string]int{schema.CatDefect: 2},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_credential": {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
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
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_port_list": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
//
// port_list_count's bare leading total, once pinned as a permanent loss by
// TestPortListCountBareTotalIsAnAcceptedLoss, is now recovered via
// overrides.go's FieldOverrides["get_port_lists:/port_list_count"] (T-23);
// that adds one more suppressed EXTRA below.
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
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_port_list": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_asset": {actionable: map[string]int{}, suppressed: map[string]int{schema.CatOmitEmpty: 2}},
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
//
// asset_count's bare leading total is now recovered via overrides.go's
// FieldOverrides["get_assets:/asset_count"] (T-23), adding one suppressed
// EXTRA below.
func TestValidateAssetGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_asset")

	wantAfter := map[string]counts{
		"create_asset": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_asset": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_assets": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_asset": {actionable: map[string]int{}, suppressed: map[string]int{schema.CatOmitEmpty: 2}},
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
			suppressed: map[string]int{schema.CatExtra: 3, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_report_format": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"verify_report_format": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
// report_format_count's bare leading total is now recovered via
// overrides.go's FieldOverrides["get_report_formats:/report_format_count"]
// (T-23), adding one more suppressed EXTRA below.
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
			suppressed: map[string]int{schema.CatExtra: 3, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_report_format": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"verify_report_format": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
		{filepath.Join(fixtures, "live", "get_report_formats", "006_response.xml"), 6},
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
			suppressed: map[string]int{schema.CatExtra: 3, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_report_config": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
// Generation closes every finding on create_report_config and
// modify_report_config except get_report_configs's allowlisted
// empty-pattern SCHEMA-DEFECT and one deliberate MISSING: FieldSuppressions
// relocates param/value/@using_default onto <param> via a matching
// FieldOverride (see TestReportConfigUsingDefaultOnParamAgreesWithFixture and
// TestReportConfigUsingDefaultDecodesOnGeneratedParam below); schema.CatMissing
// has no suppression mechanism, so the relocation surfaces as one actionable
// MISSING and one suppressed EXTRA rather than a clean close.
//
// modify_report_config's suppressed EXTRA is
// FieldOverrides["modify_report_config:/modify_report_config"] adding
// ReportConfigID, a live-gvmd-evidenced field the schema never declares (see
// the override's Reason and its CommandWarning in overrides.go).
//
// report_config_count's bare leading total is now recovered via
// overrides.go's FieldOverrides["get_report_configs:/report_config_count"]
// (T-23), adding one more suppressed EXTRA below.
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
			suppressed: map[string]int{schema.CatExtra: 3, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_report_config": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_user": {actionable: map[string]int{}, suppressed: map[string]int{schema.CatOmitEmpty: 4}},
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
// finding except get_users's two allowlisted empty-pattern SCHEMA-DEFECT
// findings on role/permissions and groups/group/permissions.
// user_count's bare leading total is now recovered via overrides.go's
// FieldOverrides["get_users:/user_count"] (T-23), adding one more suppressed
// EXTRA below.
func TestValidateUserGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_user")

	wantAfter := map[string]counts{
		"create_user": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_user": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_users": {
			actionable: map[string]int{schema.CatDefect: 2},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_user": {actionable: map[string]int{}, suppressed: map[string]int{schema.CatOmitEmpty: 4}},
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
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_target": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
// finding on get_targets except the nine allowlisted SCHEMA-DEFECT findings
// the schema itself carries (see
// internal/codegen/testdata/schema/get_targets.json's Defects array and the
// matching DefectAllowlist entries in overrides.go).
// target_count's bare leading total is now recovered via overrides.go's
// FieldOverrides["get_targets:/target_count"] (T-23), adding one more
// suppressed EXTRA below.
func TestValidateTargetGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_target")

	wantAfter := map[string]counts{
		"create_target": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_target": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_targets": {
			actionable: map[string]int{schema.CatDefect: 9},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_target": {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateTaskBeforeMatchesKnownConformance is task's counterpart to
// TestValidateTargetBeforeMatchesKnownConformance. task is adopted, so this
// is a drift detector on the committed commands/ files. Pre-adoption
// findings are recorded in docs/review-2026-08.md.
func TestValidateTaskBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_task")

	want := map[string]counts{
		"create_task": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_task": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_tasks": {
			actionable: map[string]int{schema.CatDefect: 7},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_task": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"move_task":   {actionable: map[string]int{}, suppressed: map[string]int{}},
		"resume_task": {
			actionable: map[string]int{},
			suppressed: map[string]int{schema.CatOmitEmpty: 1},
		},
		"start_task": {
			actionable: map[string]int{},
			suppressed: map[string]int{schema.CatOmitEmpty: 1},
		},
		"stop_task": {actionable: map[string]int{}, suppressed: map[string]int{}},
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
// finding on get_tasks except the seven allowlisted SCHEMA-DEFECT findings
// the schema itself carries (see
// internal/codegen/testdata/schema/get_tasks.json's Defects array and the
// matching DefectAllowlist entries in overrides.go). resume_task.report_id
// and start_task.report_id each keep one suppressed OMITEMPTY: their
// RequiredOverrides entries deliberately keep omitempty despite the schema's
// declared requiredness.
//
// task_count's bare leading total (the top-level get_tasks_response summary
// wrapper — not the excluded per-item TaskReportCount/TaskResultCount/
// TaskComplianceCount) is now recovered via overrides.go's
// FieldOverrides["get_tasks:/task_count"] (T-23), adding one more suppressed
// EXTRA below.
func TestValidateTaskGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_task")

	wantAfter := map[string]counts{
		"create_task": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_task": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_tasks": {
			actionable: map[string]int{schema.CatDefect: 7},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_task": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"move_task":   {actionable: map[string]int{}, suppressed: map[string]int{}},
		"resume_task": {
			actionable: map[string]int{},
			suppressed: map[string]int{schema.CatOmitEmpty: 1},
		},
		"start_task": {
			actionable: map[string]int{},
			suppressed: map[string]int{schema.CatOmitEmpty: 1},
		},
		"stop_task": {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateScannerBeforeMatchesKnownConformance is scanner's counterpart
// to TestValidateTaskBeforeMatchesKnownConformance. scanner is adopted, so
// this is a drift detector on the committed commands/ files. Pre-adoption
// findings are recorded in docs/review-2026-08.md.
func TestValidateScannerBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_scanner")

	want := map[string]counts{
		"create_scanner": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_scanner": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_scanners": {
			actionable: map[string]int{schema.CatDefect: 10},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_scanner": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"verify_scanner": {
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

// TestValidateScannerGeneratedMatchesOrImproves is scanner's pre-adoption
// gate, run at the same rigour every prior family's was. Generation closes
// every finding on get_scanners except the ten allowlisted SCHEMA-DEFECT
// findings the schema itself carries (see
// internal/codegen/testdata/schema/get_scanners.json's Defects array and the
// matching DefectAllowlist entries in overrides.go).
// get_scanners' scanner_count gains a suppressed EXTRA below via
// overrides.go's FieldOverrides["get_scanners:/scanner_count"] (T-23).
func TestValidateScannerGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_scanner")

	wantAfter := map[string]counts{
		"create_scanner": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"modify_scanner": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"get_scanners": {
			actionable: map[string]int{schema.CatDefect: 10},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"delete_scanner": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"verify_scanner": {
			actionable: map[string]int{},
			suppressed: map[string]int{schema.CatExtra: 1},
		},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateConfigBeforeMatchesKnownConformance is config's counterpart to
// TestValidateScannerBeforeMatchesKnownConformance. config is adopted, so
// this is a drift detector on the committed commands/ files. Pre-adoption
// findings are recorded in docs/review-2026-08.md.
func TestValidateConfigBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_config")

	want := map[string]counts{
		"create_config": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{},
		},
		"get_configs": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"modify_config": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"delete_config": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"sync_config":   {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{}},
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

// TestValidateConfigGeneratedMatchesOrImproves is config's pre-adoption gate,
// run at the same rigour every prior family's was. Generation closes every
// finding on all three commands except the two allowlisted SCHEMA-DEFECT
// findings the schema itself carries: create_config's unsupported-r at
// /create_config and get_configs' empty-pattern at
// /get_configs_response/config/tasks/task/permissions (see
// internal/codegen/testdata/schema/{create_config,get_configs}.json's
// Defects arrays and the matching DefectAllowlist entries in overrides.go).
// modify_config's three CARDINALITY/TYPE pairs close via the boolean_atoi
// TypeOverrides entry (commands.GMPAtoiBool) together with its
// pointerScalars registration — without the latter, these fields would
// regenerate as bare GMPAtoiBool values and reproduce the same CARDINALITY
// defect under the new type.
//
// get_configs' top-level config_count gains a suppressed EXTRA below via
// overrides.go's FieldOverrides["get_configs:/config_count"] (T-23) — this is
// distinct from the excluded per-family ConfigFamilyCount/ConfigNVTCount and
// per-family <nvt_count> mixed-content shape nested under <config>.
func TestValidateConfigGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_config")

	wantAfter := map[string]counts{
		"create_config": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{},
		},
		"get_configs": {
			actionable: map[string]int{schema.CatDefect: 1},
			suppressed: map[string]int{schema.CatExtra: 2, schema.CatMissing: 1, schema.CatOmitEmpty: 1},
		},
		"modify_config": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"delete_config": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"sync_config":   {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidatePortRangeBeforeMatchesKnownConformance is port_range's
// counterpart to TestValidateAssetBeforeMatchesKnownConformance. port_range
// is adopted, so BeforeGoRoots already reflects the generator's output and
// this is a drift detector on the committed files, not the pre-adoption
// MISSING finding create_port_range's <port_list> used to carry (see
// docs/review-2026-08.md's port_range section for that history).
func TestValidatePortRangeBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "create_port_range")

	want := map[string]counts{
		"create_port_range": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"delete_port_range": {actionable: map[string]int{}, suppressed: map[string]int{}},
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

// TestValidatePortRangeGeneratedMatchesOrImproves is port_range's
// pre-adoption gate. Generation reshapes <port_list> into a PortRangePortList
// struct carrying the id attribute, closing the MISSING finding with no
// override table entry needed — the schema already declared the shape.
func TestValidatePortRangeGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "create_port_range")

	wantAfter := map[string]counts{
		"create_port_range": {actionable: map[string]int{}, suppressed: map[string]int{}},
		"delete_port_range": {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}

// TestValidateTrashBeforeMatchesKnownConformance is trash's counterpart to
// TestValidatePortRangeBeforeMatchesKnownConformance. trash is adopted, so
// BeforeGoRoots already reflects the generator's output; empty_trashcan's
// SCHEMA-DEFECT is the allowlisted empty-pattern defect (docs/review-2026-08.md
// `trash`) — allowlisting unblocks generation but the finding itself is
// upstream-schema and stays actionable, same as get_notes' SCHEMA-DEFECT.
func TestValidateTrashBeforeMatchesKnownConformance(t *testing.T) {
	before, _ := validateOrSkip(t, "empty_trashcan")

	want := map[string]counts{
		"empty_trashcan": {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{}},
		"restore":        {actionable: map[string]int{}, suppressed: map[string]int{}},
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

// TestValidateTrashGeneratedMatchesOrImproves is trash's pre-adoption gate.
func TestValidateTrashGeneratedMatchesOrImproves(t *testing.T) {
	before, after := validateOrSkip(t, "empty_trashcan")

	wantAfter := map[string]counts{
		"empty_trashcan": {actionable: map[string]int{schema.CatDefect: 1}, suppressed: map[string]int{}},
		"restore":        {actionable: map[string]int{}, suppressed: map[string]int{}},
	}
	assertMatchesOrImproves(t, before, after, wantAfter)
}
