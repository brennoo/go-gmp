//nolint:dupl // Similar patterns across option types are intentional

package filtering

import (
	"fmt"
)

// WithAssetID adds an asset_id filter condition.
func WithAssetID(assetID string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("asset_id", OpEqual, assetID)
	}
}

// WithAssetName adds an asset name filter condition.
func WithAssetName(name string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("asset_name", OpEqual, name)
	}
}

// WithAssetNameLike adds an asset name filter condition using LIKE operator.
func WithAssetNameLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("asset_name", OpLike, pattern)
	}
}

// WithAssetType adds an asset type filter condition.
func WithAssetType(assetType string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("asset_type", OpEqual, assetType)
	}
}

// WithAssetTypeLike adds an asset type filter condition using LIKE operator.
func WithAssetTypeLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("asset_type", OpLike, pattern)
	}
}

// WithAssetHost adds a host filter condition.
func WithAssetHost(host string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host", OpEqual, host)
	}
}

// WithAssetHostLike adds a host filter condition using LIKE operator.
func WithAssetHostLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host", OpLike, pattern)
	}
}

// WithAssetOS adds an OS filter condition.
func WithAssetOS(os string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("os", OpEqual, os)
	}
}

// WithAssetOSLike adds an OS filter condition using LIKE operator.
func WithAssetOSLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("os", OpLike, pattern)
	}
}

// WithAssetLocation adds a location filter condition.
func WithAssetLocation(location string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("location", OpEqual, location)
	}
}

// WithAssetLocationLike adds a location filter condition using LIKE operator.
func WithAssetLocationLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("location", OpLike, pattern)
	}
}

// WithAssetHostCount adds a host_count filter condition.
func WithAssetHostCount(count int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host_count", OpEqual, fmt.Sprintf("%d", count))
	}
}

// WithAssetHostCountGreaterThan adds a host_count filter condition.
func WithAssetHostCountGreaterThan(count int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host_count", OpGreaterThan, fmt.Sprintf("%d", count))
	}
}

// WithAssetHostCountGreaterThanOrEqual adds a host_count filter condition.
func WithAssetHostCountGreaterThanOrEqual(count int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host_count", OpGreaterThanOrEqual, fmt.Sprintf("%d", count))
	}
}

// WithAssetHostCountLessThan adds a host_count filter condition.
func WithAssetHostCountLessThan(count int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host_count", OpLessThan, fmt.Sprintf("%d", count))
	}
}

// WithAssetHostCountLessThanOrEqual adds a host_count filter condition.
func WithAssetHostCountLessThanOrEqual(count int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host_count", OpLessThanOrEqual, fmt.Sprintf("%d", count))
	}
}

// WithAssetHosts adds a hosts filter condition.
func WithAssetHosts(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpEqual, fmt.Sprintf("%d", hosts))
	}
}

// WithAssetHostsGreaterThan adds a hosts filter condition.
func WithAssetHostsGreaterThan(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpGreaterThan, fmt.Sprintf("%d", hosts))
	}
}

// WithAssetHostsGreaterThanOrEqual adds a hosts filter condition.
func WithAssetHostsGreaterThanOrEqual(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpGreaterThanOrEqual, fmt.Sprintf("%d", hosts))
	}
}

// WithAssetHostsLessThan adds a hosts filter condition.
func WithAssetHostsLessThan(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpLessThan, fmt.Sprintf("%d", hosts))
	}
}

// WithAssetHostsLessThanOrEqual adds a hosts filter condition.
func WithAssetHostsLessThanOrEqual(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpLessThanOrEqual, fmt.Sprintf("%d", hosts))
	}
}

// WithAssetMaxHosts adds a max_hosts filter condition.
func WithAssetMaxHosts(maxHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("max_hosts", OpEqual, fmt.Sprintf("%d", maxHosts))
	}
}

// WithAssetMaxHostsGreaterThan adds a max_hosts filter condition.
func WithAssetMaxHostsGreaterThan(maxHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("max_hosts", OpGreaterThan, fmt.Sprintf("%d", maxHosts))
	}
}

// WithAssetMaxHostsGreaterThanOrEqual adds a max_hosts filter condition.
func WithAssetMaxHostsGreaterThanOrEqual(maxHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("max_hosts", OpGreaterThanOrEqual, fmt.Sprintf("%d", maxHosts))
	}
}

// WithAssetMaxHostsLessThan adds a max_hosts filter condition.
func WithAssetMaxHostsLessThan(maxHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("max_hosts", OpLessThan, fmt.Sprintf("%d", maxHosts))
	}
}

