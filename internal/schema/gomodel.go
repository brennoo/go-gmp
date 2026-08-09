package schema

import (
	"encoding"
	"encoding/xml"
	"reflect"
	"sort"
	"strings"
)

// The Go side is derived by reflection rather than by go/ast.
//
// Reflection is used because it reads the same thing encoding/xml reads. It
// gives the resolved reflect.StructTag, the real pointer/slice/value shape
// after named types and aliases are followed, and it crosses file and package
// boundaries for free. go/ast would give the source text of a tag and a type
// *expression*: resolving `*TargetCredential` to its fields means building a
// name table across commands/*.go and re-implementing Go's own rules for
// embedding and shadowing. Every one of those is a place for the tool to be
// wrong in a way that looks like a struct bug.
//
// What reflection cannot see is the source: field order as written, comments,
// and any type never reachable from internal/gmptest.Registry. Field order is
// irrelevant to a name-matched diff, and unreachable types are by definition
// not part of any command's shape, so neither loss matters here.
//
// This is also independent of the fidelity tooling, which never reflects: it
// marshals and diffs documents. Two oracles built the same way would agree for
// uninteresting reasons.

// GoNode is one Go struct field, or a synthetic wrapper standing in for the
// intermediate element in an `a>b` tag path.
type GoNode struct {
	Name      string // the XML name, not the Go name
	GoName    string
	Kind      Kind
	Card      Card
	Pointer   bool
	Slice     bool
	OmitEmpty bool
	GoType    string
	Synthetic bool // an `a>b` wrapper: no Go field of its own
	Children  []*GoNode

	// Chardata is set when the struct captures the element's text through a
	// `,chardata` field, which means the element is a leaf carrying a value
	// even though it is modelled as a struct.
	Chardata bool
	// CatchAll is set when the struct has a `,any` field, so unknown children
	// are captured rather than dropped and EXTRA/MISSING mean less.
	CatchAll bool
}

// GoRoot is the model of one command's request or response type.
type GoRoot struct {
	// TypeName is the Go type, for report lines that must be traceable back
	// to a declaration.
	TypeName string
	Root     *GoNode
}

var (
	xmlNameType       = reflect.TypeFor[xml.Name]()
	textMarshalerType = reflect.TypeFor[encoding.TextMarshaler]()
)

// ModelOf reflects v, which must be a pointer to a struct with an XMLName
// field, into the same shape the parser produces for the schema.
func ModelOf(v any) GoRoot {
	t := reflect.TypeOf(v)
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t == nil || t.Kind() != reflect.Struct {
		return GoRoot{}
	}
	root := &GoNode{Name: rootName(t), GoName: t.Name(), GoType: t.String()}
	fill(root, t, map[reflect.Type]bool{})
	return GoRoot{TypeName: t.String(), Root: root}
}

// rootName reads the element name off the XMLName tag, falling back to the
// type name the way encoding/xml does.
func rootName(t reflect.Type) string {
	if f, ok := t.FieldByName("XMLName"); ok && f.Type == xmlNameType {
		if name, _ := splitTag(f.Tag.Get("xml")); name != "" {
			return name
		}
	}
	return t.Name()
}

// fill walks a struct type, adding one child per XML-visible field.
//
// seen breaks recursive types. A type already being expanded higher up the
// stack is left as a leaf: the diff would otherwise not terminate, and the
// schema has no recursive definitions to compare it against anyway.
func fill(n *GoNode, t reflect.Type, seen map[reflect.Type]bool) {
	if seen[t] {
		return
	}
	seen[t] = true
	defer delete(seen, t)

	for i := range t.NumField() {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue // unexported
		}
		if f.Anonymous {
			// encoding/xml flattens an anonymous struct or *struct field's
			// children into the parent unconditionally — any xml tag on
			// the field itself (name, any, chardata, ...) is ignored. Only
			// an anonymous field of interface type is treated as named,
			// and structUnder already returns nil for anything that isn't
			// a struct, so that case falls through to addField below.
			if elem := structUnder(f.Type); elem != nil {
				fill(n, elem, seen)
				continue
			}
		}
		addField(n, f, seen)
	}
	sort.SliceStable(n.Children, func(i, j int) bool {
		if n.Children[i].Kind != n.Children[j].Kind {
			return n.Children[i].Kind < n.Children[j].Kind
		}
		return n.Children[i].Name < n.Children[j].Name
	})
}

