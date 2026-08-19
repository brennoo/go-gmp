package schema

import (
	"encoding/xml"
	"strings"
	"testing"
)

// The differ gets the same treatment as the parser: hand-written pairs small
// enough that a failure names the bug. A tool that reports 786 findings has to
// be trustworthy at this size first.

func diff(t *testing.T, body string, req, resp any) Result {
	t.Helper()
	cmd := parseOne(t, body)
	c := &Comparer{}
	return c.Compare(cmd, ModelOf(req), ModelOf(resp))
}

// only asserts the finding set for one category, ignoring suppressed ones.
func only(r Result, cat string) []Finding {
	var out []Finding
	for _, f := range r.Findings {
		if f.Category == cat && !f.Suppressed {
			out = append(out, f)
		}
	}
	return out
}

func wantOne(t *testing.T, r Result, cat, path string) Finding {
	t.Helper()
	fs := only(r, cat)
	if len(fs) != 1 {
		t.Fatalf("Got %d %s findings, want 1: %v", len(fs), cat, render(fs))
	}
	if fs[0].Path != path {
		t.Fatalf("Got %s at %q, want %q", cat, fs[0].Path, path)
	}
	return fs[0]
}

func render(fs []Finding) []string {
	out := make([]string, 0, len(fs))
	for _, f := range fs {
		out = append(out, f.Category+" "+f.Path+" "+f.Detail)
	}
	return out
}

const twoFields = `
    <name>x</name>
    <pattern><e>a</e><o><e>b</e></o></pattern>
    <ele><name>a</name><pattern>text</pattern></ele>
    <ele><name>b</name><pattern>text</pattern></ele>
    <response><pattern><attrib><name>status</name><type>status</type><required>1</required></attrib></pattern></response>`

type respOK struct {
	XMLName xml.Name `xml:"x_response"`
	Status  string   `xml:"status,attr"`
}

func TestMissingField(t *testing.T) {
	type req struct {
		XMLName xml.Name `xml:"x"`
		A       string   `xml:"a"`
	}
	r := diff(t, twoFields, &req{}, &respOK{})
	f := wantOne(t, r, CatMissing, "/x/b")
	if !strings.Contains(f.Detail, "optional elem b") {
		t.Errorf("Detail should name the cardinality and kind, got %q", f.Detail)
	}
}

func TestExtraField(t *testing.T) {
	type req struct {
		XMLName xml.Name `xml:"x"`
		A       string   `xml:"a"`
		B       string   `xml:"b,omitempty"`
		C       string   `xml:"c,omitempty"`
	}
	r := diff(t, twoFields, &req{}, &respOK{})
	wantOne(t, r, CatExtra, "/x/c")
}

