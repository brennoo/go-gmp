package implementation

import (
	"bufio"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/brennoo/go-gmp/commands"
)

// connState is the lifecycle state of a Connection. The zero value is
// stateOpen, so New has no lifecycle state to initialize and a bare literal
// reports ErrNoRawConn rather than panicking. Transitions are one-way.
type connState int

const (
	// stateOpen means the socket is usable and message-boundary position is known.
	stateOpen connState = iota
	// stateBroken means a transaction failed in a way that leaves the stream
	// position ambiguous relative to message boundaries. Reusing the socket
	// could silently cross-read another response, so it is permanently refused.
	stateBroken
	// stateClosed means Close was called. Terminal: it always outranks broken.
	stateClosed
)

// Sentinel errors describing why a Connection refused to run a transaction.
// They are returned wrapped in *commands.TransportError, so errors.Is works.
var (
	// ErrConnectionClosed is returned once Close has been called.
	ErrConnectionClosed = errors.New("connection closed")
	// ErrConnectionBroken is returned once a transaction has left the stream
	// position ambiguous (see stateBroken).
	ErrConnectionBroken = errors.New("connection broken")
	// ErrNoRawConn is returned when a Connection is used without a socket, which
	// only happens if one is built as a bare literal instead of by New.
	ErrNoRawConn = errors.New("connection has no underlying socket")
)

// Connection is a GMP transport over a single stream socket.
//
// GMP is a strict request/response protocol with no correlation IDs, so a
// complete write-then-read transaction must occupy the socket exclusively.
// Connection therefore serializes transactions: it is safe for concurrent use,
// but not parallel — callers wait. Use one Connection per goroutine (each with
// its own authentication) if genuine parallelism is needed.
//
// Waiting callers are not served in FIFO order: admission is a semaphore, not a
// queue, with no starvation bound. Callers needing fairness must arrange it
// above this type.
//
// Two synchronization primitives are used, deliberately separate:
//
//   - admission, a size-one buffered channel semaphore, serializes whole
//     transactions. A channel not a mutex, because a mutex cannot be acquired
//     under a context, so a waiting caller could not honor its own deadline.
//   - mu guards state, rawConn and reader. Separate from admission so Close can
//     act while another goroutine holds the transaction, and never held across
//     blocking socket I/O.
//
// A Connection must not be copied after first use.
type Connection struct {
	// admission is lazily created by sem(); see the note on admissionOnce.
	admission chan struct{}
	// admissionOnce keeps the zero value usable: Connection is exported, so a
	// bare literal remains constructible and a nil chan would block forever. It
	// also gives every sem() caller the happens-before edge to observe the same
	// channel safely.
	admissionOnce sync.Once

	mu      sync.Mutex
	state   connState
	rawConn net.Conn
	// reader is a long-lived buffered reader over rawConn. It must be
	// persistent (not per-call) so that no bytes are stranded in a discarded
	// buffer between transactions, and it is only ever touched by the goroutine
	// currently admitted to a transaction, which is what makes it safe.
	reader *bufio.Reader
}

// New returns a Connection that owns rawConn, together with the persistent
// buffered reader built over it.
//
// A socket can only be installed at construction. There is deliberately no
// setter: swapping it on a live Connection would race an in-flight transaction
// or resurrect one that Close or a transport failure had retired. Reconnecting
// means a new Connection, which is also what keeps reconnect safe — a stream
// that failed ambiguously must never be resumed.
func New(rawConn net.Conn) *Connection {
	return &Connection{
		rawConn: rawConn,
		reader:  bufio.NewReader(rawConn),
	}
}

// sem returns the admission semaphore, creating it on first use.
func (conn *Connection) sem() chan struct{} {
	conn.admissionOnce.Do(func() {
		conn.admission = make(chan struct{}, 1)
	})
	return conn.admission
}

// stateError reports the error a non-open connection should return, or nil if
// the connection is open. The caller must not hold mu.
func (conn *Connection) stateError() error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	return errForState(conn.state)
}

func errForState(s connState) error {
	switch s {
	case stateClosed:
		return transportError(ErrConnectionClosed)
	case stateBroken:
		return transportError(ErrConnectionBroken)
	default:
		return nil
	}
}

// checkUsable reports why the connection cannot begin a new transaction, or nil
// if it can. Run before any work on the caller's behalf, including marshalling.
// State outranks context.
//
// Refusing an already-cancelled context here rather than letting it reach the
// close-on-cancel watchdog is what stops one caller's stale context from
// destroying a connection shared by everyone else: nothing has been written, so
// there is no boundary ambiguity and nothing to poison.
func (conn *Connection) checkUsable(ctx context.Context) error {
	if err := conn.stateError(); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		// Refused, not retired: nothing has been written.
		return refusedError(err)
	}
	return nil
}

