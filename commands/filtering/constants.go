package filtering

// Task status constants based on GMP protocol.
const (
	StatusNew           = "New"
	StatusRequested     = "Requested"
	StatusQueued        = "Queued"
	StatusRunning       = "Running"
	StatusStopRequested = "Stop Requested"
	StatusStopped       = "Stopped"
	StatusInterrupted   = "Interrupted"
	StatusDone          = "Done"
	StatusContainer     = "Container"
)

// Sort order constants.
const (
	SortAsc  = "asc"
	SortDesc = "desc"
)

// Sort field constants for tasks.
const (
	SortByName               = "name"
	SortByComment            = "comment"
	SortByCreated            = "created"
	SortByModified           = "modified"
	SortByOwner              = "owner"
	SortByStatus             = "status"
	SortByTotal              = "total"
	SortByThreat             = "threat"
	SortByTrend              = "trend"
	SortBySeverity           = "severity"
	SortBySchedule           = "schedule"
	SortByNextDue            = "next_due"
	SortByFirstReportCreated = "first_report_created"
	SortByLastReportCreated  = "last_report_created"
	SortByFalsePositive      = "false_positive"
	SortByLog                = "log"
	SortByLow                = "low"
	SortByMedium             = "medium"
	SortByHigh               = "high"
	SortByCritical           = "critical"
	SortByHosts              = "hosts"
	SortByResultHosts        = "result_hosts"
	SortByTarget             = "target"
	SortByUsageType          = "usage_type"
)

// Sort field constants for results.
const (
	SortByResultName        = "name"
	SortByResultComment     = "comment"
	SortByResultCreated     = "created"
	SortByResultModified    = "modified"
	SortByResultOwner       = "owner"
	SortByResultSeverity    = "severity"
	SortByResultQoD         = "qod"
	SortByResultHost        = "host"
	SortByResultPort        = "port"
	SortByResultNVT         = "nvt"
	SortByResultThreat      = "threat"
	SortByResultTask        = "task"
	SortByResultReport      = "report"
	SortByResultLocation    = "location"
	SortByResultDescription = "description"
)

// Sort field constants for assets.
const (
	SortByAssetName          = "name"
	SortByAssetComment       = "comment"
	SortByAssetCreated       = "created"
	SortByAssetModified      = "modified"
	SortByAssetOwner         = "owner"
	SortByAssetHosts         = "hosts"
	SortByAssetOS            = "os"
	SortByAssetHostCount     = "host_count"
	SortByAssetSeverity      = "severity"
	SortByAssetHigh          = "high"
	SortByAssetMedium        = "medium"
	SortByAssetLow           = "low"
	SortByAssetLog           = "log"
	SortByAssetFalsePositive = "false_positive"
)

// Sort field constants for targets.
const (
	SortByTargetName      = "name"
	SortByTargetComment   = "comment"
	SortByTargetCreated   = "created"
	SortByTargetModified  = "modified"
	SortByTargetOwner     = "owner"
	SortByTargetHosts     = "hosts"
	SortByTargetMaxHosts  = "max_hosts"
	SortByTargetInUse     = "in_use"
	SortByTargetTasks     = "tasks"
	SortByTargetAlterable = "alterable"
	SortByTargetWritable  = "writable"
)

// Sort field constants for tickets.
const (
	SortByTicketName       = "name"
	SortByTicketComment    = "comment"
	SortByTicketCreated    = "created"
	SortByTicketModified   = "modified"
	SortByTicketOwner      = "owner"
	SortByTicketStatus     = "status"
	SortByTicketSeverity   = "severity"
	SortByTicketHost       = "host"
	SortByTicketLocation   = "location"
	SortByTicketNVT        = "nvt"
	SortByTicketTask       = "task"
	SortByTicketResult     = "result"
	SortByTicketAssignedTo = "assigned_to"
	SortByTicketAssignedBy = "assigned_by"
)

