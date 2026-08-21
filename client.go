package gmp

import (
	"context"

	"github.com/brennoo/go-gmp/commands"
	"github.com/brennoo/go-gmp/commands/filtering"
	"github.com/brennoo/go-gmp/commands/pagination"
)

type Client interface {
	// Authentication
	Authenticate(ctx context.Context, cmd *commands.Authenticate) (resp *commands.AuthenticateResponse, err error)
	DescribeAuth(ctx context.Context, cmd *commands.DescribeAuth) (resp *commands.DescribeAuthResponse, err error)
	ModifyAuth(ctx context.Context, cmd *commands.ModifyAuth) (resp *commands.ModifyAuthResponse, err error)

	// User Management
	CreateUser(ctx context.Context, cmd *commands.CreateUser) (resp *commands.CreateUserResponse, err error)
	ModifyUser(ctx context.Context, cmd *commands.ModifyUser) (resp *commands.ModifyUserResponse, err error)
	GetUsers(ctx context.Context, args ...filtering.FilterArg) (*commands.GetUsersResponse, error)
	DeleteUser(ctx context.Context, cmd *commands.DeleteUser) (resp *commands.DeleteUserResponse, err error)
	CreateRole(ctx context.Context, cmd *commands.CreateRole) (resp *commands.CreateRoleResponse, err error)
	ModifyRole(ctx context.Context, cmd *commands.ModifyRole) (resp *commands.ModifyRoleResponse, err error)
	GetRoles(ctx context.Context, cmd *commands.GetRoles) (resp *commands.GetRolesResponse, err error)
	DeleteRole(ctx context.Context, cmd *commands.DeleteRole) (resp *commands.DeleteRoleResponse, err error)
	CreatePermission(ctx context.Context, cmd *commands.CreatePermission) (resp *commands.CreatePermissionResponse, err error)
	ModifyPermission(ctx context.Context, cmd *commands.ModifyPermission) (resp *commands.ModifyPermissionResponse, err error)
	GetPermissions(ctx context.Context, cmd *commands.GetPermissions) (resp *commands.GetPermissionsResponse, err error)
	DeletePermission(ctx context.Context, cmd *commands.DeletePermission) (resp *commands.DeletePermissionResponse, err error)
	CreateGroup(ctx context.Context, cmd *commands.CreateGroup) (resp *commands.CreateGroupResponse, err error)
	ModifyGroup(ctx context.Context, cmd *commands.ModifyGroup) (resp *commands.ModifyGroupResponse, err error)
	GetGroups(ctx context.Context, cmd *commands.GetGroups) (resp *commands.GetGroupsResponse, err error)
	DeleteGroup(ctx context.Context, cmd *commands.DeleteGroup) (resp *commands.DeleteGroupResponse, err error)

	// Asset Management
	CreateAsset(ctx context.Context, cmd *commands.CreateAsset) (resp *commands.CreateAssetResponse, err error)
	ModifyAsset(ctx context.Context, cmd *commands.ModifyAsset) (resp *commands.ModifyAssetResponse, err error)
	GetAssets(ctx context.Context, args ...filtering.FilterArg) (*commands.GetAssetsResponse, error)
	DeleteAsset(ctx context.Context, cmd *commands.DeleteAsset) (resp *commands.DeleteAssetResponse, err error)

	// Credential Management
	CreateCredential(ctx context.Context, cmd *commands.CreateCredential) (resp *commands.CreateCredentialResponse, err error)
	ModifyCredential(ctx context.Context, cmd *commands.ModifyCredential) (resp *commands.ModifyCredentialResponse, err error)
	GetCredentials(ctx context.Context, cmd *commands.GetCredentials) (resp *commands.GetCredentialsResponse, err error)
	DeleteCredential(ctx context.Context, cmd *commands.DeleteCredential) (resp *commands.DeleteCredentialResponse, err error)

	// Target Management
	CreateTarget(ctx context.Context, cmd *commands.CreateTarget) (resp *commands.CreateTargetResponse, err error)
	ModifyTarget(ctx context.Context, cmd *commands.ModifyTarget) (resp *commands.ModifyTargetResponse, err error)
	GetTargets(ctx context.Context, args ...filtering.FilterArg) (*commands.GetTargetsResponse, error)
	DeleteTarget(ctx context.Context, cmd *commands.DeleteTarget) (resp *commands.DeleteTargetResponse, err error)

	// Task Management
	CreateTask(ctx context.Context, cmd *commands.CreateTask) (resp *commands.CreateTaskResponse, err error)
	ModifyTask(ctx context.Context, cmd *commands.ModifyTask) (resp *commands.ModifyTaskResponse, err error)
	GetTasks(ctx context.Context, args ...filtering.FilterArg) (*commands.GetTasksResponse, error)
	DeleteTask(ctx context.Context, cmd *commands.DeleteTask) (resp *commands.DeleteTaskResponse, err error)
	StartTask(ctx context.Context, cmd *commands.StartTask) (resp *commands.StartTaskResponse, err error)
	StopTask(ctx context.Context, cmd *commands.StopTask) (resp *commands.StopTaskResponse, err error)
	ResumeTask(ctx context.Context, cmd *commands.ResumeTask) (resp *commands.ResumeTaskResponse, err error)
	MoveTask(ctx context.Context, cmd *commands.MoveTask) (resp *commands.MoveTaskResponse, err error)

	// Schedule Management
	CreateSchedule(ctx context.Context, cmd *commands.CreateSchedule) (resp *commands.CreateScheduleResponse, err error)
	ModifySchedule(ctx context.Context, cmd *commands.ModifySchedule) (resp *commands.ModifyScheduleResponse, err error)
	GetSchedules(ctx context.Context, cmd *commands.GetSchedules) (resp *commands.GetSchedulesResponse, err error)
	DeleteSchedule(ctx context.Context, cmd *commands.DeleteSchedule) (resp *commands.DeleteScheduleResponse, err error)

	// Alert Management
	CreateAlert(ctx context.Context, cmd *commands.CreateAlert) (resp *commands.CreateAlertResponse, err error)
	GetAlerts(ctx context.Context, cmd *commands.GetAlerts) (resp *commands.GetAlertsResponse, err error)
	ModifyAlert(ctx context.Context, cmd *commands.ModifyAlert) (resp *commands.ModifyAlertResponse, err error)
	DeleteAlert(ctx context.Context, cmd *commands.DeleteAlert) (resp *commands.DeleteAlertResponse, err error)
	TestAlert(ctx context.Context, cmd *commands.TestAlert) (resp *commands.TestAlertResponse, err error)

	// Report Management
	CreateReport(ctx context.Context, cmd *commands.CreateReport) (resp *commands.CreateReportResponse, err error)
	GetReports(ctx context.Context, cmd *commands.GetReports) (resp *commands.GetReportsResponse, err error)
	DeleteReport(ctx context.Context, cmd *commands.DeleteReport) (resp *commands.DeleteReportResponse, err error)
	CreateReportFormat(ctx context.Context, cmd *commands.CreateReportFormat) (resp *commands.CreateReportFormatResponse, err error)
	ModifyReportFormat(ctx context.Context, cmd *commands.ModifyReportFormat) (resp *commands.ModifyReportFormatResponse, err error)
	GetReportFormats(ctx context.Context, cmd *commands.GetReportFormats) (resp *commands.GetReportFormatsResponse, err error)
	DeleteReportFormat(ctx context.Context, cmd *commands.DeleteReportFormat) (resp *commands.DeleteReportFormatResponse, err error)
	VerifyReportFormat(ctx context.Context, cmd *commands.VerifyReportFormat) (resp *commands.VerifyReportFormatResponse, err error)
	CreateReportConfig(ctx context.Context, cmd *commands.CreateReportConfig) (resp *commands.CreateReportConfigResponse, err error)
	ModifyReportConfig(ctx context.Context, cmd *commands.ModifyReportConfig) (resp *commands.ModifyReportConfigResponse, err error)
	GetReportConfigs(ctx context.Context, cmd *commands.GetReportConfigs) (resp *commands.GetReportConfigsResponse, err error)
	DeleteReportConfig(ctx context.Context, cmd *commands.DeleteReportConfig) (resp *commands.DeleteReportConfigResponse, err error)
	GetSystemReports(ctx context.Context, cmd *commands.GetSystemReports) (resp *commands.GetSystemReportsResponse, err error)

	// Results & Vulnerabilities
	GetResults(ctx context.Context, args ...filtering.FilterArg) (*commands.GetResultsResponse, error)
	GetVulns(ctx context.Context, args ...filtering.FilterArg) (*commands.GetVulnsResponse, error)

	// Notes, Tags, Filters, Tickets, Overrides
	CreateNote(ctx context.Context, cmd *commands.CreateNote) (resp *commands.CreateNoteResponse, err error)
	ModifyNote(ctx context.Context, cmd *commands.ModifyNote) (resp *commands.ModifyNoteResponse, err error)
	GetNotes(ctx context.Context, cmd *commands.GetNotes) (resp *commands.GetNotesResponse, err error)
	DeleteNote(ctx context.Context, cmd *commands.DeleteNote) (resp *commands.DeleteNoteResponse, err error)
	CreateTag(ctx context.Context, cmd *commands.CreateTag) (resp *commands.CreateTagResponse, err error)
	ModifyTag(ctx context.Context, cmd *commands.ModifyTag) (resp *commands.ModifyTagResponse, err error)
	GetTags(ctx context.Context, args ...filtering.FilterArg) (*commands.GetTagsResponse, error)
	DeleteTag(ctx context.Context, cmd *commands.DeleteTag) (resp *commands.DeleteTagResponse, err error)
	CreateFilter(ctx context.Context, cmd *commands.CreateFilter) (resp *commands.CreateFilterResponse, err error)
	ModifyFilter(ctx context.Context, cmd *commands.ModifyFilter) (resp *commands.ModifyFilterResponse, err error)
	GetFilters(ctx context.Context, args ...filtering.FilterArg) (*commands.GetFiltersResponse, error)
	DeleteFilter(ctx context.Context, cmd *commands.DeleteFilter) (resp *commands.DeleteFilterResponse, err error)
	CreateTicket(ctx context.Context, cmd *commands.CreateTicket) (resp *commands.CreateTicketResponse, err error)
	ModifyTicket(ctx context.Context, cmd *commands.ModifyTicket) (resp *commands.ModifyTicketResponse, err error)
	GetTickets(ctx context.Context, args ...filtering.FilterArg) (*commands.GetTicketsResponse, error)
	DeleteTicket(ctx context.Context, cmd *commands.DeleteTicket) (resp *commands.DeleteTicketResponse, err error)
	CreateOverride(ctx context.Context, cmd *commands.CreateOverride) (resp *commands.CreateOverrideResponse, err error)
	DeleteOverride(ctx context.Context, cmd *commands.DeleteOverride) (resp *commands.DeleteOverrideResponse, err error)
	GetOverrides(ctx context.Context, args ...filtering.FilterArg) (*commands.GetOverridesResponse, error)
	ModifyOverride(ctx context.Context, cmd *commands.ModifyOverride) (resp *commands.ModifyOverrideResponse, err error)

	// Configuration, Preferences, Settings
	GetConfigs(ctx context.Context, args ...filtering.FilterArg) (*commands.GetConfigsResponse, error)
	CreateConfig(ctx context.Context, cmd *commands.CreateConfig) (resp *commands.CreateConfigResponse, err error)
	ModifyConfig(ctx context.Context, cmd *commands.ModifyConfig) (resp *commands.ModifyConfigResponse, err error)
	DeleteConfig(ctx context.Context, cmd *commands.DeleteConfig) (resp *commands.DeleteConfigResponse, err error)
	SyncConfig(ctx context.Context, cmd *commands.SyncConfig) (resp *commands.SyncConfigResponse, err error)
	GetPreferences(ctx context.Context, cmd *commands.GetPreferences) (resp *commands.GetPreferencesResponse, err error)
	ModifySetting(ctx context.Context, cmd *commands.ModifySetting) (resp *commands.ModifySettingResponse, err error)
	GetSettings(ctx context.Context, args ...filtering.FilterArg) (*commands.GetSettingsResponse, error)

	// Port Lists & Ranges
	CreatePortList(ctx context.Context, cmd *commands.CreatePortList) (resp *commands.CreatePortListResponse, err error)
	ModifyPortList(ctx context.Context, cmd *commands.ModifyPortList) (resp *commands.ModifyPortListResponse, err error)
	GetPortLists(ctx context.Context, args ...filtering.FilterArg) (*commands.GetPortListsResponse, error)
	DeletePortList(ctx context.Context, cmd *commands.DeletePortList) (resp *commands.DeletePortListResponse, err error)
	CreatePortRange(ctx context.Context, cmd *commands.CreatePortRange) (resp *commands.CreatePortRangeResponse, err error)
	DeletePortRange(ctx context.Context, cmd *commands.DeletePortRange) (resp *commands.DeletePortRangeResponse, err error)

	// Scanners
	CreateScanner(ctx context.Context, cmd *commands.CreateScanner) (resp *commands.CreateScannerResponse, err error)
	ModifyScanner(ctx context.Context, cmd *commands.ModifyScanner) (resp *commands.ModifyScannerResponse, err error)
	GetScanners(ctx context.Context, args ...filtering.FilterArg) (*commands.GetScannersResponse, error)
	DeleteScanner(ctx context.Context, cmd *commands.DeleteScanner) (resp *commands.DeleteScannerResponse, err error)
	VerifyScanner(ctx context.Context, cmd *commands.VerifyScanner) (resp *commands.VerifyScannerResponse, err error)

	// TLS Certificates
	CreateTLSCertificate(ctx context.Context, cmd *commands.CreateTLSCertificate) (resp *commands.CreateTLSCertificateResponse, err error)
	ModifyTLSCertificate(ctx context.Context, cmd *commands.ModifyTLSCertificate) (resp *commands.ModifyTLSCertificateResponse, err error)
	GetTLSCertificates(ctx context.Context, args ...filtering.FilterArg) (*commands.GetTLSCertificatesResponse, error)

	// Feeds, Licenses, Features, Aggregates, Resource Names
	GetNVTs(ctx context.Context, cmd *commands.GetNVTs) (resp *commands.GetNVTsResponse, err error)
	GetNVTFamilies(ctx context.Context, cmd *commands.GetNVTFamilies) (resp *commands.GetNVTFamiliesResponse, err error)
	GetFeeds(ctx context.Context, cmd *commands.GetFeeds) (resp *commands.GetFeedsResponse, err error)
	ModifyLicense(ctx context.Context, cmd *commands.ModifyLicense) (resp *commands.ModifyLicenseResponse, err error)
	GetLicense(ctx context.Context, cmd *commands.GetLicense) (resp *commands.GetLicenseResponse, err error)
	GetFeatures(ctx context.Context, cmd *commands.GetFeatures) (resp *commands.GetFeaturesResponse, err error)
	GetAggregates(ctx context.Context, args ...filtering.FilterArg) (*commands.GetAggregatesResponse, error)
	GetResourceNames(ctx context.Context, args ...filtering.FilterArg) (*commands.GetResourceNamesResponse, error)

	// Miscellaneous
	GetInfo(ctx context.Context, cmd *commands.GetInfo) (resp *commands.GetInfoResponse, err error)
	GetVersion(ctx context.Context, cmd *commands.GetVersion) (resp *commands.GetVersionResponse, err error)
	Help(ctx context.Context, cmd *commands.Help) (resp *commands.HelpResponse, err error)
	EmptyTrashcan(ctx context.Context, cmd *commands.EmptyTrashcan) (resp *commands.EmptyTrashcanResponse, err error)
	Restore(ctx context.Context, cmd *commands.Restore) (resp *commands.RestoreResponse, err error)
	RunWizard(ctx context.Context, cmd *commands.RunWizard) (resp *commands.RunWizardResponse, err error)
	RawXML(ctx context.Context, xml string) (string, error)

	// Iterator Methods
	Tasks(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.TaskIterator
	Results(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.ResultIterator
	Assets(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.AssetIterator
	Targets(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.TargetIterator
	Tickets(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.TicketIterator
	PortLists(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.PortListIterator
	Settings(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.SettingsIterator
}
