package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// CreateCredential.AllowInsecure is *bool: unset, true, and an explicit
	// false all round-trip distinctly.
	RegisterRequestCase("create_credential", "allow_insecure_true", &commands.CreateCredential{
		Name:          "Legacy SNMP credential",
		Comment:       "Needed for older scan targets that only speak SNMPv1",
		Type:          "up",
		Login:         "svc-scan",
		Password:      "s3cr3t-password",
		AllowInsecure: commands.Bool(true),
	})

	RegisterRequestCase("create_credential", "allow_insecure_unset", &commands.CreateCredential{
		Name:     "Standard login credential",
		Type:     "up",
		Login:    "svc-scan",
		Password: "s3cr3t-password",
	})

	RegisterRequestCase("create_credential", "allow_insecure_false", &commands.CreateCredential{
		Name:          "Hardened login credential",
		Type:          "up",
		Login:         "svc-scan",
		Password:      "s3cr3t-password",
		AllowInsecure: commands.Bool(false),
	})

	// create_credential.Kdc/.Kdcs/.Key/.Privacy: the shared CredentialKdcs/
	// CredentialKey/CredentialPrivacy types, and Privacy's required
	// (non-omitempty) Algorithm/Password fields.
	RegisterRequestCase("create_credential", "kdc_kdcs_key_privacy", &commands.CreateCredential{
		Name: "cred-1",
		Type: "krb5",
		Kdc:  "kdc0.example.com",
		Kdcs: &commands.CredentialKdcs{Kdc: []string{"kdc1.example.com", "kdc2.example.com"}},
		Key:  &commands.CredentialKey{Phrase: "phrase1", Private: "priv1", Public: "pub1"},
		Privacy: &commands.CredentialPrivacy{
			Algorithm: "aes",
			Password:  "privpass", // nolint:gosec // G101: synthetic fixture value, not a real secret
		},
	})

	RegisterRequestCase("create_credential", "privacy_zero_valued", &commands.CreateCredential{
		Name:    "cred-2",
		Type:    "snmp",
		Privacy: &commands.CredentialPrivacy{},
	})

	// GetCredentials.Details/Scanners/Trash/Targets are *bool: unset, true,
	// and an explicit false all round-trip distinctly.
	RegisterRequestCase("get_credentials", "all_details", &commands.GetCredentials{
		Filter:   "rows=-1",
		Details:  commands.Bool(true),
		Scanners: commands.Bool(true),
		Trash:    commands.Bool(true),
		Targets:  commands.Bool(true),
	})

	RegisterRequestCase("get_credentials", "by_id", &commands.GetCredentials{ // nolint:gosec // G101 fires on the UUID; it is a resource id, not a secret
		CredentialID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("get_credentials", "all_details_false", &commands.GetCredentials{
		Filter:   "rows=-1",
		Details:  commands.Bool(false),
		Scanners: commands.Bool(false),
		Trash:    commands.Bool(false),
		Targets:  commands.Bool(false),
	})

	RegisterRequestCase("get_credentials", "filt_id", &commands.GetCredentials{
		FiltID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
	})

	// ModifyCredential.AllowInsecure is *bool: unset, true, and an explicit
	// false all round-trip distinctly.
	RegisterRequestCase("modify_credential", "allow_insecure_true", &commands.ModifyCredential{ // nolint:gosec // G101 fires on the UUID; it is a resource id, not a secret
		CredentialID:  "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Comment:       ptr("Now permitted to use an insecure algorithm"),
		AllowInsecure: commands.Bool(true),
	})

	RegisterRequestCase("modify_credential", "allow_insecure_unset", &commands.ModifyCredential{ // nolint:gosec // synthetic fixture values: a resource id and a deliberately fake password
		CredentialID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Comment:      ptr("Rotate the password"),
		Password:     "n3w-s3cr3t",
	})

	RegisterRequestCase("modify_credential", "allow_insecure_false", &commands.ModifyCredential{ // nolint:gosec // G101 fires on the UUID; it is a resource id, not a secret
		CredentialID:  "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Comment:       ptr("No longer permitted to use an insecure algorithm"),
		AllowInsecure: commands.Bool(false),
	})

	// ModifyCredential.Comment is *string: gvmd distinguishes an omitted
	// comment (unchanged) from an explicit empty one (clears it). These pin
	// the omitted/explicit-empty/set wire shapes side by side.
	RegisterRequestCase("modify_credential", "all_unset", &commands.ModifyCredential{ // nolint:gosec // G101 fires on the UUID; it is a resource id, not a secret
		CredentialID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("modify_credential", "comment_empty", &commands.ModifyCredential{ // nolint:gosec // G101 fires on the UUID; it is a resource id, not a secret
		CredentialID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Comment:      ptr(""),
	})

	RegisterRequestCase("modify_credential", "comment_set", &commands.ModifyCredential{ // nolint:gosec // G101 fires on the UUID; it is a resource id, not a secret
		CredentialID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Comment:      ptr("Rotating credentials ahead of the audit."),
	})

	// ModifyCredential.Kdc/.Kdcs/.Key/.Privacy: same shared types as
	// create_credential above.
	RegisterRequestCase("modify_credential", "kdc_kdcs_key_privacy", &commands.ModifyCredential{ // nolint:gosec // G101 fires on the UUID; it is a resource id, not a secret
		CredentialID: "cred-1",
		Kdc:          "kdc0.example.com",
		Kdcs:         &commands.CredentialKdcs{Kdc: []string{"kdc1.example.com", "kdc2.example.com"}},
		Key:          &commands.CredentialKey{Phrase: "phrase1", Private: "priv1", Public: "pub1"},
		Privacy: &commands.CredentialPrivacy{
			Algorithm: "aes",
			Password:  "privpass", // nolint:gosec // G101: synthetic fixture value, not a real secret
		},
	})
}
