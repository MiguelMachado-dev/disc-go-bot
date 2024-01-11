package handler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type LiveHandler struct{}

func (h *LiveHandler) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "live",
		Description: "Create a live announcement",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "username",
				Description: "Username of the live streamer",
				Required:    true,
			},
		},
	}
}

func (h *LiveHandler) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if the user has the required role
	requiredRoles := []string{"1187118777079451669", "1187550419187142757"}
	hasRole := false

	for _, role := range i.Member.Roles {
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}

		if hasRole {
			break
		}
	}

	if !hasRole {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Você não tem permissão para usar esse comando.",
			},
		})
		return
	}

	// Extract the username
	username := i.ApplicationCommandData().Options[0].StringValue()

	// Construct Twitch stream URL
	twitchURL := fmt.Sprintf("https://www.twitch.tv/%s", username)
	twitchThumbnailURL := fmt.Sprintf("https://static-cdn.jtvnw.net/previews-ttv/live_user_%s-1920x1080.jpg", username)

	// Create an embed message
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s está ao vivo!", username),
		Description: fmt.Sprintf("Clique no link abaixo para assistir a live do %s.", username),
		URL:         twitchURL, // Add the Twitch URL here
		Color:       0x00ff00,  // Change color as needed
		Image: &discordgo.MessageEmbedImage{
			URL: twitchThumbnailURL, // Set the Twitch stream thumbnail as the image
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Link da live",
				Value: twitchURL,
			},
		},
	}

	// Send embed message to a specific channel
	channelID := "1187123213382201486"
	_, err := s.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		log.Errorf("error sending embed message: %v", err)
		return
	}

	// Respond to the interaction
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Live announcement created successfully.",
		},
	})
}
