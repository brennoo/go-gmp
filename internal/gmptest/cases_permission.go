package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// The subject-only and subject+resource shapes are both common in real
	// GMP requests (global vs. resource-scoped permissions); create_permission
	// is the clearest place in the corpus to pin both.
	RegisterRequestCase("create_permission", "subject_only", &commands.CreatePermission{
		Name: "get_tasks",
		Subject: &commands.PermissionSubject{
			ID:   "d290f1ee-6c54-4b01-90e6-d701748f0851",
			Type: "user",
		},
	})

	RegisterRequestCase("create_permission", "subject_and_resource", &commands.CreatePermission{
		Name: "get_tasks",
		Subject: &commands.PermissionSubject{
			ID:   "d290f1ee-6c54-4b01-90e6-d701748f0851",
			Type: "role",
		},
		Resource: &commands.PermissionResource{
			ID:   "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31",
			Type: "task",
		},
	})
}
