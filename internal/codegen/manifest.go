package codegen

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// manifestFile is the committed manifest's name, relative to this package's
// directory. It is the single declaration of "which commands are generated":
// the CLI, the drift check and every test in this package read their command
// list from it rather than carrying a second hardcoded copy.
const manifestFile = "manifest.json"

// Family is one command family: a noun (create_note, get_notes and
// modify_note are all about a "note") and the commands that belong to it.
//
// A family is also the unit Generate is invoked on — see GenerateAll — so a
// family's structural dedup scope, its type names and its output files are
// determined entirely by its own commands. Adding a family to the manifest can
// therefore never change another family's output.
type Family struct {
	Name     string   `json:"name"`
	Commands []string `json:"commands"`

	// ReviewOnly marks a family that is generated but whose output is not
	// committed to commands/ yet: it exists as a golden under
	// testdata/golden, to be diffed against the hand-written types and
	// conformance-checked, before anyone replaces the real files.
	//
	// The flag is what makes "generate and review" a state the manifest can
	// hold. Without it, adding a family here would immediately mean two
	// things at once — "generate this" and "commands/ must already contain
	// the result" — so CheckOutput would call the still-hand-written files
	// stale and hack/generate.sh would overwrite them unreviewed. Both the
	// CLI's write mode and CheckOutput therefore work from Adopted(), while
	// the golden test and GenerateAll see every family.
	//
	// Adopting a family is then a one-line manifest edit (drop the flag) plus
	// a run of hack/generate.sh, which is exactly the reviewable diff the
	// swap-in step wants.
	ReviewOnly bool `json:"review_only,omitempty"`
}

// Manifest is the decoded manifest.json.
type Manifest struct {
	Families []Family `json:"families"`
}

// LoadManifest reads and validates the manifest at path.
//
// Validation is deliberately strict — unknown JSON fields, an unsorted command
// list or a command claimed by two families are all errors rather than
// something quietly normalised — so the checked-in file stays self-disciplined
// and reviewable as a diff.
func LoadManifest(path string) (Manifest, error) {
	data, err := os.ReadFile(path) //nolint:gosec // an operator-supplied manifest path
	if err != nil {
		return Manifest{}, fmt.Errorf("codegen: reading the manifest: %w", err)
	}
	var m Manifest
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&m); err != nil {
		return Manifest{}, fmt.Errorf("codegen: %s: %w", path, err)
	}
	// DisallowUnknownFields only guards the fields of one JSON document; a
	// second document appended after the manifest (accidental concatenation,
	// a stray merge conflict marker left as text, ...) would otherwise be
	// silently ignored, which is exactly the kind of quietly-normalised input
	// the "deliberately strict" validation above is trying to rule out. Decode
	// into `any`, not another Manifest, so this check reports "there is more
	// input" rather than DisallowUnknownFields' unrelated complaint about the
	// second document's own shape.
	var trailing any
	if err := dec.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return Manifest{}, fmt.Errorf("codegen: %s: multiple JSON documents", path)
		}
		return Manifest{}, fmt.Errorf("codegen: %s: trailing data after the manifest: %w", path, err)
	}
	if err := m.validate(); err != nil {
		return Manifest{}, fmt.Errorf("codegen: %s: %w", path, err)
	}
	return m, nil
}

func (m Manifest) validate() error {
	if len(m.Families) == 0 {
		return fmt.Errorf("the manifest declares no families")
	}
	seenFamily := map[string]bool{}
	owner := map[string]string{}
	for _, f := range m.Families {
		if f.Name == "" {
			return fmt.Errorf("a family has no name")
		}
		if seenFamily[f.Name] {
			return fmt.Errorf("family %q is declared twice", f.Name)
		}
		seenFamily[f.Name] = true
		if err := f.validateCommands(owner); err != nil {
			return err
		}
	}
	return nil
}

// validateCommands checks one family's command list and records each command's
// owning family in owner, which is what makes "a command belongs to exactly
// one family" enforceable across the whole manifest.
func (f Family) validateCommands(owner map[string]string) error {
	if len(f.Commands) == 0 {
		return fmt.Errorf("family %q declares no commands", f.Name)
	}
	if !sort.StringsAreSorted(f.Commands) {
		return fmt.Errorf("family %q: commands must be sorted alphabetically for a stable diff", f.Name)
	}
	for _, c := range f.Commands {
		if c == "" {
			return fmt.Errorf("family %q has an empty command name", f.Name)
		}
		if prev, ok := owner[c]; ok {
			return fmt.Errorf("command %q is claimed by both family %q and family %q", c, prev, f.Name)
		}
		owner[c] = f.Name
	}
	return nil
}

// SnapshotPath is where one command's Stage-A JSON lives under dir. It is the
// package-level SnapshotPath convention, exposed as a method so manifest-driven
// callers do not need to reach past the manifest for it.
func (m Manifest) SnapshotPath(dir, command string) string {
	return SnapshotPath(dir, command)
}

// SharedFile is the file holding the nested types more than one of the
// family's commands reach — the structural dedup pass's output. The name is a
// convention, not a manifest field: it must agree with what Generate derives
// from the command names, and GenerateAll fails loudly if it ever does not.
func (f Family) SharedFile() string {
	return f.Name + sharedFileSuffix
}

// CommandFile is the generated file holding one command's request and response
// types.
func CommandFile(command string) string {
	return command + ".go"
}

// ExpectedOutputFiles is the machine-readable list of paths under outDir that
// this manifest can account for: every command's file plus every family's
// shared file, sorted.
//
// It is the conventional maximum. A family whose commands share no nested
// shape produces no shared file, so the set actually emitted by GenerateAll is
// a subset of this one; CheckOutput uses that emitted set for the missing and
// modified categories and this one only to bound what may legitimately exist.
func (m Manifest) ExpectedOutputFiles(outDir string) []string {
	var out []string
	for _, f := range m.Families {
		for _, c := range f.Commands {
			out = append(out, filepath.Join(outDir, CommandFile(c)))
		}
		out = append(out, filepath.Join(outDir, f.SharedFile()))
	}
	sort.Strings(out)
	return out
}

// AllCommands is every command in the manifest, flattened and sorted, for
// feeding LoadSnapshots in one call.
func (m Manifest) AllCommands() []string {
	var out []string
	for _, f := range m.Families {
		out = append(out, f.Commands...)
	}
	sort.Strings(out)
	return out
}

// Adopted is m without its review-only families: the families whose generated
// output is committed under commands/ and is therefore the CLI's and
// CheckOutput's business. It is a value copy sharing the same Family slice
// entries, so it is safe to call repeatedly and never mutates m.
func (m Manifest) Adopted() Manifest {
	out := Manifest{Families: make([]Family, 0, len(m.Families))}
	for _, f := range m.Families {
		if !f.ReviewOnly {
			out.Families = append(out.Families, f)
		}
	}
	return out
}

// ReviewOnlyFamilies is the complement of Adopted, for the one-line notice the
// CLI prints so a skipped family is visible rather than silently absent.
func (m Manifest) ReviewOnlyFamilies() []string {
	var out []string
	for _, f := range m.Families {
		if f.ReviewOnly {
			out = append(out, f.Name)
		}
	}
	sort.Strings(out)
	return out
}

// FamilyOf reports which family owns a command.
func (m Manifest) FamilyOf(command string) (Family, bool) {
	for _, f := range m.Families {
		for _, c := range f.Commands {
			if c == command {
				return f, true
			}
		}
	}
	return Family{}, false
}
