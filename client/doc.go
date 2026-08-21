// Package client implements [github.com/brennoo/go-gmp.Client] on top of a
// [github.com/brennoo/go-gmp.Connection].
//
// [New] wraps a connection and returns a ready-to-use client:
//
//	conn, err := connections.NewTLSConnection(ctx, "gmp.example.com:9390")
//	if err != nil {
//		log.Fatal(err)
//	}
//	cli := client.New(conn)
//
// Every method takes a context.Context as its first parameter and maps
// directly onto one GMP command: it marshals a *commands.XxxCommand into XML,
// sends it over the connection, and unmarshals the response into a
// *commands.XxxResponse. List-returning methods such as GetTasks or GetAssets
// additionally accept variadic [github.com/brennoo/go-gmp/commands/filtering.FilterArg]
// values to filter and sort results, and a parallel set of iterator methods
// (Tasks, Results, Assets, ...) page through large result sets — see
// [github.com/brennoo/go-gmp/commands/pagination] for how the iterators work.
//
// A client has no state of its own beyond the connection it was built with:
// authentication, retry and reconnect logic are the caller's responsibility.
package client
