package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	RegisterRequestCase("modify_user", "password_and_roles", &commands.ModifyUser{
		UserID:  "d290f1ee-6c54-4b01-90e6-d701748f0851",
		Name:    "jdoe",
		Comment: "Promoted to admin",
		Password: &commands.UserPassword{
			Modify: "1",
			Text:   "n3w-p@ssw0rd",
		},
		Roles: []*commands.UserRole{{ID: "7a10cb42-af9f-11e1-a5d2-406186ea4fc5"}},
		Hosts: &commands.UserHosts{Allow: "1", Value: "192.0.2.0/24"},
	})
}
