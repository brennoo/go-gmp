package gmptest

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"
)

// transcribe regenerates testdata/upstream from the vendored GMP schema. It is
// off by default and skips entirely when vendor/ is absent, because vendor/ is
// gitignored and most checkouts will not have it.
var transcribe = flag.Bool("transcribe", false,
	"regenerate testdata/upstream from vendor/XML (requires vendor/, which is gitignored)")

const (
	vendorCommandsDir = "../../vendor/XML/commands"

	// upstreamRevision identifies the schema snapshot. vendor/ carries no
	// version file, so this is derived from the newest entry in changes.xml.
	// Note it is one release ahead of the compatibility claim in README.md.
	upstreamRevision = "GMP 22.7 (newest entry in vendor/XML/commands/changes.xml)"

	upstreamLicense = "AGPL-3.0-or-later, Copyright (C) Greenbone AG"

	// gmpVersionLiteral replaces @GMP_VERSION@, which CMake would substitute at
	// gvmd build time.
	gmpVersionLiteral = "22.7"
)

// schemaFile is the subset of a command_*.xml we need. innerxml keeps the
// author's exact bytes instead of re-serializing them through encoding/xml,
// which would silently normalize quoting and self-closing forms.
type schemaFile struct {
	Commands []struct {
		Name     string `xml:"name"`
		Examples []struct {
			Summary  string `xml:"summary"`
			Request  inner  `xml:"request"`
			Response inner  `xml:"response"`
		} `xml:"example"`
	} `xml:"command"`
}

type inner struct {
	XML string `xml:",innerxml"`
}

// Marker handling. GMP.xml.in is a CMake template and the substitution markers
// appear inside the example payloads, so the transcribed documents are not
// valid XML until they are resolved. We resolve as a gvmd built with CVSS3
// ratings and agents enabled, which is the default configuration.
var (
	ifNotBlock = regexp.MustCompile(`(?s)@IF_NOT_CVSS3_RATINGS@.*?@ENDIF_NOT_CVSS3_RATINGS@`)
	markerRe   = regexp.MustCompile(`@(?:END)?IF_[A-Z0-9_]+@`)
)

func resolveMarkers(s string) (string, []string) {
	var applied []string

	if ifNotBlock.MatchString(s) {
		s = ifNotBlock.ReplaceAllString(s, "")
		applied = append(applied, "cvss3_ratings=on (dropped @IF_NOT_CVSS3_RATINGS@ blocks)")
	}
	if strings.Contains(s, "@IF_CVSS3_RATINGS@") {
		applied = append(applied, "cvss3_ratings=on (kept @IF_CVSS3_RATINGS@ blocks)")
	}
	if strings.Contains(s, "@IF_ENABLE_AGENTS@") {
		applied = append(applied, "enable_agents=on (kept @IF_ENABLE_AGENTS@ blocks)")
	}
	if strings.Contains(s, "@GMP_VERSION@") {
		s = strings.ReplaceAll(s, "@GMP_VERSION@", gmpVersionLiteral)
		applied = append(applied, "gmp_version="+gmpVersionLiteral)
	}
	if markerRe.MatchString(s) {
		s = markerRe.ReplaceAllString(s, "")
	}

	sort.Strings(applied)
	return s, applied
}

