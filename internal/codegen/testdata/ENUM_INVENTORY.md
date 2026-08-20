# Enum field inventory

Regenerate with: go test ./internal/codegen -run TestEnumInventory -update

T-19: every generated named string type for an <alts>-constrained field, with
its wire name, the command(s) it belongs to, and its ordered constants —
extracted directly from the generated commands/*.go bytes so this can never
disagree with what is actually committed.

total enum types: 77

AggregateRelation        file=get_aggregates.go        wire=<relation> owner=get_aggregates
    AggregateRelationEq = "="
    AggregateRelationContains = ":"
    AggregateRelationRegex = "~"
    AggregateRelationGT = ">"
    AggregateRelationLT = "<"
AlertOrder               file=get_alerts.go            wire=<order> owner=get_alerts
    AlertOrderAscending = "ascending"
    AlertOrderDescending = "descending"
AlertRelation            file=get_alerts.go            wire=<relation> owner=get_alerts
    AlertRelationEq = "="
    AlertRelationContains = ":"
    AlertRelationRegex = "~"
    AlertRelationGT = ">"
    AlertRelationLT = "<"
AlertTimeStatus          file=get_alerts.go            wire=<time_status> owner=get_alerts
    AlertTimeStatusExpired = "expired"
    AlertTimeStatusInactive = "inactive"
    AlertTimeStatusUnknown = "unknown"
    AlertTimeStatusValid = "valid"
AssetOrder               file=get_assets.go            wire=<order> owner=get_assets
    AssetOrderAscending = "ascending"
    AssetOrderDescending = "descending"
AssetRelation            file=get_assets.go            wire=<relation> owner=get_assets
    AssetRelationEq = "="
    AssetRelationContains = ":"
    AssetRelationRegex = "~"
    AssetRelationGT = ">"
    AssetRelationLT = "<"
AuthRole                 file=authenticate.go          wire=<role> owner=authenticate
    AuthRoleUser = "User"
    AuthRoleAdmin = "Admin"
    AuthRoleSuperAdmin = "Super Admin"
    AuthRoleObserver = "Observer"
AuthTimeStatus           file=describe_auth.go         wire=<time_status> owner=describe_auth
    AuthTimeStatusExpired = "expired"
    AuthTimeStatusInactive = "inactive"
    AuthTimeStatusUnknown = "unknown"
    AuthTimeStatusValid = "valid"
ConfigOrder              file=get_configs.go           wire=<order> owner=get_configs
    ConfigOrderAscending = "ascending"
    ConfigOrderDescending = "descending"
ConfigRelation           file=get_configs.go           wire=<relation> owner=get_configs
    ConfigRelationEq = "="
    ConfigRelationContains = ":"
    ConfigRelationRegex = "~"
    ConfigRelationGT = ">"
    ConfigRelationLT = "<"
ConfigType               file=get_configs.go           wire=<type> owner=get_configs
    ConfigTypeN0 = "0"
    ConfigTypeN1 = "1"
ConfigUsageType          file=config_shared.go         wire=<usage_type> owner=create_config and get_configs
    ConfigUsageTypeScan = "scan"
    ConfigUsageTypePolicy = "policy"
CredentialAlgorithm      file=credential_shared.go     wire=<algorithm> owner=create_credential, get_credentials and modify_credential
    CredentialAlgorithmAes = "aes"
    CredentialAlgorithmDes = "des"
CredentialAuthAlgorithm  file=credential_shared.go     wire=<auth_algorithm> owner=create_credential, get_credentials and modify_credential
    CredentialAuthAlgorithmMD5 = "md5"
    CredentialAuthAlgorithmSha1 = "sha1"
CredentialFormat         file=get_credentials.go       wire=<format> owner=get_credentials
    CredentialFormatKey = "key"
    CredentialFormatRpm = "rpm"
    CredentialFormatDeb = "deb"
    CredentialFormatExe = "exe"
    CredentialFormatPem = "pem"
CredentialOrder          file=get_credentials.go       wire=<order> owner=get_credentials
    CredentialOrderAscending = "ascending"
    CredentialOrderDescending = "descending"
CredentialRelation       file=get_credentials.go       wire=<relation> owner=get_credentials
    CredentialRelationEq = "="
    CredentialRelationContains = ":"
    CredentialRelationRegex = "~"
    CredentialRelationGT = ">"
    CredentialRelationLT = "<"
CredentialTimeStatus     file=get_credentials.go       wire=<time_status> owner=get_credentials
    CredentialTimeStatusExpired = "expired"
    CredentialTimeStatusInactive = "inactive"
    CredentialTimeStatusUnknown = "unknown"
    CredentialTimeStatusValid = "valid"
CredentialType           file=credential_shared.go     wire=<type> owner=create_credential and get_credentials
    CredentialTypeCc = "cc"
    CredentialTypeKrb5 = "krb5"
    CredentialTypePgp = "pgp"
    CredentialTypePw = "pw"
    CredentialTypeSmime = "smime"
    CredentialTypeSNMP = "snmp"
    CredentialTypeUp = "up"
    CredentialTypeUsk = "usk"
FilterOrder              file=get_filters.go           wire=<order> owner=get_filters
    FilterOrderAscending = "ascending"
    FilterOrderDescending = "descending"
FilterRelation           file=get_filters.go           wire=<relation> owner=get_filters
    FilterRelationEq = "="
    FilterRelationContains = ":"
    FilterRelationRegex = "~"
    FilterRelationGT = ">"
    FilterRelationLT = "<"
GetReportFormatsTrust    file=get_report_formats.go    wire=<trust> owner=get_report_formats
    GetReportFormatsTrustYes = "yes"
    GetReportFormatsTrustNo = "no"
    GetReportFormatsTrustUnknown = "unknown"
GroupOrder               file=get_groups.go            wire=<order> owner=get_groups
    GroupOrderAscending = "ascending"
    GroupOrderDescending = "descending"
GroupRelation            file=get_groups.go            wire=<relation> owner=get_groups
    GroupRelationEq = "="
    GroupRelationContains = ":"
    GroupRelationRegex = "~"
    GroupRelationGT = ">"
    GroupRelationLT = "<"
NVTCategory              file=get_nvts.go              wire=<category> owner=get_nvts
    NVTCategoryN0 = "0"
    NVTCategoryN1 = "1"
    NVTCategoryN2 = "2"
    NVTCategoryN3 = "3"
    NVTCategoryN4 = "4"
    NVTCategoryN5 = "5"
    NVTCategoryN6 = "6"
    NVTCategoryN7 = "7"
    NVTCategoryN8 = "8"
    NVTCategoryN9 = "9"
    NVTCategoryN10 = "10"
    NVTCategoryN11 = "11"
NVTOrder                 file=get_info.go              wire=<order> owner=get_info
    NVTOrderAscending = "ascending"
    NVTOrderDescending = "descending"
NVTRelation              file=get_info.go              wire=<relation> owner=get_info
    NVTRelationEq = "="
    NVTRelationContains = ":"
    NVTRelationRegex = "~"
    NVTRelationGT = ">"
    NVTRelationLT = "<"
NoteOrder                file=get_notes.go             wire=<order> owner=get_notes
    NoteOrderAscending = "ascending"
    NoteOrderDescending = "descending"
NoteRelation             file=get_notes.go             wire=<relation> owner=get_notes
    NoteRelationEq = "="
    NoteRelationContains = ":"
    NoteRelationRegex = "~"
    NoteRelationGT = ">"
    NoteRelationLT = "<"
OverrideOrder            file=get_overrides.go         wire=<order> owner=get_overrides
    OverrideOrderAscending = "ascending"
    OverrideOrderDescending = "descending"
OverrideRelation         file=get_overrides.go         wire=<relation> owner=get_overrides
    OverrideRelationEq = "="
    OverrideRelationContains = ":"
    OverrideRelationRegex = "~"
    OverrideRelationGT = ">"
    OverrideRelationLT = "<"
PermissionOrder          file=get_permissions.go       wire=<order> owner=get_permissions
    PermissionOrderAscending = "ascending"
    PermissionOrderDescending = "descending"
PermissionRelation       file=get_permissions.go       wire=<relation> owner=get_permissions
    PermissionRelationEq = "="
    PermissionRelationContains = ":"
    PermissionRelationRegex = "~"
    PermissionRelationGT = ">"
    PermissionRelationLT = "<"
PortListOrder            file=get_port_lists.go        wire=<order> owner=get_port_lists
    PortListOrderAscending = "ascending"
    PortListOrderDescending = "descending"
PortListRelation         file=get_port_lists.go        wire=<relation> owner=get_port_lists
    PortListRelationEq = "="
    PortListRelationContains = ":"
    PortListRelationRegex = "~"
    PortListRelationGT = ">"
    PortListRelationLT = "<"
ReportApplyOverrides     file=report_shared.go         wire=<apply_overrides> owner=create_report and get_reports
    ReportApplyOverridesN0 = "0"
    ReportApplyOverridesN1 = "1"
ReportConfigOrder        file=get_report_configs.go    wire=<order> owner=get_report_configs
    ReportConfigOrderAscending = "ascending"
    ReportConfigOrderDescending = "descending"
ReportConfigParamType    file=get_report_configs.go    wire=<type> owner=get_report_configs
    ReportConfigParamTypeBoolean = "boolean"
    ReportConfigParamTypeInteger = "integer"
    ReportConfigParamTypeMultiSelection = "multi_selection"
    ReportConfigParamTypeSelection = "selection"
    ReportConfigParamTypeString = "string"
    ReportConfigParamTypeText = "text"
    ReportConfigParamTypeReportFormatList = "report_format_list"
ReportConfigRelation     file=get_report_configs.go    wire=<relation> owner=get_report_configs
    ReportConfigRelationEq = "="
    ReportConfigRelationContains = ":"
    ReportConfigRelationRegex = "~"
    ReportConfigRelationGT = ">"
    ReportConfigRelationLT = "<"
ReportFilter             file=report_shared.go         wire=<filter> owner=create_report and get_reports
    ReportFilterCritical = "Critical"
    ReportFilterHigh = "High"
    ReportFilterMedium = "Medium"
    ReportFilterLow = "Low"
    ReportFilterLog = "Log"
ReportFiltersNotes       file=report_shared.go         wire=<notes> owner=create_report and get_reports
    ReportFiltersNotesN0 = "0"
    ReportFiltersNotesN1 = "1"
ReportFiltersOverrides   file=report_shared.go         wire=<overrides> owner=create_report and get_reports
    ReportFiltersOverridesN0 = "0"
    ReportFiltersOverridesN1 = "1"
ReportFormatOrder        file=get_report_formats.go    wire=<order> owner=get_report_formats
    ReportFormatOrderAscending = "ascending"
    ReportFormatOrderDescending = "descending"
ReportFormatParamType    file=get_report_formats.go    wire=<type> owner=get_report_formats
    ReportFormatParamTypeBoolean = "boolean"
    ReportFormatParamTypeInteger = "integer"
    ReportFormatParamTypeMultiSelection = "multi_selection"
    ReportFormatParamTypeSelection = "selection"
    ReportFormatParamTypeString = "string"
    ReportFormatParamTypeText = "text"
    ReportFormatParamTypeReportFormatList = "report_format_list"
ReportFormatRelation     file=get_report_formats.go    wire=<relation> owner=get_report_formats
    ReportFormatRelationEq = "="
    ReportFormatRelationContains = ":"
    ReportFormatRelationRegex = "~"
    ReportFormatRelationGT = ">"
    ReportFormatRelationLT = "<"
ReportFormatReportType   file=get_report_formats.go    wire=<report_type> owner=get_report_formats
    ReportFormatReportTypeScan = "scan"
    ReportFormatReportTypeAudit = "audit"
    ReportFormatReportTypeAll = "all"
ReportOrder              file=report_shared.go         wire=<order> owner=create_report and get_reports
    ReportOrderAscending = "ascending"
    ReportOrderDescending = "descending"
ReportRelation           file=get_reports.go           wire=<relation> owner=get_reports
    ReportRelationEq = "="
    ReportRelationContains = ":"
    ReportRelationRegex = "~"
    ReportRelationGT = ">"
    ReportRelationLT = "<"
ReportResultHostsOnly    file=report_shared.go         wire=<result_hosts_only> owner=create_report and get_reports
    ReportResultHostsOnlyN0 = "0"
    ReportResultHostsOnlyN1 = "1"
ResultOrder              file=get_results.go           wire=<order> owner=get_results
    ResultOrderAscending = "ascending"
    ResultOrderDescending = "descending"
ResultRelation           file=get_results.go           wire=<relation> owner=get_results
    ResultRelationEq = "="
    ResultRelationContains = ":"
    ResultRelationRegex = "~"
    ResultRelationGT = ">"
    ResultRelationLT = "<"
RoleOrder                file=get_roles.go             wire=<order> owner=get_roles
    RoleOrderAscending = "ascending"
    RoleOrderDescending = "descending"
RoleRelation             file=get_roles.go             wire=<relation> owner=get_roles
    RoleRelationEq = "="
    RoleRelationContains = ":"
    RoleRelationRegex = "~"
    RoleRelationGT = ">"
    RoleRelationLT = "<"
ScannerOrder             file=get_scanners.go          wire=<order> owner=get_scanners
    ScannerOrderAscending = "ascending"
    ScannerOrderDescending = "descending"
ScannerRelation          file=get_scanners.go          wire=<relation> owner=get_scanners
    ScannerRelationEq = "="
    ScannerRelationContains = ":"
    ScannerRelationRegex = "~"
    ScannerRelationGT = ">"
    ScannerRelationLT = "<"
ScannerTimeStatus        file=get_scanners.go          wire=<time_status> owner=get_scanners
    ScannerTimeStatusExpired = "expired"
    ScannerTimeStatusInactive = "inactive"
    ScannerTimeStatusUnknown = "unknown"
    ScannerTimeStatusValid = "valid"
ScannerUsageType         file=get_scanners.go          wire=<usage_type> owner=get_scanners
    ScannerUsageTypeScan = "scan"
    ScannerUsageTypePolicy = "policy"
ScheduleOrder            file=get_schedules.go         wire=<order> owner=get_schedules
    ScheduleOrderAscending = "ascending"
    ScheduleOrderDescending = "descending"
ScheduleRelation         file=get_schedules.go         wire=<relation> owner=get_schedules
    ScheduleRelationEq = "="
    ScheduleRelationContains = ":"
    ScheduleRelationRegex = "~"
    ScheduleRelationGT = ">"
    ScheduleRelationLT = "<"
TLSCertificateOrder      file=get_tls_certificates.go  wire=<order> owner=get_tls_certificates
    TLSCertificateOrderAscending = "ascending"
    TLSCertificateOrderDescending = "descending"
TLSCertificateRelation   file=get_tls_certificates.go  wire=<relation> owner=get_tls_certificates
    TLSCertificateRelationEq = "="
    TLSCertificateRelationContains = ":"
    TLSCertificateRelationRegex = "~"
    TLSCertificateRelationGT = ">"
    TLSCertificateRelationLT = "<"
TLSCertificateTimeStatus file=get_tls_certificates.go  wire=<time_status> owner=get_tls_certificates
    TLSCertificateTimeStatusExpired = "expired"
    TLSCertificateTimeStatusInactive = "inactive"
    TLSCertificateTimeStatusUnknown = "unknown"
    TLSCertificateTimeStatusValid = "valid"
TagOrder                 file=get_tags.go              wire=<order> owner=get_tags
    TagOrderAscending = "ascending"
    TagOrderDescending = "descending"
TagRelation              file=get_tags.go              wire=<relation> owner=get_tags
    TagRelationEq = "="
    TagRelationContains = ":"
    TagRelationRegex = "~"
    TagRelationGT = ">"
    TagRelationLT = "<"
TargetOrder              file=get_targets.go           wire=<order> owner=get_targets
    TargetOrderAscending = "ascending"
    TargetOrderDescending = "descending"
TargetRelation           file=get_targets.go           wire=<relation> owner=get_targets
    TargetRelationEq = "="
    TargetRelationContains = ":"
    TargetRelationRegex = "~"
    TargetRelationGT = ">"
    TargetRelationLT = "<"
TaskApplyOverrides       file=get_tasks.go             wire=<apply_overrides> owner=get_tasks
    TaskApplyOverridesN0 = "0"
    TaskApplyOverridesN1 = "1"
TaskOrder                file=get_tasks.go             wire=<order> owner=get_tasks
    TaskOrderAscending = "ascending"
    TaskOrderDescending = "descending"
TaskRelation             file=get_tasks.go             wire=<relation> owner=get_tasks
    TaskRelationEq = "="
    TaskRelationContains = ":"
    TaskRelationRegex = "~"
    TaskRelationGT = ">"
    TaskRelationLT = "<"
TaskUsageType            file=task_shared.go           wire=<usage_type> owner=create_task and get_tasks
    TaskUsageTypeScan = "scan"
    TaskUsageTypeAudit = "audit"
TicketOrder              file=get_tickets.go           wire=<order> owner=get_tickets
    TicketOrderAscending = "ascending"
    TicketOrderDescending = "descending"
TicketRelation           file=get_tickets.go           wire=<relation> owner=get_tickets
    TicketRelationEq = "="
    TicketRelationContains = ":"
    TicketRelationRegex = "~"
    TicketRelationGT = ">"
    TicketRelationLT = "<"
UserOrder                file=get_users.go             wire=<order> owner=get_users
    UserOrderAscending = "ascending"
    UserOrderDescending = "descending"
UserRelation             file=get_users.go             wire=<relation> owner=get_users
    UserRelationEq = "="
    UserRelationContains = ":"
    UserRelationRegex = "~"
    UserRelationGT = ">"
    UserRelationLT = "<"
UserSource               file=user_shared.go           wire=<source> owner=create_user, get_users and modify_user
    UserSourceFile = "file"
    UserSourceLdapConnect = "ldap_connect"
    UserSourceRadiusConnect = "radius_connect"
VulnOrder                file=get_vulns.go             wire=<order> owner=get_vulns
    VulnOrderAscending = "ascending"
    VulnOrderDescending = "descending"
VulnRelation             file=get_vulns.go             wire=<relation> owner=get_vulns
    VulnRelationEq = "="
    VulnRelationContains = ":"
    VulnRelationRegex = "~"
    VulnRelationGT = ">"
    VulnRelationLT = "<"
