package client

import (
	"context"
	"testing"
	"time"

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

// TestIteratorsWithFilters tests all iterator methods with both string filters and functional options.
func TestIteratorsWithFilters(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()

	// Test Results with filters
	iter := cli.Results(ctx, 10, "severity>5.0")
	if iter == nil {
		t.Fatal("Results with string filters returned nil iterator")
	}
	defer iter.Close()

	iter2 := cli.Results(ctx, 10, filtering.WithSeverityGreaterThan(5.0))
	if iter2 == nil {
		t.Fatal("Results with functional options returned nil iterator")
	}
	defer iter2.Close()

	// Test Assets with filters
	iter3 := cli.Assets(ctx, 10, "name~localhost")
	if iter3 == nil {
		t.Fatal("Assets with string filters returned nil iterator")
	}
	defer iter3.Close()

	iter4 := cli.Assets(ctx, 10, filtering.WithName("test"))
	if iter4 == nil {
		t.Fatal("Assets with functional options returned nil iterator")
	}
	defer iter4.Close()

	// Test Targets with filters
	iter5 := cli.Targets(ctx, 10, "name~test")
	if iter5 == nil {
		t.Fatal("Targets with string filters returned nil iterator")
	}
	defer iter5.Close()

	iter6 := cli.Targets(ctx, 10, filtering.WithName("test"))
	if iter6 == nil {
		t.Fatal("Targets with functional options returned nil iterator")
	}
	defer iter6.Close()

	// Test Tickets with filters
	iter7 := cli.Tickets(ctx, 10, "status=Open")
	if iter7 == nil {
		t.Fatal("Tickets with string filters returned nil iterator")
	}
	defer iter7.Close()

	iter8 := cli.Tickets(ctx, 10, filtering.WithStatus("Open"))
	if iter8 == nil {
		t.Fatal("Tickets with functional options returned nil iterator")
	}
	defer iter8.Close()

	// Test PortLists with filters
	iter9 := cli.PortLists(ctx, 10, "name~default")
	if iter9 == nil {
		t.Fatal("PortLists with string filters returned nil iterator")
	}
	defer iter9.Close()

	iter10 := cli.PortLists(ctx, 10, filtering.WithName("test"))
	if iter10 == nil {
		t.Fatal("PortLists with functional options returned nil iterator")
	}
	defer iter10.Close()

	// Test Settings with filters
	iter11 := cli.Settings(ctx, 10, "name~timeout")
	if iter11 == nil {
		t.Fatal("Settings with string filters returned nil iterator")
	}
	defer iter11.Close()

	iter12 := cli.Settings(ctx, 10, filtering.WithName("test"))
	if iter12 == nil {
		t.Fatal("Settings with functional options returned nil iterator")
	}
	defer iter12.Close()
}
