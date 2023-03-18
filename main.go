package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Get the bot token from the environment variable
	token := os.Getenv("DISCORD_BOT_TOKEN")

	// Create a new Discord session using the bot token
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	// Register the messageCreate function as a callback for the MessageCreate event
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	// Change voice channel name each 5 minutes
	go ChangeVoiceChannelNamePeriodically(dg, "1086042539997536336", 5)

	// Wait here until CTRL-C or other term signal is received
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Define the command prefix
	prefix := ">"

	// Check if the message starts with the command prefix
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	// Remove the prefix from the message content and split it into arguments
	args := strings.Fields(strings.TrimPrefix(m.Content, prefix))

	// If there are no arguments, return
	if len(args) == 0 {
		return
	}

	// Check the first argument to determine the command
	switch args[0] {
	case "ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
		// Add more commands here by extending the switch statement
	case "meow":
		SendCatImage(s, m.ChannelID)
	case "auau":
		sendDogImage(s, m.ChannelID)
	case "birb":
		SendBirbImage(s, m.ChannelID)
	}
}

func sendDogImage(s *discordgo.Session, channelID string) {
	// Send a "loading" message to the channel
	loadingMessage, err := s.ChannelMessageSend(channelID, "Loading dog image...")
	if err != nil {
		fmt.Println("Error sending loading message:", err)
		return
	}

	url := "https://random.dog/woof.json"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making API request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading API response:", err)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return
	}

	imageURL := result["url"].(string)

	// Download the image data using another HTTP GET request
	imageResp, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("Error downloading image:", err)
		return
	}
	defer imageResp.Body.Close()

	// Use the image data as the Reader for discordgo.File
	attachment := discordgo.File{
		Name:   "dog_image.jpg",
		Reader: imageResp.Body,
	}

	// Send the image and delete the "loading" message
	_, err = s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Files: []*discordgo.File{&attachment},
	})

	if err != nil {
		fmt.Println("Error sending dog image:", err)
		return
	}

	err = s.ChannelMessageDelete(channelID, loadingMessage.ID)
	if err != nil {
		fmt.Println("Error deleting loading message:", err)
	}
}
