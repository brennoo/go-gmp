package client

import (
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestCreateSchedule(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.CreateSchedule{
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

	cmd := &commands.ModifySchedule{
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

	cmd := &commands.GetSchedules{
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

	cmd := &commands.DeleteSchedule{
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
