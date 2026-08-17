# GMP schema conformance report

Regenerate with: go test ./internal/schema -run TestConformanceReport -update

What this is: the GMP schema DSL, read directly, diffed against the
hand-written structs in commands/. The schema revision compared here is the
one vendored at vendor/XML/commands. Before T-16b this matched commands/
one-for-one, so a difference found here was always a transcription error;
since T-16b, commands/ is generated from a live gvmd server instead, so a
difference here can also be version drift between the vendored copy and the
live schema — see docs/review-2026-08.md's T-16b/T-17 notes for the known
cases.

This is a *schema* oracle. internal/gmptest is a *runtime* oracle: it
round-trips captured documents and records what they lose. Findings marked
[fidelity] were reached by both, from different evidence, by tools that share
no code — those are the highest-confidence class here.

Categories:
  MISSING       the schema declares it, no Go field carries it
  EXTRA         a Go field with nothing behind it in the schema
  CARDINALITY   repeated-vs-slice, or an optional whose absence Go cannot express
  KIND          an attribute modelled as an element, or the reverse
  OMITEMPTY     omitempty where the schema says required, or absent where optional
  TYPE          the Go type cannot hold the schema's declared lexical space
  SCHEMA-DEFECT the schema is wrong; nothing is asked of commands/

Findings marked [suppressed] are counted but not charged to commands/. They are
almost all request-side over-constraint: the schema marks create_report's
report/@id, config_id, format_id, extension and content_type all required,
though real payloads only ever populate one submission path, so response-side
optionality is trusted and request-side is not.

## Totals

commands compared:            122
fully conformant:             93
with actionable findings:     29
actionable findings:          256  (request 117, response 139)
suppressed findings:          439
confirmed by fidelity corpus: 53
in the schema, not in Go:     0
in Go, not in the schema:     0

### Findings by category

                 actionable suppressed
  MISSING               131          0
  EXTRA                  80          6
  CARDINALITY            22          0
  KIND                    0          0
  OMITEMPTY               0        433
  TYPE                   23          0
  SCHEMA-DEFECT         128          0

### Schema defects by class

How much of the schema is unusable as written. None of these are
transcription errors; they are upstream bugs the parser has to survive.

These counts are lower than the ones in docs/review-2026-08.md T-9, which were
taken against upstream main (155 commands). Two reasons, both expected:

  - This revision predates the commands that carry most of the damage. The 68
    dangling <e> references upstream are concentrated in get_scan_report,
    get_agent_groups, get_report_tls_certificates and get_integration_configs,
    none of which exist here.
  - Only the shape is walked. A <type> naming UUID inside a <filter_keywords>
    <column> describes filter-term syntax, not an element or an attribute, so
    it is out of scope: 28 of the 31 UUID uses are of that kind.

  dangling-e           1 sites in 1 commands
  duplicate-ele       18 sites in 5 commands
  duplicate-ref        7 sites in 2 commands
  empty-pattern       44 sites in 23 commands
  no-pattern          17 sites in 4 commands
  recursion-cut       12 sites in 1 commands
  type-for-t          23 sites in 1 commands
  undefined-type       3 sites in 3 commands
  unsupported-r        3 sites in 3 commands

### Commands by actionable findings

  get_reports                  70   MISSING=55  EXTRA=6  CARDINALITY=5  TYPE=4
  create_report                52   MISSING=38  EXTRA=5  CARDINALITY=5  TYPE=4
  get_info                     22   MISSING=12  EXTRA=7  TYPE=3
  delete_agents                18   EXTRA=18
  get_nvts                     18   MISSING=6  EXTRA=1  CARDINALITY=7  TYPE=4
  modify_agents                18   EXTRA=18
  get_aggregates               12   MISSING=5  EXTRA=5  CARDINALITY=2
  get_agents                    9   MISSING=7  TYPE=2
  get_results                   7   MISSING=1  EXTRA=4  CARDINALITY=2
  get_vulns                     4   MISSING=1  EXTRA=3
  get_feeds                     3   TYPE=3
  describe_auth                 2   MISSING=1  CARDINALITY=1
  get_credentials               2   EXTRA=2
  get_scanners                  2   EXTRA=2
  get_settings                  2   EXTRA=2
  modify_target                 2   MISSING=1  EXTRA=1
  create_target                 1   EXTRA=1
  create_task                   1   MISSING=1
  get_features                  1   TYPE=1
  get_license                   1   TYPE=1
  ... and 9 more with fewer

### Confirmed by both oracles

Paths the schema says exist and that a captured document was also seen to
lose. Fix these first: neither oracle is guessing.