// Sort field constants for port lists.
const (
	SortByPortListName      = "name"
	SortByPortListComment   = "comment"
	SortByPortListCreated   = "created"
	SortByPortListModified  = "modified"
	SortByPortListOwner     = "owner"
	SortByPortListPorts     = "ports"
	SortByPortListInUse     = "in_use"
	SortByPortListAlterable = "alterable"
	SortByPortListWritable  = "writable"
)

// Sort field constants for settings.
const (
	SortBySettingName      = "name"
	SortBySettingComment   = "comment"
	SortBySettingCreated   = "created"
	SortBySettingModified  = "modified"
	SortBySettingOwner     = "owner"
	SortBySettingValue     = "value"
	SortBySettingType      = "type"
	SortBySettingScanner   = "scanner"
	SortBySettingAlterable = "alterable"
	SortBySettingWritable  = "writable"
)

// Usage type constants.
const (
	UsageTypeScan  = "scan"
	UsageTypeAudit = "audit"
)

// Trend constants.
const (
	TrendUp   = "up"
	TrendDown = "down"
	TrendSame = "same"
)

// Threat level constants.
const (
	ThreatNone     = "None"
	ThreatLog      = "Log"
	ThreatLow      = "Low"
	ThreatMedium   = "Medium"
	ThreatHigh     = "High"
	ThreatCritical = "Critical"
)

// Common operators for filter conditions.
const (
	OpEqual              = "="
	OpNotEqual           = "!="
	OpLessThan           = "<"
	OpLessThanOrEqual    = "<="
	OpGreaterThan        = ">"
	OpGreaterThanOrEqual = ">="
	OpLike               = "~"
	OpNotLike            = "!~"
	OpRegex              = "=~"
	OpNotRegex           = "!~"
)

// Common filter columns for tasks.
const (
	TaskColumnUUID               = "uuid"
	TaskColumnName               = "name"
	TaskColumnComment            = "comment"
	TaskColumnCreated            = "created"
	TaskColumnModified           = "modified"
	TaskColumnOwner              = "owner"
	TaskColumnStatus             = "status"
	TaskColumnTotal              = "total"
	TaskColumnThreat             = "threat"
	TaskColumnTrend              = "trend"
	TaskColumnSeverity           = "severity"
	TaskColumnSchedule           = "schedule"
	TaskColumnNextDue            = "next_due"
	TaskColumnFirstReportCreated = "first_report_created"
	TaskColumnLastReportCreated  = "last_report_created"
	TaskColumnFalsePositive      = "false_positive"
	TaskColumnLog                = "log"
	TaskColumnLow                = "low"
	TaskColumnMedium             = "medium"
	TaskColumnHigh               = "high"
	TaskColumnCritical           = "critical"
	TaskColumnHosts              = "hosts"
	TaskColumnResultHosts        = "result_hosts"
	TaskColumnTarget             = "target"
	TaskColumnUsageType          = "usage_type"
	TaskColumnTag                = "tag"
	TaskColumnTagID              = "tag_id"
)

// Common filter columns for results.
const (
	ResultColumnUUID        = "uuid"
	ResultColumnName        = "name"
	ResultColumnComment     = "comment"
	ResultColumnCreated     = "created"
	ResultColumnModified    = "modified"
	ResultColumnOwner       = "owner"
	ResultColumnSeverity    = "severity"
	ResultColumnQoD         = "qod"
	ResultColumnHost        = "host"
	ResultColumnPort        = "port"
	ResultColumnNVT         = "nvt"
	ResultColumnThreat      = "threat"
	ResultColumnTask        = "task"
	ResultColumnReport      = "report"
	ResultColumnLocation    = "location"
	ResultColumnDescription = "description"
	ResultColumnTag         = "tag"
	ResultColumnTagID       = "tag_id"
)

