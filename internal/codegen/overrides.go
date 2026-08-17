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
	// boolean_atoi is NOT the same wire rule as boolean: vendor/XML's
	// <type><name>boolean_atoi</name> says the value is run through atoi(3)
	// first, "which means that an empty string is considered false" — the
	// exact opposite of <type><name>boolean</name>'s empty-is-true rule
	// above. modify_config is the first generated family to declare it
	// (family_selection/growing, family/all, family/growing), so it gets its
	// own type rather than reusing GMPBool, which would silently invert
	// every empty element.
	"boolean_atoi": {
		GoType: "commands.GMPAtoiBool",
		Reason: "T-9 config: atoi(3) semantics (empty/non-numeric=false) differ from boolean's empty=true",
	},
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

	// RelocatedFrom is the FieldSuppressions key this entry is the other half
	// of, set only when this addition exists to reproduce a field the schema
	// declared at the wrong nesting level (see FieldSuppressions). Empty for
	// every ordinary FieldOverride, which adds a field the schema simply
	// never declares anywhere.
	RelocatedFrom string
	Reason        string // cites the live/upstream evidence
}

// FieldOverrides is keyed "<command>:/<ancestor>/.../<node-name>" — the same
// addressing style as PointerOverrides's "<command>.<field>", extended with a
// path since this targets a nested struct reached from a command, not only
// a command's top-level request/response envelope. Scope is deliberately
// narrow: a structurally deduplicated shape reached by more than one
// (command, path) use is only allowed when every use resolves to the same
// effective override set (fieldOverridesFor, generate.go) — a use that
// disagrees, or has no entry at all, fails the build rather than leaking one
// use's override onto another. emitStruct and emitRequest both consult this
// table; emitResponse does not — widening it needs its own evidence, not a
// speculative entry now. emitRequest's own widening is modify_report_config's
// report_config_id below: a command's top-level request envelope needed a
// field the schema never declares anywhere, the same shape emitStruct's
// nested structs already had evidence for.
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
	"get_report_formats:/report_format": {{
		Name: "comment", GoType: "string",
		Reason: "T-9 report_format review: schema's report_format pattern never declares <comment>, but " +
			"gvmd emits one on every report format — x5 in " +
			"internal/gmptest/testdata/live/get_report_formats/001_response.xml, x6 in 006_response.xml " +
			"(after a report format copy was added to the fixture set), and x1 upstream, " +
			"all recorded as DROPPED in the matching _fidelity.txt. Every captured occurrence is empty, " +
			"so it is the element's presence that is evidenced, not its content; hence no omitempty, " +
			"which would re-marshal to the absence gvmd never sends",
	}},
	"get_report_configs:/report_config/param": {{
		Name: "using_default", Attr: true, GoType: "GMPBool",
		RelocatedFrom: "get_report_configs:/report_config/param/value",
		Reason: "T-9 report_config review: real position of using_default that " +
			"FieldSuppressions[\"get_report_configs:/report_config/param/value\"] removes — " +
			"command_get_report_configs.xml's own <example> and " +
			"internal/gmptest/testdata/upstream/get_report_configs/001_{response,fidelity}.txt both place " +
			"the attribute on <param>, not <value>",
	}},
	"modify_report_config:/modify_report_config": {{
		Name: "report_config_id", Attr: true, GoType: "string",
		Reason: "T-9 report_config review: command_modify_report_config.xml's <pattern> declares no " +
			"selector attribute at all, unlike every other modify_X command, but a live gvmd 26.36.1 " +
			"test proved gvmd requires one anyway — <modify_report_config report_config_id=\"...\"> " +
			"against a nonexistent id returned \"404 Failed to find Report Config\" (a resource lookup, " +
			"not a parse rejection), reproduced on a second connection. Without this field the generated " +
			"request could never target a specific report config",
	}},
}

