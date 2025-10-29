package pagination

import (
	"context"
	"errors"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestTargetIterator(t *testing.T) {
	// Create test client
	cli := &testClient{}

	// Create iterator
	ctx := context.Background()
	iterator := &TargetIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	// Test iteration
	var targets []*commands.Target
	for iterator.Next() {
		target := iterator.Current()
		targets = append(targets, target)
	}

	// Verify we got all targets
	if len(targets) != 3 {
		t.Errorf("Expected 3 targets, got %d", len(targets))
	}

	// Verify total count
	if iterator.Total() != 3 {
		t.Errorf("Expected total 3, got %d", iterator.Total())
	}

	// Test error handling
	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}

	// Test individual target properties
	if len(targets) > 0 {
		if targets[0].ID != "target-1" {
			t.Errorf("Expected first target ID 'target-1', got '%s'", targets[0].ID)
		}
		if targets[0].Name != "Test Target 1" {
			t.Errorf("Expected first target name 'Test Target 1', got '%s'", targets[0].Name)
		}
	}
}

func TestTargetIterator_PaginationSuccess(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &TargetIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 1},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	// Test first page
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for first page")
	}

	target := iterator.Current()
	if target == nil {
		t.Fatal("Expected Current() to return a target")
	}

	if target.ID != "target-1" {
		t.Errorf("Expected target ID 'target-1', got '%s'", target.ID)
	}

	// Test second page
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for second page")
	}

	target = iterator.Current()
	if target == nil {
		t.Fatal("Expected Current() to return a target")
	}

	if target.ID != "target-2" {
		t.Errorf("Expected target ID 'target-2', got '%s'", target.ID)
	}
}

func TestTargetIterator_EmptyFirstPageStops(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &TargetIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 0},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	// Should not advance on empty page
	if iterator.Next() {
		t.Error("Expected Next() to return false for empty page")
	}

	if iterator.Current() != nil {
		t.Error("Expected Current() to return nil for empty page")
	}
}

func TestTargetIterator_ErrorStops(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &TargetIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	// Simulate error by setting err field directly
	iterator.err = errors.New("test error")

	if iterator.Next() {
		t.Error("Expected Next() to return false when error is set")
	}

	if iterator.Err() == nil {
		t.Error("Expected Err() to return the error")
	}
}

func TestTargetIterator_ContextCanceled(t *testing.T) {
	cli := &testClient{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	iterator := &TargetIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	if iterator.Next() {
		t.Error("Expected Next() to return false when context is canceled")
	}

	if iterator.Err() == nil {
		t.Error("Expected Err() to return context canceled error")
	}
}

func TestTargetIterator_CurrentBeforeNext(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &TargetIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	// Current() should return nil before Next() is called
	if iterator.Current() != nil {
		t.Error("Expected Current() to return nil before Next() is called")
	}
}

func TestTargetIterator_CloseStopsIteration(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &TargetIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	// Close the iterator
	iterator.Close()

	// Should not advance after close
	if iterator.Next() {
		t.Error("Expected Next() to return false after Close()")
	}

	// HasMore should return false after close
	if iterator.HasMore() {
		t.Error("Expected HasMore() to return false after Close()")
	}
}
