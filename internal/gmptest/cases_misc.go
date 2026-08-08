package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// GetVersion has no fields at all: the simplest possible request.
	RegisterRequestCase("get_version", "no_args", &commands.GetVersion{})

	RegisterRequestCase("get_info", "cve_by_id", &commands.GetInfo{
		Type:   "CVE",
		InfoID: "CVE-2024-12345",
	})

	// GetSystemReports.Brief: non-pointer bool, omitempty.
	RegisterRequestCase("get_system_reports", "brief_true", &commands.GetSystemReports{
		Name:  "load",
		Brief: true,
	})

	RegisterRequestCase("get_system_reports", "brief_unset", &commands.GetSystemReports{
		Name: "load",
	})

	// GetTLSCertificates.IncludeCertificateData: non-pointer bool, omitempty.
	RegisterRequestCase("get_tls_certificates", "include_data_true", &commands.GetTLSCertificates{
		Filter:                 "rows=-1",
		IncludeCertificateData: true,
	})

	RegisterRequestCase("get_tls_certificates", "include_data_unset", &commands.GetTLSCertificates{
		TLSCertificateID: "7c8d9e0f-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
	})

	// GetVulns.Trash: non-pointer bool, omitempty.
	RegisterRequestCase("get_vulns", "trash_true", &commands.GetVulns{
		Filter: "rows=-1",
		Trash:  true,
	})

	RegisterRequestCase("get_vulns", "by_id", &commands.GetVulns{
		VulnID: "1234abcd-0000-1111-2222-333344445555",
	})

	// RunWizard.ReadOnly: non-pointer bool, omitempty attr.
	RegisterRequestCase("run_wizard", "read_only_true", &commands.RunWizard{
		Name: "quick_first_scan",
		Params: &commands.WizardParams{
			Params: []commands.WizardParam{{Name: "hosts", Value: "192.0.2.0/24"}},
		},
		ReadOnly: true,
	})

	RegisterRequestCase("run_wizard", "read_only_unset", &commands.RunWizard{
		Name: "quick_first_scan",
		Params: &commands.WizardParams{
			Params: []commands.WizardParam{{Name: "hosts", Value: "192.0.2.0/24"}},
		},
	})

	// ModifyLicense.AllowEmpty: non-pointer bool, omitempty attr.
	RegisterRequestCase("modify_license", "allow_empty_true", &commands.ModifyLicense{
		AllowEmpty: true,
	})

	RegisterRequestCase("modify_license", "allow_empty_unset", &commands.ModifyLicense{
		File: "BASE64LICENSEDATA==",
	})
}
