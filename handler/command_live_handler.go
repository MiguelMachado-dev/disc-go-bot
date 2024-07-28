package handler

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "platform",
				Description: "Streaming platform (twitch or youtube)",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Twitch",
						Value: "twitch",
					},
					{
						Name:  "YouTube",
						Value: "youtube",
					},
				},
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

	// Extract the username and platform
	options := i.ApplicationCommandData().Options
	username := options[0].StringValue()
	platform := strings.ToLower(options[1].StringValue())

	var streamURL, thumbnailURL string
	var color int

	switch platform {
	case "twitch":
		streamURL = fmt.Sprintf("https://www.twitch.tv/%s", username)
		thumbnailURL = fmt.Sprintf("https://static-cdn.jtvnw.net/previews-ttv/live_user_%s-1920x1080.jpg", username)
		color = 0x6441A4 // Twitch purple
	case "youtube":
		streamURL = fmt.Sprintf("https://www.youtube.com/channel/%s/live", username)
		thumbnailURL = fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", username)
		color = 0xFF0000 // YouTube red
	default:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid platform selected. Please choose either 'twitch' or 'youtube'.",
			},
		})
		return
	}

	// Use cases.Title() instead of strings.Title()
	caser := cases.Title(language.BrazilianPortuguese)
	platformTitle := caser.String(platform)

	// Create an embed message
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s está ao vivo no %s!", username, platformTitle),
		Description: fmt.Sprintf("Clique no link abaixo para assistir a live de %s.", username),
		URL:         streamURL,
		Color:       color,
		Image: &discordgo.MessageEmbedImage{
			URL: thumbnailURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Link da live",
				Value: streamURL,
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
