package pagination

import (
	"context"
	"errors"
	"testing"
)

// Test that constructing an iterator without factory functions yields the expected error and no iteration.
func TestIterator_RequiresFactoryFunctions(t *testing.T) {
	it := &Iterator[*struct{}]{
		Ctx:         context.Background(),
		Opts:        PaginationOptions{Page: 1, PageSize: 10},
		HasMoreData: true,
	}
	if it.Next() {
		t.Fatalf("expected Next() to be false when factory functions are missing")
	}
	err := it.Err()
	if err == nil || err.Error() == "" {
		t.Fatalf("expected a descriptive error when factory functions are missing, got: %v", err)
	}
}

// Test Close stops iteration and HasMore immediately.
func TestIterator_CloseStopsImmediately(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	it := NewTaskIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 2})
	it.Close()
	if it.Next() {
		t.Fatalf("expected Next() to be false after Close()")
	}
	if it.HasMore() {
		t.Fatalf("expected HasMore() to be false after Close()")
	}
}

// Test context cancelled before first Next prevents iteration and sets Err.
func TestIterator_ContextCancelledBeforeNext(t *testing.T) {
	cli := &testClient{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	it := NewTaskIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 2})
	if it.Next() {
		t.Fatalf("expected Next() to return false when context is cancelled")
	}
	if it.Err() == nil {
		t.Fatalf("expected context error from Err() when context is cancelled")
	}
}

// Test BuildPaginationFilter edge cases for defaults and clamping.
func TestBuildPaginationFilter_Edges(t *testing.T) {
	// Default Page to 1; default PageSize when negative
	f := BuildPaginationFilter(PaginationOptions{Page: 0, PageSize: -1}, "status=ok")
	if f == "" || f[:6] != "first=" {
		t.Fatalf("expected first/rows prefix, got: %q", f)
	}
	// rows=0 for PageSize==0
	f0 := BuildPaginationFilter(PaginationOptions{Page: 1, PageSize: 0})
	if f0 != "first=1 rows=0" {
		t.Fatalf("expected rows=0 for empty page, got: %q", f0)
	}
	// Clamp by MaxItems
	fClamp := BuildPaginationFilter(PaginationOptions{Page: 1, PageSize: 10, MaxItems: 5})
	if fClamp != "first=1 rows=5" {
		t.Fatalf("expected rows=5 when maxItems clamps, got: %q", fClamp)
	}
	// Past MaxItems → rows=0
	fPast := BuildPaginationFilter(PaginationOptions{Page: 2, PageSize: 10, MaxItems: 5})
	if fPast != "rows=0" {
		t.Fatalf("expected rows=0 when past MaxItems, got: %q", fPast)
	}
}

// Test loadPage bubbling error from executeAndParse.
func TestIterator_LoadPageError(t *testing.T) {
	errCli := &errorTestClient{}
	ctx := context.Background()
	it := NewTaskIterator(errCli, ctx, PaginationOptions{Page: 1, PageSize: 2})
	if it.Next() {
		t.Fatalf("expected Next() false on execute error")
	}
	if it.Err() == nil {
		t.Fatalf("expected error, got nil")
	}
	// ensure not context-cancelled path
	if errors.Is(it.Err(), context.Canceled) {
		t.Fatalf("unexpected context cancellation error")
	}
}
