package filtering

import (
	"testing"
)

func TestWithTargetID(t *testing.T) {
	opts := []Option{WithTargetID("target-123")}
	filter := BuildFilter(opts...)
	expected := "target_id=target-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetName(t *testing.T) {
	opts := []Option{WithTargetName("test-target")}
	filter := BuildFilter(opts...)
	expected := "target_name=test-target"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetNameLike(t *testing.T) {
	opts := []Option{WithTargetNameLike("test%")}
	filter := BuildFilter(opts...)
	expected := "target_name~test%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetHosts(t *testing.T) {
	opts := []Option{WithTargetHosts("192.168.1.0/24")}
	filter := BuildFilter(opts...)
	expected := "hosts=192.168.1.0/24"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetHostsLike(t *testing.T) {
	opts := []Option{WithTargetHostsLike("192.168.%")}
	filter := BuildFilter(opts...)
	expected := "hosts~192.168.%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetPorts(t *testing.T) {
	opts := []Option{WithTargetPorts("1-1000")}
	filter := BuildFilter(opts...)
	expected := "ports=1-1000"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetPortsLike(t *testing.T) {
	opts := []Option{WithTargetPortsLike("1-1%")}
	filter := BuildFilter(opts...)
	expected := "ports~1-1%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetPortList(t *testing.T) {
	opts := []Option{WithTargetPortList("portlist-123")}
	filter := BuildFilter(opts...)
	expected := "port_list=portlist-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetPortListLike(t *testing.T) {
	opts := []Option{WithTargetPortListLike("portlist%")}
	filter := BuildFilter(opts...)
	expected := "port_list~portlist%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetCredentials(t *testing.T) {
	opts := []Option{WithTargetCredentials("creds-123")}
	filter := BuildFilter(opts...)
	expected := "credentials=creds-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetCredentialsLike(t *testing.T) {
	opts := []Option{WithTargetCredentialsLike("creds%")}
	filter := BuildFilter(opts...)
	expected := "credentials~creds%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetAlterable(t *testing.T) {
	opts := []Option{WithTargetAlterable(true)}
	filter := BuildFilter(opts...)
	expected := "alterable=true"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetWritable(t *testing.T) {
	opts := []Option{WithTargetWritable(false)}
	filter := BuildFilter(opts...)
	expected := "writable=false"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetInUse(t *testing.T) {
	opts := []Option{WithTargetInUse(true)}
	filter := BuildFilter(opts...)
	expected := "in_use=true"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetTrash(t *testing.T) {
	opts := []Option{WithTargetTrash(false)}
	filter := BuildFilter(opts...)
	expected := "trash=false"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetTasks(t *testing.T) {
	opts := []Option{WithTargetTasks(5)}
	filter := BuildFilter(opts...)
	expected := "tasks=5"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetTasksGreaterThan(t *testing.T) {
	opts := []Option{WithTargetTasksGreaterThan(2)}
	filter := BuildFilter(opts...)
	expected := "tasks>2"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetTasksGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithTargetTasksGreaterThanOrEqual(2)}
	filter := BuildFilter(opts...)
	expected := "tasks>=2"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetTasksLessThan(t *testing.T) {
	opts := []Option{WithTargetTasksLessThan(10)}
	filter := BuildFilter(opts...)
	expected := "tasks<10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTargetTasksLessThanOrEqual(t *testing.T) {
	opts := []Option{WithTargetTasksLessThanOrEqual(10)}
	filter := BuildFilter(opts...)
	expected := "tasks<=10"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}
