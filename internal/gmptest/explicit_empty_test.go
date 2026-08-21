//go:build integration

package gmptest

import (
	"context"
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

// TestExplicitEmptySemantics validates T-25 explicit-empty audit findings by
// asserting the documented wire behavior for fields whose omit/explicit-empty
// semantics have been classified as EQUIVALENT via live gvmd probes.
//
// Run with:
//
//	go test -tags integration ./internal/gmptest -run TestExplicitEmpty -v -gmp.pass=admin -gmp.insecure
func TestExplicitEmptySemantics(t *testing.T) {
	t.Run("modify_config_comment", testModifyConfigComment)
	t.Run("modify_port_list_comment", testModifyPortListComment)
	t.Run("modify_note_hosts", testModifyNoteHosts)
	t.Run("modify_note_port", testModifyNotePort)
	t.Run("modify_override_hosts", testModifyOverrideHosts)
	t.Run("modify_override_port", testModifyOverridePort)
	t.Run("modify_override_threat", testModifyOverrideThreat)
	t.Run("modify_override_new_threat", testModifyOverrideNewThreat)
}

// modify_user.sources.source segfault reproduction is intentionally omitted:
// triggering a gvmd crash against the shared dev stack is undesirable, and the
// typed API already cannot construct the <sources><source></source></sources>
// payload (UserSource is omitempty-tagged, so empty string is dropped).

// getConfigByID fetches one config by ID via RawXML, not gmp.Client's
// GetConfigs(ctx, args ...filtering.FilterArg): that signature only ever
// consults recognized FilterArg types (string, Option, their slices), so a
// *commands.GetConfigs passed as an arg is silently ignored and the call
// falls back to gvmd's default (unfiltered) list — fetching an arbitrary,
// unrelated config instead of the one this test just created.
func getConfigByID(t *testing.T, ctx context.Context, c interface {
	RawXML(ctx context.Context, xmlStr string) (string, error)
}, id string) commands.Config {
	t.Helper()
	out, err := c.RawXML(ctx, fmt.Sprintf(`<get_configs config_id="%s"/>`, id))
	if err != nil {
		t.Fatalf("RawXML get_configs: %v", err)
	}
	var resp commands.GetConfigsResponse
	if err := xml.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("decode get_configs_response: %v", err)
	}
	if len(resp.Config) == 0 {
		t.Fatalf("get_configs returned no config for %s", id)
	}
	return resp.Config[0]
}

// getPortListByID is getConfigByID's counterpart for GetPortLists, which has
// the same filtered-args-only signature on gmp.Client.
func getPortListByID(t *testing.T, ctx context.Context, c interface {
	RawXML(ctx context.Context, xmlStr string) (string, error)
}, id string) commands.PortList {
	t.Helper()
	out, err := c.RawXML(ctx, fmt.Sprintf(`<get_port_lists port_list_id="%s"/>`, id))
	if err != nil {
		t.Fatalf("RawXML get_port_lists: %v", err)
	}
	var resp commands.GetPortListsResponse
	if err := xml.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("decode get_port_lists_response: %v", err)
	}
	if len(resp.PortList) == 0 {
		t.Fatalf("get_port_lists returned no port list for %s", id)
	}
	return resp.PortList[0]
}

// getOverrideByID is getConfigByID's counterpart for GetOverrides — same
// filtered-args-only signature, plus details="1" is required for gvmd to
// include hosts/port/threat/new_threat in the response at all: gvmd's
// default (non-detailed) view omits them entirely, which would otherwise
// read as "cleared to empty" for every field this test checks.
func getOverrideByID(t *testing.T, ctx context.Context, c interface {
	RawXML(ctx context.Context, xmlStr string) (string, error)
}, id string) commands.Override {
	t.Helper()
	out, err := c.RawXML(ctx, fmt.Sprintf(`<get_overrides override_id="%s" details="1"/>`, id))
	if err != nil {
		t.Fatalf("RawXML get_overrides: %v", err)
	}
	var resp commands.GetOverridesResponse
	if err := xml.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("decode get_overrides_response: %v", err)
	}
	if len(resp.Override) == 0 {
		t.Fatalf("get_overrides returned no override for %s", id)
	}
	return resp.Override[0]
}

