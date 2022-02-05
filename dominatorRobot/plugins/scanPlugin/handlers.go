package scanPlugin

import (
	"strconv"
	"strings"
	"time"

	"github.com/ALiwoto/StrongStringGo/strongStringGo"
	"github.com/ALiwoto/argparser/argparser"
	"github.com/ALiwoto/mdparser/mdparser"
	sibylSystemGo "github.com/ALiwoto/sibylSystemGo/sibylSystem"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/utils"
	wv "github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func scanHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	sender := ctx.EffectiveUser.Id
	if msg.ReplyToMessage == nil || msg.ReplyToMessage.From == nil {
		return ext.EndGroups
	}

	u := utils.ResolveUser(sender)
	if !utils.CanScan(u) {
		return ext.EndGroups
	}

	src := msg.GetLink()

	replied := msg.ReplyToMessage
	target := replied.From.Id
	var targetType sibylSystemGo.EntityType

	args, err := argparser.ParseArgDefault(msg.Text)
	if err != nil {
		return ext.EndGroups
	}

	force := args.HasFlag("f", "force", "force-ban")
	reason := args.GetAsStringOrRaw("r", "reason", "reason")
	original := args.HasFlag("o", "original", "origin")

	if reason == "" {
		reason = args.GetFirstNoneEmptyValue()
	}

	if original && replied.ForwardFrom != nil && replied.ForwardFrom.Id != 0 {
		if replied.ForwardFrom.IsBot {
			targetType = sibylSystemGo.EntityTypeBot
		}
		target = replied.ForwardFrom.Id
	} else if replied.From.IsBot {
		targetType = sibylSystemGo.EntityTypeBot
	}

	targetInfo, _ := wv.SibylClient.GetInfo(target)
	if targetInfo != nil && targetInfo.Banned {
		var banConfig *sibylSystemGo.BanConfig
		var reportConfig *sibylSystemGo.ReportConfig
		if force {
			banConfig = &sibylSystemGo.BanConfig{
				Message:    replied.Text,
				SrcUrl:     src,
				TargetType: targetType,
				TheToken:   u.Hash,
			}
		} else {
			reportConfig = &sibylSystemGo.ReportConfig{
				Message:    replied.Text,
				SrcUrl:     src,
				TargetType: targetType,
				TheToken:   u.Hash,
			}
		}

		pending := &pendingScanData{
			ctx:          ctx,
			targetInfo:   targetInfo,
			bot:          b,
			banConfig:    banConfig,
			reportConfig: reportConfig,
			OwnerId:      sender,
			Target:       target,
			Reason:       reason,
		}
		return showAlreadyBannedHandler(b, pending)
	}

	if force {
		_, err = wv.SibylClient.Ban(target, reason, &sibylSystemGo.BanConfig{
			Message:    replied.Text,
			SrcUrl:     src,
			TargetType: targetType,
			TheToken:   u.Hash,
		})
	} else {
		_, err = wv.SibylClient.Report(target, reason, &sibylSystemGo.ReportConfig{
			Message:    replied.Text,
			SrcUrl:     src,
			TargetType: targetType,
			TheToken:   u.Hash,
		})
	}

	if err != nil {
		_ = utils.SendAlertErr(b, msg, err)
		return ext.EndGroups
	}

	md := mdparser.GetNormal("Sending a cymatic scan request to Sibyl...")
	topMsg, err := ctx.EffectiveMessage.Reply(b, md.ToString(), &gotgbot.SendMessageOpts{
		AllowSendingWithoutReply: false,
		ParseMode:                wv.MarkdownV2,
	})
	if err != nil {
		return ext.EndGroups
	}

	time.Sleep(time.Millisecond * 600)

	if err != nil {
		_ = utils.SendAlertErr(b, msg, err)
		return ext.EndGroups
	}

	md = mdparser.GetMono("Sibyl request has been sent.")

	_, _, _ = topMsg.EditText(b, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})

	return ext.EndGroups
}

func showAlreadyBannedHandler(b *gotgbot.Bot, data *pendingScanData) error {
	data.GeneratedUniqueId()
	scanManager.AddData(data)
	msg := data.ctx.EffectiveMessage
	_, _ = msg.Reply(b, data.ParseAsMd().ToString(), &gotgbot.SendMessageOpts{
		ParseMode:                wv.MarkdownV2,
		DisableWebPagePreview:    true,
		AllowSendingWithoutReply: true,
		ReplyMarkup:              data.GetButtons(),
	})

	return ext.EndGroups
}

func revertHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	sender := ctx.EffectiveUser.Id

	u := utils.ResolveUser(sender)
	if !utils.CanScan(u) {
		return ext.EndGroups
	}

	replied := msg.ReplyToMessage
	var target int64

	args, err := argparser.ParseArgDefault(msg.Text)
	if err != nil {
		return ext.EndGroups
	}

	targetUser, ok := args.GetAsIntegerOrRaw("u", "id", "target", "user")
	original := args.HasFlag("o", "original", "origin")
	if ok {
		target = targetUser
	}

	if target == 0 && replied != nil && replied.From != nil && replied.From.Id != 0 {
		if original && replied.ForwardFrom != nil && replied.ForwardFrom.Id != 0 {
			target = replied.ForwardFrom.Id
		} else {
			target = replied.From.Id
		}
	}

	_, err = wv.SibylClient.RemoveBan(target)
	if err != nil {
		_ = utils.SendAlertErr(b, msg, err)
		return ext.EndGroups
	}

	md := mdparser.GetNormal("Sending a cymatic scan request to Sibyl...")
	topMsg, err := ctx.EffectiveMessage.Reply(b, md.ToString(), &gotgbot.SendMessageOpts{
		AllowSendingWithoutReply: false,
		ParseMode:                wv.MarkdownV2,
	})
	if err != nil {
		return ext.EndGroups
	}

	time.Sleep(time.Millisecond * 600)

	if err != nil {
		_ = utils.SendAlertErr(b, msg, err)
		return ext.EndGroups
	}

	md = mdparser.GetMono("Sibyl request has been sent.")

	_, _, _ = topMsg.EditText(b, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})

	return ext.EndGroups
}

//---------------------------------------------------------

func cancelScanCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return strings.HasPrefix(cq.Data, cancelData+sepChar)
}

func cancelScanResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	allStrs := strongStringGo.Split(query.Data, sepChar)
	msg := query.Message
	// format is cancelData + sepChar + d.getStrOwnerId() + sepChar + d.UniqueId
	if len(allStrs) < 3 {
		return ext.EndGroups
	}

	ownerId, err := strconv.ParseInt(allStrs[1], 10, 64)
	if err != nil {
		return ext.EndGroups
	}

	if ownerId != query.From.Id {
		_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "This button is not for you...",
			ShowAlert: true,
			CacheTime: 5500,
		})
		return ext.EndGroups
	}

	uniqueId := allStrs[2]
	scanManager.RemoveData(uniqueId)
	if msg == nil {
		return ext.EndGroups
	}

	md := mdparser.GetMono("Sibyl request cancelled by user.")
	_, _, _ = msg.EditText(bot, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})
	return ext.EndGroups
}

func finalScanCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return strings.HasPrefix(cq.Data, pendingData+sepChar)
}

func finalScanResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	allStrs := strongStringGo.Split(query.Data, sepChar)
	msg := query.Message
	// format is cancelData + sepChar + d.getStrOwnerId() + sepChar + d.UniqueId
	if len(allStrs) < 3 {
		return ext.EndGroups
	}

	ownerId, err := strconv.ParseInt(allStrs[1], 10, 64)
	if err != nil {
		return ext.EndGroups
	}

	if ownerId != query.From.Id {
		_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "This button is not for you...",
			ShowAlert: true,
			CacheTime: 5500,
		})
		return ext.EndGroups
	}

	uniqueId := allStrs[2]
	err = scanManager.GetScanData(uniqueId).TakeAction()
	scanManager.RemoveData(uniqueId)
	if msg == nil {
		return ext.EndGroups
	}

	if err != nil {
		md := mdparser.GetMono("Sibyl request failed: ").Mono(err.Error())
		_, _, _ = msg.EditText(bot, md.ToString(), &gotgbot.EditMessageTextOpts{
			ParseMode: wv.MarkdownV2,
		})
		return ext.EndGroups
	}

	md := mdparser.GetMono("Sibyl request cancelled by user.")
	_, _, _ = msg.EditText(bot, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})
	return ext.EndGroups
}

//---------------------------------------------------------
