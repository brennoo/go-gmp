package codegen

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/brennoo/go-gmp/internal/schema"
)

// goldenDir holds the generator's committed output. These are .go files under
// testdata, so the go tool ignores them: they are golden *bytes*, and the only
// thing that ever compiles them is TestGeneratedCodeCompiles below.
const goldenDir = "testdata/golden"

// TestGenerate is the test CI actually runs for T-9. It needs only the
// committed JSON snapshot, never vendor/XML, so stale generated output shows
// up as an ordinary test failure on any checkout.
func TestGenerate(t *testing.T) {
	got := generateOrFail(t)

	if *update {
		if err := os.RemoveAll(goldenDir); err != nil {
			t.Fatalf("Failed to clear %s: %v", goldenDir, err)
		}
		if err := os.MkdirAll(goldenDir, 0o750); err != nil {
			t.Fatalf("Failed to create %s: %v", goldenDir, err)
		}
		for name, src := range got {
			if err := os.WriteFile(filepath.Join(goldenDir, name), []byte(src), 0o600); err != nil {
				t.Fatalf("Failed to write %s: %v", name, err)
			}
		}
		return
	}

	entries, err := os.ReadDir(goldenDir)
	if err != nil {
		t.Fatalf("No committed goldens (%v).\nRun with -update to create them.", err)
	}
	committed := map[string]bool{}
	for _, e := range entries {
		committed[e.Name()] = true
		want, err := os.ReadFile(filepath.Join(goldenDir, e.Name())) //nolint:gosec // a fixed testdata path
		if err != nil {
			t.Fatalf("Failed to read %s: %v", e.Name(), err)
		}
		src, ok := got[e.Name()]
		if !ok {
			t.Errorf("%s is committed but the generator no longer emits it. Re-run with -update.", e.Name())
			continue
		}
		if src != string(want) {
			t.Errorf("%s changed. Re-run with -update and review the diff.\n"+
				"committed %d bytes, now %d bytes\n%s", e.Name(), len(want), len(src), firstDiff(string(want), src))
		}
	}
	for name := range got {
		if !committed[name] {
			t.Errorf("The generator emits %s but it is not committed. Re-run with -update.", name)
		}
	}
}

// TestGenerateIsDeterministic guards the golden's only real precondition:
// nothing in Generate may depend on map iteration order.
func TestGenerateIsDeterministic(t *testing.T) {
	a, b := generateOrFail(t), generateOrFail(t)
	if len(a) != len(b) {
		t.Fatalf("Two runs produced %d and %d files", len(a), len(b))
	}
	for name, src := range a {
		if b[name] != src {
			t.Errorf("%s differs between two runs of Generate", name)
		}
	}
}

// TestGeneratedCodeCompiles is acceptance criterion 2: the generated source
// has to build without the hand-written commands/create_note.go,
// modify_note.go or get_notes.go, or the generator could be quietly leaning on
// types it did not emit.
//
// It covers every family the manifest declares, review-only ones included, and
// compiles them all into one package — which is also the only check that two
// families never declare the same type name. Family-local dedup means each
// family names its shapes without consulting the others, so "note and override
// both describe an <nvt>" is resolved by the family prefix (NoteNVT vs
// OverrideNVT) rather than by a shared type; if that prefixing ever failed to
// separate two families, this test is where the redeclaration surfaces.
//
// Those three files still exist on disk, so the check cannot simply build the
// output as package commands: every type would be redeclared. Instead the
// files are copied into a throwaway module with only their package clause
// rewritten — that one line is the entire transformation — and the support
// types the override table targets (commands.GMPBool, commands.GMPTime) are
// re-supplied as aliases to the real commands package. Everything the
// generated code needs beyond the standard library therefore still comes from
// commands, and nothing else does.
func TestGeneratedCodeCompiles(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("Skipping: no go toolchain on PATH (%v)", err)
	}
	got := generateOrFail(t)
	// Every family's command files have to be in the build, or "it compiles"
	// would be a claim about whichever subset happened to be generated.
	for _, f := range manifestOrFail(t).Families {
		for _, c := range f.Commands {
			if _, ok := got[CommandFile(c)]; !ok {
				t.Fatalf("family %q: %s was not generated, so this test would not compile it", f.Name, CommandFile(c))
			}
		}
	}

	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to resolve the module root: %v", err)
	}
	dir := t.TempDir()
	pkgDir := filepath.Join(dir, "genproto")
	if err := os.MkdirAll(pkgDir, 0o750); err != nil {
		t.Fatalf("Failed to create the scratch package: %v", err)
	}

	goVersion, err := goDirective(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatalf("Failed to read the module's go directive: %v", err)
	}
	gomod := "module gogmpcodegencheck\n\ngo " + goVersion + "\n\n" +
		"require github.com/brennoo/go-gmp v0.0.0\n\n" +
		"replace github.com/brennoo/go-gmp => " + root + "\n"
	mustWriteFile(t, filepath.Join(dir, "go.mod"), gomod)
	// supportShim() re-exports GMPBool/GMPTime under package commands (see its
	// doc comment); the package clause is rewritten to genproto below along
	// with every generated file, so it lands in the same throwaway package.
	mustWriteFile(t, filepath.Join(pkgDir, "support.go"),
		strings.Replace(supportShim(), "package "+targetPackage+"\n", "package genproto\n", 1))

	for name, src := range got {
		rewritten := strings.Replace(src, "\npackage "+targetPackage+"\n", "\npackage genproto\n", 1)
		if rewritten == src {
			t.Fatalf("%s: expected a %q package clause to rewrite", name, targetPackage)
		}
		mustWriteFile(t, filepath.Join(pkgDir, name), rewritten)
	}

	// context.Background rather than t.Context: this module still declares
	// go 1.23, where testing.T has no Context method.
	cmd := exec.CommandContext(context.Background(), goBin, "build", "./...")
	cmd.Dir = dir
	// GOWORK=off keeps a developer's workspace file from redirecting the
	// replace; GOFLAGS is cleared so an inherited -mod flag cannot apply.
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Generated code does not compile on its own: %v\n%s", err, out)
	}
}

// TestCollisionRequiresAttributeDiscriminators is the generation-time half
// of finding 1's fix: attribute-signature dispatch (see
// collisionDiscriminators/emitCollisionMethods) can only be generated when
// every member of a same-wire-name collision carries at least one start-tag
// attribute no sibling shares. This proves the guard fires when a member has
// none — a collision the generator cannot safely disambiguate — rather than
// trusting that the real corpus never exercises that shape.
func TestCollisionRequiresAttributeDiscriminators(t *testing.T) {
	cmds, err := LoadSnapshots(snapshotDir, []string{"get_filters"})
	if err != nil {
		t.Fatalf("Failed to load the get_filters snapshot: %v", err)
	}
	cmd := cmds[0]
	occurrence := 0
	found := false
	for _, kid := range cmd.Response.Children {
		if kid.Name != "filters" {
			continue
		}
		if occurrence == 1 {
			// Strip the pagination variant's attributes (start, max),
			// leaving it with no signature at all to discriminate on.
			var kept []*schema.Node
			for _, c := range kid.Children {
				if c.Kind != schema.KindAttribute {
					kept = append(kept, c)
				}
			}
			kid.Children = kept
			found = true
		}
		occurrence++
	}
	if !found {
		t.Fatal("get_filters's response no longer has a second <filters> child to weaken; did the snapshot change?")
	}

	if _, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist); err == nil {
		t.Fatal("Generate succeeded with an attribute-less collision member; want an error")
	} else if !strings.Contains(err.Error(), "no attribute that uniquely identifies it") {
		t.Fatalf("Generate error = %q, want it to mention a missing discriminating attribute", err)
	}
}

// TestGetFiltersCollisionRoundTrips proves, against a real captured response
// and a battery of synthetic ones, that GetFiltersResponse's generated
// UnmarshalXML/MarshalXML (get_filters.go's only collided shape — see
// generate.go's emitCollisionMethods) separates the two sibling <filters>
// elements by attribute signature, not by arrival order. It reuses
// TestGeneratedCodeCompiles's scratch-module machinery, then goes one step
// further: a `go run` package main, in the same module, that imports the
// generated package and asserts on stdout (it is the only thing that can
// import the generated types, so assertions have to run inside it).
func TestGetFiltersCollisionRoundTrips(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("Skipping: no go toolchain on PATH (%v)", err)
	}
	got := generateOrFail(t)
	if _, ok := got[CommandFile("get_filters")]; !ok {
		t.Fatal("get_filters was not generated, so this test would not compile it")
	}

	fixture, err := filepath.Abs("../gmptest/testdata/live/get_filters/001_response.xml")
	if err != nil {
		t.Fatalf("Failed to resolve the fixture path: %v", err)
	}
	if _, err := os.Stat(fixture); err != nil {
		t.Fatalf("Fixture not found: %v", err)
	}

	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to resolve the module root: %v", err)
	}
	dir := t.TempDir()
	pkgDir := filepath.Join(dir, "genproto")
	if err := os.MkdirAll(pkgDir, 0o750); err != nil {
		t.Fatalf("Failed to create the scratch package: %v", err)
	}

	goVersion, err := goDirective(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatalf("Failed to read the module's go directive: %v", err)
	}
	gomod := "module gogmpcodegencheck\n\ngo " + goVersion + "\n\n" +
		"require github.com/brennoo/go-gmp v0.0.0\n\n" +
		"replace github.com/brennoo/go-gmp => " + root + "\n"
	mustWriteFile(t, filepath.Join(dir, "go.mod"), gomod)
	mustWriteFile(t, filepath.Join(pkgDir, "support.go"),
		strings.Replace(supportShim(), "package "+targetPackage+"\n", "package genproto\n", 1))
	for name, src := range got {
		rewritten := strings.Replace(src, "\npackage "+targetPackage+"\n", "\npackage genproto\n", 1)
		if rewritten == src {
			t.Fatalf("%s: expected a %q package clause to rewrite", name, targetPackage)
		}
		mustWriteFile(t, filepath.Join(pkgDir, name), rewritten)
	}
	mustWriteFile(t, filepath.Join(dir, "main.go"), filtersRoundTripProgram)

	cmd := exec.CommandContext(context.Background(), goBin, "run", ".", fixture)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Round-trip program failed: %v\n%s", err, out)
	}
	const want = "decode: Term=\"first=1 rows=10 sort=name\" Keywords=3 Start=1 Max=1000 Filters=0\n" +
		"roundtrip: Term=\"first=1 rows=10 sort=name\" Keywords=3 Start=1 Max=1000 Filters=0\n" +
		"zero-marshal-root: get_filters_response\n" +
		"reset: Status=201 FilterCount=0 Term=\"t2\"\n" +
		"extra-occurrence: GetFiltersResponse: unexpected extra <filters> element\n" +
		"reversed-order: Term=\"reversed\" Keywords=0 Start=1 Max=1000 Filters=0\n" +
		"solo-id: Term=\"solo-id\" Keywords=0 Start=0 Max=0 Filters=0\n" +
		"solo-pagination: Term=\"\" Keywords=0 Start=5 Max=50 Filters=0\n" +
		"unknown-signature: GetFiltersResponse: unrecognized <filters> element: no known variant's attributes are present\n" +
		"ambiguous-signature: GetFiltersResponse: ambiguous <filters> element: attributes match more than one known variant\n"
	if got := string(out); got != want {
		t.Fatalf("Round-trip program printed:\n%s\nwant:\n%s", got, want)
	}
}

