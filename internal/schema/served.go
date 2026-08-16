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
// no build version of its own — supplied explicitly, the same way Load's
// sourceRevision is a git revision rather than something inferred.
func LoadServed(data []byte, sourceRevision string) ([]Command, error) {
	if err := imageref.ValidateDigest(sourceRevision); err != nil {
		return nil, fmt.Errorf("sourceRevision: %w", err)
	}

	protocolXML, err := unwrapServedProtocol(data)
	if err != nil {
		return nil, err
	}

	// Parse's name parameter is the per-command file a report would cite
	// (e.g. "command_get_tasks.xml") — file-based extraction has one real
	// file per command, served extraction has none, so there is nothing
	// honest to put here. Leaving it empty also keeps a served/file parity
	// diff (T-16b) confined to Provenance, as the plan specifies, rather
	// than manufacturing a same-for-every-command "file name" that would
	// differ from the file-based side on every single command for no
	// schema-relevant reason.
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

	// There is no split step for served extraction — what Parse read is
	// exactly what the server sent, modulo the envelope unwrap, which is a
	// lossless projection rather than an edit. So, unlike Load, where
	// ParsedSHA256 exists to catch a split drifting from its source
	// document, "what was parsed" and "what was served" are identical by
	// construction and share one hash.
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
// It locates that element by its exact ancestor path — root <help_response>,
// its child <schema>, its child <protocol> — rather than by searching for a
// tag name anywhere in the byte stream. That distinction matters here: the
// envelope's own <command name="help"> entry documents itself with a
// truncated example that embeds a second, decoy
// <help_response>/<schema>/<protocol> nested many levels inside
// <command>/<example>/<response>. A name search could not tell the two
// apart in general; matching by ancestry finds the real one regardless of
// where either happens to fall in the byte stream.
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