Response side only: the fidelity corpus keeps per-path DROPPED lists for
responses and only aggregate counts for requests, so a request-side finding can
never carry this mark. That is a gap in the cross-reference, not evidence that
requests are in better shape — get_notes models none of its seven request
attributes, and one of its own captured request fixtures uses three of them.

  get_agents               /get_agents_response/agent_count
  get_agents               /get_agents_response/agents
  get_agents               /get_agents_response/filters
  get_agents               /get_agents_response/sort
  get_aggregates           /get_aggregates_response/filters
  get_info                 /get_info_response/info/cve/severity
  get_nvts                 /get_nvts_response/nvt/refs/*
  get_report_configs       /get_report_configs_response/report_config/param/value/@using_default
  get_reports              /get_reports_response/filters
  get_reports              /get_reports_response/report/#text
  get_reports              /get_reports_response/report/report/compliance_count/#text
  get_reports              /get_reports_response/report/report/compliance_count/incomplete/*
  get_reports              /get_reports_response/report/report/compliance_count/no/*
  get_reports              /get_reports_response/report/report/compliance_count/undefined/*
  get_reports              /get_reports_response/report/report/compliance_count/yes/*
  get_reports              /get_reports_response/report/report/delta/report/*
  get_reports              /get_reports_response/report/report/errors/error/host/*
  get_reports              /get_reports_response/report/report/errors/error/nvt/*
  get_reports              /get_reports_response/report/report/filters/#text
  get_reports              /get_reports_response/report/report/filters/delta/#text
  get_reports              /get_reports_response/report/report/host/asset/*
  get_reports              /get_reports_response/report/report/host/compliance_count/*
  get_reports              /get_reports_response/report/report/host/detail/source/*
  get_reports              /get_reports_response/report/report/host/port_count/*
  get_reports              /get_reports_response/report/report/host/result_count/*
  get_reports              /get_reports_response/report/report/ports/port/#text
  get_reports              /get_reports_response/report/report/result_count/critical
  get_reports              /get_reports_response/report/report/result_count/high
  get_reports              /get_reports_response/report/report/result_count/hole
  get_reports              /get_reports_response/report/report/result_count/info
  get_reports              /get_reports_response/report/report/result_count/log
  get_reports              /get_reports_response/report/report/result_count/low
  get_reports              /get_reports_response/report/report/result_count/medium
  get_reports              /get_reports_response/report/report/result_count/warning
  get_reports              /get_reports_response/report/report/results/result/delta/*
  get_reports              /get_reports_response/report/report/results/result/detection/*
  get_reports              /get_reports_response/report/report/results/result/host/asset
  get_reports              /get_reports_response/report/report/results/result/notes/*
  get_reports              /get_reports_response/report/report/results/result/nvt/epss/*
  get_reports              /get_reports_response/report/report/results/result/nvt/refs/*
  get_reports              /get_reports_response/report/report/results/result/nvt/severities/*
  get_reports              /get_reports_response/report/report/results/result/overrides/*
  get_reports              /get_reports_response/report/report/results/result/qod/*
  get_reports              /get_reports_response/report/report/results/result/tickets/*
  get_reports              /get_reports_response/report/report/sort/field/#text
  get_reports              /get_reports_response/report/report/tls_certificates/tls_certificate/certificate/*
  get_reports              /get_reports_response/report/report/tls_certificates/tls_certificate/host/*
  get_reports              /get_reports_response/report_count
  get_reports              /get_reports_response/reports
  get_reports              /get_reports_response/sort
  get_results              /get_results_response/result/nvt/severities/severity
  get_tasks                /get_tasks_response/task/hosts_ordering
  help                     /help_response/schema

## Per command

### authenticate  (command_authenticate.xml)

OMITEMPTY=2

  OMITEMPTY     response /authenticate_response/role
                schema says required, Role is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /authenticate_response/timezone
                schema says required, Timezone is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### create_alert  (command_create_alert.xml)

conformant

### create_asset  (command_create_asset.xml)

conformant

### create_config  (command_create_config.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /create_config
                unsupported-r: <r> references a response from inside a pattern

### create_credential  (command_create_credential.xml)

conformant

### create_filter  (command_create_filter.xml)

conformant

### create_group  (command_create_group.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /create_group/specials/full
                empty-pattern: pattern declares no content

### create_note  (command_create_note.xml)

conformant

### create_override  (command_create_override.xml)

conformant

### create_permission  (command_create_permission.xml)

conformant

### create_port_list  (command_create_port_list.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /create_port_list
                unsupported-r: <r> references a response from inside a pattern

### create_port_range  (command_create_port_range.xml)

conformant

### create_report  (command_create_report.xml)

MISSING=38  EXTRA=5  CARDINALITY=5  OMITEMPTY=144  TYPE=4  SCHEMA-DEFECT=12

  MISSING       request  /create_report/report/#text
                schema also declares character data (text) here; Go *commands.Report has no `,chardata` field
  MISSING       request  /create_report/report/report/compliance_count/#text
                schema also declares character data (text) here; Go *commands.ReportComplianceCount has no `,chardata` field
  MISSING       request  /create_report/report/report/compliance_count/incomplete/*
                schema declares 2 children; Go models this as scalar int
  MISSING       request  /create_report/report/report/compliance_count/no/*
                schema declares 2 children; Go models this as scalar int
  MISSING       request  /create_report/report/report/compliance_count/undefined/*
                schema declares 2 children; Go models this as scalar int
  MISSING       request  /create_report/report/report/compliance_count/yes/*
                schema declares 2 children; Go models this as scalar int
  MISSING       request  /create_report/report/report/delta/report/*
                schema declares 5 children; Go models this as scalar *commands.ReportContent
  MISSING       request  /create_report/report/report/errors/error/host/*
                schema declares 1 children; Go models this as scalar string
  MISSING       request  /create_report/report/report/errors/error/nvt/*
                schema declares 4 children; Go models this as scalar string
  MISSING       request  /create_report/report/report/filters/#text
                schema also declares character data (levels) here; Go *commands.ReportFilters has no `,chardata` field
  MISSING       request  /create_report/report/report/filters/delta/#text
                schema also declares character data (text) here; Go *commands.ReportDeltaStates has no `,chardata` field
  MISSING       request  /create_report/report/report/host/asset/*
                schema declares 1 children; Go models this as scalar string
  MISSING       request  /create_report/report/report/host/compliance_count/*
                schema declares 5 children; Go models this as scalar int
  MISSING       request  /create_report/report/report/host/detail/source/*
                schema declares 3 children; Go models this as scalar string
  MISSING       request  /create_report/report/report/host/port_count/*
                schema declares 1 children; Go models this as scalar int
  MISSING       request  /create_report/report/report/host/result_count/*
                schema declares 10 children; Go models this as scalar int
  MISSING       request  /create_report/report/report/ports/port/#text
                schema also declares character data (text) here; Go []commands.ReportPort has no `,chardata` field
  MISSING       request  /create_report/report/report/result_count/critical
                required elem critical is not modelled
  MISSING       request  /create_report/report/report/result_count/high
                required elem high is not modelled
  MISSING       request  /create_report/report/report/result_count/hole
                required elem hole is not modelled
  MISSING       request  /create_report/report/report/result_count/info
                required elem info is not modelled
  MISSING       request  /create_report/report/report/result_count/log
                required elem log is not modelled
  MISSING       request  /create_report/report/report/result_count/low
                required elem low is not modelled
  MISSING       request  /create_report/report/report/result_count/medium
                required elem medium is not modelled
  MISSING       request  /create_report/report/report/result_count/warning
                required elem warning is not modelled
  MISSING       request  /create_report/report/report/results/result/delta/*
                schema declares 4 children; Go models this as scalar string
  MISSING       request  /create_report/report/report/results/result/detection/*
                schema declares 1 children; Go models this as scalar string
  MISSING       request  /create_report/report/report/results/result/host/asset
                required elem asset is not modelled
  MISSING       request  /create_report/report/report/results/result/notes/*
                schema declares 1 children; Go models this as scalar string
  MISSING       request  /create_report/report/report/results/result/nvt/epss/*
                schema declares 2 children; Go models this as scalar string
  MISSING       request  /create_report/report/report/results/result/nvt/refs/*
                schema declares 1 children; Go models this as scalar string
  MISSING       request  /create_report/report/report/results/result/nvt/severities/*
                schema declares 2 children; Go models this as scalar string
  MISSING       request  /create_report/report/report/results/result/overrides/*
                schema declares 1 children; Go models this as scalar string
  MISSING       request  /create_report/report/report/results/result/qod/*
                schema declares 2 children; Go models this as scalar string
  MISSING       request  /create_report/report/report/results/result/tickets/*
                schema declares 1 children; Go models this as scalar string
  MISSING       request  /create_report/report/report/sort/field/#text
                schema also declares character data (text) here; Go []commands.ReportSortField has no `,chardata` field
  MISSING       request  /create_report/report/report/tls_certificates/tls_certificate/certificate/*
                schema declares 1 children; Go models this as scalar string
  MISSING       request  /create_report/report/report/tls_certificates/tls_certificate/host/*
                schema declares 2 children; Go models this as scalar string
  EXTRA         request  /create_report/report/base64
                Go field Base64 string has no counterpart in the schema
  EXTRA         request  /create_report/report/report/filters/levels
                Go field Levels string has no counterpart in the schema
  EXTRA         request  /create_report/report/report/results/result/host/@asset_id
                Go field AssetID string has no counterpart in the schema
  EXTRA         request  /create_report/report/report/tls_certificates/tls_certificate/@format
                Go field Format string has no counterpart in the schema
  EXTRA         request  /create_report/task/name
                Go field Name string has no counterpart in the schema
  CARDINALITY   request  /create_report/report/report/host/compliance_count
                schema says optional; ComplianceCount int cannot express absence (zero is a legal value)
  CARDINALITY   request  /create_report/report/report/host/result_count
                schema says optional; ResultCount int cannot express absence (zero is a legal value)
  CARDINALITY   request  /create_report/report/report/report_format/param
                schema says required, Go has a slice []commands.ReportContentFormatParam (Params)
  CARDINALITY   request  /create_report/report/report/results/result/original_severity
                schema says optional; OriginalSeverity float64 cannot express absence (zero is a legal value)
  CARDINALITY   request  /create_report/report/report/sort/field
                schema says required, Go has a slice []commands.ReportSortField (Fields)
  OMITEMPTY     request  /create_report/report/@config_id
                schema says required, ConfigID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/@content_type
                schema says required, ContentType is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/@extension
                schema says required, Extension is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/@format_id
                schema says required, FormatID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/comment
                schema says required, Comment is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/creation_time
                schema says required, CreationTime is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/in_use
                schema says required, InUse is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/modification_time
                schema says required, ModificationTime is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/owner
                schema says required, Owner is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/owner/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report
                schema says required, Report is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/compliance/filtered
                schema says required, Filtered is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/compliance/full
                schema says required, Full is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/compliance_count/filtered
                schema says required, Filtered is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/compliance_count/full
                schema says required, Full is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/compliance_count/incomplete
                schema says required, Incomplete is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/compliance_count/no
                schema says required, No is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/compliance_count/undefined
                schema says required, Undefined is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/compliance_count/yes
                schema says required, Yes is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/delta/report
                schema says required, Report is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/errors/error/description
                schema says required, Description is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/errors/error/host
                schema says required, Host is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/errors/error/nvt
                schema says required, NVT is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/errors/error/port
                schema says required, Port is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/errors/error/scan_nvt_version
                schema says required, ScanNVTVersion is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/errors/error/severity
                schema says required, Severity is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/filters
                schema says required, Filters is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/filters/apply_overrides
                schema says required, ApplyOverrides is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/filters/delta/changed
                schema says required, Changed is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/filters/delta/gone
                schema says required, Gone is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/filters/delta/new
                schema says required, New is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/filters/delta/same
                schema says required, Same is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/filters/min_qod
                schema says required, MinQod is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/filters/notes
                schema says required, Notes is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/filters/overrides
                schema says required, Overrides is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/filters/phrase
                schema says required, Phrase is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/filters/result_hosts_only
                schema says required, ResultHostsOnly is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/filters/term
                schema says required, Term is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/gmp
                schema says required, GMP is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/gmp/version
                schema says required, Version is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/host/asset
                schema says required, Asset is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/host/detail/extra
                schema says required, Extra is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/host/detail/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/host/detail/source
                schema says required, Source is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/host/detail/value
                schema says required, Value is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/host/end
                schema says required, End is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/host/ip
                schema says required, IP is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/host/port_count
                schema says required, PortCount is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/host/start
                schema says required, Start is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/hosts/@max
                schema says required, Max is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]
  OMITEMPTY     request  /create_report/report/report/hosts/@start
                schema says required, Start is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]
  OMITEMPTY     request  /create_report/report/report/permissions/permission/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/ports/@max
                schema says required, Max is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/ports/@start
                schema says required, Start is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/ports/port/host
                schema says required, Host is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/ports/port/severity
                schema says required, Severity is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/ports/port/threat
                schema says required, Threat is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/report_format
                schema says required, ReportFormat is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/report_format/param
                schema says required, Params is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/report_format/param/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/report_format/param/value
                schema says required, Value is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/result_count/filtered
                schema says required, Filtered is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/result_count/full
                schema says required, Full is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/@max
                schema says required, Max is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/@start
                schema says required, Start is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/comment
                schema says required, Comment is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/compliance
                schema says required, Compliance is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/creation_time
                schema says required, CreationTime is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/description
                schema says required, Description is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/host
                schema says required, Host is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/host/hostname
                schema says required, Hostname is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/modification_time
                schema says required, ModificationTime is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/nvt
                schema says required, NVT is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/nvt/@oid
                schema says required, OID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/nvt/cpe
                schema says required, CPE is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/nvt/cvss_base
                schema says required, CVSSBase is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/nvt/family
                schema says required, Family is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/nvt/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/nvt/refs
                schema says required, Refs is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/nvt/severities
                schema says required, Severities is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/nvt/tags
                schema says required, Tags is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/nvt/type
                schema says required, Type is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/owner
                schema says required, Owner is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/owner/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/port
                schema says required, Port is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/qod
                schema says required, QOD is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/report/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/scan_nvt_version
                schema says required, ScanNVTVersion is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/severity
                schema says required, Severity is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/task/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/task/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/threat
                schema says required, Threat is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/user_tags/count
                schema says required, Count is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/user_tags/tag/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/user_tags/tag/comment
                schema says required, Comment is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/user_tags/tag/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/results/result/user_tags/tag/value
                schema says required, Value is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/severity/filtered
                schema says required, Filtered is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/severity/full
                schema says required, Full is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/sort
                schema says required, Sort is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/sort/field
                schema says required, Fields is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/sort/field/order
                schema says required, Order is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/task/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/task/comment
                schema says required, Comment is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/task/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/task/progress
                schema says required, Progress is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/task/target
                schema says required, Target is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/task/target/comment
                schema says required, Comment is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/task/target/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/task/target/trash
                schema says required, Trash is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/task/user_tags/count
                schema says required, Count is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/task/user_tags/tag/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/task/user_tags/tag/comment
                schema says required, Comment is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/task/user_tags/tag/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/task/user_tags/tag/value
                schema says required, Value is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/tls_certificates/tls_certificate/activation_time
                schema says required, ActivationTime is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/tls_certificates/tls_certificate/certificate
                schema says required, Certificate is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/tls_certificates/tls_certificate/expiration_time
                schema says required, ExpirationTime is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/tls_certificates/tls_certificate/host
                schema says required, Host is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/tls_certificates/tls_certificate/issuer_dn
                schema says required, IssuerDN is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/tls_certificates/tls_certificate/md5_fingerprint
                schema says required, MD5Fingerprint is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/tls_certificates/tls_certificate/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/tls_certificates/tls_certificate/serial
                schema says required, Serial is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/tls_certificates/tls_certificate/sha256_fingerprint
                schema says required, SHA256Fingerprint is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/tls_certificates/tls_certificate/subject_dn
                schema says required, SubjectDN is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/tls_certificates/tls_certificate/valid
                schema says required, Valid is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/user_tags/count
                schema says required, Count is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/user_tags/tag/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/user_tags/tag/comment
                schema says required, Comment is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/user_tags/tag/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report/user_tags/tag/value
                schema says required, Value is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report_config/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report_config/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report_format/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/report_format/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/task
                schema says required, Task is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/task/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/task/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/report/writable
                schema says required, Writable is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_report/task
                schema says required, Task is tagged omitempty  [suppressed: request-side over-constraint]
  TYPE          request  /create_report/report/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          request  /create_report/report/report/task/target/trash
                schema type boolean is text ("" is true), Go Trash is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          request  /create_report/report/report/tls_certificates/tls_certificate/valid
                schema type boolean is text ("" is true), Go Valid is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          request  /create_report/report/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT request  /create_report/report/name
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT request  /create_report/report/report/apps/count
                no-pattern: <ele> has neither <pattern> nor <type>
  SCHEMA-DEFECT request  /create_report/report/report/closed_cves/count
                no-pattern: <ele> has neither <pattern> nor <type>
  SCHEMA-DEFECT request  /create_report/report/report/errors/count
                no-pattern: <ele> has neither <pattern> nor <type>
  SCHEMA-DEFECT request  /create_report/report/report/hosts/count
                no-pattern: <ele> has neither <pattern> nor <type>
  SCHEMA-DEFECT request  /create_report/report/report/os/count
                no-pattern: <ele> has neither <pattern> nor <type>
  SCHEMA-DEFECT request  /create_report/report/report/results/result/delta/notes/note/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT request  /create_report/report/report/results/result/delta/overrides/override/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT request  /create_report/report/report/results/result/notes/note/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT request  /create_report/report/report/results/result/overrides/override/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT request  /create_report/report/report/ssl_certs/count
                no-pattern: <ele> has neither <pattern> nor <type>
  SCHEMA-DEFECT request  /create_report/report/report/vulns/count
                no-pattern: <ele> has neither <pattern> nor <type>

### create_report_config  (command_create_report_config.xml)

conformant

### create_report_format  (command_create_report_format.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /create_report_format
                unsupported-r: <r> references a response from inside a pattern

### create_role  (command_create_role.xml)

conformant

### create_scanner  (command_create_scanner.xml)

conformant

### create_schedule  (command_create_schedule.xml)

conformant

### create_tag  (command_create_tag.xml)

conformant

### create_target  (command_create_target.xml)

EXTRA=1

  EXTRA         request  /create_target/alive_tests/*
                schema declares a "alive_test" leaf; Go models 1 children

### create_task  (command_create_task.xml)

MISSING=1

  MISSING       request  /create_task/hosts_ordering
                optional elem hosts_ordering is not modelled

### create_ticket  (command_create_ticket.xml)

conformant

### create_tls_certificate  (command_create_tls_certificate.xml)

conformant

### create_user  (command_create_user.xml)

conformant

### delete_agents  (command_delete_agents.xml)

EXTRA=18

  EXTRA         request  /delete_agents/agents/agent/agent_id
                Go field AgentID string has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/authorized
                Go field Authorized int has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/comment
                Go field Comment string has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/connection_status
                Go field ConnectionStatus string has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/creation_time
                Go field CreationTime commands.GMPTime has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/heartbeat_interval
                Go field HeartbeatInterval int has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/hostname
                Go field Hostname string has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/in_use
                Go field InUse bool has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/ip
                Go field IP []string has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/last_update
                Go field LastUpdate commands.GMPTime has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/min_interval
                Go field MinInterval int has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/modification_time
                Go field ModificationTime commands.GMPTime has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/owner
                Go field Owner *commands.AgentOwner has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/permissions
                Go field Permissions *commands.AgentPermissions has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/scanner
                Go field Scanner *commands.AgentScanner has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/schedule
                Go field Schedule string has no counterpart in the schema
  EXTRA         request  /delete_agents/agents/agent/writable
                Go field Writable bool has no counterpart in the schema

### delete_alert  (command_delete_alert.xml)

conformant

### delete_asset  (command_delete_asset.xml)

OMITEMPTY=2

  OMITEMPTY     request  /delete_asset/@asset_id
                schema says required, AssetID is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]
  OMITEMPTY     request  /delete_asset/@report_id
                schema says required, ReportID is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]

### delete_config  (command_delete_config.xml)

conformant

### delete_credential  (command_delete_credential.xml)

conformant

### delete_filter  (command_delete_filter.xml)

conformant

### delete_group  (command_delete_group.xml)

conformant

### delete_note  (command_delete_note.xml)

conformant

### delete_override  (command_delete_override.xml)

conformant

### delete_permission  (command_delete_permission.xml)

conformant

### delete_port_list  (command_delete_port_list.xml)

conformant

### delete_port_range  (command_delete_port_range.xml)

conformant

### delete_report  (command_delete_report.xml)

conformant

### delete_report_config  (command_delete_report_config.xml)

conformant

### delete_report_format  (command_delete_report_format.xml)

conformant

### delete_role  (command_delete_role.xml)

conformant

### delete_scanner  (command_delete_scanner.xml)

conformant

### delete_schedule  (command_delete_schedule.xml)

conformant

### delete_tag  (command_delete_tag.xml)

conformant

### delete_target  (command_delete_target.xml)

conformant

### delete_task  (command_delete_task.xml)

conformant

### delete_ticket  (command_delete_ticket.xml)

conformant

### delete_user  (command_delete_user.xml)

OMITEMPTY=4

  OMITEMPTY     request  /delete_user/@inheritor_id
                schema says required, InheritorID is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]
  OMITEMPTY     request  /delete_user/@inheritor_name
                schema says required, InheritorName is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]
  OMITEMPTY     request  /delete_user/@name
                schema says required, Name is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]
  OMITEMPTY     request  /delete_user/@user_id
                schema says required, UserID is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]

### describe_auth  (command_describe_auth.xml)

MISSING=1  CARDINALITY=1  SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /describe_auth
                empty-pattern: pattern declares no content
  MISSING       response /describe_auth_response/group/auth_conf_setting/certificate_info
                optional elem certificate_info is not modelled
  CARDINALITY   response /describe_auth_response/group
                schema says required, Go has a slice []commands.AuthGroup (Groups)

### empty_trashcan  (command_empty_trashcan.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /empty_trashcan
                empty-pattern: pattern declares no content

### get_agents  (command_get_agents.xml)

MISSING=7  OMITEMPTY=17  TYPE=2  SCHEMA-DEFECT=23

  MISSING       request  /get_agents/@agent_id
                optional attr agent_id is not modelled
  MISSING       request  /get_agents/@filt_id
                optional attr filt_id is not modelled
  MISSING       request  /get_agents/@filter
                optional attr filter is not modelled
  MISSING       response /get_agents_response/agent_count
                required elem agent_count is not modelled  [fidelity]
  MISSING       response /get_agents_response/agents
                required elem agents is not modelled  [fidelity]
  MISSING       response /get_agents_response/filters
                required elem filters is not modelled  [fidelity]
  MISSING       response /get_agents_response/sort
                required elem sort is not modelled  [fidelity]
  OMITEMPTY     response /get_agents_response/agent/agent_id
                schema says required, AgentID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_agents_response/agent/authorized
                schema says required, Authorized is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_agents_response/agent/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_agents_response/agent/connection_status
                schema says required, ConnectionStatus is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_agents_response/agent/creation_time
                schema says required, CreationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_agents_response/agent/heartbeat_interval
                schema says required, HeartbeatInterval is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_agents_response/agent/hostname
                schema says required, Hostname is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_agents_response/agent/in_use
                schema says required, InUse is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_agents_response/agent/last_update
                schema says required, LastUpdate is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_agents_response/agent/min_interval
                schema says required, MinInterval is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_agents_response/agent/modification_time
                schema says required, ModificationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_agents_response/agent/owner
                schema says required, Owner is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_agents_response/agent/permissions
                schema says required, Permissions is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_agents_response/agent/scanner
                schema says required, Scanner is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_agents_response/agent/scanner/@id
                schema says optional, ID string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_agents_response/agent/schedule
                schema says required, Schedule is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_agents_response/agent/writable
                schema says required, Writable is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_agents_response/agent/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_agents_response/agent/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_agents_response/agent/agent_id
                type-for-t: <type> used where <t> was meant: text
  SCHEMA-DEFECT response /get_agents_response/agent/authorized
                type-for-t: <type> used where <t> was meant: integer
  SCHEMA-DEFECT response /get_agents_response/agent/comment
                type-for-t: <type> used where <t> was meant: text
  SCHEMA-DEFECT response /get_agents_response/agent/connection_status
                type-for-t: <type> used where <t> was meant: text
  SCHEMA-DEFECT response /get_agents_response/agent/creation_time
                type-for-t: <type> used where <t> was meant: iso_time
  SCHEMA-DEFECT response /get_agents_response/agent/heartbeat_interval
                type-for-t: <type> used where <t> was meant: integer
  SCHEMA-DEFECT response /get_agents_response/agent/hostname
                type-for-t: <type> used where <t> was meant: text
  SCHEMA-DEFECT response /get_agents_response/agent/in_use
                type-for-t: <type> used where <t> was meant: boolean
  SCHEMA-DEFECT response /get_agents_response/agent/ip
                type-for-t: <type> used where <t> was meant: text
  SCHEMA-DEFECT response /get_agents_response/agent/last_update
                type-for-t: <type> used where <t> was meant: iso_time
  SCHEMA-DEFECT response /get_agents_response/agent/min_interval
                type-for-t: <type> used where <t> was meant: integer
  SCHEMA-DEFECT response /get_agents_response/agent/modification_time
                type-for-t: <type> used where <t> was meant: iso_time
  SCHEMA-DEFECT response /get_agents_response/agent/name
                type-for-t: <type> used where <t> was meant: name
  SCHEMA-DEFECT response /get_agents_response/agent/owner/name
                type-for-t: <type> used where <t> was meant: name
  SCHEMA-DEFECT response /get_agents_response/agent/permissions/permission/name
                type-for-t: <type> used where <t> was meant: name
  SCHEMA-DEFECT response /get_agents_response/agent/scanner/name
                type-for-t: <type> used where <t> was meant: name
  SCHEMA-DEFECT response /get_agents_response/agent/schedule
                type-for-t: <type> used where <t> was meant: text
  SCHEMA-DEFECT response /get_agents_response/agent/writable
                type-for-t: <type> used where <t> was meant: boolean
  SCHEMA-DEFECT response /get_agents_response/agent_count/filtered
                type-for-t: <type> used where <t> was meant: integer
  SCHEMA-DEFECT response /get_agents_response/agent_count/page
                type-for-t: <type> used where <t> was meant: integer
  SCHEMA-DEFECT response /get_agents_response/filters/keywords/keyword/column
                type-for-t: <type> used where <t> was meant: text
  SCHEMA-DEFECT response /get_agents_response/filters/keywords/keyword/value
                type-for-t: <type> used where <t> was meant: text
  SCHEMA-DEFECT response /get_agents_response/filters/term
                type-for-t: <type> used where <t> was meant: text

### get_aggregates  (command_get_aggregates.xml)

MISSING=5  EXTRA=5  CARDINALITY=2  OMITEMPTY=30

  MISSING       request  /get_aggregates/data_column
                repeated elem data_column is not modelled
  MISSING       request  /get_aggregates/sort
                repeated elem sort is not modelled
  MISSING       request  /get_aggregates/text_column
                repeated elem text_column is not modelled
  MISSING       response /get_aggregates_response/aggregate/overall/text
                repeated elem text is not modelled
  MISSING       response /get_aggregates_response/filters
                required elem filters is not modelled  [fidelity]
  EXTRA         response /get_aggregates_response/aggregate/group/c_sum
                Go field CSum string has no counterpart in the schema
  EXTRA         response /get_aggregates_response/aggregate/group/max
                Go field Max string has no counterpart in the schema
  EXTRA         response /get_aggregates_response/aggregate/group/mean
                Go field Mean string has no counterpart in the schema
  EXTRA         response /get_aggregates_response/aggregate/group/min
                Go field Min string has no counterpart in the schema
  EXTRA         response /get_aggregates_response/aggregate/group/sum
                Go field Sum string has no counterpart in the schema
  CARDINALITY   response /get_aggregates_response/aggregate
                schema says zero-or-more, Go has a single *commands.Aggregate (Aggregate)
  CARDINALITY   response /get_aggregates_response/aggregate/data_column
                schema says zero-or-more, Go has a single string (DataColumn)
  OMITEMPTY     response /get_aggregates_response/aggregate/column_info
                schema says required, ColumnInfo is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/column_info/aggregate_column/column
                schema says required, Column is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/column_info/aggregate_column/data_type
                schema says required, DataType is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/column_info/aggregate_column/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/column_info/aggregate_column/stat
                schema says required, Stat is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/column_info/aggregate_column/type
                schema says required, Type is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/data_type
                schema says required, DataType is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/c_count
                schema says required, CCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/count
                schema says required, Count is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/stats/c_sum
                schema says required, CSum is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/stats/max
                schema says required, Max is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/stats/mean
                schema says required, Mean is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/stats/min
                schema says required, Min is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/stats/sum
                schema says required, Sum is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/subgroup/c_count
                schema says required, CCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/subgroup/count
                schema says required, Count is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/subgroup/stats/c_sum
                schema says required, CSum is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/subgroup/stats/max
                schema says required, Max is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/subgroup/stats/mean
                schema says required, Mean is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/subgroup/stats/min
                schema says required, Min is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/subgroup/stats/sum
                schema says required, Sum is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/subgroup/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/group/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/overall/c_count
                schema says required, CCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/overall/count
                schema says required, Count is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/overall/stats/c_sum
                schema says required, CSum is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/overall/stats/max
                schema says required, Max is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/overall/stats/mean
                schema says required, Mean is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/overall/stats/min
                schema says required, Min is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_aggregates_response/aggregate/overall/stats/sum
                schema says required, Sum is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### get_alerts  (command_get_alerts.xml)

SCHEMA-DEFECT=2

  SCHEMA-DEFECT response /get_alerts_response/alert/filter/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_alerts_response/alert/tasks/task/permissions
                empty-pattern: pattern declares no content

### get_assets  (command_get_assets.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT response /get_assets_response/asset/host/routes/route/host/@distance
                undefined-type: number is undefined; aliased to integer

### get_configs  (command_get_configs.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT response /get_configs_response/config/tasks/task/permissions
                empty-pattern: pattern declares no content

### get_credentials  (command_get_credentials.xml)

EXTRA=2  SCHEMA-DEFECT=2

  EXTRA         response /get_credentials_response/credential/private_key_info
                Go field PrivateKeyInfo *commands.CredentialPrivateKeyInfo has no counterpart in the schema
  EXTRA         response /get_credentials_response/credential/public_key_info
                Go field PublicKeyInfo *commands.CredentialPublicKeyInfo has no counterpart in the schema
  SCHEMA-DEFECT response /get_credentials_response/credential/scanners/scanner/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_credentials_response/credential/targets/target/permissions
                empty-pattern: pattern declares no content

### get_features  (command_get_features.xml)

TYPE=1  SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /get_features
                empty-pattern: pattern declares no content
  TYPE          response /get_features_response/feature/@enabled
                schema type boolean is text ("" is true), Go Enabled is bool: '7' fails to decode and '1' re-emits as 'true'

### get_feeds  (command_get_feeds.xml)

OMITEMPTY=7  TYPE=3  SCHEMA-DEFECT=1

  OMITEMPTY     response /get_feeds_response/feed/description
                schema says required, Description is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_feeds_response/feed/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_feeds_response/feed/type
                schema says required, Type is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_feeds_response/feed/version
                schema says required, Version is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_feeds_response/feed_owner_set
                schema says required, FeedOwnerSet is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_feeds_response/feed_resources_access
                schema says required, FeedResourcesAccess is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_feeds_response/feed_roles_set
                schema says required, FeedRolesSet is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_feeds_response/feed_owner_set
                schema type boolean is text ("" is true), Go FeedOwnerSet is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_feeds_response/feed_resources_access
                schema type boolean is text ("" is true), Go FeedResourcesAccess is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_feeds_response/feed_roles_set
                schema type boolean is text ("" is true), Go FeedRolesSet is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_feeds_response/feed/sync_not_available/error
                no-pattern: <ele> has neither <pattern> nor <type>

### get_filters  (command_get_filters.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT response /get_filters_response/filter/alerts/alert/permissions
                empty-pattern: pattern declares no content

### get_groups  (command_get_groups.xml)

conformant

### get_info  (command_get_info.xml)

MISSING=12  EXTRA=7  OMITEMPTY=5  TYPE=3  SCHEMA-DEFECT=1

  MISSING       request  /get_info/@filt_id
                optional attr filt_id is not modelled
  MISSING       request  /get_info/@filter
                optional attr filter is not modelled
  MISSING       response /get_info_response/info/cpe/cpe_name_id
                required elem cpe_name_id is not modelled
  MISSING       response /get_info_response/info/cpe/severity
                required elem severity is not modelled
  MISSING       response /get_info_response/info/cve/cert/*
                schema declares 2 children; Go models this as scalar *commands.InfoCERT
  MISSING       response /get_info_response/info/cve/configuration_nodes/*
                schema declares 1 children; Go models this as scalar *commands.InfoConfigurationNodes
  MISSING       response /get_info_response/info/cve/epss/*
                schema declares 2 children; Go models this as scalar *commands.InfoEPSS
  MISSING       response /get_info_response/info/cve/nvts/nvt/*
                schema declares 2 children; Go models this as scalar []string
  MISSING       response /get_info_response/info/cve/references/reference/tags
                required elem tags is not modelled
  MISSING       response /get_info_response/info/cve/references/reference/url
                required elem url is not modelled
  MISSING       response /get_info_response/info/cve/severity
                required elem severity is not modelled  [fidelity]
  MISSING       response /get_info_response/info/dfn_cert_adv
                optional elem dfn_cert_adv is not modelled
  EXTRA         response /get_info_response/info/cpe/score
                Go field Score string has no counterpart in the schema
  EXTRA         response /get_info_response/info/cpe/status
                Go field Status string has no counterpart in the schema
  EXTRA         response /get_info_response/info/cve/references/reference/#text
                Go []commands.InfoReference has a `,chardata` field; schema declares no character data here
  EXTRA         response /get_info_response/info/cve/references/reference/@href
                Go field HREF string has no counterpart in the schema
  EXTRA         response /get_info_response/info/cve/score
                Go field Score string has no counterpart in the schema
  EXTRA         response /get_info_response/info/in_use
                Go field InUse bool has no counterpart in the schema
  EXTRA         response /get_info_response/info/writable
                Go field Writable bool has no counterpart in the schema
  OMITEMPTY     response /get_info_response/filters/name
                schema says optional, Name string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_info_response/info/@id
                schema says optional, ID string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_info_response/info/cpe/title
                schema says optional, Title string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_info_response/info/owner
                schema says required, Owner is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_info_response/info/permissions
                schema says required, Permissions is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_info_response/details
                schema type boolean is text ("" is true), Go Details is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_info_response/info/@id
                schema type integer, Go ID is string (want int)
  TYPE          response /get_info_response/info/cpe/deprecated
                schema type boolean is text ("" is true), Go Deprecated is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_info_response/info/cve/configuration_nodes/node/node
                undefined-type: node is undefined and has no alias

### get_license  (command_get_license.xml)

OMITEMPTY=3  TYPE=1  SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /get_license
                empty-pattern: pattern declares no content
  OMITEMPTY     response /get_license_response/license
                schema says required, License is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_license_response/license/content/keys/key/@name
                schema says optional, Name string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_license_response/license/content/signatures/signature/@name
                schema says optional, Name string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  TYPE          response /get_license_response/license/content/appliance/sensor
                schema type boolean is text ("" is true), Go Sensor is bool: '7' fails to decode and '1' re-emits as 'true'

### get_notes  (command_get_notes.xml)

EXTRA=1  SCHEMA-DEFECT=3

  EXTRA         response /get_notes_response/note/result/oci_image
                Go field OCIImage *commands.NoteOCIImage has no counterpart in the schema
  SCHEMA-DEFECT response /get_notes_response/note/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_notes_response/note/result/delta/overrides/override/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_notes_response/note/result/overrides/override/active
                duplicate-ele: redefined in the same scope

### get_nvt_families  (command_get_nvt_families.xml)

conformant

### get_nvts  (command_get_nvts.xml)

MISSING=6  EXTRA=1  CARDINALITY=7  OMITEMPTY=19  TYPE=4

  OMITEMPTY     request  /get_nvts/@nvt_oid
                schema says required, NvtOID is tagged omitempty  [suppressed: request-side over-constraint]
  MISSING       response /get_nvts_response/nvt/epss/*
                schema declares 2 children; Go models this as scalar string
  MISSING       response /get_nvts_response/nvt/preferences/#text
                schema also declares character data (text) here; Go *commands.NVTPreferences has no `,chardata` field
  MISSING       response /get_nvts_response/nvt/qod/*
                schema declares 2 children; Go models this as scalar string
  MISSING       response /get_nvts_response/nvt/refs/*
                schema declares 2 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_nvts_response/nvt/solution/*
                schema declares 2 children; Go models this as scalar string
  MISSING       response /get_nvts_response/nvt/user_tags/*
                schema declares 2 children; Go models this as scalar string
  EXTRA         response /get_nvts_response/nvt/description
                Go field Description string has no counterpart in the schema
  CARDINALITY   response /get_nvts_response/nvt
                schema says required, Go has a slice []commands.NVT (NVTs)
  CARDINALITY   response /get_nvts_response/nvt/category
                schema says optional; Category int cannot express absence (zero is a legal value)
  CARDINALITY   response /get_nvts_response/nvt/cvss_base
                schema says optional; CvssBase int cannot express absence (zero is a legal value)
  CARDINALITY   response /get_nvts_response/nvt/default_timeout
                schema says optional; DefaultTimeout int cannot express absence (zero is a legal value)
  CARDINALITY   response /get_nvts_response/nvt/preference_count
                schema says optional; PreferenceCount int cannot express absence (zero is a legal value)
  CARDINALITY   response /get_nvts_response/nvt/preferences/preference/alt
                schema says zero-or-more, Go has a single string (Alt)
  CARDINALITY   response /get_nvts_response/nvt/timeout
                schema says optional; Timeout int cannot express absence (zero is a legal value)
  OMITEMPTY     response /get_nvts_response/nvt
                schema says required, NVTs is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/preferences/default_timeout
                schema says required, DefaultTimeout is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/preferences/preference/default
                schema says required, Default is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/preferences/preference/hr_name
                schema says required, HRName is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/preferences/preference/id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/preferences/preference/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/preferences/preference/nvt
                schema says required, NVT is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/preferences/preference/nvt/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/preferences/preference/type
                schema says required, Type is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/preferences/preference/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/preferences/timeout
                schema says required, Timeout is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/severities/@score
                schema says required, Score is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/severities/severity/@type
                schema says required, Type is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/severities/severity/date
                schema says required, Date is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/severities/severity/origin
                schema says required, Origin is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/severities/severity/score
                schema says required, Score is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_nvts_response/nvt/severities/severity/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_nvts_response/nvt/cvss_base
                schema type text is text, Go CvssBase is int
  TYPE          response /get_nvts_response/nvt/default_timeout
                schema type text is text, Go DefaultTimeout is int
  TYPE          response /get_nvts_response/nvt/preference_count
                schema type text is text, Go PreferenceCount is int
  TYPE          response /get_nvts_response/nvt/timeout
                schema type text is text, Go Timeout is int

### get_overrides  (command_get_overrides.xml)

EXTRA=1  SCHEMA-DEFECT=3

  EXTRA         response /get_overrides_response/override/result/oci_image
                Go field OCIImage *commands.OverrideOCIImage has no counterpart in the schema
  SCHEMA-DEFECT response /get_overrides_response/override/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_overrides_response/override/result/delta/notes/note/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_overrides_response/override/result/notes/note/active
                duplicate-ele: redefined in the same scope

### get_permissions  (command_get_permissions.xml)

SCHEMA-DEFECT=2

  SCHEMA-DEFECT response /get_permissions_response/permission/resource/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_permissions_response/permission/subject/permissions
                empty-pattern: pattern declares no content

### get_port_lists  (command_get_port_lists.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT response /get_port_lists_response/port_list/targets/target/permissions
                empty-pattern: pattern declares no content

### get_preferences  (command_get_preferences.xml)

conformant

### get_report_configs  (command_get_report_configs.xml)

MISSING=1  EXTRA=1  SCHEMA-DEFECT=1

  MISSING       response /get_report_configs_response/report_config/param/value/@using_default
                required attr using_default is not modelled  [fidelity]
  EXTRA         response /get_report_configs_response/report_config/param/@using_default
                Go field UsingDefault commands.GMPBool has no counterpart in the schema  [suppressed: T-9 report_config review: real position of using_default that FieldSuppressions["get_report_configs:/report_config/param/value"] removes — command_get_report_configs.xml's own <example> and internal/gmptest/testdata/upstream/get_report_configs/001_{response,fidelity}.txt both place the attribute on <param>, not <value>]
  SCHEMA-DEFECT response /get_report_configs_response/report_config/report_format/permissions
                empty-pattern: pattern declares no content

### get_report_formats  (command_get_report_formats.xml)

EXTRA=1

  EXTRA         response /get_report_formats_response/report_format/comment
                Go field Comment string has no counterpart in the schema  [suppressed: T-9 report_format review: schema's report_format pattern never declares <comment>, but gvmd emits one on every report format — x5 in internal/gmptest/testdata/live/get_report_formats/001_response.xml, x6 in 006_response.xml (after a report format copy was added to the fixture set), and x1 upstream, all recorded as DROPPED in the matching _fidelity.txt. Every captured occurrence is empty, so it is the element's presence that is evidenced, not its content; hence no omitempty, which would re-marshal to the absence gvmd never sends]

### get_reports  (command_get_reports.xml)

MISSING=55  EXTRA=6  CARDINALITY=5  OMITEMPTY=143  TYPE=4  SCHEMA-DEFECT=12

  MISSING       request  /get_reports/@alert_id
                optional attr alert_id is not modelled
  MISSING       request  /get_reports/@config_id
                optional attr config_id is not modelled
  MISSING       request  /get_reports/@delta_report_id
                optional attr delta_report_id is not modelled
  MISSING       request  /get_reports/@details
                optional attr details is not modelled
  MISSING       request  /get_reports/@filt_id
                optional attr filt_id is not modelled
  MISSING       request  /get_reports/@filter
                optional attr filter is not modelled
  MISSING       request  /get_reports/@ignore_pagination
                optional attr ignore_pagination is not modelled
  MISSING       request  /get_reports/@lean
                optional attr lean is not modelled
  MISSING       request  /get_reports/@notes_details
                optional attr notes_details is not modelled
  MISSING       request  /get_reports/@overrides_details
                optional attr overrides_details is not modelled
  MISSING       request  /get_reports/@report_filt_id
                optional attr report_filt_id is not modelled
  MISSING       request  /get_reports/@report_filter
                optional attr report_filter is not modelled
  MISSING       request  /get_reports/@result_tags
                optional attr result_tags is not modelled
  EXTRA         request  /get_reports/@content_type
                Go field ContentType string has no counterpart in the schema
  EXTRA         request  /get_reports/@extension
                Go field Extension string has no counterpart in the schema
  MISSING       response /get_reports_response/filters
                optional elem filters is not modelled  [fidelity]
  MISSING       response /get_reports_response/report/#text
                schema also declares character data (text) here; Go []commands.Report has no `,chardata` field  [fidelity]
  MISSING       response /get_reports_response/report/report/compliance_count/#text
                schema also declares character data (text) here; Go *commands.ReportComplianceCount has no `,chardata` field  [fidelity]
  MISSING       response /get_reports_response/report/report/compliance_count/incomplete/*
                schema declares 2 children; Go models this as scalar int  [fidelity]
  MISSING       response /get_reports_response/report/report/compliance_count/no/*
                schema declares 2 children; Go models this as scalar int  [fidelity]
  MISSING       response /get_reports_response/report/report/compliance_count/undefined/*
                schema declares 2 children; Go models this as scalar int  [fidelity]
  MISSING       response /get_reports_response/report/report/compliance_count/yes/*
                schema declares 2 children; Go models this as scalar int  [fidelity]
  MISSING       response /get_reports_response/report/report/delta/report/*
                schema declares 5 children; Go models this as scalar *commands.ReportContent  [fidelity]
  MISSING       response /get_reports_response/report/report/errors/error/host/*
                schema declares 1 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_reports_response/report/report/errors/error/nvt/*
                schema declares 4 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_reports_response/report/report/filters/#text
                schema also declares character data (levels) here; Go *commands.ReportFilters has no `,chardata` field  [fidelity]
  MISSING       response /get_reports_response/report/report/filters/delta/#text
                schema also declares character data (text) here; Go *commands.ReportDeltaStates has no `,chardata` field  [fidelity]
  MISSING       response /get_reports_response/report/report/host/asset/*
                schema declares 1 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_reports_response/report/report/host/compliance_count/*
                schema declares 5 children; Go models this as scalar int  [fidelity]
  MISSING       response /get_reports_response/report/report/host/detail/source/*
                schema declares 3 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_reports_response/report/report/host/port_count/*
                schema declares 1 children; Go models this as scalar int  [fidelity]
  MISSING       response /get_reports_response/report/report/host/result_count/*
                schema declares 10 children; Go models this as scalar int  [fidelity]
  MISSING       response /get_reports_response/report/report/ports/port/#text
                schema also declares character data (text) here; Go []commands.ReportPort has no `,chardata` field  [fidelity]
  MISSING       response /get_reports_response/report/report/result_count/critical
                required elem critical is not modelled  [fidelity]
  MISSING       response /get_reports_response/report/report/result_count/high
                required elem high is not modelled  [fidelity]
  MISSING       response /get_reports_response/report/report/result_count/hole
                required elem hole is not modelled  [fidelity]
  MISSING       response /get_reports_response/report/report/result_count/info
                required elem info is not modelled  [fidelity]
  MISSING       response /get_reports_response/report/report/result_count/log
                required elem log is not modelled  [fidelity]
  MISSING       response /get_reports_response/report/report/result_count/low
                required elem low is not modelled  [fidelity]
  MISSING       response /get_reports_response/report/report/result_count/medium
                required elem medium is not modelled  [fidelity]
  MISSING       response /get_reports_response/report/report/result_count/warning
                required elem warning is not modelled  [fidelity]
  MISSING       response /get_reports_response/report/report/results/result/delta/*
                schema declares 4 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_reports_response/report/report/results/result/detection/*
                schema declares 1 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_reports_response/report/report/results/result/host/asset
                required elem asset is not modelled  [fidelity]
  MISSING       response /get_reports_response/report/report/results/result/notes/*
                schema declares 1 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_reports_response/report/report/results/result/nvt/epss/*
                schema declares 2 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_reports_response/report/report/results/result/nvt/refs/*
                schema declares 1 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_reports_response/report/report/results/result/nvt/severities/*
                schema declares 2 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_reports_response/report/report/results/result/overrides/*
                schema declares 1 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_reports_response/report/report/results/result/qod/*
                schema declares 2 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_reports_response/report/report/results/result/tickets/*
                schema declares 1 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_reports_response/report/report/sort/field/#text
                schema also declares character data (text) here; Go []commands.ReportSortField has no `,chardata` field  [fidelity]
  MISSING       response /get_reports_response/report/report/tls_certificates/tls_certificate/certificate/*
                schema declares 1 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_reports_response/report/report/tls_certificates/tls_certificate/host/*
                schema declares 2 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_reports_response/report_count
                optional elem report_count is not modelled  [fidelity]
  MISSING       response /get_reports_response/reports
                optional elem reports is not modelled  [fidelity]
  MISSING       response /get_reports_response/sort
                optional elem sort is not modelled  [fidelity]
  EXTRA         response /get_reports_response/report/base64
                Go field Base64 string has no counterpart in the schema
  EXTRA         response /get_reports_response/report/report/filters/levels
                Go field Levels string has no counterpart in the schema
  EXTRA         response /get_reports_response/report/report/results/result/host/@asset_id
                Go field AssetID string has no counterpart in the schema
  EXTRA         response /get_reports_response/report/report/tls_certificates/tls_certificate/@format
                Go field Format string has no counterpart in the schema
  CARDINALITY   response /get_reports_response/report/report/host/compliance_count
                schema says optional; ComplianceCount int cannot express absence (zero is a legal value)
  CARDINALITY   response /get_reports_response/report/report/host/result_count
                schema says optional; ResultCount int cannot express absence (zero is a legal value)
  CARDINALITY   response /get_reports_response/report/report/report_format/param
                schema says required, Go has a slice []commands.ReportContentFormatParam (Params)
  CARDINALITY   response /get_reports_response/report/report/results/result/original_severity
                schema says optional; OriginalSeverity float64 cannot express absence (zero is a legal value)
  CARDINALITY   response /get_reports_response/report/report/sort/field
                schema says required, Go has a slice []commands.ReportSortField (Fields)
  OMITEMPTY     response /get_reports_response/report/@config_id
                schema says required, ConfigID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/@content_type
                schema says required, ContentType is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/@extension
                schema says required, Extension is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/@format_id
                schema says required, FormatID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/creation_time
                schema says required, CreationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/in_use
                schema says required, InUse is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/modification_time
                schema says required, ModificationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/owner
                schema says required, Owner is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/owner/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report
                schema says required, Report is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/compliance/filtered
                schema says required, Filtered is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/compliance/full
                schema says required, Full is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/compliance_count/filtered
                schema says required, Filtered is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/compliance_count/full
                schema says required, Full is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/compliance_count/incomplete
                schema says required, Incomplete is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/compliance_count/no
                schema says required, No is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/compliance_count/undefined
                schema says required, Undefined is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/compliance_count/yes
                schema says required, Yes is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/delta/report
                schema says required, Report is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/errors/error/description
                schema says required, Description is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/errors/error/host
                schema says required, Host is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/errors/error/nvt
                schema says required, NVT is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/errors/error/port
                schema says required, Port is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/errors/error/scan_nvt_version
                schema says required, ScanNVTVersion is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/errors/error/severity
                schema says required, Severity is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/filters
                schema says required, Filters is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/filters/apply_overrides
                schema says required, ApplyOverrides is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/filters/delta/changed
                schema says required, Changed is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/filters/delta/gone
                schema says required, Gone is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/filters/delta/new
                schema says required, New is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/filters/delta/same
                schema says required, Same is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/filters/min_qod
                schema says required, MinQod is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/filters/notes
                schema says required, Notes is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/filters/overrides
                schema says required, Overrides is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/filters/phrase
                schema says required, Phrase is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/filters/result_hosts_only
                schema says required, ResultHostsOnly is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/filters/term
                schema says required, Term is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/gmp
                schema says required, GMP is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/gmp/version
                schema says required, Version is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/host/asset
                schema says required, Asset is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/host/detail/extra
                schema says required, Extra is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/host/detail/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/host/detail/source
                schema says required, Source is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/host/detail/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/host/end
                schema says required, End is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/host/ip
                schema says required, IP is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/host/port_count
                schema says required, PortCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/host/start
                schema says required, Start is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/hosts/@max
                schema says required, Max is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]
  OMITEMPTY     response /get_reports_response/report/report/hosts/@start
                schema says required, Start is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]
  OMITEMPTY     response /get_reports_response/report/report/permissions/permission/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/ports/@max
                schema says required, Max is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/ports/@start
                schema says required, Start is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/ports/port/host
                schema says required, Host is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/ports/port/severity
                schema says required, Severity is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/ports/port/threat
                schema says required, Threat is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/report_format
                schema says required, ReportFormat is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/report_format/param
                schema says required, Params is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/report_format/param/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/report_format/param/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/result_count/filtered
                schema says required, Filtered is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/result_count/full
                schema says required, Full is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/@max
                schema says required, Max is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/@start
                schema says required, Start is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/compliance
                schema says required, Compliance is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/creation_time
                schema says required, CreationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/description
                schema says required, Description is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/host
                schema says required, Host is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/host/hostname
                schema says required, Hostname is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/modification_time
                schema says required, ModificationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/nvt
                schema says required, NVT is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/nvt/@oid
                schema says required, OID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/nvt/cpe
                schema says required, CPE is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/nvt/cvss_base
                schema says required, CVSSBase is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/nvt/family
                schema says required, Family is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/nvt/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/nvt/refs
                schema says required, Refs is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/nvt/severities
                schema says required, Severities is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/nvt/tags
                schema says required, Tags is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/nvt/type
                schema says required, Type is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/owner
                schema says required, Owner is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/owner/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/port
                schema says required, Port is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/qod
                schema says required, QOD is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/report/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/scan_nvt_version
                schema says required, ScanNVTVersion is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/severity
                schema says required, Severity is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/task/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/task/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/threat
                schema says required, Threat is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/user_tags/count
                schema says required, Count is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/user_tags/tag/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/user_tags/tag/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/user_tags/tag/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/results/result/user_tags/tag/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/severity/filtered
                schema says required, Filtered is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/severity/full
                schema says required, Full is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/sort
                schema says required, Sort is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/sort/field
                schema says required, Fields is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/sort/field/order
                schema says required, Order is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/task/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/task/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/task/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/task/progress
                schema says required, Progress is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/task/target
                schema says required, Target is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/task/target/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/task/target/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/task/target/trash
                schema says required, Trash is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/task/user_tags/count
                schema says required, Count is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/task/user_tags/tag/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/task/user_tags/tag/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/task/user_tags/tag/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/task/user_tags/tag/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/tls_certificates/tls_certificate/activation_time
                schema says required, ActivationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/tls_certificates/tls_certificate/certificate
                schema says required, Certificate is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/tls_certificates/tls_certificate/expiration_time
                schema says required, ExpirationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/tls_certificates/tls_certificate/host
                schema says required, Host is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/tls_certificates/tls_certificate/issuer_dn
                schema says required, IssuerDN is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/tls_certificates/tls_certificate/md5_fingerprint
                schema says required, MD5Fingerprint is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/tls_certificates/tls_certificate/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/tls_certificates/tls_certificate/serial
                schema says required, Serial is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/tls_certificates/tls_certificate/sha256_fingerprint
                schema says required, SHA256Fingerprint is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/tls_certificates/tls_certificate/subject_dn
                schema says required, SubjectDN is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/tls_certificates/tls_certificate/valid
                schema says required, Valid is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/user_tags/count
                schema says required, Count is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/user_tags/tag/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/user_tags/tag/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/user_tags/tag/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report/user_tags/tag/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report_config/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report_config/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report_format/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/report_format/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/task
                schema says required, Task is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/task/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/task/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_reports_response/report/writable
                schema says required, Writable is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_reports_response/report/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_reports_response/report/report/task/target/trash
                schema type boolean is text ("" is true), Go Trash is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_reports_response/report/report/tls_certificates/tls_certificate/valid
                schema type boolean is text ("" is true), Go Valid is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_reports_response/report/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_reports_response/report/name
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_reports_response/report/report/apps/count
                no-pattern: <ele> has neither <pattern> nor <type>
  SCHEMA-DEFECT response /get_reports_response/report/report/closed_cves/count
                no-pattern: <ele> has neither <pattern> nor <type>
  SCHEMA-DEFECT response /get_reports_response/report/report/errors/count
                no-pattern: <ele> has neither <pattern> nor <type>
  SCHEMA-DEFECT response /get_reports_response/report/report/hosts/count
                no-pattern: <ele> has neither <pattern> nor <type>
  SCHEMA-DEFECT response /get_reports_response/report/report/os/count
                no-pattern: <ele> has neither <pattern> nor <type>
  SCHEMA-DEFECT response /get_reports_response/report/report/results/result/delta/notes/note/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_reports_response/report/report/results/result/delta/overrides/override/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_reports_response/report/report/results/result/notes/note/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_reports_response/report/report/results/result/overrides/override/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_reports_response/report/report/ssl_certs/count
                no-pattern: <ele> has neither <pattern> nor <type>
  SCHEMA-DEFECT response /get_reports_response/report/report/vulns/count
                no-pattern: <ele> has neither <pattern> nor <type>

### get_resource_names  (command_get_resource_names.xml)

OMITEMPTY=2  TYPE=1  SCHEMA-DEFECT=1

  OMITEMPTY     response /get_resource_names_response/resource/@id
                schema says optional, ID string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_resource_names_response/type
                schema says required, Type is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_resource_names_response/resource/@id
                schema type integer, Go ID is string (want int)
  SCHEMA-DEFECT response /get_resource_names_response/type
                undefined-type: type is undefined and has no alias

### get_results  (command_get_results.xml)

MISSING=1  EXTRA=4  CARDINALITY=2  OMITEMPTY=33  SCHEMA-DEFECT=4

  MISSING       response /get_results_response/result/nvt/severities/severity
                repeated elem severity is not modelled  [fidelity]
  EXTRA         response /get_results_response/result/delta/notes/note/result/*
                schema declares a "result" leaf; Go models 8 children
  EXTRA         response /get_results_response/result/delta/overrides/override/result/*
                schema declares a "result" leaf; Go models 8 children
  EXTRA         response /get_results_response/result/notes/note/result/*
                schema declares a "result" leaf; Go models 8 children
  EXTRA         response /get_results_response/result/overrides/override/result/*
                schema declares a "result" leaf; Go models 8 children
  CARDINALITY   response /get_results_response/result/delta/notes/note/severity
                schema says optional; Severity float64 cannot express absence (zero is a legal value)
  CARDINALITY   response /get_results_response/result/notes/note/severity
                schema says optional; Severity float64 cannot express absence (zero is a legal value)
  OMITEMPTY     response /get_results_response/result/delta/diff
                schema says optional, Diff string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/delta/notes
                schema says optional, Notes commands.ResultDeltaNotes has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/delta/notes/note/end_time
                schema says optional, EndTime string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/delta/notes/note/hosts
                schema says optional, Hosts string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/delta/notes/note/port
                schema says optional, Port string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/delta/notes/note/result
                schema says optional, Result commands.ResultNoteResult has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/delta/notes/note/severity
                schema says optional, Severity float64 has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/delta/notes/note/task
                schema says optional, Task commands.ResultNoteTask has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/delta/notes/note/user_tags
                schema says optional, UserTags commands.ResultNoteUserTags has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/delta/overrides
                schema says optional, Overrides commands.ResultDeltaOverrides has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/delta/overrides/override/end_time
                schema says optional, EndTime string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/delta/overrides/override/hosts
                schema says optional, Hosts string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/delta/overrides/override/port
                schema says optional, Port string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/delta/overrides/override/result
                schema says optional, Result commands.ResultOverrideResult has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/delta/overrides/override/task
                schema says optional, Task commands.ResultOverrideTask has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/delta/overrides/override/user_tags
                schema says optional, UserTags commands.ResultOverrideUserTags has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/notes/note/end_time
                schema says optional, EndTime string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/notes/note/hosts
                schema says optional, Hosts string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/notes/note/port
                schema says optional, Port string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/notes/note/result
                schema says optional, Result commands.ResultNoteResult has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/notes/note/severity
                schema says optional, Severity float64 has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/notes/note/task
                schema says optional, Task commands.ResultNoteTask has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/notes/note/user_tags
                schema says optional, UserTags commands.ResultNoteUserTags has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/nvt/epss/max_epss/cve/@id
                schema says optional, ID string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/nvt/epss/max_severity/cve/@id
                schema says optional, ID string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/nvt/refs/ref/@id
                schema says optional, ID string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/nvt/refs/ref/@type
                schema says optional, Type string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/overrides/override/end_time
                schema says optional, EndTime string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/overrides/override/hosts
                schema says optional, Hosts string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/overrides/override/port
                schema says optional, Port string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/overrides/override/result
                schema says optional, Result commands.ResultOverrideResult has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/overrides/override/task
                schema says optional, Task commands.ResultOverrideTask has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/overrides/override/user_tags
                schema says optional, UserTags commands.ResultOverrideUserTags has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  SCHEMA-DEFECT response /get_results_response/result/delta/notes/note/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_results_response/result/delta/overrides/override/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_results_response/result/notes/note/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_results_response/result/overrides/override/active
                duplicate-ele: redefined in the same scope

### get_roles  (command_get_roles.xml)

conformant

### get_scanners  (command_get_scanners.xml)

EXTRA=2  SCHEMA-DEFECT=10

  EXTRA         response /get_scanners_response/scanner/disabled
                Go field Disabled *commands.GMPBool has no counterpart in the schema
  EXTRA         response /get_scanners_response/scanner/tasks/task/usage_type
                Go field UsageType string has no counterpart in the schema
  SCHEMA-DEFECT response /get_scanners_response/scanner/configs/config/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_scanners_response/scanner/info/daemon/name
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_scanners_response/scanner/info/daemon/version
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_scanners_response/scanner/info/description
                no-pattern: <ele> has neither <pattern> nor <type>
  SCHEMA-DEFECT response /get_scanners_response/scanner/info/params/param
                no-pattern: <ele> has neither <pattern> nor <type>
  SCHEMA-DEFECT response /get_scanners_response/scanner/info/protocol/name
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_scanners_response/scanner/info/protocol/version
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_scanners_response/scanner/info/scanner/name
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_scanners_response/scanner/info/scanner/version
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_scanners_response/scanner/tasks/task/permissions
                empty-pattern: pattern declares no content

### get_schedules  (command_get_schedules.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT response /get_schedules_response/schedule/tasks/task/permissions
                empty-pattern: pattern declares no content

### get_settings  (command_get_settings.xml)

EXTRA=2  OMITEMPTY=11

  EXTRA         response /get_settings_response/setting/@id
                Go field ID string has no counterpart in the schema
  EXTRA         response /get_settings_response/settings/setting
                Go field Setting []commands.Setting has no counterpart in the schema
  OMITEMPTY     response /get_settings_response/filters
                schema says required, Filters is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_settings_response/filters/term
                schema says required, Term is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_settings_response/setting/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_settings_response/setting/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_settings_response/setting/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_settings_response/setting_count
                schema says required, SettingCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_settings_response/setting_count/filtered
                schema says required, Filtered is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_settings_response/setting_count/page
                schema says required, Page is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_settings_response/settings
                schema says required, Settings is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_settings_response/settings/@max
                schema says required, Max is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_settings_response/settings/@start
                schema says required, Start is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### get_system_reports  (command_get_system_reports.xml)

EXTRA=1

  EXTRA         response /get_system_reports_response/truncated
                Go field Truncated string has no counterpart in the schema

### get_tags  (command_get_tags.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT response /get_tags_response/tag/resources/resource/permissions
                empty-pattern: pattern declares no content

### get_targets  (command_get_targets.xml)

EXTRA=1  SCHEMA-DEFECT=9

  EXTRA         response /get_targets_response/target/alive_tests/*
                schema declares a "alive_test" leaf; Go models 1 children
  SCHEMA-DEFECT response /get_targets_response/target/esxi_credential/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_targets_response/target/krb5_credential/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_targets_response/target/permissions
                duplicate-ref: referenced more than once in the same pattern
  SCHEMA-DEFECT response /get_targets_response/target/port_list/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_targets_response/target/smb_credential/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_targets_response/target/snmp_credential/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_targets_response/target/ssh_credential/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_targets_response/target/ssh_elevate_credential/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_targets_response/target/tasks/task/permissions
                empty-pattern: pattern declares no content

### get_tasks  (command_get_tasks.xml)

MISSING=1  SCHEMA-DEFECT=10

  MISSING       response /get_tasks_response/task/hosts_ordering
                optional elem hosts_ordering is not modelled  [fidelity]
  SCHEMA-DEFECT response /get_tasks_response/task/alert/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_tasks_response/task/config/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_tasks_response/task/last_report/report/result_count/hole
                duplicate-ref: referenced more than once in the same pattern
  SCHEMA-DEFECT response /get_tasks_response/task/last_report/report/result_count/info
                duplicate-ref: referenced more than once in the same pattern
  SCHEMA-DEFECT response /get_tasks_response/task/last_report/report/result_count/warning
                duplicate-ref: referenced more than once in the same pattern
  SCHEMA-DEFECT response /get_tasks_response/task/name
                duplicate-ref: referenced more than once in the same pattern
  SCHEMA-DEFECT response /get_tasks_response/task/scanner/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_tasks_response/task/schedule
                duplicate-ref: referenced more than once in the same pattern
  SCHEMA-DEFECT response /get_tasks_response/task/schedule_periods
                duplicate-ref: referenced more than once in the same pattern
  SCHEMA-DEFECT response /get_tasks_response/task/target/permissions
                empty-pattern: pattern declares no content

### get_tickets  (command_get_tickets.xml)

conformant

### get_tls_certificates  (command_get_tls_certificates.xml)

EXTRA=2

  EXTRA         response /get_tls_certificates_response/tls_certificate/certificate/#text
                Go commands.TLSCertificateCertificate has a `,chardata` field; schema declares no character data here  [suppressed: T-9 finding 2: schema's certificate pattern declares only @format, no chardata token, but internal/gmptest/testdata/upstream/get_tls_certificates/001_response.xml sends real base64 certificate bytes as this element's text (see IncludeCertificateData)]
  EXTRA         response /get_tls_certificates_response/tls_certificate/valid
                Go field Valid commands.GMPBool has no counterpart in the schema  [suppressed: T-9 finding 1: schema's tls_certificate pattern never declares <valid>, but internal/gmptest/testdata/live/get_tls_certificates/003_response.xml and the upstream example both send <valid>1</valid> as real data]

### get_users  (command_get_users.xml)

SCHEMA-DEFECT=2

  SCHEMA-DEFECT response /get_users_response/user/groups/group/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_users_response/user/role/permissions
                empty-pattern: pattern declares no content

### get_version  (command_get_version.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /get_version
                empty-pattern: pattern declares no content

### get_vulns  (command_get_vulns.xml)

MISSING=1  EXTRA=3  OMITEMPTY=8

  EXTRA         request  /get_vulns/@trash
                Go field Trash *commands.GMPBool has no counterpart in the schema
  MISSING       response /get_vulns_response/sort/#text
                schema also declares character data (text) here; Go *commands.VulnsSort has no `,chardata` field
  EXTRA         response /get_vulns_response/truncated
                Go field Truncated string has no counterpart in the schema
  EXTRA         response /get_vulns_response/vuln/type
                Go field Type string has no counterpart in the schema
  OMITEMPTY     response /get_vulns_response/filters
                schema says required, Filters is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_vulns_response/filters/keywords
                schema says required, Keywords is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_vulns_response/sort
                schema says required, Sort is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_vulns_response/sort/field
                schema says required, Field is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_vulns_response/vuln/hosts
                schema says required, Hosts is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_vulns_response/vuln/results
                schema says required, Results is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_vulns_response/vuln_count
                schema says required, VulnCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_vulns_response/vulns
                schema says required, VulnsMeta is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### help  (command_help.xml)

MISSING=1  SCHEMA-DEFECT=12

  MISSING       response /help_response/schema
                optional elem schema is not modelled  [fidelity]
  SCHEMA-DEFECT response /help_response/schema/protocol/command/pattern/any/g/any
                recursion-cut: any reaches itself; the walk is cut here
  SCHEMA-DEFECT response /help_response/schema/protocol/command/pattern/any/g/o/g
                recursion-cut: g reaches itself; the walk is cut here
  SCHEMA-DEFECT response /help_response/schema/protocol/command/pattern/g/any/g
                recursion-cut: g reaches itself; the walk is cut here
  SCHEMA-DEFECT response /help_response/schema/protocol/command/pattern/g/o/g
                recursion-cut: g reaches itself; the walk is cut here
  SCHEMA-DEFECT response /help_response/schema/protocol/command/pattern/o/g/any/g
                recursion-cut: g reaches itself; the walk is cut here
  SCHEMA-DEFECT response /help_response/schema/protocol/command/pattern/o/g/o
                recursion-cut: o reaches itself; the walk is cut here
  SCHEMA-DEFECT response /help_response/schema/protocol/command/response/any/g/any
                recursion-cut: any reaches itself; the walk is cut here
  SCHEMA-DEFECT response /help_response/schema/protocol/command/response/any/g/o/g
                recursion-cut: g reaches itself; the walk is cut here
  SCHEMA-DEFECT response /help_response/schema/protocol/command/response/g/any/g
                recursion-cut: g reaches itself; the walk is cut here
  SCHEMA-DEFECT response /help_response/schema/protocol/command/response/g/o/g
                recursion-cut: g reaches itself; the walk is cut here
  SCHEMA-DEFECT response /help_response/schema/protocol/command/response/o/g/any/g
                recursion-cut: g reaches itself; the walk is cut here
  SCHEMA-DEFECT response /help_response/schema/protocol/command/response/o/g/o
                recursion-cut: o reaches itself; the walk is cut here

### modify_agents  (command_modify_agents.xml)

EXTRA=18

  EXTRA         request  /modify_agents/agents/agent/agent_id
                Go field AgentID string has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/authorized
                Go field Authorized int has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/comment
                Go field Comment string has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/connection_status
                Go field ConnectionStatus string has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/creation_time
                Go field CreationTime commands.GMPTime has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/heartbeat_interval
                Go field HeartbeatInterval int has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/hostname
                Go field Hostname string has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/in_use
                Go field InUse bool has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/ip
                Go field IP []string has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/last_update
                Go field LastUpdate commands.GMPTime has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/min_interval
                Go field MinInterval int has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/modification_time
                Go field ModificationTime commands.GMPTime has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/owner
                Go field Owner *commands.AgentOwner has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/permissions
                Go field Permissions *commands.AgentPermissions has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/scanner
                Go field Scanner *commands.AgentScanner has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/schedule
                Go field Schedule string has no counterpart in the schema
  EXTRA         request  /modify_agents/agents/agent/writable
                Go field Writable bool has no counterpart in the schema

### modify_alert  (command_modify_alert.xml)

conformant

### modify_asset  (command_modify_asset.xml)

conformant

### modify_auth  (command_modify_auth.xml)

EXTRA=1

  EXTRA         request  /modify_auth/group/@name
                Go field Name string has no counterpart in the schema

### modify_config  (command_modify_config.xml)

conformant

### modify_credential  (command_modify_credential.xml)

conformant

### modify_filter  (command_modify_filter.xml)

conformant

### modify_group  (command_modify_group.xml)

conformant

### modify_license  (command_modify_license.xml)

conformant

### modify_note  (command_modify_note.xml)

conformant

### modify_override  (command_modify_override.xml)

conformant

### modify_permission  (command_modify_permission.xml)

conformant

### modify_port_list  (command_modify_port_list.xml)

conformant

### modify_report_config  (command_modify_report_config.xml)

EXTRA=1

  EXTRA         request  /modify_report_config/@report_config_id
                Go field ReportConfigID string has no counterpart in the schema  [suppressed: T-9 report_config review: command_modify_report_config.xml's <pattern> declares no selector attribute at all, unlike every other modify_X command, but a live gvmd 26.36.1 test proved gvmd requires one anyway — <modify_report_config report_config_id="..."> against a nonexistent id returned "404 Failed to find Report Config" (a resource lookup, not a parse rejection), reproduced on a second connection. Without this field the generated request could never target a specific report config]

### modify_report_format  (command_modify_report_format.xml)

conformant

### modify_role  (command_modify_role.xml)

conformant

### modify_scanner  (command_modify_scanner.xml)

conformant

### modify_schedule  (command_modify_schedule.xml)

conformant

### modify_setting  (command_modify_setting.xml)

conformant

### modify_tag  (command_modify_tag.xml)

conformant

### modify_target  (command_modify_target.xml)

MISSING=1  EXTRA=1

  MISSING       request  /modify_target/hosts_ordering
                optional elem hosts_ordering is not modelled
  EXTRA         request  /modify_target/alive_tests/*
                schema declares a "alive_test" leaf; Go models 1 children

### modify_task  (command_modify_task.xml)

conformant

### modify_ticket  (command_modify_ticket.xml)

conformant

### modify_tls_certificate  (command_modify_tls_certificate.xml)

conformant

### modify_user  (command_modify_user.xml)

conformant

### move_task  (command_move_task.xml)

conformant

### restore  (command_restore.xml)

conformant

### resume_task  (command_resume_task.xml)

OMITEMPTY=1

  OMITEMPTY     response /resume_task_response/report_id
                schema says required, ReportID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### run_wizard  (command_run_wizard.xml)

OMITEMPTY=1  SCHEMA-DEFECT=1

  OMITEMPTY     request  /run_wizard/params
                schema says required, Params is tagged omitempty  [suppressed: request-side over-constraint]
  SCHEMA-DEFECT response /run_wizard_response/response
                empty-pattern: pattern declares no content

### start_task  (command_start_task.xml)

OMITEMPTY=1

  OMITEMPTY     response /start_task_response/report_id
                schema says required, ReportID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### stop_task  (command_stop_task.xml)

conformant

### sync_config  (command_sync_config.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /sync_config
                empty-pattern: pattern declares no content

### test_alert  (command_test_alert.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT response /test_alert_response/status_details
                dangling-e: no <ele> defines this name in scope

### verify_report_format  (command_verify_report_format.xml)

conformant

### verify_scanner  (command_verify_scanner.xml)

EXTRA=1

  EXTRA         response /verify_scanner_response/version
                Go field Version string has no counterpart in the schema  [suppressed: T-9 scanner review: command_verify_scanner.xml's response <pattern> declares only status/status_text (its <example> shows <version> but the pattern never does), yet a live gvmd 26.36.1 test against two real scanners returned a populated <version> element both times ("OpenVAS 23.50.0", "GVM/26.36.1") — real, populated wire data, so dropping it would be a behavior regression rather than a defect fix]