func testModifyConfigComment(t *testing.T) {
	c := requireLive(t)
	ctx := context.Background()

	configs, err := c.GetConfigs(ctx, &commands.GetConfigs{})
	if err != nil || len(configs.Config) == 0 {
		t.Skip("No configs available to copy")
	}
	baseCfgID := configs.Config[0].ID

	createResp, err := c.CreateConfig(ctx, &commands.CreateConfig{
		Name:      "test-explicit-empty-config",
		Comment:   "initial",
		Copy:      baseCfgID,
		UsageType: "scan",
	})
	if err != nil {
		t.Fatalf("CreateConfig: %v", err)
	}
	defer c.DeleteConfig(ctx, &commands.DeleteConfig{ConfigID: createResp.ID, Ultimate: true})

	if _, err := c.ModifyConfig(ctx, &commands.ModifyConfig{
		ConfigID: createResp.ID,
		Comment:  "Known value",
	}); err != nil {
		t.Fatalf("ModifyConfig (set): %v", err)
	}

	req := fmt.Sprintf(`<modify_config config_id="%s"></modify_config>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (omit): %v", err)
	}

	afterOmit := getConfigByID(t, ctx, c, createResp.ID).Comment

	if _, err := c.ModifyConfig(ctx, &commands.ModifyConfig{
		ConfigID: createResp.ID,
		Comment:  "Known value",
	}); err != nil {
		t.Fatalf("ModifyConfig (reset): %v", err)
	}

	req = fmt.Sprintf(`<modify_config config_id="%s"><comment></comment></modify_config>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (explicit-empty): %v", err)
	}

	afterEmpty := getConfigByID(t, ctx, c, createResp.ID).Comment

	if afterOmit != afterEmpty {
		t.Errorf("modify_config.comment: omit=%q, empty=%q; expected EQUIVALENT (both leave unchanged)", afterOmit, afterEmpty)
	}
}

func testModifyPortListComment(t *testing.T) {
	c := requireLive(t)
	ctx := context.Background()

	createResp, err := c.CreatePortList(ctx, &commands.CreatePortList{
		Name:      "test-explicit-empty-portlist",
		Comment:   "initial",
		PortRange: "T:1-10",
	})
	if err != nil {
		t.Fatalf("CreatePortList: %v", err)
	}
	defer c.DeletePortList(ctx, &commands.DeletePortList{PortListID: createResp.ID, Ultimate: true})

	if _, err := c.ModifyPortList(ctx, &commands.ModifyPortList{
		PortListID: createResp.ID,
		Comment:    "Known value",
	}); err != nil {
		t.Fatalf("ModifyPortList (set): %v", err)
	}

	req := fmt.Sprintf(`<modify_port_list port_list_id="%s"></modify_port_list>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (omit): %v", err)
	}

	afterOmit := getPortListByID(t, ctx, c, createResp.ID).Comment

	if _, err := c.ModifyPortList(ctx, &commands.ModifyPortList{
		PortListID: createResp.ID,
		Comment:    "Known value",
	}); err != nil {
		t.Fatalf("ModifyPortList (reset): %v", err)
	}

	req = fmt.Sprintf(`<modify_port_list port_list_id="%s"><comment></comment></modify_port_list>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (explicit-empty): %v", err)
	}

	afterEmpty := getPortListByID(t, ctx, c, createResp.ID).Comment

	if afterOmit != "" || afterEmpty != "" {
		t.Errorf("modify_port_list.comment: omit=%q, empty=%q; expected EQUIVALENT (both clear to empty)", afterOmit, afterEmpty)
	}
}

