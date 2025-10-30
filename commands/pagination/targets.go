//nolint:dupl

package pagination

import (
	"context"

	"github.com/brennoo/go-gmp/commands"
)

// TargetIterator is the iterator over targets.
type TargetIterator = Iterator[*commands.Target]

// NewTargetIterator constructs a TargetIterator configured for GetTargets.
func NewTargetIterator(client Client, ctx context.Context, opts PaginationOptions, filters ...string) *TargetIterator {
	return newIterator[*commands.Target, *commands.GetTargets, *commands.GetTargetsResponse](
		client,
		ctx,
		opts,
		filters,
		func(filter string) *commands.GetTargets { return &commands.GetTargets{Filter: filter} },
		func(client Client, cmd *commands.GetTargets) (*commands.GetTargetsResponse, error) {
			return client.GetTargetsRaw(cmd)
		},
		func(resp *commands.GetTargetsResponse) []*commands.Target {
			items := make([]*commands.Target, len(resp.Target))
			for i := range resp.Target {
				items[i] = &resp.Target[i]
			}
			return items
		},
		func(resp *commands.GetTargetsResponse) int {
			if resp.TargetCount != nil {
				return resp.TargetCount.Filtered
			}
			return 0
		},
	)
}
