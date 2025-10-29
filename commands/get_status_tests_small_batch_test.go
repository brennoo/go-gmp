package commands

import "testing"

func TestGetTasksResponse_GetStatus(t *testing.T) {
	r := &GetTasksResponse{Status: "200", StatusText: "OK"}
	status, text := r.GetStatus()
	if status != "200" || text != "OK" {
		t.Fatalf("unexpected status: %s %s", status, text)
	}
}

func TestGetResultsResponse_GetStatus(t *testing.T) {
	r := &GetResultsResponse{Status: "404", StatusText: "Not Found"}
	status, text := r.GetStatus()
	if status != "404" || text != "Not Found" {
		t.Fatalf("unexpected status: %s %s", status, text)
	}
}

func TestGetAssetsResponse_GetStatus(t *testing.T) {
	r := &GetAssetsResponse{Status: "401", StatusText: "Unauthorized"}
	status, text := r.GetStatus()
	if status != "401" || text != "Unauthorized" {
		t.Fatalf("unexpected status: %s %s", status, text)
	}
}

func TestGetTargetsResponse_GetStatus(t *testing.T) {
	r := &GetTargetsResponse{Status: "500", StatusText: "Internal Server Error"}
	status, text := r.GetStatus()
	if status != "500" || text != "Internal Server Error" {
		t.Fatalf("unexpected status: %s %s", status, text)
	}
}

func TestGetTicketsResponse_GetStatus(t *testing.T) {
	r := &GetTicketsResponse{Status: "403", StatusText: "Forbidden"}
	status, text := r.GetStatus()
	if status != "403" || text != "Forbidden" {
		t.Fatalf("unexpected status: %s %s", status, text)
	}
}
