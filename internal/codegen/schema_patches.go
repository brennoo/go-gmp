package codegen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/brennoo/go-gmp/internal/schema"
)

// SchemaPatch splices a node the served schema itself never declares into a
// parsed schema.Command tree, before Generate ever sees it. Unlike
// FieldOverride (whose GoType can only reference an existing Go type, for a
// field the mechanical schema->Go mapping can't produce), this adds a real
// schema.Node, so the existing walk/collect/pick/emitStruct pipeline
// generates a proper family-local struct for it — full cardinality, pointer
// and doc-comment handling, no special-cased emission path.
//
// Reserved for a field with direct GMP-wire evidence (live capture, cross-
// checked where possible against a sibling command's own declared schema for
// the same reused element) that the schema served by gvmd itself omits — not
// a place to route ordinary FieldOverride cases through. ApplySchemaPatches
// fails closed the moment the target already has a child of that name/kind,
// so a patch that the schema catches up with breaks the build loudly instead
// of silently drifting; there is no separate staleness/"consumed" table.
//
// A patch must never leak into a raw-schema comparison: LoadSnapshots stays
// unpatched for exactly that reason, and internal/codegen/drift.go's
// DetectDrift and vendor_parity_test.go's TestVendorParity both call it
// directly rather than LoadGenerationModels, so a locally-added field is
// never mistaken for upstream schema drift or a vendor/XML mismatch. Only
// LoadGenerationModels (generation, validation-after-generation, gengmp
// -check) applies patches.
type SchemaPatch struct {
	Name     string // wire element/attribute name (schema.Node.Name)
	Kind     schema.Kind
	Card     schema.Card
	Type     string // schema.Node.Type; required for every leaf, attribute, or chardata node
	Chardata bool   // schema.Node.Chardata: bare text alongside Children (mixed content)

	// InsertAfter names an existing sibling this patch is spliced after, so
	// the generated struct's field order matches the wire order live capture
	// shows rather than always landing last. Empty appends at the end.
	InsertAfter string

	// Children recursively becomes the new node's own Children, converted the
	// same way; nil for a leaf.
	Children []SchemaPatch

	Reason string // cites the live/upstream evidence for both presence and shape
}

// SchemaPatches is keyed exactly like FieldOverrides: "<command>:/<path>"
// identifying the already-schema-declared node the new child is spliced
// under. ApplySchemaPatches consults it once per LoadGenerationModels (or
// check.go's loadFamilySnapshots) call, before Generate runs, so every
// generation-side consumer sees the same patched tree.
var SchemaPatches = map[string][]SchemaPatch{
	"get_results:/result/nvt": {{
		Name: "solution", Kind: schema.KindElement, Card: schema.CardOptional,
		Type: "text", Chardata: true, InsertAfter: "tags",
		Reason: "T-9 result review, finding 1: internal/codegen/testdata/schema/get_results.json's " +
			"result/nvt node has no solution child, but all 4 live fixtures " +
			"(internal/gmptest/testdata/live/get_results/00{1,2,3,4}_response.xml) send <solution type=''> " +
			"right after </tags>, retained as DROPPED in every one of those fixtures' own _fidelity.txt. The " +
			"full shape — optional element, chardata text plus optional type/method attributes — is the exact " +
			"declaration get_nvts' OWN live schema snapshot (internal/codegen/testdata/schema/get_nvts.json's " +
			"get_nvts_response/nvt/solution) already carries for the identical nvt/solution reuse, already " +
			"proven correct as commands.NVTSolution (commands/get_nvts.go); no live get_results fixture " +
			"happens to send a non-empty type, method or text, so those parts are evidenced by get_nvts' own " +
			"schema declaration and generated type, not by get_results' own captures",
		Children: []SchemaPatch{
			{
				Name: "type", Kind: schema.KindAttribute, Card: schema.CardOptional, Type: "text",
				Reason: "see the solution entry above: get_nvts' own schema declares solution/@type " +
					"optional text; all 4 live get_results fixtures send it present but empty (type='')",
			},
			{
				Name: "method", Kind: schema.KindAttribute, Card: schema.CardOptional, Type: "text",
				Reason: "see the solution entry above: get_nvts' own schema declares solution/@method " +
					"optional text; no live get_results fixture ever sends it, matching get_nvts' own " +
					"declaration of it as optional rather than evidencing its absence as a defect",
			},
		},
	}},
}

