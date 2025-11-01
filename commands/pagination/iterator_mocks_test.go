package pagination

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

// Simple test client that implements the minimal Client interface.
type testClient struct {
	conn testConnection
}

type testConnection struct{}

//nolint:gocyclo // Mock connection handles all GMP commands
func (tc *testConnection) Execute(ctx context.Context, command interface{}, response interface{}) error {
	// Handle pagination commands
	if cmd, ok := command.(*commands.GetTasks); ok {
		resp := response.(*commands.GetTasksResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		if cmd.Filter != "" {
			// Check if page size is 0
			if cmd.Filter == "first=1 rows=0" {
				resp.Task = []commands.GetTasksResponseTask{}
			} else {
				// Parse pagination parameters
				first := 1
				rows := 3
				if strings.Contains(cmd.Filter, "first=") {
					// Extract just the pagination part before any additional filters
					parts := strings.Fields(cmd.Filter)
					for _, part := range parts {
						if strings.HasPrefix(part, "first=") {
							fmt.Sscanf(part, "first=%d", &first)
						} else if strings.HasPrefix(part, "rows=") {
							fmt.Sscanf(part, "rows=%d", &rows)
						}
					}
				}

				// Return appropriate slice based on pagination
				allTasks := []commands.GetTasksResponseTask{
					{ID: "task-1", Name: "Test Task 1"},
					{ID: "task-2", Name: "Test Task 2"},
					{ID: "task-3", Name: "Test Task 3"},
				}

				start := first - 1
				end := start + rows
				if start >= len(allTasks) {
					resp.Task = []commands.GetTasksResponseTask{}
				} else {
					if end > len(allTasks) {
						end = len(allTasks)
					}
					resp.Task = allTasks[start:end]
				}
			}
			resp.TaskCount = commands.GetTasksResponseTaskCount{Filtered: 3}
		}
	}

	if cmd, ok := command.(*commands.GetResults); ok {
		resp := response.(*commands.GetResultsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		if cmd.Filter != "" {
			// Check if page size is 0
			if cmd.Filter == "first=1 rows=0" {
				resp.Results = []commands.Result{}
			} else {
				// Parse pagination parameters
				first := 1
				rows := 3
				if strings.Contains(cmd.Filter, "first=") {
					// Extract just the pagination part before any additional filters
					parts := strings.Fields(cmd.Filter)
					for _, part := range parts {
						if strings.HasPrefix(part, "first=") {
							fmt.Sscanf(part, "first=%d", &first)
						} else if strings.HasPrefix(part, "rows=") {
							fmt.Sscanf(part, "rows=%d", &rows)
						}
					}
				}

				// Return appropriate slice based on pagination
				allResults := []commands.Result{
					{ID: "result-1", Name: "Test Result 1"},
					{ID: "result-2", Name: "Test Result 2"},
					{ID: "result-3", Name: "Test Result 3"},
				}

				start := first - 1
				end := start + rows
				if start >= len(allResults) {
					resp.Results = []commands.Result{}
				} else {
					if end > len(allResults) {
						end = len(allResults)
					}
					resp.Results = allResults[start:end]
				}
			}
			resp.Count = &commands.ResultsCount{Filtered: 3}
		}
	}

	if cmd, ok := command.(*commands.GetAssets); ok {
		resp := response.(*commands.GetAssetsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		if cmd.Filter != "" {
			// Check if page size is 0
			if cmd.Filter == "first=1 rows=0" {
				resp.Assets = []commands.Asset{}
			} else {
				// Parse pagination parameters
				first := 1
				rows := 2
				if strings.Contains(cmd.Filter, "first=") {
					fmt.Sscanf(cmd.Filter, "first=%d rows=%d", &first, &rows)
				}

				// Return appropriate slice based on pagination
				allAssets := []commands.Asset{
					{ID: "asset-1", Name: "Test Asset 1"},
					{ID: "asset-2", Name: "Test Asset 2"},
				}

				start := first - 1
				end := start + rows
				if start >= len(allAssets) {
					resp.Assets = []commands.Asset{}
				} else {
					if end > len(allAssets) {
						end = len(allAssets)
					}
					resp.Assets = allAssets[start:end]
				}
			}
			resp.AssetCount = &commands.AssetCount{Filtered: 2}
		}
	}

	if cmd, ok := command.(*commands.GetTargets); ok {
		resp := response.(*commands.GetTargetsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		if cmd.Filter != "" {
			// Check if page size is 0
			if cmd.Filter == "first=1 rows=0" {
				resp.Target = []commands.Target{}
			} else {
				// Parse pagination parameters
				first := 1
				rows := 3
				if strings.Contains(cmd.Filter, "first=") {
					// Extract just the pagination part before any additional filters
					parts := strings.Fields(cmd.Filter)
					for _, part := range parts {
						if strings.HasPrefix(part, "first=") {
							fmt.Sscanf(part, "first=%d", &first)
						} else if strings.HasPrefix(part, "rows=") {
							fmt.Sscanf(part, "rows=%d", &rows)
						}
					}
				}

				// Return appropriate slice based on pagination
				allTargets := []commands.Target{
					{ID: "target-1", Name: "Test Target 1"},
					{ID: "target-2", Name: "Test Target 2"},
					{ID: "target-3", Name: "Test Target 3"},
				}

				start := first - 1
				end := start + rows
				if start >= len(allTargets) {
					resp.Target = []commands.Target{}
				} else {
					if end > len(allTargets) {
						end = len(allTargets)
					}
					resp.Target = allTargets[start:end]
				}
			}
			resp.TargetCount = &commands.TargetCount{Filtered: 3}
		}
	}

	if cmd, ok := command.(*commands.GetTickets); ok {
		resp := response.(*commands.GetTicketsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		if cmd.Filter != "" {
			// Check if page size is 0
			if cmd.Filter == "first=1 rows=0" {
				resp.Tickets = []commands.Ticket{}
			} else {
				// Parse pagination parameters
				first := 1
				rows := 3
				if strings.Contains(cmd.Filter, "first=") {
					// Extract just the pagination part before any additional filters
					parts := strings.Fields(cmd.Filter)
					for _, part := range parts {
						if strings.HasPrefix(part, "first=") {
							fmt.Sscanf(part, "first=%d", &first)
						} else if strings.HasPrefix(part, "rows=") {
							fmt.Sscanf(part, "rows=%d", &rows)
						}
					}
				}

				// Return appropriate slice based on pagination
				allTickets := []commands.Ticket{
					{ID: "ticket-1", Name: "Test Ticket 1"},
					{ID: "ticket-2", Name: "Test Ticket 2"},
					{ID: "ticket-3", Name: "Test Ticket 3"},
				}

				start := first - 1
				end := start + rows
				if start >= len(allTickets) {
					resp.Tickets = []commands.Ticket{}
				} else {
					if end > len(allTickets) {
						end = len(allTickets)
					}
					resp.Tickets = allTickets[start:end]
				}
			}
			resp.TicketCount = &commands.TicketCount{Filtered: 3}
		}
	}

	if cmd, ok := command.(*commands.GetPortLists); ok {
		resp := response.(*commands.GetPortListsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		if cmd.Filter != "" {
			// Check if page size is 0
			if cmd.Filter == "first=1 rows=0" {
				resp.PortLists = []commands.PortList{}
			} else {
				// Parse pagination parameters
				first := 1
				rows := 3
				if strings.Contains(cmd.Filter, "first=") {
					// Extract just the pagination part before any additional filters
					parts := strings.Fields(cmd.Filter)
					for _, part := range parts {
						if strings.HasPrefix(part, "first=") {
							fmt.Sscanf(part, "first=%d", &first)
						} else if strings.HasPrefix(part, "rows=") {
							fmt.Sscanf(part, "rows=%d", &rows)
						}
					}
				}

				// Return appropriate slice based on pagination
				allPortLists := []commands.PortList{
					{ID: "portlist-1", Name: "Test PortList 1"},
					{ID: "portlist-2", Name: "Test PortList 2"},
					{ID: "portlist-3", Name: "Test PortList 3"},
				}

				start := first - 1
				end := start + rows
				if start >= len(allPortLists) {
					resp.PortLists = []commands.PortList{}
				} else {
					if end > len(allPortLists) {
						end = len(allPortLists)
					}
					resp.PortLists = allPortLists[start:end]
				}
			}
			resp.Count = &commands.PortListCount{Filtered: 3}
		}
	}

	if cmd, ok := command.(*commands.GetSettings); ok {
		resp := response.(*commands.GetSettingsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		if cmd.Filter != "" {
			// Check if page size is 0
			if cmd.Filter == "first=1 rows=0" {
				resp.SettingsList = []commands.Setting{}
			} else {
				// Parse pagination parameters
				first := 1
				rows := 3
				if strings.Contains(cmd.Filter, "first=") {
					// Extract just the pagination part before any additional filters
					parts := strings.Fields(cmd.Filter)
					for _, part := range parts {
						if strings.HasPrefix(part, "first=") {
							fmt.Sscanf(part, "first=%d", &first)
						} else if strings.HasPrefix(part, "rows=") {
							fmt.Sscanf(part, "rows=%d", &rows)
						}
					}
				}

				// Return appropriate slice based on pagination
				allSettings := []commands.Setting{
					{ID: "setting-1", Name: "Test Setting 1"},
					{ID: "setting-2", Name: "Test Setting 2"},
					{ID: "setting-3", Name: "Test Setting 3"},
				}

				start := first - 1
				end := start + rows
				if start >= len(allSettings) {
					resp.SettingsList = []commands.Setting{}
				} else {
					if end > len(allSettings) {
						end = len(allSettings)
					}
					resp.SettingsList = allSettings[start:end]
				}
			}
			resp.SettingCount = &commands.SettingCount{Filtered: 3}
		}
	}

	return nil
}

func (tc *testClient) GetTasksRaw(ctx context.Context, cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
	resp := &commands.GetTasksResponse{}
	err := tc.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (tc *testClient) GetResultsRaw(ctx context.Context, cmd *commands.GetResults) (*commands.GetResultsResponse, error) {
	resp := &commands.GetResultsResponse{}
	err := tc.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (tc *testClient) GetAssetsRaw(ctx context.Context, cmd *commands.GetAssets) (*commands.GetAssetsResponse, error) {
	resp := &commands.GetAssetsResponse{}
	err := tc.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (tc *testClient) GetTargetsRaw(ctx context.Context, cmd *commands.GetTargets) (*commands.GetTargetsResponse, error) {
	resp := &commands.GetTargetsResponse{}
	err := tc.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (tc *testClient) GetTicketsRaw(ctx context.Context, cmd *commands.GetTickets) (*commands.GetTicketsResponse, error) {
	resp := &commands.GetTicketsResponse{}
	err := tc.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (tc *testClient) GetPortListsRaw(ctx context.Context, cmd *commands.GetPortLists) (*commands.GetPortListsResponse, error) {
	resp := &commands.GetPortListsResponse{}
	err := tc.conn.Execute(ctx, cmd, resp)
	return resp, err
}

func (tc *testClient) GetSettingsRaw(ctx context.Context, cmd *commands.GetSettings) (*commands.GetSettingsResponse, error) {
	resp := &commands.GetSettingsResponse{}
	err := tc.conn.Execute(ctx, cmd, resp)
	return resp, err
}

// Test error client that returns errors.
type errorTestClient struct{}

func (etc *errorTestClient) GetTasksRaw(ctx context.Context, cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
	return nil, errors.New("test error")
}

func (etc *errorTestClient) GetResultsRaw(ctx context.Context, cmd *commands.GetResults) (*commands.GetResultsResponse, error) {
	return nil, errors.New("test error")
}

func (etc *errorTestClient) GetAssetsRaw(ctx context.Context, cmd *commands.GetAssets) (*commands.GetAssetsResponse, error) {
	return nil, errors.New("test error")
}

func (etc *errorTestClient) GetTargetsRaw(ctx context.Context, cmd *commands.GetTargets) (*commands.GetTargetsResponse, error) {
	return nil, errors.New("test error")
}

func (etc *errorTestClient) GetTicketsRaw(ctx context.Context, cmd *commands.GetTickets) (*commands.GetTicketsResponse, error) {
	return nil, errors.New("test error")
}

func (etc *errorTestClient) GetPortListsRaw(ctx context.Context, cmd *commands.GetPortLists) (*commands.GetPortListsResponse, error) {
	return nil, errors.New("test error")
}

func (etc *errorTestClient) GetSettingsRaw(ctx context.Context, cmd *commands.GetSettings) (*commands.GetSettingsResponse, error) {
	return nil, errors.New("test error")
}

// General iterator tests that apply to all iterator types.
func TestIteratorErrorHandling(t *testing.T) {
	// Test with error client
	errorCli := &errorTestClient{}
	ctx := context.Background()

	// Test TaskIterator error handling
	taskIter := NewTaskIterator(errorCli, ctx, PaginationOptions{Page: 1, PageSize: 2})

	if taskIter.Next() {
		t.Error("Expected Next() to return false when client returns error")
	}

	if taskIter.Err() == nil {
		t.Error("Expected Err() to return the error from client")
	}
}

func TestIteratorEdgeCases(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	// Test with large page size
	iterator := NewTaskIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 100})

	var tasks []*commands.GetTasksResponseTask
	for iterator.Next() {
		task := iterator.Current()
		tasks = append(tasks, task)
	}

	// Should get all 3 tasks
	if len(tasks) != 3 {
		t.Errorf("Expected 3 items with large page size, got %d", len(tasks))
	}
}

func TestIteratorWithFilters(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	// Test with filters
	iterator := NewTaskIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 2}, "name=Test")

	var tasks []*commands.GetTasksResponseTask
	for iterator.Next() {
		task := iterator.Current()
		tasks = append(tasks, task)
	}

	// Should get all 3 tasks with filters
	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks with filters, got %d", len(tasks))
	}
}

func TestIteratorPagination(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	// Test with small page size to trigger pagination
	iterator := NewTaskIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 1})

	var tasks []*commands.GetTasksResponseTask
	for iterator.Next() {
		task := iterator.Current()
		tasks = append(tasks, task)
	}

	// Should get all 3 tasks with small page size
	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks with small page size, got %d", len(tasks))
	}
}

func TestIteratorStateManagement(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	iterator := NewTaskIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 2})

	// Test first item
	if !iterator.Next() {
		t.Fatal("Expected first Next() to return true")
	}

	task := iterator.Current()
	if task == nil {
		t.Fatal("Expected first item, got nil")
	}

	// Should still have more data
	if !iterator.HasMore() {
		t.Error("Iterator should still have more data after first item")
	}

	// Test second item
	if !iterator.Next() {
		t.Fatal("Expected second Next() to return true")
	}

	task = iterator.Current()
	if task == nil {
		t.Fatal("Expected second item, got nil")
	}
}

func TestIteratorErrorScenarios(t *testing.T) {
	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	iterator := NewTaskIterator(&testClient{}, ctx, PaginationOptions{Page: 1, PageSize: 2})

	if iterator.Next() {
		t.Error("Expected Next() to return false when context is canceled")
	}

	if iterator.Err() == nil {
		t.Error("Expected Err() to return context canceled error")
	}
}

func TestIteratorEmptyCurrentSlice(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	// Test with page size 0 to get empty results
	iterator := NewTaskIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 0})

	if iterator.Next() {
		t.Error("Expected Next() to return false for empty results")
	}

	if iterator.Current() != nil {
		t.Error("Expected Current() to return nil for empty results")
	}
}

