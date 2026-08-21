package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// repoRoot is the module root, relative to this package's directory: the
// working directory the CLI's flag defaults are written for.
const repoRoot = "../../.."

func at(parts ...string) string {
	return filepath.Join(append([]string{repoRoot}, parts...)...)
}

// TestRunCheckOnCommittedTree is the CLI-level counterpart to
// codegen.CheckOutput's own regression test: `gengmp -check` must succeed
// against the tree exactly as committed.
func TestRunCheckOnCommittedTree(t *testing.T) {
	err := run([]string{
		"-check",
		"-manifest", at("internal", "codegen", "manifest.json"),
		"-schema", at("internal", "codegen", "testdata", "schema"),
		"-out", at("commands"),
	})
	if err != nil {
		t.Fatalf("gengmp -check reported the committed tree as out of date: %v", err)
	}
}

// TestRunCheckFailsOnDrift proves the exit path: a non-empty report has to
// come back as an error, or CI would pass on a stale tree.
func TestRunCheckFailsOnDrift(t *testing.T) {
	dir := t.TempDir()
	// One generated file, hand-edited; the other three simply absent.
	src, err := os.ReadFile(at("commands", "create_note.go"))
	if err != nil {
		t.Fatalf("Failed to read create_note.go: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "create_note.go"), append(src, '\n'), 0o600); err != nil {
		t.Fatalf("Failed to write the drifted copy: %v", err)
	}

	err = run([]string{
		"-check",
		"-manifest", at("internal", "codegen", "manifest.json"),
		"-schema", at("internal", "codegen", "testdata", "schema"),
		"-out", dir,
	})
	if err == nil || !strings.Contains(err.Error(), "out of date") {
		t.Fatalf("gengmp -check error = %v, want an out-of-date failure", err)
	}
}

// TestRunWriteRequiresOut keeps the destructive path's guard: without -check,
// a forgotten -out must fail rather than write into the current directory.
func TestRunWriteRequiresOut(t *testing.T) {
	if err := run(nil); err == nil || !strings.Contains(err.Error(), "-out is required") {
		t.Fatalf("run(nil) error = %v, want an -out complaint", err)
	}
}

// TestRunCheckNeedsAnExistingDir fails fast rather than reporting every
// generated file as "missing" when -out is a typo.
func TestRunCheckNeedsAnExistingDir(t *testing.T) {
	err := run([]string{
		"-check",
		"-manifest", at("internal", "codegen", "manifest.json"),
		"-schema", at("internal", "codegen", "testdata", "schema"),
		"-out", filepath.Join(t.TempDir(), "absent"),
	})
	if err == nil || !strings.Contains(err.Error(), "existing directory") {
		t.Fatalf("run error = %v, want a missing-directory complaint", err)
	}
}

// TestParseCheckDefaultsOut documents the one flag whose meaning differs
// between the two modes.
func TestParseCheckDefaultsOut(t *testing.T) {
	o, err := parse([]string{"-check"})
	if err != nil {
		t.Fatalf("parse(-check) failed: %v", err)
	}
	if o.out != checkOutDir {
		t.Errorf("parse(-check).out = %q, want %q", o.out, checkOutDir)
	}
	o, err = parse([]string{"-check", "-out", "elsewhere"})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if o.out != "elsewhere" {
		t.Errorf("An explicit -out was overridden: %q", o.out)
	}
}

func TestParseRejectsPositionalArgs(t *testing.T) {
	if _, err := parse([]string{"-check", "commands"}); err == nil ||
		!strings.Contains(err.Error(), "flags only") {
		t.Fatalf("parse error = %v, want a positional-argument complaint", err)
	}
}
