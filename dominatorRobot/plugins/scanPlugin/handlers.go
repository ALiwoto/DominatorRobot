package scanPlugin

import (
	"strings"
	"time"

	"github.com/ALiwoto/argparser/argparser"
	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/utils"
	wv "github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func scanHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	sender := ctx.EffectiveUser.Id
	if ctx.Message.ReplyToMessage == nil || ctx.Message.ReplyToMessage.From == nil {
		return ext.EndGroups
	}

	u := utils.ResolveUser(sender)
	if !utils.CanScan(u) {
		return ext.EndGroups
	}

	src := ctx.Message.GetLink()

	replied := msg.ReplyToMessage
	target := replied.From.Id

	args, err := argparser.ParseArgDefault(msg.Text)
	if err != nil {
		return ext.EndGroups
	}

	force := args.GetAsBool("f", "force", "force-ban")
	reason := args.GetAsStringOrRaw("r", "reason", "reason")
	original := args.HasFlag("o", "original", "origin")

	if reason == "" {
		reason = args.GetFirstNoneEmptyValue()
	}

	if original && replied.ForwardFrom != nil && replied.ForwardFrom.Id != 0 {
		target = replied.ForwardFrom.Id
	}

	text := strings.ReplaceAll(replied.Text, "\r", "")
	text = strings.ReplaceAll(text, "\n", "\r\n")
	text = strings.TrimSpace(text)

	if force {
		_, err = wv.SibylClient.Ban(target, reason, text, src, replied.From.IsBot)
	} else {
		_, err = wv.SibylClient.Report(target, reason, text, src, replied.From.IsBot)
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

	md = mdparser.GetMono("Sibyl request has been sent!")

	_, _ = topMsg.EditText(b, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
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

	md = mdparser.GetMono("Sibyl request has been sent!")

	_, _ = topMsg.EditText(b, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})

	return ext.EndGroups
}
