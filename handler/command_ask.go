package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MiguelMachado-dev/disc-go-bot/config"
	"github.com/MiguelMachado-dev/disc-go-bot/database"
	"github.com/bwmarrin/discordgo"
)

// AskHandler struct for Ask command
type AskHandler struct{}

func (h *AskHandler) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ask-ai",
		Description: "Ask the AI a question using Gemini API",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "question",
				Description: "Question to ask the AI",
				Required:    true,
			},
		},
	}
}

// SetGeminiKeyHandler handles the set-gemini-key command
type SetGeminiKeyHandler struct{}

func (h *SetGeminiKeyHandler) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "set-gemini-key",
		Description: "Set your Gemini API key",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "api-key",
				Description: "Your Gemini API key",
				Required:    true,
			},
		},
	}
}

func (h *SetGeminiKeyHandler) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	apiKey := i.ApplicationCommandData().Options[0].StringValue()
	userID := i.Member.User.ID
	log := config.NewLogger("SetGeminiKeyHandler")

	// Store the API key in the database
	err := database.StoreGeminiAPIKey(userID, apiKey)
	if err != nil {
		log.Errorf("Error storing API key: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "An error occurred while saving your API key. Please try again later.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Your Gemini API key has been set. You can now use the /ask-ai command.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Errorf("Error responding to interaction: %v", err)
	}
}

// GeminiResponse represents the structure of the Gemini API response
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
}

func (h *AskHandler) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	q := i.ApplicationCommandData().Options[0].StringValue()
	userID := i.Member.User.ID
	log := config.NewLogger("AskHandler")

	// Get the API key from the database
	apiKey, err := database.GetGeminiAPIKey(userID)
	if err != nil {
		log.Errorf("Error retrieving API key: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You need to set your Gemini API key first with `/set-gemini-key`",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Asking Gemini AI...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Errorf("Error responding to interaction: %v", err)
		return
	}

	// Gemini API request
	requestBody := fmt.Sprintf(`{
		"contents": [
			{
				"parts": [
					{
						"text": "%s"
					}
				]
			}
		],
		"generationConfig": {
			"temperature": 0.7,
			"maxOutputTokens": 800
		}
	}`, strings.ReplaceAll(q, "\"", "\\\""))

	req, err := http.NewRequest("POST", "https://generativelanguage.googleapis.com/v1/models/gemini-2.0-flash:generateContent?key="+apiKey, strings.NewReader(requestBody))
	if err != nil {
		log.Errorf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := netClient.Do(req)
	if err != nil {
		log.Errorf("Error making request: %v", err)
		s.ChannelMessageSend(i.ChannelID, "Error connecting to Gemini API. Please try again later.")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Gemini API error: %d - %s - API URL: %s", resp.StatusCode, resp.Status, req.URL.String())
		errorMessage := fmt.Sprintf("Gemini API error: %s. Please check your API key and ensure it has access to the Gemini models.", resp.Status)
		s.ChannelMessageSend(i.ChannelID, errorMessage)
		return
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		log.Errorf("Error decoding response: %v", err)
		s.ChannelMessageSend(i.ChannelID, "Error processing response from Gemini API.")
		return
	}

	var content string
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		content = geminiResp.Candidates[0].Content.Parts[0].Text
	} else {
		content = "No response generated."
	}

	// Send the response to the channel
	msg := "> USER: " + q + "\n\nAI: " + content

	_, err = s.ChannelMessageSend(i.ChannelID, msg)
	if err != nil {
		log.Errorf("Error sending message to channel: %v", err)
	}
}
