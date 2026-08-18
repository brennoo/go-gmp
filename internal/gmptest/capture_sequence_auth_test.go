//go:build integration

package gmptest

import (
	"slices"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestDiffAuthSettingKeys(t *testing.T) {
	setting := func(key, value string) commands.AuthConfSetting {
		return commands.AuthConfSetting{Key: key, Value: value}
	}

	tests := []struct {
		name   string
		before []commands.AuthConfSetting
		after  []commands.AuthConfSetting
		want   []string
	}{
		{
			name:   "equal settings in different orders",
			before: []commands.AuthConfSetting{setting("enable", "true"), setting("order", "1")},
			after:  []commands.AuthConfSetting{setting("order", "1"), setting("enable", "true")},
			want:   nil,
		},
		{
			name:   "changed value",
			before: []commands.AuthConfSetting{setting("enable", "true")},
			after:  []commands.AuthConfSetting{setting("enable", "false")},
			want:   []string{"enable"},
		},
		{
			name:   "missing key",
			before: []commands.AuthConfSetting{setting("enable", "true"), setting("order", "1")},
			after:  []commands.AuthConfSetting{setting("enable", "true")},
			want:   []string{"order"},
		},
		{
			name:   "added key",
			before: []commands.AuthConfSetting{setting("enable", "true")},
			after:  []commands.AuthConfSetting{setting("enable", "true"), setting("order", "1")},
			want:   []string{"order"},
		},
		{
			name:   "duplicate keys with the same multiset are a no-op",
			before: []commands.AuthConfSetting{setting("enable", "true"), setting("enable", "false")},
			after:  []commands.AuthConfSetting{setting("enable", "false"), setting("enable", "true")},
			want:   nil,
		},
		{
			name:   "duplicate keys collapsing into one value is flagged",
			before: []commands.AuthConfSetting{setting("enable", "true"), setting("enable", "false")},
			after:  []commands.AuthConfSetting{setting("enable", "false")},
			want:   []string{"enable"},
		},
		{
			name:   "duplicate key count changes with the same value is flagged",
			before: []commands.AuthConfSetting{setting("enable", "true")},
			after:  []commands.AuthConfSetting{setting("enable", "true"), setting("enable", "true")},
			want:   []string{"enable"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := diffAuthSettingKeys(tt.before, tt.after)
			if !slices.Equal(got, tt.want) {
				t.Errorf("diffAuthSettingKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
