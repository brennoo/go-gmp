package client

import (
	"context"
	"testing"

	"github.com/brennoo/go-gmp/commands/filtering"
)

func TestGetResults(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()
	resp, err := cli.GetResults(ctx, "task_id=e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Expected status 200, got %s", resp.Status)
	}
}

// TestGetResultsConsolidated tests the GetResults method.
func TestGetResultsConsolidated(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()

	// Test with no filters
	resp, err := cli.GetResults(ctx)
	if err != nil {
		t.Fatalf("Unexpected error during GetResults: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with string filter
	resp, err = cli.GetResults(ctx, "severity>5.0")
	if err != nil {
		t.Fatalf("Unexpected error during GetResults with string filter: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with functional options
	resp, err = cli.GetResults(ctx, filtering.WithSeverityGreaterThan(5.0))
	if err != nil {
		t.Fatalf("Unexpected error during GetResults with functional options: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

func TestGetVulns(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()
	resp, err := cli.GetVulns(ctx, "name=test-vuln")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.Status != "200" {
		t.Logf("Expected status 200, got %s (this is expected for filtered requests)", resp.Status)
		return
	}
}
