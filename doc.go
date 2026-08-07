// Package gmp defines the core interfaces for talking to a Greenbone/OpenVAS
// vulnerability scanner over the Greenbone Management Protocol (GMP): [Client]
// and [Connection].
//
// GMP is an XML request/response protocol exposed by gvmd (the Greenbone
// Vulnerability Manager daemon), either over TLS or over a local Unix socket.
// A client authenticates once per connection and then issues commands such as
// get_tasks, create_target or start_task, each with its own request and
// response schema.
//
// This module is organized in three layers:
//
//   - [github.com/brennoo/go-gmp/connections] dials a transport (TLS or Unix
//     socket) and returns a [Connection], which serializes GMP transactions
//     over that socket.
//   - [github.com/brennoo/go-gmp/client] wraps a [Connection] in [Client], a
//     typed method for every GMP command.
//   - [github.com/brennoo/go-gmp/commands] holds the request and response
//     structs those methods marshal and unmarshal, plus the library's error
//     types.
//
// [github.com/brennoo/go-gmp/commands/filtering] and
// [github.com/brennoo/go-gmp/commands/pagination] build on top of [Client] to
// filter and page through list responses.
//
// # Usage
//
// A minimal end-to-end interaction: connect, authenticate, and list tasks.
//
//	ctx := context.Background()
//
//	conn, err := connections.NewTLSConnection(ctx, "gmp.example.com:9390",
//		connections.WithRootCAFile("/path/to/ca.pem"),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer conn.Close()
//
//	cli := client.New(conn)
//
//	_, err = cli.Authenticate(ctx, &commands.Authenticate{
//		Credentials: &commands.AuthenticateCredentials{
//			Username: "admin",
//			Password: "secret",
//		},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	tasksResp, err := cli.GetTasks(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, task := range tasksResp.Task {
//		fmt.Printf("Task: %s (%s)\n", task.Name, task.Status)
//	}
//
// This library is pre-1.0. Breaking changes are expected until 1.0; see
// CHANGELOG.md for what has changed release to release.
package gmp
