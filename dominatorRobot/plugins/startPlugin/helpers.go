package startPlugin

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func LoadAllHandlers(d *ext.Dispatcher, t []rune) {
	dCmd := handlers.NewCommand(DCmd, dHandler)
	dominatorCmd := handlers.NewCommand(DominatorCmd, dHandler)
	dCmd.Triggers = t
	dominatorCmd.Triggers = t
	d.AddHandler(dCmd)
	d.AddHandler(dominatorCmd)
}
