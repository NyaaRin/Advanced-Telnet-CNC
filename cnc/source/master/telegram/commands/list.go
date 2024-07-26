package command

import (
	"advanced-telnet-cnc/source/niggers"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {
	Make(&Command{
		Aliases: []string{"list", "botlist"},
		Executor: func(args []string, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
			var clients string

			clients += "Here is a list of all connected clients at the moment\r\n"

			for k, v := range niggers.Distribution() {
				clients += " - " + k + ": " + fmt.Sprint(v) + "\r\n"
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(clients+"The number of currently connected devices is %d.", niggers.Count()))
			msg.ReplyToMessageID = update.Message.MessageID
			_, err := bot.Send(msg)
			return err
		},
	})
}
