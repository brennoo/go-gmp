package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	RegisterRequestCase("authenticate", "basic", &commands.Authenticate{
		Credentials: &commands.AuthenticateCredentials{
			Username: "admin",
			Password: "admin",
		},
	})
}
