package commands

import "testing"

func TestIntHelper(t *testing.T) {
	if got := Int(0); got == nil || *got != 0 {
		t.Errorf("Int(0) = %v, want pointer to 0", got)
	}
	if got := Int(42); got == nil || *got != 42 {
		t.Errorf("Int(42) = %v, want pointer to 42", got)
	}
}

func TestFloat64Helper(t *testing.T) {
	if got := Float64(0); got == nil || *got != 0 {
		t.Errorf("Float64(0) = %v, want pointer to 0", got)
	}
	if got := Float64(2.5); got == nil || *got != 2.5 {
		t.Errorf("Float64(2.5) = %v, want pointer to 2.5", got)
	}
}
