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

commands compared:            119
fully conformant:             100
with actionable findings:     19
actionable findings:          27  (request 7, response 20)
suppressed findings:          143
confirmed by fidelity corpus: 2
in the schema, not in Go:     3  (delete_agents, get_agents, modify_agents)
in Go, not in the schema:     0

### Findings by category

                 actionable suppressed
  MISSING                 5         25
  EXTRA                  19         72
  CARDINALITY             0          2
  KIND                    0          0
  OMITEMPTY               0         44
  TYPE                    3          0
  SCHEMA-DEFECT         105          0

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
  undefined-type       3 sites in 3 commands
  unsupported-r        3 sites in 3 commands

### Commands by actionable findings

  create_report                 3   EXTRA=3
  get_reports                   3   EXTRA=3
  get_credentials               2   EXTRA=2
  get_info                      2   EXTRA=1  TYPE=1
  get_scanners                  2   EXTRA=2
  modify_target                 2   MISSING=1  EXTRA=1
  create_target                 1   EXTRA=1
  create_task                   1   MISSING=1
  get_features                  1   EXTRA=1
  get_notes                     1   EXTRA=1
  get_nvts                      1   EXTRA=1
  get_overrides                 1   EXTRA=1
  get_report_configs            1   MISSING=1
  get_resource_names            1   TYPE=1
  get_results                   1   EXTRA=1
  get_system_reports            1   TYPE=1
  get_targets                   1   EXTRA=1
  get_tasks                     1   MISSING=1
  help                          1   MISSING=1

### Confirmed by both oracles

Paths the schema says exist and that a captured document was also seen to
lose. Fix these first: neither oracle is guessing.

Response side only: the fidelity corpus keeps per-path DROPPED lists for
responses and only aggregate counts for requests, so a request-side finding can
never carry this mark. That is a gap in the cross-reference, not evidence that
requests are in better shape — get_notes models none of its seven request
attributes, and one of its own captured request fixtures uses three of them.

  get_report_configs       /get_report_configs_response/report_config/param/value/@using_default
  get_tasks                /get_tasks_response/task/hosts_ordering

## Per command

### authenticate  (command_authenticate.xml)

conformant

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

EXTRA=3  OMITEMPTY=3  SCHEMA-DEFECT=12

  EXTRA         request  /create_report/report/report/cves
                Go field CVEs *commands.ReportCVEs has no counterpart in the schema
  EXTRA         request  /create_report/report/report/result_count/#text
                Go *commands.ReportResultCount has a `,chardata` field; schema declares no character data here
  EXTRA         request  /create_report/report/report/results/result/oci_image
                Go field OCIImage *commands.ReportOCIImage has no counterpart in the schema
  OMITEMPTY     request  /create_report/report/report/hosts/@max
                schema says required, Max is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]
  OMITEMPTY     request  /create_report/report/report/hosts/@start
                schema says required, Start is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]
  OMITEMPTY     request  /create_report/report/report/results/result/nvt/refs
                schema says required, Refs is tagged omitempty  [suppressed: request-side over-constraint]
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

