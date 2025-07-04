package client

import (
	"testing"

	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

type mockConn struct{}

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
