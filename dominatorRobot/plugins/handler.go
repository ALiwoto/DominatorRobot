package plugins

import (
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/plugins/scanPlugin"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func LoadAllHandlers(d *ext.Dispatcher, triggers []rune) {
	scanPlugin.LoadAllHandlers(d, triggers)
}
