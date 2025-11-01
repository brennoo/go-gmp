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
			items := make([]*commands.Result, len(resp.Results))
			for i := range resp.Results {
				items[i] = &resp.Results[i]
			}
			return items
		},
		func(resp *commands.GetResultsResponse) int {
			if resp.Count != nil {
				return resp.Count.Filtered
			}
			return 0
		},
	)
}
