package main

//nolint:errcheck // Ignoring Close errors in example

import (
	"context"
	"fmt"
	"log"

	"github.com/brennoo/go-gmp/client"
	"github.com/brennoo/go-gmp/commands"
	"github.com/brennoo/go-gmp/connections"
)

func main() {
	ctx := context.Background()

	// Connect to GMP over TLS.
	// WithInsecureSkipVerify disables certificate verification and must not be used in
	// production; use connections.WithRootCAFile with a trusted CA instead.
	gmpConn, err := connections.NewTLSConnection(ctx, "localhost:9390", connections.WithInsecureSkipVerify())
	if err != nil {
		log.Fatalf("failed to connect to GMP: %v", err)
	}
	defer gmpConn.Close() //nolint:errcheck // Ignoring close error

	cli := client.New(gmpConn)

	// Authenticate as admin:admin
	auth := &commands.Authenticate{Credentials: &commands.AuthenticateCredentials{Username: "admin", Password: "admin"}}
	respAuth, err := cli.Authenticate(ctx, auth)
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}
	if respAuth.Status != "200" {
		log.Fatalf("Authentication failed: %s - %s", respAuth.Status, respAuth.StatusText)
	}
	fmt.Println("Authenticated successfully!")

	// Create a credential
	cmd := &commands.CreateCredential{
		Name:          "Example Credential",
		Login:         "username",
		Password:      "password",
		Type:          "ssh",
		Kdcs:          &commands.CredentialKdcs{Kdc: []string{"kdc1.example.com", "kdc2.example.com"}},
		Key:           &commands.CredentialKey{Phrase: "passphrase", Private: "private-key", Public: "public-key"},
		Privacy:       &commands.CredentialPrivacy{Algorithm: "AES", Password: "snmp-password"},
		AllowInsecure: commands.Bool(false),
	}

	resp, err := cli.CreateCredential(ctx, cmd)
	if err != nil {
		log.Fatalf("CreateCredential failed: %v", err)
	}

	fmt.Printf("Credential created! ID: %s, Status: %s, StatusText: %s\n", resp.ID, resp.Status, resp.StatusText)
}
