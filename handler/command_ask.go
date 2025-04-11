package handler

import (
	"context"
	"fmt"

	"github.com/MiguelMachado-dev/disc-go-bot/config"
	"github.com/MiguelMachado-dev/disc-go-bot/database"
	"github.com/bwmarrin/discordgo"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
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
		Description: "Set your Gemini API key (Don't worry, it's stored securely)",
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

	// Create Gemini client using official library
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Errorf("Error creating Gemini client: %v", err)
		s.ChannelMessageSend(i.ChannelID, "Error connecting to Gemini API. Please try again later.")
		return
	}
	defer client.Close()

	// Configure the model
	model := client.GenerativeModel("gemini-2.0-flash")
	temperature := float32(0.7)
	maxOutputTokens := int32(800)
	model.SetTemperature(temperature)
	model.SetMaxOutputTokens(maxOutputTokens)

	// Generate content
	resp, err := model.GenerateContent(ctx, genai.Text(q))
	if err != nil {
		log.Errorf("Error generating content: %v", err)
		s.ChannelMessageSend(i.ChannelID, "Error connecting to Gemini API. Please check your API key and try again later.")
		return
	}

	var content string
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		content = fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
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
