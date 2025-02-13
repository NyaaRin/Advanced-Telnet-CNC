package command

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

var (
	// commands - Commands map, where all the commands are saved.
	commands = make(map[string]*Command)
)

// Command - Structure of the command
type Command struct {
	Aliases    []string
	AllowedIDs []int64
	Executor   func(args []string, bot *tgbotapi.BotAPI, update tgbotapi.Update) error
}


// Make - Registers the command
func Make(command *Command) {
	if _, exists := commands[command.Aliases[0]]; exists {
		log.Fatal("command already exists")
	}

	commands[command.Aliases[0]] = command
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

// Clone - Gets all cmdImpl in a slice
func Clone() []*Command {
	var commandSlice []*Command
	for _, cmd := range commands {
		commandSlice = append(commandSlice, cmd)
	}

	return commandSlice
}
