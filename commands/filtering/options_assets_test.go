package filtering

import (
	"testing"
	"time"
)

func TestWithAssetID(t *testing.T) {
	opts := []Option{WithAssetID("asset-123")}
	filter := BuildFilter(opts...)
	expected := "uuid=asset-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetName(t *testing.T) {
	opts := []Option{WithAssetName("test-asset")}
	filter := BuildFilter(opts...)
	expected := "name=test-asset"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetNameLike(t *testing.T) {
	opts := []Option{WithAssetNameLike("test%")}
	filter := BuildFilter(opts...)
	expected := "name~test%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetType(t *testing.T) {
	opts := []Option{WithAssetType("host")}
	filter := BuildFilter(opts...)
	expected := "asset_type=host"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetTypeLike(t *testing.T) {
	opts := []Option{WithAssetTypeLike("host%")}
	filter := BuildFilter(opts...)
	expected := "asset_type~host%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHost(t *testing.T) {
	opts := []Option{WithAssetHost("192.168.1.1")}
	filter := BuildFilter(opts...)
	expected := "host=192.168.1.1"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHostLike(t *testing.T) {
	opts := []Option{WithAssetHostLike("192.168.%")}
	filter := BuildFilter(opts...)
	expected := "host~192.168.%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetOS(t *testing.T) {
	opts := []Option{WithAssetOS("Linux")}
	filter := BuildFilter(opts...)
	expected := "os=Linux"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetOSLike(t *testing.T) {
	opts := []Option{WithAssetOSLike("Linux%")}
	filter := BuildFilter(opts...)
	expected := "os~Linux%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetLocation(t *testing.T) {
	opts := []Option{WithAssetLocation("datacenter-1")}
	filter := BuildFilter(opts...)
	expected := "location=datacenter-1"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetLocationLike(t *testing.T) {
	opts := []Option{WithAssetLocationLike("datacenter%")}
	filter := BuildFilter(opts...)
	expected := "location~datacenter%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHostCount(t *testing.T) {
	opts := []Option{WithAssetHostCount(10)}
	filter := BuildFilter(opts...)
	expected := "host_count=10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHostCountGreaterThan(t *testing.T) {
	opts := []Option{WithAssetHostCountGreaterThan(5)}
	filter := BuildFilter(opts...)
	expected := "host_count>5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHostCountGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetHostCountGreaterThanOrEqual(5)}
	filter := BuildFilter(opts...)
	expected := "host_count>=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHostCountLessThan(t *testing.T) {
	opts := []Option{WithAssetHostCountLessThan(20)}
	filter := BuildFilter(opts...)
	expected := "host_count<20"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHostCountLessThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetHostCountLessThanOrEqual(20)}
	filter := BuildFilter(opts...)
	expected := "host_count<=20"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHosts(t *testing.T) {
	opts := []Option{WithAssetHosts(15)}
	filter := BuildFilter(opts...)
	expected := "hosts=15"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHostsGreaterThan(t *testing.T) {
	opts := []Option{WithAssetHostsGreaterThan(10)}
	filter := BuildFilter(opts...)
	expected := "hosts>10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHostsGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetHostsGreaterThanOrEqual(10)}
	filter := BuildFilter(opts...)
	expected := "hosts>=10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHostsLessThan(t *testing.T) {
	opts := []Option{WithAssetHostsLessThan(25)}
	filter := BuildFilter(opts...)
	expected := "hosts<25"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHostsLessThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetHostsLessThanOrEqual(25)}
	filter := BuildFilter(opts...)
	expected := "hosts<=25"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetMaxHosts(t *testing.T) {
	opts := []Option{WithAssetMaxHosts(100)}
	filter := BuildFilter(opts...)
	expected := "max_hosts=100"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetMaxHostsGreaterThan(t *testing.T) {
	opts := []Option{WithAssetMaxHostsGreaterThan(50)}
	filter := BuildFilter(opts...)
	expected := "max_hosts>50"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetMaxHostsGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetMaxHostsGreaterThanOrEqual(50)}
	filter := BuildFilter(opts...)
	expected := "max_hosts>=50"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetMaxHostsLessThan(t *testing.T) {
	opts := []Option{WithAssetMaxHostsLessThan(200)}
	filter := BuildFilter(opts...)
	expected := "max_hosts<200"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetMaxHostsLessThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetMaxHostsLessThanOrEqual(200)}
	filter := BuildFilter(opts...)
	expected := "max_hosts<=200"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHosts(t *testing.T) {
	opts := []Option{WithHosts(30)}
	filter := BuildFilter(opts...)
	expected := "hosts=30"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostsGreaterThan(t *testing.T) {
	opts := []Option{WithHostsGreaterThan(20)}
	filter := BuildFilter(opts...)
	expected := "hosts>20"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetSeverity(t *testing.T) {
	opts := []Option{WithAssetSeverity(5)}
	filter := BuildFilter(opts...)
	expected := "severity=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetSeverityGreaterThan(t *testing.T) {
	opts := []Option{WithAssetSeverityGreaterThan(3)}
	filter := BuildFilter(opts...)
	expected := "severity>3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetSeverityGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetSeverityGreaterThanOrEqual(3)}
	filter := BuildFilter(opts...)
	expected := "severity>=3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetSeverityLessThan(t *testing.T) {
	opts := []Option{WithAssetSeverityLessThan(7)}
	filter := BuildFilter(opts...)
	expected := "severity<7"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetSeverityLessThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetSeverityLessThanOrEqual(7)}
	filter := BuildFilter(opts...)
	expected := "severity<=7"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHigh(t *testing.T) {
	opts := []Option{WithAssetHigh(10)}
	filter := BuildFilter(opts...)
	expected := "high=10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHighGreaterThan(t *testing.T) {
	opts := []Option{WithAssetHighGreaterThan(5)}
	filter := BuildFilter(opts...)
	expected := "high>5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHighGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetHighGreaterThanOrEqual(5)}
	filter := BuildFilter(opts...)
	expected := "high>=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHighLessThan(t *testing.T) {
	opts := []Option{WithAssetHighLessThan(15)}
	filter := BuildFilter(opts...)
	expected := "high<15"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetHighLessThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetHighLessThanOrEqual(15)}
	filter := BuildFilter(opts...)
	expected := "high<=15"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetMedium(t *testing.T) {
	opts := []Option{WithAssetMedium(20)}
	filter := BuildFilter(opts...)
	expected := "medium=20"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetMediumGreaterThan(t *testing.T) {
	opts := []Option{WithAssetMediumGreaterThan(10)}
	filter := BuildFilter(opts...)
	expected := "medium>10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetMediumGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetMediumGreaterThanOrEqual(10)}
	filter := BuildFilter(opts...)
	expected := "medium>=10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetMediumLessThan(t *testing.T) {
	opts := []Option{WithAssetMediumLessThan(30)}
	filter := BuildFilter(opts...)
	expected := "medium<30"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetMediumLessThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetMediumLessThanOrEqual(30)}
	filter := BuildFilter(opts...)
	expected := "medium<=30"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetLow(t *testing.T) {
	opts := []Option{WithAssetLow(30)}
	filter := BuildFilter(opts...)
	expected := "low=30"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetLowGreaterThan(t *testing.T) {
	opts := []Option{WithAssetLowGreaterThan(15)}
	filter := BuildFilter(opts...)
	expected := "low>15"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetLowGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetLowGreaterThanOrEqual(15)}
	filter := BuildFilter(opts...)
	expected := "low>=15"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetLowLessThan(t *testing.T) {
	opts := []Option{WithAssetLowLessThan(45)}
	filter := BuildFilter(opts...)
	expected := "low<45"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetLowLessThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetLowLessThanOrEqual(45)}
	filter := BuildFilter(opts...)
	expected := "low<=45"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetLog(t *testing.T) {
	opts := []Option{WithAssetLog(5)}
	filter := BuildFilter(opts...)
	expected := "log=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetLogGreaterThan(t *testing.T) {
	opts := []Option{WithAssetLogGreaterThan(2)}
	filter := BuildFilter(opts...)
	expected := "log>2"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetLogGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetLogGreaterThanOrEqual(2)}
	filter := BuildFilter(opts...)
	expected := "log>=2"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetLogLessThan(t *testing.T) {
	opts := []Option{WithAssetLogLessThan(8)}
	filter := BuildFilter(opts...)
	expected := "log<8"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetLogLessThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetLogLessThanOrEqual(8)}
	filter := BuildFilter(opts...)
	expected := "log<=8"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetFalsePositive(t *testing.T) {
	opts := []Option{WithAssetFalsePositive(3)}
	filter := BuildFilter(opts...)
	expected := "false_positive=3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetFalsePositiveGreaterThan(t *testing.T) {
	opts := []Option{WithAssetFalsePositiveGreaterThan(1)}
	filter := BuildFilter(opts...)
	expected := "false_positive>1"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetFalsePositiveGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetFalsePositiveGreaterThanOrEqual(1)}
	filter := BuildFilter(opts...)
	expected := "false_positive>=1"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetFalsePositiveLessThan(t *testing.T) {
	opts := []Option{WithAssetFalsePositiveLessThan(5)}
	filter := BuildFilter(opts...)
	expected := "false_positive<5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetFalsePositiveLessThanOrEqual(t *testing.T) {
	opts := []Option{WithAssetFalsePositiveLessThanOrEqual(5)}
	filter := BuildFilter(opts...)
	expected := "false_positive<=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetComment(t *testing.T) {
	opts := []Option{WithAssetComment("test comment")}
	filter := BuildFilter(opts...)
	expected := "comment=\"test comment\""
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetCommentLike(t *testing.T) {
	opts := []Option{WithAssetCommentLike("test%")}
	filter := BuildFilter(opts...)
	expected := "comment~test%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetOwner(t *testing.T) {
	opts := []Option{WithAssetOwner("admin")}
	filter := BuildFilter(opts...)
	expected := "owner=admin"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetCreated(t *testing.T) {
	opts := []Option{WithAssetCreated("2023-01-01")}
	filter := BuildFilter(opts...)
	expected := "created=2023-01-01"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetModified(t *testing.T) {
	opts := []Option{WithAssetModified("2023-01-02")}
	filter := BuildFilter(opts...)
	expected := "modified=2023-01-02"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetTag(t *testing.T) {
	opts := []Option{WithAssetTag("production")}
	filter := BuildFilter(opts...)
	expected := "tag=production"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetTagID(t *testing.T) {
	opts := []Option{WithAssetTagID("tag-123")}
	filter := BuildFilter(opts...)
	expected := "tag_id=tag-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

// OS Asset filtering tests

func TestWithOSTitle(t *testing.T) {
	opts := []Option{WithOSTitle("Windows 10")}
	filter := BuildFilter(opts...)
	expected := "title=\"Windows 10\""
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithOSTitleLike(t *testing.T) {
	opts := []Option{WithOSTitleLike("Windows%")}
	filter := BuildFilter(opts...)
	expected := "title~Windows%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithOSHosts(t *testing.T) {
	opts := []Option{WithOSHosts(5)}
	filter := BuildFilter(opts...)
	expected := "hosts=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithOSHostsGreaterThan(t *testing.T) {
	opts := []Option{WithOSHostsGreaterThan(3)}
	filter := BuildFilter(opts...)
	expected := "hosts>3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithOSHostsGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithOSHostsGreaterThanOrEqual(2)}
	filter := BuildFilter(opts...)
	expected := "hosts>=2"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithOSHostsLessThan(t *testing.T) {
	opts := []Option{WithOSHostsLessThan(10)}
	filter := BuildFilter(opts...)
	expected := "hosts<10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithOSHostsLessThanOrEqual(t *testing.T) {
	opts := []Option{WithOSHostsLessThanOrEqual(8)}
	filter := BuildFilter(opts...)
	expected := "hosts<=8"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithOSLatestSeverity(t *testing.T) {
	opts := []Option{WithOSLatestSeverity("high")}
	filter := BuildFilter(opts...)
	expected := "latest_severity=high"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithOSHighestSeverity(t *testing.T) {
	opts := []Option{WithOSHighestSeverity("critical")}
	filter := BuildFilter(opts...)
	expected := "highest_severity=critical"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithOSAverageSeverity(t *testing.T) {
	opts := []Option{WithOSAverageSeverity("medium")}
	filter := BuildFilter(opts...)
	expected := "average_severity=medium"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

// Host Asset filtering tests

func TestWithHostSeverity(t *testing.T) {
	opts := []Option{WithHostSeverity("high")}
	filter := BuildFilter(opts...)
	expected := "severity=high"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostOS(t *testing.T) {
	opts := []Option{WithHostOS("Linux")}
	filter := BuildFilter(opts...)
	expected := "os=Linux"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostOSLike(t *testing.T) {
	opts := []Option{WithHostOSLike("Linux%")}
	filter := BuildFilter(opts...)
	expected := "os~Linux%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostOSS(t *testing.T) {
	opts := []Option{WithHostOSS("Ubuntu 20.04")}
	filter := BuildFilter(opts...)
	expected := "oss=\"Ubuntu 20.04\""
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostOSSLike(t *testing.T) {
	opts := []Option{WithHostOSSLike("Ubuntu%")}
	filter := BuildFilter(opts...)
	expected := "oss~Ubuntu%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostHostname(t *testing.T) {
	opts := []Option{WithHostHostname("server01")}
	filter := BuildFilter(opts...)
	expected := "hostname=server01"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostHostnameLike(t *testing.T) {
	opts := []Option{WithHostHostnameLike("server%")}
	filter := BuildFilter(opts...)
	expected := "hostname~server%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostIP(t *testing.T) {
	opts := []Option{WithHostIP("192.168.1.100")}
	filter := BuildFilter(opts...)
	expected := "ip=192.168.1.100"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostIPLike(t *testing.T) {
	opts := []Option{WithHostIPLike("192.168.%")}
	filter := BuildFilter(opts...)
	expected := "ip~192.168.%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

// Time-based Asset filtering tests

func TestWithAssetCreatedAfter(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	opts := []Option{WithAssetCreatedAfter(testTime)}
	filter := BuildFilter(opts...)
	expected := "created>2023-01-01T00:00:00Z"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetCreatedBefore(t *testing.T) {
	testTime := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
	opts := []Option{WithAssetCreatedBefore(testTime)}
	filter := BuildFilter(opts...)
	expected := "created<2023-12-31T23:59:59Z"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetCreatedBetween(t *testing.T) {
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
	opts := []Option{WithAssetCreatedBetween(startTime, endTime)}
	filter := BuildFilter(opts...)
	expected := "created>=2023-01-01T00:00:00Z created<=2023-12-31T23:59:59Z"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetModifiedAfter(t *testing.T) {
	testTime := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	opts := []Option{WithAssetModifiedAfter(testTime)}
	filter := BuildFilter(opts...)
	expected := "modified>2023-06-01T00:00:00Z"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetModifiedBefore(t *testing.T) {
	testTime := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
	opts := []Option{WithAssetModifiedBefore(testTime)}
	filter := BuildFilter(opts...)
	expected := "modified<2023-12-31T23:59:59Z"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetModifiedBetween(t *testing.T) {
	startTime := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
	opts := []Option{WithAssetModifiedBetween(startTime, endTime)}
	filter := BuildFilter(opts...)
	expected := "modified>=2023-06-01T00:00:00Z modified<=2023-12-31T23:59:59Z"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}
