//go:build integration

package gmptest

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	gmp "github.com/brennoo/go-gmp"
	"github.com/brennoo/go-gmp/client"
	"github.com/brennoo/go-gmp/commands"
	"github.com/brennoo/go-gmp/connections"
	"github.com/brennoo/go-gmp/internal/dockerstack"
)

// This file is the machinery of the live capture; capture_sequence_test.go is
// the sequence it drives. Nothing here has ever run against a real gvmd — see
// the report accompanying this change for the list of behaviours that only a
// real run can confirm.

var (
	// captureEnabled gates the whole thing behind an explicit opt-in. The rest
	// of the integration suite only reads; this one creates and deletes
	// resources and writes to the repository, so it must never be something a
	// `go test -tags integration ./...` picks up by accident.
	captureEnabled = flag.Bool("capture", false,
		"capture live fixtures into testdata/live (writes to the repository; requires -tags integration)")

	// captureMaxBytes keeps a response that is genuinely enormous out of git.
	// get_nvts with a synced feed is tens of megabytes and get_info over the
	// SCAP database is larger still; neither is a fixture, both are a repo
	// nobody can clone. The size is recorded either way, which is itself the
	// finding.
	captureMaxBytes = flag.Int("capture.max-bytes", 1<<20,
		"skip writing any fixture whose response exceeds this many bytes")

	// captureCommandTimeout bounds one exchange. It is far shorter than the
	// readiness ladder's budget because by this point gvmd is known to be up:
	// a command that hangs here is a malformed-input probe that gvmd never
	// answered, which is a result, not a reason to burn an hour.
	captureCommandTimeout = flag.Duration("capture.command-timeout", 2*time.Minute,
		"per-command timeout for a single capture exchange")

	// captureComposeFile locates the stack backing the gvmd this run talks
	// to, so its image digest can be read from the running container
	// (internal/dockerstack) rather than inferred from a tag. The default
	// matches -gmp.cafile's own default resolution of hack/ relative to this
	// package's test working directory.
	captureComposeFile = flag.String("capture.compose-file", "../../hack/docker-compose.yml",
		"path to the Compose file backing the live gvmd, used to resolve its image digest for Meta.ImageDigest")
)

// captureOutcome is what happened to one attempted fixture.
type captureOutcome string

const (
	outcomeCaptured captureOutcome = "captured"
	outcomeFailed   captureOutcome = "failed"
	outcomeSkipped  captureOutcome = "skipped"
	// outcomeProbe is an exchange that was deliberately not written as a
	// fixture. The only current use is the unknown-command probe: the corpus
	// files fixtures by command directory and TestCorpusIsComplete requires
	// every directory to name a registered command, so a response to
	// <no_such_command/> has nowhere honest to live. Its root element is
	// reported instead.
	outcomeProbe captureOutcome = "probe"
)

// captureResult is one line of the end-of-run summary.
type captureResult struct {
	Command string
	Desc    string
	ID      string
	Outcome captureOutcome
	Status  string
	Root    string
	Bytes   int
	Notes   []string
}

// captureRun is the state of one capture.
type captureRun struct {
	t   *testing.T
	ctx context.Context

	conn gmp.Connection
	red  *Redactor

	runID       string
	gmpVersion  string
	imageDigest string
	reconnects  int

	seq     map[string]int
	written map[string]bool
	touched map[string]bool

	results  []captureResult
	notices  map[string]int
	rootDiff []string
}

// TestCaptureLive records the live half of the fixture corpus.
//
// It is not an assertion suite. Its job is to come back with as much of the
// corpus as gvmd will give it and an honest account of the rest: one command
// failing must never cost the other hundred, because a partial corpus is useful
// and an aborted one is not. The single exception is the credential sweep,
// which stops everything — see (*captureRun).write.
//
//	go test -tags integration ./internal/gmptest -run TestCaptureLive -capture \
//	  -gmp.addr=127.0.0.1:9390 -gmp.pass=… -gmp.servername=… -timeout 60m
func TestCaptureLive(t *testing.T) {
	if !*captureEnabled {
		t.Skip("Skipping: live capture writes to testdata/live and creates resources on the " +
			"server; pass -capture to opt in")
	}

	requireLive(t)
	r := newCaptureRun(t)

	// Registered first, so LIFO runs it last: after every teardown below has
	// had its own exchange recorded.
	t.Cleanup(r.reportSummary)
	t.Cleanup(r.captureEmptyTrashcan)

	r.run()
}

