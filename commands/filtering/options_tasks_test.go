package filtering

import (
	"testing"
)

func TestWithTaskID(t *testing.T) {
	opts := []Option{WithTaskID("task-123")}
	filter := BuildFilter(opts...)
	expected := "uuid=task-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskName(t *testing.T) {
	opts := []Option{WithTaskName("test-task")}
	filter := BuildFilter(opts...)
	expected := "task_name=test-task"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskNameLike(t *testing.T) {
	opts := []Option{WithTaskNameLike("test%")}
	filter := BuildFilter(opts...)
	expected := "task_name~test%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskStatus(t *testing.T) {
	opts := []Option{WithTaskStatus("Running")}
	filter := BuildFilter(opts...)
	expected := "task_status=Running"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskStatusIn(t *testing.T) {
	opts := []Option{WithTaskStatusIn("Running", "Done")}
	filter := BuildFilter(opts...)
	expected := "task_status=Running task_status=Done"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskProgress(t *testing.T) {
	opts := []Option{WithTaskProgress(50)}
	filter := BuildFilter(opts...)
	expected := "progress=50"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskProgressGreaterThan(t *testing.T) {
	opts := []Option{WithTaskProgressGreaterThan(25)}
	filter := BuildFilter(opts...)
	expected := "progress>25"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskProgressGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithTaskProgressGreaterThanOrEqual(25)}
	filter := BuildFilter(opts...)
	expected := "progress>=25"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskProgressLessThan(t *testing.T) {
	opts := []Option{WithTaskProgressLessThan(75)}
	filter := BuildFilter(opts...)
	expected := "progress<75"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskProgressLessThanOrEqual(t *testing.T) {
	opts := []Option{WithTaskProgressLessThanOrEqual(75)}
	filter := BuildFilter(opts...)
	expected := "progress<=75"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskAlterable(t *testing.T) {
	opts := []Option{WithTaskAlterable(true)}
	filter := BuildFilter(opts...)
	expected := "alterable=true"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskWritable(t *testing.T) {
	opts := []Option{WithTaskWritable(false)}
	filter := BuildFilter(opts...)
	expected := "writable=false"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskInUse(t *testing.T) {
	opts := []Option{WithTaskInUse(true)}
	filter := BuildFilter(opts...)
	expected := "in_use=true"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskTrash(t *testing.T) {
	opts := []Option{WithTaskTrash(false)}
	filter := BuildFilter(opts...)
	expected := "trash=false"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskSchedule(t *testing.T) {
	opts := []Option{WithTaskSchedule("schedule-123")}
	filter := BuildFilter(opts...)
	expected := "schedule=schedule-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskScheduleLike(t *testing.T) {
	opts := []Option{WithTaskScheduleLike("schedule%")}
	filter := BuildFilter(opts...)
	expected := "schedule~schedule%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskScanner(t *testing.T) {
	opts := []Option{WithTaskScanner("scanner-123")}
	filter := BuildFilter(opts...)
	expected := "scanner=scanner-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskScannerLike(t *testing.T) {
	opts := []Option{WithTaskScannerLike("scanner%")}
	filter := BuildFilter(opts...)
	expected := "scanner~scanner%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskTarget(t *testing.T) {
	opts := []Option{WithTaskTarget("target-123")}
	filter := BuildFilter(opts...)
	expected := "target=target-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskTargetLike(t *testing.T) {
	opts := []Option{WithTaskTargetLike("target%")}
	filter := BuildFilter(opts...)
	expected := "target~target%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskConfig(t *testing.T) {
	opts := []Option{WithTaskConfig("config-123")}
	filter := BuildFilter(opts...)
	expected := "config=config-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTaskConfigLike(t *testing.T) {
	opts := []Option{WithTaskConfigLike("config%")}
	filter := BuildFilter(opts...)
	expected := "config~config%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}
