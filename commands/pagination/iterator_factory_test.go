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

func (f *fakeAllClient) GetTasksRaw(cmd *commands.GetTasks) (*commands.GetTasksResponse, error) {
	return &commands.GetTasksResponse{}, nil
}

func (f *fakeAllClient) GetResultsRaw(cmd *commands.GetResults) (*commands.GetResultsResponse, error) {
	return f.resultsResp, nil
}

func (f *fakeAllClient) GetAssetsRaw(cmd *commands.GetAssets) (*commands.GetAssetsResponse, error) {
	return f.assetsResp, nil
}

func (f *fakeAllClient) GetTargetsRaw(cmd *commands.GetTargets) (*commands.GetTargetsResponse, error) {
	return f.targetsResp, nil
}

func (f *fakeAllClient) GetTicketsRaw(cmd *commands.GetTickets) (*commands.GetTicketsResponse, error) {
	return f.ticketsResp, nil
}

func (f *fakeAllClient) GetPortListsRaw(cmd *commands.GetPortLists) (*commands.GetPortListsResponse, error) {
	return f.portlistsResp, nil
}

func (f *fakeAllClient) GetSettingsRaw(cmd *commands.GetSettings) (*commands.GetSettingsResponse, error) {
	return f.settingsResp, nil
}

func TestNewIterator_DefaultPageSizeAndHasMore_Computed(t *testing.T) {
	// PageSize < 0 should default to 100; provide exactly 100 items to set hasMore=true path
	client := &fakeAllClient{
		assetsResp: &commands.GetAssetsResponse{Assets: make([]commands.Asset, 100)},
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

func TestAssetsIterator_TotalNilCount(t *testing.T) {
	client := &fakeAllClient{
		assetsResp: &commands.GetAssetsResponse{Assets: []commands.Asset{{ID: "a1"}, {ID: "a2"}}, AssetCount: nil},
	}
	// PageSize > len(items) to force hasMore=false
	it := NewAssetIterator(client, context.Background(), PaginationOptions{Page: 1, PageSize: 3})
	if !it.Next() {
		t.Fatalf("expected first Next() true")
	}
	_ = it.Current()
	// Next should exhaust and then stop
	for it.Next() {
		_ = it.Current()
	}
	if it.Total() != 0 {
		t.Fatalf("expected total 0 with nil AssetCount, got %d", it.Total())
	}
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

func TestTargetsIterator_TotalNilCount(t *testing.T) {
	client := &fakeAllClient{
		targetsResp: &commands.GetTargetsResponse{Target: []commands.Target{{ID: "t1"}}, TargetCount: nil},
	}
	it := NewTargetIterator(client, context.Background(), PaginationOptions{Page: 1, PageSize: 2})
	if !it.Next() {
		t.Fatalf("expected first Next() true")
	}
	if it.Total() != 0 {
		t.Fatalf("expected total 0 with nil TargetCount, got %d", it.Total())
	}
}

func TestTicketsIterator_TotalNilCount(t *testing.T) {
	client := &fakeAllClient{
		ticketsResp: &commands.GetTicketsResponse{Tickets: []commands.Ticket{{ID: "k1"}}, TicketCount: nil},
	}
	it := NewTicketIterator(client, context.Background(), PaginationOptions{Page: 1, PageSize: 2})
	if !it.Next() {
		t.Fatalf("expected first Next() true")
	}
	if it.Total() != 0 {
		t.Fatalf("expected total 0 with nil TicketCount, got %d", it.Total())
	}
}

func TestPortListsIterator_TotalNilCount(t *testing.T) {
	client := &fakeAllClient{
		portlistsResp: &commands.GetPortListsResponse{PortLists: []commands.PortList{{ID: "p1"}}, Count: nil},
	}
	it := NewPortListIterator(client, context.Background(), PaginationOptions{Page: 1, PageSize: 2})
	if !it.Next() {
		t.Fatalf("expected first Next() true")
	}
	if it.Total() != 0 {
		t.Fatalf("expected total 0 with nil PortList Count, got %d", it.Total())
	}
}

func TestSettingsIterator_Layouts_TotalNil(t *testing.T) {
	// SettingsList variant
	client1 := &fakeAllClient{
		settingsResp: &commands.GetSettingsResponse{SettingsList: []commands.Setting{{ID: "s1"}, {ID: "s2"}}, SettingCount: nil},
	}
	it1 := NewSettingsIterator(client1, context.Background(), PaginationOptions{Page: 1, PageSize: 3})
	if !it1.Next() {
		t.Fatalf("expected first Next() true")
	}
	for it1.Next() {
		_ = it1.Current()
	}
	if it1.Total() != 0 {
		t.Fatalf("expected total 0 with nil SettingCount (list), got %d", it1.Total())
	}
	// Settings.Settings variant
	client2 := &fakeAllClient{
		settingsResp: &commands.GetSettingsResponse{Settings: &commands.Settings{Setting: []commands.Setting{{ID: "s3"}}}, SettingCount: nil},
	}
	it2 := NewSettingsIterator(client2, context.Background(), PaginationOptions{Page: 1, PageSize: 2})
	if !it2.Next() {
		t.Fatalf("expected first Next() true")
	}
	if it2.Total() != 0 {
		t.Fatalf("expected total 0 with nil SettingCount (nested), got %d", it2.Total())
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