// filtersRoundTripProgram is a `go run` main for TestGetFiltersCollisionRoundTrips.
// It decodes the real fixture named on argv[1] and round-trips it, then runs
// a battery of synthetic documents through the same decoder: a zero-value
// marshal, a receiver reused across two decodes, an extra <filters>
// occurrence, the two variants in reversed order, each variant alone, an
// unrecognized attribute signature, and an ambiguous one.
const filtersRoundTripProgram = `package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"gogmpcodegencheck/genproto"
)

func report(label string, r *genproto.GetFiltersResponse) {
	fmt.Printf("%s: Term=%q Keywords=%d Start=%d Max=%d Filters=%d\n",
		label, r.Filters.Term, len(r.Filters.Keywords.Keyword), r.Filters2.Start, r.Filters2.Max, len(r.Filter))
}

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var decoded genproto.GetFiltersResponse
	if err := xml.Unmarshal(data, &decoded); err != nil {
		fmt.Fprintln(os.Stderr, "decode:", err)
		os.Exit(1)
	}
	report("decode", &decoded)

	reencoded, err := xml.Marshal(&decoded)
	if err != nil {
		fmt.Fprintln(os.Stderr, "encode:", err)
		os.Exit(1)
	}
	var roundtripped genproto.GetFiltersResponse
	if err := xml.Unmarshal(reencoded, &roundtripped); err != nil {
		fmt.Fprintln(os.Stderr, "re-decode:", err)
		os.Exit(1)
	}
	report("roundtrip", &roundtripped)

	// A zero-value marshal must not produce a nameless root element.
	zeroBytes, err := xml.Marshal(&genproto.GetFiltersResponse{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "zero-marshal:", err)
		os.Exit(1)
	}
	var zeroName struct{ XMLName xml.Name }
	if err := xml.Unmarshal(zeroBytes, &zeroName); err != nil {
		fmt.Fprintln(os.Stderr, "zero-marshal-decode:", err)
		os.Exit(1)
	}
	fmt.Printf("zero-marshal-root: %s\n", zeroName.XMLName.Local)

	// Decoding into an already-populated receiver must not leave data behind
	// from the prior decode.
	const doc1 = "<get_filters_response status=\"200\" status_text=\"OK\">" +
		"<filters id=\"a\"><term>t1</term></filters>" +
		"<sort><field>name<order>ascending</order></field></sort>" +
		"<filters start=\"1\" max=\"10\"/>" +
		"<filter_count>0<filtered>3</filtered><page>0</page></filter_count>" +
		"<filter id=\"f1\"><name>one</name></filter>" +
		"<filter id=\"f2\"><name>two</name></filter>" +
		"</get_filters_response>"
	const doc2 = "<get_filters_response status=\"201\" status_text=\"Created\">" +
		"<filters id=\"b\"><term>t2</term></filters>" +
		"<sort><field>name<order>ascending</order></field></sort>" +
		"<filters start=\"2\" max=\"20\"/>" +
		"<filter_count>0<filtered>0</filtered><page>0</page></filter_count>" +
		"</get_filters_response>"
	var reused genproto.GetFiltersResponse
	if err := xml.Unmarshal([]byte(doc1), &reused); err != nil {
		fmt.Fprintln(os.Stderr, "reset-doc1:", err)
		os.Exit(1)
	}
	if err := xml.Unmarshal([]byte(doc2), &reused); err != nil {
		fmt.Fprintln(os.Stderr, "reset-doc2:", err)
		os.Exit(1)
	}
	fmt.Printf("reset: Status=%s FilterCount=%d Term=%q\n", reused.Status, len(reused.Filter), reused.Filters.Term)

	// A third <filters> beyond the two schema-declared fields must fail
	// loudly, not be silently skipped.
	const doc3 = "<get_filters_response status=\"200\" status_text=\"OK\">" +
		"<filters id=\"a\"><term>t1</term></filters>" +
		"<filters start=\"1\" max=\"10\"/>" +
		"<filters start=\"2\" max=\"20\"/>" +
		"</get_filters_response>"
	var extra genproto.GetFiltersResponse
	if err := xml.Unmarshal([]byte(doc3), &extra); err == nil {
		fmt.Println("extra-occurrence: no error (want one)")
		os.Exit(1)
	} else {
		fmt.Printf("extra-occurrence: %v\n", err)
	}

	// Reversed schema order must still land by attribute signature, not by
	// arrival order.
	const docReversed = "<get_filters_response status=\"200\" status_text=\"OK\">" +
		"<filters start=\"1\" max=\"1000\"/>" +
		"<filters id=\"a\"><term>reversed</term></filters>" +
		"</get_filters_response>"
	var reversed genproto.GetFiltersResponse
	if err := xml.Unmarshal([]byte(docReversed), &reversed); err != nil {
		fmt.Fprintln(os.Stderr, "reversed-order:", err)
		os.Exit(1)
	}
	report("reversed-order", &reversed)

	// Either variant alone, with the other omitted, must still decode.
	const docSoloID = "<get_filters_response status=\"200\" status_text=\"OK\">" +
		"<filters id=\"a\"><term>solo-id</term></filters>" +
		"</get_filters_response>"
	var soloID genproto.GetFiltersResponse
	if err := xml.Unmarshal([]byte(docSoloID), &soloID); err != nil {
		fmt.Fprintln(os.Stderr, "solo-id:", err)
		os.Exit(1)
	}
	report("solo-id", &soloID)

	const docSoloPagination = "<get_filters_response status=\"200\" status_text=\"OK\">" +
		"<filters start=\"5\" max=\"50\"/>" +
		"</get_filters_response>"
	var soloPagination genproto.GetFiltersResponse
	if err := xml.Unmarshal([]byte(docSoloPagination), &soloPagination); err != nil {
		fmt.Fprintln(os.Stderr, "solo-pagination:", err)
		os.Exit(1)
	}
	report("solo-pagination", &soloPagination)

	// A start tag matching no known variant must fail closed.
	const docUnknown = "<get_filters_response status=\"200\" status_text=\"OK\">" +
		"<filters foo=\"bar\"/>" +
		"</get_filters_response>"
	var unknown genproto.GetFiltersResponse
	if err := xml.Unmarshal([]byte(docUnknown), &unknown); err == nil {
		fmt.Println("unknown-signature: no error (want one)")
		os.Exit(1)
	} else {
		fmt.Printf("unknown-signature: %v\n", err)
	}

	// A start tag matching both known variants must fail closed too.
	const docAmbiguous = "<get_filters_response status=\"200\" status_text=\"OK\">" +
		"<filters id=\"a\" start=\"1\"/>" +
		"</get_filters_response>"
	var ambiguous genproto.GetFiltersResponse
	if err := xml.Unmarshal([]byte(docAmbiguous), &ambiguous); err == nil {
		fmt.Println("ambiguous-signature: no error (want one)")
		os.Exit(1)
	} else {
		fmt.Printf("ambiguous-signature: %v\n", err)
	}
}
`

// TestTLSCertificateGeneratedRoundTrips proves the real generated
// GetTLSCertificatesResponse agrees with the two FieldOverride entries in
// overrides.go (T-9 finding 3) — that certificate data, format, and validity
// survive a decode/re-encode/re-decode round trip through the generated type
// for both the upstream bundled example and a redacted live capture.
//
// ModifyTLSCertificate's omit/explicit-empty/set wire shape (finding 4) used
// to be pinned here too, as a generator-level stand-in for
// internal/gmptest/cases_tls_certificate.go, which could not exist until
// tls_certificate was adopted. Now that it is adopted, that coverage lives in
// the real cases_tls_certificate.go instead.
func TestTLSCertificateGeneratedRoundTrips(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("Skipping: no go toolchain on PATH (%v)", err)
	}
	got := generateOrFail(t)
	if _, ok := got[CommandFile("get_tls_certificates")]; !ok {
		t.Fatal("get_tls_certificates was not generated, so this test would not compile it")
	}

	upstreamFixture, err := filepath.Abs("../gmptest/testdata/upstream/get_tls_certificates/001_response.xml")
	if err != nil {
		t.Fatalf("Failed to resolve the upstream fixture path: %v", err)
	}
	if _, err := os.Stat(upstreamFixture); err != nil {
		t.Fatalf("Upstream fixture not found: %v", err)
	}
	liveFixture, err := filepath.Abs("../gmptest/testdata/live/get_tls_certificates/004_response.xml")
	if err != nil {
		t.Fatalf("Failed to resolve the live fixture path: %v", err)
	}
	if _, err := os.Stat(liveFixture); err != nil {
		t.Fatalf("Live fixture not found: %v", err)
	}

	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to resolve the module root: %v", err)
	}
	dir := t.TempDir()
	pkgDir := filepath.Join(dir, "genproto")
	if err := os.MkdirAll(pkgDir, 0o750); err != nil {
		t.Fatalf("Failed to create the scratch package: %v", err)
	}

	goVersion, err := goDirective(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatalf("Failed to read the module's go directive: %v", err)
	}
	gomod := "module gogmpcodegencheck\n\ngo " + goVersion + "\n\n" +
		"require github.com/brennoo/go-gmp v0.0.0\n\n" +
		"replace github.com/brennoo/go-gmp => " + root + "\n"
	mustWriteFile(t, filepath.Join(dir, "go.mod"), gomod)
	mustWriteFile(t, filepath.Join(pkgDir, "support.go"),
		strings.Replace(supportShim(), "package "+targetPackage+"\n", "package genproto\n", 1))
	for name, src := range got {
		rewritten := strings.Replace(src, "\npackage "+targetPackage+"\n", "\npackage genproto\n", 1)
		if rewritten == src {
			t.Fatalf("%s: expected a %q package clause to rewrite", name, targetPackage)
		}
		mustWriteFile(t, filepath.Join(pkgDir, name), rewritten)
	}
	mustWriteFile(t, filepath.Join(dir, "main.go"), tlsCertificateRoundTripProgram)

	cmd := exec.CommandContext(context.Background(), goBin, "run", ".", upstreamFixture, liveFixture)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Round-trip program failed: %v\n%s", err, out)
	}

	const upstreamCert = `MIIDNjCCAp+gAwIBAgIBATANBgkqhkiG9w0BAQQFADCBqTELM[...]`
	// Must match gmptest.RedactionPlaceholder (internal/gmptest/redact.go).
	const liveCert = `[REDACTED-BY-GO-GMP-CAPTURE]`
	want := "upstream-decode: Value=\"" + upstreamCert + "\" Format=\"DER\" Valid=true\n" +
		"upstream-roundtrip: Value=\"" + upstreamCert + "\" Format=\"DER\" Valid=true\n" +
		"live-decode: Value=\"" + liveCert + "\" Format=\"PEM\" Valid=true\n" +
		"live-roundtrip: Value=\"" + liveCert + "\" Format=\"PEM\" Valid=true\n"
	if got := string(out); got != want {
		t.Fatalf("Round-trip program printed:\n%s\nwant:\n%s", got, want)
	}
}

