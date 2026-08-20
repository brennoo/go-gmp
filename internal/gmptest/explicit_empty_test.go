//go:build integration

package gmptest

import (
	"context"
	"fmt"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

// TestExplicitEmptySemantics validates T-25 explicit-empty audit findings by
// asserting the documented wire behavior for fields whose omit/explicit-empty
// semantics have been classified as EQUIVALENT via live gvmd probes.
//
// Run with: go test -tags integration ./internal/gmptest -run TestExplicitEmpty -v -gmp.pass=admin -gmp.insecure
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

	getResp, err := c.GetConfigs(ctx, &commands.GetConfigs{ConfigID: createResp.ID})
	if err != nil {
		t.Fatalf("GetConfigs (after omit): %v", err)
	}
	if len(getResp.Config) == 0 {
		t.Fatalf("GetConfigs returned empty response")
	}
	afterOmit := getResp.Config[0].Comment

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

	getResp, err = c.GetConfigs(ctx, &commands.GetConfigs{ConfigID: createResp.ID})
	if err != nil {
		t.Fatalf("GetConfigs (after empty): %v", err)
	}
	if len(getResp.Config) == 0 {
		t.Fatalf("GetConfigs returned empty response")
	}
	afterEmpty := getResp.Config[0].Comment

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

	getResp, err := c.GetPortLists(ctx, &commands.GetPortLists{PortListID: createResp.ID})
	if err != nil {
		t.Fatalf("GetPortLists (after omit): %v", err)
	}
	if len(getResp.PortList) == 0 {
		t.Fatalf("GetPortLists returned empty response")
	}
	afterOmit := getResp.PortList[0].Comment

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

	getResp, err = c.GetPortLists(ctx, &commands.GetPortLists{PortListID: createResp.ID})
	if err != nil {
		t.Fatalf("GetPortLists (after empty): %v", err)
	}
	if len(getResp.PortList) == 0 {
		t.Fatalf("GetPortLists returned empty response")
	}
	afterEmpty := getResp.PortList[0].Comment

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

	if _, err := c.ModifyNote(ctx, &commands.ModifyNote{
		NoteID: createResp.ID,
		Hosts:  "Known value",
	}); err != nil {
		t.Fatalf("ModifyNote (set): %v", err)
	}

	req := fmt.Sprintf(`<modify_note note_id="%s"></modify_note>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (omit): %v", err)
	}

	getResp, err := c.GetNotes(ctx, &commands.GetNotes{NoteID: createResp.ID})
	if err != nil {
		t.Fatalf("GetNotes (after omit): %v", err)
	}
	if len(getResp.Note) == 0 {
		t.Fatalf("GetNotes returned empty response")
	}
	afterOmit := getResp.Note[0].Hosts

	if _, err := c.ModifyNote(ctx, &commands.ModifyNote{
		NoteID: createResp.ID,
		Hosts:  "Known value",
	}); err != nil {
		t.Fatalf("ModifyNote (reset): %v", err)
	}

	req = fmt.Sprintf(`<modify_note note_id="%s"><hosts></hosts></modify_note>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (explicit-empty): %v", err)
	}

	getResp, err = c.GetNotes(ctx, &commands.GetNotes{NoteID: createResp.ID})
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
		Port:   "Known value",
	}); err != nil {
		t.Fatalf("ModifyNote (set): %v", err)
	}

	req := fmt.Sprintf(`<modify_note note_id="%s"></modify_note>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (omit): %v", err)
	}

	getResp, err := c.GetNotes(ctx, &commands.GetNotes{NoteID: createResp.ID})
	if err != nil {
		t.Fatalf("GetNotes (after omit): %v", err)
	}
	if len(getResp.Note) == 0 {
		t.Fatalf("GetNotes returned empty response")
	}
	afterOmit := getResp.Note[0].Port

	if _, err := c.ModifyNote(ctx, &commands.ModifyNote{
		NoteID: createResp.ID,
		Port:   "Known value",
	}); err != nil {
		t.Fatalf("ModifyNote (reset): %v", err)
	}

	req = fmt.Sprintf(`<modify_note note_id="%s"><port></port></modify_note>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (explicit-empty): %v", err)
	}

	getResp, err = c.GetNotes(ctx, &commands.GetNotes{NoteID: createResp.ID})
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

	if _, err := c.ModifyOverride(ctx, &commands.ModifyOverride{
		OverrideID: createResp.ID,
		Hosts:      "Known value",
	}); err != nil {
		t.Fatalf("ModifyOverride (set): %v", err)
	}

	req := fmt.Sprintf(`<modify_override override_id="%s"></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (omit): %v", err)
	}

	getResp, err := c.GetOverrides(ctx, &commands.GetOverrides{OverrideID: createResp.ID})
	if err != nil {
		t.Fatalf("GetOverrides (after omit): %v", err)
	}
	if len(getResp.Override) == 0 {
		t.Fatalf("GetOverrides returned empty response")
	}
	afterOmit := getResp.Override[0].Hosts

	if _, err := c.ModifyOverride(ctx, &commands.ModifyOverride{
		OverrideID: createResp.ID,
		Hosts:      "Known value",
	}); err != nil {
		t.Fatalf("ModifyOverride (reset): %v", err)
	}

	req = fmt.Sprintf(`<modify_override override_id="%s"><hosts></hosts></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (explicit-empty): %v", err)
	}

	getResp, err = c.GetOverrides(ctx, &commands.GetOverrides{OverrideID: createResp.ID})
	if err != nil {
		t.Fatalf("GetOverrides (after empty): %v", err)
	}
	if len(getResp.Override) == 0 {
		t.Fatalf("GetOverrides returned empty response")
	}
	afterEmpty := getResp.Override[0].Hosts

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
		Port:       "Known value",
	}); err != nil {
		t.Fatalf("ModifyOverride (set): %v", err)
	}

	req := fmt.Sprintf(`<modify_override override_id="%s"></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (omit): %v", err)
	}

	getResp, err := c.GetOverrides(ctx, &commands.GetOverrides{OverrideID: createResp.ID})
	if err != nil {
		t.Fatalf("GetOverrides (after omit): %v", err)
	}
	if len(getResp.Override) == 0 {
		t.Fatalf("GetOverrides returned empty response")
	}
	afterOmit := getResp.Override[0].Port

	if _, err := c.ModifyOverride(ctx, &commands.ModifyOverride{
		OverrideID: createResp.ID,
		Port:       "Known value",
	}); err != nil {
		t.Fatalf("ModifyOverride (reset): %v", err)
	}

	req = fmt.Sprintf(`<modify_override override_id="%s"><port></port></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (explicit-empty): %v", err)
	}

	getResp, err = c.GetOverrides(ctx, &commands.GetOverrides{OverrideID: createResp.ID})
	if err != nil {
		t.Fatalf("GetOverrides (after empty): %v", err)
	}
	if len(getResp.Override) == 0 {
		t.Fatalf("GetOverrides returned empty response")
	}
	afterEmpty := getResp.Override[0].Port

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
		Threat:     "Medium",
	}); err != nil {
		t.Fatalf("ModifyOverride (set): %v", err)
	}

	req := fmt.Sprintf(`<modify_override override_id="%s"></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (omit): %v", err)
	}

	getResp, err := c.GetOverrides(ctx, &commands.GetOverrides{OverrideID: createResp.ID})
	if err != nil {
		t.Fatalf("GetOverrides (after omit): %v", err)
	}
	if len(getResp.Override) == 0 {
		t.Fatalf("GetOverrides returned empty response")
	}
	afterOmit := getResp.Override[0].Threat

	if _, err := c.ModifyOverride(ctx, &commands.ModifyOverride{
		OverrideID: createResp.ID,
		Threat:     "Medium",
	}); err != nil {
		t.Fatalf("ModifyOverride (reset): %v", err)
	}

	req = fmt.Sprintf(`<modify_override override_id="%s"><threat></threat></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (explicit-empty): %v", err)
	}

	getResp, err = c.GetOverrides(ctx, &commands.GetOverrides{OverrideID: createResp.ID})
	if err != nil {
		t.Fatalf("GetOverrides (after empty): %v", err)
	}
	if len(getResp.Override) == 0 {
		t.Fatalf("GetOverrides returned empty response")
	}
	afterEmpty := getResp.Override[0].Threat

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
		NewThreat:  "High",
	}); err != nil {
		t.Fatalf("ModifyOverride (set): %v", err)
	}

	req := fmt.Sprintf(`<modify_override override_id="%s"></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (omit): %v", err)
	}

	getResp, err := c.GetOverrides(ctx, &commands.GetOverrides{OverrideID: createResp.ID})
	if err != nil {
		t.Fatalf("GetOverrides (after omit): %v", err)
	}
	if len(getResp.Override) == 0 {
		t.Fatalf("GetOverrides returned empty response")
	}
	afterOmit := getResp.Override[0].NewThreat

	if _, err := c.ModifyOverride(ctx, &commands.ModifyOverride{
		OverrideID: createResp.ID,
		NewThreat:  "High",
	}); err != nil {
		t.Fatalf("ModifyOverride (reset): %v", err)
	}

	req = fmt.Sprintf(`<modify_override override_id="%s"><new_threat></new_threat></modify_override>`, createResp.ID)
	if _, err := c.RawXML(ctx, req); err != nil {
		t.Fatalf("RawXML (explicit-empty): %v", err)
	}

	getResp, err = c.GetOverrides(ctx, &commands.GetOverrides{OverrideID: createResp.ID})
	if err != nil {
		t.Fatalf("GetOverrides (after empty): %v", err)
	}
	if len(getResp.Override) == 0 {
		t.Fatalf("GetOverrides returned empty response")
	}
	afterEmpty := getResp.Override[0].NewThreat

	if afterOmit != "High" || afterEmpty != "High" {
		t.Errorf("modify_override.new_threat: omit=%q, empty=%q; expected EQUIVALENT (both leave unchanged at 'High')", afterOmit, afterEmpty)
	}
}
