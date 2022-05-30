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
	// inspectorsMap contains the pending issued command by an inspector. the key used is
	// the user id.
	inspectorsMap = _getInspectorsMap()
	// multipleTargetsMap contains the pending issued command by an enforcer/inspector that
	// has pointed to multiple targets. the key used is the user id.
	multipleTargetsMap = _getMultipleTargetsMap()
)

var (
	ErrRequestTooOld = errors.New("request is too old")
)
