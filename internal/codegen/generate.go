package codegen

import (
	"fmt"
	"go/format"
	"sort"
	"strings"

	"github.com/brennoo/go-gmp/internal/schema"
)

// targetPackage is the package the generated files are written for. Support
// types named in TypeOverrides are spelled "commands.X" there so the table
// reads the same as the import a caller would write; inside the target package
// itself the qualifier has to go, and stripping it here is the one place that
// knows about it. The compile-check test re-adds it through type aliases.
const targetPackage = "commands"

// sharedFileSuffix names the file that holds nested types used by more than
// one command in the family. Emitting them once, rather than per command, is
// the whole point of the structural dedup pass: without it "nvt" would be
// declared by both create_note.go and modify_note.go and neither would build.
const sharedFileSuffix = "_shared.go"

// Generate renders Go source for every command in cmds.
//
// It returns a filename -> gofmt'd source map. The output is a pure function
// of (cmds, ov): no clock, no filesystem, no map iteration order leaks into
// it, and nothing here reads vendor/XML. That is what lets the golden test in
// generate_test.go run in CI with only the committed JSON snapshot.
func Generate(cmds []schema.Command, ov map[string]Override) (map[string]string, error) {
	if len(cmds) == 0 {
		return nil, fmt.Errorf("codegen: no commands to generate")
	}
	g := &gen{ov: ov, byHash: map[string]*structType{}, taken: map[string]bool{}}

	// Sorting decouples the output from the caller's slice order, so a
	// reordered snapshot list cannot silently renumber disambiguated names.
	g.cmds = append(g.cmds, cmds...)
	sort.Slice(g.cmds, func(i, j int) bool { return g.cmds[i].Name < g.cmds[j].Name })

	for i := range g.cmds {
		if g.cmds[i].Name == "" {
			return nil, fmt.Errorf("codegen: command %d has no name", i)
		}
		if g.cmds[i].Request == nil || g.cmds[i].Response == nil {
			return nil, fmt.Errorf("codegen: %s: snapshot is missing a request or response shape", g.cmds[i].Name)
		}
	}

	g.family = familyPrefix(g.cmds)
	g.collect()
	g.name()
	return g.emit()
}

