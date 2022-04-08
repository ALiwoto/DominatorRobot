package scanPlugin

const (
	ScanCmd   = "scan"
	RevertCmd = "revert"
)

const (
	sepChar        = "_"
	pendingData    = "pending"
	cancelData     = "cancel"
	anonCancelData = "anCanc"
	anonConfirm    = "anConfirm"
)

const (
	anonRequestScan anonRequestType = iota + 1
	anonRequestRevert
)
