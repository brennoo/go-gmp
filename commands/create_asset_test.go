package commands

import "testing"

func TestNewAsset(t *testing.T) {
	name := "Localhost"
	comment := "A test asset."
	typ := "host"
	a := NewAsset(name, comment, typ)
	if a.Name != name {
		t.Errorf("expected Name %q, got %q", name, a.Name)
	}
	if a.Comment != comment {
		t.Errorf("expected Comment %q, got %q", comment, a.Comment)
	}
	if a.Type != typ {
		t.Errorf("expected Type %q, got %q", typ, a.Type)
	}
}

func TestNewAssetReport(t *testing.T) {
	id := "report-uuid"
	filter := NewFilter("min_qod=80")
	r := NewAssetReport(id, filter)
	if r.ID != id {
		t.Errorf("expected ID %q, got %q", id, r.ID)
	}
	if r.Filter == nil || r.Filter.Term != filter.Term {
		t.Errorf("unexpected Filter: %+v", r.Filter)
	}
}

func TestNewCreateAsset(t *testing.T) {
	a := NewAsset("Localhost", "", "host")
	r := NewAssetReport("report-uuid", nil)
	cmd := NewCreateAsset(a, r)
	if cmd.Asset != a {
		t.Errorf("expected Asset %+v, got %+v", a, cmd.Asset)
	}
	if cmd.Report != r {
		t.Errorf("expected Report %+v, got %+v", r, cmd.Report)
	}
}
