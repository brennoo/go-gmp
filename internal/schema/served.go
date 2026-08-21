package schema

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/brennoo/go-gmp/internal/imageref"
)

// servedSource identifies a live gvmd's <help format="XML" type=""/>
// response as the origin of a batch of Commands, independent of which
// server instance answered it — analogous to upstreamSource identifying
// GMP.xml.in for file-based extraction. See docs/review-2026-08.md, T-16b.
const servedSource = `gvmd <help format="XML" type=""/> response`

// LoadServed extracts commands from a live gvmd's <help format="XML"
// type=""/> response (T-16b), rather than a vendored file (Load).
//
// data must be the full raw response envelope exactly as read off the wire,
// e.g. via gmp.Connection.RawXML — never a HelpResponse.Text decode, which
// only captures bare chardata and silently loses the nested <protocol>
// document entirely.
//
// sourceRevision identifies the server that served data — in practice its
// resolved container image digest (internal/dockerstack), since GMP exposes
// no build version of its own.
func LoadServed(data []byte, sourceRevision string) ([]Command, error) {
	if err := imageref.ValidateDigest(sourceRevision); err != nil {
		return nil, fmt.Errorf("sourceRevision: %w", err)
	}

	protocolXML, err := unwrapServedProtocol(data)
	if err != nil {
		return nil, err
	}

	// Served extraction has no per-command file to name here, unlike Load's
	// "command_get_tasks.xml"; leaving it empty keeps a served/file parity
	// diff confined to Provenance instead of a manufactured, always-different
	// placeholder.
	cmds, err := Parse(protocolXML, "")
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		return cmds, nil
	}

	var ver rawProtocolVersion
	if err := xml.Unmarshal(protocolXML, &ver); err != nil {
		return nil, fmt.Errorf("reading the served <protocol> version: %w", err)
	}

	// There is no split step for served extraction: the envelope unwrap is a
	// lossless projection, so unlike Load, "what was parsed" and "what was
	// served" are identical by construction and share one hash.
	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])
	prov := Provenance{
		Method:          ExtractionServed,
		Source:          servedSource,
		SourceRevision:  sourceRevision,
		DeclaredVersion: trim(ver.Version),
		SourceSHA256:    hash,
		ParsedSHA256:    hash,
		ExtractedAt:     time.Now().UTC(),
	}
	for i := range cmds {
		cmds[i].Provenance = prov
	}
	sort.Slice(cmds, func(i, j int) bool { return cmds[i].Name < cmds[j].Name })
	return cmds, nil
}

// unwrapServedProtocol returns the raw bytes of the top-level <protocol>
// element inside a <help_response><schema><protocol> envelope.
//
// It matches by exact ancestor path rather than by tag name anywhere in the
// stream: the envelope's own <command name="help"> entry documents itself
// with a truncated example that embeds a second, decoy
// <help_response>/<schema>/<protocol> nested inside
// <command>/<example>/<response>, which a plain name search could not tell
// apart from the real one.
func unwrapServedProtocol(data []byte) ([]byte, error) {
	dec := xml.NewDecoder(bytes.NewReader(data))

	var stack []string
	for {
		offset := dec.InputOffset()
		tok, err := dec.Token()
		if err != nil {
			return nil, fmt.Errorf("scanning for the served <protocol> document: %w", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			if len(stack) == 0 {
				if t.Name.Local != "help_response" {
					return nil, fmt.Errorf("unexpected root element <%s>, want <help_response>", t.Name.Local)
				}
				if status := attrValue(t.Attr, "status"); status != "" && !strings.HasPrefix(status, "2") {
					return nil, fmt.Errorf("help_response status %s, want a 2xx", status)
				}
			}
			if len(stack) == 2 && stack[0] == "help_response" && stack[1] == "schema" && t.Name.Local == "protocol" {
				if err := dec.Skip(); err != nil {
					return nil, fmt.Errorf("reading the served <protocol> document: %w", err)
				}
				return data[offset:dec.InputOffset()], nil
			}
			stack = append(stack, t.Name.Local)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		}
	}
}

func attrValue(attrs []xml.Attr, name string) string {
	for _, a := range attrs {
		if a.Name.Local == name {
			return a.Value
		}
	}
	return ""
}
