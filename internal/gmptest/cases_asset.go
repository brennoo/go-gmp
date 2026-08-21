package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	RegisterRequestCase("create_asset", "asset", &commands.CreateAsset{
		Asset: &commands.Asset{
			Name:    "192.168.1.1",
			Type:    "host",
			Comment: "Imported manually",
		},
	})

	RegisterRequestCase("create_asset", "report", &commands.CreateAsset{
		Report: &commands.AssetReport{
			ID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
			Filter: commands.AssetFilter{
				Term: "apply_overrides=0 levels=hml",
			},
		},
	})

	// GetAssets.IgnorePagination/Details are *GMPBool: unset, true, and an
	// explicit false all round-trip distinctly.
	RegisterRequestCase("get_assets", "unset", &commands.GetAssets{})

	RegisterRequestCase("get_assets", "by_id", &commands.GetAssets{
		AssetID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("get_assets", "filt_id", &commands.GetAssets{
		FiltID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("get_assets", "all_details_true", &commands.GetAssets{
		Filter:           "rows=-1",
		IgnorePagination: commands.Bool(true),
		Details:          commands.Bool(true),
	})

	RegisterRequestCase("get_assets", "all_details_false", &commands.GetAssets{
		Filter:           "rows=-1",
		IgnorePagination: commands.Bool(false),
		Details:          commands.Bool(false),
	})

	// ModifyAsset.Comment is cleared whenever omitted, not left unchanged, so
	// these cases pin the set and omitted wire shapes side by side.
	RegisterRequestCase("modify_asset", "comment_set", &commands.ModifyAsset{
		AssetID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Comment: "Re-commented after review",
	})

	RegisterRequestCase("modify_asset", "comment_omitted", &commands.ModifyAsset{
		AssetID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
	})
}
