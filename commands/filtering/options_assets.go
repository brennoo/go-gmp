package filtering

import (
	"time"
)

// WithAssetID adds an asset_id filter condition.
func WithAssetID(assetID string) Option {
	return WithField(AssetColumnUUID, assetID)
}

// WithAssetName adds an asset name filter condition.
func WithAssetName(name string) Option {
	return WithField(AssetColumnName, name)
}

// WithAssetNameLike adds an asset name filter condition using LIKE operator.
func WithAssetNameLike(pattern string) Option {
	return WithFieldLike(AssetColumnName, pattern)
}

// WithAssetType adds an asset type filter condition.
func WithAssetType(assetType string) Option {
	return WithField(AssetColumnType, assetType)
}

// WithAssetTypeLike adds an asset type filter condition using LIKE operator.
func WithAssetTypeLike(pattern string) Option {
	return WithFieldLike(AssetColumnType, pattern)
}

// WithAssetHost adds a host filter condition.
func WithAssetHost(host string) Option {
	return WithField(AssetColumnHost, host)
}

// WithAssetHostLike adds a host filter condition using LIKE operator.
func WithAssetHostLike(pattern string) Option {
	return WithFieldLike(AssetColumnHost, pattern)
}

// WithAssetOS adds an OS filter condition.
func WithAssetOS(os string) Option {
	return WithField(AssetColumnOS, os)
}

// WithAssetOSLike adds an OS filter condition using LIKE operator.
func WithAssetOSLike(pattern string) Option {
	return WithFieldLike(AssetColumnOS, pattern)
}

// WithAssetLocation adds a location filter condition.
func WithAssetLocation(location string) Option {
	return WithField(AssetColumnLocation, location)
}

// WithAssetLocationLike adds a location filter condition using LIKE operator.
func WithAssetLocationLike(pattern string) Option {
	return WithFieldLike(AssetColumnLocation, pattern)
}

// OS Asset filtering options (when type="os")

// WithOSTitle adds a title filter condition for OS assets.
func WithOSTitle(title string) Option {
	return WithField(OSColumnTitle, title)
}

// WithOSTitleLike adds a title filter condition using LIKE operator for OS assets.
func WithOSTitleLike(pattern string) Option {
	return WithFieldLike(OSColumnTitle, pattern)
}

// WithOSHosts adds a hosts count filter condition for OS assets.
func WithOSHosts(hosts int) Option {
	return WithNumericField(OSColumnHosts, hosts)
}

// WithOSHostsGreaterThan adds a hosts count filter condition for OS assets.
func WithOSHostsGreaterThan(hosts int) Option {
	return WithNumericFieldGreaterThan(OSColumnHosts, hosts)
}

// WithOSHostsGreaterThanOrEqual adds a hosts count filter condition for OS assets.
func WithOSHostsGreaterThanOrEqual(hosts int) Option {
	return WithNumericFieldGreaterThanOrEqual(OSColumnHosts, hosts)
}

// WithOSHostsLessThan adds a hosts count filter condition for OS assets.
func WithOSHostsLessThan(hosts int) Option {
	return WithNumericFieldLessThan(OSColumnHosts, hosts)
}

// WithOSHostsLessThanOrEqual adds a hosts count filter condition for OS assets.
func WithOSHostsLessThanOrEqual(hosts int) Option {
	return WithNumericFieldLessThanOrEqual(OSColumnHosts, hosts)
}

// WithOSLatestSeverity adds a latest severity filter condition for OS assets.
func WithOSLatestSeverity(severity string) Option {
	return WithField(OSColumnLatestSeverity, severity)
}

// WithOSHighestSeverity adds a highest severity filter condition for OS assets.
func WithOSHighestSeverity(severity string) Option {
	return WithField(OSColumnHighestSeverity, severity)
}

// WithOSAverageSeverity adds an average severity filter condition for OS assets.
func WithOSAverageSeverity(severity string) Option {
	return WithField(OSColumnAverageSeverity, severity)
}

// Host Asset filtering options (when type="host")

// WithHostSeverity adds a severity filter condition for host assets.
func WithHostSeverity(severity string) Option {
	return WithField(HostColumnSeverity, severity)
}

// WithHostOS adds an OS filter condition for host assets.
func WithHostOS(os string) Option {
	return WithField(HostColumnOS, os)
}

