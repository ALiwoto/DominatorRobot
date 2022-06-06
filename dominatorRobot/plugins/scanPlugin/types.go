package scanPlugin

import (
	sibylSystemGo "github.com/ALiwoto/sibylSystemGo/sibylSystem"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type anonRequestType uint8
type wrappedUserType uint8
type coreHandler func(b *gotgbot.Bot, ctx *ext.Context, forceScan, noRedirect bool, targetId int64) error

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

type ScanDataContainer struct {
	ctx      *ext.Context
	bot      *gotgbot.Bot
	UniqueId string
	OwnerId  int64
}

type anonContainer struct {
	myMessage *gotgbot.Message
	ctx       *ext.Context
	bot       *gotgbot.Bot
	request   anonRequestType
}

type inspectorContainer struct {
	myMessage     *gotgbot.Message
	ctx           *ext.Context
	bot           *gotgbot.Bot
	targetUser    int64
	originHandler coreHandler
}

type TargetUserWrapper struct {
	UserType wrappedUserType
	User     *gotgbot.User
}

type multipleTargetContainer struct {
	myMessage     *gotgbot.Message
	ctx           *ext.Context
	bot           *gotgbot.Bot
	targetUsers   []*TargetUserWrapper
	originHandler coreHandler
}
