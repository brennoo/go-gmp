package gmptest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestUpdateRefusesToRewriteEvidence pins the rule that keeps the corpus
// trustworthy. A golden test that can rewrite its own input goes green by
// destroying the observation that would have failed it, and nothing about the
// resulting diff looks wrong.
func TestUpdateRefusesToRewriteEvidence(t *testing.T) {
	refused := []string{
		"testdata/upstream/get_tasks/001_response.xml",
		"testdata/upstream/get_tasks/001_request.xml",
		"testdata/live/get_tasks/001_response.xml",
		"testdata/upstream/get_tasks/001_meta.json",
		"testdata/upstream/MANIFEST.md",
	}
	for _, p := range refused {
		if err := allowedUpdatePath(filepath.FromSlash(p)); err == nil {
			t.Errorf("Expected -update to refuse %s", p)
		}
	}

	// Outside the corpus, regardless of how the file is named.
	outside := []string{
		filepath.Join(t.TempDir(), "arbitrary_fidelity.txt"),
		filepath.Join(t.TempDir(), "requests", "x.xml"),
		"../../commands/get_tasks_fidelity.txt",
	}
	for _, p := range outside {
		if err := allowedUpdatePath(p); err == nil {
			t.Errorf("Expected -update to refuse %s (outside %s)", p, TestdataDir)
		}
	}

	allowed := []string{
		"testdata/upstream/get_tasks/001_fidelity.txt",
		"testdata/live/get_tasks/001_fidelity.txt",
		"testdata/FIDELITY-SUMMARY.txt",
		"testdata/requests/get_tasks/with_filter.xml",
	}
	for _, p := range allowed {
		if err := allowedUpdatePath(filepath.FromSlash(p)); err != nil {
			t.Errorf("Expected -update to allow %s, got: %v", p, err)
		}
	}
}

// The guard must hold at the write call too, not only in the predicate.
func TestWriteGoldenRefusesEvidence(t *testing.T) {
	// Inside the corpus, so containment passes and the evidence rule is what
	// rejects it.
	err := WriteGolden(filepath.Join(TestdataDir, "upstream", "get_tasks", "001_response.xml"), []byte("x"))
	if err == nil {
		t.Fatalf("Expected WriteGolden to refuse a response document")
	}
	if !strings.Contains(err.Error(), "evidence") {
		t.Fatalf("Expected the error to explain why, got: %v", err)
	}

	// And the file is untouched.
	if b, rerr := os.ReadFile(filepath.Join(TestdataDir, "upstream", "get_tasks", "001_response.xml")); rerr == nil && string(b) == "x" {
		t.Fatalf("WriteGolden overwrote the evidence despite refusing")
	}
}

// A symlink out of the corpus must not widen the boundary.
func TestUpdateRefusesSymlinkEscape(t *testing.T) {
	outside := t.TempDir()
	link := filepath.Join(TestdataDir, "zz_escape_test")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("Skipping: cannot create a symlink (%v)", err)
	}
	t.Cleanup(func() { _ = os.Remove(link) })

	if err := allowedUpdatePath(filepath.Join(link, "x_fidelity.txt")); err == nil {
		t.Fatalf("Expected -update to refuse a path escaping via a symlink")
	}
}