// tlsCertificateRoundTripProgram is a `go run` main for
// TestTLSCertificateGeneratedRoundTrips. argv[1] is the upstream fixture,
// argv[2] the live one; it decodes and round-trips each through the real
// generated GetTLSCertificatesResponse.
const tlsCertificateRoundTripProgram = `package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"gogmpcodegencheck/genproto"
)

func report(label string, r *genproto.GetTLSCertificatesResponse) {
	c := r.TLSCertificate[0]
	fmt.Printf("%s: Value=%q Format=%q Valid=%t\n", label, c.Certificate.Value, c.Certificate.Format, bool(c.Valid))
}

func decodeAndRoundTrip(path, decodeLabel, roundtripLabel string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var decoded genproto.GetTLSCertificatesResponse
	if err := xml.Unmarshal(data, &decoded); err != nil {
		fmt.Fprintln(os.Stderr, "decode:", err)
		os.Exit(1)
	}
	report(decodeLabel, &decoded)

	reencoded, err := xml.Marshal(&decoded)
	if err != nil {
		fmt.Fprintln(os.Stderr, "encode:", err)
		os.Exit(1)
	}
	var roundtripped genproto.GetTLSCertificatesResponse
	if err := xml.Unmarshal(reencoded, &roundtripped); err != nil {
		fmt.Fprintln(os.Stderr, "re-decode:", err)
		os.Exit(1)
	}
	report(roundtripLabel, &roundtripped)
}

func main() {
	decodeAndRoundTrip(os.Args[1], "upstream-decode", "upstream-roundtrip")
	decodeAndRoundTrip(os.Args[2], "live-decode", "live-roundtrip")
}
`

// ticketScratchModule copies the generated ticket family into a throwaway
// module under dir/genproto, the same technique
// TestTLSCertificateGeneratedRoundTrips uses, and returns the go binary path.
// internal/gmptest/cases_ticket.go still targets the hand-written commands
// package (ticket is review-only), so this is the only way to exercise the
// real generated types before adoption.
func ticketScratchModule(t *testing.T, dir string) string {
	t.Helper()
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("Skipping: no go toolchain on PATH (%v)", err)
	}
	got := generateOrFail(t)
	for _, name := range []string{CommandFile("create_ticket"), CommandFile("modify_ticket"), CommandFile("get_tickets")} {
		if _, ok := got[name]; !ok {
			t.Fatalf("%s was not generated, so this test would not compile it", name)
		}
	}

	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to resolve the module root: %v", err)
	}
	pkgDir := filepath.Join(dir, "genproto")
	if err := os.MkdirAll(pkgDir, 0o750); err != nil {
		t.Fatalf("Failed to create the scratch package: %v", err)
	}
	goVersion, err := goDirective(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatalf("Failed to read the module's go directive: %v", err)
	}
	gomod := "module gogmpcodegencheck\n\ngo " + goVersion + "\n\n" +
		"require github.com/brennoo/go-gmp v0.0.0\n\n" +
		"replace github.com/brennoo/go-gmp => " + root + "\n"
	mustWriteFile(t, filepath.Join(dir, "go.mod"), gomod)
	mustWriteFile(t, filepath.Join(pkgDir, "support.go"),
		strings.Replace(supportShim(), "package "+targetPackage+"\n", "package genproto\n", 1))
	for name, src := range got {
		rewritten := strings.Replace(src, "\npackage "+targetPackage+"\n", "\npackage genproto\n", 1)
		if rewritten == src {
			t.Fatalf("%s: expected a %q package clause to rewrite", name, targetPackage)
		}
		mustWriteFile(t, filepath.Join(pkgDir, name), rewritten)
	}
	return goBin
}

// TestTicketGeneratedRequestRoundTrips is T-9 ticket review finding 2:
// internal/gmptest/cases_ticket.go still uses the hand-written Assigned
// *TicketUser field, so it cannot exercise the generated AssignedTo shape or
// the PointerOverrides added for modify_ticket's note fields. This proves
// the real generated CreateTicket/ModifyTicket/GetTickets request types
// against that experiment's conclusions directly.
func TestTicketGeneratedRequestRoundTrips(t *testing.T) {
	dir := t.TempDir()
	goBin := ticketScratchModule(t, dir)
	mustWriteFile(t, filepath.Join(dir, "main.go"), ticketRequestRoundTripProgram)

	cmd := exec.CommandContext(context.Background(), goBin, "run", ".")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Round-trip program failed: %v\n%s", err, out)
	}

	want := "" +
		"create-ticket-assigned-to: <create_ticket><result id=\"result-1\"></result><assigned_to><user id=\"user-1\"></user></assigned_to><open_note>opened</open_note></create_ticket>\n" +
		"modify-ticket-omitted: <modify_ticket ticket_id=\"ticket-1\"></modify_ticket>\n" +
		"modify-ticket-assigned-to-set: <modify_ticket ticket_id=\"ticket-1\"><assigned_to><user id=\"user-2\"></user></assigned_to></modify_ticket>\n" +
		"modify-ticket-comment-explicit-empty: <modify_ticket ticket_id=\"ticket-1\"><comment></comment></modify_ticket>\n" +
		"modify-ticket-comment-set: <modify_ticket ticket_id=\"ticket-1\"><comment>a comment</comment></modify_ticket>\n" +
		"modify-ticket-status-set: <modify_ticket ticket_id=\"ticket-1\"><status>Open</status></modify_ticket>\n" +
		"get-tickets-trash-unset: <get_tickets></get_tickets>\n" +
		"get-tickets-trash-false: <get_tickets trash=\"0\"></get_tickets>\n" +
		"get-tickets-trash-true: <get_tickets trash=\"1\"></get_tickets>\n" +
		"get-tickets-filt-id: <get_tickets filt_id=\"filter-1\"></get_tickets>\n"
	if got := string(out); got != want {
		t.Fatalf("Round-trip program printed:\n%s\nwant:\n%s", got, want)
	}
}

// ticketRequestRoundTripProgram is a `go run` main for
// TestTicketGeneratedRequestRoundTrips. modify-ticket-omitted proves the
// omitted-vs-pointer distinction PointerOverrides relies on: nil drops the
// field entirely, a pointer to "" sends explicit-empty, matching the live
// gvmd experiment in docs/review-2026-08.md.
const ticketRequestRoundTripProgram = `package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"gogmpcodegencheck/genproto"
)

func printXML(label string, v any) {
	b, err := xml.Marshal(v)
	if err != nil {
		fmt.Fprintln(os.Stderr, label, err)
		os.Exit(1)
	}
	fmt.Printf("%s: %s\n", label, b)
}

func strp(s string) *string { return &s }
func boolp(b genproto.GMPBool) *genproto.GMPBool { return &b }

func main() {
	printXML("create-ticket-assigned-to", &genproto.CreateTicket{
		Result:     genproto.TicketResult{ID: "result-1"},
		AssignedTo: genproto.TicketAssignedTo{User: genproto.TicketUser{ID: "user-1"}},
		OpenNote:   "opened",
	})

	printXML("modify-ticket-omitted", &genproto.ModifyTicket{TicketID: "ticket-1"})
	printXML("modify-ticket-assigned-to-set", &genproto.ModifyTicket{
		TicketID:   "ticket-1",
		AssignedTo: &genproto.TicketAssignedTo{User: genproto.TicketUser{ID: "user-2"}},
	})
	printXML("modify-ticket-comment-explicit-empty", &genproto.ModifyTicket{TicketID: "ticket-1", Comment: strp("")})
	printXML("modify-ticket-comment-set", &genproto.ModifyTicket{TicketID: "ticket-1", Comment: strp("a comment")})
	printXML("modify-ticket-status-set", &genproto.ModifyTicket{TicketID: "ticket-1", Status: "Open"})

	printXML("get-tickets-trash-unset", &genproto.GetTickets{})
	printXML("get-tickets-trash-false", &genproto.GetTickets{Trash: boolp(false)})
	printXML("get-tickets-trash-true", &genproto.GetTickets{Trash: boolp(true)})
	printXML("get-tickets-filt-id", &genproto.GetTickets{FiltID: "filter-1"})
}
`

// TestTicketGeneratedResponseRoundTrips is T-9 ticket review finding 3: the
// upstream and live get_tickets fixtures on file confirm every documented
// generation change except two paths neither fixture exercises
// (user_tags/nvt), so this adds a synthetic document for those alongside the
// real fixtures.
func TestTicketGeneratedResponseRoundTrips(t *testing.T) {
	dir := t.TempDir()
	goBin := ticketScratchModule(t, dir)

	upstreamFixture, err := filepath.Abs("../gmptest/testdata/upstream/get_tickets/001_response.xml")
	if err != nil {
		t.Fatalf("Failed to resolve the upstream fixture path: %v", err)
	}
	if _, err := os.Stat(upstreamFixture); err != nil {
		t.Fatalf("Upstream fixture not found: %v", err)
	}
	liveFixture, err := filepath.Abs("../gmptest/testdata/live/get_tickets/001_response.xml")
	if err != nil {
		t.Fatalf("Failed to resolve the live fixture path: %v", err)
	}
	if _, err := os.Stat(liveFixture); err != nil {
		t.Fatalf("Live fixture not found: %v", err)
	}
	syntheticFixture := filepath.Join(dir, "synthetic.xml")
	mustWriteFile(t, syntheticFixture, ticketSyntheticResponse)

	mustWriteFile(t, filepath.Join(dir, "main.go"), ticketResponseRoundTripProgram)

	cmd := exec.CommandContext(context.Background(), goBin, "run", ".", upstreamFixture, liveFixture, syntheticFixture)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Round-trip program failed: %v\n%s", err, out)
	}

	want := "" +
		"upstream-decode: tickets=1 filtered=2 page=2\n" +
		"  ticket id=\"93cd2f71-48c3-4cf2-b542-5b256f59cae0\" assigned_to.user.id=\"33e92d3e-a379-4c46-a4cf-88c8201ab710\" writable=true in_use=false orphan=false nvt.oid=\"\"\n" +
		"upstream-roundtrip: tickets=1 filtered=2 page=2\n" +
		"  ticket id=\"93cd2f71-48c3-4cf2-b542-5b256f59cae0\" assigned_to.user.id=\"33e92d3e-a379-4c46-a4cf-88c8201ab710\" writable=true in_use=false orphan=false nvt.oid=\"\"\n" +
		"live-decode: tickets=0 filtered=0 page=0\n" +
		"live-roundtrip: tickets=0 filtered=0 page=0\n" +
		"synthetic-decode: tickets=1 filtered=1 page=1\n" +
		"  ticket id=\"ticket-syn-1\" assigned_to.user.id=\"user-syn-1\" writable=false in_use=false orphan=true nvt.oid=\"1.3.6.1.4.1.25623.1.0.12345\" user_tags.count=2 tag.comment=\"a tag comment\"\n" +
		"synthetic-roundtrip: tickets=1 filtered=1 page=1\n" +
		"  ticket id=\"ticket-syn-1\" assigned_to.user.id=\"user-syn-1\" writable=false in_use=false orphan=true nvt.oid=\"1.3.6.1.4.1.25623.1.0.12345\" user_tags.count=2 tag.comment=\"a tag comment\"\n"
	if got := string(out); got != want {
		t.Fatalf("Round-trip program printed:\n%s\nwant:\n%s", got, want)
	}
}

// ticketSyntheticResponse covers the two paths neither committed get_tickets
// fixture exercises: nvt and user_tags.
const ticketSyntheticResponse = `<get_tickets_response status="200" status_text="OK">
  <ticket id="ticket-syn-1">
    <nvt oid="1.3.6.1.4.1.25623.1.0.12345"/>
    <assigned_to><user id="user-syn-1"/></assigned_to>
    <orphan>1</orphan>
    <user_tags>
      <count>2</count>
      <tag id="tag-1"><name>tag-name</name><value>tag-value</value><comment>a tag comment</comment></tag>
    </user_tags>
  </ticket>
  <tickets start="1" max="10"/>
  <ticket_count>1<filtered>1</filtered><page>1</page></ticket_count>
</get_tickets_response>
`

