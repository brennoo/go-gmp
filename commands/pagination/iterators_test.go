package pagination

import (
	"context"
	"errors"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

// Simple test client that implements the minimal Client interface
type testClient struct {
	conn testConnection
}

type testConnection struct{}

func (tc *testConnection) Execute(command interface{}, response interface{}) error {
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
				// Always return 3 tasks for simplicity
				resp.Task = []commands.GetTasksResponseTask{
					{ID: "task-1", Name: "Test Task 1"},
					{ID: "task-2", Name: "Test Task 2"},
					{ID: "task-3", Name: "Test Task 3"},
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
			resp.Results = []commands.Result{
				{ID: "result-1", Name: "Test Result 1"},
				{ID: "result-2", Name: "Test Result 2"},
				{ID: "result-3", Name: "Test Result 3"},
			}
			resp.Count = &commands.ResultsCount{Filtered: 3}
		}
	}

	if cmd, ok := command.(*commands.GetAssets); ok {
		resp := response.(*commands.GetAssetsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		if cmd.Filter != "" {
			resp.Assets = []commands.Asset{
				{ID: "asset-1", Name: "Test Asset 1"},
				{ID: "asset-2", Name: "Test Asset 2"},
			}
			resp.AssetCount = &commands.AssetCount{Filtered: 2}
		}
	}

	if cmd, ok := command.(*commands.GetTargets); ok {
		resp := response.(*commands.GetTargetsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		if cmd.Filter != "" {
			resp.Target = []commands.Target{
				{ID: "target-1", Name: "Test Target 1"},
				{ID: "target-2", Name: "Test Target 2"},
				{ID: "target-3", Name: "Test Target 3"},
			}
			resp.TargetCount = &commands.TargetCount{Filtered: 3}
		}
	}

	if cmd, ok := command.(*commands.GetTickets); ok {
		resp := response.(*commands.GetTicketsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		if cmd.Filter != "" {
			resp.Tickets = []commands.Ticket{
				{ID: "ticket-1", Name: "Test Ticket 1"},
				{ID: "ticket-2", Name: "Test Ticket 2"},
				{ID: "ticket-3", Name: "Test Ticket 3"},
			}
			resp.TicketCount = &commands.TicketCount{Filtered: 3}
		}
	}

	if cmd, ok := command.(*commands.GetPortLists); ok {
		resp := response.(*commands.GetPortListsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		if cmd.Filter != "" {
			resp.PortLists = []commands.PortList{
				{ID: "portlist-1", Name: "Test PortList 1"},
				{ID: "portlist-2", Name: "Test PortList 2"},
				{ID: "portlist-3", Name: "Test PortList 3"},
			}
			resp.Count = &commands.PortListCount{Filtered: 3}
		}
	}

	if cmd, ok := command.(*commands.GetSettings); ok {
		resp := response.(*commands.GetSettingsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		if cmd.Filter != "" {
			resp.SettingsList = []commands.Setting{
				{ID: "setting-1", Name: "Test Setting 1"},
				{ID: "setting-2", Name: "Test Setting 2"},
				{ID: "setting-3", Name: "Test Setting 3"},
			}
			resp.SettingCount = &commands.SettingCount{Filtered: 3}
		}
	}

	return nil
}

func (tc *testConnection) Close() error {
	return nil
}

// Implement the Client interface methods
func (tc *testClient) GetTasks(cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
	resp := &commands.GetTasksResponse{}
	err := tc.conn.Execute(cmd, resp)
	return resp, err
}

func (tc *testClient) GetResults(cmd *commands.GetResults) (*commands.GetResultsResponse, error) {
	resp := &commands.GetResultsResponse{}
	err := tc.conn.Execute(cmd, resp)
	return resp, err
}

func (tc *testClient) GetAssets(cmd *commands.GetAssets) (*commands.GetAssetsResponse, error) {
	resp := &commands.GetAssetsResponse{}
	err := tc.conn.Execute(cmd, resp)
	return resp, err
}

func (tc *testClient) GetTargets(cmd *commands.GetTargets) (*commands.GetTargetsResponse, error) {
	resp := &commands.GetTargetsResponse{}
	err := tc.conn.Execute(cmd, resp)
	return resp, err
}

func (tc *testClient) GetTickets(cmd *commands.GetTickets) (*commands.GetTicketsResponse, error) {
	resp := &commands.GetTicketsResponse{}
	err := tc.conn.Execute(cmd, resp)
	return resp, err
}

func (tc *testClient) GetPortLists(cmd *commands.GetPortLists) (*commands.GetPortListsResponse, error) {
	resp := &commands.GetPortListsResponse{}
	err := tc.conn.Execute(cmd, resp)
	return resp, err
}

func (tc *testClient) GetSettings(cmd *commands.GetSettings) (*commands.GetSettingsResponse, error) {
	resp := &commands.GetSettingsResponse{}
	err := tc.conn.Execute(cmd, resp)
	return resp, err
}

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
		Page:        1,
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

	// Test close
	iterator.Close()
}

func TestResultIterator(t *testing.T) {
	// Create test client
	cli := &testClient{}

	// Create iterator
	ctx := context.Background()
	iterator := &ResultIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test iteration
	var results []*commands.Result
	for iterator.HasMore() {
		if !iterator.Next() {
			break
		}
		result := iterator.Current()
		results = append(results, result)
	}

	// Verify we got all results
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// Verify total count
	if iterator.Total() != 3 {
		t.Errorf("Expected total 3, got %d", iterator.Total())
	}

	// Test error handling
	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}

	// Test close
	iterator.Close()
}

func TestAssetIterator(t *testing.T) {
	// Create test client
	cli := &testClient{}

	// Create iterator
	ctx := context.Background()
	iterator := &AssetIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 1},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test iteration
	var assets []*commands.Asset
	for iterator.HasMore() {
		if !iterator.Next() {
			break
		}
		asset := iterator.Current()
		assets = append(assets, asset)
	}

	// Verify we got all assets
	if len(assets) != 2 {
		t.Errorf("Expected 2 assets, got %d", len(assets))
	}

	// Verify total count
	if iterator.Total() != 2 {
		t.Errorf("Expected total 2, got %d", iterator.Total())
	}

	// Test error handling
	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}

	// Test close
	iterator.Close()
}

func TestTargetIterator(t *testing.T) {
	// Create test client
	cli := &testClient{}

	// Create iterator
	ctx := context.Background()
	iterator := &TargetIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test iteration
	var targets []*commands.Target
	for iterator.HasMore() {
		if !iterator.Next() {
			break
		}
		target := iterator.Current()
		targets = append(targets, target)
	}

	// Verify we got all targets
	if len(targets) != 3 {
		t.Errorf("Expected 3 targets, got %d", len(targets))
	}

	// Test Total() method
	if iterator.Total() != 3 {
		t.Errorf("Expected total 3, got %d", iterator.Total())
	}

	// Test error handling
	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}

	// Test close
	iterator.Close()
}

