package database

import (
	"strconv"
	"strings"
	"sync"
	"time"

	sibyl "github.com/ALiwoto/sibylSystemGo/sibylSystem"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/logging"
	wConf "github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoConfig"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var SESSION *gorm.DB

func StartDatabase() error {
	// check if `SESSION` variable is already established or not.
	// if yes, check if we have got any error from it or not.
	// if there is an error in the session, it mean we have to establish
	// a new connection again.
	if SESSION != nil && SESSION.Error == nil {
		return nil
	}

	var db *gorm.DB
	var err error
	var conf *gorm.Config
	if wConf.IsDebug() {
		conf = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	} else {
		conf = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	}

	if wConf.ConfigSettings.UseSqlite {
		db, err = gorm.Open(sqlite.Open(wConf.ConfigSettings.DbName+".db"), conf)
	} else {
		db, err = gorm.Open(postgres.Open(wConf.ConfigSettings.DatabaseUrl), conf)
	}

	if err != nil {
		return err
	}

	SESSION = db
	logging.Info("Database connected")

	// Create tables if they don't exist
	err = SESSION.AutoMigrate(modelToken)
	if err != nil {
		return err
	}

	if wConf.ConfigSettings.UseSqlite {
		dbMutex = &sync.Mutex{}
	}

	tokenMapMutex = &sync.Mutex{}
	tokenDbMap = make(map[int64]*sibyl.TokenInfo)
	go cleanMaps()
	logging.Info("Auto-migrated database schema")

	return nil
}

func cleanMaps() {
	mtime := wConf.GetMaxCacheTime()
	for {
		time.Sleep(mtime)

		// please don't use len() function here, as it may return
		// `true` in some situations, but the maps may actually be
		// healthy, but they are only unused and so their caches are
		// completely deleted by cleaner.
		if tokenDbMap == nil {
			return
		}

		tokenMapMutex.Lock()
		for key, value := range tokenDbMap {
			if value == nil || value.IsExpired(mtime) {
				delete(tokenDbMap, key)
			}
		}
		tokenMapMutex.Unlock()
	}
}

func IsFirstTime() bool {
	return SESSION.Find(modelToken).RowsAffected == 0
}

func lockdb() {
	if wConf.ConfigSettings.UseSqlite {
		dbMutex.Lock()
	}
}

func unlockdb() {
	if wConf.ConfigSettings.UseSqlite {
		dbMutex.Unlock()
	}
}

func NewToken(t *sibyl.TokenInfo) {
	lockdb()
	tx := SESSION.Begin()
	tx.Save(t)
	tx.Commit()
	unlockdb()

	tokenMapMutex.Lock()
	tokenDbMap[t.UserId] = t
	tokenMapMutex.Unlock()
}

func GetTokenFromId(id int64) (*sibyl.TokenInfo, error) {
	if SESSION == nil {
		return nil, ErrNoSession
	}

	tokenMapMutex.Lock()
	t := tokenDbMap[id]
	tokenMapMutex.Unlock()
	if t != nil {
		t.SetCachedTime(time.Now())
		return t, nil
	}

	p := &sibyl.TokenInfo{}
	lockdb()
	SESSION.Where("user_id = ?", id).Take(p)
	unlockdb()
	if len(p.Hash) == 0 || p.UserId == 0 || p.UserId != id {
		// not found
		return nil, nil
	}
	p.SetCachedTime(time.Now())
	tokenMapMutex.Lock()
	tokenDbMap[p.UserId] = p
	tokenMapMutex.Unlock()

	return p, nil
}

func GetTokenFromString(token string) (*sibyl.TokenInfo, error) {
	id := GetIdFromToken(token)
	if id == 0 {
		return nil, ErrInvalidToken
	}

	u, err := GetTokenFromId(id)
	if err != nil {
		return nil, err
	}

	if u == nil || u.Hash != token {
		return nil, ErrInvalidToken
	}

	return u, nil
}

func GetIdFromToken(value string) int64 {
	if !strings.Contains(value, ":") {
		return 0
	}

	id, _ := strconv.ParseInt(strings.Split(value, ":")[0], 10, 64)
	return id
}

func GetIdAndHashFromToken(value string) (int64, string) {
	if !strings.Contains(value, ":") {
		return 0, ""
	}

	strs := strings.Split(value, ":")
	id, _ := strconv.ParseInt(strs[0], 10, 64)
	return id, strs[1]
}
