package pagination

import (
	"context"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestTaskIterator(t *testing.T) {
	// Create test client
	cli := &testClient{}

	// Create iterator
	ctx := context.Background()
	iterator := &TaskIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	// Test iteration
	var tasks []*commands.GetTasksResponseTask
	for iterator.Next() {
		task := iterator.Current()
		tasks = append(tasks, task)
	}

	// Verify we got all tasks
	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}

	// Verify total count
	if iterator.Total() != 3 {
		t.Errorf("Expected total 3, got %d", iterator.Total())
	}

	// Test error handling
	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}

	// Test individual task properties
	if len(tasks) > 0 {
		if tasks[0].ID != "task-1" {
			t.Errorf("Expected first task ID 'task-1', got '%s'", tasks[0].ID)
		}
		if tasks[0].Name != "Test Task 1" {
			t.Errorf("Expected first task name 'Test Task 1', got '%s'", tasks[0].Name)
		}
	}
}
