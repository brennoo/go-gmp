package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// GetVersion has no fields at all: the simplest possible request.
	RegisterRequestCase("get_version", "no_args", &commands.GetVersion{})

	RegisterRequestCase("get_info", "cve_by_id", &commands.GetInfo{
		Type:   "CVE",
		InfoID: "CVE-2024-12345",
	})

	// GetSystemReports.Brief is *bool: unset, true, and an explicit false all
	// round-trip distinctly.
	RegisterRequestCase("get_system_reports", "brief_true", &commands.GetSystemReports{
		Name:  "load",
		Brief: commands.Bool(true),
	})

	RegisterRequestCase("get_system_reports", "brief_unset", &commands.GetSystemReports{
		Name: "load",
	})

	RegisterRequestCase("get_system_reports", "brief_false", &commands.GetSystemReports{
		Name:  "load",
		Brief: commands.Bool(false),
	})

	// GetSystemReports.Duration is *int: unset, a nonzero value, and an
	// explicit zero all round-trip distinctly.
	RegisterRequestCase("get_system_reports", "duration_set", &commands.GetSystemReports{
		Name:     "load",
		Duration: commands.Int(3600),
	})

	RegisterRequestCase("get_system_reports", "duration_zero", &commands.GetSystemReports{
		Name:     "load",
		Duration: commands.Int(0),
	})

	// GetVulns.Trash is *bool: unset, true, and an explicit false all
	// round-trip distinctly.
	RegisterRequestCase("get_vulns", "trash_true", &commands.GetVulns{
		Filter: "rows=-1",
		Trash:  commands.Bool(true),
	})

	RegisterRequestCase("get_vulns", "by_id", &commands.GetVulns{
		VulnID: "1234abcd-0000-1111-2222-333344445555",
	})

	RegisterRequestCase("get_vulns", "trash_false", &commands.GetVulns{
		Filter: "rows=-1",
		Trash:  commands.Bool(false),
	})

	// RunWizard.ReadOnly is *bool: unset, true, and an explicit false all
	// round-trip distinctly.
	RegisterRequestCase("run_wizard", "read_only_true", &commands.RunWizard{
		Name: "quick_first_scan",
		Params: &commands.WizardParams{
			Params: []commands.WizardParam{{Name: "hosts", Value: "192.0.2.0/24"}},
		},
		ReadOnly: commands.Bool(true),
	})

	RegisterRequestCase("run_wizard", "read_only_unset", &commands.RunWizard{
		Name: "quick_first_scan",
		Params: &commands.WizardParams{
			Params: []commands.WizardParam{{Name: "hosts", Value: "192.0.2.0/24"}},
		},
	})

	RegisterRequestCase("run_wizard", "read_only_false", &commands.RunWizard{
		Name: "quick_first_scan",
		Params: &commands.WizardParams{
			Params: []commands.WizardParam{{Name: "hosts", Value: "192.0.2.0/24"}},
		},
		ReadOnly: commands.Bool(false),
	})

	// ModifyLicense.AllowEmpty is *bool: unset, true, and an explicit false
	// all round-trip distinctly.
	RegisterRequestCase("modify_license", "allow_empty_true", &commands.ModifyLicense{
		AllowEmpty: commands.Bool(true),
	})

	RegisterRequestCase("modify_license", "allow_empty_unset", &commands.ModifyLicense{
		File: "BASE64LICENSEDATA==",
	})

	RegisterRequestCase("modify_license", "allow_empty_false", &commands.ModifyLicense{
		AllowEmpty: commands.Bool(false),
		File:       "BASE64LICENSEDATA==",
	})
}
