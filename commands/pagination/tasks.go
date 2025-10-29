//nolint:dupl
// Similar patterns across iterators are intentional

package pagination

import (
	"context"

	"github.com/brennoo/go-gmp/commands"
)

// TaskIterator provides iteration over tasks with automatic pagination.
type TaskIterator struct {
	Client      Client
	Ctx         context.Context
	Opts        PaginationOptions
	Filters     []string
	current     []*commands.GetTasksResponseTask
	index       int
	Page        int
	HasMoreData bool
	total       int
	err         error
}

// Next advances the iterator to the next item and returns true if successful.
//
//nolint:dupl // Similar pattern across iterators is intentional
func (it *TaskIterator) Next() bool {
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

// Current returns the current task item.
func (it *TaskIterator) Current() *commands.GetTasksResponseTask {
	if it.index > 0 && it.index <= len(it.current) {
		return it.current[it.index-1]
	}
	return nil
}

// Err returns any error that occurred during iteration.
func (it *TaskIterator) Err() error {
	return it.err
}

// HasMore returns true if there are more items available.
func (it *TaskIterator) HasMore() bool {
	return it.index < len(it.current) || it.HasMoreData
}

// Total returns the total number of items available.
func (it *TaskIterator) Total() int {
	return it.total
}

// Close cleans up the iterator.
func (it *TaskIterator) Close() {
	it.current = nil
	it.index = 0
	it.HasMoreData = false
}

//nolint:dupl // Similar pattern across iterators is intentional
func (it *TaskIterator) loadPage() error {
	opts := PaginationOptions{
		Page:     it.Page,
		PageSize: it.Opts.PageSize,
	}
	filter := BuildPaginationFilter(opts, it.Filters...)

	cmd := &commands.GetTasks{
		Filter: filter,
	}

	resp, err := it.Client.GetTasksRaw(cmd)
	if err != nil {
		return err
	}

	// Convert tasks to typed slice
	it.current = make([]*commands.GetTasksResponseTask, len(resp.Task))
	for i, task := range resp.Task {
		it.current[i] = &task
	}

	// Update total if available (using TaskCount for now)
	if resp.TaskCount.Filtered > 0 {
		it.total = resp.TaskCount.Filtered
	}

	// Check if there are more pages - if we got fewer items than PageSize, we're done
	// If we got exactly PageSize items, there might be more pages
	it.HasMoreData = len(resp.Task) == it.Opts.PageSize

	return nil
}
