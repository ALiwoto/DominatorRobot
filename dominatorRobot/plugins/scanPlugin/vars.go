package scanPlugin

import (
	"errors"

	ws "github.com/ALiwoto/StrongStringGo/strongStringGo"
)

var (
	scansMap = ws.NewSafeEMap[string, pendingScanData]()
)

var (
	ErrRequestTooOld = errors.New("request is too old")
)
