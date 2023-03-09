package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type AuthInfo struct {
	Token   string `yaml:"token"`
	BaseUrl string `yaml:"base_url"`
}

func NewAuthInfo(filePath string) *AuthInfo {
	authInfo := new(AuthInfo)

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(fileContent, authInfo)
	if err != nil {
		panic(err)
	}

	return authInfo
}