// ticketResponseRoundTripProgram is a `go run` main for
// TestTicketGeneratedResponseRoundTrips. argv[1] is the upstream fixture,
// argv[2] the live one, argv[3] the synthetic one.
const ticketResponseRoundTripProgram = `package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"gogmpcodegencheck/genproto"
)

func report(label string, r *genproto.GetTicketsResponse) {
	fmt.Printf("%s: tickets=%d filtered=%d page=%d\n", label, len(r.Ticket), r.TicketCount.Filtered, r.TicketCount.Page)
	for _, t := range r.Ticket {
		fmt.Printf("  ticket id=%q assigned_to.user.id=%q writable=%t in_use=%t orphan=%t nvt.oid=%q",
			t.ID, t.AssignedTo.User.ID, bool(t.Writable), bool(t.InUse), bool(t.Orphan), t.NVT.OID)
		if t.UserTags != nil {
			fmt.Printf(" user_tags.count=%d", t.UserTags.Count)
			for _, tag := range t.UserTags.Tag {
				fmt.Printf(" tag.comment=%q", tag.Comment)
			}
		}
		fmt.Println()
	}
}

func decodeAndRoundTrip(path, decodeLabel, roundtripLabel string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var decoded genproto.GetTicketsResponse
	if err := xml.Unmarshal(data, &decoded); err != nil {
		fmt.Fprintln(os.Stderr, "decode:", err)
		os.Exit(1)
	}
	report(decodeLabel, &decoded)

	reencoded, err := xml.Marshal(&decoded)
	if err != nil {
		fmt.Fprintln(os.Stderr, "encode:", err)
		os.Exit(1)
	}
	var roundtripped genproto.GetTicketsResponse
	if err := xml.Unmarshal(reencoded, &roundtripped); err != nil {
		fmt.Fprintln(os.Stderr, "re-decode:", err)
		os.Exit(1)
	}
	report(roundtripLabel, &roundtripped)
}

func main() {
	decodeAndRoundTrip(os.Args[1], "upstream-decode", "upstream-roundtrip")
	decodeAndRoundTrip(os.Args[2], "live-decode", "live-roundtrip")
	decodeAndRoundTrip(os.Args[3], "synthetic-decode", "synthetic-roundtrip")
}
`

// credentialScratchModule copies the generated credential family into a
// throwaway module under dir/genproto, the same technique
// ticketScratchModule uses. internal/gmptest/cases_credential.go does not
// exist yet (credential is review-only), and the existing goldens compile
// against the hand-written commands package, not the generated one, so this
// is the only way to exercise the real generated Kdc/Kdcs/Key/Privacy shapes
// and the modify_credential.comment PointerOverride before adoption.
func credentialScratchModule(t *testing.T, dir string) string {
	t.Helper()
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("Skipping: no go toolchain on PATH (%v)", err)
	}
	got := generateOrFail(t)
	for _, name := range []string{CommandFile("create_credential"), CommandFile("modify_credential"), CommandFile("get_credentials")} {
		if _, ok := got[name]; !ok {
			t.Fatalf("%s was not generated, so this test would not compile it", name)
		}
	}

	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to resolve the module root: %v", err)
	}
	pkgDir := filepath.Join(dir, "genproto")
	if err := os.MkdirAll(pkgDir, 0o750); err != nil {
		t.Fatalf("Failed to create the scratch package: %v", err)
	}
	goVersion, err := goDirective(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatalf("Failed to read the module's go directive: %v", err)
	}
	gomod := "module gogmpcodegencheck\n\ngo " + goVersion + "\n\n" +
		"require github.com/brennoo/go-gmp v0.0.0\n\n" +
		"replace github.com/brennoo/go-gmp => " + root + "\n"
	mustWriteFile(t, filepath.Join(dir, "go.mod"), gomod)
	mustWriteFile(t, filepath.Join(pkgDir, "support.go"),
		strings.Replace(supportShim(), "package "+targetPackage+"\n", "package genproto\n", 1))
	for name, src := range got {
		rewritten := strings.Replace(src, "\npackage "+targetPackage+"\n", "\npackage genproto\n", 1)
		if rewritten == src {
			t.Fatalf("%s: expected a %q package clause to rewrite", name, targetPackage)
		}
		mustWriteFile(t, filepath.Join(pkgDir, name), rewritten)
	}
	return goBin
}

// TestCredentialGeneratedRequestRoundTrips is T-9 credential review finding
// 1: generation renames KDC/KDCs to Kdc/Kdcs, dedups CreateCredentialKDCs and
// CreateCredentialKey into the shared CredentialKdcs/CredentialKey, and drops
// CredentialPrivacy.Algorithm/.Password's omitempty (they are required
// elements, not optional ones, per the schema). None of that is exercised by
// the hand-written-type goldens, so this proves the real generated request
// shapes directly, including the modify_credential.comment PointerOverride
// added from this family's live gvmd experiment.
func TestCredentialGeneratedRequestRoundTrips(t *testing.T) {
	dir := t.TempDir()
	goBin := credentialScratchModule(t, dir)
	mustWriteFile(t, filepath.Join(dir, "main.go"), credentialRequestRoundTripProgram)

	cmd := exec.CommandContext(context.Background(), goBin, "run", ".")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Round-trip program failed: %v\n%s", err, out)
	}

	want := "" +
		"create-credential-kdc-kdcs-key-privacy: <create_credential><name>cred-1</name><kdc>kdc0.example.com</kdc><kdcs><kdc>kdc1.example.com</kdc><kdc>kdc2.example.com</kdc></kdcs><key><phrase>phrase1</phrase><private>priv1</private><public>pub1</public></key><privacy><algorithm>aes</algorithm><password>privpass</password></privacy><type>krb5</type></create_credential>\n" +
		"modify-credential-kdc-kdcs-key-privacy: <modify_credential credential_id=\"cred-1\"><kdc>kdc0.example.com</kdc><kdcs><kdc>kdc1.example.com</kdc><kdc>kdc2.example.com</kdc></kdcs><key><phrase>phrase1</phrase><private>priv1</private><public>pub1</public></key><privacy><algorithm>aes</algorithm><password>privpass</password></privacy></modify_credential>\n" +
		"create-credential-privacy-zero-valued: <create_credential><name>cred-2</name><privacy><algorithm></algorithm><password></password></privacy><type>snmp</type></create_credential>\n" +
		"create-credential-allow-insecure-unset: <create_credential><name>cred-3</name></create_credential>\n" +
		"create-credential-allow-insecure-false: <create_credential><name>cred-3</name><allow_insecure>0</allow_insecure></create_credential>\n" +
		"create-credential-allow-insecure-true: <create_credential><name>cred-3</name><allow_insecure>1</allow_insecure></create_credential>\n" +
		"modify-credential-comment-omitted: <modify_credential credential_id=\"cred-1\"></modify_credential>\n" +
		"modify-credential-comment-explicit-empty: <modify_credential credential_id=\"cred-1\"><comment></comment></modify_credential>\n" +
		"modify-credential-comment-set: <modify_credential credential_id=\"cred-1\"><comment>a comment</comment></modify_credential>\n"
	if got := string(out); got != want {
		t.Fatalf("Round-trip program printed:\n%s\nwant:\n%s", got, want)
	}
}

// credentialRequestRoundTripProgram is a `go run` main for
// TestCredentialGeneratedRequestRoundTrips. The zero-valued-privacy case
// pins that CredentialPrivacy.Algorithm/.Password are required-but-empty:
// generation dropped their omitempty tags, so a non-nil *CredentialPrivacy
// with both fields at their zero value still emits both elements rather than
// omitting them. modify-credential-comment-omitted/-explicit-empty prove the
// PointerOverride's nil-vs-pointer-to-empty-string distinction, matching the
// live gvmd experiment recorded in overrides.go and docs/review-2026-08.md.
const credentialRequestRoundTripProgram = `package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"gogmpcodegencheck/genproto"
)

func printXML(label string, v any) {
	b, err := xml.Marshal(v)
	if err != nil {
		fmt.Fprintln(os.Stderr, label, err)
		os.Exit(1)
	}
	fmt.Printf("%s: %s\n", label, b)
}

func strp(s string) *string { return &s }
func boolp(b genproto.GMPBool) *genproto.GMPBool { return &b }

func main() {
	printXML("create-credential-kdc-kdcs-key-privacy", &genproto.CreateCredential{
		Name: "cred-1",
		Type: "krb5",
		Kdc:  "kdc0.example.com",
		Kdcs: &genproto.CredentialKdcs{Kdc: []string{"kdc1.example.com", "kdc2.example.com"}},
		Key:  &genproto.CredentialKey{Phrase: "phrase1", Private: "priv1", Public: "pub1"},
		Privacy: &genproto.CredentialPrivacy{Algorithm: "aes", Password: "privpass"},
	})
	printXML("modify-credential-kdc-kdcs-key-privacy", &genproto.ModifyCredential{
		CredentialID: "cred-1",
		Kdc:          "kdc0.example.com",
		Kdcs:         &genproto.CredentialKdcs{Kdc: []string{"kdc1.example.com", "kdc2.example.com"}},
		Key:          &genproto.CredentialKey{Phrase: "phrase1", Private: "priv1", Public: "pub1"},
		Privacy:      &genproto.CredentialPrivacy{Algorithm: "aes", Password: "privpass"},
	})
	printXML("create-credential-privacy-zero-valued", &genproto.CreateCredential{
		Name: "cred-2", Type: "snmp", Privacy: &genproto.CredentialPrivacy{},
	})
	printXML("create-credential-allow-insecure-unset", &genproto.CreateCredential{Name: "cred-3"})
	printXML("create-credential-allow-insecure-false", &genproto.CreateCredential{Name: "cred-3", AllowInsecure: boolp(false)})
	printXML("create-credential-allow-insecure-true", &genproto.CreateCredential{Name: "cred-3", AllowInsecure: boolp(true)})

	printXML("modify-credential-comment-omitted", &genproto.ModifyCredential{CredentialID: "cred-1"})
	printXML("modify-credential-comment-explicit-empty", &genproto.ModifyCredential{CredentialID: "cred-1", Comment: strp("")})
	printXML("modify-credential-comment-set", &genproto.ModifyCredential{CredentialID: "cred-1", Comment: strp("a comment")})
}
`

