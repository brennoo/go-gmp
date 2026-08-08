package gmptest

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

// PathSet is a document reduced to canonical paths, their occurrence counts,
// and their values.
//
// A path is element names joined by "/", optionally ending in "/@attr" for an
// attribute or "/#text" for non-empty character data:
//
//	/get_tasks_response/task/@id
//	/get_tasks_response/task/name/#text
//
// An element in a namespace is written {uri}local. Collapsing namespaces to
// local names would merge <a:x/> and <b:x/> into one path and let a document
// that lost one of them still report clean.
//
// Paths carry no positional index. xml.Marshal emits fields in struct
// declaration order, so element order does not survive a round trip and any
// positional scheme would report reordering as loss.
type PathSet struct {
	Counts map[string]int
	// Values holds the text behind each path, in document order. Only paths
	// occurring exactly once are comparable between two documents — with no
	// positional index there is no way to align the nth occurrence of a
	// repeated path across a reordering round trip.
	Values map[string][]string
	// Namespaces lists every distinct non-empty namespace URI seen. GMP is
	// namespace-free in practice; when this is non-empty the comparison is
	// marked advisory, because encoding/xml does not round-trip namespace
	// prefixes faithfully and the resulting diff cannot be trusted either way.
	Namespaces []string
}

// Paths returns the paths in lexicographic order.
func (p PathSet) Paths() []string {
	out := make([]string, 0, len(p.Counts))
	for path := range p.Counts {
		out = append(out, path)
	}
	sort.Strings(out)
	return out
}

// SoleValue returns the value of a path that occurs exactly once.
func (p PathSet) SoleValue(path string) (string, bool) {
	if v, ok := p.Values[path]; ok && len(v) == 1 {
		return v[0], true
	}
	return "", false
}

// qualify renders an xml.Name as it appears in a path.
func qualify(n xml.Name) string {
	if n.Space == "" {
		return n.Local
	}
	return "{" + n.Space + "}" + n.Local
}

// ExtractPaths walks doc and records every element, attribute and non-empty
// text node.
//
// xmlns and xmlns:* attributes are skipped: encoding/xml synthesizes them on
// re-marshal even when the source had none, so counting them would manufacture
// differences that mean nothing. The namespaces they declare are recorded on
// the PathSet instead.
func ExtractPaths(doc []byte) (PathSet, error) {
	ps := PathSet{Counts: map[string]int{}, Values: map[string][]string{}}
	seenNS := map[string]bool{}

	noteNS := func(space string) {
		if space != "" && !seenNS[space] {
			seenNS[space] = true
			ps.Namespaces = append(ps.Namespaces, space)
		}
	}

	dec := xml.NewDecoder(bytes.NewReader(doc))
	var stack []string

	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return ps, fmt.Errorf("extract paths: %w", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			stack = append(stack, qualify(t.Name))
			here := "/" + strings.Join(stack, "/")
			ps.Counts[here]++
			noteNS(t.Name.Space)

			for _, a := range t.Attr {
				if a.Name.Local == "xmlns" || a.Name.Space == "xmlns" {
					noteNS(a.Value)
					continue
				}
				noteNS(a.Name.Space)
				p := here + "/@" + qualify(a.Name)
				ps.Counts[p]++
				ps.Values[p] = append(ps.Values[p], a.Value)
			}

		case xml.CharData:
			// Whitespace-only text is formatting, not content. Indented
			// fixtures would otherwise report text on every container element.
			if len(stack) > 0 {
				if s := strings.TrimSpace(string(t)); s != "" {
					p := "/" + strings.Join(stack, "/") + "/#text"
					ps.Counts[p]++
					ps.Values[p] = append(ps.Values[p], s)
				}
			}

		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		}
	}

	sort.Strings(ps.Namespaces)
	return ps, nil
}
