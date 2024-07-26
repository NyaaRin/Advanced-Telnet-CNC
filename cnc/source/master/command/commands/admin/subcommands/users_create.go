package subcommands

import (
	"advanced-telnet-cnc/packages/duration"
	"advanced-telnet-cnc/packages/pflag"
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/database"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func init() {
	command.Subcommand("users", &command.SubCommand{
		Aliases:   []string{"create", "add"},
		Arguments: []command.Argument{{"username", false}, {"password", false}, {"options", true}},
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session, parent *command.Command) error {
			var user = &database.UserProfile{
				Name:       args[0],
				Password:   args[1],
				Methods:    []int{-1},
				Cooldown:   60,
				MaxTime:    60,
				MaxAttacks: 50,
				Devices:    -1,
				Expiry:     time.Now().Add(24 * time.Hour),
				Admin:      false,
				Reseller:   false,
			}

			var expiry string

			flagSet := pflag.NewFlagSet("users create", pflag.ContinueOnError)
			flagSet.SetOutput(session.Conn)

			flagSet.StringVarP(&expiry, "expiry", "e", "1d", "Expiry of the users plan")

			// Only admins are allowed to use the admin & reseller flags.
			if user.Admin {
				flagSet.BoolVarP(&user.Admin, "admin", "a", user.Admin, "Allows the user to execute administrative commands")
				flagSet.BoolVarP(&user.Reseller, "reseller", "r", user.Reseller, "Allows the user to execute commands like 'users add'")
			}

			flagSet.IntVarP(&user.Devices, "devices", "d", user.Devices, "The devices the user has access to")
			flagSet.IntVarP(&user.Cooldown, "cooldown", "c", user.Cooldown, "The cooldown the user has after an attack")
			flagSet.IntVarP(&user.MaxTime, "time", "t", user.MaxTime, "The maximum attack duration the user has")
			flagSet.IntVarP(&user.MaxAttacks, "attacks", "f", user.MaxAttacks, "The maximum attacks the user has per day")
			flagSet.IntSliceVarP(&user.Methods, "methods", "m", user.Methods, "Floods the user has access to")

			err := flagSet.Parse(args[2:])
			if err != nil {
				if err == pflag.ErrHelp {
					return nil
				}

				return session.CmdError(err)
			}

			durationParsed, err := duration.ModifiedParseDuration(expiry)
			if err != nil {
				return session.CmdError(err)
			}

			user.Expiry = time.Now().Add(durationParsed)

			err = database.CreateUser(user, true)
			if err != nil {
				if err == database.ErrKnownUser {
					return session.CmdError(errors.New(fmt.Sprintf("failed to add user to the database because a user with the name %s already exists.", strconv.Quote(user.Name))))
				}

				return session.CmdError(err)
			}

			return session.Printf("%sSuccessfully added %s into the database.\r\n", config.Green, strconv.Quote(user.Name))
		},
	})
}
