package pagination

import (
	"context"
	"errors"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestResultIterator(t *testing.T) {
	// Create test client
	cli := &testClient{}

	// Create iterator
	ctx := context.Background()
	iterator := &ResultIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	// Test iteration
	var results []*commands.Result
	for iterator.Next() {
		result := iterator.Current()
		results = append(results, result)
	}

	// Verify we got all results
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// Verify total count
	if iterator.Total() != 3 {
		t.Errorf("Expected total 3, got %d", iterator.Total())
	}

	// Test error handling
	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}

	// Test individual result properties
	if len(results) > 0 {
		if results[0].ID != "result-1" {
			t.Errorf("Expected first result ID 'result-1', got '%s'", results[0].ID)
		}
		if results[0].Name != "Test Result 1" {
			t.Errorf("Expected first result name 'Test Result 1', got '%s'", results[0].Name)
		}
	}
}

func TestResultIterator_PaginationSuccess(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &ResultIterator{
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

	result := iterator.Current()
	if result == nil {
		t.Fatal("Expected Current() to return a result")
	}

	if result.ID != "result-1" {
		t.Errorf("Expected result ID 'result-1', got '%s'", result.ID)
	}

	// Test second page
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for second page")
	}

	result = iterator.Current()
	if result == nil {
		t.Fatal("Expected Current() to return a result")
	}

	if result.ID != "result-2" {
		t.Errorf("Expected result ID 'result-2', got '%s'", result.ID)
	}
}

func TestResultIterator_EmptyFirstPageStops(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &ResultIterator{
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

func TestResultIterator_ErrorStops(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &ResultIterator{
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

func TestResultIterator_ContextCanceled(t *testing.T) {
	cli := &testClient{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	iterator := &ResultIterator{
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

func TestResultIterator_CurrentBeforeNext(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &ResultIterator{
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

func TestResultIterator_CloseStopsIteration(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &ResultIterator{
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
