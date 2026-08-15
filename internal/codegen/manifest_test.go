package codegen

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/brennoo/go-gmp/internal/schema"
)

func loadManifestErr(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), manifestFile)
	writeTestFile(t, path, body)
	if _, err := LoadManifest(path); err != nil {
		return err.Error()
	}
	t.Fatalf("LoadManifest accepted %s, want an error", body)
	return ""
}

// TestCommittedManifest pins the shape of the real manifest.json, which every
// other test and the CLI source their command list from.
func TestCommittedManifest(t *testing.T) {
	m := manifestOrFail(t)

	// note, override, role, group, tag, schedule, alert, filter,
	// tls_certificate, permission, and ticket are all adopted: commands/{create_note,
	// get_notes,modify_note,note_shared,create_override,get_overrides,
	// modify_override,override_shared,create_role,get_roles,modify_role,
	// create_group,get_groups,modify_group,create_tag,get_tags,modify_tag,
	// tag_shared,create_schedule,get_schedules,modify_schedule,create_alert,
	// get_alerts,modify_alert,alert_shared,create_filter,get_filters,
	// modify_filter,create_tls_certificate,get_tls_certificates,
	// modify_tls_certificate,create_permission,get_permissions,
	// modify_permission,permission_shared,create_ticket,get_tickets,
	// modify_ticket,ticket_shared}.go are all the generator's output.
	// credential is review-only: its golden output lives under
	// testdata/golden but nothing has been written to commands/ yet.
	assertFamily(t, m, "get_notes", "note", false)
	assertFamily(t, m, "get_overrides", "override", false)
	assertFamily(t, m, "get_roles", "role", false)
	assertFamily(t, m, "get_groups", "group", false)
	assertFamily(t, m, "get_tags", "tag", false)
	assertFamily(t, m, "get_schedules", "schedule", false)
	assertFamily(t, m, "get_alerts", "alert", false)
	assertFamily(t, m, "get_filters", "filter", false)
	assertFamily(t, m, "get_tls_certificates", "tls_certificate", false)
	assertFamily(t, m, "get_permissions", "permission", false)
	assertFamily(t, m, "get_tickets", "ticket", false)
	assertFamily(t, m, "get_credentials", "credential", true)
	if _, ok := m.FamilyOf("delete_note"); ok {
		t.Error("FamilyOf(delete_note) claims a family the manifest does not declare")
	}

	want := []string{
		"create_alert", "create_credential", "create_filter", "create_group", "create_note", "create_override",
		"create_permission", "create_role", "create_schedule", "create_tag", "create_ticket",
		"create_tls_certificate", "get_alerts", "get_credentials", "get_filters", "get_groups", "get_notes",
		"get_overrides", "get_permissions", "get_roles", "get_schedules", "get_tags", "get_tickets",
		"get_tls_certificates", "modify_alert", "modify_credential", "modify_filter", "modify_group",
		"modify_note", "modify_override", "modify_permission", "modify_role", "modify_schedule", "modify_tag",
		"modify_ticket", "modify_tls_certificate",
	}
	if got := m.AllCommands(); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("AllCommands() = %v, want %v", got, want)
	}
	wantAdopted := []string{
		"create_alert", "create_filter", "create_group", "create_note", "create_override", "create_permission",
		"create_role", "create_schedule", "create_tag", "create_ticket", "create_tls_certificate", "get_alerts",
		"get_filters", "get_groups", "get_notes", "get_overrides", "get_permissions", "get_roles",
		"get_schedules", "get_tags", "get_tickets", "get_tls_certificates", "modify_alert", "modify_filter",
		"modify_group", "modify_note", "modify_override", "modify_permission", "modify_role", "modify_schedule",
		"modify_tag", "modify_ticket", "modify_tls_certificate",
	}
	if got := m.Adopted().AllCommands(); strings.Join(got, ",") != strings.Join(wantAdopted, ",") {
		t.Errorf("Adopted().AllCommands() = %v, want %v", got, wantAdopted)
	}
	if got := m.ReviewOnlyFamilies(); strings.Join(got, ",") != "credential" {
		t.Errorf("ReviewOnlyFamilies() = %v, want [credential]", got)
	}

	assertFiles(t, "ExpectedOutputFiles", m.ExpectedOutputFiles("commands"), []string{
		"alert_shared.go", "create_alert.go", "create_credential.go", "create_filter.go", "create_group.go",
		"create_note.go", "create_override.go", "create_permission.go", "create_role.go", "create_schedule.go",
		"create_tag.go", "create_ticket.go", "create_tls_certificate.go", "credential_shared.go",
		"filter_shared.go", "get_alerts.go", "get_credentials.go", "get_filters.go", "get_groups.go",
		"get_notes.go", "get_overrides.go", "get_permissions.go", "get_roles.go", "get_schedules.go",
		"get_tags.go", "get_tickets.go", "get_tls_certificates.go", "group_shared.go", "modify_alert.go",
		"modify_credential.go", "modify_filter.go", "modify_group.go", "modify_note.go", "modify_override.go",
		"modify_permission.go", "modify_role.go", "modify_schedule.go", "modify_tag.go", "modify_ticket.go",
		"modify_tls_certificate.go", "note_shared.go", "override_shared.go", "permission_shared.go",
		"role_shared.go", "schedule_shared.go", "tag_shared.go", "ticket_shared.go", "tls_certificate_shared.go",
	})
	assertFiles(t, "Adopted().ExpectedOutputFiles", m.Adopted().ExpectedOutputFiles("commands"), []string{
		"alert_shared.go", "create_alert.go", "create_filter.go", "create_group.go", "create_note.go",
		"create_override.go", "create_permission.go", "create_role.go", "create_schedule.go", "create_tag.go",
		"create_ticket.go", "create_tls_certificate.go", "filter_shared.go", "get_alerts.go", "get_filters.go",
		"get_groups.go", "get_notes.go", "get_overrides.go", "get_permissions.go", "get_roles.go",
		"get_schedules.go", "get_tags.go", "get_tickets.go", "get_tls_certificates.go", "group_shared.go",
		"modify_alert.go", "modify_filter.go", "modify_group.go", "modify_note.go", "modify_override.go",
		"modify_permission.go", "modify_role.go", "modify_schedule.go", "modify_tag.go", "modify_ticket.go",
		"modify_tls_certificate.go", "note_shared.go", "override_shared.go", "permission_shared.go",
		"role_shared.go", "schedule_shared.go", "tag_shared.go", "ticket_shared.go", "tls_certificate_shared.go",
	})

	if got, want := m.SnapshotPath("d", "get_notes"), SnapshotPath("d", "get_notes"); got != want {
		t.Errorf("Manifest.SnapshotPath = %q, want the package-level %q", got, want)
	}
}

