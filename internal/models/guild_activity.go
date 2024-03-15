package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// A user's activity within a particular guild.
type GuildActivity struct {
	ID               uuid.UUID `gorm:"primaryKey"`
	GuildID          string    `gorm:"index"`
	UserID           string    `gorm:"index"`
	MessagesSent     int
	VoiceChannelTime time.Duration
}

// Return the guild activity using the user and guild ID strings.
func GetGuildActivity(userID string, guildID string) (*GuildActivity, error) {
	m := &GuildActivity{}

	return m, DB.Find(m, "user_id = ? AND guild_id = ?", userID, guildID).Error
}

// Return the guild activity using the UUID primary key.
func GetGuildActivityByID(id uuid.UUID) (*GuildActivity, error) {
	m := &GuildActivity{}
	return m, DB.First(m, "id = ?", id).Error
}
