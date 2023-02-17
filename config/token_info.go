package config

type OpenAIInfo struct {
	Token   string
	BaseUrl string
}

func NewOpenAIInfo(token string, baseUrl string) *OpenAIInfo {
	return &OpenAIInfo{token, baseUrl}
}