// FieldSuppression removes a schema-declared attribute/element from the
// struct its FieldSuppressions key would otherwise put it on, because live
// wire evidence contradicts the schema's declared nesting for it — not
// because the field doesn't exist. RelocatedTo is required and must name the
// FieldOverrides key holding the real position; this table exists for "the
// schema puts this in the wrong place, here's where it really goes", not a
// general-purpose "drop this field" escape hatch.
type FieldSuppression struct {
	Name string // wire element/attribute name to remove
	Attr bool

	RelocatedTo string
	Reason      string
}

// FieldSuppressions is keyed exactly like FieldOverrides:
// "<command>:/<ancestor>/.../<node-name>" identifying the struct the
// suppressed attribute/element would otherwise appear on. Consulted from
// emitStruct, before that struct's fields are built from the schema.
//
// Every entry's RelocatedTo must name a FieldOverrides entry with the same
// (Name, Attr) whose own RelocatedFrom points back at this entry's own key —
// validateFieldRelocations (generate.go) enforces the round trip in both
// directions, so deleting either half without the other fails Generate for
// every family, not just this one. This leaves schemaChildCollision's
// ordinary-collision check untouched: a suppression and its paired addition
// always live at two different node paths (a child and its parent), so they
// never appear in the same schemaChildCollision call.
var FieldSuppressions = map[string][]FieldSuppression{
	"get_report_configs:/report_config/param/value": {{
		Name: "using_default", Attr: true,
		RelocatedTo: "get_report_configs:/report_config/param",
		Reason: "T-9 report_config review: schema's <pattern> declares using_default as an attribute of " +
			"<value>, but command_get_report_configs.xml's own <example> and " +
			"internal/gmptest/testdata/upstream/get_report_configs/001_{response,fidelity}.txt both place " +
			"it on the enclosing <param> instead; decoding a real response through the schema-faithful " +
			"nesting leaves UsingDefault permanently false",
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
	// modify_user: live gvmd 22.7 experiment run against review-only user's
	// deferred question (created a throwaway user with the built-in "User"
	// role, then probed each field via RawXML and read the result back
	// through get_users, or through a fresh authentication attempt for
	// password since get_users never echoes it). comment behaves like
	// modify_alert's (omit leaves the stored value unchanged, explicit-empty
	// clears it to ""). new_name did not: an explicit-empty <new_name></new_name>
	// returns 200 OK but leaves the name unchanged, indistinguishable from
	// omitting the element entirely, so it keeps the mechanical string
	// default. password's *UserPassword already models every legal
	// combination structurally (the whole element can be omitted to leave
	// the password alone; modify="0" and modify="1" with empty chardata both
	// left the stored password unchanged in this gvmd version; modify="1"
	// with real chardata changed it) so it needs no override here — see
	// docs/review-2026-08.md's `user` section for the full experiment
	// record, including the observed anomaly that the prior password kept
	// authenticating after a successful rotation on this stack.
	"modify_user.comment": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <comment> are distinct requests",
	},
	// modify_target: live gvmd 26.36.1 experiment run against review-only
	// target's deferred question (created a throwaway target with unique
	// hosts/exclude_hosts/port_range, then probed each field via RawXML and
	// read comment/name/hosts/exclude_hosts back through get_targets).
	// comment behaves like modify_alert's (omit leaves the stored value
	// unchanged, explicit-empty clears it to ""). name did not: an
	// explicit-empty <name></name> is rejected outright ("400 Name must be
	// at least one character long"), so it keeps the mechanical string
	// default (see the modify_target CommandWarning below). hosts and
	// exclude_hosts could not be tested for an omit-vs-explicit-empty
	// distinction either: gvmd requires both present together whenever
	// either is touched, and rejects both explicit-empty too, so neither
	// field has a legal standalone empty state a pointer could express (same
	// warning). hosts_ordering accepted every value tried but could not be
	// read back at all: get_targets's response never echoes it, so the
	// omit-vs-clear question is unanswerable from this client alone.
	"modify_target.comment": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <comment> are distinct requests",
	},
	// modify_task: live gvmd 26.36.1 experiment run against review-only
	// task's deferred question (created a throwaway container task and a
	// throwaway second user via create_user, then probed comment, name and
	// observers via RawXML and read the result back through get_tasks).
	// comment behaves like modify_alert's (omit leaves the stored value
	// unchanged, explicit-empty clears it to ""). observers needed the
	// throwaway user's login as a real value first (an arbitrary chardata
	// string is rejected with "400 User name error"); once set, it behaves
	// like comment (omit leaves it unchanged, explicit-empty clears it). name
	// did not distinguish the two: an explicit-empty <name></name> returns
	// 200 OK but leaves the stored name unchanged, indistinguishable from
	// omitting the element entirely, so it keeps the mechanical string
	// default. See docs/review-2026-08.md's `task` review-only section for
	// the full experiment record.
	"modify_task.comment": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <comment> are distinct requests",
	},
	"modify_task.observers": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <observers> are distinct requests",
	},
	// modify_scanner: live gvmd 22.7 experiment run against review-only
	// scanner's optional fields (created a throwaway CVE-type scanner with a
	// throwaway "cc" client-certificate credential, since every enabled
	// scanner type on this instance requires one, then probed each field via
	// RawXML and read the result back through get_scanners). comment, name
	// and ca_pub all behave like modify_alert's fields (omit leaves the
	// stored value unchanged, explicit-empty clears it to ""). relay_host
	// does too, but clearing it also resets relay_port to 0 as a side effect
	// -- gvmd treats an empty relay_host as "no relay" and drops the paired
	// port with it. relay_port did not get its own entry: gvmd validates
	// relay_port and relay_host together, rejecting an explicit-empty
	// <relay_port></relay_port> outright ("400 ... RELAY_PORT must be a valid
	// port number ... if RELAY_HOST is not a UNIX socket path") whenever
	// relay_host is still set, so there is no legal standalone empty state
	// for relay_port a pointer could express beyond what relay_host already
	// covers; it keeps the mechanical string default. See
	// docs/review-2026-08.md's `scanner` review-only section for the full
	// experiment record.
	"modify_scanner.comment": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <comment> are distinct requests",
	},
	"modify_scanner.name": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <name> are distinct requests",
	},
	"modify_scanner.ca_pub": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <ca_pub> are distinct requests",
	},
	"modify_scanner.relay_host": {
		Reason: "T-9: live gvmd test proved omitted vs. explicit-empty <relay_host> are distinct requests",
	},
}

