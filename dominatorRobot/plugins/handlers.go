package plugins

import (
	"github.com/ALiwoto/sibylSystemGo/sibylSystem"
	wv "github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoValues"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/plugins/scanPlugin"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/plugins/startPlugin"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/gotgbot/ratelimiter/ratelimiter"
)

func LoadAllHandlers(d *ext.Dispatcher, triggers []rune) {
	loadLimiter(d)
	scanPlugin.LoadAllHandlers(d, triggers)
	startPlugin.LoadAllHandlers(d, triggers)
}

func LoadAllSibylHandlers(d *sibylSystem.SibylDispatcher) {
	scanPlugin.LoadAllSibylHandlers(d)
}

func loadLimiter(d *ext.Dispatcher) {
	wv.RateLimiter = ratelimiter.NewLimiter(d, nil)
	wv.RateLimiter.TextOnly = true
	wv.RateLimiter.Start()
}
