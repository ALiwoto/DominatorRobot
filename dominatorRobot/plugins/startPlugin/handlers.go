package startPlugin

import (
	"strconv"
	"time"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/utils"
	wv "github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func dHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	sender := user.Id
	message := ctx.EffectiveMessage
	u := utils.ResolveUser(sender)
	if !utils.CanScan(u) {
		return ext.EndGroups
	}

	//"TS No: `{userid}-48`\nARDR: `005-001`\n\nInitializing\n ▯ ▯ ▯ ▯"
	md := mdparser.GetNormal("Dominator Portable Psychological Diagnosis")
	md.AppendNormalThis(" and Supression System has been activated!")
	topMsg, err := message.Reply(b, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode: wv.MarkdownV2,
	})
	if err != nil || topMsg == nil {
		return ext.EndGroups
	}

	time.Sleep(time.Millisecond * 2500)

	//"TS No: `{userid}-48`\nARDR: `005-001`\n\nInitializing\n ▯ ▯ ▯ ▯"
	md = mdparser.GetNormal("TS No: ").AppendMonoThis(strconv.FormatInt(sender, 10))
	md.AppendNormalThis("\nARDR: ").AppendMonoThis("005-001")
	prototype := md.El()
	md.AppendNormalThis("\n\nInitializing\n ")

	var mdBack mdparser.WMarkDown

	for _, frame := range initializeAnim {
		mdBack = md.AppendNormal(frame)

		topMsg, err = topMsg.EditText(b, mdBack.ToString(), &gotgbot.EditMessageTextOpts{
			ParseMode: wv.MarkdownV2,
		})
		if err != nil || topMsg == nil {
			return ext.EndGroups
		}

		time.Sleep(time.Millisecond * 2500)
	}

	var userStatus string
	if u.Permission > 0x1 {
		userStatus = "Inspector"
	} else {
		userStatus = "Enforcer"
	}

	md = prototype.AppendBoldThis("\nCRIMINAL INVESTIGATION DEPARTMENT")
	md.AppendBoldThis("\n • Name: ").AppendNormalThis(user.FirstName)
	md.AppendBoldThis("\n • RD: ").AppendNormalThis("2021-1107-T175741")
	mdBack = md.AppendBoldThis("\n • Position: ").AppendNormalThis(userStatus)
	md = md.AppendNormal("\n\nAffiliation: Public Safety Bureau, ")
	md.AppendBoldThis("Criminal Investigation Department")

	topMsg, err = topMsg.EditText(b, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})
	if err != nil || topMsg == nil {
		return ext.EndGroups
	}

	time.Sleep(time.Millisecond * 2500)

	md = mdBack.AppendNormalThis("\n\nDominator usage approval Confirmed.")
	md.AppendNormalThis(" You are a valid user!")

	_, _ = topMsg.EditText(b, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})
	if err != nil || topMsg == nil {
		return ext.EndGroups
	}

	return ext.EndGroups
}
