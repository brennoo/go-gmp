package client

import (
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestCreateAlert(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	condData := []commands.AlertData{
		commands.NewAlertData("severity", "5.5"),
	}
	evtData := []commands.AlertData{
		commands.NewAlertData("status", "Done"),
	}
	methData := []commands.AlertData{
		commands.NewAlertData("to_address", "test@example.org"),
		commands.NewAlertData("from_address", "alert@example.org"),
	}

	cmd := commands.NewCreateAlert(
		"Test Alert",
		"",
		"",
		commands.NewAlertCondition("Severity at least", condData),
		commands.NewAlertEvent("Task run status changed", evtData),
		commands.NewAlertMethod("Email", methData),
		nil, // filter
		nil, // active
	)
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

	cmd := &commands.GetAlerts{}
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
	filter := &commands.ModifyAlertFilter{ID: "7a06bd00-7e4a-4669-b7d7-8fe65ec64a41"}
	cmd := &commands.ModifyAlert{
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
	cmd := &commands.DeleteAlert{
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
	cmd := &commands.TestAlert{
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
