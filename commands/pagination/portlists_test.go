package pagination

import (
	"context"
	"testing"
)

func TestPortListIterator(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewPortListIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 3})

	var ids []string
	for iterator.Next() {
		pl := iterator.Current()
		ids = append(ids, pl.ID)
	}

	if len(ids) != 3 {
		t.Errorf("Expected 3 port lists, got %d", len(ids))
	}

	if iterator.Total() != 3 {
		t.Errorf("Expected total 3, got %d", iterator.Total())
	}
}

func TestPortListIterator_PaginationSuccess(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewPortListIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 1})

	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for first page")
	}
	if iterator.Current() == nil {
		t.Fatal("Expected Current() to return a port list")
	}
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for second page")
	}
}

func TestPortListIterator_EmptyFirstPageStops(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewPortListIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 0})

	if iterator.Next() {
		t.Error("Expected Next() to return false for empty page")
	}
	if iterator.Current() != nil {
		t.Error("Expected Current() to return nil for empty page")
	}
}
