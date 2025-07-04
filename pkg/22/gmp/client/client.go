package client

import (
	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

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
