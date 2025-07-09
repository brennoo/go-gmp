package gmp

import "github.com/brennoo/go-gmp/commands"

type Client interface {
	// Authentication
	Authenticate(cmd *commands.AuthenticateCommand) (resp *commands.AuthenticateResponse, err error)
	DescribeAuth(cmd *commands.DescribeAuthCommand) (resp *commands.DescribeAuthResponse, err error)
	ModifyAuth(cmd *commands.ModifyAuthCommand) (resp *commands.ModifyAuthResponse, err error)

	// User Management
	CreateUser(cmd *commands.CreateUserCommand) (resp *commands.CreateUserResponse, err error)
	ModifyUser(cmd *commands.ModifyUserCommand) (resp *commands.ModifyUserResponse, err error)
	GetUsers(cmd *commands.GetUsersCommand) (resp *commands.GetUsersResponse, err error)
	DeleteUser(cmd *commands.DeleteUserCommand) (resp *commands.DeleteUserResponse, err error)
	CreateRole(cmd *commands.CreateRoleCommand) (resp *commands.CreateRoleResponse, err error)
	ModifyRole(cmd *commands.ModifyRoleCommand) (resp *commands.ModifyRoleResponse, err error)
	GetRoles(cmd *commands.GetRolesCommand) (resp *commands.GetRolesResponse, err error)
	DeleteRole(cmd *commands.DeleteRoleCommand) (resp *commands.DeleteRoleResponse, err error)
	CreatePermission(cmd *commands.CreatePermissionCommand) (resp *commands.CreatePermissionResponse, err error)
	ModifyPermission(cmd *commands.ModifyPermissionCommand) (resp *commands.ModifyPermissionResponse, err error)
	GetPermissions(cmd *commands.GetPermissionsCommand) (resp *commands.GetPermissionsResponse, err error)
	DeletePermission(cmd *commands.DeletePermissionCommand) (resp *commands.DeletePermissionResponse, err error)
	CreateGroup(cmd *commands.CreateGroupCommand) (resp *commands.CreateGroupResponse, err error)
	ModifyGroup(cmd *commands.ModifyGroupCommand) (resp *commands.ModifyGroupResponse, err error)
	GetGroups(cmd *commands.GetGroupsCommand) (resp *commands.GetGroupsResponse, err error)
	DeleteGroup(cmd *commands.DeleteGroupCommand) (resp *commands.DeleteGroupResponse, err error)

	// Asset Management
	CreateAsset(cmd *commands.CreateAssetCommand) (resp *commands.CreateAssetResponse, err error)
	ModifyAsset(cmd *commands.ModifyAssetCommand) (resp *commands.ModifyAssetResponse, err error)
	GetAssets(cmd *commands.GetAssetsCommand) (resp *commands.GetAssetsResponse, err error)
	DeleteAsset(cmd *commands.DeleteAssetCommand) (resp *commands.DeleteAssetResponse, err error)

	// Credential Management
	CreateCredential(cmd *commands.CreateCredentialCommand) (resp *commands.CreateCredentialResponse, err error)
	ModifyCredential(cmd *commands.ModifyCredentialCommand) (resp *commands.ModifyCredentialResponse, err error)
	GetCredentials(cmd *commands.GetCredentialsCommand) (resp *commands.GetCredentialsResponse, err error)
	DeleteCredential(cmd *commands.DeleteCredentialCommand) (resp *commands.DeleteCredentialResponse, err error)

	// Target Management
	CreateTarget(cmd *commands.CreateTargetCommand) (resp *commands.CreateTargetResponse, err error)
	ModifyTarget(cmd *commands.ModifyTargetCommand) (resp *commands.ModifyTargetResponse, err error)
	GetTargets(cmd *commands.GetTargetsCommand) (resp *commands.GetTargetsResponse, err error)
	DeleteTarget(cmd *commands.DeleteTargetCommand) (resp *commands.DeleteTargetResponse, err error)

	// Task Management
	CreateTask(cmd *commands.CreateTaskCommand) (resp *commands.CreateTaskResponse, err error)
	ModifyTask(cmd *commands.ModifyTaskCommand) (resp *commands.ModifyTaskResponse, err error)
	GetTasks(cmd *commands.GetTasksCommand) (resp *commands.GetTasksResponse, err error)
	DeleteTask(cmd *commands.DeleteTaskCommand) (resp *commands.DeleteTaskResponse, err error)
	StartTask(cmd *commands.StartTaskCommand) (resp *commands.StartTaskResponse, err error)
	StopTask(cmd *commands.StopTaskCommand) (resp *commands.StopTaskResponse, err error)
	ResumeTask(cmd *commands.ResumeTaskCommand) (resp *commands.ResumeTaskResponse, err error)
	MoveTask(cmd *commands.MoveTaskCommand) (resp *commands.MoveTaskResponse, err error)

	// Schedule Management
	CreateSchedule(cmd *commands.CreateScheduleCommand) (resp *commands.CreateScheduleResponse, err error)
	ModifySchedule(cmd *commands.ModifyScheduleCommand) (resp *commands.ModifyScheduleResponse, err error)
	GetSchedules(cmd *commands.GetSchedulesCommand) (resp *commands.GetSchedulesResponse, err error)
	DeleteSchedule(cmd *commands.DeleteScheduleCommand) (resp *commands.DeleteScheduleResponse, err error)

	// Alert Management
	CreateAlert(cmd *commands.CreateAlertRequest) (resp *commands.CreateAlertResponse, err error)
	GetAlerts(cmd *commands.GetAlertsCommand) (resp *commands.GetAlertsResponse, err error)
	ModifyAlert(cmd *commands.ModifyAlertCommand) (resp *commands.ModifyAlertResponse, err error)
	DeleteAlert(cmd *commands.DeleteAlertCommand) (resp *commands.DeleteAlertResponse, err error)
	TestAlert(cmd *commands.TestAlertCommand) (resp *commands.TestAlertResponse, err error)

	// Report Management
	CreateReport(cmd *commands.CreateReportCommand) (resp *commands.CreateReportResponse, err error)
	GetReports(cmd *commands.GetReportsCommand) (resp *commands.GetReportsResponse, err error)
	DeleteReport(cmd *commands.DeleteReportCommand) (resp *commands.DeleteReportResponse, err error)
	CreateReportFormat(cmd *commands.CreateReportFormatCommand) (resp *commands.CreateReportFormatResponse, err error)
	ModifyReportFormat(cmd *commands.ModifyReportFormatCommand) (resp *commands.ModifyReportFormatResponse, err error)
	GetReportFormats(cmd *commands.GetReportFormatsCommand) (resp *commands.GetReportFormatsResponse, err error)
	DeleteReportFormat(cmd *commands.DeleteReportFormatCommand) (resp *commands.DeleteReportFormatResponse, err error)
	VerifyReportFormat(cmd *commands.VerifyReportFormatCommand) (resp *commands.VerifyReportFormatResponse, err error)
	CreateReportConfig(cmd *commands.CreateReportConfigCommand) (resp *commands.CreateReportConfigResponse, err error)
	ModifyReportConfig(cmd *commands.ModifyReportConfigCommand) (resp *commands.ModifyReportConfigResponse, err error)
	GetReportConfigs(cmd *commands.GetReportConfigsCommand) (resp *commands.GetReportConfigsResponse, err error)
	DeleteReportConfig(cmd *commands.DeleteReportConfigCommand) (resp *commands.DeleteReportConfigResponse, err error)
	GetSystemReports(cmd *commands.GetSystemReportsCommand) (resp *commands.GetSystemReportsResponse, err error)

	// Results & Vulnerabilities
	GetResults(cmd *commands.GetResultsCommand) (resp *commands.GetResultsResponse, err error)
	GetVulns(cmd *commands.GetVulnsCommand) (resp *commands.GetVulnsResponse, err error)

	// Notes, Tags, Filters, Tickets, Overrides
	CreateNote(cmd *commands.CreateNoteCommand) (resp *commands.CreateNoteResponse, err error)
	ModifyNote(cmd *commands.ModifyNoteCommand) (resp *commands.ModifyNoteResponse, err error)
	GetNotes(cmd *commands.GetNotesCommand) (resp *commands.GetNotesResponse, err error)
	DeleteNote(cmd *commands.DeleteNoteCommand) (resp *commands.DeleteNoteResponse, err error)
	CreateTag(cmd *commands.CreateTagCommand) (resp *commands.CreateTagResponse, err error)
	ModifyTag(cmd *commands.ModifyTagCommand) (resp *commands.ModifyTagResponse, err error)
	GetTags(cmd *commands.GetTagsCommand) (resp *commands.GetTagsResponse, err error)
	DeleteTag(cmd *commands.DeleteTagCommand) (resp *commands.DeleteTagResponse, err error)
	CreateFilter(cmd *commands.CreateFilterCommand) (resp *commands.CreateFilterResponse, err error)
	ModifyFilter(cmd *commands.ModifyFilterCommand) (resp *commands.ModifyFilterResponse, err error)
	GetFilters(cmd *commands.GetFiltersCommand) (resp *commands.GetFiltersResponse, err error)
	DeleteFilter(cmd *commands.DeleteFilterCommand) (resp *commands.DeleteFilterResponse, err error)
	CreateTicket(cmd *commands.CreateTicketCommand) (resp *commands.CreateTicketResponse, err error)
	ModifyTicket(cmd *commands.ModifyTicketCommand) (resp *commands.ModifyTicketResponse, err error)
	GetTickets(cmd *commands.GetTicketsCommand) (resp *commands.GetTicketsResponse, err error)
	DeleteTicket(cmd *commands.DeleteTicketCommand) (resp *commands.DeleteTicketResponse, err error)
	CreateOverride(cmd *commands.CreateOverrideCommand) (resp *commands.CreateOverrideResponse, err error)
	DeleteOverride(cmd *commands.DeleteOverrideCommand) (resp *commands.DeleteOverrideResponse, err error)
	GetOverrides(cmd *commands.GetOverridesCommand) (resp *commands.GetOverridesResponse, err error)
	ModifyOverride(cmd *commands.ModifyOverrideCommand) (resp *commands.ModifyOverrideResponse, err error)

	// Configuration, Preferences, Settings
	GetConfigs(cmd *commands.GetConfigsCommand) (resp *commands.GetConfigsResponse, err error)
	CreateConfig(cmd *commands.CreateConfigCommand) (resp *commands.CreateConfigResponse, err error)
	ModifyConfig(cmd *commands.ModifyConfigCommand) (resp *commands.ModifyConfigResponse, err error)
	DeleteConfig(cmd *commands.DeleteConfigCommand) (resp *commands.DeleteConfigResponse, err error)
	SyncConfig(cmd *commands.SyncConfigCommand) (resp *commands.SyncConfigResponse, err error)
	GetPreferences(cmd *commands.GetPreferencesCommand) (resp *commands.GetPreferencesResponse, err error)
	ModifySetting(cmd *commands.ModifySettingCommand) (resp *commands.ModifySettingResponse, err error)
	GetSettings(cmd *commands.GetSettingsCommand) (resp *commands.GetSettingsResponse, err error)

	// Port Lists & Ranges
	CreatePortList(cmd *commands.CreatePortListCommand) (resp *commands.CreatePortListResponse, err error)
	ModifyPortList(cmd *commands.ModifyPortListCommand) (resp *commands.ModifyPortListResponse, err error)
	GetPortLists(cmd *commands.GetPortListsCommand) (resp *commands.GetPortListsResponse, err error)
	DeletePortList(cmd *commands.DeletePortListCommand) (resp *commands.DeletePortListResponse, err error)
	CreatePortRange(cmd *commands.CreatePortRangeCommand) (resp *commands.CreatePortRangeResponse, err error)
	DeletePortRange(cmd *commands.DeletePortRangeCommand) (resp *commands.DeletePortRangeResponse, err error)

	// Scanners
	CreateScanner(cmd *commands.CreateScannerCommand) (resp *commands.CreateScannerResponse, err error)
	ModifyScanner(cmd *commands.ModifyScannerCommand) (resp *commands.ModifyScannerResponse, err error)
	GetScanners(cmd *commands.GetScannersCommand) (resp *commands.GetScannersResponse, err error)
	DeleteScanner(cmd *commands.DeleteScannerCommand) (resp *commands.DeleteScannerResponse, err error)
	VerifyScanner(cmd *commands.VerifyScannerCommand) (resp *commands.VerifyScannerResponse, err error)

	// TLS Certificates
	CreateTLSCertificate(cmd *commands.CreateTLSCertificateCommand) (resp *commands.CreateTLSCertificateResponse, err error)
	ModifyTLSCertificate(cmd *commands.ModifyTLSCertificateCommand) (resp *commands.ModifyTLSCertificateResponse, err error)
	GetTLSCertificates(cmd *commands.GetTLSCertificatesCommand) (resp *commands.GetTLSCertificatesResponse, err error)

	// Agents
	ModifyAgents(cmd *commands.ModifyAgentsCommand) (resp *commands.ModifyAgentsResponse, err error)
	GetAgents(cmd *commands.GetAgentsCommand) (resp *commands.GetAgentsResponse, err error)
	DeleteAgents(cmd *commands.DeleteAgentsCommand) (resp *commands.DeleteAgentsResponse, err error)

	// Feeds, Licenses, Features, Aggregates, Resource Names
	GetNvts(cmd *commands.GetNvtsCommand) (resp *commands.GetNvtsResponse, err error)
	GetNvtFamilies(cmd *commands.GetNvtFamiliesCommand) (resp *commands.GetNvtFamiliesResponse, err error)
	GetFeeds(cmd *commands.GetFeedsCommand) (resp *commands.GetFeedsResponse, err error)
	ModifyLicense(cmd *commands.ModifyLicenseCommand) (resp *commands.ModifyLicenseResponse, err error)
	GetLicense(cmd *commands.GetLicenseCommand) (resp *commands.GetLicenseResponse, err error)
	GetFeatures(cmd *commands.GetFeaturesCommand) (resp *commands.GetFeaturesResponse, err error)
	GetAggregates(cmd *commands.GetAggregatesCommand) (resp *commands.GetAggregatesResponse, err error)
	GetResourceNames(cmd *commands.GetResourceNamesCommand) (resp *commands.GetResourceNamesResponse, err error)

	// Miscellaneous
	GetInfo(cmd *commands.GetInfoCommand) (resp *commands.GetInfoResponse, err error)
	GetVersion(cmd *commands.GetVersionCommand) (resp *commands.GetVersionResponse, err error)
	Help(cmd *commands.HelpCommand) (resp *commands.HelpResponse, err error)
	EmptyTrashcan(cmd *commands.EmptyTrashcanCommand) (resp *commands.EmptyTrashcanResponse, err error)
	Restore(cmd *commands.RestoreCommand) (resp *commands.RestoreResponse, err error)
	RunWizard(cmd *commands.RunWizardCommand) (resp *commands.RunWizardResponse, err error)
	RawXML(xml string) (string, error)
}
