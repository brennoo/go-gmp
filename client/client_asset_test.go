package client

import (
	"context"
	"testing"

	"github.com/brennoo/go-gmp/commands"
	"github.com/brennoo/go-gmp/commands/filtering"
)

func TestCreateAsset(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.CreateAsset{
		Asset: &commands.AssetInput{
			Name:    "Test Asset",
			Comment: "A test asset",
			Type:    "host",
		},
	}
	resp, err := cli.CreateAsset(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateAsset: %s", err)
	}
	if resp.Status != "201" {
		t.Fatalf("Unexpected status. Expected: 201 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK, resource created" {
		t.Fatalf("Unexpected status text. Expected: OK, resource created Got: %s", resp.StatusText)
	}
	if resp.ID != "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
		t.Fatalf("Unexpected ID. Expected: 254cd3ef-bbe1-4d58-859d-21b8d0c046c6 Got: %s", resp.ID)
	}
}

func TestCreateAssetWithReport(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmdWithReport := &commands.CreateAsset{
		Report: &commands.AssetReport{
			ID:     "report-id",
			Filter: &commands.AssetFilter{Term: "term"},
		},
	}
	resp, err := cli.CreateAsset(cmdWithReport)
	if err != nil {
		t.Fatalf("Unexpected error during CreateAsset (report): %s", err)
	}
	if resp.Status != "201" {
		t.Fatalf("Unexpected status. Expected: 201 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK, resource created" {
		t.Fatalf("Unexpected status text. Expected: OK, resource created Got: %s", resp.StatusText)
	}
	if resp.ID != "report-asset-uuid" {
		t.Fatalf("Unexpected ID. Expected: report-asset-uuid Got: %s", resp.ID)
	}
}

func TestModifyAsset(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.ModifyAsset{
		AssetID: "914b59f8-25f5-4c8f-832c-2379cd625236",
		Comment: "New comment",
	}
	resp, err := cli.ModifyAsset(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyAsset: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
}

func TestGetAssets(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()
	resp, err := cli.GetAssets(ctx, "asset_id=b493b7a8-7489-11df-a3ec-002264764cea")
	if err != nil {
		t.Fatalf("Unexpected error during GetAssets: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if len(resp.Assets) != 2 {
		t.Fatalf("Expected 2 assets, got %d", len(resp.Assets))
	}
	// Check first asset
	asset := resp.Assets[0]
	if asset.ID != "asset-1" {
		t.Fatalf("Unexpected asset ID. Expected: asset-1 Got: %s", asset.ID)
	}
	if asset.Name != "Test Asset 1" {
		t.Fatalf("Unexpected asset name. Expected: Test Asset 1 Got: %s", asset.Name)
	}
}

func TestDeleteAsset(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DeleteAsset{
		AssetID: "267a3405-e84a-47da-97b2-5fa0d2e8995e",
	}
	resp, err := cli.DeleteAsset(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteAsset: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
}

// TestGetAssetsConsolidated tests the GetAssets method.
func TestGetAssetsConsolidated(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()

	// Test with no filters
	resp, err := cli.GetAssets(ctx)
	if err != nil {
		t.Fatalf("Unexpected error during GetAssets: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with string filters
	resp, err = cli.GetAssets(ctx, "name~localhost")
	if err != nil {
		t.Fatalf("Unexpected error during GetAssets with string filter: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with functional options
	resp, err = cli.GetAssets(ctx, filtering.WithName("test"))
	if err != nil {
		t.Fatalf("Unexpected error during GetAssets with functional options: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}