func addField(n *GoNode, f reflect.StructField, seen map[reflect.Type]bool) {
	name, opts := splitTag(f.Tag.Get("xml"))
	if name == "-" {
		return
	}
	switch {
	case f.Type == xmlNameType:
		return
	case opts["chardata"], opts["cdata"]:
		n.Chardata = true
		return
	case opts["any"]:
		n.CatchAll = true
		return
	case opts["innerxml"], opts["comment"]:
		n.Chardata = true
		return
	}
	if name == "" {
		name = f.Name // what encoding/xml falls back to
	}

	// `a>b` declares an intermediate element with no Go field of its own.
	parent := n
	parts := strings.Split(name, ">")
	for _, step := range parts[:len(parts)-1] {
		parent = wrapper(parent, step)
	}

	kid := leafNode(parts[len(parts)-1], f, opts)
	parent.Children = append(parent.Children, kid)

	if elem := structUnder(f.Type); elem != nil {
		fill(kid, elem, seen)
	}
}

// wrapper finds or creates the synthetic node for one step of an `a>b` path.
func wrapper(parent *GoNode, step string) *GoNode {
	for _, c := range parent.Children {
		if c.Name == step && c.Synthetic {
			return c
		}
	}
	w := &GoNode{Name: step, Card: CardOptional, Synthetic: true, GoType: "(wrapper)"}
	parent.Children = append(parent.Children, w)
	return w
}

func leafNode(name string, f reflect.StructField, opts map[string]bool) *GoNode {
	kid := &GoNode{
		Name:      name,
		GoName:    f.Name,
		GoType:    f.Type.String(),
		OmitEmpty: opts["omitempty"],
		Pointer:   f.Type.Kind() == reflect.Pointer,
		Slice:     f.Type.Kind() == reflect.Slice && f.Type.Elem().Kind() != reflect.Uint8,
	}
	if opts["attr"] {
		kid.Kind = KindAttribute
	}
	kid.Card = cardOf(kid)
	return kid
}

// cardOf reads the cardinality the Go declaration actually expresses.
//
// A slice is zero-or-more. A pointer is genuinely optional: nil and present
// are distinguishable. A value field with omitempty is optional only in the
// weaker sense that the zero value is not written — which is right for a
// string and wrong for a bool or an int, and cardinalityFinding is where that
// distinction is drawn.
func cardOf(n *GoNode) Card {
	switch {
	case n.Slice:
		return CardRepeated
	case n.Pointer, n.OmitEmpty:
		return CardOptional
	}
	return CardRequired
}

// structUnder returns the struct type a field ultimately holds, or nil when
// the field is a scalar. A struct that implements encoding.TextMarshaler
// (time.Time, GMPTime) is a scalar here: it is a leaf on the wire, not a
// container to descend into.
func structUnder(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice {
		if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
			return nil // []byte is base64 or raw text, not a container
		}
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct || t == xmlNameType || t.Implements(textMarshalerType) {
		return nil
	}
	return t
}

// splitTag parses an xml struct tag into its name and its option set.
func splitTag(tag string) (string, map[string]bool) {
	parts := strings.Split(tag, ",")
	opts := make(map[string]bool, len(parts))
	for _, o := range parts[1:] {
		opts[strings.TrimSpace(o)] = true
	}
	return strings.TrimSpace(parts[0]), opts
}

// IsScalarGo reports whether the Go type holds a value rather than a subtree.
func (n *GoNode) IsScalarGo() bool {
	return len(n.Children) == 0
}

// ZeroIsMeaningful reports whether the Go zero value of this field is a legal
// wire value, which is what makes `omitempty` on a non-pointer lossy. A "0"
// severity and an absent severity are different facts; an empty string and an
// absent string are, in this protocol, not.
func (n *GoNode) ZeroIsMeaningful() bool {
	if n.Pointer || n.Slice {
		return false
	}
	switch strings.TrimPrefix(n.GoType, "*") {
	case "bool", "int", "int32", "int64", "uint", "uint32", "uint64", "float32", "float64":
		return true
	}
	return false
}
