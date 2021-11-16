package plugins

import (
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

func loadLimiter(d *ext.Dispatcher) {
	wv.RateLimiter = ratelimiter.NewLimiter(d, false, false)
	wv.RateLimiter.TextOnly = true
	wv.RateLimiter.ConsiderUser = true
	wv.RateLimiter.Start()
}
