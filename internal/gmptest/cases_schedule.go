package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// GetSchedules.Trash is *GMPBool: unset, true, and an explicit false all
	// round-trip distinctly.
	RegisterRequestCase("get_schedules", "filter_and_trash_true", &commands.GetSchedules{
		Filter: "rows=-1",
		FiltID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
		Trash:  commands.Bool(true),
	})

	RegisterRequestCase("get_schedules", "by_id", &commands.GetSchedules{
		ScheduleID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
	})

	RegisterRequestCase("get_schedules", "filter_and_trash_false", &commands.GetSchedules{
		Filter: "rows=-1",
		FiltID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
		Trash:  commands.Bool(false),
	})

	// GetSchedules.Tasks is *GMPBool: unset, true, and an explicit false all
	// round-trip distinctly.
	RegisterRequestCase("get_schedules", "tasks_true", &commands.GetSchedules{
		Tasks: commands.Bool(true),
	})

	RegisterRequestCase("get_schedules", "tasks_false", &commands.GetSchedules{
		Tasks: commands.Bool(false),
	})
}
