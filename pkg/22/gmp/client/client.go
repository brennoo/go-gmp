package client

import (
	"github.com/brennoo/go-gmp/pkg/22/gmp"
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

func (cli *client) Authenticate(cmd *gmp.AuthenticateCommand) (resp *gmp.AuthenticateResponse, err error) {
	resp = new(gmp.AuthenticateResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetConfigs(cmd *gmp.GetConfigsCommand) (resp *gmp.GetConfigsResponse, err error) {
	resp = new(gmp.GetConfigsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetScanners(cmd *gmp.GetScannersCommand) (resp *gmp.GetScannersResponse, err error) {
	resp = new(gmp.GetScannersResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetPreferences(cmd *gmp.GetPreferencesCommand) (resp *gmp.GetPreferencesResponse, err error) {
	resp = new(gmp.GetPreferencesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateConfig(cmd *gmp.CreateConfigCommand) (resp *gmp.CreateConfigResponse, err error) {
	resp = new(gmp.CreateConfigResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyConfig(cmd *gmp.ModifyConfigCommand) (resp *gmp.ModifyConfigResponse, err error) {
	resp = new(gmp.ModifyConfigResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateTask(cmd *gmp.CreateTaskCommand) (resp *gmp.CreateTaskResponse, err error) {
	resp = new(gmp.CreateTaskResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyTask(cmd *gmp.ModifyTaskCommand) (resp *gmp.ModifyTaskResponse, err error) {
	resp = new(gmp.ModifyTaskResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateTarget(cmd *gmp.CreateTargetCommand) (resp *gmp.CreateTargetResponse, err error) {
	resp = new(gmp.CreateTargetResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyTarget(cmd *gmp.ModifyTargetCommand) (resp *gmp.ModifyTargetResponse, err error) {
	resp = new(gmp.ModifyTargetResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetTargets(cmd *gmp.GetTargetsCommand) (resp *gmp.GetTargetsResponse, err error) {
	resp = new(gmp.GetTargetsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteTarget(cmd *gmp.DeleteTargetCommand) (resp *gmp.DeleteTargetResponse, err error) {
	resp = new(gmp.DeleteTargetResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetPortLists(cmd *gmp.GetPortListsCommand) (resp *gmp.GetPortListsResponse, err error) {
	resp = new(gmp.GetPortListsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) StartTask(cmd *gmp.StartTaskCommand) (resp *gmp.StartTaskResponse, err error) {
	resp = new(gmp.StartTaskResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetTasks(cmd *gmp.GetTasksCommand) (resp *gmp.GetTasksResponse, err error) {
	resp = new(gmp.GetTasksResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) StopTask(cmd *gmp.StopTaskCommand) (resp *gmp.StopTaskResponse, err error) {
	resp = new(gmp.StopTaskResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteTask(cmd *gmp.DeleteTaskCommand) (resp *gmp.DeleteTaskResponse, err error) {
	resp = new(gmp.DeleteTaskResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetResults(cmd *gmp.GetResultsCommand) (resp *gmp.GetResultsResponse, err error) {
	resp = new(gmp.GetResultsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetVulns(cmd *gmp.GetVulnsCommand) (resp *gmp.GetVulnsResponse, err error) {
	resp = new(gmp.GetVulnsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateAlert(cmd *gmp.CreateAlertRequest) (resp *gmp.CreateAlertResponse, err error) {
	resp = new(gmp.CreateAlertResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetAlerts(cmd *gmp.GetAlertsCommand) (resp *gmp.GetAlertsResponse, err error) {
	resp = new(gmp.GetAlertsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyAlert(cmd *gmp.ModifyAlertCommand) (resp *gmp.ModifyAlertResponse, err error) {
	resp = new(gmp.ModifyAlertResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteAlert(cmd *gmp.DeleteAlertCommand) (resp *gmp.DeleteAlertResponse, err error) {
	resp = new(gmp.DeleteAlertResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) TestAlert(cmd *gmp.TestAlertCommand) (resp *gmp.TestAlertResponse, err error) {
	resp = new(gmp.TestAlertResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ResumeTask(cmd *gmp.ResumeTaskCommand) (resp *gmp.ResumeTaskResponse, err error) {
	resp = new(gmp.ResumeTaskResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateAsset(cmd *gmp.CreateAssetCommand) (resp *gmp.CreateAssetResponse, err error) {
	resp = new(gmp.CreateAssetResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyAsset(cmd *gmp.ModifyAssetCommand) (resp *gmp.ModifyAssetResponse, err error) {
	resp = new(gmp.ModifyAssetResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetAssets(cmd *gmp.GetAssetsCommand) (resp *gmp.GetAssetsResponse, err error) {
	resp = new(gmp.GetAssetsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteAsset(cmd *gmp.DeleteAssetCommand) (resp *gmp.DeleteAssetResponse, err error) {
	resp = new(gmp.DeleteAssetResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateSchedule(cmd *gmp.CreateScheduleCommand) (resp *gmp.CreateScheduleResponse, err error) {
	resp = new(gmp.CreateScheduleResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifySchedule(cmd *gmp.ModifyScheduleCommand) (resp *gmp.ModifyScheduleResponse, err error) {
	resp = new(gmp.ModifyScheduleResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetSchedules(cmd *gmp.GetSchedulesCommand) (resp *gmp.GetSchedulesResponse, err error) {
	resp = new(gmp.GetSchedulesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteSchedule(cmd *gmp.DeleteScheduleCommand) (resp *gmp.DeleteScheduleResponse, err error) {
	resp = new(gmp.DeleteScheduleResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateOverride(cmd *gmp.CreateOverrideCommand) (resp *gmp.CreateOverrideResponse, err error) {
	resp = new(gmp.CreateOverrideResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteOverride(cmd *gmp.DeleteOverrideCommand) (resp *gmp.DeleteOverrideResponse, err error) {
	resp = new(gmp.DeleteOverrideResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetOverrides(cmd *gmp.GetOverridesCommand) (resp *gmp.GetOverridesResponse, err error) {
	resp = new(gmp.GetOverridesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyOverride(cmd *gmp.ModifyOverrideCommand) (resp *gmp.ModifyOverrideResponse, err error) {
	resp = new(gmp.ModifyOverrideResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateReport(cmd *gmp.CreateReportCommand) (resp *gmp.CreateReportResponse, err error) {
	resp = new(gmp.CreateReportResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetReports(cmd *gmp.GetReportsCommand) (resp *gmp.GetReportsResponse, err error) {
	resp = new(gmp.GetReportsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteReport(cmd *gmp.DeleteReportCommand) (resp *gmp.DeleteReportResponse, err error) {
	resp = new(gmp.DeleteReportResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateReportFormat(cmd *gmp.CreateReportFormatCommand) (resp *gmp.CreateReportFormatResponse, err error) {
	resp = new(gmp.CreateReportFormatResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyReportFormat(cmd *gmp.ModifyReportFormatCommand) (resp *gmp.ModifyReportFormatResponse, err error) {
	resp = new(gmp.ModifyReportFormatResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetReportFormats(cmd *gmp.GetReportFormatsCommand) (resp *gmp.GetReportFormatsResponse, err error) {
	resp = new(gmp.GetReportFormatsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteReportFormat(cmd *gmp.DeleteReportFormatCommand) (resp *gmp.DeleteReportFormatResponse, err error) {
	resp = new(gmp.DeleteReportFormatResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) VerifyReportFormat(cmd *gmp.VerifyReportFormatCommand) (resp *gmp.VerifyReportFormatResponse, err error) {
	resp = new(gmp.VerifyReportFormatResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateReportConfig(cmd *gmp.CreateReportConfigCommand) (resp *gmp.CreateReportConfigResponse, err error) {
	resp = new(gmp.CreateReportConfigResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyReportConfig(cmd *gmp.ModifyReportConfigCommand) (resp *gmp.ModifyReportConfigResponse, err error) {
	resp = new(gmp.ModifyReportConfigResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetReportConfigs(cmd *gmp.GetReportConfigsCommand) (resp *gmp.GetReportConfigsResponse, err error) {
	resp = new(gmp.GetReportConfigsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteReportConfig(cmd *gmp.DeleteReportConfigCommand) (resp *gmp.DeleteReportConfigResponse, err error) {
	resp = new(gmp.DeleteReportConfigResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetSystemReports(cmd *gmp.GetSystemReportsCommand) (resp *gmp.GetSystemReportsResponse, err error) {
	resp = new(gmp.GetSystemReportsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateCredential(cmd *gmp.CreateCredentialCommand) (resp *gmp.CreateCredentialResponse, err error) {
	resp = new(gmp.CreateCredentialResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyCredential(cmd *gmp.ModifyCredentialCommand) (resp *gmp.ModifyCredentialResponse, err error) {
	resp = new(gmp.ModifyCredentialResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetCredentials(cmd *gmp.GetCredentialsCommand) (resp *gmp.GetCredentialsResponse, err error) {
	resp = new(gmp.GetCredentialsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteCredential(cmd *gmp.DeleteCredentialCommand) (resp *gmp.DeleteCredentialResponse, err error) {
	resp = new(gmp.DeleteCredentialResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateScanner(cmd *gmp.CreateScannerCommand) (resp *gmp.CreateScannerResponse, err error) {
	resp = new(gmp.CreateScannerResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyScanner(cmd *gmp.ModifyScannerCommand) (resp *gmp.ModifyScannerResponse, err error) {
	resp = new(gmp.ModifyScannerResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteScanner(cmd *gmp.DeleteScannerCommand) (resp *gmp.DeleteScannerResponse, err error) {
	resp = new(gmp.DeleteScannerResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) VerifyScanner(cmd *gmp.VerifyScannerCommand) (resp *gmp.VerifyScannerResponse, err error) {
	resp = new(gmp.VerifyScannerResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreatePortList(cmd *gmp.CreatePortListCommand) (resp *gmp.CreatePortListResponse, err error) {
	resp = new(gmp.CreatePortListResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyPortList(cmd *gmp.ModifyPortListCommand) (resp *gmp.ModifyPortListResponse, err error) {
	resp = new(gmp.ModifyPortListResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeletePortList(cmd *gmp.DeletePortListCommand) (resp *gmp.DeletePortListResponse, err error) {
	resp = new(gmp.DeletePortListResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreatePortRange(cmd *gmp.CreatePortRangeCommand) (resp *gmp.CreatePortRangeResponse, err error) {
	resp = new(gmp.CreatePortRangeResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeletePortRange(cmd *gmp.DeletePortRangeCommand) (resp *gmp.DeletePortRangeResponse, err error) {
	resp = new(gmp.DeletePortRangeResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DescribeAuth(cmd *gmp.DescribeAuthCommand) (resp *gmp.DescribeAuthResponse, err error) {
	resp = new(gmp.DescribeAuthResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetInfo(cmd *gmp.GetInfoCommand) (resp *gmp.GetInfoResponse, err error) {
	resp = new(gmp.GetInfoResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetVersion(cmd *gmp.GetVersionCommand) (resp *gmp.GetVersionResponse, err error) {
	resp = new(gmp.GetVersionResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) Help(cmd *gmp.HelpCommand) (resp *gmp.HelpResponse, err error) {
	resp = new(gmp.HelpResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteConfig(cmd *gmp.DeleteConfigCommand) (resp *gmp.DeleteConfigResponse, err error) {
	resp = new(gmp.DeleteConfigResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) SyncConfig(cmd *gmp.SyncConfigCommand) (resp *gmp.SyncConfigResponse, err error) {
	resp = new(gmp.SyncConfigResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) MoveTask(cmd *gmp.MoveTaskCommand) (resp *gmp.MoveTaskResponse, err error) {
	resp = new(gmp.MoveTaskResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateTLSCertificate(cmd *gmp.CreateTLSCertificateCommand) (resp *gmp.CreateTLSCertificateResponse, err error) {
	resp = new(gmp.CreateTLSCertificateResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyTLSCertificate(cmd *gmp.ModifyTLSCertificateCommand) (resp *gmp.ModifyTLSCertificateResponse, err error) {
	resp = new(gmp.ModifyTLSCertificateResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetTLSCertificates(cmd *gmp.GetTLSCertificatesCommand) (resp *gmp.GetTLSCertificatesResponse, err error) {
	resp = new(gmp.GetTLSCertificatesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyAgents(cmd *gmp.ModifyAgentsCommand) (resp *gmp.ModifyAgentsResponse, err error) {
	resp = new(gmp.ModifyAgentsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetAgents(cmd *gmp.GetAgentsCommand) (resp *gmp.GetAgentsResponse, err error) {
	resp = new(gmp.GetAgentsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteAgents(cmd *gmp.DeleteAgentsCommand) (resp *gmp.DeleteAgentsResponse, err error) {
	resp = new(gmp.DeleteAgentsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetNvts(cmd *gmp.GetNvtsCommand) (resp *gmp.GetNvtsResponse, err error) {
	resp = new(gmp.GetNvtsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetNvtFamilies(cmd *gmp.GetNvtFamiliesCommand) (resp *gmp.GetNvtFamiliesResponse, err error) {
	resp = new(gmp.GetNvtFamiliesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetFeeds(cmd *gmp.GetFeedsCommand) (resp *gmp.GetFeedsResponse, err error) {
	resp = new(gmp.GetFeedsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyLicense(cmd *gmp.ModifyLicenseCommand) (resp *gmp.ModifyLicenseResponse, err error) {
	resp = new(gmp.ModifyLicenseResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetLicense(cmd *gmp.GetLicenseCommand) (resp *gmp.GetLicenseResponse, err error) {
	resp = new(gmp.GetLicenseResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifySetting(cmd *gmp.ModifySettingCommand) (resp *gmp.ModifySettingResponse, err error) {
	resp = new(gmp.ModifySettingResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetSettings(cmd *gmp.GetSettingsCommand) (resp *gmp.GetSettingsResponse, err error) {
	resp = new(gmp.GetSettingsResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetResourceNames(cmd *gmp.GetResourceNamesCommand) (resp *gmp.GetResourceNamesResponse, err error) {
	resp = new(gmp.GetResourceNamesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetAggregates(cmd *gmp.GetAggregatesCommand) (resp *gmp.GetAggregatesResponse, err error) {
	resp = new(gmp.GetAggregatesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetFeatures(cmd *gmp.GetFeaturesCommand) (resp *gmp.GetFeaturesResponse, err error) {
	resp = new(gmp.GetFeaturesResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) EmptyTrashcan(cmd *gmp.EmptyTrashcanCommand) (resp *gmp.EmptyTrashcanResponse, err error) {
	resp = new(gmp.EmptyTrashcanResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) Restore(cmd *gmp.RestoreCommand) (resp *gmp.RestoreResponse, err error) {
	resp = new(gmp.RestoreResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) RunWizard(cmd *gmp.RunWizardCommand) (resp *gmp.RunWizardResponse, err error) {
	resp = new(gmp.RunWizardResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) RawXML(xmlStr string) (string, error) {
	return cli.conn.RawXML(xmlStr)
}

func (cli *client) CreateUser(cmd *gmp.CreateUserCommand) (resp *gmp.CreateUserResponse, err error) {
	resp = new(gmp.CreateUserResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyUser(cmd *gmp.ModifyUserCommand) (resp *gmp.ModifyUserResponse, err error) {
	resp = new(gmp.ModifyUserResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) GetUsers(cmd *gmp.GetUsersCommand) (resp *gmp.GetUsersResponse, err error) {
	resp = new(gmp.GetUsersResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) DeleteUser(cmd *gmp.DeleteUserCommand) (resp *gmp.DeleteUserResponse, err error) {
	resp = new(gmp.DeleteUserResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) CreateRole(cmd *gmp.CreateRoleCommand) (resp *gmp.CreateRoleResponse, err error) {
	resp = new(gmp.CreateRoleResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}

func (cli *client) ModifyRole(cmd *gmp.ModifyRoleCommand) (resp *gmp.ModifyRoleResponse, err error) {
	resp = new(gmp.ModifyRoleResponse)
	err = cli.conn.Execute(cmd, resp)
	return resp, err
}
