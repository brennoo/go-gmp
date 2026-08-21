package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// CreateTLSCertificate.Trust is *GMPBool: unset, true, and an explicit
	// false all round-trip distinctly. Certificate is schema-required
	// (non-pointer, no omitempty), so every case gives it a populated value.
	RegisterRequestCase("create_tls_certificate", "trust_true", &commands.CreateTLSCertificate{
		Name:        "Example Certificate",
		Comment:     "Imported from the primary scan target",
		Certificate: "MIIDNjCCAp+gAwIBAgIBATANBgkqhkiG9w0BAQQFADCBqTELMAkGA1UEBhMCVVM=",
		Trust:       commands.Bool(true),
	})

	RegisterRequestCase("create_tls_certificate", "trust_unset", &commands.CreateTLSCertificate{
		Certificate: "MIIDNjCCAp+gAwIBAgIBATANBgkqhkiG9w0BAQQFADCBqTELMAkGA1UEBhMCVVM=",
	})

	RegisterRequestCase("create_tls_certificate", "trust_false", &commands.CreateTLSCertificate{
		Certificate: "MIIDNjCCAp+gAwIBAgIBATANBgkqhkiG9w0BAQQFADCBqTELMAkGA1UEBhMCVVM=",
		Trust:       commands.Bool(false),
	})

	// GetTLSCertificates.IncludeCertificateData is *GMPBool: unset, true, and
	// an explicit false all round-trip distinctly. Moved here from
	// cases_misc.go now that tls_certificate has its own case file.
	RegisterRequestCase("get_tls_certificates", "include_data_true", &commands.GetTLSCertificates{
		Filter:                 "rows=-1",
		IncludeCertificateData: commands.Bool(true),
	})

	RegisterRequestCase("get_tls_certificates", "include_data_unset", &commands.GetTLSCertificates{
		TLSCertificateID: "7c8d9e0f-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
	})

	RegisterRequestCase("get_tls_certificates", "include_data_false", &commands.GetTLSCertificates{
		Filter:                 "rows=-1",
		IncludeCertificateData: commands.Bool(false),
	})

	// ModifyTLSCertificate.Name/.Comment are *string (PointerOverrides): a
	// live test against gvmd 26.36.1 proved omitting the element leaves the
	// stored value unchanged while sending it explicitly empty clears it. See
	// docs/review-2026-08.md's `tls_certificate adopted` section. These three
	// cases pin all three wire shapes side by side.
	RegisterRequestCase("modify_tls_certificate", "name_comment_unset", &commands.ModifyTLSCertificate{
		TLSCertificateID: "8a925978-59d0-494b-a837-40b271569727",
	})

	RegisterRequestCase("modify_tls_certificate", "name_comment_empty", &commands.ModifyTLSCertificate{
		TLSCertificateID: "8a925978-59d0-494b-a837-40b271569727",
		Name:             ptr(""),
		Comment:          ptr(""),
	})

	RegisterRequestCase("modify_tls_certificate", "name_comment_set", &commands.ModifyTLSCertificate{
		TLSCertificateID: "8a925978-59d0-494b-a837-40b271569727",
		Name:             ptr("Renamed certificate"),
		Comment:          ptr("Updated comment"),
	})

	// ModifyTLSCertificate.Trust is *GMPBool: unset, true, and an explicit
	// false all round-trip distinctly. Comment varies alongside it so each
	// case's golden differs from the name/comment cases above.
	RegisterRequestCase("modify_tls_certificate", "trust_true", &commands.ModifyTLSCertificate{
		TLSCertificateID: "8a925978-59d0-494b-a837-40b271569727",
		Comment:          ptr("Marked as trusted"),
		Trust:            commands.Bool(true),
	})

	RegisterRequestCase("modify_tls_certificate", "trust_unset", &commands.ModifyTLSCertificate{
		TLSCertificateID: "8a925978-59d0-494b-a837-40b271569727",
		Comment:          ptr("Trust decision deferred"),
	})

	RegisterRequestCase("modify_tls_certificate", "trust_false", &commands.ModifyTLSCertificate{
		TLSCertificateID: "8a925978-59d0-494b-a837-40b271569727",
		Comment:          ptr("Marked as untrusted"),
		Trust:            commands.Bool(false),
	})
}
