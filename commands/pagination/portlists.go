//nolint:dupl

package pagination

import (
	"context"

	"github.com/brennoo/go-gmp/commands"
)

// PortListIterator is the iterator over port lists.
type PortListIterator = Iterator[*commands.PortList]

// NewPortListIterator constructs a PortListIterator configured for GetPortLists.
func NewPortListIterator(client Client, ctx context.Context, opts PaginationOptions, filters ...string) *PortListIterator {
	return newIterator[*commands.PortList, *commands.GetPortLists, *commands.GetPortListsResponse](
		client,
		ctx,
		opts,
		filters,
		func(filter string) *commands.GetPortLists { return &commands.GetPortLists{Filter: filter} },
		func(client Client, ctx context.Context, cmd *commands.GetPortLists) (*commands.GetPortListsResponse, error) {
			return client.GetPortListsRaw(ctx, cmd)
		},
		func(resp *commands.GetPortListsResponse) []*commands.PortList {
			items := make([]*commands.PortList, len(resp.PortLists))
			for i := range resp.PortLists {
				items[i] = &resp.PortLists[i]
			}
			return items
		},
		func(resp *commands.GetPortListsResponse) int {
			if resp.Count != nil {
				return resp.Count.Filtered
			}
			return 0
		},
	)
}
