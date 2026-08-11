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

	if _, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, FieldOverrides); err == nil {
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
// GetTLSCertificatesResponse/ModifyTLSCertificate agree with the two
// FieldOverride entries in overrides.go (T-9 finding 3), and pins down
// ModifyTLSCertificate's omit/explicit-empty/set wire shape (finding 4) as a
// generator-level stand-in for internal/gmptest/cases_tls_certificate.go,
// which cannot exist until tls_certificate is adopted.
func TestTLSCertificateGeneratedRoundTrips(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("Skipping: no go toolchain on PATH (%v)", err)
	}
	got := generateOrFail(t)
	if _, ok := got[CommandFile("get_tls_certificates")]; !ok {
		t.Fatal("get_tls_certificates was not generated, so this test would not compile it")
	}
	if _, ok := got[CommandFile("modify_tls_certificate")]; !ok {
		t.Fatal("modify_tls_certificate was not generated, so this test would not compile it")
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
		"live-roundtrip: Value=\"" + liveCert + "\" Format=\"PEM\" Valid=true\n" +
		"marshal-name-nil: <modify_tls_certificate tls_certificate_id=\"tid-1\"></modify_tls_certificate>\n" +
		"marshal-name-empty: <modify_tls_certificate tls_certificate_id=\"tid-1\"><name></name></modify_tls_certificate>\n" +
		"marshal-name-set: <modify_tls_certificate tls_certificate_id=\"tid-1\"><name>new-name</name></modify_tls_certificate>\n" +
		"marshal-comment-nil: <modify_tls_certificate tls_certificate_id=\"tid-1\"></modify_tls_certificate>\n" +
		"marshal-comment-empty: <modify_tls_certificate tls_certificate_id=\"tid-1\"><comment></comment></modify_tls_certificate>\n" +
		"marshal-comment-set: <modify_tls_certificate tls_certificate_id=\"tid-1\"><comment>new-comment</comment></modify_tls_certificate>\n" +
		"marshal-trust-nil: <modify_tls_certificate tls_certificate_id=\"tid-1\"></modify_tls_certificate>\n" +
		"marshal-trust-false: <modify_tls_certificate tls_certificate_id=\"tid-1\"><trust>0</trust></modify_tls_certificate>\n" +
		"marshal-trust-true: <modify_tls_certificate tls_certificate_id=\"tid-1\"><trust>1</trust></modify_tls_certificate>\n"
	if got := string(out); got != want {
		t.Fatalf("Round-trip program printed:\n%s\nwant:\n%s", got, want)
	}
}

// tlsCertificateRoundTripProgram is a `go run` main for
// TestTLSCertificateGeneratedRoundTrips. argv[1] is the upstream fixture,
// argv[2] the live one; it decodes and round-trips each through the real
// generated GetTLSCertificatesResponse, then exercises
// ModifyTLSCertificate's marshal shape across the omit/explicit-empty/set
// matrix on Name, Comment and Trust.
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

func marshalModify(label string, m genproto.ModifyTLSCertificate) {
	out, err := xml.Marshal(&m)
	if err != nil {
		fmt.Fprintln(os.Stderr, label, "marshal:", err)
		os.Exit(1)
	}
	fmt.Printf("%s: %s\n", label, out)
}

func main() {
	decodeAndRoundTrip(os.Args[1], "upstream-decode", "upstream-roundtrip")
	decodeAndRoundTrip(os.Args[2], "live-decode", "live-roundtrip")

	empty, name, comment := "", "new-name", "new-comment"
	trustFalse, trustTrue := genproto.GMPBool(false), genproto.GMPBool(true)

	// Stand-in for internal/gmptest/cases_tls_certificate.go: that file needs
	// commands.ModifyTLSCertificate, which does not exist until adoption.
	marshalModify("marshal-name-nil", genproto.ModifyTLSCertificate{TLSCertificateID: "tid-1"})
	marshalModify("marshal-name-empty", genproto.ModifyTLSCertificate{TLSCertificateID: "tid-1", Name: &empty})
	marshalModify("marshal-name-set", genproto.ModifyTLSCertificate{TLSCertificateID: "tid-1", Name: &name})
	marshalModify("marshal-comment-nil", genproto.ModifyTLSCertificate{TLSCertificateID: "tid-1"})
	marshalModify("marshal-comment-empty", genproto.ModifyTLSCertificate{TLSCertificateID: "tid-1", Comment: &empty})
	marshalModify("marshal-comment-set", genproto.ModifyTLSCertificate{TLSCertificateID: "tid-1", Comment: &comment})
	marshalModify("marshal-trust-nil", genproto.ModifyTLSCertificate{TLSCertificateID: "tid-1"})
	marshalModify("marshal-trust-false", genproto.ModifyTLSCertificate{TLSCertificateID: "tid-1", Trust: &trustFalse})
	marshalModify("marshal-trust-true", genproto.ModifyTLSCertificate{TLSCertificateID: "tid-1", Trust: &trustTrue})
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
	out, err := GenerateAll(m, cmds, TypeOverrides, PointerOverrides, FieldOverrides)
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
	cmdA = schema.Command{Name: "cmd_a", Request: &schema.Node{Name: "cmd_a"}, Response: shape("cmd_a_response")}
	cmdB = schema.Command{Name: "cmd_b", Request: &schema.Node{Name: "cmd_b"}, Response: shape("cmd_b_response")}
	return cmdA, cmdB
}

func TestFieldOverridesCrossUseMismatchOnMissingKey(t *testing.T) {
	cmdA, cmdB := twinFieldOverrideCommands()
	fo := map[string][]FieldOverride{
		"cmd_a:/thing": {{Name: "extra", GoType: "string", Reason: "t"}},
	}
	_, err := Generate([]schema.Command{cmdA, cmdB}, TypeOverrides, PointerOverrides, fo)
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
	_, err := Generate([]schema.Command{cmdA, cmdB}, TypeOverrides, PointerOverrides, fo)
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
	out, err := Generate([]schema.Command{cmdA, cmdB}, TypeOverrides, PointerOverrides, fo)
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
	_, err := Generate([]schema.Command{cmdA}, TypeOverrides, PointerOverrides, fo)
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