func TestIteratorContextCancellation(t *testing.T) {
	cli := &testClient{}
	ctx, cancel := context.WithCancel(context.Background())

	iterator := NewTaskIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 2})

	// Cancel context after first item
	if !iterator.Next() {
		t.Fatal("Expected first Next() to return true")
	}

	cancel() // Cancel context

	// Next call should detect cancellation
	if iterator.Next() {
		t.Error("Expected Next() to return false after context cancellation")
	}

	if iterator.Err() == nil {
		t.Error("Expected Err() to return context canceled error")
	}
}

func TestIteratorLoadPageError(t *testing.T) {
	errorCli := &errorTestClient{}
	ctx := context.Background()

	iterator := NewTaskIterator(errorCli, ctx, PaginationOptions{Page: 1, PageSize: 2})

	if iterator.Next() {
		t.Error("Expected Next() to return false when loadPage returns error")
	}

	if iterator.Err() == nil {
		t.Error("Expected Err() to return the error from loadPage")
	}
}

func TestIteratorCurrentBoundaryConditions(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	iterator := NewTaskIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 2})

	// Current() should return nil before any Next() calls
	if iterator.Current() != nil {
		t.Error("Expected Current() to return nil before Next() is called")
	}

	// After Next() returns false, Current() should return the last item
	if !iterator.Next() {
		t.Fatal("Expected first Next() to return true")
	}

	task := iterator.Current()
	if task == nil {
		t.Fatal("Expected Current() to return a task after Next()")
	}

	// Continue until no more items
	for iterator.Next() {
		// Just advance
	}

	// Current() should still return the last item
	lastTask := iterator.Current()
	if lastTask == nil {
		t.Error("Expected Current() to return last item after iteration ends")
	}
}

func TestIteratorNextNoMoreData(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	iterator := NewTaskIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 2})

	// Consume all items
	for iterator.Next() {
		// Just advance
	}

	// Next() should return false when no more data
	if iterator.Next() {
		t.Error("Expected Next() to return false when no more data")
	}

	// HasMore() should return false when no more data
	if iterator.HasMore() {
		t.Error("Expected HasMore() to return false when no more data")
	}
}

func TestIteratorNextValidContext(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	iterator := NewTaskIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 2})

	// Next() should work with valid context
	if !iterator.Next() {
		t.Error("Expected Next() to return true with valid context")
	}

	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}
}

func TestIteratorNextValidContextAll(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	iterator := NewTaskIterator(cli, ctx, PaginationOptions{Page: 1, PageSize: 2})

	// Iterate through all items
	count := 0
	for iterator.Next() {
		count++
		task := iterator.Current()
		if task == nil {
			t.Errorf("Expected Current() to return task %d", count)
		}
	}

	// Should have processed all items
	if count != 3 {
		t.Errorf("Expected to process 3 items, got %d", count)
	}

	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}
}
