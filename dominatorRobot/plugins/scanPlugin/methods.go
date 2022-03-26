package scanPlugin

import (
	"time"

	ws "github.com/ALiwoto/StrongStringGo/strongStringGo"
	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/ALiwoto/sibylSystemGo/sibylSystem"
	wv "github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func (d *pendingScanData) TakeAction() error {
	if d == nil {
		return ErrRequestTooOld
	}

	var err error
	if d.banConfig != nil {
		_, err = wv.SibylClient.Ban(d.Target, d.Reason, d.banConfig)
	} else if d.reportConfig != nil {
		_, err = wv.SibylClient.Report(d.Target, d.Reason, d.reportConfig)
	}

	return err
}

func (d *pendingScanData) getOwnerStr() string {
	return ws.ToBase32(d.OwnerId)
}

func (d *pendingScanData) getStrOwnerId() string {
	return ws.ToBase10(d.OwnerId)
}

func (d *pendingScanData) getStampStr() string {
	return ws.ToBase32(time.Now().Unix())
}

func (d *pendingScanData) GeneratedUniqueId() string {
	if d.UniqueId != "" {
		return d.UniqueId
	}

	d.UniqueId = d.getStampStr() + "=" + d.getOwnerStr()
	return d.UniqueId
}

func (d *pendingScanData) GetButtons() *gotgbot.InlineKeyboardMarkup {
	markup := &gotgbot.InlineKeyboardMarkup{}
	markup.InlineKeyboard = make([][]gotgbot.InlineKeyboardButton, 2)
	markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], gotgbot.InlineKeyboardButton{
		Text:         "Scan",
		CallbackData: pendingData + sepChar + d.getStrOwnerId() + sepChar + d.UniqueId,
	})
	markup.InlineKeyboard[1] = append(markup.InlineKeyboard[1], gotgbot.InlineKeyboardButton{
		Text:         "cancel",
		CallbackData: cancelData + sepChar + d.getStrOwnerId() + sepChar + d.UniqueId,
	})

	return markup
}

func (d *pendingScanData) getOperatorMd() mdparser.WMarkDown {
	md := mdparser.GetEmpty()
	byUser, err := d.bot.GetChat(d.targetInfo.BannedBy)
	if err != nil {
		return md
	}

	byInfo, err := wv.SibylClient.GetGeneralInfo(d.targetInfo.BannedBy)
	if err != nil {
		return md
	}

	md.Bold("\n • Scanned by").Normal(": ")
	switch byInfo.Permission {
	case sibylSystem.Enforcer:
		md.Bold("enforcer ")
	case sibylSystem.Inspector:
		md.Bold("inspector ")
	}

	md.Mention(byUser.FirstName, byUser.Id)
	md.Normal("[").Mono(ws.ToBase10(byUser.Id)).Normal("]")
	return md
}

func (d *pendingScanData) ParseAsMd() mdparser.WMarkDown {
	md := mdparser.GetNormal("Target user is currently banned in Sibyl System ")
	md.Normal("with the following details:")
	user, err := d.bot.GetChat(d.Target)
	if err != nil {
		/* most likely impossible */
		user = &gotgbot.Chat{
			FirstName: "Unknown",
			Id:        d.Target,
		}
	}

	md.Bold("\n • Target").Normal(": ")
	md.Mention(user.FirstName, user.Id)
	md.Normal("[").Mono(ws.ToBase10(user.Id)).Normal("]")
	md.Bold("\n • Type").Normal(": ")
	md.Mono("User").AppendThis(d.getOperatorMd())
	md.Bold("\n • Crime Coefficient").Normal(": ")
	md.Mono(d.targetInfo.GetStringCrimeCoefficient())
	md.Bold("\n • Reason(s)").Normal(": ")
	md.AppendThis(d.targetInfo.FormatFlags())
	md.Bold("\n • Description").Normal(": ")
	md.Normal(d.targetInfo.Reason)
	md.Bold("\n • Scan Date").Normal(": ")
	md.Mono(d.targetInfo.GetDateAsShort())
	if d.targetInfo.BanSourceUrl != "" {
		md.Bold("\n • Scan source").Normal(": ")
		md.Normal(d.targetInfo.BanSourceUrl)
	}

	md.Normal("\n\n Are you sure you want to proceed with scanning?")

	return md
}

//---------------------------------------------------------

func (a *anonContainer) DeleteMessage() {
	if a != nil && a.myMessage != nil {
		_, _ = a.myMessage.Delete(a.bot)
	}
}

func (a *anonContainer) ParseAsMd() mdparser.WMarkDown {
	md := mdparser.GetNormal("Seems like you are an anonymous user.\n")
	md.Normal("Please press the button below to confirm you are a valid user registered at PSB.")
	return md
}

func (a *anonContainer) GetButtons() *gotgbot.InlineKeyboardMarkup {
	markup := &gotgbot.InlineKeyboardMarkup{}
	markup.InlineKeyboard = make([][]gotgbot.InlineKeyboardButton, 2)
	markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], gotgbot.InlineKeyboardButton{
		Text:         "Press to confirm",
		CallbackData: anonConfirm + sepChar + a.getStrChatId(),
	})
	markup.InlineKeyboard[1] = append(markup.InlineKeyboard[1], gotgbot.InlineKeyboardButton{
		Text:         "cancel",
		CallbackData: anonCancelData + sepChar + a.getStrChatId(),
	})

	return markup
}

func (a *anonContainer) getStrChatId() string {
	return ws.ToBase10(a.ctx.EffectiveChat.Id)
}

//---------------------------------------------------------
