package client

import (
	"testing"

	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

type mockConn struct{}

// nolint:gocyclo // This mock handles all GMP commands
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
