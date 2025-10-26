package client

import (
	"context"
	"testing"

	"github.com/brennoo/go-gmp/commands"
	"github.com/brennoo/go-gmp/commands/filtering"
)

func TestCreateTask(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.CreateTask{}
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
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.ModifyTask{}
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
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()
	resp, err := cli.GetTasks(ctx, "task_id=e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31")
	if err != nil {
		t.Fatalf("Unexpected error during GetTasks: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

// TestGetTasksConsolidated tests the GetTasks method.
func TestGetTasksConsolidated(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()

	// Test with no filters
	resp, err := cli.GetTasks(ctx)
	if err != nil {
		t.Fatalf("Unexpected error during GetTasks: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with string filter
	resp, err = cli.GetTasks(ctx, "status=running")
	if err != nil {
		t.Fatalf("Unexpected error during GetTasks with string filter: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with functional options
	resp, err = cli.GetTasks(ctx, filtering.WithStatus(filtering.StatusRunning))
	if err != nil {
		t.Fatalf("Unexpected error during GetTasks with functional options: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

func TestStartTask(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.StartTask{}
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
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.StopTask{}
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
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DeleteTask{}
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
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	taskID := "267a3405-e84a-47da-97b2-5fa0d2e8995e"
	cmd := &commands.ResumeTask{
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

func TestMoveTask(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.MoveTask{
		TaskID:  "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		SlaveID: "97390ade-e075-11df-9973-002264764cea",
	}
	resp, err := cli.MoveTask(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during MoveTask: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
}