// WithHostOSLike adds an OS filter condition using LIKE operator for host assets.
func WithHostOSLike(pattern string) Option {
	return WithFieldLike(HostColumnOS, pattern)
}

// WithHostOSS adds an OSS filter condition for host assets.
func WithHostOSS(oss string) Option {
	return WithField(HostColumnOSS, oss)
}

// WithHostOSSLike adds an OSS filter condition using LIKE operator for host assets.
func WithHostOSSLike(pattern string) Option {
	return WithFieldLike(HostColumnOSS, pattern)
}

// WithHostHostname adds a hostname filter condition for host assets.
func WithHostHostname(hostname string) Option {
	return WithField(HostColumnHostname, hostname)
}

// WithHostHostnameLike adds a hostname filter condition using LIKE operator for host assets.
func WithHostHostnameLike(pattern string) Option {
	return WithFieldLike(HostColumnHostname, pattern)
}

// WithHostIP adds an IP address filter condition for host assets.
func WithHostIP(ip string) Option {
	return WithField(HostColumnIP, ip)
}

// WithHostIPLike adds an IP address filter condition using LIKE operator for host assets.
func WithHostIPLike(pattern string) Option {
	return WithFieldLike(HostColumnIP, pattern)
}

// WithAssetHostCount adds a host_count filter condition.
func WithAssetHostCount(count int) Option {
	return WithNumericField("host_count", count)
}

// WithAssetHostCountGreaterThan adds a host_count filter condition.
func WithAssetHostCountGreaterThan(count int) Option {
	return WithNumericFieldGreaterThan("host_count", count)
}

// WithAssetHostCountGreaterThanOrEqual adds a host_count filter condition.
func WithAssetHostCountGreaterThanOrEqual(count int) Option {
	return WithNumericFieldGreaterThanOrEqual("host_count", count)
}

// WithAssetHostCountLessThan adds a host_count filter condition.
func WithAssetHostCountLessThan(count int) Option {
	return WithNumericFieldLessThan("host_count", count)
}

// WithAssetHostCountLessThanOrEqual adds a host_count filter condition.
func WithAssetHostCountLessThanOrEqual(count int) Option {
	return WithNumericFieldLessThanOrEqual("host_count", count)
}

// WithAssetHosts adds a hosts filter condition.
func WithAssetHosts(hosts int) Option {
	return WithNumericField("hosts", hosts)
}

// WithAssetHostsGreaterThan adds a hosts filter condition.
func WithAssetHostsGreaterThan(hosts int) Option {
	return WithNumericFieldGreaterThan("hosts", hosts)
}

// WithAssetHostsGreaterThanOrEqual adds a hosts filter condition.
func WithAssetHostsGreaterThanOrEqual(hosts int) Option {
	return WithNumericFieldGreaterThanOrEqual("hosts", hosts)
}

// WithAssetHostsLessThan adds a hosts filter condition.
func WithAssetHostsLessThan(hosts int) Option {
	return WithNumericFieldLessThan("hosts", hosts)
}

// WithAssetHostsLessThanOrEqual adds a hosts filter condition.
func WithAssetHostsLessThanOrEqual(hosts int) Option {
	return WithNumericFieldLessThanOrEqual("hosts", hosts)
}

// WithAssetMaxHosts adds a max_hosts filter condition.
func WithAssetMaxHosts(maxHosts int) Option {
	return WithNumericField("max_hosts", maxHosts)
}

// WithAssetMaxHostsGreaterThan adds a max_hosts filter condition.
func WithAssetMaxHostsGreaterThan(maxHosts int) Option {
	return WithNumericFieldGreaterThan("max_hosts", maxHosts)
}

// WithAssetMaxHostsGreaterThanOrEqual adds a max_hosts filter condition.
func WithAssetMaxHostsGreaterThanOrEqual(maxHosts int) Option {
	return WithNumericFieldGreaterThanOrEqual("max_hosts", maxHosts)
}

// WithAssetMaxHostsLessThan adds a max_hosts filter condition.
func WithAssetMaxHostsLessThan(maxHosts int) Option {
	return WithNumericFieldLessThan("max_hosts", maxHosts)
}

