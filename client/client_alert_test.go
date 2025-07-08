package client

import (
	"testing"

	"github.com/brennoo/go-gmp"
)

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