// assertFamily checks that command is owned by the named family, that the
// family's shared file follows the convention, and that its ReviewOnly flag
// matches reviewOnly.
func assertFamily(t *testing.T, m Manifest, command, family string, reviewOnly bool) {
	t.Helper()
	f, ok := m.FamilyOf(command)
	if !ok {
		t.Fatalf("FamilyOf(%s) found no family", command)
	}
	if f.Name != family {
		t.Errorf("FamilyOf(%s).Name = %q, want %s", command, f.Name, family)
	}
	if want := family + "_shared.go"; f.SharedFile() != want {
		t.Errorf("%s: SharedFile() = %q, want %s", family, f.SharedFile(), want)
	}
	if f.ReviewOnly != reviewOnly {
		t.Errorf("%s: ReviewOnly = %v, want %v", family, f.ReviewOnly, reviewOnly)
	}
}

func assertFiles(t *testing.T, what string, got, wantNames []string) {
	t.Helper()
	want := make([]string, 0, len(wantNames))
	for _, n := range wantNames {
		want = append(want, filepath.Join("commands", n))
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("%s() = %v, want %v", what, got, want)
	}
}

// TestManifestCoversGeneratedTree ties the manifest to what the generator
// actually emits: every emitted file must be one ExpectedOutputFiles predicted,
// which is what lets -check trust the convention instead of a JSON field.
func TestManifestCoversGeneratedTree(t *testing.T) {
	m := manifestOrFail(t)
	expected := map[string]bool{}
	for _, p := range m.ExpectedOutputFiles("") {
		expected[p] = true
	}
	for name := range generateOrFail(t) {
		if !expected[name] {
			t.Errorf("The generator emits %s, which ExpectedOutputFiles does not predict", name)
		}
	}
}

