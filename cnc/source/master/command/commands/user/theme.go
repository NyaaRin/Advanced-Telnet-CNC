package user

import (
	"advanced-telnet-cnc/packages/simpletable"
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"advanced-telnet-cnc/source/master/termfx/packages"
	"fmt"
	"strings"
)

func init() {
	command.Make(&command.Command{
		Aliases:     []string{"theme", "themes"},
		Description: "Sets themes",
		Executor: func(args []string, term *termfx.TermFX, session *sessions.Session) error {
			if len(args) < 1 {
				table := simpletable.New()

				table.Header = &simpletable.Header{
					Cells: []*simpletable.Cell{
						{Align: simpletable.AlignCenter, Text: "#"},
						{Align: simpletable.AlignCenter, Text: "Name"},
						{Align: simpletable.AlignCenter, Text: "Preview"},
						{Align: simpletable.AlignCenter, Text: "Primary"},
						{Align: simpletable.AlignCenter, Text: "Secondary"},
					},
				}

				for i, profile := range config.Themes.Themes {
					r := []*simpletable.Cell{
						{Align: 0, Text: fmt.Sprintf("%d", i+1)},
						{Align: 0, Text: profile.Name},
						{Align: 0, Text: fmt.Sprintf("%s    %s    \x1b[0m", packages.ColorizeBG(profile.Primary), packages.ColorizeBG(profile.Secondary))},
						{Align: 0, Text: fmt.Sprintf("%s", profile.Primary)},
						{Align: 0, Text: fmt.Sprintf("%s", profile.Secondary)},
					}

					table.Body.Cells = append(table.Body.Cells, r)
				}

				table.SetStyle(simpletable.StyleCompactLite)
				return session.Println(" " + strings.ReplaceAll(table.String(), "\n", "\r\n "))
			}

			theme := config.ThemeByName(args[0])
			if theme == nil {
				return session.Println("Theme does not exist.")
			}

			session.Theme = theme
			return session.Printf("Set theme to %s\r\n", session.Theme.Name)
		},
	})
}
