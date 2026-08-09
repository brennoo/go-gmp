package commands

import (
	"encoding/xml"
	"testing"
)

// chardataRoundTrip unmarshals in into a fresh *T, hands it to check for the
// field assertions specific to that case, then marshals it back and requires
// byte-identical output — the shared body every TestChardata*RoundTrip case
// below runs, split out so each case stays under the linter's complexity
// budget for a single test func.
func chardataRoundTrip[T any](t *testing.T, in string, check func(t *testing.T, got *T)) {
	t.Helper()
	got := new(T)
	if err := xml.Unmarshal([]byte(in), got); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	check(t, got)

	out, err := xml.Marshal(got)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if string(out) != in {
		t.Errorf("round-tripped XML = %q, want %q", out, in)
	}
}

// TestChardataRoundTripText exercises mixed content: an element that carries
// both an attribute and bare text of its own — <text excerpt="...">the
// text</text>, the shape the schema DSL spells with no dedicated keyword,
// which schema.Node's Chardata field and the generator's emitChardata now
// model as a `,chardata` field alongside the element's real attribute.
//
// The wrapper mirrors how Note actually field-tags NoteText (see
// commands/get_notes.go) rather than relying on the type's own name for the
// wire element, since NoteText is never marshaled standalone in production.
func TestChardataRoundTripText(t *testing.T) {
	type wrapper struct {
		XMLName xml.Name `xml:"note"`
		Text    NoteText `xml:"text"`
	}
	const in = `<note><text excerpt="1">This is the override text.</text></note>`

	chardataRoundTrip(t, in, func(t *testing.T, got *wrapper) {
		if got.Text.Value != "This is the override text." {
			t.Errorf("Value = %q, want %q", got.Text.Value, "This is the override text.")
		}
		if bool(got.Text.Excerpt) != true {
			t.Errorf("Excerpt = %v, want true", got.Text.Excerpt)
		}
	})
}

// TestChardataRoundTripHost is TestChardataRoundTripText's counterpart for an
// element with a real child *and* bare text — <host><hostname>...</hostname>
// the hostname</host>.
func TestChardataRoundTripHost(t *testing.T) {
	type wrapper struct {
		XMLName xml.Name `xml:"result"`
		Host    NoteHost `xml:"host"`
	}
	const in = `<result><host>192.0.2.1<asset asset_id="a1"></asset><hostname>example.com</hostname></host></result>`

	chardataRoundTrip(t, in, func(t *testing.T, got *wrapper) {
		if got.Host.Value != "192.0.2.1" {
			t.Errorf("Value = %q, want %q", got.Host.Value, "192.0.2.1")
		}
		if got.Host.Hostname != "example.com" {
			t.Errorf("Hostname = %q, want %q", got.Host.Hostname, "example.com")
		}
		if got.Host.Asset.AssetID != "a1" {
			t.Errorf("Asset.AssetID = %q, want %q", got.Host.Asset.AssetID, "a1")
		}
	})
}

// TestChardataRoundTripSort covers the third recurring shape: an element with
// bare text and a child but no attributes of its own — <sort>text<field>
// <order>...</order></field></sort>, shared by ~40 get_* commands' response.
func TestChardataRoundTripSort(t *testing.T) {
	type wrapper struct {
		XMLName xml.Name `xml:"get_notes_response"`
		Sort    NoteSort `xml:"sort"`
	}
	const in = `<get_notes_response><sort>text<field><order>ascending</order></field></sort></get_notes_response>`

	chardataRoundTrip(t, in, func(t *testing.T, got *wrapper) {
		if got.Sort.Value != "text" {
			t.Errorf("Value = %q, want %q", got.Sort.Value, "text")
		}
		if got.Sort.Field.Order != "ascending" {
			t.Errorf("Field.Order = %q, want %q", got.Sort.Field.Order, "ascending")
		}
	})
}