func TestLoadManifestRejects(t *testing.T) {
	cases := []struct {
		name, body, want string
	}{
		{"unsorted commands",
			`{"families":[{"name":"note","commands":["modify_note","create_note"]}]}`,
			"must be sorted"},
		{"command in two families",
			`{"families":[{"name":"note","commands":["create_note"]},{"name":"tag","commands":["create_note"]}]}`,
			`claimed by both family "note" and family "tag"`},
		{"empty family name",
			`{"families":[{"name":"","commands":["create_note"]}]}`,
			"has no name"},
		{"no commands",
			`{"families":[{"name":"note","commands":[]}]}`,
			"declares no commands"},
		{"empty command name",
			`{"families":[{"name":"note","commands":[""]}]}`,
			"empty command name"},
		{"duplicate family",
			`{"families":[{"name":"note","commands":["create_note"]},{"name":"note","commands":["get_notes"]}]}`,
			"declared twice"},
		{"no families",
			`{"families":[]}`,
			"declares no families"},
		{"unknown field",
			`{"families":[{"name":"note","commands":["create_note"],"shared":"note_shared.go"}]}`,
			"unknown field"},
		{"not JSON", `nope`, "invalid character"},
		{"trailing JSON document",
			`{"families":[{"name":"note","commands":["create_note"]}]}` + "\n" + `{"unexpected":"document"}`,
			"multiple JSON documents"},
		{"trailing garbage",
			`{"families":[{"name":"note","commands":["create_note"]}]}` + "garbage",
			"trailing data"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := loadManifestErr(t, tc.body); !strings.Contains(got, tc.want) {
				t.Errorf("Error = %q, want it to mention %q", got, tc.want)
			}
		})
	}
}

func TestLoadManifestMissingFile(t *testing.T) {
	if _, err := LoadManifest(filepath.Join(t.TempDir(), "absent.json")); err == nil {
		t.Fatal("LoadManifest accepted a nonexistent path")
	}
}

// TestGenerateAllIsPerFamily is the isolation property (a) buys: a second
// family in the manifest must not change the first family's bytes.
func TestGenerateAllIsPerFamily(t *testing.T) {
	m := manifestOrFail(t)
	cmds, err := LoadSnapshots(snapshotDir, m.AllCommands())
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}
	alone, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings)
	if err != nil {
		t.Fatalf("GenerateAll failed: %v", err)
	}

	// A second family built from a renamed copy of get_notes. Its shapes are
	// structurally identical to the note family's, so a global dedup pass
	// would have to disambiguate the names of one of them; a per-family pass
	// cannot see the other family at all.
	twin, ok := findCommand(cmds, "get_notes")
	if !ok {
		t.Fatal("No get_notes snapshot")
	}
	twin.Name = "get_twins"
	two := Manifest{Families: append(append([]Family{}, m.Families...), Family{Name: "twin", Commands: []string{"get_twins"}})}
	both, err := GenerateAll(two, append(append([]schema.Command{}, cmds...), twin), TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings)
	if err != nil {
		t.Fatalf("GenerateAll with a second family failed: %v", err)
	}

	for name, src := range alone {
		if both[name] != src {
			t.Errorf("%s changed when an unrelated family was added to the manifest", name)
		}
	}
	if _, ok := both["get_twins.go"]; !ok {
		t.Error("The second family emitted no get_twins.go")
	}
}

// TestRealFamiliesGenerateIndependently is TestGenerateAllIsPerFamily's
// non-synthetic counterpart, and the one that would actually have caught a
// regression when the override family was added to the manifest: it asserts
// that every real family's bytes are the same generated alone as generated
// alongside the others.
//
// The note case is the historically interesting one — Manifest{note} is
// literally the manifest as committed before override was added, so a failure
// there would mean adding override silently renamed a type or moved a struct in
// commands/{create_note,get_notes,modify_note,note_shared}.go, which are
// already committed and already called from client/ and internal/gmptest. The
// two families are a real collision risk rather than a hypothetical one: both
// describe <nvt>, <result>, <task>, <text>, <owner>, <permissions> and a dozen
// more elements with identical shapes, so a global dedup pass would have had to
// disambiguate one side of every one of them.
func TestRealFamiliesGenerateIndependently(t *testing.T) {
	m := manifestOrFail(t)
	if len(m.Families) < 2 {
		t.Fatalf("The manifest declares %d family; this test needs at least two", len(m.Families))
	}
	cmds, err := LoadSnapshots(snapshotDir, m.AllCommands())
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}
	together, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings)
	if err != nil {
		t.Fatalf("GenerateAll failed: %v", err)
	}

	accounted := 0
	for _, f := range m.Families {
		alone, err := GenerateAll(Manifest{Families: []Family{f}}, cmds, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings)
		if err != nil {
			t.Fatalf("GenerateAll(%s alone) failed: %v", f.Name, err)
		}
		if len(alone) == 0 {
			t.Fatalf("Generating %s alone emitted nothing", f.Name)
		}
		for name, src := range alone {
			accounted++
			switch got, ok := together[name]; {
			case !ok:
				t.Errorf("%s emits %s alone but not alongside the other families", f.Name, name)
			case got != src:
				t.Errorf("%s: %s differs between generating %s alone and generating the whole manifest\n%s",
					f.Name, name, f.Name, firstDiff(src, got))
			}
		}
	}
	if accounted != len(together) {
		t.Errorf("The per-family runs account for %d files, the combined run emits %d", accounted, len(together))
	}
}

