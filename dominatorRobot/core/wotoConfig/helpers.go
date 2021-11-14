package wotoConfig

import (
	"errors"
	"time"

	sibyl "github.com/ALiwoto/sibylSystemGo/sibylSystem"
	"github.com/bigkevmcd/go-configparser"
)

func ParseConfig(configFile string) (*BotConfig, error) {
	if ConfigSettings != nil {
		return ConfigSettings, nil
	}
	ConfigSettings = &BotConfig{}

	cfgContent, err := configparser.NewConfigParserFromFile(configFile)
	if err != nil {
		return nil, err
	}

	ConfigSettings.BotToken, err = cfgContent.Get("general", "bot_token")
	if err != nil {
		return nil, err
	}

	// don't ignore if sibyl token doesn't exist in config file
	ConfigSettings.SibylToken, err = cfgContent.Get("general", "sibyl_token")
	if err != nil {
		return nil, err
	}

	ConfigSettings.SibylUrl, err = cfgContent.Get("general", "sibyl_url")
	if err != nil {
		return nil, err
	}

	if ConfigSettings.SibylToken == "" {
		return nil, errors.New("sibyl token is empty")
	}

	ConfigSettings.DropUpdates, _ = cfgContent.GetBool("general", "drop_updates")

	ConfigSettings.IsDebug, _ = cfgContent.GetBool("general", "is_debug")

	ConfigSettings.DatabaseUrl, err = cfgContent.Get("database", "database_url")
	if err != nil {
		return nil, err
	}

	ConfigSettings.UseSqlite, err = cfgContent.GetBool("database", "use_sqlite")
	if err != nil {
		return nil, err
	}

	if ConfigSettings.UseSqlite {
		ConfigSettings.DbName, err = cfgContent.Get("database", "db_name")
		if err != nil {
			ConfigSettings.DbName = "dominatorRobot-db"
		}
	}

	ConfigSettings.MaxCacheTime, err = cfgContent.GetInt64("database", "max_cache_time")
	if err != nil {
		ConfigSettings.MaxCacheTime = 40
	}

	return ConfigSettings, nil
}

func LoadConfig() (*BotConfig, error) {
	return ParseConfig("config.ini")
}

func IsDebug() bool {
	if ConfigSettings != nil {
		return ConfigSettings.IsDebug
	}
	return true
}

func GetMaxCacheTime() time.Duration {
	if ConfigSettings != nil {
		return time.Duration(ConfigSettings.MaxCacheTime) * time.Minute
	}
	return time.Duration(40) * time.Minute
}

func GetBotToken() string {
	if ConfigSettings != nil {
		return ConfigSettings.BotToken
	}
	return ""
}

func DropUpdates() bool {
	if ConfigSettings != nil {
		return ConfigSettings.DropUpdates
	}
	return false
}

func GetCmdPrefixes() []rune {
	return []rune{'/', '!'}
}

func GetSibylClient() sibyl.SibylClient {
	return sibyl.NewClient(
		ConfigSettings.SibylToken,
		&sibyl.SibylConfig{
			HostUrl: ConfigSettings.SibylUrl,
		},
	)
}

func GetSibylConfig() *sibyl.SibylConfig {
	return &sibyl.SibylConfig{
		HostUrl: ConfigSettings.SibylUrl,
	}
}
