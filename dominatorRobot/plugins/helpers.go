package plugins

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ALiwoto/sibylSystemGo/sibylSystem"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/logging"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/utils"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoConfig"
	wv "github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func StartTelegramBot() error {
	token := wotoConfig.GetBotToken()
	if len(token) == 0 {
		return errors.New("bot token is empty")
	}

	b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		Client: http.Client{},
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: 3 * gotgbot.DefaultTimeout,
		},
	})
	if err != nil {
		return err
	}

	utmp := ext.NewUpdater(nil)
	updater := &utmp
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: wotoConfig.DropUpdates(),
	})
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("%s has started | ID: %d", b.Username, b.Id))

	wv.HelperBot = b
	wv.BotUpdater = updater
	wv.SibylClient = wotoConfig.GetSibylClient()

	if wv.SibylClient == nil {
		// just to make sure.
		return errors.New("sibyl client is nil")
	}

	t, err := wv.SibylClient.GetToken(utils.GetIdFromToken(wotoConfig.GetSibylToken()))
	if err != nil {
		return err
	}
	if !t.IsOwner() {
		return errors.New("sibyl token doesn't have owner permissions")
	}

	wv.SibylDispatcher = sibylSystem.GetNewDispatcher(wv.SibylClient)

	LoadAllSibylHandlers(wv.SibylDispatcher)
	LoadAllHandlers(updater.Dispatcher, wotoConfig.GetCmdPrefixes())

	wv.SibylDispatcher.Listen()
	updater.Idle()
	return nil
}
