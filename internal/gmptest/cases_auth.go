package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	RegisterRequestCase("authenticate", "basic", &commands.Authenticate{
		Credentials: commands.AuthCredentials{
			Username: "admin",
			Password: "admin",
		},
	})
}
