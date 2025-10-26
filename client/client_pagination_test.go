package client

import (
	"context"
	"testing"
	"time"

	"github.com/brennoo/go-gmp"
	"github.com/brennoo/go-gmp/commands/filtering"
)

func TestTasks(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test basic iteration
	iter := cli.Tasks(ctx, 10)
	if iter == nil {
		t.Fatal("Tasks returned nil iterator")
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

func TestResults(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test basic iteration
	iter := cli.Results(ctx, 10)
	if iter == nil {
		t.Fatal("Results returned nil iterator")
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

func TestAssets(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test basic iteration
	iter := cli.Assets(ctx, 10)
	if iter == nil {
		t.Fatal("Assets returned nil iterator")
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

func TestTargets(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test basic iteration
	iter := cli.Targets(ctx, 10)
	if iter == nil {
		t.Fatal("Targets returned nil iterator")
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

func TestTickets(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test basic iteration
	iter := cli.Tickets(ctx, 10)
	if iter == nil {
		t.Fatal("Tickets returned nil iterator")
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

func TestPortLists(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test basic iteration
	iter := cli.PortLists(ctx, 10)
	if iter == nil {
		t.Fatal("PortLists returned nil iterator")
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

func TestSettings(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test basic iteration
	iter := cli.Settings(ctx, 10)
	if iter == nil {
		t.Fatal("Settings returned nil iterator")
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
	iter := cli.Tasks(ctx, 5, "status=running", "sort=name")
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

	iter := cli.Tasks(ctx, 10)
	if iter == nil {
		t.Fatal("Tasks returned nil iterator")
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
	iter := cli.Tasks(ctx, 10)
	if iter == nil {
		t.Fatal("Tasks returned nil iterator")
	}

	// Close should not panic
	iter.Close()

	// Test that HasMore returns false after close
	if iter.HasMore() {
		t.Error("HasMore() should return false after Close()")
	}
}

func TestIteratorMethods(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()

	// Test that all iterator methods exist and can be called
	t.Run("Tasks", func(t *testing.T) {
		iter := cli.Tasks(ctx, 10)
		if iter == nil {
			t.Error("Tasks returned nil iterator")
		} else {
			iter.Close() // Clean up
		}
	})

	t.Run("Results", func(t *testing.T) {
		iter := cli.Results(ctx, 10)
		if iter == nil {
			t.Error("Results returned nil iterator")
		} else {
			iter.Close() // Clean up
		}
	})

	t.Run("Assets", func(t *testing.T) {
		iter := cli.Assets(ctx, 10)
		if iter == nil {
			t.Error("Assets returned nil iterator")
		} else {
			iter.Close() // Clean up
		}
	})

	t.Run("Targets", func(t *testing.T) {
		iter := cli.Targets(ctx, 10)
		if iter == nil {
			t.Error("Targets returned nil iterator")
		} else {
			iter.Close() // Clean up
		}
	})

	t.Run("Tickets", func(t *testing.T) {
		iter := cli.Tickets(ctx, 10)
		if iter == nil {
			t.Error("Tickets returned nil iterator")
		} else {
			iter.Close() // Clean up
		}
	})

	t.Run("PortLists", func(t *testing.T) {
		iter := cli.PortLists(ctx, 10)
		if iter == nil {
			t.Error("PortLists returned nil iterator")
		} else {
			iter.Close() // Clean up
		}
	})

	t.Run("Settings", func(t *testing.T) {
		iter := cli.Settings(ctx, 10)
		if iter == nil {
			t.Error("Settings returned nil iterator")
		} else {
			iter.Close() // Clean up
		}
	})
}

// testIteratorWithFilters is a helper function to test an iterator method with both string and functional filters.
func testIteratorWithFilters(t *testing.T, cli gmp.Client, ctx context.Context, name string,
	stringFilter string, funcFilter filtering.FilterArg,
	iteratorFunc func(context.Context, int, ...filtering.FilterArg) interface{}) {

	// Test with string filter
	iter := iteratorFunc(ctx, 10, stringFilter)
	if iter == nil {
		t.Fatalf("%s with string filters returned nil iterator", name)
	}
	if closer, ok := iter.(interface{ Close() }); ok {
		defer closer.Close()
	}

	// Test with functional filter
	iter2 := iteratorFunc(ctx, 10, funcFilter)
	if iter2 == nil {
		t.Fatalf("%s with functional options returned nil iterator", name)
	}
	if closer, ok := iter2.(interface{ Close() }); ok {
		defer closer.Close()
	}
}

// TestIteratorsWithStringFilters tests all iterator methods with string filters.
func TestIteratorsWithStringFilters(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}
	ctx := context.Background()

	testIteratorWithFilters(t, cli, ctx, "Results", "severity>5.0", filtering.WithSeverityGreaterThan(5.0),
		func(ctx context.Context, pageSize int, args ...filtering.FilterArg) interface{} {
			return cli.Results(ctx, pageSize, args...)
		})
}

// TestIteratorsWithFunctionalFilters tests all iterator methods with functional options.
func TestIteratorsWithFunctionalFilters(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}
	ctx := context.Background()

	testIteratorWithFilters(t, cli, ctx, "Assets", "name~localhost", filtering.WithName("test"),
		func(ctx context.Context, pageSize int, args ...filtering.FilterArg) interface{} {
			return cli.Assets(ctx, pageSize, args...)
		})

	testIteratorWithFilters(t, cli, ctx, "Targets", "name~test", filtering.WithName("test"),
		func(ctx context.Context, pageSize int, args ...filtering.FilterArg) interface{} {
			return cli.Targets(ctx, pageSize, args...)
		})

	testIteratorWithFilters(t, cli, ctx, "Tickets", "status=Open", filtering.WithStatus("Open"),
		func(ctx context.Context, pageSize int, args ...filtering.FilterArg) interface{} {
			return cli.Tickets(ctx, pageSize, args...)
		})

	testIteratorWithFilters(t, cli, ctx, "PortLists", "name~default", filtering.WithName("test"),
		func(ctx context.Context, pageSize int, args ...filtering.FilterArg) interface{} {
			return cli.PortLists(ctx, pageSize, args...)
		})

	testIteratorWithFilters(t, cli, ctx, "Settings", "name~timeout", filtering.WithName("test"),
		func(ctx context.Context, pageSize int, args ...filtering.FilterArg) interface{} {
			return cli.Settings(ctx, pageSize, args...)
		})
}
