//nolint:dupl

package pagination

import (
	"context"

	"github.com/brennoo/go-gmp/commands"
)

// AssetIterator is the iterator over assets.
type AssetIterator = Iterator[*commands.GetAssetsAsset]

// NewAssetIterator constructs an AssetIterator configured for GetAssets.
func NewAssetIterator(client Client, ctx context.Context, opts PaginationOptions, filters ...string) *AssetIterator {
	return newIterator[*commands.GetAssetsAsset, *commands.GetAssets, *commands.GetAssetsResponse](
		client,
		ctx,
		opts,
		filters,
		func(filter string) *commands.GetAssets { return &commands.GetAssets{Filter: filter} },
		func(client Client, ctx context.Context, cmd *commands.GetAssets) (*commands.GetAssetsResponse, error) {
			return client.GetAssetsRaw(ctx, cmd)
		},
		func(resp *commands.GetAssetsResponse) []*commands.GetAssetsAsset {
			items := make([]*commands.GetAssetsAsset, len(resp.Asset))
			for i := range resp.Asset {
				items[i] = &resp.Asset[i]
			}
			return items
		},
		func(resp *commands.GetAssetsResponse) int { return resp.AssetCount.Filtered },
	)
}