func newCaptureRun(t *testing.T) *captureRun {
	t.Helper()

	if liveConn == nil {
		t.Fatal("liveConn is nil; the readiness ladder did not complete")
	}
	if *livePass == *liveUser {
		// The sweep matches by value, not by which field the value came from.
		// gvmd echoes the username in <owner><name> on nearly every
		// authenticated response, so if the password is the same string, the
		// fatal -gmp.pass term matches that legitimate, unredacted text and
		// aborts the run — hack/gvmd-up.sh's own documented default
		// (admin/admin) triggers this. Failing here, before anything is
		// written, turns that into an actionable message instead of a
		// mid-run abort with a captured-but-orphaned partial fixture set.
		t.Fatalf("-gmp.pass (%q) must not equal -gmp.user (%q): the redaction sweep cannot tell "+
			"the password's own text apart from the username being legitimately echoed elsewhere "+
			"once they're the same string; set GVMD_ADMIN_PASSWORD to something else before running "+
			"hack/gvmd-up.sh, or pass a different -gmp.pass here if the account already has one",
			*livePass, *liveUser)
	}

	red := NewRedactor()
	// The password is fatal on sight: it has no legitimate place in any
	// response. The username is not, because gvmd echoes it in <owner><name>
	// on nearly every authenticated response — making that fatal would abort
	// the run on its first fixture and there would be no corpus at all.
	red.AddSecret("-gmp.pass", *livePass)
	red.AddIdentifier("-gmp.user", *liveUser)

	r := &captureRun{
		t:       t,
		ctx:     context.Background(),
		conn:    liveConn,
		red:     red,
		runID:   newRunID(t),
		seq:     map[string]int{},
		written: map[string]bool{},
		touched: map[string]bool{},
		notices: map[string]int{},
	}
	r.gmpVersion = r.detectVersion()
	r.imageDigest = r.resolveImageDigest()
	t.Logf("capture run id %s against GMP protocol version %s, image %s",
		r.runID, r.gmpVersion, r.imageDigest)
	return r
}

// newRunID produces the tag that goes into every created resource's name.
//
// A timestamp alone is not enough: two runs started in the same second on the
// same stack would produce colliding names and the second would fail on
// gvmd's uniqueness checks rather than on anything interesting.
func newRunID(t *testing.T) string {
	t.Helper()
	var b [3]byte
	if _, err := rand.Read(b[:]); err != nil {
		t.Fatalf("Failed to read random bytes for the run id: %v", err)
	}
	return time.Now().UTC().Format("20060102t150405") + "-" + hex.EncodeToString(b[:])
}

// resourceName is the name given to every resource this run creates.
//
// The go-gmp-fixture- prefix plus the run id is what makes a crashed run
// recoverable: whatever is left behind is identifiable by name and can be swept
// with a filter, which matters because a run that dies mid-way leaves its
// t.Cleanup teardowns unexecuted.
func (r *captureRun) resourceName(kind string) string {
	return "go-gmp-fixture-" + r.runID + "-" + kind
}

// detectVersion returns the GMP protocol version get_version reports, not a
// gvmd application/build version — GMP does not expose one.
func (r *captureRun) detectVersion() string {
	ctx, cancel := context.WithTimeout(r.ctx, *captureCommandTimeout)
	defer cancel()

	resp, err := liveClient.GetVersion(ctx, &commands.GetVersion{})
	if err != nil || resp == nil {
		return "unknown"
	}
	return resp.Version
}

// resolveImageDigest reads gvmd's running image digest off the Compose stack
// itself, rather than inferring it from a tag.
//
// Unlike detectVersion, a failure here fails the whole run instead of
// recording "unknown": loadMeta now requires every live fixture to carry
// Meta.ImageDigest, so a run that cannot resolve it would only produce a
// corpus that fails to load — better to stop before writing anything.
func (r *captureRun) resolveImageDigest() string {
	ctx, cancel := context.WithTimeout(r.ctx, *captureCommandTimeout)
	defer cancel()

	digest, err := dockerstack.ImageDigest(ctx, *captureComposeFile, "gvmd")
	if err != nil {
		r.t.Fatalf("resolving gvmd image digest: %v", err)
	}
	return digest
}