// TestCredentialGeneratedResponseRoundTrips is T-9 credential review finding
// 2: the conformance gate validates schema structure but never proves
// GetCredentialsResponse round-trips real fixtures. Committed fixtures cover
// repeated Credential, Formats, AllowInsecure/Writable/InUse, target entries
// and package chardata+format, but none exercises Kdcs, scanner entries or
// privacy/certificate_info, so this adds a synthetic document pinning every
// field of CredentialCertificateInfo and GetCredentialsPrivacy.Algorithm
// alongside the real fixtures.
func TestCredentialGeneratedResponseRoundTrips(t *testing.T) {
	dir := t.TempDir()
	goBin := credentialScratchModule(t, dir)

	upstreamRepeated, err := filepath.Abs("../gmptest/testdata/upstream/get_credentials/001_response.xml")
	if err != nil {
		t.Fatalf("Failed to resolve the upstream fixture path: %v", err)
	}
	if _, err := os.Stat(upstreamRepeated); err != nil {
		t.Fatalf("Upstream fixture not found: %v", err)
	}
	upstreamPackage, err := filepath.Abs("../gmptest/testdata/upstream/get_credentials/003_response.xml")
	if err != nil {
		t.Fatalf("Failed to resolve the upstream fixture path: %v", err)
	}
	if _, err := os.Stat(upstreamPackage); err != nil {
		t.Fatalf("Upstream fixture not found: %v", err)
	}
	liveFixture, err := filepath.Abs("../gmptest/testdata/live/get_credentials/004_response.xml")
	if err != nil {
		t.Fatalf("Failed to resolve the live fixture path: %v", err)
	}
	if _, err := os.Stat(liveFixture); err != nil {
		t.Fatalf("Live fixture not found: %v", err)
	}
	syntheticFixture := filepath.Join(dir, "synthetic.xml")
	mustWriteFile(t, syntheticFixture, credentialSyntheticResponse)

	mustWriteFile(t, filepath.Join(dir, "main.go"), credentialResponseRoundTripProgram)

	cmd := exec.CommandContext(context.Background(), goBin, "run", ".", upstreamRepeated, upstreamPackage, liveFixture, syntheticFixture)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Round-trip program failed: %v\n%s", err, out)
	}

	want := "" +
		"upstream-repeated-decode: credentials=2 filtered=0 page=0\n" +
		"  credential id=\"c33864a9-d3fd-44b3-8717-972bfb01dfcf\" type=\"usk\" formats=[key rpm deb] allow_insecure=false writable=true in_use=false\n" +
		"  credential id=\"8e305b0b-260d-450d-91a8-dadf0b144e15\" type=\"up\" formats=[exe] allow_insecure=false writable=true in_use=true\n" +
		"upstream-repeated-roundtrip: credentials=2 filtered=0 page=0\n" +
		"  credential id=\"c33864a9-d3fd-44b3-8717-972bfb01dfcf\" type=\"usk\" formats=[key rpm deb] allow_insecure=false writable=true in_use=false\n" +
		"  credential id=\"8e305b0b-260d-450d-91a8-dadf0b144e15\" type=\"up\" formats=[exe] allow_insecure=false writable=true in_use=true\n" +
		"upstream-package-decode: credentials=1 filtered=0 page=0\n" +
		"  credential id=\"8e305b0b-260d-450d-91a8-dadf0b144e15\" type=\"up\" target.id=\"1f28d970-17ef-4c69-ba8a-13827059f2b9\" package.format=\"deb\" package.len=29 allow_insecure=false writable=true in_use=true\n" +
		"upstream-package-roundtrip: credentials=1 filtered=0 page=0\n" +
		"  credential id=\"8e305b0b-260d-450d-91a8-dadf0b144e15\" type=\"up\" target.id=\"1f28d970-17ef-4c69-ba8a-13827059f2b9\" package.format=\"deb\" package.len=29 allow_insecure=false writable=true in_use=true\n" +
		"live-decode: credentials=1 filtered=1 page=1\n" +
		"  credential id=\"81c535ac-5b43-427d-ba8a-a865051eb096\" type=\"cc\" formats=[pem] allow_insecure=false writable=true in_use=false\n" +
		"live-roundtrip: credentials=1 filtered=1 page=1\n" +
		"  credential id=\"81c535ac-5b43-427d-ba8a-a865051eb096\" type=\"cc\" formats=[pem] allow_insecure=false writable=true in_use=false\n" +
		"synthetic-decode: credentials=1 filtered=1 page=1\n" +
		"  credential id=\"cred-syn-1\" type=\"krb5\" formats=[key] kdcs=[kdc1.example.com kdc2.example.com] scanner.id=\"scanner-1\" privacy.algorithm=\"aes\" cert_info.time_status=\"valid\" cert_info.activation_time=\"2020-01-01T00:00:00Z\" cert_info.expiration_time=\"2030-01-01T00:00:00Z\" cert_info.issuer=\"CN=test\" cert_info.md5=\"aa:bb\" cert_info.sha256=\"cc:dd\" cert_info.subject=\"CN=subject\" cert_info.serial=\"01\" allow_insecure=false writable=true in_use=false\n" +
		"synthetic-roundtrip: credentials=1 filtered=1 page=1\n" +
		"  credential id=\"cred-syn-1\" type=\"krb5\" formats=[key] kdcs=[kdc1.example.com kdc2.example.com] scanner.id=\"scanner-1\" privacy.algorithm=\"aes\" cert_info.time_status=\"valid\" cert_info.activation_time=\"2020-01-01T00:00:00Z\" cert_info.expiration_time=\"2030-01-01T00:00:00Z\" cert_info.issuer=\"CN=test\" cert_info.md5=\"aa:bb\" cert_info.sha256=\"cc:dd\" cert_info.subject=\"CN=subject\" cert_info.serial=\"01\" allow_insecure=false writable=true in_use=false\n"
	if got := string(out); got != want {
		t.Fatalf("Round-trip program printed:\n%s\nwant:\n%s", got, want)
	}
}

// credentialSyntheticResponse covers the paths no committed get_credentials
// fixture exercises: Kdcs, a scanner entry, privacy and a fully-populated
// certificate_info (every field, including the optional sha256_fingerprint,
// subject and serial).
const credentialSyntheticResponse = `<get_credentials_response status="200" status_text="OK">
  <credential id="cred-syn-1">
    <owner><name>admin</name></owner>
    <name>synthetic-krb5</name>
    <allow_insecure>0</allow_insecure>
    <login>krbuser</login>
    <comment></comment>
    <creation_time>2026-01-01T00:00:00Z</creation_time>
    <modification_time>2026-01-01T00:00:00Z</modification_time>
    <writable>1</writable>
    <in_use>0</in_use>
    <permissions><permission><name>Everything</name></permission></permissions>
    <type>krb5</type>
    <full_type>Kerberos</full_type>
    <formats><format>key</format></formats>
    <kdcs><kdc>kdc1.example.com</kdc><kdc>kdc2.example.com</kdc></kdcs>
    <scanners>
      <scanner id="scanner-1"><name>a scanner</name></scanner>
    </scanners>
    <privacy><algorithm>aes</algorithm></privacy>
    <certificate_info>
      <time_status>valid</time_status>
      <activation_time>2020-01-01T00:00:00Z</activation_time>
      <expiration_time>2030-01-01T00:00:00Z</expiration_time>
      <issuer>CN=test</issuer>
      <md5_fingerprint>aa:bb</md5_fingerprint>
      <sha256_fingerprint>cc:dd</sha256_fingerprint>
      <subject>CN=subject</subject>
      <serial>01</serial>
    </certificate_info>
  </credential>
  <credentials start="1" max="1000"/>
  <credential_count>1<filtered>1</filtered><page>1</page></credential_count>
</get_credentials_response>
`

// credentialResponseRoundTripProgram is a `go run` main for
// TestCredentialGeneratedResponseRoundTrips. argv[1] is the upstream
// repeated-credential fixture, argv[2] the upstream package/target fixture,
// argv[3] the live AllowInsecure/Writable/InUse fixture, argv[4] the
// synthetic one.
const credentialResponseRoundTripProgram = `package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"gogmpcodegencheck/genproto"
)

func report(label string, r *genproto.GetCredentialsResponse) {
	fmt.Printf("%s: credentials=%d filtered=%d page=%d\n", label, len(r.Credential), r.CredentialCount.Filtered, r.CredentialCount.Page)
	for _, c := range r.Credential {
		fmt.Printf("  credential id=%q type=%q", c.ID, c.Type)
		if len(c.Formats.Format) > 0 {
			fmt.Printf(" formats=%v", c.Formats.Format)
		}
		if c.Targets != nil {
			for _, tgt := range c.Targets.Target {
				fmt.Printf(" target.id=%q", tgt.ID)
			}
		}
		if c.Package != nil {
			fmt.Printf(" package.format=%q package.len=%d", c.Package.Format, len(c.Package.Value))
		}
		if c.Kdcs != nil && len(c.Kdcs.Kdc) > 0 {
			fmt.Printf(" kdcs=%v", c.Kdcs.Kdc)
		}
		if c.Scanners != nil {
			for _, s := range c.Scanners.Scanner {
				fmt.Printf(" scanner.id=%q", s.ID)
			}
		}
		if c.Privacy != nil {
			fmt.Printf(" privacy.algorithm=%q", c.Privacy.Algorithm)
		}
		if c.CertificateInfo != nil {
			ci := c.CertificateInfo
			fmt.Printf(" cert_info.time_status=%q cert_info.activation_time=%q cert_info.expiration_time=%q cert_info.issuer=%q cert_info.md5=%q cert_info.sha256=%q cert_info.subject=%q cert_info.serial=%q",
				ci.TimeStatus, ci.ActivationTime, ci.ExpirationTime, ci.Issuer, ci.MD5Fingerprint, ci.SHA256Fingerprint, ci.Subject, ci.Serial)
		}
		fmt.Printf(" allow_insecure=%t writable=%t in_use=%t", bool(c.AllowInsecure), bool(c.Writable), bool(c.InUse))
		fmt.Println()
	}
}

func decodeAndRoundTrip(path, decodeLabel, roundtripLabel string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var decoded genproto.GetCredentialsResponse
	if err := xml.Unmarshal(data, &decoded); err != nil {
		fmt.Fprintln(os.Stderr, "decode:", err)
		os.Exit(1)
	}
	report(decodeLabel, &decoded)

	reencoded, err := xml.Marshal(&decoded)
	if err != nil {
		fmt.Fprintln(os.Stderr, "encode:", err)
		os.Exit(1)
	}
	var roundtripped genproto.GetCredentialsResponse
	if err := xml.Unmarshal(reencoded, &roundtripped); err != nil {
		fmt.Fprintln(os.Stderr, "re-decode:", err)
		os.Exit(1)
	}
	report(roundtripLabel, &roundtripped)
}

func main() {
	decodeAndRoundTrip(os.Args[1], "upstream-repeated-decode", "upstream-repeated-roundtrip")
	decodeAndRoundTrip(os.Args[2], "upstream-package-decode", "upstream-package-roundtrip")
	decodeAndRoundTrip(os.Args[3], "live-decode", "live-roundtrip")
	decodeAndRoundTrip(os.Args[4], "synthetic-decode", "synthetic-roundtrip")
}
`

// portListScratchModule copies the generated port_list family into a
// throwaway module under dir/genproto, the same technique
// credentialScratchModule uses. internal/gmptest's cases compile against
// the hand-written commands package, so this is the only way to exercise
// the real generated types before adoption.
func portListScratchModule(t *testing.T, dir string) string {
	t.Helper()
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("Skipping: no go toolchain on PATH (%v)", err)
	}
	got := generateOrFail(t)
	for _, name := range []string{CommandFile("create_port_list"), CommandFile("modify_port_list"), CommandFile("get_port_lists")} {
		if _, ok := got[name]; !ok {
			t.Fatalf("%s was not generated, so this test would not compile it", name)
		}
	}

	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to resolve the module root: %v", err)
	}
	pkgDir := filepath.Join(dir, "genproto")
	if err := os.MkdirAll(pkgDir, 0o750); err != nil {
		t.Fatalf("Failed to create the scratch package: %v", err)
	}
	goVersion, err := goDirective(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatalf("Failed to read the module's go directive: %v", err)
	}
	gomod := "module gogmpcodegencheck\n\ngo " + goVersion + "\n\n" +
		"require github.com/brennoo/go-gmp v0.0.0\n\n" +
		"replace github.com/brennoo/go-gmp => " + root + "\n"
	mustWriteFile(t, filepath.Join(dir, "go.mod"), gomod)
	mustWriteFile(t, filepath.Join(pkgDir, "support.go"),
		strings.Replace(supportShim(), "package "+targetPackage+"\n", "package genproto\n", 1))
	for name, src := range got {
		rewritten := strings.Replace(src, "\npackage "+targetPackage+"\n", "\npackage genproto\n", 1)
		if rewritten == src {
			t.Fatalf("%s: expected a %q package clause to rewrite", name, targetPackage)
		}
		mustWriteFile(t, filepath.Join(pkgDir, name), rewritten)
	}
	return goBin
}

