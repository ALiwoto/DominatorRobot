package scanPlugin

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func LoadAllHandlers(d *ext.Dispatcher, t []rune) {
	scanCmd := handlers.NewCommand(ScanCmd, scanHandler)
	revertCmd := handlers.NewCommand(RevertCmd, revertHandler)
	scanCmd.Triggers = t
	revertCmd.Triggers = t
	d.AddHandler(scanCmd)
	d.AddHandler(revertCmd)
}
