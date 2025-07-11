package main

import (
	"fmt"
	"log"

	"github.com/brennoo/go-gmp/client"
	"github.com/brennoo/go-gmp/commands"
	"github.com/brennoo/go-gmp/connections"
)

func main() {
	// Connect to GMP over TLS
	gmpConn, err := connections.NewTLSConnection("localhost:9390", true) // true = skip cert verification for demo
	if err != nil {
		log.Fatalf("failed to connect to GMP: %v", err)
	}
	defer gmpConn.Close()

	cli := client.New(gmpConn)

	// Authenticate as admin:admin
	auth := &commands.Authenticate{Credentials: &commands.AuthenticateCredentials{Username: "admin", Password: "admin"}}
	respAuth, err := cli.Authenticate(auth)
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
		KDCs:          &commands.CreateCredentialKDCs{KDC: []string{"kdc1.example.com", "kdc2.example.com"}},
		Key:           &commands.CreateCredentialKey{Phrase: "passphrase", Private: "private-key", Public: "public-key"},
		Privacy:       &commands.CreateCredentialPrivacy{Algorithm: "AES", Password: "snmp-password"},
		AllowInsecure: "0",
	}

	resp, err := cli.CreateCredential(cmd)
	if err != nil {
		log.Fatalf("CreateCredential failed: %v", err)
	}

	fmt.Printf("Credential created! ID: %s, Status: %s, StatusText: %s\n", resp.ID, resp.Status, resp.StatusText)
}