// WithAssetMaxHostsLessThanOrEqual adds a max_hosts filter condition.
func WithAssetMaxHostsLessThanOrEqual(maxHosts int) Option {
	return WithNumericFieldLessThanOrEqual("max_hosts", maxHosts)
}

// WithAssetSeverity adds a severity filter condition.
func WithAssetSeverity(severity int) Option {
	return WithNumericField("severity", severity)
}

// WithAssetSeverityGreaterThan adds a severity filter condition.
func WithAssetSeverityGreaterThan(severity int) Option {
	return WithNumericFieldGreaterThan("severity", severity)
}

// WithAssetSeverityGreaterThanOrEqual adds a severity filter condition.
func WithAssetSeverityGreaterThanOrEqual(severity int) Option {
	return WithNumericFieldGreaterThanOrEqual("severity", severity)
}

// WithAssetSeverityLessThan adds a severity filter condition.
func WithAssetSeverityLessThan(severity int) Option {
	return WithNumericFieldLessThan("severity", severity)
}

// WithAssetSeverityLessThanOrEqual adds a severity filter condition.
func WithAssetSeverityLessThanOrEqual(severity int) Option {
	return WithNumericFieldLessThanOrEqual("severity", severity)
}

// WithAssetHigh adds a high filter condition.
func WithAssetHigh(high int) Option {
	return WithNumericField("high", high)
}

// WithAssetHighGreaterThan adds a high filter condition.
func WithAssetHighGreaterThan(high int) Option {
	return WithNumericFieldGreaterThan("high", high)
}

// WithAssetHighGreaterThanOrEqual adds a high filter condition.
func WithAssetHighGreaterThanOrEqual(high int) Option {
	return WithNumericFieldGreaterThanOrEqual("high", high)
}

// WithAssetHighLessThan adds a high filter condition.
func WithAssetHighLessThan(high int) Option {
	return WithNumericFieldLessThan("high", high)
}

// WithAssetHighLessThanOrEqual adds a high filter condition.
func WithAssetHighLessThanOrEqual(high int) Option {
	return WithNumericFieldLessThanOrEqual("high", high)
}

// WithAssetMedium adds a medium filter condition.
func WithAssetMedium(medium int) Option {
	return WithNumericField("medium", medium)
}

// WithAssetMediumGreaterThan adds a medium filter condition.
func WithAssetMediumGreaterThan(medium int) Option {
	return WithNumericFieldGreaterThan("medium", medium)
}

// WithAssetMediumGreaterThanOrEqual adds a medium filter condition.
func WithAssetMediumGreaterThanOrEqual(medium int) Option {
	return WithNumericFieldGreaterThanOrEqual("medium", medium)
}

// WithAssetMediumLessThan adds a medium filter condition.
func WithAssetMediumLessThan(medium int) Option {
	return WithNumericFieldLessThan("medium", medium)
}

// WithAssetMediumLessThanOrEqual adds a medium filter condition.
func WithAssetMediumLessThanOrEqual(medium int) Option {
	return WithNumericFieldLessThanOrEqual("medium", medium)
}

// WithAssetLow adds a low filter condition.
func WithAssetLow(low int) Option {
	return WithNumericField("low", low)
}

// WithAssetLowGreaterThan adds a low filter condition.
func WithAssetLowGreaterThan(low int) Option {
	return WithNumericFieldGreaterThan("low", low)
}

// WithAssetLowGreaterThanOrEqual adds a low filter condition.
func WithAssetLowGreaterThanOrEqual(low int) Option {
	return WithNumericFieldGreaterThanOrEqual("low", low)
}

// WithAssetLowLessThan adds a low filter condition.
func WithAssetLowLessThan(low int) Option {
	return WithNumericFieldLessThan("low", low)
}

// WithAssetLowLessThanOrEqual adds a low filter condition.
func WithAssetLowLessThanOrEqual(low int) Option {
	return WithNumericFieldLessThanOrEqual("low", low)
}

// WithAssetLog adds a log filter condition.
func WithAssetLog(log int) Option {
	return WithNumericField("log", log)
}

// WithAssetLogGreaterThan adds a log filter condition.
func WithAssetLogGreaterThan(log int) Option {
	return WithNumericFieldGreaterThan("log", log)
}

