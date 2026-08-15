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

// PointerOverride keeps one command's optional field pointer-typed even
// though its Go type is not in pointerScalars (generate.go), because live
// protocol evidence contradicts pointerScalars' generator-wide default
// assumption for that specific field.
type PointerOverride struct {
	Reason string
}

// PointerOverrides is keyed by "<command>.<element or attribute name>", and
// only ever applies to a command's own top-level request/response fields
// (fieldOf's command argument is empty for a nested struct's fields, which
// have no single owning command, so no key here can ever match one).
//
// pointerScalars documents a generator-wide default: "in this protocol an
// empty string and an absent element mean the same thing", so an optional
// string field is never a pointer. A live experiment against gvmd 26.36.1
// (Greenbone Community Edition), run to answer the question review-only
// `alert` had left open, proved that default wrong for modify_alert's name
// and comment: omitting the element leaves the stored value unchanged, while
// sending it empty (`<comment></comment>`) clears it. A bare `string` with
// `,omitempty` can only ever send the first form, so adoption's mechanical
// string-not-pointer conversion silently removed a request `alert`'s
// hand-written `*string` could express. See docs/review-2026-08.md's `alert
// adopted` section for the experiment and the reasoning below.
//
// The table stayed narrow on purpose until the post-tls_certificate T-9
// revisit pass ran the same live experiment (gvmd 26.36.1) against every
// remaining family with optional text fields: modify_role, modify_group,
// modify_tag and modify_schedule. Results, one per family below.
//
// modify_role and modify_group: name, comment and users are NOT like
// modify_alert's — omitting an element does not leave the stored value
// unchanged. A bare `<modify_role role_id="..."/>` with every optional
// element absent returns 200 OK but clears name, comment and users
// unconditionally, and sending only `<comment>` still wipes the other two
// even though they were never mentioned. Omitted and explicit-empty are
// equivalent here, just not to "unchanged" — a pointer type could not
// express a third wire state that does not exist, so no override is
// warranted. The real hazard is behavioral, not a typing gap: a caller who
// means to change only one field silently blanks the rest. See
// docs/review-2026-08.md's post-tls_certificate revisit section.
//
// modify_tag: name behaves like modify_role/modify_group's fields do NOT —
// omitted and explicit-empty both leave the stored name unchanged, so it
// keeps the mechanical string default. value and comment behave like
// modify_alert's (omit leaves the stored value unchanged, explicit-empty
// clears it), so both get an entry below.
//
// modify_schedule: name, comment, icalendar and timezone all behave like
// modify_tag's name — omitted and explicit-empty are indistinguishable,
// both leaving the stored value unchanged — so none gets an entry. A
// separate, unrelated wire quirk surfaced while testing timezone: gvmd
// rejects a request that sends `<timezone>` without also sending
// `<icalendar>` ("ICALENDAR element is required"), even though the schema
// declares both independently optional. That is a cross-field requirement
// the server enforces and the schema does not describe, not an
// omit-vs-explicit-empty distinction, so it does not belong in this table
// either.
// FieldOverride adds a field to a generated struct that the schema's own
// declared pattern does not, and structurally cannot, produce — a real,
// evidence-backed protocol field the schema simply never declares.
type FieldOverride struct {
	Name      string // wire element/attribute name; ignored when Chardata
	Attr      bool
	Chardata  bool // this struct's own `,chardata` field; mutually exclusive with Name/Attr
	GoType    string
	OmitEmpty bool
	Reason    string // cites the live/upstream evidence
}

