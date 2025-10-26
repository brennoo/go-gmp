package filtering

import (
	"testing"
)

func TestWithReportID(t *testing.T) {
	opts := []Option{WithReportID("report-123")}
	filter := BuildFilter(opts...)
	expected := "report_id=report-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportName(t *testing.T) {
	opts := []Option{WithReportName("test-report")}
	filter := BuildFilter(opts...)
	expected := "report_name=test-report"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportNameLike(t *testing.T) {
	opts := []Option{WithReportNameLike("test%")}
	filter := BuildFilter(opts...)
	expected := "report_name~test%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportFormat(t *testing.T) {
	opts := []Option{WithReportFormat("XML")}
	filter := BuildFilter(opts...)
	expected := "report_format=XML"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportFormatLike(t *testing.T) {
	opts := []Option{WithReportFormatLike("XML%")}
	filter := BuildFilter(opts...)
	expected := "report_format~XML%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportTask(t *testing.T) {
	opts := []Option{WithReportTask("task-123")}
	filter := BuildFilter(opts...)
	expected := "task=task-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportTaskLike(t *testing.T) {
	opts := []Option{WithReportTaskLike("task%")}
	filter := BuildFilter(opts...)
	expected := "task~task%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportScanner(t *testing.T) {
	opts := []Option{WithReportScanner("scanner-123")}
	filter := BuildFilter(opts...)
	expected := "scanner=scanner-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportScannerLike(t *testing.T) {
	opts := []Option{WithReportScannerLike("scanner%")}
	filter := BuildFilter(opts...)
	expected := "scanner~scanner%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportHosts(t *testing.T) {
	opts := []Option{WithReportHosts(10)}
	filter := BuildFilter(opts...)
	expected := "hosts=10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportHostsGreaterThan(t *testing.T) {
	opts := []Option{WithReportHostsGreaterThan(5)}
	filter := BuildFilter(opts...)
	expected := "hosts>5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportHostsGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithReportHostsGreaterThanOrEqual(5)}
	filter := BuildFilter(opts...)
	expected := "hosts>=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportHostsLessThan(t *testing.T) {
	opts := []Option{WithReportHostsLessThan(20)}
	filter := BuildFilter(opts...)
	expected := "hosts<20"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportHostsLessThanOrEqual(t *testing.T) {
	opts := []Option{WithReportHostsLessThanOrEqual(20)}
	filter := BuildFilter(opts...)
	expected := "hosts<=20"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportResultHosts(t *testing.T) {
	opts := []Option{WithReportResultHosts(8)}
	filter := BuildFilter(opts...)
	expected := "result_hosts=8"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportResultHostsGreaterThan(t *testing.T) {
	opts := []Option{WithReportResultHostsGreaterThan(3)}
	filter := BuildFilter(opts...)
	expected := "result_hosts>3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportResultHostsGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithReportResultHostsGreaterThanOrEqual(3)}
	filter := BuildFilter(opts...)
	expected := "result_hosts>=3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportResultHostsLessThan(t *testing.T) {
	opts := []Option{WithReportResultHostsLessThan(15)}
	filter := BuildFilter(opts...)
	expected := "result_hosts<15"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportResultHostsLessThanOrEqual(t *testing.T) {
	opts := []Option{WithReportResultHostsLessThanOrEqual(15)}
	filter := BuildFilter(opts...)
	expected := "result_hosts<=15"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportTotal(t *testing.T) {
	opts := []Option{WithReportTotal(100)}
	filter := BuildFilter(opts...)
	expected := "total=100"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportTotalGreaterThan(t *testing.T) {
	opts := []Option{WithReportTotalGreaterThan(50)}
	filter := BuildFilter(opts...)
	expected := "total>50"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportTotalGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithReportTotalGreaterThanOrEqual(50)}
	filter := BuildFilter(opts...)
	expected := "total>=50"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportTotalLessThan(t *testing.T) {
	opts := []Option{WithReportTotalLessThan(200)}
	filter := BuildFilter(opts...)
	expected := "total<200"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportTotalLessThanOrEqual(t *testing.T) {
	opts := []Option{WithReportTotalLessThanOrEqual(200)}
	filter := BuildFilter(opts...)
	expected := "total<=200"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportCritical(t *testing.T) {
	opts := []Option{WithReportCritical(5)}
	filter := BuildFilter(opts...)
	expected := "critical=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportCriticalGreaterThan(t *testing.T) {
	opts := []Option{WithReportCriticalGreaterThan(2)}
	filter := BuildFilter(opts...)
	expected := "critical>2"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportCriticalGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithReportCriticalGreaterThanOrEqual(2)}
	filter := BuildFilter(opts...)
	expected := "critical>=2"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportCriticalLessThan(t *testing.T) {
	opts := []Option{WithReportCriticalLessThan(10)}
	filter := BuildFilter(opts...)
	expected := "critical<10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportCriticalLessThanOrEqual(t *testing.T) {
	opts := []Option{WithReportCriticalLessThanOrEqual(10)}
	filter := BuildFilter(opts...)
	expected := "critical<=10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportHigh(t *testing.T) {
	opts := []Option{WithReportHigh(15)}
	filter := BuildFilter(opts...)
	expected := "high=15"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportHighGreaterThan(t *testing.T) {
	opts := []Option{WithReportHighGreaterThan(8)}
	filter := BuildFilter(opts...)
	expected := "high>8"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportHighGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithReportHighGreaterThanOrEqual(8)}
	filter := BuildFilter(opts...)
	expected := "high>=8"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportHighLessThan(t *testing.T) {
	opts := []Option{WithReportHighLessThan(25)}
	filter := BuildFilter(opts...)
	expected := "high<25"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportHighLessThanOrEqual(t *testing.T) {
	opts := []Option{WithReportHighLessThanOrEqual(25)}
	filter := BuildFilter(opts...)
	expected := "high<=25"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportMedium(t *testing.T) {
	opts := []Option{WithReportMedium(30)}
	filter := BuildFilter(opts...)
	expected := "medium=30"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportMediumGreaterThan(t *testing.T) {
	opts := []Option{WithReportMediumGreaterThan(15)}
	filter := BuildFilter(opts...)
	expected := "medium>15"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportMediumGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithReportMediumGreaterThanOrEqual(15)}
	filter := BuildFilter(opts...)
	expected := "medium>=15"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportMediumLessThan(t *testing.T) {
	opts := []Option{WithReportMediumLessThan(50)}
	filter := BuildFilter(opts...)
	expected := "medium<50"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportMediumLessThanOrEqual(t *testing.T) {
	opts := []Option{WithReportMediumLessThanOrEqual(50)}
	filter := BuildFilter(opts...)
	expected := "medium<=50"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportLow(t *testing.T) {
	opts := []Option{WithReportLow(50)}
	filter := BuildFilter(opts...)
	expected := "low=50"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportLowGreaterThan(t *testing.T) {
	opts := []Option{WithReportLowGreaterThan(25)}
	filter := BuildFilter(opts...)
	expected := "low>25"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportLowGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithReportLowGreaterThanOrEqual(25)}
	filter := BuildFilter(opts...)
	expected := "low>=25"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportLowLessThan(t *testing.T) {
	opts := []Option{WithReportLowLessThan(75)}
	filter := BuildFilter(opts...)
	expected := "low<75"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportLowLessThanOrEqual(t *testing.T) {
	opts := []Option{WithReportLowLessThanOrEqual(75)}
	filter := BuildFilter(opts...)
	expected := "low<=75"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportLog(t *testing.T) {
	opts := []Option{WithReportLog(20)}
	filter := BuildFilter(opts...)
	expected := "log=20"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportLogGreaterThan(t *testing.T) {
	opts := []Option{WithReportLogGreaterThan(10)}
	filter := BuildFilter(opts...)
	expected := "log>10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportLogGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithReportLogGreaterThanOrEqual(10)}
	filter := BuildFilter(opts...)
	expected := "log>=10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportLogLessThan(t *testing.T) {
	opts := []Option{WithReportLogLessThan(30)}
	filter := BuildFilter(opts...)
	expected := "log<30"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportLogLessThanOrEqual(t *testing.T) {
	opts := []Option{WithReportLogLessThanOrEqual(30)}
	filter := BuildFilter(opts...)
	expected := "log<=30"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportAlterable(t *testing.T) {
	opts := []Option{WithReportAlterable(true)}
	filter := BuildFilter(opts...)
	expected := "alterable=true"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportWritable(t *testing.T) {
	opts := []Option{WithReportWritable(false)}
	filter := BuildFilter(opts...)
	expected := "writable=false"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportInUse(t *testing.T) {
	opts := []Option{WithReportInUse(true)}
	filter := BuildFilter(opts...)
	expected := "in_use=true"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithReportTrash(t *testing.T) {
	opts := []Option{WithReportTrash(false)}
	filter := BuildFilter(opts...)
	expected := "trash=false"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}
