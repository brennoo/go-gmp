package commands

import (
	"encoding/xml"
	"testing"
)

type gmpBoolElement struct {
	XMLName xml.Name `xml:"root"`
	V       GMPBool  `xml:"v"`
}

type gmpBoolPtrElement struct {
	XMLName xml.Name `xml:"root"`
	V       *GMPBool `xml:"v,omitempty"`
}

type gmpBoolAttr struct {
	XMLName xml.Name `xml:"root"`
	V       GMPBool  `xml:"v,attr"`
}

type gmpBoolPtrAttr struct {
	XMLName xml.Name `xml:"root"`
	V       *GMPBool `xml:"v,attr,omitempty"`
}

func TestGMPBoolMarshalElement(t *testing.T) {
	tests := []struct {
		name string
		v    GMPBool
		want string
	}{
		{"false", GMPBool(false), `<root><v>0</v></root>`},
		{"true", GMPBool(true), `<root><v>1</v></root>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := xml.Marshal(gmpBoolElement{V: tt.v})
			if err != nil {
				t.Fatalf("Marshal returned error: %v", err)
			}
			if string(out) != tt.want {
				t.Errorf("Marshal(%v) = %q, want %q", tt.v, out, tt.want)
			}
		})
	}
}

func TestGMPBoolMarshalAttr(t *testing.T) {
	tests := []struct {
		name string
		v    GMPBool
		want string
	}{
		{"false", GMPBool(false), `<root v="0"></root>`},
		{"true", GMPBool(true), `<root v="1"></root>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := xml.Marshal(gmpBoolAttr{V: tt.v})
			if err != nil {
				t.Fatalf("Marshal returned error: %v", err)
			}
			if string(out) != tt.want {
				t.Errorf("Marshal(%v) = %q, want %q", tt.v, out, tt.want)
			}
		})
	}
}

func TestGMPBoolUnmarshalElement(t *testing.T) {
	tests := []struct {
		name string
		xml  string
		want bool
	}{
		{"zero is false", `<root><v>0</v></root>`, false},
		{"one is true", `<root><v>1</v></root>`, true},
		{"arbitrary digit is true", `<root><v>7</v></root>`, true},
		{"empty element is true", `<root><v></v></root>`, true},
		{"self-closing element is true", `<root><v/></root>`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got gmpBoolElement
			if err := xml.Unmarshal([]byte(tt.xml), &got); err != nil {
				t.Fatalf("Unmarshal returned error: %v", err)
			}
			if bool(got.V) != tt.want {
				t.Errorf("Unmarshal(%q) = %v, want %v", tt.xml, got.V, tt.want)
			}
		})
	}
}

func TestGMPBoolUnmarshalAttr(t *testing.T) {
	tests := []struct {
		name string
		xml  string
		want bool
	}{
		{"zero is false", `<root v="0"></root>`, false},
		{"one is true", `<root v="1"></root>`, true},
		{"arbitrary digit is true", `<root v="7"></root>`, true},
		{"empty value is true", `<root v=""></root>`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got gmpBoolAttr
			if err := xml.Unmarshal([]byte(tt.xml), &got); err != nil {
				t.Fatalf("Unmarshal returned error: %v", err)
			}
			if bool(got.V) != tt.want {
				t.Errorf("Unmarshal(%q) = %v, want %v", tt.xml, got.V, tt.want)
			}
		})
	}
}

func TestGMPBoolPointerRoundTrip(t *testing.T) {
	// *GMPBool is the shape actually used on request-root fields, so verify
	// the pointer field goes through the same custom marshal/unmarshal path
	// as the value type (rather than encoding/xml's default bool handling).
	marshaled, err := xml.Marshal(gmpBoolPtrElement{V: Bool(false)})
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if want := `<root><v>0</v></root>`; string(marshaled) != want {
		t.Errorf("Marshal(Bool(false)) = %q, want %q", marshaled, want)
	}

	var decoded gmpBoolPtrElement
	if err := xml.Unmarshal(marshaled, &decoded); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if decoded.V == nil || bool(*decoded.V) != false {
		t.Errorf("round-tripped value = %v, want false", decoded.V)
	}

	marshaledAttr, err := xml.Marshal(gmpBoolPtrAttr{V: Bool(true)})
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if want := `<root v="1"></root>`; string(marshaledAttr) != want {
		t.Errorf("Marshal(Bool(true)) = %q, want %q", marshaledAttr, want)
	}
}

func TestGMPBoolUnmarshalMalformedXML(t *testing.T) {
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
			var got gmpBoolElement
			if err := xml.Unmarshal([]byte(tt.xml), &got); err == nil {
				t.Errorf("Unmarshal(%q) returned no error, want an error", tt.xml)
			}
		})
	}
}

func TestBoolHelper(t *testing.T) {
	if got := Bool(false); got == nil || bool(*got) != false {
		t.Errorf("Bool(false) = %v, want pointer to false", got)
	}
	if got := Bool(true); got == nil || bool(*got) != true {
		t.Errorf("Bool(true) = %v, want pointer to true", got)
	}
}
