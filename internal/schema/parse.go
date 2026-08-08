// Package schema reads the GMP protocol schema DSL and compares the model it
// declares against the hand-written structs in commands/.
//
// It is a schema oracle. internal/gmptest is a runtime oracle: it round-trips
// real documents through the structs and records what is lost. This package
// never looks at a document. It reads what the schema says a command's shape
// is and reports where the Go types disagree, which finds fields no fixture
// happens to exercise.
//
// The input is vendor/XML/commands/command_*.xml, which is gitignored and
// AGPL-licensed. Nothing here vendors, copies or emits schema prose: only
// element names, cardinalities and type names cross the boundary, and only
// into a report. See docs/review-2026-08.md, T-9.
package schema

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Card is how many times an item may appear. The DSL has exactly three
// cardinalities: there is no "one or more" and no minOccurs/maxOccurs.
type Card uint8

// The three cardinalities the DSL can express.
const (
	CardRequired Card = iota // <e>x</e>
	CardOptional             // <o><e>x</e></o>, or an <attrib> without <required>1</required>
	CardRepeated             // <any><e>x</e></any>
)

func (c Card) String() string {
	switch c {
	case CardRequired:
		return "required"
	case CardOptional:
		return "optional"
	case CardRepeated:
		return "repeated"
	}
	return "?"
}

// Kind separates XML attributes from child elements. The DSL is unambiguous
// about this: <attrib> is an attribute, everything else is an element.
type Kind uint8

// The two node kinds.
const (
	KindElement Kind = iota
	KindAttribute
)

func (k Kind) String() string {
	if k == KindAttribute {
		return "attr"
	}
	return "elem"
}

// Node is one element or attribute in a command's declared shape.
type Node struct {
	Name     string
	Kind     Kind
	Card     Card
	InChoice bool // reached through <or>; its cardinality is advisory
	Type     string
	Alts     []string
	Children []*Node
}

// Path returns the node's location under parent, in the same syntax the
// fidelity reports in internal/gmptest use, so findings can be cross-checked
// against them by string equality.
func (n *Node) Path(parent string) string {
	if n.Kind == KindAttribute {
		return parent + "/@" + n.Name
	}
	return parent + "/" + n.Name
}

// Command is one command_*.xml file: the request shape, the response shape,
// and whatever was wrong with the file on the way through.
type Command struct {
	Name     string
	File     string
	Request  *Node
	Response *Node
	Defects  []Defect
}

// Defect is a place where the schema does not say what it means. These are
// upstream bugs, not transcription errors, and the report keeps them in their
// own category so the two are never confused.
type Defect struct {
	Class  string
	Path   string
	Detail string
}

// Defect classes, in the order they appear in the report.
const (
	DefectTypeForT      = "type-for-t"     // <pattern><type>x</type></pattern> where <t>x</t> was meant
	DefectDanglingRef   = "dangling-e"     // <e>x</e> with no <ele> defining x
	DefectDuplicateEle  = "duplicate-ele"  // two <ele> definitions of one name in one scope
	DefectEmptyPattern  = "empty-pattern"  // <pattern/>, declaring nothing at all
	DefectNoPattern     = "no-pattern"     // <ele> with neither <pattern> nor <type>
	DefectUndefinedType = "undefined-type" // a type name with no <type> definition
	DefectUnsupportedR  = "unsupported-r"  // <r>, a reference to a response, inside a pattern
	DefectRecursion     = "recursion-cut"  // a shared <element> that reaches itself
)

// typeAliases maps the type names the schema uses but never defines onto the
// ones it does. Every entry is an upstream typo or synonym, and each use is
// still counted as a DefectUndefinedType so the report says how often the
// table is load-bearing.
var typeAliases = map[string]string{
	"UUID":   "uuid",    // wrong case
	"real":   "decimal", // never defined
	"number": "integer", // never defined
}

// --- raw XML bindings -------------------------------------------------------
//
// These mirror the DSL rather than the IR. Order between different child tags
// is not preserved, which is fine: the IR is matched by name, and the DSL has
// no ordered content model that a Go struct could express anyway.

type rawProtocol struct {
	Types    []rawTypeDef `xml:"type"`
	Elements []rawEle     `xml:"element"`
	Commands []rawCommand `xml:"command"`
}

type rawTypeDef struct {
	Name    string      `xml:"name"`
	Pattern *rawPattern `xml:"pattern"`
}

