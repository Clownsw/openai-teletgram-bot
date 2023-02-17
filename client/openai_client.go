package client

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"net/http"
	"openai-teletgram-bot/config"
	"openai-teletgram-bot/domain/openai"
)

type OpenAIClient struct {
	Client    *req.Client
	ApiDomain string
}

func (openAIClient *OpenAIClient) Query(query string) (*openai.Answers, error) {
	question := openai.NewWithDefault("text-davinci-003", query, "bigduu")
	endpoint := fmt.Sprintf("%s/v1/completions", openAIClient.ApiDomain)

	response, err := openAIClient.Client.R().SetBodyJsonMarshal(question).Post(endpoint)
	if err != nil {
		return nil, err
	}

	answers := openai.Answers{}
	err = json.Unmarshal([]byte(response.String()), &answers)
	if err != nil {
		return nil, err
	}

	return &answers, nil
}

func NewOpenAIClient(openAIInfo *config.OpenAIInfo) *OpenAIClient {
	openAIClient := OpenAIClient{}
	openAIClient.Client = req.NewClient()

	openAIClient.Client.Headers = http.Header{}
	openAIClient.Client.Headers.Add("Authorization", fmt.Sprintf("Bearer %s", openAIInfo.Token))
	openAIClient.Client.Headers.Add("Content-Type", "application/json")

	openAIClient.ApiDomain = openAIInfo.BaseUrl

	return &openAIClient
}