func TestTicketIterator(t *testing.T) {
	// Create test client
	cli := &testClient{}

	// Create iterator
	ctx := context.Background()
	iterator := &TicketIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test iteration
	var tickets []*commands.Ticket
	for iterator.HasMore() {
		if !iterator.Next() {
			break
		}
		ticket := iterator.Current()
		tickets = append(tickets, ticket)
	}

	// Verify we got all tickets
	if len(tickets) != 3 {
		t.Errorf("Expected 3 tickets, got %d", len(tickets))
	}

	// Verify total count
	if iterator.Total() != 3 {
		t.Errorf("Expected total 3, got %d", iterator.Total())
	}

	// Test error handling
	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}

	// Test close
	iterator.Close()
}

func TestPortListIterator(t *testing.T) {
	// Create test client
	cli := &testClient{}

	// Create iterator
	ctx := context.Background()
	iterator := &PortListIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test iteration
	var portLists []*commands.PortList
	for iterator.HasMore() {
		if !iterator.Next() {
			break
		}
		portList := iterator.Current()
		portLists = append(portLists, portList)
	}

	// Verify we got all port lists
	if len(portLists) != 3 {
		t.Errorf("Expected 3 port lists, got %d", len(portLists))
	}

	// Verify total count
	if iterator.Total() != 3 {
		t.Errorf("Expected total 3, got %d", iterator.Total())
	}

	// Test error handling
	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}

	// Test close
	iterator.Close()
}

func TestSettingsIterator(t *testing.T) {
	// Create test client
	cli := &testClient{}

	// Create iterator
	ctx := context.Background()
	iterator := &SettingsIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test iteration
	var settings []*commands.Setting
	for iterator.HasMore() {
		if !iterator.Next() {
			break
		}
		setting := iterator.Current()
		settings = append(settings, setting)
	}

	// Verify we got all settings
	if len(settings) != 3 {
		t.Errorf("Expected 3 settings, got %d", len(settings))
	}

	// Verify total count
	if iterator.Total() != 3 {
		t.Errorf("Expected total 3, got %d", iterator.Total())
	}

	// Test error handling
	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}

	// Test close
	iterator.Close()
}

