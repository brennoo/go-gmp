package schema

import (
	"sort"
	"strings"
	"testing"
)

// These run against hand-written DSL snippets small enough to reason about
// exactly. Pointing the parser at the real schema first would make a parser
// bug indistinguishable from a struct bug, which is the one thing this whole
// package exists to tell apart.

// doc wraps a command body in the <protocol> envelope every schema file uses.
func doc(body string) string {
	return `<?xml version="1.0"?><protocol>
  <type><name>uuid</name><pattern>text</pattern></type>
  <type><name>boolean</name><pattern>text</pattern></type>
  <type><name>integer</name><pattern>text</pattern></type>
  <type><name>status</name><pattern>text</pattern></type>
  <command>` + body + `</command></protocol>`
}

func parseOne(t *testing.T, body string) Command {
	t.Helper()
	cmds, err := Parse([]byte(doc(body)), "test.xml")
	if err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if len(cmds) != 1 {
		t.Fatalf("Got %d commands, want 1", len(cmds))
	}
	return cmds[0]
}

// find locates a direct child by name, failing the test when it is absent.
func find(t *testing.T, n *Node, name string) *Node {
	t.Helper()
	for _, c := range n.Children {
		if c.Name == name {
			return c
		}
	}
	t.Fatalf("No child %q under %q; have %v", name, n.Name, names(n))
	return nil
}

func names(n *Node) []string {
	out := make([]string, 0, len(n.Children))
	for _, c := range n.Children {
		out = append(out, c.Kind.String()+":"+c.Name)
	}
	return out
}

func defectClasses(c Command) []string {
	out := make([]string, 0, len(c.Defects))
	for _, d := range c.Defects {
		out = append(out, d.Class)
	}
	sort.Strings(out)
	return out
}

// --- the constructs ---------------------------------------------------------

func TestCardinalityConstructs(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern>
      <e>req</e>
      <o><e>opt</e></o>
      <any><e>rep</e></any>
    </pattern>
    <ele><name>req</name><pattern>text</pattern></ele>
    <ele><name>opt</name><pattern>text</pattern></ele>
    <ele><name>rep</name><pattern>text</pattern></ele>
    <response><pattern/></response>`)

	want := map[string]Card{"req": CardRequired, "opt": CardOptional, "rep": CardRepeated}
	for name, card := range want {
		if got := find(t, cmd.Request, name).Card; got != card {
			t.Errorf("%s: got %v, want %v", name, got, card)
		}
	}
}

func TestNestedOptionalOfRepeatedStaysRepeated(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern><o><any><e>a</e></any></o></pattern>
    <ele><name>a</name><pattern>text</pattern></ele>
    <response><pattern/></response>`)

	if got := find(t, cmd.Request, "a").Card; got != CardRepeated {
		t.Errorf("Got %v, want repeated: zero-or-more inside optional is still zero-or-more", got)
	}
}

func TestGroupIsTransparentAndCarriesAttributes(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern>
      <g>
        <attrib><name>gattr</name><type>text</type><required>1</required></attrib>
        <e>a</e>
        <o><e>b</e></o>
      </g>
    </pattern>
    <ele><name>a</name><pattern>text</pattern></ele>
    <ele><name>b</name><pattern>text</pattern></ele>
    <response><pattern/></response>`)

	if got := find(t, cmd.Request, "a").Card; got != CardRequired {
		t.Errorf("a: got %v, want required: <g> must not change cardinality", got)
	}
	if got := find(t, cmd.Request, "b").Card; got != CardOptional {
		t.Errorf("b: got %v, want optional", got)
	}
	ga := find(t, cmd.Request, "gattr")
	if ga.Kind != KindAttribute || ga.Card != CardRequired {
		t.Errorf("gattr: got %v/%v, want attr/required", ga.Kind, ga.Card)
	}
}

func TestChoiceMakesEveryBranchOptional(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern>
      <or>
        <e>a</e>
        <g><e>b</e></g>
      </or>
    </pattern>
    <ele><name>a</name><pattern>text</pattern></ele>
    <ele><name>b</name><pattern>text</pattern></ele>
    <response><pattern/></response>`)

	for _, name := range []string{"a", "b"} {
		n := find(t, cmd.Request, name)
		if n.Card != CardOptional {
			t.Errorf("%s: got %v, want optional inside <or>", name, n.Card)
		}
		if !n.InChoice {
			t.Errorf("%s: expected InChoice to be set", name)
		}
	}
}