CARDINALITY=1  OMITEMPTY=1  SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /describe_auth
                empty-pattern: pattern declares no content
  CARDINALITY   response /describe_auth_response/group
                schema says required, Go has a slice []commands.AuthGroup (Group)  [suppressed: T-9 auth review: command_describe_auth.xml's response <pattern> declares <e>group</e> (exactly one), but internal/gmptest/testdata/live/describe_auth/001_response.xml confirms gvmd returns three sibling <group> elements (method:file, method:ldap_connect, method:radius_connect) in one unfiltered describe_auth_response]
  OMITEMPTY     response /describe_auth_response/group
                schema says required, Group is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### empty_trashcan  (command_empty_trashcan.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /empty_trashcan
                empty-pattern: pattern declares no content

### get_aggregates  (command_get_aggregates.xml)

EXTRA=7

  EXTRA         response /get_aggregates_response/aggregate/group/c_sum
                Go field CSum string has no counterpart in the schema  [suppressed: T-9 aggregate review: see the min entry above; same flat <c_sum> shape confirmed by the same evidence]
  EXTRA         response /get_aggregates_response/aggregate/group/max
                Go field Max string has no counterpart in the schema  [suppressed: T-9 aggregate review: see the min entry above; same flat <max> shape confirmed by the same evidence]
  EXTRA         response /get_aggregates_response/aggregate/group/mean
                Go field Mean string has no counterpart in the schema  [suppressed: T-9 aggregate review: see the min entry above; same flat <mean> shape confirmed by the same evidence]
  EXTRA         response /get_aggregates_response/aggregate/group/min
                Go field Min string has no counterpart in the schema  [suppressed: T-9 aggregate review: command_get_aggregates.xml's <group> pattern declares only value/subgroup/count/c_count/stats/text — per-column stats wrapped in <stats column="..."> — but the command's own <example> and internal/gmptest/testdata/upstream/get_aggregates/001_response.xml both show a flat <min> directly under <group> for the single-data_column case, the same pattern-vs-example disagreement as verify_scanner's version]
  EXTRA         response /get_aggregates_response/aggregate/group/sum
                Go field Sum string has no counterpart in the schema  [suppressed: T-9 aggregate review: see the min entry above; same flat <sum> shape confirmed by the same evidence]
  EXTRA         response /get_aggregates_response/aggregate/group/text/@column
                Go field Column string has no counterpart in the schema  [suppressed: T-9 aggregate review: command_get_aggregates.xml's schema declares <text>'s lone attribute as "name", but internal/gmptest/testdata/live/get_aggregates/004_response.xml sends a real gvmd 26.36.1 <text column="name"> — a second, undeclared attribute the schema never mentions at all, kept alongside the schema-declared (and seemingly unused) Name field since there is no evidence the schema's own declaration is categorically wrong rather than an alternate/legacy form. AggregateText is one structurally-deduplicated type shared by both aggregate/group/text and aggregate/overall/text, so both must carry this same entry]
  EXTRA         response /get_aggregates_response/aggregate/overall/text/@column
                Go field Column string has no counterpart in the schema  [suppressed: T-9 aggregate review: command_get_aggregates.xml's schema declares <text>'s lone attribute as "name", but internal/gmptest/testdata/live/get_aggregates/004_response.xml sends a real gvmd 26.36.1 <text column="name"> — a second, undeclared attribute the schema never mentions at all, kept alongside the schema-declared (and seemingly unused) Name field since there is no evidence the schema's own declaration is categorically wrong rather than an alternate/legacy form. AggregateText is one structurally-deduplicated type shared by both aggregate/group/text and aggregate/overall/text, so both must carry this same entry]

### get_alerts  (command_get_alerts.xml)

MISSING=1  EXTRA=2  OMITEMPTY=1  SCHEMA-DEFECT=2

  MISSING       response /get_alerts_response/sort/#text
                schema also declares character data (text) here; Go commands.AlertSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_alerts_response/alert_count/#text
                Go *commands.AlertCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_alerts_response/sort/field/#text
                Go commands.AlertField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_alerts_response/alert_count
                schema says required, AlertCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  SCHEMA-DEFECT response /get_alerts_response/alert/filter/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_alerts_response/alert/tasks/task/permissions
                empty-pattern: pattern declares no content

### get_assets  (command_get_assets.xml)

MISSING=1  EXTRA=2  OMITEMPTY=1  SCHEMA-DEFECT=1

  MISSING       response /get_assets_response/sort/#text
                schema also declares character data (text) here; Go commands.AssetSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_assets_response/asset_count/#text
                Go *commands.AssetCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_assets_response/sort/field/#text
                Go commands.AssetField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_assets_response/asset_count
                schema says required, AssetCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  SCHEMA-DEFECT response /get_assets_response/asset/host/routes/route/host/@distance
                undefined-type: number is undefined; aliased to integer

### get_configs  (command_get_configs.xml)

MISSING=1  EXTRA=2  OMITEMPTY=1  SCHEMA-DEFECT=1

  MISSING       response /get_configs_response/sort/#text
                schema also declares character data (text) here; Go commands.ConfigSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_configs_response/config_count/#text
                Go *commands.ConfigCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_configs_response/sort/field/#text
                Go commands.ConfigField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_configs_response/config_count
                schema says required, ConfigCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  SCHEMA-DEFECT response /get_configs_response/config/tasks/task/permissions
                empty-pattern: pattern declares no content

### get_credentials  (command_get_credentials.xml)

MISSING=1  EXTRA=4  OMITEMPTY=1  SCHEMA-DEFECT=2

  MISSING       response /get_credentials_response/sort/#text
                schema also declares character data (text) here; Go commands.CredentialSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_credentials_response/credential/private_key_info
                Go field PrivateKeyInfo *commands.CredentialPrivateKeyInfo has no counterpart in the schema
  EXTRA         response /get_credentials_response/credential/public_key_info
                Go field PublicKeyInfo *commands.CredentialPublicKeyInfo has no counterpart in the schema
  EXTRA         response /get_credentials_response/credential_count/#text
                Go *commands.CredentialCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_credentials_response/sort/field/#text
                Go commands.CredentialField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_credentials_response/credential_count
                schema says required, CredentialCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  SCHEMA-DEFECT response /get_credentials_response/credential/scanners/scanner/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_credentials_response/credential/targets/target/permissions
                empty-pattern: pattern declares no content

### get_features  (command_get_features.xml)

EXTRA=1  SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /get_features
                empty-pattern: pattern declares no content
  EXTRA         response /get_features_response/feature/@compiled_in
                Go field CompiledIn commands.GMPBool has no counterpart in the schema

### get_feeds  (command_get_feeds.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT response /get_feeds_response/feed/sync_not_available/error
                no-pattern: <ele> has neither <pattern> nor <type>

### get_filters  (command_get_filters.xml)

MISSING=1  EXTRA=2  OMITEMPTY=1  SCHEMA-DEFECT=1

  MISSING       response /get_filters_response/sort/#text
                schema also declares character data (text) here; Go commands.FilterSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_filters_response/filter_count/#text
                Go *commands.FilterCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_filters_response/sort/field/#text
                Go commands.FilterField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_filters_response/filter_count
                schema says required, FilterCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  SCHEMA-DEFECT response /get_filters_response/filter/alerts/alert/permissions
                empty-pattern: pattern declares no content

### get_groups  (command_get_groups.xml)

MISSING=1  EXTRA=2  OMITEMPTY=1

  MISSING       response /get_groups_response/sort/#text
                schema also declares character data (text) here; Go commands.GroupSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_groups_response/group_count/#text
                Go *commands.GroupCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_groups_response/sort/field/#text
                Go commands.GroupField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_groups_response/group_count
                schema says required, GroupCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### get_info  (command_get_info.xml)

MISSING=1  EXTRA=4  OMITEMPTY=1  TYPE=1  SCHEMA-DEFECT=1

  MISSING       response /get_info_response/sort/#text
                schema also declares character data (text) here; Go commands.NVTSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_info_response/info
                Go field Info2 *commands.GetInfoInfo has no counterpart in the schema  [suppressed: T-26: get_info's pagination wrapper (<info start max>) shares its element name with the data item but is schema-undeclared — see docs/review-2026-08.md's T-26 entry]
  EXTRA         response /get_info_response/info/web_application_vt
                Go field WebApplicationVT *commands.NVTWebApplicationVT has no counterpart in the schema
  EXTRA         response /get_info_response/info_count/#text
                Go *commands.NVTInfoCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_info_response/sort/field/#text
                Go commands.NVTField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_info_response/info_count
                schema says required, InfoCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  TYPE          response /get_info_response/info/@id
                schema type integer, Go ID is string (want int)
  SCHEMA-DEFECT response /get_info_response/info/cve/configuration_nodes/node/node
                undefined-type: node is undefined and has no alias

### get_license  (command_get_license.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /get_license
                empty-pattern: pattern declares no content

### get_notes  (command_get_notes.xml)

MISSING=1  EXTRA=3  OMITEMPTY=2  SCHEMA-DEFECT=3

  MISSING       response /get_notes_response/sort/#text
                schema also declares character data (text) here; Go commands.NoteSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_notes_response/note/result/oci_image
                Go field OCIImage *commands.NoteOCIImage has no counterpart in the schema
  EXTRA         response /get_notes_response/note_count/#text
                Go *commands.NoteCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_notes_response/sort/field/#text
                Go commands.NoteField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_notes_response/note/result/nvt/refs
                schema says required, Refs is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_notes_response/note_count
                schema says required, NoteCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  SCHEMA-DEFECT response /get_notes_response/note/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_notes_response/note/result/delta/overrides/override/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_notes_response/note/result/overrides/override/active
                duplicate-ele: redefined in the same scope

### get_nvt_families  (command_get_nvt_families.xml)

conformant

### get_nvts  (command_get_nvts.xml)

EXTRA=1  CARDINALITY=1  OMITEMPTY=2

  OMITEMPTY     request  /get_nvts/@nvt_oid
                schema says required, NVTOID is tagged omitempty  [suppressed: request-side over-constraint]
  EXTRA         response /get_nvts_response/nvt/discovery
                Go field Discovery *int has no counterpart in the schema
  CARDINALITY   response /get_nvts_response/nvt
                schema says required, Go has a slice []commands.GetNVTsNVT (NVT)  [suppressed: T-9 nvt review: command_get_nvts.xml's response <pattern> declares <e>nvt</e> (exactly one), contradicting its own <summary>Get one or many NVTs</summary> and its own <example> (two sibling <nvt> elements in one get_nvts_response with no nvt_oid filter, lines 4121-4159) — internal/gmptest/testdata/upstream/get_nvts/001_response.xml confirms it live, returning three sibling <nvt> elements for one unfiltered request]
  OMITEMPTY     response /get_nvts_response/nvt
                schema says required, NVT is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### get_overrides  (command_get_overrides.xml)

MISSING=1  EXTRA=3  OMITEMPTY=2  SCHEMA-DEFECT=3

  MISSING       response /get_overrides_response/sort/#text
                schema also declares character data (text) here; Go commands.OverrideSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_overrides_response/override/result/oci_image
                Go field OCIImage *commands.OverrideOCIImage has no counterpart in the schema
  EXTRA         response /get_overrides_response/override_count/#text
                Go *commands.OverrideCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_overrides_response/sort/field/#text
                Go commands.OverrideField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_overrides_response/override/result/nvt/refs
                schema says required, Refs is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  OMITEMPTY     response /get_overrides_response/override_count
                schema says required, OverrideCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  SCHEMA-DEFECT response /get_overrides_response/override/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_overrides_response/override/result/delta/notes/note/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_overrides_response/override/result/notes/note/active
                duplicate-ele: redefined in the same scope

### get_permissions  (command_get_permissions.xml)

MISSING=1  EXTRA=2  OMITEMPTY=1  SCHEMA-DEFECT=2

  MISSING       response /get_permissions_response/sort/#text
                schema also declares character data (text) here; Go commands.PermissionSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_permissions_response/permission_count/#text
                Go *commands.PermissionCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_permissions_response/sort/field/#text
                Go commands.PermissionField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_permissions_response/permission_count
                schema says required, PermissionCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  SCHEMA-DEFECT response /get_permissions_response/permission/resource/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_permissions_response/permission/subject/permissions
                empty-pattern: pattern declares no content

### get_port_lists  (command_get_port_lists.xml)

MISSING=1  EXTRA=2  OMITEMPTY=1  SCHEMA-DEFECT=1

  MISSING       response /get_port_lists_response/sort/#text
                schema also declares character data (text) here; Go commands.PortListSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_port_lists_response/port_list_count/#text
                Go *commands.PortListCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_port_lists_response/sort/field/#text
                Go commands.PortListField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_port_lists_response/port_list_count
                schema says required, PortListCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  SCHEMA-DEFECT response /get_port_lists_response/port_list/targets/target/permissions
                empty-pattern: pattern declares no content

### get_preferences  (command_get_preferences.xml)

conformant

### get_report_configs  (command_get_report_configs.xml)

MISSING=2  EXTRA=3  OMITEMPTY=1  SCHEMA-DEFECT=1

  MISSING       response /get_report_configs_response/report_config/param/value/@using_default
                required attr using_default is not modelled  [fidelity]
  MISSING       response /get_report_configs_response/sort/#text
                schema also declares character data (text) here; Go commands.ReportConfigSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_report_configs_response/report_config/param/@using_default
                Go field UsingDefault commands.GMPBool has no counterpart in the schema  [suppressed: T-9 report_config review: real position of using_default that FieldSuppressions["get_report_configs:/report_config/param/value"] removes — command_get_report_configs.xml's own <example> and internal/gmptest/testdata/upstream/get_report_configs/001_{response,fidelity}.txt both place the attribute on <param>, not <value>]
  EXTRA         response /get_report_configs_response/report_config_count/#text
                Go *commands.ReportConfigCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_report_configs_response/sort/field/#text
                Go commands.ReportConfigField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_report_configs_response/report_config_count
                schema says required, ReportConfigCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  SCHEMA-DEFECT response /get_report_configs_response/report_config/report_format/permissions
                empty-pattern: pattern declares no content

### get_report_formats  (command_get_report_formats.xml)

MISSING=1  EXTRA=3  OMITEMPTY=1

  MISSING       response /get_report_formats_response/sort/#text
                schema also declares character data (text) here; Go commands.ReportFormatSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_report_formats_response/report_format/comment
                Go field Comment string has no counterpart in the schema  [suppressed: T-9 report_format review: schema's report_format pattern never declares <comment>, but gvmd emits one on every report format — x5 in internal/gmptest/testdata/live/get_report_formats/001_response.xml, x6 in 006_response.xml (after a report format copy was added to the fixture set), and x1 upstream, all recorded as DROPPED in the matching _fidelity.txt. Every captured occurrence is empty, so it is the element's presence that is evidenced, not its content; hence no omitempty, which would re-marshal to the absence gvmd never sends]
  EXTRA         response /get_report_formats_response/report_format_count/#text
                Go *commands.ReportFormatCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_report_formats_response/sort/field/#text
                Go commands.ReportFormatField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_report_formats_response/report_format_count
                schema says required, ReportFormatCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### get_reports  (command_get_reports.xml)

MISSING=1  EXTRA=5  OMITEMPTY=3  SCHEMA-DEFECT=12

  MISSING       response /get_reports_response/sort/#text
                schema also declares character data (text) here; Go *commands.GetReportsSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_reports_response/report/report/cves
                Go field CVEs *commands.ReportCVEs has no counterpart in the schema
  EXTRA         response /get_reports_response/report/report/result_count/#text
                Go *commands.ReportResultCount has a `,chardata` field; schema declares no character data here
  EXTRA         response /get_reports_response/report/report/results/result/oci_image
                Go field OCIImage *commands.ReportOCIImage has no counterpart in the schema
  EXTRA         response /get_reports_response/report_count/#text
                Go *commands.ReportCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_reports_response/sort/field/#text
                Go commands.ReportSortField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_reports_response/report/report/hosts/@max
                schema says required, Max is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]
  OMITEMPTY     response /get_reports_response/report/report/hosts/@start
                schema says required, Start is tagged omitempty  [suppressed: inside <or>; only one branch is ever present]
  OMITEMPTY     response /get_reports_response/report/report/results/result/nvt/refs
                schema says required, Refs is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
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

TYPE=1  SCHEMA-DEFECT=1

  TYPE          response /get_resource_names_response/resource/@id
                schema type integer, Go ID is string (want int)
  SCHEMA-DEFECT response /get_resource_names_response/type
                undefined-type: type is undefined and has no alias

### get_results  (command_get_results.xml)

MISSING=1  EXTRA=4  OMITEMPTY=1  SCHEMA-DEFECT=4

  MISSING       response /get_results_response/sort/#text
                schema also declares character data (text) here; Go commands.ResultSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_results_response/result/nvt/solution
                Go field Solution *commands.ResultSolution has no counterpart in the schema  [suppressed: T-9 result review, finding 1: internal/codegen/testdata/schema/get_results.json's result/nvt node has no solution child, but all 4 live fixtures (internal/gmptest/testdata/live/get_results/00{1,2,3,4}_response.xml) send <solution type=''> right after </tags>, retained as DROPPED in every one of those fixtures' own _fidelity.txt. The full shape — optional element, chardata text plus optional type/method attributes — is the exact declaration get_nvts' OWN live schema snapshot (internal/codegen/testdata/schema/get_nvts.json's get_nvts_response/nvt/solution) already carries for the identical nvt/solution reuse, already proven correct as commands.NVTSolution (commands/get_nvts.go); no live get_results fixture happens to send a non-empty type, method or text, so those parts are evidenced by get_nvts' own schema declaration and generated type, not by get_results' own captures]
  EXTRA         response /get_results_response/result/oci_image
                Go field OCIImage *commands.ResultOCIImage has no counterpart in the schema
  EXTRA         response /get_results_response/result_count/#text
                Go *commands.ResultCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_results_response/sort/field/#text
                Go commands.ResultField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_results_response/result/nvt/refs
                schema says required, Refs is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  SCHEMA-DEFECT response /get_results_response/result/delta/notes/note/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_results_response/result/delta/overrides/override/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_results_response/result/notes/note/active
                duplicate-ele: redefined in the same scope
  SCHEMA-DEFECT response /get_results_response/result/overrides/override/active
                duplicate-ele: redefined in the same scope

### get_roles  (command_get_roles.xml)

MISSING=1  EXTRA=2  OMITEMPTY=1

  MISSING       response /get_roles_response/sort/#text
                schema also declares character data (text) here; Go commands.RoleSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_roles_response/role_count/#text
                Go *commands.RoleCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: command_get_roles.xml's role_count <pattern> declares only <e>filtered</e>/<e>page</e>, no chardata, but the command's own baked <example> shows a bare leading total ("1<filtered>1</filtered><page>1</page>") and every live fixture agrees — internal/gmptest/testdata/live/get_roles/001_response.xml sends <role_count>6<filtered>6</filtered><page>6</page></role_count>. A zero total is equally real: internal/gmptest/testdata/live/get_alerts/001_response.xml's matching <alert_count>0<filtered>0</filtered><page>0</page></alert_count> proves the leading text can be "0", which rules out modeling this with omitempty (an int chardata field with omitempty added to its xml tag drops the whole element on marshal, zero or not — confirmed by a standalone encoding/xml experiment — so this and every entry below omits OmitEmpty entirely)]
  EXTRA         response /get_roles_response/sort/field/#text
                Go commands.RoleField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_roles_response/role_count
                schema says required, RoleCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### get_scanners  (command_get_scanners.xml)

MISSING=1  EXTRA=4  OMITEMPTY=1  SCHEMA-DEFECT=10

  MISSING       response /get_scanners_response/sort/#text
                schema also declares character data (text) here; Go commands.ScannerSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_scanners_response/scanner/disabled
                Go field Disabled *commands.GMPBool has no counterpart in the schema
  EXTRA         response /get_scanners_response/scanner/tasks/task/usage_type
                Go field UsageType commands.ScannerUsageType has no counterpart in the schema
  EXTRA         response /get_scanners_response/scanner_count/#text
                Go *commands.ScannerCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_scanners_response/sort/field/#text
                Go commands.ScannerField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_scanners_response/scanner_count
                schema says required, ScannerCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
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

MISSING=1  EXTRA=2  OMITEMPTY=1  SCHEMA-DEFECT=1

  MISSING       response /get_schedules_response/sort/#text
                schema also declares character data (text) here; Go commands.ScheduleSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_schedules_response/schedule_count/#text
                Go *commands.ScheduleCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  EXTRA         response /get_schedules_response/sort/field/#text
                Go commands.ScheduleField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  OMITEMPTY     response /get_schedules_response/schedule_count
                schema says required, ScheduleCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  SCHEMA-DEFECT response /get_schedules_response/schedule/tasks/task/permissions
                empty-pattern: pattern declares no content

### get_settings  (command_get_settings.xml)

EXTRA=2

  EXTRA         response /get_settings_response/setting/@id
                Go field ID string has no counterpart in the schema  [suppressed: T-9 settings review: command_get_settings.xml's formal <setting> pattern declares only name/comment/value, no id attribute, but internal/gmptest/testdata/live/get_settings/{001,002,003}_response.xml all show every <setting id="..."> carrying a real UUID — the schema's own baked <example> agrees a setting has an id, it just nests <setting> inside <settings> in a way no live capture ever confirms]
  EXTRA         response /get_settings_response/settings/setting
                Go field Setting []commands.SettingsSetting has no counterpart in the schema  [suppressed: T-9 settings review: command_get_settings.xml's <settings> pattern declares only the start/max pagination attributes, but the command's own <example> — reproduced verbatim in testdata/upstream/get_settings/001_response.xml — nests a real <setting id="..."> inside <settings>, the same pattern-vs-example disagreement as verify_scanner's version; without this field the round trip silently drops that nested setting (confirmed by the upstream fixture's DROPPED settings/setting finding). No live capture against current gvmd ever produces this nested shape — every live fixture uses the flat top-level <setting> siblings already modelled by GetSettingsResponse.Setting — so this is an alternate/legacy wire form kept alongside the flat one, not a replacement for it]

### get_system_reports  (command_get_system_reports.xml)

EXTRA=3  TYPE=1

  EXTRA         response /get_system_reports_response/system_report/report/@end_time
                Go field EndTime *commands.GMPTime has no counterpart in the schema  [suppressed: T-9 system_report review: same undeclared-but-real, genuinely optional attribute as start_time above, same evidence and same pointer-typing rationale]
  EXTRA         response /get_system_reports_response/system_report/report/@start_time
                Go field StartTime *commands.GMPTime has no counterpart in the schema  [suppressed: T-9 system_report review: command_get_system_reports.xml's <report> pattern never declares start_time, but internal/gmptest/testdata/live/get_system_reports/001_response.xml sends a real start_time="0001-01-01T00:00:00Z" attribute — undeclared but genuine wire data, confirmed dropped by that fixture's own pre-adoption fidelity report; testdata/upstream/get_system_reports/002_response.xml proves the attribute is also genuinely absent on other reports, so it is pointer-typed (unlike this codebase's usual non-pointer GMPTime convention) to actually honor omitempty on re-encode rather than fabricate a zero timestamp for a report that never had one]
  EXTRA         response /get_system_reports_response/truncated
                Go field Truncated string has no counterpart in the schema  [suppressed: T-9 system_report review: command_get_system_reports.xml's response <pattern> never declares <truncated>, but the command's own <example> — reproduced verbatim in testdata/upstream/get_system_reports/001_response.xml — shows a real <truncated>...</truncated> sibling of <system_report>, the same pattern-vs-example disagreement as get_settings' nested setting; no live capture against current gvmd ever populates it (truncation only occurs on long-running reports the pinned stack's fixtures never triggered), so this is weaker evidence than most restored fields here, but it follows the same precedent of trusting the schema's own baked example over a formal pattern that has already been shown to omit real fields elsewhere in this command]
  TYPE          response /get_system_reports_response/system_report/report/@duration
                schema type integer, Go Duration is string (want int)

### get_tags  (command_get_tags.xml)

MISSING=1  EXTRA=2  OMITEMPTY=1  SCHEMA-DEFECT=1

  MISSING       response /get_tags_response/sort/#text
                schema also declares character data (text) here; Go commands.TagSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_tags_response/sort/field/#text
                Go commands.TagField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_tags_response/tag_count/#text
                Go *commands.TagCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  OMITEMPTY     response /get_tags_response/tag_count
                schema says required, TagCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  SCHEMA-DEFECT response /get_tags_response/tag/resources/resource/permissions
                empty-pattern: pattern declares no content

### get_targets  (command_get_targets.xml)

MISSING=1  EXTRA=3  OMITEMPTY=1  SCHEMA-DEFECT=9

  MISSING       response /get_targets_response/sort/#text
                schema also declares character data (text) here; Go commands.TargetSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_targets_response/sort/field/#text
                Go commands.TargetField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_targets_response/target/alive_tests/*
                schema declares a "alive_test" leaf; Go models 1 children
  EXTRA         response /get_targets_response/target_count/#text
                Go *commands.TargetCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  OMITEMPTY     response /get_targets_response/target_count
                schema says required, TargetCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
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

MISSING=2  EXTRA=2  OMITEMPTY=1  SCHEMA-DEFECT=10

  MISSING       response /get_tasks_response/sort/#text
                schema also declares character data (text) here; Go commands.TaskSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  MISSING       response /get_tasks_response/task/hosts_ordering
                optional elem hosts_ordering is not modelled  [fidelity]
  EXTRA         response /get_tasks_response/sort/field/#text
                Go commands.TaskField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_tasks_response/task_count/#text
                Go *commands.TaskCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  OMITEMPTY     response /get_tasks_response/task_count
                schema says required, TaskCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
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

MISSING=1  EXTRA=2  OMITEMPTY=1

  MISSING       response /get_tickets_response/sort/#text
                schema also declares character data (text) here; Go commands.TicketSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_tickets_response/sort/field/#text
                Go commands.TicketField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_tickets_response/ticket_count/#text
                Go *commands.TicketCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  OMITEMPTY     response /get_tickets_response/ticket_count
                schema says required, TicketCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### get_tls_certificates  (command_get_tls_certificates.xml)

MISSING=1  EXTRA=4  OMITEMPTY=1

  MISSING       response /get_tls_certificates_response/sort/#text
                schema also declares character data (text) here; Go commands.TLSCertificateSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_tls_certificates_response/sort/field/#text
                Go commands.TLSCertificateField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_tls_certificates_response/tls_certificate/certificate/#text
                Go commands.TLSCertificateCertificate has a `,chardata` field; schema declares no character data here  [suppressed: T-9 finding 2: schema's certificate pattern declares only @format, no chardata token, but internal/gmptest/testdata/upstream/get_tls_certificates/001_response.xml sends real base64 certificate bytes as this element's text (see IncludeCertificateData)]
  EXTRA         response /get_tls_certificates_response/tls_certificate/valid
                Go field Valid commands.GMPBool has no counterpart in the schema  [suppressed: T-9 finding 1: schema's tls_certificate pattern never declares <valid>, but internal/gmptest/testdata/live/get_tls_certificates/003_response.xml and the upstream example both send <valid>1</valid> as real data]
  EXTRA         response /get_tls_certificates_response/tls_certificate_count/#text
                Go *commands.TLSCertificateCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  OMITEMPTY     response /get_tls_certificates_response/tls_certificate_count
                schema says required, TLSCertificateCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### get_users  (command_get_users.xml)

MISSING=1  EXTRA=2  OMITEMPTY=1  SCHEMA-DEFECT=2

  MISSING       response /get_users_response/sort/#text
                schema also declares character data (text) here; Go commands.UserSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_users_response/sort/field/#text
                Go commands.UserField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_users_response/user_count/#text
                Go *commands.UserCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  OMITEMPTY     response /get_users_response/user_count
                schema says required, UserCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]
  SCHEMA-DEFECT response /get_users_response/user/groups/group/permissions
                empty-pattern: pattern declares no content
  SCHEMA-DEFECT response /get_users_response/user/role/permissions
                empty-pattern: pattern declares no content

### get_version  (command_get_version.xml)

SCHEMA-DEFECT=1

  SCHEMA-DEFECT request  /get_version
                empty-pattern: pattern declares no content

### get_vulns  (command_get_vulns.xml)

MISSING=1  EXTRA=3  OMITEMPTY=1

  MISSING       response /get_vulns_response/sort/#text
                schema also declares character data (text) here; Go commands.VulnSort has no `,chardata` field  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_vulns_response/sort/field/#text
                Go commands.VulnField has a `,chardata` field; schema declares no character data here  [suppressed: T-18a: sort-field chardata is nested under <field>]
  EXTRA         response /get_vulns_response/vuln/type
                Go field Type string has no counterpart in the schema  [suppressed: T-9 vuln review: command_get_vulns.xml's vuln <pattern> list omits <e>type</e>, but the command still declares a full <ele><name>type</name></ele> body for it a few lines down (alts nvt/cve) — a dangling definition the pattern forgot to reference, not evidence the field doesn't exist; every fixture (live and upstream) sends a real <type>nvt</type>, the same pattern-omits-a-fully-defined-element class as verify_scanner's version]
  EXTRA         response /get_vulns_response/vuln_count/#text
                Go *commands.VulnCount has a `,chardata` field; schema declares no character data here  [suppressed: T-23: count-wrapper chardata total is schema-undeclared]
  OMITEMPTY     response /get_vulns_response/vuln_count
                schema says required, VulnCount is tagged omitempty  [suppressed: inert: responses are decoded, and omitempty only affects encoding]

### help  (command_help.xml)

MISSING=1  SCHEMA-DEFECT=12

  MISSING       response /help_response/schema
                optional elem schema is not modelled
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

### modify_alert  (command_modify_alert.xml)

conformant

### modify_asset  (command_modify_asset.xml)

conformant

### modify_auth  (command_modify_auth.xml)

EXTRA=1

  EXTRA         request  /modify_auth/group/@name
                Go field Name string has no counterpart in the schema  [suppressed: T-9 auth review: command_modify_auth.xml's request <pattern> for <group> declares no name attrib at all, but the command's own <example> requires <group name="Foo"> to identify which group to modify — the schema's pattern and its own example disagree, and without this field the generated request could never target a specific group]

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

