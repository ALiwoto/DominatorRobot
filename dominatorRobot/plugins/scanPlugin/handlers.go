package scanPlugin

import (
	"strconv"
	"strings"
	"time"

	"github.com/ALiwoto/argparser/argparser"
	"github.com/ALiwoto/mdparser/mdparser"
	sibyl "github.com/ALiwoto/sibylSystemGo/sibylSystem"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/utils"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoConfig"
	wv "github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoValues"
	ws "github.com/AnimeKaizoku/ssg/ssg"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// scanHandler is just a wrapper around coreScanHandler function.
func scanHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	return coreScanHandler(b, ctx, false, false, 0)
}

func coreScanHandler(b *gotgbot.Bot, ctx *ext.Context, forceScan, noRedirect bool, targetId int64) error {
	msg := ctx.EffectiveMessage
	sender := ctx.EffectiveUser.Id
	if msg.ReplyToMessage == nil || msg.ReplyToMessage.From == nil {
		return ext.EndGroups
	}

	// check for anon admin
	if sender == 1087968824 {
		if !wotoConfig.SupportAnon() || noRedirect {
			// just return from the handler if supporting anon admin is disabled
			// or we are disallowed to redirect from this function.
			return ext.EndGroups
		}

		// delete the old message (this method is nil-safe)
		anonsMap.Get(msg.Chat.Id).DeleteMessage()

		return sendAnonMessageHandler(b, &anonContainer{
			bot:     b,
			ctx:     ctx,
			request: anonRequestScan,
		})
	}

	u := utils.ResolveUser(sender)
	if !utils.CanScan(u) {
		return ext.EndGroups
	}

	src := msg.GetLink()

	replied := msg.ReplyToMessage
	target := targetId
	if target == 0 {
		target = replied.From.Id
	}
	hasMultipleTarget := false
	var targetType sibyl.EntityType

	args, err := argparser.ParseArgDefault(msg.Text)
	if err != nil {
		return ext.EndGroups
	}

	force := forceScan || args.HasFlag("f", "force", "force-ban")
	reason := args.GetAsStringOrRaw("r", "reason", "reason")
	original := args.HasFlag("o", "original", "origin")
	noPanel := args.HasFlag("no-panel", "noPanel")

	if reason == "" {
		reason = args.GetFirstNoneEmptyValue()
	}

	if targetId == 0 {
		if replied.ForwardFrom != nil && replied.ForwardFrom.Id != 0 {
			if original {
				if replied.ForwardFrom.IsBot {
					targetType = sibyl.EntityTypeBot
				}
				target = replied.ForwardFrom.Id
			} else {
				hasMultipleTarget = true
			}

		} else if replied.From.IsBot {
			targetType = sibyl.EntityTypeBot
		}
	}

	if target == b.Id {
		// maybe add a message or a warning or a cool quote here? dunno
		return ext.EndGroups
	}

	if hasMultipleTarget && !noPanel {
		targetUsers := []*TargetUserWrapper{
			{
				UserType: wrappedUserTypeForwarder,
				User:     replied.From,
			},
			{
				UserType: wrappedUserTypeOriginalSender,
				User:     replied.ForwardFrom,
			},
		}
		container := &multipleTargetContainer{
			ctx:           ctx,
			bot:           b,
			originHandler: coreScanHandler,
			targetUsers:   targetUsers,
		}
		return sendMultipleTargetPanelHandler(b, container)
	}

	targetInfo, _ := wv.SibylClient.GetInfo(target)
	if targetInfo != nil && targetInfo.Banned {
		var banConfig *sibyl.BanConfig
		var reportConfig *sibyl.ReportConfig
		if force {
			banConfig = &sibyl.BanConfig{
				Message:    replied.Text,
				SrcUrl:     src,
				TargetType: targetType,
				TheToken:   u.Hash,
			}
		} else {
			reportConfig = &sibyl.ReportConfig{
				Message:    replied.Text,
				SrcUrl:     src,
				TargetType: targetType,
				TheToken:   u.Hash,
				PollingId:  wv.SibylDispatcher.PollingId,
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

	var scanUniqueId string

	if force {
		_, err = wv.SibylClient.Ban(target, reason, &sibyl.BanConfig{
			Message:    replied.Text,
			SrcUrl:     src,
			TargetType: targetType,
			TheToken:   u.Hash,
		})
	} else {
		if u.Permission.CanBan() && !noRedirect && !noPanel {
			container := &inspectorContainer{
				ctx:           ctx,
				bot:           b,
				targetUser:    target,
				originHandler: coreScanHandler,
			}
			return sendInspectorScanPanelHandler(b, container)
		} else {
			scanUniqueId, err = wv.SibylClient.Report(target, reason, &sibyl.ReportConfig{
				Message:    replied.Text,
				SrcUrl:     src,
				TargetType: targetType,
				TheToken:   u.Hash,
				PollingId:  wv.SibylDispatcher.PollingId,
			})
			if scanUniqueId != "" {
				scanDataMap.Add(scanUniqueId, &ScanDataContainer{
					ctx:      ctx,
					bot:      b,
					UniqueId: scanUniqueId,
					OwnerId:  sender,
				})
			}
		}
	}

	if err != nil {
		_ = utils.SendAlertErr(b, msg, err)
		return ext.EndGroups
	}

	if noRedirect {
		// noRedirect is passed as true, we shouldn't show any
		// animation message.
		return ext.EndGroups
	}

	md := mdparser.GetNormal("An on-demand Cymatic scan request sent!")
	topMsg, err := ctx.EffectiveMessage.Reply(b, md.ToString(), &gotgbot.SendMessageOpts{
		AllowSendingWithoutReply: false,
		ParseMode:                wv.MarkdownV2,
	})
	if err != nil {
		return ext.EndGroups
	}

	time.Sleep(time.Millisecond * 600)

	md = mdparser.GetMono("Scan request sent.")

	_, _, _ = topMsg.EditText(b, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})

	return ext.EndGroups
}

func sendMultipleTargetPanelHandler(b *gotgbot.Bot, container *multipleTargetContainer) error {
	multipleTargetsMap.Add(container.ctx.EffectiveSender.Id(), container)
	msg := container.ctx.EffectiveMessage
	container.myMessage, _ = msg.Reply(b, container.ParseAsMd().ToString(), &gotgbot.SendMessageOpts{
		ParseMode:                wv.MarkdownV2,
		DisableWebPagePreview:    true,
		AllowSendingWithoutReply: true,
		ReplyMarkup:              container.GetButtons(),
	})

	return ext.EndGroups
}

func sendInspectorScanPanelHandler(b *gotgbot.Bot, container *inspectorContainer) error {
	inspectorsMap.Add(container.ctx.EffectiveSender.Id(), container)
	msg := container.ctx.EffectiveMessage
	container.myMessage, _ = msg.Reply(b, container.ParseAsMd().ToString(), &gotgbot.SendMessageOpts{
		ParseMode:                wv.MarkdownV2,
		DisableWebPagePreview:    true,
		AllowSendingWithoutReply: true,
		ReplyMarkup:              container.GetButtons(),
	})

	return ext.EndGroups
}

func sendAnonMessageHandler(b *gotgbot.Bot, container *anonContainer) error {
	anonsMap.Add(container.ctx.EffectiveChat.Id, container)
	msg := container.ctx.EffectiveMessage
	container.myMessage, _ = msg.Reply(b, container.ParseAsMd().ToString(), &gotgbot.SendMessageOpts{
		ParseMode:                wv.MarkdownV2,
		DisableWebPagePreview:    true,
		AllowSendingWithoutReply: true,
		ReplyMarkup:              container.GetButtons(),
	})

	return ext.EndGroups
}

func showAlreadyBannedHandler(b *gotgbot.Bot, data *pendingScanData) error {
	data.GeneratedUniqueId()
	scansMap.Add(data.UniqueId, data)
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

	// check for anon admin
	if sender == 1087968824 {
		if !wotoConfig.SupportAnon() {
			return ext.EndGroups
		}

		// delete the old message (this method is nil-safe)
		anonsMap.Get(msg.Chat.Id).DeleteMessage()

		return sendAnonMessageHandler(b, &anonContainer{
			bot:     b,
			ctx:     ctx,
			request: anonRequestRevert,
		})
	}

	requesterToken := utils.ResolveUser(sender)
	if !utils.CanScan(requesterToken) {
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

	_, err = wv.SibylClient.RemoveBan(target, "", &sibyl.RevertConfig{
		TheToken: requesterToken.Hash,
		SrcUrl:   msg.GetLink(),
	})
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

	md = mdparser.GetMono("Scan request has been sent.")

	_, _, _ = topMsg.EditText(b, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})

	return ext.EndGroups
}

func fullRevertHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	sender := ctx.EffectiveUser.Id

	// check for anon admin
	if sender == 1087968824 {
		if !wotoConfig.SupportAnon() {
			return ext.EndGroups
		}

		// delete the old message (this method is nil-safe)
		anonsMap.Get(msg.Chat.Id).DeleteMessage()

		return sendAnonMessageHandler(b, &anonContainer{
			bot:     b,
			ctx:     ctx,
			request: anonRequestRevert,
		})
	}

	requesterToken := utils.ResolveUser(sender)
	if !utils.CanScan(requesterToken) {
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

	_, err = wv.SibylClient.RemoveBan(target, "", &sibyl.RevertConfig{
		TheToken: requesterToken.Hash,
		SrcUrl:   msg.GetLink(),
	})
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

	md = mdparser.GetMono("Scan request has been sent.")

	_, _, _ = topMsg.EditText(b, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})

	return ext.EndGroups
}

//---------------------------------------------------------

func cancelScanCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return strings.HasPrefix(cq.Data, cancelData+sepChar)
}

func cancelAnonCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return strings.HasPrefix(cq.Data, anonCancelData+sepChar)
}

func confirmAnonCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return strings.HasPrefix(cq.Data, anonConfirm+sepChar)
}

func inspectorsCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return strings.HasPrefix(cq.Data, inspectorActionData+sepChar)
}

func multiTargetCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return strings.HasPrefix(cq.Data, multipleTargetData+sepChar)
}

func cancelScanResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	allStrs := ws.Split(query.Data, sepChar)
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
			Text:      "Unauthorized user!",
			ShowAlert: true,
			CacheTime: 5500,
		})
		return ext.EndGroups
	}

	uniqueId := allStrs[2]
	scansMap.Delete(uniqueId)
	if msg == nil {
		return ext.EndGroups
	}

	md := mdparser.GetMono("Scan request cancelled by user.")
	_, _, _ = msg.EditText(bot, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})
	return ext.EndGroups
}

func cancelAnonResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	allStrs := ws.Split(query.Data, sepChar)
	// format is anonCancelData + sepChar + d.getStrChatId()
	if len(allStrs) < 2 {
		return ext.EndGroups
	}

	chatId, err := strconv.ParseInt(allStrs[1], 10, 64)
	if err != nil {
		return ext.EndGroups
	}

	u := utils.ResolveUser(query.From.Id)
	if !utils.CanScan(u) {
		return ext.EndGroups
	}

	anonsMap.Delete(chatId)
	_, _ = ctx.EffectiveMessage.Delete(bot, nil)

	return ext.EndGroups
}

func confirmAnonResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	// a simple data is "anConfirm_-1001632556172"
	allStrs := ws.Split(query.Data, sepChar)
	// format is anonConfirm + sepChar + d.getStrChatId()
	if len(allStrs) < 2 {
		return ext.EndGroups
	}

	chatId, err := strconv.ParseInt(allStrs[1], 10, 64)
	if err != nil {
		return ext.EndGroups
	}

	u := utils.ResolveUser(query.From.Id)
	if !utils.CanScan(u) {
		return ext.EndGroups
	}

	container := anonsMap.Get(chatId)
	anonsMap.Delete(chatId)
	if container == nil {
		_, _ = ctx.EffectiveMessage.Delete(bot, nil)
		return nil
	}

	container.FastDeleteMessage()

	// hacky way to reduce amount of code (or rather, reuse the previously written code)
	container.ctx.EffectiveSender = ctx.EffectiveSender
	container.ctx.EffectiveUser = ctx.EffectiveUser

	switch container.request {
	case anonRequestScan:
		return scanHandler(bot, container.ctx)
	case anonRequestRevert:
		return revertHandler(bot, container.ctx)
	case anonRequestFullRevert:
		return fullRevertHandler(bot, container.ctx)
	default:
		// hm? unknown request type, sounds like not implemented or something
		// like that
		// anyway, too lazy to log this
		return nil
	}
}

func inspectorsResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	// a simple data is "insAc_confirm_1341091260"
	allStrs := ws.Split(query.Data, sepChar)
	// format is inspectorActionData + sepChar + forceData + sepChar + i.getStrOwnerId()
	if len(allStrs) < 3 {
		return ext.EndGroups
	}

	ownerId, err := strconv.ParseInt(allStrs[2], 10, 64)
	if err != nil {
		return ext.EndGroups
	}

	u := utils.ResolveUser(query.From.Id)
	if !utils.CanScan(u) {
		// user has lost the ability to scan.
		_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "Unauthorized user!",
			ShowAlert: true,
			CacheTime: 5500,
		})
		_, _ = ctx.EffectiveMessage.Delete(bot, nil)
		return ext.EndGroups
	}

	if ownerId != query.From.Id {
		_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "Unauthorized user!",
			ShowAlert: true,
			CacheTime: 5500,
		})
		return ext.EndGroups
	}

	container := inspectorsMap.Get(ownerId)
	inspectorsMap.Delete(ownerId)
	if container == nil {
		_, _ = ctx.EffectiveMessage.Delete(bot, nil)
		return nil
	}

	data := allStrs[1]
	switch data {
	case forceData:
		md := mdparser.GetMono("Enforcement action forced!")
		_, _, _ = ctx.EffectiveMessage.EditText(bot, md.ToString(), &gotgbot.EditMessageTextOpts{
			ParseMode: wv.MarkdownV2,
		})
		return coreScanHandler(bot, container.ctx, true, true, container.targetUser)
	case confirmData:
		md := mdparser.GetMono("Cymatic scan requested for user.")
		_, _, _ = ctx.EffectiveMessage.EditText(bot, md.ToString(), &gotgbot.EditMessageTextOpts{
			ParseMode: wv.MarkdownV2,
		})
		return coreScanHandler(bot, container.ctx, false, true, container.targetUser)
	case cancelData:
		md := mdparser.GetMono("Scan request cancelled.")
		_, _, _ = ctx.EffectiveMessage.EditText(bot, md.ToString(), &gotgbot.EditMessageTextOpts{
			ParseMode: wv.MarkdownV2,
		})
		return nil
	default:
		// hm? unknown data, sounds like not implemented or something
		// like that
		// anyway, too lazy to log this
		return nil
	}
}

func multiTargetPanelResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	// a simple data is "mulTi_1341091260_1221024018"
	allStrs := ws.Split(query.Data, sepChar)
	// format is multipleTargetData + sepChar + m.getStrOwnerId() + sepChar + ssg.ToBase10(id)
	if len(allStrs) < 3 {
		return ext.EndGroups
	}

	ownerId, err := strconv.ParseInt(allStrs[1], 10, 64)
	if err != nil {
		return ext.EndGroups
	}

	u := utils.ResolveUser(query.From.Id)
	if !utils.CanScan(u) {
		// user has lost the ability to scan.
		_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "Unauthorized user!",
			ShowAlert: true,
			CacheTime: 5500,
		})
		_, _ = ctx.EffectiveMessage.Delete(bot, nil)
		return ext.EndGroups
	}

	if ownerId != query.From.Id {
		_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "Unauthorized user!",
			ShowAlert: true,
			CacheTime: 5500,
		})
		return ext.EndGroups
	}

	container := multipleTargetsMap.Get(ownerId)
	multipleTargetsMap.Delete(ownerId)
	if container == nil {
		_, _ = ctx.EffectiveMessage.Delete(bot, nil)
		return nil
	}

	container.FastDeleteMessage()

	// TODO: add support for "all", to multiple all of them at once.
	targetId := ws.ToInt64(allStrs[2])
	if targetId == 0 {
		_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "Invalid 0 data specified in callback button... please try another option.",
			ShowAlert: true,
		})
		return ext.EndGroups
	}

	return coreScanHandler(bot, container.ctx, false, false, targetId)
}

func finalScanCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return strings.HasPrefix(cq.Data, pendingData+sepChar)
}

func finalScanResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	allStrs := ws.Split(query.Data, sepChar)
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
			Text:      "Unauthorized user!",
			ShowAlert: true,
			CacheTime: 5500,
		})
		return ext.EndGroups
	}

	uniqueId := allStrs[2]
	err = scansMap.Get(uniqueId).TakeAction()
	scansMap.Delete(uniqueId)
	if msg == nil {
		return ext.EndGroups
	}

	if err != nil {
		md := mdparser.GetMono("Scan request failed: ").Mono(err.Error())
		_, _, _ = msg.EditText(bot, md.ToString(), &gotgbot.EditMessageTextOpts{
			ParseMode: wv.MarkdownV2,
		})
		return ext.EndGroups
	}

	md := mdparser.GetMono("Scan request has been sent.")
	_, _, _ = msg.EditText(bot, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})
	return ext.EndGroups
}

//---------------------------------------------------------

func sibylScanApprovedHandler(client sibyl.SibylClient, ctx *sibyl.SibylUpdateContext) error {
	approvedData := ctx.ScanRequestApproved
	data := scanDataMap.Get(approvedData.UniqueId)
	if data == nil {
		// this scan is not sent by us.
		return nil
	}
	scanDataMap.Delete(approvedData.UniqueId)

	msg := data.ctx.EffectiveMessage
	md := mdparser.GetNormal("Crime coefficient ").Bold("over 300, ")
	md.Normal("user is a target for enforcement action!\n")
	md.Bold("Enforcement mode: ").Normal("Lethal Eliminator")
	// md2 := mdparser.GetMono("Crime coefficient over $$$, user is a target for enforcement action!\nEnforcement mode: Lethal Eliminator")
	if approvedData.AgentReason != "" {
		md.Bold("\nApproved reason: ").Mono(approvedData.AgentReason)
	}
	_, _ = msg.Reply(data.bot, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode:             wv.MarkdownV2,
		DisableWebPagePreview: true,
	})
	return nil
}

func sibylScanRejectedHandler(client sibyl.SibylClient, ctx *sibyl.SibylUpdateContext) error {
	rejectedData := ctx.ScanRequestRejected
	data := scanDataMap.Get(rejectedData.UniqueId)
	if data == nil {
		// this scan is not sent by us.
		return nil
	}
	scanDataMap.Delete(rejectedData.UniqueId)

	msg := data.ctx.EffectiveMessage
	md := mdparser.GetNormal("Crime Coefficient is under 100.\n")
	md.Normal("Not a target for enforcement action.\n")
	md.Normal("Trigger of dominator will be locked!")
	// md := mdparser.GetMono("Crime Coefficient is under $$$.\nNot a target for enforcement action.\nTrigger of dominator will be locked!")
	if rejectedData.AgentReason != "" {
		md.Bold("\nReason: ").Mono(rejectedData.AgentReason)
	}
	_, _ = msg.Reply(data.bot, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode:             wv.MarkdownV2,
		DisableWebPagePreview: true,
	})
	return nil
}

//---------------------------------------------------------
