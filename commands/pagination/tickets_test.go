package pagination

import (
	"context"
	"errors"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestTicketIterator(t *testing.T) {
	// Create test client
	cli := &testClient{}

	// Create iterator
	ctx := context.Background()
	iterator := &TicketIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	// Test iteration
	var tickets []*commands.Ticket
	for iterator.Next() {
		ticket := iterator.Current()
		tickets = append(tickets, ticket)
	}

	// Verify we got all tickets
	if len(tickets) != 3 {
		t.Errorf("Expected 3 tickets, got %d", len(tickets))
	}

	// Verify total count
	if iterator.Total() != 3 {
		t.Errorf("Expected total 3, got %d", iterator.Total())
	}

	// Test error handling
	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}

	// Test individual ticket properties
	if len(tickets) > 0 {
		if tickets[0].ID != "ticket-1" {
			t.Errorf("Expected first ticket ID 'ticket-1', got '%s'", tickets[0].ID)
		}
		if tickets[0].Name != "Test Ticket 1" {
			t.Errorf("Expected first ticket name 'Test Ticket 1', got '%s'", tickets[0].Name)
		}
	}
}

func TestTicketIterator_PaginationSuccess(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &TicketIterator{
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

	ticket := iterator.Current()
	if ticket == nil {
		t.Fatal("Expected Current() to return a ticket")
	}

	if ticket.ID != "ticket-1" {
		t.Errorf("Expected ticket ID 'ticket-1', got '%s'", ticket.ID)
	}

	// Test second page
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for second page")
	}

	ticket = iterator.Current()
	if ticket == nil {
		t.Fatal("Expected Current() to return a ticket")
	}

	if ticket.ID != "ticket-2" {
		t.Errorf("Expected ticket ID 'ticket-2', got '%s'", ticket.ID)
	}
}

func TestTicketIterator_EmptyFirstPageStops(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &TicketIterator{
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

func TestTicketIterator_ErrorStops(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &TicketIterator{
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

func TestTicketIterator_ContextCanceled(t *testing.T) {
	cli := &testClient{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	iterator := &TicketIterator{
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

func TestTicketIterator_CurrentBeforeNext(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &TicketIterator{
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

func TestTicketIterator_CloseStopsIteration(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &TicketIterator{
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