func testModifyNoteHosts(t *testing.T) {
	c := requireLive(t)
	ctx := context.Background()

	nvtsResp, err := c.GetNVTs(ctx, &commands.GetNVTs{Details: commands.Bool(false)})
	if err != nil || len(nvtsResp.NVT) == 0 {
		t.Skip("No NVTs available")
	}
	nvtOID := nvtsResp.NVT[0].OID

	createResp, err := c.CreateNote(ctx, &commands.CreateNote{
		Text:  "test note",
		NVT:   commands.NoteNVT{OID: nvtOID},
		Hosts: "192.168.1.1",
	})
	if err != nil {
		t.Fatalf("CreateNote: %v", err)
	}
	defer c.DeleteNote(ctx, &commands.DeleteNote{NoteID: createResp.ID, Ultimate: true})

	// Text is schema-required on every modify_note request (gvmd rejects an
	// empty <text></text> with "A TEXT entity is required"), so every set/reset
	// below must resend it alongside the field actually under test.
	if _, err := c.ModifyNote(ctx, &commands.ModifyNote{
		NoteID: createResp.ID,
		Text:   "test note",
		Hosts:  "Known value",
	}); err != nil {
		t.Fatalf("ModifyNote (set): %v", err)
	}

	req := fmt.Sprintf(`<modify_note note_id="%s"><text>test note</text></modify_note>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (omit): %v", err)
	}

	trueBool := commands.GMPBool(true)
	getResp, err := c.GetNotes(ctx, &commands.GetNotes{NoteID: createResp.ID, Details: &trueBool})
	if err != nil {
		t.Fatalf("GetNotes (after omit): %v", err)
	}
	if len(getResp.Note) == 0 {
		t.Fatalf("GetNotes returned empty response")
	}
	afterOmit := getResp.Note[0].Hosts

	if _, err := c.ModifyNote(ctx, &commands.ModifyNote{
		NoteID: createResp.ID,
		Text:   "test note",
		Hosts:  "Known value",
	}); err != nil {
		t.Fatalf("ModifyNote (reset): %v", err)
	}

	req = fmt.Sprintf(`<modify_note note_id="%s"><text>test note</text><hosts></hosts></modify_note>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (explicit-empty): %v", err)
	}

	getResp, err = c.GetNotes(ctx, &commands.GetNotes{NoteID: createResp.ID, Details: &trueBool})
	if err != nil {
		t.Fatalf("GetNotes (after empty): %v", err)
	}
	if len(getResp.Note) == 0 {
		t.Fatalf("GetNotes returned empty response")
	}
	afterEmpty := getResp.Note[0].Hosts

	if afterOmit != "" || afterEmpty != "" {
		t.Errorf("modify_note.hosts: omit=%q, empty=%q; expected EQUIVALENT (both clear to empty)", afterOmit, afterEmpty)
	}
}

func testModifyNotePort(t *testing.T) {
	c := requireLive(t)
	ctx := context.Background()

	nvtsResp, err := c.GetNVTs(ctx, &commands.GetNVTs{Details: commands.Bool(false)})
	if err != nil || len(nvtsResp.NVT) == 0 {
		t.Skip("No NVTs available")
	}
	nvtOID := nvtsResp.NVT[0].OID

	createResp, err := c.CreateNote(ctx, &commands.CreateNote{
		Text: "test note port",
		NVT:  commands.NoteNVT{OID: nvtOID},
		Port: "22/tcp",
	})
	if err != nil {
		t.Fatalf("CreateNote: %v", err)
	}
	defer c.DeleteNote(ctx, &commands.DeleteNote{NoteID: createResp.ID, Ultimate: true})

	if _, err := c.ModifyNote(ctx, &commands.ModifyNote{
		NoteID: createResp.ID,
		Text:   "test note port",
		Port:   "80/tcp",
	}); err != nil {
		t.Fatalf("ModifyNote (set): %v", err)
	}

	req := fmt.Sprintf(`<modify_note note_id="%s"><text>test note port</text></modify_note>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (omit): %v", err)
	}

	trueBool := commands.GMPBool(true)
	getResp, err := c.GetNotes(ctx, &commands.GetNotes{NoteID: createResp.ID, Details: &trueBool})
	if err != nil {
		t.Fatalf("GetNotes (after omit): %v", err)
	}
	if len(getResp.Note) == 0 {
		t.Fatalf("GetNotes returned empty response")
	}
	afterOmit := getResp.Note[0].Port

	if _, err := c.ModifyNote(ctx, &commands.ModifyNote{
		NoteID: createResp.ID,
		Text:   "test note port",
		Port:   "80/tcp",
	}); err != nil {
		t.Fatalf("ModifyNote (reset): %v", err)
	}

	req = fmt.Sprintf(`<modify_note note_id="%s"><text>test note port</text><port></port></modify_note>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (explicit-empty): %v", err)
	}

	getResp, err = c.GetNotes(ctx, &commands.GetNotes{NoteID: createResp.ID, Details: &trueBool})
	if err != nil {
		t.Fatalf("GetNotes (after empty): %v", err)
	}
	if len(getResp.Note) == 0 {
		t.Fatalf("GetNotes returned empty response")
	}
	afterEmpty := getResp.Note[0].Port

	if afterOmit != "" || afterEmpty != "" {
		t.Errorf("modify_note.port: omit=%q, empty=%q; expected EQUIVALENT (both clear to empty)", afterOmit, afterEmpty)
	}
}

