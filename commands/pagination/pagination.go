package pagination

import (
	"fmt"

	"github.com/brennoo/go-gmp/commands"
)

// Client is the minimal interface required by the pagination package.
type Client interface {
	GetTasksRaw(cmd *commands.GetTasks) (*commands.GetTasksResponse, error)
	GetResultsRaw(cmd *commands.GetResults) (*commands.GetResultsResponse, error)
	GetAssetsRaw(cmd *commands.GetAssets) (*commands.GetAssetsResponse, error)
	GetTargetsRaw(cmd *commands.GetTargets) (*commands.GetTargetsResponse, error)
	GetTicketsRaw(cmd *commands.GetTickets) (*commands.GetTicketsResponse, error)
	GetPortListsRaw(cmd *commands.GetPortLists) (*commands.GetPortListsResponse, error)
	GetSettingsRaw(cmd *commands.GetSettings) (*commands.GetSettingsResponse, error)
}

// PaginationOptions configures paging behavior.
type PaginationOptions struct {
	Page     int // 1-based page number
	PageSize int // Number of items per page (0 means rows=0)
	MaxItems int // Maximum total items to fetch (0 = no limit)
}

// DefaultPaginationOptions returns default paging values.
func DefaultPaginationOptions() PaginationOptions {
	return PaginationOptions{
		Page:     1,
		PageSize: 100,
		MaxItems: 0,
	}
}

// BuildPaginationFilter builds the "first/rows" filter and appends extra terms.
func BuildPaginationFilter(opts PaginationOptions, additionalFilters ...string) string {
	page := opts.Page
	if page <= 0 {
		page = 1
	}

	rows := opts.PageSize
	if rows < 0 {
		rows = 100
	}
	// rows == 0 is intentional and should remain 0 (empty page)

	first := (page-1)*max(rows, 1) + 1 // avoid zero in arithmetic; corrected by rows=0 check below

	if opts.MaxItems > 0 && first > opts.MaxItems {
		return "rows=0"
	}

	if rows > 0 && opts.MaxItems > 0 && first+rows-1 > opts.MaxItems {
		rows = opts.MaxItems - first + 1
	}

	filter := fmt.Sprintf("first=%d rows=%d", first, rows)

	for _, f := range additionalFilters {
		if f != "" {
			filter += " " + f
		}
	}

	return filter
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Pager is a minimal iterator interface used by tests.
type Pager[T any] interface {
	Next() bool
	Current() T
	Err() error
	HasMore() bool
	Total() int
	Close()
}
