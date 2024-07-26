package admin

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
	command.Make(&command.Command{
		Aliases:     []string{"running", "ongoing"},
		Description: "Displays a list of all ongoing attacks.",
		Admin:       false,
		Reseller:    true,
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session) error {
			table := simpletable.New()

			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "#"},
					{Align: simpletable.AlignCenter, Text: "Username"},
					{Align: simpletable.AlignCenter, Text: "Target"},
					{Align: simpletable.AlignCenter, Text: "Method"},
					{Align: simpletable.AlignCenter, Text: "Duration"},
					{Align: simpletable.AlignCenter, Text: "Created"},
					{Align: simpletable.AlignCenter, Text: "End"},
				},
			}

			floods, err := database.RunningAttacks()
			if err != nil {
				return err
			}

			for i, row := range floods {
				user, err2 := database.UserFromId(row.UserId)
				if err2 != nil {
					return err2
				}

				r := []*simpletable.Cell{
					{Align: 0, Text: fmt.Sprintf("%d", i+1)},
					{Align: 0, Text: user.Name},
					{Align: 0, Text: row.Target},
					{Align: 0, Text: row.Method},
					{Align: 0, Text: fmt.Sprint(row.Duration)},
					{Align: 0, Text: format.Format(row.Created)},
					{Align: 0, Text: "In " + time.Until(row.End).String()},
				}

				table.Body.Cells = append(table.Body.Cells, r)
			}

			table.SetStyle(simpletable.StyleCompactLite)

			return session.Println(" " + strings.ReplaceAll(table.String(), "\n", "\r\n "))
		},
	})
}