type rawCommand struct {
	Name     string      `xml:"name"`
	Pattern  *rawPattern `xml:"pattern"`
	Defs     []rawEle    `xml:"ele"`
	Response *rawEle     `xml:"response"`
}

// rawEle covers <ele>, <element> and <response>: all three are "a pattern plus
// the definitions of the names that pattern refers to".
type rawEle struct {
	Name    string      `xml:"name"`
	Pattern *rawPattern `xml:"pattern"`
	Type    *rawTypeRef `xml:"type"`
	Defs    []rawEle    `xml:"ele"`
}

// rawTypeRef is <ele><type>, which is either a type name or an inline
// enumeration.
type rawTypeRef struct {
	Text string   `xml:",chardata"`
	Alts []string `xml:"alts>alt"`
}

// rawPattern is the content model. Bare text is a type name ("text",
// "boolean", "timezone"); <t> is the same thing spelled properly; <type> in
// this position is DefectTypeForT.
type rawPattern struct {
	Text    string      `xml:",chardata"`
	T       []string    `xml:"t"`
	Type    []string    `xml:"type"`
	Alts    []string    `xml:"alts>alt"`
	Attribs []rawAttrib `xml:"attrib"`
	rawGroup
}

// rawGroup is the part of a pattern that carries child references, shared by
// <pattern> and by the <o>/<any>/<g>/<or> wrappers.
type rawGroup struct {
	E   []string  `xml:"e"`
	R   []string  `xml:"r"`
	O   []rawNest `xml:"o"`
	Any []rawNest `xml:"any"`
	G   []rawNest `xml:"g"`
	Or  []rawNest `xml:"or"`
}

// rawNest is a wrapper that may itself carry attributes, which <g> and <or>
// both do.
type rawNest struct {
	Attribs []rawAttrib `xml:"attrib"`
	rawGroup
}

type rawAttrib struct {
	Name     string      `xml:"name"`
	Type     string      `xml:"type"`
	Required string      `xml:"required"`
	Alts     []string    `xml:"alts>alt"`
	Pattern  *rawPattern `xml:"pattern"`
}

// --- parsing ----------------------------------------------------------------

// Load parses every command_*.xml in dir, sorted by command name.
func Load(dir string) ([]Command, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "command_*.xml"))
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)

	var cmds []Command
	for _, p := range paths {
		got, err := LoadFile(p)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", filepath.Base(p), err)
		}
		cmds = append(cmds, got...)
	}
	sort.Slice(cmds, func(i, j int) bool { return cmds[i].Name < cmds[j].Name })
	return cmds, nil
}

// LoadFile parses one schema file. A file normally holds one <command>, but
// the shape allows more and nothing here assumes otherwise.
func LoadFile(path string) ([]Command, error) {
	data, err := os.ReadFile(path) //nolint:gosec // reading a caller-named schema file is the point
	if err != nil {
		return nil, err
	}
	return Parse(data, filepath.Base(path))
}

// Parse reads a <protocol> document. name is used only in reports.
func Parse(data []byte, name string) ([]Command, error) {
	var p rawProtocol
	if err := xml.Unmarshal(data, &p); err != nil {
		return nil, err
	}

	e := &env{types: map[string]*rawTypeDef{}, elements: map[string]*rawEle{}}
	for i := range p.Types {
		// A duplicate <type> is an upstream wart ("levels" is defined twice);
		// first definition wins, which is what a reader would do.
		if n := trim(p.Types[i].Name); n != "" && e.types[n] == nil {
			e.types[n] = &p.Types[i]
		}
	}
	for i := range p.Elements {
		if n := trim(p.Elements[i].Name); n != "" && e.elements[n] == nil {
			e.elements[n] = &p.Elements[i]
		}
	}

	out := make([]Command, 0, len(p.Commands))
	for i := range p.Commands {
		out = append(out, e.command(&p.Commands[i], name))
	}
	return out, nil
}

// env is the per-document type and shared-element table.
type env struct {
	types    map[string]*rawTypeDef
	elements map[string]*rawEle
	defects  []Defect
	// expanding guards against a shared <element> that reaches itself.
	expanding map[string]bool
}

