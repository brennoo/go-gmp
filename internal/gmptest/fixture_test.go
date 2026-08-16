package gmptest

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// writeMetaFixture writes meta as f's sidecar so loadMeta has something to read.
func writeMetaFixture(t *testing.T, f Fixture, meta Meta) {
	t.Helper()
	b, err := json.Marshal(meta)
	if err != nil {
		t.Fatalf("marshal meta: %v", err)
	}
	if err := os.WriteFile(f.MetaPath(), b, 0o600); err != nil {
		t.Fatalf("write meta: %v", err)
	}
}

func TestLoadMetaRejectsRuntimeProvenanceGaps(t *testing.T) {
	tests := []struct {
		name    string
		prov    Provenance
		meta    Meta
		wantErr string
	}{
		{
			name:    "live missing gmp_version",
			prov:    Live,
			meta:    Meta{Command: "get_version", Provenance: Live, GmpVersionSource: GmpVersionCaptured},
			wantErr: "missing gmp_version",
		},
		{
			name:    "live unknown gmp_version_source",
			prov:    Live,
			meta:    Meta{Command: "get_version", Provenance: Live, GmpVersion: "22.7", GmpVersionSource: "guessed"},
			wantErr: "unknown gmp_version_source",
		},
		{
			name:    "live unknown image_digest_source",
			prov:    Live,
			meta:    Meta{Command: "get_version", Provenance: Live, GmpVersion: "22.7", GmpVersionSource: GmpVersionCaptured},
			wantErr: "unknown image_digest_source",
		},
		{
			name: "live resolved but image_digest empty",
			prov: Live,
			meta: Meta{
				Command: "get_version", Provenance: Live,
				GmpVersion: "22.7", GmpVersionSource: GmpVersionCaptured,
				ImageDigestSource: ImageDigestResolved,
			},
			wantErr: "image_digest_source is \"resolved\" but image_digest is invalid",
		},
		{
			name: "live resolved but image_digest is garbage",
			prov: Live,
			meta: Meta{
				Command: "get_version", Provenance: Live,
				GmpVersion: "22.7", GmpVersionSource: GmpVersionCaptured,
				ImageDigestSource: ImageDigestResolved,
				ImageDigest:       "stable",
			},
			wantErr: "image_digest_source is \"resolved\" but image_digest is invalid",
		},
		{
			name: "live resolved but image_digest has no repository",
			prov: Live,
			meta: Meta{
				Command: "get_version", Provenance: Live,
				GmpVersion: "22.7", GmpVersionSource: GmpVersionCaptured,
				ImageDigestSource: ImageDigestResolved,
				ImageDigest:       "sha256:deadbeef",
			},
			wantErr: "image_digest_source is \"resolved\" but image_digest is invalid",
		},
		{
			name: "live legacy but image_digest set",
			prov: Live,
			meta: Meta{
				Command: "get_version", Provenance: Live,
				GmpVersion: "22.7", GmpVersionSource: GmpVersionCaptured,
				ImageDigestSource: ImageDigestLegacy,
				ImageDigest:       "example.com/gvmd@sha256:" + strings.Repeat("a", 64),
			},
			wantErr: "image_digest_source is \"legacy-unpinned\" but image_digest is set",
		},
		{
			name: "live valid resolved",
			prov: Live,
			meta: Meta{
				Command: "get_version", Provenance: Live,
				GmpVersion: "22.7", GmpVersionSource: GmpVersionCaptured,
				ImageDigestSource: ImageDigestResolved,
				ImageDigest:       "example.com/gvmd@sha256:" + strings.Repeat("a", 64),
			},
		},
		{
			name: "live valid legacy",
			prov: Live,
			meta: Meta{
				Command: "get_version", Provenance: Live,
				GmpVersion: "22.7", GmpVersionSource: GmpVersionCaptured,
				ImageDigestSource: ImageDigestLegacy,
			},
		},
		{
			name: "upstream claims gmp_version",
			prov: Upstream,
			meta: Meta{
				Command: "get_version", Provenance: Upstream,
				Source: "x", SourceRev: "y", License: "z",
				GmpVersion: "22.7",
			},
			wantErr: "must not claim runtime provenance",
		},
		{
			name: "upstream claims image_digest_source",
			prov: Upstream,
			meta: Meta{
				Command: "get_version", Provenance: Upstream,
				Source: "x", SourceRev: "y", License: "z",
				ImageDigestSource: ImageDigestLegacy,
			},
			wantErr: "must not claim runtime provenance",
		},
		{
			name: "upstream valid",
			prov: Upstream,
			meta: Meta{
				Command: "get_version", Provenance: Upstream,
				Source: "x", SourceRev: "y", License: "z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			f := Fixture{Command: "get_version", Provenance: tt.prov, ID: "001", Dir: dir}
			writeMetaFixture(t, f, tt.meta)

			err := f.loadMeta()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("loadMeta() = %v, want nil", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("loadMeta() = %v, want error containing %q", err, tt.wantErr)
			}
		})
	}
}
