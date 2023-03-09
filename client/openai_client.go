package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"openai-teletgram-bot/config"
	"openai-teletgram-bot/domain/openai"

	"github.com/sirupsen/logrus"
)

type OpenAIClient struct {
	Config *ClientConfig
	Logger *logrus.Logger
}

func (openAIClient *OpenAIClient) Query(query string) (string, error) {
	question := openai.NewOpenAIWithDefault("text-davinci-003", query, "bigduu")
	endpoint := fmt.Sprintf("%s/v1/completions", openAIClient.Config.ApiDomain)

	response, err := openAIClient.Config.ReqClient.R().SetBodyJsonMarshal(question).Post(endpoint)
	if err != nil {
		openAIClient.Logger.Error(fmt.Sprintf("openai_client-Query: err: %s", err))
		return config.EmptyString, err
	}

	openAIClient.Logger.Error(fmt.Sprintf("openai_client-Query: response: %s", response.String()))

	answers := openai.Answers{}
	err = json.Unmarshal([]byte(response.String()), &answers)
	if err != nil {
		return config.EmptyString, err
	}

	if answers.Error != nil {
		return config.EmptyString, errors.New(answers.Error.Message)
	}

	return answers.Choices[0].ToText(), nil
}

func (openAIClient *OpenAIClient) SetConfig(config *ClientConfig) {
	openAIClient.Config = config
}

func (openAIClient *OpenAIClient) SetLoggger(logger *logrus.Logger) {
	openAIClient.Logger = logger
}
