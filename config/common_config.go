package config

import "errors"

const (
	EmptyString  = ""
	ErrorMessage = "这个问题我目前无法回答!"
)

var (
	OpenAiQueryError = errors.New("OpenAi Query Error")
)
