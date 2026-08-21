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

	// port_range has no get_X/modify_X (schema only declares create/delete);
	// report has no modify_X; trash has only empty_trashcan/restore.
	assertFamily(t, m, "authenticate", "auth")
	assertFamily(t, m, "describe_auth", "auth")
	assertFamily(t, m, "modify_auth", "auth")
	assertFamily(t, m, "get_aggregates", "aggregate")
	assertFamily(t, m, "get_feeds", "feed")
	assertFamily(t, m, "get_features", "feature")
	assertFamily(t, m, "get_system_reports", "system_report")
	assertFamily(t, m, "get_license", "license")
	assertFamily(t, m, "modify_license", "license")
	assertFamily(t, m, "get_resource_names", "settings")
	assertFamily(t, m, "get_settings", "settings")
	assertFamily(t, m, "modify_setting", "settings")
	assertFamily(t, m, "get_version", "version")
	assertFamily(t, m, "get_vulns", "vuln")
	assertFamily(t, m, "get_results", "result")
	assertFamily(t, m, "empty_trashcan", "trash")
	assertFamily(t, m, "restore", "trash")
	assertFamily(t, m, "get_notes", "note")
	assertFamily(t, m, "get_reports", "report")
	assertFamily(t, m, "get_overrides", "override")
	assertFamily(t, m, "get_roles", "role")
	assertFamily(t, m, "get_groups", "group")
	assertFamily(t, m, "get_tags", "tag")
	assertFamily(t, m, "get_schedules", "schedule")
	assertFamily(t, m, "get_alerts", "alert")
	assertFamily(t, m, "get_filters", "filter")
	assertFamily(t, m, "get_tls_certificates", "tls_certificate")
	assertFamily(t, m, "get_permissions", "permission")
	assertFamily(t, m, "get_tickets", "ticket")
	assertFamily(t, m, "get_credentials", "credential")
	assertFamily(t, m, "get_port_lists", "port_list")
	assertFamily(t, m, "get_assets", "asset")
	assertFamily(t, m, "get_report_formats", "report_format")
	assertFamily(t, m, "get_report_configs", "report_config")
	assertFamily(t, m, "get_users", "user")
	assertFamily(t, m, "get_targets", "target")
	assertFamily(t, m, "get_tasks", "task")
	assertFamily(t, m, "get_scanners", "scanner")
	assertFamily(t, m, "get_configs", "config")
	assertFamily(t, m, "create_port_range", "port_range")
	assertFamily(t, m, "delete_note", "note")
	assertFamily(t, m, "delete_report", "report")
	assertFamily(t, m, "delete_override", "override")
	assertFamily(t, m, "delete_role", "role")
	assertFamily(t, m, "delete_group", "group")
	assertFamily(t, m, "delete_tag", "tag")
	assertFamily(t, m, "delete_schedule", "schedule")
	assertFamily(t, m, "delete_alert", "alert")
	assertFamily(t, m, "delete_filter", "filter")
	assertFamily(t, m, "delete_permission", "permission")
	assertFamily(t, m, "delete_ticket", "ticket")
	assertFamily(t, m, "delete_credential", "credential")
	assertFamily(t, m, "delete_port_list", "port_list")
	assertFamily(t, m, "delete_report_format", "report_format")
	assertFamily(t, m, "delete_report_config", "report_config")
	assertFamily(t, m, "delete_scanner", "scanner")
	assertFamily(t, m, "delete_config", "config")
	assertFamily(t, m, "delete_port_range", "port_range")
	assertFamily(t, m, "delete_asset", "asset")
	assertFamily(t, m, "delete_target", "target")
	assertFamily(t, m, "delete_task", "task")
	assertFamily(t, m, "delete_user", "user")

	wantAdopted := []string{
		"authenticate", "create_alert", "create_asset", "create_config", "create_credential", "create_filter", "create_group",
		"create_note", "create_override", "create_permission", "create_port_list", "create_port_range", "create_report", "create_report_config",
		"create_report_format", "create_role", "create_scanner", "create_schedule", "create_tag", "create_target",
		"create_task", "create_ticket", "create_tls_certificate", "create_user", "delete_alert", "delete_asset", "delete_config",
		"delete_credential", "delete_filter", "delete_group", "delete_note", "delete_override",
		"delete_permission", "delete_port_list", "delete_port_range", "delete_report", "delete_report_config", "delete_report_format", "delete_role",
		"delete_scanner", "delete_schedule", "delete_tag", "delete_target", "delete_task", "delete_ticket",
		"delete_user", "describe_auth", "empty_trashcan", "get_aggregates", "get_alerts", "get_assets",
		"get_configs", "get_credentials", "get_features", "get_feeds", "get_filters", "get_groups", "get_info", "get_license", "get_notes",
		"get_nvt_families", "get_nvts", "get_overrides",
		"get_permissions", "get_port_lists", "get_preferences", "get_report_configs", "get_report_formats", "get_reports", "get_resource_names", "get_results", "get_roles",
		"get_scanners", "get_schedules", "get_settings", "get_system_reports", "get_tags", "get_targets", "get_tasks", "get_tickets",
		"get_tls_certificates", "get_users", "get_version", "get_vulns", "modify_alert", "modify_asset", "modify_auth", "modify_config", "modify_credential",
		"modify_filter", "modify_group", "modify_license", "modify_note", "modify_override", "modify_permission", "modify_port_list",
		"modify_report_config", "modify_report_format", "modify_role", "modify_scanner", "modify_schedule", "modify_setting",
		"modify_tag", "modify_target", "modify_task", "modify_ticket", "modify_tls_certificate", "modify_user",
		"move_task", "restore", "resume_task", "start_task", "stop_task",
		"sync_config", "test_alert", "verify_report_format", "verify_scanner",
	}
	want := []string{
		"authenticate", "create_alert", "create_asset", "create_config", "create_credential", "create_filter", "create_group",
		"create_note", "create_override", "create_permission", "create_port_list", "create_port_range", "create_report", "create_report_config",
		"create_report_format", "create_role", "create_scanner", "create_schedule", "create_tag", "create_target",
		"create_task", "create_ticket", "create_tls_certificate", "create_user", "delete_alert", "delete_asset", "delete_config",
		"delete_credential", "delete_filter", "delete_group", "delete_note", "delete_override",
		"delete_permission", "delete_port_list", "delete_port_range", "delete_report", "delete_report_config", "delete_report_format", "delete_role",
		"delete_scanner", "delete_schedule", "delete_tag", "delete_target", "delete_task", "delete_ticket",
		"delete_user", "describe_auth", "empty_trashcan", "get_aggregates", "get_alerts", "get_assets",
		"get_configs", "get_credentials", "get_features", "get_feeds", "get_filters", "get_groups", "get_info", "get_license", "get_notes",
		"get_nvt_families", "get_nvts", "get_overrides",
		"get_permissions", "get_port_lists", "get_preferences", "get_report_configs", "get_report_formats", "get_reports", "get_resource_names", "get_results", "get_roles",
		"get_scanners", "get_schedules", "get_settings", "get_system_reports", "get_tags", "get_targets", "get_tasks", "get_tickets",
		"get_tls_certificates", "get_users", "get_version", "get_vulns", "modify_alert", "modify_asset", "modify_auth", "modify_config", "modify_credential",
		"modify_filter", "modify_group", "modify_license", "modify_note", "modify_override", "modify_permission", "modify_port_list",
		"modify_report_config", "modify_report_format", "modify_role", "modify_scanner", "modify_schedule", "modify_setting",
		"modify_tag", "modify_target", "modify_task", "modify_ticket", "modify_tls_certificate", "modify_user",
		"move_task", "restore", "resume_task", "start_task", "stop_task",
		"sync_config", "test_alert", "verify_report_format", "verify_scanner",
	}
	if got := m.AllCommands(); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("AllCommands() = %v, want %v", got, want)
	}
	if got := m.Adopted().AllCommands(); strings.Join(got, ",") != strings.Join(wantAdopted, ",") {
		t.Errorf("Adopted().AllCommands() = %v, want %v", got, wantAdopted)
	}
	if got := m.ReviewOnlyFamilies(); len(got) != 0 {
		t.Errorf("ReviewOnlyFamilies() = %v, want none", got)
	}

	assertFiles(t, "ExpectedOutputFiles", m.ExpectedOutputFiles("commands"), []string{
		"aggregate_shared.go", "alert_shared.go", "asset_shared.go", "auth_shared.go", "authenticate.go", "config_shared.go", "create_alert.go", "create_asset.go",
		"create_config.go", "create_credential.go", "create_filter.go", "create_group.go", "create_note.go",
		"create_override.go", "create_permission.go", "create_port_list.go", "create_port_range.go", "create_report.go", "create_report_config.go",
		"create_report_format.go", "create_role.go", "create_scanner.go", "create_schedule.go", "create_tag.go",
		"create_target.go", "create_task.go", "create_ticket.go", "create_tls_certificate.go", "create_user.go",
		"credential_shared.go", "delete_alert.go", "delete_asset.go", "delete_config.go", "delete_credential.go",
		"delete_filter.go",
		"delete_group.go", "delete_note.go", "delete_override.go", "delete_permission.go", "delete_port_list.go",
		"delete_port_range.go", "delete_report.go", "delete_report_config.go", "delete_report_format.go", "delete_role.go", "delete_scanner.go",
		"delete_schedule.go", "delete_tag.go", "delete_target.go", "delete_task.go", "delete_ticket.go",
		"delete_user.go", "describe_auth.go", "empty_trashcan.go", "feature_shared.go", "feed_shared.go", "filter_shared.go", "get_aggregates.go", "get_alerts.go",
		"get_assets.go", "get_configs.go", "get_credentials.go", "get_features.go", "get_feeds.go", "get_filters.go", "get_groups.go", "get_info.go", "get_license.go", "get_notes.go",
		"get_nvt_families.go", "get_nvts.go",
		"get_overrides.go", "get_permissions.go", "get_port_lists.go", "get_preferences.go", "get_report_configs.go",
		"get_report_formats.go", "get_reports.go", "get_resource_names.go", "get_results.go", "get_roles.go", "get_scanners.go", "get_schedules.go", "get_settings.go",
		"get_system_reports.go", "get_tags.go",
		"get_targets.go", "get_tasks.go", "get_tickets.go", "get_tls_certificates.go", "get_users.go", "get_version.go", "get_vulns.go",
		"group_shared.go", "license_shared.go", "modify_alert.go", "modify_asset.go", "modify_auth.go", "modify_config.go", "modify_credential.go",
		"modify_filter.go", "modify_group.go", "modify_license.go", "modify_note.go", "modify_override.go", "modify_permission.go",
		"modify_port_list.go", "modify_report_config.go", "modify_report_format.go", "modify_role.go",
		"modify_scanner.go", "modify_schedule.go", "modify_setting.go", "modify_tag.go", "modify_target.go", "modify_task.go",
		"modify_ticket.go", "modify_tls_certificate.go", "modify_user.go", "move_task.go", "note_shared.go", "nvt_shared.go", "override_shared.go",
		"permission_shared.go", "port_list_shared.go", "port_range_shared.go", "report_config_shared.go", "report_format_shared.go", "report_shared.go",
		"restore.go", "result_shared.go", "resume_task.go", "role_shared.go", "scanner_shared.go", "schedule_shared.go", "settings_shared.go", "start_task.go", "stop_task.go",
		"sync_config.go", "system_report_shared.go", "tag_shared.go",
		"target_shared.go", "task_shared.go", "test_alert.go", "ticket_shared.go", "tls_certificate_shared.go",
		"trash_shared.go", "user_shared.go", "verify_report_format.go", "verify_scanner.go", "version_shared.go", "vuln_shared.go",
	})
	assertFiles(t, "Adopted().ExpectedOutputFiles", m.Adopted().ExpectedOutputFiles("commands"), []string{
		"aggregate_shared.go", "alert_shared.go", "asset_shared.go", "auth_shared.go", "authenticate.go", "config_shared.go", "create_alert.go", "create_asset.go",
		"create_config.go", "create_credential.go", "create_filter.go", "create_group.go", "create_note.go",
		"create_override.go", "create_permission.go", "create_port_list.go", "create_port_range.go", "create_report.go", "create_report_config.go",
		"create_report_format.go", "create_role.go", "create_scanner.go", "create_schedule.go", "create_tag.go",
		"create_target.go", "create_task.go", "create_ticket.go", "create_tls_certificate.go", "create_user.go",
		"credential_shared.go", "delete_alert.go", "delete_asset.go", "delete_config.go", "delete_credential.go",
		"delete_filter.go",
		"delete_group.go", "delete_note.go", "delete_override.go", "delete_permission.go", "delete_port_list.go",
		"delete_port_range.go", "delete_report.go", "delete_report_config.go", "delete_report_format.go", "delete_role.go", "delete_scanner.go",
		"delete_schedule.go", "delete_tag.go", "delete_target.go", "delete_task.go", "delete_ticket.go",
		"delete_user.go", "describe_auth.go", "empty_trashcan.go", "feature_shared.go", "feed_shared.go", "filter_shared.go", "get_aggregates.go", "get_alerts.go",
		"get_assets.go", "get_configs.go", "get_credentials.go", "get_features.go", "get_feeds.go", "get_filters.go", "get_groups.go", "get_info.go", "get_license.go", "get_notes.go",
		"get_nvt_families.go", "get_nvts.go",
		"get_overrides.go", "get_permissions.go", "get_port_lists.go", "get_preferences.go", "get_report_configs.go",
		"get_report_formats.go", "get_reports.go", "get_resource_names.go", "get_results.go", "get_roles.go", "get_scanners.go", "get_schedules.go", "get_settings.go",
		"get_system_reports.go", "get_tags.go",
		"get_targets.go", "get_tasks.go", "get_tickets.go", "get_tls_certificates.go", "get_users.go", "get_version.go", "get_vulns.go",
		"group_shared.go", "license_shared.go", "modify_alert.go", "modify_asset.go", "modify_auth.go", "modify_config.go", "modify_credential.go",
		"modify_filter.go", "modify_group.go", "modify_license.go", "modify_note.go", "modify_override.go", "modify_permission.go",
		"modify_port_list.go", "modify_report_config.go", "modify_report_format.go", "modify_role.go",
		"modify_scanner.go", "modify_schedule.go", "modify_setting.go", "modify_tag.go", "modify_target.go", "modify_task.go",
		"modify_ticket.go", "modify_tls_certificate.go", "modify_user.go", "move_task.go", "note_shared.go", "nvt_shared.go", "override_shared.go",
		"permission_shared.go", "port_list_shared.go", "port_range_shared.go", "report_config_shared.go", "report_format_shared.go", "report_shared.go",
		"restore.go", "result_shared.go", "resume_task.go", "role_shared.go", "scanner_shared.go", "schedule_shared.go", "settings_shared.go", "start_task.go", "stop_task.go",
		"sync_config.go", "system_report_shared.go", "tag_shared.go",
		"target_shared.go", "task_shared.go", "test_alert.go", "ticket_shared.go", "tls_certificate_shared.go",
		"trash_shared.go", "user_shared.go", "verify_report_format.go", "verify_scanner.go", "version_shared.go", "vuln_shared.go",
	})

	if got, want := m.SnapshotPath("d", "get_notes"), SnapshotPath("d", "get_notes"); got != want {
		t.Errorf("Manifest.SnapshotPath = %q, want the package-level %q", got, want)
	}
}