func (e *env) command(c *rawCommand, file string) Command {
	e.defects = nil
	e.expanding = map[string]bool{}
	name := trim(c.Name)

	req := &Node{Name: name}
	e.fill(req, c.Pattern, c.Defs, "/"+name)

	resp := &Node{Name: name + "_response"}
	if c.Response != nil {
		e.fill(resp, c.Response.Pattern, c.Response.Defs, "/"+resp.Name)
	}

	defects := e.defects
	sort.SliceStable(defects, func(i, j int) bool {
		if defects[i].Class != defects[j].Class {
			return defects[i].Class < defects[j].Class
		}
		return defects[i].Path < defects[j].Path
	})
	return Command{Name: name, File: file, Request: req, Response: resp, Defects: defects}
}

// fill populates n from pattern, resolving <e> references against defs.
//
// Scoping is strictly local: an <ele> defines a name for the pattern of its
// own enclosing node and nothing else. That is not an assumption — resolving
// with local scope only leaves two unresolved references across the whole
// corpus, and both are genuine upstream bugs.
func (e *env) fill(n *Node, pattern *rawPattern, defs []rawEle, path string) {
	if pattern == nil {
		return
	}
	scope := e.scope(defs, path)

	n.Type, n.Alts = e.leafOf(pattern, path)
	e.noteType(n.Type, path)

	for i := range pattern.Attribs {
		n.Children = append(n.Children, e.attrib(&pattern.Attribs[i], path))
	}
	e.refs(n, &pattern.rawGroup, scope, CardRequired, false, path)
	e.expandShared(n, path)

	if n.Type == "" && len(n.Children) == 0 && len(n.Alts) == 0 {
		e.defect(DefectEmptyPattern, path, "pattern declares no content")
	}
}

// expandShared pulls in a top-level <element> definition when a node's leaf
// type names one.
//
// The review claims those 16 definitions are never referenced by any command.
// In this revision twelve command sites do reference them, and they are the
// largest subtrees in the protocol: result, report, note, override.
func (e *env) expandShared(n *Node, path string) {
	def, ok := e.elements[n.Type]
	if !ok || e.expanding[n.Type] {
		return
	}
	e.expanding[n.Type] = true
	e.fill(n, def.Pattern, def.Defs, path)
	delete(e.expanding, n.Type)
	n.Type = ""
}

// scope indexes an element's own <ele> definitions, reporting redefinitions.
func (e *env) scope(defs []rawEle, path string) map[string]*rawEle {
	scope := make(map[string]*rawEle, len(defs))
	for i := range defs {
		nm := trim(defs[i].Name)
		if nm == "" {
			continue
		}
		if _, dup := scope[nm]; dup {
			e.defect(DefectDuplicateEle, path+"/"+nm, "redefined in the same scope")
			continue // first definition wins
		}
		scope[nm] = &defs[i]
	}
	return scope
}

// refs walks a pattern's child references, composing cardinality as it
// descends through the <o>/<any>/<g>/<or> wrappers.
func (e *env) refs(n *Node, g *rawGroup, scope map[string]*rawEle, card Card, choice bool, path string) {
	for _, ref := range g.E {
		e.child(n, trim(ref), scope, card, choice, path)
	}
	for range g.R {
		e.defect(DefectUnsupportedR, path, "<r> references a response from inside a pattern")
	}
	for i := range g.O {
		e.nested(n, &g.O[i], scope, weaken(card, CardOptional), choice, path)
	}
	for i := range g.Any {
		e.nested(n, &g.Any[i], scope, CardRepeated, choice, path)
	}
	for i := range g.G {
		// <g> is a group: it interleaves its members without changing them.
		e.nested(n, &g.G[i], scope, card, choice, path)
	}
	for i := range g.Or {
		// <or> is a choice: at most one branch is present, so every member is
		// optional however it is spelled inside the branch.
		e.nested(n, &g.Or[i], scope, weaken(card, CardOptional), true, path)
	}
}

func (e *env) nested(n *Node, w *rawNest, scope map[string]*rawEle, card Card, choice bool, path string) {
	for i := range w.Attribs {
		a := e.attrib(&w.Attribs[i], path)
		a.InChoice = a.InChoice || choice
		n.Children = append(n.Children, a)
	}
	e.refs(n, &w.rawGroup, scope, card, choice, path)
}

// weaken combines an outer cardinality with an inner one. Repeated absorbs
// everything: <o><any><e>x</e></any></o> is still zero-or-more.
func weaken(outer, inner Card) Card {
	switch {
	case outer == CardRepeated || inner == CardRepeated:
		return CardRepeated
	case outer == CardOptional || inner == CardOptional:
		return CardOptional
	}
	return CardRequired
}

