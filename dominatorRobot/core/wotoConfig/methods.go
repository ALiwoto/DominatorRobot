package wotoConfig

func (c *BotConfig) GetBotToken() string {
	return c.BotToken
}

func (c *BotConfig) GetDropUpdates() bool {
	return c.DropUpdates
}

func (c *BotConfig) GetDatabaseURL() string {
	return c.DatabaseUrl
}

func (c *BotConfig) GetIsDebug() bool {
	return c.IsDebug
}
