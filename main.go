package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"openai-teletgram-bot/config"
	"openai-teletgram-bot/server"
)

func main() {
	newServer := server.NewServer(func(server *server.Server, update *tgbotapi.Update) {
		text := update.Message.Text[5:]
		answers, err := server.Api.OpenAIClient.Query(text)

		server.Logger.Info(fmt.Sprintf("[%d-%s] query reply: ", update.Message.From.ID, update.Message.From.UserName), answers)

		if err != nil {
			server.SendMessage(update.Message.Chat.ID, config.ErrorMessage, &update.Message.MessageID)
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
		} else {
			if len(answers.Choices) < 1 {
				server.SendMessage(update.Message.Chat.ID, config.ErrorMessage, &update.Message.MessageID)
				server.SendMessage(
					server.Bot.Owner,
					fmt.Sprintf(
						"[%d-%s] send message: %s, err: %s",
						update.Message.From.ID,
						update.Message.From.UserName,
						update.Message.Text,
						answers.Error.Message,
					),
					nil,
				)
				return
			}

			server.SendMessage(
				update.Message.Chat.ID,
				answers.Choices[0].ToText(),
				&update.Message.MessageID,
			)
		}
	})
	newServer.App()
}
