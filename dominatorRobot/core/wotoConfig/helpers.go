package wotoConfig

import (
	"net/http"
	"time"

	sibyl "github.com/ALiwoto/sibylSystemGo/sibylSystem"
	"github.com/AnimeKaizoku/ssg/ssg"
)

func ParseConfig(filename string) (*BotConfig, error) {
	if ConfigSettings != nil {
		return ConfigSettings, nil
	}
	config := &BotConfig{}

	err := ssg.ParseConfig(config, filename)
	if err != nil {
		return nil, err
	}

	ConfigSettings = config

	return ConfigSettings, nil
}

func LoadConfig() (*BotConfig, error) {
	return ParseConfig("config.ini")
}

func IsDebug() bool {
	if ConfigSettings != nil {
		return ConfigSettings.IsDebug
	}
	return true
}

func GetBotToken() string {
	if ConfigSettings != nil {
		return ConfigSettings.BotToken
	}
	return ""
}

func SupportAnon() bool {
	if ConfigSettings != nil {
		return ConfigSettings.SupportAnon
	}
	return false
}

func DropUpdates() bool {
	if ConfigSettings != nil {
		return ConfigSettings.DropUpdates
	}
	return false
}

func GetCmdPrefixes() []rune {
	return []rune{'/', '!'}
}

func GetSibylToken() string {
	if ConfigSettings == nil {
		return ""
	}

	return ConfigSettings.SibylToken
}

func GetSibylClient() sibyl.SibylClient {
	if ConfigSettings == nil {
		return nil
	}

	return sibyl.NewClient(
		ConfigSettings.SibylToken,
		&sibyl.SibylConfig{
			HostUrl: ConfigSettings.SibylUrl,
			HttpClient: &http.Client{
				Timeout: time.Second * 35,
			},
		},
	)
}

func GetSibylConfig() *sibyl.SibylConfig {
	return &sibyl.SibylConfig{
		HostUrl: ConfigSettings.SibylUrl,
	}
}
