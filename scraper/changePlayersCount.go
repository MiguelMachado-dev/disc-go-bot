package scraper

import (
	"fmt"
	"time"

	"github.com/MiguelMachado-dev/disc-go-bot/utils"
	"github.com/bwmarrin/discordgo"
)

func ChangePlayersCount(s *discordgo.Session, channelID string) {
	playersCh := make(chan string)
	go hunt(playersCh) // Call the scraper in a goroutine (concurrently)

	var players string
	select {
	case players = <-playersCh:
		// Successfully received players count from the scraper
	case <-time.After(30 * time.Second):
		// Timeout after 30 seconds if the scraper doesn't respond
		log.Warnln("Scraper timed out")
	}

	newName := fmt.Sprintf("Hunt p: %s", players)

	// Change the voice channel name
	utils.ChangeVoiceChannelName(s, channelID, newName)
}
