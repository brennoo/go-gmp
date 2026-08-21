package gmptest

import (
	"encoding/xml"
	"strings"
	"testing"
)

// These exercise the differ against hand-written types small enough to reason
// about exactly, so a differ bug can't be mistaken for a library bug.

func TestExtractPaths(t *testing.T) {
	tests := []struct {
		name string
		doc  string
		want map[string]int
	}{
		{
			name: "element attribute and text",
			doc:  `<a id="1"><b>x</b></a>`,
			want: map[string]int{"/a": 1, "/a/@id": 1, "/a/b": 1, "/a/b/#text": 1},
		},
		{
			name: "repeated siblings are counted, not indexed",
			doc:  `<a><b/><b/><b/></a>`,
			want: map[string]int{"/a": 1, "/a/b": 3},
		},
		{
			name: "whitespace-only text is formatting, not content",
			doc:  "<a>\n  <b>x</b>\n</a>",
			want: map[string]int{"/a": 1, "/a/b": 1, "/a/b/#text": 1},
		},
		{
			name: "self-closing and empty are indistinguishable",
			doc:  `<a><b/><c></c></a>`,
			want: map[string]int{"/a": 1, "/a/b": 1, "/a/c": 1},
		},
		{
			// The xmlns declaration itself is not a path, but it does put the
			// element into a namespace, which the path records.
			name: "xmlns attributes are skipped but the namespace qualifies the name",
			doc:  `<a xmlns="urn:x" id="1"/>`,
			want: map[string]int{"/{urn:x}a": 1, "/{urn:x}a/@id": 1},
		},
		{
			name: "mixed content records both text and child",
			doc:  `<a>text<b>y</b></a>`,
			want: map[string]int{"/a": 1, "/a/#text": 1, "/a/b": 1, "/a/b/#text": 1},
		},
		{
			name: "same name at different depths are distinct paths",
			doc:  `<a><b/><c><b/></c></a>`,
			want: map[string]int{"/a": 1, "/a/b": 1, "/a/c": 1, "/a/c/b": 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ps, err := ExtractPaths([]byte(tc.doc))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if len(ps.Counts) != len(tc.want) {
				t.Fatalf("Got %d paths %v, want %d %v", len(ps.Counts), ps.Paths(), len(tc.want), tc.want)
			}
			for path, n := range tc.want {
				if ps.Counts[path] != n {
					t.Errorf("Path %q: got count %d, want %d", path, ps.Counts[path], n)
				}
			}
		})
	}
}

func TestExtractPathsRecordsNamespaces(t *testing.T) {
	ps, err := ExtractPaths([]byte(`<a xmlns="urn:example"><b/></a>`))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(ps.Namespaces) != 1 || ps.Namespaces[0] != "urn:example" {
		t.Fatalf("Expected the namespace to be recorded, got %v", ps.Namespaces)
	}
}

func TestExtractPathsRejectsMalformed(t *testing.T) {
	if _, err := ExtractPaths([]byte(`<a><b></a>`)); err == nil {
		t.Fatalf("Expected an error for mismatched tags")
	}
}

// --- the four buckets -------------------------------------------------------

type dropRoot struct {
	XMLName xml.Name `xml:"a"`
	B       string   `xml:"b"`
	// no field for <c>
}

func TestCompareDropped(t *testing.T) {
	f := Compare([]byte(`<a><b>1</b><c>2</c></a>`), &dropRoot{})
	if f.DecodeErr != nil {
		t.Fatalf("Unexpected decode error: %v", f.DecodeErr)
	}
	if f.Clean() {
		t.Fatalf("Expected loss to be reported")
	}
	if len(f.Dropped) != 1 || f.Dropped[0].Path != "/a/c" {
		t.Fatalf("Expected /a/c dropped, got %+v", f.Dropped)
	}
	if f.Dropped[0].Suppressed != 1 {
		t.Fatalf("Expected /a/c/#text folded into the parent, got %d", f.Dropped[0].Suppressed)
	}
}

type truncRoot struct {
	XMLName xml.Name `xml:"a"`
	B       *struct {
		V string `xml:",chardata"`
	} `xml:"b"`
}

