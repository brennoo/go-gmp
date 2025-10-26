package client

import (
	"context"
	"testing"
	"time"
)

func TestGetTasksPaged(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Test basic pagination
	resp, err := cli.GetTasksPaged(1, 10)
	if err != nil {
		t.Fatalf("Unexpected error during GetTasksPaged: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with filters
	resp, err = cli.GetTasksPaged(1, 5, "status=running", "sort=name")
	if err != nil {
		t.Fatalf("Unexpected error during GetTasksPaged with filters: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

func TestGetResultsPaged(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Test basic pagination
	resp, err := cli.GetResultsPaged(1, 10)
	if err != nil {
		t.Fatalf("Unexpected error during GetResultsPaged: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with filters
	resp, err = cli.GetResultsPaged(2, 5, "severity>5.0", "status=Done")
	if err != nil {
		t.Fatalf("Unexpected error during GetResultsPaged with filters: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

func TestGetAssetsPaged(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Test basic pagination
	resp, err := cli.GetAssetsPaged(1, 10)
	if err != nil {
		t.Fatalf("Unexpected error during GetAssetsPaged: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with filters
	resp, err = cli.GetAssetsPaged(1, 5, "name~localhost")
	if err != nil {
		t.Fatalf("Unexpected error during GetAssetsPaged with filters: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

func TestGetTargetsPaged(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Test basic pagination
	resp, err := cli.GetTargetsPaged(1, 10)
	if err != nil {
		t.Fatalf("Unexpected error during GetTargetsPaged: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with filters
	resp, err = cli.GetTargetsPaged(1, 5, "name~test")
	if err != nil {
		t.Fatalf("Unexpected error during GetTargetsPaged with filters: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

func TestGetTicketsPaged(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Test basic pagination
	resp, err := cli.GetTicketsPaged(1, 10)
	if err != nil {
		t.Fatalf("Unexpected error during GetTicketsPaged: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with filters
	resp, err = cli.GetTicketsPaged(1, 5, "status=Open")
	if err != nil {
		t.Fatalf("Unexpected error during GetTicketsPaged with filters: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

func TestGetPortListsPaged(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Test basic pagination
	resp, err := cli.GetPortListsPaged(1, 10)
	if err != nil {
		t.Fatalf("Unexpected error during GetPortListsPaged: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with filters
	resp, err = cli.GetPortListsPaged(1, 5, "name~default")
	if err != nil {
		t.Fatalf("Unexpected error during GetPortListsPaged with filters: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

func TestGetSettingsPaged(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Test basic pagination
	resp, err := cli.GetSettingsPaged(1, 10)
	if err != nil {
		t.Fatalf("Unexpected error during GetSettingsPaged: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with filters
	resp, err = cli.GetSettingsPaged(1, 5, "name~timeout")
	if err != nil {
		t.Fatalf("Unexpected error during GetSettingsPaged with filters: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

func TestGetTasksIter(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test basic iteration
	iter := cli.GetTasksIter(ctx, 10)
	if iter == nil {
		t.Fatal("GetTasksIter returned nil iterator")
	}
	defer iter.Close()

	// Test iteration
	if !iter.Next() {
		t.Fatal("Expected Next() to return true")
	}

	task := iter.Current()
	if task == nil {
		t.Fatal("Current() returned nil task")
	}
	_ = task // Use the variable to avoid unused variable error

	// Test iterator state
	if !iter.HasMore() {
		t.Error("HasMore() should return true")
	}
	if iter.Err() != nil {
		t.Errorf("Error() should return nil, got: %v", iter.Err())
	}
	if iter.Total() < 0 {
		t.Error("Total() should return non-negative value")
	}
}

func TestGetResultsIter(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test basic iteration
	iter := cli.GetResultsIter(ctx, 10)
	if iter == nil {
		t.Fatal("GetResultsIter returned nil iterator")
	}
	defer iter.Close()

	// Test iteration
	if !iter.Next() {
		t.Fatal("Expected Next() to return true")
	}

	result := iter.Current()
	if result == nil {
		t.Fatal("Current() returned nil result")
	}
	_ = result // Use the variable to avoid unused variable error

	// Test iterator state
	if !iter.HasMore() {
		t.Error("HasMore() should return true")
	}
	if iter.Err() != nil {
		t.Errorf("Error() should return nil, got: %v", iter.Err())
	}
	if iter.Total() < 0 {
		t.Error("Total() should return non-negative value")
	}
}

func TestGetAssetsIter(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test basic iteration
	iter := cli.GetAssetsIter(ctx, 10)
	if iter == nil {
		t.Fatal("GetAssetsIter returned nil iterator")
	}
	defer iter.Close()

	// Test iteration
	if !iter.Next() {
		t.Fatal("Expected Next() to return true")
	}

	asset := iter.Current()
	if asset == nil {
		t.Fatal("Current() returned nil asset")
	}
	_ = asset // Use the variable to avoid unused variable error

	// Test iterator state
	if !iter.HasMore() {
		t.Error("HasMore() should return true")
	}
	if iter.Err() != nil {
		t.Errorf("Error() should return nil, got: %v", iter.Err())
	}
	if iter.Total() < 0 {
		t.Error("Total() should return non-negative value")
	}
}

func TestGetTargetsIter(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test basic iteration
	iter := cli.GetTargetsIter(ctx, 10)
	if iter == nil {
		t.Fatal("GetTargetsIter returned nil iterator")
	}
	defer iter.Close()

	// Test iteration
	if !iter.Next() {
		t.Fatal("Expected Next() to return true")
	}

	target := iter.Current()
	if target == nil {
		t.Fatal("Current() returned nil target")
	}
	_ = target // Use the variable to avoid unused variable error

	// Test iterator state
	if !iter.HasMore() {
		t.Error("HasMore() should return true")
	}
	if iter.Err() != nil {
		t.Errorf("Error() should return nil, got: %v", iter.Err())
	}
	if iter.Total() < 0 {
		t.Error("Total() should return non-negative value")
	}
}

func TestGetTicketsIter(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test basic iteration
	iter := cli.GetTicketsIter(ctx, 10)
	if iter == nil {
		t.Fatal("GetTicketsIter returned nil iterator")
	}
	defer iter.Close()

	// Test iteration
	if !iter.Next() {
		t.Fatal("Expected Next() to return true")
	}

	ticket := iter.Current()
	if ticket == nil {
		t.Fatal("Current() returned nil ticket")
	}
	_ = ticket // Use the variable to avoid unused variable error

	// Test iterator state
	if !iter.HasMore() {
		t.Error("HasMore() should return true")
	}
	if iter.Err() != nil {
		t.Errorf("Error() should return nil, got: %v", iter.Err())
	}
	if iter.Total() < 0 {
		t.Error("Total() should return non-negative value")
	}
}

func TestGetPortListsIter(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test basic iteration
	iter := cli.GetPortListsIter(ctx, 10)
	if iter == nil {
		t.Fatal("GetPortListsIter returned nil iterator")
	}
	defer iter.Close()

	// Test iteration
	if !iter.Next() {
		t.Fatal("Expected Next() to return true")
	}

	portList := iter.Current()
	if portList == nil {
		t.Fatal("Current() returned nil portList")
	}
	_ = portList // Use the variable to avoid unused variable error

	// Test iterator state
	if !iter.HasMore() {
		t.Error("HasMore() should return true")
	}
	if iter.Err() != nil {
		t.Errorf("Error() should return nil, got: %v", iter.Err())
	}
	if iter.Total() < 0 {
		t.Error("Total() should return non-negative value")
	}
}

func TestGetSettingsIter(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test basic iteration
	iter := cli.GetSettingsIter(ctx, 10)
	if iter == nil {
		t.Fatal("GetSettingsIter returned nil iterator")
	}
	defer iter.Close()

	// Test iteration
	if !iter.Next() {
		t.Fatal("Expected Next() to return true")
	}

	setting := iter.Current()
	if setting == nil {
		t.Fatal("Current() returned nil setting")
	}
	_ = setting // Use the variable to avoid unused variable error

	// Test iterator state
	if !iter.HasMore() {
		t.Error("HasMore() should return true")
	}
	if iter.Err() != nil {
		t.Errorf("Error() should return nil, got: %v", iter.Err())
	}
	if iter.Total() < 0 {
		t.Error("Total() should return non-negative value")
	}
}

func TestIteratorWithFilters(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test iterator with filters
	iter := cli.GetTasksIter(ctx, 5, "status=running", "sort=name")
	if iter == nil {
		t.Fatal("GetTasksIter with filters returned nil iterator")
	}
	defer iter.Close()

	// Test iteration
	if !iter.Next() {
		t.Fatal("Expected Next() to return true")
	}

	task := iter.Current()
	if task == nil {
		t.Fatal("Current() with filters returned nil task")
	}
	_ = task // Use the variable to avoid unused variable error
}

func TestIteratorContextCancellation(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	iter := cli.GetTasksIter(ctx, 10)
	if iter == nil {
		t.Fatal("GetTasksIter returned nil iterator")
	}
	defer iter.Close()

	// The iterator should handle cancelled context gracefully
	// We can't easily test the actual cancellation without a real connection,
	// but we can test that the iterator is created successfully
	if iter == nil {
		t.Fatal("Iterator should be created even with cancelled context")
	}
}

func TestIteratorClose(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()

	// Test that Close() doesn't panic
	iter := cli.GetTasksIter(ctx, 10)
	if iter == nil {
		t.Fatal("GetTasksIter returned nil iterator")
	}

	// Close should not panic
	iter.Close()

	// Test that HasMore returns false after close
	if iter.HasMore() {
		t.Error("HasMore() should return false after Close()")
	}
}

func TestPaginationHelperMethods(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Test that all pagination helper methods exist and can be called
	tests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "GetTasksPaged",
			fn: func() error {
				_, err := cli.GetTasksPaged(1, 10)
				return err
			},
		},
		{
			name: "GetResultsPaged",
			fn: func() error {
				_, err := cli.GetResultsPaged(1, 10)
				return err
			},
		},
		{
			name: "GetAssetsPaged",
			fn: func() error {
				_, err := cli.GetAssetsPaged(1, 10)
				return err
			},
		},
		{
			name: "GetTargetsPaged",
			fn: func() error {
				_, err := cli.GetTargetsPaged(1, 10)
				return err
			},
		},
		{
			name: "GetTicketsPaged",
			fn: func() error {
				_, err := cli.GetTicketsPaged(1, 10)
				return err
			},
		},
		{
			name: "GetPortListsPaged",
			fn: func() error {
				_, err := cli.GetPortListsPaged(1, 10)
				return err
			},
		},
		{
			name: "GetSettingsPaged",
			fn: func() error {
				_, err := cli.GetSettingsPaged(1, 10)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); err != nil {
				t.Errorf("%s failed: %v", tt.name, err)
			}
		})
	}
}

func TestIteratorMethods(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()

	// Test that all iterator methods exist and can be called
	t.Run("GetTasksIter", func(t *testing.T) {
		iter := cli.GetTasksIter(ctx, 10)
		if iter == nil {
			t.Error("GetTasksIter returned nil iterator")
		} else {
			iter.Close() // Clean up
		}
	})

	t.Run("GetResultsIter", func(t *testing.T) {
		iter := cli.GetResultsIter(ctx, 10)
		if iter == nil {
			t.Error("GetResultsIter returned nil iterator")
		} else {
			iter.Close() // Clean up
		}
	})

	t.Run("GetAssetsIter", func(t *testing.T) {
		iter := cli.GetAssetsIter(ctx, 10)
		if iter == nil {
			t.Error("GetAssetsIter returned nil iterator")
		} else {
			iter.Close() // Clean up
		}
	})

	t.Run("GetTargetsIter", func(t *testing.T) {
		iter := cli.GetTargetsIter(ctx, 10)
		if iter == nil {
			t.Error("GetTargetsIter returned nil iterator")
		} else {
			iter.Close() // Clean up
		}
	})

	t.Run("GetTicketsIter", func(t *testing.T) {
		iter := cli.GetTicketsIter(ctx, 10)
		if iter == nil {
			t.Error("GetTicketsIter returned nil iterator")
		} else {
			iter.Close() // Clean up
		}
	})

	t.Run("GetPortListsIter", func(t *testing.T) {
		iter := cli.GetPortListsIter(ctx, 10)
		if iter == nil {
			t.Error("GetPortListsIter returned nil iterator")
		} else {
			iter.Close() // Clean up
		}
	})

	t.Run("GetSettingsIter", func(t *testing.T) {
		iter := cli.GetSettingsIter(ctx, 10)
		if iter == nil {
			t.Error("GetSettingsIter returned nil iterator")
		} else {
			iter.Close() // Clean up
		}
	})
}
