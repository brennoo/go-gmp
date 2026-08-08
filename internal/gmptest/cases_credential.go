package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// CreateCredential.AllowInsecure: non-pointer bool, omitempty.
	RegisterRequestCase("create_credential", "allow_insecure_true", &commands.CreateCredential{
		Name:          "Legacy SNMP credential",
		Comment:       "Needed for older scan targets that only speak SNMPv1",
		Type:          "up",
		Login:         "svc-scan",
		Password:      "s3cr3t-password",
		AllowInsecure: true,
	})

	RegisterRequestCase("create_credential", "allow_insecure_unset", &commands.CreateCredential{
		Name:     "Standard login credential",
		Type:     "up",
		Login:    "svc-scan",
		Password: "s3cr3t-password",
	})

	// GetCredentials.Details/Scanners/Trash/Targets: non-pointer bools, omitempty.
	RegisterRequestCase("get_credentials", "all_details", &commands.GetCredentials{
		Filter:   "rows=-1",
		Details:  true,
		Scanners: true,
		Trash:    true,
		Targets:  true,
	})

	RegisterRequestCase("get_credentials", "by_id", &commands.GetCredentials{ // nolint:gosec // G101 fires on the UUID; it is a resource id, not a secret
		CredentialID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
	})

	// ModifyCredential.AllowInsecure: non-pointer bool, omitempty.
	RegisterRequestCase("modify_credential", "allow_insecure_true", &commands.ModifyCredential{ // nolint:gosec // G101 fires on the UUID; it is a resource id, not a secret
		CredentialID:  "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Comment:       "Now permitted to use an insecure algorithm",
		AllowInsecure: true,
	})

	RegisterRequestCase("modify_credential", "allow_insecure_unset", &commands.ModifyCredential{ // nolint:gosec // synthetic fixture values: a resource id and a deliberately fake password
		CredentialID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Comment:      "Rotate the password",
		Password:     "n3w-s3cr3t",
	})
}
