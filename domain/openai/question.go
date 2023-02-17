package openai

type Question struct {
	Model       string      `json:"model"`
	Prompt      string      `json:"prompt"`
	MaxTokens   int64       `json:"max_tokens"`
	Temperature int64       `json:"temperature"`
	TopP        int64       `json:"top_p"`
	N           int64       `json:"n"`
	Stream      bool        `json:"stream"`
	Logprobs    interface{} `json:"logprobs"`
	Stop        string      `json:"stop"`
	User        string      `json:"user"`
}

func NewWithDefault(model, prompt, user string) *Question {
	return &Question{
		model,
		prompt,
		4000,
		1,
		1,
		1,
		false,
		nil,
		"###",
		user,
	}
}
