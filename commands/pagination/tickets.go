//nolint:dupl

package pagination

import (
	"context"

	"github.com/brennoo/go-gmp/commands"
)

// TicketIterator is the iterator over tickets.
type TicketIterator = Iterator[*commands.Ticket]

// NewTicketIterator constructs a TicketIterator configured for GetTickets.
func NewTicketIterator(client Client, ctx context.Context, opts PaginationOptions, filters ...string) *TicketIterator {
	return newIterator[*commands.Ticket, *commands.GetTickets, *commands.GetTicketsResponse](
		client,
		ctx,
		opts,
		filters,
		func(filter string) *commands.GetTickets { return &commands.GetTickets{Filter: filter} },
		func(client Client, ctx context.Context, cmd *commands.GetTickets) (*commands.GetTicketsResponse, error) {
			return client.GetTicketsRaw(ctx, cmd)
		},
		func(resp *commands.GetTicketsResponse) []*commands.Ticket {
			items := make([]*commands.Ticket, len(resp.Tickets))
			for i := range resp.Tickets {
				items[i] = &resp.Tickets[i]
			}
			return items
		},
		func(resp *commands.GetTicketsResponse) int {
			if resp.TicketCount != nil {
				return resp.TicketCount.Filtered
			}
			return 0
		},
	)
}
