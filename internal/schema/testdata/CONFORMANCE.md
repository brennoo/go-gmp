# GMP schema conformance report

Regenerate with: go test ./internal/schema -run TestConformanceReport -update

What this is: the GMP schema DSL, read directly, diffed against the
hand-written structs in commands/. The schema revision compared here is the one
vendored at vendor/XML/commands, which matches commands/ one-for-one, so a
difference found here is a transcription error and not version drift.

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
almost all request-side over-constraint: the schema marks create_task's config,
target and scanner required even though the agent and OCI variants replace
them, so response-side optionality is trusted and request-side is not.

## Totals

commands compared:            122
fully conformant:             47
with actionable findings:     75
actionable findings:          667  (request 291, response 376)
suppressed findings:          809
confirmed by fidelity corpus: 73
in the schema, not in Go:     0
in Go, not in the schema:     0

### Findings by category

                 actionable suppressed
  MISSING               236          0
  EXTRA                 242          0
  CARDINALITY            43          0
  KIND                   17          0
  OMITEMPTY              10        809
  TYPE                  119          0
  SCHEMA-DEFECT         109          0

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

  dangling-e           2 sites in 2 commands
  duplicate-ele       11 sites in 6 commands
  empty-pattern       44 sites in 23 commands
  no-pattern          17 sites in 4 commands
  recursion-cut        6 sites in 1 commands
  type-for-t          23 sites in 1 commands
  undefined-type       3 sites in 3 commands
  unsupported-r        3 sites in 3 commands

