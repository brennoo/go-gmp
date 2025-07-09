package client

import (
	"github.com/brennoo/go-gmp"
	"github.com/brennoo/go-gmp/commands"
)

type mockConn struct{}

// nolint:gocyclo // This mock handles all GMP commands
// gocyclo:ignore
func (m *mockConn) Execute(command interface{}, response interface{}) error {
	if cmd, ok := command.(*commands.Authenticate); ok {
		if cmd.Credentials.Username == "openvas" && cmd.Credentials.Password == "123" {
			(*response.(*commands.AuthenticateResponse)).Status = "200"
		} else {
			(*response.(*commands.AuthenticateResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.GetConfigs); ok {
		if cmd.ConfigID == "bde773f3-2b3d-4fe6-81cb-6321ae2cc629" {
			(*response.(*commands.GetConfigsResponse)).Status = "200"
		} else {
			(*response.(*commands.GetConfigsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.GetScanners); ok {
		if cmd.ScannerID == "ee0311e7-3247-4425-bb9c-866d59f1e0e9" {
			(*response.(*commands.GetScannersResponse)).Status = "200"
		} else {
			(*response.(*commands.GetScannersResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.GetPreferences); ok {
		if cmd.ConfigID == "4b49617e-d1d8-44b8-af81-f4675b56f837" {
			(*response.(*commands.GetPreferencesResponse)).Status = "200"
		} else {
			(*response.(*commands.GetPreferencesResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.CreateConfig); ok {
		if cmd.Name == "New Config" {
			(*response.(*commands.CreateConfigResponse)).Status = "200"
		} else {
			(*response.(*commands.CreateConfigResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.ModifyConfig); ok {
		if cmd.Name == "Modified Config" {
			(*response.(*commands.ModifyConfigResponse)).Status = "200"
		} else {
			(*response.(*commands.ModifyConfigResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.CreateTask); ok {
		if cmd.Name == "New Task" {
			(*response.(*commands.CreateTaskResponse)).Status = "200"
		} else {
			(*response.(*commands.CreateTaskResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.ModifyTask); ok {
		if cmd.Comment == "Modified Task Comment" {
			(*response.(*commands.ModifyTaskResponse)).Status = "200"
		} else {
			(*response.(*commands.ModifyTaskResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.CreateTarget); ok {
		if cmd.Name == "New Target" {
			(*response.(*commands.CreateTargetResponse)).Status = "200"
		} else {
			(*response.(*commands.CreateTargetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.ModifyTarget); ok {
		if cmd.Name == "Modified Target" {
			(*response.(*commands.ModifyTargetResponse)).Status = "200"
		} else {
			(*response.(*commands.ModifyTargetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.GetTargets); ok {
		if cmd.TargetID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*commands.GetTargetsResponse)).Status = "200"
		} else {
			(*response.(*commands.GetTargetsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.DeleteTarget); ok {
		if cmd.TargetID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*commands.DeleteTargetResponse)).Status = "200"
		} else {
			(*response.(*commands.DeleteTargetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.GetPortLists); ok {
		if cmd.PortListID == "33d0cd82-57c6-11e1-8ed1-406186ea4fc5" {
			(*response.(*commands.GetPortListsResponse)).Status = "200"
		} else {
			(*response.(*commands.GetPortListsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.StartTask); ok {
		if cmd.TaskID == "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31" {
			(*response.(*commands.StartTaskResponse)).Status = "200"
		} else {
			(*response.(*commands.StartTaskResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.GetTasks); ok {
		if cmd.TaskID == "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31" {
			(*response.(*commands.GetTasksResponse)).Status = "200"
		} else {
			(*response.(*commands.GetTasksResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.StopTask); ok {
		if cmd.TaskID == "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31" {
			(*response.(*commands.StopTaskResponse)).Status = "200"
		} else {
			(*response.(*commands.StopTaskResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.DeleteTask); ok {
		if cmd.TaskID == "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31" {
			(*response.(*commands.DeleteTaskResponse)).Status = "200"
		} else {
			(*response.(*commands.DeleteTaskResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.GetResults); ok {
		if cmd.TaskID == "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31" {
			(*response.(*commands.GetResultsResponse)).Status = "200"
		} else {
			(*response.(*commands.GetResultsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.GetVulns); ok {
		if cmd.VulnID == "1.3.6.1.4.1.25623.1.0.808160" {
			(*response.(*commands.GetVulnsResponse)).Status = "200"
		} else {
			(*response.(*commands.GetVulnsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.CreateAlert); ok {
		if cmd.Name == "Test Alert" {
			(*response.(*commands.CreateAlertResponse)).Status = "201"
			(*response.(*commands.CreateAlertResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateAlertResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*commands.CreateAlertResponse)).Status = "400"
		}
	}

	if _, ok := command.(*commands.GetAlerts); ok {
		(*response.(*commands.GetAlertsResponse)).Status = "200"
		(*response.(*commands.GetAlertsResponse)).StatusText = "OK"
	}

	if _, ok := command.(*commands.ModifyAlert); ok {
		(*response.(*commands.ModifyAlertResponse)).Status = "200"
		(*response.(*commands.ModifyAlertResponse)).StatusText = "OK"
	}

	if _, ok := command.(*commands.DeleteAlert); ok {
		(*response.(*commands.DeleteAlertResponse)).Status = "200"
		(*response.(*commands.DeleteAlertResponse)).StatusText = "OK"
	}

	if _, ok := command.(*commands.TestAlert); ok {
		(*response.(*commands.TestAlertResponse)).Status = "200"
		(*response.(*commands.TestAlertResponse)).StatusText = "OK"
	}

	if _, ok := command.(*commands.ResumeTask); ok {
		(*response.(*commands.ResumeTaskResponse)).Status = "200"
		(*response.(*commands.ResumeTaskResponse)).StatusText = "OK"
		(*response.(*commands.ResumeTaskResponse)).ReportID = "330ee785-c2c0-4d4c-ab96-725142c9b789"
	}

	if cmd, ok := command.(*commands.CreateAsset); ok {
		if cmd.Asset != nil && cmd.Asset.Name == "Localhost" && cmd.Asset.Type == "host" {
			(*response.(*commands.CreateAssetResponse)).Status = "201"
			(*response.(*commands.CreateAssetResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateAssetResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else if cmd.Report != nil && cmd.Report.ID == "report-uuid" && cmd.Report.Filter != nil && cmd.Report.Filter.Term == "min_qod=70" {
			(*response.(*commands.CreateAssetResponse)).Status = "201"
			(*response.(*commands.CreateAssetResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateAssetResponse)).ID = "report-asset-uuid"
		} else {
			(*response.(*commands.CreateAssetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.ModifyAsset); ok {
		if cmd.AssetID == "914b59f8-25f5-4c8f-832c-2379cd625236" && cmd.Comment == "New comment" {
			(*response.(*commands.ModifyAssetResponse)).Status = "200"
			(*response.(*commands.ModifyAssetResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyAssetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.GetAssets); ok {
		if cmd.AssetID == "b493b7a8-7489-11df-a3ec-002264764cea" {
			(*response.(*commands.GetAssetsResponse)).Status = "200"
			(*response.(*commands.GetAssetsResponse)).StatusText = "OK"
			(*response.(*commands.GetAssetsResponse)).Assets = []commands.Asset{
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
			(*response.(*commands.GetAssetsResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.DeleteAsset); ok {
		if cmd.AssetID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" {
			(*response.(*commands.DeleteAssetResponse)).Status = "200"
			(*response.(*commands.DeleteAssetResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteAssetResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.CreateSchedule); ok {
		if cmd.Name == "Monthly Scan" && cmd.Timezone == "America/New_York" && cmd.ICalendar != "" {
			(*response.(*commands.CreateScheduleResponse)).Status = "201"
			(*response.(*commands.CreateScheduleResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateScheduleResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*commands.CreateScheduleResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.ModifySchedule); ok {
		if cmd.ScheduleID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" && cmd.Name == "Weekly Scan" {
			(*response.(*commands.ModifyScheduleResponse)).Status = "200"
			(*response.(*commands.ModifyScheduleResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyScheduleResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.GetSchedules); ok {
		if cmd.ScheduleID == "ddda859a-45be-4c58-85b3-517c66230232" {
			(*response.(*commands.GetSchedulesResponse)).Status = "200"
			(*response.(*commands.GetSchedulesResponse)).StatusText = "OK"
			(*response.(*commands.GetSchedulesResponse)).Schedules = []commands.Schedule{
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
			(*response.(*commands.GetSchedulesResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.DeleteSchedule); ok {
		if cmd.ScheduleID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" {
			(*response.(*commands.DeleteScheduleResponse)).Status = "200"
			(*response.(*commands.DeleteScheduleResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteScheduleResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.CreateOverride); ok {
		if cmd.Text == "This is actually of little concern." && cmd.NVT.OID == "1.3.6.1.4.1.25623.1.0.10330" && cmd.NewThreat == "Low" && cmd.Result == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*commands.CreateOverrideResponse)).Status = "201"
			(*response.(*commands.CreateOverrideResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateOverrideResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*commands.CreateOverrideResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.DeleteOverride); ok {
		if cmd.OverrideID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" {
			(*response.(*commands.DeleteOverrideResponse)).Status = "200"
			(*response.(*commands.DeleteOverrideResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteOverrideResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.GetOverrides); ok {
		if cmd.OverrideID == "b76b81a7-9df8-42df-afff-baa9d4620128" {
			(*response.(*commands.GetOverridesResponse)).Status = "200"
			(*response.(*commands.GetOverridesResponse)).StatusText = "OK"
			(*response.(*commands.GetOverridesResponse)).Overrides = []commands.Override{
				{
					ID:          "b76b81a7-9df8-42df-afff-baa9d4620128",
					Text:        commands.OverrideText{Text: "This is the full text of the override."},
					NewThreat:   "Log",
					NewSeverity: "0.0",
					Orphan:      true,
					Active:      true,
					NVT:         commands.OverrideNVT{OID: "1.3.6.1.4.1.25623.1.0.75", Name: "Test NVT: long lines"},
				},
			}
		} else {
			(*response.(*commands.GetOverridesResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.ModifyOverride); ok {
		if cmd.OverrideID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" && cmd.Text == "This issue is less important in our setup." && cmd.NewThreat == "Low" {
			(*response.(*commands.ModifyOverrideResponse)).Status = "200"
			(*response.(*commands.ModifyOverrideResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyOverrideResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.CreateReport); ok {
		if cmd.Report != nil && cmd.Report.ID == "report-uuid" && cmd.Task != nil && cmd.Task.ID == "task-uuid" {
			(*response.(*commands.CreateReportResponse)).Status = "201"
			(*response.(*commands.CreateReportResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateReportResponse)).ID = "created-report-id"
		} else {
			(*response.(*commands.CreateReportResponse)).Status = "400"
			(*response.(*commands.CreateReportResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.GetReports); ok {
		if cmd.ReportID == "report-uuid" && cmd.FormatID == "format-uuid" {
			(*response.(*commands.GetReportsResponse)).Status = "200"
			(*response.(*commands.GetReportsResponse)).StatusText = "OK"
			(*response.(*commands.GetReportsResponse)).Reports = []commands.ReportWrapper{
				{
					ID:          "report-uuid",
					FormatID:    "format-uuid",
					Extension:   "xml",
					ContentType: "text/xml",
				},
			}
		} else {
			(*response.(*commands.GetReportsResponse)).Status = "404"
			(*response.(*commands.GetReportsResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.DeleteReport); ok {
		if cmd.ReportID == "report-uuid" {
			(*response.(*commands.DeleteReportResponse)).Status = "200"
			(*response.(*commands.DeleteReportResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteReportResponse)).Status = "404"
			(*response.(*commands.DeleteReportResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.CreateReportFormat); ok {
		if cmd.Name == "Test Format" && cmd.Copy == "copy-uuid" {
			(*response.(*commands.CreateReportFormatResponse)).Status = "201"
			(*response.(*commands.CreateReportFormatResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateReportFormatResponse)).ID = "created-format-id"
		} else {
			(*response.(*commands.CreateReportFormatResponse)).Status = "400"
			(*response.(*commands.CreateReportFormatResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.ModifyReportFormat); ok {
		if cmd.ReportFormatID == "format-uuid" && cmd.Name == "New Name" {
			(*response.(*commands.ModifyReportFormatResponse)).Status = "200"
			(*response.(*commands.ModifyReportFormatResponse)).StatusText = "OK"
		} else if cmd.ReportFormatID == "format-uuid" && cmd.Param != nil && cmd.Param.Name == "Background Colour" && cmd.Param.Value == "red" {
			(*response.(*commands.ModifyReportFormatResponse)).Status = "200"
			(*response.(*commands.ModifyReportFormatResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyReportFormatResponse)).Status = "400"
			(*response.(*commands.ModifyReportFormatResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.GetReportFormats); ok {
		if cmd.ReportFormatID == "format-uuid" {
			(*response.(*commands.GetReportFormatsResponse)).Status = "200"
			(*response.(*commands.GetReportFormatsResponse)).StatusText = "OK"
			(*response.(*commands.GetReportFormatsResponse)).ReportFormats = []commands.ReportFormatWrapper{
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
			(*response.(*commands.GetReportFormatsResponse)).Status = "404"
			(*response.(*commands.GetReportFormatsResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.DeleteReportFormat); ok {
		if cmd.ReportFormatID == "format-uuid" && cmd.Ultimate == "1" {
			(*response.(*commands.DeleteReportFormatResponse)).Status = "200"
			(*response.(*commands.DeleteReportFormatResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteReportFormatResponse)).Status = "404"
			(*response.(*commands.DeleteReportFormatResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.VerifyReportFormat); ok {
		if cmd.ReportFormatID == "format-uuid" {
			(*response.(*commands.VerifyReportFormatResponse)).Status = "200"
			(*response.(*commands.VerifyReportFormatResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.VerifyReportFormatResponse)).Status = "404"
			(*response.(*commands.VerifyReportFormatResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.CreateReportConfig); ok {
		if cmd.Name == "Test config" && cmd.ReportFormat.ID == "format-uuid" {
			(*response.(*commands.CreateReportConfigResponse)).Status = "201"
			(*response.(*commands.CreateReportConfigResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateReportConfigResponse)).ID = "created-config-id"
		} else {
			(*response.(*commands.CreateReportConfigResponse)).Status = "400"
			(*response.(*commands.CreateReportConfigResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.ModifyReportConfig); ok {
		if len(cmd.Params) > 0 && cmd.Params[0].Name == "Node Distance" && cmd.Params[0].UseDefault == "1" {
			(*response.(*commands.ModifyReportConfigResponse)).Status = "200"
			(*response.(*commands.ModifyReportConfigResponse)).StatusText = "OK"
			(*response.(*commands.ModifyReportConfigResponse)).ID = "modified-config-id"
		} else if cmd.Name == "Renamed config" {
			(*response.(*commands.ModifyReportConfigResponse)).Status = "200"
			(*response.(*commands.ModifyReportConfigResponse)).StatusText = "OK"
			(*response.(*commands.ModifyReportConfigResponse)).ID = "modified-config-id"
		} else {
			(*response.(*commands.ModifyReportConfigResponse)).Status = "400"
			(*response.(*commands.ModifyReportConfigResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.GetReportConfigs); ok {
		if cmd.ReportConfigID == "config-uuid" {
			(*response.(*commands.GetReportConfigsResponse)).Status = "200"
			(*response.(*commands.GetReportConfigsResponse)).StatusText = "OK"
			(*response.(*commands.GetReportConfigsResponse)).ReportConfigs = []commands.ReportConfigWrapper{
				{
					ID:               "config-uuid",
					Name:             "Test config",
					Owner:            &commands.ReportConfigOwner{Name: "admin"},
					Comment:          "Test comment",
					CreationTime:     "2024-01-23T09:43:03Z",
					ModificationTime: "2024-01-26T14:11:54Z",
					Writable:         "1",
					InUse:            "0",
					Permissions: &commands.ReportConfigPermissions{
						Permissions: []commands.ReportConfigPermission{{Name: "Everything"}},
					},
					ReportFormat: &commands.ReportConfigFormat{
						ID:   "format-uuid",
						Name: "Topology SVG",
					},
					Params: []commands.ReportConfigParam{
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
			(*response.(*commands.GetReportConfigsResponse)).Status = "404"
			(*response.(*commands.GetReportConfigsResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.DeleteReportConfig); ok {
		if cmd.ReportConfigID == "config-uuid" && cmd.Ultimate == "1" {
			(*response.(*commands.DeleteReportConfigResponse)).Status = "200"
			(*response.(*commands.DeleteReportConfigResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteReportConfigResponse)).Status = "404"
			(*response.(*commands.DeleteReportConfigResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.GetSystemReports); ok {
		if cmd.Name == "proc" {
			(*response.(*commands.GetSystemReportsResponse)).Status = "200"
			(*response.(*commands.GetSystemReportsResponse)).StatusText = "OK"
			(*response.(*commands.GetSystemReportsResponse)).SystemReports = []commands.SystemReportWrapper{
				{
					Name:  "proc",
					Title: "Processes",
					Report: &commands.SystemReportContent{
						Format:   "png",
						Duration: "86400",
						Value:    "iVBORw0KGgoAAAANSUhEUgAAArkAAAE...2bEdAAAAAElFTkSuQmCC",
					},
				},
			}
		} else if cmd.Brief == "1" {
			(*response.(*commands.GetSystemReportsResponse)).Status = "200"
			(*response.(*commands.GetSystemReportsResponse)).StatusText = "OK"
			(*response.(*commands.GetSystemReportsResponse)).SystemReports = []commands.SystemReportWrapper{
				{Name: "proc", Title: "Processes"},
				{Name: "load", Title: "System Load"},
			}
		} else {
			(*response.(*commands.GetSystemReportsResponse)).Status = "404"
			(*response.(*commands.GetSystemReportsResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.CreateCredential); ok {
		if cmd.Name == "Test Credential" {
			(*response.(*commands.CreateCredentialResponse)).Status = "201"
			(*response.(*commands.CreateCredentialResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateCredentialResponse)).ID = "created-credential-id"
		} else {
			(*response.(*commands.CreateCredentialResponse)).Status = "400"
			(*response.(*commands.CreateCredentialResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.ModifyCredential); ok {
		if cmd.CredentialID == "cred-uuid" {
			(*response.(*commands.ModifyCredentialResponse)).Status = "200"
			(*response.(*commands.ModifyCredentialResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyCredentialResponse)).Status = "400"
			(*response.(*commands.ModifyCredentialResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.GetCredentials); ok {
		if cmd.CredentialID == "" {
			(*response.(*commands.GetCredentialsResponse)).Status = "200"
			(*response.(*commands.GetCredentialsResponse)).StatusText = "OK"
			(*response.(*commands.GetCredentialsResponse)).Credentials = []commands.CredentialWrapper{
				{
					ID:       "cred-uuid-1",
					Name:     "sally",
					Login:    "sally",
					Writable: "1",
					InUse:    "0",
					Type:     "usk",
					FullType: "username + SSH key",
					Formats:  &commands.CredentialFormats{Formats: []string{"key", "rpm", "deb"}},
				},
				{
					ID:       "cred-uuid-2",
					Name:     "bob",
					Login:    "bob",
					Writable: "1",
					InUse:    "1",
					Type:     "up",
					FullType: "username + password",
					Formats:  &commands.CredentialFormats{Formats: []string{"exe"}},
				},
			}
		} else if cmd.CredentialID == "cred-uuid-1" {
			(*response.(*commands.GetCredentialsResponse)).Status = "200"
			(*response.(*commands.GetCredentialsResponse)).StatusText = "OK"
			(*response.(*commands.GetCredentialsResponse)).Credentials = []commands.CredentialWrapper{
				{
					ID:       "cred-uuid-1",
					Name:     "sally",
					Login:    "sally",
					Writable: "1",
					InUse:    "0",
					Type:     "usk",
					FullType: "username + SSH key",
					Formats:  &commands.CredentialFormats{Formats: []string{"key", "rpm", "deb"}},
				},
			}
		} else {
			(*response.(*commands.GetCredentialsResponse)).Status = "404"
			(*response.(*commands.GetCredentialsResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.DeleteCredential); ok {
		if cmd.CredentialID == "cred-uuid-1" && cmd.Ultimate == "1" {
			(*response.(*commands.DeleteCredentialResponse)).Status = "200"
			(*response.(*commands.DeleteCredentialResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteCredentialResponse)).Status = "404"
			(*response.(*commands.DeleteCredentialResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.CreateScanner); ok {
		if cmd.Name == "Default Scanner" && cmd.Host == "localhost" && cmd.Port == "9391" && cmd.Type == "2" && cmd.CAPub != "" && cmd.Credential != nil && cmd.Credential.ID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*commands.CreateScannerResponse)).Status = "201"
			(*response.(*commands.CreateScannerResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateScannerResponse)).ID = "814cd30f-dee1-4d58-851d-21b8d0c048e3"
		} else {
			(*response.(*commands.CreateScannerResponse)).Status = "400"
			(*response.(*commands.CreateScannerResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.ModifyScanner); ok {
		if cmd.ScannerID == "scanner-uuid" && cmd.Name == "Updated Scanner" {
			(*response.(*commands.ModifyScannerResponse)).Status = "200"
			(*response.(*commands.ModifyScannerResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyScannerResponse)).Status = "400"
			(*response.(*commands.ModifyScannerResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.DeleteScanner); ok {
		if cmd.ScannerID == "scanner-uuid" && cmd.Ultimate == "1" {
			(*response.(*commands.DeleteScannerResponse)).Status = "200"
			(*response.(*commands.DeleteScannerResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteScannerResponse)).Status = "404"
			(*response.(*commands.DeleteScannerResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.VerifyScanner); ok {
		if cmd.ScannerID == "scanner-uuid" {
			(*response.(*commands.VerifyScannerResponse)).Status = "200"
			(*response.(*commands.VerifyScannerResponse)).StatusText = "OK"
			(*response.(*commands.VerifyScannerResponse)).Version = "OTP/2.0"
		} else {
			(*response.(*commands.VerifyScannerResponse)).Status = "404"
			(*response.(*commands.VerifyScannerResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.CreatePortList); ok {
		if cmd.Name == "All TCP" && cmd.PortRange == "T:1-65535" {
			(*response.(*commands.CreatePortListResponse)).Status = "201"
			(*response.(*commands.CreatePortListResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreatePortListResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*commands.CreatePortListResponse)).Status = "400"
			(*response.(*commands.CreatePortListResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.ModifyPortList); ok {
		if cmd.PortListID == "27140836-05ae-4e8b-9abf-f725ddc2888f" {
			(*response.(*commands.ModifyPortListResponse)).Status = "200"
			(*response.(*commands.ModifyPortListResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyPortListResponse)).Status = "400"
			(*response.(*commands.ModifyPortListResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.DeletePortList); ok {
		if cmd.PortListID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" && cmd.Ultimate == "1" {
			(*response.(*commands.DeletePortListResponse)).Status = "200"
			(*response.(*commands.DeletePortListResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeletePortListResponse)).Status = "404"
			(*response.(*commands.DeletePortListResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.CreatePortRange); ok {
		if cmd.PortList.ID == "354cd3ef-bbe1-4d58-859d-21b8d0c046c4" && cmd.Start == "777" && cmd.End == "779" && cmd.Type == "TCP" {
			(*response.(*commands.CreatePortRangeResponse)).Status = "201"
			(*response.(*commands.CreatePortRangeResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreatePortRangeResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*commands.CreatePortRangeResponse)).Status = "400"
			(*response.(*commands.CreatePortRangeResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.DeletePortRange); ok {
		if cmd.PortRangeID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" {
			(*response.(*commands.DeletePortRangeResponse)).Status = "200"
			(*response.(*commands.DeletePortRangeResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeletePortRangeResponse)).Status = "404"
			(*response.(*commands.DeletePortRangeResponse)).StatusText = "Not found"
		}
	}

	if _, ok := command.(*commands.DescribeAuth); ok {
		resp := response.(*commands.DescribeAuthResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Groups = []commands.DescribeAuthGroup{
			{
				Name: "method:file",
				Settings: []commands.DescribeAuthSetting{
					{Key: "enable", Value: "true"},
					{Key: "order", Value: "1"},
				},
			},
		}
	}

	if _, ok := command.(*commands.GetVersion); ok {
		resp := response.(*commands.GetVersionResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Version = "1.0"
	}

	if _, ok := command.(*commands.Help); ok {
		resp := response.(*commands.HelpResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Text = "AUTHENTICATE           Authenticate with the manager.\nCREATE_ALERT           Create an alert.\n...\nVERIFY_SCANNER         Verify a scanner."
	}

	if cmd, ok := command.(*commands.GetInfo); ok {
		if cmd.Name == "CVE-2011-0018" && cmd.Type == "cve" {
			resp := response.(*commands.GetInfoResponse)
			resp.Status = "200"
			resp.StatusText = "OK"
			resp.Infos = []commands.GetInfoInfo{{ID: "CVE-2011-0018", Name: "CVE-2011-0018"}}
		} else {
			resp := response.(*commands.GetInfoResponse)
			resp.Status = "404"
			resp.StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.DeleteConfig); ok {
		if cmd.ConfigID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" && cmd.Ultimate == "1" {
			(*response.(*commands.DeleteConfigResponse)).Status = "200"
			(*response.(*commands.DeleteConfigResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteConfigResponse)).Status = "404"
			(*response.(*commands.DeleteConfigResponse)).StatusText = "Not found"
		}
	}

	if _, ok := command.(*commands.SyncConfig); ok {
		(*response.(*commands.SyncConfigResponse)).Status = "202"
		(*response.(*commands.SyncConfigResponse)).StatusText = "OK, request submitted"
	}

	if cmd, ok := command.(*commands.MoveTask); ok {
		if cmd.TaskID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" && cmd.SlaveID == "97390ade-e075-11df-9973-002264764cea" {
			(*response.(*commands.MoveTaskResponse)).Status = "200"
			(*response.(*commands.MoveTaskResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.MoveTaskResponse)).Status = "404"
			(*response.(*commands.MoveTaskResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.CreateTLSCertificate); ok {
		if cmd.Name == "Example Certificate" && cmd.Certificate == "MIIDNjCCAp+gAwIBAgIBATANBgkqhkiG9w0BAQQFADCBqTELM[...]" {
			(*response.(*commands.CreateTLSCertificateResponse)).Status = "201"
			(*response.(*commands.CreateTLSCertificateResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateTLSCertificateResponse)).ID = "8a925978-59d0-494b-a837-40b271569727"
		} else {
			(*response.(*commands.CreateTLSCertificateResponse)).Status = "400"
			(*response.(*commands.CreateTLSCertificateResponse)).StatusText = "Bad request"
			(*response.(*commands.CreateTLSCertificateResponse)).ID = ""
		}
	}

	if cmd, ok := command.(*commands.ModifyTLSCertificate); ok {
		if cmd.TLSCertificateID == "8a925978-59d0-494b-a837-40b271569727" && cmd.Name == "Renamed Example Certificate" {
			(*response.(*commands.ModifyTLSCertificateResponse)).Status = "200"
			(*response.(*commands.ModifyTLSCertificateResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyTLSCertificateResponse)).Status = "404"
			(*response.(*commands.ModifyTLSCertificateResponse)).StatusText = "Not found"
		}
	}

	if _, ok := command.(*commands.GetTLSCertificates); ok {
		resp := response.(*commands.GetTLSCertificatesResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.TLSCertificates = []commands.TLSCertificate{
			{
				ID:    "ba36ed15-92fa-4ae0-af53-bad8ce472f18",
				Owner: &commands.Owner{Name: "admin"},
				Name:  "Example Certificate",
			},
		}
	}

	if cmd, ok := command.(*commands.ModifyAgents); ok {
		if len(cmd.Agents) == 1 && cmd.Agents[0].ID == "fb6451bf-ec5a-45a8-8bab-5cf4b862e51b" &&
			cmd.Authorized == "1" && cmd.MinInterval == "1000" && cmd.HeartbeatInterval == "0" &&
			cmd.Schedule == "@every 12h" && cmd.Comment == "example update" {
			(*response.(*commands.ModifyAgentsResponse)).Status = "200"
			(*response.(*commands.ModifyAgentsResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyAgentsResponse)).Status = "400"
			(*response.(*commands.ModifyAgentsResponse)).StatusText = "Bad request"
		}
	}

	if _, ok := command.(*commands.GetAgents); ok {
		resp := response.(*commands.GetAgentsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Agents = []commands.Agent{
			{
				ID:      "62462fe0-5834-4630-afc2-0d040c63487c",
				AgentID: "GAT-29-p0MPX0FT",
			},
		}
	}

	if cmd, ok := command.(*commands.DeleteAgents); ok {
		if len(cmd.Agents) == 1 && cmd.Agents[0].ID == "c6f1cdc3-8c2c-4b2e-9f43-139d23c7cfd4" {
			(*response.(*commands.DeleteAgentsResponse)).Status = "200"
			(*response.(*commands.DeleteAgentsResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteAgentsResponse)).Status = "404"
			(*response.(*commands.DeleteAgentsResponse)).StatusText = "Not found"
		}
	}

	if cmd, ok := command.(*commands.GetNvts); ok {
		if cmd.Details == "1" {
			(*response.(*commands.GetNvtsResponse)).Status = "200"
			(*response.(*commands.GetNvtsResponse)).StatusText = "OK"
			(*response.(*commands.GetNvtsResponse)).Nvts = []commands.NVT{
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
			(*response.(*commands.GetNvtsResponse)).Status = "200"
			(*response.(*commands.GetNvtsResponse)).StatusText = "OK"
			(*response.(*commands.GetNvtsResponse)).Nvts = []commands.NVT{
				{
					OID:  "1.3.6.1.4.1.25623.1.0.10330",
					Name: "Services",
				},
			}
		} else {
			(*response.(*commands.GetNvtsResponse)).Status = "404"
			(*response.(*commands.GetNvtsResponse)).StatusText = "Not found"
		}
	}

	if _, ok := command.(*commands.GetNvtFamilies); ok {
		(*response.(*commands.GetNvtFamiliesResponse)).Status = "200"
		(*response.(*commands.GetNvtFamiliesResponse)).StatusText = "OK"
		(*response.(*commands.GetNvtFamiliesResponse)).Families = []commands.Family{
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

	if _, ok := command.(*commands.GetFeeds); ok {
		resp := response.(*commands.GetFeedsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.FeedOwnerSet = "1"
		resp.FeedRolesSet = "1"
		resp.FeedResourcesAccess = "1"
		resp.Feeds = []commands.Feed{
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

	if cmd, ok := command.(*commands.ModifyLicense); ok {
		if cmd.File != "" {
			(*response.(*commands.ModifyLicenseResponse)).Status = "200"
			(*response.(*commands.ModifyLicenseResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyLicenseResponse)).Status = "400"
			(*response.(*commands.ModifyLicenseResponse)).StatusText = "Missing file"
		}
	}

	if _, ok := command.(*commands.GetLicense); ok {
		resp := response.(*commands.GetLicenseResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.License = &commands.License{
			Status: "active",
			Content: &commands.LicenseContent{
				Meta: &commands.LicenseMeta{
					ID:           "4711",
					Version:      "1.0.0",
					Comment:      "Test License",
					Type:         "trial",
					CustomerName: "Jane Doe",
					Created:      "2021-08-27T06:05:21Z",
					Begins:       "2021-08-27T07:05:21Z",
					Expires:      "2021-09-04T07:05:21Z",
				},
				Appliance: &commands.LicenseAppliance{
					Model:     "trial",
					ModelType: "virtual",
					Sensor:    "0",
				},
				Keys: &commands.LicenseKeys{
					Keys: []commands.LicenseKey{{Name: "feed", Value: "*base64 GSF key*"}},
				},
				Signatures: &commands.LicenseSignatures{
					Signatures: []commands.LicenseSignature{{Name: "license", Value: "*base64 signature*"}},
				},
			},
		}
	}

	if cmd, ok := command.(*commands.ModifySetting); ok {
		if cmd.Name == "Timezone" && cmd.Value == "QWZyaWNhL0pvaGFubmVzYnVyZw==" {
			(*response.(*commands.ModifySettingResponse)).Status = "200"
			(*response.(*commands.ModifySettingResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifySettingResponse)).Status = "400"
			(*response.(*commands.ModifySettingResponse)).StatusText = "Bad request"
		}
	}

	if _, ok := command.(*commands.GetSettings); ok {
		resp := response.(*commands.GetSettingsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Settings = &commands.SettingsBlock{
			Setting: []commands.Setting{{
				ID:    "5f5a8712-8017-11e1-8556-406186ea4fc5",
				Name:  "Rows Per Page",
				Value: "15",
			}},
		}
	}

	if _, ok := command.(*commands.GetResourceNames); ok {
		resp := response.(*commands.GetResourceNamesResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Type = "os"
		resp.Resources = []commands.ResourceName{
			{ID: "5b6b6aef-b320-42ca-899f-3161ee2a4195", Name: "cpe:/h:hp:jetdirect"},
			{ID: "5be25864-9249-431e-8a91-039e334371ad", Name: "cpe:/o:canonical:ubuntu_linux:18.04"},
		}
	}

	if _, ok := command.(*commands.GetAggregates); ok {
		resp := response.(*commands.GetAggregatesResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Aggregate = &commands.AggregateBlock{
			DataType:    "nvt",
			DataColumn:  "severity",
			GroupColumn: "family",
			Groups: []commands.AggregateGroup{
				{Value: "AIX Local Security Checks", Count: 1, CCount: 1, Min: 3.3, Max: 3.3, Mean: 3.3, Sum: 3.3, CSum: 3.3},
				{Value: "Brute force attacks", Count: 8, CCount: 9, Min: 0, Max: 7.8, Mean: 6.275, Sum: 50.2, CSum: 53.5},
			},
		}
	}

	if _, ok := command.(*commands.GetFeatures); ok {
		resp := response.(*commands.GetFeaturesResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Features = []commands.FeatureInfo{{
			Enabled: "0",
			Name:    "OPENVASD",
		}}
	}

	if _, ok := command.(*commands.EmptyTrashcan); ok {
		resp := response.(*commands.EmptyTrashcanResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
	}

	if _, ok := command.(*commands.Restore); ok {
		resp := response.(*commands.RestoreResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
	}

	if _, ok := command.(*commands.RunWizard); ok {
		resp := response.(*commands.RunWizardResponse)
		resp.Status = "202"
		resp.StatusText = "OK, request submitted"
		resp.Response = &commands.WizardResponse{
			StartTaskResponse: &commands.StartTaskResponse{
				Status:     "202",
				StatusText: "OK, request submitted",
				ReportID:   "a06d21f7-8e2f-4d7f-bceb-1df852d8b37d",
			},
		}
	}

	if cmd, ok := command.(*commands.CreateUser); ok {
		if cmd.Name == "testuser" && cmd.Password == "testpass" && len(cmd.Roles) > 0 && cmd.Roles[0].ID == "role-uuid" {
			(*response.(*commands.CreateUserResponse)).Status = "201"
			(*response.(*commands.CreateUserResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateUserResponse)).ID = "created-user-id"
		} else {
			(*response.(*commands.CreateUserResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.ModifyUser); ok {
		if cmd.UserID == "user-uuid" && cmd.Name == "testuser" && cmd.NewName == "newuser" && cmd.Comment == "Updated user" && cmd.Password != nil && cmd.Password.Modify == "1" && cmd.Password.Text == "newpass" && len(cmd.Roles) > 0 && cmd.Roles[0].ID == "role-uuid" {
			(*response.(*commands.ModifyUserResponse)).Status = "200"
			(*response.(*commands.ModifyUserResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyUserResponse)).Status = "400"
		}
	}
	if _, ok := command.(*commands.GetUsers); ok {
		(*response.(*commands.GetUsersResponse)).Status = "200"
		(*response.(*commands.GetUsersResponse)).StatusText = "OK"
		(*response.(*commands.GetUsersResponse)).Users = []commands.UserWrapper{
			{
				Name:    "foobar",
				Role:    &commands.UserRole{ID: "8d453140-b74d-11e2-b0be-406186ea4fc5", Name: "User"},
				Hosts:   &commands.CreateUserHosts{Allow: "2"},
				Sources: &commands.CreateUserSources{Source: &commands.CreateUserSource{Type: "file"}},
			},
		}
	}

	if cmd, ok := command.(*commands.DeleteUser); ok {
		if cmd.Name == "foobar" {
			(*response.(*commands.DeleteUserResponse)).Status = "200"
			(*response.(*commands.DeleteUserResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteUserResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.CreateGroup); ok {
		if cmd.Name == "Managers" {
			(*response.(*commands.CreateGroupResponse)).Status = "201"
			(*response.(*commands.CreateGroupResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateGroupResponse)).ID = "d94211b6-ba40-11e3-bcb1-406186ea4fc5"
		} else {
			(*response.(*commands.CreateGroupResponse)).Status = "400"
			(*response.(*commands.CreateGroupResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.CreateRole); ok {
		if cmd.Name == "Test Role" && cmd.Comment == "A test role" {
			(*response.(*commands.CreateRoleResponse)).Status = "201"
			(*response.(*commands.CreateRoleResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateRoleResponse)).ID = "created-role-id"
		} else {
			(*response.(*commands.CreateRoleResponse)).Status = "400"
			(*response.(*commands.CreateRoleResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.ModifyRole); ok {
		if cmd.RoleID == "role-uuid" && cmd.Name == "Updated Role" && cmd.Comment == "Updated comment" && cmd.Users == "user1,user2" {
			(*response.(*commands.ModifyRoleResponse)).Status = "200"
			(*response.(*commands.ModifyRoleResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyRoleResponse)).Status = "400"
			(*response.(*commands.ModifyRoleResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.GetRoles); ok {
		if mockFail, ok := any(cmd).(interface{ Fail() bool }); ok && mockFail.Fail() {
			(*response.(*commands.GetRolesResponse)).Status = "400"
			(*response.(*commands.GetRolesResponse)).StatusText = "Bad request"
			(*response.(*commands.GetRolesResponse)).Roles = nil
		} else {
			(*response.(*commands.GetRolesResponse)).Status = "200"
			(*response.(*commands.GetRolesResponse)).StatusText = "OK"
			(*response.(*commands.GetRolesResponse)).Roles = []commands.Role{
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

	if cmd, ok := command.(*commands.DeleteRole); ok {
		if cmd.RoleID == "b64c81b2-b9de-11e3-a2e9-406186ea4fc5" && cmd.Ultimate == "1" {
			(*response.(*commands.DeleteRoleResponse)).Status = "200"
			(*response.(*commands.DeleteRoleResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteRoleResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.CreatePermission); ok {
		if cmd.Name == "get_targets" && cmd.Subject != nil && cmd.Subject.ID == "66abe5ce-c011-11e3-b96e-406186ea4fc5" && cmd.Subject.Type == "user" && cmd.Resource != nil && cmd.Resource.ID == "b493b7a8-7489-11df-a3ec-002264764cea" {
			(*response.(*commands.CreatePermissionResponse)).Status = "201"
			(*response.(*commands.CreatePermissionResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreatePermissionResponse)).ID = "created-permission-id"
		} else {
			(*response.(*commands.CreatePermissionResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.ModifyPermission); ok {
		if cmd.PermissionID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" && cmd.Subject != nil && cmd.Subject.ID == "76e47468-c095-11e3-9285-406186ea4fc5" && cmd.Subject.Type == "user" {
			(*response.(*commands.ModifyPermissionResponse)).Status = "200"
			(*response.(*commands.ModifyPermissionResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyPermissionResponse)).Status = "400"
		}
	}

	if _, ok := command.(*commands.GetPermissions); ok {
		(*response.(*commands.GetPermissionsResponse)).Status = "200"
		(*response.(*commands.GetPermissionsResponse)).StatusText = "OK"
		(*response.(*commands.GetPermissionsResponse)).Permissions = []commands.PermissionInfo{
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

	if cmd, ok := command.(*commands.DeletePermission); ok {
		if cmd.PermissionID == "267a3405-e84a-47da-97b2-5fa0d2e8995e" {
			(*response.(*commands.DeletePermissionResponse)).Status = "200"
			(*response.(*commands.DeletePermissionResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeletePermissionResponse)).Status = "400"
			(*response.(*commands.DeletePermissionResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.ModifyGroup); ok {
		if cmd.GroupID == "d94211b6-ba40-11e3-bcb1-406186ea4fc5" && cmd.Name == "Line Managers" {
			(*response.(*commands.ModifyGroupResponse)).Status = "200"
			(*response.(*commands.ModifyGroupResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyGroupResponse)).Status = "400"
			(*response.(*commands.ModifyGroupResponse)).StatusText = "Bad request"
		}
	}

	if _, ok := command.(*commands.GetGroups); ok {
		(*response.(*commands.GetGroupsResponse)).Status = "200"
		(*response.(*commands.GetGroupsResponse)).StatusText = "OK"
		(*response.(*commands.GetGroupsResponse)).Groups = []commands.GroupInfo{
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

	if cmd, ok := command.(*commands.DeleteGroup); ok {
		if cmd.GroupID == "d94211b6-ba40-11e3-bcb1-406186ea4fc5" {
			(*response.(*commands.DeleteGroupResponse)).Status = "200"
			(*response.(*commands.DeleteGroupResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteGroupResponse)).Status = "400"
			(*response.(*commands.DeleteGroupResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.CreateTag); ok {
		if cmd.Name == "geo:long" && cmd.Resources != nil && cmd.Resources.Resource != nil && cmd.Resources.Resource.ID == "b493b7a8-7489-11df-a3ec-002264764cea" {
			(*response.(*commands.CreateTagResponse)).Status = "201"
			(*response.(*commands.CreateTagResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateTagResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*commands.CreateTagResponse)).Status = "400"
			(*response.(*commands.CreateTagResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.ModifyTag); ok {
		if cmd.TagID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" && cmd.Active == "0" {
			(*response.(*commands.ModifyTagResponse)).Status = "200"
			(*response.(*commands.ModifyTagResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyTagResponse)).Status = "400"
			(*response.(*commands.ModifyTagResponse)).StatusText = "Bad request"
		}
	}

	if _, ok := command.(*commands.GetTags); ok {
		(*response.(*commands.GetTagsResponse)).Status = "200"
		(*response.(*commands.GetTagsResponse)).StatusText = "OK"
		(*response.(*commands.GetTagsResponse)).Tags = []commands.TagInfo{
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
				Resources: []commands.TagResource{
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

	if cmd, ok := command.(*commands.DeleteTag); ok {
		if cmd.TagID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*commands.DeleteTagResponse)).Status = "200"
			(*response.(*commands.DeleteTagResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteTagResponse)).Status = "400"
			(*response.(*commands.DeleteTagResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.CreateNote); ok {
		if cmd.Text == "This issue should be resolved after the upgrade." && cmd.NVT != nil && cmd.NVT.OID == "1.3.6.1.4.1.25623.1.0.10330" && cmd.Result != nil && cmd.Result.ID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*commands.CreateNoteResponse)).Status = "202"
			(*response.(*commands.CreateNoteResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateNoteResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*commands.CreateNoteResponse)).Status = "400"
			(*response.(*commands.CreateNoteResponse)).StatusText = "Bad request"
		}
	}

	if cmd, ok := command.(*commands.ModifyNote); ok {
		if cmd.NoteID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" && cmd.Text == "This issue should be resolved after the upgrade." && cmd.Result != nil && cmd.Result.ID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*commands.ModifyNoteResponse)).Status = "200"
			(*response.(*commands.ModifyNoteResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyNoteResponse)).Status = "400"
			(*response.(*commands.ModifyNoteResponse)).StatusText = "Bad request"
		}
	}

	if _, ok := command.(*commands.GetNotes); ok {
		(*response.(*commands.GetNotesResponse)).Status = "200"
		(*response.(*commands.GetNotesResponse)).StatusText = "OK"
		(*response.(*commands.GetNotesResponse)).Notes = []commands.NoteInfo{
			{
				ID:               "b76b81a7-9df8-42df-afff-baa9d4620128",
				NVT:              &commands.GetNoteNVT{OID: "1.3.6.1.4.1.25623.1.0.75", Name: "Test NVT: long lines"},
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

	if cmd, ok := command.(*commands.DeleteNote); ok {
		if cmd.NoteID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*commands.DeleteNoteResponse)).Status = "200"
			(*response.(*commands.DeleteNoteResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteNoteResponse)).Status = "400"
			(*response.(*commands.DeleteNoteResponse)).StatusText = "Not Found"
		}
	}

	if cmd, ok := command.(*commands.CreateFilter); ok {
		if cmd.Name == "Single Targets" {
			(*response.(*commands.CreateFilterResponse)).Status = "201"
			(*response.(*commands.CreateFilterResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateFilterResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c7"
		} else {
			(*response.(*commands.CreateFilterResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.ModifyFilter); ok {
		if cmd.FilterID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c7" {
			(*response.(*commands.ModifyFilterResponse)).Status = "200"
			(*response.(*commands.ModifyFilterResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyFilterResponse)).Status = "400"
			(*response.(*commands.ModifyFilterResponse)).StatusText = "Not Found"
		}
	}

	if _, ok := command.(*commands.GetFilters); ok {
		resp := response.(*commands.GetFiltersResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Filters = []commands.FilterEntry{
			{
				ID:               "b493b7a8-7489-11df-a3ec-001164764cea",
				Name:             "Single Targets",
				Comment:          "Targets with only one host",
				Term:             "ips=1 first=1 rows=-2",
				Type:             "target",
				InUse:            "1",
				Writable:         "1",
				CreationTime:     "2015-07-15T15:05:55Z",
				ModificationTime: "2015-07-15T15:05:55Z",
			},
		}
	}

	if cmd, ok := command.(*commands.DeleteFilter); ok {
		if cmd.FilterID == "b493b7a8-7489-11df-a3ec-001164764cea" {
			(*response.(*commands.DeleteFilterResponse)).Status = "200"
			(*response.(*commands.DeleteFilterResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteFilterResponse)).Status = "400"
			(*response.(*commands.DeleteFilterResponse)).StatusText = "Not Found"
		}
	}

	if cmd, ok := command.(*commands.CreateTicket); ok {
		if cmd.Result.ID == "138c1216-4acb-4ded-bef3-7fab80eac8c7" && cmd.Assigned.User.ID == "33e92d3e-a379-4c46-a4cf-88c8201ab710" {
			(*response.(*commands.CreateTicketResponse)).Status = "201"
			(*response.(*commands.CreateTicketResponse)).StatusText = "OK, resource created"
			(*response.(*commands.CreateTicketResponse)).ID = "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"
		} else {
			(*response.(*commands.CreateTicketResponse)).Status = "400"
		}
	}

	if _, ok := command.(*commands.GetTickets); ok {
		resp := response.(*commands.GetTicketsResponse)
		resp.Status = "200"
		resp.StatusText = "OK"
		resp.Tickets = []commands.TicketEntry{
			{
				ID:               "93cd2f71-48c3-4cf2-b542-5b256f59cae0",
				Name:             "OpenSSH Denial of Service Vulnerability - Jan16",
				Comment:          "",
				CreationTime:     "2018-11-29T16:18:56Z",
				ModificationTime: "2018-11-29T16:18:56Z",
				Writable:         "1",
				InUse:            "0",
				Status:           "Open",
				OpenNote:         "Probably the new version fixes this",
			},
		}
	}

	if cmd, ok := command.(*commands.ModifyTicket); ok {
		if cmd.TicketID == "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
			(*response.(*commands.ModifyTicketResponse)).Status = "200"
			(*response.(*commands.ModifyTicketResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyTicketResponse)).Status = "400"
			(*response.(*commands.ModifyTicketResponse)).StatusText = "Not Found"
		}
	}

	if cmd, ok := command.(*commands.DeleteTicket); ok {
		if cmd.TicketID == "ticket-uuid" && cmd.Ultimate == "1" {
			(*response.(*commands.DeleteTicketResponse)).Status = "200"
			(*response.(*commands.DeleteTicketResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.DeleteTicketResponse)).Status = "400"
		}
	}

	if cmd, ok := command.(*commands.ModifyAuth); ok {
		if cmd.Group.Name == "method:file" && len(cmd.Group.Settings) > 0 && cmd.Group.Settings[0].Key == "enable" {
			(*response.(*commands.ModifyAuthResponse)).Status = "200"
			(*response.(*commands.ModifyAuthResponse)).StatusText = "OK"
		} else {
			(*response.(*commands.ModifyAuthResponse)).Status = "400"
			(*response.(*commands.ModifyAuthResponse)).StatusText = "Bad request"
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