// dial opens a second, unauthenticated connection.
//
// Needed because the pre-auth and unauthenticated-failure fixtures cannot be
// taken on liveConn: it authenticated during the readiness ladder and GMP has
// no way to un-authenticate a session.
func (r *captureRun) dial(command, desc string) (gmp.Connection, bool) {
	opts, err := dialOptions()
	if err != nil {
		r.fail(command, desc, "dial options: "+err.Error())
		return nil, false
	}
	conn, err := connections.NewTLSConnection(r.ctx, *liveAddr, opts...)
	if err != nil {
		r.fail(command, desc, "dial: "+err.Error())
		return nil, false
	}
	return conn, true
}

// send marshals req from a commands.* value and captures the exchange on the
// authenticated connection.
//
// Marshalling from the real request type rather than a hand-written string is
// deliberate: it makes every captured NNN_request.xml a byte-exact record of
// what this library actually puts on the wire, so a request struct that cannot
// express a valid command shows up as a 400 from gvmd instead of being papered
// over by a string literal that happens to be right.
func (r *captureRun) send(command, desc string, req any) (string, bool) {
	return r.sendOn(r.conn, command, desc, req)
}

func (r *captureRun) sendOn(conn gmp.Connection, command, desc string, req any) (string, bool) {
	b, err := xml.Marshal(req)
	if err != nil {
		r.fail(command, desc, "marshal: "+err.Error())
		return "", false
	}
	return r.exchange(conn, command, desc, string(b))
}

// sendRaw captures an exchange whose request cannot come from a commands.*
// value: a deliberately malformed document, or a payload shape no struct in the
// package can express.
func (r *captureRun) sendRaw(command, desc, reqXML string) (string, bool) {
	return r.exchange(r.conn, command, desc, reqXML)
}

// rawAttempt is one fallback request shape.
type rawAttempt struct {
	desc string
	xml  string
}

// sendWithFallbacks sends the typed request and, if gvmd does not answer 2xx,
// works through hand-written alternatives until one does.
//
// This exists for the handful of commands where the Go type is suspected of
// being unable to express a valid request — create_task's container form,
// where TaskTarget has a non-omitempty <name>. Both the rejection and the
// working shape are captured, so the pair is the evidence for the struct fix.
func (r *captureRun) sendWithFallbacks(command, desc string, req any, fallbacks []rawAttempt) string {
	resp, ok := r.send(command, desc, req)
	if ok && statusOK(resp) {
		return resp
	}
	for _, fb := range fallbacks {
		resp, ok = r.sendRaw(command, fb.desc, fb.xml)
		if ok && statusOK(resp) {
			return resp
		}
	}
	// The last response, whatever it was: the caller pulls an id out of it and
	// gets "" if there is none, which is already the "dependency missing" path.
	return resp
}

// exchange runs one request/response and records the outcome.
//
// gvmd is observed to close the socket outright after answering some
// requests (e.g. a disabled command's 503), which the client correctly
// treats as an unrecoverable stream-boundary ambiguity rather than resuming
// it — see TransportError.ConnectionUsable. Left unhandled, that turns one
// command's connection loss into a cascade of "connection broken" for every
// command after it in the sequence, which is exactly the outcome
// TestCaptureLive's own doc comment says must never happen. So when the
// broken connection is the shared r.conn (not a throwaway pre-auth dial),
// reconnect once and retry this exchange on the fresh, reauthenticated
// connection before giving up on it.
func (r *captureRun) exchange(conn gmp.Connection, command, desc, reqXML string) (string, bool) {
	if _, known := Registry[command]; !known {
		r.fail(command, desc, "not a registered GMP command; the corpus has no directory for it")
		return "", false
	}

	resp, err := r.rawXML(conn, reqXML)
	if err != nil && conn == r.conn && unusableTransport(err) {
		if rerr := r.reconnect(); rerr != nil {
			r.fail(command, desc, fmt.Sprintf(
				"transport: %v (reconnect failed, giving up: %v)", err, rerr))
			return "", false
		}
		r.reconnects++
		resp, err = r.rawXML(r.conn, reqXML)
	}
	if err != nil {
		r.fail(command, desc, "transport: "+err.Error())
		return "", false
	}
	r.write(command, desc, reqXML, resp)
	return resp, true
}

