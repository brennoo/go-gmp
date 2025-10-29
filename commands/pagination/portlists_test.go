package pagination

import (
	"context"
	"errors"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

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
		Page:        0,
		HasMoreData: true,
	}

	// Test iteration
	var portLists []*commands.PortList
	for iterator.Next() {
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

	// Test individual port list properties
	if len(portLists) > 0 {
		if portLists[0].ID != "portlist-1" {
			t.Errorf("Expected first port list ID 'portlist-1', got '%s'", portLists[0].ID)
		}
		if portLists[0].Name != "Test PortList 1" {
			t.Errorf("Expected first port list name 'Test PortList 1', got '%s'", portLists[0].Name)
		}
	}
}

func TestPortListIterator_PaginationSuccess(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &PortListIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 1},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	// Test first page
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for first page")
	}

	portList := iterator.Current()
	if portList == nil {
		t.Fatal("Expected Current() to return a port list")
	}

	if portList.ID != "portlist-1" {
		t.Errorf("Expected port list ID 'portlist-1', got '%s'", portList.ID)
	}

	// Test second page
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for second page")
	}

	portList = iterator.Current()
	if portList == nil {
		t.Fatal("Expected Current() to return a port list")
	}

	if portList.ID != "portlist-2" {
		t.Errorf("Expected port list ID 'portlist-2', got '%s'", portList.ID)
	}
}

func TestPortListIterator_EmptyFirstPageStops(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &PortListIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 0},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	// Should not advance on empty page
	if iterator.Next() {
		t.Error("Expected Next() to return false for empty page")
	}

	if iterator.Current() != nil {
		t.Error("Expected Current() to return nil for empty page")
	}
}

func TestPortListIterator_ErrorStops(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &PortListIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	// Simulate error by setting err field directly
	iterator.err = errors.New("test error")

	if iterator.Next() {
		t.Error("Expected Next() to return false when error is set")
	}

	if iterator.Err() == nil {
		t.Error("Expected Err() to return the error")
	}
}

func TestPortListIterator_ContextCanceled(t *testing.T) {
	cli := &testClient{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	iterator := &PortListIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	if iterator.Next() {
		t.Error("Expected Next() to return false when context is canceled")
	}

	if iterator.Err() == nil {
		t.Error("Expected Err() to return context canceled error")
	}
}

func TestPortListIterator_CurrentBeforeNext(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &PortListIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	// Current() should return nil before Next() is called
	if iterator.Current() != nil {
		t.Error("Expected Current() to return nil before Next() is called")
	}
}

func TestPortListIterator_CloseStopsIteration(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &PortListIterator{
		Client:      cli,
		Ctx:         ctx,
		Opts:        PaginationOptions{Page: 1, PageSize: 2},
		Filters:     []string{},
		Page:        0,
		HasMoreData: true,
	}

	// Close the iterator
	iterator.Close()

	// Should not advance after close
	if iterator.Next() {
		t.Error("Expected Next() to return false after Close()")
	}

	// HasMore should return false after close
	if iterator.HasMore() {
		t.Error("Expected HasMore() to return false after Close()")
	}
}
