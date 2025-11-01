package pagination

import (
	"context"
	"fmt"
)

// newIterator is a generic factory function that constructs a configured Iterator.
// It takes generic type parameters for the Item (T), Command (CmdT), and Response (RespT).
// It accepts functions to handle the type-specific logic: building a command,
// executing it, and extracting data from the response.
func newIterator[T any, CmdT interface{}, RespT interface{}](
	client Client,
	ctx context.Context,
	opts PaginationOptions,
	filters []string,
	buildCommand func(filter string) CmdT,
	executeCommand func(client Client, ctx context.Context, cmd CmdT) (RespT, error),
	extractItems func(resp RespT) []T,
	extractTotal func(resp RespT) int,
) *Iterator[T] {
	// Ensure sensible defaults but preserve explicit empty page (PageSize==0)
	if opts.PageSize < 0 {
		opts.PageSize = DefaultPaginationOptions().PageSize
	}

	return &Iterator[T]{
		Client:      client,
		Ctx:         ctx,
		Opts:        opts,
		Filters:     filters,
		current:     make([]T, 0),
		index:       0,
		Page:        0, // Start at page 0 so first loadPage targets page 1
		HasMoreData: true,
		total:       0,
		err:         nil,

		// buildCommand now just returns the result of the provided builder
		buildCommand: func(filter string) interface{} {
			return buildCommand(filter)
		},

		// executeAndParse is now a generic closure that uses the provided functions
		executeAndParse: func(cmd interface{}, pageOpts PaginationOptions) (items []T, totalFiltered int, hasMore bool, err error) {
			// Type assertion to the generic command type
			typedCmd, ok := cmd.(CmdT)
			if !ok {
				return nil, 0, false, fmt.Errorf("internal error: command type mismatch")
			}

			// Call the specific execute func
			resp, err := executeCommand(client, ctx, typedCmd)
			if err != nil {
				return nil, 0, false, err
			}

			// Call the specific extract funcs
			items = extractItems(resp)
			totalFiltered = extractTotal(resp)

			// Common logic for determining if more pages exist
			hasMore = len(items) == pageOpts.PageSize && pageOpts.PageSize > 0

			return items, totalFiltered, hasMore, nil
		},
	}
}