// A single-valued field against a repeated element keeps one occurrence and
// silently discards the rest — loss a "dropped" check can't see, hence the
// separate TRUNCATED bucket.
func TestCompareTruncated(t *testing.T) {
	f := Compare([]byte(`<a><b>1</b><b>2</b><b>3</b></a>`), &truncRoot{})
	if f.DecodeErr != nil {
		t.Fatalf("Unexpected decode error: %v", f.DecodeErr)
	}
	// Both /a/b and its text are truncated 3 -> 1. Descendants are deliberately
	// not folded away here, unlike DROPPED: a surviving occurrence may still
	// carry all of its children, so the counts are independently meaningful.
	want := map[string]bool{"/a/b": false, "/a/b/#text": false}
	for _, got := range f.Truncated {
		if _, ok := want[got.Path]; !ok {
			t.Errorf("Unexpected truncation %+v", got)
			continue
		}
		want[got.Path] = true
		if got.Src != 3 || got.Got != 1 {
			t.Errorf("Expected %s 3 -> 1, got %d -> %d", got.Path, got.Src, got.Got)
		}
	}
	for path, seen := range want {
		if !seen {
			t.Errorf("Expected %s to be reported truncated", path)
		}
	}
	if f.Clean() {
		t.Fatalf("Truncation is loss and must not be reported clean")
	}
}

type addRoot struct {
	XMLName xml.Name `xml:"a"`
	B       string   `xml:"b"`
	// Always marshals, even when the source had no <c>: the exact shape of the
	// non-pointer no-omitempty fields in GetTasksResponse.
	C string `xml:"c"`
}

func TestCompareAdded(t *testing.T) {
	f := Compare([]byte(`<a><b>1</b></a>`), &addRoot{})
	if f.DecodeErr != nil {
		t.Fatalf("Unexpected decode error: %v", f.DecodeErr)
	}
	if len(f.Added) != 1 || f.Added[0].Path != "/a/c" {
		t.Fatalf("Expected /a/c added, got %+v", f.Added)
	}
	if !f.Clean() {
		t.Fatalf("Added elements are not loss and must still report clean")
	}
}

type innerRoot struct {
	XMLName  xml.Name `xml:"a"`
	InnerXML []byte   `xml:",innerxml"`
	B        string   `xml:"b"`
}

// ,innerxml re-emits the raw bytes and the sibling field, so children appear
// twice. RunWizardResponse is the one place this happens for real.
func TestCompareDuplicated(t *testing.T) {
	f := Compare([]byte(`<a><b>1</b></a>`), &innerRoot{})
	if f.DecodeErr != nil {
		t.Fatalf("Unexpected decode error: %v", f.DecodeErr)
	}
	if len(f.Duplicated) == 0 {
		t.Fatalf("Expected duplication, got %+v", f)
	}
	if !f.Clean() {
		t.Fatalf("Duplication is not loss and must still report clean")
	}
}

type mixedRoot struct {
	XMLName xml.Name `xml:"a"`
	Value   string   `xml:",chardata"`
	B       string   `xml:"b"`
}

// The GetTasksResponseSort shape: chardata alongside a child element. Both
// survive, so nothing should be reported.
func TestCompareMixedContentIsClean(t *testing.T) {
	f := Compare([]byte(`<a>text<b>y</b></a>`), &mixedRoot{})
	if f.DecodeErr != nil {
		t.Fatalf("Unexpected decode error: %v", f.DecodeErr)
	}
	if !f.Clean() {
		t.Fatalf("Expected a clean round trip, got dropped=%+v truncated=%+v", f.Dropped, f.Truncated)
	}
}

type intRoot struct {
	XMLName xml.Name `xml:"a"`
	Count   int      `xml:"count"`
}

// The <report_count>2<finished>2</finished></report_count> shape: the int
// parses, the nested element vanishes. Exactly the signal we want.
func TestCompareIntFieldDropsNestedElement(t *testing.T) {
	f := Compare([]byte(`<a><count>2<finished>2</finished></count></a>`), &intRoot{})
	if f.DecodeErr != nil {
		t.Fatalf("Unexpected decode error: %v", f.DecodeErr)
	}
	if len(f.Dropped) != 1 || f.Dropped[0].Path != "/a/count/finished" {
		t.Fatalf("Expected /a/count/finished dropped, got %+v", f.Dropped)
	}
}

type boolRoot struct {
	XMLName xml.Name `xml:"a"`
	InUse   bool     `xml:"in_use,omitempty"`
}

