package scanPlugin

import "errors"

var (
	scanManager = _getScanManager()
)

var (
	ErrRequestTooOld = errors.New("request is too old")
)
