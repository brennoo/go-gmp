package client

import (
	"context"
	"strings"
	"testing"

	"github.com/brennoo/go-gmp/commands"
	"github.com/brennoo/go-gmp/commands/filtering"
)

func TestNew(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}
}

func TestGetPreferences(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetPreferences{}
	cmd.ConfigID = "4b49617e-d1d8-44b8-af81-f4675b56f837"
	resp, err := cli.GetPreferences(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetPreferences: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

// nolint:gocyclo
// gocyclo:ignore
func TestGetSystemReports(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case: get a specific report
	cmd := &commands.GetSystemReports{
		Name: "proc",
	}
	resp, err := cli.GetSystemReports(cmd)
	if err != nil {
		t.Fatalf("GetSystemReports failed: %v", err)
	}
	if resp.Status != "200" || resp.StatusText != "OK" || len(resp.SystemReports) != 1 {
		t.Errorf("unexpected response: %+v", resp)
	}
	if resp.SystemReports[0].Name != "proc" || resp.SystemReports[0].Title != "Processes" {
		t.Errorf("unexpected system report: %+v", resp.SystemReports[0])
	}
	if resp.SystemReports[0].Report == nil || resp.SystemReports[0].Report.Format != "png" {
		t.Errorf("unexpected report content: %+v", resp.SystemReports[0].Report)
	}

	// Success case: get brief listing
	cmd = &commands.GetSystemReports{
		Brief: true,
	}
	resp, err = cli.GetSystemReports(cmd)
	if err != nil {
		t.Fatalf("GetSystemReports failed: %v", err)
	}
	if resp.Status != "200" || resp.StatusText != "OK" || len(resp.SystemReports) < 2 {
		t.Errorf("unexpected response: %+v", resp)
	}

	// Failure case
	cmd = &commands.GetSystemReports{
		Name: "notfound",
	}
	resp, err = cli.GetSystemReports(cmd)
	if err != nil {
		t.Fatalf("GetSystemReports failed: %v", err)
	}
	if resp.Status != "404" || resp.StatusText != "Not found" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestGetInfo(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetInfo{
		Name: "CVE-2011-0018",
		Type: "cve",
	}
	resp, err := cli.GetInfo(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetInfo: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if len(resp.Infos) != 1 {
		t.Fatalf("Expected 1 info, got %d", len(resp.Infos))
	}
	info := resp.Infos[0]
	if info.ID != "CVE-2011-0018" || info.Name != "CVE-2011-0018" {
		t.Fatalf("Unexpected info: %+v", info)
	}
}

func TestGetVersion(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetVersion{}
	resp, err := cli.GetVersion(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetVersion: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if resp.Version != "1.0" {
		t.Fatalf("Unexpected version. Expected: 1.0 Got: %s", resp.Version)
	}
}

func TestHelp(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.Help{}
	resp, err := cli.Help(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during Help: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if resp.Text == "" || !contains(resp.Text, "AUTHENTICATE") {
		t.Fatalf("Expected help text to contain 'AUTHENTICATE', got: %s", resp.Text)
	}
}

func TestGetFeeds(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetFeeds{}
	resp, err := cli.GetFeeds(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetFeeds: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if resp.FeedOwnerSet != true {
		t.Fatalf("Unexpected feed_owner_set. Expected: true Got: %t", resp.FeedOwnerSet)
	}
	if resp.FeedRolesSet != true {
		t.Fatalf("Unexpected feed_roles_set. Expected: true Got: %t", resp.FeedRolesSet)
	}
	if resp.FeedResourcesAccess != true {
		t.Fatalf("Unexpected feed_resources_access. Expected: true Got: %t", resp.FeedResourcesAccess)
	}
	if len(resp.Feeds) != 3 {
		t.Fatalf("Expected 3 feeds, got %d", len(resp.Feeds))
	}
	if resp.Feeds[0].Type != "NVT" || resp.Feeds[0].Name != "Greenbone Security Feed" {
		t.Fatalf("Unexpected first feed: %+v", resp.Feeds[0])
	}
	if resp.Feeds[1].Type != "CERT" || resp.Feeds[1].Name != "Greenbone CERT Feed" {
		t.Fatalf("Unexpected second feed: %+v", resp.Feeds[1])
	}
	if resp.Feeds[2].Type != "SCAP" || resp.Feeds[2].Name != "Greenbone SCAP Feed" {
		t.Fatalf("Unexpected third feed: %+v", resp.Feeds[2])
	}
}

func TestModifyLicense(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.ModifyLicense{
		File: "YmFzZTY0bGljZW5zZQ==", // "baselicenses" in base64
	}
	resp, err := cli.ModifyLicense(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyLicense: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
}

// nolint:gocyclo
// gocyclo:ignore
func TestGetLicense(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetLicense{}
	resp, err := cli.GetLicense(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetLicense: %s", err)
	}

	hasErr := false
	if resp == nil || resp.License == nil || resp.License.Content == nil || resp.License.Content.Meta == nil ||
		resp.License.Content.Appliance == nil || resp.License.Content.Keys == nil || resp.License.Content.Signatures == nil {
		t.Error("Missing required license fields")
		hasErr = true
	}
	if resp.Status != "200" {
		t.Error("Unexpected status:", resp.Status)
		hasErr = true
	}
	if resp.StatusText != "OK" {
		t.Error("Unexpected status text:", resp.StatusText)
		hasErr = true
	}
	if resp.License.Status != "active" {
		t.Error("Unexpected license status:", resp.License.Status)
		hasErr = true
	}
	if resp.License.Content.Meta.ID != "4711" {
		t.Error("Unexpected license id:", resp.License.Content.Meta.ID)
		hasErr = true
	}
	if resp.License.Content.Meta.CustomerName != "Jane Doe" {
		t.Error("Unexpected customer name:", resp.License.Content.Meta.CustomerName)
		hasErr = true
	}
	if resp.License.Content.Appliance.Model != "trial" {
		t.Error("Unexpected appliance model:", resp.License.Content.Appliance.Model)
		hasErr = true
	}
	if len(resp.License.Content.Keys.Keys) != 1 || resp.License.Content.Keys.Keys[0].Name != "feed" {
		t.Error("Unexpected keys:", resp.License.Content.Keys)
		hasErr = true
	}
	if len(resp.License.Content.Signatures.Signatures) != 1 || resp.License.Content.Signatures.Signatures[0].Name != "license" {
		t.Error("Unexpected signatures:", resp.License.Content.Signatures)
		hasErr = true
	}
	if hasErr {
		t.FailNow()
	}
}

func TestModifySetting(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.ModifySetting{
		Name:  "Timezone",
		Value: "QWZyaWNhL0pvaGFubmVzYnVyZw==", // "Africa/Johannesburg" in base64
	}
	resp, err := cli.ModifySetting(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifySetting: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
}

func TestGetSettings(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()
	resp, err := cli.GetSettings(ctx)
	if err != nil {
		t.Fatalf("Unexpected error during GetSettings: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if resp.Settings == nil || len(resp.Settings.Setting) == 0 {
		t.Fatalf("Expected at least one setting, got %+v", resp.Settings)
	}
	if resp.Settings.Setting[0].ID != "5f5a8712-8017-11e1-8556-406186ea4fc5" {
		t.Fatalf("Unexpected setting id. Expected: 5f5a8712-8017-11e1-8556-406186ea4fc5 Got: %s", resp.Settings.Setting[0].ID)
	}
	if resp.Settings.Setting[0].Name != "Rows Per Page" {
		t.Fatalf("Unexpected setting name. Expected: Rows Per Page Got: %s", resp.Settings.Setting[0].Name)
	}
	if resp.Settings.Setting[0].Value != "15" {
		t.Fatalf("Unexpected setting value. Expected: 15 Got: %s", resp.Settings.Setting[0].Value)
	}
}

func TestGetResourceNames(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetResourceNames{Type: "os"}
	resp, err := cli.GetResourceNames(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetResourceNames: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if resp.Type != "os" {
		t.Fatalf("Unexpected type. Expected: os Got: %s", resp.Type)
	}
	if len(resp.Resources) != 2 {
		t.Fatalf("Expected 2 resources, got %d", len(resp.Resources))
	}
	if resp.Resources[0].ID != "5b6b6aef-b320-42ca-899f-3161ee2a4195" {
		t.Fatalf("Unexpected resource id. Expected: 5b6b6aef-b320-42ca-899f-3161ee2a4195 Got: %s", resp.Resources[0].ID)
	}
	if resp.Resources[0].Name != "cpe:/h:hp:jetdirect" {
		t.Fatalf("Unexpected resource name. Expected: cpe:/h:hp:jetdirect Got: %s", resp.Resources[0].Name)
	}
	if resp.Resources[1].ID != "5be25864-9249-431e-8a91-039e334371ad" {
		t.Fatalf("Unexpected resource id. Expected: 5be25864-9249-431e-8a91-039e334371ad Got: %s", resp.Resources[1].ID)
	}
	if resp.Resources[1].Name != "cpe:/o:canonical:ubuntu_linux:18.04" {
		t.Fatalf("Unexpected resource name. Expected: cpe:/o:canonical:ubuntu_linux:18.04 Got: %s", resp.Resources[1].Name)
	}
}

func TestGetAggregates(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetAggregates{Type: "nvt", GroupColumn: "family", DataColumn: "severity"}
	resp, err := cli.GetAggregates(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetAggregates: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if resp.Aggregate == nil {
		t.Fatalf("Expected aggregate block, got nil")
	}
	if resp.Aggregate.DataType != "nvt" {
		t.Fatalf("Unexpected data_type. Expected: nvt Got: %s", resp.Aggregate.DataType)
	}
	if resp.Aggregate.DataColumn != "severity" {
		t.Fatalf("Unexpected data_column. Expected: severity Got: %s", resp.Aggregate.DataColumn)
	}
	if resp.Aggregate.GroupColumn != "family" {
		t.Fatalf("Unexpected group_column. Expected: family Got: %s", resp.Aggregate.GroupColumn)
	}
	if len(resp.Aggregate.Groups) != 2 {
		t.Fatalf("Expected 2 groups, got %d", len(resp.Aggregate.Groups))
	}
	if resp.Aggregate.Groups[0].Value != "AIX Local Security Checks" {
		t.Fatalf("Unexpected group value. Expected: AIX Local Security Checks Got: %s", resp.Aggregate.Groups[0].Value)
	}
	if resp.Aggregate.Groups[0].Count != 1 {
		t.Fatalf("Unexpected group count. Expected: 1 Got: %d", resp.Aggregate.Groups[0].Count)
	}
	if resp.Aggregate.Groups[1].Value != "Brute force attacks" {
		t.Fatalf("Unexpected group value. Expected: Brute force attacks Got: %s", resp.Aggregate.Groups[1].Value)
	}
	if resp.Aggregate.Groups[1].Count != 8 {
		t.Fatalf("Unexpected group count. Expected: 8 Got: %d", resp.Aggregate.Groups[1].Count)
	}
}

func TestGetFeatures(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.GetFeatures{}
	resp, err := cli.GetFeatures(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetFeatures: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if len(resp.Features) != 1 {
		t.Fatalf("Expected 1 feature, got %d", len(resp.Features))
	}
	if resp.Features[0].Enabled != false {
		t.Fatalf("Unexpected feature enabled. Expected: false Got: %t", resp.Features[0].Enabled)
	}
	if resp.Features[0].Name != "OPENVASD" {
		t.Fatalf("Unexpected feature name. Expected: OPENVASD Got: %s", resp.Features[0].Name)
	}
}

func TestEmptyTrashcan(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.EmptyTrashcan{}
	resp, err := cli.EmptyTrashcan(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during EmptyTrashcan: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
}

func TestRestore(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.Restore{ID: "97390ade-e075-11df-9973-002264764cea"}
	resp, err := cli.Restore(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during Restore: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
}

func TestRunWizard(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.RunWizard{
		Name: "quick_first_scan",
		Params: &commands.WizardParams{
			Params: []commands.WizardParam{{
				Name:  "hosts",
				Value: "localhost",
			}},
		},
	}
	resp, err := cli.RunWizard(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during RunWizard: %s", err)
	}
	if resp.Status != "202" {
		t.Fatalf("Unexpected status. Expected: 202 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK, request submitted" {
		t.Fatalf("Unexpected status text. Expected: OK, request submitted Got: %s", resp.StatusText)
	}
	if resp.Response == nil || resp.Response.StartTaskResponse == nil {
		t.Fatalf("Expected StartTaskResponse in wizard response, got: %+v", resp.Response)
	}
	if resp.Response.StartTaskResponse.ReportID != "a06d21f7-8e2f-4d7f-bceb-1df852d8b37d" {
		t.Fatalf("Unexpected report_id. Expected: a06d21f7-8e2f-4d7f-bceb-1df852d8b37d Got: %s", resp.Response.StartTaskResponse.ReportID)
	}
}

func TestRawXML(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Test a generic command
	xml := `<get_version/>`
	resp, err := cli.RawXML(xml)
	if err != nil {
		t.Fatalf("Unexpected error from RawXML: %v", err)
	}
	if !strings.Contains(resp, `status="200"`) || !strings.Contains(resp, "mocked") {
		t.Errorf("Unexpected RawXML response: %s", resp)
	}

	// Test an undocumented command
	undoc := `<undocumented_command/>`
	resp, err = cli.RawXML(undoc)
	if err != nil {
		t.Fatalf("Unexpected error from RawXML (undocumented): %v", err)
	}
	if !strings.Contains(resp, "<undocumented_command_response") || !strings.Contains(resp, "mocked") {
		t.Errorf("Unexpected RawXML response for undocumented: %s", resp)
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// TestGetSettingsConsolidated tests the GetSettings method.
func TestGetSettingsConsolidated(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()

	// Test with no filters
	resp, err := cli.GetSettings(ctx)
	if err != nil {
		t.Fatalf("Unexpected error during GetSettings: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with string filters
	resp, err = cli.GetSettings(ctx, "name~timeout")
	if err != nil {
		t.Fatalf("Unexpected error during GetSettings with string filter: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}

	// Test with functional options
	resp, err = cli.GetSettings(ctx, filtering.WithName("test"))
	if err != nil {
		t.Fatalf("Unexpected error during GetSettings with functional options: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
}

func (m *mockConn) RawXML(xmlStr string) (string, error) {
	// For testing, just echo back a canned response or the input for now.
	if strings.Contains(xmlStr, "undocumented_command") {
		return `<undocumented_command_response status="200" status_text="OK">mocked</undocumented_command_response>`, nil
	}
	return `<response status="200" status_text="OK">mocked</response>`, nil
}