// A false boolean tagged omitempty cannot round-trip: it decodes to false and
// then vanishes. 103 fields in commands/ have this shape, and it is why the
// report is asserted unchanged rather than empty.
func TestCompareFalseBoolWithOmitemptyIsDropped(t *testing.T) {
	f := Compare([]byte(`<a><in_use>0</in_use></a>`), &boolRoot{})
	if f.DecodeErr != nil {
		t.Fatalf("Unexpected decode error: %v", f.DecodeErr)
	}
	if len(f.Dropped) != 1 || f.Dropped[0].Path != "/a/in_use" {
		t.Fatalf("Expected /a/in_use dropped, got %+v", f.Dropped)
	}

	// The true case keeps the element but rewrites the value: GMP booleans are
	// "0 is false, anything else is true", and Go emits the Go spelling. The
	// path survives, so only a value comparison can see this.
	g := Compare([]byte(`<a><in_use>1</in_use></a>`), &boolRoot{})
	if len(g.Dropped) != 0 {
		t.Fatalf("Expected a true boolean to keep its element, got %+v", g.Dropped)
	}
	if len(g.Changed) != 1 || g.Changed[0].SrcValue != "1" || g.Changed[0].GotValue != "true" {
		t.Fatalf("Expected /a/in_use/#text \"1\" -> \"true\", got %+v", g.Changed)
	}
}

type attrRoot struct {
	XMLName xml.Name `xml:"a"`
	ID      string   `xml:"id,attr"`
}

func TestCompareDroppedAttribute(t *testing.T) {
	f := Compare([]byte(`<a id="1" status="200"/>`), &attrRoot{})
	if f.DecodeErr != nil {
		t.Fatalf("Unexpected decode error: %v", f.DecodeErr)
	}
	if len(f.Dropped) != 1 || f.Dropped[0].Path != "/a/@status" {
		t.Fatalf("Expected /a/@status dropped, got %+v", f.Dropped)
	}
}

type wrongRoot struct {
	XMLName xml.Name `xml:"b"`
}

// XMLName enforcement makes a mismatched root a hard decode failure, as seen
// in one upstream get_preferences example.
func TestCompareRecordsDecodeFailure(t *testing.T) {
	f := Compare([]byte(`<a/>`), &wrongRoot{})
	if f.DecodeErr == nil {
		t.Fatalf("Expected a decode error for a mismatched root element")
	}
	if f.Clean() {
		t.Fatalf("A fixture that does not decode must not report clean")
	}
	if r := f.Report("# h\n"); !strings.Contains(r, "decode: error") {
		t.Fatalf("Expected the report to record the decode failure, got:\n%s", r)
	}
}

// --- report rendering -------------------------------------------------------

func TestReportIsDeterministic(t *testing.T) {
	f := Compare([]byte(`<a><b>1</b><c>2</c><d>3</d></a>`), &dropRoot{})
	first := f.Report("# fixture\n")
	for i := range 20 {
		if got := f.Report("# fixture\n"); got != first {
			t.Fatalf("Report is not deterministic; run %d differed:\n%s\nvs\n%s", i, got, first)
		}
	}
	for _, want := range []string{"DROPPED", "TRUNCATED", "ADDED", "DUPLICATED", "decode: ok"} {
		if !strings.Contains(first, want) {
			t.Errorf("Report is missing the %s section:\n%s", want, first)
		}
	}
}

func TestReportShowsNoneForEmptySections(t *testing.T) {
	f := Compare([]byte(`<a><b>1</b></a>`), &struct {
		XMLName xml.Name `xml:"a"`
		B       string   `xml:"b"`
	}{})
	if !strings.Contains(f.Report("# h\n"), "(none)") {
		t.Fatalf("Expected empty sections to render as (none)")
	}
}

// --- regressions from review ------------------------------------------------

type hyphenRoot struct {
	XMLName xml.Name `xml:"a"`
	// no field for <b>, <b-c> or <b>'s children
}

