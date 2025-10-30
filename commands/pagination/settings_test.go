package pagination

import (
	"context"
	"errors"
	"testing"
)

func TestSettingsIterator(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewSettingsIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 3})

	var ids []string
	for iterator.Next() {
		s := iterator.Current()
		ids = append(ids, s.ID)
	}

	if len(ids) != 3 {
		t.Errorf("Expected 3 settings, got %d", len(ids))
	}

	if iterator.Total() != 3 {
		t.Errorf("Expected total 3, got %d", iterator.Total())
	}
}

func TestSettingsIterator_PaginationSuccess(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewSettingsIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 1})

	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for first page")
	}
	if iterator.Current() == nil {
		t.Fatal("Expected Current() to return a setting")
	}
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for second page")
	}
}

func TestSettingsIterator_EmptyFirstPageStops(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewSettingsIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 0})

	if iterator.Next() {
		t.Error("Expected Next() to return false for empty page")
	}
	if iterator.Current() != nil {
		t.Error("Expected Current() to return nil for empty page")
	}
}

func TestSettingsIterator_ErrorStops(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewSettingsIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 2})

	// Simulate error by setting err field directly
	iterator.err = errors.New("test error")

	if iterator.Next() {
		t.Error("Expected Next() to return false when error is set")
	}

	if iterator.Err() == nil {
		t.Error("Expected Err() to return the error")
	}
}

func TestSettingsIterator_ContextCanceled(t *testing.T) {
	cli := &testClient{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	iterator := NewSettingsIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 2})

	if iterator.Next() {
		t.Error("Expected Next() to return false when context is canceled")
	}

	if iterator.Err() == nil {
		t.Error("Expected Err() to return context canceled error")
	}
}

func TestSettingsIterator_CurrentBeforeNext(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewSettingsIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 2})

	// Current() should return nil before Next() is called
	if iterator.Current() != nil {
		t.Error("Expected Current() to return nil before Next() is called")
	}
}

func TestSettingsIterator_CloseStopsIteration(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewSettingsIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 2})

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
