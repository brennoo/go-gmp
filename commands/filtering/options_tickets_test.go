package filtering

import (
	"testing"
)

func TestWithTicketID(t *testing.T) {
	opts := []Option{WithTicketID("ticket-123")}
	filter := BuildFilter(opts...)
	expected := "ticket_id=ticket-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketName(t *testing.T) {
	opts := []Option{WithTicketName("test-ticket")}
	filter := BuildFilter(opts...)
	expected := "ticket_name=test-ticket"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketNameLike(t *testing.T) {
	opts := []Option{WithTicketNameLike("test%")}
	filter := BuildFilter(opts...)
	expected := "ticket_name~test%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketStatus(t *testing.T) {
	opts := []Option{WithTicketStatus("Open")}
	filter := BuildFilter(opts...)
	expected := "ticket_status=Open"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketStatusIn(t *testing.T) {
	opts := []Option{WithTicketStatusIn("Open", "Closed")}
	filter := BuildFilter(opts...)
	expected := "ticket_status=Open ticket_status=Closed"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketAssignedTo(t *testing.T) {
	opts := []Option{WithTicketAssignedTo("admin")}
	filter := BuildFilter(opts...)
	expected := "assigned_to=admin"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketAssignedBy(t *testing.T) {
	opts := []Option{WithTicketAssignedBy("user1")}
	filter := BuildFilter(opts...)
	expected := "assigned_by=user1"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketResult(t *testing.T) {
	opts := []Option{WithTicketResult("result-123")}
	filter := BuildFilter(opts...)
	expected := "result=result-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketResultLike(t *testing.T) {
	opts := []Option{WithTicketResultLike("result%")}
	filter := BuildFilter(opts...)
	expected := "result~result%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketTask(t *testing.T) {
	opts := []Option{WithTicketTask("task-123")}
	filter := BuildFilter(opts...)
	expected := "task=task-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketTaskLike(t *testing.T) {
	opts := []Option{WithTicketTaskLike("task%")}
	filter := BuildFilter(opts...)
	expected := "task~task%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketHost(t *testing.T) {
	opts := []Option{WithTicketHost("192.168.1.1")}
	filter := BuildFilter(opts...)
	expected := "host=192.168.1.1"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketHostLike(t *testing.T) {
	opts := []Option{WithTicketHostLike("192.168.%")}
	filter := BuildFilter(opts...)
	expected := "host~192.168.%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketPort(t *testing.T) {
	opts := []Option{WithTicketPort("80")}
	filter := BuildFilter(opts...)
	expected := "port=80"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketPortLike(t *testing.T) {
	opts := []Option{WithTicketPortLike("8%")}
	filter := BuildFilter(opts...)
	expected := "port~8%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketNVT(t *testing.T) {
	opts := []Option{WithTicketNVT("nvt-123")}
	filter := BuildFilter(opts...)
	expected := "nvt=nvt-123"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketNVTLike(t *testing.T) {
	opts := []Option{WithTicketNVTLike("nvt%")}
	filter := BuildFilter(opts...)
	expected := "nvt~nvt%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketAlterable(t *testing.T) {
	opts := []Option{WithTicketAlterable(true)}
	filter := BuildFilter(opts...)
	expected := "alterable=true"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketWritable(t *testing.T) {
	opts := []Option{WithTicketWritable(false)}
	filter := BuildFilter(opts...)
	expected := "writable=false"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketInUse(t *testing.T) {
	opts := []Option{WithTicketInUse(true)}
	filter := BuildFilter(opts...)
	expected := "in_use=true"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTicketTrash(t *testing.T) {
	opts := []Option{WithTicketTrash(false)}
	filter := BuildFilter(opts...)
	expected := "trash=false"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}
