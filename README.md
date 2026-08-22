# Greenbone Management Protocol (GMP) Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/brennoo/go-gmp.svg)](https://pkg.go.dev/github.com/brennoo/go-gmp)
[![Go Report Card](https://goreportcard.com/badge/github.com/brennoo/go-gmp)](https://goreportcard.com/report/github.com/brennoo/go-gmp)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/594fb36107f64d72b7bc1680fdac964a)](https://app.codacy.com/gh/brennoo/go-gmp/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

A Go client for the Greenbone Management Protocol (GMP) to interact with Greenbone/OpenVAS vulnerability management servers.

> [!CAUTION]
> Expect breaking changes as the library matures.

## Compatibility

This library is tested against GMP 22.7, as served by gvmd 26.36.1
(`ghcr.io/brennoo/gvm-mirror/gvmd@sha256:cb82c501ca93809405a8d0ed4c21846703466802ce6462536e985625f102b497`,
the digest pinned in `hack/docker-compose.yml`). Other gvmd versions are best-effort: GMP has
historically been additive in the areas this library touches, but only the pinned version above is
verified. Drift between this pin and upstream `gvmd main` is tracked, not silently absorbed — see
[`docs/review-2026-08.md`](docs/review-2026-08.md)'s T-35 entry for the current count of unadopted
newer commands and fields. The agent-controller, OCI-image-target, and web-application-target
surfaces are not modeled; see T-32 in the same document for the evidence required to reintroduce them.

| | |
|---|---|
| GMP protocol | 22.7 |
| gvmd | 26.36.1 |
| gvmd image digest | `sha256:cb82c501ca93809405a8d0ed4c21846703466802ce6462536e985625f102b497` |

See [CHANGELOG.md](CHANGELOG.md) for release history and [CONTRIBUTING.md](CONTRIBUTING.md) for the
versioning policy.

## Requirements

- Go 1.23+
- Access to a Greenbone/OpenVAS server with GMP enabled

> [!NOTE]
> An example `docker-compose.yml` is provided in `./hack` for local testing.

## Quick Start

```go
package main

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

	// Connect, verifying the server certificate against a trusted CA.
	conn, err := connections.NewTLSConnection(ctx, "gmp.example.com:9390",
		connections.WithRootCAFile("/path/to/ca.pem"),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Authenticate
	cli := client.New(conn)
	_, err = cli.Authenticate(ctx, &commands.Authenticate{
		Credentials: &commands.AuthenticateCredentials{
			Username: "admin",
			Password: "secret",
		},
	})
	if err != nil {
		log.Fatalf("auth failed: %v", err)
	}

	// Query
	tasksResp, err := cli.GetTasks(ctx)
	if err != nil {
		log.Fatalf("get tasks failed: %v", err)
	}
	for _, task := range tasksResp.Task {
		fmt.Printf("Task: %s (%s)\n", task.Name, task.Status)
	}
}
```

`NewTLSConnection` verifies the server certificate and requires TLS 1.3 by default. For a server
with a certificate that isn't in the system trust store, pass its CA certificate with
`connections.WithRootCAFile`. For a server that only supports TLS 1.2, opt in explicitly with
`connections.WithMinVersion(tls.VersionTLS12)`.

> [!WARNING]
> `connections.WithInsecureSkipVerify()` disables certificate verification entirely and must
> only be used for local development against a server with a self-signed certificate — never in
> production:
>
> ```go
> conn, err := connections.NewTLSConnection(ctx, "localhost:9390",
> 	connections.WithInsecureSkipVerify(),
> )
> ```

## Filtering

The library provides two filtering approaches:

### String Filters

```go
tasksResp, err := cli.GetTasks(ctx, "status=running", "sort=name")
```

### Functional Options

```go
import "github.com/brennoo/go-gmp/commands/filtering"

tasksResp, err := cli.GetTasks(ctx,
	filtering.WithStatus(filtering.StatusRunning),
	filtering.WithSort(filtering.SortByName, filtering.SortAsc),
)
```

## Pagination

The library provides two approaches for handling large datasets:

| Resource | One-shot | Iterator |
|---------|----------|----------|
| Tasks | `GetTasks` | `Tasks` |
| Results | `GetResults` | `Results` |
| Assets | `GetAssets` | `Assets` |
| Targets | `GetTargets` | `Targets` |
| Tickets | `GetTickets` | `Tickets` |
| Port Lists | `GetPortLists` | `PortLists` |
| Settings | `GetSettings` | `Settings` |

### One-shot Retrieval

```go
// Get first page with default limit
tasksResp, err := cli.GetTasks(ctx)

// Get with custom limit and filters
tasksResp, err := cli.GetTasks(ctx, "rows=50", "status=running")
```

### Iterators

```go
taskIter := cli.Tasks(ctx, 10)
defer taskIter.Close()

for taskIter.Next() {
	task := taskIter.Current()
	fmt.Printf("Task: %s\n", task.Name)
}
if err := taskIter.Err(); err != nil {
	log.Printf("Error: %v", err)
}
```

## Error Handling

The library provides structured error handling using Go's `errors.As` and `errors.Is`. Every
error returned by a `Client` or `Connection` method, or by an iterator's `Err()`, is one of five
kinds:

| Type | Means | Connection usable after? | Retry? |
|------|-------|---------------------------|--------|
| `*commands.GMPError` | The server accepted and framed the request but reported a failure status | Yes | Depends on the status (`gmpErr.Type`) |
| `*commands.TransportError` | The transport did not complete the transaction: a socket/framing failure, or the request was refused outright | **Check `ConnectionUsable`** — false means reconnect; true means the request was refused before the socket was touched and the connection is fine | Only after reconnecting, if unusable |
| `*commands.ProtocolError` | The response framed cleanly but failed to decode (schema/version drift) | Yes | No — deterministic, retrying won't change the outcome |
| `*commands.RequestEncodingError` | The command itself could not be marshalled to XML | Yes | No — this is a caller bug (bad field value/type) |
| `*commands.ConnectionNotInitializedError` | A `Connection` was built as a bare struct literal instead of by a constructor, so it has no socket | n/a — it never had one | No — this is a caller bug |

### GMP status classification (`GMPError.Type`)

`GMPError.Type` further classifies **server-reported statuses only** — it says nothing about
transport or protocol errors, which have no status to classify:

| Error Type | Description | HTTP-like Status |
|------------|-------------|-------------------|
| `ErrorTypeAuthentication` | Authentication failed | 401 |
| `ErrorTypeNotFound` | Resource not found | 404 |
| `ErrorTypePermission` | Insufficient permissions | 403 |
| `ErrorTypeInvalidRequest` | Invalid request format | 400 |
| `ErrorTypeServerError` | Server internal error | 500+ |

### Matching errors

Switch on the kinds with `errors.As`, and only look at `gmpErr.Type` inside the `GMPError` arm.
`ConnectionNotInitializedError` is omitted below: it can only result from constructing a `Connection`
incorrectly, so it is a bug to fix rather than a condition to handle.

```go
import (
	"errors"
	"log"

	"github.com/brennoo/go-gmp/commands"
)

tasksResp, err := cli.GetTasks(ctx)
if err != nil {
	var gmpErr *commands.GMPError
	var transportErr *commands.TransportError
	var protocolErr *commands.ProtocolError
	var encodingErr *commands.RequestEncodingError
	switch {
	case errors.As(err, &gmpErr):
		switch gmpErr.Type {
		case commands.ErrorTypeAuthentication:
			log.Fatal("Authentication failed:", gmpErr.StatusText)
		case commands.ErrorTypeNotFound:
			log.Println("No tasks found, continuing...")
			return nil // Not necessarily an error
		default:
			log.Printf("GMP error: %v", gmpErr)
		}
	case errors.As(err, &transportErr):
		if transportErr.ConnectionUsable {
			// Refused before the socket was touched — usually an already-cancelled
			// context. Connections are shared, so do not discard this one.
			log.Println("Request refused, connection still usable:", transportErr)
		} else {
			log.Println("Connection is gone, reconnect:", transportErr)
			// Dial a new Connection; this one cannot be revived.
		}
	case errors.As(err, &protocolErr):
		log.Printf("Protocol error (schema drift?), connection still usable: %v", protocolErr)
	case errors.As(err, &encodingErr):
		log.Printf("Command failed to encode, this is a bug in the caller: %v", encodingErr)
	default:
		log.Printf("Unexpected error: %v", err)
	}
}
```

### Matching through the wrap chain

`GMPError` is a leaf error and does not implement `Unwrap` — the server's status is the entire
failure, with nothing underneath. `TransportError`, `ProtocolError` and `RequestEncodingError` all
implement `Unwrap`, so `errors.Is` and `errors.As` reach through them to whatever caused the
failure:

```go
if errors.Is(err, context.DeadlineExceeded) {
	// The context passed to the call timed out.
}
if errors.Is(err, context.Canceled) {
	// The context passed to the call was cancelled.
}
if errors.Is(err, connections.ErrConnectionBroken) {
	// A previous transaction poisoned this Connection; it must be replaced.
}

var netErr net.Error
if errors.As(err, &netErr) && netErr.Timeout() {
	// The underlying socket operation itself timed out.
}
```

### Iterators

An iterator's `Err()` method returns the error that stopped iteration, once `Next()` returns
`false`:

```go
taskIter := cli.Tasks(ctx, 10)
defer taskIter.Close()

for taskIter.Next() {
	task := taskIter.Current()
	fmt.Printf("Task: %s\n", task.Name)
}
if err := taskIter.Err(); err != nil {
	log.Printf("Error: %v", err)
}
```

`Err()` returns the same error kinds described above, so the matching pattern is identical
to a one-shot call.

## Raw XML

Send custom or unsupported commands:

```go
raw := `<get_custom_resource foo="bar"/>`
resp, err := cli.RawXML(ctx, raw)
```

## Extending

All GMP commands are available as typed structs matching the protocol XML. Add new commands by following the existing pattern.

## See Also

- [Greenbone GMP Protocol Documentation](https://docs.greenbone.net/API/GMP/gmp-22.7.html)

## Credits

Based on [github.com/filewalkwithme/go-gmp](https://github.com/filewalkwithme/go-gmp).
