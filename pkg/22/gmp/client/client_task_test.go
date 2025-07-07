package client

import (
	"testing"

	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

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
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
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
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
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
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
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
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
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
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
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
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
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
