package utils

import (
	"github.com/bwmarrin/discordgo"
)

func ChangeVoiceChannelName(s *discordgo.Session, channelID string, newName string) {
	channelEdit := &discordgo.ChannelEdit{Name: newName}
	_, err := s.ChannelEditComplex(channelID, channelEdit)
	if err != nil {
		log.Errorln("Erro ao alterar o nome do canal de voz:", err)
	} else {
		log.Infoln("Nome do canal de voz '%s' alterado para '%s'\n", channelID, newName)
	}
}
