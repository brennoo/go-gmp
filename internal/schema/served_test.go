package schema

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"
	"testing"
	"time"
)

const servedTestdata = "testdata/served/help_response.xml"

// fakeSourceRevision is a validly-shaped stand-in for a resolved image
// digest, for tests that don't care about its exact value.
var fakeSourceRevision = "example.test/gvmd@sha256:" + strings.Repeat("ab", 32)

// TestLoadServedIgnoresTheDecoy proves the unwrap step finds the real,
// top-level <protocol> document rather than the truncated example nested
// inside the help command's own <example> block — the trap found in the
// T-16b spike (see docs/review-2026-08.md, T-16b). A naive implementation
// that searched for the first <protocol> tag anywhere in the byte stream
// would still pass here by coincidence; this asserts the actual shape
// extracted, which a coincidence cannot fake.
func TestLoadServedIgnoresTheDecoy(t *testing.T) {
	data, err := os.ReadFile(servedTestdata)
	if err != nil {
		t.Fatalf("reading testdata: %v", err)
	}

	cmds, err := LoadServed(data, fakeSourceRevision)
	if err != nil {
		t.Fatalf("LoadServed: %v", err)
	}
	if len(cmds) != 2 {
		t.Fatalf("got %d commands, want 2 (the decoy's stub get_fixture must not be counted)", len(cmds))
	}

	byName := map[string]Command{}
	for _, c := range cmds {
		byName[c.Name] = c
	}

	fixture, ok := byName["get_fixture"]
	if !ok {
		t.Fatal("get_fixture is missing")
	}
	if len(fixture.Defects) != 0 {
		t.Errorf("get_fixture has defects %+v, want none", fixture.Defects)
	}
	if len(fixture.Request.Children) != 1 || fixture.Request.Children[0].Name != "id" {
		t.Errorf("get_fixture request children = %+v, want a single %q child", fixture.Request.Children, "id")
	}
	if len(fixture.Response.Children) != 1 || fixture.Response.Children[0].Name != "id" {
		t.Errorf("get_fixture response children = %+v, want a single %q child", fixture.Response.Children, "id")
	}

	if _, ok := byName["help"]; !ok {
		t.Fatal("help is missing")
	}
}

// TestLoadServedProvenance checks every Provenance field LoadServed stamps,
// independently of LoadServed's own computation.
func TestLoadServedProvenance(t *testing.T) {
	data, err := os.ReadFile(servedTestdata)
	if err != nil {
		t.Fatalf("reading testdata: %v", err)
	}

	cmds, err := LoadServed(data, fakeSourceRevision)
	if err != nil {
		t.Fatalf("LoadServed: %v", err)
	}

	sum := sha256.Sum256(data)
	wantHash := hex.EncodeToString(sum[:])

	first := cmds[0].Provenance
	for _, c := range cmds[1:] {
		if c.Provenance != first {
			t.Fatalf("%s: Provenance = %+v, want %+v (same as %s)", c.Name, c.Provenance, first, cmds[0].Name)
		}
	}

	if first.Method != ExtractionServed {
		t.Errorf("Method = %q, want %q", first.Method, ExtractionServed)
	}
	if first.Source != servedSource {
		t.Errorf("Source = %q, want %q", first.Source, servedSource)
	}
	if first.SourceRevision != fakeSourceRevision {
		t.Errorf("SourceRevision = %q, want %q", first.SourceRevision, fakeSourceRevision)
	}
	if first.DeclaredVersion != "22.7" {
		t.Errorf("DeclaredVersion = %q, want the real protocol's own <version>, not the decoy's %q", first.DeclaredVersion, "trunc")
	}
	if first.SourceSHA256 != wantHash {
		t.Errorf("SourceSHA256 = %q, want %q (independently computed over the full envelope)", first.SourceSHA256, wantHash)
	}
	if first.ParsedSHA256 != first.SourceSHA256 {
		t.Errorf("ParsedSHA256 = %q, want it to equal SourceSHA256: there is no split step for served extraction", first.ParsedSHA256)
	}
	if first.ExtractedAt.Location() != time.UTC {
		t.Errorf("ExtractedAt location = %v, want UTC", first.ExtractedAt.Location())
	}
	if since := time.Since(first.ExtractedAt); since < 0 || since > time.Minute {
		t.Errorf("ExtractedAt = %v, want a timestamp from this test run (age %v)", first.ExtractedAt, since)
	}
}

func TestLoadServedRejectsWrongRoot(t *testing.T) {
	_, err := LoadServed([]byte(`<get_version_response status="200"/>`), fakeSourceRevision)
	if err == nil || !strings.Contains(err.Error(), "want <help_response>") {
		t.Fatalf("LoadServed err = %v, want an unexpected-root-element error", err)
	}
}

func TestLoadServedRejectsNonOKStatus(t *testing.T) {
	_, err := LoadServed([]byte(`<help_response status="400" status_text="Bad Request"/>`), fakeSourceRevision)
	if err == nil || !strings.Contains(err.Error(), "status 400") {
		t.Fatalf("LoadServed err = %v, want a status error", err)
	}
}

func TestLoadServedRejectsMissingProtocol(t *testing.T) {
	_, err := LoadServed([]byte(`<help_response status="200"><schema format="XML"></schema></help_response>`), fakeSourceRevision)
	if err == nil {
		t.Fatal("LoadServed err = nil, want an error: the envelope has no <protocol> element")
	}
}

// TestLoadServedRejectsInvalidSourceRevision checks that every shape short of
// a resolved, repository-qualified digest is rejected before data is even
// parsed — the fail-closed claim in LoadServed's doc comment must be enforced
// by the function itself, not by caller discipline.
func TestLoadServedRejectsInvalidSourceRevision(t *testing.T) {
	data, err := os.ReadFile(servedTestdata)
	if err != nil {
		t.Fatalf("reading testdata: %v", err)
	}

	tests := []struct {
		name string
		rev  string
	}{
		{"empty", ""},
		{"tag only", "example.test/gvmd:stable"},
		{"short digest", "example.test/gvmd@sha256:" + strings.Repeat("a", 10)},
		{"malformed hex", "example.test/gvmd@sha256:" + strings.Repeat("g", 64)},
		{"mixed-case hex", "example.test/gvmd@sha256:" + strings.Repeat("aB", 32)},
		{"uppercase hex", "example.test/gvmd@sha256:" + strings.Repeat("AB", 32)},
		{"leading colon (empty repository, looks like a bare port)", ":stable@sha256:" + strings.Repeat("ab", 32)},
		{"embedded tag before the digest", "example.com/gvmd:stable@sha256:" + strings.Repeat("ab", 32)},
		{"leading slash (empty path component)", "/gvmd@sha256:" + strings.Repeat("ab", 32)},
		{"doubled slash (empty path component)", "example.com//gvmd@sha256:" + strings.Repeat("ab", 32)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := LoadServed(data, tt.rev); err == nil {
				t.Fatalf("LoadServed(_, %q) err = nil, want an error", tt.rev)
			}
		})
	}
}
