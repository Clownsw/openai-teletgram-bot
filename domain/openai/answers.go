package openai

import (
	"encoding/json"
	"strings"
)

type Choice struct {
	Text         string     `json:"text"`
	Index        int64      `json:"index"`
	Logprobs     json.Token `json:"logprobs"`
	FinishReason string     `json:"finish_reason"`
}

func (choice *Choice) ToText() string {
	splitArray := strings.Split(choice.Text, "\n")
	if len(splitArray) == 3 {
		return splitArray[2]
	}
	return choice.Text
}

type Usage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type Answers struct {
	Id      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Choices []*Choice `json:"choices"`
	Usage   *Usage    `json:"usage"`
	Error   *Error    `json:"error"`
}
