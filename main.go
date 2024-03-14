package main

import (
	"guildmaster/bot"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting application")

	err := bot.Run()
	if err != nil {
		log.Error(err)
	}
}
