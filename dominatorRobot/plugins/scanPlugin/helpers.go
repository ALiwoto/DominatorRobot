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

func _getAnonsMap() *ws.SafeEMap[int64, anonContainer] {
	m := ws.NewSafeEMap[int64, anonContainer]()

	m.SetInterval(10 * time.Minute)
	m.SetExpiration(5 * time.Minute)
	m.SetOnExpired(func(key int64, value anonContainer) {
		value.DeleteMessage()
	})
	m.EnableChecking()

	return m
}

func LoadAllHandlers(d *ext.Dispatcher, t []rune) {
	scanCmd := handlers.NewCommand(ScanCmd, scanHandler)
	revertCmd := handlers.NewCommand(RevertCmd, revertHandler)
	cancelScanCb := handlers.NewCallback(cancelScanCallBackQuery, cancelScanResponse)
	finalScanCb := handlers.NewCallback(finalScanCallBackQuery, finalScanResponse)
	cancelAnonCb := handlers.NewCallback(cancelAnonCallBackQuery, cancelAnonResponse)
	//confirmAnonCb := handlers.NewCallback(confirmAnonCallBackQuery, confirmAnonResponse)

	scanCmd.Triggers = t
	revertCmd.Triggers = t

	d.AddHandler(cancelAnonCb)
	d.AddHandler(cancelScanCb)
	d.AddHandler(finalScanCb)
	d.AddHandler(scanCmd)
	d.AddHandler(revertCmd)
}