// poison performs the compare-and-set open -> broken. It is a no-op from any
// other state, which is what makes closed terminal and outranking broken.
func (conn *Connection) poison() {
	conn.mu.Lock()
	if conn.state == stateOpen {
		conn.state = stateBroken
	}
	conn.mu.Unlock()
}

// transportError reports a failure that leaves the connection unusable.
func transportError(cause error) error {
	return &commands.TransportError{Cause: cause}
}

// refusedError reports a request turned away before the socket was touched.
// Connections are shared, so reporting this as fatal would have every other
// user discard a transport that is byte-for-byte where it was.
func refusedError(cause error) error {
	return &commands.TransportError{Cause: cause, ConnectionUsable: true}
}

// transact runs one complete request/response exchange and returns the raw
// bytes of the framed response document.
//
// The order is fixed, and is the acceptance criterion for the race detector:
// acquire admission; under the state lock verify open and snapshot
// rawConn/reader; drop the lock; do all I/O through the snapshot, never the
// struct fields; retake the lock only to classify; release admission.
//
// Every error returned here is boundary-unknown and has already poisoned the
// connection. Boundary-known failures are the caller's to classify.
func (conn *Connection) transact(ctx context.Context, request []byte) ([]byte, error) {
	// Refuse a closed or broken connection, or an already-cancelled context,
	// before waiting for admission the caller could not use anyway.
	if err := conn.checkUsable(ctx); err != nil {
		return nil, err
	}

	// Step 1: admission.
	if err := conn.acquire(ctx); err != nil {
		return nil, err
	}
	defer conn.release()

	// Steps 2 and 3: verify state, snapshot the socket and reader, drop the lock.
	c, r, err := conn.snapshot()
	if err != nil {
		return nil, err
	}

	join := conn.armCancelWatchdog(ctx, c)

	// Step 4: all I/O through the snapshot, never the struct fields.
	doc, ioErr := performRequest(ctx, c, r, request)

	// Step 5: classify. join() must run first: it settles the cancellation race
	// and guarantees the watchdog has finished touching shared state.
	wasCancelled := join()
	if err := conn.classify(ctx, ioErr, wasCancelled); err != nil {
		return nil, err
	}
	return doc, nil
}

// acquire takes the admission semaphore, honoring ctx while the caller waits.
//
// A non-blocking acquire is tried first. Without it, a caller whose context was
// cancelled between checkUsable and here would race between "admitted" and
// "cancelled while waiting" in a single select, making the two rows of the
// cancellation table nondeterministic. Splitting it means the ctx.Done arm is
// only reachable when the caller genuinely had to wait.
//
// This grants admission, not a place in line: see the note on fairness in the
// Connection doc comment.
func (conn *Connection) acquire(ctx context.Context) error {
	sem := conn.sem()
	select {
	case sem <- struct{}{}:
		return nil
	default:
	}

	select {
	case sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		// Cancelled while waiting: the socket was never touched, so there is no
		// state transition to make and the connection is not retired.
		return refusedError(ctx.Err())
	}
}

// release hands admission to the next queued caller.
func (conn *Connection) release() {
	<-conn.sem()
}

// snapshot verifies the connection is open and copies the socket and reader
// into locals, so the transaction owner never re-reads the struct fields while
// Close or a cancellation watchdog may be mutating them.
func (conn *Connection) snapshot() (net.Conn, *bufio.Reader, error) {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	if err := errForState(conn.state); err != nil {
		return nil, nil, err
	}
	if conn.rawConn == nil || conn.reader == nil {
		return nil, nil, &commands.ConnectionNotInitializedError{Cause: ErrNoRawConn}
	}
	return conn.rawConn, conn.reader, nil
}

