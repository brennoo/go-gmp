package schema

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/brennoo/go-gmp/commands"
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

// internal/codegen's FieldOverrides mechanism adds a real, evidence-backed
// field the schema's own pattern never declares, and records why in a
// private gmpoverride tag. ModelOf must recover that reason onto the leaf's
// GoNode, or conform.go would have nothing to key its suppression off.
func TestModelOfReadsGmpoverrideOnLeafField(t *testing.T) {
	type Root struct {
		XMLName xml.Name `xml:"root"`
		A       string   `xml:"a" gmpoverride:"evidence for a"`
	}
	root := ModelOf(&Root{}).Root
	if len(root.Children) != 1 {
		t.Fatalf("got %d children, want 1: %+v", len(root.Children), root.Children)
	}
	if got := root.Children[0].OverrideReason; got != "evidence for a" {
		t.Errorf("got OverrideReason %q, want %q", got, "evidence for a")
	}
}

// A ",chardata" field sets Chardata on its *parent* node rather than
// creating a child of its own (see addField), so a gmpoverride tag on a
// chardata field must land on the parent as ChardataOverrideReason instead
// of OverrideReason.
func TestModelOfReadsGmpoverrideOnChardataField(t *testing.T) {
	type Root struct {
		XMLName xml.Name `xml:"root"`
		Value   string   `xml:",chardata" gmpoverride:"evidence for chardata"`
	}
	root := ModelOf(&Root{}).Root
	if root.OverrideReason != "" {
		t.Errorf("got OverrideReason %q on the root, want empty", root.OverrideReason)
	}
	if got := root.ChardataOverrideReason; got != "evidence for chardata" {
		t.Errorf("got ChardataOverrideReason %q, want %q", got, "evidence for chardata")
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

// ZeroIsMeaningful must recognize commands.GMPBool and commands.GMPAtoiBool
// as reflection actually spells them: package-qualified (f.Type.String()
// returns "commands.GMPBool", not "GMPBool"). A prior version only matched
// the unqualified name; using the real commands types via reflection here,
// rather than a manually constructed GoNode, is what would catch that again.
func TestZeroIsMeaningfulRecognizesQualifiedGMPBool(t *testing.T) {
	type Root struct {
		XMLName xml.Name         `xml:"root"`
		Active  commands.GMPBool `xml:"active"`
	}
	root := ModelOf(&Root{}).Root
	if len(root.Children) != 1 {
		t.Fatalf("got %d children, want 1: %+v", len(root.Children), root.Children)
	}
	if got := root.Children[0].ZeroIsMeaningful(); !got {
		t.Errorf("ZeroIsMeaningful() on non-pointer commands.GMPBool = %v, want true", got)
	}
}

func TestZeroIsMeaningfulRecognizesQualifiedGMPAtoiBool(t *testing.T) {
	type Root struct {
		XMLName xml.Name             `xml:"root"`
		Active  commands.GMPAtoiBool `xml:"active"`
	}
	root := ModelOf(&Root{}).Root
	if len(root.Children) != 1 {
		t.Fatalf("got %d children, want 1: %+v", len(root.Children), root.Children)
	}
	if got := root.Children[0].ZeroIsMeaningful(); !got {
		t.Errorf("ZeroIsMeaningful() on non-pointer commands.GMPAtoiBool = %v, want true", got)
	}
}

func TestZeroIsMeaningfulIgnoresPointerGMPBool(t *testing.T) {
	type Root struct {
		XMLName xml.Name          `xml:"root"`
		Active  *commands.GMPBool `xml:"active"`
	}
	root := ModelOf(&Root{}).Root
	if len(root.Children) != 1 {
		t.Fatalf("got %d children, want 1: %+v", len(root.Children), root.Children)
	}
	if got := root.Children[0].ZeroIsMeaningful(); got {
		t.Errorf("ZeroIsMeaningful() on pointer *commands.GMPBool = %v, want false", got)
	}
}

func TestZeroIsMeaningfulIgnoresPointerGMPAtoiBool(t *testing.T) {
	type Root struct {
		XMLName xml.Name              `xml:"root"`
		Active  *commands.GMPAtoiBool `xml:"active"`
	}
	root := ModelOf(&Root{}).Root
	if len(root.Children) != 1 {
		t.Fatalf("got %d children, want 1: %+v", len(root.Children), root.Children)
	}
	if got := root.Children[0].ZeroIsMeaningful(); got {
		t.Errorf("ZeroIsMeaningful() on pointer *commands.GMPAtoiBool = %v, want false", got)
	}
}

// T-34: a `,innerxml` field captures every child of its element verbatim
// (see commands.HelpResponse.Text), so ModelOf must mark the enclosing node
// CatchAll the same way it already does for `,any` — not just Chardata.
func TestModelOfInnerxmlSetsCatchAll(t *testing.T) {
	type Root struct {
		XMLName xml.Name `xml:"root"`
		Text    string   `xml:",innerxml"`
	}
	root := ModelOf(&Root{}).Root
	if !root.Chardata {
		t.Errorf("Chardata = false, want true")
	}
	if !root.CatchAll {
		t.Errorf("CatchAll = false, want true")
	}
	if !strings.Contains(root.CatchAllReason, "innerxml") {
		t.Errorf("CatchAllReason = %q, want it to mention innerxml", root.CatchAllReason)
	}
}

// A `,comment` field captures XML comments, never child elements, so it must
// not earn innerxml's CatchAll treatment.
func TestModelOfCommentDoesNotSetCatchAll(t *testing.T) {
	type Root struct {
		XMLName xml.Name `xml:"root"`
		Note    string   `xml:",comment"`
	}
	root := ModelOf(&Root{}).Root
	if !root.Chardata {
		t.Errorf("Chardata = false, want true")
	}
	if root.CatchAll {
		t.Errorf("CatchAll = true, want false: `,comment` does not capture elements")
	}
}

// The named-field counterpart to TestModelOfReadsGmpoverrideOnLeafField:
// gmpmissing/gmpmissingreason on a struct-typed field must land on that
// field's own node, since unmatchedSchema consults them while walking that
// node's children.
func TestModelOfReadsGmpmissingOnStructField(t *testing.T) {
	type Value struct {
		ReportFormat string `xml:"report_format"`
	}
	type Root struct {
		XMLName xml.Name `xml:"root"`
		Value   Value    `xml:"value" gmpmissing:"attr:using_default" gmpmissingreason:"relocated to param"`
	}
	root := ModelOf(&Root{}).Root
	if len(root.Children) != 1 {
		t.Fatalf("got %d children, want 1: %+v", len(root.Children), root.Children)
	}
	value := root.Children[0]
	if got := value.MissingChildKeys; len(got) != 1 || got[0] != "attr:using_default" {
		t.Errorf("MissingChildKeys = %v, want [\"attr:using_default\"]", got)
	}
	if got := value.MissingChildReason; got != "relocated to param" {
		t.Errorf("MissingChildReason = %q, want %q", got, "relocated to param")
	}
}
