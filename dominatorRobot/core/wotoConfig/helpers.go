package wotoConfig

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	sibyl "github.com/ALiwoto/sibylSystemGo/sibylSystem"
	"github.com/AnimeKaizoku/ssg/ssg"
	"github.com/PaulSonOfLars/gotgbot/v2"
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

	if ConfigSettings.AddedToChatMessage != "" {
		ConfigSettings.botWelcomeMessage = GetMessageByLink(ConfigSettings.AddedToChatMessage)
	}

	return ConfigSettings, nil
}

func GetMessageByLink(link string) *gotgbot.Message {
	link = strings.ReplaceAll(link, "telegram.me", "t.me")
	link = strings.ReplaceAll(link, "telegram.dog", "t.me")
	link = strings.ReplaceAll(link, "https://", "")
	link = strings.ReplaceAll(link, "http://", "")
	if !strings.Contains(link, "t.me") {
		return nil
	}

	var chatId int64 = 0
	var messageId int64 = 0

	if strings.Contains(link, "/c/") {
		myStrs := strings.Split(link, "/c/")
		if len(myStrs) < 2 {
			return nil
		}
		myStrs = strings.Split(myStrs[1], "/")
		if len(myStrs) < 2 {
			return nil
		}
		chatId, _ = strconv.ParseInt("-100"+myStrs[0], 10, 64)
		messageId, _ = strconv.ParseInt(myStrs[1], 10, 64)
	} else {
		myStrs := strings.Split(link, "/")
		if len(myStrs) < 3 {
			return nil
		}
		chatId, _ = strconv.ParseInt(myStrs[1], 10, 64)
		messageId, _ = strconv.ParseInt(myStrs[2], 10, 64)
	}

	if chatId == 0 {
		return nil
	}

	return &gotgbot.Message{
		MessageId: messageId,
		Chat: gotgbot.Chat{
			Id: chatId,
		},
	}
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

func GetBotAddedMessage() *gotgbot.Message {
	if ConfigSettings == nil {
		return nil
	}

	return ConfigSettings.botWelcomeMessage
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
