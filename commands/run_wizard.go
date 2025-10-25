package commands

import (
	"encoding/xml"
)

// RunWizard represents a run_wizard command request.
type RunWizard struct {
	XMLName  xml.Name      `xml:"run_wizard"`
	Mode     string        `xml:"mode,omitempty"`
	Name     string        `xml:"name"`
	Params   *WizardParams `xml:"params,omitempty"`
	ReadOnly bool          `xml:"read_only,attr,omitempty"`
}

// RunWizardResponse represents a run_wizard command response.
type RunWizardResponse struct {
	XMLName    xml.Name        `xml:"run_wizard_response"`
	Status     string          `xml:"status,attr"`
	StatusText string          `xml:"status_text,attr"`
	Response   *WizardResponse `xml:"response,omitempty"`
}

type WizardParams struct {
	Params []WizardParam `xml:"param"`
}

type WizardParam struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

type WizardResponse struct {
	CreateTargetResponse     *CreateTargetResponse     `xml:"create_target_response,omitempty"`
	CreateTaskResponse       *CreateTaskResponse       `xml:"create_task_response,omitempty"`
	StartTaskResponse        *StartTaskResponse        `xml:"start_task_response,omitempty"`
	GetTasksResponse         *GetTasksResponse         `xml:"get_tasks_response,omitempty"`
	DeleteTaskResponse       *DeleteTaskResponse       `xml:"delete_task_response,omitempty"`
	DeleteTargetResponse     *DeleteTargetResponse     `xml:"delete_target_response,omitempty"`
	DeleteScheduleResponse   *DeleteScheduleResponse   `xml:"delete_schedule_response,omitempty"`
	CreateAlertResponse      *CreateAlertResponse      `xml:"create_alert_response,omitempty"`
	CreateCredentialResponse *CreateCredentialResponse `xml:"create_credential_response,omitempty"`
	DeleteCredentialResponse *DeleteCredentialResponse `xml:"delete_credential_response,omitempty"`
	CreateScheduleResponse   *CreateScheduleResponse   `xml:"create_schedule_response,omitempty"`
	ModifyTaskResponse       *ModifyTaskResponse       `xml:"modify_task_response,omitempty"`
	DeleteReportResponse     *DeleteReportResponse     `xml:"delete_report_response,omitempty"`
	GetVersionResponse       *GetVersionResponse       `xml:"get_version_response,omitempty"`
	GetCredentialsResponse   *GetCredentialsResponse   `xml:"get_credentials_response,omitempty"`
	GetReportsResponse       *GetReportsResponse       `xml:"get_reports_response,omitempty"`
	GetReportFormatsResponse *GetReportFormatsResponse `xml:"get_report_formats_response,omitempty"`
	GetConfigsResponse       *GetConfigsResponse       `xml:"get_configs_response,omitempty"`
	GetScannersResponse      *GetScannersResponse      `xml:"get_scanners_response,omitempty"`
	GetTargetsResponse       *GetTargetsResponse       `xml:"get_targets_response,omitempty"`
	GetUsersResponse         *GetUsersResponse         `xml:"get_users_response,omitempty"`
	GetTicketsResponse       *GetTicketsResponse       `xml:"get_tickets_response,omitempty"`
	GetNotesResponse         *GetNotesResponse         `xml:"get_notes_response,omitempty"`
	GetGroupsResponse        *GetGroupsResponse        `xml:"get_groups_response,omitempty"`
	GetTagsResponse          *GetTagsResponse          `xml:"get_tags_response,omitempty"`
	GetAssetsResponse        *GetAssetsResponse        `xml:"get_assets_response,omitempty"`
	GetAlertsResponse        *GetAlertsResponse        `xml:"get_alerts_response,omitempty"`
	GetInfoResponse          *GetInfoResponse          `xml:"get_info_response,omitempty"`
	GetVulnsResponse         *GetVulnsResponse         `xml:"get_vulns_response,omitempty"`
	GetResultsResponse       *GetResultsResponse       `xml:"get_results_response,omitempty"`
	GetPortListsResponse     *GetPortListsResponse     `xml:"get_port_lists_response,omitempty"`
	GetReportConfigsResponse *GetReportConfigsResponse `xml:"get_report_configs_response,omitempty"`
	GetSystemReportsResponse *GetSystemReportsResponse `xml:"get_system_reports_response,omitempty"`
	GetResourceNamesResponse *GetResourceNamesResponse `xml:"get_resource_names_response,omitempty"`
	GetPreferencesResponse   *GetPreferencesResponse   `xml:"get_preferences_response,omitempty"`
	GetFeaturesResponse      *GetFeaturesResponse      `xml:"get_features_response,omitempty"`
	GetLicenseResponse       *GetLicenseResponse       `xml:"get_license_response,omitempty"`
	GetAggregatesResponse    *GetAggregatesResponse    `xml:"get_aggregates_response,omitempty"`
	Status                   string                    `xml:"status,omitempty"`
	StatusText               string                    `xml:"status_text,omitempty"`
	InnerXML                 string                    `xml:",innerxml"` // For any unknown or extra content
}