// assertFamily checks that command is owned by the named family, that the
// family's shared file follows the convention, and that the family is
// adopted (not review-only).
func assertFamily(t *testing.T, m Manifest, command, family string) {
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
	if f.ReviewOnly {
		t.Errorf("%s: ReviewOnly = true, want false", family)
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
	cmds, err := LoadGenerationModels(snapshotDir, m.AllCommands())
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}
	alone, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist, NameOverrides, FieldTypeOverrides, nil, EnumConstantNames)
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
	// Defects are provenance-tied to the real command they were found on;
	// a renamed copy has no allowlist entry under its new name and
	// checkDefects would reject it for a defect that isn't really its own.
	twin.Defects = nil
	two := Manifest{Families: append(append([]Family{}, m.Families...), Family{Name: "twin", Commands: []string{"get_twins"}})}
	both, err := GenerateAll(two, append(append([]schema.Command{}, cmds...), twin), TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist, NameOverrides, FieldTypeOverrides, nil, EnumConstantNames)
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
// non-synthetic counterpart: every real family's bytes must be the same
// generated alone as generated alongside the others. note and override are a
// real collision risk, not a hypothetical one — both describe <nvt>,
// <result>, <task>, <text>, <owner>, <permissions>, and more with identical
// shapes, so a global dedup pass would have to disambiguate one side of each.
func TestRealFamiliesGenerateIndependently(t *testing.T) {
	m := manifestOrFail(t)
	if len(m.Families) < 2 {
		t.Fatalf("The manifest declares %d family; this test needs at least two", len(m.Families))
	}
	cmds, err := LoadGenerationModels(snapshotDir, m.AllCommands())
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}
	together, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist, NameOverrides, FieldTypeOverrides, nil, EnumConstantNames)
	if err != nil {
		t.Fatalf("GenerateAll failed: %v", err)
	}

	byName := make(map[string]schema.Command, len(cmds))
	for _, c := range cmds {
		byName[c.Name] = c
	}

	accounted := 0
	for _, f := range m.Families {
		fam := make([]schema.Command, 0, len(f.Commands))
		for _, name := range f.Commands {
			c, ok := byName[name]
			if !ok {
				t.Fatalf("%s: no snapshot loaded for %s", f.Name, name)
			}
			fam = append(fam, c)
		}
		// Generate directly rather than GenerateAll: GenerateAll's
		// EnumConstantNames staleness check is only meaningful over the whole
		// corpus (a family with no <relation> field legitimately uses none of
		// the operator entries), and this test's isolation claim is about
		// generated bytes, not override-table completeness.
		files, err := Generate(fam, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions,
			CommandWarnings, DefectAllowlist, NameOverrides, FieldTypeOverrides, f.Name, nil, EnumConstantNames, map[string]bool{})
		if err != nil {
			t.Fatalf("Generate(%s alone) failed: %v", f.Name, err)
		}
		alone := map[string]string{}
		if err := mergeFamilyFiles(alone, map[string]string{}, f, files); err != nil {
			t.Fatalf("Generate(%s alone): %v", f.Name, err)
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
// generate.go's initialisms table: without that entry goName would render
// the field "Icalendar" instead of "ICalendar" (RFC 5545's own casing), a
// public API regression. It also proves alone-vs-together consistency, since
// initialisms is package-level and TestRealFamiliesGenerateIndependently
// doesn't check the spelling itself is correct.
func TestScheduleICalendarNamingIsStable(t *testing.T) {
	m := manifestOrFail(t)
	cmds, err := LoadGenerationModels(snapshotDir, m.AllCommands())
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}

	fam, ok := m.FamilyOf("get_schedules")
	if !ok {
		t.Fatal("The manifest no longer declares a schedule family")
	}
	alone, err := GenerateAll(Manifest{Families: []Family{fam}}, cmds, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist, NameOverrides, FieldTypeOverrides, nil, EnumConstantNames)
	if err != nil {
		t.Fatalf("GenerateAll(schedule alone) failed: %v", err)
	}
	together, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist, NameOverrides, FieldTypeOverrides, nil, EnumConstantNames)
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
// generate.go's initialisms table: without those entries goName would render
// get_alerts's fingerprint fields Md5Fingerprint/Sha256Fingerprint, breaking
// consistency with the existing MD5Fingerprint/SHA256Fingerprint spelling in
// commands/get_reports.go, get_scanners.go, and get_tls_certificates.go. Like
// TestScheduleICalendarNamingIsStable, it also pins alone-vs-together
// consistency.
func TestAlertFingerprintNamingIsStable(t *testing.T) {
	m := manifestOrFail(t)
	cmds, err := LoadGenerationModels(snapshotDir, m.AllCommands())
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}

	fam, ok := m.FamilyOf("get_alerts")
	if !ok {
		t.Fatal("The manifest no longer declares an alert family")
	}
	alone, err := GenerateAll(Manifest{Families: []Family{fam}}, cmds, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist, NameOverrides, FieldTypeOverrides, nil, EnumConstantNames)
	if err != nil {
		t.Fatalf("GenerateAll(alert alone) failed: %v", err)
	}
	together, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist, NameOverrides, FieldTypeOverrides, nil, EnumConstantNames)
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
// initialisms table: without it goName would render SubjectDn/IssuerDn
// instead of SubjectDN/IssuerDN, regressing against the casing
// tls_certificate's own adoption locked in (T-9) and commands/get_reports.go
// already uses. Also pins alone-vs-together consistency, per
// TestScheduleICalendarNamingIsStable.
func TestTLSCertificateDNNamingIsStable(t *testing.T) {
	m := manifestOrFail(t)
	cmds, err := LoadGenerationModels(snapshotDir, m.AllCommands())
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}

	fam, ok := m.FamilyOf("get_tls_certificates")
	if !ok {
		t.Fatal("The manifest no longer declares a tls_certificate family")
	}
	alone, err := GenerateAll(Manifest{Families: []Family{fam}}, cmds, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist, NameOverrides, FieldTypeOverrides, nil, EnumConstantNames)
	if err != nil {
		t.Fatalf("GenerateAll(tls_certificate alone) failed: %v", err)
	}
	together, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist, NameOverrides, FieldTypeOverrides, nil, EnumConstantNames)
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