### Commands by actionable findings

  get_reports                  64   MISSING=49  EXTRA=6  CARDINALITY=5  TYPE=4
  create_report                46   MISSING=32  EXTRA=5  CARDINALITY=5  TYPE=4
  get_tasks                    41   MISSING=10  EXTRA=15  CARDINALITY=7  TYPE=9
  get_overrides                36   MISSING=27  EXTRA=3  TYPE=6
  get_targets                  29   MISSING=11  EXTRA=10  KIND=3  TYPE=5
  get_report_formats           28   MISSING=17  EXTRA=2  CARDINALITY=1  KIND=2  TYPE=6
  get_info                     21   MISSING=12  EXTRA=6  TYPE=3
  modify_target                19   EXTRA=19
  delete_agents                18   EXTRA=18
  get_tags                     18   MISSING=8  EXTRA=2  CARDINALITY=4  TYPE=4
  modify_agents                18   EXTRA=18
  create_target                17   EXTRA=17
  get_nvts                     17   MISSING=5  EXTRA=1  CARDINALITY=7  TYPE=4
  get_report_configs           17   MISSING=11  EXTRA=2  KIND=2  TYPE=2
  create_task                  16   EXTRA=15  CARDINALITY=1
  get_tickets                  15   MISSING=5  EXTRA=4  KIND=3  TYPE=3
  get_tls_certificates         14   MISSING=5  EXTRA=4  KIND=2  TYPE=3
  modify_config                14   EXTRA=6  CARDINALITY=3  OMITEMPTY=2  TYPE=3
  get_roles                    13   MISSING=11  TYPE=2
  modify_task                  13   MISSING=1  EXTRA=12
  ... and 55 more with fewer

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
  get_configs              /get_configs_response/config_count/filtered
  get_info                 /get_info_response/info/cve/severity
  get_nvts                 /get_nvts_response/nvt/refs/*
  get_port_lists           /get_port_lists_response/filters/keywords/*
  get_report_configs       /get_report_configs_response/filters/keywords/*
  get_report_formats       /get_report_formats_response/filters/keywords/*
  get_report_formats       /get_report_formats_response/report_format/owner
  get_report_formats       /get_report_formats_response/report_format/permissions
  get_reports              /get_reports_response/filters
  get_reports              /get_reports_response/report/report/compliance_count/incomplete/*
  get_reports              /get_reports_response/report/report/compliance_count/no/*
  get_reports              /get_reports_response/report/report/compliance_count/undefined/*
  get_reports              /get_reports_response/report/report/compliance_count/yes/*
  get_reports              /get_reports_response/report/report/delta/report/*
  get_reports              /get_reports_response/report/report/errors/error/host/*
  get_reports              /get_reports_response/report/report/errors/error/nvt/*
  get_reports              /get_reports_response/report/report/host/asset/*
  get_reports              /get_reports_response/report/report/host/compliance_count/*
  get_reports              /get_reports_response/report/report/host/detail/source/*
  get_reports              /get_reports_response/report/report/host/port_count/*
  get_reports              /get_reports_response/report/report/host/result_count/*
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
  get_reports              /get_reports_response/report/report/tls_certificates/tls_certificate/certificate/*
  get_reports              /get_reports_response/report/report/tls_certificates/tls_certificate/host/*
  get_reports              /get_reports_response/report_count
  get_reports              /get_reports_response/reports
  get_reports              /get_reports_response/sort
  get_results              /get_results_response/result/nvt/severities/severity
  get_roles                /get_roles_response/filters
  get_roles                /get_roles_response/role/owner
  get_roles                /get_roles_response/role/permissions
  get_roles                /get_roles_response/role_count
  get_roles                /get_roles_response/roles
  get_roles                /get_roles_response/sort
  get_tags                 /get_tags_response/filters
  get_tags                 /get_tags_response/sort
  get_tags                 /get_tags_response/tag/resources
  get_tags                 /get_tags_response/tag_count
  get_tags                 /get_tags_response/tags
  get_targets              /get_targets_response/filters/keywords
  get_targets              /get_targets_response/target/esxi_credential/trash
  get_targets              /get_targets_response/target/krb5_credential/trash
  get_targets              /get_targets_response/target/port_list/trash
  get_targets              /get_targets_response/target/smb_credential/trash
  get_targets              /get_targets_response/target/snmp_credential/trash
  get_targets              /get_targets_response/target/ssh_credential/trash
  get_targets              /get_targets_response/target/ssh_elevate_credential/trash
  get_tasks                /get_tasks_response/task/last_report/report/result_count
  get_tasks                /get_tasks_response/task/report_count/*
  get_tickets              /get_tickets_response/filters/keywords
  help                     /help_response/schema

## Per command

### authenticate  (command_authenticate.xml)

OMITEMPTY=2

  OMITEMPTY     response /authenticate_response/role
                schema says required, Role is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /authenticate_response/timezone
                schema says required, Timezone is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### create_alert  (command_create_alert.xml)

OMITEMPTY=1

  OMITEMPTY     request  /create_alert/filter
                schema says required, Filter is tagged omitempty  [suppressed: request-side over-constraint]

### create_asset  (command_create_asset.xml)

OMITEMPTY=1

  OMITEMPTY     request  /create_asset/report/filter
                schema says required, Filter is tagged omitempty  [suppressed: request-side over-constraint]

### create_config  (command_create_config.xml)

EXTRA=1  OMITEMPTY=1  SCHEMA-DEFECT=1

  EXTRA         request  /create_config/get_configs_response
                Go field GetConfigsResponse *commands.GetConfigsResponse has no counterpart in the schema
  OMITEMPTY     request  /create_config/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  SCHEMA-DEFECT request  /create_config
                unsupported-r: <r> references a response from inside a pattern

### create_credential  (command_create_credential.xml)

OMITEMPTY=2

  OMITEMPTY     request  /create_credential/privacy/algorithm
                schema says required, Algorithm is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_credential/privacy/password
                schema says required, Password is tagged omitempty  [suppressed: request-side over-constraint]

### create_filter  (command_create_filter.xml)

conformant

### create_group  (command_create_group.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /create_group/specials/full
                empty-pattern: pattern declares no content

### create_note  (command_create_note.xml)

conformant

### create_override  (command_create_override.xml)

MISSING=2  EXTRA=2

  MISSING       request  /create_override/result/*
                schema declares 1 children; Go models this as scalar string
  MISSING       request  /create_override/task/*
                schema declares 1 children; Go models this as scalar string
  EXTRA         request  /create_override/nvt/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_override/nvt/type
                Go field Type string has no counterpart in the schema

### create_permission  (command_create_permission.xml)

EXTRA=8  OMITEMPTY=5

  EXTRA         request  /create_permission/resource/deleted
                Go field Deleted bool has no counterpart in the schema
  EXTRA         request  /create_permission/resource/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_permission/resource/permissions
                Go field Permissions *commands.PermissionList has no counterpart in the schema
  EXTRA         request  /create_permission/resource/trash
                Go field Trash bool has no counterpart in the schema
  EXTRA         request  /create_permission/subject/deleted
                Go field Deleted bool has no counterpart in the schema
  EXTRA         request  /create_permission/subject/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_permission/subject/permissions
                Go field Permissions *commands.PermissionList has no counterpart in the schema
  EXTRA         request  /create_permission/subject/trash
                Go field Trash bool has no counterpart in the schema
  OMITEMPTY     request  /create_permission/resource/@id
                schema says optional, ID string has no omitempty and is always emitted
  OMITEMPTY     request  /create_permission/resource/type
                schema says required, Type is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_permission/subject/@id
                schema says optional, ID string has no omitempty and is always emitted
  OMITEMPTY     request  /create_permission/subject/type
                schema says required, Type is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     response /create_permission_response/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### create_port_list  (command_create_port_list.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /create_port_list
                unsupported-r: <r> references a response from inside a pattern

### create_port_range  (command_create_port_range.xml)

MISSING=1

  MISSING       request  /create_port_range/port_list/*
                schema declares 1 children; Go models this as scalar string

### create_report  (command_create_report.xml)

MISSING=32  EXTRA=5  CARDINALITY=5  OMITEMPTY=144  TYPE=4  SCHEMA-DEFECT=10

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
                schema says required, Go has a slice []commands.ReportFormatParam (Params)
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
  SCHEMA-DEFECT request  /create_report/report/report/ssl_certs/count
                no-pattern: <ele> has neither <pattern> nor <type>
  SCHEMA-DEFECT request  /create_report/report/report/vulns/count
                no-pattern: <ele> has neither <pattern> nor <type>

### create_report_config  (command_create_report_config.xml)

MISSING=1  EXTRA=4

  MISSING       request  /create_report_config/param/value/*
                schema declares 1 children; Go models this as scalar string
  EXTRA         request  /create_report_config/param/@using_default
                Go field UsingDefault bool has no counterpart in the schema
  EXTRA         request  /create_report_config/param/default
                Go field Default string has no counterpart in the schema
  EXTRA         request  /create_report_config/param/type
                Go field Type string has no counterpart in the schema
  EXTRA         request  /create_report_config/report_format/name
                Go field Name string has no counterpart in the schema

### create_report_format  (command_create_report_format.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /create_report_format
                unsupported-r: <r> references a response from inside a pattern

### create_role  (command_create_role.xml)

conformant

### create_scanner  (command_create_scanner.xml)

EXTRA=3  TYPE=2

  EXTRA         request  /create_scanner/credential/login
                Go field Login string has no counterpart in the schema
  EXTRA         request  /create_scanner/credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_scanner/credential/trash
                Go field Trash bool has no counterpart in the schema
  TYPE          request  /create_scanner/port
                schema type text is text, Go Port is int
  TYPE          request  /create_scanner/relay_port
                schema type text is text, Go RelayPort is int

### create_schedule  (command_create_schedule.xml)

conformant

### create_tag  (command_create_tag.xml)

EXTRA=3  OMITEMPTY=2

  EXTRA         request  /create_tag/resources/resource/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_tag/resources/resource/trash
                Go field Trash bool has no counterpart in the schema
  EXTRA         request  /create_tag/resources/resource/type
                Go field Type string has no counterpart in the schema
  OMITEMPTY     request  /create_tag/resources
                schema says required, Resources is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_tag/resources/resource/@id
                schema says optional, ID string has no omitempty and is always emitted

### create_target  (command_create_target.xml)

EXTRA=17  OMITEMPTY=12

  EXTRA         request  /create_target/esxi_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_target/esxi_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         request  /create_target/esxi_lsc_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_target/esxi_lsc_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         request  /create_target/krb5_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_target/krb5_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         request  /create_target/port_list/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_target/smb_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_target/smb_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         request  /create_target/smb_lsc_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_target/smb_lsc_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         request  /create_target/snmp_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_target/snmp_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         request  /create_target/ssh_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_target/ssh_elevate_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_target/ssh_elevate_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         request  /create_target/ssh_lsc_credential/name
                Go field Name string has no counterpart in the schema
  OMITEMPTY     request  /create_target/asset_hosts/@filter
                schema says required, Filter is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_target/esxi_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_target/esxi_lsc_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_target/krb5_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_target/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_target/port_list/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_target/smb_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_target/smb_lsc_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_target/snmp_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_target/ssh_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_target/ssh_elevate_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_target/ssh_lsc_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]

### create_task  (command_create_task.xml)

EXTRA=15  CARDINALITY=1  OMITEMPTY=4

  EXTRA         request  /create_task/alert/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_task/config/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_task/config/trash
                Go field Trash bool has no counterpart in the schema
  EXTRA         request  /create_task/config/type
                Go field Type string has no counterpart in the schema
  EXTRA         request  /create_task/observers/group/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_task/observers/role
                Go field Roles []commands.ObserverRole has no counterpart in the schema
  EXTRA         request  /create_task/observers/user_list
                Go field UserList string has no counterpart in the schema
  EXTRA         request  /create_task/preferences/preference/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_task/scanner/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_task/scanner/trash
                Go field Trash bool has no counterpart in the schema
  EXTRA         request  /create_task/scanner/type
                Go field Type string has no counterpart in the schema
  EXTRA         request  /create_task/schedule/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_task/schedule/trash
                Go field Trash bool has no counterpart in the schema
  EXTRA         request  /create_task/target/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_task/target/trash
                Go field Trash bool has no counterpart in the schema
  CARDINALITY   request  /create_task/alert
                schema says zero-or-more, Go has a single *commands.TaskAlert (Alert)
  OMITEMPTY     request  /create_task/config
                schema says required, Config is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_task/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_task/scanner
                schema says required, Scanner is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /create_task/target
                schema says required, Target is tagged omitempty  [suppressed: request-side over-constraint]

### create_ticket  (command_create_ticket.xml)

conformant

### create_tls_certificate  (command_create_tls_certificate.xml)

conformant

### create_user  (command_create_user.xml)

EXTRA=2  CARDINALITY=1  OMITEMPTY=1

  EXTRA         request  /create_user/groups/group/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /create_user/role/name
                Go field Name string has no counterpart in the schema
  CARDINALITY   request  /create_user/sources/source
                schema says optional, Go has a slice []string (Source)
  OMITEMPTY     request  /create_user/hosts/@allow
                schema says optional, Allow string has no omitempty and is always emitted

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

TYPE=1

  TYPE          request  /delete_alert/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

### delete_asset  (command_delete_asset.xml)

OMITEMPTY=2

  OMITEMPTY     request  /delete_asset/@asset_id
                schema says required, AssetID is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]
  OMITEMPTY     request  /delete_asset/@report_id
                schema says required, ReportID is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]

### delete_config  (command_delete_config.xml)

TYPE=1

  TYPE          request  /delete_config/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

### delete_credential  (command_delete_credential.xml)

TYPE=1

  TYPE          request  /delete_credential/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

### delete_filter  (command_delete_filter.xml)

TYPE=1

  TYPE          request  /delete_filter/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

### delete_group  (command_delete_group.xml)

TYPE=1

  TYPE          request  /delete_group/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

### delete_note  (command_delete_note.xml)

TYPE=1

  TYPE          request  /delete_note/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

### delete_override  (command_delete_override.xml)

TYPE=1

  TYPE          request  /delete_override/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

### delete_permission  (command_delete_permission.xml)

TYPE=1

  TYPE          request  /delete_permission/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

### delete_port_list  (command_delete_port_list.xml)

TYPE=1

  TYPE          request  /delete_port_list/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

### delete_port_range  (command_delete_port_range.xml)

conformant

### delete_report  (command_delete_report.xml)

conformant

### delete_report_config  (command_delete_report_config.xml)

TYPE=1

  TYPE          request  /delete_report_config/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

### delete_report_format  (command_delete_report_format.xml)

TYPE=1

  TYPE          request  /delete_report_format/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

### delete_role  (command_delete_role.xml)

TYPE=1

  TYPE          request  /delete_role/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

### delete_scanner  (command_delete_scanner.xml)

TYPE=1

  TYPE          request  /delete_scanner/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

### delete_schedule  (command_delete_schedule.xml)

TYPE=1

  TYPE          request  /delete_schedule/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

### delete_tag  (command_delete_tag.xml)

TYPE=1

  TYPE          request  /delete_tag/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

### delete_target  (command_delete_target.xml)

OMITEMPTY=2

  OMITEMPTY     request  /delete_target/@target_id
                schema says required, TargetID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /delete_target/@ultimate
                schema says required, Ultimate is tagged omitempty  [suppressed: request-side over-constraint]

### delete_task  (command_delete_task.xml)

OMITEMPTY=2

  OMITEMPTY     request  /delete_task/@task_id
                schema says required, TaskID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /delete_task/@ultimate
                schema says required, Ultimate is tagged omitempty  [suppressed: request-side over-constraint]

### delete_ticket  (command_delete_ticket.xml)

TYPE=1

  TYPE          request  /delete_ticket/@ultimate
                schema type boolean is text ("" is true), Go Ultimate is bool: '7' fails to decode and '1' re-emits as 'true'

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

MISSING=2  OMITEMPTY=7  TYPE=4  SCHEMA-DEFECT=2

  MISSING       response /get_alerts_response/alert/method/data/certificate_info
                optional elem certificate_info is not modelled
  MISSING       response /get_alerts_response/alert/method/data/credential
                optional elem credential is not modelled
  OMITEMPTY     response /get_alerts_response/alert/filter/@id
                schema says optional, ID string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_alerts_response/alert/filter/trash
                schema says required, Trash is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_alerts_response/alert/tasks
                schema says optional,  (wrapper) has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_alerts_response/alert_count
                schema says required, AlertCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_alerts_response/alerts
                schema says required, AlertsInfo is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_alerts_response/filters
                schema says required, Filters is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_alerts_response/sort
                schema says required, Sort is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_alerts_response/alert/active
                schema type boolean is text ("" is true), Go Active is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_alerts_response/alert/filter/trash
                schema type boolean is text ("" is true), Go Trash is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_alerts_response/alert/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_alerts_response/alert/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_alerts_response/alert/filter/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_alerts_response/alert/tasks/task/permissions
                empty-pattern: pattern declares no content

### get_assets  (command_get_assets.xml)

EXTRA=6  CARDINALITY=1  OMITEMPTY=9  TYPE=4  SCHEMA-DEFECT=1

  EXTRA         response /get_assets_response/asset/esxi_credential
                Go field ESXICredential *commands.AssetCredential has no counterpart in the schema
  EXTRA         response /get_assets_response/asset/hosts
                Go field Hosts string has no counterpart in the schema
  EXTRA         response /get_assets_response/asset/max_hosts
                Go field MaxHosts int has no counterpart in the schema
  EXTRA         response /get_assets_response/asset/smb_credential
                Go field SMBCredential *commands.AssetCredential has no counterpart in the schema
  EXTRA         response /get_assets_response/asset/ssh_credential
                Go field SSHCredential *commands.AssetCredential has no counterpart in the schema
  EXTRA         response /get_assets_response/asset/tasks
                Go field Tasks *commands.AssetTasks has no counterpart in the schema
  CARDINALITY   response /get_assets_response/asset/identifiers/identifier
                schema says required, Go has a slice []commands.AssetIdentifier (Identifier)
  OMITEMPTY     response /get_assets_response/asset/host/routes/route/host/@distance
                schema says optional, Distance string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_assets_response/asset/host/routes/route/host/@id
                schema says optional, ID string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_assets_response/asset/host/routes/route/host/@same_source
                schema says optional, SameSource string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_assets_response/asset/in_use
                schema says required, InUse is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_assets_response/asset/writable
                schema says required, Writable is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_assets_response/asset_count
                schema says required, AssetCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_assets_response/assets
                schema says required, AssetsInfo is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_assets_response/filters
                schema says required, Filters is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_assets_response/sort
                schema says required, Sort is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_assets_response/asset/host/routes/route/host/@distance
                schema type integer, Go Distance is string (want int)
  TYPE          response /get_assets_response/asset/host/severity/value
                schema type severity, Go Value is string (want float64)
  TYPE          response /get_assets_response/asset/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_assets_response/asset/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_assets_response/asset/host/routes/route/host/@distance
                undefined-type: number is undefined; aliased to integer

### get_configs  (command_get_configs.xml)

MISSING=2  EXTRA=1  OMITEMPTY=3  TYPE=2  SCHEMA-DEFECT=1

  MISSING       response /get_configs_response/config/user_tags/tag
                repeated elem tag is not modelled
  MISSING       response /get_configs_response/config_count/filtered
                required elem filtered is not modelled  [fidelity]
  EXTRA         response /get_configs_response/config/families/family/type
                Go field Type int has no counterpart in the schema
  OMITEMPTY     response /get_configs_response/config/tasks
                schema says required, Tasks is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_configs_response/config/user_tags
                schema says optional, UserTags commands.ConfigUserTags has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_configs_response/filters/name
                schema says optional, Name string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  TYPE          response /get_configs_response/config/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_configs_response/config/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_configs_response/config/tasks/task/permissions
                empty-pattern: pattern declares no content

### get_credentials  (command_get_credentials.xml)

MISSING=2  EXTRA=1  OMITEMPTY=2  TYPE=3  SCHEMA-DEFECT=2

  MISSING       response /get_credentials_response/credential/certificate_info/*
                schema declares 8 children; Go models this as scalar *commands.CredentialCertificateInfo
  MISSING       response /get_credentials_response/credential/privacy/*
                schema declares 1 children; Go models this as scalar *commands.CredentialPrivacy
  EXTRA         response /get_credentials_response/truncated
                Go field Truncated string has no counterpart in the schema
  OMITEMPTY     response /get_credentials_response/credential/formats
                schema says required, Formats is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_credentials_response/filters/name
                schema says optional, Name string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  TYPE          response /get_credentials_response/credential/allow_insecure
                schema type boolean is text ("" is true), Go AllowInsecure is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_credentials_response/credential/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_credentials_response/credential/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
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

OMITEMPTY=1  TYPE=2  SCHEMA-DEFECT=3

  OMITEMPTY     response /get_filters_response/filters/name
                schema says optional, Name string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  TYPE          response /get_filters_response/filter/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_filters_response/filter/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_filters_response/filter/alerts/alert/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_filters_response/filters
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_filters_response/filters2
                dangling-e: no <ele> defines this name in scope

### get_groups  (command_get_groups.xml)

OMITEMPTY=1  TYPE=2

  OMITEMPTY     response /get_groups_response/filters/name
                schema says optional, Name string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  TYPE          response /get_groups_response/group/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_groups_response/group/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'

### get_info  (command_get_info.xml)

MISSING=12  EXTRA=6  OMITEMPTY=5  TYPE=3  SCHEMA-DEFECT=1

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

SCHEMA-DEFECT=2

  SCHEMA-DEFECT response /get_notes_response/note/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_notes_response/note/result/delta/overrides/override/active
                duplicate-ele: redefined in the same scope

### get_nvt_families  (command_get_nvt_families.xml)

conformant

### get_nvts  (command_get_nvts.xml)

MISSING=5  EXTRA=1  CARDINALITY=7  OMITEMPTY=19  TYPE=4

  OMITEMPTY     request  /get_nvts/@nvt_oid
                schema says required, NvtOID is tagged omitempty  [suppressed: request-side over-constraint]
  MISSING       response /get_nvts_response/nvt/epss/*
                schema declares 2 children; Go models this as scalar string
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

MISSING=27  EXTRA=3  OMITEMPTY=40  TYPE=6  SCHEMA-DEFECT=2

  MISSING       request  /get_overrides/@nvt_oid
                optional attr nvt_oid is not modelled
  MISSING       request  /get_overrides/@task_id
                optional attr task_id is not modelled
  EXTRA         request  /get_overrides/@trash
                Go field Trash string has no counterpart in the schema
  MISSING       response /get_overrides_response/override/result/comment
                required elem comment is not modelled
  MISSING       response /get_overrides_response/override/result/compliance
                required elem compliance is not modelled
  MISSING       response /get_overrides_response/override/result/creation_time
                required elem creation_time is not modelled
  MISSING       response /get_overrides_response/override/result/delta
                optional elem delta is not modelled
  MISSING       response /get_overrides_response/override/result/detection
                optional elem detection is not modelled
  MISSING       response /get_overrides_response/override/result/host/hostname
                required elem hostname is not modelled
  MISSING       response /get_overrides_response/override/result/modification_time
                required elem modification_time is not modelled
  MISSING       response /get_overrides_response/override/result/name
                required elem name is not modelled
  MISSING       response /get_overrides_response/override/result/notes
                optional elem notes is not modelled
  MISSING       response /get_overrides_response/override/result/nvt/cpe
                required elem cpe is not modelled
  MISSING       response /get_overrides_response/override/result/nvt/epss
                optional elem epss is not modelled
  MISSING       response /get_overrides_response/override/result/nvt/family
                required elem family is not modelled
  MISSING       response /get_overrides_response/override/result/nvt/refs
                required elem refs is not modelled
  MISSING       response /get_overrides_response/override/result/nvt/severities
                required elem severities is not modelled
  MISSING       response /get_overrides_response/override/result/nvt/tags
                required elem tags is not modelled
  MISSING       response /get_overrides_response/override/result/original_severity
                optional elem original_severity is not modelled
  MISSING       response /get_overrides_response/override/result/original_threat
                optional elem original_threat is not modelled
  MISSING       response /get_overrides_response/override/result/overrides
                optional elem overrides is not modelled
  MISSING       response /get_overrides_response/override/result/owner
                required elem owner is not modelled
  MISSING       response /get_overrides_response/override/result/path
                optional elem path is not modelled
  MISSING       response /get_overrides_response/override/result/report
                optional elem report is not modelled
  MISSING       response /get_overrides_response/override/result/scan_nvt_version
                required elem scan_nvt_version is not modelled
  MISSING       response /get_overrides_response/override/result/task
                optional elem task is not modelled
  MISSING       response /get_overrides_response/override/result/tickets
                optional elem tickets is not modelled
  MISSING       response /get_overrides_response/override/result/user_tags
                optional elem user_tags is not modelled
  EXTRA         response /get_overrides_response/override/result/nvt/bid
                Go field BID int has no counterpart in the schema
  EXTRA         response /get_overrides_response/override/result/nvt/cve
                Go field CVE string has no counterpart in the schema
  OMITEMPTY     response /get_overrides_response/filters
                schema says required, Filters is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/active
                schema says required, Active is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/creation_time
                schema says required, CreationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/in_use
                schema says required, InUse is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/modification_time
                schema says required, ModificationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/new_severity
                schema says required, NewSeverity is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/new_threat
                schema says required, NewThreat is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/nvt
                schema says required, NVT is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/nvt/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/nvt/type
                schema says required, Type is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/orphan
                schema says required, Orphan is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/owner
                schema says required, Owner is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/permissions
                schema says required, Permissions is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/result/description
                schema says required, Description is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/result/host
                schema says required, Host is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/result/host/asset
                schema says required, Asset is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/result/nvt
                schema says required, NVT is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/result/nvt/cvss_base
                schema says required, CVSSBase is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/result/nvt/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/result/nvt/type
                schema says required, Type is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/result/port
                schema says required, Port is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/result/qod
                schema says required, QOD is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/result/qod/type
                schema says required, Type is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/result/qod/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/result/severity
                schema says required, Severity is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/result/threat
                schema says required, Threat is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/severity
                schema says required, Severity is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/task/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/task/trash
                schema says required, Trash is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/text
                schema says required, Text is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/text/@excerpt
                schema says required, Excerpt is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/threat
                schema says required, Threat is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/user_tags/count
                schema says required, Count is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/user_tags/tag/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/user_tags/tag/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/user_tags/tag/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override/writable
                schema says required, Writable is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override_count
                schema says required, Count is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/overrides
                schema says required, Pagination is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/sort
                schema says required, Sort is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_overrides_response/override/active
                schema type boolean is text ("" is true), Go Active is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_overrides_response/override/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_overrides_response/override/orphan
                schema type boolean is text ("" is true), Go Orphan is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_overrides_response/override/result/nvt/cvss_base
                schema type integer, Go CVSSBase is string (want int)
  TYPE          response /get_overrides_response/override/result/severity
                schema type severity, Go Severity is string (want float64)
  TYPE          response /get_overrides_response/override/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_overrides_response/override/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_overrides_response/override/result/delta/notes/note/active
                duplicate-ele: redefined in the same scope

### get_permissions  (command_get_permissions.xml)

EXTRA=1  OMITEMPTY=29  TYPE=5  SCHEMA-DEFECT=2

  EXTRA         response /get_permissions_response/permission/subject/deleted
                Go field Deleted bool has no counterpart in the schema
  OMITEMPTY     response /get_permissions_response/filters
                schema says required, Filters is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/creation_time
                schema says required, CreationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/in_use
                schema says required, InUse is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/modification_time
                schema says required, ModificationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/owner
                schema says required, Owner is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/owner/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/permissions
                schema says required, Permissions is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/permissions/permission/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/resource
                schema says required, Resource is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/resource/@id
                schema says optional, ID string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_permissions_response/permission/resource/deleted
                schema says required, Deleted is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/resource/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/resource/trash
                schema says required, Trash is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/resource/type
                schema says required, Type is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/subject
                schema says required, Subject is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/subject/@id
                schema says optional, ID string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_permissions_response/permission/subject/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/subject/trash
                schema says required, Trash is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/subject/type
                schema says required, Type is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/user_tags/count
                schema says required, Count is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/user_tags/tag/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/user_tags/tag/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/user_tags/tag/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission/writable
                schema says required, Writable is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permission_count
                schema says required, Count is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/permissions
                schema says required, Pagination is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_permissions_response/sort
                schema says required, Sort is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_permissions_response/permission/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_permissions_response/permission/resource/deleted
                schema type boolean is text ("" is true), Go Deleted is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_permissions_response/permission/resource/trash
                schema type boolean is text ("" is true), Go Trash is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_permissions_response/permission/subject/trash
                schema type boolean is text ("" is true), Go Trash is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_permissions_response/permission/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_permissions_response/permission/resource/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_permissions_response/permission/subject/permissions
                empty-pattern: pattern declares no content

### get_port_lists  (command_get_port_lists.xml)

MISSING=2  EXTRA=1  CARDINALITY=1  KIND=3  OMITEMPTY=36  TYPE=4  SCHEMA-DEFECT=1

  MISSING       response /get_port_lists_response/filters/keywords/*
                schema declares 1 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_port_lists_response/sort/field/*
                schema declares 1 children; Go models this as scalar string
  EXTRA         response /get_port_lists_response/sort/order
                Go field Order string has no counterpart in the schema
  CARDINALITY   response /get_port_lists_response/port_list/deprecated
                schema says optional; Deprecated bool cannot express absence (zero is a legal value)
  KIND          response /get_port_lists_response/port_list/user_tags/count
                schema declares elem, Go models it as attr (Count int)
  KIND          response /get_port_lists_response/port_list_count/filtered
                schema declares elem, Go models it as attr (Filtered int)
  KIND          response /get_port_lists_response/port_list_count/page
                schema declares elem, Go models it as attr (Page int)
  OMITEMPTY     response /get_port_lists_response/filters
                schema says required, Filters is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/filters/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/filters/keywords
                schema says required, Keywords is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/filters/term
                schema says required, Term is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/creation_time
                schema says required, CreationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/in_use
                schema says required, InUse is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/modification_time
                schema says required, ModificationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/owner
                schema says required, Owner is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/owner/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/permissions
                schema says required, Permissions is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/permissions/permission/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/port_count
                schema says required, PortCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/port_count/all
                schema says required, All is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/port_count/tcp
                schema says required, TCP is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/port_count/udp
                schema says required, UDP is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/port_ranges/port_range/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/port_ranges/port_range/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/port_ranges/port_range/end
                schema says required, End is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/port_ranges/port_range/start
                schema says required, Start is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/port_ranges/port_range/type
                schema says required, Type is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/predefined
                schema says required, Predefined is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/targets/target/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/targets/target/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/user_tags/tag/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/user_tags/tag/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/user_tags/tag/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/user_tags/tag/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list/writable
                schema says required, Writable is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_list_count
                schema says required, Count is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_lists
                schema says required, Pagination is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_lists/@max
                schema says required, Max is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/port_lists/@start
                schema says required, Start is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/sort
                schema says required, Sort is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_port_lists_response/sort/field
                schema says required, Field is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_port_lists_response/port_list/deprecated
                schema type boolean is text ("" is true), Go Deprecated is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_port_lists_response/port_list/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_port_lists_response/port_list/predefined
                schema type boolean is text ("" is true), Go Predefined is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_port_lists_response/port_list/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_port_lists_response/port_list/targets/target/permissions
                empty-pattern: pattern declares no content

### get_preferences  (command_get_preferences.xml)

conformant

### get_report_configs  (command_get_report_configs.xml)

MISSING=11  EXTRA=2  KIND=2  OMITEMPTY=15  TYPE=2  SCHEMA-DEFECT=1

  MISSING       request  /get_report_configs/@details
                optional attr details is not modelled
  MISSING       request  /get_report_configs/@filt_id
                optional attr filt_id is not modelled
  MISSING       request  /get_report_configs/@filter
                optional attr filter is not modelled
  MISSING       request  /get_report_configs/@trash
                optional attr trash is not modelled
  MISSING       response /get_report_configs_response/filters/keywords/*
                schema declares 1 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_report_configs_response/report_config/orphan
                optional elem orphan is not modelled
  MISSING       response /get_report_configs_response/report_config/param/default/*
                schema declares 1 children; Go models this as scalar string
  MISSING       response /get_report_configs_response/report_config/param/type/*
                schema declares 2 children; Go models this as scalar string
  MISSING       response /get_report_configs_response/report_config/param/value/*
                schema declares 2 children; Go models this as scalar string
  MISSING       response /get_report_configs_response/report_config/user_tags
                optional elem user_tags is not modelled
  MISSING       response /get_report_configs_response/sort/field/*
                schema declares 1 children; Go models this as scalar string
  EXTRA         response /get_report_configs_response/report_config/param/@using_default
                Go field UsingDefault bool has no counterpart in the schema
  EXTRA         response /get_report_configs_response/sort/order
                Go field Order string has no counterpart in the schema
  KIND          response /get_report_configs_response/report_config_count/filtered
                schema declares elem, Go models it as attr (Filtered int)
  KIND          response /get_report_configs_response/report_config_count/page
                schema declares elem, Go models it as attr (Page int)
  OMITEMPTY     response /get_report_configs_response/filters
                schema says required, Filters is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_configs_response/filters/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_configs_response/filters/keywords
                schema says required, Keywords is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_configs_response/filters/term
                schema says required, Term is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_configs_response/report_config/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_configs_response/report_config/owner
                schema says required, Owner is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_configs_response/report_config/param/options
                schema says optional,  (wrapper) has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_report_configs_response/report_config/permissions
                schema says required, Permissions is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_configs_response/report_config/report_format/name
                schema says optional, Name string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_report_configs_response/report_config_count
                schema says required, Count is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_configs_response/report_configs
                schema says required, Pagination is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_configs_response/report_configs/@max
                schema says required, Max is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_configs_response/report_configs/@start
                schema says required, Start is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_configs_response/sort
                schema says required, Sort is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_configs_response/sort/field
                schema says required, Field is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_report_configs_response/report_config/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_report_configs_response/report_config/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_report_configs_response/report_config/report_format/permissions
                empty-pattern: pattern declares no content

### get_report_formats  (command_get_report_formats.xml)

MISSING=17  EXTRA=2  CARDINALITY=1  KIND=2  OMITEMPTY=14  TYPE=6

  MISSING       request  /get_report_formats/@alerts
                optional attr alerts is not modelled
  MISSING       request  /get_report_formats/@filt_id
                optional attr filt_id is not modelled
  MISSING       request  /get_report_formats/@filter
                optional attr filter is not modelled
  MISSING       request  /get_report_formats/@params
                optional attr params is not modelled
  MISSING       request  /get_report_formats/@report_configs
                optional attr report_configs is not modelled
  MISSING       request  /get_report_formats/@trash
                optional attr trash is not modelled
  MISSING       response /get_report_formats_response/filters/keywords/*
                schema declares 1 children; Go models this as scalar string  [fidelity]
  MISSING       response /get_report_formats_response/report_format/alerts
                optional elem alerts is not modelled
  MISSING       response /get_report_formats_response/report_format/invisible_alerts
                optional elem invisible_alerts is not modelled
  MISSING       response /get_report_formats_response/report_format/invisible_report_configs
                optional elem invisible_report_configs is not modelled
  MISSING       response /get_report_formats_response/report_format/owner
                required elem owner is not modelled  [fidelity]
  MISSING       response /get_report_formats_response/report_format/param
                repeated elem param is not modelled
  MISSING       response /get_report_formats_response/report_format/permissions
                required elem permissions is not modelled  [fidelity]
  MISSING       response /get_report_formats_response/report_format/report_configs
                optional elem report_configs is not modelled
  MISSING       response /get_report_formats_response/report_format/trust/@name
                required attr name is not modelled
  MISSING       response /get_report_formats_response/report_format/user_tags
                optional elem user_tags is not modelled
  MISSING       response /get_report_formats_response/sort/field/*
                schema declares 1 children; Go models this as scalar string
  EXTRA         response /get_report_formats_response/report_format/comment
                Go field Comment string has no counterpart in the schema
  EXTRA         response /get_report_formats_response/sort/order
                Go field Order string has no counterpart in the schema
  CARDINALITY   response /get_report_formats_response/report_format/deprecated
                schema says optional; Deprecated bool cannot express absence (zero is a legal value)
  KIND          response /get_report_formats_response/report_format_count/filtered
                schema declares elem, Go models it as attr (Filtered int)
  KIND          response /get_report_formats_response/report_format_count/page
                schema declares elem, Go models it as attr (Page int)
  OMITEMPTY     response /get_report_formats_response/filters
                schema says required, Filters is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_formats_response/filters/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_formats_response/filters/keywords
                schema says required, Keywords is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_formats_response/filters/term
                schema says required, Term is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_formats_response/report_format/configurable
                schema says required, Configurable is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_formats_response/report_format/predefined
                schema says required, Predefined is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_formats_response/report_format/report_type
                schema says required, ReportType is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_formats_response/report_format/trust
                schema says required, Trust is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_formats_response/report_format_count
                schema says required, Count is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_formats_response/report_formats
                schema says required, Pagination is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_formats_response/report_formats/@max
                schema says required, Max is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_formats_response/report_formats/@start
                schema says required, Start is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_formats_response/sort
                schema says required, Sort is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_report_formats_response/sort/field
                schema says required, Field is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_report_formats_response/report_format/active
                schema type boolean is text ("" is true), Go Active is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_report_formats_response/report_format/configurable
                schema type boolean is text ("" is true), Go Configurable is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_report_formats_response/report_format/deprecated
                schema type boolean is text ("" is true), Go Deprecated is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_report_formats_response/report_format/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_report_formats_response/report_format/predefined
                schema type boolean is text ("" is true), Go Predefined is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_report_formats_response/report_format/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'

### get_reports  (command_get_reports.xml)

MISSING=49  EXTRA=6  CARDINALITY=5  OMITEMPTY=143  TYPE=4  SCHEMA-DEFECT=10

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
                schema says required, Go has a slice []commands.ReportFormatParam (Params)
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

MISSING=1  EXTRA=4  CARDINALITY=1  OMITEMPTY=20  SCHEMA-DEFECT=2

  MISSING       response /get_results_response/result/nvt/severities/severity
                repeated elem severity is not modelled  [fidelity]
  EXTRA         response /get_results_response/result/delta/notes/note/result/*
                schema declares a "result" leaf; Go models 8 children
  EXTRA         response /get_results_response/result/delta/overrides/override/result/*
                schema declares a "result" leaf; Go models 8 children
  EXTRA         response /get_results_response/result/notes/note/*
                schema declares a "note" leaf; Go models 18 children
  EXTRA         response /get_results_response/result/overrides/override/*
                schema declares a "override" leaf; Go models 21 children
  CARDINALITY   response /get_results_response/result/delta/notes/note/severity
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
  OMITEMPTY     response /get_results_response/result/nvt/epss/max_epss/cve/@id
                schema says optional, ID string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/nvt/epss/max_severity/cve/@id
                schema says optional, ID string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/nvt/refs/ref/@id
                schema says optional, ID string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_results_response/result/nvt/refs/ref/@type
                schema says optional, Type string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  SCHEMA-DEFECT response /get_results_response/result/delta/notes/note/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_results_response/result/delta/overrides/override/active
                duplicate-ele: redefined in the same scope

### get_roles  (command_get_roles.xml)

MISSING=11  OMITEMPTY=6  TYPE=2

  MISSING       request  /get_roles/@filt_id
                optional attr filt_id is not modelled
  MISSING       request  /get_roles/@filter
                optional attr filter is not modelled
  MISSING       request  /get_roles/@role_id
                optional attr role_id is not modelled
  MISSING       request  /get_roles/@trash
                optional attr trash is not modelled
  MISSING       response /get_roles_response/filters
                required elem filters is not modelled  [fidelity]
  MISSING       response /get_roles_response/role/owner
                required elem owner is not modelled  [fidelity]
  MISSING       response /get_roles_response/role/permissions
                required elem permissions is not modelled  [fidelity]
  MISSING       response /get_roles_response/role/user_tags
                optional elem user_tags is not modelled
  MISSING       response /get_roles_response/role_count
                required elem role_count is not modelled  [fidelity]
  MISSING       response /get_roles_response/roles
                required elem roles is not modelled  [fidelity]
  MISSING       response /get_roles_response/sort
                required elem sort is not modelled  [fidelity]
  OMITEMPTY     response /get_roles_response/role/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_roles_response/role/creation_time
                schema says required, CreationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_roles_response/role/in_use
                schema says required, InUse is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_roles_response/role/modification_time
                schema says required, ModificationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_roles_response/role/users
                schema says required, Users is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_roles_response/role/writable
                schema says required, Writable is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_roles_response/role/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_roles_response/role/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'

### get_scanners  (command_get_scanners.xml)

EXTRA=5  OMITEMPTY=10  TYPE=3  SCHEMA-DEFECT=10

  EXTRA         response /get_scanners_response/scanner/certificate_info/*
                schema declares a "certificate_info" leaf; Go models 8 children
  EXTRA         response /get_scanners_response/scanner/credential/login
                Go field Login string has no counterpart in the schema
  EXTRA         response /get_scanners_response/scanner/key_pub
                Go field KeyPub string has no counterpart in the schema
  EXTRA         response /get_scanners_response/scanner/truncated
                Go field Truncated string has no counterpart in the schema
  EXTRA         response /get_scanners_response/truncate
                Go field Truncate string has no counterpart in the schema
  OMITEMPTY     response /get_scanners_response/filters/name
                schema says optional, Name string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_scanners_response/scanner/configs
                schema says optional, Configs commands.ScannerConfigs has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_scanners_response/scanner/copy
                schema says optional, Copy string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_scanners_response/scanner/info
                schema says optional, Info commands.ScannerInfo has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_scanners_response/scanner/info/daemon
                schema says optional, Daemon commands.ScannerInfoDaemon has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_scanners_response/scanner/info/params
                schema says optional, Params commands.ScannerInfoParams has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_scanners_response/scanner/info/protocol
                schema says optional, Protocol commands.ScannerInfoProtocol has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_scanners_response/scanner/info/scanner
                schema says optional, Scanner commands.ScannerInfoScanner has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_scanners_response/scanner/tasks
                schema says optional, Tasks commands.ScannerTasks has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_scanners_response/scanner/user_tags
                schema says optional, UserTags commands.ScannerUserTags has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  TYPE          response /get_scanners_response/scanner/credential/trash
                schema type boolean is text ("" is true), Go Trash is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_scanners_response/scanner/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_scanners_response/scanner/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
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

OMITEMPTY=17  TYPE=2  SCHEMA-DEFECT=1

  OMITEMPTY     response /get_schedules_response/filters
                schema says required, Filters is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_schedules_response/filters/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_schedules_response/filters/keywords
                schema says required, Keywords is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_schedules_response/filters/keywords/keyword/column
                schema says required, Column is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_schedules_response/filters/keywords/keyword/relation
                schema says required, Relation is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_schedules_response/filters/keywords/keyword/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_schedules_response/filters/term
                schema says required, Term is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_schedules_response/schedule/tasks
                schema says optional,  (wrapper) has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_schedules_response/schedule_count
                schema says required, ScheduleCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_schedules_response/schedule_count/filtered
                schema says required, Filtered is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_schedules_response/schedule_count/page
                schema says required, Page is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_schedules_response/schedules
                schema says required, ScheduleInfo is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_schedules_response/schedules/@max
                schema says required, Max is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_schedules_response/schedules/@start
                schema says required, Start is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_schedules_response/sort
                schema says required, Sort is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_schedules_response/sort/field
                schema says required, Field is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_schedules_response/sort/field/order
                schema says required, Order is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_schedules_response/schedule/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_schedules_response/schedule/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
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

MISSING=8  EXTRA=2  CARDINALITY=4  OMITEMPTY=2  TYPE=4  SCHEMA-DEFECT=1

  MISSING       request  /get_tags/@filt_id
                optional attr filt_id is not modelled
  MISSING       request  /get_tags/@names_only
                optional attr names_only is not modelled
  MISSING       request  /get_tags/@trash
                optional attr trash is not modelled
  MISSING       response /get_tags_response/filters
                required elem filters is not modelled  [fidelity]
  MISSING       response /get_tags_response/sort
                required elem sort is not modelled  [fidelity]
  MISSING       response /get_tags_response/tag/resources
                optional elem resources is not modelled  [fidelity]
  MISSING       response /get_tags_response/tag_count
                required elem tag_count is not modelled  [fidelity]
  MISSING       response /get_tags_response/tags
                required elem tags is not modelled  [fidelity]
  EXTRA         response /get_tags_response/tag/resource
                Go field Resources []commands.TagResource has no counterpart in the schema
  EXTRA         response /get_tags_response/truncated
                Go field Truncated string has no counterpart in the schema
  CARDINALITY   response /get_tags_response/tag/active
                schema says optional; Active bool cannot express absence (zero is a legal value)
  CARDINALITY   response /get_tags_response/tag/in_use
                schema says optional; InUse bool cannot express absence (zero is a legal value)
  CARDINALITY   response /get_tags_response/tag/orphan
                schema says optional; Orphan bool cannot express absence (zero is a legal value)
  CARDINALITY   response /get_tags_response/tag/writable
                schema says optional; Writable bool cannot express absence (zero is a legal value)
  OMITEMPTY     response /get_tags_response/tag/@id
                schema says optional, ID string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_tags_response/tag/permissions
                schema says optional,  (wrapper) has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  TYPE          response /get_tags_response/tag/active
                schema type boolean is text ("" is true), Go Active is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_tags_response/tag/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_tags_response/tag/orphan
                schema type boolean is text ("" is true), Go Orphan is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_tags_response/tag/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_tags_response/tag/resources/resource/permissions
                empty-pattern: pattern declares no content

### get_targets  (command_get_targets.xml)

MISSING=11  EXTRA=10  KIND=3  OMITEMPTY=55  TYPE=5  SCHEMA-DEFECT=8

  MISSING       response /get_targets_response/filters/keywords
                required elem keywords is not modelled  [fidelity]
  MISSING       response /get_targets_response/sort/field/*
                schema declares 1 children; Go models this as scalar string
  MISSING       response /get_targets_response/target/esxi_credential/trash
                required elem trash is not modelled  [fidelity]
  MISSING       response /get_targets_response/target/krb5_credential/trash
                required elem trash is not modelled  [fidelity]
  MISSING       response /get_targets_response/target/port_list/trash
                required elem trash is not modelled  [fidelity]
  MISSING       response /get_targets_response/target/smb_credential/trash
                required elem trash is not modelled  [fidelity]
  MISSING       response /get_targets_response/target/snmp_credential/trash
                required elem trash is not modelled  [fidelity]
  MISSING       response /get_targets_response/target/ssh_credential/trash
                required elem trash is not modelled  [fidelity]
  MISSING       response /get_targets_response/target/ssh_elevate_credential/trash
                required elem trash is not modelled  [fidelity]
  MISSING       response /get_targets_response/target/tasks/task/name
                required elem name is not modelled
  MISSING       response /get_targets_response/target/user_tags/tag/comment
                required elem comment is not modelled
  EXTRA         response /get_targets_response/filters/comment
                Go field Comment string has no counterpart in the schema
  EXTRA         response /get_targets_response/sort/order
                Go field Order string has no counterpart in the schema
  EXTRA         response /get_targets_response/target/esxi_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         response /get_targets_response/target/krb5_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         response /get_targets_response/target/port_range/*
                schema declares a "name" leaf; Go models 2 children
  EXTRA         response /get_targets_response/target/smb_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         response /get_targets_response/target/snmp_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         response /get_targets_response/target/ssh_elevate_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         response /get_targets_response/target/tasks/@count
                Go field Count int has no counterpart in the schema
  EXTRA         response /get_targets_response/truncated
                Go field Truncated string has no counterpart in the schema
  KIND          response /get_targets_response/target/user_tags/count
                schema declares elem, Go models it as attr (Count int)
  KIND          response /get_targets_response/target_count/filtered
                schema declares elem, Go models it as attr (Filtered int)
  KIND          response /get_targets_response/target_count/page
                schema declares elem, Go models it as attr (Page int)
  OMITEMPTY     response /get_targets_response/filters
                schema says required, Filters is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/filters/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/filters/term
                schema says required, Term is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/sort
                schema says required, Sort is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/sort/field
                schema says required, Field is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/alive_tests
                schema says required, AliveTests is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/allow_simultaneous_ips
                schema says required, AllowSimultaneousIPs is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/creation_time
                schema says required, CreationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/esxi_credential
                schema says required, ESXICredential is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/esxi_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/esxi_credential/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/exclude_hosts
                schema says required, ExcludeHosts is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/hosts
                schema says required, Hosts is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/in_use
                schema says required, InUse is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/krb5_credential
                schema says required, Krb5Credential is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/krb5_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/krb5_credential/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/max_hosts
                schema says required, MaxHosts is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/modification_time
                schema says required, ModificationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/owner
                schema says required, Owner is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/owner/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/permissions
                schema says required, Permissions is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/permissions
                schema says required, Permissions is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/permissions/permission/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/permissions/permission/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/port_list
                schema says required, PortList is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/port_list/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/port_list/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/port_range
                schema says required, PortRange is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/reverse_lookup_only
                schema says required, ReverseLookupOnly is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/reverse_lookup_unify
                schema says required, ReverseLookupUnify is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/smb_credential
                schema says required, SMBCredential is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/smb_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/smb_credential/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/snmp_credential
                schema says required, SNMPCredential is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/snmp_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/snmp_credential/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/ssh_credential
                schema says required, SSHCredential is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/ssh_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/ssh_credential/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/ssh_credential/port
                schema says required, Port is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/ssh_elevate_credential
                schema says required, SSHElevateCredential is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/ssh_elevate_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/ssh_elevate_credential/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/tasks/task/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/user_tags/tag/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/user_tags/tag/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/user_tags/tag/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target/writable
                schema says required, Writable is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/target_count
                schema says required, TargetCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/targets
                schema says required, Targets is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/targets/@max
                schema says required, Max is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_targets_response/targets/@start
                schema says required, Start is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_targets_response/target/allow_simultaneous_ips
                schema type boolean is text ("" is true), Go AllowSimultaneousIPs is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_targets_response/target/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_targets_response/target/reverse_lookup_only
                schema type boolean is text ("" is true), Go ReverseLookupOnly is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_targets_response/target/reverse_lookup_unify
                schema type boolean is text ("" is true), Go ReverseLookupUnify is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_targets_response/target/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_targets_response/target/esxi_credential/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_targets_response/target/krb5_credential/permissions
                empty-pattern: pattern declares no content
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

MISSING=10  EXTRA=15  CARDINALITY=7  OMITEMPTY=9  TYPE=9  SCHEMA-DEFECT=4

  MISSING       response /get_tasks_response/task/alert/trash
                required elem trash is not modelled
  MISSING       response /get_tasks_response/task/current_report/report
                required elem report is not modelled
  MISSING       response /get_tasks_response/task/last_report/report/compliance_count
                optional elem compliance_count is not modelled
  MISSING       response /get_tasks_response/task/last_report/report/result_count
                optional elem result_count is not modelled  [fidelity]
  MISSING       response /get_tasks_response/task/report_count/*
                schema declares 1 children; Go models this as scalar int  [fidelity]
  MISSING       response /get_tasks_response/task/schedule/icalendar
                optional elem icalendar is not modelled
  MISSING       response /get_tasks_response/task/schedule/icalendar
                optional elem icalendar is not modelled
  MISSING       response /get_tasks_response/task/schedule/timezone
                optional elem timezone is not modelled
  MISSING       response /get_tasks_response/task/schedule/timezone
                optional elem timezone is not modelled
  MISSING       response /get_tasks_response/task/user_tags
                optional elem user_tags is not modelled
  EXTRA         response /get_tasks_response/task/config/type
                Go field Type string has no counterpart in the schema
  EXTRA         response /get_tasks_response/task/current_report/@id
                Go field ID string has no counterpart in the schema
  EXTRA         response /get_tasks_response/task/current_report/scan_end
                Go field ScanEnd commands.GMPTime has no counterpart in the schema
  EXTRA         response /get_tasks_response/task/current_report/scan_run_status
                Go field ScanRunStatus string has no counterpart in the schema
  EXTRA         response /get_tasks_response/task/current_report/scan_start
                Go field ScanStart commands.GMPTime has no counterpart in the schema
  EXTRA         response /get_tasks_response/task/current_report/severity
                Go field Severity float64 has no counterpart in the schema
  EXTRA         response /get_tasks_response/task/current_report/timestamp
                Go field Timestamp commands.GMPTime has no counterpart in the schema
  EXTRA         response /get_tasks_response/task/last_report/report/scan_run_status
                Go field ScanRunStatus string has no counterpart in the schema
  EXTRA         response /get_tasks_response/task/last_report/report/scan_start
                Go field ScanStart commands.GMPTime has no counterpart in the schema
  EXTRA         response /get_tasks_response/task/observers/user_list
                Go field UserList string has no counterpart in the schema
  EXTRA         response /get_tasks_response/task/reports
                Go field Reports *commands.TaskReports has no counterpart in the schema
  EXTRA         response /get_tasks_response/task/scanner/trash
                Go field Trash bool has no counterpart in the schema
  EXTRA         response /get_tasks_response/task/schedule_periods/*
                schema declares a "integer" leaf; Go models 1 children
  EXTRA         response /get_tasks_response/task/schedule_periods/*
                schema declares a "integer" leaf; Go models 1 children
  EXTRA         response /get_tasks_response/truncated
                Go field Truncated string has no counterpart in the schema
  CARDINALITY   response /get_tasks_response/task/alterable
                schema says optional; Alterable bool cannot express absence (zero is a legal value)
  CARDINALITY   response /get_tasks_response/task/in_use
                schema says optional; InUse bool cannot express absence (zero is a legal value)
  CARDINALITY   response /get_tasks_response/task/last_report/report/severity
                schema says optional; Severity float64 cannot express absence (zero is a legal value)
  CARDINALITY   response /get_tasks_response/task/progress
                schema says optional; Progress int cannot express absence (zero is a legal value)
  CARDINALITY   response /get_tasks_response/task/report_count
                schema says optional; ReportCount int cannot express absence (zero is a legal value)
  CARDINALITY   response /get_tasks_response/task/result_count
                schema says optional; ResultCount int cannot express absence (zero is a legal value)
  CARDINALITY   response /get_tasks_response/task/writable
                schema says optional; Writable bool cannot express absence (zero is a legal value)
  OMITEMPTY     response /get_tasks_response/filters/name
                schema says optional, Name string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_tasks_response/task/config/trash
                schema says required, Trash is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tasks_response/task/last_report/report/scan_end
                schema says required, ScanEnd is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tasks_response/task/last_report/report/timestamp
                schema says required, Timestamp is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tasks_response/task/name
                schema says optional, Name string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_tasks_response/task/name
                schema says optional, Name string has no omitempty and is always emitted  [suppressed: inert on decode; shows up as 'added' in the fidelity corpus]
  OMITEMPTY     response /get_tasks_response/task/schedule/trash
                schema says required, Trash is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tasks_response/task/schedule/trash
                schema says required, Trash is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tasks_response/task/target/trash
                schema says required, Trash is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_tasks_response/task/alterable
                schema type boolean is text ("" is true), Go Alterable is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_tasks_response/task/config/trash
                schema type boolean is text ("" is true), Go Trash is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_tasks_response/task/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_tasks_response/task/result_count
                schema type text is text, Go ResultCount is int
  TYPE          response /get_tasks_response/task/scanner/type
                schema type integer, Go Type is string (want int)
  TYPE          response /get_tasks_response/task/schedule/trash
                schema type boolean is text ("" is true), Go Trash is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_tasks_response/task/schedule/trash
                schema type boolean is text ("" is true), Go Trash is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_tasks_response/task/target/trash
                schema type boolean is text ("" is true), Go Trash is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_tasks_response/task/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_tasks_response/task/alert/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_tasks_response/task/config/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_tasks_response/task/scanner/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_tasks_response/task/target/permissions
                empty-pattern: pattern declares no content

### get_tickets  (command_get_tickets.xml)

MISSING=5  EXTRA=4  KIND=3  OMITEMPTY=23  TYPE=3

  MISSING       request  /get_tickets/@filt_id
                optional attr filt_id is not modelled
  MISSING       request  /get_tickets/@trash
                optional attr trash is not modelled
  MISSING       response /get_tickets_response/filters/keywords
                required elem keywords is not modelled  [fidelity]
  MISSING       response /get_tickets_response/sort/field/*
                schema declares 1 children; Go models this as scalar string
  MISSING       response /get_tickets_response/ticket/user_tags/tag/comment
                required elem comment is not modelled
  EXTRA         response /get_tickets_response/filters/comment
                Go field Comment string has no counterpart in the schema
  EXTRA         response /get_tickets_response/sort/order
                Go field Order string has no counterpart in the schema
  EXTRA         response /get_tickets_response/ticket/nvt/name
                Go field Name string has no counterpart in the schema
  EXTRA         response /get_tickets_response/truncated
                Go field Truncated string has no counterpart in the schema
  KIND          response /get_tickets_response/ticket/user_tags/count
                schema declares elem, Go models it as attr (Count int)
  KIND          response /get_tickets_response/ticket_count/filtered
                schema declares elem, Go models it as attr (Filtered int)
  KIND          response /get_tickets_response/ticket_count/page
                schema declares elem, Go models it as attr (Page int)
  OMITEMPTY     response /get_tickets_response/filters
                schema says required, Filters is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/filters/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/filters/term
                schema says required, Term is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/sort
                schema says required, Sort is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/sort/field
                schema says required, Field is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/ticket/assigned_to
                schema says required, AssignedTo is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/ticket/host
                schema says required, Host is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/ticket/location
                schema says required, Location is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/ticket/nvt
                schema says required, NVT is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/ticket/open_note
                schema says required, OpenNote is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/ticket/open_time
                schema says required, OpenTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/ticket/orphan
                schema says required, Orphan is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/ticket/owner
                schema says required, Owner is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/ticket/permissions
                schema says required, Permissions is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/ticket/severity
                schema says required, Severity is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/ticket/solution_type
                schema says required, SolutionType is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/ticket/user_tags/tag/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/ticket/user_tags/tag/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/ticket/user_tags/tag/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/ticket_count
                schema says required, TicketCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/tickets
                schema says required, TicketsInfo is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/tickets/@max
                schema says required, Max is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tickets_response/tickets/@start
                schema says required, Start is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_tickets_response/ticket/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_tickets_response/ticket/orphan
                schema type boolean is text ("" is true), Go Orphan is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_tickets_response/ticket/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'

### get_tls_certificates  (command_get_tls_certificates.xml)

MISSING=5  EXTRA=4  KIND=2  OMITEMPTY=33  TYPE=3

  MISSING       request  /get_tls_certificates/@filt_id
                optional attr filt_id is not modelled
  MISSING       response /get_tls_certificates_response/tls_certificate/sources/source/location
                optional elem location is not modelled
  MISSING       response /get_tls_certificates_response/tls_certificate/sources/source/origin
                required elem origin is not modelled
  MISSING       response /get_tls_certificates_response/tls_certificate/sources/source/timestamp
                required elem timestamp is not modelled
  MISSING       response /get_tls_certificates_response/tls_certificate/sources/source/tls_versions
                required elem tls_versions is not modelled
  EXTRA         response /get_tls_certificates_response/filters/comment
                Go field Comment string has no counterpart in the schema
  EXTRA         response /get_tls_certificates_response/tls_certificate/sources/source/@type
                Go field Type string has no counterpart in the schema
  EXTRA         response /get_tls_certificates_response/tls_certificate/valid
                Go field Valid bool has no counterpart in the schema
  EXTRA         response /get_tls_certificates_response/truncated
                Go field Truncated string has no counterpart in the schema
  KIND          response /get_tls_certificates_response/tls_certificate_count/filtered
                schema declares elem, Go models it as attr (Filtered int)
  KIND          response /get_tls_certificates_response/tls_certificate_count/page
                schema declares elem, Go models it as attr (Page int)
  OMITEMPTY     response /get_tls_certificates_response/filters
                schema says required, Filters is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/filters/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/filters/keywords
                schema says required, Keywords is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/filters/keywords/keyword/column
                schema says required, Column is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/filters/keywords/keyword/relation
                schema says required, Relation is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/filters/keywords/keyword/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/filters/term
                schema says required, Term is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/sort
                schema says required, Sort is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/sort/field
                schema says required, Field is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/sort/field/order
                schema says required, Order is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/activation_time
                schema says required, ActivationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/certificate
                schema says required, Certificate is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/creation_time
                schema says required, CreationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/expiration_time
                schema says required, ExpirationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/in_use
                schema says required, InUse is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/issuer_dn
                schema says required, IssuerDN is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/last_seen
                schema says required, LastSeen is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/md5_fingerprint
                schema says required, MD5Fingerprint is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/modification_time
                schema says required, ModificationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/name
                schema says required, Name is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/owner
                schema says required, Owner is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/permissions
                schema says required, Permissions is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/serial
                schema says required, Serial is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/sha256_fingerprint
                schema says required, SHA256Fingerprint is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/subject_dn
                schema says required, SubjectDN is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/time_status
                schema says required, TimeStatus is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/trust
                schema says required, Trust is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate/writable
                schema says required, Writable is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate_count
                schema says required, TLSCertificateCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificates
                schema says required, TLSCertificatesInfo is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificates/@max
                schema says required, Max is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificates/@start
                schema says required, Start is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_tls_certificates_response/tls_certificate/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_tls_certificates_response/tls_certificate/trust
                schema type boolean is text ("" is true), Go Trust is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_tls_certificates_response/tls_certificate/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'

### get_users  (command_get_users.xml)

MISSING=1  EXTRA=2  CARDINALITY=2  KIND=2  OMITEMPTY=25  TYPE=2  SCHEMA-DEFECT=2

  MISSING       request  /get_users/@filt_id
                optional attr filt_id is not modelled
  EXTRA         response /get_users_response/filters/comment
                Go field Comment string has no counterpart in the schema
  EXTRA         response /get_users_response/truncated
                Go field Truncated string has no counterpart in the schema
  CARDINALITY   response /get_users_response/user/groups/group
                schema says required, Go has a slice []commands.UserGroup (Group)
  CARDINALITY   response /get_users_response/user/role
                schema says zero-or-more, Go has a single *commands.UserRole (Role)
  KIND          response /get_users_response/user_count/filtered
                schema declares elem, Go models it as attr (Filtered int)
  KIND          response /get_users_response/user_count/page
                schema declares elem, Go models it as attr (Page int)
  OMITEMPTY     response /get_users_response/filters
                schema says required, Filters is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/filters/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/filters/keywords
                schema says required, Keywords is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/filters/keywords/keyword/column
                schema says required, Column is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/filters/keywords/keyword/relation
                schema says required, Relation is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/filters/keywords/keyword/value
                schema says required, Value is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/filters/term
                schema says required, Term is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/sort
                schema says required, Sort is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/sort/field
                schema says required, Field is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/sort/field/order
                schema says required, Order is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/user/@id
                schema says required, ID is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/user/comment
                schema says required, Comment is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/user/creation_time
                schema says required, CreationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/user/groups
                schema says required, Groups is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/user/hosts
                schema says required, Hosts is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/user/in_use
                schema says required, InUse is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/user/modification_time
                schema says required, ModificationTime is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/user/owner
                schema says required, Owner is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/user/permissions
                schema says required, Permissions is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/user/sources
                schema says required, Sources is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/user/writable
                schema says required, Writable is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/user_count
                schema says required, UserCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/users
                schema says required, UsersInfo is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/users/@max
                schema says required, Max is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_users_response/users/@start
                schema says required, Start is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_users_response/user/in_use
                schema type boolean is text ("" is true), Go InUse is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          response /get_users_response/user/writable
                schema type boolean is text ("" is true), Go Writable is bool: '7' fails to decode and '1' re-emits as 'true'
  SCHEMA-DEFECT response /get_users_response/user/groups/group/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_users_response/user/role/permissions
                empty-pattern: pattern declares no content

### get_version  (command_get_version.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /get_version
                empty-pattern: pattern declares no content

### get_vulns  (command_get_vulns.xml)

EXTRA=3  OMITEMPTY=8

  EXTRA         request  /get_vulns/@trash
                Go field Trash *commands.GMPBool has no counterpart in the schema
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

MISSING=1  SCHEMA-DEFECT=6

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

EXTRA=6  CARDINALITY=3  OMITEMPTY=6  TYPE=3

  EXTRA         request  /modify_config/preference/alt
                Go field Alt []string has no counterpart in the schema
  EXTRA         request  /modify_config/preference/default
                Go field Default string has no counterpart in the schema
  EXTRA         request  /modify_config/preference/hr_name
                Go field HRName string has no counterpart in the schema
  EXTRA         request  /modify_config/preference/id
                Go field ID string has no counterpart in the schema
  EXTRA         request  /modify_config/preference/nvt/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_config/preference/type
                Go field Type string has no counterpart in the schema
  CARDINALITY   request  /modify_config/family_selection/family/all
                schema says optional; All bool cannot express absence (zero is a legal value)
  CARDINALITY   request  /modify_config/family_selection/family/growing
                schema says optional; Growing bool cannot express absence (zero is a legal value)
  CARDINALITY   request  /modify_config/family_selection/growing
                schema says optional; Growing bool cannot express absence (zero is a legal value)
  OMITEMPTY     request  /modify_config/@config_id
                schema says required, ConfigID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_config/family_selection/family/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_config/nvt_selection/family
                schema says required, Family is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_config/nvt_selection/nvt/@oid
                schema says required, OID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_config/preference/nvt
                schema says optional, NVT commands.ConfigPreferenceNVT has no omitempty and is always emitted
  OMITEMPTY     request  /modify_config/preference/value
                schema says optional, Value string has no omitempty and is always emitted
  TYPE          request  /modify_config/family_selection/family/all
                schema type boolean_atoi is text (atoi first, so "" is false), Go All is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          request  /modify_config/family_selection/family/growing
                schema type boolean_atoi is text (atoi first, so "" is false), Go Growing is bool: '7' fails to decode and '1' re-emits as 'true'
  TYPE          request  /modify_config/family_selection/growing
                schema type boolean_atoi is text (atoi first, so "" is false), Go Growing is bool: '7' fails to decode and '1' re-emits as 'true'

### modify_credential  (command_modify_credential.xml)

OMITEMPTY=2

  OMITEMPTY     request  /modify_credential/privacy/algorithm
                schema says required, Algorithm is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_credential/privacy/password
                schema says required, Password is tagged omitempty  [suppressed: request-side over-constraint]

### modify_filter  (command_modify_filter.xml)

conformant

### modify_group  (command_modify_group.xml)

conformant

### modify_license  (command_modify_license.xml)

conformant

### modify_note  (command_modify_note.xml)

conformant

### modify_override  (command_modify_override.xml)

MISSING=2  EXTRA=2

  MISSING       request  /modify_override/result/*
                schema declares 1 children; Go models this as scalar string
  MISSING       request  /modify_override/task/*
                schema declares 1 children; Go models this as scalar string
  EXTRA         request  /modify_override/nvt/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_override/nvt/type
                Go field Type string has no counterpart in the schema

### modify_permission  (command_modify_permission.xml)

EXTRA=8  OMITEMPTY=4

  EXTRA         request  /modify_permission/resource/deleted
                Go field Deleted bool has no counterpart in the schema
  EXTRA         request  /modify_permission/resource/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_permission/resource/permissions
                Go field Permissions *commands.PermissionList has no counterpart in the schema
  EXTRA         request  /modify_permission/resource/trash
                Go field Trash bool has no counterpart in the schema
  EXTRA         request  /modify_permission/subject/deleted
                Go field Deleted bool has no counterpart in the schema
  EXTRA         request  /modify_permission/subject/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_permission/subject/permissions
                Go field Permissions *commands.PermissionList has no counterpart in the schema
  EXTRA         request  /modify_permission/subject/trash
                Go field Trash bool has no counterpart in the schema
  OMITEMPTY     request  /modify_permission/resource/@id
                schema says optional, ID string has no omitempty and is always emitted
  OMITEMPTY     request  /modify_permission/resource/type
                schema says required, Type is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_permission/subject/@id
                schema says optional, ID string has no omitempty and is always emitted
  OMITEMPTY     request  /modify_permission/subject/type
                schema says required, Type is tagged omitempty  [suppressed: request-side over-constraint]

### modify_port_list  (command_modify_port_list.xml)

conformant

### modify_report_config  (command_modify_report_config.xml)

MISSING=1  EXTRA=3

  MISSING       request  /modify_report_config/param/value/*
                schema declares 1 children; Go models this as scalar string
  EXTRA         request  /modify_report_config/param/@using_default
                Go field UsingDefault bool has no counterpart in the schema
  EXTRA         request  /modify_report_config/param/default
                Go field Default string has no counterpart in the schema
  EXTRA         request  /modify_report_config/param/type
                Go field Type string has no counterpart in the schema

### modify_report_format  (command_modify_report_format.xml)

OMITEMPTY=1

  OMITEMPTY     request  /modify_report_format/param/name
                schema says required, Name is tagged omitempty  [suppressed: request-side over-constraint]

### modify_role  (command_modify_role.xml)

conformant

### modify_scanner  (command_modify_scanner.xml)

EXTRA=3  OMITEMPTY=3  TYPE=2

  EXTRA         request  /modify_scanner/credential/login
                Go field Login string has no counterpart in the schema
  EXTRA         request  /modify_scanner/credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_scanner/credential/trash
                Go field Trash bool has no counterpart in the schema
  OMITEMPTY     request  /modify_scanner/host
                schema says required, Host is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_scanner/port
                schema says required, Port is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_scanner/type
                schema says required, Type is tagged omitempty  [suppressed: request-side over-constraint]
  TYPE          request  /modify_scanner/port
                schema type text is text, Go Port is int
  TYPE          request  /modify_scanner/relay_port
                schema type text is text, Go RelayPort is int

### modify_schedule  (command_modify_schedule.xml)

conformant

### modify_setting  (command_modify_setting.xml)

conformant

### modify_tag  (command_modify_tag.xml)

MISSING=1  EXTRA=3  OMITEMPTY=1

  MISSING       request  /modify_tag/resources/@action
                optional attr action is not modelled
  EXTRA         request  /modify_tag/resources/resource/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_tag/resources/resource/trash
                Go field Trash bool has no counterpart in the schema
  EXTRA         request  /modify_tag/resources/resource/type
                Go field Type string has no counterpart in the schema
  OMITEMPTY     request  /modify_tag/resources/resource/@id
                schema says optional, ID string has no omitempty and is always emitted

### modify_target  (command_modify_target.xml)

EXTRA=19  OMITEMPTY=10

  EXTRA         request  /modify_target/esxi_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_target/esxi_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         request  /modify_target/esxi_lsc_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_target/esxi_lsc_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         request  /modify_target/krb5_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_target/krb5_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         request  /modify_target/port_list/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_target/smb_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_target/smb_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         request  /modify_target/smb_lsc_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_target/smb_lsc_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         request  /modify_target/snmp_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_target/snmp_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         request  /modify_target/ssh_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_target/ssh_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         request  /modify_target/ssh_elevate_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_target/ssh_elevate_credential/port
                Go field Port string has no counterpart in the schema
  EXTRA         request  /modify_target/ssh_lsc_credential/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_target/ssh_lsc_credential/port
                Go field Port string has no counterpart in the schema
  OMITEMPTY     request  /modify_target/esxi_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_target/esxi_lsc_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_target/krb5_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_target/port_list/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_target/smb_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_target/smb_lsc_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_target/snmp_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_target/ssh_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_target/ssh_elevate_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]
  OMITEMPTY     request  /modify_target/ssh_lsc_credential/@id
                schema says required, ID is tagged omitempty  [suppressed: request-side over-constraint]

### modify_task  (command_modify_task.xml)

MISSING=1  EXTRA=12

  MISSING       request  /modify_task/file
                optional elem file is not modelled
  EXTRA         request  /modify_task/alert/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_task/config/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_task/config/trash
                Go field Trash bool has no counterpart in the schema
  EXTRA         request  /modify_task/config/type
                Go field Type string has no counterpart in the schema
  EXTRA         request  /modify_task/preferences/preference/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_task/scanner/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_task/scanner/trash
                Go field Trash bool has no counterpart in the schema
  EXTRA         request  /modify_task/scanner/type
                Go field Type string has no counterpart in the schema
  EXTRA         request  /modify_task/schedule/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_task/schedule/trash
                Go field Trash bool has no counterpart in the schema
  EXTRA         request  /modify_task/target/name
                Go field Name string has no counterpart in the schema
  EXTRA         request  /modify_task/target/trash
                Go field Trash bool has no counterpart in the schema

### modify_ticket  (command_modify_ticket.xml)

conformant

### modify_tls_certificate  (command_modify_tls_certificate.xml)

conformant

### modify_user  (command_modify_user.xml)

EXTRA=1  CARDINALITY=1  OMITEMPTY=1

  EXTRA         request  /modify_user/role/name
                Go field Name string has no counterpart in the schema
  CARDINALITY   request  /modify_user/sources/source
                schema says optional, Go has a slice []string (Source)
  OMITEMPTY     request  /modify_user/hosts/@allow
                schema says optional, Allow string has no omitempty and is always emitted

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

  OMITEMPTY     request  /start_task/@task_id
                schema says required, TaskID is tagged omitempty  [suppressed: request-side over-constraint]

### stop_task  (command_stop_task.xml)

OMITEMPTY=1

  OMITEMPTY     request  /stop_task/@task_id
                schema says required, TaskID is tagged omitempty  [suppressed: request-side over-constraint]

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
                Go field Version string has no counterpart in the schema