// CommandWarning attaches a public-facing behavioral-hazard note to a
// generated command's doc comment: wire behavior a caller could not infer
// from the struct's field types alone, and that emitRequest/emitResponse
// would otherwise leave undocumented anywhere but this table and an internal
// review doc.
type CommandWarning struct {
	Text   string // one or more sentences, emitted verbatim as doc-comment lines
	Reason string // cites the live/upstream evidence
}

// CommandWarnings is keyed by command name for a request-side hazard, or by
// "<command>_response" for a response-side one (attached by emitResponse to
// the *Response struct's doc comment — see get_assets_response below for a
// field-typing hazard that only affects decoding). Unlike PointerOverrides
// and FieldOverrides, a warning changes no generated field or type — it
// exists only to surface a hazard the wire shape alone does not reveal.

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
	"modify_port_list": {
		Text: "Name and Comment are cleared whenever either is omitted from a request, not left " +
			"unchanged: modifying one alone still blanks the other, so both must be resent " +
			"together to preserve them.",
		Reason: "T-9: live gvmd 26.36.1 test against the running instance confirmed omitted " +
			"Name/Comment are cleared rather than left unchanged (both together, and Name alone " +
			"while only Comment was resent); see docs/review-2026-08.md's `port_list added in " +
			"review-only mode` section",
	},
	"modify_asset": {
		Text: "Comment is cleared whenever it is omitted from a request, not left unchanged: " +
			"there is no way to leave a stored comment untouched other than resending it.",
		Reason: "T-9: live gvmd 26.36.1 test against the running instance confirmed an omitted " +
			"Comment is cleared rather than left unchanged; see docs/review-2026-08.md's `asset " +
			"added in review-only mode` section",
	},
	"get_assets_response": {
		Text: "Asset.Identifiers.Identifier.Source.Name follows the schema's boolean type, but " +
			"current gvmd sends a real username (for example \"admin\") in that element. GMPBool " +
			"decodes any non-\"0\" text as true, so the username is silently discarded rather than " +
			"preserved.",
		Reason: "T-9: live gvmd 26.36.1 test against the running instance returned a username in " +
			"this element; see docs/review-2026-08.md's `asset added in review-only mode` section",
	},
	"modify_report_config": {
		Text: "Sending an id attribute instead of ReportConfigID's report_config_id crashes gvmd's " +
			"request handler outright (confirmed SIGSEGV), rather than failing the request. Whether " +
			"omitting vs. explicit-emptying Name or Comment leaves the stored value unchanged (like " +
			"modify_alert) or clears it (like modify_role) could not be tested: this environment has no " +
			"configurable report format and no pre-existing report config to create or modify against.",
		Reason: "T-9 report_config review: live gvmd 26.36.1 test — <modify_report_config " +
			"id=\"...\"> crashed the request-handling process (container logs recorded a SIGSEGV " +
			"backtrace through libglib's markup parser); all 5 installed report formats returned " +
			"configurable=false and get_report_configs returned zero existing configs, so no valid " +
			"target existed to exercise the Name/Comment omit-vs-explicit-empty question against",
	},
	"modify_target": {
		Text: "Name rejects an explicit-empty <name></name> outright (\"400 Name must be at " +
			"least one character long\") rather than accepting it as a clear; omitting Name leaves " +
			"the stored value unchanged. Hosts and ExcludeHosts must always be sent together: " +
			"omitting either while the other is present in the same request fails (\"400 HOSTS " +
			"requires an EXCLUDE_HOSTS\" / \"400 EXCLUDE_HOSTS requires a HOSTS\"), even when the " +
			"omitted field's stored value is not changing, and sending both explicit-empty is " +
			"rejected too (\"400 HOSTS must be at least one character long\").",
		Reason: "T-9: live gvmd 26.36.1 test against the running instance confirmed the name " +
			"rejection and the hosts/exclude_hosts pairing requirement; see docs/review-2026-08.md's " +
			"`target` review-only section",
	},
}

