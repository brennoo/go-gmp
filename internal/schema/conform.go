package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Finding categories, in report order.
const (
	CatMissing     = "MISSING"       // the schema declares it, no Go field carries it
	CatExtra       = "EXTRA"         // a Go field with nothing behind it in the schema
	CatCardinality = "CARDINALITY"   // repeated-vs-slice, or an optional whose absence Go cannot express
	CatKind        = "KIND"          // an attribute modelled as an element, or the reverse
	CatOmitEmpty   = "OMITEMPTY"     // omitempty where the schema says required, or missing where it says optional
	CatType        = "TYPE"          // the Go type cannot hold the schema's lexical space
	CatDefect      = "SCHEMA-DEFECT" // the schema is wrong; nothing is asked of the Go side
)

// Categories in the order the report prints them.
var Categories = []string{CatMissing, CatExtra, CatCardinality, CatKind, CatOmitEmpty, CatType, CatDefect}

// Finding is one disagreement between the schema and the structs.
type Finding struct {
	Command  string
	Side     string // "request" or "response"
	Category string
	Path     string
	Detail   string

	// Confirmed marks a finding the fidelity corpus independently reached:
	// the same path is DROPPED when a real document round-trips. Two oracles
	// built on different evidence agreeing is the strongest signal here.
	Confirmed bool

	// Suppressed findings are counted and explained but not treated as struct
	// bugs — see suppressRequestOverConstraint.
	Suppressed bool
	Reason     string
}

// Result is one command's comparison.
type Result struct {
	Command  string
	File     string
	Findings []Finding
	// Compared is false when no Go type exists for the command at all.
	Compared bool
}

// Conformant reports whether the command has no actionable findings. Schema
// defects and suppressed findings do not count against it: neither asks for a
// change to commands/.
func (r Result) Conformant() bool {
	for _, f := range r.Findings {
		if f.Category != CatDefect && !f.Suppressed {
			return false
		}
	}
	return true
}

// Comparer holds the corroborating evidence a diff is run against.
type Comparer struct {
	// Dropped is per-command the set of paths the fidelity corpus saw a real
	// document lose. Optional: an empty map just means no cross-references.
	Dropped map[string]map[string]bool
}

// Compare diffs one command's schema shape against its Go types. Either root
// may be zero, which is itself reported.
func (c *Comparer) Compare(cmd Command, req, resp GoRoot) Result {
	res := Result{Command: cmd.Name, File: cmd.File, Compared: true}

	for _, d := range cmd.Defects {
		res.Findings = append(res.Findings, Finding{
			Command:  cmd.Name,
			Side:     sideOf(cmd, d.Path),
			Category: CatDefect,
			Path:     d.Path,
			Detail:   d.Class + ": " + d.Detail,
		})
	}

	suspect := suspectPaths(cmd)
	if req.Root != nil {
		c.walk(&res, cmd, "request", cmd.Request, req.Root, "/"+cmd.Request.Name, suspect)
	} else {
		res.add(cmd, "request", CatMissing, "/"+schemaRootName(cmd.Request, cmd.Name),
			"no Go type reflects to a request root at all")
	}
	if resp.Root != nil {
		c.walk(&res, cmd, "response", cmd.Response, resp.Root, "/"+cmd.Response.Name, suspect)
	} else {
		res.add(cmd, "response", CatMissing, "/"+schemaRootName(cmd.Response, cmd.Name),
			"no Go type reflects to a response root at all")
	}

	sortFindings(res.Findings)
	return res
}

// schemaRootName names a missing-Go-root finding's path. The schema side is
// itself absent for the same command in some malformed cases, so it falls
// back to the command name rather than dereferencing a nil *Node.
func schemaRootName(n *Node, cmdName string) string {
	if n == nil {
		return cmdName
	}
	return n.Name
}

