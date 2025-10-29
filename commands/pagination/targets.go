//nolint:dupl
// Similar patterns across iterators are intentional

package pagination

import (
	"context"

	"github.com/brennoo/go-gmp/commands"
)

// TargetIterator provides iteration over targets with automatic pagination.
type TargetIterator struct {
	Client      Client
	Ctx         context.Context
	Opts        PaginationOptions
	Filters     []string
	current     []*commands.Target
	index       int
	Page        int
	HasMoreData bool
	total       int
	err         error
}

// Next advances the iterator to the next item and returns true if successful.
//
//nolint:dupl // Similar pattern across iterators is intentional
func (it *TargetIterator) Next() bool {
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

// Current returns the current target item.
func (it *TargetIterator) Current() *commands.Target {
	if it.index > 0 && it.index <= len(it.current) {
		return it.current[it.index-1]
	}
	return nil
}

// Err returns any error that occurred during iteration.
func (it *TargetIterator) Err() error {
	return it.err
}

// HasMore returns true if there are more items available.
func (it *TargetIterator) HasMore() bool {
	return it.index < len(it.current) || it.HasMoreData
}

// Total returns the total number of items available.
func (it *TargetIterator) Total() int {
	return it.total
}

// Close cleans up the iterator.
func (it *TargetIterator) Close() {
	it.current = nil
	it.index = 0
	it.HasMoreData = false
}

func (it *TargetIterator) loadPage() error {
	opts := PaginationOptions{
		Page:     it.Page,
		PageSize: it.Opts.PageSize,
	}
	filter := BuildPaginationFilter(opts, it.Filters...)

	cmd := &commands.GetTargets{
		Filter: filter,
	}

	resp, err := it.Client.GetTargetsRaw(cmd)
	if err != nil {
		return err
	}

	// Convert targets to typed slice
	it.current = make([]*commands.Target, len(resp.Target))
	for i, target := range resp.Target {
		it.current[i] = &target
	}

	// Update total if available
	if resp.TargetCount != nil && resp.TargetCount.Filtered > 0 {
		it.total = resp.TargetCount.Filtered
	}

	// Check if there are more pages - if we got fewer items than PageSize, we're done
	// If we got exactly PageSize items, there might be more pages
	it.HasMoreData = len(resp.Target) == it.Opts.PageSize

	return nil
}
