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

// RedactionPlaceholder replaces every value a redaction spec names.
//
// It is deliberately not empty and not valid base64 or PEM. An empty element is
// indistinguishable from a server that sent nothing, and a placeholder that
// still parsed as key material would let a reviewer's eye slide over it. This
// string is also what to grep for when auditing the corpus.
const RedactionPlaceholder = "[REDACTED-BY-GO-GMP-CAPTURE]"

// secretBearingCommands is the deny-by-default list: every command whose
// request or response can carry material that must never reach git.
//
// Membership is a claim about the protocol, not about the current specs. A
// command listed here with no entry in redactionSpecs makes the capture tool
// refuse to write its fixture, because "no spec" must mean "nobody has looked
// yet", never "nothing to redact". redactSpecsAreConsistent in redact_test.go
// pins the two maps to each other.
//
// Derived by reading commands/ rather than by name matching:
//
//   - authenticate carries the caller's own password on the way in.
//   - create_credential and modify_credential carry <password>, <community>,
//     <key><phrase>, <key><private> and <privacy><password>. The passphrase is
//     included even though it is not in the usual list: CreateCredentialKey has
//     a Phrase field and a passphrase is a secret by definition.
//   - create_user carries <password>; modify_user carries it as a
//     <password modify="1"> element with chardata.
//   - modify_auth and describe_auth carry <auth_conf_setting> values, which is
//     where an LDAP bind password or a RADIUS shared secret lives.
//   - get_credentials returns exportable credential material: <package> (the
//     downloadable key/installer blob), <certificate> and <public_key>.
//   - get_tls_certificates returns the certificate itself.
//   - get_scanners returns <ca_pub> and <key_pub>; create_scanner and
//     modify_scanner send <ca_pub>.
//   - create_tls_certificate sends a certificate blob.
//
// certificate_info is deliberately absent: gvmd derives it (issuer, subject,
// expiry) from the certificate rather than echoing key material, so redacting
// it would cost fixture value for no gain. If that ever stops being true, the
// literal-value sweep is the backstop.
var secretBearingCommands = map[string]bool{
	"authenticate":           true,
	"create_credential":      true,
	"modify_credential":      true,
	"create_user":            true,
	"modify_user":            true,
	"modify_auth":            true,
	"describe_auth":          true,
	"get_credentials":        true,
	"get_tls_certificates":   true,
	"create_tls_certificate": true,
	"get_scanners":           true,
	"create_scanner":         true,
	"modify_scanner":         true,
}

// redactionSpecs maps a GMP command to the canonical paths whose values are
// rewritten before its fixture is written.
//
// Paths are exactly what ExtractPaths produces, so a spec can be checked
// against a real document by eye. One spec covers both directions of a command:
// request paths are rooted at <command> and response paths at
// <command_response>, so they cannot collide.
//
// A path that matches nothing in a given document is not an error. That is what
// makes it safe to list both the element form and the attribute form of the
// same value when it is not yet known which one gvmd actually emits — see
// describe_auth, where commands/describe_auth.go models auth_conf_setting's key
// and value as child elements and gvmd is believed to send attributes.
var redactionSpecs = map[string][]string{
	"authenticate": {
		"/authenticate/credentials/password/#text",
	},
	"create_credential": {
		"/create_credential/password/#text",
		"/create_credential/community/#text",
		"/create_credential/key/phrase/#text",
		"/create_credential/key/private/#text",
		"/create_credential/privacy/password/#text",
	},
	"modify_credential": {
		"/modify_credential/password/#text",
		"/modify_credential/community/#text",
		"/modify_credential/key/phrase/#text",
		"/modify_credential/key/private/#text",
		"/modify_credential/privacy/password/#text",
	},
	"create_user": {
		"/create_user/password/#text",
	},
	"modify_user": {
		"/modify_user/password/#text",
	},
	"modify_auth": {
		"/modify_auth/group/auth_conf_setting/value/#text",
		"/modify_auth/group/auth_conf_setting/@value",
	},
	"describe_auth": {
		"/describe_auth_response/group/auth_conf_setting/value/#text",
		"/describe_auth_response/group/auth_conf_setting/@value",
	},
	"get_credentials": {
		"/get_credentials_response/credential/package/#text",
		"/get_credentials_response/credential/certificate/#text",
		"/get_credentials_response/credential/public_key/#text",
	},
	"get_tls_certificates": {
		"/get_tls_certificates_response/tls_certificate/certificate/#text",
	},
	"create_tls_certificate": {
		"/create_tls_certificate/certificate/#text",
	},
	"get_scanners": {
		"/get_scanners_response/scanner/ca_pub/#text",
		"/get_scanners_response/scanner/key_pub/#text",
	},
	"create_scanner": {
		"/create_scanner/ca_pub/#text",
	},
	"modify_scanner": {
		"/modify_scanner/ca_pub/#text",
	},
}