// Test error scenarios
func TestIteratorErrorHandling(t *testing.T) {
	// Test with invalid context (cancelled)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	cli := &testClient{}
	iterator := &TaskIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Try to get next item - should handle cancelled context gracefully
	if iterator.Next() {
		t.Error("Expected Next() to return false due to cancelled context")
	}
	if iterator.Err() != nil && iterator.Err() != context.Canceled {
		t.Errorf("Expected context.Canceled error, got: %v", iterator.Err())
	}

	iterator.Close()
}

// Test edge cases
func TestIteratorEdgeCases(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	// Test with page size of 0 - should still work but get no data
	iterator := &TaskIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 0},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Try to get items with page size 0
	var count int
	for iterator.HasMore() {
		if !iterator.Next() {
			break
		}
		_ = iterator.Current() // Use the variable to avoid unused variable error
		count++
	}

	// Should get no items with page size 0
	if count != 0 {
		t.Errorf("Expected 0 items with page size 0, got %d", count)
	}
	iterator.Close()

	// Test with very large page size
	iterator = &TaskIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 1000},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	var count2 int
	for iterator.HasMore() {
		if !iterator.Next() {
			break
		}
		_ = iterator.Current() // Use the variable to avoid unused variable error
		count2++
	}

	// Should get all available items
	if count2 != 3 {
		t.Errorf("Expected 3 items with large page size, got %d", count2)
	}
	iterator.Close()
}

// Test iterator with filters
func TestIteratorWithFilters(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	// Test with filters
	iterator := &TaskIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{"status=running", "sort=name"},
		Page:        1,
		HasMoreData: true,
	}

	var tasks []*commands.GetTasksResponseTask
	for iterator.HasMore() {
		if !iterator.Next() {
			break
		}
		task := iterator.Current()
		tasks = append(tasks, task)
	}

	// Should still get all tasks (mock doesn't filter)
	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks with filters, got %d", len(tasks))
	}

	iterator.Close()
}

// Test pagination scenarios
func TestIteratorPagination(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	// Test with small page size to force pagination
	iterator := &TaskIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 1},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	var tasks []*commands.GetTasksResponseTask
	for iterator.HasMore() {
		if !iterator.Next() {
			break
		}
		task := iterator.Current()
		tasks = append(tasks, task)
	}

	// Should get all tasks even with small page size
	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks with small page size, got %d", len(tasks))
	}

	iterator.Close()
}

// Test iterator state management
func TestIteratorStateManagement(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	iterator := &TaskIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test initial state
	if !iterator.HasMore() {
		t.Error("Iterator should have more data initially")
	}

	// Get first item
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for first item")
	}
	task := iterator.Current()
	if task == nil {
		t.Error("Expected first item, got nil")
	}

	// Should still have more data
	if !iterator.HasMore() {
		t.Error("Iterator should still have more data after first item")
	}

	// Get remaining items
	for iterator.HasMore() {
		if !iterator.Next() {
			break
		}
		_ = iterator.Current() // Use the variable to avoid unused variable error
	}

	// Should not have more data now
	if iterator.HasMore() {
		t.Error("Iterator should not have more data after consuming all items")
	}

	iterator.Close()
}

// Test generic Next method with type assertion
func TestNextTaskMethod(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	iterator := &TaskIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test Next method with typed Current()
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true")
	}

	task := iterator.Current()
	if task == nil {
		t.Error("Expected task, got nil")
	}
	if task.ID != "task-1" {
		t.Errorf("Expected task-1, got %s", task.ID)
	}

	iterator.Close()
}

func TestNextResultMethod(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	iterator := &ResultIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test Next method with typed Current()
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true")
	}

	result := iterator.Current()
	if result == nil {
		t.Error("Expected result, got nil")
	}
	if result.ID != "result-1" {
		t.Errorf("Expected result-1, got %s", result.ID)
	}

	iterator.Close()
}

func TestNextAssetMethod(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	iterator := &AssetIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test Next method with typed Current()
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true")
	}

	asset := iterator.Current()
	if asset == nil {
		t.Error("Expected asset, got nil")
	}
	if asset.ID != "asset-1" {
		t.Errorf("Expected asset-1, got %s", asset.ID)
	}

	iterator.Close()
}

