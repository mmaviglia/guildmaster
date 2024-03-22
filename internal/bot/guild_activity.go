package bot

import (
	"guildmaster/internal/models"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// GuildActivityTracker tracks voice channel join timestamps for each user in a guild.
type GuildActivityTracker struct {
	// Mutex for concurrent access to the map
	mu sync.Mutex

	// Map to store voice channel join timestamps
	vcJoinTimes map[string]time.Time
}

// Create a new instance of GuildActivityTracker.
func NewGuildActivityTracker(session *discordgo.Session) *GuildActivityTracker {
	return &GuildActivityTracker{
		vcJoinTimes: make(map[string]time.Time),
	}
}

// Record the timestamp when a user joins a voice channel.
func (tracker *GuildActivityTracker) VoiceChannelJoin(userID string) {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	tracker.vcJoinTimes[userID] = time.Now()
}

// Updates GuildActivity when a user leaves a voice channel.
func (tracker *GuildActivityTracker) VoiceChannelLeave(userID string, guildID string) {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	joinTime, ok := tracker.vcJoinTimes[userID]
	if !ok {
		return
	}

	// Increment duration
	if err := models.IncrementVoiceChannelDuration(userID, guildID, int(time.Since(joinTime).Seconds())); err != nil {
		log.Error(err)
	}

	// Remove user join time from map
	delete(tracker.vcJoinTimes, userID)
}