func (e *env) child(n *Node, name string, scope map[string]*rawEle, card Card, choice bool, path string) {
	if name == "" {
		return
	}
	kid := &Node{Name: name, Card: card, InChoice: choice}
	n.Children = append(n.Children, kid)

	def, ok := scope[name]
	if !ok {
		// Local scope first, then the document's top-level <element> table.
		// help's response describes the DSL itself and refers to <e>c</e>,
		// <e>o</e>, <e>g</e> and friends, which are defined only there.
		shared, global := e.elements[name]
		switch {
		case global && e.expanding[name]:
			// help's response describes the DSL, and the DSL is recursive: a
			// <g> holds an <o> which holds a <g>. Cutting the walk here is not
			// a schema bug, but the subtree below the cut is unmodelled, so it
			// is recorded rather than dropped silently.
			kid.Type = name
			e.defect(DefectRecursion, kid.Path(path), name+" reaches itself; the walk is cut here")
		case global:
			e.expanding[name] = true
			e.fill(kid, shared.Pattern, shared.Defs, kid.Path(path))
			delete(e.expanding, name)
		default:
			e.defect(DefectDanglingRef, kid.Path(path), "no <ele> defines this name in scope")
		}
		return
	}
	switch {
	case def.Type != nil:
		kid.Type, kid.Alts = trim(def.Type.Text), trimAll(def.Type.Alts)
		e.noteType(kid.Type, kid.Path(path))
		e.expandShared(kid, kid.Path(path))
	case def.Pattern != nil:
		e.fill(kid, def.Pattern, def.Defs, kid.Path(path))
	default:
		e.defect(DefectNoPattern, kid.Path(path), "<ele> has neither <pattern> nor <type>")
	}
}

func (e *env) attrib(a *rawAttrib, path string) *Node {
	n := &Node{
		Name: trim(a.Name),
		Kind: KindAttribute,
		Card: CardOptional,
		Type: trim(a.Type),
		Alts: trimAll(a.Alts),
	}
	// An attribute is required only when it says so. Absent <required> reads
	// the same as <required>0</required>.
	if trim(a.Required) == "1" {
		n.Card = CardRequired
	}
	if n.Type == "" && a.Pattern != nil {
		n.Type, _ = e.leafOf(a.Pattern, path)
	}
	e.noteType(n.Type, n.Path(path))
	return n
}

// leafOf extracts a leaf type name from a pattern, tolerating the three
// spellings the schema actually uses.
func (e *env) leafOf(p *rawPattern, path string) (string, []string) {
	alts := trimAll(p.Alts)
	switch {
	case len(p.T) > 0:
		return trim(p.T[0]), alts
	case len(p.Type) > 0:
		// <pattern><type>text</type></pattern>. The generating XSL has no
		// branch for this and emits a literal ERROR; read directly it is
		// simply a misspelt <t>.
		e.defect(DefectTypeForT, path, "<type> used where <t> was meant: "+trim(p.Type[0]))
		return trim(p.Type[0]), alts
	}
	if txt := trim(p.Text); txt != "" && len(p.Attribs) == 0 && p.empty() {
		return txt, alts
	}
	return "", alts
}

func (g *rawGroup) empty() bool {
	return len(g.E) == 0 && len(g.R) == 0 && len(g.O) == 0 &&
		len(g.Any) == 0 && len(g.G) == 0 && len(g.Or) == 0
}

// noteType records a type name that the document never defines.
func (e *env) noteType(name, path string) {
	switch {
	case name == "" || name == "text":
		return
	case e.types[name] != nil, e.elements[name] != nil:
		return
	case strings.HasPrefix(name, "xsd:"):
		return // an inline RELAX NG pattern, not a name
	}
	if alias, ok := typeAliases[name]; ok {
		e.defect(DefectUndefinedType, path, name+" is undefined; aliased to "+alias)
		return
	}
	e.defect(DefectUndefinedType, path, name+" is undefined and has no alias")
}

func (e *env) defect(class, path, detail string) {
	e.defects = append(e.defects, Defect{Class: class, Path: path, Detail: detail})
}

// Resolve follows the alias table. Unknown names are returned unchanged.
func Resolve(name string) string {
	if a, ok := typeAliases[name]; ok {
		return a
	}
	return name
}

func trim(s string) string { return strings.TrimSpace(s) }

func trimAll(ss []string) []string {
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		if t := trim(s); t != "" {
			out = append(out, t)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
