package handler

import (
	"bufio"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type AskHandler struct{}

func (h *AskHandler) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ask-ai",
		Description: "Ask the AI a question",
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

type Event struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason interface{} `json:"finish_reason"`
	} `json:"choices"`
}

func (h *AskHandler) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	q := i.ApplicationCommandData().Options[0].StringValue()

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Starting to ask AI...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("Error responding to interaction", err)
		return
	}

	// Make API request
	body := strings.NewReader(`{
		"model": "lmstudio-community/Meta-Llama-3-8B-Instruct-GGUF",
		"messages": [
			{ "role": "system", "content": "You are a helpful, smart, kind, and efficient AI assistant. You always fulfill the user's requests to the best of your ability." },
			{ "role": "user", "content": "` + q + `"}
		],
		"temperature": 0.7,
		"max_tokens": -1,
		"stream": true
	}`)

	req, err := http.NewRequest("POST", "http://192.168.3.2:1234/v1/chat/completions", body)
	if err != nil {
		log.Println("Error creating request", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := netClient.Do(req)

	if err != nil {
		log.Println("Error making request", err)
		s.ChannelMessageSend(i.ChannelID, "Error making request. Maybe the AI server is down?")
		return
	}

	defer resp.Body.Close()

	var content string

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "data: ") {
			jsonData := strings.TrimPrefix(line, "data: ")

			var js json.RawMessage
			err = json.Unmarshal([]byte(jsonData), &js)
			if err != nil {
				log.Println("Received data is not valid JSON, skipping unmarshalling")
				continue
			}

			var event Event
			err = json.Unmarshal([]byte(jsonData), &event)
			if err != nil {
				log.Fatalln("Error unmarshalling response", err)
			}

			// Send all event.Choices[0].Delta.Content to the channel as one message
			for _, choice := range event.Choices {
				content += choice.Delta.Content
			}

		} else if line == "data: [DONE]" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln("Error reading response body", err)
	}

	var msg discordgo.MessageSend

	msg.Content = "> USER: " + q + "\n\nAI: " + content

	_, err = s.ChannelMessageSend(i.ChannelID, msg.Content)

	if err != nil {
		log.Println("Error sending message to channel", err)
	}
}
