package gmptest

import (
	"encoding/xml"
	"fmt"
	"slices"
	"sort"
	"strings"
)

// Fidelity is what a document loses, gains, or alters on a round trip through
// a Go type.
//
// It is computed by unmarshalling the document into the type and marshalling it
// back, then diffing canonical path multisets and values — not by reflecting
// over struct tags. Reflection would mean reimplementing encoding/xml's
// matching rules (first-match-wins, local-name matching, embedded structs,
// ,any, ,chardata), and every divergence from the real decoder is a silent
// false result. The round trip measures the property that actually matters:
// did information survive.
type Fidelity struct {
	// DecodeErr is set when the document does not unmarshal at all. The other
	// fields are then meaningless. This is recorded rather than treated as a
	// failure: a response the library cannot decode is the most valuable thing
	// the corpus can tell us, and it belongs in the inventory.
	DecodeErr error

	Dropped    []Finding // present in source, absent after the round trip
	Truncated  []Finding // survived, but fewer times than the source had
	Added      []Finding // invented by the round trip; source had none
	Duplicated []Finding // emitted more often than the source had
	Changed    []Finding // same path, different value

	Namespaces []string
}

// Finding is one path and how it changed.
type Finding struct {
	Path string
	Src  int
	Got  int
	// Suppressed counts descendant paths folded into this one. Only set on
	// Dropped: if an element is absent from the output then so is its whole
	// subtree, and listing every descendant buries the finding.
	Suppressed int
	// SrcValue and GotValue are set on Changed.
	SrcValue string
	GotValue string
}

// Advisory reports whether the comparison should be trusted. Namespaced
// documents are not: encoding/xml does not round-trip prefixes and URIs
// faithfully, so both a clean result and a dirty one may be artifacts.
func (f Fidelity) Advisory() bool { return len(f.Namespaces) > 0 }

// Clean reports whether the round trip preserved everything. Added and
// Duplicated do not count as loss; Changed does, because a value the library
// cannot reproduce is information it did not keep.
//
// An advisory comparison is never clean — reporting "nothing lost" from a
// diff we have just said is untrustworthy would be the worst of both.
func (f Fidelity) Clean() bool {
	return f.DecodeErr == nil && !f.Advisory() &&
		len(f.Dropped) == 0 && len(f.Truncated) == 0 && len(f.Changed) == 0
}

// Compare round-trips doc through v and reports what changed.
//
// v must be a pointer to a freshly zeroed value; reusing one across calls
// leaks state between fixtures.
func Compare(doc []byte, v any) Fidelity {
	var f Fidelity

	src, err := ExtractPaths(doc)
	if err != nil {
		f.DecodeErr = fmt.Errorf("source is not well-formed XML: %w", err)
		return f
	}
	f.Namespaces = src.Namespaces

	if err := xml.Unmarshal(doc, v); err != nil {
		f.DecodeErr = err
		return f
	}

	out, err := xml.Marshal(v)
	if err != nil {
		f.DecodeErr = fmt.Errorf("re-marshal: %w", err)
		return f
	}
	got, err := ExtractPaths(out)
	if err != nil {
		f.DecodeErr = fmt.Errorf("re-marshalled output is not well-formed: %w", err)
		return f
	}
	for _, ns := range got.Namespaces {
		if !slices.Contains(f.Namespaces, ns) {
			f.Namespaces = append(f.Namespaces, ns)
		}
	}
	sort.Strings(f.Namespaces)

	f.diff(src, got)
	return f
}

// diff fills the buckets from the two path sets.
func (f *Fidelity) diff(src, got PathSet) {
	for _, path := range src.Paths() {
		s, g := src.Counts[path], got.Counts[path]
		switch {
		case g == 0:
			f.Dropped = append(f.Dropped, Finding{Path: path, Src: s})
		case g < s:
			f.Truncated = append(f.Truncated, Finding{Path: path, Src: s, Got: g})
		case g > s:
			f.Duplicated = append(f.Duplicated, Finding{Path: path, Src: s, Got: g})
		}

		// Values are only comparable where the path occurs exactly once in
		// both documents. With no positional index there is no way to align
		// the nth occurrence of a repeated path across a reordering round trip.
		sv, sok := src.SoleValue(path)
		gv, gok := got.SoleValue(path)
		if sok && gok && sv != gv {
			f.Changed = append(f.Changed, Finding{Path: path, SrcValue: sv, GotValue: gv})
		}
	}
	for _, path := range got.Paths() {
		if src.Counts[path] == 0 {
			f.Added = append(f.Added, Finding{Path: path, Got: got.Counts[path]})
		}
	}
	f.Dropped = suppressDescendants(f.Dropped)
}

