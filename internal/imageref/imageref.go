// Package imageref parses and validates Docker/OCI-style image references —
// the "repository[:tag][@digest]" shape used throughout hack/docker-compose.yml
// and this repo's own provenance metadata. It exists so the things that need
// to reason about that shape — internal/dockerstack's digest resolution,
// internal/cmd/pinimages' Compose-file rewriting, and internal/schema's and
// internal/gmptest's provenance validation — share one implementation instead
// of several independently-drifting ones.
package imageref

import (
	"fmt"
	"strings"
)

// Repository returns ref's repository, stripping a tag or digest — the
// standard Docker reference-parsing rule: a "@" always introduces a digest,
// and a ":" only introduces a tag when it occurs after the last "/"
// (otherwise it's part of a "host:port" registry address that has no tag at
// all).
func Repository(ref string) string {
	if repo, _, ok := strings.Cut(ref, "@"); ok {
		return repo
	}
	if slash, colon := strings.LastIndexByte(ref, '/'), strings.LastIndexByte(ref, ':'); colon > slash {
		return ref[:colon]
	}
	return ref
}

const hexDigits = "0123456789abcdef"

// ValidateDigest checks that ref is a complete, repository-qualified content
// digest — repository@sha256:<64 lowercase hex> — the exact shape
// internal/dockerstack.ImageDigest always produces and the only shape this
// codebase accepts as a resolved runtime image identity. Unlike a bare
// "contains @sha256:" check, it also rejects a syntactically invalid
// repository: an embedded tag, empty path components, or a malformed
// registry domain.
func ValidateDigest(ref string) error {
	repo, rest, ok := strings.Cut(ref, "@")
	if !ok {
		return fmt.Errorf("%q is not a resolved image digest; want repository@sha256:<64 lowercase hex>", ref)
	}
	hex, ok := strings.CutPrefix(rest, "sha256:")
	if !ok || len(hex) != 64 || strings.Trim(hex, hexDigits) != "" {
		return fmt.Errorf("%q is not a resolved image digest; want repository@sha256:<64 lowercase hex>", ref)
	}
	if err := validateRepository(repo); err != nil {
		return fmt.Errorf("%q has an invalid repository: %w", ref, err)
	}
	return nil
}

func validateRepository(repo string) error {
	if repo == "" {
		return fmt.Errorf("empty repository")
	}
	for i, c := range strings.Split(repo, "/") {
		if c == "" {
			return fmt.Errorf("empty path component in %q (a leading, trailing, or doubled slash)", repo)
		}
		if i == 0 && looksLikeDomain(c) {
			if err := validateDomain(c); err != nil {
				return err
			}
			continue
		}
		if err := validatePathComponent(c); err != nil {
			return err
		}
	}
	return nil
}

// looksLikeDomain applies Docker's own heuristic for telling a reference's
// leading path component from a registry domain: it contains a "." or ":",
// or is exactly "localhost".
func looksLikeDomain(c string) bool {
	return strings.ContainsAny(c, ".:") || c == "localhost"
}

func validateDomain(domain string) error {
	host, port, hasPort := strings.Cut(domain, ":")
	if hasPort && (port == "" || strings.Trim(port, "0123456789") != "") {
		return fmt.Errorf("invalid port in domain %q", domain)
	}
	for _, label := range strings.Split(host, ".") {
		if label == "" || !isDomainLabel(label) {
			return fmt.Errorf("invalid domain label %q in %q", label, domain)
		}
	}
	return nil
}

func isDomainLabel(s string) bool {
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-':
		default:
			return false
		}
	}
	return true
}

func validatePathComponent(c string) error {
	if strings.ContainsRune(c, ':') {
		return fmt.Errorf("path component %q contains ':' (an embedded tag)", c)
	}
	for _, r := range c {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '.', r == '_', r == '-':
		default:
			return fmt.Errorf("path component %q has an invalid character %q", c, r)
		}
	}
	return nil
}
