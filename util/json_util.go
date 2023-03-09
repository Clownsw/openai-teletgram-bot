package util

import (
	"openai-teletgram-bot/config"

	"github.com/bytedance/sonic"
)

func JsonGetString(json string, path ...interface{}) (string, error) {
	result, err := sonic.Get(StringToByteSlice(json), path...)
	if err != nil {
		return config.EmptyString, err
	}
	return result.String()
}
