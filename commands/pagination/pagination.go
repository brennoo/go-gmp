//nolint:dupl // Similar patterns across iterators are intentional

package pagination

import (
	"fmt"

	"github.com/brennoo/go-gmp/commands"
)

// Client defines the minimal interface needed for pagination.
type Client interface {
	GetTasksRaw(cmd *commands.GetTasks) (*commands.GetTasksResponse, error)
	GetResultsRaw(cmd *commands.GetResults) (*commands.GetResultsResponse, error)
	GetAssetsRaw(cmd *commands.GetAssets) (*commands.GetAssetsResponse, error)
	GetTargetsRaw(cmd *commands.GetTargets) (*commands.GetTargetsResponse, error)
	GetTicketsRaw(cmd *commands.GetTickets) (*commands.GetTicketsResponse, error)
	GetPortListsRaw(cmd *commands.GetPortLists) (*commands.GetPortListsResponse, error)
	GetSettingsRaw(cmd *commands.GetSettings) (*commands.GetSettingsResponse, error)
}

// PaginationOptions represents pagination configuration.
type PaginationOptions struct {
	Page     int // 1-based page number
	PageSize int // Number of items per page
	MaxItems int // Maximum total items to fetch (0 = no limit)
}

// DefaultPaginationOptions returns sensible default pagination options.
func DefaultPaginationOptions() PaginationOptions {
	return PaginationOptions{
		Page:     1,
		PageSize: 100,
		MaxItems: 0, // No limit
	}
}

// BuildPaginationFilter constructs the filter string for pagination.
func BuildPaginationFilter(opts PaginationOptions, additionalFilters ...string) string {
	first := (opts.Page-1)*opts.PageSize + 1
	rows := opts.PageSize

	if opts.MaxItems > 0 && first > opts.MaxItems {
		return "rows=0" // No results if we're past the max
	}

	if opts.MaxItems > 0 && first+rows-1 > opts.MaxItems {
		rows = opts.MaxItems - first + 1
	}

	filter := fmt.Sprintf("first=%d rows=%d", first, rows)

	// Add additional filters
	for _, f := range additionalFilters {
		if f != "" {
			filter += " " + f
		}
	}

	return filter
}

// Iterator defines the generic iterator interface.
type Iterator[T any] interface {
	Next() bool
	Current() T
	Err() error
	HasMore() bool
	Total() int
	Close()
}
