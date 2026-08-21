package imageref

import (
	"strings"
	"testing"
)

func TestRepository(t *testing.T) {
	tests := []struct {
		ref  string
		want string
	}{
		{"gvmd", "gvmd"},
		{"gvmd:23.0.0", "gvmd"},
		{"gvmd@sha256:cb82c501ca93809405a8d0ed4c21846703466802ce6462536e985625f102b497", "gvmd"},
		{"registry.community.greenbone.net:5000/community/gvmd:23.0.0", "registry.community.greenbone.net:5000/community/gvmd"},
		{"registry.community.greenbone.net:5000/community/gvmd", "registry.community.greenbone.net:5000/community/gvmd"},
		{"localhost:5000/community/gvmd:stable", "localhost:5000/community/gvmd"},
		{"localhost:5000/community/gvmd", "localhost:5000/community/gvmd"},
	}
	for _, tt := range tests {
		t.Run(tt.ref, func(t *testing.T) {
			if got := Repository(tt.ref); got != tt.want {
				t.Errorf("Repository(%q) = %q, want %q", tt.ref, got, tt.want)
			}
		})
	}
}

func TestValidateDigestAccepts(t *testing.T) {
	hex := strings.Repeat("ab", 32)
	tests := []string{
		"gvmd@sha256:" + hex,
		"example.test/gvmd@sha256:" + hex,
		"registry.community.greenbone.net/community/gvmd@sha256:" + hex,
		"registry.community.greenbone.net:5000/community/gvmd@sha256:" + hex,
		"localhost:5000/gvmd@sha256:" + hex,
	}
	for _, ref := range tests {
		t.Run(ref, func(t *testing.T) {
			if err := ValidateDigest(ref); err != nil {
				t.Errorf("ValidateDigest(%q) = %v, want nil", ref, err)
			}
		})
	}
}

func TestValidateDigestRejects(t *testing.T) {
	hex := strings.Repeat("ab", 32)
	tests := []struct {
		name string
		ref  string
	}{
		{"empty", ""},
		{"tag only, no digest", "gvmd:stable"},
		{"short digest", "gvmd@sha256:" + strings.Repeat("a", 10)},
		{"malformed hex", "gvmd@sha256:" + strings.Repeat("g", 64)},
		{"uppercase hex", "gvmd@sha256:" + strings.Repeat("AB", 32)},
		{"leading colon (empty repository, looks like a bare port)", ":stable@sha256:" + hex},
		{"embedded tag before the digest", "example.com/gvmd:stable@sha256:" + hex},
		{"leading slash (empty path component)", "/gvmd@sha256:" + hex},
		{"doubled slash (empty path component)", "example.com//gvmd@sha256:" + hex},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateDigest(tt.ref); err == nil {
				t.Errorf("ValidateDigest(%q) = nil, want an error", tt.ref)
			}
		})
	}
}