func TestAttributeRequiredness(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern>
      <attrib><name>must</name><type>uuid</type><required>1</required></attrib>
      <attrib><name>may</name><type>text</type><required>0</required></attrib>
      <attrib><name>silent</name><type>text</type></attrib>
    </pattern>
    <response><pattern/></response>`)

	want := map[string]Card{"must": CardRequired, "may": CardOptional, "silent": CardOptional}
	for name, card := range want {
		n := find(t, cmd.Request, name)
		if n.Kind != KindAttribute {
			t.Errorf("%s: got kind %v, want attr", name, n.Kind)
		}
		if n.Card != card {
			t.Errorf("%s: got %v, want %v", name, n.Card, card)
		}
	}
	if got := find(t, cmd.Request, "must").Type; got != "uuid" {
		t.Errorf("must: got type %q, want uuid", got)
	}
}

func TestRecursiveEleScoping(t *testing.T) {
	// "type" is defined twice, at two depths, with different content. Local
	// scoping must give each occurrence the definition alongside it.
	cmd := parseOne(t, `
    <name>x</name>
    <pattern><e>resources</e><e>type</e></pattern>
    <ele>
      <name>resources</name>
      <pattern><any><e>resource</e></any><e>type</e></pattern>
      <ele>
        <name>resource</name>
        <pattern><attrib><name>id</name><type>uuid</type></attrib></pattern>
      </ele>
      <ele><name>type</name><pattern>text</pattern></ele>
    </ele>
    <ele><name>type</name><pattern><t>integer</t></pattern></ele>
    <response><pattern/></response>`)

	outer := find(t, cmd.Request, "type")
	if outer.Type != "integer" {
		t.Errorf("Outer type: got %q, want integer", outer.Type)
	}
	res := find(t, cmd.Request, "resources")
	inner := find(t, res, "type")
	if inner.Type != "text" {
		t.Errorf("Inner type: got %q, want text", inner.Type)
	}
	kid := find(t, res, "resource")
	if kid.Card != CardRepeated {
		t.Errorf("resource: got %v, want repeated", kid.Card)
	}
	if id := find(t, kid, "id"); id.Kind != KindAttribute {
		t.Errorf("resource/@id: got kind %v, want attr", id.Kind)
	}
}

func TestLeafTypeSpellings(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern><e>bare</e><e>tagged</e><e>enum</e><e>typedele</e></pattern>
    <ele><name>bare</name><pattern>boolean</pattern></ele>
    <ele><name>tagged</name><pattern><t>uuid</t></pattern></ele>
    <ele><name>enum</name><pattern><alts><alt>a</alt><alt>b</alt></alts></pattern></ele>
    <ele><name>typedele</name><type><alts><alt>User</alt></alts></type></ele>
    <response><pattern/></response>`)

	if got := find(t, cmd.Request, "bare").Type; got != "boolean" {
		t.Errorf("bare: got %q, want boolean", got)
	}
	if got := find(t, cmd.Request, "tagged").Type; got != "uuid" {
		t.Errorf("tagged: got %q, want uuid", got)
	}
	if got := find(t, cmd.Request, "enum").Alts; len(got) != 2 || got[0] != "a" {
		t.Errorf("enum: got alts %v, want [a b]", got)
	}
	if got := find(t, cmd.Request, "typedele").Alts; len(got) != 1 || got[0] != "User" {
		t.Errorf("typedele: got alts %v, want [User]", got)
	}
}

// --- mixed content -----------------------------------------------------

