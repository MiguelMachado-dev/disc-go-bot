package handler

import "github.com/bwmarrin/discordgo"

type PingHandler struct{}

func (h *PingHandler) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Ping the bot",
	}
}

func (h *PingHandler) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
	if err != nil {
		log.Errorf("error responding to ping: %v", err)
		return
	}
}
