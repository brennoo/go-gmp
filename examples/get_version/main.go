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

	// WithInsecureSkipVerify disables certificate verification and must not be used in
	// production; use connections.WithRootCAFile with a trusted CA instead.
	conn, err := connections.NewTLSConnection(ctx, "localhost:9390", connections.WithInsecureSkipVerify())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close() //nolint:errcheck // Ignoring close error

	cli := client.New(conn)
	resp, err := cli.GetVersion(ctx, &commands.GetVersion{})
	if err != nil {
		log.Fatalf("get version failed: %v", err)
	}
	fmt.Printf("GMP Version: %s\n", resp.Version)
}
