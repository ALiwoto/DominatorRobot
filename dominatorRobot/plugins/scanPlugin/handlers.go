package scanPlugin

import (
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

	md := mdparser.GetNormal("Sending a cymatic scan request to Sibyl...")
	topMsg, err := ctx.EffectiveMessage.Reply(b, md.ToString(), &gotgbot.SendMessageOpts{
		AllowSendingWithoutReply: false,
		ParseMode:                wv.MarkdownV2,
	})
	if err != nil {
		return ext.EndGroups
	}

	time.Sleep(time.Millisecond * 600)

	if force {
		_, err = wv.SibylClient.Ban(target, reason, replied.Text, "this is source", replied.From.IsBot)
	} else {
		_, err = wv.SibylClient.Report(target, reason, replied.Text, "the source", replied.From.IsBot)
	}

	if err != nil {
		md = mdparser.GetMono(err.Error())
		_, _ = ctx.EffectiveMessage.Reply(b, md.ToString(), &gotgbot.SendMessageOpts{
			AllowSendingWithoutReply: false,
			ParseMode:                wv.MarkdownV2,
		})
		return ext.EndGroups
	}

	md = mdparser.GetMono("Sibyl request has been sent!")

	_, _ = topMsg.EditText(b, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})

	return ext.EndGroups
}