// rawXML runs one request/response with the per-command timeout, without any
// retry policy — exchange is the only caller that needs one.
func (r *captureRun) rawXML(conn gmp.Connection, reqXML string) (string, error) {
	ctx, cancel := context.WithTimeout(r.ctx, *captureCommandTimeout)
	defer cancel()
	return conn.RawXML(ctx, reqXML)
}

// unusableTransport reports whether err is a TransportError that has already
// retired the connection it came from, as opposed to one where the request
// was merely refused (e.g. an already-cancelled context) with the connection
// left exactly as it was.
func unusableTransport(err error) bool {
	var te *commands.TransportError
	return errors.As(err, &te) && !te.ConnectionUsable
}

// reconnect replaces r.conn with a freshly dialed, reauthenticated
// connection, closing the old one. Used after exchange sees a TransportError
// with ConnectionUsable false on the shared connection: a new Connection is
// the only way to recover, per its own doc comment ("a stream that failed
// ambiguously must never be resumed").
func (r *captureRun) reconnect() error {
	opts, err := dialOptions()
	if err != nil {
		return fmt.Errorf("dial options: %w", err)
	}
	conn, err := connections.NewTLSConnection(r.ctx, *liveAddr, opts...)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	creds := &commands.AuthenticateCredentials{Username: *liveUser, Password: *livePass}
	if _, err := client.New(conn).Authenticate(r.ctx, &commands.Authenticate{Credentials: creds}); err != nil {
		_ = conn.Close()
		return fmt.Errorf("authenticate: %w", err)
	}
	if r.conn != nil {
		_ = r.conn.Close()
	}
	r.conn = conn
	return nil
}

// unfiled sends a request that has no home in the corpus and records only what
// came back.
//
// Two things need this and neither is a fixture: a command gvmd knows but
// Registry does not (delete_tls_certificate, which commands/ never modelled),
// and a command nobody knows at all (the unknown-command probe for T-7b).
// TestCorpusIsComplete requires every directory under testdata to name a
// registered command, so filing either one would mean inventing a directory
// that lies about what it holds.
//
// expectRoot is the response element the caller expected, or "" if there is no
// expectation; a mismatch or an unknown expectation is reported under T-7b.
func (r *captureRun) unfiled(conn gmp.Connection, desc, reqXML, expectRoot string) {
	ctx, cancel := context.WithTimeout(r.ctx, *captureCommandTimeout)
	defer cancel()

	res := captureResult{Command: "-", Desc: desc, Outcome: outcomeProbe}

	resp, err := conn.RawXML(ctx, reqXML)
	if err != nil {
		res.Notes = append(res.Notes, "transport: "+err.Error())
		r.results = append(r.results, res)
		return
	}

	root, status, statusText := describeResponse(resp)
	res.Root, res.Status, res.Bytes = root, status, len(resp)
	res.Notes = append(res.Notes, fmt.Sprintf("root <%s>, status %s %q (not written: no registered command)",
		root, status, statusText))
	r.results = append(r.results, res)

	if expectRoot == "" || root != expectRoot {
		r.rootDiff = append(r.rootDiff,
			fmt.Sprintf("%s -> <%s> (%s %s)", desc, root, status, statusText))
	}
}