// TestPortListGeneratedRequestRoundTrips pins the GMP lexical "0"/"1" wire
// shape for GetPortLists.Details/Targets/Trash, plus the two ModifyPortList
// wire shapes the modify_port_list CommandWarning documents: Name and
// Comment set together, and either one omitted while the other is resent.
func TestPortListGeneratedRequestRoundTrips(t *testing.T) {
	dir := t.TempDir()
	goBin := portListScratchModule(t, dir)
	mustWriteFile(t, filepath.Join(dir, "main.go"), portListRequestRoundTripProgram)

	cmd := exec.CommandContext(context.Background(), goBin, "run", ".")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Round-trip program failed: %v\n%s", err, out)
	}

	want := "" +
		"create-port-list: <create_port_list><name>list-1</name><port_range>T:1-100</port_range><comment>a comment</comment><copy>copy-src-id</copy></create_port_list>\n" +
		"modify-port-list-both-set: <modify_port_list port_list_id=\"pl-1\"><name>renamed</name><comment>updated comment</comment></modify_port_list>\n" +
		"modify-port-list-name-omitted: <modify_port_list port_list_id=\"pl-1\"><comment>updated comment</comment></modify_port_list>\n" +
		"modify-port-list-comment-omitted: <modify_port_list port_list_id=\"pl-1\"><name>renamed</name></modify_port_list>\n" +
		"get-port-lists-unset: <get_port_lists port_list_id=\"pl-1\" filter=\"name=foo\" filt_id=\"filt-1\"></get_port_lists>\n" +
		"get-port-lists-false: <get_port_lists details=\"0\" targets=\"0\" trash=\"0\"></get_port_lists>\n" +
		"get-port-lists-true: <get_port_lists details=\"1\" targets=\"1\" trash=\"1\"></get_port_lists>\n"
	if got := string(out); got != want {
		t.Fatalf("Round-trip program printed:\n%s\nwant:\n%s", got, want)
	}
}

// portListRequestRoundTripProgram is a `go run` main for
// TestPortListGeneratedRequestRoundTrips.
const portListRequestRoundTripProgram = `package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"gogmpcodegencheck/genproto"
)

func printXML(label string, v any) {
	b, err := xml.Marshal(v)
	if err != nil {
		fmt.Fprintln(os.Stderr, label, err)
		os.Exit(1)
	}
	fmt.Printf("%s: %s\n", label, b)
}

func boolp(b genproto.GMPBool) *genproto.GMPBool { return &b }

func main() {
	printXML("create-port-list", &genproto.CreatePortList{
		Name: "list-1", PortRange: "T:1-100", Comment: "a comment", Copy: "copy-src-id",
	})
	printXML("modify-port-list-both-set", &genproto.ModifyPortList{
		PortListID: "pl-1", Name: "renamed", Comment: "updated comment",
	})
	printXML("modify-port-list-name-omitted", &genproto.ModifyPortList{
		PortListID: "pl-1", Comment: "updated comment",
	})
	printXML("modify-port-list-comment-omitted", &genproto.ModifyPortList{
		PortListID: "pl-1", Name: "renamed",
	})
	printXML("get-port-lists-unset", &genproto.GetPortLists{
		PortListID: "pl-1", Filter: "name=foo", FiltID: "filt-1",
	})
	printXML("get-port-lists-false", &genproto.GetPortLists{
		Details: boolp(false), Targets: boolp(false), Trash: boolp(false),
	})
	printXML("get-port-lists-true", &genproto.GetPortLists{
		Details: boolp(true), Targets: boolp(true), Trash: boolp(true),
	})
}
`

// TestPortListGeneratedResponseRoundTrips adds a synthetic document to cover
// deprecated/user_tags, which no committed fixture exercises. The schema
// declares chardata on sort, not sort/field, so Value captures formatting
// whitespace rather than the observed field name (see
// TestPortListSortFieldNameIsAnAcceptedLoss).
func TestPortListGeneratedResponseRoundTrips(t *testing.T) {
	dir := t.TempDir()
	goBin := portListScratchModule(t, dir)

	upstreamFixture, err := filepath.Abs("../gmptest/testdata/upstream/get_port_lists/001_response.xml")
	if err != nil {
		t.Fatalf("Failed to resolve the upstream fixture path: %v", err)
	}
	if _, err := os.Stat(upstreamFixture); err != nil {
		t.Fatalf("Upstream fixture not found: %v", err)
	}
	liveFixture, err := filepath.Abs("../gmptest/testdata/live/get_port_lists/001_response.xml")
	if err != nil {
		t.Fatalf("Failed to resolve the live fixture path: %v", err)
	}
	if _, err := os.Stat(liveFixture); err != nil {
		t.Fatalf("Live fixture not found: %v", err)
	}
	syntheticFixture := filepath.Join(dir, "synthetic.xml")
	mustWriteFile(t, syntheticFixture, portListSyntheticResponse)

	mustWriteFile(t, filepath.Join(dir, "main.go"), portListResponseRoundTripProgram)

	cmd := exec.CommandContext(context.Background(), goBin, "run", ".", upstreamFixture, liveFixture, syntheticFixture)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Round-trip program failed: %v\n%s", err, out)
	}

	want := "" +
		"upstream-decode: port_lists=1 filtered=0 page=0 sort_value=\"\" sort_field_order=\"\"\n" +
		"  port_list id=\"b993b6f5-f9fb-4e6e-9c94-dd46c00e058d\" name=\"All TCP\" writable=false in_use=true predefined=false deprecated=false port_ranges=[TCP:1-65535] targets=[1f28d970-17ef-4c69-ba8a-13827059f2b9]\n" +
		"upstream-roundtrip: port_lists=1 filtered=0 page=0 sort_value=\"\" sort_field_order=\"\"\n" +
		"  port_list id=\"b993b6f5-f9fb-4e6e-9c94-dd46c00e058d\" name=\"All TCP\" writable=false in_use=true predefined=false deprecated=false port_ranges=[TCP:1-65535] targets=[1f28d970-17ef-4c69-ba8a-13827059f2b9]\n" +
		"live-decode: port_lists=3 filtered=3 page=3 sort_value=\"\" sort_field_order=\"ascending\" keywords=[first=1 rows=10 sort=name]\n" +
		"  port_list id=\"33d0cd82-57c6-11e1-8ed1-406186ea4fc5\" name=\"All IANA assigned TCP\" writable=true in_use=false predefined=true deprecated=false\n" +
		"  port_list id=\"4a4717fe-57d2-11e1-9a26-406186ea4fc5\" name=\"All IANA assigned TCP and UDP\" writable=true in_use=false predefined=true deprecated=false\n" +
		"  port_list id=\"730ef368-57e2-11e1-a90f-406186ea4fc5\" name=\"All TCP and Nmap top 100 UDP\" writable=true in_use=false predefined=true deprecated=false\n" +
		"live-roundtrip: port_lists=3 filtered=3 page=3 sort_value=\"\" sort_field_order=\"ascending\" keywords=[first=1 rows=10 sort=name]\n" +
		"  port_list id=\"33d0cd82-57c6-11e1-8ed1-406186ea4fc5\" name=\"All IANA assigned TCP\" writable=true in_use=false predefined=true deprecated=false\n" +
		"  port_list id=\"4a4717fe-57d2-11e1-9a26-406186ea4fc5\" name=\"All IANA assigned TCP and UDP\" writable=true in_use=false predefined=true deprecated=false\n" +
		"  port_list id=\"730ef368-57e2-11e1-a90f-406186ea4fc5\" name=\"All TCP and Nmap top 100 UDP\" writable=true in_use=false predefined=true deprecated=false\n" +
		"synthetic-decode: port_lists=1 filtered=1 page=1 sort_value=\"\" sort_field_order=\"ascending\"\n" +
		"  port_list id=\"port-list-syn-1\" name=\"synthetic-port-list\" writable=true in_use=false predefined=false deprecated=true user_tags=[env:prod]\n" +
		"synthetic-roundtrip: port_lists=1 filtered=1 page=1 sort_value=\"\" sort_field_order=\"ascending\"\n" +
		"  port_list id=\"port-list-syn-1\" name=\"synthetic-port-list\" writable=true in_use=false predefined=false deprecated=true user_tags=[env:prod]\n"
	if got := string(out); got != want {
		t.Fatalf("Round-trip program printed:\n%s\nwant:\n%s", got, want)
	}
}

// portListSyntheticResponse covers the paths no committed get_port_lists
// fixture exercises: deprecated and a populated user_tags block.
const portListSyntheticResponse = `<get_port_lists_response status="200" status_text="OK">
  <port_list id="port-list-syn-1">
    <owner><name>admin</name></owner>
    <name>synthetic-port-list</name>
    <comment>covers deprecated and user_tags</comment>
    <creation_time>2026-01-01T00:00:00Z</creation_time>
    <modification_time>2026-01-01T00:00:00Z</modification_time>
    <writable>1</writable>
    <in_use>0</in_use>
    <permissions><permission><name>Everything</name></permission></permissions>
    <port_count><all>1</all><tcp>1</tcp><udp>0</udp></port_count>
    <predefined>0</predefined>
    <user_tags>
      <count>1</count>
      <tag id="tag-1"><name>env</name><value>prod</value><comment></comment></tag>
    </user_tags>
    <deprecated>1</deprecated>
  </port_list>
  <filters id=""><term>first=1 rows=10 sort=name</term><keywords/></filters>
  <sort><field>name<order>ascending</order></field></sort>
  <port_lists start="1" max="1000"/>
  <port_list_count>1<filtered>1</filtered><page>1</page></port_list_count>
</get_port_lists_response>
`