// SweepTerm is one literal value the final pass looks for in the bytes about to
// be written.
//
// Fatal separates a secret from an identifier. A password is fatal: it has no
// legitimate reason to appear, so a hit means a spec missed a path and the run
// must stop. A username is not: gvmd echoes it in <owner><name> on nearly every
// authenticated response, so treating it as fatal would abort every capture on
// its first fixture and there would be no corpus at all. It is still swept, and
// every hit is reported, so a reviewer can see where the identity leaked.
type SweepTerm struct {
	Label string
	Value string
	Fatal bool
}

// Redactor rewrites secret material out of captured fixtures, refuses to write
// a secret-bearing command that nobody has written a spec for, and sweeps the
// finished bytes for literal secrets the specs might have missed.
type Redactor struct {
	specs  map[string][]string
	secret map[string]bool
	terms  []SweepTerm
}

// NewRedactor returns a Redactor backed by the committed specs.
func NewRedactor() *Redactor {
	return newRedactor(redactionSpecs, secretBearingCommands)
}

// newRedactor exists so the tests can drive the deny-by-default rule with maps
// they control. Testing it through NewRedactor is impossible by construction:
// redact_test.go asserts that every secret-bearing command has a spec, so the
// refusal path is unreachable with the real maps, which is the point.
func newRedactor(specs map[string][]string, secret map[string]bool) *Redactor {
	return &Redactor{specs: specs, secret: secret}
}

// AddSecret registers a value whose appearance in any fixture is fatal.
func (r *Redactor) AddSecret(label, value string) {
	r.addTerm(label, value, true)
}

// AddIdentifier registers a value that is swept and reported but not fatal.
func (r *Redactor) AddIdentifier(label, value string) {
	r.addTerm(label, value, false)
}

func (r *Redactor) addTerm(label, value string, fatal bool) {
	if value == "" {
		return
	}
	for _, t := range r.terms {
		if t.Label == label && t.Value == value {
			return
		}
	}
	r.terms = append(r.terms, SweepTerm{Label: label, Value: value, Fatal: fatal})
}

// PathsFor returns the redaction spec for a command, or an error if the command
// is known to carry secrets and nobody has written one.
func (r *Redactor) PathsFor(command string) ([]string, error) {
	if paths, ok := r.specs[command]; ok {
		return paths, nil
	}
	if r.secret[command] {
		return nil, fmt.Errorf(
			"refusing to handle a %q fixture: the command is listed as secret-bearing but has "+
				"no redaction spec, and an absent spec must never be read as 'nothing to redact'; "+
				"add one to redactionSpecs in internal/gmptest/redact.go", command)
	}
	return nil, nil
}

// Redact rewrites doc according to command's spec, returning the new bytes and
// the paths that actually matched. It returns doc unchanged, and no applied
// paths, when the command has no spec or the spec matched nothing — so a caller
// can use len(applied) == 0 as "these are still the wire bytes".
func (r *Redactor) Redact(command string, doc []byte) ([]byte, []string, error) {
	paths, err := r.PathsFor(command)
	if err != nil {
		return nil, nil, err
	}
	return RedactPaths(doc, paths)
}

// Sweep looks for every registered term in doc.
//
// It reports labels, never values: an error message that quoted the password it
// found would put it straight into the test log and from there into a CI
// transcript, which is the same disclosure the sweep exists to prevent.
func (r *Redactor) Sweep(doc []byte) (fatal, notices []string) {
	for _, t := range r.terms {
		if !bytes.Contains(doc, []byte(t.Value)) {
			continue
		}
		if t.Fatal {
			fatal = append(fatal, t.Label)
		} else {
			notices = append(notices, t.Label)
		}
	}
	sort.Strings(fatal)
	sort.Strings(notices)
	return fatal, notices
}

// Terms returns the registered sweep terms, for tests and for reporting what
// the sweep actually covered.
func (r *Redactor) Terms() []SweepTerm {
	out := make([]SweepTerm, len(r.terms))
	copy(out, r.terms)
	return out
}

// span is a half-open byte range of doc to be replaced.
type span struct{ start, end int }