// WithAssetLogGreaterThanOrEqual adds a log filter condition.
func WithAssetLogGreaterThanOrEqual(log int) Option {
	return WithNumericFieldGreaterThanOrEqual("log", log)
}

// WithAssetLogLessThan adds a log filter condition.
func WithAssetLogLessThan(log int) Option {
	return WithNumericFieldLessThan("log", log)
}

// WithAssetLogLessThanOrEqual adds a log filter condition.
func WithAssetLogLessThanOrEqual(log int) Option {
	return WithNumericFieldLessThanOrEqual("log", log)
}

// WithAssetFalsePositive adds a false_positive filter condition.
func WithAssetFalsePositive(falsePositive int) Option {
	return WithNumericField("false_positive", falsePositive)
}

// WithAssetFalsePositiveGreaterThan adds a false_positive filter condition.
func WithAssetFalsePositiveGreaterThan(falsePositive int) Option {
	return WithNumericFieldGreaterThan("false_positive", falsePositive)
}

// WithAssetFalsePositiveGreaterThanOrEqual adds a false_positive filter condition.
func WithAssetFalsePositiveGreaterThanOrEqual(falsePositive int) Option {
	return WithNumericFieldGreaterThanOrEqual("false_positive", falsePositive)
}

// WithAssetFalsePositiveLessThan adds a false_positive filter condition.
func WithAssetFalsePositiveLessThan(falsePositive int) Option {
	return WithNumericFieldLessThan("false_positive", falsePositive)
}

// WithAssetFalsePositiveLessThanOrEqual adds a false_positive filter condition.
func WithAssetFalsePositiveLessThanOrEqual(falsePositive int) Option {
	return WithNumericFieldLessThanOrEqual("false_positive", falsePositive)
}

// WithAssetComment adds a comment filter condition.
func WithAssetComment(comment string) Option {
	return WithField(AssetColumnComment, comment)
}

// WithAssetCommentLike adds a comment filter condition using LIKE operator.
func WithAssetCommentLike(pattern string) Option {
	return WithFieldLike(AssetColumnComment, pattern)
}

// WithAssetOwner adds an owner filter condition.
func WithAssetOwner(owner string) Option {
	return WithField(AssetColumnOwner, owner)
}

// WithAssetCreated adds a created filter condition.
func WithAssetCreated(created string) Option {
	return WithField(AssetColumnCreated, created)
}

// WithAssetModified adds a modified filter condition.
func WithAssetModified(modified string) Option {
	return WithField(AssetColumnModified, modified)
}

// WithAssetCreatedAfter adds a creation time filter condition.
func WithAssetCreatedAfter(t time.Time) Option {
	return WithTimeFieldAfter(AssetColumnCreated, t)
}

// WithAssetCreatedBefore adds a creation time filter condition.
func WithAssetCreatedBefore(t time.Time) Option {
	return WithTimeFieldBefore(AssetColumnCreated, t)
}

// WithAssetCreatedBetween adds a creation time filter condition.
func WithAssetCreatedBetween(start, end time.Time) Option {
	return WithTimeFieldBetween(AssetColumnCreated, start, end)
}

// WithAssetModifiedAfter adds a modification time filter condition.
func WithAssetModifiedAfter(t time.Time) Option {
	return WithTimeFieldAfter(AssetColumnModified, t)
}

// WithAssetModifiedBefore adds a modification time filter condition.
func WithAssetModifiedBefore(t time.Time) Option {
	return WithTimeFieldBefore(AssetColumnModified, t)
}

// WithAssetModifiedBetween adds a modification time filter condition.
func WithAssetModifiedBetween(start, end time.Time) Option {
	return WithTimeFieldBetween(AssetColumnModified, start, end)
}

// WithAssetTag adds a tag filter condition.
func WithAssetTag(tag string) Option {
	return WithField(AssetColumnTag, tag)
}

// WithAssetTagID adds a tag_id filter condition.
func WithAssetTagID(tagID string) Option {
	return WithField(AssetColumnTagID, tagID)
}

// WithHosts adds a hosts filter condition (alias for WithAssetHosts).
func WithHosts(hosts int) Option {
	return WithNumericField("hosts", hosts)
}

// WithHostsGreaterThan adds a hosts filter condition.
func WithHostsGreaterThan(hosts int) Option {
	return WithNumericFieldGreaterThan("hosts", hosts)
}

