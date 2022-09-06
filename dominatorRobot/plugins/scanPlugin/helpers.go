package scanPlugin

import (
	"time"

	sibyl "github.com/ALiwoto/sibylSystemGo/sibylSystem"
	ws "github.com/AnimeKaizoku/ssg/ssg"
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

func _getInspectorsMap() *ws.SafeEMap[int64, inspectorContainer] {
	m := ws.NewSafeEMap[int64, inspectorContainer]()

	m.SetInterval(10 * time.Minute)
	m.SetExpiration(5 * time.Minute)
	m.SetOnExpired(func(key int64, value inspectorContainer) {
		if value.myMessage == nil {
			return
		}

		_, _ = value.myMessage.Delete(value.bot, nil)

		if value.originHandler != nil {
			_ = value.originHandler(value.bot, value.ctx, false, true, 0)
		}
	})
	m.EnableChecking()

	return m
}

func _getMultipleTargetsMap() *ws.SafeEMap[int64, multipleTargetContainer] {
	m := ws.NewSafeEMap[int64, multipleTargetContainer]()

	m.SetInterval(10 * time.Minute)
	m.SetExpiration(15 * time.Minute)
	m.SetOnExpired(func(key int64, value multipleTargetContainer) {
		if value.myMessage == nil {
			return
		}

		_, _ = value.myMessage.Delete(value.bot, nil)
	})
	m.EnableChecking()

	return m
}

func LoadAllHandlers(d *ext.Dispatcher, t []rune) {
	scanCmd := handlers.NewCommand(ScanCmd, scanHandler)
	revertCmd := handlers.NewCommand(RevertCmd, revertHandler)
	fullRevertCmd := handlers.NewCommand(FullRevertCmd, fullRevertHandler)
	cancelScanCb := handlers.NewCallback(cancelScanCallBackQuery, cancelScanResponse)
	fullRevertCb := handlers.NewCallback(fullRevertCallBackQuery, fullRevertBtnResponse)
	finalScanCb := handlers.NewCallback(finalScanCallBackQuery, finalScanResponse)
	cancelAnonCb := handlers.NewCallback(cancelAnonCallBackQuery, cancelAnonResponse)
	confirmAnonCb := handlers.NewCallback(confirmAnonCallBackQuery, confirmAnonResponse)
	inspectorsCb := handlers.NewCallback(inspectorsCallBackQuery, inspectorsResponse)
	multiTargetCb := handlers.NewCallback(multiTargetCallBackQuery, multiTargetPanelResponse)

	scanCmd.Triggers = t
	revertCmd.Triggers = t
	fullRevertCmd.Triggers = t

	d.AddHandler(cancelAnonCb)
	d.AddHandler(confirmAnonCb)
	d.AddHandler(inspectorsCb)
	d.AddHandler(multiTargetCb)
	d.AddHandler(cancelScanCb)
	d.AddHandler(fullRevertCb)
	d.AddHandler(finalScanCb)
	d.AddHandler(scanCmd)
	d.AddHandler(revertCmd)
	d.AddHandler(fullRevertCmd)
}

func LoadAllSibylHandlers(d *sibyl.SibylDispatcher) {
	d.AddHandler(sibyl.UpdateTypeScanRequestApproved, sibylScanApprovedHandler)
	d.AddHandler(sibyl.UpdateTypeScanRequestRejected, sibylScanRejectedHandler)
}
