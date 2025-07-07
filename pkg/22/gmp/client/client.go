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

func (cli *client) CreateCredential(cmd *gmp.CreateCredentialCommand) (*gmp.CreateCredentialResponse, error) {
	var resp gmp.CreateCredentialResponse
	err := cli.conn.Execute(cmd, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (cli *client) ModifyCredential(cmd *gmp.ModifyCredentialCommand) (*gmp.ModifyCredentialResponse, error) {
	var resp gmp.ModifyCredentialResponse
	err := cli.conn.Execute(cmd, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (cli *client) GetCredentials(cmd *gmp.GetCredentialsCommand) (*gmp.GetCredentialsResponse, error) {
	var resp gmp.GetCredentialsResponse
	err := cli.conn.Execute(cmd, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (cli *client) DeleteCredential(cmd *gmp.DeleteCredentialCommand) (*gmp.DeleteCredentialResponse, error) {
	var resp gmp.DeleteCredentialResponse
	err := cli.conn.Execute(cmd, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
