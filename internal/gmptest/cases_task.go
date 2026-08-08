package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// CreateTask.Alterable is *bool: unset, true, and an explicit false all
	// round-trip distinctly.
	RegisterRequestCase("create_task", "alterable_true", &commands.CreateTask{
		Name:      "Weekly external scan",
		Comment:   "Runs every Monday against the DMZ",
		Alterable: commands.Bool(true),
		Config:    &commands.TaskConfig{ID: "daba56c8-73ec-11df-a475-002264764cea"},
		Target:    &commands.TaskTarget{ID: "1a2b3c4d-5e6f-7081-92a3-b4c5d6e7f8a9"},
		Scanner:   &commands.TaskScanner{ID: "08b69003-5fc2-4037-a479-93b440211c73"},
	})

	RegisterRequestCase("create_task", "alterable_unset", &commands.CreateTask{
		Name:    "One-off scan",
		Config:  &commands.TaskConfig{ID: "daba56c8-73ec-11df-a475-002264764cea"},
		Target:  &commands.TaskTarget{ID: "1a2b3c4d-5e6f-7081-92a3-b4c5d6e7f8a9"},
		Scanner: &commands.TaskScanner{ID: "08b69003-5fc2-4037-a479-93b440211c73"},
	})

	RegisterRequestCase("create_task", "alterable_false", &commands.CreateTask{
		Name:      "Locked scan",
		Alterable: commands.Bool(false),
		Config:    &commands.TaskConfig{ID: "daba56c8-73ec-11df-a475-002264764cea"},
		Target:    &commands.TaskTarget{ID: "1a2b3c4d-5e6f-7081-92a3-b4c5d6e7f8a9"},
		Scanner:   &commands.TaskScanner{ID: "08b69003-5fc2-4037-a479-93b440211c73"},
	})

	// GetTasks.Details is *bool: unset, true, and an explicit false all
	// round-trip distinctly. Taken verbatim from the plan's own worked example.
	RegisterRequestCase("get_tasks", "with_filter", &commands.GetTasks{Filter: "rows=2 first=1", Details: commands.Bool(true)})
	RegisterRequestCase("get_tasks", "by_id", &commands.GetTasks{TaskID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6"})
	RegisterRequestCase("get_tasks", "with_filter_false", &commands.GetTasks{Filter: "rows=2 first=1", Details: commands.Bool(false)})

	// DeleteTask.Ultimate is *bool: unset, true, and an explicit false all
	// round-trip distinctly.
	RegisterRequestCase("delete_task", "ultimate_true", &commands.DeleteTask{
		TaskID:   "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Ultimate: commands.Bool(true),
	})

	RegisterRequestCase("delete_task", "ultimate_unset", &commands.DeleteTask{
		TaskID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
	})

	RegisterRequestCase("delete_task", "ultimate_false", &commands.DeleteTask{
		TaskID:   "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Ultimate: commands.Bool(false),
	})
}