// internal/codegen's FieldOverrides mechanism adds a real, evidence-backed
// field the schema's own pattern never declares (e.g. tls_certificate's
// <valid>), recording the evidence in a gmpoverride struct tag. The
// resulting EXTRA finding must still fire — the schema still declares no
// such field — but come back suppressed, with Reason carrying the evidence.
func TestExtraFieldWithOverrideReasonIsSuppressed(t *testing.T) {
	type req struct {
		XMLName xml.Name `xml:"x"`
		A       string   `xml:"a"`
		B       string   `xml:"b,omitempty"`
		C       string   `xml:"c,omitempty" gmpoverride:"live evidence for c"`
	}
	r := diff(t, twoFields, &req{}, &respOK{})
	if got := only(r, CatExtra); len(got) != 0 {
		t.Errorf("A suppressed finding must not show up in the actionable count, got %v", render(got))
	}
	var found bool
	for _, f := range r.Findings {
		if f.Path == "/x/c" && f.Category == CatExtra && f.Suppressed && f.Reason == "live evidence for c" {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected a suppressed CatExtra finding at /x/c with the override's reason, got %v", render(r.Findings))
	}
}

// The chardata counterpart to TestExtraFieldWithOverrideReasonIsSuppressed:
// FieldOverrides can also add a ",chardata" field the schema declares no
// text token for (e.g. tls_certificate's <certificate> bytes), via
// checkChardata's ChardataOverrideReason path rather than walk()'s
// OverrideReason path.
func TestChardataWithOverrideReasonIsSuppressed(t *testing.T) {
	body := `
    <name>x</name>
    <pattern><e>a</e></pattern>
    <ele><name>a</name><pattern><attrib><name>fmt</name><type>text</type></attrib></pattern></ele>
    <response><pattern><attrib><name>status</name><type>status</type><required>1</required></attrib></pattern></response>`
	type req struct {
		XMLName xml.Name `xml:"x"`
		A       struct {
			Value string `xml:",chardata" gmpoverride:"live evidence for a's text"`
			Fmt   string `xml:"fmt,attr"`
		} `xml:"a"`
	}
	r := diff(t, body, &req{}, &respOK{})
	if got := only(r, CatExtra); len(got) != 0 {
		t.Errorf("A suppressed finding must not show up in the actionable count, got %v", render(got))
	}
	var found bool
	for _, f := range r.Findings {
		if f.Path == "/x/a/#text" && f.Category == CatExtra && f.Suppressed && f.Reason == "live evidence for a's text" {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected a suppressed CatExtra finding at /x/a/#text with the override's reason, got %v", render(r.Findings))
	}
}

// T-18a's shape: the schema wrongly declares an element's chardata on the
// wrapper (e.g. <sort>) when the wire actually carries it one level deeper,
// on a child (e.g. <sort><field>). internal/codegen's FieldSuppressions
// relocates it by suppressing the wrapper's chardata and adding an override
// one on the child; the wrapper field itself carries the relocation's
// gmpoverride tag, since the wrapper struct has no chardata field left of
// its own to carry it. checkChardata must suppress both halves: the
// wrapper's now-missing chardata (via g.OverrideReason) and the child's now-
// extra chardata (via g.ChardataOverrideReason), the same way
// TestChardataWithOverrideReasonIsSuppressed does for a lone addition.
func TestChardataRelocationSuppressesBothHalves(t *testing.T) {
	body := `
    <name>x</name>
    <pattern><e>a</e></pattern>
    <ele>
      <name>a</name>
      <pattern>
        <t>text</t>
        <e>field</e>
      </pattern>
      <ele>
        <name>field</name>
        <pattern><e>order</e></pattern>
        <ele><name>order</name><pattern><t>text</t></pattern></ele>
      </ele>
    </ele>
    <response><pattern><attrib><name>status</name><type>status</type><required>1</required></attrib></pattern></response>`
	type req struct {
		XMLName xml.Name `xml:"x"`
		A       struct {
			Field struct {
				Value string `xml:",chardata" gmpoverride:"relocated from a"`
				Order string `xml:"order"`
			} `xml:"field"`
		} `xml:"a" gmpoverride:"relocated to a/field"`
	}
	r := diff(t, body, &req{}, &respOK{})
	if got := only(r, CatMissing); len(got) != 0 {
		t.Errorf("A suppressed MISSING must not show up in the actionable count, got %v", render(got))
	}
	if got := only(r, CatExtra); len(got) != 0 {
		t.Errorf("A suppressed EXTRA must not show up in the actionable count, got %v", render(got))
	}

	wantSuppressed := map[string]struct{ cat, reason string }{
		"/x/a/#text":       {CatMissing, "relocated to a/field"},
		"/x/a/field/#text": {CatExtra, "relocated from a"},
	}
	for _, f := range r.Findings {
		want, ok := wantSuppressed[f.Path]
		if !ok {
			continue
		}
		delete(wantSuppressed, f.Path)
		if f.Category != want.cat || !f.Suppressed || f.Reason != want.reason {
			t.Errorf("At %s: got category %s suppressed=%v reason=%q, want %s suppressed=true reason=%q",
				f.Path, f.Category, f.Suppressed, f.Reason, want.cat, want.reason)
		}
	}
	for path := range wantSuppressed {
		t.Errorf("Expected a suppressed finding at %s, got none: %v", path, render(r.Findings))
	}
}

// The counterpart to TestChardataRelocationSuppressesBothHalves: a relocation
// missing its source-side gmpoverride tag is not a validated relocation at
// all, just a struct that lost its chardata field, and must fail closed as an
// ordinary actionable MISSING rather than being silently accepted.
func TestUnpairedChardataMoveStaysActionable(t *testing.T) {
	body := `
    <name>x</name>
    <pattern><e>a</e></pattern>
    <ele>
      <name>a</name>
      <pattern>
        <t>text</t>
        <e>field</e>
      </pattern>
      <ele>
        <name>field</name>
        <pattern><e>order</e></pattern>
        <ele><name>order</name><pattern><t>text</t></pattern></ele>
      </ele>
    </ele>
    <response><pattern/></response>`
	type req struct {
		XMLName xml.Name `xml:"x"`
		A       struct {
			Field struct {
				Value string `xml:",chardata" gmpoverride:"relocated from a"`
				Order string `xml:"order"`
			} `xml:"field"`
		} `xml:"a"` // no gmpoverride: the source half of the move is undocumented
	}
	r := diff(t, body, &req{}, &respOK{})
	wantOne(t, r, CatMissing, "/x/a/#text")
}

func TestConformantPairHasNoFindings(t *testing.T) {
	type req struct {
		XMLName xml.Name `xml:"x"`
		A       string   `xml:"a"`
		B       string   `xml:"b,omitempty"`
	}
	r := diff(t, twoFields, &req{}, &respOK{})
	if !r.Conformant() {
		t.Errorf("Expected no findings, got %v", render(r.Findings))
	}
}

func TestCardinalityRepeatedNeedsSlice(t *testing.T) {
	body := `
    <name>x</name>
    <pattern><any><e>a</e></any></pattern>
    <ele><name>a</name><pattern>text</pattern></ele>
    <response><pattern><attrib><name>status</name><type>status</type><required>1</required></attrib></pattern></response>`
	type req struct {
		XMLName xml.Name `xml:"x"`
		A       string   `xml:"a,omitempty"`
	}
	f := wantOne(t, diff(t, body, &req{}, &respOK{}), CatCardinality, "/x/a")
	if !strings.Contains(f.Detail, "zero-or-more") {
		t.Errorf("Got %q", f.Detail)
	}
}

func TestCardinalitySliceForSingle(t *testing.T) {
	type req struct {
		XMLName xml.Name `xml:"x"`
		A       []string `xml:"a"`
		B       string   `xml:"b,omitempty"`
	}
	wantOne(t, diff(t, twoFields, &req{}, &respOK{}), CatCardinality, "/x/a")
}

func TestCardinalityOptionalScalarCannotExpressAbsence(t *testing.T) {
	body := `
    <name>x</name>
    <pattern><o><e>n</e></o><o><e>s</e></o></pattern>
    <ele><name>n</name><pattern><t>integer</t></pattern></ele>
    <ele><name>s</name><pattern>text</pattern></ele>
    <response><pattern><attrib><name>status</name><type>status</type><required>1</required></attrib></pattern></response>`
	type req struct {
		XMLName xml.Name `xml:"x"`
		N       int      `xml:"n,omitempty"`
		S       string   `xml:"s,omitempty"`
	}
	// An optional string is fine: "" and absent mean the same on this wire.
	// An optional int is not: 0 is a value.
	wantOne(t, diff(t, body, &req{}, &respOK{}), CatCardinality, "/x/n")
}

func TestKindAttributeModelledAsElement(t *testing.T) {
	body := `
    <name>x</name>
    <pattern><attrib><name>a</name><type>text</type></attrib></pattern>
    <response><pattern><attrib><name>status</name><type>status</type><required>1</required></attrib></pattern></response>`
	type req struct {
		XMLName xml.Name `xml:"x"`
		A       string   `xml:"a,omitempty"`
	}
	f := wantOne(t, diff(t, body, &req{}, &respOK{}), CatKind, "/x/@a")
	if !strings.Contains(f.Detail, "schema declares attr") {
		t.Errorf("Got %q", f.Detail)
	}
	// It must not also be reported as MISSING and EXTRA.
	if got := only(diff(t, body, &req{}, &respOK{}), CatMissing); len(got) != 0 {
		t.Errorf("A kind mismatch must be one finding, also got %v", render(got))
	}
}

func TestOmitEmptyOnlyActionableOnRequests(t *testing.T) {
	body := `
    <name>x</name>
    <pattern><o><e>a</e></o></pattern>
    <ele><name>a</name><pattern>text</pattern></ele>
    <response>
      <pattern>
        <attrib><name>status</name><type>status</type><required>1</required></attrib>
        <o><e>b</e></o>
      </pattern>
      <ele><name>b</name><pattern>text</pattern></ele>
    </response>`
	type req struct {
		XMLName xml.Name `xml:"x"`
		A       string   `xml:"a"` // optional in the schema, always emitted
	}
	type resp struct {
		XMLName xml.Name `xml:"x_response"`
		Status  string   `xml:"status,attr"`
		B       string   `xml:"b"` // same shape, but decoding only
	}
	r := diff(t, body, &req{}, &resp{})
	wantOne(t, r, CatOmitEmpty, "/x/a")

	var suppressed int
	for _, f := range r.Findings {
		if f.Category == CatOmitEmpty && f.Suppressed {
			suppressed++
			if f.Path != "/x_response/b" {
				t.Errorf("Unexpected suppressed finding at %s", f.Path)
			}
		}
	}
	if suppressed != 1 {
		t.Errorf("Got %d suppressed omitempty findings, want 1", suppressed)
	}
}

func TestRequiredWithOmitEmptyIsSuppressedOnRequests(t *testing.T) {
	type req struct {
		XMLName xml.Name `xml:"x"`
		A       string   `xml:"a,omitempty"` // schema says required
		B       string   `xml:"b,omitempty"`
	}
	r := diff(t, twoFields, &req{}, &respOK{})
	if !r.Conformant() {
		t.Errorf("Request-side over-constraint must not be charged to commands/: %v", render(r.Findings))
	}
	var found bool
	for _, f := range r.Findings {
		if f.Path == "/x/a" && f.Suppressed && f.Reason == "request-side over-constraint" {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected the finding to be recorded and suppressed, got %v", render(r.Findings))
	}
}

func TestBooleanIsNotBool(t *testing.T) {
	body := `
    <name>x</name>
    <pattern><o><e>flag</e></o></pattern>
    <ele><name>flag</name><pattern>boolean</pattern></ele>
    <response><pattern><attrib><name>status</name><type>status</type><required>1</required></attrib></pattern></response>`
	type req struct {
		XMLName xml.Name `xml:"x"`
		Flag    bool     `xml:"flag,omitempty"`
	}
	f := wantOne(t, diff(t, body, &req{}, &respOK{}), CatType, "/x/flag")
	if !strings.Contains(f.Detail, "boolean") || !strings.Contains(f.Detail, "bool") {
		t.Errorf("Got %q", f.Detail)
	}
}

func TestBooleanAsStringIsNotFlagged(t *testing.T) {
	body := `
    <name>x</name>
    <pattern><o><e>flag</e></o></pattern>
    <ele><name>flag</name><pattern>boolean</pattern></ele>
    <response><pattern><attrib><name>status</name><type>status</type><required>1</required></attrib></pattern></response>`
	type req struct {
		XMLName xml.Name `xml:"x"`
		Flag    string   `xml:"flag,omitempty"`
	}
	if got := only(diff(t, body, &req{}, &respOK{}), CatType); len(got) != 0 {
		t.Errorf("A string holds the schema's lexical space exactly, got %v", render(got))
	}
}

func TestSchemaDefectsAreNotChargedToGo(t *testing.T) {
	body := `
    <name>x</name>
    <pattern><e>ghost</e><e>hollow</e></pattern>
    <ele><name>hollow</name><pattern/></ele>
    <response><pattern><attrib><name>status</name><type>status</type><required>1</required></attrib></pattern></response>`
	type req struct {
		XMLName xml.Name `xml:"x"`
	}
	r := diff(t, body, &req{}, &respOK{})
	if !r.Conformant() {
		t.Errorf("Nothing here asks anything of commands/, got %v", render(r.Findings))
	}
	if got := only(r, CatDefect); len(got) != 2 {
		t.Errorf("Got %d defects, want 2: %v", len(got), render(got))
	}
}

func TestNestedTagPathBuildsAWrapper(t *testing.T) {
	body := `
    <name>x</name>
    <pattern><e>agents</e></pattern>
    <ele>
      <name>agents</name>
      <pattern><any><e>agent</e></any></pattern>
      <ele><name>agent</name><pattern>text</pattern></ele>
    </ele>
    <response><pattern><attrib><name>status</name><type>status</type><required>1</required></attrib></pattern></response>`
	type req struct {
		XMLName xml.Name `xml:"x"`
		Agents  []string `xml:"agents>agent"`
	}
	r := diff(t, body, &req{}, &respOK{})
	if !r.Conformant() {
		t.Errorf("`agents>agent` models exactly this shape, got %v", render(r.Findings))
	}
}

// get_filters's two same-wire-name <filters> definitions, disambiguated by
// <id> (see TestIDDisambiguatesSharedWireName): both nodes are named "x", so
// walk's occurrence pairing rather than a plain name lookup is what tells
// them apart.
const duplicateWireName = `
    <name>x</name>
    <pattern><e>x</e><e>x2</e></pattern>
    <ele>
      <name>x</name>
      <pattern><e>a</e></pattern>
      <ele><name>a</name><pattern>text</pattern></ele>
    </ele>
    <ele>
      <id>x2</id>
      <name>x</name>
      <pattern><e>b</e></pattern>
      <ele><name>b</name><pattern>text</pattern></ele>
    </ele>
    <response><pattern><attrib><name>status</name><type>status</type><required>1</required></attrib></pattern></response>`

func TestDuplicateWireNamePairsByOccurrence(t *testing.T) {
	type first struct {
		A string `xml:"a"`
	}
	type second struct {
		B string `xml:"b"`
	}
	type req struct {
		XMLName xml.Name `xml:"x"`
		First   first    `xml:"x"`
		Second  second   `xml:"-" gmpwire:"x"`
	}
	r := diff(t, duplicateWireName, &req{}, &respOK{})
	if !r.Conformant() {
		t.Errorf("Declaration order matches the schema's, got %v", render(r.Findings))
	}
}

// If the comparer merely checked that "a" and "b" appear somewhere among the
// two "x" occurrences, swapping which Go field holds which shape would still
// look conformant. It doesn't: each occurrence is pairwise wrong.
func TestDuplicateWireNamePairingIsPositional(t *testing.T) {
	type first struct {
		A string `xml:"a"`
	}
	type second struct {
		B string `xml:"b"`
	}
	type req struct {
		XMLName xml.Name `xml:"x"`
		First   second   `xml:"x"`
		Second  first    `xml:"-" gmpwire:"x"`
	}
	r := diff(t, duplicateWireName, &req{}, &respOK{})
	if got := only(r, CatMissing); len(got) != 2 {
		t.Errorf("Got %d MISSING, want 2 (one per swapped occurrence): %v", len(got), render(got))
	}
	if got := only(r, CatExtra); len(got) != 2 {
		t.Errorf("Got %d EXTRA, want 2 (one per swapped occurrence): %v", len(got), render(got))
	}
}

// A third Go field sharing the wire name has no schema occurrence left to
// pair with once the first two are claimed, and must be reported EXTRA
// rather than waved through because the name exists in the schema.
func TestDuplicateWireNameExcessGoOccurrenceIsExtra(t *testing.T) {
	type first struct {
		A string `xml:"a"`
	}
	type second struct {
		B string `xml:"b"`
	}
	type third struct {
		C string `xml:"c"`
	}
	type req struct {
		XMLName xml.Name `xml:"x"`
		First   first    `xml:"x"`
		Second  second   `xml:"-" gmpwire:"x"`
		Third   third    `xml:"-" gmpwire:"x"`
	}
	wantOne(t, diff(t, duplicateWireName, &req{}, &respOK{}), CatExtra, "/x/x")
}

func TestFidelityCrossReference(t *testing.T) {
	body := `
    <name>x</name>
    <pattern/>
    <response>
      <pattern>
        <attrib><name>status</name><type>status</type><required>1</required></attrib>
        <e>lost</e>
      </pattern>
      <ele><name>lost</name><pattern>text</pattern></ele>
    </response>`
	cmd := parseOne(t, body)
	c := &Comparer{Dropped: map[string]map[string]bool{"x": {"/x_response/lost": true}}}
	r := c.Compare(cmd, ModelOf(&struct {
		XMLName xml.Name `xml:"x"`
	}{}), ModelOf(&respOK{}))

	f := wantOne(t, r, CatMissing, "/x_response/lost")
	if !f.Confirmed {
		t.Errorf("Expected the fidelity corpus to confirm this path")
	}
}

// Compare's doc comment says either root may be zero and that is itself
// reported. ModelOf(nil) is how that happens in practice: a registry entry
// whose constructor does not reflect to a struct at all.
func TestMissingGoRootIsReported(t *testing.T) {
	r := diff(t, twoFields, nil, nil)
	cmd := parseOne(t, twoFields)

	findAt := func(path string) (Finding, bool) {
		for _, f := range only(r, CatMissing) {
			if f.Path == path {
				return f, true
			}
		}
		return Finding{}, false
	}

	req, ok := findAt("/" + cmd.Request.Name)
	if !ok || req.Side != "request" || !strings.Contains(req.Detail, "no Go type") {
		t.Errorf("got %v, want a request-side no-Go-type finding at /%s", render(r.Findings), cmd.Request.Name)
	}

	resp, ok := findAt("/" + cmd.Response.Name)
	if !ok || resp.Side != "response" || !strings.Contains(resp.Detail, "no Go type") {
		t.Errorf("got %v, want a response-side no-Go-type finding at /%s", render(r.Findings), cmd.Response.Name)
	}
}

func TestFidelityCrossReferenceFoldsThroughAncestors(t *testing.T) {
	// A dropped parent folds its descendants in the fidelity report, so an
	// ancestor's presence has to confirm the child too.
	c := &Comparer{Dropped: map[string]map[string]bool{"x": {"/x_response/a": true}}}
	if !c.confirmed("x", "/x_response/a/b/c") {
		t.Errorf("A dropped ancestor must confirm its descendants")
	}
	if c.confirmed("x", "/x_response/ab") {
		t.Errorf("Prefix matching must respect path separators")
	}
	if c.confirmed("other", "/x_response/a") {
		t.Errorf("Confirmation must not leak across commands")
	}
}

func TestParseDroppedReadsAFidelityReport(t *testing.T) {
	report := `# fidelity: get_results / upstream / 002
decode: ok

DROPPED
  /get_results_response/result_count/#text        x1
  /get_results_response/sort/field                x1  (+1 descendant paths)

TRUNCATED
  (none)

ADDED
  /get_results_response/result/name               x1
`
	got := parseDropped(report)
	want := []string{"/get_results_response/result_count", "/get_results_response/sort/field"}
	if len(got) != len(want) {
		t.Fatalf("Got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Got %q, want %q", got[i], want[i])
		}
	}
}
