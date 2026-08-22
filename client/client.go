package client

import (
	"context"
	"fmt"

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

// execute runs cmd against conn and decodes the response into a freshly
// allocated *Resp, collapsing the repetitive "new response, execute, return"
// pattern shared by nearly every wrapper method below.
func execute[Resp any](ctx context.Context, conn gmp.Connection, cmd any) (*Resp, error) {
	resp := new(Resp)
	err := conn.Execute(ctx, cmd, resp)
	return resp, err
}

// Close releases the underlying Connection.
func (cli *client) Close() error {
	return cli.conn.Close()
}

func (cli *client) Authenticate(ctx context.Context, cmd *commands.Authenticate) (*commands.AuthenticateResponse, error) {
	return execute[commands.AuthenticateResponse](ctx, cli.conn, cmd)
}

// GetConfigsRaw provides the low-level method for direct command execution.
func (cli *client) GetConfigsRaw(ctx context.Context, cmd *commands.GetConfigs) (*commands.GetConfigsResponse, error) {
	return execute[commands.GetConfigsResponse](ctx, cli.conn, cmd)
}

// GetConfigs retrieves configs with optional filters.
func (cli *client) GetConfigs(ctx context.Context, args ...filtering.FilterArg) (*commands.GetConfigsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetConfigs: %w", err)
	}
	cmd := &commands.GetConfigs{
		Filter: filter,
	}
	return cli.GetConfigsRaw(ctx, cmd)
}

// GetScannersRaw provides the low-level method for direct command execution.
func (cli *client) GetScannersRaw(ctx context.Context, cmd *commands.GetScanners) (*commands.GetScannersResponse, error) {
	return execute[commands.GetScannersResponse](ctx, cli.conn, cmd)
}

// GetScanners retrieves scanners with optional filters.
func (cli *client) GetScanners(ctx context.Context, args ...filtering.FilterArg) (*commands.GetScannersResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetScanners: %w", err)
	}
	cmd := &commands.GetScanners{
		Filter: filter,
	}
	return cli.GetScannersRaw(ctx, cmd)
}

