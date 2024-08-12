package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/MiguelMachado-dev/disc-go-bot/config"
	"github.com/MiguelMachado-dev/disc-go-bot/utils"
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

	var streamURL, thumbnailURL, streamTitle, game, profilePic string
	var color int

	switch platform {
	case "twitch":
		streamURL = fmt.Sprintf("https://www.twitch.tv/%s", username)
		color = 0x6441A4 // Twitch purple

		// Fetch Twitch channel info
		accessToken, err := utils.GetTwitchAccessToken()
		if err != nil {
			log.Errorf("Error getting Twitch access token: %v", err)
			break
		}

		channelInfo, err := h.getTwitchChannelInfo(username, accessToken)
		if err != nil {
			log.Errorf("Error getting Twitch channel info: %v", err)
			break
		}

		if channelInfo == nil {
			log.Infof("No channel found for user: %s", username)
			break
		}

		log.Infof("Channel info for %s: %+v", username, channelInfo)

		// Check if the channel is live
		streamInfo, err := h.getTwitchStreamInfo(username, accessToken)
		if err != nil {
			log.Errorf("Error getting Twitch stream info: %v", err)
		}

		isLive := streamInfo != nil && streamInfo.Type == "live"
		if isLive {
			streamTitle = streamInfo.Title
			thumbnailURL = strings.Replace(streamInfo.ThumbnailURL, "{width}x{height}", "1280x720", 1)
			game = streamInfo.GameName
			profilePic = channelInfo.ProfileImageURL
		} else {
			thumbnailURL = channelInfo.ProfileImageURL
		}
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
		Title:       fmt.Sprintf("%s no %s", username, platformTitle),
		Description: fmt.Sprintf("Confira o canal de %s", username),
		URL:         streamURL,
		Color:       color,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: thumbnailURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Link do canal",
				Value: streamURL,
			},
		},
	}

	if platform == "twitch" {
		if streamTitle != "" {
			embed.Title = streamTitle
		}
		embed.Description = fmt.Sprintf("Jogando %s.", game)
		embed.Image = &discordgo.MessageEmbedImage{
			URL: thumbnailURL,
		}
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: profilePic,
		}
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

type TwitchChannelInfo struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
}

func (h *LiveHandler) getTwitchChannelInfo(username, accessToken string) (*TwitchChannelInfo, error) {
	TwitchClientID := config.GetEnv().TWITCH_CLIENT_ID

	url := fmt.Sprintf("https://api.twitch.tv/helix/users?login=%s", username)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Client-ID", TwitchClientID)
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var result struct {
		Data []TwitchChannelInfo `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	log.Infof("Raw Twitch API response for channel info: %s", string(body))

	if len(result.Data) > 0 {
		return &result.Data[0], nil
	}

	return nil, nil
}

type TwitchStreamInfo struct {
	ID           string   `json:"id"`
	UserID       string   `json:"user_id"`
	UserName     string   `json:"user_name"`
	GameID       string   `json:"game_id"`
	GameName     string   `json:"game_name"`
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	ViewerCount  int      `json:"viewer_count"`
	StartedAt    string   `json:"started_at"`
	Language     string   `json:"language"`
	ThumbnailURL string   `json:"thumbnail_url"`
	TagIDs       []string `json:"tag_ids"`
	Tags         []string `json:"tags"`
	IsMature     bool     `json:"is_mature"`
}

func (h *LiveHandler) getTwitchStreamInfo(username, accessToken string) (*TwitchStreamInfo, error) {
	TwitchClientID := config.GetEnv().TWITCH_CLIENT_ID

	url := fmt.Sprintf("https://api.twitch.tv/helix/streams?user_login=%s", username)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Client-ID", TwitchClientID)
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var result struct {
		Data []TwitchStreamInfo `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	log.Infof("Raw Twitch API response: %s", string(body))

	if len(result.Data) > 0 {
		return &result.Data[0], nil
	}

	return nil, nil
}
