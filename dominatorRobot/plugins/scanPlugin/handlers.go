package scanPlugin

import (
	"github.com/ALiwoto/argparser/argparser"
	"github.com/ALiwoto/mdparser/mdparser"
	wv "github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func scanHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	args, err := argparser.ParseArgDefault(msg.Text)
	if err != nil {
		return ext.EndGroups
	}

	force := args.GetAsBool("f", "force", "force-ban")
	reason := args.GetAsStringOrRaw("r", "reason", "reason")

	if reason == "" {
		reason = args.GetFirstNoneEmptyValue()
	}

	md := mdparser.GetNormal("Sending a cymatic scan request to Sibyl..." + reason)
	_, err = ctx.EffectiveMessage.Reply(b, md.ToString(), &gotgbot.SendMessageOpts{
		AllowSendingWithoutReply: false,
		ParseMode:                wv.MarkdownV2,
	})
	if err != nil {
		return ext.EndGroups
	}

	if force {
		wv.SibylClient.Ban(0, reason, "message", "", false)
	}

	return ext.EndGroups
}
