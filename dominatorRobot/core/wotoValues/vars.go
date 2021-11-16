package wotoValues

import (
	sibyl "github.com/ALiwoto/sibylSystemGo/sibylSystem"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/gotgbot/ratelimiter/ratelimiter"
)

var (
	HelperBot   *gotgbot.Bot
	BotUpdater  *ext.Updater
	SibylClient sibyl.SibylClient
	RateLimiter *ratelimiter.Limiter
)
