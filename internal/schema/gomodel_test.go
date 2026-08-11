package schema

import (
	"encoding/xml"
	"testing"
)

// encoding/xml flattens an anonymous struct field's children into the
// parent element whether the field is a value or a pointer. ModelOf must
// agree, or the differ blames commands/ for a field nesting encoding/xml
// never produces.

func TestModelOfFlattensAnonymousStruct(t *testing.T) {
	type Embedded struct {
		A string `xml:"a"`
	}
	type Root struct {
		XMLName xml.Name `xml:"root"`
		Embedded
	}

	root := ModelOf(&Root{}).Root
	if len(root.Children) != 1 {
		t.Fatalf("got %d children, want 1: %+v", len(root.Children), root.Children)
	}
	if got := root.Children[0].Path("/root"); got != "/root/a" {
		t.Fatalf("got path %q, want /root/a", got)
	}
}

func TestModelOfFlattensAnonymousPointerStruct(t *testing.T) {
	type Embedded struct {
		A string `xml:"a"`
	}
	type Root struct {
		XMLName xml.Name `xml:"root"`
		*Embedded
	}

	root := ModelOf(&Root{}).Root
	if len(root.Children) != 1 {
		t.Fatalf("got %d children, want 1: %+v", len(root.Children), root.Children)
	}
	if got := root.Children[0].Path("/root"); got != "/root/a" {
		t.Fatalf("got path %q, want /root/a (not /root/Embedded/a)", got)
	}
}

// encoding/xml flattens an anonymous struct field even when it carries an
// xml tag: per the encoding/xml doc, only an anonymous field of interface
// type is treated as named. A tag on a struct-kind anonymous field is
// ignored for marshaling purposes (confirmed against encoding/xml directly,
// not just assumed), so ModelOf must ignore it too.
func TestModelOfIgnoresPlainExcludedField(t *testing.T) {
	type Root struct {
		XMLName xml.Name `xml:"root"`
		A       string   `xml:"-"`
	}
	if root := ModelOf(&Root{}).Root; len(root.Children) != 0 {
		t.Fatalf("got %d children, want 0: a plain xml:\"-\" field must stay invisible: %+v", len(root.Children), root.Children)
	}
}

// internal/codegen's duplicate-wire-name collision codec also uses xml:"-"
// (the only tag go vet's structtag check accepts on a repeated wire name),
// but carries the field's real identity in a private gmpwire tag that go vet
// doesn't inspect. ModelOf must recover that identity, or the comparer would
// see the field as absent rather than as the wire element it actually is.
func TestModelOfReadsGmpwireOnExcludedField(t *testing.T) {
	type Root struct {
		XMLName xml.Name `xml:"root"`
		A       string   `xml:"-" gmpwire:"x"`
	}
	root := ModelOf(&Root{}).Root
	if len(root.Children) != 1 {
		t.Fatalf("got %d children, want 1: %+v", len(root.Children), root.Children)
	}
	if got := root.Children[0].Name; got != "x" {
		t.Errorf("got name %q, want %q", got, "x")
	}
}

func TestModelOfFlattensAnonymousFieldEvenWithATag(t *testing.T) {
	type Embedded struct {
		A string `xml:"a"`
	}
	type Root struct {
		XMLName  xml.Name `xml:"root"`
		Embedded `xml:"wrapped"`
	}

	root := ModelOf(&Root{}).Root
	if len(root.Children) != 1 {
		t.Fatalf("got %d children, want 1: %+v", len(root.Children), root.Children)
	}
	if got := root.Children[0].Path("/root"); got != "/root/a" {
		t.Fatalf("got path %q, want /root/a", got)
	}
}