// suspectPaths collects the places where the schema declares nothing usable,
// so the diff can stay quiet there instead of blaming the Go side for not
// modelling a hole.
func suspectPaths(cmd Command) map[string]bool {
	out := map[string]bool{}
	for _, d := range cmd.Defects {
		switch d.Class {
		case DefectDanglingRef, DefectEmptyPattern, DefectNoPattern, DefectRecursion:
			out[d.Path] = true
		}
	}
	return out
}

func sideOf(cmd Command, path string) string {
	if strings.HasPrefix(path, "/"+cmd.Name+"_response") {
		return "response"
	}
	return "request"
}

// walk compares one schema node against one Go node and recurses.
//
// byName maps a (kind, name) key to a *slice* of Go nodes: get_filters's two
// same-wire-name <filters> fields (see internal/codegen's collision codec and
// gomodel.go's gmpwire tag) give schema.Children two same-key entries too,
// and each must be paired with its own Go field rather than both resolving
// to whichever one a plain map kept. Pairing is by occurrence order — the
// Nth schema child named X against the Nth Go field named X — which is sound
// because both sides preserve declaration order and gvmd's schema declares
// get_filters's two definitions in the same relative order the fields are
// generated in. matchedCount tracks, per key, how many leading Go occurrences
// are claimed; outside get_filters it's always 0 or 1, so this is a no-op
// generalization of the single-valued lookup it replaces.
func (c *Comparer) walk(res *Result, cmd Command, side string, s *Node, g *GoNode, path string, suspect map[string]bool) {
	byName := map[string][]*GoNode{}
	for _, k := range g.Children {
		kk := key(k.Kind, k.Name)
		byName[kk] = append(byName[kk], k)
	}
	matchedCount := map[string]int{}

	for _, sc := range s.Children {
		kp := sc.Path(path)
		kk := key(sc.Kind, sc.Name)
		if suspect[kp] {
			// The schema, not the struct, is what is wrong here — but a Go
			// field that would have paired with this occurrence must not
			// then look unpaired and be reported EXTRA on top of it.
			matchedCount[kk]++
			continue
		}
		idx := matchedCount[kk]
		list := byName[kk]
		if idx >= len(list) {
			c.unmatchedSchema(res, cmd, side, sc, g, byName, matchedCount, kp)
			continue
		}
		matchedCount[kk] = idx + 1
		c.compareNode(res, cmd, side, sc, list[idx], kp, suspect)
	}

	occSeen := map[string]int{}
	for _, gc := range g.Children {
		kk := key(gc.Kind, gc.Name)
		idx := occSeen[kk]
		occSeen[kk] = idx + 1
		if idx < matchedCount[kk] || gc.Synthetic && len(gc.Children) > 0 {
			continue
		}
		f := Finding{Command: cmd.Name, Side: side, Category: CatExtra, Path: gc.Path(path), Detail: fmt.Sprintf(
			"Go field %s %s has no counterpart in the schema", gc.GoName, gc.GoType)}
		if gc.OverrideReason != "" {
			f.Suppressed = true
			f.Reason = gc.OverrideReason
		}
		res.Findings = append(res.Findings, f)
	}
}

// unmatchedSchema handles a schema child with no same-kind Go field: either
// the kind is wrong, or nothing models it at all.
func (c *Comparer) unmatchedSchema(res *Result, cmd Command, side string, sc *Node, g *GoNode,
	byName map[string][]*GoNode, matchedCount map[string]int, kp string) {
	// A Go field of the opposite kind is a KIND finding only if nothing else
	// has claimed it. get_aggregates declares data_column twice — once as an
	// attribute and once as <any><e>data_column</e></any> — and the attribute
	// is modelled correctly; what is missing is the repeated element, not the
	// attribute spelled wrongly.
	other := otherKind(sc.Kind)
	otherKey := key(other, sc.Name)
	if list := byName[otherKey]; matchedCount[otherKey] < len(list) {
		idx := matchedCount[otherKey]
		matchedCount[otherKey] = idx + 1
		gc := list[idx]
		res.add(cmd, side, CatKind, kp, fmt.Sprintf(
			"schema declares %s, Go models it as %s (%s %s)", sc.Kind, other, gc.GoName, gc.GoType))
		return
	}
	detail := fmt.Sprintf("%s %s %s is not modelled", sc.Card, sc.Kind, sc.Name)
	if g.CatchAll {
		detail += "; the parent has a `,any` field, so it is captured but not addressable"
	}
	f := Finding{Command: cmd.Name, Side: side, Category: CatMissing, Path: kp, Detail: detail}
	f.Confirmed = c.confirmed(cmd.Name, kp)
	res.Findings = append(res.Findings, f)
}

