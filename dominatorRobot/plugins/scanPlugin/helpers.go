package scanPlugin

import (
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func _getScanManager() *pendingScanManager {
	return &pendingScanManager{
		pendingMutex: &sync.Mutex{},
		pendingMap:   make(map[string]*pendingScanData),
	}
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
