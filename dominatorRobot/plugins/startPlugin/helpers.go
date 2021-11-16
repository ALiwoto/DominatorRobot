package startPlugin

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func LoadAllHandlers(d *ext.Dispatcher, t []rune) {
	dCmd := handlers.NewCommand(dCmd, dHandler)
	dCmd.Triggers = t
	d.AddHandler(dCmd)
}
