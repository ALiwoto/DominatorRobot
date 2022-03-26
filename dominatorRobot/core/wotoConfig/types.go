package wotoConfig

type BotConfig struct {
	BotToken    string `section:"general" key:"bot_token"`
	SibylToken  string `section:"general" key:"sibyl_token"`
	DropUpdates bool   `section:"general" key:"drop_updates"`
	SibylUrl    string `section:"general" key:"sibyl_url"`
	IsDebug     bool   `section:"general" key:"debug"`
	SupportAnon bool   `section:"general" key:"support_anon" default:"true"`
}
