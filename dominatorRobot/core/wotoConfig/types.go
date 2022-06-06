package wotoConfig

import "github.com/PaulSonOfLars/gotgbot/v2"

type BotConfig struct {
	BotToken           string `section:"general" key:"bot_token"`
	SibylToken         string `section:"general" key:"sibyl_token"`
	DropUpdates        bool   `section:"general" key:"drop_updates"`
	SibylUrl           string `section:"general" key:"sibyl_url"`
	AddedToChatMessage string `section:"general" key:"added_to_chat_message"`
	IsDebug            bool   `section:"general" key:"debug"`
	SupportAnon        bool   `section:"general" key:"support_anon" default:"true"`

	botWelcomeMessage *gotgbot.Message
}
