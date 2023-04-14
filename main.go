package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var serverVoiceChannelIDs = make(map[string]string)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Get the bot token from the environment variable
	token := os.Getenv("DISCORD_BOT_TOKEN")
	guiltyGearStriveChannelID := os.Getenv("GG_CHANNEL_ID")
	caliberChannelID := os.Getenv("CALIBER_CHANNEL_ID")

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

	// Set the bot's presence to "Streaming on Twitch"
	dg.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: "Migtito on Twitch",
				Type: discordgo.ActivityTypeStreaming,
				URL:  "https://www.twitch.tv/Migtito",
			},
		},
		Status: "online",
	})

	// Change voice channel name each 30 minutes
	go updateGGStrivePlayerCountChannel(dg, guiltyGearStriveChannelID, 30)
	go updateCaliberPlayerCountChannel(dg, caliberChannelID, 30)

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
	prefix := os.Getenv("DISCORD_BOT_PREFIX")

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
	case "birb":
		SendBirbImage(s, m.ChannelID)
	case "setchannel":
		if len(args) != 2 {
			usageMessage := fmt.Sprintf("Usage: %ssetchannel [channelID]", prefix)
			s.ChannelMessageSend(m.ChannelID, usageMessage)
			return
		}

		if isUserAdmin(s, m) {
			serverVoiceChannelIDs[m.GuildID] = args[1]
			s.ChannelMessageSend(m.ChannelID, "Voice channel updated! - Comando em beta, não faz nada.")
		} else {
			s.ChannelMessageSend(m.ChannelID, "You don't have permission to set the voice channel.")
		}
	}
}