// suppressDescendants keeps only the shallowest dropped path of each subtree,
// recording how many descendants it stands for.
//
// This checks every ancestor prefix rather than only the previously emitted
// finding. Lexicographic order does not place a subtree contiguously: '-' sorts
// before '/', so a sibling like /a/b-c falls between /a/b and /a/b/d, and
// comparing only against the last emitted path would leave /a/b/d unfolded.
// Element names containing '-' are real — the corpus has <cpe-item>.
//
// Folding is sound only because it runs on fully dropped paths: if /a/b has
// zero occurrences in the output then no /a/b/c can exist there either, since
// it would need a surviving parent.
func suppressDescendants(in []Finding) []Finding {
	if len(in) == 0 {
		return in
	}
	sort.Slice(in, func(i, j int) bool { return in[i].Path < in[j].Path })

	var out []Finding
	at := map[string]int{} // emitted path -> index in out

	for _, f := range in {
		if i, ok := shallowestAncestor(f.Path, at); ok {
			out[i].Suppressed++
			continue
		}
		at[f.Path] = len(out)
		out = append(out, f)
	}
	return out
}

// shallowestAncestor finds the outermost already-emitted path that contains
// path. Deeper ancestors cannot be in the map: they would themselves have been
// folded into that same outermost one.
func shallowestAncestor(path string, at map[string]int) (int, bool) {
	for i := 1; i < len(path); i++ {
		if path[i] != '/' {
			continue
		}
		if idx, ok := at[path[:i]]; ok {
			return idx, true
		}
	}
	return 0, false
}

// Report renders f as the deterministic text that gets committed as a golden.
//
// The assertion on these files is "unchanged", not "empty". 103 bool fields
// tagged omitempty cannot round-trip a false value, so an empty report is
// unreachable until that is fixed, and a threshold would need permanent
// tuning. Committing the report instead makes volume irrelevant: a struct fix
// shows up as a shrinkage diff, and a server that starts sending a new element
// shows up as one added line.
func (f Fidelity) Report(header string) string {
	var b strings.Builder
	b.WriteString(header)

	if f.DecodeErr != nil {
		fmt.Fprintf(&b, "decode: error\nerror: %s\n", f.DecodeErr)
		return b.String()
	}
	fmt.Fprintf(&b, "decode: ok\nnamespaces-seen: %s\n", strings.Join(f.Namespaces, " "))
	if f.Advisory() {
		b.WriteString("ADVISORY: this document uses namespaces, which encoding/xml does not\n" +
			"round-trip faithfully. Treat every count below as unreliable.\n")
	}

	section := func(name string, findings []Finding, render func(Finding) string) {
		fmt.Fprintf(&b, "\n%s\n", name)
		if len(findings) == 0 {
			b.WriteString("  (none)\n")
			return
		}
		lines := make([]string, len(findings))
		for i, x := range findings {
			lines[i] = render(x)
		}
		sort.Strings(lines)
		for _, l := range lines {
			fmt.Fprintf(&b, "  %s\n", l)
		}
	}

	section("DROPPED", f.Dropped, func(x Finding) string {
		s := fmt.Sprintf("%-72s x%d", x.Path, x.Src)
		if x.Suppressed > 0 {
			s += fmt.Sprintf("  (+%d descendant paths)", x.Suppressed)
		}
		return s
	})
	section("TRUNCATED", f.Truncated, func(x Finding) string {
		return fmt.Sprintf("%-72s %d -> %d", x.Path, x.Src, x.Got)
	})
	section("CHANGED", f.Changed, func(x Finding) string {
		return fmt.Sprintf("%-72s %q -> %q", x.Path, x.SrcValue, x.GotValue)
	})
	section("ADDED", f.Added, func(x Finding) string {
		return fmt.Sprintf("%-72s x%d", x.Path, x.Got)
	})
	section("DUPLICATED", f.Duplicated, func(x Finding) string {
		return fmt.Sprintf("%-72s %d -> %d", x.Path, x.Src, x.Got)
	})

	return b.String()
}
