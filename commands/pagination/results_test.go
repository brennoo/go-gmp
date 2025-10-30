package pagination

import (
	"context"
	"testing"
)

func TestResultIterator(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewResultIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 3})

	var results []string
	for iterator.Next() {
		result := iterator.Current()
		results = append(results, result.ID)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	if iterator.Total() != 3 {
		t.Errorf("Expected total 3, got %d", iterator.Total())
	}
}

func TestResultIterator_PaginationSuccess(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewResultIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 1})

	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for first page")
	}
	if iterator.Current() == nil {
		t.Fatal("Expected Current() to return a result")
	}
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for second page")
	}
}

func TestResultIterator_EmptyFirstPageStops(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewResultIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 0})

	if iterator.Next() {
		t.Error("Expected Next() to return false for empty page")
	}
	if iterator.Current() != nil {
		t.Error("Expected Current() to return nil for empty page")
	}
}
