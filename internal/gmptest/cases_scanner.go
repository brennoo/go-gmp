package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	RegisterRequestCase("create_scanner", "basic", &commands.CreateScanner{
		Name:  "Corporate OpenVAS scanner",
		Host:  "scanner.internal.example.com",
		Port:  "9391",
		Type:  "22",
		CAPub: "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----",
		Credential: commands.ScannerCredential{
			ID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
		},
	})

	RegisterRequestCase("create_scanner", "with_relay_port", &commands.CreateScanner{
		Name:  "Relayed OpenVAS scanner",
		Host:  "scanner.internal.example.com",
		Port:  "9391",
		Type:  "22",
		CAPub: "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----",
		Credential: commands.ScannerCredential{
			ID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
		},
		RelayHost: "relay.internal.example.com",
		RelayPort: "9392",
	})

	// ModifyScanner.Host/Port/Type are required, non-omitempty strings.
	RegisterRequestCase("modify_scanner", "unset", &commands.ModifyScanner{
		ScannerID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
		Host:      "scanner.internal.example.com",
		Port:      "9391",
		Type:      "22",
	})

	RegisterRequestCase("modify_scanner", "optional_fields_set", &commands.ModifyScanner{
		ScannerID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
		Host:      "scanner.internal.example.com",
		Port:      "9391",
		Type:      "22",
		Name:      ptr("Renamed scanner"),
		Comment:   ptr("Updated during credential rotation"),
		CAPub:     ptr("-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----"),
		Credential: &commands.ScannerCredential{
			ID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
		},
		RelayHost: ptr("relay.internal.example.com"),
		RelayPort: "9392",
	})
}
