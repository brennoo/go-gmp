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

	// ModifyCredential.AllowInsecure is *bool: unset, true, and an explicit
	// false all round-trip distinctly.
	RegisterRequestCase("modify_credential", "allow_insecure_true", &commands.ModifyCredential{ // nolint:gosec // G101 fires on the UUID; it is a resource id, not a secret
		CredentialID:  "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Comment:       "Now permitted to use an insecure algorithm",
		AllowInsecure: commands.Bool(true),
	})

	RegisterRequestCase("modify_credential", "allow_insecure_unset", &commands.ModifyCredential{ // nolint:gosec // synthetic fixture values: a resource id and a deliberately fake password
		CredentialID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Comment:      "Rotate the password",
		Password:     "n3w-s3cr3t",
	})

	RegisterRequestCase("modify_credential", "allow_insecure_false", &commands.ModifyCredential{ // nolint:gosec // G101 fires on the UUID; it is a resource id, not a secret
		CredentialID:  "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Comment:       "No longer permitted to use an insecure algorithm",
		AllowInsecure: commands.Bool(false),
	})
}
