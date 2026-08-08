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
