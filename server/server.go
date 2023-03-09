package server

import (
	"openai-teletgram-bot/client"
	"openai-teletgram-bot/config"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type CallBack func(server *Server, update *tgbotapi.Update) (bool, error)
type CallBackRetry func(server *Server, update tgbotapi.Update)

type Api struct {
	Token            string `yaml:"token"`
	BaseUrl          string `yaml:"baseUrl"`
	Client           client.Client
	MaxReplyCount    int32  `yaml:"maxReplyCount"`
	SendErrorMessage bool   `yaml:"sendErrorMessage"`
	Model            string `yaml:"model"`
}

type Bot struct {
	Token string `json:"token"`
	Bot   *tgbotapi.BotAPI
	Owner int64
}

type Log struct {
	Level string `json:"level"`
}

type Server struct {
	Api           Api `json:"api"`
	Bot           Bot `json:"bot"`
	Log           Log `json:"log"`
	Logger        *logrus.Logger
	CallBack      CallBack
	CallBackRetry CallBackRetry
}

func (server *Server) App() {
	bot, err := tgbotapi.NewBotAPI(server.Bot.Token)
	if err != nil {
		panic(err)
	}

	server.Bot.Bot = bot

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		panic(err)
	}

	server.Logger.Info("bot start success!")

	for update := range updates {
		if update.Message != nil {
			server.Logger.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if len(update.Message.Text) > 5 && update.Message.Text[0:5] == "look " {
				go server.CallBackRetry(server, update)
			}
		}
	}
}

func (server *Server) SendMessage(chatId int64, msg string, replyMessageId *int) error {
	message := tgbotapi.NewMessage(chatId, msg)
	message.ParseMode = tgbotapi.ModeMarkdown

	if replyMessageId != nil {
		message.ReplyToMessageID = *replyMessageId
	}

	_, err := server.Bot.Bot.Send(message)
	if err != nil {

		// Telegram Message Markdown error
		if strings.Index(err.Error(), "can't parse entities") != -1 {
			message.ParseMode = ""
			_, _ = server.Bot.Bot.Send(message)
			return nil
		}
		return err
	}
	return nil
}

func NewServer(callback CallBack, callBackRetry CallBackRetry) *Server {
	args := os.Args
	if len(args) != 2 {
		panic("not found config file")
	}

	configFile := args[1]
	fileContent, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	server := new(Server)

	err = yaml.Unmarshal(fileContent, &server)
	if err != nil {
		panic(err)
	}

	server.Logger = logrus.New()

	level, err := logrus.ParseLevel(server.Log.Level)
	if err != nil {
		panic(err)
	}

	server.Api.Client = client.NewClient(server.Api.Model, &config.AuthInfo{Token: server.Api.Token, BaseUrl: server.Api.BaseUrl})
	server.Logger.SetLevel(level)
	server.Api.Client.SetLoggger(server.Logger)
	server.CallBack = callback
	server.CallBackRetry = callBackRetry
	return server
}