// TestChardataBareTextSpelling covers the DSL's first mixed-content
// spelling: bare identifier text sitting next to an <attrib>, the shape
// behind <text excerpt="...">actual text</text> in note/override.
func TestChardataBareTextSpelling(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern><e>text</e></pattern>
    <ele>
      <name>text</name>
      <pattern>
        text
        <attrib><name>excerpt</name><type>boolean</type></attrib>
      </pattern>
    </ele>
    <response><pattern/></response>`)

	n := find(t, cmd.Request, "text")
	if n.Type != "text" {
		t.Errorf("Got type %q, want text", n.Type)
	}
	if !n.Chardata {
		t.Error("Expected Chardata to be set for bare text beside an <attrib>")
	}
	find(t, n, "excerpt")
}

// TestChardataTaggedTypeSpelling covers the DSL's second mixed-content
// spelling: an explicit <t> alongside an <e> reference, the shape behind
// config/nvt_count and config/family_count in command_get_configs.xml,
// where a real example response shows bare integer text next to a
// <growing> child (e.g. <nvt_count>88005<growing>1</growing></nvt_count>).
// leafOf's T branch returns the type unconditionally, so unlike the bare-
// text spelling above, this one is never routed through bareChardata.
func TestChardataTaggedTypeSpelling(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern><e>nvt_count</e></pattern>
    <ele>
      <name>nvt_count</name>
      <pattern>
        <t>integer</t>
        <e>growing</e>
      </pattern>
      <ele><name>growing</name><pattern><t>boolean</t></pattern></ele>
    </ele>
    <response><pattern/></response>`)

	n := find(t, cmd.Request, "nvt_count")
	if n.Type != "integer" {
		t.Errorf("Got type %q, want integer", n.Type)
	}
	if !n.Chardata {
		t.Error("Expected Chardata to be set for <t> beside an <e> reference")
	}
	find(t, n, "growing")
}

// TestChardataDoesNotFireForSharedElementReference guards the boundary the
// tagged-type spelling introduces: a <t> can also name a shared top-level
// <element> rather than a scalar type. If that pattern also carries a local
// <e> reference of its own, n.Type is non-empty and Children is non-empty
// at the exact point fill() checks — indistinguishable from genuine mixed
// content unless the shared-element case is excluded. expandShared is about
// to replace n's content with the referenced element's own pattern, so
// there is no bare text here at all.
func TestChardataDoesNotFireForSharedElementReference(t *testing.T) {
	body := `<?xml version="1.0"?><protocol>
      <element>
        <name>nvt</name>
        <pattern><e>family</e></pattern>
        <ele><name>family</name><pattern>text</pattern></ele>
      </element>
      <command>
        <name>x</name>
        <pattern><e>a</e></pattern>
        <ele>
          <name>a</name>
          <pattern><t>nvt</t><e>extra</e></pattern>
          <ele><name>extra</name><pattern>text</pattern></ele>
        </ele>
        <response><pattern/></response>
      </command></protocol>`

	cmds, err := Parse([]byte(body), "test.xml")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	a := find(t, cmds[0].Request, "a")
	if a.Chardata {
		t.Error("A shared-element reference must not be marked Chardata")
	}
	find(t, a, "family")
	find(t, a, "extra")
}

func TestResponseIsParsedSeparately(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern><e>a</e></pattern>
    <ele><name>a</name><pattern>text</pattern></ele>
    <response>
      <pattern>
        <attrib><name>status</name><type>text</type><required>1</required></attrib>
        <e>role</e>
      </pattern>
      <ele><name>role</name><pattern>text</pattern></ele>
    </response>`)

	if cmd.Response.Name != "x_response" {
		t.Errorf("Got response root %q, want x_response", cmd.Response.Name)
	}
	find(t, cmd.Response, "role")
	find(t, cmd.Response, "status")
	// The request's <ele> must not leak into the response scope.
	for _, c := range cmd.Response.Children {
		if c.Name == "a" {
			t.Errorf("Request element leaked into the response scope")
		}
	}
}

func TestPathUsesFidelitySyntax(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern><e>a</e></pattern>
    <ele>
      <name>a</name>
      <pattern><attrib><name>id</name><type>uuid</type></attrib></pattern>
    </ele>
    <response><pattern/></response>`)

	a := find(t, cmd.Request, "a")
	if got := a.Path("/x"); got != "/x/a" {
		t.Errorf("Got %q, want /x/a", got)
	}
	if got := find(t, a, "id").Path("/x/a"); got != "/x/a/@id" {
		t.Errorf("Got %q, want /x/a/@id", got)
	}
}

