package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// The subject-only and subject+resource shapes are both common in real
	// GMP requests (global vs. resource-scoped permissions); create_permission
	// is the clearest place in the corpus to pin both.
	RegisterRequestCase("create_permission", "subject_only", &commands.CreatePermission{
		Name: "get_tasks",
		Subject: commands.PermissionSubject{
			ID:   "d290f1ee-6c54-4b01-90e6-d701748f0851",
			Type: "user",
		},
	})

	RegisterRequestCase("create_permission", "subject_and_resource", &commands.CreatePermission{
		Name: "get_tasks",
		Subject: commands.PermissionSubject{
			ID:   "d290f1ee-6c54-4b01-90e6-d701748f0851",
			Type: "role",
		},
		Resource: &commands.PermissionResource{
			ID:   "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31",
			Type: "task",
		},
	})

	// ModifyPermission.Name/.Comment stay plain string,omitempty rather than
	// PointerOverrides: a live gvmd experiment could not distinguish omitted
	// from explicit-empty because every real modify_permission call crashes
	// the gvmd worker (see docs/review-2026-08.md's `permission adopted`
	// section). These cases pin the generated shape as it stands today —
	// populated Name/Comment, a populated Subject, and Resource both
	// populated and omitted — not an omit/explicit-empty/set matrix, since
	// that distinction cannot be represented by this field type.
	RegisterRequestCase("modify_permission", "name_comment_subject", &commands.ModifyPermission{
		PermissionID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Name:         "get_tasks",
		Comment:      "Renamed permission",
		Subject: &commands.PermissionSubject{
			ID:   "d290f1ee-6c54-4b01-90e6-d701748f0851",
			Type: "user",
		},
	})

	RegisterRequestCase("modify_permission", "subject_and_resource", &commands.ModifyPermission{
		PermissionID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Comment:      "Rescoped to a different resource",
		Subject: &commands.PermissionSubject{
			ID:   "d290f1ee-6c54-4b01-90e6-d701748f0851",
			Type: "role",
		},
		Resource: &commands.PermissionResource{
			ID:   "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31",
			Type: "task",
		},
	})

	RegisterRequestCase("modify_permission", "resource_omitted", &commands.ModifyPermission{
		PermissionID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Comment:      "Global permission, no resource",
		Subject: &commands.PermissionSubject{
			ID:   "d290f1ee-6c54-4b01-90e6-d701748f0851",
			Type: "user",
		},
	})
}
