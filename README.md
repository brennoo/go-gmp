# Greenbone Management Protocol (GMP) Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/brennoo/go-gmp.svg)](https://pkg.go.dev/github.com/brennoo/go-gmp)
[![Go Report Card](https://goreportcard.com/badge/github.com/brennoo/go-gmp)](https://goreportcard.com/report/github.com/brennoo/go-gmp)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/594fb36107f64d72b7bc1680fdac964a)](https://app.codacy.com/gh/brennoo/go-gmp/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

A modern (Go) client for the Greenbone Management Protocol (GMP), used to interact with Greenbone/OpenVAS vulnerability management servers.

> [!CAUTION]
> Expect some breaking changes/refactors as the library matures.

## Features

- Full GMP protocol coverage (commands for tasks, targets, reports, assets, credentials, etc.)
- Strongly-typed request/response structs matching the protocol XML
- Idiomatic Go API: simple, consistent, and easy to extend
- Supports filtering, pagination, and error handling as per GMP
- Works with both TLS and Unix socket connections

## Requirements

- Go 1.23+
- Access to a running Greenbone/OpenVAS server with GMP enabled
- Access to the GMP server (typically via TLS or Unix socket)

---

> [!NOTE]
> An example `docker-compose.yml` is provided in the `./hack` directory. You can use it to spin up a full GMP test environment with TLS enabled for local development and testing.

---

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/brennoo/go-gmp/client"
	"github.com/brennoo/go-gmp/commands"
	"github.com/brennoo/go-gmp/connections"
)

func main() {
	// Create a TLS connection (set insecure to true for self-signed certs)
	conn, err := connections.NewTLSConnection("gmp.example.com:9390", true)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	cli := client.New(conn)

	// Authenticate
	auth := &commands.Authenticate{
		Credentials: &commands.AuthenticateCredentials{
			Username: "admin",
			Password: "secret",
		},
	}
	authResp, err := cli.Authenticate(auth)
	if err != nil {
		log.Fatalf("auth failed: %v", err)
	}
	if authResp.Status != "200" {
		log.Fatalf("auth failed: %s (%s)", authResp.Status, authResp.StatusText)
	}
	fmt.Printf("Authenticated! Role: %s, Timezone: %s\n", authResp.Role, authResp.Timezone)

	// Basic query: get tasks
	tasksResp, err := cli.GetTasks(&commands.GetTasks{})
	if err != nil {
		log.Fatalf("get tasks failed: %v", err)
	}
	if tasksResp.Status != "200" {
		log.Fatalf("get tasks failed: %s (%s)", tasksResp.Status, tasksResp.StatusText)
	}
	for _, task := range tasksResp.Task {
		fmt.Printf("Task: %s (ID: %s, Status: %s)\n", task.Name, task.ID, task.Status)
	}
}
```

## How It Works

- **Connection**: You create a `gmp.Connection` (TLS or Unix socket). This handles the low-level protocol.
- **Client**: Pass the connection to `client.New()` to get a `gmp.Client`.
- **Commands**: Each GMP command is a Go struct (e.g., `GetTasks`).
- **Responses**: Each command returns a strongly-typed response struct (e.g., `GetTasksResponse`).
- **Error Handling**: All client methods return a response and an error. Protocol errors are mapped to Go errors.

## Filtering

GMP supports powerful filtering for most list/get commands. You can use the `Filter` field:

```go
// Get only running tasks, sorted by name
tasksResp, err := cli.GetTasks(&commands.GetTasks{
	Filter: "status=running sort=name",
})
if err != nil {
	log.Fatalf("get tasks failed: %v", err)
}
if tasksResp.Status != "200" {
	log.Fatalf("get tasks failed: %s (%s)", tasksResp.Status, tasksResp.StatusText)
}
for _, task := range tasksResp.Task {
	fmt.Println(task.Name, task.Status)
}
```

## Pagination

> [!NOTE]
> By default, the GMP server returns 10 rows per page if you do not specify the `rows` filter.

For large result sets, use `Filter` with `rows` and `first`:

```go
// Get the first 10 results for a task
resultsResp, err := cli.GetResults(&commands.GetResults{
	TaskID: "<task-id>",
	Filter: "rows=10 first=1",
})
if err != nil {
	log.Fatalf("get results failed: %v", err)
}
if resultsResp.Status != "200" {
	log.Fatalf("get results failed: %s (%s)", resultsResp.Status, resultsResp.StatusText)
}
for _, result := range resultsResp.Results {
	fmt.Println(result.ID, result.Name)
}
```

## Error Handling

All client methods return a response and an error. If the server returns a protocol error, it is mapped to a Go error:

```go
resp, err := cli.DeleteTask(&commands.DeleteTask{TaskID: "bad-id"})
if err != nil {
	log.Printf("delete failed: %v", err)
}
if resp.Status != "200" {
	log.Printf("delete failed: %s (%s)", resp.Status, resp.StatusText)
}
```

## RawXML: Custom/Unsupported Requests

If you need to send a custom or undocumented GMP command, you can use the `RawXML` method to send raw XML and receive the raw XML response:

```go
raw := `<get_custom_resource foo="bar"/>`
resp, err := cli.RawXML(raw)
if err != nil {
	log.Fatalf("raw xml failed: %v", err)
}
fmt.Println("Raw response:", resp)
```

This is useful for experimenting with custom or not-yet-supported GMP commands.

## Extending the client

- All major GMP commands are available.
- Each command/response struct is documented and matches the protocol XML.
- You can add new commands by following the same pattern.

## Examples

A set of example programs is provided in the [examples/](./examples/) directory

## See Also

- [Greenbone GMP Protocol Documentation](https://docs.greenbone.net/API/GMP/gmp-22.6.html)
- [Examples](./examples/)

## Credits

This package was based on [github.com/filewalkwithme/go-gmp](https://github.com/filewalkwithme/go-gmp). 