// ApplySchemaPatches splices every SchemaPatches entry into the matching
// command's tree in cmds, in place. Called once per LoadGenerationModels (and
// loadFamilySnapshots) call, before anything downstream sees the tree.
func ApplySchemaPatches(cmds []schema.Command) error {
	if err := validateSchemaPatches(); err != nil {
		return err
	}

	keys := make([]string, 0, len(SchemaPatches))
	for k := range SchemaPatches {
		keys = append(keys, k)
	}
	sort.Strings(keys) // deterministic application order and error messages

	index := make(map[string]int, len(cmds))
	for i := range cmds {
		index[cmds[i].Name] = i
	}

	for _, key := range keys {
		command, path, err := splitSchemaPatchKey(key)
		if err != nil {
			return fmt.Errorf("codegen: SchemaPatches[%q]: %w", key, err)
		}
		i, ok := index[command]
		if !ok {
			continue // command not in this call's scope
		}
		if err := applySchemaPatchesTo(&cmds[i], path, SchemaPatches[key]); err != nil {
			return fmt.Errorf("codegen: SchemaPatches[%q]: %w", key, err)
		}
	}
	return nil
}

// applySchemaPatchesTo locates path's target node in c's tree — trying
// Request and Response — and appends each patch as a new child, failing
// closed rather than silently no-op'ing or guessing.
func applySchemaPatchesTo(c *schema.Command, path []string, patches []SchemaPatch) error {
	reqNode, reqErr := descendSchemaPath(c.Request, path)
	respNode, respErr := descendSchemaPath(c.Response, path)

	var target *schema.Node
	switch {
	case reqErr == nil && respErr == nil:
		return fmt.Errorf("path %q matches both the request and response tree; ambiguous target", "/"+strings.Join(path, "/"))
	case reqErr == nil:
		target = reqNode
	case respErr == nil:
		target = respNode
	default:
		return fmt.Errorf("no node at %q in either request or response (request: %w; response: %w)",
			"/"+strings.Join(path, "/"), reqErr, respErr)
	}

	for _, p := range patches {
		if findSchemaChild(target, p.Name) != nil {
			return fmt.Errorf("%s already has a child %q — the schema has caught up; remove this patch", target.Name, p.Name)
		}
		node := schemaPatchNode(p)
		if p.InsertAfter == "" {
			target.Children = append(target.Children, node)
			continue
		}
		at := -1
		for i, kid := range target.Children {
			if kid.Name == p.InsertAfter {
				at = i
				break
			}
		}
		if at == -1 {
			return fmt.Errorf("%s: InsertAfter %q names no existing child", p.Name, p.InsertAfter)
		}
		target.Children = append(target.Children[:at+1],
			append([]*schema.Node{node}, target.Children[at+1:]...)...)
	}
	return nil
}

// descendSchemaPath walks n's Children by name for each path segment,
// returning the final node — the same segment semantics fieldOverrideKey's
// path already uses (the root's own name is never a segment).
func descendSchemaPath(n *schema.Node, path []string) (*schema.Node, error) {
	cur := n
	for _, seg := range path {
		next := findSchemaChild(cur, seg)
		if next == nil {
			return nil, fmt.Errorf("%q has no child %q", cur.Name, seg)
		}
		cur = next
	}
	return cur, nil
}

