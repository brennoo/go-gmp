//nolint:dupl

package pagination

import (
	"context"

	"github.com/brennoo/go-gmp/commands"
)

// AssetIterator is the iterator over assets.
type AssetIterator = Iterator[*commands.Asset]

// NewAssetIterator constructs an AssetIterator configured for GetAssets.
func NewAssetIterator(client Client, ctx context.Context, opts PaginationOptions, filters ...string) *AssetIterator {
	return newIterator[*commands.Asset, *commands.GetAssets, *commands.GetAssetsResponse](
		client,
		ctx,
		opts,
		filters,
		func(filter string) *commands.GetAssets { return &commands.GetAssets{Filter: filter} },
		func(client Client, cmd *commands.GetAssets) (*commands.GetAssetsResponse, error) {
			return client.GetAssetsRaw(cmd)
		},
		func(resp *commands.GetAssetsResponse) []*commands.Asset {
			items := make([]*commands.Asset, len(resp.Assets))
			for i := range resp.Assets {
				items[i] = &resp.Assets[i]
			}
			return items
		},
		func(resp *commands.GetAssetsResponse) int {
			if resp.AssetCount != nil {
				return resp.AssetCount.Filtered
			}
			return 0
		},
	)
}
