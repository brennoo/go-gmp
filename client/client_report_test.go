package client

import (
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestCreateReport(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &commands.CreateReportCommand{
		Report: &commands.ReportWrapper{
			ID:          "report-uuid",
			FormatID:    "format-uuid",
			Extension:   "xml",
			ContentType: "text/xml",
		},
		Task: &commands.CreateReportTask{
			ID: "task-uuid",
		},
		InAssets: nil,
	}
	resp, err := cli.CreateReport(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateReport: %s", err)
	}
	if resp.Status != "201" {
		t.Errorf("Expected status 201, got %s", resp.Status)
	}
	if resp.StatusText != "OK, resource created" {
		t.Errorf("Expected status text 'OK, resource created', got '%s'", resp.StatusText)
	}
	if resp.ID != "created-report-id" {
		t.Errorf("Expected ID 'created-report-id', got '%s'", resp.ID)
	}

	// Failure case
	cmdFail := &commands.CreateReportCommand{
		Report: &commands.ReportWrapper{
			ID: "wrong-id",
		},
		Task: &commands.CreateReportTask{
			ID: "wrong-task",
		},
	}
	respFail, err := cli.CreateReport(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during CreateReport (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
	if respFail.StatusText != "Bad request" {
		t.Errorf("Expected status text 'Bad request', got '%s'", respFail.StatusText)
	}
}

func TestGetReports(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &commands.GetReportsCommand{
		ReportID:    "report-uuid",
		FormatID:    "format-uuid",
		Extension:   "xml",
		ContentType: "text/xml",
	}
	resp, err := cli.GetReports(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetReports: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got '%s'", resp.StatusText)
	}
	if len(resp.Reports) != 1 {
		t.Errorf("Expected 1 report, got %d", len(resp.Reports))
	} else {
		report := resp.Reports[0]
		if report.ID != "report-uuid" {
			t.Errorf("Expected report ID 'report-uuid', got '%s'", report.ID)
		}
		if report.FormatID != "format-uuid" {
			t.Errorf("Expected format ID 'format-uuid', got '%s'", report.FormatID)
		}
		if report.Extension != "xml" {
			t.Errorf("Expected extension 'xml', got '%s'", report.Extension)
		}
		if report.ContentType != "text/xml" {
			t.Errorf("Expected content type 'text/xml', got '%s'", report.ContentType)
		}
	}

	// Failure case
	cmdFail := &commands.GetReportsCommand{
		ReportID: "wrong-id",
	}
	respFail, err := cli.GetReports(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during GetReports (fail): %s", err)
	}
	if respFail.Status != "404" {
		t.Errorf("Expected status 404, got %s", respFail.Status)
	}
	if respFail.StatusText != "Not found" {
		t.Errorf("Expected status text 'Not found', got '%s'", respFail.StatusText)
	}
}

func TestDeleteReport(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DeleteReportCommand{ReportID: "report-uuid"}
	resp, err := cli.DeleteReport(cmd)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Expected status text 'OK', got '%s'", resp.StatusText)
	}
}

func TestGetReportFormats(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &commands.GetReportFormatsCommand{ReportFormatID: "format-uuid"}
	resp, err := cli.GetReportFormats(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during GetReportFormats: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Expected status text 'OK', got '%s'", resp.StatusText)
	}
	if len(resp.ReportFormats) != 1 {
		t.Fatalf("Expected 1 report format, got %d", len(resp.ReportFormats))
	}

	// Failure case
	cmdFail := &commands.GetReportFormatsCommand{ReportFormatID: "notfound"}
	respFail, err := cli.GetReportFormats(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during GetReportFormats (fail): %s", err)
	}
	if respFail.Status != "404" {
		t.Fatalf("Expected status 404, got %s", respFail.Status)
	}
	if respFail.StatusText != "Not found" {
		t.Fatalf("Expected status text 'Not found', got '%s'", respFail.StatusText)
	}
}

func TestDeleteReportFormat(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DeleteReportFormatCommand{ReportFormatID: "format-uuid", Ultimate: "1"}
	resp, err := cli.DeleteReportFormat(cmd)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Expected status text 'OK', got '%s'", resp.StatusText)
	}
}

func TestVerifyReportFormat(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.VerifyReportFormatCommand{ReportFormatID: "format-uuid"}
	resp, err := cli.VerifyReportFormat(cmd)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Expected status text 'OK', got '%s'", resp.StatusText)
	}
}

