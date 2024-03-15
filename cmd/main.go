package main

import (
	"guildmaster/internal/bot"
	"guildmaster/internal/models"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting application")

	if err := models.SetupDB(); err != nil {
		log.Error(err)
	}

	if err := bot.Run(); err != nil {
		log.Error(err)
	}
}
