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

	// CreateScanner.RelayPort is *int: unset, a nonzero value, and an
	// explicit zero all round-trip distinctly.
	RegisterRequestCase("create_scanner", "with_relay_port", &commands.CreateScanner{
		Name:      "Relayed OpenVAS scanner",
		Host:      "scanner.internal.example.com",
		Port:      9391,
		Type:      "22",
		CAPub:     "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----",
		RelayHost: "relay.internal.example.com",
		RelayPort: commands.Int(9392),
	})

	RegisterRequestCase("create_scanner", "relay_port_zero", &commands.CreateScanner{
		Name:      "Relayed OpenVAS scanner",
		Host:      "scanner.internal.example.com",
		Port:      9391,
		Type:      "22",
		CAPub:     "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----",
		RelayHost: "relay.internal.example.com",
		RelayPort: commands.Int(0),
	})

	// ModifyScanner.Port/RelayPort are *int: unset, a nonzero value, and an
	// explicit zero all round-trip distinctly.
	RegisterRequestCase("modify_scanner", "unset", &commands.ModifyScanner{
		ScannerID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
		Name:      "Renamed scanner",
	})

	RegisterRequestCase("modify_scanner", "port_and_relay_port_set", &commands.ModifyScanner{
		ScannerID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
		Port:      commands.Int(9391),
		RelayPort: commands.Int(9392),
	})

	RegisterRequestCase("modify_scanner", "port_and_relay_port_zero", &commands.ModifyScanner{
		ScannerID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
		Port:      commands.Int(0),
		RelayPort: commands.Int(0),
	})
}
