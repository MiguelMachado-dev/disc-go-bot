package main

import (
	"fmt"
	"time"

	"github.com/MiguelMachado-dev/disc-go-bot/scraper"
	"github.com/bwmarrin/discordgo"
)

func updateGGStrivePlayerCountChannel(s *discordgo.Session, channelID string, intervalMinutes int) {
	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		playersCh := make(chan string)
		go scraper.GuiltyGear(playersCh) // Call the scraper in a goroutine (concurrently)

		var players string
		select {
		case players = <-playersCh:
			// Successfully received players count from the scraper
		case <-time.After(30 * time.Second):
			// Timeout after 30 seconds if the scraper doesn't respond
			fmt.Println("Scraper timed out")
			continue
		}

		newName := fmt.Sprintf("GG Strive p: %s", players)

		// Change the voice channel name
		changeVoiceChannelName(s, channelID, newName)
	}
}
