package client

import (
	"errors"
	"fmt"
	"openai-teletgram-bot/config"
	"openai-teletgram-bot/domain/openai"
	"openai-teletgram-bot/util"

	"github.com/bytedance/sonic"
	"github.com/sirupsen/logrus"
)

type GptClient struct {
	Config *ClientConfig
	Logger *logrus.Logger
}

func (gptClient *GptClient) Query(query string) (string, error) {
	message := new(openai.GptMessage)
	message.Role = "user"
	message.Content = query

	question := openai.NewGptWithDefault(
		"gpt-3.5-turbo",
		[]*openai.GptMessage{message},
		"bigduu",
	)
	endpoint := fmt.Sprintf("%s/v1/chat/completions", gptClient.Config.ApiDomain)

	response, err := gptClient.Config.ReqClient.R().SetBodyJsonMarshal(question).Post(endpoint)
	if err != nil {
		gptClient.Logger.Error(fmt.Sprintf("gpt_client-Query: err: %s", err))
		return config.EmptyString, err
	}

	gptClient.Logger.Error(fmt.Sprintf("gpt_client-Query: response: %s", response.String()))

	choicesNode, err := sonic.Get(util.StringToByteSlice(response.String()), "choices")
	if err != nil {
		return config.EmptyString, err
	}

	choices, err := choicesNode.ArrayUseNode()
	if err != nil {
		return config.EmptyString, err
	}

	if len(choices) == 0 {
		return config.EmptyString, errors.New("choices size == 0")
	}

	messageNode := choices[0].Get("message")
	if messageNode != nil {
		contentNode := messageNode.Get("content")
		if contentNode != nil {
			return contentNode.String()
		}
	}

	return config.EmptyString, errors.New("unknown error")
}

func (gptClient *GptClient) SetConfig(config *ClientConfig) {
	gptClient.Config = config
}

func (gptClient *GptClient) SetLoggger(logger *logrus.Logger) {
	gptClient.Logger = logger
}