func TestNextTargetMethod(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	iterator := &TargetIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test Next method with typed Current()
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true")
	}

	target := iterator.Current()
	if target == nil {
		t.Error("Expected target, got nil")
	}
	if target.ID != "target-1" {
		t.Errorf("Expected target-1, got %s", target.ID)
	}

	iterator.Close()
}

func TestNextTicketMethod(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	iterator := &TicketIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test Next method with typed Current()
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true")
	}

	ticket := iterator.Current()
	if ticket == nil {
		t.Error("Expected ticket, got nil")
	}
	if ticket.ID != "ticket-1" {
		t.Errorf("Expected ticket-1, got %s", ticket.ID)
	}

	iterator.Close()
}

func TestNextPortListMethod(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	iterator := &PortListIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test Next method with typed Current()
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true")
	}

	portList := iterator.Current()
	if portList == nil {
		t.Error("Expected portList, got nil")
	}
	if portList.ID != "portlist-1" {
		t.Errorf("Expected portlist-1, got %s", portList.ID)
	}

	iterator.Close()
}

func TestNextSettingMethod(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	iterator := &SettingsIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test Next method with typed Current()
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true")
	}

	setting := iterator.Current()
	if setting == nil {
		t.Error("Expected setting, got nil")
	}
	if setting.ID != "setting-1" {
		t.Errorf("Expected setting-1, got %s", setting.ID)
	}

	iterator.Close()
}

