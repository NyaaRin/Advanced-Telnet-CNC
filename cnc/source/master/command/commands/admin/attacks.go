package admin

import (
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
)

func init() {
	command.Make(&command.Command{
		Aliases:     []string{"attacks"},
		Description: "Enables or disables the execution of attacks.",
		Admin:       true,
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session) error {
			config.Master.Attacks = !config.Master.Attacks
			config.Rewrite("master.toml")

			if config.Master.Attacks {
				return session.Println(config.Green, "Attacks have been enabled.")
			}

			return session.Println(config.Red, "Attacks have been disabled.")
		},
	})
}
