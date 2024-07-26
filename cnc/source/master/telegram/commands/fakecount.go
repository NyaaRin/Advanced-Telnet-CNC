package command

import (
	"advanced-telnet-cnc/source/niggers"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {
	Make(&Command{
		Aliases: []string{"fakecounter", "fakecount"},
		Executor: func(args []string, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
			niggers.FakeCounting = !niggers.FakeCounting

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Fake Counting: %v", niggers.FakeCounting))
			msg.ReplyToMessageID = update.Message.MessageID
			_, err := bot.Send(msg)
			return err
		},
	})
}
