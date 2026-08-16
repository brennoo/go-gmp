package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	RegisterRequestCase("create_report_config", "name", &commands.CreateReportConfig{
		Name:         "Test Config",
		ReportFormat: commands.ReportConfigReportFormat{ID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f"},
	})

	RegisterRequestCase("create_report_config", "copy", &commands.CreateReportConfig{
		Name:         "Test Config",
		ReportFormat: commands.ReportConfigReportFormat{ID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f"},
		Copy:         "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("create_report_config", "param_set", &commands.CreateReportConfig{
		Name:         "Test Config",
		ReportFormat: commands.ReportConfigReportFormat{ID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f"},
		Comment:      "Test comment",
		Param: []commands.ReportConfigParam{{
			Name:  "Node Distance",
			Value: commands.ReportConfigValue{Value: "8", UseDefault: commands.Bool(false)},
		}},
	})

	// GetReportConfigs' Trash/Details are *GMPBool: unset, true, and an
	// explicit false all round-trip distinctly.
	RegisterRequestCase("get_report_configs", "unset", &commands.GetReportConfigs{})

	RegisterRequestCase("get_report_configs", "by_id", &commands.GetReportConfigs{
		ReportConfigID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("get_report_configs", "filt_id", &commands.GetReportConfigs{
		FiltID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("get_report_configs", "all_details_true", &commands.GetReportConfigs{
		Filter:  "rows=-1",
		Trash:   commands.Bool(true),
		Details: commands.Bool(true),
	})

	RegisterRequestCase("get_report_configs", "all_details_false", &commands.GetReportConfigs{
		Filter:  "rows=-1",
		Trash:   commands.Bool(false),
		Details: commands.Bool(false),
	})

	// ModifyReportConfig.ReportConfigID is required in practice (see the
	// field's gmpoverride note), so every case sets it.
	RegisterRequestCase("modify_report_config", "name", &commands.ModifyReportConfig{
		ReportConfigID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Name:           "Renamed config",
	})

	RegisterRequestCase("modify_report_config", "param_set", &commands.ModifyReportConfig{
		ReportConfigID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Param: []commands.ReportConfigParam{{
			Name:  "Node Distance",
			Value: commands.ReportConfigValue{UseDefault: commands.Bool(true)},
		}},
	})

	RegisterRequestCase("modify_report_config", "param_omitted", &commands.ModifyReportConfig{
		ReportConfigID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Comment:        "Updated comment",
	})
}