func testModifyOverrideHosts(t *testing.T) {
	c := requireLive(t)
	ctx := context.Background()

	nvtsResp, err := c.GetNVTs(ctx, &commands.GetNVTs{Details: commands.Bool(false)})
	if err != nil || len(nvtsResp.NVT) == 0 {
		t.Skip("No NVTs available")
	}
	nvtOID := nvtsResp.NVT[0].OID

	createResp, err := c.CreateOverride(ctx, &commands.CreateOverride{
		Text:      "test override",
		NVT:       commands.OverrideNVT{OID: nvtOID},
		Hosts:     "192.168.1.1",
		NewThreat: "Log",
	})
	if err != nil {
		t.Fatalf("CreateOverride: %v", err)
	}
	defer c.DeleteOverride(ctx, &commands.DeleteOverride{OverrideID: createResp.ID, Ultimate: true})

	// Text is schema-required on every modify_override request, same as
	// modify_note. new_threat is also resent throughout: gvmd additionally
	// requires new_severity/new_threat whenever hosts/threat is present on a
	// modify_override request, a cross-field constraint the schema itself
	// does not declare.
	if _, err := c.ModifyOverride(ctx, &commands.ModifyOverride{
		OverrideID: createResp.ID,
		Text:       "test override",
		Hosts:      "10.0.0.1",
		NewThreat:  "Log",
	}); err != nil {
		t.Fatalf("ModifyOverride (set): %v", err)
	}

	req := fmt.Sprintf(`<modify_override override_id="%s"><text>test override</text><new_threat>Log</new_threat></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (omit): %v", err)
	}

	afterOmit := getOverrideByID(t, ctx, c, createResp.ID).Hosts

	if _, err := c.ModifyOverride(ctx, &commands.ModifyOverride{
		OverrideID: createResp.ID,
		Text:       "test override",
		Hosts:      "10.0.0.1",
		NewThreat:  "Log",
	}); err != nil {
		t.Fatalf("ModifyOverride (reset): %v", err)
	}

	req = fmt.Sprintf(`<modify_override override_id="%s"><text>test override</text><hosts></hosts><new_threat>Log</new_threat></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (explicit-empty): %v", err)
	}

	afterEmpty := getOverrideByID(t, ctx, c, createResp.ID).Hosts

	if afterOmit != "" || afterEmpty != "" {
		t.Errorf("modify_override.hosts: omit=%q, empty=%q; expected EQUIVALENT (both clear to empty)", afterOmit, afterEmpty)
	}
}

func testModifyOverridePort(t *testing.T) {
	c := requireLive(t)
	ctx := context.Background()

	nvtsResp, err := c.GetNVTs(ctx, &commands.GetNVTs{Details: commands.Bool(false)})
	if err != nil || len(nvtsResp.NVT) == 0 {
		t.Skip("No NVTs available")
	}
	nvtOID := nvtsResp.NVT[0].OID

	createResp, err := c.CreateOverride(ctx, &commands.CreateOverride{
		Text:      "test override port",
		NVT:       commands.OverrideNVT{OID: nvtOID},
		Port:      "22/tcp",
		NewThreat: "Log",
	})
	if err != nil {
		t.Fatalf("CreateOverride: %v", err)
	}
	defer c.DeleteOverride(ctx, &commands.DeleteOverride{OverrideID: createResp.ID, Ultimate: true})

	if _, err := c.ModifyOverride(ctx, &commands.ModifyOverride{
		OverrideID: createResp.ID,
		Text:       "test override port",
		Port:       "80/tcp",
		NewThreat:  "Log",
	}); err != nil {
		t.Fatalf("ModifyOverride (set): %v", err)
	}

	req := fmt.Sprintf(`<modify_override override_id="%s"><text>test override port</text><new_threat>Log</new_threat></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (omit): %v", err)
	}

	afterOmit := getOverrideByID(t, ctx, c, createResp.ID).Port

	if _, err := c.ModifyOverride(ctx, &commands.ModifyOverride{
		OverrideID: createResp.ID,
		Text:       "test override port",
		Port:       "80/tcp",
		NewThreat:  "Log",
	}); err != nil {
		t.Fatalf("ModifyOverride (reset): %v", err)
	}

	req = fmt.Sprintf(`<modify_override override_id="%s"><text>test override port</text><port></port><new_threat>Log</new_threat></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (explicit-empty): %v", err)
	}

	afterEmpty := getOverrideByID(t, ctx, c, createResp.ID).Port

	if afterOmit != "" || afterEmpty != "" {
		t.Errorf("modify_override.port: omit=%q, empty=%q; expected EQUIVALENT (both clear to empty)", afterOmit, afterEmpty)
	}
}

