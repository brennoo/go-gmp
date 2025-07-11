package main

import (
	"fmt"
	"log"

	"github.com/brennoo/go-gmp/client"
	"github.com/brennoo/go-gmp/commands"
	"github.com/brennoo/go-gmp/connections"
)

func main() {
	conn, err := connections.NewTLSConnection("localhost:9390", true)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	cli := client.New(conn)
	resp, err := cli.GetVersion(&commands.GetVersion{})
	if err != nil {
		log.Fatalf("get version failed: %v", err)
	}
	fmt.Printf("GMP Version: %s\n", resp.Version)
}