// WithAssetMaxHostsLessThanOrEqual adds a max_hosts filter condition.
func WithAssetMaxHostsLessThanOrEqual(maxHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("max_hosts", OpLessThanOrEqual, fmt.Sprintf("%d", maxHosts))
	}
}

// WithHosts adds a hosts filter condition.
func WithHosts(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpEqual, fmt.Sprintf("%d", hosts))
	}
}

// WithHostsGreaterThan adds a hosts filter condition.
func WithHostsGreaterThan(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpGreaterThan, fmt.Sprintf("%d", hosts))
	}
}

// WithHostsGreaterThanOrEqual adds a hosts filter condition.
func WithHostsGreaterThanOrEqual(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpGreaterThanOrEqual, fmt.Sprintf("%d", hosts))
	}
}

// WithHostsLessThan adds a hosts filter condition.
func WithHostsLessThan(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpLessThan, fmt.Sprintf("%d", hosts))
	}
}

// WithHostsLessThanOrEqual adds a hosts filter condition.
func WithHostsLessThanOrEqual(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpLessThanOrEqual, fmt.Sprintf("%d", hosts))
	}
}

// WithOS adds an OS filter condition.
func WithOS(os string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("os", OpEqual, os)
	}
}

// WithOSLike adds an OS filter condition using LIKE operator.
func WithOSLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("os", OpLike, pattern)
	}
}

// WithLocation adds a location filter condition.
func WithLocation(location string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("location", OpEqual, location)
	}
}

// WithLocationLike adds a location filter condition using LIKE operator.
func WithLocationLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("location", OpLike, pattern)
	}
}

// WithHostCount adds a host count filter condition.
func WithHostCount(count int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host_count", OpEqual, fmt.Sprintf("%d", count))
	}
}

// WithHostCountGreaterThan adds a host count filter condition.
func WithHostCountGreaterThan(count int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host_count", OpGreaterThan, fmt.Sprintf("%d", count))
	}
}

// WithHostCountGreaterThanOrEqual adds a host count filter condition.
func WithHostCountGreaterThanOrEqual(count int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host_count", OpGreaterThanOrEqual, fmt.Sprintf("%d", count))
	}
}

// WithHostCountLessThan adds a host count filter condition.
func WithHostCountLessThan(count int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host_count", OpLessThan, fmt.Sprintf("%d", count))
	}
}

// WithHostCountLessThanOrEqual adds a host count filter condition.
func WithHostCountLessThanOrEqual(count int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host_count", OpLessThanOrEqual, fmt.Sprintf("%d", count))
	}
}

// WithMaxHosts adds a max hosts filter condition.
func WithMaxHosts(maxHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("max_hosts", OpEqual, fmt.Sprintf("%d", maxHosts))
	}
}

// WithMaxHostsGreaterThan adds a max hosts filter condition.
func WithMaxHostsGreaterThan(maxHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("max_hosts", OpGreaterThan, fmt.Sprintf("%d", maxHosts))
	}
}

// WithMaxHostsGreaterThanOrEqual adds a max hosts filter condition.
func WithMaxHostsGreaterThanOrEqual(maxHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("max_hosts", OpGreaterThanOrEqual, fmt.Sprintf("%d", maxHosts))
	}
}

// WithMaxHostsLessThan adds a max hosts filter condition.
func WithMaxHostsLessThan(maxHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("max_hosts", OpLessThan, fmt.Sprintf("%d", maxHosts))
	}
}

// WithMaxHostsLessThanOrEqual adds a max hosts filter condition.
func WithMaxHostsLessThanOrEqual(maxHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("max_hosts", OpLessThanOrEqual, fmt.Sprintf("%d", maxHosts))
	}
}

// WithPorts adds a ports filter condition.
func WithPorts(ports int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ports", OpEqual, fmt.Sprintf("%d", ports))
	}
}

// WithPortsGreaterThan adds a ports filter condition.
func WithPortsGreaterThan(ports int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ports", OpGreaterThan, fmt.Sprintf("%d", ports))
	}
}

// WithPortsGreaterThanOrEqual adds a ports filter condition.
func WithPortsGreaterThanOrEqual(ports int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ports", OpGreaterThanOrEqual, fmt.Sprintf("%d", ports))
	}
}

// WithPortsLessThan adds a ports filter condition.
func WithPortsLessThan(ports int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ports", OpLessThan, fmt.Sprintf("%d", ports))
	}
}

// WithPortsLessThanOrEqual adds a ports filter condition.
func WithPortsLessThanOrEqual(ports int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ports", OpLessThanOrEqual, fmt.Sprintf("%d", ports))
	}
}