// '-' sorts before '/', so /a/b-c falls between /a/b and /a/b/d; comparing
// each finding only against the previously emitted one left /a/b/d unfolded.
// <cpe-item> makes this reachable with real documents, not just contrived ones.
func TestSuppressDescendantsAcrossSortingGap(t *testing.T) {
	f := Compare([]byte(`<a><b><d>1</d></b><b-c>2</b-c></a>`), &hyphenRoot{})
	if f.DecodeErr != nil {
		t.Fatalf("Unexpected decode error: %v", f.DecodeErr)
	}

	got := map[string]int{}
	for _, d := range f.Dropped {
		got[d.Path] = d.Suppressed
	}
	if _, folded := got["/a/b/d"]; folded {
		t.Errorf("/a/b/d should have been folded into /a/b, but was reported separately: %+v", f.Dropped)
	}
	if n, ok := got["/a/b"]; !ok {
		t.Errorf("Expected /a/b reported dropped, got %+v", f.Dropped)
	} else if n != 2 { // /a/b/d and /a/b/d/#text
		t.Errorf("Expected /a/b to absorb 2 descendants, got %d", n)
	}
	if _, ok := got["/a/b-c"]; !ok {
		t.Errorf("Expected the sibling /a/b-c reported separately, got %+v", f.Dropped)
	}
}

type nsRoot struct {
	XMLName xml.Name `xml:"a"`
	X       string   `xml:"x"`
}

// Two same-local-name elements in different namespaces used to collapse onto
// one path, so losing one could still report clean. Paths are now
// namespace-qualified, and any namespaced document is marked advisory rather
// than certified.
func TestNamespacedDocumentIsNeverClean(t *testing.T) {
	doc := []byte(`<a xmlns:p="urn:p" xmlns:q="urn:q"><p:x>1</p:x><q:x>2</q:x></a>`)

	ps, err := ExtractPaths(doc)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if ps.Counts["/a/{urn:p}x"] != 1 || ps.Counts["/a/{urn:q}x"] != 1 {
		t.Fatalf("Expected the two namespaces to stay distinct, got %v", ps.Paths())
	}

	f := Compare(doc, &nsRoot{})
	if !f.Advisory() {
		t.Fatalf("Expected a namespaced document to be advisory")
	}
	if !strings.Contains(f.Report("# h\n"), "ADVISORY") {
		t.Fatalf("Expected the report to carry the advisory banner")
	}

	// Isolate the rule: a comparison with no findings at all must still refuse
	// to certify a namespaced document, because the diff it is based on is not
	// one encoding/xml can be trusted to produce faithfully.
	if (Fidelity{Namespaces: []string{"urn:x"}}).Clean() {
		t.Fatalf("An advisory comparison must never report clean, even with no findings")
	}
	if !(Fidelity{}).Clean() {
		t.Fatalf("A namespace-free comparison with no findings must report clean")
	}
}

type coerceRoot struct {
	XMLName xml.Name `xml:"a"`
	N       float64  `xml:"n"`
	S       string   `xml:"s"`
}

// Structural retention is not fidelity: a value the library can't reproduce
// is information it didn't keep, even though the path survived.
func TestCompareReportsChangedValues(t *testing.T) {
	f := Compare([]byte(`<a><n>5.0</n><s>keep</s></a>`), &coerceRoot{})
	if f.DecodeErr != nil {
		t.Fatalf("Unexpected decode error: %v", f.DecodeErr)
	}
	if len(f.Changed) != 1 {
		t.Fatalf("Expected one changed value, got %+v", f.Changed)
	}
	if c := f.Changed[0]; c.Path != "/a/n/#text" || c.SrcValue != "5.0" || c.GotValue != "5" {
		t.Fatalf("Expected /a/n/#text 5.0 -> 5, got %+v", c)
	}
	if f.Clean() {
		t.Fatalf("A changed value is loss and must not report clean")
	}
}

type attrValRoot struct {
	XMLName xml.Name `xml:"a"`
	Flag    bool     `xml:"flag,attr"`
}

func TestCompareReportsChangedAttributeValues(t *testing.T) {
	f := Compare([]byte(`<a flag="1"/>`), &attrValRoot{})
	if f.DecodeErr != nil {
		t.Fatalf("Unexpected decode error: %v", f.DecodeErr)
	}
	if len(f.Changed) != 1 || f.Changed[0].Path != "/a/@flag" {
		t.Fatalf("Expected /a/@flag changed, got %+v", f.Changed)
	}
}

// Repeated paths are deliberately not value-compared: without a positional
// index there is no way to align the nth occurrence across a reordering.
func TestCompareSkipsValuesOfRepeatedPaths(t *testing.T) {
	f := Compare([]byte(`<a><s>x</s><s>y</s></a>`), &struct {
		XMLName xml.Name `xml:"a"`
		S       []string `xml:"s"`
	}{})
	if len(f.Changed) != 0 {
		t.Fatalf("Expected no value comparison on a repeated path, got %+v", f.Changed)
	}
}
