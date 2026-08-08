package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	RegisterRequestCase("create_scanner", "basic", &commands.CreateScanner{
		Name:  "Corporate OpenVAS scanner",
		Host:  "scanner.internal.example.com",
		Port:  9391,
		Type:  "22",
		CAPub: "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----",
		Credential: &commands.ScannerCredential{
			ID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
		},
	})
}
