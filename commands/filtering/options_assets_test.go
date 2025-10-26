package filtering

import (
	"testing"
)

func TestWithAssetID(t *testing.T) {
	opts := []Option{WithAssetID("asset-123")}
	filter := BuildFilter(opts...)
	expected := "asset_id=asset-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetName(t *testing.T) {
	opts := []Option{WithAssetName("test-asset")}
	filter := BuildFilter(opts...)
	expected := "asset_name=test-asset"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssetNameLike(t *testing.T) {
	opts := []Option{WithAssetNameLike("test%")}
	filter := BuildFilter(opts...)
	expected := "asset_name~test%"
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
