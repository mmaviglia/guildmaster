package bot

import (
	"fmt"
	"guildmaster/internal/config"
	"guildmaster/internal/models"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Start the Discord bot.
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

// Initiate the connection with Discord with proper intents and handlers added.
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
		// TODO: add switch statement for different interactions
	})

	// Initiate websocker connection with Discord
	if err := session.Open(); err != nil {
		return nil, fmt.Errorf("open session: %w", err)
	}

	return session, nil
}

// Called when the websocket connection with Discord has been successfully opened.
func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Add anything that needs to be run on first connection with Discord
}

// Called whenever the bot initially connects to a guild.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	log.Infof("Connected to Guild: %v\n", event.Guild.Name)
}

// Called any time a message is posted in a channel the bot is allowed to see.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages from the bot itself, and those that do not have the command prefix
	if m.Author.ID == s.State.User.ID {
		return
	}

	if err := models.IncrementMessagesSent(m.Author.ID, m.GuildID); err != nil {
		log.Error(err)
	}
}

// Called any time a voice state updates.
func voiceStateUpdate(t *GuildActivityTracker, s *discordgo.Session, v *discordgo.VoiceStateUpdate) {

	if v.ChannelID != "" && v.BeforeUpdate == nil {
		log.Info("joining")
		t.VoiceChannelJoin(v.UserID)
	}

	if v.ChannelID == "" && v.BeforeUpdate != nil {
		log.Info("leaving")
		t.VoiceChannelLeave(v.UserID, v.GuildID)
	}
}
