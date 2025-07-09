package commands

import "testing"

func TestNewPermissionSubject(t *testing.T) {
	sub := NewPermissionSubject("id-123", "user")
	if sub.ID != "id-123" {
		t.Errorf("expected ID 'id-123', got '%s'", sub.ID)
	}
	if sub.Type != "user" {
		t.Errorf("expected Type 'user', got '%s'", sub.Type)
	}
}

func TestNewPermissionResource(t *testing.T) {
	res := NewPermissionResource("id-456", "target")
	if res.ID != "id-456" {
		t.Errorf("expected ID 'id-456', got '%s'", res.ID)
	}
	if res.Type != "target" {
		t.Errorf("expected Type 'target', got '%s'", res.Type)
	}
}