// write redacts, sweeps and persists one exchange.
//
// The order is the whole safety argument and must not be rearranged: redact
// from the per-command spec, then sweep the resulting bytes for literal
// secrets, and only then touch the filesystem. Nothing is written between the
// two checks, so a spec that missed a path cannot produce a file on disk.
func (r *captureRun) write(command, desc, reqXML, respXML string) {
	res := captureResult{Command: command, Desc: desc, Bytes: len(respXML)}
	root, status, statusText := describeResponse(respXML)
	res.Root, res.Status = root, status

	// T-7b: does gvmd ever answer under a root element other than
	// <command>_response? Recorded here for every single exchange, so the
	// answer comes from the whole corpus rather than from the four probes that
	// were expected to show it.
	if want := command + "_response"; root != "" && root != want {
		note := fmt.Sprintf("ROOT MISMATCH: <%s>, not <%s> (status %s %q)", root, want, status, statusText)
		res.Notes = append(res.Notes, note)
		r.rootDiff = append(r.rootDiff, command+": "+desc+" -> "+note)
	}

	if len(respXML) > *captureMaxBytes {
		res.Outcome = outcomeSkipped
		res.Notes = append(res.Notes, fmt.Sprintf(
			"response is %d bytes, over -capture.max-bytes=%d; raise the flag to keep it",
			len(respXML), *captureMaxBytes))
		r.results = append(r.results, res)
		return
	}

	reqOut, respOut, applied, err := r.redactPair(command, reqXML, respXML)
	if err != nil {
		res.Outcome = outcomeFailed
		res.Notes = append(res.Notes, err.Error())
		r.results = append(r.results, res)
		return
	}

	r.sweepOrDie(command, desc, reqOut, respOut)

	res.ID = r.nextID(command)
	if err := r.persist(command, desc, res.ID, reqOut, respOut, applied); err != nil {
		res.Outcome = outcomeFailed
		res.Notes = append(res.Notes, err.Error())
		r.results = append(r.results, res)
		return
	}

	res.Outcome = outcomeCaptured
	if len(applied) > 0 {
		res.Notes = append(res.Notes, "redacted: "+strings.Join(applied, " "))
	}
	r.results = append(r.results, res)
}

// redactPair applies the command's spec to both directions and returns the
// union of the paths that matched.
func (r *captureRun) redactPair(command, reqXML, respXML string) (req, resp []byte, applied []string, err error) {
	req, reqApplied, err := r.red.Redact(command, []byte(reqXML))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("redact request: %w", err)
	}
	resp, respApplied, err := r.red.Redact(command, []byte(respXML))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("redact response: %w", err)
	}

	seen := map[string]bool{}
	for _, p := range append(reqApplied, respApplied...) {
		if !seen[p] {
			seen[p] = true
			applied = append(applied, p)
		}
	}
	sort.Strings(applied)
	return req, resp, applied, nil
}

// sweepOrDie is the backstop for a path the spec missed.
//
// A fatal hit stops the whole run with t.Fatalf. That is not an overreaction:
// the specs are the primary defence and a hit means one of them is wrong, so
// every fixture already on disk from this run is suspect too.
//
// Aborting from the sequence still runs the registered teardowns. Aborting from
// inside a teardown does not run the remaining ones — t.Fatalf is a Goexit and
// the cleanup stack unwinds with it — which is the other reason every created
// resource carries the run id in its name.
func (r *captureRun) sweepOrDie(command, desc string, docs ...[]byte) {
	var fatal, notices []string
	for _, d := range docs {
		f, n := r.red.Sweep(d)
		fatal = append(fatal, f...)
		notices = append(notices, n...)
	}
	for _, n := range notices {
		r.notices[n]++
	}
	if len(fatal) == 0 {
		return
	}

	// The labels, never the values: an error message that quoted the secret
	// would put it in the test log and from there into a CI transcript.
	r.results = append(r.results, captureResult{
		Command: command, Desc: desc, Outcome: outcomeFailed,
		Notes: []string{"SWEEP: " + strings.Join(fatal, ", ")},
	})
	r.t.Fatalf("Aborting the capture: %s (%s) was about to be written with the value of %s in it. "+
		"A redaction spec in internal/gmptest/redact.go missed a path; nothing was written for this "+
		"exchange, but review everything this run already wrote before committing any of it.",
		command, desc, strings.Join(fatal, ", "))
}

