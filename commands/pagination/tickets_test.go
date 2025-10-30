package pagination

import (
	"context"
	"testing"
)

func TestTicketIterator(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewTicketIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 3})

	var tickets []string
	for iterator.Next() {
		t := iterator.Current()
		tickets = append(tickets, t.ID)
	}

	if len(tickets) != 3 {
		t.Errorf("Expected 3 tickets, got %d", len(tickets))
	}

	if iterator.Total() != 3 {
		t.Errorf("Expected total 3, got %d", iterator.Total())
	}
}

func TestTicketIterator_PaginationSuccess(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewTicketIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 1})

	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for first page")
	}
	if iterator.Current() == nil {
		t.Fatal("Expected Current() to return a ticket")
	}
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for second page")
	}
}

func TestTicketIterator_EmptyFirstPageStops(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewTicketIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 0})

	if iterator.Next() {
		t.Error("Expected Next() to return false for empty page")
	}
	if iterator.Current() != nil {
		t.Error("Expected Current() to return nil for empty page")
	}
}
