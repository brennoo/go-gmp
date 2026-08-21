// Package commands defines the request and response types for the Greenbone
// Management Protocol, and the error types returned when a command fails.
//
// Every GMP command has a matching pair of structs named after the protocol
// element, for example [Authenticate] / [AuthenticateResponse] or [GetTasks] /
// [GetTasksResponse]. Fields carry `xml` struct tags so a command can be
// marshalled straight to the wire format gvmd expects, and a response
// unmarshalled straight back. These types are consumed by
// [github.com/brennoo/go-gmp/client], which offers one typed method per
// command; most callers never construct or decode them directly.
//
// Responses that carry a GMP status code implement [ResponseWithStatus], which
// [ValidateResponse] uses to turn a failure status into a [GMPError].
//
// # Errors
//
// A call to the client can fail in four distinct ways:
//
//   - [GMPError] — the server accepted and framed the request but reported a
//     failure status (for example authentication failure or a not-found
//     resource). The connection is unaffected; [GMPError.Type] classifies the
//     status.
//   - [TransportError] — the transport did not complete the transaction: the
//     socket-level exchange failed, the response could not be framed as a
//     complete XML document, or the request was refused outright. Check
//     [TransportError.ConnectionUsable] before discarding the connection: a
//     request refused before the socket was touched, such as one made with an
//     already-cancelled context, leaves it usable.
//   - [ConnectionNotInitializedError] — a Connection was used without an
//     underlying socket, which means it was built as a bare struct literal
//     instead of by a constructor. A misuse of the library, not a runtime
//     failure.
//   - [ProtocolError] — the response framed as a complete XML document but
//     failed to decode into the expected type, typically schema or version
//     drift between client and server. The connection is unaffected and
//     remains usable; [ProtocolError.Document] retains a capped copy of the
//     offending response for diagnosis.
//   - [RequestEncodingError] — the command itself could not be marshalled to
//     XML. The connection is never touched.
//
// Every error struct except GMPError implements Unwrap, so errors.Is and
// errors.As reach through them to an underlying cause such as
// context.DeadlineExceeded or a net.Error.
package commands
