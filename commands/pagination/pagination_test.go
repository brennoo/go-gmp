package pagination

import (
	"context"
	"errors"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

// fakeClient provides a test client for pagination tests.
type fakeClient struct {
	getTasksFunc func(*commands.GetTasks) (*commands.GetTasksResponse, error)
}

func (f *fakeClient) GetTasks(cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
	if f.getTasksFunc != nil {
		return f.getTasksFunc(cmd)
	}
	panic("GetTasks called unexpectedly")
}

func (f *fakeClient) GetResults(cmd *commands.GetResults) (*commands.GetResultsResponse, error) {
	panic("GetResults not implemented for this test")
}

func (f *fakeClient) GetAssets(cmd *commands.GetAssets) (*commands.GetAssetsResponse, error) {
	panic("GetAssets not implemented for this test")
}

func (f *fakeClient) GetTargets(cmd *commands.GetTargets) (*commands.GetTargetsResponse, error) {
	panic("GetTargets not implemented for this test")
}

func (f *fakeClient) GetTickets(cmd *commands.GetTickets) (*commands.GetTicketsResponse, error) {
	panic("GetTickets not implemented for this test")
}

func (f *fakeClient) GetPortLists(cmd *commands.GetPortLists) (*commands.GetPortListsResponse, error) {
	panic("GetPortLists not implemented for this test")
}

func (f *fakeClient) GetSettings(cmd *commands.GetSettings) (*commands.GetSettingsResponse, error) {
	panic("GetSettings not implemented for this test")
}

// Helper to create a TaskIterator with defaults.
func newTestTaskIterator(client Client, ctx context.Context) *TaskIterator {
	return &TaskIterator{
		Client:      client,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		HasMoreData: true,
		Page:        1,
	}
}

// Task iterator tests

func TestTaskIterator_PaginationSuccess(t *testing.T) {
	ctx := context.Background()
	pages := [][]commands.GetTasksResponseTask{
		{{ID: "1"}, {ID: "2"}},
		{{ID: "3"}},
	}

	callIndex := 0
	client := &fakeClient{
		getTasksFunc: func(cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
			if callIndex >= len(pages) {
				return &commands.GetTasksResponse{Task: []commands.GetTasksResponseTask{}, TaskCount: commands.GetTasksResponseTaskCount{Filtered: 3}}, nil
			}
			resp := &commands.GetTasksResponse{
				Task:      pages[callIndex],
				TaskCount: commands.GetTasksResponseTaskCount{Filtered: 3},
			}
			callIndex++
			return resp, nil
		},
	}

	it := newTestTaskIterator(client, ctx)

	ids := []string{}
	for it.Next() {
		task := it.Current()
		if task == nil {
			t.Fatal("Current() returned nil after Next()")
		}
		ids = append(ids, task.ID)
	}

	if len(ids) != 3 || ids[0] != "1" || ids[1] != "2" || ids[2] != "3" {
		t.Fatalf("Unexpected task IDs: %v", ids)
	}

	if it.Total() != 3 {
		t.Fatalf("Expected total=3, got %d", it.Total())
	}

	if it.Err() != nil {
		t.Fatalf("Unexpected error: %v", it.Err())
	}
}

func TestTaskIterator_EmptyFirstPageStops(t *testing.T) {
	client := &fakeClient{
		getTasksFunc: func(cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
			return &commands.GetTasksResponse{Task: []commands.GetTasksResponseTask{}, TaskCount: commands.GetTasksResponseTaskCount{Filtered: 0}}, nil
		},
	}

	it := newTestTaskIterator(client, context.Background())

	if it.Next() {
		t.Fatal("Expected no results")
	}
}

func TestTaskIterator_ErrorStops(t *testing.T) {
	client := &fakeClient{
		getTasksFunc: func(cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
			return nil, errors.New("load failed")
		},
	}

	it := newTestTaskIterator(client, context.Background())

	if it.Next() {
		t.Fatal("Expected Next to return false on error")
	}

	if it.Err() == nil || it.Err().Error() != "load failed" {
		t.Fatalf("Expected 'load failed' error, got %v", it.Err())
	}
}

func TestTaskIterator_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := &fakeClient{
		getTasksFunc: func(cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
			t.Fatal("Client should not be called when context is canceled")
			return nil, nil
		},
	}

	it := newTestTaskIterator(client, ctx)

	if it.Next() {
		t.Fatal("Expected Next to return false when context is canceled")
	}

	if !errors.Is(it.Err(), context.Canceled) {
		t.Fatalf("Expected context.Canceled, got %v", it.Err())
	}
}

func TestTaskIterator_CurrentBeforeNext(t *testing.T) {
	client := &fakeClient{
		getTasksFunc: func(cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
			return &commands.GetTasksResponse{Task: []commands.GetTasksResponseTask{{ID: "1"}}, TaskCount: commands.GetTasksResponseTaskCount{Filtered: 1}}, nil
		},
	}

	it := newTestTaskIterator(client, context.Background())

	if it.Current() != nil {
		t.Fatal("Expected Current() to be nil before Next()")
	}
}

func TestTaskIterator_CloseStopsIteration(t *testing.T) {
	client := &fakeClient{
		getTasksFunc: func(cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
			return &commands.GetTasksResponse{Task: []commands.GetTasksResponseTask{{ID: "1"}}, TaskCount: commands.GetTasksResponseTaskCount{Filtered: 1}}, nil
		},
	}

	it := newTestTaskIterator(client, context.Background())
	it.Close()

	if it.Next() {
		t.Fatal("Expected Next() to return false after Close()")
	}
}

func TestDefaultPaginationOptions(t *testing.T) {
	opts := DefaultPaginationOptions()
	if opts.Page != 1 {
		t.Errorf("Expected Page=1, got %d", opts.Page)
	}
	if opts.PageSize != 100 {
		t.Errorf("Expected PageSize=100, got %d", opts.PageSize)
	}
	if opts.MaxItems != 0 {
		t.Errorf("Expected MaxItems=0, got %d", opts.MaxItems)
	}
}
