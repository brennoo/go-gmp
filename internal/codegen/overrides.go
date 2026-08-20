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
// use's override onto another. emitStruct, emitRequest, and emitResponse all
// consult this table. emitRequest's own widening is modify_report_config's
// report_config_id below: a command's top-level request envelope needed a
// field the schema never declares anywhere, the same shape emitStruct's
// nested structs already had evidence for. emitResponse's widening is
// verify_scanner's version below: live gvmd traffic confirmed the field is a
// real, populated wire value the schema's <pattern> never declares (only its
// <example> shows it) — dropping it would be a real behavior regression, not
// a defect fix.
// t18aSortFieldReason is the shared Reason for every T-18a sort-chardata
// relocation below, in both this table and FieldSuppressions: every affected
// get_X command, and the fidelity evidence behind each one, is already
// listed in docs/review-2026-08.md's T-18a entry, so the table itself only
// needs to name the defect, not repeat that evidence 25 times over.
const t18aSortFieldReason = "T-18a: sort-field chardata is nested under <field>"

// t23CountTotalReason is the shared Reason for every T-23 count-wrapper
// chardata addition below: every affected get_X command, the schema/wire
// evidence, and the zero-total citation are recorded once in
// docs/review-2026-08.md's T-23 entry rather than repeated 25 times over.
const t23CountTotalReason = "T-23: count-wrapper chardata total is schema-undeclared"

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
	"verify_scanner:/verify_scanner_response": {{
		Name: "version", GoType: "string", OmitEmpty: true,
		Reason: "T-9 scanner review: command_verify_scanner.xml's response <pattern> declares only " +
			"status/status_text (its <example> shows <version> but the pattern never does), yet a live " +
			"gvmd 26.36.1 test against two real scanners returned a populated <version> element both " +
			"times (\"OpenVAS 23.50.0\", \"GVM/26.36.1\") — real, populated wire data, so dropping it " +
			"would be a behavior regression rather than a defect fix",
	}},
	"modify_auth:/group": {{
		Name: "name", Attr: true, GoType: "string",
		Reason: "T-9 auth review: command_modify_auth.xml's request <pattern> for <group> declares no " +
			"name attrib at all, but the command's own <example> requires <group name=\"Foo\"> to " +
			"identify which group to modify — the schema's pattern and its own example disagree, and " +
			"without this field the generated request could never target a specific group",
	}},
	"get_aggregates:/aggregate/group": {
		{
			Name: "min", GoType: "string", OmitEmpty: true,
			Reason: "T-9 aggregate review: command_get_aggregates.xml's <group> pattern declares only " +
				"value/subgroup/count/c_count/stats/text — per-column stats wrapped in <stats column=\"...\">" +
				" — but the command's own <example> and internal/gmptest/testdata/upstream/get_aggregates/" +
				"001_response.xml both show a flat <min> directly under <group> for the single-data_column " +
				"case, the same pattern-vs-example disagreement as verify_scanner's version",
		},
		{
			Name: "max", GoType: "string", OmitEmpty: true,
			Reason: "T-9 aggregate review: see the min entry above; same flat <max> shape confirmed by the " +
				"same evidence",
		},
		{
			Name: "mean", GoType: "string", OmitEmpty: true,
			Reason: "T-9 aggregate review: see the min entry above; same flat <mean> shape confirmed by " +
				"the same evidence",
		},
		{
			Name: "sum", GoType: "string", OmitEmpty: true,
			Reason: "T-9 aggregate review: see the min entry above; same flat <sum> shape confirmed by the " +
				"same evidence",
		},
		{
			Name: "c_sum", GoType: "string", OmitEmpty: true,
			Reason: "T-9 aggregate review: see the min entry above; same flat <c_sum> shape confirmed by " +
				"the same evidence",
		},
	},
	"get_aggregates:/aggregate/group/text": {{
		Name: "column", Attr: true, GoType: "string", OmitEmpty: true,
		Reason: "T-9 aggregate review: command_get_aggregates.xml's schema declares <text>'s lone " +
			"attribute as \"name\", but internal/gmptest/testdata/live/get_aggregates/004_response.xml " +
			"sends a real gvmd 26.36.1 <text column=\"name\"> — a second, undeclared attribute the schema " +
			"never mentions at all, kept alongside the schema-declared (and seemingly unused) Name field " +
			"since there is no evidence the schema's own declaration is categorically wrong rather than " +
			"an alternate/legacy form. AggregateText is one structurally-deduplicated type shared by both " +
			"aggregate/group/text and aggregate/overall/text, so both must carry this same entry",
	}},
	"get_aggregates:/aggregate/overall/text": {{
		Name: "column", Attr: true, GoType: "string", OmitEmpty: true,
		Reason: "T-9 aggregate review: command_get_aggregates.xml's schema declares <text>'s lone " +
			"attribute as \"name\", but internal/gmptest/testdata/live/get_aggregates/004_response.xml " +
			"sends a real gvmd 26.36.1 <text column=\"name\"> — a second, undeclared attribute the schema " +
			"never mentions at all, kept alongside the schema-declared (and seemingly unused) Name field " +
			"since there is no evidence the schema's own declaration is categorically wrong rather than " +
			"an alternate/legacy form. AggregateText is one structurally-deduplicated type shared by both " +
			"aggregate/group/text and aggregate/overall/text, so both must carry this same entry",
	}},
	"get_settings:/setting": {{
		Name: "id", Attr: true, GoType: "string",
		Reason: "T-9 settings review: command_get_settings.xml's formal <setting> pattern declares only " +
			"name/comment/value, no id attribute, but internal/gmptest/testdata/live/get_settings/" +
			"{001,002,003}_response.xml all show every <setting id=\"...\"> carrying a real UUID — the " +
			"schema's own baked <example> agrees a setting has an id, it just nests <setting> inside " +
			"<settings> in a way no live capture ever confirms",
	}},
	"get_settings:/settings": {{
		Name: "setting", GoType: "[]SettingsSetting", OmitEmpty: true,
		Reason: "T-9 settings review: command_get_settings.xml's <settings> pattern declares only the " +
			"start/max pagination attributes, but the command's own <example> — reproduced verbatim in " +
			"testdata/upstream/get_settings/001_response.xml — nests a real <setting id=\"...\"> inside " +
			"<settings>, the same pattern-vs-example disagreement as verify_scanner's version; without " +
			"this field the round trip silently drops that nested setting (confirmed by the upstream " +
			"fixture's DROPPED settings/setting finding). No live capture against current gvmd ever " +
			"produces this nested shape — every live fixture uses the flat top-level <setting> siblings " +
			"already modelled by GetSettingsResponse.Setting — so this is an alternate/legacy wire form " +
			"kept alongside the flat one, not a replacement for it",
	}},
	"get_system_reports:/get_system_reports_response": {{
		Name: "truncated", GoType: "string", OmitEmpty: true,
		Reason: "T-9 system_report review: command_get_system_reports.xml's response <pattern> never " +
			"declares <truncated>, but the command's own <example> — reproduced verbatim in " +
			"testdata/upstream/get_system_reports/001_response.xml — shows a real <truncated>...</truncated> " +
			"sibling of <system_report>, the same pattern-vs-example disagreement as get_settings' nested " +
			"setting; no live capture against current gvmd ever populates it (truncation only occurs on " +
			"long-running reports the pinned stack's fixtures never triggered), so this is weaker evidence " +
			"than most restored fields here, but it follows the same precedent of trusting the schema's own " +
			"baked example over a formal pattern that has already been shown to omit real fields elsewhere " +
			"in this command",
	}},
	"get_system_reports:/system_report/report": {{
		Name: "start_time", Attr: true, GoType: "*GMPTime", OmitEmpty: true,
		Reason: "T-9 system_report review: command_get_system_reports.xml's <report> pattern never " +
			"declares start_time, but internal/gmptest/testdata/live/get_system_reports/001_response.xml " +
			"sends a real start_time=\"0001-01-01T00:00:00Z\" attribute — undeclared but genuine wire data, " +
			"confirmed dropped by that fixture's own pre-adoption fidelity report; " +
			"testdata/upstream/get_system_reports/002_response.xml proves the attribute is also genuinely " +
			"absent on other reports, so it is pointer-typed (unlike this codebase's usual non-pointer " +
			"GMPTime convention) to actually honor omitempty on re-encode rather than fabricate a zero " +
			"timestamp for a report that never had one",
	}, {
		Name: "end_time", Attr: true, GoType: "*GMPTime", OmitEmpty: true,
		Reason: "T-9 system_report review: same undeclared-but-real, genuinely optional attribute as " +
			"start_time above, same evidence and same pointer-typing rationale",
	}},
	"get_vulns:/vuln": {{
		Name: "type", GoType: "string",
		Reason: "T-9 vuln review: command_get_vulns.xml's vuln <pattern> list omits <e>type</e>, but the " +
			"command still declares a full <ele><name>type</name></ele> body for it a few lines down (alts " +
			"nvt/cve) — a dangling definition the pattern forgot to reference, not evidence the field " +
			"doesn't exist; every fixture (live and upstream) sends a real <type>nvt</type>, the same " +
			"pattern-omits-a-fully-defined-element class as verify_scanner's version",
	}},

	"get_alerts:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_alerts:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_assets:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_assets:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_configs:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_configs:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_credentials:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_credentials:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_filters:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_filters:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_groups:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_groups:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_info:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_info:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_notes:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_notes:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_overrides:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_overrides:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_permissions:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_permissions:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_port_lists:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_port_lists:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_report_configs:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_report_configs:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_report_formats:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_report_formats:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_reports:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_reports:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_results:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_results:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_roles:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_roles:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_scanners:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_scanners:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_schedules:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_schedules:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_tags:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_tags:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_targets:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_targets:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_tasks:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_tasks:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_tickets:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_tickets:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_tls_certificates:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_tls_certificates:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_users:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_users:/sort",
		Reason: t18aSortFieldReason,
	}},
	"get_vulns:/sort/field": {{
		Chardata: true, GoType: "string", RelocatedFrom: "get_vulns:/sort",
		Reason: t18aSortFieldReason,
	}},

	"get_roles:/role_count": {{
		Chardata: true, GoType: "int",
		Reason: "T-23: command_get_roles.xml's role_count <pattern> declares only <e>filtered</e>/" +
			"<e>page</e>, no chardata, but the command's own baked <example> shows a bare leading total " +
			"(\"1<filtered>1</filtered><page>1</page>\") and every live fixture agrees — " +
			"internal/gmptest/testdata/live/get_roles/001_response.xml sends " +
			"<role_count>6<filtered>6</filtered><page>6</page></role_count>. A zero total is equally real: " +
			"internal/gmptest/testdata/live/get_alerts/001_response.xml's matching " +
			"<alert_count>0<filtered>0</filtered><page>0</page></alert_count> proves the leading text can be " +
			"\"0\", which rules out modeling this with omitempty (an int chardata field with omitempty " +
			"added to its xml tag drops the whole element on marshal, zero or not — confirmed by a " +
			"standalone encoding/xml experiment — so this and every entry below omits OmitEmpty entirely)",
	}},
	// The remaining 24 commands below share get_roles' exact defect: same
	// pattern-declares-only-filtered/page shape, same undeclared leading
	// total, confirmed independently by each command's own schema and
	// fidelity/live evidence.
	"get_alerts:/alert_count":                     {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_assets:/asset_count":                     {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_configs:/config_count":                   {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_credentials:/credential_count":           {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_filters:/filter_count":                   {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_groups:/group_count":                     {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_info:/info_count":                        {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_notes:/note_count":                       {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_overrides:/override_count":               {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_permissions:/permission_count":           {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_port_lists:/port_list_count":             {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_report_configs:/report_config_count":     {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_report_formats:/report_format_count":     {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_reports:/report_count":                   {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_results:/result_count":                   {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_scanners:/scanner_count":                 {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_schedules:/schedule_count":               {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_tags:/tag_count":                         {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_targets:/target_count":                   {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_tasks:/task_count":                       {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_tickets:/ticket_count":                   {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_tls_certificates:/tls_certificate_count": {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_users:/user_count":                       {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
	"get_vulns:/vuln_count":                       {{Chardata: true, GoType: "int", Reason: t23CountTotalReason}},
}

// FieldTypeOverride corrects the Go type the mechanical mapping would assign
// to a schema-declared attribute/element, at one specific nesting path. It
// changes nothing about the field's presence, position or wire name — for
// that, see FieldOverride (adds a field) and FieldSuppression (relocates
// one). This exists because TypeOverrides is keyed by the schema's <type>
// name globally, which is too broad when the schema declares the *same*
// element name with two different, inconsistent types depending on where it
// is nested: overriding "integer" everywhere would silently retype every
// other integer field in the corpus.
type FieldTypeOverride struct {
	Name   string // wire element/attribute name
	Attr   bool
	GoType string
	Reason string
}

// FieldTypeOverrides is keyed exactly like FieldOverrides/FieldSuppressions:
// "<command>:/<ancestor>/.../<node-name>" identifying the struct the retyped
// field lives on. Consulted from fieldOf, after the mechanical type mapping,
// so the corrected type still goes through the normal cardinality/pointer
// logic.
var FieldTypeOverrides = map[string][]FieldTypeOverride{
	"get_system_reports:/system_report/report": {{
		Name: "duration", Attr: true, GoType: "string",
		Reason: "T-9 system_report review: command_get_system_reports.xml declares report/@duration as a " +
			"required integer, but internal/gmptest/testdata/live/get_system_reports/001_response.xml sends " +
			"duration=\"\" (empty string) for it — a plain int silently decodes that as zero and then " +
			"re-encodes as \"0\", changing the wire value on round trip (confirmed by that fixture's own " +
			"pre-adoption CHANGED finding); same class of schema/wire type disagreement as get_info's id " +
			"below, retyped to string to preserve the exact captured value",
	}},
	"create_report:/report/report/results/result/nvt": {{
		Name: "cvss_base", GoType: "float64",
		Reason: "T-9 report review: command_get_reports.xml's schema declares " +
			"results/result/nvt/cvss_base as type integer, but " +
			"internal/gmptest/testdata/live/create_report/001_request.xml sends a real decimal value " +
			"(\"5.0\") for it — the same element name is declared type severity (float64) a few levels " +
			"over at errors/error/nvt/cvss_base, so this is the schema itself disagreeing with its own " +
			"declaration, not a defect in the generated int",
	}},
	"get_reports:/report/report/results/result/nvt": {{
		Name: "cvss_base", GoType: "float64",
		Reason: "T-9 report review: same schema/wire disagreement as create_report's entry above — " +
			"internal/gmptest/testdata/upstream/get_reports/002_response.xml sends the same decimal " +
			"cvss_base shape at this path",
	}},
	"get_results:/result/nvt": {{
		Name: "cvss_base", GoType: "float64",
		Reason: "T-9 result review, finding 2: same schema/wire disagreement as create_report/get_reports " +
			"above — command_get_results.xml's served schema declares result/nvt/cvss_base as type integer, " +
			"but internal/gmptest/testdata/upstream/get_results/001_response.xml sends a real decimal value " +
			"(\"5.0\") for it, which a plain int would truncate on round trip. Not every fixture agrees on " +
			"shape — upstream 002 sends a whole number (\"9\") and all 4 live fixtures send \"0.0\" — but " +
			"upstream 001's decimal alone is sufficient to rule out int",
	}},
	"get_info:/info/cve/configuration_nodes/node": {{
		Name: "node", GoType: "NVTNode",
		Reason: "T-9 nvt review: node/node is a genuine self-reference (see DefectAllowlist's " +
			"undefined-type entry for this path) that the mechanical mapping falls back to string for; " +
			"Go structs may reference themselves through a slice field, so retyping it back to its own " +
			"generated struct name is a correct closed-form fix, not a hand patch",
	}},
	"get_info:/info": {{
		Name: "id", Attr: true, GoType: "string",
		Reason: "T-9 nvt review: command_get_info.xml declares info/@id as type integer, and the " +
			"live-captured schema snapshot agrees, but " +
			"internal/gmptest/testdata/live/get_info/001_response.xml's type=cve lookups return the CVE " +
			"name itself (\"CVE-1999-0001\") as @id, not a numeric ID — real wire data the generated *int " +
			"can't even decode",
	}},
	"get_resource_names:/resource": {{
		Name: "id", Attr: true, GoType: "string",
		Reason: "T-9 settings review: command_get_resource_names.xml declares resource/@id as type " +
			"integer, but its own baked <example> and internal/gmptest/testdata/live/get_resource_names/" +
			"{001,002}_response.xml and testdata/upstream/get_resource_names/001_response.xml all show a " +
			"UUID string (e.g. \"5b6b6aef-b320-42ca-899f-3161ee2a4195\") — real wire data the generated " +
			"*int can't even decode, same class as get_info's id entry above",
	}},
}

// FieldSuppression removes a schema-declared attribute/element, or a
// struct's own schema-declared chardata, from the struct its
// FieldSuppressions key would otherwise put it on, because live wire
// evidence contradicts the schema's declared nesting for it — not because
// the field doesn't exist. RelocatedTo is required and must name the
// FieldOverrides key holding the real position; this table exists for "the
// schema puts this in the wrong place, here's where it really goes", not a
// general-purpose "drop this field" escape hatch.
type FieldSuppression struct {
	Name     string // wire element/attribute name to remove; ignored when Chardata
	Attr     bool
	Chardata bool // this struct's own schema-declared `,chardata` field; mutually exclusive with Name/Attr

	RelocatedTo string
	Reason      string
}

// FieldSuppressions is keyed exactly like FieldOverrides:
// "<command>:/<ancestor>/.../<node-name>" identifying the struct the
// suppressed attribute/element would otherwise appear on. Consulted from
// emitStruct, before that struct's fields are built from the schema.
//
// Every entry's RelocatedTo must name a FieldOverrides entry with the same
// (Name, Attr, Chardata) whose own RelocatedFrom points back at this entry's
// own key — validateFieldRelocations (generate.go) enforces the round trip
// in both directions, so deleting either half without the other fails
// Generate for every family, not just this one. This leaves
// schemaChildCollision's ordinary-collision check untouched: a suppression
// and its paired addition always live at two different node paths (a child
// and its parent), so they never appear in the same schemaChildCollision
// call.
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

	"get_alerts:/sort": {{
		Chardata: true, RelocatedTo: "get_alerts:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_assets:/sort": {{
		Chardata: true, RelocatedTo: "get_assets:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_configs:/sort": {{
		Chardata: true, RelocatedTo: "get_configs:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_credentials:/sort": {{
		Chardata: true, RelocatedTo: "get_credentials:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_filters:/sort": {{
		Chardata: true, RelocatedTo: "get_filters:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_groups:/sort": {{
		Chardata: true, RelocatedTo: "get_groups:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_info:/sort": {{
		Chardata: true, RelocatedTo: "get_info:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_notes:/sort": {{
		Chardata: true, RelocatedTo: "get_notes:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_overrides:/sort": {{
		Chardata: true, RelocatedTo: "get_overrides:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_permissions:/sort": {{
		Chardata: true, RelocatedTo: "get_permissions:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_port_lists:/sort": {{
		Chardata: true, RelocatedTo: "get_port_lists:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_report_configs:/sort": {{
		Chardata: true, RelocatedTo: "get_report_configs:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_report_formats:/sort": {{
		Chardata: true, RelocatedTo: "get_report_formats:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_reports:/sort": {{
		Chardata: true, RelocatedTo: "get_reports:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_results:/sort": {{
		Chardata: true, RelocatedTo: "get_results:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_roles:/sort": {{
		Chardata: true, RelocatedTo: "get_roles:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_scanners:/sort": {{
		Chardata: true, RelocatedTo: "get_scanners:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_schedules:/sort": {{
		Chardata: true, RelocatedTo: "get_schedules:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_tags:/sort": {{
		Chardata: true, RelocatedTo: "get_tags:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_targets:/sort": {{
		Chardata: true, RelocatedTo: "get_targets:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_tasks:/sort": {{
		Chardata: true, RelocatedTo: "get_tasks:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_tickets:/sort": {{
		Chardata: true, RelocatedTo: "get_tickets:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_tls_certificates:/sort": {{
		Chardata: true, RelocatedTo: "get_tls_certificates:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_users:/sort": {{
		Chardata: true, RelocatedTo: "get_users:/sort/field",
		Reason: t18aSortFieldReason,
	}},
	"get_vulns:/sort": {{
		Chardata: true, RelocatedTo: "get_vulns:/sort/field",
		Reason: t18aSortFieldReason,
	}},
}

// NameOverride assigns an explicit Go type name to a structurally-collected
// shape, bypassing pick's candidate ladder (generate.go) entirely.
//
// The ladder's qualify() step treats a node already starting with the family
// noun as self-qualified and leaves it bare — the rule that keeps
// <note_count> as NoteCount rather than NoteNoteCount. That rule is correct
// for every other family only because English happens not to prefix an
// unrelated sibling family's name with theirs. "report" breaks it:
// report_format and report_config are also the top-level entity names of two
// OTHER adopted families (get_report_formats.go's ReportFormat,
// get_report_configs.go's ReportConfig). qualify()'s short-circuit hands the
// report family's own nested <report_format>/<report_config> elements those
// same bare names, which is a real Go "redeclared in this block" error
// between two files of the same commands package, not just an awkward name.
//
// GenerateAll's family-local isolation (its own doc comment, and
// TestGenerateAllIsPerFamily/TestRealFamiliesGenerateIndependently) rules out
// fixing this with a smarter ladder heuristic or a cross-family name
// registry: a family's generated names may depend on nothing but that
// family's own commands, and no local string rule can distinguish this
// collision from note_count's legitimate collapse without regressing it (see
// docs/review-2026-08.md `report`). This table is the human-curated
// exception - naming a specific shape from knowledge the generator itself is
// not allowed to consult.
type NameOverride struct {
	GoName string
	Reason string
}

// NameOverrides is keyed exactly like FieldOverrides (fieldOverrideKey):
// "<command>:/<ancestor>/.../<node-name>". Consulted by pick (generate.go)
// before the candidate ladder runs at all.
var NameOverrides = map[string]NameOverride{
	"create_report:/report/report_format": {
		GoName: "ReportReportFormat",
		Reason: "T-9; docs/review-2026-08.md `report` — the top-level <report_format> id+name " +
			"reference would otherwise take the bare name get_report_formats.go's own ReportFormat " +
			"entity type already declares",
	},
	"get_reports:/report/report_format": {
		GoName: "ReportReportFormat",
		Reason: "T-9; docs/review-2026-08.md `report` — same collision as create_report's",
	},
	"create_report:/report/report/report_format": {
		GoName: "ReportReportFormatParams",
		Reason: "T-9; docs/review-2026-08.md `report` — the nested <report/report_format> param list " +
			"would otherwise also contend for the bare \"ReportFormat\" name; named for its actual " +
			"param-list shape to stay distinct from the id+name reference above",
	},
	"get_reports:/report/report/report_format": {
		GoName: "ReportReportFormatParams",
		Reason: "T-9; docs/review-2026-08.md `report` — same collision as create_report's",
	},
	"create_report:/report/report_config": {
		GoName: "ReportReportConfig",
		Reason: "T-9; docs/review-2026-08.md `report` — the <report_config> id+name reference would " +
			"otherwise take the bare name get_report_configs.go's own ReportConfig entity type " +
			"already declares",
	},
	"get_reports:/report/report_config": {
		GoName: "ReportReportConfig",
		Reason: "T-9; docs/review-2026-08.md `report` — same collision as create_report's",
	},
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

// RequiredOverride keeps a schema-required field's Go tag carrying `,omitempty`
// despite that, because live protocol evidence contradicts the schema's
// declared requiredness for that specific field: gvmd genuinely omits it under
// some response status, so a non-omitempty Go field would fabricate the
// element on re-encode. Unlike PointerOverride, which only ever adjusts the
// pointer-vs-value choice for a field the schema already marks optional, this
// overrides fieldOf's requiredness-derived omitempty decision itself.
//
// When the named field is structured (has children), fieldOf also makes it a
// pointer — the wrapper's own zero value (e.g. a *_count summary struct with
// Value/Filtered/Page all 0) is otherwise indistinguishable from "genuinely
// absent," so a plain (non-pointer) struct field would fabricate the element
// on re-encode exactly like a missing omitempty would for a scalar.
type RequiredOverride struct {
	Reason string
}

// RequiredOverrides is keyed exactly like PointerOverrides:
// "<command>.<element or attribute name>", a command's own top-level
// request/response fields only.
var RequiredOverrides = map[string]RequiredOverride{
	"resume_task.report_id": {
		Reason: "T-9 task lifecycle review: command_resume_task.xml's response <pattern> marks " +
			"report_id required, but internal/gmptest/testdata/live/resume_task/001_response.xml " +
			"(a genuine 400 \"Task must be in Stopped or Interrupted state\") omits it entirely — " +
			"status-dependent presence, not unconditional requiredness",
	},
	"start_task.report_id": {
		Reason: "T-9 task lifecycle review: command_start_task.xml's response <pattern> marks " +
			"report_id required, but internal/gmptest/testdata/live/start_task/001_response.xml " +
			"(a genuine 500 \"Internal error\") omits it entirely — status-dependent presence, not " +
			"unconditional requiredness",
	},
	"get_nvts.nvt_oid": {
		Reason: "T-9 nvt review: command_get_nvts.xml's own <summary>Get one or many NVTs</summary> and " +
			"request <pattern> disagree with each other — nvt_oid is marked required, but " +
			"internal/gmptest/testdata/upstream/get_nvts/001_request.xml omits it entirely and gets back " +
			"three NVTs, the many-NVTs case the summary promises",
	},

	"get_filters.filter_count": {
		Reason: "T-23 (review-flagged expansion): command_get_filters.xml's response <pattern> declares " +
			"filter_count required (bare <e>filter_count</e>), but the command's own single-filter-lookup " +
			"<example> (\"Get single filter, including alerts that use the filter\", filter_id lookup) " +
			"omits <filter_count>, <sort> and <filters> entirely — confirmed live by " +
			"internal/gmptest/testdata/upstream/get_filters/002_response.xml, the identical fixture. A " +
			"plain (non-pointer) FilterCount field cannot represent that absence: its own zero value " +
			"(Value/Filtered/Page all 0) is indistinguishable from \"the wrapper genuinely wasn't sent,\" so " +
			"re-encoding a filter-lookup response fabricated <filter_count>0<filtered>0</filtered>" +
			"<page>0</page></filter_count> where the wire sent nothing — the same class of defect T-23's " +
			"chardata-total fix does not by itself solve. Scoped to count wrappers only: the identically " +
			"required, identically sometimes-absent sibling <sort>/<filters> wrappers on this and every " +
			"other command below share this defect too, but are deliberately left as-is for now (see " +
			"docs/review-2026-08.md's T-23 entry) — only *_count is in scope here.",
	},
	// The remaining 22 commands below share get_filters' exact defect: same
	// pattern-declares-required shape, same single-item-lookup examples/fixtures
	// that omit the wrapper outright, confirmed independently by each command's
	// own schema and fidelity/live evidence. get_reports and get_results are
	// absent from this table on purpose — their schema <pattern> already wraps
	// report_count/result_count in an <o>, so fieldOf's existing CardOptional
	// auto-pointer logic already makes them *ReportCount/*ResultCount.
	"get_alerts.alert_count":                     {Reason: t23CountWrapperReason},
	"get_assets.asset_count":                     {Reason: t23CountWrapperReason},
	"get_configs.config_count":                   {Reason: t23CountWrapperReason},
	"get_credentials.credential_count":           {Reason: t23CountWrapperReason},
	"get_groups.group_count":                     {Reason: t23CountWrapperReason},
	"get_info.info_count":                        {Reason: t23CountWrapperReason},
	"get_notes.note_count":                       {Reason: t23CountWrapperReason},
	"get_overrides.override_count":               {Reason: t23CountWrapperReason},
	"get_permissions.permission_count":           {Reason: t23CountWrapperReason},
	"get_port_lists.port_list_count":             {Reason: t23CountWrapperReason},
	"get_report_configs.report_config_count":     {Reason: t23CountWrapperReason},
	"get_report_formats.report_format_count":     {Reason: t23CountWrapperReason},
	"get_roles.role_count":                       {Reason: t23CountWrapperReason},
	"get_scanners.scanner_count":                 {Reason: t23CountWrapperReason},
	"get_schedules.schedule_count":               {Reason: t23CountWrapperReason},
	"get_tags.tag_count":                         {Reason: t23CountWrapperReason},
	"get_targets.target_count":                   {Reason: t23CountWrapperReason},
	"get_tasks.task_count":                       {Reason: t23CountWrapperReason},
	"get_tickets.ticket_count":                   {Reason: t23CountWrapperReason},
	"get_tls_certificates.tls_certificate_count": {Reason: t23CountWrapperReason},
	"get_users.user_count":                       {Reason: t23CountWrapperReason},
	"get_vulns.vuln_count":                       {Reason: t23CountWrapperReason},
}

// t23CountWrapperReason is the shared Reason for every T-23 count-wrapper
// RequiredOverrides entry above get_filters.filter_count: same schema-required
// wire shape, same single-item-lookup absence, same pointer-for-presence fix.
const t23CountWrapperReason = "T-23 (review-flagged expansion): same required-but-sometimes-absent " +
	"count-wrapper defect as get_filters.filter_count above, confirmed independently by this command's own " +
	"schema <pattern> and fidelity/live evidence"

// CardinalityOverride keeps a schema-singular field Go-slice-typed despite
// that, because live protocol evidence contradicts the schema's declared
// cardinality for that specific field. The mirror image of RequiredOverride:
// that one distrusts the schema's requiredness, this one distrusts its Card.
type CardinalityOverride struct {
	Reason string
}

// CardinalityOverrides is keyed exactly like RequiredOverrides:
// "<command>.<element or attribute name>", a command's own top-level
// request/response fields only.
var CardinalityOverrides = map[string]CardinalityOverride{
	"get_nvts.nvt": {
		Reason: "T-9 nvt review: command_get_nvts.xml's response <pattern> declares <e>nvt</e> (exactly " +
			"one), contradicting its own <summary>Get one or many NVTs</summary> and its own <example> " +
			"(two sibling <nvt> elements in one get_nvts_response with no nvt_oid filter, lines " +
			"4121-4159) — internal/gmptest/testdata/upstream/get_nvts/001_response.xml confirms it live, " +
			"returning three sibling <nvt> elements for one unfiltered request",
	},
	"describe_auth.group": {
		Reason: "T-9 auth review: command_describe_auth.xml's response <pattern> declares <e>group</e> " +
			"(exactly one), but internal/gmptest/testdata/live/describe_auth/001_response.xml confirms " +
			"gvmd returns three sibling <group> elements (method:file, method:ldap_connect, " +
			"method:radius_connect) in one unfiltered describe_auth_response",
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
	"modify_user": {
		Text: "An explicit-empty <sources><source></source></sources> reproducibly segfaults gvmd " +
			"(crash inside glib's g_markup_parse_context_parse). The typed API already cannot " +
			"construct this payload: UserSources.Source is UserSource (a T-19 named string type) " +
			"with xml:\"source,omitempty\", so an empty UserSource(\"\") is dropped by omitempty, " +
			"not encoded as <source></source>. Omission leaves the stored source unchanged.",
		Reason: "T-25: live gvmd probe confirmed the segfault; see docs/review-2026-08.md's T-25 section",
	},
	"modify_note": {
		Text: "Hosts and Port are both cleared whenever either is omitted from a request, not left " +
			"unchanged: modifying one alone still blanks the other, so both must be resent together " +
			"to preserve them. Explicit-empty produces the identical cleared result.",
		Reason: "T-25: live gvmd probe confirmed omitted fields are cleared rather than left unchanged; " +
			"see docs/review-2026-08.md's T-25 section",
	},
	"modify_override": {
		Text: "Hosts, Port and Threat are all cleared whenever any is omitted from a request, not " +
			"left unchanged: modifying one alone blanks the others. NewThreat behaves oppositely: " +
			"omission and explicit-empty both leave it unchanged. This asymmetry within one command " +
			"means a caller cannot assume uniform full-replace across modify_override's fields.",
		Reason: "T-25: live gvmd probe confirmed the Hosts/Port/Threat full-replace behavior and the " +
			"NewThreat leave-unchanged behavior; see docs/review-2026-08.md's T-25 section",
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
	{Command: "get_results", Class: "duplicate-ele", Path: "/get_results_response/result/delta/notes/note/active"}: {
		Reason: "T-9 result review; same upstream `<ele>` redefined in one scope as get_notes/get_overrides, docs/review-2026-08.md T-15 defect table",
	},
	{Command: "get_results", Class: "duplicate-ele", Path: "/get_results_response/result/delta/overrides/override/active"}: {
		Reason: "T-9 result review; same upstream `<ele>` redefined in one scope as get_notes/get_overrides, docs/review-2026-08.md T-15 defect table",
	},
	{Command: "get_results", Class: "duplicate-ele", Path: "/get_results_response/result/notes/note/active"}: {
		Reason: "T-9 result review; same upstream `<ele>` redefined in one scope as get_notes/get_overrides, docs/review-2026-08.md T-15 defect table",
	},
	{Command: "get_results", Class: "duplicate-ele", Path: "/get_results_response/result/overrides/override/active"}: {
		Reason: "T-9 result review; same upstream `<ele>` redefined in one scope as get_notes/get_overrides, docs/review-2026-08.md T-15 defect table",
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
	{Command: "sync_config", Class: "empty-pattern", Path: "/sync_config"}: {
		Reason: "T-9; docs/review-2026-08.md `config` — sync_config's request is a bare <pattern/> " +
			"(no fields at all), and this gvmd build rejects the command outright " +
			"(400 \"Bogus command name\"), likely reflecting upstream deprecation in favor of " +
			"openvasd-based scanners; the schema still declares the command, so it is adopted as-is",
	},
	{Command: "test_alert", Class: "dangling-e", Path: "/test_alert_response/status_details"}: {
		Reason: "T-9; docs/review-2026-08.md `alert` — the response schema references " +
			"<e>status_details</e> but no <ele>/<elem> resolves it in scope, an upstream schema " +
			"scoping quirk rather than an adoption risk",
	},
	{Command: "create_report", Class: "empty-pattern", Path: "/create_report/report/name"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — the shared report type's <name> is " +
			"documented \"Creation time as compatibility placeholder\" with a deliberately empty " +
			"<pattern/>, not an omission",
	},
	{Command: "get_reports", Class: "empty-pattern", Path: "/get_reports_response/report/name"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type as create_report's " +
			"<name> compatibility placeholder",
	},
	{Command: "create_report", Class: "no-pattern", Path: "/create_report/report/report/apps/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — the shared report type's per-category " +
			"<count> child declares neither <pattern> nor <type>, one of several identically-shaped " +
			"count elements in this type",
	},
	{Command: "create_report", Class: "no-pattern", Path: "/create_report/report/report/closed_cves/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type <count> defect as apps/count",
	},
	{Command: "create_report", Class: "no-pattern", Path: "/create_report/report/report/cves/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type <count> defect as apps/count",
	},
	{Command: "create_report", Class: "no-pattern", Path: "/create_report/report/report/errors/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type <count> defect as apps/count",
	},
	{Command: "create_report", Class: "no-pattern", Path: "/create_report/report/report/hosts/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type <count> defect as apps/count",
	},
	{Command: "create_report", Class: "no-pattern", Path: "/create_report/report/report/os/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type <count> defect as apps/count",
	},
	{Command: "create_report", Class: "no-pattern", Path: "/create_report/report/report/ssl_certs/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type <count> defect as apps/count",
	},
	{Command: "create_report", Class: "no-pattern", Path: "/create_report/report/report/vulns/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type <count> defect as apps/count",
	},
	{Command: "get_reports", Class: "no-pattern", Path: "/get_reports_response/report/report/apps/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type <count> defect as create_report",
	},
	{Command: "get_reports", Class: "no-pattern", Path: "/get_reports_response/report/report/closed_cves/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type <count> defect as create_report",
	},
	{Command: "get_reports", Class: "no-pattern", Path: "/get_reports_response/report/report/cves/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type <count> defect as create_report",
	},
	{Command: "get_reports", Class: "no-pattern", Path: "/get_reports_response/report/report/errors/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type <count> defect as create_report",
	},
	{Command: "get_reports", Class: "no-pattern", Path: "/get_reports_response/report/report/hosts/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type <count> defect as create_report",
	},
	{Command: "get_reports", Class: "no-pattern", Path: "/get_reports_response/report/report/os/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type <count> defect as create_report",
	},
	{Command: "get_reports", Class: "no-pattern", Path: "/get_reports_response/report/report/ssl_certs/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type <count> defect as create_report",
	},
	{Command: "get_reports", Class: "no-pattern", Path: "/get_reports_response/report/report/vulns/count"}: {
		Reason: "T-9; docs/review-2026-08.md `report` — same shared report type <count> defect as create_report",
	},
	{Command: "create_report", Class: "duplicate-ele", Path: "/create_report/report/report/results/result/delta/notes/note/active"}: {
		Reason: "T-9; upstream `<ele>` redefined in one scope, docs/review-2026-08.md `report` — " +
			"same shared result type as get_notes'/get_overrides' delta-path `active` defect",
	},
	{Command: "create_report", Class: "duplicate-ele", Path: "/create_report/report/report/results/result/delta/overrides/override/active"}: {
		Reason: "T-9; upstream `<ele>` redefined in one scope, docs/review-2026-08.md `report` — " +
			"same shared result type as get_notes'/get_overrides' delta-path `active` defect",
	},
	{Command: "create_report", Class: "duplicate-ele", Path: "/create_report/report/report/results/result/notes/note/active"}: {
		Reason: "T-9; upstream `<ele>` redefined in one scope, docs/review-2026-08.md `report` — " +
			"same shared result type as get_notes'/get_overrides' `active` defect",
	},
	{Command: "create_report", Class: "duplicate-ele", Path: "/create_report/report/report/results/result/overrides/override/active"}: {
		Reason: "T-9; upstream `<ele>` redefined in one scope, docs/review-2026-08.md `report` — " +
			"same shared result type as get_notes'/get_overrides' `active` defect",
	},
	{Command: "get_reports", Class: "duplicate-ele", Path: "/get_reports_response/report/report/results/result/delta/notes/note/active"}: {
		Reason: "T-9; upstream `<ele>` redefined in one scope, docs/review-2026-08.md `report` — " +
			"same shared result type as create_report's delta-path `active` defect",
	},
	{Command: "get_reports", Class: "duplicate-ele", Path: "/get_reports_response/report/report/results/result/delta/overrides/override/active"}: {
		Reason: "T-9; upstream `<ele>` redefined in one scope, docs/review-2026-08.md `report` — " +
			"same shared result type as create_report's delta-path `active` defect",
	},
	{Command: "get_reports", Class: "duplicate-ele", Path: "/get_reports_response/report/report/results/result/notes/note/active"}: {
		Reason: "T-9; upstream `<ele>` redefined in one scope, docs/review-2026-08.md `report` — " +
			"same shared result type as create_report's `active` defect",
	},
	{Command: "get_reports", Class: "duplicate-ele", Path: "/get_reports_response/report/report/results/result/overrides/override/active"}: {
		Reason: "T-9; upstream `<ele>` redefined in one scope, docs/review-2026-08.md `report` — " +
			"same shared result type as create_report's `active` defect",
	},
	{Command: "empty_trashcan", Class: "empty-pattern", Path: "/empty_trashcan"}: {
		Reason: "T-9; docs/review-2026-08.md `trash` — the request pattern declares no content, " +
			"same generic upstream-schema class already carried by note/override/create_group",
	},
	{Command: "get_info", Class: "no-pattern", Path: "/get_info_response/info/web_application_vt/refs/ref"}: {
		Reason: "T-9; docs/review-2026-08.md `nvt` — <ref> under web_application_vt/refs declares " +
			"neither <pattern> nor <type>, same generic upstream-schema class as get_scanners' params/param",
	},
	{Command: "get_info", Class: "undefined-type", Path: "/get_info_response/info/cve/configuration_nodes/node/node"}: {
		Reason: "T-9; docs/review-2026-08.md `nvt` — <node> nests itself recursively, a self-reference " +
			"the schema cannot express as a named type; FieldTypeOverrides retypes it back to its own " +
			"generated struct, which Go allows through a slice field",
	},
	{Command: "describe_auth", Class: "empty-pattern", Path: "/describe_auth"}: {
		Reason: "T-9; docs/review-2026-08.md `auth` — the request pattern declares no content, same " +
			"generic upstream-schema class already carried by note/override/create_group/empty_trashcan",
	},
	{Command: "get_license", Class: "empty-pattern", Path: "/get_license"}: {
		Reason: "T-9; docs/review-2026-08.md `license` — the request pattern declares no content, same " +
			"generic upstream-schema class already carried by describe_auth/empty_trashcan",
	},
	{Command: "get_version", Class: "empty-pattern", Path: "/get_version"}: {
		Reason: "T-9; docs/review-2026-08.md `version` — the request pattern declares no content, same " +
			"generic upstream-schema class already carried by describe_auth/get_license/empty_trashcan",
	},
	{Command: "get_feeds", Class: "no-pattern", Path: "/get_feeds_response/feed/sync_not_available/error"}: {
		Reason: "T-9; docs/review-2026-08.md `feed` — <error>'s <ele> declares neither <pattern> nor " +
			"<type>, same class already carried by get_scanners/create_report/get_reports/get_info",
	},
	{Command: "get_features", Class: "empty-pattern", Path: "/get_features"}: {
		Reason: "T-9; docs/review-2026-08.md `feature` — the request pattern declares no content, same " +
			"generic upstream-schema class already carried by describe_auth/get_license/get_version",
	},
	{Command: "get_resource_names", Class: "undefined-type", Path: "/get_resource_names_response/type"}: {
		Reason: "T-9; docs/review-2026-08.md `settings` — <type>'s own declared <type> is itself " +
			"`type`, a self-reference with no defined alias, same class as get_assets' @distance and " +
			"get_info's recursive <node>",
	},
}

// EnumConstantName names the constant suffix for one <alts>-declared wire
// value that goName (generate.go) cannot turn into a Go identifier fragment
// on its own — a value that is nothing but punctuation has no letters for
// goName to capitalize and no separator to split on, so it comes back
// unchanged and invalid as an exported constant name.
type EnumConstantName struct {
	Suffix string
	Reason string
}

// EnumConstantNames is keyed by the literal wire value, not by command: a GMP
// <relation> filter operator means the same thing in every command it
// appears in (T-19; docs/review-2026-08.md), so one entry stands in for every
// occurrence across the corpus instead of a family-scoped table repeating the
// same fix. Checked for staleness once across the whole GenerateAll run
// rather than per family, since a family that never happens to use a given
// operator is not evidence the entry is unused corpus-wide.
var EnumConstantNames = map[string]EnumConstantName{
	"=": {
		Suffix: "Eq",
		Reason: "T-19: <relation> equality operator; punctuation-only, no word for goName to derive a name from",
	},
	":": {
		Suffix: "Contains",
		Reason: "T-19: <relation> substring-match operator; goName treats ':' itself as a word separator, " +
			"which would otherwise silently produce the misleading generic name \"Field\"",
	},
	"~": {
		Suffix: "Regex",
		Reason: "T-19: <relation> regular-expression-match operator; punctuation-only",
	},
	">": {
		Suffix: "GT",
		Reason: "T-19: <relation> greater-than operator; punctuation-only",
	},
	"<": {
		Suffix: "LT",
		Reason: "T-19: <relation> less-than operator; punctuation-only",
	},
}
