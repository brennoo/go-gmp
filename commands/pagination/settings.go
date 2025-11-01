//nolint:dupl

package pagination

import (
	"context"

	"github.com/brennoo/go-gmp/commands"
)

// SettingsIterator is the iterator over settings.
type SettingsIterator = Iterator[*commands.Setting]

// NewSettingsIterator constructs a SettingsIterator configured for GetSettings.
func NewSettingsIterator(client Client, ctx context.Context, opts PaginationOptions, filters ...string) *SettingsIterator {
	return newIterator[*commands.Setting, *commands.GetSettings, *commands.GetSettingsResponse](
		client,
		ctx,
		opts,
		filters,
		func(filter string) *commands.GetSettings { return &commands.GetSettings{Filter: filter} },
		func(client Client, ctx context.Context, cmd *commands.GetSettings) (*commands.GetSettingsResponse, error) {
			return client.GetSettingsRaw(ctx, cmd)
		},
		func(resp *commands.GetSettingsResponse) []*commands.Setting {
			var list []commands.Setting
			if len(resp.SettingsList) > 0 {
				list = resp.SettingsList
			} else if resp.Settings != nil {
				list = resp.Settings.Setting
			}
			items := make([]*commands.Setting, len(list))
			for i := range list {
				items[i] = &list[i]
			}
			return items
		},
		func(resp *commands.GetSettingsResponse) int {
			if resp.SettingCount != nil {
				return resp.SettingCount.Filtered
			}
			return 0
		},
	)
}
