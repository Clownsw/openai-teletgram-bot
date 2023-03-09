package main

import (
	"fmt"
	"openai-teletgram-bot/config"
	"openai-teletgram-bot/server"
	"sync/atomic"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	newServer := server.NewServer(func(server *server.Server, update *tgbotapi.Update) (bool, error) {
		text := update.Message.Text[5:]
		result, err := server.Api.Client.Query(text)

		if err != nil {
			return false, err
		}

		server.Logger.Info(fmt.Sprintf("[%d-%s] query reply: ", update.Message.From.ID, update.Message.From.UserName), result)

		err = server.SendMessage(
			update.Message.Chat.ID,
			result,
			&update.Message.MessageID,
		)
		if err != nil {
			return false, err
		}

		return true, nil
	}, func(server *server.Server, update tgbotapi.Update) {
		counter := new(atomic.Int32)

		for true {
			result, err := server.CallBack(server, &update)
			if result {
				break
			}

			server.Logger.Info(
				fmt.Sprintf(
					"messageId: %d, form: %d, counter: %d",
					update.Message.MessageID,
					update.Message.From.ID,
					counter.Load(),
				),
			)

			if counter.Load() >= server.Api.MaxReplyCount {
				_ = server.SendMessage(
					update.Message.Chat.ID,
					config.ErrorMessage,
					&update.Message.MessageID,
				)
				break
			}

			if server.Api.SendErrorMessage {
				_ = server.SendMessage(
					server.Bot.Owner,
					fmt.Sprintf(
						"[%d-%s] send message: %s, err: %s",
						update.Message.From.ID,
						update.Message.From.UserName,
						update.Message.Text, err.Error(),
					),
					nil,
				)
			}

			counter.Add(1)
		}
	})
	newServer.App()
}