// TestScheduleICalendarNamingIsStable pins the fix for "icalendar" in
// generate.go's initialisms table: the schema element and the adopted
// commands/{create_schedule,modify_schedule,get_schedules}.go all spell the
// Go field "ICalendar" (RFC 5545's own casing), and without the table entry
// goName would instead render it "Icalendar" — a public API regression.
//
// Generating schedule alone and generating it alongside every other real
// family in the committed manifest must produce byte-identical output for
// this reason: initialisms is a package-level table with no per-family
// state, so nothing about combining families could plausibly move this
// spelling — but TestRealFamiliesGenerateIndependently above only proves
// alone-vs-together consistency in general, not that the agreed-upon spelling
// is the *correct* one. This test pins both: the exact identifier, and that
// no other family's presence changes it.
func TestScheduleICalendarNamingIsStable(t *testing.T) {
	m := manifestOrFail(t)
	cmds, err := LoadSnapshots(snapshotDir, m.AllCommands())
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}

	fam, ok := m.FamilyOf("get_schedules")
	if !ok {
		t.Fatal("The manifest no longer declares a schedule family")
	}
	alone, err := GenerateAll(Manifest{Families: []Family{fam}}, cmds, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings)
	if err != nil {
		t.Fatalf("GenerateAll(schedule alone) failed: %v", err)
	}
	together, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings)
	if err != nil {
		t.Fatalf("GenerateAll(full manifest) failed: %v", err)
	}

	for _, name := range []string{"create_schedule.go", "modify_schedule.go", "get_schedules.go"} {
		aloneSrc, ok := alone[name]
		if !ok {
			t.Fatalf("Generating schedule alone did not emit %s", name)
		}
		togetherSrc, ok := together[name]
		if !ok {
			t.Fatalf("Generating the full manifest did not emit %s", name)
		}
		if aloneSrc != togetherSrc {
			t.Errorf("%s differs between generating schedule alone and generating the full manifest\n%s",
				name, firstDiff(aloneSrc, togetherSrc))
		}
		if !strings.Contains(aloneSrc, "ICalendar") {
			t.Errorf("%s does not spell the field ICalendar", name)
		}
		if strings.Contains(aloneSrc, "Icalendar ") || strings.Contains(aloneSrc, "Icalendar\t") {
			t.Errorf("%s spells the field Icalendar (single capital) instead of ICalendar", name)
		}
	}
}

// TestAlertFingerprintNamingIsStable pins the fix for "md5"/"sha256" in
// generate.go's initialisms table: get_alerts's certificate_info element
// exposes md5_fingerprint/sha256_fingerprint, and the hand-written
// commands/get_reports.go and commands/get_scanners.go, plus the adopted
// commands/get_tls_certificates.go, all already spell these fields
// MD5Fingerprint/SHA256Fingerprint. Without the table entries goName would
// instead render them Md5Fingerprint/Sha256Fingerprint — an inconsistent
// public API that alert's adoption (commands/get_alerts.go is now the
// generator's output) must not carry.
//
// As with TestScheduleICalendarNamingIsStable, generating alert alone and
// generating it alongside every other real family in the committed manifest
// must produce byte-identical output for this reason: initialisms is a
// package-level table with no per-family state, so nothing about combining
// families could plausibly move this spelling. This test pins both: the
// exact identifiers, and that no other family's presence changes them.
func TestAlertFingerprintNamingIsStable(t *testing.T) {
	m := manifestOrFail(t)
	cmds, err := LoadSnapshots(snapshotDir, m.AllCommands())
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}

	fam, ok := m.FamilyOf("get_alerts")
	if !ok {
		t.Fatal("The manifest no longer declares an alert family")
	}
	alone, err := GenerateAll(Manifest{Families: []Family{fam}}, cmds, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings)
	if err != nil {
		t.Fatalf("GenerateAll(alert alone) failed: %v", err)
	}
	together, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings)
	if err != nil {
		t.Fatalf("GenerateAll(full manifest) failed: %v", err)
	}

	const name = "get_alerts.go"
	aloneSrc, ok := alone[name]
	if !ok {
		t.Fatalf("Generating alert alone did not emit %s", name)
	}
	togetherSrc, ok := together[name]
	if !ok {
		t.Fatalf("Generating the full manifest did not emit %s", name)
	}
	if aloneSrc != togetherSrc {
		t.Errorf("%s differs between generating alert alone and generating the full manifest\n%s",
			name, firstDiff(aloneSrc, togetherSrc))
	}
	for _, want := range []string{"MD5Fingerprint", "SHA256Fingerprint"} {
		if !strings.Contains(aloneSrc, want) {
			t.Errorf("%s does not spell the field %s", name, want)
		}
	}
	for _, unwanted := range []string{"Md5Fingerprint", "Sha256Fingerprint"} {
		if strings.Contains(aloneSrc, unwanted) {
			t.Errorf("%s spells the field %s instead of the fixed-casing form", name, unwanted)
		}
	}
}

