package client

import (
	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

type mockConn struct{}

// nolint:gocyclo // This mock handles all GMP commands
// gocyclo:ignore
func (m *mockConn) Execute(command interface{}, response interface{}) error {
	if cmd, ok := command.(*gmp.AuthenticateCommand); ok {
		if cmd.Credentials.Username == "openvas" && cmd.Credentials.Password == "123" {
			(*response.(*gmp.AuthenticateResponse)).Status = "200"
		} else {
			(*response.(*gmp.AuthenticateResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetConfigsCommand); ok {
		if cmd.ConfigID == "bde773f3-2b3d-4fe6-81cb-6321ae2cc629" {
			(*response.(*gmp.GetConfigsResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetConfigsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetScannersCommand); ok {
		if cmd.ScannerID == "ee0311e7-3247-4425-bb9c-866d59f1e0e9" {
			(*response.(*gmp.GetScannersResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetScannersResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetPreferencesCommand); ok {
		if cmd.ConfigID == "4b49617e-d1d8-44b8-af81-f4675b56f837" {
			(*response.(*gmp.GetPreferencesResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetPreferencesResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreateConfigCommand); ok {
		if cmd.Name == "New Config" {
			(*response.(*gmp.CreateConfigResponse)).Status = "200"
		} else {
			(*response.(*gmp.CreateConfigResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.ModifyConfigCommand); ok {
		if cmd.Name == "Modified Config" {
			(*response.(*gmp.ModifyConfigResponse)).Status = "200"
		} else {
			(*response.(*gmp.ModifyConfigResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreateTaskCommand); ok {
		if cmd.Name == "New Task" {
			(*response.(*gmp.CreateTaskResponse)).Status = "200"
		} else {
			(*response.(*gmp.CreateTaskResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.ModifyTaskCommand); ok {
		if cmd.Comment == "Modified Task Comment" {
			(*response.(*gmp.ModifyTaskResponse)).Status = "200"
		} else {
			(*response.(*gmp.ModifyTaskResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreateTargetCommand); ok {
		if cmd.Name == "New Target" {
			(*response.(*gmp.CreateTargetResponse)).Status = "200"
		} else {
			(*response.(*gmp.CreateTargetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.ModifyTargetCommand); ok {
		if cmd.Name == "Modified Target" {
			(*response.(*gmp.ModifyTargetResponse)).Status = "200"
		} else {
			(*response.(*gmp.ModifyTargetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetTargetsCommand); ok {
		if cmd.TargetID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*gmp.GetTargetsResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetTargetsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.DeleteTargetCommand); ok {
		if cmd.TargetID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*gmp.DeleteTargetResponse)).Status = "200"
		} else {
			(*response.(*gmp.DeleteTargetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetPortListsCommand); ok {
		if cmd.PortListID == "33d0cd82-57c6-11e1-8ed1-406186ea4fc5" {
			(*response.(*gmp.GetPortListsResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetPortListsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.StartTaskCommand); ok {
		if cmd.TaskID == "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31" {
			(*response.(*gmp.StartTaskResponse)).Status = "200"
		} else {
			(*response.(*gmp.StartTaskResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetTasksCommand); ok {
		if cmd.TaskID == "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31" {
			(*response.(*gmp.GetTasksResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetTasksResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.StopTaskCommand); ok {
		if cmd.TaskID == "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31" {
			(*response.(*gmp.StopTaskResponse)).Status = "200"
		} else {
			(*response.(*gmp.StopTaskResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.DeleteTaskCommand); ok {
		if cmd.TaskID == "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31" {
			(*response.(*gmp.DeleteTaskResponse)).Status = "200"
		} else {
			(*response.(*gmp.DeleteTaskResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetResultsCommand); ok {
		if cmd.TaskID == "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31" {
			(*response.(*gmp.GetResultsResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetResultsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetVulnsCommand); ok {
		if cmd.VulnID == "1.3.6.1.4.1.25623.1.0.808160" {
			(*response.(*gmp.GetVulnsResponse)).Status = "200"
		} else {
			(*response.(*gmp.GetVulnsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreateAlertRequest); ok {
		if cmd.Name == "Test Alert" {
			(*response.(*gmp.CreateAlertResponse)).Status = "201"
			(*response.(*gmp.CreateAlertResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateAlertResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*gmp.CreateAlertResponse)).Status = "400"
		}
	}

	if _, ok := command.(*gmp.GetAlertsCommand); ok {
		(*response.(*gmp.GetAlertsResponse)).Status = "200"
		(*response.(*gmp.GetAlertsResponse)).StatusText = "OK"
	}

	if _, ok := command.(*gmp.ModifyAlertCommand); ok {
		(*response.(*gmp.ModifyAlertResponse)).Status = "200"
		(*response.(*gmp.ModifyAlertResponse)).StatusText = "OK"
	}

	if _, ok := command.(*gmp.DeleteAlertCommand); ok {
		(*response.(*gmp.DeleteAlertResponse)).Status = "200"
		(*response.(*gmp.DeleteAlertResponse)).StatusText = "OK"
	}

	if _, ok := command.(*gmp.TestAlertCommand); ok {
		(*response.(*gmp.TestAlertResponse)).Status = "200"
		(*response.(*gmp.TestAlertResponse)).StatusText = "OK"
	}

	if _, ok := command.(*gmp.ResumeTaskCommand); ok {
		(*response.(*gmp.ResumeTaskResponse)).Status = "200"
		(*response.(*gmp.ResumeTaskResponse)).StatusText = "OK"
		(*response.(*gmp.ResumeTaskResponse)).ReportID = "330ee785-c2c0-4d4c-ab96-725142c9b789"
	}

	if cmd, ok := command.(*gmp.CreateAssetCommand); ok {
		if cmd.Asset != nil && cmd.Asset.Name == "Localhost" && cmd.Asset.Type == "host" {
			(*response.(*gmp.CreateAssetResponse)).Status = "201"
			(*response.(*gmp.CreateAssetResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateAssetResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else if cmd.Report != nil && cmd.Report.ID == "report-uuid" && cmd.Report.Filter != nil && cmd.Report.Filter.Term == "min_qod=70" {
			(*response.(*gmp.CreateAssetResponse)).Status = "201"
			(*response.(*gmp.CreateAssetResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateAssetResponse)).ID = "report-asset-uuid"
		} else {
			(*response.(*gmp.CreateAssetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.ModifyAssetCommand); ok {
		if cmd.AssetID == "914b59f8-25f5-4c8f-832c-2379cd625236" && cmd.Comment == "New comment" {
			(*response.(*gmp.ModifyAssetResponse)).Status = "200"
			(*response.(*gmp.ModifyAssetResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyAssetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetAssetsCommand); ok {
		if cmd.AssetID == "b493b7a8-7489-11df-a3ec-002264764cea" {
			(*response.(*gmp.GetAssetsResponse)).Status = "200"
			(*response.(*gmp.GetAssetsResponse)).StatusText = "OK"
			(*response.(*gmp.GetAssetsResponse)).Assets = []gmp.Asset{
				{
					ID:               "b493b7a8-7489-11df-a3ec-002264764cea",
					Name:             "Localhost",
					Comment:          "",
					CreationTime:     "2018-08-29T20:21:33Z",
					ModificationTime: "2018-08-29T20:21:33Z",
					Type:             "host",
				},
			}
		} else {
			(*response.(*gmp.GetAssetsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.DeleteAssetCommand); ok {
		if cmd.AssetID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" {
			(*response.(*gmp.DeleteAssetResponse)).Status = "200"
			(*response.(*gmp.DeleteAssetResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteAssetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreateScheduleCommand); ok {
		if cmd.Name == "Monthly Scan" && cmd.Timezone == "America/New_York" && cmd.ICalendar != "" {
			(*response.(*gmp.CreateScheduleResponse)).Status = "201"
			(*response.(*gmp.CreateScheduleResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateScheduleResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*gmp.CreateScheduleResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.ModifyScheduleCommand); ok {
		if cmd.ScheduleID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" && cmd.Name == "Weekly Scan" {
			(*response.(*gmp.ModifyScheduleResponse)).Status = "200"
			(*response.(*gmp.ModifyScheduleResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyScheduleResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetSchedulesCommand); ok {
		if cmd.ScheduleID == "ddda859a-45be-4c58-85b3-517c66230232" {
			(*response.(*gmp.GetSchedulesResponse)).Status = "200"
			(*response.(*gmp.GetSchedulesResponse)).StatusText = "OK"
			(*response.(*gmp.GetSchedulesResponse)).Schedules = []gmp.Schedule{
				{
					ID:               "ddda859a-45be-4c58-85b3-517c66230232",
					Name:             "Every day",
					Comment:          "",
					CreationTime:     "2020-06-03T16:27:05Z",
					ModificationTime: "2020-06-03T16:27:05Z",
					Writable:         "1",
					InUse:            "0",
					ICalendar:        "DTSTART:20200603T162600Z DURATION:PT0S RRULE:FREQ=DAILY",
					Timezone:         "UTC",
				},
			}
		} else {
			(*response.(*gmp.GetSchedulesResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.DeleteScheduleCommand); ok {
		if cmd.ScheduleID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" {
			(*response.(*gmp.DeleteScheduleResponse)).Status = "200"
			(*response.(*gmp.DeleteScheduleResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteScheduleResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreateOverrideCommand); ok {
		if cmd.Text == "This is actually of little concern." && cmd.NVT.OID == "1.3.6.1.4.1.25623.1.0.10330" && cmd.NewThreat == "Low" && cmd.Result == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*gmp.CreateOverrideResponse)).Status = "201"
			(*response.(*gmp.CreateOverrideResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateOverrideResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*gmp.CreateOverrideResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.DeleteOverrideCommand); ok {
		if cmd.OverrideID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" {
			(*response.(*gmp.DeleteOverrideResponse)).Status = "200"
			(*response.(*gmp.DeleteOverrideResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteOverrideResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.GetOverridesCommand); ok {
		if cmd.OverrideID == "b76b81a7-9df8-42df-afff-baa9d4620128" {
			(*response.(*gmp.GetOverridesResponse)).Status = "200"
			(*response.(*gmp.GetOverridesResponse)).StatusText = "OK"
			(*response.(*gmp.GetOverridesResponse)).Overrides = []gmp.Override{
				{
					ID:          "b76b81a7-9df8-42df-afff-baa9d4620128",
					Text:        gmp.OverrideText{Text: "This is the full text of the override."},
					NewThreat:   "Log",
					NewSeverity: "0.0",
					Orphan:      true,
					Active:      true,
					NVT:         gmp.OverrideNVT{OID: "1.3.6.1.4.1.25623.1.0.75", Name: "Test NVT: long lines"},
				},
			}
		} else {
			(*response.(*gmp.GetOverridesResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.ModifyOverrideCommand); ok {
		if cmd.OverrideID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" && cmd.Text == "This issue is less important in our setup." && cmd.NewThreat == "Low" {
			(*response.(*gmp.ModifyOverrideResponse)).Status = "200"
			(*response.(*gmp.ModifyOverrideResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyOverrideResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreateReportCommand); ok {
		if cmd.Report != nil && cmd.Report.ID == "report-uuid" && cmd.Task != nil && cmd.Task.ID == "task-uuid" {
			(*response.(*gmp.CreateReportResponse)).Status = "201"
			(*response.(*gmp.CreateReportResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateReportResponse)).ID = "created-report-id"
		} else {
			(*response.(*gmp.CreateReportResponse)).Status = "400"
			(*response.(*gmp.CreateReportResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.GetReportsCommand); ok {
		if cmd.ReportID == "report-uuid" && cmd.FormatID == "format-uuid" {
			(*response.(*gmp.GetReportsResponse)).Status = "200"
			(*response.(*gmp.GetReportsResponse)).StatusText = "OK"
			(*response.(*gmp.GetReportsResponse)).Reports = []gmp.ReportWrapper{
				{
					ID:          "report-uuid",
					FormatID:    "format-uuid",
					Extension:   "xml",
					ContentType: "text/xml",
				},
			}
		} else {
			(*response.(*gmp.GetReportsResponse)).Status = "404"
			(*response.(*gmp.GetReportsResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.DeleteReportCommand); ok {
		if cmd.ReportID == "report-uuid" {
			(*response.(*gmp.DeleteReportResponse)).Status = "200"
			(*response.(*gmp.DeleteReportResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteReportResponse)).Status = "404"
			(*response.(*gmp.DeleteReportResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.CreateReportFormatCommand); ok {
		if cmd.Name == "Test Format" && cmd.Copy == "copy-uuid" {
			(*response.(*gmp.CreateReportFormatResponse)).Status = "201"
			(*response.(*gmp.CreateReportFormatResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateReportFormatResponse)).ID = "created-format-id"
		} else {
			(*response.(*gmp.CreateReportFormatResponse)).Status = "400"
			(*response.(*gmp.CreateReportFormatResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.ModifyReportFormatCommand); ok {
		if cmd.ReportFormatID == "format-uuid" && cmd.Name == "New Name" {
			(*response.(*gmp.ModifyReportFormatResponse)).Status = "200"
			(*response.(*gmp.ModifyReportFormatResponse)).StatusText = "OK"
		} else if cmd.ReportFormatID == "format-uuid" && cmd.Param != nil && cmd.Param.Name == "Background Colour" && cmd.Param.Value == "red" {
			(*response.(*gmp.ModifyReportFormatResponse)).Status = "200"
			(*response.(*gmp.ModifyReportFormatResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyReportFormatResponse)).Status = "400"
			(*response.(*gmp.ModifyReportFormatResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.GetReportFormatsCommand); ok {
		if cmd.ReportFormatID == "format-uuid" {
			(*response.(*gmp.GetReportFormatsResponse)).Status = "200"
			(*response.(*gmp.GetReportFormatsResponse)).StatusText = "OK"
			(*response.(*gmp.GetReportFormatsResponse)).ReportFormats = []gmp.ReportFormatWrapper{
				{
					ID:               "format-uuid",
					Name:             "HTML",
					Extension:        "html",
					ContentType:      "text/html",
					Summary:          "Single page HTML report.",
					Description:      "A single HTML page listing results of a scan.",
					CreationTime:     "2013-01-31T16:46:32+01:00",
					ModificationTime: "2013-01-31T16:46:32+01:00",
					Writable:         "1",
					InUse:            "0",
					Active:           "1",
				},
			}
		} else {
			(*response.(*gmp.GetReportFormatsResponse)).Status = "404"
			(*response.(*gmp.GetReportFormatsResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.DeleteReportFormatCommand); ok {
		if cmd.ReportFormatID == "format-uuid" && cmd.Ultimate == "1" {
			(*response.(*gmp.DeleteReportFormatResponse)).Status = "200"
			(*response.(*gmp.DeleteReportFormatResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteReportFormatResponse)).Status = "404"
			(*response.(*gmp.DeleteReportFormatResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.VerifyReportFormatCommand); ok {
		if cmd.ReportFormatID == "format-uuid" {
			(*response.(*gmp.VerifyReportFormatResponse)).Status = "200"
			(*response.(*gmp.VerifyReportFormatResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.VerifyReportFormatResponse)).Status = "404"
			(*response.(*gmp.VerifyReportFormatResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.CreateReportConfigCommand); ok {
		if cmd.Name == "Test config" && cmd.ReportFormat.ID == "format-uuid" {
			(*response.(*gmp.CreateReportConfigResponse)).Status = "201"
			(*response.(*gmp.CreateReportConfigResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateReportConfigResponse)).ID = "created-config-id"
		} else {
			(*response.(*gmp.CreateReportConfigResponse)).Status = "400"
			(*response.(*gmp.CreateReportConfigResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.ModifyReportConfigCommand); ok {
		if len(cmd.Params) > 0 && cmd.Params[0].Name == "Node Distance" && cmd.Params[0].UseDefault == "1" {
			(*response.(*gmp.ModifyReportConfigResponse)).Status = "200"
			(*response.(*gmp.ModifyReportConfigResponse)).StatusText = "OK"
			(*response.(*gmp.ModifyReportConfigResponse)).ID = "modified-config-id"
		} else if cmd.Name == "Renamed config" {
			(*response.(*gmp.ModifyReportConfigResponse)).Status = "200"
			(*response.(*gmp.ModifyReportConfigResponse)).StatusText = "OK"
			(*response.(*gmp.ModifyReportConfigResponse)).ID = "modified-config-id"
		} else {
			(*response.(*gmp.ModifyReportConfigResponse)).Status = "400"
			(*response.(*gmp.ModifyReportConfigResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.GetReportConfigsCommand); ok {
		if cmd.ReportConfigID == "config-uuid" {
			(*response.(*gmp.GetReportConfigsResponse)).Status = "200"
			(*response.(*gmp.GetReportConfigsResponse)).StatusText = "OK"
			(*response.(*gmp.GetReportConfigsResponse)).ReportConfigs = []gmp.ReportConfigWrapper{
				{
					ID:               "config-uuid",
					Name:             "Test config",
					Owner:            &gmp.ReportConfigOwner{Name: "admin"},
					Comment:          "Test comment",
					CreationTime:     "2024-01-23T09:43:03Z",
					ModificationTime: "2024-01-26T14:11:54Z",
					Writable:         "1",
					InUse:            "0",
					Permissions: &gmp.ReportConfigPermissions{
						Permissions: []gmp.ReportConfigPermission{{Name: "Everything"}},
					},
					ReportFormat: &gmp.ReportConfigFormat{
						ID:   "format-uuid",
						Name: "Topology SVG",
					},
					Params: []gmp.ReportConfigParam{
						{
							UsingDefault: "0",
							Name:         "Graph Type",
							Type:         "selection",
							Value:        "dot",
							Default:      "twopi",
							Options:      []string{"circo", "dot", "twopi"},
						},
						{
							UsingDefault: "1",
							Name:         "Node Distance",
							Type:         "integer<min>1</min><max>20</max>",
							Value:        "8",
							Default:      "8",
						},
					},
				},
			}
		} else {
			(*response.(*gmp.GetReportConfigsResponse)).Status = "404"
			(*response.(*gmp.GetReportConfigsResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.DeleteReportConfigCommand); ok {
		if cmd.ReportConfigID == "config-uuid" && cmd.Ultimate == "1" {
			(*response.(*gmp.DeleteReportConfigResponse)).Status = "200"
			(*response.(*gmp.DeleteReportConfigResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteReportConfigResponse)).Status = "404"
			(*response.(*gmp.DeleteReportConfigResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.GetSystemReportsCommand); ok {
		if cmd.Name == "proc" {
			(*response.(*gmp.GetSystemReportsResponse)).Status = "200"
			(*response.(*gmp.GetSystemReportsResponse)).StatusText = "OK"
			(*response.(*gmp.GetSystemReportsResponse)).SystemReports = []gmp.SystemReportWrapper{
				{
					Name:  "proc",
					Title: "Processes",
					Report: &gmp.SystemReportContent{
						Format:   "png",
						Duration: "86400",
						Value:    "iVBORw0KGgoAAAANSUhEUgAAArkAAAE...2bEdAAAAAElFTkSuQmCC",
					},
				},
			}
		} else if cmd.Brief == "1" {
			(*response.(*gmp.GetSystemReportsResponse)).Status = "200"
			(*response.(*gmp.GetSystemReportsResponse)).StatusText = "OK"
			(*response.(*gmp.GetSystemReportsResponse)).SystemReports = []gmp.SystemReportWrapper{
				{Name: "proc", Title: "Processes"},
				{Name: "load", Title: "System Load"},
			}
		} else {
			(*response.(*gmp.GetSystemReportsResponse)).Status = "404"
			(*response.(*gmp.GetSystemReportsResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.CreateCredentialCommand); ok {
		if cmd.Name == "Test Credential" {
			(*response.(*gmp.CreateCredentialResponse)).Status = "201"
			(*response.(*gmp.CreateCredentialResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateCredentialResponse)).ID = "created-credential-id"
		} else {
			(*response.(*gmp.CreateCredentialResponse)).Status = "400"
			(*response.(*gmp.CreateCredentialResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.ModifyCredentialCommand); ok {
		if cmd.CredentialID == "cred-uuid" {
			(*response.(*gmp.ModifyCredentialResponse)).Status = "200"
			(*response.(*gmp.ModifyCredentialResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyCredentialResponse)).Status = "400"
			(*response.(*gmp.ModifyCredentialResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.GetCredentialsCommand); ok {
		if cmd.CredentialID == "" {
			(*response.(*gmp.GetCredentialsResponse)).Status = "200"
			(*response.(*gmp.GetCredentialsResponse)).StatusText = "OK"
			(*response.(*gmp.GetCredentialsResponse)).Credentials = []gmp.CredentialWrapper{
				{
					ID:       "cred-uuid-1",
					Name:     "sally",
					Login:    "sally",
					Writable: "1",
					InUse:    "0",
					Type:     "usk",
					FullType: "username + SSH key",
					Formats:  &gmp.CredentialFormats{Formats: []string{"key", "rpm", "deb"}},
				},
				{
					ID:       "cred-uuid-2",
					Name:     "bob",
					Login:    "bob",
					Writable: "1",
					InUse:    "1",
					Type:     "up",
					FullType: "username + password",
					Formats:  &gmp.CredentialFormats{Formats: []string{"exe"}},
				},
			}
		} else if cmd.CredentialID == "cred-uuid-1" {
			(*response.(*gmp.GetCredentialsResponse)).Status = "200"
			(*response.(*gmp.GetCredentialsResponse)).StatusText = "OK"
			(*response.(*gmp.GetCredentialsResponse)).Credentials = []gmp.CredentialWrapper{
				{
					ID:       "cred-uuid-1",
					Name:     "sally",
					Login:    "sally",
					Writable: "1",
					InUse:    "0",
					Type:     "usk",
					FullType: "username + SSH key",
					Formats:  &gmp.CredentialFormats{Formats: []string{"key", "rpm", "deb"}},
				},
			}
		} else {
			(*response.(*gmp.GetCredentialsResponse)).Status = "404"
			(*response.(*gmp.GetCredentialsResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.DeleteCredentialCommand); ok {
		if cmd.CredentialID == "cred-uuid-1" && cmd.Ultimate == "1" {
			(*response.(*gmp.DeleteCredentialResponse)).Status = "200"
			(*response.(*gmp.DeleteCredentialResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteCredentialResponse)).Status = "404"
			(*response.(*gmp.DeleteCredentialResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.CreateScannerCommand); ok {
		if cmd.Name == "Default Scanner" && cmd.Host == "localhost" && cmd.Port == "9391" && cmd.Type == "2" && cmd.CAPub != "" && cmd.Credential != nil && cmd.Credential.ID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*gmp.CreateScannerResponse)).Status = "201"
			(*response.(*gmp.CreateScannerResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateScannerResponse)).ID = "814cd30f-dee1-4d58-851d-21b8d0c048e3"
		} else {
			(*response.(*gmp.CreateScannerResponse)).Status = "400"
			(*response.(*gmp.CreateScannerResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.ModifyScannerCommand); ok {
		if cmd.ScannerID == "scanner-uuid" && cmd.Name == "Updated Scanner" {
			(*response.(*gmp.ModifyScannerResponse)).Status = "200"
			(*response.(*gmp.ModifyScannerResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyScannerResponse)).Status = "400"
			(*response.(*gmp.ModifyScannerResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.DeleteScannerCommand); ok {
		if cmd.ScannerID == "scanner-uuid" && cmd.Ultimate == "1" {
			(*response.(*gmp.DeleteScannerResponse)).Status = "200"
			(*response.(*gmp.DeleteScannerResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteScannerResponse)).Status = "404"
			(*response.(*gmp.DeleteScannerResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.VerifyScannerCommand); ok {
		if cmd.ScannerID == "scanner-uuid" {
			(*response.(*gmp.VerifyScannerResponse)).Status = "200"
			(*response.(*gmp.VerifyScannerResponse)).StatusText = "OK"
			(*response.(*gmp.VerifyScannerResponse)).Version = "OTP/2.0"
		} else {
			(*response.(*gmp.VerifyScannerResponse)).Status = "404"
			(*response.(*gmp.VerifyScannerResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.CreatePortListCommand); ok {
		if cmd.Name == "All TCP" && cmd.PortRange == "T:1-65535" {
			(*response.(*gmp.CreatePortListResponse)).Status = "201"
			(*response.(*gmp.CreatePortListResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreatePortListResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*gmp.CreatePortListResponse)).Status = "400"
			(*response.(*gmp.CreatePortListResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.ModifyPortListCommand); ok {
		if cmd.PortListID == "27140836-05ae-4e8b-9abf-f725ddc2888f" {
			(*response.(*gmp.ModifyPortListResponse)).Status = "200"
			(*response.(*gmp.ModifyPortListResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyPortListResponse)).Status = "400"
			(*response.(*gmp.ModifyPortListResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.DeletePortListCommand); ok {
		if cmd.PortListID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" && cmd.Ultimate == "1" {
			(*response.(*gmp.DeletePortListResponse)).Status = "200"
			(*response.(*gmp.DeletePortListResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeletePortListResponse)).Status = "404"
			(*response.(*gmp.DeletePortListResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.CreatePortRangeCommand); ok {
		if cmd.PortList.ID == "354cd3ef-bbe1-4d58-859d-21b8d0c046c4" && cmd.Start == "777" && cmd.End == "779" && cmd.Type == "TCP" {
			(*response.(*gmp.CreatePortRangeResponse)).Status = "201"
			(*response.(*gmp.CreatePortRangeResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreatePortRangeResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*gmp.CreatePortRangeResponse)).Status = "400"
			(*response.(*gmp.CreatePortRangeResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.DeletePortRangeCommand); ok {
		if cmd.PortRangeID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" {
			(*response.(*gmp.DeletePortRangeResponse)).Status = "200"
			(*response.(*gmp.DeletePortRangeResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeletePortRangeResponse)).Status = "404"
			(*response.(*gmp.DeletePortRangeResponse)).StatusText = "Not found"
		}
	}

	if _, ok := command.(*gmp.DescribeAuthCommand); ok {
		resp := response.(*gmp.DescribeAuthResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Groups = []gmp.DescribeAuthGroup{
			{
				Name: "method:file",
				Settings: []gmp.DescribeAuthSetting{
					{Key: "enable", Value: "true"},
					{Key: "order", Value: "1"},
				},
			},
		}
	}

	if _, ok := command.(*gmp.GetVersionCommand); ok {
		resp := response.(*gmp.GetVersionResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Version = "1.0"
	}

	if _, ok := command.(*gmp.HelpCommand); ok {
		resp := response.(*gmp.HelpResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Text = "AUTHENTICATE           Authenticate with the manager.\nCREATE_ALERT           Create an alert.\n...\nVERIFY_SCANNER         Verify a scanner."
	}

	if cmd, ok := command.(*gmp.GetInfoCommand); ok {
		if cmd.Name == "CVE-2011-0018" && cmd.Type == "cve" {
			resp := response.(*gmp.GetInfoResponse)
			resp.Status = "200"
			resp.StatusText = "OK"
			resp.Infos = []gmp.GetInfoInfo{{ID: "CVE-2011-0018", Name: "CVE-2011-0018"}}
		} else {
			resp := response.(*gmp.GetInfoResponse)
			resp.Status = "404"
			resp.StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.DeleteConfigCommand); ok {
		if cmd.ConfigID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" && cmd.Ultimate == "1" {
			(*response.(*gmp.DeleteConfigResponse)).Status = "200"
			(*response.(*gmp.DeleteConfigResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteConfigResponse)).Status = "404"
			(*response.(*gmp.DeleteConfigResponse)).StatusText = "Not found"
		}
	}

	if _, ok := command.(*gmp.SyncConfigCommand); ok {
		(*response.(*gmp.SyncConfigResponse)).Status = "202"
		(*response.(*gmp.SyncConfigResponse)).StatusText = "OK, request submitted"
	}

	if cmd, ok := command.(*gmp.MoveTaskCommand); ok {
		if cmd.TaskID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" && cmd.SlaveID == "97390ade-e075-11df-9973-002264764cea" {
			(*response.(*gmp.MoveTaskResponse)).Status = "200"
			(*response.(*gmp.MoveTaskResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.MoveTaskResponse)).Status = "404"
			(*response.(*gmp.MoveTaskResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.CreateTLSCertificateCommand); ok {
		if cmd.Name == "Example Certificate" && cmd.Certificate == "MIIDNjCCAp+gAwIBAgIBATANBgkqhkiG9w0BAQQFADCBqTELM[...]" {
			(*response.(*gmp.CreateTLSCertificateResponse)).Status = "201"
			(*response.(*gmp.CreateTLSCertificateResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateTLSCertificateResponse)).ID = "8a925978-59d0-494b-a837-40b271569727"
		} else {
			(*response.(*gmp.CreateTLSCertificateResponse)).Status = "400"
			(*response.(*gmp.CreateTLSCertificateResponse)).StatusText = "Bad request"
			(*response.(*gmp.CreateTLSCertificateResponse)).ID = ""
		}
	}

	if cmd, ok := command.(*gmp.ModifyTLSCertificateCommand); ok {
		if cmd.TLSCertificateID == "8a925978-59d0-494b-a837-40b271569727" && cmd.Name == "Renamed Example Certificate" {
			(*response.(*gmp.ModifyTLSCertificateResponse)).Status = "200"
			(*response.(*gmp.ModifyTLSCertificateResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyTLSCertificateResponse)).Status = "404"
			(*response.(*gmp.ModifyTLSCertificateResponse)).StatusText = "Not found"
		}
	}

	if _, ok := command.(*gmp.GetTLSCertificatesCommand); ok {
		resp := response.(*gmp.GetTLSCertificatesResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.TLSCertificates = []gmp.TLSCertificate{
			{
				ID:    "ba36ed15-92fa-4ae0-af53-bad8ce472f18",
				Owner: &gmp.Owner{Name: "admin"},
				Name:  "Example Certificate",
			},
		}
	}

	if cmd, ok := command.(*gmp.ModifyAgentsCommand); ok {
		if len(cmd.Agents) == 1 && cmd.Agents[0].ID == "fb6451bf-ec5a-45a8-8bab-5cf4b862e51b" &&
			cmd.Authorized == "1" && cmd.MinInterval == "1000" && cmd.HeartbeatInterval == "0" &&
			cmd.Schedule == "@every 12h" && cmd.Comment == "example update" {
			(*response.(*gmp.ModifyAgentsResponse)).Status = "200"
			(*response.(*gmp.ModifyAgentsResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyAgentsResponse)).Status = "400"
			(*response.(*gmp.ModifyAgentsResponse)).StatusText = "Bad request"
		}
	}

	if _, ok := command.(*gmp.GetAgentsCommand); ok {
		resp := response.(*gmp.GetAgentsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Agents = []gmp.Agent{
			{
				ID:      "62462fe0-5834-4630-afc2-0d040c63487c",
				AgentID: "GAT-29-p0MPX0FT",
			},
		}
	}

	if cmd, ok := command.(*gmp.DeleteAgentsCommand); ok {
		if len(cmd.Agents) == 1 && cmd.Agents[0].ID == "c6f1cdc3-8c2c-4b2e-9f43-139d23c7cfd4" {
			(*response.(*gmp.DeleteAgentsResponse)).Status = "200"
			(*response.(*gmp.DeleteAgentsResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteAgentsResponse)).Status = "404"
			(*response.(*gmp.DeleteAgentsResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*gmp.GetNvtsCommand); ok {
		if cmd.Details == "1" {
			(*response.(*gmp.GetNvtsResponse)).Status = "200"
			(*response.(*gmp.GetNvtsResponse)).StatusText = "OK"
			(*response.(*gmp.GetNvtsResponse)).Nvts = []gmp.NVT{
				{
					OID:              "1.3.6.1.4.1.25623.1.7.13005",
					Name:             "Services",
					CreationTime:     "2011-01-14T10:12:23+01:00",
					ModificationTime: "2012-09-19T20:56:15+02:00",
					Category:         "3",
					Family:           "Service detection",
					CvssBase:         "",
					Severities:       "0",
					Refs:             "",
					Tags:             "NOTAG",
					PreferenceCount:  "-1",
					Timeout:          "",
				},
				{
					OID:              "1.3.6.1.4.1.25623.1.7.13006",
					Name:             "FooBar 21",
					CreationTime:     "2011-01-14T10:12:23+01:00",
					ModificationTime: "2012-09-19T20:56:15+02:00",
					Category:         "3",
					Family:           "Service detection",
				},
			}
		} else if cmd.NvtOID == "1.3.6.1.4.1.25623.1.0.10330" {
			(*response.(*gmp.GetNvtsResponse)).Status = "200"
			(*response.(*gmp.GetNvtsResponse)).StatusText = "OK"
			(*response.(*gmp.GetNvtsResponse)).Nvts = []gmp.NVT{
				{
					OID:  "1.3.6.1.4.1.25623.1.0.10330",
					Name: "Services",
				},
			}
		} else {
			(*response.(*gmp.GetNvtsResponse)).Status = "404"
			(*response.(*gmp.GetNvtsResponse)).StatusText = "Not found"
		}
	}

	if _, ok := command.(*gmp.GetNvtFamiliesCommand); ok {
		(*response.(*gmp.GetNvtFamiliesResponse)).Status = "200"
		(*response.(*gmp.GetNvtFamiliesResponse)).StatusText = "OK"
		(*response.(*gmp.GetNvtFamiliesResponse)).Families = []gmp.Family{
			{
				Name:        "Credentials",
				MaxNvtCount: 8,
			},
			{
				Name:        "Port scanners",
				MaxNvtCount: 12,
			},
		}
	}

	if _, ok := command.(*gmp.GetFeedsCommand); ok {
		resp := response.(*gmp.GetFeedsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.FeedOwnerSet = "1"
		resp.FeedRolesSet = "1"
		resp.FeedResourcesAccess = "1"
		resp.Feeds = []gmp.Feed{
			{
				Type:        "NVT",
				Name:        "Greenbone Security Feed",
				Version:     "201608180124",
				Description: "This script synchronizes an NVT collection with...",
			},
			{
				Type:        "CERT",
				Name:        "Greenbone CERT Feed",
				Version:     "201609130000",
				Description: "This script synchronizes a CERT collection with...",
			},
			{
				Type:        "SCAP",
				Name:        "Greenbone SCAP Feed",
				Version:     "201608172300",
				Description: "This script synchronizes a SCAP collection with...",
			},
		}
	}

	if cmd, ok := command.(*gmp.ModifyLicenseCommand); ok {
		if cmd.File != "" {
			(*response.(*gmp.ModifyLicenseResponse)).Status = "200"
			(*response.(*gmp.ModifyLicenseResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyLicenseResponse)).Status = "400"
			(*response.(*gmp.ModifyLicenseResponse)).StatusText = "Missing file"
		}
	}

	if _, ok := command.(*gmp.GetLicenseCommand); ok {
		resp := response.(*gmp.GetLicenseResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.License = &gmp.License{
			Status: "active",
			Content: &gmp.LicenseContent{
				Meta: &gmp.LicenseMeta{
					ID:           "4711",
					Version:      "1.0.0",
					Comment:      "Test License",
					Type:         "trial",
					CustomerName: "Jane Doe",
					Created:      "2021-08-27T06:05:21Z",
					Begins:       "2021-08-27T07:05:21Z",
					Expires:      "2021-09-04T07:05:21Z",
				},
				Appliance: &gmp.LicenseAppliance{
					Model:     "trial",
					ModelType: "virtual",
					Sensor:    "0",
				},
				Keys: &gmp.LicenseKeys{
					Keys: []gmp.LicenseKey{{Name: "feed", Value: "*base64 GSF key*"}},
				},
				Signatures: &gmp.LicenseSignatures{
					Signatures: []gmp.LicenseSignature{{Name: "license", Value: "*base64 signature*"}},
				},
			},
		}
	}

	if cmd, ok := command.(*gmp.ModifySettingCommand); ok {
		if cmd.Name == "Timezone" && cmd.Value == "QWZyaWNhL0pvaGFubmVzYnVyZw==" {
			(*response.(*gmp.ModifySettingResponse)).Status = "200"
			(*response.(*gmp.ModifySettingResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifySettingResponse)).Status = "400"
			(*response.(*gmp.ModifySettingResponse)).StatusText = "Bad request"
		}
	}

	if _, ok := command.(*gmp.GetSettingsCommand); ok {
		resp := response.(*gmp.GetSettingsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Settings = &gmp.SettingsBlock{
			Setting: []gmp.Setting{{
				ID:    "5f5a8712-8017-11e1-8556-406186ea4fc5",
				Name:  "Rows Per Page",
				Value: "15",
			}},
		}
	}

	if _, ok := command.(*gmp.GetResourceNamesCommand); ok {
		resp := response.(*gmp.GetResourceNamesResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Type = "os"
		resp.Resources = []gmp.ResourceName{
			{ID: "5b6b6aef-b320-42ca-899f-3161ee2a4195", Name: "cpe:/h:hp:jetdirect"},
			{ID: "5be25864-9249-431e-8a91-039e334371ad", Name: "cpe:/o:canonical:ubuntu_linux:18.04"},
		}
	}

	if _, ok := command.(*gmp.GetAggregatesCommand); ok {
		resp := response.(*gmp.GetAggregatesResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Aggregate = &gmp.AggregateBlock{
			DataType:    "nvt",
			DataColumn:  "severity",
			GroupColumn: "family",
			Groups: []gmp.AggregateGroup{
				{Value: "AIX Local Security Checks", Count: 1, CCount: 1, Min: 3.3, Max: 3.3, Mean: 3.3, Sum: 3.3, CSum: 3.3},
				{Value: "Brute force attacks", Count: 8, CCount: 9, Min: 0, Max: 7.8, Mean: 6.275, Sum: 50.2, CSum: 53.5},
			},
		}
	}

	if _, ok := command.(*gmp.GetFeaturesCommand); ok {
		resp := response.(*gmp.GetFeaturesResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Features = []gmp.FeatureInfo{{
			Enabled: "0",
			Name:    "OPENVASD",
		}}
	}

	if _, ok := command.(*gmp.EmptyTrashcanCommand); ok {
		resp := response.(*gmp.EmptyTrashcanResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
	}

	if _, ok := command.(*gmp.RestoreCommand); ok {
		resp := response.(*gmp.RestoreResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
	}

	if _, ok := command.(*gmp.RunWizardCommand); ok {
		resp := response.(*gmp.RunWizardResponse)
		resp.Status = "202"
		resp.StatusText = "OK, request submitted"
		resp.Response = &gmp.WizardResponse{
			StartTaskResponse: &gmp.StartTaskResponse{
				Status:     "202",
				StatusText: "OK, request submitted",
				ReportID:   "a06d21f7-8e2f-4d7f-bceb-1df852d8b37d",
			},
		}
	}

	if cmd, ok := command.(*gmp.CreateUserCommand); ok {
		if cmd.Name == "testuser" && cmd.Password == "testpass" && len(cmd.Roles) > 0 && cmd.Roles[0].ID == "role-uuid" {
			(*response.(*gmp.CreateUserResponse)).Status = "201"
			(*response.(*gmp.CreateUserResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateUserResponse)).ID = "created-user-id"
		} else {
			(*response.(*gmp.CreateUserResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.ModifyUserCommand); ok {
		if cmd.UserID == "user-uuid" && cmd.Name == "testuser" && cmd.NewName == "newuser" && cmd.Comment == "Updated user" && cmd.Password != nil && cmd.Password.Modify == "1" && cmd.Password.Text == "newpass" && len(cmd.Roles) > 0 && cmd.Roles[0].ID == "role-uuid" {
			(*response.(*gmp.ModifyUserResponse)).Status = "200"
			(*response.(*gmp.ModifyUserResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyUserResponse)).Status = "400"
		}
	}
	if _, ok := command.(*gmp.GetUsersCommand); ok {
		(*response.(*gmp.GetUsersResponse)).Status = "200"
		(*response.(*gmp.GetUsersResponse)).StatusText = "OK"
		(*response.(*gmp.GetUsersResponse)).Users = []gmp.UserWrapper{
			{
				Name:    "foobar",
				Role:    &gmp.UserRole{ID: "8d453140-b74d-11e2-b0be-406186ea4fc5", Name: "User"},
				Hosts:   &gmp.CreateUserHosts{Allow: "2"},
				Sources: &gmp.CreateUserSources{Source: &gmp.CreateUserSource{Type: "file"}},
			},
		}
	}

	if cmd, ok := command.(*gmp.DeleteUserCommand); ok {
		if cmd.Name == "foobar" {
			(*response.(*gmp.DeleteUserResponse)).Status = "200"
			(*response.(*gmp.DeleteUserResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteUserResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreateGroupCommand); ok {
		if cmd.Name == "Managers" {
			(*response.(*gmp.CreateGroupResponse)).Status = "201"
			(*response.(*gmp.CreateGroupResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateGroupResponse)).ID = "d94211b6-ba40-11e3-bcb1-406186ea4fc5"
		} else {
			(*response.(*gmp.CreateGroupResponse)).Status = "400"
			(*response.(*gmp.CreateGroupResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.CreateRoleCommand); ok {
		if cmd.Name == "Test Role" && cmd.Comment == "A test role" {
			(*response.(*gmp.CreateRoleResponse)).Status = "201"
			(*response.(*gmp.CreateRoleResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateRoleResponse)).ID = "created-role-id"
		} else {
			(*response.(*gmp.CreateRoleResponse)).Status = "400"
			(*response.(*gmp.CreateRoleResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.ModifyRoleCommand); ok {
		if cmd.RoleID == "role-uuid" && cmd.Name == "Updated Role" && cmd.Comment == "Updated comment" && cmd.Users == "user1,user2" {
			(*response.(*gmp.ModifyRoleResponse)).Status = "200"
			(*response.(*gmp.ModifyRoleResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyRoleResponse)).Status = "400"
			(*response.(*gmp.ModifyRoleResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.GetRolesCommand); ok {
		if mockFail, ok := any(cmd).(interface{ Fail() bool }); ok && mockFail.Fail() {
			(*response.(*gmp.GetRolesResponse)).Status = "400"
			(*response.(*gmp.GetRolesResponse)).StatusText = "Bad request"
			(*response.(*gmp.GetRolesResponse)).Roles = nil
		} else {
			(*response.(*gmp.GetRolesResponse)).Status = "200"
			(*response.(*gmp.GetRolesResponse)).StatusText = "OK"
			(*response.(*gmp.GetRolesResponse)).Roles = []gmp.Role{
				{
					ID:               "b493b7a8-7489-11df-a3ec-002264764cea",
					Name:             "Management",
					Comment:          "Managers",
					CreationTime:     "2018-08-29T20:21:33Z",
					ModificationTime: "2018-08-29T20:21:33Z",
					Writable:         "1",
					InUse:            "0",
					Users:            "sarah, frank",
				},
			}
		}
	}

	if cmd, ok := command.(*gmp.DeleteRoleCommand); ok {
		if cmd.RoleID == "b64c81b2-b9de-11e3-a2e9-406186ea4fc5" && cmd.Ultimate == "1" {
			(*response.(*gmp.DeleteRoleResponse)).Status = "200"
			(*response.(*gmp.DeleteRoleResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteRoleResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.CreatePermissionCommand); ok {
		if cmd.Name == "get_targets" && cmd.Subject != nil && cmd.Subject.ID == "66abe5ce-c011-11e3-b96e-406186ea4fc5" && cmd.Subject.Type == "user" && cmd.Resource != nil && cmd.Resource.ID == "b493b7a8-7489-11df-a3ec-002264764cea" {
			(*response.(*gmp.CreatePermissionResponse)).Status = "201"
			(*response.(*gmp.CreatePermissionResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreatePermissionResponse)).ID = "created-permission-id"
		} else {
			(*response.(*gmp.CreatePermissionResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*gmp.ModifyPermissionCommand); ok {
		if cmd.PermissionID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" && cmd.Subject != nil && cmd.Subject.ID == "76e47468-c095-11e3-9285-406186ea4fc5" && cmd.Subject.Type == "user" {
			(*response.(*gmp.ModifyPermissionResponse)).Status = "200"
			(*response.(*gmp.ModifyPermissionResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyPermissionResponse)).Status = "400"
		}
	}

	if _, ok := command.(*gmp.GetPermissionsCommand); ok {
		(*response.(*gmp.GetPermissionsResponse)).Status = "200"
		(*response.(*gmp.GetPermissionsResponse)).StatusText = "OK"
		(*response.(*gmp.GetPermissionsResponse)).Permissions = []gmp.PermissionInfo{
			{
				ID:               "b493b7a8-7489-11df-a3ec-002264764cea",
				Name:             "Management",
				Comment:          "Managers",
				CreationTime:     "2018-08-29T20:21:33Z",
				ModificationTime: "2018-08-29T20:21:33Z",
				Writable:         "1",
				InUse:            "0",
				Users:            "sarah, frank",
			},
		}
	}

	if cmd, ok := command.(*gmp.DeletePermissionCommand); ok {
		if cmd.PermissionID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" {
			(*response.(*gmp.DeletePermissionResponse)).Status = "200"
			(*response.(*gmp.DeletePermissionResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeletePermissionResponse)).Status = "400"
			(*response.(*gmp.DeletePermissionResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.ModifyGroupCommand); ok {
		if cmd.GroupID == "d94211b6-ba40-11e3-bcb1-406186ea4fc5" && cmd.Name == "Line Managers" {
			(*response.(*gmp.ModifyGroupResponse)).Status = "200"
			(*response.(*gmp.ModifyGroupResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyGroupResponse)).Status = "400"
			(*response.(*gmp.ModifyGroupResponse)).StatusText = "Bad request"
		}
	}

	if _, ok := command.(*gmp.GetGroupsCommand); ok {
		(*response.(*gmp.GetGroupsResponse)).Status = "200"
		(*response.(*gmp.GetGroupsResponse)).StatusText = "OK"
		(*response.(*gmp.GetGroupsResponse)).Groups = []gmp.GroupInfo{
			{
				ID:               "b493b7a8-7489-11df-a3ec-002264764cea",
				Name:             "Management",
				Comment:          "Managers",
				CreationTime:     "2018-08-29T20:21:33Z",
				ModificationTime: "2018-08-29T20:21:33Z",
				Writable:         "1",
				InUse:            "0",
				Users:            "sarah, frank",
			},
		}
	}

	if cmd, ok := command.(*gmp.DeleteGroupCommand); ok {
		if cmd.GroupID == "d94211b6-ba40-11e3-bcb1-406186ea4fc5" {
			(*response.(*gmp.DeleteGroupResponse)).Status = "200"
			(*response.(*gmp.DeleteGroupResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteGroupResponse)).Status = "400"
			(*response.(*gmp.DeleteGroupResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.CreateTagCommand); ok {
		if cmd.Name == "geo:long" && cmd.Resources != nil && cmd.Resources.Resource != nil && cmd.Resources.Resource.ID == "b493b7a8-7489-11df-a3ec-002264764cea" {
			(*response.(*gmp.CreateTagResponse)).Status = "201"
			(*response.(*gmp.CreateTagResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateTagResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*gmp.CreateTagResponse)).Status = "400"
			(*response.(*gmp.CreateTagResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.ModifyTagCommand); ok {
		if cmd.TagID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" && cmd.Active == "0" {
			(*response.(*gmp.ModifyTagResponse)).Status = "200"
			(*response.(*gmp.ModifyTagResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyTagResponse)).Status = "400"
			(*response.(*gmp.ModifyTagResponse)).StatusText = "Bad request"
		}
	}

	if _, ok := command.(*gmp.GetTagsCommand); ok {
		(*response.(*gmp.GetTagsResponse)).Status = "200"
		(*response.(*gmp.GetTagsResponse)).StatusText = "OK"
		(*response.(*gmp.GetTagsResponse)).Tags = []gmp.TagInfo{
			{
				ID:               "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
				Name:             "geo:long",
				Comment:          "",
				CreationTime:     "2016-03-03T11:46:56Z",
				ModificationTime: "2016-03-03T11:46:56Z",
				Writable:         "1",
				InUse:            "0",
				Value:            "52.2788",
				Active:           "1",
				Resources: []gmp.TagResource{
					{
						ID:    "b493b7a8-7489-11df-a3ec-002264764cea",
						Type:  "target",
						Name:  "Server 1",
						Trash: "0",
					},
				},
			},
		}
	}

	if cmd, ok := command.(*gmp.DeleteTagCommand); ok {
		if cmd.TagID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*gmp.DeleteTagResponse)).Status = "200"
			(*response.(*gmp.DeleteTagResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.DeleteTagResponse)).Status = "400"
			(*response.(*gmp.DeleteTagResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.CreateNoteCommand); ok {
		if cmd.Text == "This issue should be resolved after the upgrade." && cmd.NVT != nil && cmd.NVT.OID == "1.3.6.1.4.1.25623.1.0.10330" && cmd.Result != nil && cmd.Result.ID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*gmp.CreateNoteResponse)).Status = "202"
			(*response.(*gmp.CreateNoteResponse)).StatusText = "OK, resource created"
			(*response.(*gmp.CreateNoteResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*gmp.CreateNoteResponse)).Status = "400"
			(*response.(*gmp.CreateNoteResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*gmp.ModifyNoteCommand); ok {
		if cmd.NoteID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" && cmd.Text == "This issue should be resolved after the upgrade." && cmd.Result != nil && cmd.Result.ID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*gmp.ModifyNoteResponse)).Status = "200"
			(*response.(*gmp.ModifyNoteResponse)).StatusText = "OK"
		} else {
			(*response.(*gmp.ModifyNoteResponse)).Status = "400"
			(*response.(*gmp.ModifyNoteResponse)).StatusText = "Bad request"
		}
	}

	if _, ok := command.(*gmp.GetNotesCommand); ok {
		(*response.(*gmp.GetNotesResponse)).Status = "200"
		(*response.(*gmp.GetNotesResponse)).StatusText = "OK"
		(*response.(*gmp.GetNotesResponse)).Notes = []gmp.NoteInfo{
			{
				ID:               "b76b81a7-9df8-42df-afff-baa9d4620128",
				NVT:              &gmp.GetNoteNVT{OID: "1.3.6.1.4.1.25623.1.0.75", Name: "Test NVT: long lines"},
				CreationTime:     "2013-01-09T09:47:41+01:00",
				ModificationTime: "2013-01-09T09:47:41+01:00",
				Writable:         "1",
				InUse:            "0",
				Active:           "1",
				Orphan:           "1",
				Text:             "This is the full text of the note.",
			},
		}
	}

	return nil
}

func (m *mockConn) Close() error {
	return nil
}

func mockedConnection() gmp.Connection {
	return &mockConn{}
}
