package gmptest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Provenance says where a fixture's bytes came from, which governs what may be
// asserted about them.
type Provenance string

const (
	// Upstream fixtures are transcribed from gvmd's GMP schema. They are
	// AGPL-derived, live in their own directory so they can be removed as a
	// unit, and are documentation formatting rather than wire bytes — so no
	// byte-exactness assertion may run against them.
	Upstream Provenance = "upstream"
	// Live fixtures are captured from a real gvmd over the wire.
	Live Provenance = "live"
)

// TestdataDir is where the corpus lives, relative to this package.
const TestdataDir = "testdata"

// Meta is the sidecar recorded next to every fixture. Provenance is deliberately
// stored twice — here and in the directory layout — so that neither a moved
// file nor a lost sidecar can silently misattribute AGPL-derived content.
type Meta struct {
	Command    string     `json:"command"`
	Provenance Provenance `json:"provenance"`
	Source     string     `json:"source"`
	SourceRev  string     `json:"source_upstream_revision,omitempty"`
	License    string     `json:"source_license,omitempty"`
	Transforms []string   `json:"transforms,omitempty"`
	// ByteFaithful reports whether the bytes are exactly what crossed the wire.
	// False for upstream fixtures, whose formatting is the schema author's and
	// whose content is edited by marker stripping.
	ByteFaithful bool     `json:"byte_faithful"`
	CapturedAt   string   `json:"captured_at,omitempty"`
	GvmdVersion  string   `json:"gvmd_version,omitempty"`
	Redactions   []string `json:"redactions,omitempty"`
}

// Fixture is one request/response pair on disk.
type Fixture struct {
	Command    string
	Provenance Provenance
	ID         string // "001"
	Dir        string
	Meta       Meta
}

// RequestPath and ResponsePath are the evidence files. Nothing driven by
// -update may ever write to them; see allowedUpdatePath.
func (f Fixture) RequestPath() string  { return filepath.Join(f.Dir, f.ID+"_request.xml") }
func (f Fixture) ResponsePath() string { return filepath.Join(f.Dir, f.ID+"_response.xml") }
func (f Fixture) FidelityPath() string { return filepath.Join(f.Dir, f.ID+"_fidelity.txt") }
func (f Fixture) MetaPath() string     { return filepath.Join(f.Dir, f.ID+"_meta.json") }

// Response reads the captured response document.
func (f Fixture) Response() ([]byte, error) { return os.ReadFile(f.ResponsePath()) }

// Request reads the captured request document. Not every fixture has one.
func (f Fixture) Request() ([]byte, error) { return os.ReadFile(f.RequestPath()) }

// HasRequest reports whether a request document was recorded.
func (f Fixture) HasRequest() bool {
	_, err := os.Stat(f.RequestPath())
	return err == nil
}

// LoadFixtures discovers every fixture under root, walking the tree rather than
// consulting any list. That is what makes `rm -rf testdata/upstream` a
// supported operation needing no code change.
func LoadFixtures(root string) ([]Fixture, error) {
	var out []Fixture

	for _, prov := range []Provenance{Upstream, Live} {
		base := filepath.Join(root, string(prov))
		if _, err := os.Stat(base); err != nil {
			continue // this provenance is simply absent
		}

		cmdDirs, err := os.ReadDir(base)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", base, err)
		}
		for _, cd := range cmdDirs {
			if !cd.IsDir() {
				continue
			}
			dir := filepath.Join(base, cd.Name())
			entries, err := os.ReadDir(dir)
			if err != nil {
				return nil, fmt.Errorf("read %s: %w", dir, err)
			}
			for _, e := range entries {
				id, ok := strings.CutSuffix(e.Name(), "_response.xml")
				if !ok {
					continue
				}
				f := Fixture{Command: cd.Name(), Provenance: prov, ID: id, Dir: dir}
				if err := f.loadMeta(); err != nil {
					return nil, err
				}
				out = append(out, f)
			}
		}
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Provenance != out[j].Provenance {
			return out[i].Provenance < out[j].Provenance
		}
		if out[i].Command != out[j].Command {
			return out[i].Command < out[j].Command
		}
		return out[i].ID < out[j].ID
	})
	return out, nil
}

func (f *Fixture) loadMeta() error {
	b, err := os.ReadFile(f.MetaPath())
	if err != nil {
		return fmt.Errorf("read meta for %s/%s: %w", f.Command, f.ID, err)
	}
	if err := json.Unmarshal(b, &f.Meta); err != nil {
		return fmt.Errorf("parse meta for %s/%s: %w", f.Command, f.ID, err)
	}
	// Both identity fields are checked against the directory layout. A sidecar
	// copied between commands would otherwise let a fixture be decoded with the
	// wrong registry constructor while its metadata still claimed the original,
	// which is exactly the kind of drift the corpus exists to rule out.
	if f.Meta.Provenance != f.Provenance {
		return fmt.Errorf("%s/%s: meta says provenance %q but it is filed under %q",
			f.Command, f.ID, f.Meta.Provenance, f.Provenance)
	}
	if f.Meta.Command != f.Command {
		return fmt.Errorf("%s/%s: meta says command %q but it is filed under %q",
			f.Command, f.ID, f.Meta.Command, f.Command)
	}
	if f.Provenance == Upstream {
		for name, val := range map[string]string{
			"source": f.Meta.Source, "source_upstream_revision": f.Meta.SourceRev,
			"source_license": f.Meta.License,
		} {
			if strings.TrimSpace(val) == "" {
				return fmt.Errorf("%s/%s: upstream fixture is missing %s, so its provenance is unverifiable",
					f.Command, f.ID, name)
			}
		}
	}
	return nil
}

// WriteMeta persists m next to the fixture.
func (f Fixture) WriteMeta(m Meta) error {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.MetaPath(), append(b, '\n'), 0o600)
}
