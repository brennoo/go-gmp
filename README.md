# Greenbone Management Protocol (GMP) Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/brennoo/go-gmp.svg)](https://pkg.go.dev/github.com/brennoo/go-gmp)
[![Go Report Card](https://goreportcard.com/badge/github.com/brennoo/go-gmp)](https://goreportcard.com/report/github.com/brennoo/go-gmp)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/594fb36107f64d72b7bc1680fdac964a)](https://app.codacy.com/gh/brennoo/go-gmp/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

A Go client for the Greenbone Management Protocol (GMP) to interact with Greenbone/OpenVAS vulnerability management servers.

> [!CAUTION]
> Expect breaking changes as the library matures.

## Requirements

- Go 1.23+
- Access to a Greenbone/OpenVAS server with GMP enabled

> [!NOTE]
> An example `docker-compose.yml` is provided in `./hack` for local testing.

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
	// Connect
	conn, err := connections.NewTLSConnection("gmp.example.com:9390", true)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Authenticate
	cli := client.New(conn)
	authResp, err := cli.Authenticate(&commands.Authenticate{
		Credentials: &commands.AuthenticateCredentials{
			Username: "admin",
			Password: "secret",
		},
	})
	if err != nil || authResp.Status != "200" {
		log.Fatalf("auth failed: %v", err)
	}

	// Query
	tasksResp, err := cli.GetTasks(&commands.GetTasks{})
	if err != nil || tasksResp.Status != "200" {
		log.Fatalf("get tasks failed: %v", err)
	}
	for _, task := range tasksResp.Task {
		fmt.Printf("Task: %s (%s)\n", task.Name, task.Status)
	}
}
```

## Core Concepts

**Connection**: Create a `gmp.Connection` (TLS or Unix socket) to handle low-level protocol.  
**Client**: Wrap the connection with `client.New()` to get a `gmp.Client`.  
**Commands**: Each GMP command is a Go struct (e.g., `GetTasks`, `CreateTask`).  
**Responses**: Strongly-typed response structs (e.g., `GetTasksResponse`).  
**Errors**: All methods return `(response, error)`.

## Filtering

Use the `Filter` field on list/get commands:

```go
tasksResp, err := cli.GetTasks(&commands.GetTasks{
	Filter: "status=running sort=name",
})
```

## Pagination

Three approaches available:

| Resource | Basic | Paged | Iterator |
|---------|-------|-------|----------|
| Tasks | `GetTasks` | `GetTasksPaged` | `GetTasksIter` |
| Results | `GetResults` | `GetResultsPaged` | `GetResultsIter` |
| Assets | `GetAssets` | `GetAssetsPaged` | `GetAssetsIter` |
| Targets | `GetTargets` | `GetTargetsPaged` | `GetTargetsIter` |
| Tickets | `GetTickets` | `GetTicketsPaged` | `GetTicketsIter` |
| Port Lists | `GetPortLists` | `GetPortListsPaged` | `GetPortListsIter` |
| Settings | `GetSettings` | `GetSettingsPaged` | `GetSettingsIter` |

### Paged Helpers

```go
// Get page 1 with 10 tasks
tasksResp, err := cli.GetTasksPaged(1, 10)

// Get page 2 with filters
resultsResp, err := cli.GetResultsPaged(2, 5, "severity>5.0", "status=Done")
```

### Iterators

```go
taskIter := cli.GetTasksIter(ctx, 10)
defer taskIter.Close()

for taskIter.Next() {
	task := taskIter.Current()
	fmt.Printf("Task: %s\n", task.Name)
}
if err := taskIter.Err(); err != nil {
	log.Printf("Error: %v", err)
}
```

## Raw XML

Send custom or unsupported commands:

```go
raw := `<get_custom_resource foo="bar"/>`
resp, err := cli.RawXML(raw)
```

## Extending

All GMP commands are available as typed structs matching the protocol XML. Add new commands by following the existing pattern.

## See Also

- [Greenbone GMP Protocol Documentation](https://docs.greenbone.net/API/GMP/gmp-22.6.html)

## Credits

Based on [github.com/filewalkwithme/go-gmp](https://github.com/filewalkwithme/go-gmp). 