// TestOCIImageNamingIsStable pins the fix for "oci" in generate.go's
// initialisms table (T-16 review, docs/review-2026-08.md): without it
// get_notes's and get_overrides's oci_image elements would render as
// OciImage/NoteOciImage/OverrideOciImage instead of matching the other
// full-caps acronyms (API, CVE, TLS, ...) in the generated surface. Also
// pins alone-vs-together consistency, per TestScheduleICalendarNamingIsStable.
func TestOCIImageNamingIsStable(t *testing.T) {
	m := manifestOrFail(t)
	cmds, err := LoadGenerationModels(snapshotDir, m.AllCommands())
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}
	together, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist, NameOverrides, FieldTypeOverrides, nil, EnumConstantNames)
	if err != nil {
		t.Fatalf("GenerateAll(full manifest) failed: %v", err)
	}

	for _, tt := range []struct {
		familyName string
		fileName   string
		want       []string
		unwanted   []string
	}{
		{"note", "get_notes.go", []string{"OCIImage", "NoteOCIImage"}, []string{"OciImage", "NoteOciImage"}},
		{"override", "get_overrides.go", []string{"OCIImage", "OverrideOCIImage"}, []string{"OciImage", "OverrideOciImage"}},
	} {
		fam, ok := m.FamilyOf("get_" + tt.familyName + "s")
		if !ok {
			t.Fatalf("The manifest no longer declares a %s family", tt.familyName)
		}
		alone, err := GenerateAll(Manifest{Families: []Family{fam}}, cmds, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist, NameOverrides, FieldTypeOverrides, nil, EnumConstantNames)
		if err != nil {
			t.Fatalf("GenerateAll(%s alone) failed: %v", tt.familyName, err)
		}

		aloneSrc, ok := alone[tt.fileName]
		if !ok {
			t.Fatalf("Generating %s alone did not emit %s", tt.familyName, tt.fileName)
		}
		togetherSrc, ok := together[tt.fileName]
		if !ok {
			t.Fatalf("Generating the full manifest did not emit %s", tt.fileName)
		}
		if aloneSrc != togetherSrc {
			t.Errorf("%s differs between generating %s alone and generating the full manifest\n%s",
				tt.fileName, tt.familyName, firstDiff(aloneSrc, togetherSrc))
		}
		for _, want := range tt.want {
			if !strings.Contains(aloneSrc, want) {
				t.Errorf("%s does not spell the field %s", tt.fileName, want)
			}
		}
		for _, unwanted := range tt.unwanted {
			if strings.Contains(aloneSrc, unwanted) {
				t.Errorf("%s spells the field %s instead of the fixed OCI casing", tt.fileName, unwanted)
			}
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
	// One family's own commands, not m.AllCommands(): merging every family
	// into one batch defeats GenerateAll's family-local isolation and can
	// collide structurally-identical types before this test's actual
	// naming-convention check ever runs.
	noteCommands := []string{"create_note", "delete_note", "get_notes", "modify_note"}
	cmds, err := LoadGenerationModels(snapshotDir, noteCommands)
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}
	wrong := Manifest{Families: []Family{{Name: "annotation", Commands: noteCommands}}}
	if _, err := GenerateAll(wrong, cmds, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist, NameOverrides, FieldTypeOverrides, nil, EnumConstantNames); err == nil ||
		!strings.Contains(err.Error(), "naming convention") {
		t.Fatalf("GenerateAll error = %v, want a naming-convention complaint", err)
	}
}

func TestGenerateAllNeedsEverySnapshot(t *testing.T) {
	m := manifestOrFail(t)
	cmds, err := LoadGenerationModels(snapshotDir, []string{"create_note"})
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}
	if _, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist, NameOverrides, FieldTypeOverrides, nil, EnumConstantNames); err == nil ||
		!strings.Contains(err.Error(), "no snapshot loaded for delete_note") {
		t.Fatalf("GenerateAll error = %v, want a missing-snapshot complaint", err)
	}
}
