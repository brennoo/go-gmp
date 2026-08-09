// Package codegen turns the GMP schema's declared command shapes into Go
// source for commands/.
//
// It is split in two stages at the AGPL boundary described in
// docs/review-2026-08.md, T-9:
//
//   - Stage A (snapshot_test.go, -update only, never run in CI) reads
//     vendor/XML/commands through internal/schema and writes a committed JSON
//     snapshot under testdata/schema. Only names, kinds, cardinalities and
//     type names cross the boundary: schema.Command carries no <summary> or
//     <example> prose to begin with.
//   - Stage B (this file, generate.go) is pure. It reads the committed
//     snapshot, never vendor/XML, and is the part CI runs.
package codegen

// Override replaces the generator's mechanical default Go type for a named
// schema type. Every entry must cite the T-9/T-5 defect it addresses so an
// auditor can trace why the mechanical DSL->Go mapping wasn't trusted here.
type Override struct {
	GoType string
	Reason string
}

// TypeOverrides is keyed by the schema's declared <type> name (schema.Node.Type),
// after internal/schema's alias table has been applied.
//
// The table is deliberately tiny. Every entry exists because the schema's own
// type definition says something the mechanical mapping would get wrong, not
// because a Go type happened to look nicer.
var TypeOverrides = map[string]Override{
	// vendor/XML's <type><name>boolean</name> declares <pattern>text</pattern>
	// and documents "zero is false, anything else is true. As a result, the
	// empty string is considered true". Go's encoding/xml bool would reject
	// "" and accept only "true"/"false"/"1"/"0".
	"boolean": {
		GoType: "commands.GMPBool",
		Reason: `T-9 defect #6: lexical free text ("0"=false, anything else incl. empty=true), not xsd:boolean`,
	},
	// boolean_atoi is deliberately NOT mapped to GMPBool: vendor/XML's
	// <type><name>boolean_atoi</name> says the value is run through atoi(3)
	// first, "which means that an empty string is considered false" — the
	// exact opposite of <type><name>boolean</name>'s empty-is-true rule
	// above (T-9 defect #6's two boolean flavours disagree on the empty
	// case), so reusing GMPBool would silently invert every empty element.
	//
	// Neither generated family uses it. Inventoried by loading all three
	// note and all three override commands through schema.Load against
	// vendor/XML/commands and walking every leaf node: zero nodes with a
	// resolved Type of boolean_atoi, and no entry in internal/schema's
	// typeAliases maps onto it either — the only boolean_atoi text in those
	// files is the type definition each split schema file repeats. So there
	// is still no fixture to verify a dedicated type against, and the name
	// keeps falling through to the mechanical string default rather than
	// gaining a plausible-looking but unverified override. Add
	// commands.GMPAtoiBool (empty/non-numeric -> false, "0" -> false,
	// anything atoi reads as non-zero -> true) when the first family that
	// actually declares one is generated.
	//
	// gvmd emits ctime(3) strings and empty elements under fields the schema
	// labels iso_time; time.Time's RFC 3339 UnmarshalText fails on both.
	"iso_time": {
		GoType: "commands.GMPTime",
		Reason: "T-5e: strict RFC3339 time.Time unmarshal fails on real gvmd output shapes",
	},
	"ctime": {
		GoType: "commands.GMPTime",
		Reason: "T-5e: same as iso_time",
	},
}
