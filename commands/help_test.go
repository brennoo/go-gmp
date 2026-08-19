package commands

import (
	"encoding/xml"
	"strings"
	"testing"
)

// Help's format/type attributes are both declared without a schema
// <required> flag (vendor/XML/GMP.xml.in's "Required help format" summary is
// documentation prose, not an enforced wire constraint), so an unset Help
// must marshal with neither attribute present rather than fabricating an
// empty one.
func TestHelpRequestAttributesAbsent(t *testing.T) {
	out, err := xml.Marshal(&Help{})
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if got := string(out); got != `<help></help>` {
		t.Fatalf("Marshal = %s, want <help></help> with no attributes", got)
	}
}

// A set Format/Type round-trips verbatim. Note that a live-server
// introspection call sends explicit-empty type="" — the omitempty-tagged
// Type field here cannot represent that (an empty Go string is simply
// omitted, not serialized as type=""); this only covers the non-empty case.
func TestHelpRequestAttributesPresent(t *testing.T) {
	req := Help{Format: "XML", Type: "brief"}

	out, err := xml.Marshal(&req)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if got := string(out); !strings.Contains(got, `format="XML"`) || !strings.Contains(got, `type="brief"`) {
		t.Fatalf("Marshal = %s, want both format and type present", got)
	}
}

// The plain-text format's response has no nested elements, so innerxml and
// chardata agree here; this pins that switching to innerxml (for the XML
// format's sake, below) didn't regress the plain-text case.
func TestHelpResponseTextRoundTripsPlainText(t *testing.T) {
	in := `<help_response status="200" status_text="OK">plain help text, any format</help_response>`

	var r HelpResponse
	if err := xml.Unmarshal([]byte(in), &r); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if r.Text != "plain help text, any format" {
		t.Fatalf("Text = %q, want the chardata preserved verbatim", r.Text)
	}
}

// The XML format's response nests a <schema><protocol>...</protocol></schema>
// subtree (see internal/gmptest/testdata/upstream/help/002_response.xml).
// chardata drops this subtree entirely (confirmed by that fixture's
// 002_fidelity.txt, which records /help_response/schema and 28 descendant
// paths as DROPPED); innerxml preserves it verbatim as opaque raw XML.
func TestHelpResponseTextRoundTripsXMLSchema(t *testing.T) {
	in := `<help_response status="200" status_text="OK">` +
		`<schema format="XML" extension="xml" content_type="text/xml">` +
		`<protocol><name>Greenbone Management Protocol</name><abbreviation>GMP</abbreviation>` +
		`<command><name>authenticate</name><summary>Authenticate with the manager</summary></command>` +
		`</protocol></schema></help_response>`

	var r HelpResponse
	if err := xml.Unmarshal([]byte(in), &r); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !strings.Contains(r.Text, "<schema") || !strings.Contains(r.Text, "<command><name>authenticate</name>") {
		t.Fatalf("Text = %q, want the nested <schema> subtree preserved verbatim as raw XML", r.Text)
	}
}