// Test error scenarios for better coverage
func TestIteratorErrorScenarios(t *testing.T) {
	// Test with error client
	errorClient := &errorTestClient{}
	ctx := context.Background()

	iterator := &TaskIterator{
		Client:      errorClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Try to get next item - should get error
	if iterator.Next() {
		t.Error("Expected Next() to return false on error")
	}
	if iterator.Err() == nil {
		t.Error("Expected error, got nil")
	}
	if iterator.Current() != nil {
		t.Error("Expected nil current item on error, got item")
	}

	// Error should be stored
	if iterator.Err() == nil {
		t.Error("Expected error to be stored")
	}

	iterator.Close()
}

// Error test client
type errorTestClient struct{}

func (etc *errorTestClient) GetTasks(cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
	return nil, errors.New("mock error")
}

func (etc *errorTestClient) GetResults(cmd *commands.GetResults) (*commands.GetResultsResponse, error) {
	return nil, errors.New("mock error")
}

func (etc *errorTestClient) GetAssets(cmd *commands.GetAssets) (*commands.GetAssetsResponse, error) {
	return nil, errors.New("mock error")
}

func (etc *errorTestClient) GetTargets(cmd *commands.GetTargets) (*commands.GetTargetsResponse, error) {
	return nil, errors.New("mock error")
}

func (etc *errorTestClient) GetTickets(cmd *commands.GetTickets) (*commands.GetTicketsResponse, error) {
	return nil, errors.New("mock error")
}

func (etc *errorTestClient) GetPortLists(cmd *commands.GetPortLists) (*commands.GetPortListsResponse, error) {
	return nil, errors.New("mock error")
}

func (etc *errorTestClient) GetSettings(cmd *commands.GetSettings) (*commands.GetSettingsResponse, error) {
	return nil, errors.New("mock error")
}

// Test error scenarios for NextXXX methods
func TestNextTaskMethodError(t *testing.T) {
	errorClient := &errorTestClient{}
	ctx := context.Background()

	iterator := &TaskIterator{
		Client:      errorClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test Next method with error
	if iterator.Next() {
		t.Error("Expected Next() to return false on error")
	}
	if iterator.Err() == nil {
		t.Error("Expected error, got nil")
	}
	if iterator.Current() != nil {
		t.Error("Expected nil current item on error, got item")
	}

	iterator.Close()
}

func TestNextResultMethodError(t *testing.T) {
	errorClient := &errorTestClient{}
	ctx := context.Background()

	iterator := &ResultIterator{
		Client:      errorClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test Next method with error
	if iterator.Next() {
		t.Error("Expected Next() to return false on error")
	}
	if iterator.Err() == nil {
		t.Error("Expected error, got nil")
	}
	if iterator.Current() != nil {
		t.Error("Expected nil current item on error, got item")
	}

	iterator.Close()
}

func TestNextAssetMethodError(t *testing.T) {
	errorClient := &errorTestClient{}
	ctx := context.Background()

	iterator := &AssetIterator{
		Client:      errorClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test Next method with error
	if iterator.Next() {
		t.Error("Expected Next() to return false on error")
	}
	if iterator.Err() == nil {
		t.Error("Expected error, got nil")
	}
	if iterator.Current() != nil {
		t.Error("Expected nil current item on error, got item")
	}

	iterator.Close()
}

func TestNextTargetMethodError(t *testing.T) {
	errorClient := &errorTestClient{}
	ctx := context.Background()

	iterator := &TargetIterator{
		Client:      errorClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test Next method with error
	if iterator.Next() {
		t.Error("Expected Next() to return false on error")
	}
	if iterator.Err() == nil {
		t.Error("Expected error, got nil")
	}
	if iterator.Current() != nil {
		t.Error("Expected nil current item on error, got item")
	}

	iterator.Close()
}

func TestNextTicketMethodError(t *testing.T) {
	errorClient := &errorTestClient{}
	ctx := context.Background()

	iterator := &TicketIterator{
		Client:      errorClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test Next method with error
	if iterator.Next() {
		t.Error("Expected Next() to return false on error")
	}
	if iterator.Err() == nil {
		t.Error("Expected error, got nil")
	}
	if iterator.Current() != nil {
		t.Error("Expected nil current item on error, got item")
	}

	iterator.Close()
}

func TestNextPortListMethodError(t *testing.T) {
	errorClient := &errorTestClient{}
	ctx := context.Background()

	iterator := &PortListIterator{
		Client:      errorClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test Next method with error
	if iterator.Next() {
		t.Error("Expected Next() to return false on error")
	}
	if iterator.Err() == nil {
		t.Error("Expected error, got nil")
	}
	if iterator.Current() != nil {
		t.Error("Expected nil current item on error, got item")
	}

	iterator.Close()
}

func TestNextSettingMethodError(t *testing.T) {
	errorClient := &errorTestClient{}
	ctx := context.Background()

	iterator := &SettingsIterator{
		Client:      errorClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}

	// Test Next method with error
	if iterator.Next() {
		t.Error("Expected Next() to return false on error")
	}
	if iterator.Err() == nil {
		t.Error("Expected error, got nil")
	}
	if iterator.Current() != nil {
		t.Error("Expected nil current item on error, got item")
	}

	iterator.Close()
}

// Test BuildPaginationFilter additional edge cases
func TestBuildPaginationFilterAdditional(t *testing.T) {
	// Test with MaxItems
	opts := PaginationOptions{
		Page:     2,
		PageSize: 10,
		MaxItems: 15,
	}
	filter := BuildPaginationFilter(opts, "status=running")
	expected := "first=11 rows=5 status=running"
	if filter != expected {
		t.Errorf("Expected %s, got %s", expected, filter)
	}

	// Test with MaxItems exceeded
	opts = PaginationOptions{
		Page:     3,
		PageSize: 10,
		MaxItems: 15,
	}
	filter = BuildPaginationFilter(opts)
	expected = "rows=0"
	if filter != expected {
		t.Errorf("Expected %s, got %s", expected, filter)
	}

	// Test with empty additional filters
	opts = PaginationOptions{
		Page:     1,
		PageSize: 5,
	}
	filter = BuildPaginationFilter(opts, "", "status=running", "")
	expected = "first=1 rows=5 status=running"
	if filter != expected {
		t.Errorf("Expected %s, got %s", expected, filter)
	}
}

// Test empty current slice scenario for all iterators
func TestIteratorEmptyCurrentSlice(t *testing.T) {
	// Create a custom test client that returns empty results
	emptyClient := &emptyTestClient{}

	ctx := context.Background()

	// Test TaskIterator
	taskIterator := &TaskIterator{
		Client:      emptyClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if taskIterator.Next() {
		t.Error("Expected TaskIterator Next() to return false for empty current slice")
	}
	if taskIterator.HasMoreData {
		t.Error("Expected TaskIterator HasMoreData to be false after empty loadPage()")
	}
	taskIterator.Close()

	// Test ResultIterator
	resultIterator := &ResultIterator{
		Client:      emptyClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if resultIterator.Next() {
		t.Error("Expected ResultIterator Next() to return false for empty current slice")
	}
	if resultIterator.HasMoreData {
		t.Error("Expected ResultIterator HasMoreData to be false after empty loadPage()")
	}
	resultIterator.Close()

	// Test AssetIterator
	assetIterator := &AssetIterator{
		Client:      emptyClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if assetIterator.Next() {
		t.Error("Expected AssetIterator Next() to return false for empty current slice")
	}
	if assetIterator.HasMoreData {
		t.Error("Expected AssetIterator HasMoreData to be false after empty loadPage()")
	}
	assetIterator.Close()

	// Test TargetIterator
	targetIterator := &TargetIterator{
		Client:      emptyClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if targetIterator.Next() {
		t.Error("Expected TargetIterator Next() to return false for empty current slice")
	}
	if targetIterator.HasMoreData {
		t.Error("Expected TargetIterator HasMoreData to be false after empty loadPage()")
	}
	targetIterator.Close()

	// Test TicketIterator
	ticketIterator := &TicketIterator{
		Client:      emptyClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if ticketIterator.Next() {
		t.Error("Expected TicketIterator Next() to return false for empty current slice")
	}
	if ticketIterator.HasMoreData {
		t.Error("Expected TicketIterator HasMoreData to be false after empty loadPage()")
	}
	ticketIterator.Close()

	// Test PortListIterator
	portListIterator := &PortListIterator{
		Client:      emptyClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if portListIterator.Next() {
		t.Error("Expected PortListIterator Next() to return false for empty current slice")
	}
	if portListIterator.HasMoreData {
		t.Error("Expected PortListIterator HasMoreData to be false after empty loadPage()")
	}
	portListIterator.Close()

	// Test SettingsIterator
	settingsIterator := &SettingsIterator{
		Client:      emptyClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if settingsIterator.Next() {
		t.Error("Expected SettingsIterator Next() to return false for empty current slice")
	}
	if settingsIterator.HasMoreData {
		t.Error("Expected SettingsIterator HasMoreData to be false after empty loadPage()")
	}
	settingsIterator.Close()
}

// Test context cancellation for all iterators
func TestIteratorContextCancellation(t *testing.T) {
	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	cli := &testClient{}

	// Test TaskIterator
	taskIterator := &TaskIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if taskIterator.Next() {
		t.Error("Expected TaskIterator Next() to return false due to cancelled context")
	}
	if taskIterator.Err() != context.Canceled {
		t.Errorf("Expected TaskIterator context.Canceled error, got: %v", taskIterator.Err())
	}
	taskIterator.Close()

	// Test ResultIterator
	resultIterator := &ResultIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if resultIterator.Next() {
		t.Error("Expected ResultIterator Next() to return false due to cancelled context")
	}
	if resultIterator.Err() != context.Canceled {
		t.Errorf("Expected ResultIterator context.Canceled error, got: %v", resultIterator.Err())
	}
	resultIterator.Close()

	// Test AssetIterator
	assetIterator := &AssetIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if assetIterator.Next() {
		t.Error("Expected AssetIterator Next() to return false due to cancelled context")
	}
	if assetIterator.Err() != context.Canceled {
		t.Errorf("Expected AssetIterator context.Canceled error, got: %v", assetIterator.Err())
	}
	assetIterator.Close()

	// Test TargetIterator
	targetIterator := &TargetIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if targetIterator.Next() {
		t.Error("Expected TargetIterator Next() to return false due to cancelled context")
	}
	if targetIterator.Err() != context.Canceled {
		t.Errorf("Expected TargetIterator context.Canceled error, got: %v", targetIterator.Err())
	}
	targetIterator.Close()

	// Test TicketIterator
	ticketIterator := &TicketIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if ticketIterator.Next() {
		t.Error("Expected TicketIterator Next() to return false due to cancelled context")
	}
	if ticketIterator.Err() != context.Canceled {
		t.Errorf("Expected TicketIterator context.Canceled error, got: %v", ticketIterator.Err())
	}
	ticketIterator.Close()

	// Test PortListIterator
	portListIterator := &PortListIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if portListIterator.Next() {
		t.Error("Expected PortListIterator Next() to return false due to cancelled context")
	}
	if portListIterator.Err() != context.Canceled {
		t.Errorf("Expected PortListIterator context.Canceled error, got: %v", portListIterator.Err())
	}
	portListIterator.Close()

	// Test SettingsIterator
	settingsIterator := &SettingsIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if settingsIterator.Next() {
		t.Error("Expected SettingsIterator Next() to return false due to cancelled context")
	}
	if settingsIterator.Err() != context.Canceled {
		t.Errorf("Expected SettingsIterator context.Canceled error, got: %v", settingsIterator.Err())
	}
	settingsIterator.Close()
}

// Test loadPage error handling for all iterators
func TestIteratorLoadPageError(t *testing.T) {
	// Create an error test client that returns errors
	errorClient := &errorTestClient{}

	ctx := context.Background()

	// Test TaskIterator
	taskIterator := &TaskIterator{
		Client:      errorClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if taskIterator.Next() {
		t.Error("Expected TaskIterator Next() to return false due to loadPage error")
	}
	if taskIterator.Err() == nil {
		t.Error("Expected TaskIterator to have error from loadPage")
	}
	taskIterator.Close()

	// Test ResultIterator
	resultIterator := &ResultIterator{
		Client:      errorClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if resultIterator.Next() {
		t.Error("Expected ResultIterator Next() to return false due to loadPage error")
	}
	if resultIterator.Err() == nil {
		t.Error("Expected ResultIterator to have error from loadPage")
	}
	resultIterator.Close()

	// Test AssetIterator
	assetIterator := &AssetIterator{
		Client:      errorClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if assetIterator.Next() {
		t.Error("Expected AssetIterator Next() to return false due to loadPage error")
	}
	if assetIterator.Err() == nil {
		t.Error("Expected AssetIterator to have error from loadPage")
	}
	assetIterator.Close()

	// Test TargetIterator
	targetIterator := &TargetIterator{
		Client:      errorClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if targetIterator.Next() {
		t.Error("Expected TargetIterator Next() to return false due to loadPage error")
	}
	if targetIterator.Err() == nil {
		t.Error("Expected TargetIterator to have error from loadPage")
	}
	targetIterator.Close()

	// Test TicketIterator
	ticketIterator := &TicketIterator{
		Client:      errorClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if ticketIterator.Next() {
		t.Error("Expected TicketIterator Next() to return false due to loadPage error")
	}
	if ticketIterator.Err() == nil {
		t.Error("Expected TicketIterator to have error from loadPage")
	}
	ticketIterator.Close()

	// Test PortListIterator
	portListIterator := &PortListIterator{
		Client:      errorClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if portListIterator.Next() {
		t.Error("Expected PortListIterator Next() to return false due to loadPage error")
	}
	if portListIterator.Err() == nil {
		t.Error("Expected PortListIterator to have error from loadPage")
	}
	portListIterator.Close()

	// Test SettingsIterator
	settingsIterator := &SettingsIterator{
		Client:      errorClient,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
	}
	if settingsIterator.Next() {
		t.Error("Expected SettingsIterator Next() to return false due to loadPage error")
	}
	if settingsIterator.Err() == nil {
		t.Error("Expected SettingsIterator to have error from loadPage")
	}
	settingsIterator.Close()
}

// Test Current() method boundary conditions
func TestIteratorCurrentBoundaryConditions(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	// Test TaskIterator Current() with out-of-bounds index
	taskIterator := &TaskIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
		index:       0, // Out of bounds
		current:     []*commands.GetTasksResponseTask{},
	}

	// Current() should return nil when index is 0
	if taskIterator.Current() != nil {
		t.Error("Expected Current() to return nil when index is 0")
	}

	// Test with index > len(current)
	taskIterator.index = 5
	taskIterator.current = []*commands.GetTasksResponseTask{
		{ID: "task-1", Name: "Test Task 1"},
	}

	// Current() should return nil when index > len(current)
	if taskIterator.Current() != nil {
		t.Error("Expected Current() to return nil when index > len(current)")
	}

	taskIterator.Close()
}

// Test Next() method when HasMoreData is false
func TestIteratorNextNoMoreData(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	// Test TaskIterator Next() when HasMoreData is false
	taskIterator := &TaskIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: false, // No more data
		index:       2,     // Index beyond current slice
		current: []*commands.GetTasksResponseTask{
			{ID: "task-1", Name: "Test Task 1"},
		},
	}

	// Next() should return false when HasMoreData is false and index >= len(current)
	if taskIterator.Next() {
		t.Error("Expected Next() to return false when HasMoreData is false")
	}

	taskIterator.Close()
}

// Test Next() method with valid context (default case)
func TestIteratorNextValidContext(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	// Test TaskIterator Next() with valid context
	taskIterator := &TaskIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
		index:       0,
		current:     []*commands.GetTasksResponseTask{},
	}

	// Next() should return true and load data
	if !taskIterator.Next() {
		t.Error("Expected Next() to return true with valid context")
	}

	// Should have loaded data
	if len(taskIterator.current) == 0 {
		t.Error("Expected current to have data after Next()")
	}

	taskIterator.Close()
}

// Test Next() method with valid context for all iterators
func TestIteratorNextValidContextAll(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()

	// Test ResultIterator
	resultIterator := &ResultIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
		index:       0,
		current:     []*commands.Result{},
	}
	if !resultIterator.Next() {
		t.Error("Expected ResultIterator Next() to return true with valid context")
	}
	resultIterator.Close()

	// Test AssetIterator
	assetIterator := &AssetIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
		index:       0,
		current:     []*commands.Asset{},
	}
	if !assetIterator.Next() {
		t.Error("Expected AssetIterator Next() to return true with valid context")
	}
	assetIterator.Close()

	// Test TargetIterator
	targetIterator := &TargetIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
		index:       0,
		current:     []*commands.Target{},
	}
	if !targetIterator.Next() {
		t.Error("Expected TargetIterator Next() to return true with valid context")
	}
	targetIterator.Close()

	// Test TicketIterator
	ticketIterator := &TicketIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
		index:       0,
		current:     []*commands.Ticket{},
	}
	if !ticketIterator.Next() {
		t.Error("Expected TicketIterator Next() to return true with valid context")
	}
	ticketIterator.Close()

	// Test PortListIterator
	portListIterator := &PortListIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
		index:       0,
		current:     []*commands.PortList{},
	}
	if !portListIterator.Next() {
		t.Error("Expected PortListIterator Next() to return true with valid context")
	}
	portListIterator.Close()

	// Test SettingsIterator
	settingsIterator := &SettingsIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        1,
		HasMoreData: true,
		index:       0,
		current:     []*commands.Setting{},
	}
	if !settingsIterator.Next() {
		t.Error("Expected SettingsIterator Next() to return true with valid context")
	}
	settingsIterator.Close()
}

// emptyTestClient returns empty results to test empty current slice scenario
type emptyTestClient struct{}

func (etc *emptyTestClient) GetTasks(cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
	return &commands.GetTasksResponse{
		Status:     "200",
		StatusText: "OK",
		Task:       []commands.GetTasksResponseTask{}, // Empty slice
		TaskCount:  commands.GetTasksResponseTaskCount{Filtered: 0},
	}, nil
}

func (etc *emptyTestClient) GetResults(cmd *commands.GetResults) (*commands.GetResultsResponse, error) {
	return &commands.GetResultsResponse{
		Status:  "200",
		Results: []commands.Result{},
		Count:   &commands.ResultsCount{Filtered: 0},
	}, nil
}

func (etc *emptyTestClient) GetAssets(cmd *commands.GetAssets) (*commands.GetAssetsResponse, error) {
	return &commands.GetAssetsResponse{
		Status:     "200",
		Assets:     []commands.Asset{},
		AssetCount: &commands.AssetCount{Filtered: 0},
	}, nil
}

func (etc *emptyTestClient) GetTargets(cmd *commands.GetTargets) (*commands.GetTargetsResponse, error) {
	return &commands.GetTargetsResponse{
		Status:      "200",
		Target:      []commands.Target{},
		TargetCount: &commands.TargetCount{Filtered: 0},
	}, nil
}

func (etc *emptyTestClient) GetTickets(cmd *commands.GetTickets) (*commands.GetTicketsResponse, error) {
	return &commands.GetTicketsResponse{
		Status:      "200",
		Tickets:     []commands.Ticket{},
		TicketCount: &commands.TicketCount{Filtered: 0},
	}, nil
}

func (etc *emptyTestClient) GetPortLists(cmd *commands.GetPortLists) (*commands.GetPortListsResponse, error) {
	return &commands.GetPortListsResponse{
		Status:    "200",
		PortLists: []commands.PortList{},
		Count:     &commands.PortListCount{Filtered: 0},
	}, nil
}

func (etc *emptyTestClient) GetSettings(cmd *commands.GetSettings) (*commands.GetSettingsResponse, error) {
	return &commands.GetSettingsResponse{
		Status:       "200",
		SettingsList: []commands.Setting{},
		SettingCount: &commands.SettingCount{Filtered: 0},
	}, nil
}
