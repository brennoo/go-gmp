package pagination

import (
	"context"
	"errors"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

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
		Page:        0,
		HasMoreData: true,
	}

	// Test iteration
	var settings []*commands.Setting
	for iterator.Next() {
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

	// Test individual setting properties
	if len(settings) > 0 {
		if settings[0].ID != "setting-1" {
			t.Errorf("Expected first setting ID 'setting-1', got '%s'", settings[0].ID)
		}
		if settings[0].Name != "Test Setting 1" {
			t.Errorf("Expected first setting name 'Test Setting 1', got '%s'", settings[0].Name)
		}
	}
}

func TestSettingsIterator_PaginationSuccess(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &SettingsIterator{
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

	setting := iterator.Current()
	if setting == nil {
		t.Fatal("Expected Current() to return a setting")
	}

	if setting.ID != "setting-1" {
		t.Errorf("Expected setting ID 'setting-1', got '%s'", setting.ID)
	}

	// Test second page
	if !iterator.Next() {
		t.Fatal("Expected Next() to return true for second page")
	}

	setting = iterator.Current()
	if setting == nil {
		t.Fatal("Expected Current() to return a setting")
	}

	if setting.ID != "setting-2" {
		t.Errorf("Expected setting ID 'setting-2', got '%s'", setting.ID)
	}
}

func TestSettingsIterator_EmptyFirstPageStops(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &SettingsIterator{
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

func TestSettingsIterator_ErrorStops(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &SettingsIterator{
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

func TestSettingsIterator_ContextCanceled(t *testing.T) {
	cli := &testClient{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	iterator := &SettingsIterator{
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

func TestSettingsIterator_CurrentBeforeNext(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &SettingsIterator{
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

func TestSettingsIterator_CloseStopsIteration(t *testing.T) {
	cli := &testClient{}
	ctx := context.Background()
	iterator := &SettingsIterator{
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
