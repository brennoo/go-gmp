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
// GMP has no correlation IDs, so a write-then-read transaction must occupy the
// socket exclusively. Connection serializes transactions accordingly: safe for
// concurrent use, but not parallel — callers wait, with no FIFO fairness
// guarantee. Use one Connection per goroutine for genuine parallelism.
//
// admission is a channel rather than a mutex so a waiting caller can honor
// ctx cancellation (a mutex can't be acquired under a context); mu guards
// state, rawConn and reader separately, so Close can act while another
// goroutine holds the transaction, and is never held across blocking I/O.
//
// A Connection must not be copied after first use.
type Connection struct {
	// admission is lazily created by sem(); see admissionOnce.
	admission chan struct{}
	// admissionOnce keeps the zero value usable: Connection is exported, so a
	// bare literal must stay constructible (a nil chan would block forever).
	admissionOnce sync.Once

	mu      sync.Mutex
	state   connState
	rawConn net.Conn
	// reader must be persistent, not per-call, so no bytes are stranded in a
	// discarded buffer between transactions.
	reader *bufio.Reader
}

// New returns a Connection that owns rawConn, together with the persistent
// buffered reader built over it.
//
// There is deliberately no setter to swap the socket later: that would race
// an in-flight transaction or resurrect one Close had retired. Reconnecting
// means building a new Connection.
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
// Order is fixed: acquire admission, snapshot rawConn/reader under the state
// lock, drop the lock, do all I/O through the snapshot (never the struct
// fields), then retake the lock only to classify.
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
// The non-blocking attempt is tried first so the ctx.Done arm is only
// reachable once the caller genuinely had to wait — otherwise a context
// cancelled between checkUsable and here would race "admitted" against
// "cancelled" nondeterministically.
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

	// A plain select on ctx.Done() vs. done isn't enough: when both are ready Go
	// picks at random, so a cancellation arriving after a successful transaction
	// could still win and destroy a healthy connection. finished/outcomeMu
	// resolve that race to one winner deterministically.
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
		// Must join before admission is released: a straggling cancel goroutine
		// must never close the socket during the next transaction.
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
		// An error even when ioErr is nil: the response can frame completely just
		// before cancellation lands, leaving a good document but a connection the
		// watchdog already destroyed. ctx.Err() is wrapped outermost so errors.Is
		// finds the cancellation, not the platform error it provoked.
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
// end-of-message would hang forever on a response that's an exact multiple of
// the buffer size. A fresh xml.Decoder per call is fine — it consumes directly
// from an io.ByteReader rather than wrapping it, so discarding it strands no
// bytes from the persistent reader underneath.
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
// decoder pulls through it — exactly the bytes consumed, since the decoder
// reads byte-at-a-time via ReadByte rather than through buffered Read calls.
// Only ever used by the goroutine holding the transaction.
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
		// gvmd sometimes reports failures under a generic <gmp_response/> root
		// instead of the command's own typed one (T-7b). ErrorFromDocument only
		// matches that root, so it's safe to try unconditionally here.
		if gmpErr := commands.ErrorFromDocument(doc); gmpErr != nil {
			return gmpErr
		}
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
// leading declaration/comments/whitespace before the root are included,
// trailing whitespace after it is not (nothing is trimmed; it's simply never
// consumed). The result is therefore not guaranteed to be standalone
// well-formed XML.
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
