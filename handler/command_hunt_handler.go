package handler

import (
	"github.com/MiguelMachado-dev/disc-go-bot/config"
	"github.com/MiguelMachado-dev/disc-go-bot/scraper"
	"github.com/bwmarrin/discordgo"
)

type HuntHandler struct{}

func (h *HuntHandler) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "hunt-player-update",
		Description: "Update Hunt: Showdown players count",
	}
}

func (h *HuntHandler) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	huntChannelID := config.GetEnv().HUNT_CHANNEL_ID
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Players count updated!",
		},
	})

	scraper.ChangePlayersCount(s, huntChannelID)

	if err != nil {
		log.Errorf("Error updating Hunt: Showdown players count: %v", err)
		return
	}
}
