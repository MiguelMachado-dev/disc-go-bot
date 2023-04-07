package main

import (
	"fmt"
	"time"

	"github.com/MiguelMachado-dev/disc-go-bot/scraper"
	"github.com/bwmarrin/discordgo"
)

func changeVoiceChannelName(s *discordgo.Session, channelID string, newName string) {
	channelEdit := &discordgo.ChannelEdit{Name: newName}
	_, err := s.ChannelEditComplex(channelID, channelEdit)
	if err != nil {
		fmt.Println("Erro ao alterar o nome do canal de voz:", err)
	} else {
		fmt.Printf("Nome do canal de voz '%s' alterado para '%s'\n", channelID, newName)
	}
}

func ChangeVoiceChannelNamePeriodically(s *discordgo.Session, channelID string, intervalMinutes int) {
	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Second)
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

		newName := fmt.Sprintf("Guilty Gear Players: %s", players)

		// Change the voice channel name
		changeVoiceChannelName(s, channelID, newName)
	}
}