// armCancelWatchdog starts close-on-cancel for a transaction this caller now
// owns, and returns a join function reporting whether the watchdog fired.
//
// It is armed here rather than earlier because a watchdog armed while merely
// queued would let one caller's cancellation close the socket out from under
// whoever is currently mid-request.
func (conn *Connection) armCancelWatchdog(ctx context.Context, c net.Conn) func() bool {
	if ctx.Done() == nil {
		return func() bool { return false }
	}

	// These three variables resolve the race between "the watchdog observed
	// cancellation" and "the owner finished its I/O" to exactly one winner. A
	// plain select on ctx.Done() vs. done is not enough: when both are ready Go
	// picks at random, so a cancellation arriving strictly *after* a successful
	// transaction could still win and destroy a healthy connection. The owner
	// publishes finished before closing done, so any cancellation observed from
	// that point on is a no-op.
	var (
		outcomeMu sync.Mutex
		cancelled bool
		finished  bool
	)
	done := make(chan struct{})
	guard := make(chan struct{})

	go func() {
		defer close(guard)
		select {
		case <-ctx.Done():
			outcomeMu.Lock()
			if finished {
				// The owner's I/O already completed; this cancellation happened
				// after the fact and must not touch the socket.
				outcomeMu.Unlock()
				return
			}
			cancelled = true
			outcomeMu.Unlock()
			// Acting on behalf of the transaction owner: CAS open -> broken so a
			// concurrent Close cannot be clobbered, then interrupt the blocked
			// socket operation.
			conn.poison()
			_ = c.Close()
		case <-done:
		}
	}()

	return func() bool {
		outcomeMu.Lock()
		finished = true
		wasCancelled := cancelled
		outcomeMu.Unlock()
		// Joining the watchdog before admission is released is essential: a
		// straggling cancel goroutine must never close the socket during the
		// *next* transaction. It also guarantees that if the watchdog won, its
		// poison() has landed before the outcome is classified.
		close(done)
		<-guard
		return wasCancelled
	}
}

// classify records the outcome of a transaction under the state lock and maps
// it onto the error the caller sees.
//
// Every ioErr reaching here leaves the message-boundary position unknown, so it
// poisons. Boundary-known failures (a document that framed cleanly but did not
// decode into the caller's type) never reach this function.
func (conn *Connection) classify(ctx context.Context, ioErr error, wasCancelled bool) error {
	conn.mu.Lock()
	if ioErr != nil && conn.state == stateOpen {
		conn.state = stateBroken
	}
	st := conn.state
	conn.mu.Unlock()

	if ioErr == nil && !wasCancelled {
		return nil
	}

	switch {
	case st == stateClosed:
		// Close outranks everything, including a cancellation landing after it.
		return transportError(joinCause(ErrConnectionClosed, ioErr))

	case wasCancelled:
		// An error even when ioErr is nil: the response can frame completely
		// just before the cancellation lands, leaving the caller a good document
		// and a connection the watchdog has already destroyed. Reporting success
		// there would make every later call fail for no visible reason.
		//
		// Any platform error is wrapped underneath so errors.Is finds the
		// cancellation, not the "use of closed connection" the watchdog provoked.
		return transportError(joinCause(ctx.Err(), ioErr))

	default:
		return transportError(ioErr)
	}
}

// joinCause wraps cause underneath reason, keeping both matchable by errors.Is.
func joinCause(reason, cause error) error {
	if cause == nil {
		return reason
	}
	return fmt.Errorf("%w: %w", reason, cause)
}

// performRequest writes request and frames exactly one XML document from the
// response stream. It touches only the snapshot passed in.
func performRequest(ctx context.Context, c net.Conn, r *bufio.Reader, request []byte) ([]byte, error) {
	deadline, hasDeadline := ctx.Deadline()
	if hasDeadline {
		_ = c.SetWriteDeadline(deadline)
	} else {
		_ = c.SetWriteDeadline(time.Time{})
	}

	n, err := c.Write(request)
	if err != nil {
		// Any write error poisons: a partial transmission cannot be ruled out,
		// so the server may or may not act on a fragment of the command.
		return nil, err
	}
	if n != len(request) {
		// A short write with no error is ambiguous about what the server
		// received. Do not retry the remainder — resuming a half-sent command
		// is only safe under deliberately designed retry semantics.
		return nil, fmt.Errorf("%w: wrote %d of %d bytes", io.ErrShortWrite, n, len(request))
	}

	if hasDeadline {
		_ = c.SetReadDeadline(deadline)
	} else {
		// Clear any deadline left by a previous transaction.
		_ = c.SetReadDeadline(time.Time{})
	}
	return readDocument(r)
}

// readDocument frames exactly one XML document off r and returns every byte
// consumed while doing so.
//
// Framing is structural, not length-based: treating a short read as
// end-of-message hangs forever on a response that is an exact multiple of the
// buffer size.
//
// A fresh xml.Decoder per call is fine; the *persistent* reader is what
// matters. xml.NewDecoder consumes directly from an io.ByteReader rather than
// wrapping it, so discarding the decoder strands no bytes.
func readDocument(r *bufio.Reader) ([]byte, error) {
	rec := &recordingByteReader{src: r}
	dec := xml.NewDecoder(rec)

	depth := 0
	started := false
	for {
		tok, err := dec.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if !started && len(rec.buf) == 0 {
					// Peer closed without sending anything.
					return nil, io.EOF
				}
				return nil, io.ErrUnexpectedEOF
			}
			return nil, err
		}

		switch tok.(type) {
		case xml.StartElement:
			// Only start counting once the root has actually been seen. A
			// decoder legally emits ProcInst, Directive, Comment and
			// whitespace CharData before the root while depth is still zero.
			started = true
			depth++
		case xml.EndElement:
			depth--
			if started && depth == 0 {
				return rec.buf, nil
			}
		}
	}
}

