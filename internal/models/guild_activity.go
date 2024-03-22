package models

import (
	"fmt"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// A user's activity within a particular guild.
type GuildActivity struct {
	ID                  uuid.UUID `gorm:"primaryKey"`
	GuildID             string    `gorm:"index"`
	UserID              string    `gorm:"index"`
	MessagesSent        int       `gorm:"type:bigint"`
	VoiceChannelSeconds int       `gorm:"type:bigint"`
}

// Return the guild activity using the user and guild ID strings.
func GetGuildActivity(userID, guildID string) (*GuildActivity, error) {
	m := &GuildActivity{}
	return m, DB.Find(m, "id = ?", guildActivityID(userID, guildID)).Error
}

// Return the guild activity using the UUID primary key.
func GetGuildActivityByID(id uuid.UUID) (*GuildActivity, error) {
	m := &GuildActivity{}
	return m, DB.First(m, "id = ?", id).Error
}

// Create a GuildActivity record in the database.
func CreateGuildActivity(activity *GuildActivity) error {
	activity.ID = guildActivityID(activity.UserID, activity.GuildID)
	return DB.Create(activity).Error
}

// Return the ID of a GuildActivity record based on the user and guild IDs.
func guildActivityID(userID, guildID string) uuid.UUID {
	return uuid.NewV5(namespace, fmt.Sprintf(userID+guildID))
}

// Increment the number of messages sent by the user within the given guild.
func IncrementMessagesSent(userID, guildID string) error {
	tx := DB.Model(&GuildActivity{}).Where("id = ?", guildActivityID(userID, guildID))
	result := tx.UpdateColumn("messages_sent", gorm.Expr("messages_sent + 1"))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		err := CreateGuildActivity(&GuildActivity{
			UserID:       userID,
			GuildID:      guildID,
			MessagesSent: 1,
		})
		if err != nil {
			return fmt.Errorf("create guild activity: %w", err)
		}
	}

	return nil
}

// Increment the duration spent within voice channels by the user within the given guild.
func IncrementVoiceChannelDuration(userID, guildID string, duration int) error {
	tx := DB.Model(&GuildActivity{}).Where("id = ?", guildActivityID(userID, guildID))
	result := tx.UpdateColumn("voice_channel_seconds", gorm.Expr("voice_channel_seconds + ?", duration))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		err := CreateGuildActivity(&GuildActivity{
			UserID:              userID,
			GuildID:             guildID,
			VoiceChannelSeconds: duration,
		})
		if err != nil {
			return fmt.Errorf("create guild activity: %w", err)
		}
	}

	return nil
}