// --- schema defects ---------------------------------------------------------

// Defect 1: <pattern><type>text</type></pattern> where <t>text</t> was meant.
func TestDefectTypeUsedForT(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern><e>a</e></pattern>
    <ele><name>a</name><pattern><type>text</type></pattern></ele>
    <response><pattern/></response>`)

	if got := find(t, cmd.Request, "a").Type; got != "text" {
		t.Errorf("Got type %q, want text: <type> must be read as <t>", got)
	}
	if !hasDefect(cmd, DefectTypeForT, "/x/a") {
		t.Errorf("Expected %s at /x/a, got %v", DefectTypeForT, cmd.Defects)
	}
}

// Defect 2: an <e> reference with no matching <ele> definition.
func TestDefectDanglingReference(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern><e>ghost</e><e>real</e></pattern>
    <ele><name>real</name><pattern>text</pattern></ele>
    <response><pattern/></response>`)

	// The parser must survive it and still produce the node, so the diff can
	// report the Go side honestly instead of the command vanishing.
	g := find(t, cmd.Request, "ghost")
	if g.Type != "" || len(g.Children) != 0 {
		t.Errorf("Expected an empty node for the dangling reference, got %+v", g)
	}
	find(t, cmd.Request, "real")
	if !hasDefect(cmd, DefectDanglingRef, "/x/ghost") {
		t.Errorf("Expected %s at /x/ghost, got %v", DefectDanglingRef, cmd.Defects)
	}
}

// Defect 4: a completely empty <pattern/>.
func TestDefectEmptyPattern(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern><e>hollow</e></pattern>
    <ele><name>hollow</name><pattern/></ele>
    <response><pattern/></response>`)

	h := find(t, cmd.Request, "hollow")
	if len(h.Children) != 0 || h.Type != "" {
		t.Errorf("Expected a contentless node, got %+v", h)
	}
	if !hasDefect(cmd, DefectEmptyPattern, "/x/hollow") {
		t.Errorf("Expected %s at /x/hollow, got %v", DefectEmptyPattern, cmd.Defects)
	}
	// The response's own <pattern/> is empty too and must be reported.
	if !hasDefect(cmd, DefectEmptyPattern, "/x_response") {
		t.Errorf("Expected %s at /x_response, got %v", DefectEmptyPattern, cmd.Defects)
	}
}

// Defect 3: a name defined twice in one scope. First definition wins.
func TestDefectDuplicateDefinition(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern><e>dup</e></pattern>
    <ele><name>dup</name><pattern><t>uuid</t></pattern></ele>
    <ele><name>dup</name><pattern><t>integer</t></pattern></ele>
    <response><pattern/></response>`)

	if got := find(t, cmd.Request, "dup").Type; got != "uuid" {
		t.Errorf("Got %q, want uuid: the first definition wins", got)
	}
	if !hasDefect(cmd, DefectDuplicateEle, "/x/dup") {
		t.Errorf("Expected %s, got %v", DefectDuplicateEle, cmd.Defects)
	}
}

