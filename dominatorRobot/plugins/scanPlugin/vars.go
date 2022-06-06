package scanPlugin

import (
	"errors"

	"github.com/AnimeKaizoku/ssg/ssg"
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
	// scanDataMap contains all sent scans by using their scan unique id as the key.
	// values of this map has to be removed from the map once we receive scan_approved or
	// scan_rejected event from API.
	scanDataMap = ssg.NewSafeMap[string, ScanDataContainer]()
)

var (
	ErrRequestTooOld = errors.New("request is too old")
)