// GenerateAll renders every family in m, one Generate call per family, and
// merges the results into a single filename -> source map.
//
// Per-family rather than one global call is a deliberate choice. Generate's
// structural dedup, its type-name disambiguation ladder and its shared-file
// decision are all scoped to the commands it is handed; giving it exactly one
// family's commands is what guarantees that adding an unrelated family to the
// manifest can never rename a type or move a struct in an existing one. A
// single global call would have to reproduce that isolation internally, and
// nothing in the output needs to be shared across families to make it work.
//
// Stated as the rule it is: structurally-equal shapes are deduplicated
// family-locally and never globally. Two families that happen to describe the
// same <nvt> shape each get their own type for it, in their own shared file,
// under their own family-qualified name. The duplication is the price of the
// property that matters more — a family's generated names, files and bytes
// depend on that family's commands and nothing else — because the alternative
// makes every future manifest addition a potential rename of already-committed,
// already-called types across commands/, client/ and internal/gmptest.
// Two tests hold the line: TestGenerateAllIsPerFamily (a synthetic twin family
// whose shapes are identical to note's, the adversarial case) and
// TestRealFamiliesGenerateIndependently (the real note and override families,
// each generated alone and compared byte-for-byte against its slice of the
// combined run).
//
// cmds may contain commands in any order and must cover every command the
// manifest names; a command not claimed by any family is an error, since it
// would otherwise be silently dropped.
func GenerateAll(m Manifest, cmds []schema.Command, ov map[string]Override) (map[string]string, error) {
	byName := make(map[string]schema.Command, len(cmds))
	for _, c := range cmds {
		byName[c.Name] = c
	}

	out := map[string]string{}
	from := map[string]string{} // filename -> the family that emitted it
	for _, f := range m.Families {
		fam := make([]schema.Command, 0, len(f.Commands))
		for _, name := range f.Commands {
			c, ok := byName[name]
			if !ok {
				return nil, fmt.Errorf("codegen: family %q: no snapshot loaded for %s", f.Name, name)
			}
			fam = append(fam, c)
		}
		files, err := Generate(fam, ov)
		if err != nil {
			return nil, fmt.Errorf("codegen: family %q: %w", f.Name, err)
		}
		if err := mergeFamilyFiles(out, from, f, files); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// mergeFamilyFiles folds one family's Generate output into out.
//
// Both checks it makes should be unreachable by construction — a command
// belongs to one family (LoadManifest enforces it) and the shared file is named
// after the family — but "should be unreachable" is exactly the property worth
// asserting: silently overwriting one family's file with another's would
// produce a plausible-looking tree that is missing types.
func mergeFamilyFiles(out, from map[string]string, f Family, files map[string]string) error {
	allowed := map[string]bool{f.SharedFile(): true}
	for _, name := range f.Commands {
		allowed[CommandFile(name)] = true
	}
	for name, src := range files {
		if !allowed[name] {
			return fmt.Errorf("codegen: family %q emitted %s, which the manifest's naming convention does not "+
				"account for (expected %s or one of the family's command files); is the family's name wrong?",
				f.Name, name, f.SharedFile())
		}
		if other, ok := from[name]; ok {
			return fmt.Errorf("codegen: families %q and %q both emit %s", other, f.Name, name)
		}
		from[name] = f.Name
		out[name] = src
	}
	return nil
}

// --- collection -------------------------------------------------------------

// structType is one generated Go struct: a nested element shape, plus every
// place in the family that shape was reached from.
type structType struct {
	node *schema.Node // the first occurrence; all others are structurally equal
	uses []use

	goName string
	file   string
}

// use is one path a shape was reached by, e.g. get_notes /note/result/nvt.
type use struct {
	command string
	path    []string // ancestor element names, outermost first, excluding the node
}

type gen struct {
	cmds   []schema.Command
	ov     map[string]Override
	family string

	byHash map[string]*structType
	types  []*structType   // first-occurrence order
	taken  map[string]bool // every Go type name already handed out
}

// collect walks every command's request and response, recording one entry per
// distinct nested struct shape.
func (g *gen) collect() {
	for i := range g.cmds {
		c := &g.cmds[i]
		g.walk(c.Name, c.Request, nil)
		g.walk(c.Name, c.Response, nil)
	}
}

func (g *gen) walk(command string, n *schema.Node, path []string) {
	for _, kid := range n.Children {
		if len(kid.Children) == 0 {
			continue // a scalar leaf needs no type of its own
		}
		h := typeHash(kid)
		st, ok := g.byHash[h]
		if !ok {
			st = &structType{node: kid}
			g.byHash[h] = st
			g.types = append(g.types, st)
		}
		st.uses = append(st.uses, use{command: command, path: append([]string(nil), path...)})
		// A fresh slice per branch: sharing the backing array would let a
		// later sibling overwrite an ancestor path already recorded above.
		down := make([]string, len(path), len(path)+1)
		copy(down, path)
		g.walk(command, kid, append(down, kid.Name))
	}
}

// typeHash identifies a generated struct by the shape it would have.
//
// The node's own name is part of the hash: <result id=""/> and <task id=""/>
// are the same shape but different concepts, and collapsing them would produce
// a type no reader could name. The node's own Card and Kind are deliberately
// *not* part of it: those describe the field that holds the struct, not the
// struct, which is why create_note's required <nvt> and modify_note's optional
// <nvt> share one type.
func typeHash(n *schema.Node) string {
	var b strings.Builder
	writeHash(&b, n)
	return b.String()
}

func writeHash(b *strings.Builder, n *schema.Node) {
	b.WriteString(n.Name)
	if n.Chardata {
		// A node with mixed content gets a Value field the same shape's
		// pure-element twin does not; without this marker two commands whose
		// children are identical but which disagree on bare chardata would
		// collapse into one struct and one of them would be wrong.
		b.WriteString("#text:")
		b.WriteString(n.Type)
	}
	b.WriteByte('{')
	for _, c := range n.Children {
		fmt.Fprintf(b, "%s|%s|%s|%s|%s|", c.Name, c.Kind, c.Card, c.Type, strings.Join(c.Alts, "\x1f"))
		writeHash(b, c)
		b.WriteByte(';')
	}
	b.WriteByte('}')
}

// --- naming -----------------------------------------------------------------

// familyPrefix derives the noun the commands are about: create_note,
// modify_note and get_notes are all about a "note". Dropping the leading verb
// and the plural is enough for GMP, whose command names are uniformly
// <verb>_<noun>. If the commands disagree there is no family, and every nested
// type falls back to being named after its own command.
func familyPrefix(cmds []schema.Command) string {
	noun := ""
	for _, c := range cmds {
		parts := strings.SplitN(c.Name, "_", 2)
		if len(parts) != 2 || parts[1] == "" {
			return ""
		}
		n := strings.TrimSuffix(parts[1], "s")
		if noun == "" {
			noun = n
		} else if noun != n {
			return ""
		}
	}
	return goName(noun)
}

// name assigns every collected shape a unique Go type name, in
// first-occurrence order so the assignment is stable.
func (g *gen) name() {
	for i := range g.cmds {
		g.taken[goName(g.cmds[i].Name)] = true
		g.taken[goName(g.cmds[i].Name)+"Response"] = true
	}
	for _, st := range g.types {
		st.goName = g.pick(st)
		g.taken[st.goName] = true
		st.file = g.fileFor(st)
	}
}

// pick walks a ladder of increasingly qualified candidates and takes the first
// one free: the bare element name, then the element name prefixed by one more
// ancestor at a time, then the same qualified by the command, then a counter.
// Short names win when nothing contends for them, and the collisions that do
// happen (get_notes has three structurally different <nvt>s) are resolved with
// the context that actually distinguishes them.
func (g *gen) pick(st *structType) string {
	first := st.uses[0]
	for _, c := range g.candidates(st, first) {
		if !g.taken[c] {
			return c
		}
	}
	base := g.qualify(goName(st.node.Name))
	for i := 2; ; i++ {
		if c := fmt.Sprintf("%s%d", base, i); !g.taken[c] {
			return c
		}
	}
}

func (g *gen) candidates(st *structType, first use) []string {
	var out []string
	add := func(s string) {
		for _, seen := range out {
			if seen == s {
				return
			}
		}
		out = append(out, s)
	}

	self := goName(st.node.Name)
	// Family-qualified, from the bare name outward through the ancestors.
	add(g.qualify(self))
	for i := len(first.path) - 1; i >= 0; i-- {
		add(g.qualify(goName(strings.Join(first.path[i:], "_")) + self))
	}
	// Then the same ladder qualified by the command that first reached it.
	cmd := goName(first.command)
	add(cmd + self)
	for i := len(first.path) - 1; i >= 0; i-- {
		add(cmd + goName(strings.Join(first.path[i:], "_")) + self)
	}
	return out
}

// qualify prefixes a name with the family noun, unless it already starts with
// it: <note_count> is NoteCount, not NoteNoteCount.
func (g *gen) qualify(s string) string {
	if g.family == "" || strings.HasPrefix(s, g.family) {
		return s
	}
	return g.family + s
}

// fileFor puts a shape used by one command in that command's file, and a shape
// reached from more than one command in the family's shared file.
func (g *gen) fileFor(st *structType) string {
	first := st.uses[0].command
	for _, u := range st.uses[1:] {
		if u.command != first {
			base := g.family
			if base == "" {
				base = "commands"
			}
			return strings.ToLower(base) + sharedFileSuffix
		}
	}
	return first + ".go"
}

// --- Go type mapping ---------------------------------------------------------

// scalarTypes maps schema type names onto Go types for the leaves the
// mechanical rule can decide on its own. Anything absent is a string: the DSL's
// remaining leaf types (uuid, oid, name, status, threat, compliance_status,
// ...) are all xsd:token or xsd:string patterns, and so is an unnamed leaf.
//
// severity is here rather than in TypeOverrides because it is not an override
// of anything: its declared pattern only ever matches a decimal number, so
// float64 *is* the mechanical reading of it.
var scalarTypes = map[string]string{
	"text":     "string",
	"integer":  "int",
	"decimal":  "float64",
	"severity": "float64",
}

// pointerScalars are the Go types whose zero value is a legal, meaningful wire
// value, so an optional field of that type has to be a pointer for "absent"
// and "zero" to stay distinguishable. This is T-5's rule. A string is not one
// of them: in this protocol an empty string and an absent element mean the
// same thing, which is why the hand-written code uses `string` + omitempty.
var pointerScalars = map[string]bool{"int": true, "float64": true, "GMPBool": true}

// goTypeOf resolves a node's Go type name, unqualified for targetPackage.
func (g *gen) goTypeOf(n *schema.Node) string {
	if len(n.Children) > 0 {
		return g.byHash[typeHash(n)].goName
	}
	return g.scalarGoType(n.Type)
}

// scalarGoType is goTypeOf's leaf-resolution rule, factored out so
// emitChardata can share it directly: a chardata node's own Type is always a
// lexical leaf type (e.g. "text", "base64"), never a nested shape, so it has
// no use for goTypeOf's structural branch.
func (g *gen) scalarGoType(t string) string {
	t = schema.Resolve(t)
	if o, ok := g.ov[t]; ok {
		return strings.TrimPrefix(o.GoType, targetPackage+".")
	}
	if s, ok := scalarTypes[t]; ok {
		return s
	}
	return "string"
}

type field struct {
	name   string
	goType string
	tag    string
}

// fieldOf applies the DSL -> Go cardinality mapping.
func (g *gen) fieldOf(n *schema.Node, used map[string]bool) field {
	base := g.goTypeOf(n)
	goType, opts := base, ""

	switch n.Card {
	case schema.CardRepeated:
		goType = "[]" + base
		opts = ",omitempty"
	case schema.CardOptional:
		if len(n.Children) > 0 || pointerScalars[base] {
			goType = "*" + base
		}
		opts = ",omitempty"
	case schema.CardRequired:
		// A required item is present by definition: no pointer, no omitempty.
	}

	tag := n.Name
	if n.Kind == schema.KindAttribute {
		tag += ",attr"
	}
	tag += opts

	return field{name: uniqueName(goName(n.Name), used), goType: goType, tag: tag}
}

// emitChardata declares the field for a mixed-content node's bare text, if it
// has any. It is called before emitFields so the chardata field reads first,
// matching the hand-written override types this generator is meant to match
// (e.g. commands.OverrideText's Text field precedes its Excerpt attribute).
func (g *gen) emitChardata(b *strings.Builder, n *schema.Node, used map[string]bool) {
	if !n.Chardata {
		return
	}
	name := uniqueName("Value", used)
	fmt.Fprintf(b, "%s %s `xml:\",chardata\"`\n", name, g.scalarGoType(n.Type))
}

func uniqueName(name string, used map[string]bool) string {
	out := name
	for i := 2; used[out]; i++ {
		out = fmt.Sprintf("%s%d", name, i)
	}
	used[out] = true
	return out
}

// --- emission ---------------------------------------------------------------

func (g *gen) emit() (map[string]string, error) {
	files := map[string]*strings.Builder{}
	needsXML := map[string]bool{}
	get := func(name string) *strings.Builder {
		b, ok := files[name]
		if !ok {
			b = &strings.Builder{}
			files[name] = b
		}
		return b
	}

	for i := range g.cmds {
		c := &g.cmds[i]
		file := c.Name + ".go"
		needsXML[file] = true
		b := get(file)
		g.emitRequest(b, c)
		g.emitResponse(b, c)
	}
	for _, st := range g.types {
		g.emitStruct(get(st.file), st)
	}

	out := make(map[string]string, len(files))
	for name, b := range files {
		var src strings.Builder
		src.WriteString(GeneratedMarker + "\n\n")
		src.WriteString("package " + targetPackage + "\n\n")
		if needsXML[name] {
			src.WriteString("import \"encoding/xml\"\n\n")
		}
		src.WriteString(b.String())

		formatted, err := format.Source([]byte(src.String()))
		if err != nil {
			return nil, fmt.Errorf("codegen: %s: generated source does not parse: %w", name, err)
		}
		out[name] = string(formatted)
	}
	return out, nil
}

func (g *gen) emitRequest(b *strings.Builder, c *schema.Command) {
	name := goName(c.Name)
	fmt.Fprintf(b, "// %s is a %s command request. See the GMP protocol docs for details.\n", name, c.Name)
	fmt.Fprintf(b, "type %s struct {\n", name)
	used := map[string]bool{"XMLName": true}
	fmt.Fprintf(b, "XMLName xml.Name `xml:%q`\n", c.Name)
	g.emitChardata(b, c.Request, used)
	g.emitFields(b, c.Request.Children, used, nil)
	b.WriteString("}\n\n")
}

func (g *gen) emitResponse(b *strings.Builder, c *schema.Command) {
	name := goName(c.Name) + "Response"
	root := c.Response.Name
	if root == "" {
		root = c.Name + "_response"
	}
	fmt.Fprintf(b, "// %s is a %s command response. See the GMP protocol docs for details.\n", name, c.Name)
	fmt.Fprintf(b, "type %s struct {\n", name)
	used := map[string]bool{"XMLName": true, "Status": true, "StatusText": true}
	fmt.Fprintf(b, "XMLName xml.Name `xml:%q`\n", root)
	// status/status_text are gvmd's universal response envelope, declared by
	// every command's schema and modelled identically by all 122 hand-written
	// response types. They are emitted first, from here rather than from the
	// schema children, so their order and spelling cannot drift per command.
	b.WriteString("Status string `xml:\"status,attr\"`\n")
	b.WriteString("StatusText string `xml:\"status_text,attr\"`\n")
	g.emitChardata(b, c.Response, used)
	g.emitFields(b, c.Response.Children, used, envelopeAttrs)
	b.WriteString("}\n\n")

	fmt.Fprintf(b, "// GetStatus returns the status and status text from the response.\n")
	fmt.Fprintf(b, "func (r *%s) GetStatus() (string, string) {\n\treturn r.Status, r.StatusText\n}\n\n", name)
}

// envelopeAttrs are the response attributes emitRequest/emitResponse already
// declare by hand, and must not be declared a second time from the schema.
var envelopeAttrs = map[string]bool{"status": true, "status_text": true}

func (g *gen) emitStruct(b *strings.Builder, st *structType) {
	fmt.Fprintf(b, "// %s is the <%s> element of %s. See the GMP protocol docs for details.\n",
		st.goName, st.node.Name, g.ownerOf(st))
	fmt.Fprintf(b, "type %s struct {\n", st.goName)
	used := map[string]bool{}
	g.emitChardata(b, st.node, used)
	g.emitFields(b, st.node.Children, used, nil)
	b.WriteString("}\n\n")
}

// ownerOf describes, for the doc comment, which commands reach a shape.
func (g *gen) ownerOf(st *structType) string {
	var names []string
	seen := map[string]bool{}
	for _, u := range st.uses {
		if !seen[u.command] {
			seen[u.command] = true
			names = append(names, u.command)
		}
	}
	sort.Strings(names)
	if len(names) == 1 {
		return names[0]
	}
	return strings.Join(names[:len(names)-1], ", ") + " and " + names[len(names)-1]
}

func (g *gen) emitFields(b *strings.Builder, kids []*schema.Node, used, skipAttrs map[string]bool) {
	for _, kid := range kids {
		if kid.Kind == schema.KindAttribute && skipAttrs[kid.Name] {
			continue
		}
		f := g.fieldOf(kid, used)
		fmt.Fprintf(b, "%s %s `xml:%q`\n", f.name, f.goType, f.tag)
	}
}

// --- identifiers -------------------------------------------------------------

// initialisms are the acronyms and other fixed-casing exceptions the
// hand-written commands/ package already spells this way — most are full-caps
// acronyms, but "icalendar" is the one mixed-case exception (RFC 5545's own
// spelling, not an acronym). Keeping the generator's spelling identical is
// what lets generated types drop in beside them without renaming every call
// site.
var initialisms = map[string]string{
	"api":       "API",
	"cpe":       "CPE",
	"cve":       "CVE",
	"cvss":      "CVSS",
	"dns":       "DNS",
	"http":      "HTTP",
	"icalendar": "ICalendar",
	"icmp":      "ICMP",
	"id":        "ID",
	"ip":        "IP",
	"nvt":       "NVT",
	"oid":       "OID",
	"os":        "OS",
	"qod":       "QOD",
	"smb":       "SMB",
	"snmp":      "SNMP",
	"ssh":       "SSH",
	"tcp":       "TCP",
	"tls":       "TLS",
	"udp":       "UDP",
	"url":       "URL",
	"uuid":      "UUID",
	"xml":       "XML",
}

// goName turns a schema element or command name into an exported Go
// identifier: snake_case to CamelCase, with known acronyms capitalised.
func goName(s string) string {
	var b strings.Builder
	for _, part := range strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == '.' || r == ':' || r == ' '
	}) {
		lower := strings.ToLower(part)
		if up, ok := initialisms[lower]; ok {
			b.WriteString(up)
			continue
		}
		b.WriteString(strings.ToUpper(lower[:1]) + lower[1:])
	}
	if b.Len() == 0 {
		return "Field"
	}
	return b.String()
}