// DefectKey identifies one schema.Defect (internal/schema/parse.go's
// Command.Defects) precisely enough to decide whether to generate through it:
// the command whose snapshot carries it, the defect's class
// (schema.DefectEmptyPattern and friends), and its path.
type DefectKey struct {
	Command string
	Class   string
	Path    string
}

// DefectAllowlistEntry records why a schema defect is generated through
// rather than refused. Unlike Override/FieldOverride/PointerOverride, an
// entry does not alter the parser-derived node or the generated output; it
// only permits generation despite the recorded defect, so a human decision is
// on record instead of the generator silently letting it through.
type DefectAllowlistEntry struct {
	Reason string
}

// DefectAllowlist is Generate's fail-closed gate on schema.Command.Defects
// (checkDefects, generate.go): a defect whose {Command, Class, Path} is not a
// key here fails generation, and a stale entry fails it too
// (checkDefectAllowlistConsumed). Each entry records a specific, reviewed
// defect a family's adoption accepted; see docs/review-2026-08.md's T-15
// entry for the full evidence log — per T-20, a Reason here states the
// conclusion and cites the section, one line.
var DefectAllowlist = map[DefectKey]DefectAllowlistEntry{
	{Command: "create_group", Class: "empty-pattern", Path: "/create_group/specials/full"}: {
		Reason: "T-15; docs/review-2026-08.md `group adopted`",
	},
	{Command: "get_tags", Class: "empty-pattern", Path: "/get_tags_response/tag/resources/resource/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `tag`",
	},
	{Command: "get_schedules", Class: "empty-pattern", Path: "/get_schedules_response/schedule/tasks/task/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `schedule`",
	},
	{Command: "get_alerts", Class: "empty-pattern", Path: "/get_alerts_response/alert/filter/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `alert`",
	},
	{Command: "get_alerts", Class: "empty-pattern", Path: "/get_alerts_response/alert/tasks/task/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `alert`",
	},
	{Command: "get_filters", Class: "empty-pattern", Path: "/get_filters_response/filter/alerts/alert/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `filter`",
	},
	{Command: "get_permissions", Class: "empty-pattern", Path: "/get_permissions_response/permission/resource/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `permission`",
	},
	{Command: "get_permissions", Class: "empty-pattern", Path: "/get_permissions_response/permission/subject/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `permission`",
	},
	{Command: "get_credentials", Class: "empty-pattern", Path: "/get_credentials_response/credential/scanners/scanner/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `credential`",
	},
	{Command: "get_credentials", Class: "empty-pattern", Path: "/get_credentials_response/credential/targets/target/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `credential`",
	},
	{Command: "get_port_lists", Class: "empty-pattern", Path: "/get_port_lists_response/port_list/targets/target/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `port_list`",
	},
	{Command: "get_notes", Class: "duplicate-ele", Path: "/get_notes_response/note/active"}: {
		Reason: "T-15; upstream `<ele>` redefined in one scope, docs/review-2026-08.md T-15 defect table",
	},
	{Command: "get_notes", Class: "duplicate-ele", Path: "/get_notes_response/note/result/delta/overrides/override/active"}: {
		Reason: "T-15; upstream `<ele>` redefined in one scope, docs/review-2026-08.md T-15 defect table",
	},
	{Command: "get_overrides", Class: "duplicate-ele", Path: "/get_overrides_response/override/active"}: {
		Reason: "T-15; upstream `<ele>` redefined in one scope, docs/review-2026-08.md T-15 defect table",
	},
	{Command: "get_overrides", Class: "duplicate-ele", Path: "/get_overrides_response/override/result/delta/notes/note/active"}: {
		Reason: "T-15; upstream `<ele>` redefined in one scope, docs/review-2026-08.md T-15 defect table",
	},
	{Command: "get_notes", Class: "duplicate-ele", Path: "/get_notes_response/note/result/overrides/override/active"}: {
		Reason: "T-21; same upstream redefinition as the delta-path override reference, docs/review-2026-08.md T-21",
	},
	{Command: "get_overrides", Class: "duplicate-ele", Path: "/get_overrides_response/override/result/notes/note/active"}: {
		Reason: "T-21; same upstream redefinition as the delta-path note reference, docs/review-2026-08.md T-21",
	},
	{Command: "get_assets", Class: "undefined-type", Path: "/get_assets_response/asset/host/routes/route/host/@distance"}: {
		Reason: "T-15; docs/review-2026-08.md `asset`",
	},
	{Command: "create_port_list", Class: "unsupported-r", Path: "/create_port_list"}: {
		Reason: "T-15; docs/review-2026-08.md `port_list`",
	},
	{Command: "create_report_format", Class: "unsupported-r", Path: "/create_report_format"}: {
		Reason: "T-15; docs/review-2026-08.md `report_format`",
	},
	{Command: "get_report_configs", Class: "empty-pattern", Path: "/get_report_configs_response/report_config/report_format/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `report_config`",
	},
	{Command: "get_users", Class: "empty-pattern", Path: "/get_users_response/user/role/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `user`",
	},
	{Command: "get_users", Class: "empty-pattern", Path: "/get_users_response/user/groups/group/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `user`",
	},
	{Command: "get_targets", Class: "duplicate-ref", Path: "/get_targets_response/target/permissions"}: {
		Reason: "T-9; docs/review-2026-08.md `target`",
	},
	{Command: "get_targets", Class: "empty-pattern", Path: "/get_targets_response/target/esxi_credential/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `target`",
	},
	{Command: "get_targets", Class: "empty-pattern", Path: "/get_targets_response/target/krb5_credential/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `target`",
	},
	{Command: "get_targets", Class: "empty-pattern", Path: "/get_targets_response/target/port_list/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `target`",
	},
	{Command: "get_targets", Class: "empty-pattern", Path: "/get_targets_response/target/smb_credential/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `target`",
	},
	{Command: "get_targets", Class: "empty-pattern", Path: "/get_targets_response/target/snmp_credential/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `target`",
	},
	{Command: "get_targets", Class: "empty-pattern", Path: "/get_targets_response/target/ssh_credential/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `target`",
	},
	{Command: "get_targets", Class: "empty-pattern", Path: "/get_targets_response/target/ssh_elevate_credential/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `target`",
	},
	{Command: "get_targets", Class: "empty-pattern", Path: "/get_targets_response/target/tasks/task/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `target`",
	},
	{Command: "get_tasks", Class: "duplicate-ref", Path: "/get_tasks_response/task/name"}: {
		Reason: "T-9; docs/review-2026-08.md `task`",
	},
	{Command: "get_tasks", Class: "duplicate-ref", Path: "/get_tasks_response/task/schedule"}: {
		Reason: "T-9; docs/review-2026-08.md `task`",
	},
	{Command: "get_tasks", Class: "duplicate-ref", Path: "/get_tasks_response/task/schedule_periods"}: {
		Reason: "T-9; docs/review-2026-08.md `task`",
	},
	{Command: "get_tasks", Class: "empty-pattern", Path: "/get_tasks_response/task/alert/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `task`",
	},
	{Command: "get_tasks", Class: "empty-pattern", Path: "/get_tasks_response/task/config/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `task`",
	},
	{Command: "get_tasks", Class: "empty-pattern", Path: "/get_tasks_response/task/scanner/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `task`",
	},
	{Command: "get_tasks", Class: "empty-pattern", Path: "/get_tasks_response/task/target/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `task`",
	},
	{Command: "get_scanners", Class: "empty-pattern", Path: "/get_scanners_response/scanner/configs/config/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `scanner`",
	},
	{Command: "get_scanners", Class: "empty-pattern", Path: "/get_scanners_response/scanner/tasks/task/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `scanner`",
	},
	{Command: "get_scanners", Class: "empty-pattern", Path: "/get_scanners_response/scanner/info/scanner/name"}: {
		Reason: "T-9; docs/review-2026-08.md `scanner`",
	},
	{Command: "get_scanners", Class: "empty-pattern", Path: "/get_scanners_response/scanner/info/scanner/version"}: {
		Reason: "T-9; docs/review-2026-08.md `scanner`",
	},
	{Command: "get_scanners", Class: "empty-pattern", Path: "/get_scanners_response/scanner/info/daemon/name"}: {
		Reason: "T-9; docs/review-2026-08.md `scanner`",
	},
	{Command: "get_scanners", Class: "empty-pattern", Path: "/get_scanners_response/scanner/info/daemon/version"}: {
		Reason: "T-9; docs/review-2026-08.md `scanner`",
	},
	{Command: "get_scanners", Class: "empty-pattern", Path: "/get_scanners_response/scanner/info/protocol/name"}: {
		Reason: "T-9; docs/review-2026-08.md `scanner`",
	},
	{Command: "get_scanners", Class: "empty-pattern", Path: "/get_scanners_response/scanner/info/protocol/version"}: {
		Reason: "T-9; docs/review-2026-08.md `scanner`",
	},
	{Command: "get_scanners", Class: "no-pattern", Path: "/get_scanners_response/scanner/info/description"}: {
		Reason: "T-9; docs/review-2026-08.md `scanner`",
	},
	{Command: "get_scanners", Class: "no-pattern", Path: "/get_scanners_response/scanner/info/params/param"}: {
		Reason: "T-9; docs/review-2026-08.md `scanner`",
	},
	{Command: "create_config", Class: "unsupported-r", Path: "/create_config"}: {
		Reason: "T-15; docs/review-2026-08.md `config`",
	},
	{Command: "get_configs", Class: "empty-pattern", Path: "/get_configs_response/config/tasks/task/permissions"}: {
		Reason: "T-15; docs/review-2026-08.md `config`",
	},
}
