package commands

import (
	"encoding/xml"
	"testing"
)

type gmpAtoiBoolElement struct {
	XMLName xml.Name    `xml:"root"`
	V       GMPAtoiBool `xml:"v"`
}

type gmpAtoiBoolPtrElement struct {
	XMLName xml.Name     `xml:"root"`
	V       *GMPAtoiBool `xml:"v,omitempty"`
}

type gmpAtoiBoolAttr struct {
	XMLName xml.Name    `xml:"root"`
	V       GMPAtoiBool `xml:"v,attr"`
}

type gmpAtoiBoolPtrAttr struct {
	XMLName xml.Name     `xml:"root"`
	V       *GMPAtoiBool `xml:"v,attr,omitempty"`
}

func TestGMPAtoiBoolMarshalElement(t *testing.T) {
	tests := []struct {
		name string
		v    GMPAtoiBool
		want string
	}{
		{"false", GMPAtoiBool(false), `<root><v>0</v></root>`},
		{"true", GMPAtoiBool(true), `<root><v>1</v></root>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := xml.Marshal(gmpAtoiBoolElement{V: tt.v})
			if err != nil {
				t.Fatalf("Marshal returned error: %v", err)
			}
			if string(out) != tt.want {
				t.Errorf("Marshal(%v) = %q, want %q", tt.v, out, tt.want)
			}
		})
	}
}

func TestGMPAtoiBoolMarshalAttr(t *testing.T) {
	tests := []struct {
		name string
		v    GMPAtoiBool
		want string
	}{
		{"false", GMPAtoiBool(false), `<root v="0"></root>`},
		{"true", GMPAtoiBool(true), `<root v="1"></root>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := xml.Marshal(gmpAtoiBoolAttr{V: tt.v})
			if err != nil {
				t.Fatalf("Marshal returned error: %v", err)
			}
			if string(out) != tt.want {
				t.Errorf("Marshal(%v) = %q, want %q", tt.v, out, tt.want)
			}
		})
	}
}

// atoiBoolCases lists decode inputs shared by the element and attribute
// tests, so the two paths are exercised against the same atoi(3) policy.
var atoiBoolCases = []struct {
	name string
	text string
	want bool
}{
	{"zero is false", "0", false},
	{"one is true", "1", true},
	{"arbitrary digit is true", "7", true},
	{"empty text is false", "", false},
	{"nonnumeric text is false", "abc", false},
	{"leading whitespace then digit is true", "  1", true},
	{"leading plus sign is true", "+1", true},
	{"leading minus sign nonzero is true", "-1", true},
	{"negative zero is false", "-0", false},
	{"numeric prefix before letters is true", "1x", true},
}

func TestGMPAtoiBoolUnmarshalElement(t *testing.T) {
	for _, tt := range atoiBoolCases {
		t.Run(tt.name, func(t *testing.T) {
			doc := `<root><v>` + tt.text + `</v></root>`
			var got gmpAtoiBoolElement
			if err := xml.Unmarshal([]byte(doc), &got); err != nil {
				t.Fatalf("Unmarshal returned error: %v", err)
			}
			if bool(got.V) != tt.want {
				t.Errorf("Unmarshal(%q) = %v, want %v", doc, got.V, tt.want)
			}
		})
	}

	t.Run("self-closing element is false", func(t *testing.T) {
		var got gmpAtoiBoolElement
		if err := xml.Unmarshal([]byte(`<root><v/></root>`), &got); err != nil {
			t.Fatalf("Unmarshal returned error: %v", err)
		}
		if bool(got.V) != false {
			t.Errorf("Unmarshal(self-closing) = %v, want false", got.V)
		}
	})
}

func TestGMPAtoiBoolUnmarshalAttr(t *testing.T) {
	for _, tt := range atoiBoolCases {
		t.Run(tt.name, func(t *testing.T) {
			doc := `<root v="` + tt.text + `"></root>`
			var got gmpAtoiBoolAttr
			if err := xml.Unmarshal([]byte(doc), &got); err != nil {
				t.Fatalf("Unmarshal returned error: %v", err)
			}
			if bool(got.V) != tt.want {
				t.Errorf("Unmarshal(%q) = %v, want %v", doc, got.V, tt.want)
			}
		})
	}
}

func TestGMPAtoiBoolPointerRoundTrip(t *testing.T) {
	// *GMPAtoiBool is the shape actually used on request-root fields, so
	// verify it goes through the same custom marshal/unmarshal path as the
	// value type rather than encoding/xml's default bool handling.
	marshaled, err := xml.Marshal(gmpAtoiBoolPtrElement{V: AtoiBool(false)})
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if want := `<root><v>0</v></root>`; string(marshaled) != want {
		t.Errorf("Marshal(AtoiBool(false)) = %q, want %q", marshaled, want)
	}

	var decoded gmpAtoiBoolPtrElement
	if err := xml.Unmarshal(marshaled, &decoded); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if decoded.V == nil || bool(*decoded.V) != false {
		t.Errorf("round-tripped value = %v, want false", decoded.V)
	}

	marshaledAttr, err := xml.Marshal(gmpAtoiBoolPtrAttr{V: AtoiBool(true)})
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if want := `<root v="1"></root>`; string(marshaledAttr) != want {
		t.Errorf("Marshal(AtoiBool(true)) = %q, want %q", marshaledAttr, want)
	}
}

func TestGMPAtoiBoolPointerOmitEmpty(t *testing.T) {
	// A nil pointer must be omitted, distinguishing "not set" from an
	// explicit false — this is the reason AtoiBool exists.
	out, err := xml.Marshal(gmpAtoiBoolPtrElement{V: nil})
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if want := `<root></root>`; string(out) != want {
		t.Errorf("Marshal(nil) = %q, want %q", out, want)
	}

	out, err = xml.Marshal(gmpAtoiBoolPtrElement{V: AtoiBool(false)})
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if want := `<root><v>0</v></root>`; string(out) != want {
		t.Errorf("Marshal(AtoiBool(false)) = %q, want %q", out, want)
	}
}

func TestGMPAtoiBoolUnmarshalMalformedXML(t *testing.T) {
	tests := []struct {
		name string
		xml  string
	}{
		{"unclosed element", `<root><v>0</root>`},
		{"unclosed attribute value", `<root v="0></root>`},
		{"not xml at all", `not xml`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got gmpAtoiBoolElement
			if err := xml.Unmarshal([]byte(tt.xml), &got); err == nil {
				t.Errorf("Unmarshal(%q) returned no error, want an error", tt.xml)
			}
		})
	}
}

func TestAtoiBoolHelper(t *testing.T) {
	if got := AtoiBool(false); got == nil || bool(*got) != false {
		t.Errorf("AtoiBool(false) = %v, want pointer to false", got)
	}
	if got := AtoiBool(true); got == nil || bool(*got) != true {
		t.Errorf("AtoiBool(true) = %v, want pointer to true", got)
	}
}
