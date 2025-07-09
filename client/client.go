package client

import (
	"github.com/brennoo/go-gmp"
	"github.com/brennoo/go-gmp/commands"
)

// client implements the gmp.Client interface and provides methods to interact with the GMP server.
type client struct {
	conn gmp.Connection
}

// New returns a new instance of a `gmp.Client`. You need to pass a valid `gmp.Connection`.
func New(conn gmp.Connection) gmp.Client {
	cli := &client{
		conn: conn,
	}
	return cli
}

func (cli *client) Authenticate(cmd *commands.AuthenticateCommand) (resp *commands.AuthenticateResponse, err error) {
	resp = new(commands.AuthenticateResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetConfigs(cmd *commands.GetConfigsCommand) (resp *commands.GetConfigsResponse, err error) {
	resp = new(commands.GetConfigsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetScanners(cmd *commands.GetScannersCommand) (resp *commands.GetScannersResponse, err error) {
	resp = new(commands.GetScannersResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetPreferences(cmd *commands.GetPreferencesCommand) (resp *commands.GetPreferencesResponse, err error) {
	resp = new(commands.GetPreferencesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateConfig(cmd *commands.CreateConfigCommand) (resp *commands.CreateConfigResponse, err error) {
	resp = new(commands.CreateConfigResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyConfig(cmd *commands.ModifyConfigCommand) (resp *commands.ModifyConfigResponse, err error) {
	resp = new(commands.ModifyConfigResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateTask(cmd *commands.CreateTaskCommand) (resp *commands.CreateTaskResponse, err error) {
	resp = new(commands.CreateTaskResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyTask(cmd *commands.ModifyTaskCommand) (resp *commands.ModifyTaskResponse, err error) {
	resp = new(commands.ModifyTaskResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateTarget(cmd *commands.CreateTargetCommand) (resp *commands.CreateTargetResponse, err error) {
	resp = new(commands.CreateTargetResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyTarget(cmd *commands.ModifyTargetCommand) (resp *commands.ModifyTargetResponse, err error) {
	resp = new(commands.ModifyTargetResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetTargets(cmd *commands.GetTargetsCommand) (resp *commands.GetTargetsResponse, err error) {
	resp = new(commands.GetTargetsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteTarget(cmd *commands.DeleteTargetCommand) (resp *commands.DeleteTargetResponse, err error) {
	resp = new(commands.DeleteTargetResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetPortLists(cmd *commands.GetPortListsCommand) (resp *commands.GetPortListsResponse, err error) {
	resp = new(commands.GetPortListsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) StartTask(cmd *commands.StartTaskCommand) (resp *commands.StartTaskResponse, err error) {
	resp = new(commands.StartTaskResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetTasks(cmd *commands.GetTasksCommand) (resp *commands.GetTasksResponse, err error) {
	resp = new(commands.GetTasksResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) StopTask(cmd *commands.StopTaskCommand) (resp *commands.StopTaskResponse, err error) {
	resp = new(commands.StopTaskResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteTask(cmd *commands.DeleteTaskCommand) (resp *commands.DeleteTaskResponse, err error) {
	resp = new(commands.DeleteTaskResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetResults(cmd *commands.GetResultsCommand) (resp *commands.GetResultsResponse, err error) {
	resp = new(commands.GetResultsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetVulns(cmd *commands.GetVulnsCommand) (resp *commands.GetVulnsResponse, err error) {
	resp = new(commands.GetVulnsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateAlert(cmd *commands.CreateAlertRequest) (resp *commands.CreateAlertResponse, err error) {
	resp = new(commands.CreateAlertResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetAlerts(cmd *commands.GetAlertsCommand) (resp *commands.GetAlertsResponse, err error) {
	resp = new(commands.GetAlertsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyAlert(cmd *commands.ModifyAlertCommand) (resp *commands.ModifyAlertResponse, err error) {
	resp = new(commands.ModifyAlertResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteAlert(cmd *commands.DeleteAlertCommand) (resp *commands.DeleteAlertResponse, err error) {
	resp = new(commands.DeleteAlertResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) TestAlert(cmd *commands.TestAlertCommand) (resp *commands.TestAlertResponse, err error) {
	resp = new(commands.TestAlertResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ResumeTask(cmd *commands.ResumeTaskCommand) (resp *commands.ResumeTaskResponse, err error) {
	resp = new(commands.ResumeTaskResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateAsset(cmd *commands.CreateAssetCommand) (resp *commands.CreateAssetResponse, err error) {
	resp = new(commands.CreateAssetResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyAsset(cmd *commands.ModifyAssetCommand) (resp *commands.ModifyAssetResponse, err error) {
	resp = new(commands.ModifyAssetResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetAssets(cmd *commands.GetAssetsCommand) (resp *commands.GetAssetsResponse, err error) {
	resp = new(commands.GetAssetsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteAsset(cmd *commands.DeleteAssetCommand) (resp *commands.DeleteAssetResponse, err error) {
	resp = new(commands.DeleteAssetResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateSchedule(cmd *commands.CreateScheduleCommand) (resp *commands.CreateScheduleResponse, err error) {
	resp = new(commands.CreateScheduleResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifySchedule(cmd *commands.ModifyScheduleCommand) (resp *commands.ModifyScheduleResponse, err error) {
	resp = new(commands.ModifyScheduleResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetSchedules(cmd *commands.GetSchedulesCommand) (resp *commands.GetSchedulesResponse, err error) {
	resp = new(commands.GetSchedulesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteSchedule(cmd *commands.DeleteScheduleCommand) (resp *commands.DeleteScheduleResponse, err error) {
	resp = new(commands.DeleteScheduleResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateOverride(cmd *commands.CreateOverrideCommand) (resp *commands.CreateOverrideResponse, err error) {
	resp = new(commands.CreateOverrideResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteOverride(cmd *commands.DeleteOverrideCommand) (resp *commands.DeleteOverrideResponse, err error) {
	resp = new(commands.DeleteOverrideResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetOverrides(cmd *commands.GetOverridesCommand) (resp *commands.GetOverridesResponse, err error) {
	resp = new(commands.GetOverridesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyOverride(cmd *commands.ModifyOverrideCommand) (resp *commands.ModifyOverrideResponse, err error) {
	resp = new(commands.ModifyOverrideResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateReport(cmd *commands.CreateReportCommand) (resp *commands.CreateReportResponse, err error) {
	resp = new(commands.CreateReportResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetReports(cmd *commands.GetReportsCommand) (resp *commands.GetReportsResponse, err error) {
	resp = new(commands.GetReportsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteReport(cmd *commands.DeleteReportCommand) (resp *commands.DeleteReportResponse, err error) {
	resp = new(commands.DeleteReportResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateReportFormat(cmd *commands.CreateReportFormatCommand) (resp *commands.CreateReportFormatResponse, err error) {
	resp = new(commands.CreateReportFormatResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyReportFormat(cmd *commands.ModifyReportFormatCommand) (resp *commands.ModifyReportFormatResponse, err error) {
	resp = new(commands.ModifyReportFormatResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetReportFormats(cmd *commands.GetReportFormatsCommand) (resp *commands.GetReportFormatsResponse, err error) {
	resp = new(commands.GetReportFormatsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteReportFormat(cmd *commands.DeleteReportFormatCommand) (resp *commands.DeleteReportFormatResponse, err error) {
	resp = new(commands.DeleteReportFormatResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) VerifyReportFormat(cmd *commands.VerifyReportFormatCommand) (resp *commands.VerifyReportFormatResponse, err error) {
	resp = new(commands.VerifyReportFormatResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateReportConfig(cmd *commands.CreateReportConfigCommand) (resp *commands.CreateReportConfigResponse, err error) {
	resp = new(commands.CreateReportConfigResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyReportConfig(cmd *commands.ModifyReportConfigCommand) (resp *commands.ModifyReportConfigResponse, err error) {
	resp = new(commands.ModifyReportConfigResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetReportConfigs(cmd *commands.GetReportConfigsCommand) (resp *commands.GetReportConfigsResponse, err error) {
	resp = new(commands.GetReportConfigsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteReportConfig(cmd *commands.DeleteReportConfigCommand) (resp *commands.DeleteReportConfigResponse, err error) {
	resp = new(commands.DeleteReportConfigResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetSystemReports(cmd *commands.GetSystemReportsCommand) (resp *commands.GetSystemReportsResponse, err error) {
	resp = new(commands.GetSystemReportsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateCredential(cmd *commands.CreateCredentialCommand) (resp *commands.CreateCredentialResponse, err error) {
	resp = new(commands.CreateCredentialResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyCredential(cmd *commands.ModifyCredentialCommand) (resp *commands.ModifyCredentialResponse, err error) {
	resp = new(commands.ModifyCredentialResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetCredentials(cmd *commands.GetCredentialsCommand) (resp *commands.GetCredentialsResponse, err error) {
	resp = new(commands.GetCredentialsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteCredential(cmd *commands.DeleteCredentialCommand) (resp *commands.DeleteCredentialResponse, err error) {
	resp = new(commands.DeleteCredentialResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateScanner(cmd *commands.CreateScannerCommand) (resp *commands.CreateScannerResponse, err error) {
	resp = new(commands.CreateScannerResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyScanner(cmd *commands.ModifyScannerCommand) (resp *commands.ModifyScannerResponse, err error) {
	resp = new(commands.ModifyScannerResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteScanner(cmd *commands.DeleteScannerCommand) (resp *commands.DeleteScannerResponse, err error) {
	resp = new(commands.DeleteScannerResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) VerifyScanner(cmd *commands.VerifyScannerCommand) (resp *commands.VerifyScannerResponse, err error) {
	resp = new(commands.VerifyScannerResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreatePortList(cmd *commands.CreatePortListCommand) (resp *commands.CreatePortListResponse, err error) {
	resp = new(commands.CreatePortListResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyPortList(cmd *commands.ModifyPortListCommand) (resp *commands.ModifyPortListResponse, err error) {
	resp = new(commands.ModifyPortListResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeletePortList(cmd *commands.DeletePortListCommand) (resp *commands.DeletePortListResponse, err error) {
	resp = new(commands.DeletePortListResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreatePortRange(cmd *commands.CreatePortRangeCommand) (resp *commands.CreatePortRangeResponse, err error) {
	resp = new(commands.CreatePortRangeResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeletePortRange(cmd *commands.DeletePortRangeCommand) (resp *commands.DeletePortRangeResponse, err error) {
	resp = new(commands.DeletePortRangeResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DescribeAuth(cmd *commands.DescribeAuthCommand) (resp *commands.DescribeAuthResponse, err error) {
	resp = new(commands.DescribeAuthResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyAuth(cmd *commands.ModifyAuthCommand) (resp *commands.ModifyAuthResponse, err error) {
	resp = new(commands.ModifyAuthResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetInfo(cmd *commands.GetInfoCommand) (resp *commands.GetInfoResponse, err error) {
	resp = new(commands.GetInfoResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetVersion(cmd *commands.GetVersionCommand) (resp *commands.GetVersionResponse, err error) {
	resp = new(commands.GetVersionResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) Help(cmd *commands.HelpCommand) (resp *commands.HelpResponse, err error) {
	resp = new(commands.HelpResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteConfig(cmd *commands.DeleteConfigCommand) (resp *commands.DeleteConfigResponse, err error) {
	resp = new(commands.DeleteConfigResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) SyncConfig(cmd *commands.SyncConfigCommand) (resp *commands.SyncConfigResponse, err error) {
	resp = new(commands.SyncConfigResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) MoveTask(cmd *commands.MoveTaskCommand) (resp *commands.MoveTaskResponse, err error) {
	resp = new(commands.MoveTaskResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateTLSCertificate(cmd *commands.CreateTLSCertificateCommand) (resp *commands.CreateTLSCertificateResponse, err error) {
	resp = new(commands.CreateTLSCertificateResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyTLSCertificate(cmd *commands.ModifyTLSCertificateCommand) (resp *commands.ModifyTLSCertificateResponse, err error) {
	resp = new(commands.ModifyTLSCertificateResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetTLSCertificates(cmd *commands.GetTLSCertificatesCommand) (resp *commands.GetTLSCertificatesResponse, err error) {
	resp = new(commands.GetTLSCertificatesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyAgents(cmd *commands.ModifyAgentsCommand) (resp *commands.ModifyAgentsResponse, err error) {
	resp = new(commands.ModifyAgentsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetAgents(cmd *commands.GetAgentsCommand) (resp *commands.GetAgentsResponse, err error) {
	resp = new(commands.GetAgentsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteAgents(cmd *commands.DeleteAgentsCommand) (resp *commands.DeleteAgentsResponse, err error) {
	resp = new(commands.DeleteAgentsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetNvts(cmd *commands.GetNvtsCommand) (resp *commands.GetNvtsResponse, err error) {
	resp = new(commands.GetNvtsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetNvtFamilies(cmd *commands.GetNvtFamiliesCommand) (resp *commands.GetNvtFamiliesResponse, err error) {
	resp = new(commands.GetNvtFamiliesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetFeeds(cmd *commands.GetFeedsCommand) (resp *commands.GetFeedsResponse, err error) {
	resp = new(commands.GetFeedsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyLicense(cmd *commands.ModifyLicenseCommand) (resp *commands.ModifyLicenseResponse, err error) {
	resp = new(commands.ModifyLicenseResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetLicense(cmd *commands.GetLicenseCommand) (resp *commands.GetLicenseResponse, err error) {
	resp = new(commands.GetLicenseResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifySetting(cmd *commands.ModifySettingCommand) (resp *commands.ModifySettingResponse, err error) {
	resp = new(commands.ModifySettingResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetSettings(cmd *commands.GetSettingsCommand) (resp *commands.GetSettingsResponse, err error) {
	resp = new(commands.GetSettingsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetResourceNames(cmd *commands.GetResourceNamesCommand) (resp *commands.GetResourceNamesResponse, err error) {
	resp = new(commands.GetResourceNamesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetAggregates(cmd *commands.GetAggregatesCommand) (resp *commands.GetAggregatesResponse, err error) {
	resp = new(commands.GetAggregatesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetFeatures(cmd *commands.GetFeaturesCommand) (resp *commands.GetFeaturesResponse, err error) {
	resp = new(commands.GetFeaturesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) EmptyTrashcan(cmd *commands.EmptyTrashcanCommand) (resp *commands.EmptyTrashcanResponse, err error) {
	resp = new(commands.EmptyTrashcanResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) Restore(cmd *commands.RestoreCommand) (resp *commands.RestoreResponse, err error) {
	resp = new(commands.RestoreResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) RunWizard(cmd *commands.RunWizardCommand) (resp *commands.RunWizardResponse, err error) {
	resp = new(commands.RunWizardResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) RawXML(xmlStr string) (string, error) {
	return cli.conn.RawXML(xmlStr)
}

func (cli *client) CreateUser(cmd *commands.CreateUserCommand) (resp *commands.CreateUserResponse, err error) {
	resp = new(commands.CreateUserResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyUser(cmd *commands.ModifyUserCommand) (resp *commands.ModifyUserResponse, err error) {
	resp = new(commands.ModifyUserResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetUsers(cmd *commands.GetUsersCommand) (resp *commands.GetUsersResponse, err error) {
	resp = new(commands.GetUsersResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteUser(cmd *commands.DeleteUserCommand) (resp *commands.DeleteUserResponse, err error) {
	resp = new(commands.DeleteUserResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateRole(cmd *commands.CreateRoleCommand) (resp *commands.CreateRoleResponse, err error) {
	resp = new(commands.CreateRoleResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyRole(cmd *commands.ModifyRoleCommand) (resp *commands.ModifyRoleResponse, err error) {
	resp = new(commands.ModifyRoleResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetRoles(cmd *commands.GetRolesCommand) (resp *commands.GetRolesResponse, err error) {
	resp = new(commands.GetRolesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteRole(cmd *commands.DeleteRoleCommand) (resp *commands.DeleteRoleResponse, err error) {
	resp = new(commands.DeleteRoleResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreatePermission(cmd *commands.CreatePermissionCommand) (resp *commands.CreatePermissionResponse, err error) {
	resp = new(commands.CreatePermissionResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyPermission(cmd *commands.ModifyPermissionCommand) (resp *commands.ModifyPermissionResponse, err error) {
	resp = new(commands.ModifyPermissionResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetPermissions(cmd *commands.GetPermissionsCommand) (resp *commands.GetPermissionsResponse, err error) {
	resp = new(commands.GetPermissionsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeletePermission(cmd *commands.DeletePermissionCommand) (resp *commands.DeletePermissionResponse, err error) {
	resp = new(commands.DeletePermissionResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateGroup(cmd *commands.CreateGroupCommand) (resp *commands.CreateGroupResponse, err error) {
	resp = new(commands.CreateGroupResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyGroup(cmd *commands.ModifyGroupCommand) (resp *commands.ModifyGroupResponse, err error) {
	resp = new(commands.ModifyGroupResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetGroups(cmd *commands.GetGroupsCommand) (resp *commands.GetGroupsResponse, err error) {
	resp = new(commands.GetGroupsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteGroup(cmd *commands.DeleteGroupCommand) (resp *commands.DeleteGroupResponse, err error) {
	resp = new(commands.DeleteGroupResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateTag(cmd *commands.CreateTagCommand) (resp *commands.CreateTagResponse, err error) {
	resp = new(commands.CreateTagResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyTag(cmd *commands.ModifyTagCommand) (resp *commands.ModifyTagResponse, err error) {
	resp = new(commands.ModifyTagResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetTags(cmd *commands.GetTagsCommand) (resp *commands.GetTagsResponse, err error) {
	resp = new(commands.GetTagsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteTag(cmd *commands.DeleteTagCommand) (resp *commands.DeleteTagResponse, err error) {
	resp = new(commands.DeleteTagResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateNote(cmd *commands.CreateNoteCommand) (resp *commands.CreateNoteResponse, err error) {
	resp = new(commands.CreateNoteResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyNote(cmd *commands.ModifyNoteCommand) (resp *commands.ModifyNoteResponse, err error) {
	resp = new(commands.ModifyNoteResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetNotes(cmd *commands.GetNotesCommand) (resp *commands.GetNotesResponse, err error) {
	resp = new(commands.GetNotesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteNote(cmd *commands.DeleteNoteCommand) (resp *commands.DeleteNoteResponse, err error) {
	resp = new(commands.DeleteNoteResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateFilter(cmd *commands.CreateFilterCommand) (resp *commands.CreateFilterResponse, err error) {
	resp = new(commands.CreateFilterResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyFilter(cmd *commands.ModifyFilterCommand) (resp *commands.ModifyFilterResponse, err error) {
	resp = new(commands.ModifyFilterResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetFilters(cmd *commands.GetFiltersCommand) (resp *commands.GetFiltersResponse, err error) {
	resp = new(commands.GetFiltersResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteFilter(cmd *commands.DeleteFilterCommand) (resp *commands.DeleteFilterResponse, err error) {
	resp = new(commands.DeleteFilterResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateTicket(cmd *commands.CreateTicketCommand) (resp *commands.CreateTicketResponse, err error) {
	resp = new(commands.CreateTicketResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyTicket(cmd *commands.ModifyTicketCommand) (resp *commands.ModifyTicketResponse, err error) {
	resp = new(commands.ModifyTicketResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetTickets(cmd *commands.GetTicketsCommand) (resp *commands.GetTicketsResponse, err error) {
	resp = new(commands.GetTicketsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteTicket(cmd *commands.DeleteTicketCommand) (resp *commands.DeleteTicketResponse, err error) {
	resp = new(commands.DeleteTicketResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}