// WithHostsGreaterThanOrEqual adds a hosts filter condition.
func WithHostsGreaterThanOrEqual(hosts int) Option {
	return WithNumericFieldGreaterThanOrEqual("hosts", hosts)
}

// WithHostsLessThan adds a hosts filter condition.
func WithHostsLessThan(hosts int) Option {
	return WithNumericFieldLessThan("hosts", hosts)
}

// WithHostsLessThanOrEqual adds a hosts filter condition.
func WithHostsLessThanOrEqual(hosts int) Option {
	return WithNumericFieldLessThanOrEqual("hosts", hosts)
}

// WithLocation adds a location filter condition (alias for WithAssetLocation).
func WithLocation(location string) Option {
	return WithField(AssetColumnLocation, location)
}

// WithLocationLike adds a location filter condition using LIKE operator.
func WithLocationLike(pattern string) Option {
	return WithFieldLike(AssetColumnLocation, pattern)
}

// WithOS adds an OS filter condition (alias for WithAssetOS).
func WithOS(os string) Option {
	return WithField(AssetColumnOS, os)
}

// WithOSLike adds an OS filter condition using LIKE operator.
func WithOSLike(pattern string) Option {
	return WithFieldLike(AssetColumnOS, pattern)
}

// WithHostCount adds a host_count filter condition (alias for WithAssetHostCount).
func WithHostCount(count int) Option {
	return WithNumericField("host_count", count)
}

// WithHostCountGreaterThan adds a host_count filter condition.
func WithHostCountGreaterThan(count int) Option {
	return WithNumericFieldGreaterThan("host_count", count)
}

// WithHostCountGreaterThanOrEqual adds a host_count filter condition.
func WithHostCountGreaterThanOrEqual(count int) Option {
	return WithNumericFieldGreaterThanOrEqual("host_count", count)
}

// WithHostCountLessThan adds a host_count filter condition.
func WithHostCountLessThan(count int) Option {
	return WithNumericFieldLessThan("host_count", count)
}

// WithHostCountLessThanOrEqual adds a host_count filter condition.
func WithHostCountLessThanOrEqual(count int) Option {
	return WithNumericFieldLessThanOrEqual("host_count", count)
}

// WithMaxHosts adds a max_hosts filter condition (alias for WithAssetMaxHosts).
func WithMaxHosts(maxHosts int) Option {
	return WithNumericField("max_hosts", maxHosts)
}

// WithMaxHostsGreaterThan adds a max_hosts filter condition.
func WithMaxHostsGreaterThan(maxHosts int) Option {
	return WithNumericFieldGreaterThan("max_hosts", maxHosts)
}

// WithMaxHostsGreaterThanOrEqual adds a max_hosts filter condition.
func WithMaxHostsGreaterThanOrEqual(maxHosts int) Option {
	return WithNumericFieldGreaterThanOrEqual("max_hosts", maxHosts)
}

// WithMaxHostsLessThan adds a max_hosts filter condition.
func WithMaxHostsLessThan(maxHosts int) Option {
	return WithNumericFieldLessThan("max_hosts", maxHosts)
}

// WithMaxHostsLessThanOrEqual adds a max_hosts filter condition.
func WithMaxHostsLessThanOrEqual(maxHosts int) Option {
	return WithNumericFieldLessThanOrEqual("max_hosts", maxHosts)
}

// WithPorts adds a ports filter condition.
func WithPorts(ports int) Option {
	return WithNumericField("ports", ports)
}

// WithPortsGreaterThan adds a ports filter condition.
func WithPortsGreaterThan(ports int) Option {
	return WithNumericFieldGreaterThan("ports", ports)
}

// WithPortsGreaterThanOrEqual adds a ports filter condition.
func WithPortsGreaterThanOrEqual(ports int) Option {
	return WithNumericFieldGreaterThanOrEqual("ports", ports)
}

// WithPortsLessThan adds a ports filter condition.
func WithPortsLessThan(ports int) Option {
	return WithNumericFieldLessThan("ports", ports)
}

// WithPortsLessThanOrEqual adds a ports filter condition.
func WithPortsLessThanOrEqual(ports int) Option {
	return WithNumericFieldLessThanOrEqual("ports", ports)
}
