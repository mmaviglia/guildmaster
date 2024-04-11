package bot

import (
	"fmt"
	"guildmaster/internal/config"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "info",
		Description: "Displays general information about the bot.",
	},
	{
		Name:        "leaderboard",
		Description: "Displays the current leaderboard for the server",
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"info": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       config.BOT_EMBED_COLOR,
						Title:       "Info Panel",
						Description: "Discord bot focused on general functionality and basic leveling mechanics.",
						Author: &discordgo.MessageEmbedAuthor{
							Name:    s.State.User.Username,
							URL:     config.BOT_URL,
							IconURL: s.State.User.AvatarURL(""),
						},
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: config.BOT_EMBED_THUMBNAIL_URL,
						},
						Timestamp: time.Now().Format(ISO8601),
						Fields: []*discordgo.MessageEmbedField{
							{},
							{
								Name:   "Servers:",
								Value:  fmt.Sprintf("%d", len(s.State.Guilds)),
								Inline: false,
							},
							{
								Name:   "Memory Usage:",
								Value:  memoryUsageString(),
								Inline: true,
							},
							{
								Name:   "Running For:",
								Value:  runningDurationString(),
								Inline: true,
							},
							{},
						},
					},
				},
			},
		})
		if err != nil {
			log.Error(err)
		}
	},
	"leaderboard": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		leaderboard, err := GuildLeaderboardFields(s, i.GuildID)
		if err != nil {
			log.Error(fmt.Errorf("guild leaderboard fields: %w", err))
		}
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       config.BOT_EMBED_COLOR,
						Title:       "Leaderboard",
						Description: "Members in the server are ranked by their overall XP.",
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: config.BOT_EMBED_THUMBNAIL_URL,
						},
						Fields: leaderboard,
					},
				},
			},
		})
		if err != nil {
			log.Error(err)
		}
	},
}

// Update the Guild's registered application commands to match the codebase.
func syncServerCommands(s *discordgo.Session, guildID string) error {
	existingCommands, err := s.ApplicationCommands(s.State.User.ID, guildID)
	if err != nil {
		return err
	}

	existingCommandMap := make(map[string]bool)
	for _, cmd := range existingCommands {
		existingCommandMap[cmd.Name] = true
	}

	commandMap := make(map[string]bool)
	for _, cmd := range commands {
		commandMap[cmd.Name] = true
	}

	// Remove any commands that are not present in the codebase.
	for _, cmd := range existingCommands {
		if _, ok := commandMap[cmd.Name]; !ok {
			err := s.ApplicationCommandDelete(s.State.User.ID, guildID, cmd.ID)
			if err != nil {
				return err
			}
		}
	}

	// Register any new commands with the guild.
	for _, cmd := range commands {
		if _, ok := existingCommandMap[cmd.Name]; !ok {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