// compareNode checks a matched pair, then descends.
func (c *Comparer) compareNode(res *Result, cmd Command, side string, s *Node, g *GoNode, path string, suspect map[string]bool) {
	c.checkCardinality(res, cmd, side, s, g, path)
	checkOmitEmpty(res, cmd, side, s, g, path)
	checkType(res, cmd, side, s, g, path)

	switch {
	case len(s.Children) > 0 && len(g.Children) == 0:
		f := Finding{Command: cmd.Name, Side: side, Category: CatMissing, Path: path + "/*",
			Detail: fmt.Sprintf("schema declares %d children; Go models this as scalar %s",
				len(s.Children), g.GoType)}
		f.Confirmed = c.confirmed(cmd.Name, path)
		res.Findings = append(res.Findings, f)
	case len(s.Children) > 0:
		c.walk(res, cmd, side, s, g, path, suspect)
		c.checkChardata(res, cmd, side, s, g, path)
	case len(g.Children) > 0 && !g.Chardata:
		res.add(cmd, side, CatExtra, path+"/*", fmt.Sprintf(
			"schema declares a %q leaf; Go models %d children", displayType(s), len(g.Children)))
	}
}

// checkChardata compares mixed content: an element with real children that
// also carries bare text of its own (e.g. <text excerpt="...">the text
// </text>). This is the reverse of the len(g.Children) > 0 && !g.Chardata
// case above, which only fires when the schema is a pure leaf; here the
// schema side already has children, so the chardata comparison has to run on
// its own rather than falling out of the existing switch.
func (c *Comparer) checkChardata(res *Result, cmd Command, side string, s *Node, g *GoNode, path string) {
	switch {
	case s.Chardata && !g.Chardata:
		f := Finding{Command: cmd.Name, Side: side, Category: CatMissing, Path: path + "/#text",
			Detail: fmt.Sprintf("schema also declares character data (%s) here; Go %s has no `,chardata` field",
				displayType(s), g.GoType)}
		f.Confirmed = c.confirmed(cmd.Name, path)
		res.Findings = append(res.Findings, f)
	case !s.Chardata && g.Chardata:
		f := Finding{Command: cmd.Name, Side: side, Category: CatExtra, Path: path + "/#text", Detail: fmt.Sprintf(
			"Go %s has a `,chardata` field; schema declares no character data here", g.GoType)}
		if g.ChardataOverrideReason != "" {
			f.Suppressed = true
			f.Reason = g.ChardataOverrideReason
		}
		res.Findings = append(res.Findings, f)
	}
}

func (c *Comparer) checkCardinality(res *Result, cmd Command, side string, s *Node, g *GoNode, path string) {
	switch {
	case s.Card == CardRepeated && !g.Slice:
		res.add(cmd, side, CatCardinality, path, fmt.Sprintf(
			"schema says zero-or-more, Go has a single %s (%s)", g.GoType, g.GoName))
	case s.Card != CardRepeated && g.Slice:
		res.add(cmd, side, CatCardinality, path, fmt.Sprintf(
			"schema says %s, Go has a slice %s (%s)", s.Card, g.GoType, g.GoName))
	case s.Card == CardOptional && g.ZeroIsMeaningful():
		// The T-5 case. A bool or an int cannot distinguish "absent" from
		// "present and zero", and with omitempty it cannot even send zero.
		res.add(cmd, side, CatCardinality, path, fmt.Sprintf(
			"schema says optional; %s %s cannot express absence (zero is a legal value)", g.GoName, g.GoType))
	}
}

