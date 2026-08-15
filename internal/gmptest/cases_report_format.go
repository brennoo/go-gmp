package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	RegisterRequestCase("create_report_format", "name", &commands.CreateReportFormat{
		Name: "Test Format",
	})

	RegisterRequestCase("create_report_format", "copy", &commands.CreateReportFormat{
		Copy: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
	})

	// GetReportFormats' Trash/Alerts/Params/ReportConfigs/Details are all
	// *GMPBool: unset, true, and an explicit false all round-trip distinctly.
	RegisterRequestCase("get_report_formats", "unset", &commands.GetReportFormats{})

	RegisterRequestCase("get_report_formats", "by_id", &commands.GetReportFormats{
		ReportFormatID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("get_report_formats", "filt_id", &commands.GetReportFormats{
		FiltID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("get_report_formats", "all_details_true", &commands.GetReportFormats{
		Filter:        "rows=-1",
		Trash:         commands.Bool(true),
		Alerts:        commands.Bool(true),
		Params:        commands.Bool(true),
		ReportConfigs: commands.Bool(true),
		Details:       commands.Bool(true),
	})

	RegisterRequestCase("get_report_formats", "all_details_false", &commands.GetReportFormats{
		Filter:        "rows=-1",
		Trash:         commands.Bool(false),
		Alerts:        commands.Bool(false),
		Params:        commands.Bool(false),
		ReportConfigs: commands.Bool(false),
		Details:       commands.Bool(false),
	})

	// ModifyReportFormat.Param is a single element, unlike ModifyAsset's
	// scalar Comment, but the same set/omitted contrast applies.
	RegisterRequestCase("modify_report_format", "param_set", &commands.ModifyReportFormat{
		ReportFormatID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Param: &commands.ModifyReportFormatParam{
			Name:  "Background Colour",
			Value: "red",
		},
	})

	RegisterRequestCase("modify_report_format", "param_omitted", &commands.ModifyReportFormat{
		ReportFormatID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Name:           "New Name",
	})
}
