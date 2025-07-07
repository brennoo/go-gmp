package client

import (
	"testing"

	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

type mockConn struct{}

// nolint:gocyclo // This mock handles all GMP commands
// gocyclo:ignore
func (m *mockConn) Execute(command interface{}, response interface{}) error {
	if cmd, ok := command.(*gmp.AuthenticateCommand); ok {
		if cmd.Credentials.Username == "openvas" && cmd.Credentials.Password == "123" {
			(*response.(*gmp.AuthenticateResponse)).Status = "200"
		} else {
			(*response.(*gmp.AuthenticateResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetConfigsCommand); ok {
		if cmd.ConfigID == "bde773f3-2b3d-4fe6-81cb-6321ae2cc629" {
			(*response.(*gmp.GetConfigsResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetConfigsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetScannersCommand); ok {
		if cmd.ScannerID == "ee0311e7-3247-4425-bb9c-866d59f1e0e9" {
			(*response.(*gmp.GetScannersResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetScannersResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetPreferencesCommand); ok {
		if cmd.ConfigID == "4b49617e-d1d8-44b8-af81-f4675b56f837" {
			(*response.(*gmp.GetPreferencesResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetPreferencesResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreateConfigCommand); ok {
		if cmd.Name == "New Config" {
			(*response.(*gmp.CreateConfigResponse)).Status = "200"
		} else {
			(*response.(*gmp.CreateConfigResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.ModifyConfigCommand); ok {
		if cmd.Name == "Modified Config" {
			(*response.(*gmp.ModifyConfigResponse)).Status = "200"
		} else {
			(*response.(*gmp.ModifyConfigResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreateTaskCommand); ok {
		if cmd.Name == "New Task" {
			(*response.(*gmp.CreateTaskResponse)).Status = "200"
		} else {
			(*response.(*gmp.CreateTaskResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.ModifyTaskCommand); ok {
		if cmd.Comment == "Modified Task Comment" {
			(*response.(*gmp.ModifyTaskResponse)).Status = "200"
		} else {
			(*response.(*gmp.ModifyTaskResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreateTargetCommand); ok {
		if cmd.Name == "New Target" {
			(*response.(*gmp.CreateTargetResponse)).Status = "200"
		} else {
			(*response.(*gmp.CreateTargetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.ModifyTargetCommand); ok {
		if cmd.Name == "Modified Target" {
			(*response.(*gmp.ModifyTargetResponse)).Status = "200"
		} else {
			(*response.(*gmp.ModifyTargetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetTargetsCommand); ok {
		if cmd.TargetID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*gmp.GetTargetsResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetTargetsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.DeleteTargetCommand); ok {
		if cmd.TargetID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*gmp.DeleteTargetResponse)).Status = "200"
		} else {
			(*response.(*gmp.DeleteTargetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetPortListsCommand); ok {
		if cmd.PortListID == "33d0cd82-57c6-11e1-8ed1-406186ea4fc5" {
			(*response.(*gmp.GetPortListsResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetPortListsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.StartTaskCommand); ok {
		if cmd.TaskID == "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31" {
			(*response.(*gmp.StartTaskResponse)).Status = "200"
		} else {
			(*response.(*gmp.StartTaskResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetTasksCommand); ok {
		if cmd.TaskID == "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31" {
			(*response.(*gmp.GetTasksResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetTasksResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.StopTaskCommand); ok {
		if cmd.TaskID == "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31" {
			(*response.(*gmp.StopTaskResponse)).Status = "200"
		} else {
			(*response.(*gmp.StopTaskResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.DeleteTaskCommand); ok {
		if cmd.TaskID == "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31" {
			(*response.(*gmp.DeleteTaskResponse)).Status = "200"
		} else {
			(*response.(*gmp.DeleteTaskResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetResultsCommand); ok {
		if cmd.TaskID == "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31" {
			(*response.(*gmp.GetResultsResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetResultsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetVulnsCommand); ok {
		if cmd.VulnID == "1.3.6.1.4.1.25623.1.0.808160" {
			(*response.(*gmp.GetVulnsResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetVulnsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreateAlertRequest); ok {
		if cmd.Name == "Test Alert" {
			(*response.(*gmp.CreateAlertResponse)).Status = "201"
			(*response.(*gmp.CreateAlertResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateAlertResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*gmp.CreateAlertResponse)).Status = "400"
		}
	}

	if _, ok := command.(*gmp.GetAlertsCommand); ok {
		(*response.(*gmp.GetAlertsResponse)).Status = "200"
		(*response.(*gmp.GetAlertsResponse)).StatusText = "OK"
	}

	if _, ok := command.(*gmp.ModifyAlertCommand); ok {
		(*response.(*gmp.ModifyAlertResponse)).Status = "200"
		(*response.(*gmp.ModifyAlertResponse)).StatusText = "OK"
	}

	if _, ok := command.(*gmp.DeleteAlertCommand); ok {
		(*response.(*gmp.DeleteAlertResponse)).Status = "200"
		(*response.(*gmp.DeleteAlertResponse)).StatusText = "OK"
	}

	if _, ok := command.(*gmp.TestAlertCommand); ok {
		(*response.(*gmp.TestAlertResponse)).Status = "200"
		(*response.(*gmp.TestAlertResponse)).StatusText = "OK"
	}

	if _, ok := command.(*gmp.ResumeTaskCommand); ok {
		(*response.(*gmp.ResumeTaskResponse)).Status = "200"
		(*response.(*gmp.ResumeTaskResponse)).StatusText = "OK"
		(*response.(*gmp.ResumeTaskResponse)).ReportID = "330ee785-c2c0-4d4c-ab96-725142c9b789"
	}

	if cmd, ok := command.(*gmp.CreateAssetCommand); ok {
		if cmd.Asset != nil && cmd.Asset.Name == "Localhost" && cmd.Asset.Type == "host" {
			(*response.(*gmp.CreateAssetResponse)).Status = "201"
			(*response.(*gmp.CreateAssetResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateAssetResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else if cmd.Report != nil && cmd.Report.ID == "report-uuid" && cmd.Report.Filter != nil && cmd.Report.Filter.Term == "min_qod=70" {
			(*response.(*gmp.CreateAssetResponse)).Status = "201"
			(*response.(*gmp.CreateAssetResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateAssetResponse)).ID = "report-asset-uuid"
		} else {
			(*response.(*gmp.CreateAssetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.ModifyAssetCommand); ok {
		if cmd.AssetID == "914b59f8-25f5-4c8f-832c-2379cd625236" && cmd.Comment == "New comment" {
			(*response.(*gmp.ModifyAssetResponse)).Status = "200"
			(*response.(*gmp.ModifyAssetResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyAssetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetAssetsCommand); ok {
		if cmd.AssetID == "b493b7a8-7489-11df-a3ec-002264764cea" {
			(*response.(*gmp.GetAssetsResponse)).Status = "200"
			(*response.(*gmp.GetAssetsResponse)).StatusText = "OK"
			(*response.(*gmp.GetAssetsResponse)).Assets = []gmp.Asset{
				{
					ID:               "b493b7a8-7489-11df-a3ec-002264764cea",
					Name:             "Localhost",
					Comment:          "",
					CreationTime:     "2018-08-29T20:21:33Z",
					ModificationTime: "2018-08-29T20:21:33Z",
					Type:             "host",
				},
			}
		} else {
			(*response.(*gmp.GetAssetsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.DeleteAssetCommand); ok {
		if cmd.AssetID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" {
			(*response.(*gmp.DeleteAssetResponse)).Status = "200"
			(*response.(*gmp.DeleteAssetResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteAssetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreateScheduleCommand); ok {
		if cmd.Name == "Monthly Scan" && cmd.Timezone == "America/New_York" && cmd.ICalendar != "" {
			(*response.(*gmp.CreateScheduleResponse)).Status = "201"
			(*response.(*gmp.CreateScheduleResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateScheduleResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*gmp.CreateScheduleResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.ModifyScheduleCommand); ok {
		if cmd.ScheduleID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" && cmd.Name == "Weekly Scan" {
			(*response.(*gmp.ModifyScheduleResponse)).Status = "200"
			(*response.(*gmp.ModifyScheduleResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyScheduleResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetSchedulesCommand); ok {
		if cmd.ScheduleID == "ddda859a-45be-4c58-85b3-517c66230232" {
			(*response.(*gmp.GetSchedulesResponse)).Status = "200"
			(*response.(*gmp.GetSchedulesResponse)).StatusText = "OK"
			(*response.(*gmp.GetSchedulesResponse)).Schedules = []gmp.Schedule{
				{
					ID:               "ddda859a-45be-4c58-85b3-517c66230232",
					Name:             "Every day",
					Comment:          "",
					CreationTime:     "2020-06-03T16:27:05Z",
					ModificationTime: "2020-06-03T16:27:05Z",
					Writable:         "1",
					InUse:            "0",
					ICalendar:        "DTSTART:20200603T162600Z DURATION:PT0S RRULE:FREQ=DAILY",
					Timezone:         "UTC",
				},
			}
		} else {
			(*response.(*gmp.GetSchedulesResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.DeleteScheduleCommand); ok {
		if cmd.ScheduleID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" {
			(*response.(*gmp.DeleteScheduleResponse)).Status = "200"
			(*response.(*gmp.DeleteScheduleResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteScheduleResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreateOverrideCommand); ok {
		if cmd.Text == "This is actually of little concern." && cmd.NVT.OID == "1.3.6.1.4.1.25623.1.0.10330" && cmd.NewThreat == "Low" && cmd.Result == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*gmp.CreateOverrideResponse)).Status = "201"
			(*response.(*gmp.CreateOverrideResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateOverrideResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*gmp.CreateOverrideResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.DeleteOverrideCommand); ok {
		if cmd.OverrideID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" {
			(*response.(*gmp.DeleteOverrideResponse)).Status = "200"
			(*response.(*gmp.DeleteOverrideResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteOverrideResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetOverridesCommand); ok {
		if cmd.OverrideID == "b76b81a7-9df8-42df-afff-baa9d4620128" {
			(*response.(*gmp.GetOverridesResponse)).Status = "200"
			(*response.(*gmp.GetOverridesResponse)).StatusText = "OK"
			(*response.(*gmp.GetOverridesResponse)).Overrides = []gmp.Override{
				{
					ID:          "b76b81a7-9df8-42df-afff-baa9d4620128",
					Text:        gmp.OverrideText{Text: "This is the full text of the override."},
					NewThreat:   "Log",
					NewSeverity: "0.0",
					Orphan:      true,
					Active:      true,
					NVT:         gmp.OverrideNVT{OID: "1.3.6.1.4.1.25623.1.0.75", Name: "Test NVT: long lines"},
				},
			}
		} else {
			(*response.(*gmp.GetOverridesResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.ModifyOverrideCommand); ok {
		if cmd.OverrideID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" && cmd.Text == "This issue is less important in our setup." && cmd.NewThreat == "Low" {
			(*response.(*gmp.ModifyOverrideResponse)).Status = "200"
			(*response.(*gmp.ModifyOverrideResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyOverrideResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreateReportCommand); ok {
		if cmd.Report != nil && cmd.Report.ID == "report-uuid" && cmd.Task != nil && cmd.Task.ID == "task-uuid" {
			(*response.(*gmp.CreateReportResponse)).Status = "201"
			(*response.(*gmp.CreateReportResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateReportResponse)).ID = "created-report-id"
		} else {
			(*response.(*gmp.CreateReportResponse)).Status = "400"
			(*response.(*gmp.CreateReportResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.GetReportsCommand); ok {
		if cmd.ReportID == "report-uuid" && cmd.FormatID == "format-uuid" {
			(*response.(*gmp.GetReportsResponse)).Status = "200"
			(*response.(*gmp.GetReportsResponse)).StatusText = "OK"
			(*response.(*gmp.GetReportsResponse)).Reports = []gmp.ReportWrapper{
				{
					ID:          "report-uuid",
					FormatID:    "format-uuid",
					Extension:   "xml",
					ContentType: "text/xml",
					RawXML:      []byte("<report>...</report>"),
				},
			}
		} else {
			(*response.(*gmp.GetReportsResponse)).Status = "404"
			(*response.(*gmp.GetReportsResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.DeleteReportCommand); ok {
		if cmd.ReportID == "report-uuid" {
			(*response.(*gmp.DeleteReportResponse)).Status = "200"
			(*response.(*gmp.DeleteReportResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteReportResponse)).Status = "404"
			(*response.(*gmp.DeleteReportResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.CreateReportFormatCommand); ok {
		if cmd.Name == "Test Format" && cmd.Copy == "copy-uuid" {
			(*response.(*gmp.CreateReportFormatResponse)).Status = "201"
			(*response.(*gmp.CreateReportFormatResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateReportFormatResponse)).ID = "created-format-id"
		} else {
			(*response.(*gmp.CreateReportFormatResponse)).Status = "400"
			(*response.(*gmp.CreateReportFormatResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.ModifyReportFormatCommand); ok {
		if cmd.ReportFormatID == "format-uuid" && cmd.Name == "New Name" {
			(*response.(*gmp.ModifyReportFormatResponse)).Status = "200"
			(*response.(*gmp.ModifyReportFormatResponse)).StatusText = "OK"
		} else if cmd.ReportFormatID == "format-uuid" && cmd.Param != nil && cmd.Param.Name == "Background Colour" && cmd.Param.Value == "red" {
			(*response.(*gmp.ModifyReportFormatResponse)).Status = "200"
			(*response.(*gmp.ModifyReportFormatResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyReportFormatResponse)).Status = "400"
			(*response.(*gmp.ModifyReportFormatResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.GetReportFormatsCommand); ok {
		if cmd.ReportFormatID == "format-uuid" {
			(*response.(*gmp.GetReportFormatsResponse)).Status = "200"
			(*response.(*gmp.GetReportFormatsResponse)).StatusText = "OK"
			(*response.(*gmp.GetReportFormatsResponse)).ReportFormats = []gmp.ReportFormatWrapper{
				{
					ID:               "format-uuid",
					Name:             "HTML",
					Extension:        "html",
					ContentType:      "text/html",
					Summary:          "Single page HTML report.",
					Description:      "A single HTML page listing results of a scan.",
					CreationTime:     "2013-01-31T16:46:32+01:00",
					ModificationTime: "2013-01-31T16:46:32+01:00",
					Writable:         "1",
					InUse:            "0",
					Active:           "1",
				},
			}
		} else {
			(*response.(*gmp.GetReportFormatsResponse)).Status = "404"
			(*response.(*gmp.GetReportFormatsResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.DeleteReportFormatCommand); ok {
		if cmd.ReportFormatID == "format-uuid" && cmd.Ultimate == "1" {
			(*response.(*gmp.DeleteReportFormatResponse)).Status = "200"
			(*response.(*gmp.DeleteReportFormatResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteReportFormatResponse)).Status = "404"
			(*response.(*gmp.DeleteReportFormatResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.VerifyReportFormatCommand); ok {
		if cmd.ReportFormatID == "format-uuid" {
			(*response.(*gmp.VerifyReportFormatResponse)).Status = "200"
			(*response.(*gmp.VerifyReportFormatResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.VerifyReportFormatResponse)).Status = "404"
			(*response.(*gmp.VerifyReportFormatResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.CreateReportConfigCommand); ok {
		if cmd.Name == "Test config" && cmd.ReportFormat.ID == "format-uuid" {
			(*response.(*gmp.CreateReportConfigResponse)).Status = "201"
			(*response.(*gmp.CreateReportConfigResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateReportConfigResponse)).ID = "created-config-id"
		} else {
			(*response.(*gmp.CreateReportConfigResponse)).Status = "400"
			(*response.(*gmp.CreateReportConfigResponse)).StatusText = "Bad request"
		}
	}

	return nil
}

func (m *mockConn) Close() error {
	return nil
}

func mockedConnection() gmp.Connection {
	return &mockConn{}
}

func TestNew(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}
}

func TestAuthenticate(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.AuthenticateCommand{}
	cmd.Credentials.Username = "openvas"
	cmd.Credentials.Password = "123"
	resp, err := cli.Authenticate(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during Authenticate: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestGetConfigs(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetConfigsCommand{}
	cmd.ConfigID = "bde773f3-2b3d-4fe6-81cb-6321ae2cc629"
	resp, err := cli.GetConfigs(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetConfigs: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestGetScanners(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetScannersCommand{}
	cmd.ScannerID = "ee0311e7-3247-4425-bb9c-866d59f1e0e9"
	resp, err := cli.GetScanners(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetScanners: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestGetPreferences(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetPreferencesCommand{}
	cmd.ConfigID = "4b49617e-d1d8-44b8-af81-f4675b56f837"
	resp, err := cli.GetPreferences(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetPreferences: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestCreateConfig(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateConfigCommand{}
	cmd.Name = "New Config"
	resp, err := cli.CreateConfig(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateConfig: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestModifyConfig(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.ModifyConfigCommand{}
	cmd.Name = "Modified Config"
	resp, err := cli.ModifyConfig(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyConfig: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestCreateTask(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateTaskCommand{}
	cmd.Name = "New Task"
	resp, err := cli.CreateTask(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateTask: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestModifyTask(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.ModifyTaskCommand{}
	cmd.TaskID = "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31"
	cmd.Comment = "Modified Task Comment"
	resp, err := cli.ModifyTask(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyTask: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestCreateTarget(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateTargetCommand{}
	cmd.Name = "New Target"
	resp, err := cli.CreateTarget(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateTarget: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestStartTask(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.StartTaskCommand{}
	cmd.TaskID = "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31"
	resp, err := cli.StartTask(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during StartTask: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestGetTasks(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetTasksCommand{}
	cmd.TaskID = "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31"
	resp, err := cli.GetTasks(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetTasks: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestStopTask(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.StopTaskCommand{}
	cmd.TaskID = "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31"
	resp, err := cli.StopTask(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during StopTask: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestDeleteTask(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.DeleteTaskCommand{}
	cmd.TaskID = "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31"
	resp, err := cli.DeleteTask(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteTask: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestGetResults(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetResultsCommand{}
	cmd.TaskID = "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31"
	resp, err := cli.GetResults(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetResults: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestGetVulns(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetVulnsCommand{}
	cmd.VulnID = "1.3.6.1.4.1.25623.1.0.808160"
	resp, err := cli.GetVulns(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetVulns: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestModifyTarget(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.ModifyTargetCommand{}
	cmd.TargetID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
	cmd.Name = "Modified Target"
	resp, err := cli.ModifyTarget(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyTarget: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestGetTargets(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetTargetsCommand{}
	cmd.TargetID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
	resp, err := cli.GetTargets(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetTargets: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestDeleteTarget(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.DeleteTargetCommand{}
	cmd.TargetID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
	resp, err := cli.DeleteTarget(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteTarget: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestGetPortLists(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetPortListsCommand{}
	cmd.PortListID = "33d0cd82-57c6-11e1-8ed1-406186ea4fc5"
	resp, err := cli.GetPortLists(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetPortLists: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestCreateAlert(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateAlertRequest{}
	cmd.Name = "Test Alert"
	cmd.Condition = gmp.AlertCondition{
		Text: "Severity at least",
		Data: []gmp.AlertData{
			{Name: "severity", Text: "5.5"},
		},
	}
	cmd.Event = gmp.AlertEvent{
		Text: "Task run status changed",
		Data: []gmp.AlertData{
			{Name: "status", Text: "Done"},
		},
	}
	cmd.Method = gmp.AlertMethod{
		Text: "Email",
		Data: []gmp.AlertData{
			{Name: "to_address", Text: "test@example.org"},
			{Name: "from_address", Text: "alert@example.org"},
		},
	}
	resp, err := cli.CreateAlert(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateAlert: %s", err)
	}

	if resp.Status != "201" {
		t.Fatalf("Unexpected status. \nExpected: 201 \nGot: %s", resp.Status)
	}

	if resp.StatusText != "OK, resource created" {
		t.Fatalf("Unexpected status text. \nExpected: OK, resource created \nGot: %s", resp.StatusText)
	}

	if resp.ID != "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
		t.Fatalf("Unexpected ID. \nExpected: 254cd3ef-bbe1-4d58-859d-21b8d0c046c6 \nGot: %s", resp.ID)
	}
}

func TestGetAlerts(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetAlertsCommand{}
	resp, err := cli.GetAlerts(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Status != "200" {
		t.Errorf("expected status '200', got %s", resp.Status)
	}
}

func TestModifyAlert(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	alertID := "914b59f8-25f5-4c8f-832c-2379cd625236"
	name := "Updated Alert Name"
	comment := "Updated comment"
	filter := &gmp.ModifyAlertFilter{ID: "7a06bd00-7e4a-4669-b7d7-8fe65ec64a41"}
	cmd := &gmp.ModifyAlertCommand{
		AlertID: alertID,
		Name:    &name,
		Comment: &comment,
		Filter:  filter,
	}
	resp, err := cli.ModifyAlert(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "200" {
		t.Errorf("expected status '200', got %s", resp.Status)
	}
}

func TestDeleteAlert(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	alertID := "267a3405-e84a-47da-97b2-5fa0d2e8995e"
	cmd := &gmp.DeleteAlertCommand{
		AlertID:  alertID,
		Ultimate: true,
	}
	resp, err := cli.DeleteAlert(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "200" {
		t.Errorf("expected status '200', got %s", resp.Status)
	}
}

func TestTestAlert(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	alertID := "97390ade-e075-11df-9973-002264764cea"
	cmd := &gmp.TestAlertCommand{
		AlertID: alertID,
	}
	resp, err := cli.TestAlert(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "200" {
		t.Errorf("expected status '200', got %s", resp.Status)
	}
}

func TestResumeTask(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	taskID := "267a3405-e84a-47da-97b2-5fa0d2e8995e"
	cmd := &gmp.ResumeTaskCommand{
		TaskID: taskID,
	}
	resp, err := cli.ResumeTask(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "200" {
		t.Errorf("expected status '200', got %s", resp.Status)
	}
	if resp.ReportID != "330ee785-c2c0-4d4c-ab96-725142c9b789" {
		t.Errorf("expected report_id '330ee785-c2c0-4d4c-ab96-725142c9b789', got %s", resp.ReportID)
	}
}

func TestCreateAsset(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateAssetCommand{
		Asset: &gmp.CreateAssetAsset{
			Name:    "Localhost",
			Type:    "host",
			Comment: "Test asset",
		},
	}
	resp, err := cli.CreateAsset(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateAsset: %s", err)
	}
	if resp.Status != "201" {
		t.Fatalf("Unexpected status. Expected: 201 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK, resource created" {
		t.Fatalf("Unexpected status text. Expected: OK, resource created Got: %s", resp.StatusText)
	}
	if resp.ID != "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
		t.Fatalf("Unexpected ID. Expected: 254cd3ef-bbe1-4d58-859d-21b8d0c046c6 Got: %s", resp.ID)
	}
}

func TestCreateAssetWithReport(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateAssetCommand{
		Report: &gmp.CreateAssetReport{
			ID: "report-uuid",
			Filter: &gmp.CreateAssetReportFilter{
				Term: "min_qod=70",
			},
		},
	}
	resp, err := cli.CreateAsset(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateAsset (report): %s", err)
	}
	if resp.Status != "201" {
		t.Fatalf("Unexpected status. Expected: 201 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK, resource created" {
		t.Fatalf("Unexpected status text. Expected: OK, resource created Got: %s", resp.StatusText)
	}
	if resp.ID != "report-asset-uuid" {
		t.Fatalf("Unexpected ID. Expected: report-asset-uuid Got: %s", resp.ID)
	}
}

func TestModifyAsset(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.ModifyAssetCommand{
		AssetID: "914b59f8-25f5-4c8f-832c-2379cd625236",
		Comment: "New comment",
	}
	resp, err := cli.ModifyAsset(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyAsset: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
}

func TestGetAssets(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetAssetsCommand{
		AssetID: "b493b7a8-7489-11df-a3ec-002264764cea",
	}
	resp, err := cli.GetAssets(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetAssets: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if len(resp.Assets) != 1 {
		t.Fatalf("Expected 1 asset, got %d", len(resp.Assets))
	}
	asset := resp.Assets[0]
	if asset.ID != "b493b7a8-7489-11df-a3ec-002264764cea" {
		t.Fatalf("Unexpected asset ID. Expected: b493b7a8-7489-11df-a3ec-002264764cea Got: %s", asset.ID)
	}
	if asset.Name != "Localhost" {
		t.Fatalf("Unexpected asset name. Expected: Localhost Got: %s", asset.Name)
	}
}

func TestDeleteAsset(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.DeleteAssetCommand{
		AssetID: "267a3405-e84a-47da-97b2-5fa0d2e8995e",
	}
	resp, err := cli.DeleteAsset(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteAsset: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
}

func TestCreateSchedule(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateScheduleCommand{
		Name:      "Monthly Scan",
		Timezone:  "America/New_York",
		ICalendar: "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//XXX//NONSGML//EN\nBEGIN:VEVENT\nDTSTART;TZID=\"America/New_York\":20221214T000100\nRRULE:FREQ=MONTHLY;BYDAY=3WE;\nEND:VEVENT\nEND:VCALENDAR",
	}
	resp, err := cli.CreateSchedule(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateSchedule: %s", err)
	}
	if resp.Status != "201" {
		t.Errorf("Expected status 201, got %s", resp.Status)
	}
	if resp.StatusText != "OK, resource created" {
		t.Errorf("Expected status text 'OK, resource created', got '%s'", resp.StatusText)
	}
	if resp.ID != "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
		t.Errorf("Expected ID '254cd3ef-bbe1-4d58-859d-21b8d0c046c6', got '%s'", resp.ID)
	}
}

func TestModifySchedule(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.ModifyScheduleCommand{
		ScheduleID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Name:       "Weekly Scan",
	}
	resp, err := cli.ModifySchedule(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifySchedule: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got '%s'", resp.StatusText)
	}
}

func TestGetSchedules(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetSchedulesCommand{
		ScheduleID: "ddda859a-45be-4c58-85b3-517c66230232",
	}
	resp, err := cli.GetSchedules(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetSchedules: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got '%s'", resp.StatusText)
	}
	if len(resp.Schedules) != 1 {
		t.Errorf("Expected 1 schedule, got %d", len(resp.Schedules))
	}
	if resp.Schedules[0].Name != "Every day" {
		t.Errorf("Expected schedule name 'Every day', got '%s'", resp.Schedules[0].Name)
	}
	if resp.Schedules[0].Timezone != "UTC" {
		t.Errorf("Expected timezone 'UTC', got '%s'", resp.Schedules[0].Timezone)
	}
}

func TestDeleteSchedule(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.DeleteScheduleCommand{
		ScheduleID: "267a3405-e84a-47da-97b2-5fa0d2e8995e",
		Ultimate:   "0",
	}
	resp, err := cli.DeleteSchedule(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteSchedule: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got '%s'", resp.StatusText)
	}
}

func TestCreateOverride(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateOverrideCommand{
		Text:      "This is actually of little concern.",
		NVT:       gmp.CreateOverrideNVT{OID: "1.3.6.1.4.1.25623.1.0.10330"},
		NewThreat: "Low",
		Result:    "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
	}
	resp, err := cli.CreateOverride(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateOverride: %s", err)
	}
	if resp.Status != "201" {
		t.Errorf("Expected status 201, got %s", resp.Status)
	}
	if resp.StatusText != "OK, resource created" {
		t.Errorf("Expected status text 'OK, resource created', got '%s'", resp.StatusText)
	}
	if resp.ID != "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
		t.Errorf("Expected ID '254cd3ef-bbe1-4d58-859d-21b8d0c046c6', got '%s'", resp.ID)
	}
}

func TestDeleteOverride(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.DeleteOverrideCommand{
		OverrideID: "267a3405-e84a-47da-97b2-5fa0d2e8995e",
		Ultimate:   "0",
	}
	resp, err := cli.DeleteOverride(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteOverride: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text OK, got %s", resp.StatusText)
	}
}

func TestGetOverrides(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.GetOverridesCommand{
		OverrideID: "b76b81a7-9df8-42df-afff-baa9d4620128",
	}
	resp, err := cli.GetOverrides(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetOverrides: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text OK, got %s", resp.StatusText)
	}
	if len(resp.Overrides) != 1 {
		t.Errorf("Expected 1 override, got %d", len(resp.Overrides))
	}
	if resp.Overrides[0].ID != "b76b81a7-9df8-42df-afff-baa9d4620128" {
		t.Errorf("Expected override ID b76b81a7-9df8-42df-afff-baa9d4620128, got %s", resp.Overrides[0].ID)
	}
}

func TestModifyOverride(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.ModifyOverrideCommand{
		OverrideID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Text:       "This issue is less important in our setup.",
		NewThreat:  "Low",
	}
	resp, err := cli.ModifyOverride(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyOverride: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text OK, got %s", resp.StatusText)
	}
}

func TestCreateReport(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.CreateReportCommand{
		Report: &gmp.ReportWrapper{
			ID:          "report-uuid",
			FormatID:    "format-uuid",
			Extension:   "xml",
			ContentType: "text/xml",
			RawXML:      []byte("<report>...</report>"),
		},
		Task: &gmp.CreateReportTask{
			ID: "task-uuid",
		},
		InAssets: nil,
	}
	resp, err := cli.CreateReport(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateReport: %s", err)
	}
	if resp.Status != "201" {
		t.Errorf("Expected status 201, got %s", resp.Status)
	}
	if resp.StatusText != "OK, resource created" {
		t.Errorf("Expected status text 'OK, resource created', got '%s'", resp.StatusText)
	}
	if resp.ID != "created-report-id" {
		t.Errorf("Expected ID 'created-report-id', got '%s'", resp.ID)
	}

	// Failure case
	cmdFail := &gmp.CreateReportCommand{
		Report: &gmp.ReportWrapper{
			ID: "wrong-id",
		},
		Task: &gmp.CreateReportTask{
			ID: "wrong-task",
		},
	}
	respFail, err := cli.CreateReport(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during CreateReport (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
	if respFail.StatusText != "Bad request" {
		t.Errorf("Expected status text 'Bad request', got '%s'", respFail.StatusText)
	}
}

func TestGetReports(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.GetReportsCommand{
		ReportID:    "report-uuid",
		FormatID:    "format-uuid",
		Extension:   "xml",
		ContentType: "text/xml",
	}
	resp, err := cli.GetReports(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetReports: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got '%s'", resp.StatusText)
	}
	if len(resp.Reports) != 1 {
		t.Errorf("Expected 1 report, got %d", len(resp.Reports))
	} else {
		report := resp.Reports[0]
		if report.ID != "report-uuid" {
			t.Errorf("Expected report ID 'report-uuid', got '%s'", report.ID)
		}
		if report.FormatID != "format-uuid" {
			t.Errorf("Expected format ID 'format-uuid', got '%s'", report.FormatID)
		}
		if report.Extension != "xml" {
			t.Errorf("Expected extension 'xml', got '%s'", report.Extension)
		}
		if report.ContentType != "text/xml" {
			t.Errorf("Expected content type 'text/xml', got '%s'", report.ContentType)
		}
	}

	// Failure case
	cmdFail := &gmp.GetReportsCommand{
		ReportID: "wrong-id",
	}
	respFail, err := cli.GetReports(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during GetReports (fail): %s", err)
	}
	if respFail.Status != "404" {
		t.Errorf("Expected status 404, got %s", respFail.Status)
	}
	if respFail.StatusText != "Not found" {
		t.Errorf("Expected status text 'Not found', got '%s'", respFail.StatusText)
	}
}

func TestDeleteReport(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.DeleteReportCommand{
		ReportID: "report-uuid",
	}
	resp, err := cli.DeleteReport(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteReport: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got '%s'", resp.StatusText)
	}

	// Failure case
	cmdFail := &gmp.DeleteReportCommand{
		ReportID: "wrong-id",
	}
	respFail, err := cli.DeleteReport(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteReport (fail): %s", err)
	}
	if respFail.Status != "404" {
		t.Errorf("Expected status 404, got %s", respFail.Status)
	}
	if respFail.StatusText != "Not found" {
		t.Errorf("Expected status text 'Not found', got '%s'", respFail.StatusText)
	}
}

func TestCreateReportFormat(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.CreateReportFormatCommand{
		Name: "Test Format",
		Copy: "copy-uuid",
	}
	resp, err := cli.CreateReportFormat(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateReportFormat: %s", err)
	}
	if resp.Status != "201" {
		t.Errorf("Expected status 201, got %s", resp.Status)
	}
	if resp.StatusText != "OK, resource created" {
		t.Errorf("Expected status text 'OK, resource created', got '%s'", resp.StatusText)
	}
	if resp.ID != "created-format-id" {
		t.Errorf("Expected ID 'created-format-id', got '%s'", resp.ID)
	}

	// Failure case
	cmdFail := &gmp.CreateReportFormatCommand{
		Name: "Wrong Format",
		Copy: "wrong-copy",
	}
	respFail, err := cli.CreateReportFormat(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during CreateReportFormat (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
	if respFail.StatusText != "Bad request" {
		t.Errorf("Expected status text 'Bad request', got '%s'", respFail.StatusText)
	}
}

func TestModifyReportFormat(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case: modify name
	cmd := &gmp.ModifyReportFormatCommand{
		ReportFormatID: "format-uuid",
		Name:           "New Name",
	}
	resp, err := cli.ModifyReportFormat(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyReportFormat: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got '%s'", resp.StatusText)
	}

	// Success case: modify param
	cmdParam := &gmp.ModifyReportFormatCommand{
		ReportFormatID: "format-uuid",
		Param: &gmp.ModifyReportFormatParam{
			Name:  "Background Colour",
			Value: "red",
		},
	}
	respParam, err := cli.ModifyReportFormat(cmdParam)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyReportFormat (param): %s", err)
	}
	if respParam.Status != "200" {
		t.Errorf("Expected status 200, got %s", respParam.Status)
	}
	if respParam.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got '%s'", respParam.StatusText)
	}

	// Failure case
	cmdFail := &gmp.ModifyReportFormatCommand{
		ReportFormatID: "wrong-uuid",
		Name:           "Invalid",
	}
	respFail, err := cli.ModifyReportFormat(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyReportFormat (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
	if respFail.StatusText != "Bad request" {
		t.Errorf("Expected status text 'Bad request', got '%s'", respFail.StatusText)
	}
}

func TestGetReportFormats(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.GetReportFormatsCommand{
		ReportFormatID: "format-uuid",
	}
	resp, err := cli.GetReportFormats(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetReportFormats: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got '%s'", resp.StatusText)
	}
	if len(resp.ReportFormats) != 1 {
		t.Errorf("Expected 1 report format, got %d", len(resp.ReportFormats))
	} else {
		format := resp.ReportFormats[0]
		if format.ID != "format-uuid" {
			t.Errorf("Expected report format ID 'format-uuid', got '%s'", format.ID)
		}
		if format.Name != "HTML" {
			t.Errorf("Expected name 'HTML', got '%s'", format.Name)
		}
		if format.Extension != "html" {
			t.Errorf("Expected extension 'html', got '%s'", format.Extension)
		}
		if format.ContentType != "text/html" {
			t.Errorf("Expected content type 'text/html', got '%s'", format.ContentType)
		}
	}

	// Failure case
	cmdFail := &gmp.GetReportFormatsCommand{
		ReportFormatID: "wrong-id",
	}
	respFail, err := cli.GetReportFormats(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during GetReportFormats (fail): %s", err)
	}
	if respFail.Status != "404" {
		t.Errorf("Expected status 404, got %s", respFail.Status)
	}
	if respFail.StatusText != "Not found" {
		t.Errorf("Expected status text 'Not found', got '%s'", respFail.StatusText)
	}
}

func TestDeleteReportFormat(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.DeleteReportFormatCommand{
		ReportFormatID: "format-uuid",
		Ultimate:       "1",
	}
	resp, err := cli.DeleteReportFormat(cmd)
	if err != nil {
		t.Fatalf("DeleteReportFormat failed: %v", err)
	}
	if resp.Status != "200" || resp.StatusText != "OK" {
		t.Errorf("unexpected response: %+v", resp)
	}

	// Failure case
	cmd = &gmp.DeleteReportFormatCommand{
		ReportFormatID: "notfound",
		Ultimate:       "1",
	}
	resp, err = cli.DeleteReportFormat(cmd)
	if err != nil {
		t.Fatalf("DeleteReportFormat failed: %v", err)
	}
	if resp.Status != "404" || resp.StatusText != "Not found" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestVerifyReportFormat(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.VerifyReportFormatCommand{
		ReportFormatID: "format-uuid",
	}
	resp, err := cli.VerifyReportFormat(cmd)
	if err != nil {
		t.Fatalf("VerifyReportFormat failed: %v", err)
	}
	if resp.Status != "200" || resp.StatusText != "OK" {
		t.Errorf("unexpected response: %+v", resp)
	}

	// Failure case
	cmd = &gmp.VerifyReportFormatCommand{
		ReportFormatID: "notfound",
	}
	resp, err = cli.VerifyReportFormat(cmd)
	if err != nil {
		t.Fatalf("VerifyReportFormat failed: %v", err)
	}
	if resp.Status != "404" || resp.StatusText != "Not found" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestCreateReportConfig(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &gmp.CreateReportConfigCommand{
		Name: "Test config",
		ReportFormat: gmp.CreateReportConfigFormat{
			ID: "format-uuid",
		},
		Params: []gmp.CreateReportConfigParam{
			{Name: "Node Distance", Value: "10"},
		},
		Comment: "Test comment",
	}
	resp, err := cli.CreateReportConfig(cmd)
	if err != nil {
		t.Fatalf("CreateReportConfig failed: %v", err)
	}
	if resp.Status != "201" || resp.StatusText != "OK, resource created" || resp.ID != "created-config-id" {
		t.Errorf("unexpected response: %+v", resp)
	}

	// Failure case
	cmd = &gmp.CreateReportConfigCommand{
		Name: "",
		ReportFormat: gmp.CreateReportConfigFormat{
			ID: "",
		},
	}
	resp, err = cli.CreateReportConfig(cmd)
	if err != nil {
		t.Fatalf("CreateReportConfig failed: %v", err)
	}
	if resp.Status != "400" || resp.StatusText != "Bad request" {
		t.Errorf("unexpected response: %+v", resp)
	}
}
