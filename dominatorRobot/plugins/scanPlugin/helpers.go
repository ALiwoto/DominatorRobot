package scanPlugin

import (
	"time"

	ws "github.com/ALiwoto/StrongStringGo/strongStringGo"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func _getScansMap() *ws.SafeEMap[string, pendingScanData] {
	m := ws.NewSafeEMap[string, pendingScanData]()

	m.SetInterval(15 * time.Minute)
	m.SetExpiration(10 * time.Minute)
	m.EnableChecking()

	return m
}

func LoadAllHandlers(d *ext.Dispatcher, t []rune) {
	scanCmd := handlers.NewCommand(ScanCmd, scanHandler)
	revertCmd := handlers.NewCommand(RevertCmd, revertHandler)
	cancelScanCb := handlers.NewCallback(cancelScanCallBackQuery, cancelScanResponse)
	finalScanCb := handlers.NewCallback(finalScanCallBackQuery, finalScanResponse)

	scanCmd.Triggers = t
	revertCmd.Triggers = t

	d.AddHandler(cancelScanCb)
	d.AddHandler(finalScanCb)
	d.AddHandler(scanCmd)
	d.AddHandler(revertCmd)
}
