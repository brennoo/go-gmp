//nolint:dupl

package pagination

import (
	"context"

	"github.com/brennoo/go-gmp/commands"
)

// ResultIterator is the iterator over results.
type ResultIterator = Iterator[*commands.Result]

// NewResultIterator constructs a ResultIterator configured for GetResults.
func NewResultIterator(client Client, ctx context.Context, opts PaginationOptions, filters ...string) *ResultIterator {
	return newIterator[*commands.Result, *commands.GetResults, *commands.GetResultsResponse](
		client,
		ctx,
		opts,
		filters,
		func(filter string) *commands.GetResults { return &commands.GetResults{Filter: filter} },
		func(client Client, ctx context.Context, cmd *commands.GetResults) (*commands.GetResultsResponse, error) {
			return client.GetResultsRaw(ctx, cmd)
		},
		func(resp *commands.GetResultsResponse) []*commands.Result {
			items := make([]*commands.Result, len(resp.Result))
			for i := range resp.Result {
				items[i] = &resp.Result[i]
			}
			return items
		},
		func(resp *commands.GetResultsResponse) int {
			if resp.ResultCount != nil {
				return resp.ResultCount.Filtered
			}
			return 0
		},
	)
}