func TestCreateReportConfig(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &commands.CreateReportConfigCommand{
		Name:         "Test config",
		ReportFormat: commands.CreateReportConfigFormat{ID: "format-uuid"},
	}
	resp, err := cli.CreateReportConfig(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateReportConfig: %s", err)
	}
	if resp.Status != "201" {
		t.Fatalf("Expected status 201, got %s", resp.Status)
	}
	if resp.StatusText != "OK, resource created" {
		t.Fatalf("Expected status text 'OK, resource created', got '%s'", resp.StatusText)
	}
	if resp.ID != "created-config-id" {
		t.Fatalf("Expected ID 'created-config-id', got '%s'", resp.ID)
	}

	// Failure case
	cmdFail := &commands.CreateReportConfigCommand{
		Name:         "Wrong config",
		ReportFormat: commands.CreateReportConfigFormat{ID: "wrong-uuid"},
	}
	respFail, err := cli.CreateReportConfig(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during CreateReportConfig (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Fatalf("Expected status 400, got %s", respFail.Status)
	}
	if respFail.StatusText != "Bad request" {
		t.Fatalf("Expected status text 'Bad request', got '%s'", respFail.StatusText)
	}
}

func TestModifyReportConfig(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case: modify name
	cmd := &commands.ModifyReportConfigCommand{
		Name: "Renamed config",
	}
	resp, err := cli.ModifyReportConfig(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyReportConfig: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Expected status text 'OK', got '%s'", resp.StatusText)
	}
	if resp.ID != "modified-config-id" {
		t.Fatalf("Expected ID 'modified-config-id', got '%s'", resp.ID)
	}

	// Success case: modify param
	cmdParam := &commands.ModifyReportConfigCommand{
		Params: []commands.ModifyReportConfigParam{{Name: "Node Distance", UseDefault: "1"}},
	}
	respParam, err := cli.ModifyReportConfig(cmdParam)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyReportConfig (param): %s", err)
	}
	if respParam.Status != "200" {
		t.Fatalf("Expected status 200, got %s", respParam.Status)
	}
	if respParam.StatusText != "OK" {
		t.Fatalf("Expected status text 'OK', got '%s'", respParam.StatusText)
	}
	if respParam.ID != "modified-config-id" {
		t.Fatalf("Expected ID 'modified-config-id', got '%s'", respParam.ID)
	}

	// Failure case
	cmdFail := &commands.ModifyReportConfigCommand{
		Name: "Invalid",
	}
	respFail, err := cli.ModifyReportConfig(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyReportConfig (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Fatalf("Expected status 400, got %s", respFail.Status)
	}
	if respFail.StatusText != "Bad request" {
		t.Fatalf("Expected status text 'Bad request', got '%s'", respFail.StatusText)
	}
}

func TestCreateReportFormat(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &commands.CreateReportFormatCommand{
		Name: "Test Format",
		Copy: "copy-uuid",
	}
	resp, err := cli.CreateReportFormat(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateReportFormat: %s", err)
	}
	if resp.Status != "201" {
		t.Errorf("Expected status 201, got %s", resp.Status)
	}
	if resp.StatusText != "OK, resource created" {
		t.Errorf("Expected status text 'OK, resource created', got '%s'", resp.StatusText)
	}
	if resp.ID != "created-format-id" {
		t.Errorf("Expected ID 'created-format-id', got '%s'", resp.ID)
	}

	// Failure case
	cmdFail := &commands.CreateReportFormatCommand{
		Name: "Wrong Format",
		Copy: "wrong-copy",
	}
	respFail, err := cli.CreateReportFormat(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during CreateReportFormat (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
	if respFail.StatusText != "Bad request" {
		t.Errorf("Expected status text 'Bad request', got '%s'", respFail.StatusText)
	}
}

func TestModifyReportFormat(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case: modify name
	cmd := &commands.ModifyReportFormatCommand{
		ReportFormatID: "format-uuid",
		Name:           "New Name",
	}
	resp, err := cli.ModifyReportFormat(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyReportFormat: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got '%s'", resp.StatusText)
	}

	// Success case: modify param
	cmdParam := &commands.ModifyReportFormatCommand{
		ReportFormatID: "format-uuid",
		Param: &commands.ModifyReportFormatParam{
			Name:  "Background Colour",
			Value: "red",
		},
	}
	respParam, err := cli.ModifyReportFormat(cmdParam)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyReportFormat (param): %s", err)
	}
	if respParam.Status != "200" {
		t.Errorf("Expected status 200, got %s", respParam.Status)
	}
	if respParam.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got '%s'", respParam.StatusText)
	}

	// Failure case
	cmdFail := &commands.ModifyReportFormatCommand{
		ReportFormatID: "wrong-uuid",
		Name:           "Invalid",
	}
	respFail, err := cli.ModifyReportFormat(cmdFail)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyReportFormat (fail): %s", err)
	}
	if respFail.Status != "400" {
		t.Errorf("Expected status 400, got %s", respFail.Status)
	}
	if respFail.StatusText != "Bad request" {
		t.Errorf("Expected status text 'Bad request', got '%s'", respFail.StatusText)
	}
}

func TestGetReportConfigs(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &commands.GetReportConfigsCommand{
		ReportConfigID: "config-uuid",
	}
	resp, err := cli.GetReportConfigs(cmd)
	if err != nil {
		t.Fatalf("GetReportConfigs failed: %v", err)
	}
	if resp.Status != "200" || resp.StatusText != "OK" || len(resp.ReportConfigs) != 1 {
		t.Errorf("unexpected response: %+v", resp)
	}
	if resp.ReportConfigs[0].ID != "config-uuid" || resp.ReportConfigs[0].Name != "Test config" {
		t.Errorf("unexpected report config: %+v", resp.ReportConfigs[0])
	}

	// Failure case
	cmd = &commands.GetReportConfigsCommand{
		ReportConfigID: "notfound",
	}
	resp, err = cli.GetReportConfigs(cmd)
	if err != nil {
		t.Fatalf("GetReportConfigs failed: %v", err)
	}
	if resp.Status != "404" || resp.StatusText != "Not found" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestDeleteReportConfig(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	// Success case
	cmd := &commands.DeleteReportConfigCommand{
		ReportConfigID: "config-uuid",
		Ultimate:       "1",
	}
	resp, err := cli.DeleteReportConfig(cmd)
	if err != nil {
		t.Fatalf("DeleteReportConfig failed: %v", err)
	}
	if resp.Status != "200" || resp.StatusText != "OK" {
		t.Errorf("unexpected response: %+v", resp)
	}

	// Failure case
	cmd = &commands.DeleteReportConfigCommand{
		ReportConfigID: "notfound",
		Ultimate:       "1",
	}
	resp, err = cli.DeleteReportConfig(cmd)
	if err != nil {
		t.Fatalf("DeleteReportConfig failed: %v", err)
	}
	if resp.Status != "404" || resp.StatusText != "Not found" {
		t.Errorf("unexpected response: %+v", resp)
	}
}
