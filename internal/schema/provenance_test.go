package schema

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// TestLoadStampsUTCProvenance proves ExtractedAt is a real, current UTC
// timestamp, not a zero value or one from some other location.
func TestLoadStampsUTCProvenance(t *testing.T) {
	cmds := loadSchemaOrSkip(t)

	got := cmds[0].Provenance.ExtractedAt
	if got.IsZero() {
		t.Fatal("ExtractedAt is zero")
	}
	if got.Location() != time.UTC {
		t.Fatalf("ExtractedAt location = %v, want UTC", got.Location())
	}
	if since := time.Since(got); since < 0 || since > time.Minute {
		t.Fatalf("ExtractedAt = %v, want a timestamp from this test run (age %v)", got, since)
	}
}

// TestLoadSharesOneProvenanceAcrossBatch proves every command from one Load
// call carries byte-identical Provenance; split command files have no
// version marker of their own to disagree with.
func TestLoadSharesOneProvenanceAcrossBatch(t *testing.T) {
	cmds := loadSchemaOrSkip(t)

	first := cmds[0].Provenance
	for _, c := range cmds[1:] {
		if c.Provenance != first {
			t.Fatalf("%s: Provenance = %+v, want %+v (same as %s)", c.Name, c.Provenance, first, cmds[0].Name)
		}
	}
}

// rootVersionPattern finds the first <version> element, independent of
// extractionProvenance's own xml.Unmarshal-based approach.
var rootVersionPattern = regexp.MustCompile(`<version>([^<]*)</version>`)

// TestProvenanceMatchesSourceDocument proves SourceSHA256 and DeclaredVersion
// match vendor/XML/GMP.xml.in, computed here independently of extractionProvenance.
func TestProvenanceMatchesSourceDocument(t *testing.T) {
	cmds := loadSchemaOrSkip(t)

	data, err := os.ReadFile("../../vendor/XML/GMP.xml.in")
	if err != nil {
		t.Fatalf("reading the source document independently: %v", err)
	}

	sum := sha256.Sum256(data)
	wantHash := hex.EncodeToString(sum[:])

	m := rootVersionPattern.FindSubmatch(data)
	if m == nil {
		t.Fatal("no <version> element found in the source document")
	}
	wantVersion := string(m[1])

	got := cmds[0].Provenance
	if got.SourceSHA256 != wantHash {
		t.Errorf("SourceSHA256 = %q, want %q (independently computed)", got.SourceSHA256, wantHash)
	}
	if got.DeclaredVersion != wantVersion {
		t.Errorf("DeclaredVersion = %q, want %q (independently parsed)", got.DeclaredVersion, wantVersion)
	}
	if got.Method != ExtractionFile {
		t.Errorf("Method = %q, want %q", got.Method, ExtractionFile)
	}
	if got.Source != upstreamSource {
		t.Errorf("Source = %q, want %q", got.Source, upstreamSource)
	}
}

// extractionTestdata is an original, minimal schema tree committed so that
// extraction itself is exercised on a clean checkout. Every other test here
// skips without vendor/XML, which would otherwise leave the code that
// *produces* provenance with no CI coverage at all — only the codegen-side
// code that consumes it.
const extractionTestdata = "testdata/extraction/commands"

func TestExtractionProvenanceWithoutVendoredSchema(t *testing.T) {
	cmds, err := Load(extractionTestdata, "v0.0.0-testdata")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cmds) != 1 || cmds[0].Name != "get_fixture" {
		t.Fatalf("Load returned %d commands, want just get_fixture", len(cmds))
	}

	got := cmds[0].Provenance
	if got.DeclaredVersion != "0.0-testdata" {
		t.Errorf("DeclaredVersion = %q, want the protocol's own <version>, not a nested one", got.DeclaredVersion)
	}
	if got.SourceRevision != "v0.0.0-testdata" {
		t.Errorf("SourceRevision = %q, want the value passed to Load", got.SourceRevision)
	}
	if !sha256Hex.MatchString(got.SourceSHA256) {
		t.Errorf("SourceSHA256 = %q, want 64 lowercase hex characters", got.SourceSHA256)
	}
	if !sha256Hex.MatchString(got.ParsedSHA256) {
		t.Errorf("ParsedSHA256 = %q, want 64 lowercase hex characters", got.ParsedSHA256)
	}
	if got.SourceSHA256 == got.ParsedSHA256 {
		t.Error("SourceSHA256 and ParsedSHA256 are equal; they hash different inputs")
	}
}

var sha256Hex = regexp.MustCompile(`^[0-9a-f]{64}$`)

// TestParsedSHA256TracksTheFilesActuallyParsed is the property SourceSHA256
// cannot provide: nothing binds split output back to the document it came
// from, so an edited or renamed split has to move a digest of its own.
func TestParsedSHA256TracksTheFilesActuallyParsed(t *testing.T) {
	base, err := Load(extractionTestdata, "")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	dir := t.TempDir()
	copyInto := func(name string, edit func([]byte) []byte) string {
		t.Helper()
		commands := filepath.Join(dir, name, "commands")
		if err := os.MkdirAll(commands, 0o750); err != nil {
			t.Fatal(err)
		}
		for _, f := range []struct{ from, to string }{
			{"testdata/extraction/GMP.xml.in", filepath.Join(dir, name, "GMP.xml.in")},
			{filepath.Join(extractionTestdata, "command_get_fixture.xml"), filepath.Join(commands, "command_get_fixture.xml")},
		} {
			data, err := os.ReadFile(f.from) //nolint:gosec // fixed testdata paths
			if err != nil {
				t.Fatal(err)
			}
			if edit != nil && strings.HasPrefix(filepath.Base(f.to), "command_") {
				data = edit(data)
			}
			if err := os.WriteFile(f.to, data, 0o600); err != nil {
				t.Fatal(err)
			}
		}
		return commands
	}

	unchanged, err := Load(copyInto("same", nil), "")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if unchanged[0].Provenance.ParsedSHA256 != base[0].Provenance.ParsedSHA256 {
		t.Error("ParsedSHA256 changed for identical files; it must depend only on their names and bytes")
	}

	edit := func(b []byte) []byte {
		return []byte(strings.Replace(string(b), "<name>version</name>", "<name>version </name>", 1))
	}
	edited, err := Load(copyInto("edited", edit), "")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if edited[0].Provenance.SourceSHA256 != base[0].Provenance.SourceSHA256 {
		t.Fatal("the test edited a command file but SourceSHA256 moved; it should hash only GMP.xml.in")
	}
	if edited[0].Provenance.ParsedSHA256 == base[0].Provenance.ParsedSHA256 {
		t.Error("an edited command file left ParsedSHA256 unchanged; the split is not attested")
	}
}
