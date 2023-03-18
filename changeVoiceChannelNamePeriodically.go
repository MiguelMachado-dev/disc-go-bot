package main

import (
	"fmt"
	"math/rand"
	"time"

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
	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		// Gere um novo nome para o canal de voz
		newName := fmt.Sprintf("Voice Channel %d", rand.Intn(1000))

		// Altere o nome do canal de voz
		changeVoiceChannelName(s, channelID, newName)
	}
}
