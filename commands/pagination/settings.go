//nolint:dupl

package pagination

import (
	"context"

	"github.com/brennoo/go-gmp/commands"
)

// SettingsIterator is the iterator over settings.
type SettingsIterator = Iterator[*commands.SettingsSetting]

// NewSettingsIterator constructs a SettingsIterator configured for GetSettings.
func NewSettingsIterator(client Client, ctx context.Context, opts PaginationOptions, filters ...string) *SettingsIterator {
	return newIterator[*commands.SettingsSetting, *commands.GetSettings, *commands.GetSettingsResponse](
		client,
		ctx,
		opts,
		filters,
		func(filter string) *commands.GetSettings { return &commands.GetSettings{Filter: filter} },
		func(client Client, ctx context.Context, cmd *commands.GetSettings) (*commands.GetSettingsResponse, error) {
			return client.GetSettingsRaw(ctx, cmd)
		},
		func(resp *commands.GetSettingsResponse) []*commands.SettingsSetting {
			// Prefer the live top-level layout; fall back to the schema example's nested layout.
			list := resp.Setting
			if len(list) == 0 {
				list = resp.Settings.Setting
			}
			items := make([]*commands.SettingsSetting, len(list))
			for i := range list {
				items[i] = &list[i]
			}
			return items
		},
		func(resp *commands.GetSettingsResponse) int { return resp.SettingCount.Filtered },
	)
}
