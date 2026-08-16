package main

import (
	"errors"
	"testing"
)

func TestParseImageLine(t *testing.T) {
	tests := []struct {
		name           string
		line           string
		wantRepo       string
		wantTagComment string
		wantPinned     bool
		wantOK         bool
	}{
		{
			name:   "not an image line",
			line:   "    build: .",
			wantOK: false,
		},
		{
			name:           "unpinned with tag",
			line:           "    image: registry.community.greenbone.net/community/gvmd:stable",
			wantRepo:       "registry.community.greenbone.net/community/gvmd",
			wantTagComment: "stable",
			wantPinned:     false,
			wantOK:         true,
		},
		{
			name:       "unpinned without tag",
			line:       "    image: registry.community.greenbone.net/community/gvm-tools",
			wantRepo:   "registry.community.greenbone.net/community/gvm-tools",
			wantPinned: false,
			wantOK:     true,
		},
		{
			name:           "unpinned with tag behind a host:port registry",
			line:           "    image: localhost:5000/community/gvmd:stable",
			wantRepo:       "localhost:5000/community/gvmd",
			wantTagComment: "stable",
			wantPinned:     false,
			wantOK:         true,
		},
		{
			name:       "unpinned without tag behind a host:port registry",
			line:       "    image: localhost:5000/community/gvm-tools",
			wantRepo:   "localhost:5000/community/gvm-tools",
			wantPinned: false,
			wantOK:     true,
		},
		{
			name:           "pinned with was-tag comment",
			line:           "    image: gvmd@sha256:aaaa # was stable",
			wantRepo:       "gvmd",
			wantTagComment: "was stable",
			wantPinned:     true,
			wantOK:         true,
		},
		{
			name:           "pinned with untagged comment",
			line:           "    image: gvm-tools@sha256:bbbb # untagged, resolved from the implicit :latest",
			wantRepo:       "gvm-tools",
			wantTagComment: "untagged, resolved from the implicit :latest",
			wantPinned:     true,
			wantOK:         true,
		},
		{
			name:       "pinned with missing comment",
			line:       "    image: gvmd@sha256:aaaa",
			wantRepo:   "gvmd",
			wantPinned: true,
			wantOK:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, tagOrComment, pinned, ok := parseImageLine(tt.line)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !ok {
				return
			}
			if repo != tt.wantRepo || tagOrComment != tt.wantTagComment || pinned != tt.wantPinned {
				t.Errorf("parseImageLine(%q) = (%q, %q, %v), want (%q, %q, %v)",
					tt.line, repo, tagOrComment, pinned, tt.wantRepo, tt.wantTagComment, tt.wantPinned)
			}
		})
	}
}

// fakeResolve returns digest for ref, or errUnresolvable if ref is not a key
// of digests — a stand-in for resolveDigest that needs no Docker.
func fakeResolve(digests map[string]string) func(ref string) (string, error) {
	return func(ref string) (string, error) {
		digest, ok := digests[ref]
		if !ok {
			return "", errUnresolvable
		}
		return digest, nil
	}
}

var errUnresolvable = errors.New("no such ref")

func TestRefreshLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		digests map[string]string
		want    string // empty means "unchanged", checked only when neither wantErr nor wantAnyErr is set
		wantErr error  // set to assert errors.Is(err, wantErr)
		anyErr  bool   // set to assert only that refreshLine returned some error
	}{
		{
			name: "non-image line passes through unchanged",
			line: "    restart: unless-stopped",
		},
		{
			name:    "unpinned with tag is resolved and pinned",
			line:    "    image: gvmd:stable",
			digests: map[string]string{"gvmd:stable": "aaaa"},
			want:    "    image: gvmd@sha256:aaaa # was stable",
		},
		{
			name:    "unpinned without tag is resolved and pinned as untagged",
			line:    "    image: gvm-tools",
			digests: map[string]string{"gvm-tools": "bbbb"},
			want:    "    image: gvm-tools@sha256:bbbb # untagged, resolved from the implicit :latest",
		},
		{
			name:    "unpinned with tag behind a host:port registry is resolved and pinned",
			line:    "    image: localhost:5000/community/gvmd:stable",
			digests: map[string]string{"localhost:5000/community/gvmd:stable": "eeee"},
			want:    "    image: localhost:5000/community/gvmd@sha256:eeee # was stable",
		},
		{
			name:    "already pinned is re-resolved and rewritten when the digest moved",
			line:    "    image: gvmd@sha256:aaaa # was stable",
			digests: map[string]string{"gvmd:stable": "cccc"},
			want:    "    image: gvmd@sha256:cccc # was stable",
		},
		{
			name:    "already pinned is byte-identical when the digest has not moved",
			line:    "    image: gvmd@sha256:aaaa # was stable",
			digests: map[string]string{"gvmd:stable": "aaaa"},
		},
		{
			name:    "already pinned untagged is re-resolved against the bare repository",
			line:    "    image: gvm-tools@sha256:bbbb # untagged, resolved from the implicit :latest",
			digests: map[string]string{"gvm-tools": "dddd"},
			want:    "    image: gvm-tools@sha256:dddd # untagged, resolved from the implicit :latest",
		},
		{
			name:   "already pinned with a missing comment fails closed",
			line:   "    image: gvmd@sha256:aaaa",
			anyErr: true,
		},
		{
			name:   "already pinned with an unparseable comment fails closed",
			line:   "    image: gvmd@sha256:aaaa # who knows",
			anyErr: true,
		},
		{
			name:    "resolve failure propagates",
			line:    "    image: gvmd:stable",
			wantErr: errUnresolvable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := refreshLine(tt.line, fakeResolve(tt.digests))
			switch {
			case tt.anyErr:
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
			case tt.wantErr != nil:
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err = %v, want it to wrap %v", err, tt.wantErr)
				}
			default:
				if err != nil {
					t.Fatalf("refreshLine: %v", err)
				}
				want := tt.want
				if want == "" {
					want = tt.line
				}
				if got != want {
					t.Fatalf("got %q, want %q", got, want)
				}
			}
		})
	}
}
