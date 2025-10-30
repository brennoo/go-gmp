package client

import (
	"context"

	"github.com/brennoo/go-gmp"
	"github.com/brennoo/go-gmp/commands"
	"github.com/brennoo/go-gmp/commands/filtering"
	"github.com/brennoo/go-gmp/commands/pagination"
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

func (cli *client) Authenticate(cmd *commands.Authenticate) (resp *commands.AuthenticateResponse, err error) {
	resp = new(commands.AuthenticateResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetConfigsRaw provides the low-level method for direct command execution.
func (cli *client) GetConfigsRaw(cmd *commands.GetConfigs) (resp *commands.GetConfigsResponse, err error) {
	resp = new(commands.GetConfigsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetConfigs retrieves configs with optional filters.
func (cli *client) GetConfigs(ctx context.Context, args ...filtering.FilterArg) (*commands.GetConfigsResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetConfigs{
		Filter: filter,
	}
	resp := new(commands.GetConfigsResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

// GetScannersRaw provides the low-level method for direct command execution.
func (cli *client) GetScannersRaw(cmd *commands.GetScanners) (resp *commands.GetScannersResponse, err error) {
	resp = new(commands.GetScannersResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetScanners retrieves scanners with optional filters.
func (cli *client) GetScanners(ctx context.Context, args ...filtering.FilterArg) (*commands.GetScannersResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetScanners{
		Filter: filter,
	}
	resp := new(commands.GetScannersResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (cli *client) GetPreferences(cmd *commands.GetPreferences) (resp *commands.GetPreferencesResponse, err error) {
	resp = new(commands.GetPreferencesResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateConfig(cmd *commands.CreateConfig) (resp *commands.CreateConfigResponse, err error) {
	resp = new(commands.CreateConfigResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyConfig(cmd *commands.ModifyConfig) (resp *commands.ModifyConfigResponse, err error) {
	resp = new(commands.ModifyConfigResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateTask(cmd *commands.CreateTask) (resp *commands.CreateTaskResponse, err error) {
	resp = new(commands.CreateTaskResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyTask(cmd *commands.ModifyTask) (resp *commands.ModifyTaskResponse, err error) {
	resp = new(commands.ModifyTaskResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateTarget(cmd *commands.CreateTarget) (resp *commands.CreateTargetResponse, err error) {
	resp = new(commands.CreateTargetResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyTarget(cmd *commands.ModifyTarget) (resp *commands.ModifyTargetResponse, err error) {
	resp = new(commands.ModifyTargetResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetTargetsRaw provides the low-level method for direct command execution.
func (cli *client) GetTargetsRaw(cmd *commands.GetTargets) (resp *commands.GetTargetsResponse, err error) {
	resp = new(commands.GetTargetsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetTargets retrieves targets with optional filters.
func (cli *client) GetTargets(ctx context.Context, args ...filtering.FilterArg) (*commands.GetTargetsResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetTargets{
		Filter: filter,
	}
	resp := new(commands.GetTargetsResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (cli *client) DeleteTarget(cmd *commands.DeleteTarget) (resp *commands.DeleteTargetResponse, err error) {
	resp = new(commands.DeleteTargetResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetPortListsRaw provides the low-level method for direct command execution.
func (cli *client) GetPortListsRaw(cmd *commands.GetPortLists) (resp *commands.GetPortListsResponse, err error) {
	resp = new(commands.GetPortListsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetPortLists retrieves port lists with optional filters.
func (cli *client) GetPortLists(ctx context.Context, args ...filtering.FilterArg) (*commands.GetPortListsResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetPortLists{
		Filter: filter,
	}
	resp := new(commands.GetPortListsResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (cli *client) StartTask(cmd *commands.StartTask) (resp *commands.StartTaskResponse, err error) {
	resp = new(commands.StartTaskResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetTasksRaw provides the low-level method for direct command execution.
func (cli *client) GetTasksRaw(cmd *commands.GetTasks) (resp *commands.GetTasksResponse, err error) {
	resp = new(commands.GetTasksResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetTasks retrieves tasks with optional filters.
func (cli *client) GetTasks(ctx context.Context, args ...filtering.FilterArg) (*commands.GetTasksResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetTasks{
		Filter: filter,
	}
	resp := new(commands.GetTasksResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (cli *client) StopTask(cmd *commands.StopTask) (resp *commands.StopTaskResponse, err error) {
	resp = new(commands.StopTaskResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeleteTask(cmd *commands.DeleteTask) (resp *commands.DeleteTaskResponse, err error) {
	resp = new(commands.DeleteTaskResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetResultsRaw provides the low-level method for direct command execution.
func (cli *client) GetResultsRaw(cmd *commands.GetResults) (resp *commands.GetResultsResponse, err error) {
	resp = new(commands.GetResultsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetResults retrieves results with optional filters.
func (cli *client) GetResults(ctx context.Context, args ...filtering.FilterArg) (*commands.GetResultsResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetResults{
		Filter: filter,
	}
	resp := new(commands.GetResultsResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

// GetVulnsRaw provides the low-level method for direct command execution.
func (cli *client) GetVulnsRaw(cmd *commands.GetVulns) (resp *commands.GetVulnsResponse, err error) {
	resp = new(commands.GetVulnsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetVulns retrieves vulnerabilities with optional filters.
func (cli *client) GetVulns(ctx context.Context, args ...filtering.FilterArg) (*commands.GetVulnsResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetVulns{
		Filter: filter,
	}
	resp := new(commands.GetVulnsResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (cli *client) CreateAlert(cmd *commands.CreateAlert) (resp *commands.CreateAlertResponse, err error) {
	resp = new(commands.CreateAlertResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetAlerts(cmd *commands.GetAlerts) (resp *commands.GetAlertsResponse, err error) {
	resp = new(commands.GetAlertsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyAlert(cmd *commands.ModifyAlert) (resp *commands.ModifyAlertResponse, err error) {
	resp = new(commands.ModifyAlertResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeleteAlert(cmd *commands.DeleteAlert) (resp *commands.DeleteAlertResponse, err error) {
	resp = new(commands.DeleteAlertResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) TestAlert(cmd *commands.TestAlert) (resp *commands.TestAlertResponse, err error) {
	resp = new(commands.TestAlertResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ResumeTask(cmd *commands.ResumeTask) (resp *commands.ResumeTaskResponse, err error) {
	resp = new(commands.ResumeTaskResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateAsset(cmd *commands.CreateAsset) (resp *commands.CreateAssetResponse, err error) {
	resp = new(commands.CreateAssetResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyAsset(cmd *commands.ModifyAsset) (resp *commands.ModifyAssetResponse, err error) {
	resp = new(commands.ModifyAssetResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetAssetsRaw provides the low-level method for direct command execution.
func (cli *client) GetAssetsRaw(cmd *commands.GetAssets) (resp *commands.GetAssetsResponse, err error) {
	resp = new(commands.GetAssetsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetAssets retrieves assets with optional filters.
func (cli *client) GetAssets(ctx context.Context, args ...filtering.FilterArg) (*commands.GetAssetsResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetAssets{
		Filter: filter,
	}
	resp := new(commands.GetAssetsResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (cli *client) DeleteAsset(cmd *commands.DeleteAsset) (resp *commands.DeleteAssetResponse, err error) {
	resp = new(commands.DeleteAssetResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateSchedule(cmd *commands.CreateSchedule) (resp *commands.CreateScheduleResponse, err error) {
	resp = new(commands.CreateScheduleResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifySchedule(cmd *commands.ModifySchedule) (resp *commands.ModifyScheduleResponse, err error) {
	resp = new(commands.ModifyScheduleResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetSchedules(cmd *commands.GetSchedules) (resp *commands.GetSchedulesResponse, err error) {
	resp = new(commands.GetSchedulesResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeleteSchedule(cmd *commands.DeleteSchedule) (resp *commands.DeleteScheduleResponse, err error) {
	resp = new(commands.DeleteScheduleResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateOverride(cmd *commands.CreateOverride) (resp *commands.CreateOverrideResponse, err error) {
	resp = new(commands.CreateOverrideResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeleteOverride(cmd *commands.DeleteOverride) (resp *commands.DeleteOverrideResponse, err error) {
	resp = new(commands.DeleteOverrideResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetOverridesRaw provides the low-level method for direct command execution.
func (cli *client) GetOverridesRaw(cmd *commands.GetOverrides) (resp *commands.GetOverridesResponse, err error) {
	resp = new(commands.GetOverridesResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetOverrides retrieves overrides with optional filters.
func (cli *client) GetOverrides(ctx context.Context, args ...filtering.FilterArg) (*commands.GetOverridesResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetOverrides{
		Filter: filter,
	}
	resp := new(commands.GetOverridesResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (cli *client) ModifyOverride(cmd *commands.ModifyOverride) (resp *commands.ModifyOverrideResponse, err error) {
	resp = new(commands.ModifyOverrideResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateReport(cmd *commands.CreateReport) (resp *commands.CreateReportResponse, err error) {
	resp = new(commands.CreateReportResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetReports(cmd *commands.GetReports) (resp *commands.GetReportsResponse, err error) {
	resp = new(commands.GetReportsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeleteReport(cmd *commands.DeleteReport) (resp *commands.DeleteReportResponse, err error) {
	resp = new(commands.DeleteReportResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateReportFormat(cmd *commands.CreateReportFormat) (resp *commands.CreateReportFormatResponse, err error) {
	resp = new(commands.CreateReportFormatResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyReportFormat(cmd *commands.ModifyReportFormat) (resp *commands.ModifyReportFormatResponse, err error) {
	resp = new(commands.ModifyReportFormatResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetReportFormats(cmd *commands.GetReportFormats) (resp *commands.GetReportFormatsResponse, err error) {
	resp = new(commands.GetReportFormatsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeleteReportFormat(cmd *commands.DeleteReportFormat) (resp *commands.DeleteReportFormatResponse, err error) {
	resp = new(commands.DeleteReportFormatResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) VerifyReportFormat(cmd *commands.VerifyReportFormat) (resp *commands.VerifyReportFormatResponse, err error) {
	resp = new(commands.VerifyReportFormatResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateReportConfig(cmd *commands.CreateReportConfig) (resp *commands.CreateReportConfigResponse, err error) {
	resp = new(commands.CreateReportConfigResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyReportConfig(cmd *commands.ModifyReportConfig) (resp *commands.ModifyReportConfigResponse, err error) {
	resp = new(commands.ModifyReportConfigResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetReportConfigs(cmd *commands.GetReportConfigs) (resp *commands.GetReportConfigsResponse, err error) {
	resp = new(commands.GetReportConfigsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeleteReportConfig(cmd *commands.DeleteReportConfig) (resp *commands.DeleteReportConfigResponse, err error) {
	resp = new(commands.DeleteReportConfigResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetSystemReports(cmd *commands.GetSystemReports) (resp *commands.GetSystemReportsResponse, err error) {
	resp = new(commands.GetSystemReportsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateCredential(cmd *commands.CreateCredential) (resp *commands.CreateCredentialResponse, err error) {
	resp = new(commands.CreateCredentialResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyCredential(cmd *commands.ModifyCredential) (resp *commands.ModifyCredentialResponse, err error) {
	resp = new(commands.ModifyCredentialResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetCredentials(cmd *commands.GetCredentials) (resp *commands.GetCredentialsResponse, err error) {
	resp = new(commands.GetCredentialsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeleteCredential(cmd *commands.DeleteCredential) (resp *commands.DeleteCredentialResponse, err error) {
	resp = new(commands.DeleteCredentialResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateScanner(cmd *commands.CreateScanner) (resp *commands.CreateScannerResponse, err error) {
	resp = new(commands.CreateScannerResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyScanner(cmd *commands.ModifyScanner) (resp *commands.ModifyScannerResponse, err error) {
	resp = new(commands.ModifyScannerResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeleteScanner(cmd *commands.DeleteScanner) (resp *commands.DeleteScannerResponse, err error) {
	resp = new(commands.DeleteScannerResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) VerifyScanner(cmd *commands.VerifyScanner) (resp *commands.VerifyScannerResponse, err error) {
	resp = new(commands.VerifyScannerResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreatePortList(cmd *commands.CreatePortList) (resp *commands.CreatePortListResponse, err error) {
	resp = new(commands.CreatePortListResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyPortList(cmd *commands.ModifyPortList) (resp *commands.ModifyPortListResponse, err error) {
	resp = new(commands.ModifyPortListResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeletePortList(cmd *commands.DeletePortList) (resp *commands.DeletePortListResponse, err error) {
	resp = new(commands.DeletePortListResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreatePortRange(cmd *commands.CreatePortRange) (resp *commands.CreatePortRangeResponse, err error) {
	resp = new(commands.CreatePortRangeResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeletePortRange(cmd *commands.DeletePortRange) (resp *commands.DeletePortRangeResponse, err error) {
	resp = new(commands.DeletePortRangeResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DescribeAuth(cmd *commands.DescribeAuth) (resp *commands.DescribeAuthResponse, err error) {
	resp = new(commands.DescribeAuthResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyAuth(cmd *commands.ModifyAuth) (resp *commands.ModifyAuthResponse, err error) {
	resp = new(commands.ModifyAuthResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetInfo(cmd *commands.GetInfo) (resp *commands.GetInfoResponse, err error) {
	resp = new(commands.GetInfoResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetVersion(cmd *commands.GetVersion) (resp *commands.GetVersionResponse, err error) {
	resp = new(commands.GetVersionResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) Help(cmd *commands.Help) (resp *commands.HelpResponse, err error) {
	resp = new(commands.HelpResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeleteConfig(cmd *commands.DeleteConfig) (resp *commands.DeleteConfigResponse, err error) {
	resp = new(commands.DeleteConfigResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) SyncConfig(cmd *commands.SyncConfig) (resp *commands.SyncConfigResponse, err error) {
	resp = new(commands.SyncConfigResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) MoveTask(cmd *commands.MoveTask) (resp *commands.MoveTaskResponse, err error) {
	resp = new(commands.MoveTaskResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateTLSCertificate(cmd *commands.CreateTLSCertificate) (resp *commands.CreateTLSCertificateResponse, err error) {
	resp = new(commands.CreateTLSCertificateResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyTLSCertificate(cmd *commands.ModifyTLSCertificate) (resp *commands.ModifyTLSCertificateResponse, err error) {
	resp = new(commands.ModifyTLSCertificateResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetTLSCertificatesRaw provides the low-level method for direct command execution.
func (cli *client) GetTLSCertificatesRaw(cmd *commands.GetTLSCertificates) (resp *commands.GetTLSCertificatesResponse, err error) {
	resp = new(commands.GetTLSCertificatesResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetTLSCertificates retrieves TLS certificates with optional filters.
func (cli *client) GetTLSCertificates(ctx context.Context, args ...filtering.FilterArg) (*commands.GetTLSCertificatesResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetTLSCertificates{
		Filter: filter,
	}
	resp := new(commands.GetTLSCertificatesResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (cli *client) ModifyAgents(cmd *commands.ModifyAgents) (resp *commands.ModifyAgentsResponse, err error) {
	resp = new(commands.ModifyAgentsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetAgents(cmd *commands.GetAgents) (resp *commands.GetAgentsResponse, err error) {
	resp = new(commands.GetAgentsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeleteAgents(cmd *commands.DeleteAgents) (resp *commands.DeleteAgentsResponse, err error) {
	resp = new(commands.DeleteAgentsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetNVTs(cmd *commands.GetNVTs) (resp *commands.GetNVTsResponse, err error) {
	resp = new(commands.GetNVTsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetNVTFamilies(cmd *commands.GetNVTFamilies) (resp *commands.GetNVTFamiliesResponse, err error) {
	resp = new(commands.GetNVTFamiliesResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetFeeds(cmd *commands.GetFeeds) (resp *commands.GetFeedsResponse, err error) {
	resp = new(commands.GetFeedsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyLicense(cmd *commands.ModifyLicense) (resp *commands.ModifyLicenseResponse, err error) {
	resp = new(commands.ModifyLicenseResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetLicense(cmd *commands.GetLicense) (resp *commands.GetLicenseResponse, err error) {
	resp = new(commands.GetLicenseResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifySetting(cmd *commands.ModifySetting) (resp *commands.ModifySettingResponse, err error) {
	resp = new(commands.ModifySettingResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetSettingsRaw provides the low-level method for direct command execution.
func (cli *client) GetSettingsRaw(cmd *commands.GetSettings) (resp *commands.GetSettingsResponse, err error) {
	resp = new(commands.GetSettingsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetSettings retrieves settings with optional filters.
func (cli *client) GetSettings(ctx context.Context, args ...filtering.FilterArg) (*commands.GetSettingsResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetSettings{
		Filter: filter,
	}
	resp := new(commands.GetSettingsResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

// GetResourceNamesRaw provides the low-level method for direct command execution.
func (cli *client) GetResourceNamesRaw(cmd *commands.GetResourceNames) (resp *commands.GetResourceNamesResponse, err error) {
	resp = new(commands.GetResourceNamesResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetResourceNames retrieves resource names with optional filters.
func (cli *client) GetResourceNames(ctx context.Context, args ...filtering.FilterArg) (*commands.GetResourceNamesResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetResourceNames{
		Filter: filter,
	}
	resp := new(commands.GetResourceNamesResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

// GetAggregatesRaw provides the low-level method for direct command execution.
func (cli *client) GetAggregatesRaw(cmd *commands.GetAggregates) (resp *commands.GetAggregatesResponse, err error) {
	resp = new(commands.GetAggregatesResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetAggregates retrieves aggregates with optional filters.
func (cli *client) GetAggregates(ctx context.Context, args ...filtering.FilterArg) (*commands.GetAggregatesResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetAggregates{
		Filter: filter,
	}
	resp := new(commands.GetAggregatesResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (cli *client) GetFeatures(cmd *commands.GetFeatures) (resp *commands.GetFeaturesResponse, err error) {
	resp = new(commands.GetFeaturesResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) EmptyTrashcan(cmd *commands.EmptyTrashcan) (resp *commands.EmptyTrashcanResponse, err error) {
	resp = new(commands.EmptyTrashcanResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) Restore(cmd *commands.Restore) (resp *commands.RestoreResponse, err error) {
	resp = new(commands.RestoreResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) RunWizard(cmd *commands.RunWizard) (resp *commands.RunWizardResponse, err error) {
	resp = new(commands.RunWizardResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) RawXML(xmlStr string) (string, error) {
	return cli.conn.RawXML(xmlStr)
}

func (cli *client) CreateUser(cmd *commands.CreateUser) (resp *commands.CreateUserResponse, err error) {
	resp = new(commands.CreateUserResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyUser(cmd *commands.ModifyUser) (resp *commands.ModifyUserResponse, err error) {
	resp = new(commands.ModifyUserResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetUsersRaw provides the low-level method for direct command execution.
func (cli *client) GetUsersRaw(cmd *commands.GetUsers) (resp *commands.GetUsersResponse, err error) {
	resp = new(commands.GetUsersResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetUsers retrieves users with optional filters.
func (cli *client) GetUsers(ctx context.Context, args ...filtering.FilterArg) (*commands.GetUsersResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetUsers{
		Filter: filter,
	}
	resp := new(commands.GetUsersResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (cli *client) DeleteUser(cmd *commands.DeleteUser) (resp *commands.DeleteUserResponse, err error) {
	resp = new(commands.DeleteUserResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateRole(cmd *commands.CreateRole) (resp *commands.CreateRoleResponse, err error) {
	resp = new(commands.CreateRoleResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyRole(cmd *commands.ModifyRole) (resp *commands.ModifyRoleResponse, err error) {
	resp = new(commands.ModifyRoleResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetRoles(cmd *commands.GetRoles) (resp *commands.GetRolesResponse, err error) {
	resp = new(commands.GetRolesResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeleteRole(cmd *commands.DeleteRole) (resp *commands.DeleteRoleResponse, err error) {
	resp = new(commands.DeleteRoleResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreatePermission(cmd *commands.CreatePermission) (resp *commands.CreatePermissionResponse, err error) {
	resp = new(commands.CreatePermissionResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyPermission(cmd *commands.ModifyPermission) (resp *commands.ModifyPermissionResponse, err error) {
	resp = new(commands.ModifyPermissionResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetPermissions(cmd *commands.GetPermissions) (resp *commands.GetPermissionsResponse, err error) {
	resp = new(commands.GetPermissionsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeletePermission(cmd *commands.DeletePermission) (resp *commands.DeletePermissionResponse, err error) {
	resp = new(commands.DeletePermissionResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateGroup(cmd *commands.CreateGroup) (resp *commands.CreateGroupResponse, err error) {
	resp = new(commands.CreateGroupResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyGroup(cmd *commands.ModifyGroup) (resp *commands.ModifyGroupResponse, err error) {
	resp = new(commands.ModifyGroupResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetGroups(cmd *commands.GetGroups) (resp *commands.GetGroupsResponse, err error) {
	resp = new(commands.GetGroupsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeleteGroup(cmd *commands.DeleteGroup) (resp *commands.DeleteGroupResponse, err error) {
	resp = new(commands.DeleteGroupResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateTag(cmd *commands.CreateTag) (resp *commands.CreateTagResponse, err error) {
	resp = new(commands.CreateTagResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyTag(cmd *commands.ModifyTag) (resp *commands.ModifyTagResponse, err error) {
	resp = new(commands.ModifyTagResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetTagsRaw provides the low-level method for direct command execution.
func (cli *client) GetTagsRaw(cmd *commands.GetTags) (resp *commands.GetTagsResponse, err error) {
	resp = new(commands.GetTagsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetTags retrieves tags with optional filters.
func (cli *client) GetTags(ctx context.Context, args ...filtering.FilterArg) (*commands.GetTagsResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetTags{
		Filter: filter,
	}
	resp := new(commands.GetTagsResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (cli *client) DeleteTag(cmd *commands.DeleteTag) (resp *commands.DeleteTagResponse, err error) {
	resp = new(commands.DeleteTagResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateNote(cmd *commands.CreateNote) (resp *commands.CreateNoteResponse, err error) {
	resp = new(commands.CreateNoteResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyNote(cmd *commands.ModifyNote) (resp *commands.ModifyNoteResponse, err error) {
	resp = new(commands.ModifyNoteResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) GetNotes(cmd *commands.GetNotes) (resp *commands.GetNotesResponse, err error) {
	resp = new(commands.GetNotesResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) DeleteNote(cmd *commands.DeleteNote) (resp *commands.DeleteNoteResponse, err error) {
	resp = new(commands.DeleteNoteResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateFilter(cmd *commands.CreateFilter) (resp *commands.CreateFilterResponse, err error) {
	resp = new(commands.CreateFilterResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyFilter(cmd *commands.ModifyFilter) (resp *commands.ModifyFilterResponse, err error) {
	resp = new(commands.ModifyFilterResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetFiltersRaw provides the low-level method for direct command execution.
func (cli *client) GetFiltersRaw(cmd *commands.GetFilters) (resp *commands.GetFiltersResponse, err error) {
	resp = new(commands.GetFiltersResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetFilters retrieves filters with optional filters.
func (cli *client) GetFilters(ctx context.Context, args ...filtering.FilterArg) (*commands.GetFiltersResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetFilters{
		Filter: filter,
	}
	resp := new(commands.GetFiltersResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (cli *client) DeleteFilter(cmd *commands.DeleteFilter) (resp *commands.DeleteFilterResponse, err error) {
	resp = new(commands.DeleteFilterResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) CreateTicket(cmd *commands.CreateTicket) (resp *commands.CreateTicketResponse, err error) {
	resp = new(commands.CreateTicketResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

func (cli *client) ModifyTicket(cmd *commands.ModifyTicket) (resp *commands.ModifyTicketResponse, err error) {
	resp = new(commands.ModifyTicketResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetTicketsRaw provides the low-level method for direct command execution.
func (cli *client) GetTicketsRaw(cmd *commands.GetTickets) (resp *commands.GetTicketsResponse, err error) {
	resp = new(commands.GetTicketsResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// GetTickets retrieves tickets with optional filters.
func (cli *client) GetTickets(ctx context.Context, args ...filtering.FilterArg) (*commands.GetTicketsResponse, error) {
	filter := filtering.BuildFilterFromArgs(args...)
	cmd := &commands.GetTickets{
		Filter: filter,
	}
	resp := new(commands.GetTicketsResponse)
	err := cli.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (cli *client) DeleteTicket(cmd *commands.DeleteTicket) (resp *commands.DeleteTicketResponse, err error) {
	resp = new(commands.DeleteTicketResponse)
	err = cli.conn.Execute(context.Background(), cmd, resp)
	return resp, err
}

// Tasks returns an iterator for tasks.
func (cli *client) Tasks(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.TaskIterator {
	filter := filtering.BuildFilterFromArgs(args...)
	var filters []string
	if filter != "" {
		filters = []string{filter}
	}
	opts := pagination.PaginationOptions{Page: 1, PageSize: pageSize}
	return pagination.NewTaskIterator(cli, ctx, opts, filters...)
}

// Results returns an iterator for results.
func (cli *client) Results(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.ResultIterator {
	filter := filtering.BuildFilterFromArgs(args...)
	var filters []string
	if filter != "" {
		filters = []string{filter}
	}
	opts := pagination.PaginationOptions{Page: 1, PageSize: pageSize}
	return pagination.NewResultIterator(cli, ctx, opts, filters...)
}

// Assets returns an iterator for assets.
func (cli *client) Assets(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.AssetIterator {
	filter := filtering.BuildFilterFromArgs(args...)
	var filters []string
	if filter != "" {
		filters = []string{filter}
	}
	opts := pagination.PaginationOptions{Page: 1, PageSize: pageSize}
	return pagination.NewAssetIterator(cli, ctx, opts, filters...)
}

// Targets returns an iterator for targets.
func (cli *client) Targets(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.TargetIterator {
	filter := filtering.BuildFilterFromArgs(args...)
	var filters []string
	if filter != "" {
		filters = []string{filter}
	}
	opts := pagination.PaginationOptions{Page: 1, PageSize: pageSize}
	return pagination.NewTargetIterator(cli, ctx, opts, filters...)
}

// Tickets returns an iterator for tickets.
func (cli *client) Tickets(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.TicketIterator {
	filter := filtering.BuildFilterFromArgs(args...)
	var filters []string
	if filter != "" {
		filters = []string{filter}
	}
	opts := pagination.PaginationOptions{Page: 1, PageSize: pageSize}
	return pagination.NewTicketIterator(cli, ctx, opts, filters...)
}

// PortLists returns an iterator for port lists.
func (cli *client) PortLists(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.PortListIterator {
	filter := filtering.BuildFilterFromArgs(args...)
	var filters []string
	if filter != "" {
		filters = []string{filter}
	}
	opts := pagination.PaginationOptions{Page: 1, PageSize: pageSize}
	return pagination.NewPortListIterator(cli, ctx, opts, filters...)
}

// Settings returns an iterator for settings.
func (cli *client) Settings(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.SettingsIterator {
	filter := filtering.BuildFilterFromArgs(args...)
	var filters []string
	if filter != "" {
		filters = []string{filter}
	}
	opts := pagination.PaginationOptions{Page: 1, PageSize: pageSize}
	return pagination.NewSettingsIterator(cli, ctx, opts, filters...)
}
