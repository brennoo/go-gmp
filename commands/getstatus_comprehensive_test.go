package commands

import "testing"

// Test all GetStatus() methods for comprehensive coverage.
func TestGetStatusMethods(t *testing.T) {
	testCases := []struct {
		name       string
		response   interface{ GetStatus() (string, string) }
		status     string
		statusText string
	}{
		// Authentication
		{"AuthenticateResponse", &AuthenticateResponse{Status: "200", StatusText: "OK"}, "200", "OK"},

		// Create operations
		{"CreateAlertResponse", &CreateAlertResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateAssetResponse", &CreateAssetResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateConfigResponse", &CreateConfigResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateCredentialResponse", &CreateCredentialResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateFilterResponse", &CreateFilterResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateGroupResponse", &CreateGroupResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateNoteResponse", &CreateNoteResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateOverrideResponse", &CreateOverrideResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreatePermissionResponse", &CreatePermissionResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreatePortListResponse", &CreatePortListResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreatePortRangeResponse", &CreatePortRangeResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateReportResponse", &CreateReportResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateReportConfigResponse", &CreateReportConfigResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateReportFormatResponse", &CreateReportFormatResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateRoleResponse", &CreateRoleResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateScannerResponse", &CreateScannerResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateScheduleResponse", &CreateScheduleResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateTagResponse", &CreateTagResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateTargetResponse", &CreateTargetResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateTaskResponse", &CreateTaskResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateTicketResponse", &CreateTicketResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateTLSCertificateResponse", &CreateTLSCertificateResponse{Status: "201", StatusText: "Created"}, "201", "Created"},
		{"CreateUserResponse", &CreateUserResponse{Status: "201", StatusText: "Created"}, "201", "Created"},

		// Delete operations
		{"DeleteAgentsResponse", &DeleteAgentsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteAlertResponse", &DeleteAlertResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteAssetResponse", &DeleteAssetResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteConfigResponse", &DeleteConfigResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteCredentialResponse", &DeleteCredentialResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteFilterResponse", &DeleteFilterResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteGroupResponse", &DeleteGroupResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteNoteResponse", &DeleteNoteResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteOverrideResponse", &DeleteOverrideResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeletePermissionResponse", &DeletePermissionResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeletePortListResponse", &DeletePortListResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeletePortRangeResponse", &DeletePortRangeResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteReportResponse", &DeleteReportResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteReportConfigResponse", &DeleteReportConfigResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteReportFormatResponse", &DeleteReportFormatResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteRoleResponse", &DeleteRoleResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteScannerResponse", &DeleteScannerResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteScheduleResponse", &DeleteScheduleResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteTagResponse", &DeleteTagResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteTargetResponse", &DeleteTargetResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteTaskResponse", &DeleteTaskResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteTicketResponse", &DeleteTicketResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DeleteUserResponse", &DeleteUserResponse{Status: "200", StatusText: "OK"}, "200", "OK"},

		// Get operations
		{"GetAgentsResponse", &GetAgentsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetAggregatesResponse", &GetAggregatesResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetAlertsResponse", &GetAlertsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetConfigsResponse", &GetConfigsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetCredentialsResponse", &GetCredentialsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetFeaturesResponse", &GetFeaturesResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetFeedsResponse", &GetFeedsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetFiltersResponse", &GetFiltersResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetGroupsResponse", &GetGroupsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetInfoResponse", &GetInfoResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetLicenseResponse", &GetLicenseResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetNotesResponse", &GetNotesResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetNVTFamiliesResponse", &GetNVTFamiliesResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetNVTsResponse", &GetNVTsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetOverridesResponse", &GetOverridesResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetPermissionsResponse", &GetPermissionsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetPortListsResponse", &GetPortListsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetPreferencesResponse", &GetPreferencesResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetReportConfigsResponse", &GetReportConfigsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetReportFormatsResponse", &GetReportFormatsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetReportsResponse", &GetReportsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetResourceNamesResponse", &GetResourceNamesResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetRolesResponse", &GetRolesResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetScannersResponse", &GetScannersResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetSchedulesResponse", &GetSchedulesResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetSettingsResponse", &GetSettingsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetSystemReportsResponse", &GetSystemReportsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetTagsResponse", &GetTagsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetTLSCertificatesResponse", &GetTLSCertificatesResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetUsersResponse", &GetUsersResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetVersionResponse", &GetVersionResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"GetVulnsResponse", &GetVulnsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"HelpResponse", &HelpResponse{Status: "200", StatusText: "OK"}, "200", "OK"},

		// Modify operations
		{"ModifyAgentsResponse", &ModifyAgentsResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyAlertResponse", &ModifyAlertResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyAssetResponse", &ModifyAssetResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyAuthResponse", &ModifyAuthResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyConfigResponse", &ModifyConfigResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyCredentialResponse", &ModifyCredentialResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyFilterResponse", &ModifyFilterResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyGroupResponse", &ModifyGroupResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyLicenseResponse", &ModifyLicenseResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyNoteResponse", &ModifyNoteResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyOverrideResponse", &ModifyOverrideResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyPermissionResponse", &ModifyPermissionResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyPortListResponse", &ModifyPortListResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyReportConfigResponse", &ModifyReportConfigResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyReportFormatResponse", &ModifyReportFormatResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyRoleResponse", &ModifyRoleResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyScannerResponse", &ModifyScannerResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyScheduleResponse", &ModifyScheduleResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifySettingResponse", &ModifySettingResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyTagResponse", &ModifyTagResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyTargetResponse", &ModifyTargetResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyTaskResponse", &ModifyTaskResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyTicketResponse", &ModifyTicketResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyTLSCertificateResponse", &ModifyTLSCertificateResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ModifyUserResponse", &ModifyUserResponse{Status: "200", StatusText: "OK"}, "200", "OK"},

		// Other operations
		{"MoveTaskResponse", &MoveTaskResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"RestoreResponse", &RestoreResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"ResumeTaskResponse", &ResumeTaskResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"StartTaskResponse", &StartTaskResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"StopTaskResponse", &StopTaskResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"SyncConfigResponse", &SyncConfigResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"TestAlertResponse", &TestAlertResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"VerifyReportFormatResponse", &VerifyReportFormatResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"VerifyScannerResponse", &VerifyScannerResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"EmptyTrashcanResponse", &EmptyTrashcanResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"DescribeAuthResponse", &DescribeAuthResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"RunWizardResponse", &RunWizardResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
		{"WizardResponse", &WizardResponse{Status: "200", StatusText: "OK"}, "200", "OK"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, statusText := tc.response.GetStatus()
			if status != tc.status {
				t.Errorf("Expected status %q, got %q", tc.status, status)
			}
			if statusText != tc.statusText {
				t.Errorf("Expected statusText %q, got %q", tc.statusText, statusText)
			}
		})
	}
}

// Test error scenarios for GetStatus methods.
func TestGetStatusMethods_ErrorScenarios(t *testing.T) {
	testCases := []struct {
		name       string
		response   interface{ GetStatus() (string, string) }
		status     string
		statusText string
	}{
		{"AuthenticateResponse_Error", &AuthenticateResponse{Status: "401", StatusText: "Unauthorized"}, "401", "Unauthorized"},
		{"CreateAssetResponse_Error", &CreateAssetResponse{Status: "400", StatusText: "Bad Request"}, "400", "Bad Request"},
		{"GetTasksResponse_Error", &GetTasksResponse{Status: "404", StatusText: "Not Found"}, "404", "Not Found"},
		{"GetResultsResponse_Error", &GetResultsResponse{Status: "500", StatusText: "Internal Server Error"}, "500", "Internal Server Error"},
		{"GetAssetsResponse_Error", &GetAssetsResponse{Status: "403", StatusText: "Forbidden"}, "403", "Forbidden"},
		{"GetTargetsResponse_Error", &GetTargetsResponse{Status: "503", StatusText: "Service Unavailable"}, "503", "Service Unavailable"},
		{"GetTicketsResponse_Error", &GetTicketsResponse{Status: "502", StatusText: "Bad Gateway"}, "502", "Bad Gateway"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, statusText := tc.response.GetStatus()
			if status != tc.status {
				t.Errorf("Expected status %q, got %q", tc.status, status)
			}
			if statusText != tc.statusText {
				t.Errorf("Expected statusText %q, got %q", tc.statusText, statusText)
			}
		})
	}
}
