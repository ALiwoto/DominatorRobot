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
	multipleTargetData  = "mulTi"
)

const (
	anonRequestScan anonRequestType = iota + 1
	anonRequestRevert
)

const (
	wrappedUserTypeForwarder = iota + 1
	wrappedUserTypeOriginalSender
)
