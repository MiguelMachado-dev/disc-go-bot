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
	playersCh := make(chan string)

	// Start the scraper in a separate goroutine
	go scraper.GuiltyGear(playersCh)
	// Receive the players count from the channel
	players := <-playersCh

	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		newName := fmt.Sprintf("Guilty Gear Players: %s", players)

		// Altere o nome do canal de voz
		changeVoiceChannelName(s, channelID, newName)
	}
}
