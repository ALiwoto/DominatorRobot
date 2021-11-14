package database

import (
	"errors"
	"sync"

	sibyl "github.com/ALiwoto/sibylSystemGo/sibylSystem"
)

var (
	ErrInvalidToken   = errors.New("token is invalid")
	ErrNoSession      = errors.New("database session is not initialized")
	ErrTooManyRevokes = errors.New("token has been revoked too many times")
)

var (
	dbMutex       *sync.Mutex
	tokenMapMutex *sync.Mutex
	tokenDbMap    map[int64]*sibyl.TokenInfo
	modelToken    *sibyl.TokenInfo = &sibyl.TokenInfo{}
)