func (cli *client) GetPreferences(ctx context.Context, cmd *commands.GetPreferences) (*commands.GetPreferencesResponse, error) {
	return execute[commands.GetPreferencesResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateConfig(ctx context.Context, cmd *commands.CreateConfig) (*commands.CreateConfigResponse, error) {
	return execute[commands.CreateConfigResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyConfig(ctx context.Context, cmd *commands.ModifyConfig) (*commands.ModifyConfigResponse, error) {
	return execute[commands.ModifyConfigResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateTask(ctx context.Context, cmd *commands.CreateTask) (*commands.CreateTaskResponse, error) {
	return execute[commands.CreateTaskResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyTask(ctx context.Context, cmd *commands.ModifyTask) (*commands.ModifyTaskResponse, error) {
	return execute[commands.ModifyTaskResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateTarget(ctx context.Context, cmd *commands.CreateTarget) (*commands.CreateTargetResponse, error) {
	return execute[commands.CreateTargetResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyTarget(ctx context.Context, cmd *commands.ModifyTarget) (*commands.ModifyTargetResponse, error) {
	return execute[commands.ModifyTargetResponse](ctx, cli.conn, cmd)
}

// GetTargetsRaw provides the low-level method for direct command execution.
func (cli *client) GetTargetsRaw(ctx context.Context, cmd *commands.GetTargets) (*commands.GetTargetsResponse, error) {
	return execute[commands.GetTargetsResponse](ctx, cli.conn, cmd)
}

// GetTargets retrieves targets with optional filters.
func (cli *client) GetTargets(ctx context.Context, args ...filtering.FilterArg) (*commands.GetTargetsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetTargets: %w", err)
	}
	cmd := &commands.GetTargets{
		Filter: filter,
	}
	return cli.GetTargetsRaw(ctx, cmd)
}

func (cli *client) DeleteTarget(ctx context.Context, cmd *commands.DeleteTarget) (*commands.DeleteTargetResponse, error) {
	return execute[commands.DeleteTargetResponse](ctx, cli.conn, cmd)
}

// GetPortListsRaw provides the low-level method for direct command execution.
func (cli *client) GetPortListsRaw(ctx context.Context, cmd *commands.GetPortLists) (*commands.GetPortListsResponse, error) {
	return execute[commands.GetPortListsResponse](ctx, cli.conn, cmd)
}

// GetPortLists retrieves port lists with optional filters.
func (cli *client) GetPortLists(ctx context.Context, args ...filtering.FilterArg) (*commands.GetPortListsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetPortLists: %w", err)
	}
	cmd := &commands.GetPortLists{
		Filter: filter,
	}
	return cli.GetPortListsRaw(ctx, cmd)
}

func (cli *client) StartTask(ctx context.Context, cmd *commands.StartTask) (*commands.StartTaskResponse, error) {
	return execute[commands.StartTaskResponse](ctx, cli.conn, cmd)
}

// GetTasksRaw provides the low-level method for direct command execution.
func (cli *client) GetTasksRaw(ctx context.Context, cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
	return execute[commands.GetTasksResponse](ctx, cli.conn, cmd)
}

// GetTasks retrieves tasks with optional filters.
func (cli *client) GetTasks(ctx context.Context, args ...filtering.FilterArg) (*commands.GetTasksResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetTasks: %w", err)
	}
	cmd := &commands.GetTasks{
		Filter: filter,
	}
	return cli.GetTasksRaw(ctx, cmd)
}

func (cli *client) StopTask(ctx context.Context, cmd *commands.StopTask) (*commands.StopTaskResponse, error) {
	return execute[commands.StopTaskResponse](ctx, cli.conn, cmd)
}

func (cli *client) DeleteTask(ctx context.Context, cmd *commands.DeleteTask) (*commands.DeleteTaskResponse, error) {
	return execute[commands.DeleteTaskResponse](ctx, cli.conn, cmd)
}

// GetResultsRaw provides the low-level method for direct command execution.
func (cli *client) GetResultsRaw(ctx context.Context, cmd *commands.GetResults) (*commands.GetResultsResponse, error) {
	return execute[commands.GetResultsResponse](ctx, cli.conn, cmd)
}

// GetResults retrieves results with optional filters.
func (cli *client) GetResults(ctx context.Context, args ...filtering.FilterArg) (*commands.GetResultsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetResults: %w", err)
	}
	cmd := &commands.GetResults{
		Filter: filter,
	}
	return cli.GetResultsRaw(ctx, cmd)
}

// GetVulnsRaw provides the low-level method for direct command execution.
func (cli *client) GetVulnsRaw(ctx context.Context, cmd *commands.GetVulns) (*commands.GetVulnsResponse, error) {
	return execute[commands.GetVulnsResponse](ctx, cli.conn, cmd)
}

// GetVulns retrieves vulnerabilities with optional filters.
func (cli *client) GetVulns(ctx context.Context, args ...filtering.FilterArg) (*commands.GetVulnsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetVulns: %w", err)
	}
	cmd := &commands.GetVulns{
		Filter: filter,
	}
	return cli.GetVulnsRaw(ctx, cmd)
}

func (cli *client) CreateAlert(ctx context.Context, cmd *commands.CreateAlert) (*commands.CreateAlertResponse, error) {
	return execute[commands.CreateAlertResponse](ctx, cli.conn, cmd)
}

// GetAlertsRaw provides the low-level method for direct command execution.
func (cli *client) GetAlertsRaw(ctx context.Context, cmd *commands.GetAlerts) (*commands.GetAlertsResponse, error) {
	return execute[commands.GetAlertsResponse](ctx, cli.conn, cmd)
}

// GetAlerts retrieves alerts with optional filters.
func (cli *client) GetAlerts(ctx context.Context, args ...filtering.FilterArg) (*commands.GetAlertsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetAlerts: %w", err)
	}
	cmd := &commands.GetAlerts{
		Filter: filter,
	}
	return cli.GetAlertsRaw(ctx, cmd)
}

func (cli *client) ModifyAlert(ctx context.Context, cmd *commands.ModifyAlert) (*commands.ModifyAlertResponse, error) {
	return execute[commands.ModifyAlertResponse](ctx, cli.conn, cmd)
}

func (cli *client) DeleteAlert(ctx context.Context, cmd *commands.DeleteAlert) (*commands.DeleteAlertResponse, error) {
	return execute[commands.DeleteAlertResponse](ctx, cli.conn, cmd)
}

func (cli *client) TestAlert(ctx context.Context, cmd *commands.TestAlert) (*commands.TestAlertResponse, error) {
	return execute[commands.TestAlertResponse](ctx, cli.conn, cmd)
}

func (cli *client) ResumeTask(ctx context.Context, cmd *commands.ResumeTask) (*commands.ResumeTaskResponse, error) {
	return execute[commands.ResumeTaskResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateAsset(ctx context.Context, cmd *commands.CreateAsset) (*commands.CreateAssetResponse, error) {
	return execute[commands.CreateAssetResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyAsset(ctx context.Context, cmd *commands.ModifyAsset) (*commands.ModifyAssetResponse, error) {
	return execute[commands.ModifyAssetResponse](ctx, cli.conn, cmd)
}

// GetAssetsRaw provides the low-level method for direct command execution.
func (cli *client) GetAssetsRaw(ctx context.Context, cmd *commands.GetAssets) (*commands.GetAssetsResponse, error) {
	return execute[commands.GetAssetsResponse](ctx, cli.conn, cmd)
}

// GetAssets retrieves assets with optional filters.
func (cli *client) GetAssets(ctx context.Context, args ...filtering.FilterArg) (*commands.GetAssetsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetAssets: %w", err)
	}
	cmd := &commands.GetAssets{
		Filter: filter,
	}
	return cli.GetAssetsRaw(ctx, cmd)
}

func (cli *client) DeleteAsset(ctx context.Context, cmd *commands.DeleteAsset) (*commands.DeleteAssetResponse, error) {
	return execute[commands.DeleteAssetResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateSchedule(ctx context.Context, cmd *commands.CreateSchedule) (*commands.CreateScheduleResponse, error) {
	return execute[commands.CreateScheduleResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifySchedule(ctx context.Context, cmd *commands.ModifySchedule) (*commands.ModifyScheduleResponse, error) {
	return execute[commands.ModifyScheduleResponse](ctx, cli.conn, cmd)
}

// GetSchedulesRaw provides the low-level method for direct command execution.
func (cli *client) GetSchedulesRaw(ctx context.Context, cmd *commands.GetSchedules) (*commands.GetSchedulesResponse, error) {
	return execute[commands.GetSchedulesResponse](ctx, cli.conn, cmd)
}

// GetSchedules retrieves schedules with optional filters.
func (cli *client) GetSchedules(ctx context.Context, args ...filtering.FilterArg) (*commands.GetSchedulesResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetSchedules: %w", err)
	}
	cmd := &commands.GetSchedules{
		Filter: filter,
	}
	return cli.GetSchedulesRaw(ctx, cmd)
}

func (cli *client) DeleteSchedule(ctx context.Context, cmd *commands.DeleteSchedule) (*commands.DeleteScheduleResponse, error) {
	return execute[commands.DeleteScheduleResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateOverride(ctx context.Context, cmd *commands.CreateOverride) (*commands.CreateOverrideResponse, error) {
	return execute[commands.CreateOverrideResponse](ctx, cli.conn, cmd)
}

func (cli *client) DeleteOverride(ctx context.Context, cmd *commands.DeleteOverride) (*commands.DeleteOverrideResponse, error) {
	return execute[commands.DeleteOverrideResponse](ctx, cli.conn, cmd)
}

// GetOverridesRaw provides the low-level method for direct command execution.
func (cli *client) GetOverridesRaw(ctx context.Context, cmd *commands.GetOverrides) (*commands.GetOverridesResponse, error) {
	return execute[commands.GetOverridesResponse](ctx, cli.conn, cmd)
}

// GetOverrides retrieves overrides with optional filters.
func (cli *client) GetOverrides(ctx context.Context, args ...filtering.FilterArg) (*commands.GetOverridesResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetOverrides: %w", err)
	}
	cmd := &commands.GetOverrides{
		Filter: filter,
	}
	return cli.GetOverridesRaw(ctx, cmd)
}

func (cli *client) ModifyOverride(ctx context.Context, cmd *commands.ModifyOverride) (*commands.ModifyOverrideResponse, error) {
	return execute[commands.ModifyOverrideResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateReport(ctx context.Context, cmd *commands.CreateReport) (*commands.CreateReportResponse, error) {
	return execute[commands.CreateReportResponse](ctx, cli.conn, cmd)
}

// GetReportsRaw provides the low-level method for direct command execution.
func (cli *client) GetReportsRaw(ctx context.Context, cmd *commands.GetReports) (*commands.GetReportsResponse, error) {
	return execute[commands.GetReportsResponse](ctx, cli.conn, cmd)
}

// GetReports retrieves reports with optional filters.
func (cli *client) GetReports(ctx context.Context, args ...filtering.FilterArg) (*commands.GetReportsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetReports: %w", err)
	}
	cmd := &commands.GetReports{
		Filter: filter,
	}
	return cli.GetReportsRaw(ctx, cmd)
}

func (cli *client) DeleteReport(ctx context.Context, cmd *commands.DeleteReport) (*commands.DeleteReportResponse, error) {
	return execute[commands.DeleteReportResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateReportFormat(ctx context.Context, cmd *commands.CreateReportFormat) (*commands.CreateReportFormatResponse, error) {
	return execute[commands.CreateReportFormatResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyReportFormat(ctx context.Context, cmd *commands.ModifyReportFormat) (*commands.ModifyReportFormatResponse, error) {
	return execute[commands.ModifyReportFormatResponse](ctx, cli.conn, cmd)
}

// GetReportFormatsRaw provides the low-level method for direct command execution.
func (cli *client) GetReportFormatsRaw(ctx context.Context, cmd *commands.GetReportFormats) (*commands.GetReportFormatsResponse, error) {
	return execute[commands.GetReportFormatsResponse](ctx, cli.conn, cmd)
}

// GetReportFormats retrieves report formats with optional filters.
func (cli *client) GetReportFormats(ctx context.Context, args ...filtering.FilterArg) (*commands.GetReportFormatsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetReportFormats: %w", err)
	}
	cmd := &commands.GetReportFormats{
		Filter: filter,
	}
	return cli.GetReportFormatsRaw(ctx, cmd)
}

func (cli *client) DeleteReportFormat(ctx context.Context, cmd *commands.DeleteReportFormat) (*commands.DeleteReportFormatResponse, error) {
	return execute[commands.DeleteReportFormatResponse](ctx, cli.conn, cmd)
}

func (cli *client) VerifyReportFormat(ctx context.Context, cmd *commands.VerifyReportFormat) (*commands.VerifyReportFormatResponse, error) {
	return execute[commands.VerifyReportFormatResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateReportConfig(ctx context.Context, cmd *commands.CreateReportConfig) (*commands.CreateReportConfigResponse, error) {
	return execute[commands.CreateReportConfigResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyReportConfig(ctx context.Context, cmd *commands.ModifyReportConfig) (*commands.ModifyReportConfigResponse, error) {
	return execute[commands.ModifyReportConfigResponse](ctx, cli.conn, cmd)
}

// GetReportConfigsRaw provides the low-level method for direct command execution.
func (cli *client) GetReportConfigsRaw(ctx context.Context, cmd *commands.GetReportConfigs) (*commands.GetReportConfigsResponse, error) {
	return execute[commands.GetReportConfigsResponse](ctx, cli.conn, cmd)
}

// GetReportConfigs retrieves report configs with optional filters.
func (cli *client) GetReportConfigs(ctx context.Context, args ...filtering.FilterArg) (*commands.GetReportConfigsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetReportConfigs: %w", err)
	}
	cmd := &commands.GetReportConfigs{
		Filter: filter,
	}
	return cli.GetReportConfigsRaw(ctx, cmd)
}

func (cli *client) DeleteReportConfig(ctx context.Context, cmd *commands.DeleteReportConfig) (*commands.DeleteReportConfigResponse, error) {
	return execute[commands.DeleteReportConfigResponse](ctx, cli.conn, cmd)
}

func (cli *client) GetSystemReports(ctx context.Context, cmd *commands.GetSystemReports) (*commands.GetSystemReportsResponse, error) {
	return execute[commands.GetSystemReportsResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateCredential(ctx context.Context, cmd *commands.CreateCredential) (*commands.CreateCredentialResponse, error) {
	return execute[commands.CreateCredentialResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyCredential(ctx context.Context, cmd *commands.ModifyCredential) (*commands.ModifyCredentialResponse, error) {
	return execute[commands.ModifyCredentialResponse](ctx, cli.conn, cmd)
}

// GetCredentialsRaw provides the low-level method for direct command execution.
func (cli *client) GetCredentialsRaw(ctx context.Context, cmd *commands.GetCredentials) (*commands.GetCredentialsResponse, error) {
	return execute[commands.GetCredentialsResponse](ctx, cli.conn, cmd)
}

// GetCredentials retrieves credentials with optional filters.
func (cli *client) GetCredentials(ctx context.Context, args ...filtering.FilterArg) (*commands.GetCredentialsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetCredentials: %w", err)
	}
	cmd := &commands.GetCredentials{
		Filter: filter,
	}
	return cli.GetCredentialsRaw(ctx, cmd)
}

func (cli *client) DeleteCredential(ctx context.Context, cmd *commands.DeleteCredential) (*commands.DeleteCredentialResponse, error) {
	return execute[commands.DeleteCredentialResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateScanner(ctx context.Context, cmd *commands.CreateScanner) (*commands.CreateScannerResponse, error) {
	return execute[commands.CreateScannerResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyScanner(ctx context.Context, cmd *commands.ModifyScanner) (*commands.ModifyScannerResponse, error) {
	return execute[commands.ModifyScannerResponse](ctx, cli.conn, cmd)
}

func (cli *client) DeleteScanner(ctx context.Context, cmd *commands.DeleteScanner) (*commands.DeleteScannerResponse, error) {
	return execute[commands.DeleteScannerResponse](ctx, cli.conn, cmd)
}

func (cli *client) VerifyScanner(ctx context.Context, cmd *commands.VerifyScanner) (*commands.VerifyScannerResponse, error) {
	return execute[commands.VerifyScannerResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreatePortList(ctx context.Context, cmd *commands.CreatePortList) (*commands.CreatePortListResponse, error) {
	return execute[commands.CreatePortListResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyPortList(ctx context.Context, cmd *commands.ModifyPortList) (*commands.ModifyPortListResponse, error) {
	return execute[commands.ModifyPortListResponse](ctx, cli.conn, cmd)
}

func (cli *client) DeletePortList(ctx context.Context, cmd *commands.DeletePortList) (*commands.DeletePortListResponse, error) {
	return execute[commands.DeletePortListResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreatePortRange(ctx context.Context, cmd *commands.CreatePortRange) (*commands.CreatePortRangeResponse, error) {
	return execute[commands.CreatePortRangeResponse](ctx, cli.conn, cmd)
}

func (cli *client) DeletePortRange(ctx context.Context, cmd *commands.DeletePortRange) (*commands.DeletePortRangeResponse, error) {
	return execute[commands.DeletePortRangeResponse](ctx, cli.conn, cmd)
}

func (cli *client) DescribeAuth(ctx context.Context, cmd *commands.DescribeAuth) (*commands.DescribeAuthResponse, error) {
	return execute[commands.DescribeAuthResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyAuth(ctx context.Context, cmd *commands.ModifyAuth) (*commands.ModifyAuthResponse, error) {
	return execute[commands.ModifyAuthResponse](ctx, cli.conn, cmd)
}

func (cli *client) GetInfo(ctx context.Context, cmd *commands.GetInfo) (*commands.GetInfoResponse, error) {
	return execute[commands.GetInfoResponse](ctx, cli.conn, cmd)
}

func (cli *client) GetVersion(ctx context.Context, cmd *commands.GetVersion) (*commands.GetVersionResponse, error) {
	return execute[commands.GetVersionResponse](ctx, cli.conn, cmd)
}

func (cli *client) Help(ctx context.Context, cmd *commands.Help) (*commands.HelpResponse, error) {
	return execute[commands.HelpResponse](ctx, cli.conn, cmd)
}

func (cli *client) DeleteConfig(ctx context.Context, cmd *commands.DeleteConfig) (*commands.DeleteConfigResponse, error) {
	return execute[commands.DeleteConfigResponse](ctx, cli.conn, cmd)
}

func (cli *client) SyncConfig(ctx context.Context, cmd *commands.SyncConfig) (*commands.SyncConfigResponse, error) {
	return execute[commands.SyncConfigResponse](ctx, cli.conn, cmd)
}

func (cli *client) MoveTask(ctx context.Context, cmd *commands.MoveTask) (*commands.MoveTaskResponse, error) {
	return execute[commands.MoveTaskResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateTLSCertificate(ctx context.Context, cmd *commands.CreateTLSCertificate) (*commands.CreateTLSCertificateResponse, error) {
	return execute[commands.CreateTLSCertificateResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyTLSCertificate(ctx context.Context, cmd *commands.ModifyTLSCertificate) (*commands.ModifyTLSCertificateResponse, error) {
	return execute[commands.ModifyTLSCertificateResponse](ctx, cli.conn, cmd)
}

// GetTLSCertificatesRaw provides the low-level method for direct command execution.
func (cli *client) GetTLSCertificatesRaw(ctx context.Context, cmd *commands.GetTLSCertificates) (*commands.GetTLSCertificatesResponse, error) {
	return execute[commands.GetTLSCertificatesResponse](ctx, cli.conn, cmd)
}

// GetTLSCertificates retrieves TLS certificates with optional filters.
func (cli *client) GetTLSCertificates(ctx context.Context, args ...filtering.FilterArg) (*commands.GetTLSCertificatesResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetTLSCertificates: %w", err)
	}
	cmd := &commands.GetTLSCertificates{
		Filter: filter,
	}
	return cli.GetTLSCertificatesRaw(ctx, cmd)
}

func (cli *client) GetNVTs(ctx context.Context, cmd *commands.GetNVTs) (*commands.GetNVTsResponse, error) {
	return execute[commands.GetNVTsResponse](ctx, cli.conn, cmd)
}

func (cli *client) GetNVTFamilies(ctx context.Context, cmd *commands.GetNVTFamilies) (*commands.GetNVTFamiliesResponse, error) {
	return execute[commands.GetNVTFamiliesResponse](ctx, cli.conn, cmd)
}

func (cli *client) GetFeeds(ctx context.Context, cmd *commands.GetFeeds) (*commands.GetFeedsResponse, error) {
	return execute[commands.GetFeedsResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyLicense(ctx context.Context, cmd *commands.ModifyLicense) (*commands.ModifyLicenseResponse, error) {
	return execute[commands.ModifyLicenseResponse](ctx, cli.conn, cmd)
}

func (cli *client) GetLicense(ctx context.Context, cmd *commands.GetLicense) (*commands.GetLicenseResponse, error) {
	return execute[commands.GetLicenseResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifySetting(ctx context.Context, cmd *commands.ModifySetting) (*commands.ModifySettingResponse, error) {
	return execute[commands.ModifySettingResponse](ctx, cli.conn, cmd)
}

// GetSettingsRaw provides the low-level method for direct command execution.
func (cli *client) GetSettingsRaw(ctx context.Context, cmd *commands.GetSettings) (*commands.GetSettingsResponse, error) {
	return execute[commands.GetSettingsResponse](ctx, cli.conn, cmd)
}

// GetSettings retrieves settings with optional filters.
func (cli *client) GetSettings(ctx context.Context, args ...filtering.FilterArg) (*commands.GetSettingsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetSettings: %w", err)
	}
	cmd := &commands.GetSettings{
		Filter: filter,
	}
	return cli.GetSettingsRaw(ctx, cmd)
}

func (cli *client) GetResourceNames(ctx context.Context, cmd *commands.GetResourceNames) (*commands.GetResourceNamesResponse, error) {
	return execute[commands.GetResourceNamesResponse](ctx, cli.conn, cmd)
}

func (cli *client) GetAggregates(ctx context.Context, cmd *commands.GetAggregates) (*commands.GetAggregatesResponse, error) {
	return execute[commands.GetAggregatesResponse](ctx, cli.conn, cmd)
}

func (cli *client) GetFeatures(ctx context.Context, cmd *commands.GetFeatures) (*commands.GetFeaturesResponse, error) {
	return execute[commands.GetFeaturesResponse](ctx, cli.conn, cmd)
}

func (cli *client) EmptyTrashcan(ctx context.Context, cmd *commands.EmptyTrashcan) (*commands.EmptyTrashcanResponse, error) {
	return execute[commands.EmptyTrashcanResponse](ctx, cli.conn, cmd)
}

func (cli *client) Restore(ctx context.Context, cmd *commands.Restore) (*commands.RestoreResponse, error) {
	return execute[commands.RestoreResponse](ctx, cli.conn, cmd)
}

func (cli *client) RunWizard(ctx context.Context, cmd *commands.RunWizard) (*commands.RunWizardResponse, error) {
	return execute[commands.RunWizardResponse](ctx, cli.conn, cmd)
}

func (cli *client) RawXML(ctx context.Context, xmlStr string) (string, error) {
	return cli.conn.RawXML(ctx, xmlStr)
}

func (cli *client) CreateUser(ctx context.Context, cmd *commands.CreateUser) (*commands.CreateUserResponse, error) {
	return execute[commands.CreateUserResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyUser(ctx context.Context, cmd *commands.ModifyUser) (*commands.ModifyUserResponse, error) {
	return execute[commands.ModifyUserResponse](ctx, cli.conn, cmd)
}

// GetUsersRaw provides the low-level method for direct command execution.
func (cli *client) GetUsersRaw(ctx context.Context, cmd *commands.GetUsers) (*commands.GetUsersResponse, error) {
	return execute[commands.GetUsersResponse](ctx, cli.conn, cmd)
}

// GetUsers retrieves users with optional filters.
func (cli *client) GetUsers(ctx context.Context, args ...filtering.FilterArg) (*commands.GetUsersResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetUsers: %w", err)
	}
	cmd := &commands.GetUsers{
		Filter: filter,
	}
	return cli.GetUsersRaw(ctx, cmd)
}

func (cli *client) DeleteUser(ctx context.Context, cmd *commands.DeleteUser) (*commands.DeleteUserResponse, error) {
	return execute[commands.DeleteUserResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateRole(ctx context.Context, cmd *commands.CreateRole) (*commands.CreateRoleResponse, error) {
	return execute[commands.CreateRoleResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyRole(ctx context.Context, cmd *commands.ModifyRole) (*commands.ModifyRoleResponse, error) {
	return execute[commands.ModifyRoleResponse](ctx, cli.conn, cmd)
}

// GetRolesRaw provides the low-level method for direct command execution.
func (cli *client) GetRolesRaw(ctx context.Context, cmd *commands.GetRoles) (*commands.GetRolesResponse, error) {
	return execute[commands.GetRolesResponse](ctx, cli.conn, cmd)
}

// GetRoles retrieves roles with optional filters.
func (cli *client) GetRoles(ctx context.Context, args ...filtering.FilterArg) (*commands.GetRolesResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetRoles: %w", err)
	}
	cmd := &commands.GetRoles{
		Filter: filter,
	}
	return cli.GetRolesRaw(ctx, cmd)
}

func (cli *client) DeleteRole(ctx context.Context, cmd *commands.DeleteRole) (*commands.DeleteRoleResponse, error) {
	return execute[commands.DeleteRoleResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreatePermission(ctx context.Context, cmd *commands.CreatePermission) (*commands.CreatePermissionResponse, error) {
	return execute[commands.CreatePermissionResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyPermission(ctx context.Context, cmd *commands.ModifyPermission) (*commands.ModifyPermissionResponse, error) {
	return execute[commands.ModifyPermissionResponse](ctx, cli.conn, cmd)
}

// GetPermissionsRaw provides the low-level method for direct command execution.
func (cli *client) GetPermissionsRaw(ctx context.Context, cmd *commands.GetPermissions) (*commands.GetPermissionsResponse, error) {
	return execute[commands.GetPermissionsResponse](ctx, cli.conn, cmd)
}

// GetPermissions retrieves permissions with optional filters.
func (cli *client) GetPermissions(ctx context.Context, args ...filtering.FilterArg) (*commands.GetPermissionsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetPermissions: %w", err)
	}
	cmd := &commands.GetPermissions{
		Filter: filter,
	}
	return cli.GetPermissionsRaw(ctx, cmd)
}

func (cli *client) DeletePermission(ctx context.Context, cmd *commands.DeletePermission) (*commands.DeletePermissionResponse, error) {
	return execute[commands.DeletePermissionResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateGroup(ctx context.Context, cmd *commands.CreateGroup) (*commands.CreateGroupResponse, error) {
	return execute[commands.CreateGroupResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyGroup(ctx context.Context, cmd *commands.ModifyGroup) (*commands.ModifyGroupResponse, error) {
	return execute[commands.ModifyGroupResponse](ctx, cli.conn, cmd)
}

// GetGroupsRaw provides the low-level method for direct command execution.
func (cli *client) GetGroupsRaw(ctx context.Context, cmd *commands.GetGroups) (*commands.GetGroupsResponse, error) {
	return execute[commands.GetGroupsResponse](ctx, cli.conn, cmd)
}

// GetGroups retrieves groups with optional filters.
func (cli *client) GetGroups(ctx context.Context, args ...filtering.FilterArg) (*commands.GetGroupsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetGroups: %w", err)
	}
	cmd := &commands.GetGroups{
		Filter: filter,
	}
	return cli.GetGroupsRaw(ctx, cmd)
}

func (cli *client) DeleteGroup(ctx context.Context, cmd *commands.DeleteGroup) (*commands.DeleteGroupResponse, error) {
	return execute[commands.DeleteGroupResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateTag(ctx context.Context, cmd *commands.CreateTag) (*commands.CreateTagResponse, error) {
	return execute[commands.CreateTagResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyTag(ctx context.Context, cmd *commands.ModifyTag) (*commands.ModifyTagResponse, error) {
	return execute[commands.ModifyTagResponse](ctx, cli.conn, cmd)
}

// GetTagsRaw provides the low-level method for direct command execution.
func (cli *client) GetTagsRaw(ctx context.Context, cmd *commands.GetTags) (*commands.GetTagsResponse, error) {
	return execute[commands.GetTagsResponse](ctx, cli.conn, cmd)
}

// GetTags retrieves tags with optional filters.
func (cli *client) GetTags(ctx context.Context, args ...filtering.FilterArg) (*commands.GetTagsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetTags: %w", err)
	}
	cmd := &commands.GetTags{
		Filter: filter,
	}
	return cli.GetTagsRaw(ctx, cmd)
}

func (cli *client) DeleteTag(ctx context.Context, cmd *commands.DeleteTag) (*commands.DeleteTagResponse, error) {
	return execute[commands.DeleteTagResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateNote(ctx context.Context, cmd *commands.CreateNote) (*commands.CreateNoteResponse, error) {
	return execute[commands.CreateNoteResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyNote(ctx context.Context, cmd *commands.ModifyNote) (*commands.ModifyNoteResponse, error) {
	return execute[commands.ModifyNoteResponse](ctx, cli.conn, cmd)
}

// GetNotesRaw provides the low-level method for direct command execution.
func (cli *client) GetNotesRaw(ctx context.Context, cmd *commands.GetNotes) (*commands.GetNotesResponse, error) {
	return execute[commands.GetNotesResponse](ctx, cli.conn, cmd)
}

// GetNotes retrieves notes with optional filters.
func (cli *client) GetNotes(ctx context.Context, args ...filtering.FilterArg) (*commands.GetNotesResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetNotes: %w", err)
	}
	cmd := &commands.GetNotes{
		Filter: filter,
	}
	return cli.GetNotesRaw(ctx, cmd)
}

func (cli *client) DeleteNote(ctx context.Context, cmd *commands.DeleteNote) (*commands.DeleteNoteResponse, error) {
	return execute[commands.DeleteNoteResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateFilter(ctx context.Context, cmd *commands.CreateFilter) (*commands.CreateFilterResponse, error) {
	return execute[commands.CreateFilterResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyFilter(ctx context.Context, cmd *commands.ModifyFilter) (*commands.ModifyFilterResponse, error) {
	return execute[commands.ModifyFilterResponse](ctx, cli.conn, cmd)
}

// GetFiltersRaw provides the low-level method for direct command execution.
func (cli *client) GetFiltersRaw(ctx context.Context, cmd *commands.GetFilters) (*commands.GetFiltersResponse, error) {
	return execute[commands.GetFiltersResponse](ctx, cli.conn, cmd)
}

// GetFilters retrieves filters with optional filters.
func (cli *client) GetFilters(ctx context.Context, args ...filtering.FilterArg) (*commands.GetFiltersResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetFilters: %w", err)
	}
	cmd := &commands.GetFilters{
		Filter: filter,
	}
	return cli.GetFiltersRaw(ctx, cmd)
}

func (cli *client) DeleteFilter(ctx context.Context, cmd *commands.DeleteFilter) (*commands.DeleteFilterResponse, error) {
	return execute[commands.DeleteFilterResponse](ctx, cli.conn, cmd)
}

func (cli *client) CreateTicket(ctx context.Context, cmd *commands.CreateTicket) (*commands.CreateTicketResponse, error) {
	return execute[commands.CreateTicketResponse](ctx, cli.conn, cmd)
}

func (cli *client) ModifyTicket(ctx context.Context, cmd *commands.ModifyTicket) (*commands.ModifyTicketResponse, error) {
	return execute[commands.ModifyTicketResponse](ctx, cli.conn, cmd)
}

// GetTicketsRaw provides the low-level method for direct command execution.
func (cli *client) GetTicketsRaw(ctx context.Context, cmd *commands.GetTickets) (*commands.GetTicketsResponse, error) {
	return execute[commands.GetTicketsResponse](ctx, cli.conn, cmd)
}

// GetTickets retrieves tickets with optional filters.
func (cli *client) GetTickets(ctx context.Context, args ...filtering.FilterArg) (*commands.GetTicketsResponse, error) {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return nil, fmt.Errorf("GetTickets: %w", err)
	}
	cmd := &commands.GetTickets{
		Filter: filter,
	}
	return cli.GetTicketsRaw(ctx, cmd)
}

func (cli *client) DeleteTicket(ctx context.Context, cmd *commands.DeleteTicket) (*commands.DeleteTicketResponse, error) {
	return execute[commands.DeleteTicketResponse](ctx, cli.conn, cmd)
}

// Tasks returns an iterator for tasks.
func (cli *client) Tasks(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.TaskIterator {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return pagination.NewFailedIterator[*commands.Task](fmt.Errorf("Tasks: %w", err))
	}
	var filters []string
	if filter != "" {
		filters = []string{filter}
	}
	opts := pagination.PaginationOptions{Page: 1, PageSize: pageSize}
	return pagination.NewTaskIterator(cli, ctx, opts, filters...)
}

// Results returns an iterator for results.
func (cli *client) Results(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.ResultIterator {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return pagination.NewFailedIterator[*commands.Result](fmt.Errorf("Results: %w", err))
	}
	var filters []string
	if filter != "" {
		filters = []string{filter}
	}
	opts := pagination.PaginationOptions{Page: 1, PageSize: pageSize}
	return pagination.NewResultIterator(cli, ctx, opts, filters...)
}

// Assets returns an iterator for assets.
func (cli *client) Assets(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.AssetIterator {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return pagination.NewFailedIterator[*commands.GetAssetsAsset](fmt.Errorf("Assets: %w", err))
	}
	var filters []string
	if filter != "" {
		filters = []string{filter}
	}
	opts := pagination.PaginationOptions{Page: 1, PageSize: pageSize}
	return pagination.NewAssetIterator(cli, ctx, opts, filters...)
}

// Targets returns an iterator for targets.
func (cli *client) Targets(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.TargetIterator {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return pagination.NewFailedIterator[*commands.Target](fmt.Errorf("Targets: %w", err))
	}
	var filters []string
	if filter != "" {
		filters = []string{filter}
	}
	opts := pagination.PaginationOptions{Page: 1, PageSize: pageSize}
	return pagination.NewTargetIterator(cli, ctx, opts, filters...)
}

// Tickets returns an iterator for tickets.
func (cli *client) Tickets(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.TicketIterator {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return pagination.NewFailedIterator[*commands.Ticket](fmt.Errorf("Tickets: %w", err))
	}
	var filters []string
	if filter != "" {
		filters = []string{filter}
	}
	opts := pagination.PaginationOptions{Page: 1, PageSize: pageSize}
	return pagination.NewTicketIterator(cli, ctx, opts, filters...)
}

// PortLists returns an iterator for port lists.
func (cli *client) PortLists(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.PortListIterator {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return pagination.NewFailedIterator[*commands.PortList](fmt.Errorf("PortLists: %w", err))
	}
	var filters []string
	if filter != "" {
		filters = []string{filter}
	}
	opts := pagination.PaginationOptions{Page: 1, PageSize: pageSize}
	return pagination.NewPortListIterator(cli, ctx, opts, filters...)
}

// Settings returns an iterator for settings.
func (cli *client) Settings(ctx context.Context, pageSize int, args ...filtering.FilterArg) *pagination.SettingsIterator {
	filter, err := filtering.BuildFilterFromArgs(args...)
	if err != nil {
		return pagination.NewFailedIterator[*commands.SettingsSetting](fmt.Errorf("Settings: %w", err))
	}
	var filters []string
	if filter != "" {
		filters = []string{filter}
	}
	opts := pagination.PaginationOptions{Page: 1, PageSize: pageSize}
	return pagination.NewSettingsIterator(cli, ctx, opts, filters...)
}
