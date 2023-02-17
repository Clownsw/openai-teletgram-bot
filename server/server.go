package server

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"openai-teletgram-bot/client"
	"openai-teletgram-bot/config"
	"os"
)

type Api struct {
	Token        string `yaml:"token"`
	BaseUrl      string `yaml:"baseUrl"`
	OpenAIClient *client.OpenAIClient
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
	Api      Api `json:"api"`
	Bot      Bot `json:"bot"`
	Log      Log `json:"log"`
	Logger   *logrus.Logger
	CallBack func(server *Server, update *tgbotapi.Update)
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
				go server.CallBack(server, &update)
			}
		}
	}
}

func (server *Server) SendMessage(chatId int64, msg string, replyMessageId *int) {
	message := tgbotapi.NewMessage(chatId, msg)
	if replyMessageId != nil {
		message.ReplyToMessageID = *replyMessageId
	}

	_, err := server.Bot.Bot.Send(message)
	if err != nil {
		server.Logger.Error("send message err: ", err.Error())
	}
}

func NewServer(callback func(server *Server, update *tgbotapi.Update)) *Server {
	args := os.Args
	if len(args) != 2 {
		panic("not found config file")
	}

	configFile := args[1]
	fileContent, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	server := Server{}

	err = yaml.Unmarshal(fileContent, &server)
	if err != nil {
		panic(err)
	}

	server.Logger = logrus.New()

	level, err := logrus.ParseLevel(server.Log.Level)
	if err != nil {
		panic(err)
	}

	server.Api.OpenAIClient = client.NewOpenAIClient(config.NewOpenAIInfo(server.Api.Token, server.Api.BaseUrl))
	server.Logger.SetLevel(level)
	server.CallBack = callback
	return &server
}