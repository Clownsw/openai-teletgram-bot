package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"openai-teletgram-bot/config"
	"openai-teletgram-bot/server"
	"sync/atomic"
)

func main() {
	newServer := server.NewServer(func(server *server.Server, update *tgbotapi.Update) (bool, error) {
		text := update.Message.Text[5:]
		answers, err := server.Api.OpenAIClient.Query(text)

		if err != nil {
			return false, err
		}

		server.Logger.Info(fmt.Sprintf("[%d-%s] query reply: ", update.Message.From.ID, update.Message.From.UserName), answers)

		if len(answers.Choices) < 1 {
			return false, config.OpenAiQueryError
		}

		server.SendMessage(
			update.Message.Chat.ID,
			answers.Choices[0].ToText(),
			&update.Message.MessageID,
		)

		return true, nil
	}, func(server *server.Server, update tgbotapi.Update) {
		counter := new(atomic.Int32)

		for true {
			result, err := server.CallBack(server, &update)
			if result {
				break
			}

			if counter.Load() >= server.Api.MaxReplyCount {
				server.SendMessage(
					update.Message.Chat.ID,
					config.ErrorMessage,
					&update.Message.MessageID,
				)
				break
			}

			if server.Api.SendErrorMessage {
				server.SendMessage(
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