// TestTLSCertificateDNNamingIsStable pins the fix for "dn" in generate.go's
// initialisms table: get_tls_certificates's subject_dn/issuer_dn elements and
// the hand-written commands/get_reports.go both already spell the Go fields
// SubjectDN/IssuerDN, and the adopted commands/get_tls_certificates.go must
// match. Without the table entry goName would instead render them
// SubjectDn/IssuerDn — a public API regression against the casing
// tls_certificate's own adoption locked in (T-9).
//
// As with TestScheduleICalendarNamingIsStable and
// TestAlertFingerprintNamingIsStable, generating tls_certificate alone and
// generating it alongside every other family in the committed manifest must
// produce byte-identical output for this reason: initialisms is a
// package-level table with no per-family state, so nothing about combining
// families could plausibly move this spelling. This test pins both: the
// exact identifiers, and that no other family's presence changes them.
func TestTLSCertificateDNNamingIsStable(t *testing.T) {
	m := manifestOrFail(t)
	cmds, err := LoadSnapshots(snapshotDir, m.AllCommands())
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}

	fam, ok := m.FamilyOf("get_tls_certificates")
	if !ok {
		t.Fatal("The manifest no longer declares a tls_certificate family")
	}
	alone, err := GenerateAll(Manifest{Families: []Family{fam}}, cmds, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings)
	if err != nil {
		t.Fatalf("GenerateAll(tls_certificate alone) failed: %v", err)
	}
	together, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings)
	if err != nil {
		t.Fatalf("GenerateAll(full manifest) failed: %v", err)
	}

	const name = "get_tls_certificates.go"
	aloneSrc, ok := alone[name]
	if !ok {
		t.Fatalf("Generating tls_certificate alone did not emit %s", name)
	}
	togetherSrc, ok := together[name]
	if !ok {
		t.Fatalf("Generating the full manifest did not emit %s", name)
	}
	if aloneSrc != togetherSrc {
		t.Errorf("%s differs between generating tls_certificate alone and generating the full manifest\n%s",
			name, firstDiff(aloneSrc, togetherSrc))
	}
	for _, want := range []string{"SubjectDN", "IssuerDN"} {
		if !strings.Contains(aloneSrc, want) {
			t.Errorf("%s does not spell the field %s", name, want)
		}
	}
	for _, unwanted := range []string{"SubjectDn ", "SubjectDn\t", "IssuerDn ", "IssuerDn\t"} {
		if strings.Contains(aloneSrc, unwanted) {
			t.Errorf("%s contains %q; want the fixed DN casing", name, unwanted)
		}
	}
}

func findCommand(cmds []schema.Command, name string) (schema.Command, bool) {
	for _, c := range cmds {
		if c.Name == name {
			return c, true
		}
	}
	return schema.Command{}, false
}

// TestGenerateAllRejectsAnUnknownFamilyName guards the one way the
// file-naming convention can go wrong: a family whose declared name is not the
// noun the generator derives from its commands. Silently accepting it would
// put the shared types in a file ExpectedOutputFiles never predicts, and
// -check would then call that file obsolete.
func TestGenerateAllRejectsAnUnknownFamilyName(t *testing.T) {
	m := manifestOrFail(t)
	cmds, err := LoadSnapshots(snapshotDir, m.AllCommands())
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}
	wrong := Manifest{Families: []Family{{Name: "annotation", Commands: m.AllCommands()}}}
	if _, err := GenerateAll(wrong, cmds, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings); err == nil ||
		!strings.Contains(err.Error(), "naming convention") {
		t.Fatalf("GenerateAll error = %v, want a naming-convention complaint", err)
	}
}

func TestGenerateAllNeedsEverySnapshot(t *testing.T) {
	m := manifestOrFail(t)
	cmds, err := LoadSnapshots(snapshotDir, []string{"create_note"})
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}
	if _, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings); err == nil ||
		!strings.Contains(err.Error(), "no snapshot loaded for get_notes") {
		t.Fatalf("GenerateAll error = %v, want a missing-snapshot complaint", err)
	}
}
