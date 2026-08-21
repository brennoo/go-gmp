//nolint:dupl

package pagination

import (
	"context"

	"github.com/brennoo/go-gmp/commands"
)

// TaskIterator is the iterator over tasks.
type TaskIterator = Iterator[*commands.Task]

// NewTaskIterator constructs a TaskIterator configured for GetTasks.
func NewTaskIterator(client Client, ctx context.Context, opts PaginationOptions, filters ...string) *TaskIterator {
	return newIterator[*commands.Task, *commands.GetTasks, *commands.GetTasksResponse](
		client,
		ctx,
		opts,
		filters,
		func(filter string) *commands.GetTasks { return &commands.GetTasks{Filter: filter} },
		func(client Client, ctx context.Context, cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
			return client.GetTasksRaw(ctx, cmd)
		},
		func(resp *commands.GetTasksResponse) []*commands.Task {
			items := make([]*commands.Task, len(resp.Task))
			for i := range resp.Task {
				items[i] = &resp.Task[i]
			}
			return items
		},
		func(resp *commands.GetTasksResponse) int {
			if resp.TaskCount != nil {
				return resp.TaskCount.Filtered
			}
			return 0
		},
	)
}
