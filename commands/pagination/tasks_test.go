package pagination

import (
	"context"
	"testing"
)

func TestTaskIterator(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := NewTaskIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 3})

	var tasks []string
	for iterator.Next() {
		task := iterator.Current()
		tasks = append(tasks, task.ID)
	}

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}

	if iterator.Total() != 3 {
		t.Errorf("Expected total 3, got %d", iterator.Total())
	}
}