// checkOmitEmpty compares declared optionality against the tag.
//
// Only one of the four combinations is a bug in commands/. omitempty governs
// encoding only, so on a response type — which this library decodes and never
// sends — it is inert either way, and the two response cases are recorded and
// suppressed. On a request it decides what goes on the wire: a missing
// omitempty makes the client emit an empty element it was never asked to send,
// which is the actionable case.
func checkOmitEmpty(res *Result, cmd Command, side string, s *Node, g *GoNode, path string) {
	switch {
	case s.Card == CardRequired && g.OmitEmpty:
		f := Finding{Command: cmd.Name, Side: side, Category: CatOmitEmpty, Path: path,
			Detail: fmt.Sprintf("schema says required, %s is tagged omitempty", g.GoName)}
		switch {
		case s.InChoice:
			f.Suppressed = true
			f.Reason = "inside <or>; only one branch is ever present"
		case side == "request":
			// Requests are over-constrained: create_report marks report/@id,
			// config_id, format_id, extension and content_type all required
			// even though real payloads only ever populate one submission
			// path. The schema is the less trustworthy side here.
			f.Suppressed = true
			f.Reason = "request-side over-constraint"
		default:
			// Responses are over-constrained too. authenticate declares role
			// and timezone required, and its own second example is a 400 with
			// neither.
			f.Suppressed = true
			f.Reason = "inert: responses are decoded, and omitempty only affects encoding"
		}
		res.Findings = append(res.Findings, f)

	case s.Card == CardOptional && !g.OmitEmpty && !g.Pointer && !g.Slice:
		f := Finding{Command: cmd.Name, Side: side, Category: CatOmitEmpty, Path: path,
			Detail: fmt.Sprintf("schema says optional, %s %s has no omitempty and is always emitted",
				g.GoName, g.GoType)}
		if side == "response" {
			// This is the fidelity corpus's 'added' bucket, which its summary
			// records and explicitly does not count as loss.
			f.Suppressed = true
			f.Reason = "inert on decode; shows up as 'added' in the fidelity corpus"
		}
		res.Findings = append(res.Findings, f)
	}
}

// checkType compares the schema's lexical space against what the Go type can
// hold.
func checkType(res *Result, cmd Command, side string, s *Node, g *GoNode, path string) {
	t := Resolve(s.Type)
	if t == "" || !g.IsScalarGo() {
		return
	}
	base := strings.TrimPrefix(g.GoType, "*")

	// The headline. "boolean" is declared as text: zero is false and anything
	// else, including the empty string, is true. boolean_atoi is the opposite
	// — it runs atoi first, so "" is false. Neither is Go's bool, and both
	// need a custom marshaler.
	if t == "boolean" || t == "boolean_atoi" {
		if base == "bool" {
			res.add(cmd, side, CatType, path, fmt.Sprintf(
				"schema type %s is text (%s), Go %s is bool: '7' fails to decode and '1' re-emits as 'true'",
				t, booleanNote(t), g.GoName))
		}
		return
	}
	if want, ok := numericTypes[t]; ok && base == "string" {
		res.add(cmd, side, CatType, path, fmt.Sprintf(
			"schema type %s, Go %s is string (want %s)", t, g.GoName, want))
		return
	}
	if textual[t] && (base == "bool" || base == "int" || base == "float64") {
		res.add(cmd, side, CatType, path, fmt.Sprintf(
			"schema type %s is text, Go %s is %s", t, g.GoName, base))
	}
}

func booleanNote(t string) string {
	if t == "boolean_atoi" {
		return "atoi first, so \"\" is false"
	}
	return "\"\" is true"
}

var numericTypes = map[string]string{
	"integer":  "int",
	"decimal":  "float64",
	"severity": "float64",
}

