package main

import (
	"fmt"

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
