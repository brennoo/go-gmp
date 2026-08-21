package gmptest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/brennoo/go-gmp/internal/imageref"
)

// Provenance says where a fixture's bytes came from, which governs what may be
// asserted about them.
type Provenance string

const (
	// Upstream fixtures are transcribed from gvmd's GMP schema. AGPL-derived
	// and documentation formatting rather than wire bytes, so no
	// byte-exactness assertion may run against them.
	Upstream Provenance = "upstream"
	// Live fixtures are captured from a real gvmd over the wire.
	Live Provenance = "live"
)

// TestdataDir is where the corpus lives, relative to this package.
const TestdataDir = "testdata"

// GmpVersionSource says how a fixture's GmpVersion was determined.
type GmpVersionSource string

const (
	// GmpVersionCaptured means this run's own get_version call recorded it.
	GmpVersionCaptured GmpVersionSource = "captured"
	// GmpVersionInferred means it was backfilled rather than recorded by the
	// capturing run. Nothing produces it yet; exists so a future backfill
	// has to say so instead of passing itself off as captured.
	GmpVersionInferred GmpVersionSource = "inferred-from-corpus"
)

var validGmpVersionSources = map[GmpVersionSource]bool{
	GmpVersionCaptured: true,
	GmpVersionInferred: true,
}

// ImageDigestSource says how a live fixture's ImageDigest field was
// determined, and is required on every live fixture regardless of whether
// ImageDigest itself is set.
type ImageDigestSource string

const (
	// ImageDigestResolved means internal/dockerstack read a real digest off
	// the running container at capture time; ImageDigest must be set.
	ImageDigestResolved ImageDigestSource = "resolved"
	// ImageDigestLegacy marks a fixture captured before T-16c pinned the
	// integration stack. Permanent classification, not a placeholder to
	// backfill — there's no way to recover what image a past capture ran
	// against. ImageDigest must be empty.
	ImageDigestLegacy ImageDigestSource = "legacy-unpinned"
)

var validImageDigestSources = map[ImageDigestSource]bool{
	ImageDigestResolved: true,
	ImageDigestLegacy:   true,
}

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
	ByteFaithful bool   `json:"byte_faithful"`
	CapturedAt   string `json:"captured_at,omitempty"`
	// GmpVersion is the GMP protocol version get_version reported, not a gvmd
	// build version — GMP exposes no such thing. See docs/review-2026-08.md, T-16a.
	GmpVersion string `json:"gmp_version,omitempty"`
	// GmpVersionSource is required whenever GmpVersion is set; see GmpVersionSource's doc.
	GmpVersionSource GmpVersionSource `json:"gmp_version_source,omitempty"`
	// ImageDigest is the resolved gvmd image digest the fixture was captured
	// against, read from the running container rather than inferred from a
	// tag; see internal/dockerstack. Set exactly when ImageDigestSource is
	// ImageDigestResolved.
	ImageDigest string `json:"image_digest,omitempty"`
	// ImageDigestSource is required on every live fixture; see its doc.
	ImageDigestSource ImageDigestSource `json:"image_digest_source,omitempty"`
	Redactions        []string          `json:"redactions,omitempty"`
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
	// Guards against a sidecar copied between commands, which would decode
	// with the wrong registry constructor while claiming the original.
	if f.Meta.Provenance != f.Provenance {
		return fmt.Errorf("%s/%s: meta says provenance %q but it is filed under %q",
			f.Command, f.ID, f.Meta.Provenance, f.Provenance)
	}
	if f.Meta.Command != f.Command {
		return fmt.Errorf("%s/%s: meta says command %q but it is filed under %q",
			f.Command, f.ID, f.Meta.Command, f.Command)
	}
	switch f.Provenance {
	case Upstream:
		return f.validateUpstreamMeta()
	case Live:
		return f.validateLiveMeta()
	}
	return nil
}

func (f *Fixture) validateUpstreamMeta() error {
	for name, val := range map[string]string{
		"source": f.Meta.Source, "source_upstream_revision": f.Meta.SourceRev,
		"source_license": f.Meta.License,
	} {
		if strings.TrimSpace(val) == "" {
			return fmt.Errorf("%s/%s: upstream fixture is missing %s, so its provenance is unverifiable",
				f.Command, f.ID, name)
		}
	}
	if f.Meta.GmpVersion != "" || f.Meta.GmpVersionSource != "" || f.Meta.ImageDigest != "" || f.Meta.ImageDigestSource != "" {
		return fmt.Errorf("%s/%s: upstream fixture must not claim runtime provenance", f.Command, f.ID)
	}
	return nil
}

func (f *Fixture) validateLiveMeta() error {
	if strings.TrimSpace(f.Meta.GmpVersion) == "" {
		return fmt.Errorf("%s/%s: live fixture is missing gmp_version", f.Command, f.ID)
	}
	if !validGmpVersionSources[f.Meta.GmpVersionSource] {
		return fmt.Errorf("%s/%s: live fixture has unknown gmp_version_source %q",
			f.Command, f.ID, f.Meta.GmpVersionSource)
	}
	if !validImageDigestSources[f.Meta.ImageDigestSource] {
		return fmt.Errorf("%s/%s: live fixture has unknown image_digest_source %q",
			f.Command, f.ID, f.Meta.ImageDigestSource)
	}
	switch f.Meta.ImageDigestSource {
	case ImageDigestResolved:
		if err := imageref.ValidateDigest(f.Meta.ImageDigest); err != nil {
			return fmt.Errorf("%s/%s: image_digest_source is %q but image_digest is invalid: %w",
				f.Command, f.ID, ImageDigestResolved, err)
		}
	case ImageDigestLegacy:
		if f.Meta.ImageDigest != "" {
			return fmt.Errorf("%s/%s: image_digest_source is %q but image_digest is set; use %q instead",
				f.Command, f.ID, ImageDigestLegacy, ImageDigestResolved)
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
