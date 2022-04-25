package scanPlugin

import (
	"errors"
)

var (
	// scansMap contains the pending scan data by using their unique id as the key.
	scansMap = _getScansMap()
	// anonsMap contains the pending issued command by an anon admin. the key used is
	// the group id.
	anonsMap = _getAnonsMap()
	// inspectorsMap contains the
	inspectorsMap = _getInspectorsMap()
)

var (
	ErrRequestTooOld = errors.New("request is too old")
)
