package gmptest

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Update is set by -update and regenerates derived goldens in place.
var Update = flag.Bool("update", false,
	"regenerate derived goldens (fidelity reports and constructed request goldens)")

// requestsDir holds request goldens built from Go values rather than captured,
// so unlike the other trees they are derived and may be regenerated.
const requestsDir = "requests"

// WriteGolden writes a derived golden, refusing any path that is evidence
// rather than a derivation.
//
// The rule this enforces: -update may rewrite a fidelity report or a
// constructed request golden, and may never rewrite a captured or transcribed
// NNN_request.xml or NNN_response.xml. Those bytes are the observation the
// whole corpus rests on. Without the check, the failure mode is a golden test
// going green because it quietly rewrote the fixture to match a regression —
// silent, and it destroys the evidence that would have revealed it.
func WriteGolden(path string, content []byte) error {
	if err := allowedUpdatePath(path); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o600)
}

// allowedUpdatePath reports whether path is a derived golden inside the
// corpus.
//
// Containment is checked before the name, so a matching suffix somewhere else
// on the filesystem is not enough: only files under TestdataDir qualify. A
// guarantee described as an evidence boundary should not be satisfiable by
// naming a file in /tmp.
func allowedUpdatePath(path string) error {
	root, err := filepath.Abs(TestdataDir)
	if err != nil {
		return fmt.Errorf("resolve %s: %w", TestdataDir, err)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve %s: %w", path, err)
	}

	// EvalSymlinks on the nearest existing ancestor, so a symlink cannot point
	// the writer out of the tree. The leaf itself may legitimately not exist.
	resolved, err := resolveExistingPrefix(abs)
	if err != nil {
		return err
	}
	if realRoot, err := filepath.EvalSymlinks(root); err == nil {
		root = realRoot
	}

	rel, err := filepath.Rel(root, resolved)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("refusing to write %s: outside the fixture corpus at %s", path, TestdataDir)
	}

	base := filepath.Base(abs)
	slash := filepath.ToSlash(rel)
	switch {
	case strings.HasSuffix(base, "_fidelity.txt"):
		return nil
	case slash == "FIDELITY-SUMMARY.txt":
		return nil
	case slash == "REQUEST-FIDELITY.txt":
		return nil
	case strings.HasPrefix(slash, requestsDir+"/"):
		return nil
	}
	return fmt.Errorf(
		"refusing to write %s under -update: only fidelity reports and constructed "+
			"request goldens are derived; captured request/response documents are evidence", path)
}

// resolveExistingPrefix resolves symlinks on the deepest existing ancestor of
// path and re-joins the remainder.
func resolveExistingPrefix(abs string) (string, error) {
	missing := []string{}
	cur := abs
	for {
		if _, err := os.Lstat(cur); err == nil {
			real, err := filepath.EvalSymlinks(cur)
			if err != nil {
				return "", fmt.Errorf("resolve %s: %w", cur, err)
			}
			return filepath.Join(append([]string{real}, reverse(missing)...)...), nil
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return abs, nil // nothing on this path exists; nothing to spoof
		}
		missing = append(missing, filepath.Base(cur))
		cur = parent
	}
}

func reverse(in []string) []string {
	out := make([]string, len(in))
	for i, s := range in {
		out[len(in)-1-i] = s
	}
	return out
}