// persist writes the four files that make up one fixture.
//
// The request and response go through os.WriteFile directly, not WriteGolden.
// That is deliberate and is the one place in the package where it is correct:
// WriteGolden refuses evidence files by design, because -update must never be
// able to rewrite a captured document to match a regression. Capture is not an
// -update path — it is gated on -capture, it never runs from a golden test, and
// it is the only thing that may create these bytes in the first place.
func (r *captureRun) persist(command, desc, id string, req, resp []byte, applied []string) error {
	dir := filepath.Join(TestdataDir, string(Live), command)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}
	r.touched[dir] = true

	f := Fixture{Command: command, Provenance: Live, ID: id, Dir: dir}
	f.Meta = Meta{
		Command:    command,
		Provenance: Live,
		Source:     desc,
		// Byte-faithful exactly when nothing was rewritten. A redacted fixture
		// is still evidence, but it is no longer the bytes gvmd sent, and the
		// distinction has to survive into the sidecar or a later byte-exactness
		// assertion would be asserting against a placeholder.
		ByteFaithful:      len(applied) == 0,
		CapturedAt:        time.Now().UTC().Format(time.RFC3339),
		GmpVersion:        r.gmpVersion,
		GmpVersionSource:  GmpVersionCaptured,
		ImageDigest:       r.imageDigest,
		ImageDigestSource: ImageDigestResolved,
		Redactions:        applied,
	}

	if err := os.WriteFile(f.RequestPath(), req, 0o600); err != nil {
		return fmt.Errorf("write request: %w", err)
	}
	if err := os.WriteFile(f.ResponsePath(), resp, 0o600); err != nil {
		return fmt.Errorf("write response: %w", err)
	}
	if err := f.WriteMeta(f.Meta); err != nil {
		return fmt.Errorf("write meta: %w", err)
	}

	// TestCorpusIsComplete requires a fidelity report next to every fixture, so
	// it is generated here rather than left for a later -update run: a capture
	// that leaves the corpus failing its own completeness check is not done.
	report := Compare(resp, Registry[command].NewResp()).Report(fidelityHeader(f))
	if err := WriteGolden(f.FidelityPath(), []byte(report)); err != nil {
		return fmt.Errorf("write fidelity report: %w", err)
	}

	for _, p := range []string{f.RequestPath(), f.ResponsePath(), f.MetaPath(), f.FidelityPath()} {
		r.written[p] = true
	}
	return nil
}

func (r *captureRun) nextID(command string) string {
	r.seq[command]++
	return fmt.Sprintf("%03d", r.seq[command])
}

func (r *captureRun) fail(command, desc, note string) {
	r.results = append(r.results, captureResult{
		Command: command, Desc: desc, Outcome: outcomeFailed, Notes: []string{note},
	})
}

func (r *captureRun) skip(command, desc, note string) {
	r.results = append(r.results, captureResult{
		Command: command, Desc: desc, Outcome: outcomeSkipped, Notes: []string{note},
	})
}

// reportSummary is the deliverable of a run that did not complete cleanly, and
// the review checklist for one that did.
func (r *captureRun) reportSummary() {
	t := r.t
	var counts = map[captureOutcome]int{}
	for _, res := range r.results {
		counts[res.Outcome]++
	}

	var b strings.Builder
	fmt.Fprintf(&b, "\ncapture summary (run %s, GMP protocol version %s)\n", r.runID, r.gmpVersion)
	fmt.Fprintf(&b, "captured=%d failed=%d skipped=%d probes=%d exchanges=%d reconnects=%d\n\n",
		counts[outcomeCaptured], counts[outcomeFailed], counts[outcomeSkipped],
		counts[outcomeProbe], len(r.results), r.reconnects)

	for _, res := range r.results {
		fmt.Fprintf(&b, "%-10s %-24s %-3s %-5s %8dB  %s\n",
			res.Outcome, res.Command, res.ID, res.Status, res.Bytes, res.Desc)
		for _, n := range res.Notes {
			fmt.Fprintf(&b, "               %s\n", n)
		}
	}

	b.WriteString("\nT-7b (does gvmd ever answer under a different root element?)\n")
	if len(r.rootDiff) == 0 {
		b.WriteString("  no divergence observed in this run\n")
	} else {
		for _, d := range r.rootDiff {
			fmt.Fprintf(&b, "  %s\n", d)
		}
	}

	if len(r.notices) > 0 {
		b.WriteString("\nsweep notices (non-fatal; an identifier appeared in the written bytes)\n")
		for _, label := range sortedKeys(r.notices) {
			fmt.Fprintf(&b, "  %s in %d written documents\n", label, r.notices[label])
		}
	}

	if stale := r.staleFiles(); len(stale) > 0 {
		b.WriteString("\npre-existing fixture files this run did not rewrite\n")
		b.WriteString("  (an earlier capture produced more fixtures for these commands than this\n")
		b.WriteString("   one did; nothing is deleted automatically, review and remove by hand)\n")
		for _, p := range stale {
			fmt.Fprintf(&b, "  %s\n", p)
		}
	}

	b.WriteString("\nBefore committing: read the diff. Grep it for " + RedactionPlaceholder + " to see\n")
	b.WriteString("what was rewritten, and for anything that looks like key material that was not.\n")
	t.Log(b.String())
}

