package startPlugin

import (
	"strconv"
	"time"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/utils"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoConfig"
	wv "github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func dHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if registering {
		return ext.ContinueGroups
	} else {
		registering = true
		defer func() {
			registering = false
		}()
	}
	user := ctx.EffectiveUser
	sender := user.Id
	message := ctx.EffectiveMessage
	u := utils.ResolveUser(sender)
	if !utils.CanScan(u) {
		return ext.EndGroups
	}

	wv.RateLimiter.AddCustomIgnore(sender, time.Minute*5, false)
	defer wv.RateLimiter.RemoveCustomIgnore(sender)

	//"TS No: `{userId}-48`\nARDR: `005-001`\n\nInitializing\n ▯ ▯ ▯ ▯"
	md := mdparser.GetNormal("Dominator Portable Psychological Diagnosis")
	md.AppendNormalThis(" and Supression System has been activated!")
	topMsg, err := message.Reply(b, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode: wv.MarkdownV2,
	})
	if err != nil || topMsg == nil {
		return ext.EndGroups
	}

	time.Sleep(sleepTime)

	//"TS No: `{userId}-48`\nARDR: `005-001`\n\nInitializing\n ▯ ▯ ▯ ▯"
	md = mdparser.GetNormal("TS No: ").AppendMonoThis(strconv.FormatInt(sender, 10))
	md.AppendNormalThis("\nARDR: ").AppendMonoThis("005-001")
	prototype := md.El()
	md.AppendNormalThis("\n\nInitializing\n ")

	var mdBack mdparser.WMarkDown

	for _, frame := range initializeAnim {
		mdBack = md.AppendNormal(frame)

		topMsg, _, err = topMsg.EditText(b, mdBack.ToString(), &gotgbot.EditMessageTextOpts{
			ParseMode: wv.MarkdownV2,
		})
		if err != nil || topMsg == nil {
			return ext.EndGroups
		}

		time.Sleep(sleepTime)
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

	topMsg, _, err = topMsg.EditText(b, mdBack.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})
	if err != nil || topMsg == nil {
		return ext.EndGroups
	}

	time.Sleep(sleepTime)

	md = md.AppendNormal("\n\nAffiliation: Public Safety Bureau, ")
	md.AppendBoldThis("Criminal Investigation Department")

	topMsg, _, err = topMsg.EditText(b, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})
	if err != nil || topMsg == nil {
		return ext.EndGroups
	}

	time.Sleep(sleepTime)

	md = mdBack.AppendNormalThis("\n\nDominator usage approval Confirmed.")
	md.AppendNormalThis(" You are a valid user!")

	_, _, _ = topMsg.EditText(b, md.ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: wv.MarkdownV2,
	})
	if err != nil || topMsg == nil {
		return ext.EndGroups
	}

	return ext.EndGroups
}

func chatMemberFilter(u *gotgbot.ChatMemberUpdated) bool {
	return u.NewChatMember.GetUser().Id == wv.HelperBot.Id
}

func chatMemberResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	chatMember := ctx.MyChatMember.NewChatMember
	chat := ctx.EffectiveChat
	if chatMember == nil {
		return nil
	}

	status := chatMember.GetStatus()
	// if status isn't equal to "member", it means there has been some other
	// operations being done on the group (such as a new admin being promoted or
	// something, idk).
	if status != "member" {
		return ext.EndGroups
	}

	botAddedMessage := wotoConfig.GetBotAddedMessage()
	if botAddedMessage != nil {
		_, _ = botAddedMessage.Copy(b, chat.Id, nil)
	}

	return nil
}
