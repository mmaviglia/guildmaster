package bot

import (
	"fmt"
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

// Return a list of embed fields containing the leaderboard data.
func GuildLeaderboardFields(s *discordgo.Session, guildID string) ([]*discordgo.MessageEmbedField, error) {
	activities, err := models.GetGuildActivities(guildID)
	if err != nil {
		return nil, fmt.Errorf("get guild activities: %w", err)
	}

	names, messages, voiceTime := "", "", ""
	for i, a := range activities {
		// The leaderboard will show a maximum of 10 members
		if i == 10 {
			break
		}
		member, err := s.GuildMember(guildID, a.UserID)
		if err != nil {
			return nil, fmt.Errorf("get guild member: %w", err)
		}

		// Display a medal for the top three members; Otherwise their number
		rank := ""
		switch i {
		case 0:
			rank = "🥇"
		case 1:
			rank = "🥈"
		case 2:
			rank = "🥉"
		default:
			rank = fmt.Sprintf("#%v", i+1)
		}

		names += fmt.Sprintf("%v - %v\n", rank, member.DisplayName())
		messages += fmt.Sprintf("%v\n", a.MessagesSent)
		voiceTime += fmt.Sprintf("%v\n", formatDuration(time.Duration(a.VoiceChannelSeconds)*time.Second))
	}

	return []*discordgo.MessageEmbedField{
		{
			Name:   "Members",
			Value:  names,
			Inline: true,
		},
		{
			Name:   "Messages Sent",
			Value:  messages,
			Inline: true,
		},
		{
			Name:   "Voice Channel Time",
			Value:  voiceTime,
			Inline: true,
		},
	}, nil

}
