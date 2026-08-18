package pagination

import (
	"context"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

// Minimal fake client that implements Client and returns controlled responses.
type fakeAllClient struct {
	assetsResp    *commands.GetAssetsResponse
	resultsResp   *commands.GetResultsResponse
	targetsResp   *commands.GetTargetsResponse
	ticketsResp   *commands.GetTicketsResponse
	portlistsResp *commands.GetPortListsResponse
	settingsResp  *commands.GetSettingsResponse
}

func (f *fakeAllClient) GetTasksRaw(ctx context.Context, cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
	return &commands.GetTasksResponse{}, nil
}

func (f *fakeAllClient) GetResultsRaw(ctx context.Context, cmd *commands.GetResults) (*commands.GetResultsResponse, error) {
	return f.resultsResp, nil
}

func (f *fakeAllClient) GetAssetsRaw(ctx context.Context, cmd *commands.GetAssets) (*commands.GetAssetsResponse, error) {
	return f.assetsResp, nil
}

func (f *fakeAllClient) GetTargetsRaw(ctx context.Context, cmd *commands.GetTargets) (*commands.GetTargetsResponse, error) {
	return f.targetsResp, nil
}

func (f *fakeAllClient) GetTicketsRaw(ctx context.Context, cmd *commands.GetTickets) (*commands.GetTicketsResponse, error) {
	return f.ticketsResp, nil
}

func (f *fakeAllClient) GetPortListsRaw(ctx context.Context, cmd *commands.GetPortLists) (*commands.GetPortListsResponse, error) {
	return f.portlistsResp, nil
}

func (f *fakeAllClient) GetSettingsRaw(ctx context.Context, cmd *commands.GetSettings) (*commands.GetSettingsResponse, error) {
	return f.settingsResp, nil
}

func TestNewIterator_DefaultPageSizeAndHasMore_Computed(t *testing.T) {
	// PageSize < 0 should default to 100; provide exactly 100 items to set hasMore=true path
	client := &fakeAllClient{
		assetsResp: &commands.GetAssetsResponse{Asset: make([]commands.GetAssetsAsset, 100)},
	}
	it := NewAssetIterator(client, context.Background(), PaginationOptions{Page: 1, PageSize: -1})
	if !it.Next() {
		t.Fatalf("expected Next() true to load first page with defaulted page size")
	}
	if it.Err() != nil {
		t.Fatalf("unexpected error: %v", it.Err())
	}
	// do not continue iterating; hasMore could remain true with fixed response
}

func TestResultsIterator_TotalNilCount(t *testing.T) {
	client := &fakeAllClient{
		resultsResp: &commands.GetResultsResponse{Results: []commands.Result{{ID: "r1"}, {ID: "r2"}}, Count: nil},
	}
	it := NewResultIterator(client, context.Background(), PaginationOptions{Page: 1, PageSize: 3})
	if !it.Next() {
		t.Fatalf("expected first Next() true")
	}
	_ = it.Current()
	for it.Next() {
		_ = it.Current()
	}
	if it.Total() != 0 {
		t.Fatalf("expected total 0 with nil Count, got %d", it.Total())
	}
}

func TestSettingsIterator_TotalZero(t *testing.T) {
	client := &fakeAllClient{
		settingsResp: &commands.GetSettingsResponse{Setting: []commands.SettingsSetting{{ID: "s1"}, {ID: "s2"}}},
	}
	it := NewSettingsIterator(client, context.Background(), PaginationOptions{Page: 1, PageSize: 3})
	if !it.Next() {
		t.Fatalf("expected first Next() true")
	}
	for it.Next() {
		_ = it.Current()
	}
	if it.Total() != 0 {
		t.Fatalf("expected total 0 with zero-value SettingCount, got %d", it.Total())
	}
}

// Type mismatch in internal executeAndParse: call with wrong cmd type to hit assertion path.
func TestNewIterator_TypeMismatch_Internal(t *testing.T) {
	client := &fakeAllClient{resultsResp: &commands.GetResultsResponse{Results: []commands.Result{{ID: "r1"}}}}
	it := NewResultIterator(client, context.Background(), PaginationOptions{Page: 1, PageSize: 1})
	// Directly call the internal closure with a wrong command type to exercise the error path.
	_, _, _, err := it.executeAndParse(&commands.GetTasks{}, PaginationOptions{Page: 1, PageSize: 1})
	if err == nil {
		t.Fatalf("expected type mismatch error, got nil")
	}
}
