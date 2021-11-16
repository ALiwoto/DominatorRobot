package utils

import (
	sibyl "github.com/ALiwoto/sibylSystemGo/sibylSystem"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoValues"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/database"
)

func ResolveUser(id int64) *sibyl.TokenInfo {
	return GetTokenFromServer(id, false)
}

func GetTokenFromServer(id int64, cache bool) *sibyl.TokenInfo {
	t, err := wotoValues.SibylClient.GetToken(id)
	if err != nil || t == nil {
		return nil
	}

	if cache {
		database.NewToken(t)
	}
	return t
}

func CanScan(t *sibyl.TokenInfo) bool {
	return t != nil && t.Permission > 0x0
}

func CanForceScan(t *sibyl.TokenInfo) bool {
	return t != nil && t.Permission > 0x1
}
