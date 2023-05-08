package scraper

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func UpdateHuntPlayerCountTicker(s *discordgo.Session, channelID string, intervalMinutes int) {
	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		ChangePlayersCount(s, channelID)
	}
}