// RedactPaths rewrites the value at each canonical path in doc and leaves every
// other byte exactly as it was.
//
// Splicing bytes rather than decoding and re-encoding is the whole point. A
// round trip through xml.Encoder would rewrite <a/> as <a></a>, normalize
// attribute quoting and expand numeric character references across the entire
// document, so every live fixture would stop being wire bytes the moment one
// element in it needed redacting — and fixtures that needed no redaction would
// have to take a different code path to stay faithful, which is exactly the
// kind of asymmetry that hides bugs.
//
// The byte spans come from Decoder.InputOffset(): the offset after the previous
// token is the offset of this one, because the decoder yields a token for every
// byte range including the whitespace between elements.
func RedactPaths(doc []byte, paths []string) ([]byte, []string, error) {
	if len(paths) == 0 {
		return doc, nil, nil
	}
	want := make(map[string]bool, len(paths))
	for _, p := range paths {
		if err := ValidateRedactionPath(p); err != nil {
			return nil, nil, err
		}
		want[p] = true
	}

	spans, hit, err := findRedactionSpans(doc, want)
	if err != nil {
		return nil, nil, err
	}
	if len(spans) == 0 {
		return doc, nil, nil
	}

	applied := make([]string, 0, len(hit))
	for p := range hit {
		applied = append(applied, p)
	}
	sort.Strings(applied)

	return spliceSpans(doc, spans), applied, nil
}

// findRedactionSpans walks doc once and collects the byte ranges to replace.
func findRedactionSpans(doc []byte, want map[string]bool) ([]span, map[string]bool, error) {
	dec := xml.NewDecoder(bytes.NewReader(doc))
	var (
		stack []string
		spans []span
	)
	hit := map[string]bool{}
	prev := 0

	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("redact: document is not well-formed: %w", err)
		}
		cur := int(dec.InputOffset())

		switch t := tok.(type) {
		case xml.StartElement:
			stack = append(stack, qualify(t.Name))
			found, err := attrSpans(doc[prev:cur], prev, pathOf(stack), t.Attr, want, hit)
			if err != nil {
				return nil, nil, err
			}
			spans = append(spans, found...)

		case xml.CharData:
			// Whitespace-only text is formatting, matching ExtractPaths, which
			// is what the specs are written against. A spec path is therefore
			// exactly a path from a fidelity report, with no separate rules.
			if len(stack) > 0 && strings.TrimSpace(string(t)) != "" {
				p := pathOf(stack) + "/#text"
				if want[p] {
					spans = append(spans, span{prev, cur})
					hit[p] = true
				}
			}

		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		}
		prev = cur
	}
	return spans, hit, nil
}

// attrSpans locates the value of every wanted attribute inside one start tag.
//
// raw is the tag's bytes and base is where they start in the document. A wanted
// attribute that cannot be located is a hard error rather than a skip: the only
// way to reach it is a bug in the scanner, and silently declining to redact is
// the one outcome this package must never produce.
func attrSpans(raw []byte, base int, here string, attrs []xml.Attr, want, hit map[string]bool) ([]span, error) {
	var out []span
	for _, a := range attrs {
		if a.Name.Local == "xmlns" || a.Name.Space == "xmlns" {
			continue
		}
		p := here + "/@" + qualify(a.Name)
		if !want[p] {
			continue
		}
		start, end, ok := attrValueSpan(raw, a.Name)
		if !ok {
			return nil, fmt.Errorf("redact: attribute %s is present but its value could not be "+
				"located in the start tag; refusing to write a possibly unredacted fixture", p)
		}
		out = append(out, span{base + start, base + end})
		hit[p] = true
	}
	return out, nil
}

func pathOf(stack []string) string { return "/" + strings.Join(stack, "/") }

// spliceSpans replaces each span with the placeholder, merging spans that abut.
//
// Merging matters because one logical text node can arrive as several CharData
// tokens: a CDATA section is its own token, so <p><![CDATA[a]]>b</p> yields two
// adjacent spans that would otherwise produce two placeholders in a row.
func spliceSpans(doc []byte, spans []span) []byte {
	sort.Slice(spans, func(i, j int) bool { return spans[i].start < spans[j].start })

	var out bytes.Buffer
	last := 0
	for i := 0; i < len(spans); i++ {
		s := spans[i]
		for i+1 < len(spans) && spans[i+1].start <= s.end {
			i++
			if spans[i].end > s.end {
				s.end = spans[i].end
			}
		}
		out.Write(doc[last:s.start])
		out.WriteString(RedactionPlaceholder)
		last = s.end
	}
	out.Write(doc[last:])
	return out.Bytes()
}

