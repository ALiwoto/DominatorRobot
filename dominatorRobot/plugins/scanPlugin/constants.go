package scanPlugin

const (
	ScanCmd   = "scan"
	RevertCmd = "revert"
)

const (
	sepChar             = "_"
	pendingData         = "pending"
	cancelData          = "cancel"
	anonCancelData      = "anCanc"
	anonConfirm         = "anConfirm"
	forceData           = "force"
	confirmData         = "confirm"
	inspectorActionData = "insAc"
)

const (
	anonRequestScan anonRequestType = iota + 1
	anonRequestRevert
)
