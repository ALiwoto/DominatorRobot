package wotoConfig

type BotConfig struct {
	BotToken     string
	DropUpdates  bool
	DatabaseUrl  string
	IsDebug      bool
	UseSqlite    bool
	DbName       string
	SibylToken   string
	MaxCacheTime int64
}