func testModifyOverrideThreat(t *testing.T) {
	c := requireLive(t)
	ctx := context.Background()

	nvtsResp, err := c.GetNVTs(ctx, &commands.GetNVTs{Details: commands.Bool(false)})
	if err != nil || len(nvtsResp.NVT) == 0 {
		t.Skip("No NVTs available")
	}
	nvtOID := nvtsResp.NVT[0].OID

	createResp, err := c.CreateOverride(ctx, &commands.CreateOverride{
		Text:      "test override threat",
		NVT:       commands.OverrideNVT{OID: nvtOID},
		Threat:    "High",
		NewThreat: "Low",
	})
	if err != nil {
		t.Fatalf("CreateOverride: %v", err)
	}
	defer c.DeleteOverride(ctx, &commands.DeleteOverride{OverrideID: createResp.ID, Ultimate: true})

	if _, err := c.ModifyOverride(ctx, &commands.ModifyOverride{
		OverrideID: createResp.ID,
		Text:       "test override threat",
		Threat:     "Medium",
		NewThreat:  "Low",
	}); err != nil {
		t.Fatalf("ModifyOverride (set): %v", err)
	}

	req := fmt.Sprintf(`<modify_override override_id="%s"><text>test override threat</text><new_threat>Low</new_threat></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (omit): %v", err)
	}

	afterOmit := getOverrideByID(t, ctx, c, createResp.ID).Threat

	if _, err := c.ModifyOverride(ctx, &commands.ModifyOverride{
		OverrideID: createResp.ID,
		Text:       "test override threat",
		Threat:     "Medium",
		NewThreat:  "Low",
	}); err != nil {
		t.Fatalf("ModifyOverride (reset): %v", err)
	}

	req = fmt.Sprintf(`<modify_override override_id="%s"><text>test override threat</text><threat></threat><new_threat>Low</new_threat></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (explicit-empty): %v", err)
	}

	afterEmpty := getOverrideByID(t, ctx, c, createResp.ID).Threat

	if afterOmit != "" || afterEmpty != "" {
		t.Errorf("modify_override.threat: omit=%q, empty=%q; expected EQUIVALENT (both clear to empty)", afterOmit, afterEmpty)
	}
}

func testModifyOverrideNewThreat(t *testing.T) {
	c := requireLive(t)
	ctx := context.Background()

	nvtsResp, err := c.GetNVTs(ctx, &commands.GetNVTs{Details: commands.Bool(false)})
	if err != nil || len(nvtsResp.NVT) == 0 {
		t.Skip("No NVTs available")
	}
	nvtOID := nvtsResp.NVT[0].OID

	createResp, err := c.CreateOverride(ctx, &commands.CreateOverride{
		Text:      "test override new threat",
		NVT:       commands.OverrideNVT{OID: nvtOID},
		NewThreat: "Low",
	})
	if err != nil {
		t.Fatalf("CreateOverride: %v", err)
	}
	defer c.DeleteOverride(ctx, &commands.DeleteOverride{OverrideID: createResp.ID, Ultimate: true})

	if _, err := c.ModifyOverride(ctx, &commands.ModifyOverride{
		OverrideID: createResp.ID,
		Text:       "test override new threat",
		NewThreat:  "High",
	}); err != nil {
		t.Fatalf("ModifyOverride (set): %v", err)
	}

	req := fmt.Sprintf(`<modify_override override_id="%s"><text>test override new threat</text></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (omit): %v", err)
	}

	afterOmit := getOverrideByID(t, ctx, c, createResp.ID).NewThreat

	if _, err := c.ModifyOverride(ctx, &commands.ModifyOverride{
		OverrideID: createResp.ID,
		Text:       "test override new threat",
		NewThreat:  "High",
	}); err != nil {
		t.Fatalf("ModifyOverride (reset): %v", err)
	}

	req = fmt.Sprintf(`<modify_override override_id="%s"><text>test override new threat</text><new_threat></new_threat></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (explicit-empty): %v", err)
	}

	afterEmpty := getOverrideByID(t, ctx, c, createResp.ID).NewThreat

	if afterOmit != "High" || afterEmpty != "High" {
		t.Errorf("modify_override.new_threat: omit=%q, empty=%q; expected EQUIVALENT (both leave unchanged at 'High')", afterOmit, afterEmpty)
	}
}
