package commands

import "testing"

func TestNewCreateAlert_BasicFields(t *testing.T) {
	name := "Test Alert"
	comment := "A test alert."
	copy := "some-uuid"
	cond := condition{Text: "x", Data: nil}
	evt := event{Text: "x", Data: nil}
	meth := method{Text: "x", Data: nil}
	filt := &filter{ID: "filter-uuid"}
	active := new(bool)
	*active = true

	cmd := NewCreateAlert(name, comment, copy, cond, evt, meth, filt, active)

	if cmd.Name != name {
		t.Errorf("expected Name %q, got %q", name, cmd.Name)
	}
	if cmd.Comment != comment {
		t.Errorf("expected Comment %q, got %q", comment, cmd.Comment)
	}
	if cmd.Copy != copy {
		t.Errorf("expected Copy %q, got %q", copy, cmd.Copy)
	}
}

func TestNewCreateAlert_Condition(t *testing.T) {
	cond := condition{
		Text: "Severity at least",
		Data: []AlertData{{Name: "severity", Text: "5.5"}},
	}
	cmd := NewCreateAlert("n", "c", "y", cond, event{}, method{}, nil, nil)
	if cmd.Condition.Text != cond.Text || len(cmd.Condition.Data) != 1 || cmd.Condition.Data[0] != cond.Data[0] {
		t.Errorf("unexpected Condition: %+v", cmd.Condition)
	}
}

func TestNewCreateAlert_Event(t *testing.T) {
	evt := event{
		Text: "Task run status changed",
		Data: []AlertData{{Name: "status", Text: "Done"}},
	}
	cmd := NewCreateAlert("n", "c", "y", condition{}, evt, method{}, nil, nil)
	if cmd.Event.Text != evt.Text || len(cmd.Event.Data) != 1 || cmd.Event.Data[0] != evt.Data[0] {
		t.Errorf("unexpected Event: %+v", cmd.Event)
	}
}

func TestNewCreateAlert_Method(t *testing.T) {
	meth := method{
		Text: "Email",
		Data: []AlertData{{Name: "to_address", Text: "sally@example.org"}},
	}
	cmd := NewCreateAlert("n", "c", "y", condition{}, event{}, meth, nil, nil)
	if cmd.Method.Text != meth.Text || len(cmd.Method.Data) != 1 || cmd.Method.Data[0] != meth.Data[0] {
		t.Errorf("unexpected Method: %+v", cmd.Method)
	}
}

func TestNewCreateAlert_Filter(t *testing.T) {
	filt := &filter{ID: "filter-uuid"}
	cmd := NewCreateAlert("n", "c", "y", condition{}, event{}, method{}, filt, nil)
	if cmd.Filter == nil || cmd.Filter.ID != filt.ID {
		t.Errorf("unexpected Filter: %+v", cmd.Filter)
	}
}

func TestNewCreateAlert_Active(t *testing.T) {
	active := new(bool)
	*active = true
	cmd := NewCreateAlert("n", "c", "y", condition{}, event{}, method{}, nil, active)
	if cmd.Active == nil || *cmd.Active != *active {
		t.Errorf("unexpected Active: %+v", cmd.Active)
	}
}

func TestNewAlertData(t *testing.T) {
	name := "severity"
	text := "5.5"
	data := NewAlertData(name, text)
	if data.Name != name {
		t.Errorf("expected Name %q, got %q", name, data.Name)
	}
	if data.Text != text {
		t.Errorf("expected Text %q, got %q", text, data.Text)
	}
}

func TestNewAlertCondition(t *testing.T) {
	text := "Severity at least"
	data := []AlertData{{Name: "severity", Text: "5.5"}}
	cond := NewAlertCondition(text, data)
	if cond.Text != text {
		t.Errorf("expected Text %q, got %q", text, cond.Text)
	}
	if len(cond.Data) != 1 || cond.Data[0] != data[0] {
		t.Errorf("unexpected Data: %+v", cond.Data)
	}
}

func TestNewAlertEvent(t *testing.T) {
	text := "Task run status changed"
	data := []AlertData{{Name: "status", Text: "Done"}}
	evt := NewAlertEvent(text, data)
	if evt.Text != text {
		t.Errorf("expected Text %q, got %q", text, evt.Text)
	}
	if len(evt.Data) != 1 || evt.Data[0] != data[0] {
		t.Errorf("unexpected Data: %+v", evt.Data)
	}
}

func TestNewAlertMethod(t *testing.T) {
	text := "Email"
	data := []AlertData{{Name: "to_address", Text: "sally@example.org"}}
	meth := NewAlertMethod(text, data)
	if meth.Text != text {
		t.Errorf("expected Text %q, got %q", text, meth.Text)
	}
	if len(meth.Data) != 1 || meth.Data[0] != data[0] {
		t.Errorf("unexpected Data: %+v", meth.Data)
	}
}