var textual = map[string]bool{
	"text": true, "uuid": true, "uuid_or_empty": true, "name": true,
	"iso_time": true, "ctime": true, "oid": true, "status": true,
}

// confirmed reports whether the fidelity corpus already recorded this path as
// dropped from a real document.
func (c *Comparer) confirmed(cmd, path string) bool {
	set := c.Dropped[cmd]
	if set == nil {
		return false
	}
	if set[path] {
		return true
	}
	// A dropped parent folds its descendants in the fidelity report, so a
	// dropped ancestor confirms this path too.
	for p := path; ; {
		i := strings.LastIndex(p, "/")
		if i <= 0 {
			return false
		}
		p = p[:i]
		if set[p] {
			return true
		}
	}
}

func (r *Result) add(cmd Command, side, cat, path, detail string) {
	r.Findings = append(r.Findings, Finding{
		Command: cmd.Name, Side: side, Category: cat, Path: path, Detail: detail,
	})
}

func otherKind(k Kind) Kind {
	if k == KindAttribute {
		return KindElement
	}
	return KindAttribute
}

func key(k Kind, name string) string { return k.String() + ":" + name }

func displayType(s *Node) string {
	if s.Type != "" {
		return s.Type
	}
	if len(s.Alts) > 0 {
		return "enum"
	}
	return "text"
}

// Path is where the Go node sits under parent, in fidelity-report syntax.
func (n *GoNode) Path(parent string) string {
	if n.Kind == KindAttribute {
		return parent + "/@" + n.Name
	}
	return parent + "/" + n.Name
}

func sortFindings(fs []Finding) {
	rank := map[string]int{}
	for i, c := range Categories {
		rank[c] = i
	}
	sort.SliceStable(fs, func(i, j int) bool {
		a, b := fs[i], fs[j]
		switch {
		case a.Side != b.Side:
			return a.Side < b.Side // request before response
		case rank[a.Category] != rank[b.Category]:
			return rank[a.Category] < rank[b.Category]
		case a.Path != b.Path:
			return a.Path < b.Path
		}
		return a.Detail < b.Detail
	})
}

// --- the corroborating oracle ----------------------------------------------

// LoadDroppedPaths reads the per-fixture fidelity reports written by
// internal/gmptest and returns, per command, the response paths that a real
// document lost on the way through the Go types.
//
// This package only reads that tree. A MISSING finding whose path appears here
// was reached twice over, from the schema and from a captured document, by two
// tools that share no code.
func LoadDroppedPaths(testdataDir string) (map[string]map[string]bool, error) {
	matches, err := filepath.Glob(filepath.Join(testdataDir, "*", "*", "*_fidelity.txt"))
	if err != nil {
		return nil, err
	}
	sort.Strings(matches)

	out := map[string]map[string]bool{}
	for _, m := range matches {
		cmd := filepath.Base(filepath.Dir(m))
		data, err := os.ReadFile(m) //nolint:gosec // a report in the repo's own testdata
		if err != nil {
			return nil, err
		}
		for _, p := range parseDropped(string(data)) {
			if out[cmd] == nil {
				out[cmd] = map[string]bool{}
			}
			out[cmd][p] = true
		}
	}
	return out, nil
}

// parseDropped pulls the paths out of the DROPPED block of a fidelity report.
func parseDropped(s string) []string {
	var out []string
	in := false
	for _, line := range strings.Split(s, "\n") {
		switch {
		case strings.HasPrefix(line, "DROPPED"):
			in = true
			continue
		case line != "" && !strings.HasPrefix(line, " "):
			in = false
			continue
		}
		if !in {
			continue
		}
		f := strings.Fields(line)
		if len(f) > 0 && strings.HasPrefix(f[0], "/") {
			// Trailing /#text marks a lost value, not a lost element; the
			// element path is what a MISSING finding names.
			out = append(out, strings.TrimSuffix(f[0], "/#text"))
		}
	}
	return out
}