// <ele><id> lets two definitions share one wire name: get_filters's second
// <filters> definition carries <id>filters2</id> so a distinct <e>filters2</e>
// reference can pick it out. This must resolve to two distinct nodes, both
// named after the shared wire name rather than the disambiguating token, and
// must not trip the id-less duplicate-definition or dangling-reference
// defects that would otherwise fire here.
func TestIDDisambiguatesSharedWireName(t *testing.T) {
	cmd := parseOne(t, `
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
    <response><pattern/></response>`)

	if len(cmd.Request.Children) != 2 {
		t.Fatalf("Got %d children, want 2: <id> must not collapse the two defs into one", len(cmd.Request.Children))
	}
	first, second := cmd.Request.Children[0], cmd.Request.Children[1]
	if first.Name != "x" || second.Name != "x" {
		t.Errorf("Got names %q/%q, want both %q: the wire name, not the disambiguating token", first.Name, second.Name, "x")
	}
	find(t, first, "a")
	find(t, second, "b")
	if hasDefect(cmd, DefectDuplicateEle, "/x/x") {
		t.Errorf("<id> disambiguation must not be flagged as a duplicate definition: %v", cmd.Defects)
	}
	if hasDefect(cmd, DefectDanglingRef, "/x/x2") {
		t.Errorf("<e>x2</e> must resolve through <id>, not be flagged dangling: %v", cmd.Defects)
	}
	// "x" and "x2" are distinct tokens as far as the seen set is concerned;
	// they must never be treated as the same reference repeated.
	if hasDefect(cmd, DefectDuplicateRef, "/x/x") || hasDefect(cmd, DefectDuplicateRef, "/x/x2") {
		t.Errorf("Distinct tokens aliased to the same wire name must not be flagged as a duplicate reference: %v", cmd.Defects)
	}
}

// The same <e> token referenced twice in one pattern: one definition, but two
// occurrences. Distinct from DefectDuplicateEle, which is about two
// colliding definitions.
func TestDefectDuplicateReference(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern><e>dup</e><e>dup</e></pattern>
    <ele><name>dup</name><pattern>text</pattern></ele>
    <response><pattern/></response>`)

	got := 0
	for _, c := range cmd.Request.Children {
		if c.Name == "dup" {
			got++
		}
	}
	if got != 1 {
		t.Errorf("Got %d children named dup, want 1: the repeated reference must not double the node", got)
	}
	if !hasDefect(cmd, DefectDuplicateRef, "/x/dup") {
		t.Errorf("Expected %s at /x/dup, got %v", DefectDuplicateRef, cmd.Defects)
	}
}

// seen is shared across one fill call's whole nested walk rather than reset
// per wrapper, so a token repeated inside an <o>/<any>/<g>/<or> must be caught
// exactly like one repeated at the top level.
func TestDefectDuplicateReferenceAcrossWrappers(t *testing.T) {
	wrappers := []struct{ open, close string }{
		{"<o>", "</o>"},
		{"<any>", "</any>"},
		{"<g>", "</g>"},
		{"<or>", "</or>"},
	}
	for _, w := range wrappers {
		t.Run(w.open, func(t *testing.T) {
			cmd := parseOne(t, `
        <name>x</name>
        <pattern><e>dup</e>`+w.open+`<e>dup</e>`+w.close+`</pattern>
        <ele><name>dup</name><pattern>text</pattern></ele>
        <response><pattern/></response>`)

			got := 0
			for _, c := range cmd.Request.Children {
				if c.Name == "dup" {
					got++
				}
			}
			if got != 1 {
				t.Errorf("Got %d children named dup, want 1", got)
			}
			if !hasDefect(cmd, DefectDuplicateRef, "/x/dup") {
				t.Errorf("Expected %s at /x/dup, got %v", DefectDuplicateRef, cmd.Defects)
			}
		})
	}
}

// seen is local to one fill call: two elements' own patterns may each
// reference the same shared child name without colliding, since expanding
// into an <ele>'s own <pattern> starts a fresh seen set.
func TestDuplicateRefSeenResetsPerScope(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern><e>a</e><e>b</e></pattern>
    <ele>
      <name>a</name>
      <pattern><e>shared</e></pattern>
      <ele><name>shared</name><pattern>text</pattern></ele>
    </ele>
    <ele>
      <name>b</name>
      <pattern><e>shared</e></pattern>
      <ele><name>shared</name><pattern>text</pattern></ele>
    </ele>
    <response><pattern/></response>`)

	find(t, find(t, cmd.Request, "a"), "shared")
	find(t, find(t, cmd.Request, "b"), "shared")
	for _, d := range cmd.Defects {
		if d.Class == DefectDuplicateRef {
			t.Errorf("A name reused across separate elements' own scopes must not be flagged a duplicate reference: %v", d)
		}
	}
}

