package bot

import (
	"fmt"
	"guildmaster/internal/config"
	"guildmaster/internal/models"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// ISO8601 is the timestamp format required by Discord.
const ISO8601 = "2006-01-02T15:04:05Z07:00"

// Run starts the Discord bot.
func Run() error {
	session, err := createSession()
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}

	// Hang indefinitely unless CTRL-C is received
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-done

	// Clean up the Discord session before terminating
	return session.Close()
}

// createSession initiates the connection with Discord with proper intents and handlers added.
func createSession() (s *discordgo.Session, err error) {
	// Create a new Discord session using the provided bot token
	session, err := discordgo.New("Bot " + config.DISCORD_TOKEN)
	if err != nil {
		return nil, fmt.Errorf("discord new: %w", err)
	}

	tracker := NewGuildActivityTracker(session)

	// Declare all intents
	session.Identify.Intents = discordgo.IntentsAll

	// Add handler functions that will be called for each event
	session.AddHandler(ready)
	session.AddHandler(guildCreate)
	session.AddHandler(messageCreate)
	session.AddHandler(func(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
		voiceStateUpdate(tracker, s, v)
	})
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		}
	})

	// Initiate websocket connection with Discord
	if err := session.Open(); err != nil {
		return nil, fmt.Errorf("open session: %w", err)
	}

	// Track uptime from this point forward
	runningSince = time.Now()

	return session, nil
}

// ready should be called when the websocket connection with Discord has been successfully opened.
func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Add anything that needs to be run on first connection with Discord
}

// guildCreate should be called whenever the bot initially connects to a guild.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	log.Infof("Connected to Guild: %v\n", event.Guild.Name)

	err := syncServerCommands(s, event.Guild.ID)
	if err != nil {
		log.Error(fmt.Errorf("sync server commands: %w", err))
	}
}

// messageCreate should be called any time a message is posted in a channel the bot is allowed to see.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages from the bot itself, and those that do not have the command prefix
	if m.Author.ID == s.State.User.ID {
		return
	}

	if err := models.IncrementMessagesSent(m.Author.ID, m.GuildID); err != nil {
		log.Error(err)
	}
}

// voiceStateUpdate should be called any time a voice state updates.
func voiceStateUpdate(t *GuildActivityTracker, _ *discordgo.Session, v *discordgo.VoiceStateUpdate) {

	if v.ChannelID != "" && v.BeforeUpdate == nil {
		t.VoiceChannelJoin(v.UserID)
	}

	if v.ChannelID == "" && v.BeforeUpdate != nil {
		t.VoiceChannelLeave(v.UserID, v.GuildID)
	}
}