func TestTranscribeUpstream(t *testing.T) {
	if !*transcribe {
		t.Skip("Skipping: pass -transcribe to regenerate testdata/upstream")
	}
	if _, err := os.Stat(vendorCommandsDir); err != nil {
		t.Skipf("Skipping: %s is absent (vendor/ is gitignored)", vendorCommandsDir)
	}

	files, err := filepath.Glob(filepath.Join(vendorCommandsDir, "command_*.xml"))
	if err != nil || len(files) == 0 {
		t.Fatalf("Found no command files in %s (%v)", vendorCommandsDir, err)
	}

	root := filepath.Join(TestdataDir, string(Upstream))
	if err := os.RemoveAll(root); err != nil {
		t.Fatalf("Failed to clear %s: %v", root, err)
	}

	var written, skipped int
	for _, path := range files {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("Failed to read %s: %v", path, err)
		}

		var sf schemaFile
		if err := xml.Unmarshal(raw, &sf); err != nil {
			t.Errorf("Failed to parse %s: %v", path, err)
			continue
		}

		for _, cmd := range sf.Commands {
			name := strings.ToLower(strings.TrimSpace(cmd.Name))
			if _, known := Registry[name]; !known {
				// Not a root command with a matching *_response type. Nothing
				// to compare it against, so there is no fixture to write.
				skipped++
				continue
			}

			for i, ex := range cmd.Examples {
				if err := writeExample(root, name, i, filepath.Base(path), cmd.Name, ex.Request.XML, ex.Response.XML); err != nil {
					t.Fatalf("Failed to write %s example %d: %v", name, i+1, err)
				}
				written++
			}
		}
	}

	if err := writeUpstreamManifest(root, written); err != nil {
		t.Fatalf("Failed to write the manifest: %v", err)
	}

	t.Logf("Transcribed %d fixtures (%d commands had no registry entry)", written, skipped)
	if written == 0 {
		t.Fatalf("Transcription produced nothing")
	}
}

// writeExample transcribes one example into the fixture tree.
func writeExample(root, name string, idx int, srcFile, srcCmd, rawReq, rawResp string) error {
	id := fmt.Sprintf("%03d", idx+1)
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}

	resp, transforms := resolveMarkers(rawResp)
	req, reqTransforms := resolveMarkers(rawReq)
	for _, tr := range reqTransforms {
		if !slices.Contains(transforms, tr) {
			transforms = append(transforms, tr)
		}
	}
	sort.Strings(transforms)

	f := Fixture{Command: name, Provenance: Upstream, ID: id, Dir: dir}
	for path, body := range map[string]string{f.ResponsePath(): resp, f.RequestPath(): req} {
		if s := strings.TrimSpace(body); s != "" {
			if err := os.WriteFile(path, []byte(s+"\n"), 0o600); err != nil {
				return err
			}
		}
	}

	return f.WriteMeta(Meta{
		Command:      name,
		Provenance:   Upstream,
		Source:       fmt.Sprintf("%s, command %q, example %d", srcFile, srcCmd, idx+1),
		SourceRev:    upstreamRevision,
		License:      upstreamLicense,
		Transforms:   transforms,
		ByteFaithful: false,
	})
}

func writeUpstreamManifest(root string, n int) error {
	// A nested go.mod excludes this tree from the parent module's zip. testdata
	// is otherwise published inside the module tarball, and internal/ restricts
	// importability, not packaging — so without this the AGPL-derived corpus
	// would ship inside an MIT module.
	mod := "// This nested module exists solely to exclude this directory from the\n" +
		"// parent module's distribution zip. It is never built or imported.\n" +
		"module github.com/brennoo/go-gmp/internal/gmptest/testdata/upstream\n\ngo 1.23\n"
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte(mod), 0o644); err != nil {
		return err
	}

	manifest := fmt.Sprintf(`# Upstream fixture corpus

%d request/response examples transcribed from gvmd's GMP schema by
`+"`go test ./internal/gmptest -run TestTranscribeUpstream -transcribe`"+`.

## Provenance

- Source: `+"`vendor/XML/commands/command_*.xml`"+`, derived from gvmd's `+"`GMP.xml.in`"+`.
- Revision: %s
- License: %s

This directory is **AGPL-derived** and is deliberately isolated so it can be
removed as a unit without touching any code — the fixture loader discovers
files by walking the tree and never names one. The nested `+"`go.mod`"+` keeps it out
of the published module zip.

## These are not wire bytes

The examples are documentation. Their indentation and self-closing forms are
the schema author's, not gvmd's, and marker resolution edits them further, so
every fixture here records `+"`byte_faithful: false`"+`. **No byte-exactness
assertion may run against this tree** — only path-level fidelity, which is
insensitive to formatting.

GMP.xml.in is a CMake template and its substitution markers appear inside the
example payloads, so the raw examples are not valid XML. Transcription resolves
them as a gvmd built with CVSS3 ratings and agents enabled, and each fixture's
`+"`meta.json`"+` records exactly which transforms were applied.
`, n, upstreamRevision, upstreamLicense)

	return os.WriteFile(filepath.Join(root, "MANIFEST.md"), []byte(manifest), 0o644)
}