func findSchemaChild(n *schema.Node, name string) *schema.Node {
	for _, c := range n.Children {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// schemaPatchNode converts one SchemaPatch (and its Children, recursively)
// into the *schema.Node the rest of the pipeline consumes. InChoice and Alts
// have no SchemaPatch equivalent: every use so far is an unconditional
// addition, not an <or> branch or an alternately-typed field.
func schemaPatchNode(p SchemaPatch) *schema.Node {
	kids := make([]*schema.Node, len(p.Children))
	for i, c := range p.Children {
		kids[i] = schemaPatchNode(c)
	}
	return &schema.Node{
		Name: p.Name, Kind: p.Kind, Card: p.Card, Type: p.Type,
		Chardata: p.Chardata, Children: kids,
	}
}

func splitSchemaPatchKey(key string) (command string, path []string, err error) {
	if err := validateOverrideKey(key); err != nil { // reuse the shared "<command>:/<path>" shape check
		return "", nil, err
	}
	command, rest, _ := strings.Cut(key, ":/")
	return command, strings.Split(rest, "/"), nil
}

// validateSchemaPatches checks SchemaPatches' own shape, independent of
// which commands are in scope for this call — the same unconditional check
// validateFieldOverrides runs for FieldOverrides. There is no GoType
// equivalent to validate here: a SchemaPatch produces a schema.Node, not a Go
// type directly — the existing mechanical mapping (fieldOf/goTypeOf) derives
// the Go type from the spliced node the same way it does for any
// organically-declared one.
func validateSchemaPatches() error {
	keys := make([]string, 0, len(SchemaPatches))
	for k := range SchemaPatches {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if err := validateOverrideKey(key); err != nil {
			return fmt.Errorf("codegen: SchemaPatches[%q]: %w", key, err)
		}
		if err := validateSchemaPatchList(key, SchemaPatches[key]); err != nil {
			return err
		}
	}
	return nil
}

// validateSchemaPatchList rejects two siblings sharing a wire name regardless
// of Kind: applySchemaPatchesTo/findSchemaChild resolve an existing child by
// name alone, so an element and an attribute both named "x" would pass this
// check but only one of them would ever actually collide against a schema
// that "catches up" — the other would insert silently, an inconsistency this
// closes at table-validation time instead of leaving it to be found later.
func validateSchemaPatchList(key string, patches []SchemaPatch) error {
	seen := map[string]bool{}
	for i, p := range patches {
		if err := validateSchemaPatchShape(key, i, p); err != nil {
			return err
		}
		if seen[p.Name] {
			return fmt.Errorf("codegen: SchemaPatches[%q]: duplicate child %q", key, p.Name)
		}
		seen[p.Name] = true
		if err := validateSchemaPatchList(key+"/"+p.Name, p.Children); err != nil {
			return err
		}
	}
	return nil
}

// validateSchemaPatchShape checks one SchemaPatch's own fields, independent
// of its siblings — the duplicate-child and recursion checks stay in
// validateSchemaPatchList, which needs the whole slice.
func validateSchemaPatchShape(key string, i int, p SchemaPatch) error {
	if p.Name == "" {
		return fmt.Errorf("codegen: SchemaPatches[%q][%d]: Name is empty", key, i)
	}
	if p.Kind != schema.KindElement && p.Kind != schema.KindAttribute {
		return fmt.Errorf("codegen: SchemaPatches[%q][%d]: Kind %v is neither KindElement nor KindAttribute", key, i, p.Kind)
	}
	if p.Card != schema.CardRequired && p.Card != schema.CardOptional && p.Card != schema.CardRepeated {
		return fmt.Errorf("codegen: SchemaPatches[%q][%d]: Card %v is not a valid schema.Card", key, i, p.Card)
	}
	if err := validateSchemaPatchLeafShape(key, i, p); err != nil {
		return err
	}
	if p.Reason == "" {
		return fmt.Errorf("codegen: SchemaPatches[%q][%d]: Reason is empty", key, i)
	}
	return nil
}

// validateSchemaPatchLeafShape checks the constraints that depend on whether
// p is an attribute, a chardata node, or a pure struct element.
func validateSchemaPatchLeafShape(key string, i int, p SchemaPatch) error {
	if p.Kind == schema.KindAttribute && len(p.Children) > 0 {
		return fmt.Errorf("codegen: SchemaPatches[%q][%d]: attribute %q cannot have Children", key, i, p.Name)
	}
	if p.Kind == schema.KindAttribute && p.Chardata {
		return fmt.Errorf("codegen: SchemaPatches[%q][%d]: attribute %q cannot set Chardata", key, i, p.Name)
	}
	// Type is required wherever the node carries a lexical value of its own:
	// every attribute, every chardata node, and every pure leaf element. It
	// may be empty only for a pure struct element (real Children, no
	// Chardata).
	isLeafOrMixed := p.Kind == schema.KindAttribute || p.Chardata || len(p.Children) == 0
	if isLeafOrMixed && p.Type == "" {
		return fmt.Errorf("codegen: SchemaPatches[%q][%d]: %q needs a Type (leaf or mixed-content node)", key, i, p.Name)
	}
	return nil
}