// Common filter columns for assets.
const (
	AssetColumnUUID          = "uuid"
	AssetColumnName          = "name"
	AssetColumnComment       = "comment"
	AssetColumnCreated       = "created"
	AssetColumnModified      = "modified"
	AssetColumnOwner         = "owner"
	AssetColumnHosts         = "hosts"
	AssetColumnOS            = "os"
	AssetColumnHostCount     = "host_count"
	AssetColumnSeverity      = "severity"
	AssetColumnHigh          = "high"
	AssetColumnMedium        = "medium"
	AssetColumnLow           = "low"
	AssetColumnLog           = "log"
	AssetColumnFalsePositive = "false_positive"
	AssetColumnTag           = "tag"
	AssetColumnTagID         = "tag_id"
	AssetColumnType          = "asset_type"
	AssetColumnHost          = "host"
	AssetColumnLocation      = "location"
)

// Common filter columns for targets.
const (
	TargetColumnUUID      = "uuid"
	TargetColumnName      = "name"
	TargetColumnComment   = "comment"
	TargetColumnCreated   = "created"
	TargetColumnModified  = "modified"
	TargetColumnOwner     = "owner"
	TargetColumnHosts     = "hosts"
	TargetColumnMaxHosts  = "max_hosts"
	TargetColumnInUse     = "in_use"
	TargetColumnTasks     = "tasks"
	TargetColumnAlterable = "alterable"
	TargetColumnWritable  = "writable"
	TargetColumnTag       = "tag"
	TargetColumnTagID     = "tag_id"
)

// Common filter columns for tickets.
const (
	TicketColumnUUID       = "uuid"
	TicketColumnName       = "name"
	TicketColumnComment    = "comment"
	TicketColumnCreated    = "created"
	TicketColumnModified   = "modified"
	TicketColumnOwner      = "owner"
	TicketColumnStatus     = "status"
	TicketColumnSeverity   = "severity"
	TicketColumnHost       = "host"
	TicketColumnLocation   = "location"
	TicketColumnNVT        = "nvt"
	TicketColumnTask       = "task"
	TicketColumnResult     = "result"
	TicketColumnAssignedTo = "assigned_to"
	TicketColumnAssignedBy = "assigned_by"
	TicketColumnTag        = "tag"
	TicketColumnTagID      = "tag_id"
)

// Common filter columns for port lists.
const (
	PortListColumnUUID      = "uuid"
	PortListColumnName      = "name"
	PortListColumnComment   = "comment"
	PortListColumnCreated   = "created"
	PortListColumnModified  = "modified"
	PortListColumnOwner     = "owner"
	PortListColumnPorts     = "ports"
	PortListColumnInUse     = "in_use"
	PortListColumnAlterable = "alterable"
	PortListColumnWritable  = "writable"
	PortListColumnTag       = "tag"
	PortListColumnTagID     = "tag_id"
)

// Common filter columns for settings.
const (
	SettingColumnUUID      = "uuid"
	SettingColumnName      = "name"
	SettingColumnComment   = "comment"
	SettingColumnCreated   = "created"
	SettingColumnModified  = "modified"
	SettingColumnOwner     = "owner"
	SettingColumnValue     = "value"
	SettingColumnType      = "type"
	SettingColumnScanner   = "scanner"
	SettingColumnAlterable = "alterable"
	SettingColumnWritable  = "writable"
	SettingColumnTag       = "tag"
	SettingColumnTagID     = "tag_id"
)

// OS-specific filter columns (when type="os").
const (
	OSColumnTitle           = "title"
	OSColumnHosts           = "hosts"
	OSColumnLatestSeverity  = "latest_severity"
	OSColumnHighestSeverity = "highest_severity"
	OSColumnAverageSeverity = "average_severity"
)

// Host-specific filter columns (when type="host").
const (
	HostColumnSeverity = "severity"
	HostColumnOS       = "os"
	HostColumnOSS      = "oss"
	HostColumnHostname = "hostname"
	HostColumnIP       = "ip"
)