// staleFiles lists fixture files in the directories this run touched that this
// run did not write, which is how a shrinking capture leaves orphans behind.
//
// They are reported, never deleted. A tool that removes captured evidence
// because it did not happen to reproduce it this time is a tool that quietly
// loses the one fixture that mattered.
func (r *captureRun) staleFiles() []string {
	var out []string
	for dir := range r.touched {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			p := filepath.Join(dir, e.Name())
			if !e.IsDir() && !r.written[p] {
				out = append(out, p)
			}
		}
	}
	sort.Strings(out)
	return out
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// describeResponse reads the root element name and the status attributes off a
// response without decoding it into anything typed.
//
// Deliberately untyped: the point is to be able to describe a response whose
// root element is not the one the registry would have picked, which is exactly
// the case a typed decode cannot report.
func describeResponse(doc string) (root, status, statusText string) {
	dec := xml.NewDecoder(strings.NewReader(doc))
	for {
		tok, err := dec.Token()
		if err != nil {
			// Including io.EOF: a document with no element at all has no root
			// to report, and an empty root is how the caller learns that.
			return root, status, statusText
		}
		se, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		root = se.Name.Local
		for _, a := range se.Attr {
			switch a.Name.Local {
			case "status":
				status = a.Value
			case "status_text":
				statusText = a.Value
			}
		}
		return root, status, statusText
	}
}

// statusOK reports whether a GMP status attribute is a success.
//
// Prefix matching rather than a set: GMP uses 200, 201 and 202 for success and
// the 2xx family is the contract, so enumerating today's three would make a
// future 203 look like a failure.
func statusOK(doc string) bool {
	_, status, _ := describeResponse(doc)
	return strings.HasPrefix(status, "2")
}

// responseID pulls the id attribute off a create_* response.
func responseID(doc string) string {
	dec := xml.NewDecoder(strings.NewReader(doc))
	for {
		tok, err := dec.Token()
		if err != nil {
			return ""
		}
		se, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		for _, a := range se.Attr {
			if a.Name.Local == "id" {
				return a.Value
			}
		}
		return ""
	}
}

// firstID returns the id attribute of the first <element> in doc, which is
// how an id is lifted out of a list response without a typed decode.
func firstID(doc, element string) string {
	dec := xml.NewDecoder(strings.NewReader(doc))
	for {
		tok, err := dec.Token()
		if err != nil {
			return ""
		}
		se, ok := tok.(xml.StartElement)
		if !ok || se.Name.Local != element {
			continue
		}
		for _, a := range se.Attr {
			if a.Name.Local == "id" {
				return a.Value
			}
		}
	}
}

// setStringField sets a named string field on a request value, reporting
// whether the field existed.
//
// Reflection rather than 60 hand-written variants. The read-only sweep has to
// send every get_* twice, bare and filtered, and a hand-written table would
// have to be edited every time the registry grows — worse, a command whose
// request type simply has no Filter field (get_roles and get_notes have no
// fields at all) would be indistinguishable from one somebody forgot. Here the
// false return is the finding, and it lands in the summary as a skip.
func setStringField(req any, field, value string) bool {
	v := reflect.ValueOf(req)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return false
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return false
	}
	f := v.FieldByName(field)
	if !f.IsValid() || f.Kind() != reflect.String || !f.CanSet() {
		return false
	}
	f.SetString(value)
	return true
}
