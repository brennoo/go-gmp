package gmp

import (
	"context"

	"github.com/brennoo/go-gmp/commands"
	"github.com/brennoo/go-gmp/commands/pagination"
)

type Client interface {
	// Authentication
	Authenticate(cmd *commands.Authenticate) (resp *commands.AuthenticateResponse, err error)
	DescribeAuth(cmd *commands.DescribeAuth) (resp *commands.DescribeAuthResponse, err error)
	ModifyAuth(cmd *commands.ModifyAuth) (resp *commands.ModifyAuthResponse, err error)

	// User Management
	CreateUser(cmd *commands.CreateUser) (resp *commands.CreateUserResponse, err error)
	ModifyUser(cmd *commands.ModifyUser) (resp *commands.ModifyUserResponse, err error)
	GetUsers(cmd *commands.GetUsers) (resp *commands.GetUsersResponse, err error)
	DeleteUser(cmd *commands.DeleteUser) (resp *commands.DeleteUserResponse, err error)
	CreateRole(cmd *commands.CreateRole) (resp *commands.CreateRoleResponse, err error)
	ModifyRole(cmd *commands.ModifyRole) (resp *commands.ModifyRoleResponse, err error)
	GetRoles(cmd *commands.GetRoles) (resp *commands.GetRolesResponse, err error)
	DeleteRole(cmd *commands.DeleteRole) (resp *commands.DeleteRoleResponse, err error)
	CreatePermission(cmd *commands.CreatePermission) (resp *commands.CreatePermissionResponse, err error)
	ModifyPermission(cmd *commands.ModifyPermission) (resp *commands.ModifyPermissionResponse, err error)
	GetPermissions(cmd *commands.GetPermissions) (resp *commands.GetPermissionsResponse, err error)
	DeletePermission(cmd *commands.DeletePermission) (resp *commands.DeletePermissionResponse, err error)
	CreateGroup(cmd *commands.CreateGroup) (resp *commands.CreateGroupResponse, err error)
	ModifyGroup(cmd *commands.ModifyGroup) (resp *commands.ModifyGroupResponse, err error)
	GetGroups(cmd *commands.GetGroups) (resp *commands.GetGroupsResponse, err error)
	DeleteGroup(cmd *commands.DeleteGroup) (resp *commands.DeleteGroupResponse, err error)

	// Asset Management
	CreateAsset(cmd *commands.CreateAsset) (resp *commands.CreateAssetResponse, err error)
	ModifyAsset(cmd *commands.ModifyAsset) (resp *commands.ModifyAssetResponse, err error)
	GetAssets(cmd *commands.GetAssets) (resp *commands.GetAssetsResponse, err error)
	DeleteAsset(cmd *commands.DeleteAsset) (resp *commands.DeleteAssetResponse, err error)

	// Credential Management
	CreateCredential(cmd *commands.CreateCredential) (resp *commands.CreateCredentialResponse, err error)
	ModifyCredential(cmd *commands.ModifyCredential) (resp *commands.ModifyCredentialResponse, err error)
	GetCredentials(cmd *commands.GetCredentials) (resp *commands.GetCredentialsResponse, err error)
	DeleteCredential(cmd *commands.DeleteCredential) (resp *commands.DeleteCredentialResponse, err error)

	// Target Management
	CreateTarget(cmd *commands.CreateTarget) (resp *commands.CreateTargetResponse, err error)
	ModifyTarget(cmd *commands.ModifyTarget) (resp *commands.ModifyTargetResponse, err error)
	GetTargets(cmd *commands.GetTargets) (resp *commands.GetTargetsResponse, err error)
	DeleteTarget(cmd *commands.DeleteTarget) (resp *commands.DeleteTargetResponse, err error)

	// Task Management
	CreateTask(cmd *commands.CreateTask) (resp *commands.CreateTaskResponse, err error)
	ModifyTask(cmd *commands.ModifyTask) (resp *commands.ModifyTaskResponse, err error)
	GetTasks(cmd *commands.GetTasks) (resp *commands.GetTasksResponse, err error)
	DeleteTask(cmd *commands.DeleteTask) (resp *commands.DeleteTaskResponse, err error)
	StartTask(cmd *commands.StartTask) (resp *commands.StartTaskResponse, err error)
	StopTask(cmd *commands.StopTask) (resp *commands.StopTaskResponse, err error)
	ResumeTask(cmd *commands.ResumeTask) (resp *commands.ResumeTaskResponse, err error)
	MoveTask(cmd *commands.MoveTask) (resp *commands.MoveTaskResponse, err error)

	// Schedule Management
	CreateSchedule(cmd *commands.CreateSchedule) (resp *commands.CreateScheduleResponse, err error)
	ModifySchedule(cmd *commands.ModifySchedule) (resp *commands.ModifyScheduleResponse, err error)
	GetSchedules(cmd *commands.GetSchedules) (resp *commands.GetSchedulesResponse, err error)
	DeleteSchedule(cmd *commands.DeleteSchedule) (resp *commands.DeleteScheduleResponse, err error)

	// Alert Management
	CreateAlert(cmd *commands.CreateAlert) (resp *commands.CreateAlertResponse, err error)
	GetAlerts(cmd *commands.GetAlerts) (resp *commands.GetAlertsResponse, err error)
	ModifyAlert(cmd *commands.ModifyAlert) (resp *commands.ModifyAlertResponse, err error)
	DeleteAlert(cmd *commands.DeleteAlert) (resp *commands.DeleteAlertResponse, err error)
	TestAlert(cmd *commands.TestAlert) (resp *commands.TestAlertResponse, err error)

	// Report Management
	CreateReport(cmd *commands.CreateReport) (resp *commands.CreateReportResponse, err error)
	GetReports(cmd *commands.GetReports) (resp *commands.GetReportsResponse, err error)
	DeleteReport(cmd *commands.DeleteReport) (resp *commands.DeleteReportResponse, err error)
	CreateReportFormat(cmd *commands.CreateReportFormat) (resp *commands.CreateReportFormatResponse, err error)
	ModifyReportFormat(cmd *commands.ModifyReportFormat) (resp *commands.ModifyReportFormatResponse, err error)
	GetReportFormats(cmd *commands.GetReportFormats) (resp *commands.GetReportFormatsResponse, err error)
	DeleteReportFormat(cmd *commands.DeleteReportFormat) (resp *commands.DeleteReportFormatResponse, err error)
	VerifyReportFormat(cmd *commands.VerifyReportFormat) (resp *commands.VerifyReportFormatResponse, err error)
	CreateReportConfig(cmd *commands.CreateReportConfig) (resp *commands.CreateReportConfigResponse, err error)
	ModifyReportConfig(cmd *commands.ModifyReportConfig) (resp *commands.ModifyReportConfigResponse, err error)
	GetReportConfigs(cmd *commands.GetReportConfigs) (resp *commands.GetReportConfigsResponse, err error)
	DeleteReportConfig(cmd *commands.DeleteReportConfig) (resp *commands.DeleteReportConfigResponse, err error)
	GetSystemReports(cmd *commands.GetSystemReports) (resp *commands.GetSystemReportsResponse, err error)

	// Results & Vulnerabilities
	GetResults(cmd *commands.GetResults) (resp *commands.GetResultsResponse, err error)
	GetVulns(cmd *commands.GetVulns) (resp *commands.GetVulnsResponse, err error)

	// Notes, Tags, Filters, Tickets, Overrides
	CreateNote(cmd *commands.CreateNote) (resp *commands.CreateNoteResponse, err error)
	ModifyNote(cmd *commands.ModifyNote) (resp *commands.ModifyNoteResponse, err error)
	GetNotes(cmd *commands.GetNotes) (resp *commands.GetNotesResponse, err error)
	DeleteNote(cmd *commands.DeleteNote) (resp *commands.DeleteNoteResponse, err error)
	CreateTag(cmd *commands.CreateTag) (resp *commands.CreateTagResponse, err error)
	ModifyTag(cmd *commands.ModifyTag) (resp *commands.ModifyTagResponse, err error)
	GetTags(cmd *commands.GetTags) (resp *commands.GetTagsResponse, err error)
	DeleteTag(cmd *commands.DeleteTag) (resp *commands.DeleteTagResponse, err error)
	CreateFilter(cmd *commands.CreateFilter) (resp *commands.CreateFilterResponse, err error)
	ModifyFilter(cmd *commands.ModifyFilter) (resp *commands.ModifyFilterResponse, err error)
	GetFilters(cmd *commands.GetFilters) (resp *commands.GetFiltersResponse, err error)
	DeleteFilter(cmd *commands.DeleteFilter) (resp *commands.DeleteFilterResponse, err error)
	CreateTicket(cmd *commands.CreateTicket) (resp *commands.CreateTicketResponse, err error)
	ModifyTicket(cmd *commands.ModifyTicket) (resp *commands.ModifyTicketResponse, err error)
	GetTickets(cmd *commands.GetTickets) (resp *commands.GetTicketsResponse, err error)
	DeleteTicket(cmd *commands.DeleteTicket) (resp *commands.DeleteTicketResponse, err error)
	CreateOverride(cmd *commands.CreateOverride) (resp *commands.CreateOverrideResponse, err error)
	DeleteOverride(cmd *commands.DeleteOverride) (resp *commands.DeleteOverrideResponse, err error)
	GetOverrides(cmd *commands.GetOverrides) (resp *commands.GetOverridesResponse, err error)
	ModifyOverride(cmd *commands.ModifyOverride) (resp *commands.ModifyOverrideResponse, err error)

	// Configuration, Preferences, Settings
	GetConfigs(cmd *commands.GetConfigs) (resp *commands.GetConfigsResponse, err error)
	CreateConfig(cmd *commands.CreateConfig) (resp *commands.CreateConfigResponse, err error)
	ModifyConfig(cmd *commands.ModifyConfig) (resp *commands.ModifyConfigResponse, err error)
	DeleteConfig(cmd *commands.DeleteConfig) (resp *commands.DeleteConfigResponse, err error)
	SyncConfig(cmd *commands.SyncConfig) (resp *commands.SyncConfigResponse, err error)
	GetPreferences(cmd *commands.GetPreferences) (resp *commands.GetPreferencesResponse, err error)
	ModifySetting(cmd *commands.ModifySetting) (resp *commands.ModifySettingResponse, err error)
	GetSettings(cmd *commands.GetSettings) (resp *commands.GetSettingsResponse, err error)

	// Port Lists & Ranges
	CreatePortList(cmd *commands.CreatePortList) (resp *commands.CreatePortListResponse, err error)
	ModifyPortList(cmd *commands.ModifyPortList) (resp *commands.ModifyPortListResponse, err error)
	GetPortLists(cmd *commands.GetPortLists) (resp *commands.GetPortListsResponse, err error)
	DeletePortList(cmd *commands.DeletePortList) (resp *commands.DeletePortListResponse, err error)
	CreatePortRange(cmd *commands.CreatePortRange) (resp *commands.CreatePortRangeResponse, err error)
	DeletePortRange(cmd *commands.DeletePortRange) (resp *commands.DeletePortRangeResponse, err error)

	// Scanners
	CreateScanner(cmd *commands.CreateScanner) (resp *commands.CreateScannerResponse, err error)
	ModifyScanner(cmd *commands.ModifyScanner) (resp *commands.ModifyScannerResponse, err error)
	GetScanners(cmd *commands.GetScanners) (resp *commands.GetScannersResponse, err error)
	DeleteScanner(cmd *commands.DeleteScanner) (resp *commands.DeleteScannerResponse, err error)
	VerifyScanner(cmd *commands.VerifyScanner) (resp *commands.VerifyScannerResponse, err error)

	// TLS Certificates
	CreateTLSCertificate(cmd *commands.CreateTLSCertificate) (resp *commands.CreateTLSCertificateResponse, err error)
	ModifyTLSCertificate(cmd *commands.ModifyTLSCertificate) (resp *commands.ModifyTLSCertificateResponse, err error)
	GetTLSCertificates(cmd *commands.GetTLSCertificates) (resp *commands.GetTLSCertificatesResponse, err error)

	// Agents
	ModifyAgents(cmd *commands.ModifyAgents) (resp *commands.ModifyAgentsResponse, err error)
	GetAgents(cmd *commands.GetAgents) (resp *commands.GetAgentsResponse, err error)
	DeleteAgents(cmd *commands.DeleteAgents) (resp *commands.DeleteAgentsResponse, err error)

	// Feeds, Licenses, Features, Aggregates, Resource Names
	GetNVTs(cmd *commands.GetNVTs) (resp *commands.GetNVTsResponse, err error)
	GetNVTFamilies(cmd *commands.GetNVTFamilies) (resp *commands.GetNVTFamiliesResponse, err error)
	GetFeeds(cmd *commands.GetFeeds) (resp *commands.GetFeedsResponse, err error)
	ModifyLicense(cmd *commands.ModifyLicense) (resp *commands.ModifyLicenseResponse, err error)
	GetLicense(cmd *commands.GetLicense) (resp *commands.GetLicenseResponse, err error)
	GetFeatures(cmd *commands.GetFeatures) (resp *commands.GetFeaturesResponse, err error)
	GetAggregates(cmd *commands.GetAggregates) (resp *commands.GetAggregatesResponse, err error)
	GetResourceNames(cmd *commands.GetResourceNames) (resp *commands.GetResourceNamesResponse, err error)

	// Miscellaneous
	GetInfo(cmd *commands.GetInfo) (resp *commands.GetInfoResponse, err error)
	GetVersion(cmd *commands.GetVersion) (resp *commands.GetVersionResponse, err error)
	Help(cmd *commands.Help) (resp *commands.HelpResponse, err error)
	EmptyTrashcan(cmd *commands.EmptyTrashcan) (resp *commands.EmptyTrashcanResponse, err error)
	Restore(cmd *commands.Restore) (resp *commands.RestoreResponse, err error)
	RunWizard(cmd *commands.RunWizard) (resp *commands.RunWizardResponse, err error)
	RawXML(xml string) (string, error)

	// Pagination Methods
	GetTasksPaged(page, pageSize int, filters ...string) (*commands.GetTasksResponse, error)
	GetResultsPaged(page, pageSize int, filters ...string) (*commands.GetResultsResponse, error)
	GetAssetsPaged(page, pageSize int, filters ...string) (*commands.GetAssetsResponse, error)
	GetTargetsPaged(page, pageSize int, filters ...string) (*commands.GetTargetsResponse, error)
	GetTicketsPaged(page, pageSize int, filters ...string) (*commands.GetTicketsResponse, error)
	GetPortListsPaged(page, pageSize int, filters ...string) (*commands.GetPortListsResponse, error)
	GetSettingsPaged(page, pageSize int, filters ...string) (*commands.GetSettingsResponse, error)

	// Iterator Methods
	GetTasksIter(ctx context.Context, pageSize int, filters ...string) *pagination.TaskIterator
	GetResultsIter(ctx context.Context, pageSize int, filters ...string) *pagination.ResultIterator
	GetAssetsIter(ctx context.Context, pageSize int, filters ...string) *pagination.AssetIterator
	GetTargetsIter(ctx context.Context, pageSize int, filters ...string) *pagination.TargetIterator
	GetTicketsIter(ctx context.Context, pageSize int, filters ...string) *pagination.TicketIterator
	GetPortListsIter(ctx context.Context, pageSize int, filters ...string) *pagination.PortListIterator
	GetSettingsIter(ctx context.Context, pageSize int, filters ...string) *pagination.SettingsIterator
}