// FieldOverrides is keyed "<command>:/<ancestor>/.../<node-name>" — the same
// addressing style as PointerOverrides's "<command>.<field>", extended with a
// path since this targets a nested struct reached from a command, not only
// a command's top-level request/response envelope. Scope is deliberately
// narrow: a structurally deduplicated shape reached by more than one
// (command, path) use is only allowed when every use resolves to the same
// effective override set (fieldOverridesFor, generate.go) — a use that
// disagrees, or has no entry at all, fails the build rather than leaking one
// use's override onto another. Only emitStruct consults this table (not
// emitRequest/emitResponse) — widening either needs its own evidence, not a
// speculative entry now.
var FieldOverrides = map[string][]FieldOverride{
	"get_tls_certificates:/tls_certificate": {{
		Name: "valid", GoType: "GMPBool",
		Reason: "T-9 finding 1: schema's tls_certificate pattern never declares <valid>, but " +
			"internal/gmptest/testdata/live/get_tls_certificates/003_response.xml and the upstream " +
			"example both send <valid>1</valid> as real data",
	}},
	"get_tls_certificates:/tls_certificate/certificate": {{
		Chardata: true, GoType: "string",
		Reason: "T-9 finding 2: schema's certificate pattern declares only @format, no chardata token, " +
			"but internal/gmptest/testdata/upstream/get_tls_certificates/001_response.xml sends real " +
			"base64 certificate bytes as this element's text (see IncludeCertificateData)",
	}},
}

var PointerOverrides = map[string]PointerOverride{
	"modify_alert.name": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <name> are distinct requests",
	},
	"modify_alert.comment": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <comment> are distinct requests",
	},
	// modify_filter: a second live experiment (T-9, gvmd 26.36.1), run to
	// answer the review finding that filter's review-only doc had asserted
	// "no PointerOverrides entry needed" without testing it. Result: name and
	// comment behave exactly like modify_alert's (omit leaves the stored
	// value unchanged, explicit-empty clears it); term also distinguishes the
	// two, though explicit-empty doesn't clear it to "" but resets it to
	// gvmd's default filter term (e.g. "first=1 rows=11 sort=name") — still a
	// real distinction a bare string+omitempty could never send. type showed
	// no distinction (omit and explicit-empty both leave it unchanged), so it
	// keeps the mechanical default and has no entry here.
	"modify_filter.name": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <name> are distinct requests",
	},
	"modify_filter.comment": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <comment> are distinct requests",
	},
	"modify_filter.term": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <term> are distinct requests (explicit-empty resets to gvmd's default term rather than leaving it unchanged)",
	},
	// modify_tls_certificate: a third live experiment (T-9, gvmd 26.36.1), run
	// to answer the review-only doc's unsubstantiated "no PointerOverrides
	// entry needed" claim for name/comment/copy, the same claim already
	// disproven once for modify_filter. Result: name and comment behave
	// exactly like modify_alert's/modify_filter's (omit leaves the stored
	// value unchanged, explicit-empty clears it to ""). copy did not: omit,
	// explicit-empty, and setting it to a second real certificate's UUID all
	// returned 200 OK but left the target certificate's name, comment, and
	// both fingerprints completely unchanged in every case, so unlike
	// create_tls_certificate's copy (which clones at creation) this field is
	// observably inert on modify_tls_certificate in this gvmd version -- omit
	// and explicit-empty are not merely equivalent, all three tested inputs
	// are. No override is warranted for a field with no observable effect.
	"modify_tls_certificate.name": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <name> are distinct requests",
	},
	"modify_tls_certificate.comment": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <comment> are distinct requests",
	},
	// modify_tag: the post-tls_certificate T-9 revisit pass's live experiment.
	// value and comment behave like modify_alert's; name did not (see the
	// doc comment above this table) and keeps the mechanical default.
	"modify_tag.value": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <value> are distinct requests",
	},
	"modify_tag.comment": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <comment> are distinct requests",
	},
	// modify_ticket: live gvmd 26.36.1 experiment run against review-only
	// ticket's deferred question. comment, open_note, fixed_note and
	// closed_note behave like modify_alert's fields (omit leaves the stored
	// value unchanged, explicit-empty clears it to ""). status does not: an
	// explicit-empty <status></status> is rejected outright with
	// "400 Error in status" rather than accepted as a clear, so there is no
	// legal empty wire state for a pointer to express and it keeps the
	// mechanical string default (see the modify_ticket CommandWarning below).
	"modify_ticket.comment": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <comment> are distinct requests",
	},
	"modify_ticket.open_note": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <open_note> are distinct requests",
	},
	"modify_ticket.fixed_note": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <fixed_note> are distinct requests",
	},
	"modify_ticket.closed_note": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <closed_note> are distinct requests",
	},
	// modify_credential: live gvmd 26.36.1 experiment run against review-only
	// credential's optional string fields (up, cc, snmp and krb5 credentials
	// created and modified directly against the running instance). comment
	// behaves like modify_alert's (omit leaves the stored value unchanged,
	// explicit-empty clears it to ""). name, login and auth_algorithm did
	// not: explicit-empty left each stored value unchanged, indistinguishable
	// from omitting the element, so they keep the mechanical string default.
	// certificate, kdc, community and password are not covered here at all —
	// see the modify_credential CommandWarning below for why kdc and realm
	// can't take a pointer either, and docs/review-2026-08.md's `credential
	// added in review-only mode` section for the full experiment record,
	// including why certificate, community and password could not be tested.
	"modify_credential.comment": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <comment> are distinct requests",
	},
}

