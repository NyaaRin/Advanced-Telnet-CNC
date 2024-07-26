package command

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {
	Make(&Command{
		Aliases: []string{"chatid", "id"},
		Executor: func(args []string, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("[*] Chat ID: %d", update.Message.Chat.ID))
			msg.ReplyToMessageID = update.Message.MessageID
			_, err := bot.Send(msg)
			return err
		},
	})
}
