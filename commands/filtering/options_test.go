package filtering

import (
	"testing"
	"time"
)

func TestWithNameLike(t *testing.T) {
	opts := []Option{WithNameLike("test%")}
	filter := BuildFilter(opts...)
	expected := "name~test%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithComment(t *testing.T) {
	opts := []Option{WithComment("test comment")}
	filter := BuildFilter(opts...)
	expected := "comment=\"test comment\""
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithCommentLike(t *testing.T) {
	opts := []Option{WithCommentLike("test%")}
	filter := BuildFilter(opts...)
	expected := "comment~test%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithOwner(t *testing.T) {
	opts := []Option{WithOwner("admin")}
	filter := BuildFilter(opts...)
	expected := "owner=admin"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithStatusIn(t *testing.T) {
	opts := []Option{WithStatusIn("Running", "Done")}
	filter := BuildFilter(opts...)
	expected := "status=Running status=Done"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithSeverityGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithSeverityGreaterThanOrEqual(5.0)}
	filter := BuildFilter(opts...)
	expected := "severity>=5.0"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithSeverityLessThanOrEqual(t *testing.T) {
	opts := []Option{WithSeverityLessThanOrEqual(8.0)}
	filter := BuildFilter(opts...)
	expected := "severity<=8.0"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithThreat(t *testing.T) {
	opts := []Option{WithThreat("High")}
	filter := BuildFilter(opts...)
	expected := "threat=High"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTrend(t *testing.T) {
	opts := []Option{WithTrend("up")}
	filter := BuildFilter(opts...)
	expected := "trend=up"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithCreatedBetween(t *testing.T) {
	now := time.Now()
	opts := []Option{WithCreatedBetween(now, now.Add(24*time.Hour))}
	filter := BuildFilter(opts...)
	if !containsHelper(filter, "created>") || !containsHelper(filter, "created<") {
		t.Errorf("Expected filter to contain created> and created<, got %q", filter)
	}
}

func TestWithModifiedAfter(t *testing.T) {
	now := time.Now()
	opts := []Option{WithModifiedAfter(now)}
	filter := BuildFilter(opts...)
	if !containsHelper(filter, "modified>") {
		t.Errorf("Expected filter to contain modified>, got %q", filter)
	}
}

func TestWithModifiedBefore(t *testing.T) {
	now := time.Now()
	opts := []Option{WithModifiedBefore(now)}
	filter := BuildFilter(opts...)
	if !containsHelper(filter, "modified<") {
		t.Errorf("Expected filter to contain modified<, got %q", filter)
	}
}

func TestWithModifiedBetween(t *testing.T) {
	now := time.Now()
	opts := []Option{WithModifiedBetween(now, now.Add(24*time.Hour))}
	filter := BuildFilter(opts...)
	if !containsHelper(filter, "modified>") || !containsHelper(filter, "modified<") {
		t.Errorf("Expected filter to contain modified> and modified<, got %q", filter)
	}
}

func TestWithSort(t *testing.T) {
	opts := []Option{WithSort("name", "asc")}
	filter := BuildFilter(opts...)
	expected := "sort=name"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithSortDesc(t *testing.T) {
	opts := []Option{WithSortDesc(SortByName)}
	filter := BuildFilter(opts...)
	expected := "sort-reverse=name"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithOffset(t *testing.T) {
	opts := []Option{WithOffset(10)}
	filter := BuildFilter(opts...)
	expected := "first=10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTag(t *testing.T) {
	opts := []Option{WithTag("test-tag")}
	filter := BuildFilter(opts...)
	expected := "tag=test-tag"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTagID(t *testing.T) {
	opts := []Option{WithTagID("tag-uuid")}
	filter := BuildFilter(opts...)
	expected := "tag_id=tag-uuid"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithUsageType(t *testing.T) {
	opts := []Option{WithUsageType("scan")}
	filter := BuildFilter(opts...)
	expected := "usage_type=scan"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTrash(t *testing.T) {
	opts := []Option{WithTrash(true)}
	filter := BuildFilter(opts...)
	expected := "trash=true"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithApplyOverrides(t *testing.T) {
	opts := []Option{WithApplyOverrides(true)}
	filter := BuildFilter(opts...)
	expected := "apply_overrides=true"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithMinQoD(t *testing.T) {
	opts := []Option{WithMinQoD(70)}
	filter := BuildFilter(opts...)
	expected := "min_qod>=70"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHost(t *testing.T) {
	opts := []Option{WithHost("192.168.1.1")}
	filter := BuildFilter(opts...)
	expected := "host=192.168.1.1"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostLike(t *testing.T) {
	opts := []Option{WithHostLike("192.168.%")}
	filter := BuildFilter(opts...)
	expected := "host~192.168.%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithPort(t *testing.T) {
	opts := []Option{WithPort("80")}
	filter := BuildFilter(opts...)
	expected := "port=80"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithPortLike(t *testing.T) {
	opts := []Option{WithPortLike("80%")}
	filter := BuildFilter(opts...)
	expected := "port~80%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithNVT(t *testing.T) {
	opts := []Option{WithNVT("1.3.6.1.4.1.25623.1.0.10330")}
	filter := BuildFilter(opts...)
	expected := "nvt=1.3.6.1.4.1.25623.1.0.10330"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithNVTLike(t *testing.T) {
	opts := []Option{WithNVTLike("1.3.6.1.4.1.25623.1.0.%")}
	filter := BuildFilter(opts...)
	expected := "nvt~1.3.6.1.4.1.25623.1.0.%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTask(t *testing.T) {
	opts := []Option{WithTask("task-uuid")}
	filter := BuildFilter(opts...)
	expected := "task=task-uuid"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskLike(t *testing.T) {
	opts := []Option{WithTaskLike("task-%")}
	filter := BuildFilter(opts...)
	expected := "task~task-%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReport(t *testing.T) {
	opts := []Option{WithReport("report-uuid")}
	filter := BuildFilter(opts...)
	expected := "report=report-uuid"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportLike(t *testing.T) {
	opts := []Option{WithReportLike("report-%")}
	filter := BuildFilter(opts...)
	expected := "report~report-%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithLocation(t *testing.T) {
	opts := []Option{WithLocation("location-uuid")}
	filter := BuildFilter(opts...)
	expected := "location=location-uuid"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithLocationLike(t *testing.T) {
	opts := []Option{WithLocationLike("location-%")}
	filter := BuildFilter(opts...)
	expected := "location~location-%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithDescription(t *testing.T) {
	opts := []Option{WithDescription("test description")}
	filter := BuildFilter(opts...)
	expected := "description=\"test description\""
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithDescriptionLike(t *testing.T) {
	opts := []Option{WithDescriptionLike("test%")}
	filter := BuildFilter(opts...)
	expected := "description~test%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithQoD(t *testing.T) {
	opts := []Option{WithQoD(80)}
	filter := BuildFilter(opts...)
	expected := "qod=80.0"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithQoDGreaterThan(t *testing.T) {
	opts := []Option{WithQoDGreaterThan(70)}
	filter := BuildFilter(opts...)
	expected := "qod>70.0"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithQoDLessThan(t *testing.T) {
	opts := []Option{WithQoDLessThan(90)}
	filter := BuildFilter(opts...)
	expected := "qod<90.0"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithQoDLessThanOrEqual(t *testing.T) {
	opts := []Option{WithQoDLessThanOrEqual(85)}
	filter := BuildFilter(opts...)
	expected := "qod<=85.0"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostsGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithHostsGreaterThanOrEqual(5)}
	filter := BuildFilter(opts...)
	expected := "hosts>=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostsLessThan(t *testing.T) {
	opts := []Option{WithHostsLessThan(20)}
	filter := BuildFilter(opts...)
	expected := "hosts<20"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostsLessThanOrEqual(t *testing.T) {
	opts := []Option{WithHostsLessThanOrEqual(15)}
	filter := BuildFilter(opts...)
	expected := "hosts<=15"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithOS(t *testing.T) {
	opts := []Option{WithOS("Linux")}
	filter := BuildFilter(opts...)
	expected := "os=Linux"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithOSLike(t *testing.T) {
	opts := []Option{WithOSLike("Linux%")}
	filter := BuildFilter(opts...)
	expected := "os~Linux%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostCount(t *testing.T) {
	opts := []Option{WithHostCount(10)}
	filter := BuildFilter(opts...)
	expected := "host_count=10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostCountGreaterThan(t *testing.T) {
	opts := []Option{WithHostCountGreaterThan(5)}
	filter := BuildFilter(opts...)
	expected := "host_count>5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostCountGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithHostCountGreaterThanOrEqual(8)}
	filter := BuildFilter(opts...)
	expected := "host_count>=8"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostCountLessThan(t *testing.T) {
	opts := []Option{WithHostCountLessThan(20)}
	filter := BuildFilter(opts...)
	expected := "host_count<20"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHostCountLessThanOrEqual(t *testing.T) {
	opts := []Option{WithHostCountLessThanOrEqual(15)}
	filter := BuildFilter(opts...)
	expected := "host_count<=15"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithMaxHosts(t *testing.T) {
	opts := []Option{WithMaxHosts(100)}
	filter := BuildFilter(opts...)
	expected := "max_hosts=100"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithMaxHostsGreaterThan(t *testing.T) {
	opts := []Option{WithMaxHostsGreaterThan(50)}
	filter := BuildFilter(opts...)
	expected := "max_hosts>50"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithMaxHostsGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithMaxHostsGreaterThanOrEqual(75)}
	filter := BuildFilter(opts...)
	expected := "max_hosts>=75"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithMaxHostsLessThan(t *testing.T) {
	opts := []Option{WithMaxHostsLessThan(200)}
	filter := BuildFilter(opts...)
	expected := "max_hosts<200"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithMaxHostsLessThanOrEqual(t *testing.T) {
	opts := []Option{WithMaxHostsLessThanOrEqual(150)}
	filter := BuildFilter(opts...)
	expected := "max_hosts<=150"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTasks(t *testing.T) {
	opts := []Option{WithTasks(5)}
	filter := BuildFilter(opts...)
	expected := "tasks=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTasksGreaterThan(t *testing.T) {
	opts := []Option{WithTasksGreaterThan(3)}
	filter := BuildFilter(opts...)
	expected := "tasks>3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTasksGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithTasksGreaterThanOrEqual(4)}
	filter := BuildFilter(opts...)
	expected := "tasks>=4"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTasksLessThan(t *testing.T) {
	opts := []Option{WithTasksLessThan(10)}
	filter := BuildFilter(opts...)
	expected := "tasks<10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTasksLessThanOrEqual(t *testing.T) {
	opts := []Option{WithTasksLessThanOrEqual(8)}
	filter := BuildFilter(opts...)
	expected := "tasks<=8"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithPorts(t *testing.T) {
	opts := []Option{WithPorts(25)}
	filter := BuildFilter(opts...)
	expected := "ports=25"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithPortsGreaterThan(t *testing.T) {
	opts := []Option{WithPortsGreaterThan(20)}
	filter := BuildFilter(opts...)
	expected := "ports>20"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithPortsGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithPortsGreaterThanOrEqual(22)}
	filter := BuildFilter(opts...)
	expected := "ports>=22"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithPortsLessThan(t *testing.T) {
	opts := []Option{WithPortsLessThan(30)}
	filter := BuildFilter(opts...)
	expected := "ports<30"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithPortsLessThanOrEqual(t *testing.T) {
	opts := []Option{WithPortsLessThanOrEqual(28)}
	filter := BuildFilter(opts...)
	expected := "ports<=28"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithValue(t *testing.T) {
	opts := []Option{WithValue("test-value")}
	filter := BuildFilter(opts...)
	expected := "value=test-value"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithValueLike(t *testing.T) {
	opts := []Option{WithValueLike("test%")}
	filter := BuildFilter(opts...)
	expected := "value~test%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithType(t *testing.T) {
	opts := []Option{WithType("scan")}
	filter := BuildFilter(opts...)
	expected := "type=scan"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTypeLike(t *testing.T) {
	opts := []Option{WithTypeLike("scan%")}
	filter := BuildFilter(opts...)
	expected := "type~scan%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithScanner(t *testing.T) {
	opts := []Option{WithScanner("scanner-uuid")}
	filter := BuildFilter(opts...)
	expected := "scanner=scanner-uuid"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithScannerLike(t *testing.T) {
	opts := []Option{WithScannerLike("scanner-%")}
	filter := BuildFilter(opts...)
	expected := "scanner~scanner-%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssignedTo(t *testing.T) {
	opts := []Option{WithAssignedTo("user-uuid")}
	filter := BuildFilter(opts...)
	expected := "assigned_to=user-uuid"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithAssignedBy(t *testing.T) {
	opts := []Option{WithAssignedBy("user-uuid")}
	filter := BuildFilter(opts...)
	expected := "assigned_by=user-uuid"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithResult(t *testing.T) {
	opts := []Option{WithResult("result-uuid")}
	filter := BuildFilter(opts...)
	expected := "result=result-uuid"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithResultLike(t *testing.T) {
	opts := []Option{WithResultLike("result-%")}
	filter := BuildFilter(opts...)
	expected := "result~result-%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTotal(t *testing.T) {
	opts := []Option{WithTotal(100)}
	filter := BuildFilter(opts...)
	expected := "total=100"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTotalGreaterThan(t *testing.T) {
	opts := []Option{WithTotalGreaterThan(50)}
	filter := BuildFilter(opts...)
	expected := "total>50"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTotalGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithTotalGreaterThanOrEqual(75)}
	filter := BuildFilter(opts...)
	expected := "total>=75"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTotalLessThan(t *testing.T) {
	opts := []Option{WithTotalLessThan(200)}
	filter := BuildFilter(opts...)
	expected := "total<200"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTotalLessThanOrEqual(t *testing.T) {
	opts := []Option{WithTotalLessThanOrEqual(150)}
	filter := BuildFilter(opts...)
	expected := "total<=150"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithFalsePositive(t *testing.T) {
	opts := []Option{WithFalsePositive(true)}
	filter := BuildFilter(opts...)
	expected := "false_positive=true"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithFalsePositiveGreaterThan(t *testing.T) {
	opts := []Option{WithFalsePositiveGreaterThan(3)}
	filter := BuildFilter(opts...)
	expected := "false_positive>3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithFalsePositiveGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithFalsePositiveGreaterThanOrEqual(4)}
	filter := BuildFilter(opts...)
	expected := "false_positive>=4"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithFalsePositiveLessThan(t *testing.T) {
	opts := []Option{WithFalsePositiveLessThan(10)}
	filter := BuildFilter(opts...)
	expected := "false_positive<10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithFalsePositiveLessThanOrEqual(t *testing.T) {
	opts := []Option{WithFalsePositiveLessThanOrEqual(8)}
	filter := BuildFilter(opts...)
	expected := "false_positive<=8"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithLog(t *testing.T) {
	opts := []Option{WithLog(5)}
	filter := BuildFilter(opts...)
	expected := "log=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithLogGreaterThan(t *testing.T) {
	opts := []Option{WithLogGreaterThan(3)}
	filter := BuildFilter(opts...)
	expected := "log>3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithLogGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithLogGreaterThanOrEqual(4)}
	filter := BuildFilter(opts...)
	expected := "log>=4"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithLogLessThan(t *testing.T) {
	opts := []Option{WithLogLessThan(10)}
	filter := BuildFilter(opts...)
	expected := "log<10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithLogLessThanOrEqual(t *testing.T) {
	opts := []Option{WithLogLessThanOrEqual(8)}
	filter := BuildFilter(opts...)
	expected := "log<=8"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithLow(t *testing.T) {
	opts := []Option{WithLow(5)}
	filter := BuildFilter(opts...)
	expected := "low=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithLowGreaterThan(t *testing.T) {
	opts := []Option{WithLowGreaterThan(3)}
	filter := BuildFilter(opts...)
	expected := "low>3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithLowGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithLowGreaterThanOrEqual(4)}
	filter := BuildFilter(opts...)
	expected := "low>=4"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithLowLessThan(t *testing.T) {
	opts := []Option{WithLowLessThan(10)}
	filter := BuildFilter(opts...)
	expected := "low<10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithLowLessThanOrEqual(t *testing.T) {
	opts := []Option{WithLowLessThanOrEqual(8)}
	filter := BuildFilter(opts...)
	expected := "low<=8"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithMedium(t *testing.T) {
	opts := []Option{WithMedium(5)}
	filter := BuildFilter(opts...)
	expected := "medium=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithMediumGreaterThan(t *testing.T) {
	opts := []Option{WithMediumGreaterThan(3)}
	filter := BuildFilter(opts...)
	expected := "medium>3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithMediumGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithMediumGreaterThanOrEqual(4)}
	filter := BuildFilter(opts...)
	expected := "medium>=4"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithMediumLessThan(t *testing.T) {
	opts := []Option{WithMediumLessThan(10)}
	filter := BuildFilter(opts...)
	expected := "medium<10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithMediumLessThanOrEqual(t *testing.T) {
	opts := []Option{WithMediumLessThanOrEqual(8)}
	filter := BuildFilter(opts...)
	expected := "medium<=8"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHigh(t *testing.T) {
	opts := []Option{WithHigh(5)}
	filter := BuildFilter(opts...)
	expected := "high=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHighGreaterThan(t *testing.T) {
	opts := []Option{WithHighGreaterThan(3)}
	filter := BuildFilter(opts...)
	expected := "high>3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHighGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithHighGreaterThanOrEqual(4)}
	filter := BuildFilter(opts...)
	expected := "high>=4"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHighLessThan(t *testing.T) {
	opts := []Option{WithHighLessThan(10)}
	filter := BuildFilter(opts...)
	expected := "high<10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithHighLessThanOrEqual(t *testing.T) {
	opts := []Option{WithHighLessThanOrEqual(8)}
	filter := BuildFilter(opts...)
	expected := "high<=8"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithCritical(t *testing.T) {
	opts := []Option{WithCritical(5)}
	filter := BuildFilter(opts...)
	expected := "critical=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithCriticalGreaterThan(t *testing.T) {
	opts := []Option{WithCriticalGreaterThan(3)}
	filter := BuildFilter(opts...)
	expected := "critical>3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithCriticalGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithCriticalGreaterThanOrEqual(4)}
	filter := BuildFilter(opts...)
	expected := "critical>=4"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithCriticalLessThan(t *testing.T) {
	opts := []Option{WithCriticalLessThan(10)}
	filter := BuildFilter(opts...)
	expected := "critical<10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithCriticalLessThanOrEqual(t *testing.T) {
	opts := []Option{WithCriticalLessThanOrEqual(8)}
	filter := BuildFilter(opts...)
	expected := "critical<=8"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithResultHosts(t *testing.T) {
	opts := []Option{WithResultHosts(5)}
	filter := BuildFilter(opts...)
	expected := "result_hosts=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithResultHostsGreaterThan(t *testing.T) {
	opts := []Option{WithResultHostsGreaterThan(3)}
	filter := BuildFilter(opts...)
	expected := "result_hosts>3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithResultHostsGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithResultHostsGreaterThanOrEqual(4)}
	filter := BuildFilter(opts...)
	expected := "result_hosts>=4"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithResultHostsLessThan(t *testing.T) {
	opts := []Option{WithResultHostsLessThan(10)}
	filter := BuildFilter(opts...)
	expected := "result_hosts<10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithResultHostsLessThanOrEqual(t *testing.T) {
	opts := []Option{WithResultHostsLessThanOrEqual(8)}
	filter := BuildFilter(opts...)
	expected := "result_hosts<=8"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithSchedule(t *testing.T) {
	opts := []Option{WithSchedule("schedule-uuid")}
	filter := BuildFilter(opts...)
	expected := "schedule=schedule-uuid"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithScheduleLike(t *testing.T) {
	opts := []Option{WithScheduleLike("schedule-%")}
	filter := BuildFilter(opts...)
	expected := "schedule~schedule-%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithNextDueAfter(t *testing.T) {
	now := time.Now()
	opts := []Option{WithNextDueAfter(now)}
	filter := BuildFilter(opts...)
	if !containsHelper(filter, "next_due>") {
		t.Errorf("Expected filter to contain next_due>, got %q", filter)
	}
}

func TestWithNextDueBefore(t *testing.T) {
	now := time.Now()
	opts := []Option{WithNextDueBefore(now)}
	filter := BuildFilter(opts...)
	if !containsHelper(filter, "next_due<") {
		t.Errorf("Expected filter to contain next_due<, got %q", filter)
	}
}

func TestWithNextDueBetween(t *testing.T) {
	now := time.Now()
	opts := []Option{WithNextDueBetween(now, now.Add(24*time.Hour))}
	filter := BuildFilter(opts...)
	if !containsHelper(filter, "next_due>") || !containsHelper(filter, "next_due<") {
		t.Errorf("Expected filter to contain next_due> and next_due<, got %q", filter)
	}
}

func TestWithFirstReportCreatedAfter(t *testing.T) {
	now := time.Now()
	opts := []Option{WithFirstReportCreatedAfter(now)}
	filter := BuildFilter(opts...)
	if !containsHelper(filter, "first_report_created>") {
		t.Errorf("Expected filter to contain first_report_created>, got %q", filter)
	}
}

func TestWithFirstReportCreatedBefore(t *testing.T) {
	now := time.Now()
	opts := []Option{WithFirstReportCreatedBefore(now)}
	filter := BuildFilter(opts...)
	if !containsHelper(filter, "first_report_created<") {
		t.Errorf("Expected filter to contain first_report_created<, got %q", filter)
	}
}

func TestWithFirstReportCreatedBetween(t *testing.T) {
	now := time.Now()
	opts := []Option{WithFirstReportCreatedBetween(now, now.Add(24*time.Hour))}
	filter := BuildFilter(opts...)
	if !containsHelper(filter, "first_report_created>") || !containsHelper(filter, "first_report_created<") {
		t.Errorf("Expected filter to contain first_report_created> and first_report_created<, got %q", filter)
	}
}

func TestWithLastReportCreatedAfter(t *testing.T) {
	now := time.Now()
	opts := []Option{WithLastReportCreatedAfter(now)}
	filter := BuildFilter(opts...)
	if !containsHelper(filter, "last_report_created>") {
		t.Errorf("Expected filter to contain last_report_created>, got %q", filter)
	}
}

func TestWithLastReportCreatedBefore(t *testing.T) {
	now := time.Now()
	opts := []Option{WithLastReportCreatedBefore(now)}
	filter := BuildFilter(opts...)
	if !containsHelper(filter, "last_report_created<") {
		t.Errorf("Expected filter to contain last_report_created<, got %q", filter)
	}
}

func TestWithLastReportCreatedBetween(t *testing.T) {
	now := time.Now()
	opts := []Option{WithLastReportCreatedBetween(now, now.Add(24*time.Hour))}
	filter := BuildFilter(opts...)
	if !containsHelper(filter, "last_report_created>") || !containsHelper(filter, "last_report_created<") {
		t.Errorf("Expected filter to contain last_report_created> and last_report_created<, got %q", filter)
	}
}
