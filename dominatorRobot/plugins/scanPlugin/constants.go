package scanPlugin

const (
	ScanCmd       = "scan"
	RevertCmd     = "revert"
	FullRevertCmd = "fullRevert"
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
	fullRevertBtnData   = "fRe"
)

const (
	anonRequestScan anonRequestType = iota + 1
	anonRequestRevert
	anonRequestFullRevert
)

const (
	wrappedUserTypeForwarder = iota + 1
	wrappedUserTypeOriginalSender
)
