package main

import (
	"log"

	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/logging"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/core/wotoConfig"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/database"
	"github.com/AnimeKaizoku/DominatorRobot/dominatorRobot/plugins"
)

func main() {
	_, err := wotoConfig.LoadConfig()
	if err != nil {
		log.Fatal("Error parsing config file", err)
	}

	f := logging.LoadLogger()
	if f != nil {
		defer f()
	}

	err = database.StartDatabase()
	if err != nil {
		logging.Fatal("Error starting database", err)
	}

	err = plugins.StartTelegramBot()
	if err != nil {
		logging.Fatal("Failed to start the bot bot: ", err)
	}
}
