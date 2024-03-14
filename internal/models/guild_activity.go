package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type GuildActivity struct {
	ID               uuid.UUID `gorm:"primaryKey"`
	GuildID          string    `gorm:"index"`
	UserID           string    `gorm:"index"`
	MessagesSent     int
	VoiceChannelTime time.Duration
}