// CommandWarning attaches a public-facing behavioral-hazard note to a
// generated command's doc comment: wire behavior a caller could not infer
// from the struct's field types alone, and that emitRequest would otherwise
// leave undocumented anywhere but this table and an internal review doc.
type CommandWarning struct {
	Text   string // one or more sentences, emitted verbatim as doc-comment lines
	Reason string // cites the live/upstream evidence
}

// CommandWarnings is keyed by command name. Unlike PointerOverrides and
// FieldOverrides, a warning changes no generated field or type — it exists
// only to surface a hazard the wire shape alone does not reveal.
var CommandWarnings = map[string]CommandWarning{
	"modify_role": {
		Text: "Name, Comment and Users are cleared whenever any of them is omitted from a " +
			"request, not left unchanged: gvmd requires all three together, so a request meant " +
			"to change only one silently blanks the other two.",
		Reason: "T-9: live gvmd 26.36.1 test proved omitted fields are cleared rather than left " +
			"unchanged; see docs/review-2026-08.md's post-tls_certificate revisit section",
	},
	"modify_group": {
		Text: "Name, Comment and Users are cleared whenever any of them is omitted from a " +
			"request, not left unchanged: gvmd requires all three together, so a request meant " +
			"to change only one silently blanks the other two.",
		Reason: "T-9: live gvmd 26.36.1 test proved omitted fields are cleared rather than left " +
			"unchanged; see docs/review-2026-08.md's post-tls_certificate revisit section",
	},
	"modify_permission": {
		Text: "Subject is required by current gvmd even though the schema marks it optional: " +
			"omitting it fails the request with \"400 Error in SUBJECT\", even when only Name or " +
			"Comment is being changed.",
		Reason: "T-9: live gvmd test against the running instance confirmed the omitted-subject " +
			"rejection; see docs/review-2026-08.md's `permission adopted` section",
	},
	"modify_ticket": {
		Text: "Status rejects an explicit-empty <status></status> outright (\"400 Error in " +
			"status\") rather than accepting it as a clear; omitting Status leaves the stored " +
			"value unchanged.",
		Reason: "T-9: live gvmd 26.36.1 test against the running instance confirmed the " +
			"explicit-empty-status rejection; see docs/review-2026-08.md's `ticket` review-only section",
	},
	"modify_credential": {
		Text: "Realm and Kdc each reject an explicit-empty element outright (\"400 Invalid or " +
			"empty kerberos realm value\" / \"400 Invalid or empty kdc value(s)\") rather than " +
			"accepting it as a clear; omitting either leaves the stored value unchanged, so " +
			"neither can take a PointerOverride. Certificate, Community and Password could not be " +
			"tested: gvmd's get_credentials response never echoes their stored values back (with " +
			"or without details=1), so no live test can observe whether omission and " +
			"explicit-empty differ for them.",
		Reason: "T-9: live gvmd 26.36.1 test against the running instance confirmed the " +
			"explicit-empty rejections for realm and kdc, and the absence of certificate, " +
			"community and password from every get_credentials response shape; see " +
			"docs/review-2026-08.md's `credential added in review-only mode` section",
	},
}
