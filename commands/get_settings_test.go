package commands

import (
	"encoding/xml"
	"testing"
)

// TestGetSettingsDualLayout is the direct behavioral proof for the
// FieldOverrides["get_settings:/settings"] entry: gvmd's own baked <example>
// nests <setting> inside <settings> (no live capture ever confirms this
// shape), while every live gvmd 26.36.1 capture sends <setting> as a flat
// top-level sibling of <settings>. Both must decode without loss so future
// generation cannot silently drop either.
func TestGetSettingsDualLayout(t *testing.T) {
	t.Run("flat top-level setting (live gvmd shape)", func(t *testing.T) {
		in := `<get_settings_response status="200" status_text="OK">` +
			`<setting id="5f5a8712-8017-11e1-8556-406186ea4fc5"><name>Rows Per Page</name><value>15</value></setting>` +
			`</get_settings_response>`
		var r GetSettingsResponse
		if err := xml.Unmarshal([]byte(in), &r); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if len(r.Setting) != 1 || r.Setting[0].ID != "5f5a8712-8017-11e1-8556-406186ea4fc5" || r.Setting[0].Name != "Rows Per Page" {
			t.Fatalf("top-level Setting = %+v, want one populated entry", r.Setting)
		}
		if len(r.Settings.Setting) != 0 {
			t.Fatalf("Settings.Setting = %+v, want empty for the flat layout", r.Settings.Setting)
		}
	})

	t.Run("nested setting inside settings (schema's own baked example)", func(t *testing.T) {
		in := `<get_settings_response status="200" status_text="OK">` +
			`<settings><setting id="5f5a8712-8017-11e1-8556-406186ea4fc5"><name>Rows Per Page</name><value>15</value></setting></settings>` +
			`</get_settings_response>`
		var r GetSettingsResponse
		if err := xml.Unmarshal([]byte(in), &r); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if len(r.Settings.Setting) != 1 || r.Settings.Setting[0].ID != "5f5a8712-8017-11e1-8556-406186ea4fc5" || r.Settings.Setting[0].Name != "Rows Per Page" {
			t.Fatalf("Settings.Setting = %+v, want one populated entry", r.Settings.Setting)
		}
		if len(r.Setting) != 0 {
			t.Fatalf("top-level Setting = %+v, want empty for the nested layout", r.Setting)
		}
	})
}