// Defect 5: undefined type names resolved through the alias table.
func TestDefectUndefinedTypeAliased(t *testing.T) {
	cmd := parseOne(t, `
    <name>x</name>
    <pattern><e>a</e><e>b</e></pattern>
    <ele><name>a</name><pattern><t>UUID</t></pattern></ele>
    <ele><name>b</name><pattern><t>nonesuch</t></pattern></ele>
    <response><pattern/></response>`)

	if !hasDefect(cmd, DefectUndefinedType, "/x/a") || !hasDefect(cmd, DefectUndefinedType, "/x/b") {
		t.Fatalf("Expected both to be reported undefined, got %v", cmd.Defects)
	}
	if Resolve("UUID") != "uuid" || Resolve("real") != "decimal" || Resolve("number") != "integer" {
		t.Errorf("Alias table is wrong")
	}
	if Resolve("nonesuch") != "nonesuch" {
		t.Errorf("Unaliased names must pass through unchanged")
	}
	for _, d := range cmd.Defects {
		if d.Path == "/x/a" && !strings.Contains(d.Detail, "aliased to uuid") {
			t.Errorf("Expected the alias to be named in the detail, got %q", d.Detail)
		}
	}
}

func TestSharedElementIsExpanded(t *testing.T) {
	// A leaf type naming a top-level <element> is a reference to it.
	body := `<?xml version="1.0"?><protocol>
      <type><name>uuid</name><pattern>text</pattern></type>
      <element>
        <name>nvt</name>
        <pattern>
          <attrib><name>oid</name><type>text</type><required>1</required></attrib>
          <e>family</e>
        </pattern>
        <ele><name>family</name><pattern>text</pattern></ele>
      </element>
      <command>
        <name>x</name>
        <pattern><e>a</e></pattern>
        <ele><name>a</name><type>nvt</type></ele>
        <response>
          <pattern><attrib><name>status</name><type>text</type><required>1</required></attrib></pattern>
        </response>
      </command></protocol>`

	cmds, err := Parse([]byte(body), "test.xml")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	a := find(t, cmds[0].Request, "a")
	if a.Type != "" {
		t.Errorf("Expected the shared type to be replaced by its content, got %q", a.Type)
	}
	find(t, a, "family")
	find(t, a, "oid")
	if got := defectClasses(cmds[0]); len(got) != 0 {
		t.Errorf("A resolved shared element is not a defect, got %v", got)
	}
}

func TestSharedElementCycleTerminates(t *testing.T) {
	body := `<?xml version="1.0"?><protocol>
      <element>
        <name>loop</name>
        <pattern><e>self</e></pattern>
        <ele><name>self</name><type>loop</type></ele>
      </element>
      <command>
        <name>x</name>
        <pattern><e>a</e></pattern>
        <ele><name>a</name><type>loop</type></ele>
        <response><pattern/></response>
      </command></protocol>`

	cmds, err := Parse([]byte(body), "test.xml")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Terminating at all is the assertion; the guard stops the second visit.
	a := find(t, cmds[0].Request, "a")
	find(t, a, "self")
}

func TestParseRejectsMalformedXML(t *testing.T) {
	if _, err := Parse([]byte(`<protocol><command></protocol>`), "bad.xml"); err == nil {
		t.Fatalf("Expected an error for mismatched tags")
	}
}

func hasDefect(c Command, class, path string) bool {
	for _, d := range c.Defects {
		if d.Class == class && d.Path == path {
			return true
		}
	}
	return false
}
