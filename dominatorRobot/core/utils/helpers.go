package utils

import (
	"log"
	"strconv"
	"strings"

	"github.com/ALiwoto/mdparser/mdparser"
	sibyl "github.com/ALiwoto/sibylSystemGo/sibylSystem"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/logging"
	wv "github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func ResolveUser(id int64) *sibyl.TokenInfo {
	return GetTokenFromServer(id)
}

func GetIdFromToken(token string) int64 {
	if !strings.Contains(token, ":") {
		return 0
	}

	i, err := strconv.ParseInt(strings.Split(token, ":")[0], 10, 64)
	if err != nil {
		return 0
	}

	return i
}

func GetTokenFromServer(id int64) *sibyl.TokenInfo {
	t, err := wv.SibylClient.GetToken(id)
	if err != nil || t == nil {
		return nil
	}

	return t
}

func SendAlert(b *gotgbot.Bot, m *gotgbot.Message, md mdparser.WMarkDown) error {
	str := md.ToString()
	str = strings.ReplaceAll(str, b.GetToken(), "")
	_, err := m.Reply(b, str, &gotgbot.SendMessageOpts{ParseMode: MarkDownV2})
	if err != nil {
		log.Println(err)
	}

	return nil
}

func SendAlertErr(b *gotgbot.Bot, m *gotgbot.Message, e error) error {
	if e == nil {
		return nil
	}
	md := mdparser.GetItalic("Sibyl System says: \n")
	md = md.AppendNormal(e.Error())

	return SendAlert(b, m, md)
}

func SafeReply(b *gotgbot.Bot, ctx *ext.Context, output string) error {
	msg := ctx.EffectiveMessage
	if len(output) < 4096 {
		_, err := msg.Reply(b, output,
			&gotgbot.SendMessageOpts{ParseMode: MarkDownV2})
		if err != nil {
			logging.Error("got an error when trying to send results: ", err)
			return err
		}
	} else {
		_, err := b.SendDocument(ctx.EffectiveChat.Id, []byte(output), &gotgbot.SendDocumentOpts{
			ReplyToMessageId: msg.MessageId,
		})
		if err != nil {
			logging.Error("got an error when trying to send document: ", err)
			return err
		}
	}

	return nil
}

func CanScan(t *sibyl.TokenInfo) bool {
	return t != nil && t.Permission > 0x0
}

func CanForceScan(t *sibyl.TokenInfo) bool {
	return t != nil && t.Permission > 0x1
}

func CanRevert(t *sibyl.TokenInfo) bool {
	return t != nil && t.Permission > 0x1
}
