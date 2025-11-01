//nolint:dupl

package pagination

import (
	"context"

	"github.com/brennoo/go-gmp/commands"
)

// TaskIterator is the iterator over tasks.
type TaskIterator = Iterator[*commands.GetTasksResponseTask]

// NewTaskIterator constructs a TaskIterator configured for GetTasks.
func NewTaskIterator(client Client, ctx context.Context, opts PaginationOptions, filters ...string) *TaskIterator {
	return newIterator[*commands.GetTasksResponseTask, *commands.GetTasks, *commands.GetTasksResponse](
		client,
		ctx,
		opts,
		filters,
		func(filter string) *commands.GetTasks { return &commands.GetTasks{Filter: filter} },
		func(client Client, ctx context.Context, cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
			return client.GetTasksRaw(ctx, cmd)
		},
		func(resp *commands.GetTasksResponse) []*commands.GetTasksResponseTask {
			items := make([]*commands.GetTasksResponseTask, len(resp.Task))
			for i := range resp.Task {
				items[i] = &resp.Task[i]
			}
			return items
		},
		func(resp *commands.GetTasksResponse) int { return resp.TaskCount.Filtered },
	)
}