// portListResponseRoundTripProgram is a `go run` main for
// TestPortListGeneratedResponseRoundTrips. argv[1] is the upstream
// port_ranges/targets fixture, argv[2] the live repeated-port_list fixture,
// argv[3] the synthetic deprecated/user_tags fixture.
const portListResponseRoundTripProgram = `package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"gogmpcodegencheck/genproto"
)

func report(label string, r *genproto.GetPortListsResponse) {
	fmt.Printf("%s: port_lists=%d filtered=%d page=%d sort_value=%q sort_field_order=%q",
		label, len(r.PortList), r.PortListCount.Filtered, r.PortListCount.Page, r.Sort.Value, r.Sort.Field.Order)
	if len(r.Filters.Keywords.Keyword) > 0 {
		fmt.Printf(" keywords=[")
		for i, kw := range r.Filters.Keywords.Keyword {
			if i > 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("%s%s%s", kw.Column, kw.Relation, kw.Value)
		}
		fmt.Printf("]")
	}
	fmt.Println()
	for _, p := range r.PortList {
		fmt.Printf("  port_list id=%q name=%q writable=%t in_use=%t predefined=%t deprecated=%t",
			p.ID, p.Name, bool(p.Writable), bool(p.InUse), bool(p.Predefined), p.Deprecated != nil && bool(*p.Deprecated))
		if p.PortRanges != nil && len(p.PortRanges.PortRange) > 0 {
			fmt.Printf(" port_ranges=[")
			for i, pr := range p.PortRanges.PortRange {
				if i > 0 {
					fmt.Printf(" ")
				}
				fmt.Printf("%s:%d-%d", pr.Type, pr.Start, pr.End)
			}
			fmt.Printf("]")
		}
		if p.Targets != nil && len(p.Targets.Target) > 0 {
			fmt.Printf(" targets=[")
			for i, tgt := range p.Targets.Target {
				if i > 0 {
					fmt.Printf(" ")
				}
				fmt.Printf("%s", tgt.ID)
			}
			fmt.Printf("]")
		}
		if p.UserTags != nil && len(p.UserTags.Tag) > 0 {
			fmt.Printf(" user_tags=[")
			for i, tag := range p.UserTags.Tag {
				if i > 0 {
					fmt.Printf(" ")
				}
				fmt.Printf("%s:%s", tag.Name, tag.Value)
			}
			fmt.Printf("]")
		}
		fmt.Println()
	}
}

func decodeAndRoundTrip(path, decodeLabel, roundtripLabel string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var decoded genproto.GetPortListsResponse
	if err := xml.Unmarshal(data, &decoded); err != nil {
		fmt.Fprintln(os.Stderr, "decode:", err)
		os.Exit(1)
	}
	report(decodeLabel, &decoded)

	reencoded, err := xml.Marshal(&decoded)
	if err != nil {
		fmt.Fprintln(os.Stderr, "encode:", err)
		os.Exit(1)
	}
	var roundtripped genproto.GetPortListsResponse
	if err := xml.Unmarshal(reencoded, &roundtripped); err != nil {
		fmt.Fprintln(os.Stderr, "re-decode:", err)
		os.Exit(1)
	}
	report(roundtripLabel, &roundtripped)
}

func main() {
	decodeAndRoundTrip(os.Args[1], "upstream-decode", "upstream-roundtrip")
	decodeAndRoundTrip(os.Args[2], "live-decode", "live-roundtrip")
	decodeAndRoundTrip(os.Args[3], "synthetic-decode", "synthetic-roundtrip")
}
`

// assetScratchModule mirrors portListScratchModule for the asset family,
// which is still review-only: the goldens it compiles come from
// testdata/golden, not commands/.
func assetScratchModule(t *testing.T, dir string) string {
	t.Helper()
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("Skipping: no go toolchain on PATH (%v)", err)
	}
	got := generateOrFail(t)
	for _, name := range []string{CommandFile("create_asset"), CommandFile("modify_asset"), CommandFile("get_assets")} {
		if _, ok := got[name]; !ok {
			t.Fatalf("%s was not generated, so this test would not compile it", name)
		}
	}

	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to resolve the module root: %v", err)
	}
	pkgDir := filepath.Join(dir, "genproto")
	if err := os.MkdirAll(pkgDir, 0o750); err != nil {
		t.Fatalf("Failed to create the scratch package: %v", err)
	}
	goVersion, err := goDirective(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatalf("Failed to read the module's go directive: %v", err)
	}
	gomod := "module gogmpcodegencheck\n\ngo " + goVersion + "\n\n" +
		"require github.com/brennoo/go-gmp v0.0.0\n\n" +
		"replace github.com/brennoo/go-gmp => " + root + "\n"
	mustWriteFile(t, filepath.Join(dir, "go.mod"), gomod)
	mustWriteFile(t, filepath.Join(pkgDir, "support.go"),
		strings.Replace(supportShim(), "package "+targetPackage+"\n", "package genproto\n", 1))
	for name, src := range got {
		rewritten := strings.Replace(src, "\npackage "+targetPackage+"\n", "\npackage genproto\n", 1)
		if rewritten == src {
			t.Fatalf("%s: expected a %q package clause to rewrite", name, targetPackage)
		}
		mustWriteFile(t, filepath.Join(pkgDir, name), rewritten)
	}
	return goBin
}

// TestAssetGeneratedRequestRoundTrips proves the generated request shape for
// create_asset (both the asset and report choice branches), get_assets
// (IgnorePagination/Details unset/false/true, plus AssetID and FiltID
// together — the schema declares them as independent optional attributes),
// and modify_asset (comment set and comment omitted, the two wire shapes the
// modify_asset CommandWarning documents).
func TestAssetGeneratedRequestRoundTrips(t *testing.T) {
	dir := t.TempDir()
	goBin := assetScratchModule(t, dir)
	mustWriteFile(t, filepath.Join(dir, "main.go"), assetRequestRoundTripProgram)

	cmd := exec.CommandContext(context.Background(), goBin, "run", ".")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Round-trip program failed: %v\n%s", err, out)
	}

	want := "" +
		"create-asset-host: <create_asset><asset><name>203.0.113.42</name><type>host</type><comment>a comment</comment></asset></create_asset>\n" +
		"create-asset-report: <create_asset><report id=\"report-1\"><filter><term>apply_overrides=0</term></filter></report></create_asset>\n" +
		"get-assets-unset: <get_assets filter=\"name=foo\" type=\"host\"></get_assets>\n" +
		"get-assets-false: <get_assets ignore_pagination=\"0\" details=\"0\"></get_assets>\n" +
		"get-assets-true: <get_assets ignore_pagination=\"1\" details=\"1\"></get_assets>\n" +
		"get-assets-id-and-filtid: <get_assets asset_id=\"asset-1\" filt_id=\"filter-1\"></get_assets>\n" +
		"modify-asset-comment-set: <modify_asset asset_id=\"asset-1\"><comment>updated comment</comment></modify_asset>\n" +
		"modify-asset-comment-omitted: <modify_asset asset_id=\"asset-1\"></modify_asset>\n"
	if got := string(out); got != want {
		t.Fatalf("Round-trip program printed:\n%s\nwant:\n%s", got, want)
	}
}

// assetRequestRoundTripProgram is a `go run` main for
// TestAssetGeneratedRequestRoundTrips.
const assetRequestRoundTripProgram = `package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"gogmpcodegencheck/genproto"
)

func printXML(label string, v any) {
	b, err := xml.Marshal(v)
	if err != nil {
		fmt.Fprintln(os.Stderr, label, err)
		os.Exit(1)
	}
	fmt.Printf("%s: %s\n", label, b)
}

func boolp(b genproto.GMPBool) *genproto.GMPBool { return &b }

func main() {
	printXML("create-asset-host", &genproto.CreateAsset{
		Asset: &genproto.Asset{Name: "203.0.113.42", Type: "host", Comment: "a comment"},
	})
	printXML("create-asset-report", &genproto.CreateAsset{
		Report: &genproto.AssetReport{ID: "report-1", Filter: genproto.AssetFilter{Term: "apply_overrides=0"}},
	})
	printXML("get-assets-unset", &genproto.GetAssets{
		Type: "host", Filter: "name=foo",
	})
	printXML("get-assets-false", &genproto.GetAssets{
		IgnorePagination: boolp(false), Details: boolp(false),
	})
	printXML("get-assets-true", &genproto.GetAssets{
		IgnorePagination: boolp(true), Details: boolp(true),
	})
	printXML("get-assets-id-and-filtid", &genproto.GetAssets{
		AssetID: "asset-1", FiltID: "filter-1",
	})
	printXML("modify-asset-comment-set", &genproto.ModifyAsset{
		AssetID: "asset-1", Comment: "updated comment",
	})
	printXML("modify-asset-comment-omitted", &genproto.ModifyAsset{
		AssetID: "asset-1",
	})
}
`

// TestAssetGeneratedResponseRoundTrips decodes the committed live fixtures
// plus a synthetic one; see docs/review-2026-08.md's `asset added in
// review-only mode` section for what each fixture covers.
func TestAssetGeneratedResponseRoundTrips(t *testing.T) {
	dir := t.TempDir()
	goBin := assetScratchModule(t, dir)

	liveFixture1, err := filepath.Abs("../gmptest/testdata/live/get_assets/001_response.xml")
	if err != nil {
		t.Fatalf("Failed to resolve the live fixture path: %v", err)
	}
	if _, err := os.Stat(liveFixture1); err != nil {
		t.Fatalf("Live fixture not found: %v", err)
	}
	liveFixture2, err := filepath.Abs("../gmptest/testdata/live/get_assets/002_response.xml")
	if err != nil {
		t.Fatalf("Failed to resolve the live fixture path: %v", err)
	}
	if _, err := os.Stat(liveFixture2); err != nil {
		t.Fatalf("Live fixture not found: %v", err)
	}
	syntheticFixture := filepath.Join(dir, "synthetic.xml")
	mustWriteFile(t, syntheticFixture, assetSyntheticResponse)

	mustWriteFile(t, filepath.Join(dir, "main.go"), assetResponseRoundTripProgram)

	cmd := exec.CommandContext(context.Background(), goBin, "run", ".", liveFixture1, liveFixture2, syntheticFixture)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Round-trip program failed: %v\n%s", err, out)
	}

	want := "" +
		"live1-decode: assets=0 filtered=0 page=0\n" +
		"live1-roundtrip: assets=0 filtered=0 page=0\n" +
		"live2-decode: assets=0 filtered=0 page=0\n" +
		"live2-roundtrip: assets=0 filtered=0 page=0\n" +
		"synthetic-decode: assets=1 filtered=1 page=1\n" +
		"  asset id=\"asset-syn-1\" name=\"203.0.113.42\" writable=true in_use=false type=\"host\"" +
		" identifier_name=\"ip\" identifier_source_name=true" +
		" host_severity=5 routes=[203.0.113.1(distance=1)] detail=[MAC:00:11:22:33:44:55]" +
		" os_hosts=[asset-syn-2(sev=5.0)]\n" +
		"synthetic-roundtrip: assets=1 filtered=1 page=1\n" +
		"  asset id=\"asset-syn-1\" name=\"203.0.113.42\" writable=true in_use=false type=\"host\"" +
		" identifier_name=\"ip\" identifier_source_name=true" +
		" host_severity=5 routes=[203.0.113.1(distance=1)] detail=[MAC:00:11:22:33:44:55]" +
		" os_hosts=[asset-syn-2(sev=5.0)]\n"
	if got := string(out); got != want {
		t.Fatalf("Round-trip program printed:\n%s\nwant:\n%s", got, want)
	}
}