// ValidateRedactionPath rejects a spec path that cannot possibly match.
//
// A typo in a spec is the failure mode with no symptom: the path matches
// nothing, the fixture is written, and the secret is in it. Every spec path
// goes through this before any document is touched.
func ValidateRedactionPath(p string) error {
	switch {
	case !strings.HasPrefix(p, "/"):
		return fmt.Errorf("redaction path %q: must start with /", p)
	case strings.HasSuffix(p, "/"):
		return fmt.Errorf("redaction path %q: must not end with /", p)
	case strings.Contains(p, "//"):
		return fmt.Errorf("redaction path %q: contains an empty segment", p)
	case strings.ContainsAny(p, " \t\n"):
		return fmt.Errorf("redaction path %q: contains whitespace", p)
	}
	last := p[strings.LastIndex(p, "/")+1:]
	if last != "#text" && !strings.HasPrefix(last, "@") {
		return fmt.Errorf("redaction path %q: must end in /#text or /@attr; a path naming an "+
			"element rewrites nothing, because an element has no value of its own", p)
	}
	if strings.HasPrefix(last, "@") && len(last) == 1 {
		return fmt.Errorf("redaction path %q: names an attribute with no name", p)
	}
	return nil
}

// attrValueSpan finds the byte range of one attribute's value inside a start
// tag, returning offsets relative to raw.
//
// A hand-rolled scan rather than a regexp because the value may contain
// anything except its own quote character, and because the decoder has already
// proved the tag is well-formed, so the scan only has to agree with a parse
// that is known to have succeeded.
func attrValueSpan(raw []byte, want xml.Name) (int, int, bool) {
	s := &attrScanner{raw: raw}
	s.skipElementName()
	for {
		name, start, end, ok := s.next()
		if !ok {
			return 0, 0, false
		}
		if attrNameMatches(name, want) {
			return start, end, true
		}
	}
}

// attrScanner walks the attributes of a single start tag.
type attrScanner struct {
	raw []byte
	i   int
}

func (s *attrScanner) skipElementName() {
	if s.i < len(s.raw) && s.raw[s.i] == '<' {
		s.i++
	}
	for s.i < len(s.raw) && !isXMLSpace(s.raw[s.i]) && s.raw[s.i] != '>' && s.raw[s.i] != '/' {
		s.i++
	}
}

func (s *attrScanner) skipSpace() {
	for s.i < len(s.raw) && isXMLSpace(s.raw[s.i]) {
		s.i++
	}
}

// next returns the next attribute's literal name and the byte range of its
// value, or ok=false once the tag is exhausted or stops making sense.
func (s *attrScanner) next() (name string, valueStart, valueEnd int, ok bool) {
	name, ok = s.readName()
	if !ok {
		return "", 0, 0, false
	}
	if !s.expectEquals() {
		return "", 0, 0, false
	}
	valueStart, valueEnd, ok = s.readQuotedValue()
	if !ok {
		return "", 0, 0, false
	}
	return name, valueStart, valueEnd, true
}

func (s *attrScanner) readName() (string, bool) {
	s.skipSpace()
	if s.i >= len(s.raw) || s.raw[s.i] == '>' || s.raw[s.i] == '/' {
		return "", false
	}
	start := s.i
	for s.i < len(s.raw) && s.raw[s.i] != '=' && !isXMLSpace(s.raw[s.i]) {
		s.i++
	}
	return string(s.raw[start:s.i]), true
}

func (s *attrScanner) expectEquals() bool {
	s.skipSpace()
	if s.i >= len(s.raw) || s.raw[s.i] != '=' {
		return false
	}
	s.i++
	return true
}

func (s *attrScanner) readQuotedValue() (int, int, bool) {
	s.skipSpace()
	if s.i >= len(s.raw) || (s.raw[s.i] != '"' && s.raw[s.i] != '\'') {
		return 0, 0, false
	}
	quote := s.raw[s.i]
	s.i++

	start := s.i
	for s.i < len(s.raw) && s.raw[s.i] != quote {
		s.i++
	}
	if s.i >= len(s.raw) {
		return 0, 0, false
	}
	end := s.i
	s.i++
	return start, end, true
}

// attrNameMatches compares the literal name as written in the document against
// the namespace-translated xml.Name the decoder handed back.
func attrNameMatches(literal string, want xml.Name) bool {
	if literal == want.Local {
		return want.Space == ""
	}
	prefix, local, ok := strings.Cut(literal, ":")
	return ok && prefix != "" && local == want.Local && want.Space != ""
}

func isXMLSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}
