package scanPlugin

import (
	sibylSystemGo "github.com/ALiwoto/sibylSystemGo/sibylSystem"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type anonRequestType uint8

type pendingScanData struct {
	ctx          *ext.Context
	bot          *gotgbot.Bot
	UniqueId     string
	OwnerId      int64
	Target       int64
	Reason       string
	targetInfo   *sibylSystemGo.GetInfoResult
	banConfig    *sibylSystemGo.BanConfig
	reportConfig *sibylSystemGo.ReportConfig
}

type anonContainer struct {
	myMessage *gotgbot.Message
	ctx       *ext.Context
	bot       *gotgbot.Bot
	request   anonRequestType
}