// assetSyntheticResponse pins schema-derived fields absent from the
// asset-less live fixtures, including the accepted source/name loss.
const assetSyntheticResponse = `<get_assets_response status="200" status_text="OK">
  <asset id="asset-syn-1">
    <owner><name>admin</name></owner>
    <name>203.0.113.42</name>
    <comment></comment>
    <creation_time>2026-01-01T00:00:00Z</creation_time>
    <modification_time>2026-01-01T00:00:00Z</modification_time>
    <writable>1</writable>
    <in_use>0</in_use>
    <permissions><permission><name>Everything</name></permission></permissions>
    <type>host</type>
    <user_tags>
      <count>1</count>
      <tag id="tag-1"><name>env</name><value>prod</value><comment></comment></tag>
    </user_tags>
    <identifiers>
      <identifier id="ident-1">
        <name>ip</name>
        <value>203.0.113.42</value>
        <creation_time>2026-01-01T00:00:00Z</creation_time>
        <modification_time>2026-01-01T00:00:00Z</modification_time>
        <source id="source-1"><type>User</type><data></data><deleted>0</deleted><name>admin</name></source>
      </identifier>
    </identifiers>
    <host>
      <severity><value>5</value></severity>
      <routes>
        <route>
          <host id="host-1" distance="1" same_source="1"><ip>203.0.113.1</ip></host>
        </route>
      </routes>
      <detail>
        <name>MAC</name>
        <value>00:11:22:33:44:55</value>
        <source><id>source-2</id><type>Report</type></source>
      </detail>
    </host>
    <os>
      <title>Linux Kernel</title>
      <installs>1</installs>
      <all_installs>1</all_installs>
      <latest_severity><value>5.0</value></latest_severity>
      <highest_severity><value>5.0</value></highest_severity>
      <average_severity><value>5.0</value></average_severity>
      <hosts>
        <asset id="asset-syn-2"><name>203.0.113.43</name><severity><value>5.0</value></severity></asset>
      </hosts>
    </os>
  </asset>
  <filters id=""><term>first=1 rows=10 sort=name</term><keywords/></filters>
  <sort><field>name<order>ascending</order></field></sort>
  <assets start="1" max="1000"/>
  <asset_count>1<filtered>1</filtered><page>1</page></asset_count>
</get_assets_response>
`

// assetResponseRoundTripProgram is a `go run` main for
// TestAssetGeneratedResponseRoundTrips. argv[1] and argv[2] are the two
// committed live (asset-less) fixtures, argv[3] the synthetic populated one.
const assetResponseRoundTripProgram = `package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"gogmpcodegencheck/genproto"
)

func report(label string, r *genproto.GetAssetsResponse) {
	fmt.Printf("%s: assets=%d filtered=%d page=%d\n",
		label, len(r.Asset), r.AssetCount.Filtered, r.AssetCount.Page)
	for _, a := range r.Asset {
		fmt.Printf("  asset id=%q name=%q writable=%t in_use=%t type=%q",
			a.ID, a.Name, bool(a.Writable), bool(a.InUse), a.Type)
		if a.Identifiers != nil {
			fmt.Printf(" identifier_name=%q identifier_source_name=%t",
				a.Identifiers.Identifier.Name, bool(a.Identifiers.Identifier.Source.Name))
		}
		if a.Host != nil {
			fmt.Printf(" host_severity=%v", a.Host.Severity.Value)
			if a.Host.Routes != nil {
				fmt.Printf(" routes=[")
				for i, rt := range a.Host.Routes.Route {
					for j, h := range rt.Host {
						if i+j > 0 {
							fmt.Printf(" ")
						}
						distance := 0
						if h.Distance != nil {
							distance = *h.Distance
						}
						fmt.Printf("%s(distance=%d)", h.IP, distance)
					}
				}
				fmt.Printf("]")
			}
			if len(a.Host.Detail) > 0 {
				fmt.Printf(" detail=[")
				for i, d := range a.Host.Detail {
					if i > 0 {
						fmt.Printf(" ")
					}
					fmt.Printf("%s:%s", d.Name, d.Value)
				}
				fmt.Printf("]")
			}
		}
		if a.OS != nil {
			fmt.Printf(" os_hosts=[")
			for i, h := range a.OS.Hosts.Asset {
				if i > 0 {
					fmt.Printf(" ")
				}
				fmt.Printf("%s(sev=%s)", h.ID, h.Severity.Value)
			}
			fmt.Printf("]")
		}
		fmt.Println()
	}
}

func decodeAndRoundTrip(path, decodeLabel, roundtripLabel string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var decoded genproto.GetAssetsResponse
	if err := xml.Unmarshal(data, &decoded); err != nil {
		fmt.Fprintln(os.Stderr, "decode:", err)
		os.Exit(1)
	}
	report(decodeLabel, &decoded)

	reencoded, err := xml.Marshal(&decoded)
	if err != nil {
		fmt.Fprintln(os.Stderr, "encode:", err)
		os.Exit(1)
	}
	var roundtripped genproto.GetAssetsResponse
	if err := xml.Unmarshal(reencoded, &roundtripped); err != nil {
		fmt.Fprintln(os.Stderr, "re-decode:", err)
		os.Exit(1)
	}
	report(roundtripLabel, &roundtripped)
}

func main() {
	decodeAndRoundTrip(os.Args[1], "live1-decode", "live1-roundtrip")
	decodeAndRoundTrip(os.Args[2], "live2-decode", "live2-roundtrip")
	decodeAndRoundTrip(os.Args[3], "synthetic-decode", "synthetic-roundtrip")
}
`

// mustWriteFile is writeFile (scratch.go) with the *testing.T ceremony a test
// wants but a reusable helper shouldn't carry.
func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := writeFile(path, content); err != nil {
		t.Fatalf("Failed to write %s: %v", path, err)
	}
}

// generateOrFail renders every family the manifest declares, review-only ones
// included: the goldens are where a review-only family's output lives, so the
// golden test has to see it.
func generateOrFail(t *testing.T) map[string]string {
	t.Helper()
	return generateManifest(t, manifestOrFail(t))
}

// generateAdoptedOrFail renders only the families whose output is committed
// under commands/ — what the CLI writes and what CheckOutput inspects.
func generateAdoptedOrFail(t *testing.T) map[string]string {
	t.Helper()
	return generateManifest(t, manifestOrFail(t).Adopted())
}

func generateManifest(t *testing.T, m Manifest) map[string]string {
	t.Helper()
	cmds, err := LoadSnapshots(snapshotDir, m.AllCommands())
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}
	out, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist)
	if err != nil {
		t.Fatalf("GenerateAll failed: %v", err)
	}
	return out
}

// manifestOrFail loads the committed manifest, which every test in this
// package sources its command list from so none of them carries a second
// hardcoded copy.
func manifestOrFail(t *testing.T) Manifest {
	t.Helper()
	m, err := LoadManifest(manifestFile)
	if err != nil {
		t.Fatalf("Failed to load the manifest: %v", err)
	}
	return m
}

// --- fieldOverridesFor cross-use consistency (bug #2 regression tests) -----

// twinFieldOverrideCommands builds two commands, cmd_a and cmd_b, whose
// responses each declare a structurally identical nested <thing> shape (same
// name, same <leaf> child), so Generate's structural dedup collapses them into
// one shared structType reached by two different (command, path) uses — the
// exact situation fieldOverridesFor's cross-use consistency check exists for.
// Built directly rather than loaded from a snapshot, same spirit as
// TestGenerateAllIsPerFamily's synthetic twin family.
func twinFieldOverrideCommands() (cmdA, cmdB schema.Command) {
	shape := func(rootName string) *schema.Node {
		return &schema.Node{
			Name: rootName,
			Children: []*schema.Node{
				{
					Name: "thing", Kind: schema.KindElement, Card: schema.CardRequired,
					Children: []*schema.Node{
						{Name: "leaf", Kind: schema.KindElement, Card: schema.CardRequired, Type: "text"},
					},
				},
			},
		}
	}
	cmdA = schema.Command{Name: "cmd_a", Request: &schema.Node{Name: "cmd_a"}, Response: shape("cmd_a_response"), Provenance: testProvenance()}
	cmdB = schema.Command{Name: "cmd_b", Request: &schema.Node{Name: "cmd_b"}, Response: shape("cmd_b_response"), Provenance: testProvenance()}
	return cmdA, cmdB
}

func TestFieldOverridesCrossUseMismatchOnMissingKey(t *testing.T) {
	cmdA, cmdB := twinFieldOverrideCommands()
	fo := map[string][]FieldOverride{
		"cmd_a:/thing": {{Name: "extra", GoType: "string", Reason: "t"}},
	}
	_, err := Generate([]schema.Command{cmdA, cmdB}, TypeOverrides, PointerOverrides, fo, nil, nil, nil)
	if err == nil {
		t.Fatal("Generate succeeded; want an error: cmd_b's use has no matching override while cmd_a's does")
	}
	if !strings.Contains(err.Error(), "cmd_a:/thing") || !strings.Contains(err.Error(), "cmd_b:/thing") {
		t.Fatalf("Generate error = %q, want it to name both cmd_a:/thing and cmd_b:/thing", err.Error())
	}
}

func TestFieldOverridesCrossUseMismatchOnDifferentEntries(t *testing.T) {
	cmdA, cmdB := twinFieldOverrideCommands()
	fo := map[string][]FieldOverride{
		"cmd_a:/thing": {{Name: "extra", GoType: "string", Reason: "t"}},
		"cmd_b:/thing": {{Name: "other", GoType: "string", Reason: "t"}},
	}
	_, err := Generate([]schema.Command{cmdA, cmdB}, TypeOverrides, PointerOverrides, fo, nil, nil, nil)
	if err == nil {
		t.Fatal("Generate succeeded; want an error: cmd_a and cmd_b disagree on the override for the shared <thing> shape")
	}
	if !strings.Contains(err.Error(), "cmd_a:/thing") || !strings.Contains(err.Error(), "cmd_b:/thing") {
		t.Fatalf("Generate error = %q, want it to name both cmd_a:/thing and cmd_b:/thing", err.Error())
	}
}

func TestFieldOverridesCrossUseIdenticalEntriesSucceed(t *testing.T) {
	cmdA, cmdB := twinFieldOverrideCommands()
	fo := map[string][]FieldOverride{
		"cmd_a:/thing": {{Name: "extra", GoType: "string", Reason: "t"}},
		"cmd_b:/thing": {{Name: "extra", GoType: "string", Reason: "t"}},
	}
	out, err := Generate([]schema.Command{cmdA, cmdB}, TypeOverrides, PointerOverrides, fo, nil, nil, nil)
	if err != nil {
		t.Fatalf("Generate failed with identical overrides on both uses: %v", err)
	}
	var all strings.Builder
	for _, src := range out {
		all.WriteString(src)
	}
	if got := strings.Count(all.String(), `gmpoverride:"t"`); got != 1 {
		t.Fatalf(`gmpoverride:"t" appears %d times across the generated output, want exactly 1 `+
			`(the structurally-shared struct must be emitted once, not once per use)`, got)
	}
}

func TestFieldOverridesSchemaChildCollision(t *testing.T) {
	cmdA, _ := twinFieldOverrideCommands()
	fo := map[string][]FieldOverride{
		"cmd_a:/thing": {{Name: "leaf", GoType: "string", Reason: "t"}},
	}
	_, err := Generate([]schema.Command{cmdA}, TypeOverrides, PointerOverrides, fo, nil, nil, nil)
	if err == nil {
		t.Fatal("Generate succeeded; want an error: the override's Name duplicates <thing>'s schema-declared <leaf> child")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("Generate error = %q, want it to mention the duplicate", err.Error())
	}
}

// firstDiff points at the first differing line, so a failure names the change
// instead of dumping thirteen kilobytes of struct definitions.
func firstDiff(want, got string) string {
	w, g := strings.Split(want, "\n"), strings.Split(got, "\n")
	for i := 0; i < len(w) && i < len(g); i++ {
		if w[i] != g[i] {
			return "first difference at line " + strconv.Itoa(i+1) + ":\n  committed: " + w[i] + "\n  now:       " + g[i]
		}
	}
	return "one file is a prefix of the other"
}
