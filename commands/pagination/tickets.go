//nolint:dupl
// Similar patterns across iterators are intentional

package pagination

import (
	"context"

	"github.com/brennoo/go-gmp/commands"
)

// TicketIterator provides iteration over tickets with automatic pagination.
type TicketIterator struct {
	Client      Client
	Ctx         context.Context
	Opts        PaginationOptions
	Filters     []string
	current     []*commands.Ticket
	index       int
	Page        int
	HasMoreData bool
	total       int
	err         error
}

// Next advances the iterator to the next item and returns true if successful.
//
//nolint:dupl // Similar pattern across iterators is intentional
func (it *TicketIterator) Next() bool {
	if it.err != nil {
		return false
	}

	// Check context
	if it.Ctx != nil {
		select {
		case <-it.Ctx.Done():
			it.err = it.Ctx.Err()
			return false
		default:
		}
	}

	// Check if we need to load more data
	if it.index >= len(it.current) {
		if !it.HasMoreData {
			return false // No more data
		}

		// Load the next page
		it.Page++
		it.index = 0
		it.err = it.loadPage()
		if it.err != nil {
			return false
		}

		// If we got no items, return false
		if len(it.current) == 0 {
			return false
		}
	}

	// Return the next item
	it.index++
	return true
}

// Current returns the current ticket item.
func (it *TicketIterator) Current() *commands.Ticket {
	if it.index > 0 && it.index <= len(it.current) {
		return it.current[it.index-1]
	}
	return nil
}

// Err returns any error that occurred during iteration.
func (it *TicketIterator) Err() error {
	return it.err
}

// HasMore returns true if there are more items available.
func (it *TicketIterator) HasMore() bool {
	return it.index < len(it.current) || it.HasMoreData
}

// Total returns the total number of items available.
func (it *TicketIterator) Total() int {
	return it.total
}

// Close cleans up the iterator.
func (it *TicketIterator) Close() {
	it.current = nil
	it.index = 0
	it.HasMoreData = false
}

func (it *TicketIterator) loadPage() error {
	opts := PaginationOptions{
		Page:     it.Page,
		PageSize: it.Opts.PageSize,
	}
	filter := BuildPaginationFilter(opts, it.Filters...)

	cmd := &commands.GetTickets{
		Filter: filter,
	}

	resp, err := it.Client.GetTicketsRaw(cmd)
	if err != nil {
		return err
	}

	// Convert tickets to typed slice
	it.current = make([]*commands.Ticket, len(resp.Tickets))
	for i, ticket := range resp.Tickets {
		it.current[i] = &ticket
	}

	// Update total if available
	if resp.TicketCount.Filtered > 0 {
		it.total = resp.TicketCount.Filtered
	}

	// Check if there are more pages - if we got fewer items than PageSize, we're done
	// If we got exactly PageSize items, there might be more pages
	it.HasMoreData = len(resp.Tickets) == it.Opts.PageSize

	return nil
}
