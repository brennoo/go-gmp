package pagination

import (
	"context"
	"errors"
)

// Iterator is a generic, stateful paginator over GMP list endpoints.
type Iterator[T any] struct {
	Client      Client
	Ctx         context.Context
	Opts        PaginationOptions
	Filters     []string
	current     []T
	index       int
	Page        int  // Current page number (1-based) fetched, next page to fetch is Page+1
	HasMoreData bool // Flag indicating if more pages might exist
	total       int
	err         error
	closed      bool

	// Type-specific functions provided by the factory
	buildCommand    func(filter string) interface{}
	executeAndParse func(cmd interface{}, pageOpts PaginationOptions) (items []T, totalFiltered int, hasMore bool, err error)
}

// loadPage fetches the next page using configured functions.
func (it *Iterator[T]) loadPage() error {
	// Ensure type-specific functions are provided
	if it.buildCommand == nil || it.executeAndParse == nil {
		err := errors.New("iterator created without required buildCommand or executeAndParse functions; use a factory function (e.g., NewTaskIterator)")
		it.err = err
		return err
	}

	// Build pagination filter for the next page
	nextPage := it.Page + 1
	opts := PaginationOptions{
		Page:     nextPage,
		PageSize: it.Opts.PageSize,
		MaxItems: it.Opts.MaxItems,
	}
	filter := BuildPaginationFilter(opts, it.Filters...)

	// Build and execute command
	cmd := it.buildCommand(filter)
	items, totalFiltered, hasMore, err := it.executeAndParse(cmd, opts)
	if err != nil {
		it.err = err
		return err
	}

	// Update internal state
	it.current = items
	it.total = totalFiltered
	it.HasMoreData = hasMore
	it.Page = nextPage
	it.index = 0
	return nil
}

// Next advances to the next item; returns false when exhausted or on error.
func (it *Iterator[T]) Next() bool {
	if it.closed || it.err != nil {
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

	if it.index >= len(it.current) {
		if !it.HasMoreData && it.Page > 0 {
			return false
		}
		if err := it.loadPage(); err != nil {
			return false
		}
		if len(it.current) == 0 {
			it.HasMoreData = false
			return false
		}
	}

	it.index++
	return true
}

// Current returns the current item (zero value if out of bounds).
func (it *Iterator[T]) Current() T {
	var zero T
	if it.index > 0 && it.index <= len(it.current) {
		return it.current[it.index-1]
	}
	return zero
}

// Err returns the last error (context error has precedence).
func (it *Iterator[T]) Err() error {
	if it.Ctx != nil && it.Ctx.Err() != nil {
		return it.Ctx.Err()
	}
	return it.err
}

// HasMore reports whether more items may be available.
func (it *Iterator[T]) HasMore() bool {
	if it.closed {
		return false
	}
	return (it.index < len(it.current)) || it.HasMoreData
}

// Total returns the total number of filtered items if known.
func (it *Iterator[T]) Total() int { return it.total }

// Close releases internal buffers and stops further iteration.
func (it *Iterator[T]) Close() {
	it.current = nil
	it.index = 0
	it.HasMoreData = false
	it.closed = true
}