// recordingByteReader wraps the persistent reader and records every byte the
// decoder pulls through it.
//
// Because the decoder consumes byte-at-a-time from an io.ByteReader, the
// recording is exactly the bytes the decoder consumed — no InputOffset
// arithmetic and no tee-buffer slicing, both of which record what the
// *reader* pulled rather than what the decoder used.
//
// It is only ever used by the goroutine holding the transaction.
type recordingByteReader struct {
	src *bufio.Reader
	buf []byte
}

func (r *recordingByteReader) ReadByte() (byte, error) {
	b, err := r.src.ReadByte()
	if err != nil {
		return 0, err
	}
	r.buf = append(r.buf, b)
	return b, nil
}

// Read exists to satisfy io.Reader for xml.NewDecoder's parameter type; the
// decoder itself only calls ReadByte once it detects io.ByteReader.
func (r *recordingByteReader) Read(p []byte) (int, error) {
	n, err := r.src.Read(p)
	if n > 0 {
		r.buf = append(r.buf, p[:n]...)
	}
	return n, err
}

// Execute marshals command, performs one serialized transaction, and decodes
// the framed response document into response.
//
// A closed or broken connection, or an already-cancelled context, is refused
// before command is marshalled. That ordering is contractual: the lifecycle
// error wins over a marshalling error, and a custom MarshalXML never runs its
// side effects on a connection that has stopped accepting work.
//
// A response that frames cleanly but fails to decode leaves the connection
// usable and is reported as *commands.ProtocolError — the stream position is
// known, and the failure is schema drift, not a transport fault.
func (conn *Connection) Execute(ctx context.Context, command interface{}, response interface{}) error {
	// Rechecked inside transact, where the state may have moved on. Here it only
	// establishes precedence over marshalling.
	if err := conn.checkUsable(ctx); err != nil {
		return err
	}

	cmdBuf, err := xml.Marshal(command)
	if err != nil {
		return &commands.RequestEncodingError{Cause: err}
	}

	doc, err := conn.transact(ctx, cmdBuf)
	if err != nil {
		return err
	}

	// Boundary-known from here on: the root element's matching EndElement was
	// consumed, so nothing below may poison the connection.
	if err := xml.Unmarshal(doc, response); err != nil {
		return commands.NewProtocolError(doc, err)
	}

	if rs, ok := response.(commands.ResponseWithStatus); ok {
		status, statusText := rs.GetStatus()
		if gmpErr := commands.ValidateResponse(status, statusText); gmpErr != nil {
			return gmpErr
		}
	}

	return nil
}

// RawXML sends xmlStr verbatim and returns the response as a wire segment.
//
// The result is every byte consumed while locating and framing the response:
// any leading declaration, comments, directives and whitespace before the root
// are included, trailing whitespace after it is not. That asymmetry is
// deliberate — nothing is trimmed, since trimming would be neither exact wire
// fidelity nor clean document extraction — but it means the result is not
// guaranteed to be a standalone well-formed document.
//
// RawXML shares Execute's admission gate, so mixing the two is safe.
func (conn *Connection) RawXML(ctx context.Context, xmlStr string) (string, error) {
	doc, err := conn.transact(ctx, []byte(xmlStr))
	if err != nil {
		return "", err
	}
	return string(doc), nil
}

// Close marks the connection closed and closes the socket.
//
// Safe to call while a request is in flight, and does not wait for it: the
// state transition is published under the lock and the socket closed outside
// it, since net.Conn.Close is itself I/O and may block or re-enter. The
// in-flight request fails; later callers get ErrConnectionClosed before any
// write.
//
// Idempotent. An already-closed socket, e.g. one the cancel watchdog closed,
// is not reported as a Close failure.
func (conn *Connection) Close() error {
	conn.mu.Lock()
	previous := conn.state
	conn.state = stateClosed
	c := conn.rawConn
	conn.mu.Unlock()

	if previous == stateClosed || c == nil {
		return nil
	}

	err := c.Close()
	if err == nil || errors.Is(err, net.ErrClosed) || errors.Is(err, io.ErrClosedPipe) {
		return nil
	}
	return err
}
