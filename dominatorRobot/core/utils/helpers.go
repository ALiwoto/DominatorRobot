package utils

import (
	sibyl "github.com/ALiwoto/sibylSystemGo/sibylSystem"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoValues"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/database"
)

func ResolveUser(id int64) *sibyl.TokenInfo {
	t, err := database.GetTokenFromId(id)
	if err != nil || t == nil {
		return nil
	}

	return GetTokenFromServer(id)
}

func GetTokenFromServer(id int64) *sibyl.TokenInfo {
	t, err := wotoValues.SibylClient.GetToken(id)
	if err != nil || t == nil {
		return nil
	}

	return t
}

func CanScan(t *sibyl.TokenInfo) bool {
	return t != nil && t.Permission > 0x0
}

func CanForceScan(t *sibyl.TokenInfo) bool {
	return t != nil && t.Permission > 0x1
}
